f"""
WebSocket consumer for shared application state between Django and Go.

This module provides a WebSocket consumer that connects to the Go WebSocket
state manager to enable real-time state synchronization between the Django
backend and the Go components of agent_runtime.
"""

import json
import logging
import asyncio
import uuid
from typing import Dict, Any, Optional, List
from channels.generic.websocket import AsyncWebsocketConsumer
from channels.db import database_sync_to_async
from django.conf import settings

logger = logging.getLogger(__name__)


class StateType:
    """State types that match the Go implementation."""
    TASK = "task"
    AGENT = "agent"
    LIFECYCLE = "lifecycle"
    SHARED = "shared"


class SharedStateConsumer(AsyncWebsocketConsumer):
    """
    WebSocket consumer for shared application state.

    This consumer handles WebSocket connections for shared state between
    Django and Go components, enabling real-time updates and synchronization.
    """

    connections: Dict[str, List['SharedStateConsumer']] = {}

    async def connect(self):
        """Handle WebSocket connection."""
        self.state_type = self.scope['url_route']['kwargs'].get('state_type', StateType.SHARED)
        self.state_id = self.scope['url_route']['kwargs'].get('state_id', 'default')
        self.connection_id = str(uuid.uuid4())

        key = f"{self.state_type}:{self.state_id}"
        if key not in self.connections:
            self.connections[key] = []
        self.connections[key].append(self)

        await self.accept()

        initial_state = await database_sync_to_async(self.get_initial_state)()
        if initial_state:
            await self.send(text_data=json.dumps({
                'type': 'state_update',
                'state_type': self.state_type,
                'state_id': self.state_id,
                'data': initial_state,
            }))

        logger.info(f"WebSocket connection established for {self.state_type} state with ID {self.state_id}")

    async def disconnect(self, close_code):
        """Handle WebSocket disconnection."""
        key = f"{self.state_type}:{self.state_id}"
        if key in self.connections:
            if self in self.connections[key]:
                self.connections[key].remove(self)
            if not self.connections[key]:
                del self.connections[key]

        logger.info(f"WebSocket connection closed for {self.state_type} state with ID {self.state_id}")

    async def receive(self, text_data):
        """
        Handle incoming WebSocket messages.

        This method processes messages from clients and updates the shared state
        in the Go state manager.
        """
        try:
            data = json.loads(text_data)
            message_type = data.get('type')

            if message_type == 'update_state':
                success = await database_sync_to_async(self.update_state)(data.get('data', {}))

                if success:
                    await self.broadcast_state_update(data.get('data', {}))

            elif message_type == 'get_state':
                state = await database_sync_to_async(self.get_initial_state)()
                await self.send(text_data=json.dumps({
                    'type': 'state_update',
                    'state_type': self.state_type,
                    'state_id': self.state_id,
                    'data': state,
                }))

            else:
                logger.warning(f"Unknown message type: {message_type}")

        except json.JSONDecodeError:
            logger.error("Failed to decode JSON message")
        except Exception as e:
            logger.error(f"Error processing WebSocket message: {e}")

    async def state_update(self, event):
        """
        Handle state update events.

        This method is called when a state update is received from the Go
        state manager or another Django component.
        """
        await self.send(text_data=json.dumps({
            'type': 'state_update',
            'state_type': event['state_type'],
            'state_id': event['state_id'],
            'data': event['data'],
        }))

    def get_initial_state(self) -> Optional[Dict[str, Any]]:
        """
        Get the initial state from the Go state manager.

        This method retrieves the current state for the specified state type
        and ID from the Go state manager.
        """
        try:
            from api.grpc_bridge import grpc_bridge

            response = grpc_bridge.get_state(
                state_type=self.state_type,
                state_id=self.state_id
            )

            if response and 'data' in response:
                return response['data']

            return None

        except Exception as e:
            logger.error(f"Error getting initial state: {e}")
            return None

    def update_state(self, data: Dict[str, Any]) -> bool:
        """
        Update the state in the Go state manager.

        This method sends a state update to the Go state manager through
        the gRPC bridge.
        """
        try:
            from api.grpc_bridge import grpc_bridge

            response = grpc_bridge.update_state(
                state_type=self.state_type,
                state_id=self.state_id,
                data=data
            )

            return response.get('status') == 'success'

        except Exception as e:
            logger.error(f"Error updating state: {e}")
            return False

    async def broadcast_state_update(self, data: Dict[str, Any]):
        """
        Broadcast a state update to all connected clients.

        This method sends a state update to all WebSocket clients connected
        to the same state type and ID.
        """
        key = f"{self.state_type}:{self.state_id}"
        if key in self.connections:
            for connection in self.connections[key]:
                await connection.send(text_data=json.dumps({
                    'type': 'state_update',
                    'state_type': self.state_type,
                    'state_id': self.state_id,
                    'data': data,
                }))



def get_shared_state(state_id: str = 'default') -> Optional[Dict[str, Any]]:
    """
    Get shared state from the Go state manager.

    This utility function retrieves the current shared state for the
    specified state ID from the Go state manager.
    """
    try:
        from api.grpc_bridge import grpc_bridge

        response = grpc_bridge.get_state(
            state_type=StateType.SHARED,
            state_id=state_id
        )

        if response and 'data' in response:
            return response['data']

        return None

    except Exception as e:
        logger.error(f"Error getting shared state: {e}")
        return None


async def update_shared_state(state_id: str, data: Dict[str, Any]) -> bool:
    """
    Update shared state in the Go state manager.

    This utility function sends a state update to the Go state manager
    through the gRPC bridge and broadcasts the update to all connected clients.
    """
    try:
        # Update state in Go state manager
        from api.grpc_bridge import grpc_bridge
        
        loop = asyncio.get_event_loop()
        response = await loop.run_in_executor(
            None,
            lambda: grpc_bridge.update_state(
                state_type=StateType.SHARED,
                state_id=state_id,
                data=data
            )
        )

        # Broadcast to all connected WebSocket clients
        key = f"{StateType.SHARED}:{state_id}"
        if key in SharedStateConsumer.connections:
            for connection in SharedStateConsumer.connections[key]:
                await connection.send(text_data=json.dumps({
                    'type': 'state_update',
                    'state_type': StateType.SHARED,
                    'state_id': state_id,
                    'data': data,
                }))

        return response.get('status') == 'success'

    except Exception as e:
        logger.error(f"Error updating shared state: {e}")
        return False

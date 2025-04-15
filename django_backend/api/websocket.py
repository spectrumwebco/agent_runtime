"""
WebSocket implementation for real-time communication.

This module provides WebSocket endpoints for real-time communication between
the frontend and backend components of the agent runtime system. It enables
two-way communication for agent state updates, task progress, and other
real-time events.
"""

import json
import logging
from typing import Dict, Set
from channels.generic.websocket import AsyncWebsocketConsumer
from django.conf import settings
import sys

sys.path.append(str(settings.SRC_DIR))

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


class AgentConsumer(AsyncWebsocketConsumer):
    """
    WebSocket consumer for agent runtime communication.

    This consumer handles WebSocket connections for real-time communication
    between the frontend and backend components of the agent runtime system.
    It enables two-way communication for agent state updates, task progress,
    and other real-time events.
    """

    connected_clients: Set[str] = set()
    task_groups: Dict[str, Set[str]] = {}

    async def connect(self):
        """Handle WebSocket connection."""
        self.client_id = self.scope['url_route']['kwargs'].get(
            'client_id', 'anonymous')
        self.task_id = self.scope['url_route']['kwargs'].get('task_id', None)

        await self.accept()

        AgentConsumer.connected_clients.add(self.client_id)

        if self.task_id:
            if self.task_id not in AgentConsumer.task_groups:
                AgentConsumer.task_groups[self.task_id] = set()
            AgentConsumer.task_groups[self.task_id].add(self.client_id)

            await self.channel_layer.group_add(
                f"task_{self.task_id}",
                self.channel_name
            )

        await self.channel_layer.group_add(
            "broadcast",
            self.channel_name
        )

        await self.send(text_data=json.dumps({
            'type': 'connection_established',
            'message': f'Connected as {self.client_id}',
            'client_id': self.client_id,
            'task_id': self.task_id
        }))

        logger.info(
            f"WebSocket connection established for client {self.client_id}")

    async def disconnect(self, close_code):
        """Handle WebSocket disconnection."""
        if self.client_id in AgentConsumer.connected_clients:
            AgentConsumer.connected_clients.remove(self.client_id)

        if self.task_id and self.task_id in AgentConsumer.task_groups:
            if self.client_id in AgentConsumer.task_groups[self.task_id]:
                AgentConsumer.task_groups[self.task_id].remove(self.client_id)

            if not AgentConsumer.task_groups[self.task_id]:
                del AgentConsumer.task_groups[self.task_id]

            await self.channel_layer.group_discard(
                f"task_{self.task_id}",
                self.channel_name
            )

        await self.channel_layer.group_discard(
            "broadcast",
            self.channel_name
        )

        logger.info(f"WebSocket connection closed for client {self.client_id}")

    async def receive(self, text_data):
        """Handle incoming WebSocket messages."""
        try:
            data = json.loads(text_data)
            message_type = data.get('type')

            if message_type == 'task_update':
                await self.handle_task_update(data)
            elif message_type == 'agent_command':
                await self.handle_agent_command(data)
            else:
                logger.warning(f"Unknown message type: {message_type}")
                await self.send(text_data=json.dumps({
                    'type': 'error',
                    'message': f'Unknown message type: {message_type}'
                }))
        except json.JSONDecodeError:
            logger.error(f"Invalid JSON received: {text_data}")
            await self.send(text_data=json.dumps({
                'type': 'error',
                'message': 'Invalid JSON'
            }))
        except Exception as e:
            logger.error(f"Error handling message: {str(e)}")
            await self.send(text_data=json.dumps({
                'type': 'error',
                'message': f'Error: {str(e)}'
            }))

    async def handle_task_update(self, data):
        """Handle task update message."""
        task_id = data.get('task_id')
        status = data.get('status')
        message = data.get('message')

        if not task_id:
            await self.send(text_data=json.dumps({
                'type': 'error',
                'message': 'Missing task_id'
            }))
            return

        await self.channel_layer.group_send(
            f"task_{task_id}",
            {
                'type': 'task_update',
                'task_id': task_id,
                'status': status,
                'message': message,
                'sender': self.client_id
            }
        )

    async def handle_agent_command(self, data):
        """Handle agent command message."""
        command = data.get('command')
        params = data.get('params', {})

        if not command:
            await self.send(text_data=json.dumps({
                'type': 'error',
                'message': 'Missing command'
            }))
            return

        await self.send(text_data=json.dumps({
            'type': 'command_received',
            'command': command,
            'params': params,
            'message': f'Command {command} received'
        }))

    async def task_update(self, event):
        """Handle task update event from channel layer."""
        await self.send(text_data=json.dumps({
            'type': 'task_update',
            'task_id': event['task_id'],
            'status': event['status'],
            'message': event['message'],
            'sender': event['sender']
        }))

    async def broadcast_message(self, event):
        """Handle broadcast message event from channel layer."""
        await self.send(text_data=json.dumps({
            'type': 'broadcast',
            'message': event['message'],
            'sender': event['sender']
        }))


async def send_task_update(task_id, status, message):
    """Send a task update to all clients in the task group."""
    from channels.layers import get_channel_layer
    channel_layer = get_channel_layer()

    await channel_layer.group_send(
        f"task_{task_id}",
        {
            'type': 'task_update',
            'task_id': task_id,
            'status': status,
            'message': message,
            'sender': 'system'
        }
    )


async def broadcast_message(message):
    """Broadcast a message to all connected clients."""
    from channels.layers import get_channel_layer
    channel_layer = get_channel_layer()

    await channel_layer.group_send(
        "broadcast",
        {
            'type': 'broadcast_message',
            'message': message,
            'sender': 'system'
        }
    )

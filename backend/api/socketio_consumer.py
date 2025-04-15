"""
Socket.IO implementation for the OpenHands frontend.

This module provides a Socket.IO server implementation that matches the
frontend's Socket.IO client, enabling real-time communication between
the frontend and backend components.
"""

import json
import logging
import asyncio
import uuid
from typing import Dict, Any, List, Set, Optional
from channels.generic.websocket import AsyncWebsocketConsumer
from django.conf import settings
from asgiref.sync import sync_to_async
from channels.db import database_sync_to_async

logger = logging.getLogger(__name__)


class OpenHandsSocketIOConsumer(AsyncWebsocketConsumer):
    """
    Socket.IO consumer for OpenHands frontend communication.

    This consumer handles WebSocket connections using the Socket.IO protocol
    to match the frontend's Socket.IO client implementation. It enables
    real-time communication for agent state updates, messages, and events.
    """

    connected_clients: Dict[str, 'OpenHandsSocketIOConsumer'] = {}
    conversation_groups: Dict[str, Set[str]] = {}
    
    latest_event_ids: Dict[str, int] = {}
    
    async def connect(self):
        """Handle WebSocket connection."""
        query_string = self.scope.get('query_string', b'').decode()
        query_params = dict(param.split('=') for param in query_string.split('&') if '=' in param)
        
        self.conversation_id = query_params.get('conversation_id', 'default')
        self.latest_event_id = int(query_params.get('latest_event_id', -1))
        self.client_id = str(uuid.uuid4())
        
        await self.accept()
        
        OpenHandsSocketIOConsumer.connected_clients[self.client_id] = self
        
        if self.conversation_id not in OpenHandsSocketIOConsumer.conversation_groups:
            OpenHandsSocketIOConsumer.conversation_groups[self.conversation_id] = set()
        OpenHandsSocketIOConsumer.conversation_groups[self.conversation_id].add(self.client_id)
        
        await self.channel_layer.group_add(
            f"conversation_{self.conversation_id}",
            self.channel_name
        )
        
        await self.send_event({
            'type': 'connection_established',
            'client_id': self.client_id,
            'conversation_id': self.conversation_id
        })
        
        await self.send_missed_events()
        
        logger.info(f"Socket.IO connection established for client {self.client_id} in conversation {self.conversation_id}")

    async def disconnect(self, close_code):
        """Handle WebSocket disconnection."""
        if self.client_id in OpenHandsSocketIOConsumer.connected_clients:
            del OpenHandsSocketIOConsumer.connected_clients[self.client_id]
        
        if self.conversation_id in OpenHandsSocketIOConsumer.conversation_groups:
            if self.client_id in OpenHandsSocketIOConsumer.conversation_groups[self.conversation_id]:
                OpenHandsSocketIOConsumer.conversation_groups[self.conversation_id].remove(self.client_id)
            
            if not OpenHandsSocketIOConsumer.conversation_groups[self.conversation_id]:
                del OpenHandsSocketIOConsumer.conversation_groups[self.conversation_id]
        
        await self.channel_layer.group_discard(
            f"conversation_{self.conversation_id}",
            self.channel_name
        )
        
        logger.info(f"Socket.IO connection closed for client {self.client_id}")

    async def receive(self, text_data):
        """Handle incoming WebSocket messages."""
        try:
            if text_data.startswith('42'):
                socketio_data = json.loads(text_data[2:])
                event_name = socketio_data[0]
                event_data = socketio_data[1] if len(socketio_data) > 1 else {}
                
                if event_name == 'oh_user_action':
                    await self.handle_user_action(event_data)
                else:
                    logger.warning(f"Unknown Socket.IO event: {event_name}")
            else:
                if text_data == '2':  # ping
                    await self.send(text_data='3')  # pong
                
        except json.JSONDecodeError:
            logger.error(f"Invalid JSON received: {text_data}")
        except Exception as e:
            logger.error(f"Error handling Socket.IO message: {str(e)}")
            await self.send_error(f"Error: {str(e)}")

    async def handle_user_action(self, data):
        """Handle user action event."""
        action_type = data.get('type')
        
        if action_type == 'message':
            await self.handle_user_message(data)
        elif action_type == 'agent_state_change':
            await self.handle_agent_state_change(data)
        else:
            await self.forward_to_agent_runtime(data)

    async def handle_user_message(self, data):
        """Handle user message event."""
        event_id = await self.get_next_event_id()
        
        event = {
            'id': str(event_id),
            'source': 'user',
            'type': 'message',
            'message': data.get('message', ''),
            'timestamp': data.get('timestamp', int(asyncio.get_event_loop().time() * 1000))
        }
        
        await self.store_event(event)
        
        await self.channel_layer.group_send(
            f"conversation_{self.conversation_id}",
            {
                'type': 'oh_event',
                'event': event
            }
        )
        
        await self.forward_to_agent_runtime(data)

    async def handle_agent_state_change(self, data):
        """Handle agent state change event."""
        event_id = await self.get_next_event_id()
        
        event = {
            'id': str(event_id),
            'source': 'user',
            'type': 'agent_state_change',
            'state': data.get('state'),
            'timestamp': data.get('timestamp', int(asyncio.get_event_loop().time() * 1000))
        }
        
        await self.store_event(event)
        
        await self.channel_layer.group_send(
            f"conversation_{self.conversation_id}",
            {
                'type': 'oh_event',
                'event': event
            }
        )
        
        await self.forward_to_agent_runtime(data)

    async def forward_to_agent_runtime(self, data):
        """Forward user action to the agent runtime."""
        try:
            from api.views import execute_agent_task_async
            
            data['conversation_id'] = self.conversation_id
            data['client_id'] = self.client_id
            
            asyncio.create_task(execute_agent_task_async(data))
            
        except Exception as e:
            logger.error(f"Error forwarding to agent runtime: {str(e)}")
            await self.send_error(f"Error: {str(e)}")

    async def oh_event(self, event):
        """Handle oh_event from channel layer."""
        event_data = event['event']
        await self.send_event(event_data)

    async def send_event(self, event_data):
        """Send an event to the client using Socket.IO protocol."""
        socketio_message = f'42["oh_event", {json.dumps(event_data)}]'
        await self.send(text_data=socketio_message)

    async def send_error(self, message):
        """Send an error event to the client."""
        await self.send_event({
            'type': 'error',
            'message': message,
            'timestamp': int(asyncio.get_event_loop().time() * 1000)
        })

    async def get_next_event_id(self) -> int:
        """Get the next event ID for the conversation."""
        if self.conversation_id not in OpenHandsSocketIOConsumer.latest_event_ids:
            OpenHandsSocketIOConsumer.latest_event_ids[self.conversation_id] = 0
        
        OpenHandsSocketIOConsumer.latest_event_ids[self.conversation_id] += 1
        return OpenHandsSocketIOConsumer.latest_event_ids[self.conversation_id]

    @database_sync_to_async
    def store_event(self, event):
        """Store the event in the database."""
        from api.models import ConversationEvent
        
        ConversationEvent.objects.create(
            conversation_id=self.conversation_id,
            event_id=event['id'],
            event_type=event['type'],
            source=event['source'],
            content=json.dumps(event),
            timestamp=event.get('timestamp')
        )

    async def send_missed_events(self):
        """Send missed events to the client."""
        if self.latest_event_id < 0:
            return
        
        from api.models import ConversationEvent
        
        events = await database_sync_to_async(self._get_missed_events)()
        
        for event in events:
            await self.send_event(json.loads(event.content))

    def _get_missed_events(self):
        """Get missed events from the database."""
        from api.models import ConversationEvent
        
        return list(ConversationEvent.objects.filter(
            conversation_id=self.conversation_id,
            event_id__gt=self.latest_event_id
        ).order_by('event_id'))



async def send_agent_message(conversation_id: str, message: str, extras: Optional[Dict[str, Any]] = None):
    """Send an agent message to all clients in a conversation."""
    from channels.layers import get_channel_layer
    channel_layer = get_channel_layer()
    
    event_id = await _get_next_event_id(conversation_id)
    
    event = {
        'id': str(event_id),
        'source': 'agent',
        'type': 'message',
        'message': message,
        'timestamp': int(asyncio.get_event_loop().time() * 1000)
    }
    
    if extras:
        event['extras'] = extras
    
    await _store_event(conversation_id, event)
    
    await channel_layer.group_send(
        f"conversation_{conversation_id}",
        {
            'type': 'oh_event',
            'event': event
        }
    )

async def send_agent_observation(conversation_id: str, observation: str, observation_type: str, extras: Optional[Dict[str, Any]] = None):
    """Send an agent observation to all clients in a conversation."""
    from channels.layers import get_channel_layer
    channel_layer = get_channel_layer()
    
    event_id = await _get_next_event_id(conversation_id)
    
    event = {
        'id': str(event_id),
        'source': 'agent',
        'type': 'observation',
        'observation': observation_type,
        'message': observation,
        'timestamp': int(asyncio.get_event_loop().time() * 1000)
    }
    
    if extras:
        event['extras'] = extras
    
    await _store_event(conversation_id, event)
    
    await channel_layer.group_send(
        f"conversation_{conversation_id}",
        {
            'type': 'oh_event',
            'event': event
        }
    )

async def send_agent_error(conversation_id: str, message: str, error_id: Optional[str] = None):
    """Send an error observation to all clients in a conversation."""
    extras = {'error_id': error_id} if error_id else None
    await send_agent_observation(conversation_id, message, 'error', extras)

async def send_agent_state_update(conversation_id: str, state: str):
    """Send an agent state update to all clients in a conversation."""
    from channels.layers import get_channel_layer
    channel_layer = get_channel_layer()
    
    event_id = await _get_next_event_id(conversation_id)
    
    event = {
        'id': str(event_id),
        'source': 'agent',
        'type': 'agent_state_update',
        'state': state,
        'timestamp': int(asyncio.get_event_loop().time() * 1000)
    }
    
    await _store_event(conversation_id, event)
    
    await channel_layer.group_send(
        f"conversation_{conversation_id}",
        {
            'type': 'oh_event',
            'event': event
        }
    )

async def _get_next_event_id(conversation_id: str) -> int:
    """Get the next event ID for a conversation."""
    if conversation_id not in OpenHandsSocketIOConsumer.latest_event_ids:
        OpenHandsSocketIOConsumer.latest_event_ids[conversation_id] = 0
    
    OpenHandsSocketIOConsumer.latest_event_ids[conversation_id] += 1
    return OpenHandsSocketIOConsumer.latest_event_ids[conversation_id]

@database_sync_to_async
def _store_event(conversation_id: str, event):
    """Store an event in the database."""
    from api.models import ConversationEvent
    
    ConversationEvent.objects.create(
        conversation_id=conversation_id,
        event_id=event['id'],
        event_type=event['type'],
        source=event['source'],
        content=json.dumps(event),
        timestamp=event.get('timestamp')
    )

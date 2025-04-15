"""
WebSocket consumers for the API.

This module provides WebSocket consumers for the API that integrate with
the WebSocket manager and handle real-time communication with clients.
"""

import json
import logging
import uuid
from typing import Dict, Any, List, Optional
from channels.generic.websocket import AsyncWebsocketConsumer
from django.contrib.auth.models import AnonymousUser
from rest_framework.authtoken.models import Token

from api.websocket_manager import get_manager

logger = logging.getLogger(__name__)


class BaseWebSocketConsumer(AsyncWebsocketConsumer):
    """Base WebSocket consumer for the API."""
    
    async def connect(self):
        """Handle WebSocket connection."""
        self.consumer_id = str(uuid.uuid4())
        self.user = self.scope.get('user', AnonymousUser())
        self.groups = []
        
        if self.user and self.user.is_authenticated:
            self.groups.append(f"user_{self.user.id}")
        
        manager = get_manager()
        manager.register_consumer(self.consumer_id, self, self.groups)
        
        await self.accept()
        
        await self.send(text_data=json.dumps({
            'type': 'connection_established',
            'consumer_id': self.consumer_id,
        }))
    
    async def disconnect(self, close_code):
        """Handle WebSocket disconnection."""
        manager = get_manager()
        manager.unregister_consumer(self.consumer_id)
    
    async def receive(self, text_data):
        """Handle incoming WebSocket messages."""
        try:
            data = json.loads(text_data)
            message_type = data.get('type')
            
            if message_type == 'ping':
                await self.send(text_data=json.dumps({
                    'type': 'pong',
                    'timestamp': data.get('timestamp'),
                }))
            else:
                await self.handle_message(data)
        except json.JSONDecodeError:
            logger.error(f"Invalid JSON received: {text_data}")
            await self.send(text_data=json.dumps({
                'type': 'error',
                'message': 'Invalid JSON',
            }))
        except Exception as e:
            logger.error(f"Error handling message: {e}")
            await self.send(text_data=json.dumps({
                'type': 'error',
                'message': str(e),
            }))
    
    async def handle_message(self, data: Dict[str, Any]):
        """Handle a specific message type."""
        message_type = data.get('type')
        
        if message_type == 'subscribe':
            await self.handle_subscribe(data)
        elif message_type == 'unsubscribe':
            await self.handle_unsubscribe(data)
        elif message_type == 'event':
            await self.handle_event(data)
        else:
            logger.warning(f"Unknown message type: {message_type}")
            await self.send(text_data=json.dumps({
                'type': 'error',
                'message': f'Unknown message type: {message_type}',
            }))
    
    async def handle_subscribe(self, data: Dict[str, Any]):
        """Handle subscription to events."""
        event_types = data.get('event_types', [])
        
        if not event_types:
            await self.send(text_data=json.dumps({
                'type': 'error',
                'message': 'No event types specified',
            }))
            return
        
        manager = get_manager()
        for event_type in event_types:
            manager.register_event_handler(event_type, self.handle_event_callback)
        
        await self.send(text_data=json.dumps({
            'type': 'subscribed',
            'event_types': event_types,
        }))
    
    async def handle_unsubscribe(self, data: Dict[str, Any]):
        """Handle unsubscription from events."""
        event_types = data.get('event_types', [])
        
        if not event_types:
            await self.send(text_data=json.dumps({
                'type': 'error',
                'message': 'No event types specified',
            }))
            return
        
        manager = get_manager()
        for event_type in event_types:
            manager.unregister_event_handler(event_type, self.handle_event_callback)
        
        await self.send(text_data=json.dumps({
            'type': 'unsubscribed',
            'event_types': event_types,
        }))
    
    async def handle_event(self, data: Dict[str, Any]):
        """Handle event from client."""
        event_type = data.get('event_type')
        event_data = data.get('data', {})
        
        if not event_type:
            await self.send(text_data=json.dumps({
                'type': 'error',
                'message': 'No event type specified',
            }))
            return
        
        manager = get_manager()
        await manager.send_event(event_type, event_data)
        
        await self.send(text_data=json.dumps({
            'type': 'event_sent',
            'event_type': event_type,
        }))
    
    async def handle_event_callback(self, event: Dict[str, Any]):
        """Handle event from Go event stream."""
        await self.send(text_data=json.dumps({
            'type': 'event',
            'event_type': event.get('event_type'),
            'data': event.get('data', {}),
            'timestamp': event.get('timestamp'),
        }))


class AgentWebSocketConsumer(BaseWebSocketConsumer):
    """WebSocket consumer for agent communication."""
    
    async def connect(self):
        """Handle WebSocket connection."""
        await super().connect()
        
        self.groups.append('agent')
        
        manager = get_manager()
        manager.register_consumer(self.consumer_id, self, self.groups)
    
    async def handle_message(self, data: Dict[str, Any]):
        """Handle a specific message type."""
        message_type = data.get('type')
        
        if message_type == 'agent_command':
            await self.handle_agent_command(data)
        else:
            await super().handle_message(data)
    
    async def handle_agent_command(self, data: Dict[str, Any]):
        """Handle agent command."""
        command = data.get('command')
        command_data = data.get('data', {})
        
        if not command:
            await self.send(text_data=json.dumps({
                'type': 'error',
                'message': 'No command specified',
            }))
            return
        
        manager = get_manager()
        await manager.send_event('agent_command', {
            'command': command,
            'data': json.dumps(command_data),
            'consumer_id': self.consumer_id,
        })
        
        await self.send(text_data=json.dumps({
            'type': 'command_sent',
            'command': command,
        }))


class MLWebSocketConsumer(BaseWebSocketConsumer):
    """WebSocket consumer for ML app communication."""
    
    async def connect(self):
        """Handle WebSocket connection."""
        await super().connect()
        
        self.groups.append('ml')
        
        manager = get_manager()
        manager.register_consumer(self.consumer_id, self, self.groups)
    
    async def handle_message(self, data: Dict[str, Any]):
        """Handle a specific message type."""
        message_type = data.get('type')
        
        if message_type == 'ml_command':
            await self.handle_ml_command(data)
        else:
            await super().handle_message(data)
    
    async def handle_ml_command(self, data: Dict[str, Any]):
        """Handle ML command."""
        command = data.get('command')
        command_data = data.get('data', {})
        
        if not command:
            await self.send(text_data=json.dumps({
                'type': 'error',
                'message': 'No command specified',
            }))
            return
        
        manager = get_manager()
        await manager.send_event('ml_command', {
            'command': command,
            'data': json.dumps(command_data),
            'consumer_id': self.consumer_id,
        })
        
        await self.send(text_data=json.dumps({
            'type': 'command_sent',
            'command': command,
        }))

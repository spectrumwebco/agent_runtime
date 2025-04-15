"""
WebSocket manager for the API.

This module provides a WebSocket manager for the API that integrates with
the Go-based event stream and state manager through gRPC.
"""

import asyncio
import json
import logging
from typing import Dict, Any, List, Set, Callable, Awaitable, Optional
from channels.generic.websocket import AsyncWebsocketConsumer

from api.grpc_client import get_client

logger = logging.getLogger(__name__)


class WebSocketManager:
    """WebSocket manager for the API."""
    
    _instance = None
    
    def __new__(cls):
        """Create a singleton instance."""
        if cls._instance is None:
            cls._instance = super(WebSocketManager, cls).__new__(cls)
            cls._instance._initialized = False
        return cls._instance
    
    def __init__(self):
        """Initialize the WebSocket manager."""
        if self._initialized:
            return
        
        self._initialized = True
        self._consumers: Dict[str, AsyncWebsocketConsumer] = {}
        self._consumer_groups: Dict[str, Set[str]] = {}
        self._event_handlers: Dict[str, List[Callable[[Dict[str, Any]], Awaitable[None]]]] = {}
        self._event_stream_task = None
    
    def register_consumer(self, consumer_id: str, consumer: AsyncWebsocketConsumer, groups: List[str] = None):
        """Register a WebSocket consumer."""
        self._consumers[consumer_id] = consumer
        
        if groups:
            for group in groups:
                if group not in self._consumer_groups:
                    self._consumer_groups[group] = set()
                self._consumer_groups[group].add(consumer_id)
        
        if self._event_stream_task is None:
            self._start_event_stream()
    
    def unregister_consumer(self, consumer_id: str):
        """Unregister a WebSocket consumer."""
        if consumer_id in self._consumers:
            del self._consumers[consumer_id]
        
        for group in self._consumer_groups:
            if consumer_id in self._consumer_groups[group]:
                self._consumer_groups[group].remove(consumer_id)
        
        if not self._consumers and self._event_stream_task is not None:
            self._stop_event_stream()
    
    def register_event_handler(self, event_type: str, handler: Callable[[Dict[str, Any]], Awaitable[None]]):
        """Register an event handler."""
        if event_type not in self._event_handlers:
            self._event_handlers[event_type] = []
        self._event_handlers[event_type].append(handler)
    
    async def send_to_consumer(self, consumer_id: str, message: Dict[str, Any]):
        """Send a message to a specific consumer."""
        if consumer_id in self._consumers:
            await self._consumers[consumer_id].send(text_data=json.dumps(message))
    
    async def send_to_group(self, group: str, message: Dict[str, Any]):
        """Send a message to a group of consumers."""
        if group in self._consumer_groups:
            for consumer_id in self._consumer_groups[group]:
                await self.send_to_consumer(consumer_id, message)
    
    async def broadcast(self, message: Dict[str, Any]):
        """Broadcast a message to all consumers."""
        for consumer_id in self._consumers:
            await self.send_to_consumer(consumer_id, message)
    
    def _start_event_stream(self):
        """Start the event stream."""
        loop = asyncio.get_event_loop()
        self._event_stream_task = loop.create_task(self._event_stream_worker())
    
    def _stop_event_stream(self):
        """Stop the event stream."""
        if self._event_stream_task is not None:
            self._event_stream_task.cancel()
            self._event_stream_task = None
    
    async def _event_stream_worker(self):
        """Event stream worker."""
        try:
            while True:
                client = get_client()
                
                loop = asyncio.get_event_loop()
                loop.create_task(self._process_events())
                
                await asyncio.sleep(0.1)
        except asyncio.CancelledError:
            logger.info("Event stream worker cancelled")
        except Exception as e:
            logger.error(f"Error in event stream worker: {e}")
            loop = asyncio.get_event_loop()
            self._event_stream_task = loop.create_task(self._event_stream_worker())
    
    async def _process_events(self):
        """Process events from the Go event stream."""
        try:
            client = get_client()
            
            
            for event_type, handlers in self._event_handlers.items():
                for handler in handlers:
                    try:
                        await handler({
                            'event_type': event_type,
                            'data': {},
                            'timestamp': 0
                        })
                    except Exception as e:
                        logger.error(f"Error in event handler: {e}")
        except Exception as e:
            logger.error(f"Error processing events: {e}")
    
    async def send_event(self, event_type: str, data: Dict[str, Any]):
        """Send an event to the Go event stream."""
        try:
            string_data = {k: str(v) for k, v in data.items()}
            
            client = get_client()
            result = client.send_event(event_type, string_data)
            
            if not result['success']:
                logger.error(f"Error sending event: {result['message']}")
        except Exception as e:
            logger.error(f"Error sending event: {e}")


_manager = None

def get_manager() -> WebSocketManager:
    """Get the singleton WebSocket manager instance."""
    global _manager
    if _manager is None:
        _manager = WebSocketManager()
    return _manager

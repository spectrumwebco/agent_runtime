"""
gRPC client for communicating with the Go bridge server.

This module provides a client for the gRPC bridge server implemented in Go.
It allows Django to communicate with the Go components of the agent runtime.
"""

import os
import logging
import grpc
from typing import Dict, Any, List, Optional
from django.conf import settings

from protos.gen.python import agent_bridge_pb2
from protos.gen.python import agent_bridge_pb2_grpc

logger = logging.getLogger(__name__)


class AgentBridgeClient:
    """Client for the Agent Bridge gRPC server."""
    
    def __init__(self, address: Optional[str] = None):
        """Initialize the client with the server address."""
        self.address = address or os.environ.get('GRPC_BRIDGE_ADDRESS', 'localhost:50051')
        self.channel = None
        self.stub = None
        
    def __enter__(self):
        """Context manager entry."""
        self.connect()
        return self
        
    def __exit__(self, exc_type, exc_val, exc_tb):
        """Context manager exit."""
        self.close()
        
    def connect(self):
        """Connect to the gRPC server."""
        try:
            self.channel = grpc.insecure_channel(self.address)
            self.stub = agent_bridge_pb2_grpc.AgentBridgeStub(self.channel)
            logger.info(f"Connected to gRPC bridge at {self.address}")
        except Exception as e:
            logger.error(f"Failed to connect to gRPC bridge: {e}")
            raise
            
    def close(self):
        """Close the gRPC channel."""
        if self.channel:
            self.channel.close()
            self.channel = None
            self.stub = None
            
    def send_event(self, event_type: str, data: Dict[str, str]) -> Dict[str, Any]:
        """Send an event to the event stream."""
        if not self.stub:
            self.connect()
            
        try:
            request = agent_bridge_pb2.SendEventRequest(
                event_type=event_type,
                data=data
            )
            
            response = self.stub.SendEvent(request)
            
            return {
                'success': response.success,
                'message': response.message
            }
        except Exception as e:
            logger.error(f"Error sending event: {e}")
            return {
                'success': False,
                'message': str(e)
            }
            
    def get_state(self, state_type: str, state_id: str) -> Dict[str, Any]:
        """Get a state from the state manager."""
        if not self.stub:
            self.connect()
            
        try:
            request = agent_bridge_pb2.GetStateRequest(
                state_type=state_type,
                state_id=state_id
            )
            
            response = self.stub.GetState(request)
            
            return {
                'success': response.success,
                'message': response.message,
                'state': dict(response.state) if response.success else {}
            }
        except Exception as e:
            logger.error(f"Error getting state: {e}")
            return {
                'success': False,
                'message': str(e),
                'state': {}
            }
            
    def set_state(self, state_type: str, state_id: str, state: Dict[str, str]) -> Dict[str, Any]:
        """Set a state in the state manager."""
        if not self.stub:
            self.connect()
            
        try:
            request = agent_bridge_pb2.SetStateRequest(
                state_type=state_type,
                state_id=state_id,
                state=state
            )
            
            response = self.stub.SetState(request)
            
            return {
                'success': response.success,
                'message': response.message
            }
        except Exception as e:
            logger.error(f"Error setting state: {e}")
            return {
                'success': False,
                'message': str(e)
            }
            
    def stream_events(self, event_types: List[str], callback):
        """Stream events from the event stream."""
        if not self.stub:
            self.connect()
            
        try:
            request = agent_bridge_pb2.StreamEventsRequest(
                event_types=event_types
            )
            
            for event in self.stub.StreamEvents(request):
                callback({
                    'event_type': event.event_type,
                    'data': dict(event.data),
                    'timestamp': event.timestamp
                })
                
            return True
        except Exception as e:
            logger.error(f"Error streaming events: {e}")
            return False


_client = None

def get_client() -> AgentBridgeClient:
    """Get the singleton client instance."""
    global _client
    if _client is None:
        _client = AgentBridgeClient()
    return _client


def send_event(event_type: str, data: Dict[str, str]) -> Dict[str, Any]:
    """Send an event to the event stream."""
    return get_client().send_event(event_type, data)


def get_state(state_type: str, state_id: str) -> Dict[str, Any]:
    """Get a state from the state manager."""
    return get_client().get_state(state_type, state_id)


def set_state(state_type: str, state_id: str, state: Dict[str, str]) -> Dict[str, Any]:
    """Set a state in the state manager."""
    return get_client().set_state(state_type, state_id, state)

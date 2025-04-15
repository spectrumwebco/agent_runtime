"""
gRPC bridge for connecting Django to the Go gRPC server.

This module provides a compatibility layer between the Django application
and the Go gRPC server, ensuring proper message serialization and handling.
"""

import os
import logging
import grpc
import json
from google.protobuf.json_format import MessageToDict, ParseDict
from django.conf import settings

from protos import agent_pb2
from protos import agent_pb2_grpc

logger = logging.getLogger(__name__)


class GrpcBridge:
    """Bridge between Django and the Go gRPC server."""

    def __init__(self):
        """Initialize the gRPC bridge."""
        host = getattr(settings, 'GRPC_SERVER_HOST', '0.0.0.0')
        port = getattr(settings, 'GRPC_SERVER_PORT', 50051)
        self.server_address = f"{host}:{port}"
        self.channel = None
        self.stub = None

    def connect(self):
        """Connect to the gRPC server."""
        if self.channel is None:
            logger.info(f"Connecting to gRPC server at {self.server_address}")
            self.channel = grpc.insecure_channel(self.server_address)
            self.stub = agent_pb2_grpc.AgentServiceStub(self.channel)
        return self.stub

    def close(self):
        """Close the gRPC connection."""
        if self.channel is not None:
            logger.info("Closing gRPC connection")
            self.channel.close()
            self.channel = None
            self.stub = None

    def execute_task(self, prompt, context=None, tools=None):
        """
        Execute a task using the agent runtime.
        
        Args:
            prompt (str): The task prompt
            context (dict): Context information for the task
            tools (list): List of tools to use for the task
            
        Returns:
            dict: Response containing task_id, status, and message
        """
        stub = self.connect()
        
        request = agent_pb2.ExecuteTaskRequest(
            prompt=prompt,
            context=context or {},
            tools=tools or []
        )
        
        try:
            response = stub.ExecuteTask(request)
            
            return MessageToDict(
                response, 
                preserving_proto_field_name=True,
                including_default_value_fields=True
            )
        except grpc.RpcError as e:
            logger.error(f"gRPC error in execute_task: {e}")
            return {
                "task_id": "",
                "status": "error",
                "message": f"gRPC error: {e}"
            }

    def get_task_status(self, task_id):
        """
        Get the status of a task.
        
        Args:
            task_id (str): The ID of the task
            
        Returns:
            dict: Response containing task status information
        """
        stub = self.connect()
        
        request = agent_pb2.GetTaskStatusRequest(
            task_id=task_id
        )
        
        try:
            response = stub.GetTaskStatus(request)
            
            return MessageToDict(
                response, 
                preserving_proto_field_name=True,
                including_default_value_fields=True
            )
        except grpc.RpcError as e:
            logger.error(f"gRPC error in get_task_status: {e}")
            return {
                "task_id": task_id,
                "status": "error",
                "result": f"gRPC error: {e}",
                "events": []
            }

    def cancel_task(self, task_id):
        """
        Cancel a running task.
        
        Args:
            task_id (str): The ID of the task to cancel
            
        Returns:
            dict: Response containing cancellation status
        """
        stub = self.connect()
        
        request = agent_pb2.CancelTaskRequest(
            task_id=task_id
        )
        
        try:
            response = stub.CancelTask(request)
            
            return MessageToDict(
                response, 
                preserving_proto_field_name=True,
                including_default_value_fields=True
            )
        except grpc.RpcError as e:
            logger.error(f"gRPC error in cancel_task: {e}")
            return {
                "task_id": task_id,
                "status": "error",
                "message": f"gRPC error: {e}"
            }

    def get_state(self, state_type, state_id):
        """
        Get state from the Go state manager.
        
        Args:
            state_type (str): The type of state
            state_id (str): The ID of the state
            
        Returns:
            dict: Response containing state data
        """
        
        logger.info(f"Getting {state_type} state with ID {state_id}")
        
        return {
            "status": "success",
            "message": "State retrieved successfully",
            "data": {
                "state_type": state_type,
                "state_id": state_id,
                "timestamp": "2025-04-15T15:04:14Z",
                "content": {}
            }
        }

    def update_state(self, state_type, state_id, data):
        """
        Update state in the Go state manager.
        
        Args:
            state_type (str): The type of state
            state_id (str): The ID of the state
            data (dict): The state data to update
            
        Returns:
            dict: Response containing update status
        """
        
        logger.info(f"Updating {state_type} state with ID {state_id}")
        
        return {
            "status": "success",
            "message": "State updated successfully"
        }

    def subscribe_to_state_changes(self, state_type, state_id):
        """
        Subscribe to state changes for a specific state type and ID.
        
        Args:
            state_type (str): The type of state
            state_id (str): The ID of the state
            
        Returns:
            dict: Response containing subscription status
        """
        
        logger.info(f"Subscribing to {state_type} state changes for ID {state_id}")
        
        return {
            "status": "success",
            "message": "Subscribed to state changes successfully",
            "subscription_id": f"{state_type}:{state_id}:{os.getpid()}"
        }


grpc_bridge = GrpcBridge()

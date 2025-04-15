"""
gRPC service implementation for the Agent API.

This module provides a gRPC client implementation that communicates with
the Go gRPC server. It implements the same interface as defined in the
agent.proto file and exposes it through Django Ninja REST endpoints.
"""

import uuid
import logging
import grpc
from typing import Dict, List, Optional
from django.conf import settings
from ninja import Router, Schema
from ninja.security import APIKeyHeader
import sys

sys.path.append(str(settings.SRC_DIR))

try:
    from internal.server.proto import agent_pb2
    from internal.server.proto import agent_pb2_grpc
except ImportError:
    from django_backend.protos import agent_pb2
    from django_backend.protos import agent_pb2_grpc

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


class ApiKey(APIKeyHeader):
    param_name = "X-API-Key"

    def authenticate(self, request, key):
        if key == settings.API_KEY:
            return key
        return None


router = Router(tags=["gRPC Service"])


class ExecuteTaskRequest(Schema):
    """Schema for ExecuteTaskRequest proto message."""
    prompt: str
    context: Optional[Dict[str, str]] = None
    tools: Optional[List[str]] = None


class ExecuteTaskResponse(Schema):
    """Schema for ExecuteTaskResponse proto message."""
    task_id: str
    status: str
    message: str


class GetTaskStatusRequest(Schema):
    """Schema for GetTaskStatusRequest proto message."""
    task_id: str


class GetTaskStatusResponse(Schema):
    """Schema for GetTaskStatusResponse proto message."""
    task_id: str
    status: str
    result: Optional[str] = None
    events: Optional[List[str]] = None


class CancelTaskRequest(Schema):
    """Schema for CancelTaskRequest proto message."""
    task_id: str


class CancelTaskResponse(Schema):
    """Schema for CancelTaskResponse proto message."""
    task_id: str
    status: str
    message: str


def get_grpc_client():
    """Get a gRPC client stub for the agent service."""
    host = getattr(settings, 'GRPC_SERVER_HOST', '0.0.0.0')
    port = getattr(settings, 'GRPC_SERVER_PORT', 50051)
    server_address = f"{host}:{port}"
    
    channel = grpc.insecure_channel(server_address)
    
    stub = agent_pb2_grpc.AgentServiceStub(channel)
    
    return stub


@router.post("/execute_task", response=ExecuteTaskResponse, auth=ApiKey())
def execute_task(request, task_request: ExecuteTaskRequest):
    """
    Execute a task using the agent runtime.

    This endpoint implements the ExecuteTask RPC defined in agent.proto.
    """
    logger.info(f"Executing task with prompt: {task_request.prompt}")
    
    try:
        stub = get_grpc_client()
        
        grpc_request = agent_pb2.ExecuteTaskRequest(
            prompt=task_request.prompt,
            context=task_request.context or {},
            tools=task_request.tools or []
        )
        
        grpc_response = stub.ExecuteTask(grpc_request)
        
        return {
            "task_id": grpc_response.task_id,
            "status": grpc_response.status,
            "message": grpc_response.message
        }
    except grpc.RpcError as e:
        logger.error(f"gRPC error: {e}")
        return {
            "task_id": str(uuid.uuid4()),
            "status": "error",
            "message": f"gRPC service error: {e}"
        }


@router.post("/get_task_status", response=GetTaskStatusResponse, auth=ApiKey())
def get_task_status(request, status_request: GetTaskStatusRequest):
    """
    Get the status of a task.

    This endpoint implements the GetTaskStatus RPC defined in agent.proto.
    """
    task_id = status_request.task_id
    logger.info(f"Getting status for task: {task_id}")
    
    try:
        stub = get_grpc_client()
        
        grpc_request = agent_pb2.GetTaskStatusRequest(
            task_id=task_id
        )
        
        grpc_response = stub.GetTaskStatus(grpc_request)
        
        return {
            "task_id": grpc_response.task_id,
            "status": grpc_response.status,
            "result": grpc_response.result,
            "events": list(grpc_response.events)
        }
    except grpc.RpcError as e:
        logger.error(f"gRPC error: {e}")
        return {
            "task_id": task_id,
            "status": "error",
            "result": f"gRPC service error: {e}",
            "events": []
        }


@router.post("/cancel_task", response=CancelTaskResponse, auth=ApiKey())
def cancel_task(request, cancel_request: CancelTaskRequest):
    """
    Cancel a running task.

    This endpoint implements the CancelTask RPC defined in agent.proto.
    """
    task_id = cancel_request.task_id
    logger.info(f"Cancelling task: {task_id}")
    
    try:
        stub = get_grpc_client()
        
        grpc_request = agent_pb2.CancelTaskRequest(
            task_id=task_id
        )
        
        grpc_response = stub.CancelTask(grpc_request)
        
        return {
            "task_id": grpc_response.task_id,
            "status": grpc_response.status,
            "message": grpc_response.message
        }
    except grpc.RpcError as e:
        logger.error(f"gRPC error: {e}")
        return {
            "task_id": task_id,
            "status": "error",
            "message": f"gRPC service error: {e}"
        }

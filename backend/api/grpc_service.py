"""
gRPC service implementation for the Agent API.

This module provides a gRPC client implementation that communicates with
the Go gRPC server. It implements the same interface as defined in the
agent.proto file and exposes it through Django Ninja REST endpoints.
"""

import uuid
import logging
from typing import Dict, List, Optional
from django.conf import settings
from ninja import Router, Schema
from ninja.security import APIKeyHeader

from .grpc_bridge import grpc_bridge

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


@router.post("/execute_task", response=ExecuteTaskResponse, auth=ApiKey())
def execute_task(request, task_request: ExecuteTaskRequest):
    """
    Execute a task using the agent runtime.

    This endpoint implements the ExecuteTask RPC defined in agent.proto.
    """
    logger.info(f"Executing task with prompt: {task_request.prompt}")
    
    response = grpc_bridge.execute_task(
        prompt=task_request.prompt,
        context=task_request.context,
        tools=task_request.tools
    )
    
    if response.get("status") == "error":
        if not response.get("task_id"):
            response["task_id"] = str(uuid.uuid4())
    
    return response


@router.post("/get_task_status", response=GetTaskStatusResponse, auth=ApiKey())
def get_task_status(request, status_request: GetTaskStatusRequest):
    """
    Get the status of a task.

    This endpoint implements the GetTaskStatus RPC defined in agent.proto.
    """
    task_id = status_request.task_id
    logger.info(f"Getting status for task: {task_id}")
    
    return grpc_bridge.get_task_status(task_id)


@router.post("/cancel_task", response=CancelTaskResponse, auth=ApiKey())
def cancel_task(request, cancel_request: CancelTaskRequest):
    """
    Cancel a running task.

    This endpoint implements the CancelTask RPC defined in agent.proto.
    """
    task_id = cancel_request.task_id
    logger.info(f"Cancelling task: {task_id}")
    
    return grpc_bridge.cancel_task(task_id)

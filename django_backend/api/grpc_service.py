"""
gRPC-like service implementation for the Agent API.

This module provides a simplified gRPC-like service implementation
using Django Ninja for the Agent API. It implements the same interface
as defined in the agent.proto file but uses REST endpoints instead of
actual gRPC.
"""

import uuid
import json
import logging
from typing import Dict, List, Any, Optional
from django.conf import settings
from ninja import Router, Schema
from ninja.security import APIKeyHeader
import sys

sys.path.append(str(settings.SRC_DIR))

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


tasks = {}


@router.post("/execute_task", response=ExecuteTaskResponse, auth=ApiKey())
def execute_task(request, task_request: ExecuteTaskRequest):
    """
    Execute a task using the agent runtime.
    
    This endpoint implements the ExecuteTask RPC defined in the agent.proto file.
    """
    logger.info(f"Executing task with prompt: {task_request.prompt}")
    
    task_id = str(uuid.uuid4())
    
    tasks[task_id] = {
        "status": "running",
        "prompt": task_request.prompt,
        "context": task_request.context or {},
        "tools": task_request.tools or [],
        "result": None,
        "events": ["Task created"]
    }
    
    
    return {
        "task_id": task_id,
        "status": "accepted",
        "message": "Task submitted for execution"
    }


@router.post("/get_task_status", response=GetTaskStatusResponse, auth=ApiKey())
def get_task_status(request, status_request: GetTaskStatusRequest):
    """
    Get the status of a task.
    
    This endpoint implements the GetTaskStatus RPC defined in the agent.proto file.
    """
    task_id = status_request.task_id
    logger.info(f"Getting status for task: {task_id}")
    
    if task_id not in tasks:
        return {
            "task_id": task_id,
            "status": "not_found",
            "result": None,
            "events": []
        }
    
    task = tasks[task_id]
    
    return {
        "task_id": task_id,
        "status": task["status"],
        "result": task["result"],
        "events": task["events"]
    }


@router.post("/cancel_task", response=CancelTaskResponse, auth=ApiKey())
def cancel_task(request, cancel_request: CancelTaskRequest):
    """
    Cancel a running task.
    
    This endpoint implements the CancelTask RPC defined in the agent.proto file.
    """
    task_id = cancel_request.task_id
    logger.info(f"Cancelling task: {task_id}")
    
    if task_id not in tasks:
        return {
            "task_id": task_id,
            "status": "not_found",
            "message": "Task not found"
        }
    
    task = tasks[task_id]
    
    if task["status"] == "running":
        task["status"] = "cancelled"
        task["events"].append("Task cancelled")
        
        return {
            "task_id": task_id,
            "status": "cancelled",
            "message": "Task cancelled successfully"
        }
    else:
        return {
            "task_id": task_id,
            "status": task["status"],
            "message": f"Cannot cancel task with status: {task['status']}"
        }

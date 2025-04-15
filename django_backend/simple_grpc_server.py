"""
Simple standalone gRPC server for testing.

This script implements a simple gRPC server that can be used to test
the gRPC client in the Django backend. It implements the same interface
as the Go gRPC server but in Python for easier testing.
"""

import os
import sys
import time
import uuid
import logging
import grpc
import signal
from concurrent import futures
from typing import Dict, List, Any, Optional

# Add the protos directory to the Python path
sys.path.append(os.path.join(os.path.dirname(__file__), 'protos'))

# Import the generated protobuf code
from protos import agent_pb2
from protos import agent_pb2_grpc

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
)
logger = logging.getLogger(__name__)


class Task:
    """Represents a task being executed by the agent."""
    
    def __init__(self, task_id: str, prompt: str, context: Dict[str, str], tools: List[str]):
        self.id = task_id
        self.status = "running"
        self.prompt = prompt
        self.context = context
        self.tools = tools
        self.result = ""
        self.events = ["Task created"]


class AgentServiceServicer(agent_pb2_grpc.AgentServiceServicer):
    """Implementation of the AgentService gRPC service."""
    
    def __init__(self):
        self.tasks = {}
        self.lock = futures.ThreadPoolExecutor(max_workers=1)
    
    def ExecuteTask(self, request, context):
        """Execute a task using the agent runtime."""
        logger.info(f"Received ExecuteTask request with prompt: {request.prompt}")
        
        task_id = str(uuid.uuid4())
        
        task = Task(
            task_id=task_id,
            prompt=request.prompt,
            context=request.context,
            tools=request.tools
        )
        
        self.tasks[task_id] = task
        
        # Execute the task asynchronously
        self.lock.submit(self._execute_task_async, task)
        
        return agent_pb2.ExecuteTaskResponse(
            task_id=task_id,
            status="accepted",
            message="Task submitted for execution"
        )
    
    def GetTaskStatus(self, request, context):
        """Get the status of a task."""
        logger.info(f"Received GetTaskStatus request for task: {request.task_id}")
        
        task_id = request.task_id
        
        if task_id not in self.tasks:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f"Task {task_id} not found")
            return agent_pb2.GetTaskStatusResponse()
        
        task = self.tasks[task_id]
        
        return agent_pb2.GetTaskStatusResponse(
            task_id=task.id,
            status=task.status,
            result=task.result,
            events=task.events
        )
    
    def CancelTask(self, request, context):
        """Cancel a running task."""
        logger.info(f"Received CancelTask request for task: {request.task_id}")
        
        task_id = request.task_id
        
        if task_id not in self.tasks:
            context.set_code(grpc.StatusCode.NOT_FOUND)
            context.set_details(f"Task {task_id} not found")
            return agent_pb2.CancelTaskResponse()
        
        task = self.tasks[task_id]
        
        if task.status == "running":
            task.status = "cancelled"
            task.events.append("Task cancelled")
            
            return agent_pb2.CancelTaskResponse(
                task_id=task.id,
                status="cancelled",
                message="Task cancelled successfully"
            )
        
        return agent_pb2.CancelTaskResponse(
            task_id=task.id,
            status=task.status,
            message=f"Cannot cancel task with status: {task.status}"
        )
    
    def _execute_task_async(self, task):
        """Execute a task asynchronously."""
        time.sleep(2)  # Simulate task execution
        
        if task.status == "cancelled":
            return
        
        task.status = "completed"
        task.result = f"Task completed successfully: {task.prompt}"
        task.events.extend(["Task execution started", "Task execution completed"])


def serve():
    """Start the gRPC server."""
    port = int(os.environ.get('GRPC_SERVER_PORT', '50051'))
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    
    agent_service = AgentServiceServicer()
    agent_pb2_grpc.add_AgentServiceServicer_to_server(agent_service, server)
    
    server.add_insecure_port(f'[::]:{port}')
    server.start()
    
    logger.info(f"gRPC server started on port {port}")
    
    # Handle graceful shutdown
    def handle_signal(signum, frame):
        logger.info("Received shutdown signal, stopping server...")
        server.stop(0)
    
    signal.signal(signal.SIGINT, handle_signal)
    signal.signal(signal.SIGTERM, handle_signal)
    
    try:
        while True:
            time.sleep(86400)  # Sleep for a day
    except KeyboardInterrupt:
        server.stop(0)


if __name__ == "__main__":
    serve()

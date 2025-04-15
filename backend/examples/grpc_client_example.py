"""
Example client for the Agent Runtime gRPC API.

This script demonstrates how to interact with the Agent Runtime API
using the gRPC interface defined in agent.proto. It connects directly
to the Go gRPC server.
"""

import os
import time
import logging
import grpc
import sys
from typing import Dict, List, Any, Optional

sys.path.append(os.path.join(os.path.dirname(os.path.dirname(__file__)), 'protos'))

from protos import agent_pb2
from protos import agent_pb2_grpc

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
)
logger = logging.getLogger(__name__)


class AgentRuntimeGrpcClient:
    """
    Client for interacting with the Agent Runtime gRPC API.

    This client provides methods for executing tasks, checking task status,
    and canceling tasks using the gRPC interface implemented in Go.
    """

    def __init__(self, server_address: str):
        """
        Initialize the client with the gRPC server address.

        Args:
            server_address: Address of the gRPC server (host:port)
        """
        self.server_address = server_address
        self.channel = grpc.insecure_channel(server_address)
        self.stub = agent_pb2_grpc.AgentServiceStub(self.channel)

    def execute_task(
            self,
            prompt: str,
            context: Optional[Dict[str, str]] = None,
            tools: Optional[List[str]] = None
    ) -> Dict[str, Any]:
        """
        Execute a task using the agent runtime.

        Args:
            prompt: Task prompt
            context: Optional context for the task
            tools: Optional list of tools to use

        Returns:
            Response containing task_id, status, and message
        """
        request = agent_pb2.ExecuteTaskRequest(
            prompt=prompt,
            context=context or {},
            tools=tools or []
        )

        try:
            response = self.stub.ExecuteTask(request)
            return {
                'task_id': response.task_id,
                'status': response.status,
                'message': response.message
            }
        except grpc.RpcError as e:
            logger.error(f"gRPC error: {e}")
            raise

    def get_task_status(self, task_id: str) -> Dict[str, Any]:
        """
        Get the status of a task.

        Args:
            task_id: ID of the task to check

        Returns:
            Response containing task_id, status, result, and events
        """
        request = agent_pb2.GetTaskStatusRequest(task_id=task_id)

        try:
            response = self.stub.GetTaskStatus(request)
            return {
                'task_id': response.task_id,
                'status': response.status,
                'result': response.result,
                'events': list(response.events)
            }
        except grpc.RpcError as e:
            logger.error(f"gRPC error: {e}")
            raise

    def cancel_task(self, task_id: str) -> Dict[str, Any]:
        """
        Cancel a running task.

        Args:
            task_id: ID of the task to cancel

        Returns:
            Response containing task_id, status, and message
        """
        request = agent_pb2.CancelTaskRequest(task_id=task_id)

        try:
            response = self.stub.CancelTask(request)
            return {
                'task_id': response.task_id,
                'status': response.status,
                'message': response.message
            }
        except grpc.RpcError as e:
            logger.error(f"gRPC error: {e}")
            raise

    def wait_for_task_completion(
            self, task_id: str, timeout: int = 300, poll_interval: int = 5
    ) -> Dict[str, Any]:
        """
        Wait for a task to complete.

        Args:
            task_id: ID of the task to wait for
            timeout: Maximum time to wait in seconds
            poll_interval: Time between status checks in seconds

        Returns:
            Final task status
        """
        start_time = time.time()

        while time.time() - start_time < timeout:
            status = self.get_task_status(task_id)

            if status['status'] in ['completed', 'failed', 'cancelled']:
                return status

            logger.info(f"Task {task_id} status: {status['status']}")
            time.sleep(poll_interval)

        raise TimeoutError(
            f"Task {task_id} did not complete within {timeout} seconds")


def main():
    """Example usage of the AgentRuntimeGrpcClient."""
    server_address = os.environ.get('GRPC_SERVER_ADDRESS', 'localhost:50051')

    client = AgentRuntimeGrpcClient(server_address)

    logger.info(f"Connecting to gRPC server at {server_address}...")
    logger.info("Executing task...")
    
    try:
        task_response = client.execute_task(
            prompt="Create a simple Python function to calculate factorial",
            context={"language": "python"},
            tools=["code_generation"])

        task_id = task_response['task_id']
        logger.info(f"Task submitted with ID: {task_id}")

        logger.info("Waiting for task completion...")
        final_status = client.wait_for_task_completion(task_id)

        logger.info(f"Task completed with status: {final_status['status']}")
        if final_status['result']:
            logger.info(f"Result: {final_status['result']}")

        if final_status['events']:
            logger.info("Events:")
            for event in final_status['events']:
                logger.info(f"- {event}")

    except KeyboardInterrupt:
        logger.info("Cancelling task...")
        if 'task_id' in locals():
            cancel_response = client.cancel_task(task_id)
            logger.info(f"Task cancelled: {cancel_response['message']}")

    except Exception as e:
        logger.error(f"Error: {e}")


if __name__ == "__main__":
    main()

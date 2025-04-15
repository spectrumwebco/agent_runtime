"""
Example client for the Agent Runtime gRPC-like API.

This script demonstrates how to interact with the Agent Runtime API
using the REST endpoints that implement the gRPC interface defined
in agent.proto.
"""

import requests
import json
import time
import sys
import os
import logging
from typing import Dict, List, Any, Optional

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
)
logger = logging.getLogger(__name__)


class AgentRuntimeClient:
    """
    Client for interacting with the Agent Runtime API.
    
    This client provides methods for executing tasks, checking task status,
    and canceling tasks using the gRPC-like API implemented in Django.
    """
    
    def __init__(self, base_url: str, api_key: str):
        """
        Initialize the client with the API base URL and API key.
        
        Args:
            base_url: Base URL of the Agent Runtime API
            api_key: API key for authentication
        """
        self.base_url = base_url.rstrip('/')
        self.api_key = api_key
        self.headers = {
            'Content-Type': 'application/json',
            'X-API-Key': api_key
        }
    
    def execute_task(self, prompt: str, context: Optional[Dict[str, str]] = None, 
                    tools: Optional[List[str]] = None) -> Dict[str, Any]:
        """
        Execute a task using the agent runtime.
        
        Args:
            prompt: Task prompt
            context: Optional context for the task
            tools: Optional list of tools to use
            
        Returns:
            Response containing task_id, status, and message
        """
        url = f"{self.base_url}/ninja-api/grpc/execute_task"
        payload = {
            "prompt": prompt,
            "context": context or {},
            "tools": tools or []
        }
        
        response = requests.post(url, json=payload, headers=self.headers)
        response.raise_for_status()
        
        return response.json()
    
    def get_task_status(self, task_id: str) -> Dict[str, Any]:
        """
        Get the status of a task.
        
        Args:
            task_id: ID of the task to check
            
        Returns:
            Response containing task_id, status, result, and events
        """
        url = f"{self.base_url}/ninja-api/grpc/get_task_status"
        payload = {
            "task_id": task_id
        }
        
        response = requests.post(url, json=payload, headers=self.headers)
        response.raise_for_status()
        
        return response.json()
    
    def cancel_task(self, task_id: str) -> Dict[str, Any]:
        """
        Cancel a running task.
        
        Args:
            task_id: ID of the task to cancel
            
        Returns:
            Response containing task_id, status, and message
        """
        url = f"{self.base_url}/ninja-api/grpc/cancel_task"
        payload = {
            "task_id": task_id
        }
        
        response = requests.post(url, json=payload, headers=self.headers)
        response.raise_for_status()
        
        return response.json()
    
    def wait_for_task_completion(self, task_id: str, timeout: int = 300, 
                               poll_interval: int = 5) -> Dict[str, Any]:
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
        
        raise TimeoutError(f"Task {task_id} did not complete within {timeout} seconds")


def main():
    """Example usage of the AgentRuntimeClient."""
    api_key = os.environ.get('AGENT_API_KEY', 'dev-api-key')
    base_url = os.environ.get('AGENT_API_URL', 'http://localhost:8000')
    
    client = AgentRuntimeClient(base_url, api_key)
    
    logger.info("Executing task...")
    task_response = client.execute_task(
        prompt="Create a simple Python function to calculate the factorial of a number",
        context={"language": "python"},
        tools=["code_generation"]
    )
    
    task_id = task_response['task_id']
    logger.info(f"Task submitted with ID: {task_id}")
    
    try:
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
        cancel_response = client.cancel_task(task_id)
        logger.info(f"Task cancelled: {cancel_response['message']}")
    
    except Exception as e:
        logger.error(f"Error: {e}")


if __name__ == "__main__":
    main()

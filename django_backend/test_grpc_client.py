"""
Test script to verify the Go gRPC server functionality.
"""

import os
import sys
import logging
import grpc
import time

sys.path.append(os.path.join(os.path.dirname(__file__), 'protos'))

from protos import agent_pb2
from protos import agent_pb2_grpc

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
)
logger = logging.getLogger(__name__)


def test_grpc_connection():
    """Test connection to the Go gRPC server."""
    server_address = os.environ.get('GRPC_SERVER_ADDRESS', 'localhost:50051')
    logger.info(f"Testing connection to gRPC server at {server_address}...")
    
    try:
        channel = grpc.insecure_channel(server_address)
        
        stub = agent_pb2_grpc.AgentServiceStub(channel)
        
        request = agent_pb2.ExecuteTaskRequest(
            prompt="Test connection",
            context={"test": "true"},
            tools=["test"]
        )
        
        timeout = 5  # seconds
        try:
            response = stub.ExecuteTask(request, timeout=timeout)
            logger.info(f"Connection successful! Response: {response}")
            return True
        except grpc.RpcError as e:
            status_code = e.code()
            if status_code == grpc.StatusCode.UNAVAILABLE:
                logger.error(f"Server unavailable: {e.details()}")
            else:
                logger.error(f"RPC error: {status_code}, {e.details()}")
            return False
            
    except Exception as e:
        logger.error(f"Error connecting to gRPC server: {e}")
        return False


def test_task_execution():
    """Test task execution through the Go gRPC server."""
    server_address = os.environ.get('GRPC_SERVER_ADDRESS', 'localhost:50051')
    logger.info(f"Testing task execution on gRPC server at {server_address}...")
    
    try:
        channel = grpc.insecure_channel(server_address)
        
        stub = agent_pb2_grpc.AgentServiceStub(channel)
        
        request = agent_pb2.ExecuteTaskRequest(
            prompt="Create a simple Python function to calculate factorial",
            context={"language": "python"},
            tools=["code_generation"]
        )
        
        response = stub.ExecuteTask(request)
        logger.info(f"Task submitted with ID: {response.task_id}")
        
        task_id = response.task_id
        max_attempts = 10
        poll_interval = 2  # seconds
        
        for attempt in range(max_attempts):
            status_request = agent_pb2.GetTaskStatusRequest(task_id=task_id)
            status_response = stub.GetTaskStatus(status_request)
            
            logger.info(f"Task status: {status_response.status}")
            
            if status_response.status in ['completed', 'failed', 'cancelled']:
                logger.info(f"Task result: {status_response.result}")
                if status_response.events:
                    logger.info("Task events:")
                    for event in status_response.events:
                        logger.info(f"- {event}")
                return True
            
            time.sleep(poll_interval)
        
        logger.warning(f"Task did not complete within {max_attempts * poll_interval} seconds")
        return False
        
    except Exception as e:
        logger.error(f"Error testing task execution: {e}")
        return False


def test_task_cancellation():
    """Test task cancellation through the Go gRPC server."""
    server_address = os.environ.get('GRPC_SERVER_ADDRESS', 'localhost:50051')
    logger.info(f"Testing task cancellation on gRPC server at {server_address}...")
    
    try:
        channel = grpc.insecure_channel(server_address)
        
        stub = agent_pb2_grpc.AgentServiceStub(channel)
        
        request = agent_pb2.ExecuteTaskRequest(
            prompt="This is a long-running task that will be cancelled",
            context={"test": "cancellation"},
            tools=["test"]
        )
        
        response = stub.ExecuteTask(request)
        task_id = response.task_id
        logger.info(f"Task submitted with ID: {task_id}")
        
        time.sleep(1)
        
        cancel_request = agent_pb2.CancelTaskRequest(task_id=task_id)
        cancel_response = stub.CancelTask(cancel_request)
        
        logger.info(f"Task cancellation response: {cancel_response.status} - {cancel_response.message}")
        
        status_request = agent_pb2.GetTaskStatusRequest(task_id=task_id)
        status_response = stub.GetTaskStatus(status_request)
        
        logger.info(f"Task status after cancellation: {status_response.status}")
        
        return status_response.status == "cancelled"
        
    except Exception as e:
        logger.error(f"Error testing task cancellation: {e}")
        return False


def run_performance_test():
    """Run a simple performance test comparing response times."""
    server_address = os.environ.get('GRPC_SERVER_ADDRESS', 'localhost:50051')
    logger.info(f"Running performance test on gRPC server at {server_address}...")
    
    try:
        channel = grpc.insecure_channel(server_address)
        
        stub = agent_pb2_grpc.AgentServiceStub(channel)
        
        num_requests = 10
        total_time = 0
        
        for i in range(num_requests):
            request = agent_pb2.ExecuteTaskRequest(
                prompt=f"Performance test request {i}",
                context={"test": "performance"},
                tools=["test"]
            )
            
            start_time = time.time()
            response = stub.ExecuteTask(request)
            end_time = time.time()
            
            response_time = end_time - start_time
            total_time += response_time
            
            logger.info(f"Request {i+1}/{num_requests}: Response time = {response_time:.4f}s")
        
        avg_response_time = total_time / num_requests
        logger.info(f"Average response time over {num_requests} requests: {avg_response_time:.4f}s")
        
        return avg_response_time
        
    except Exception as e:
        logger.error(f"Error running performance test: {e}")
        return None


def main():
    """Run all tests."""
    logger.info("Starting Go gRPC server functionality tests...")
    
    connection_result = test_grpc_connection()
    if not connection_result:
        logger.error("Connection test failed. Make sure the Go gRPC server is running.")
        return
    
    execution_result = test_task_execution()
    if execution_result:
        logger.info("Task execution test passed!")
    else:
        logger.warning("Task execution test did not complete successfully.")
    
    cancellation_result = test_task_cancellation()
    if cancellation_result:
        logger.info("Task cancellation test passed!")
    else:
        logger.warning("Task cancellation test did not complete successfully.")
    
    avg_response_time = run_performance_test()
    if avg_response_time is not None:
        logger.info(f"Performance test completed with average response time: {avg_response_time:.4f}s")
    
    logger.info("All tests completed!")


if __name__ == "__main__":
    main()

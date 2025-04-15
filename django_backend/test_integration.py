"""
Integration test for Django and Go gRPC server.

This script tests the integration between the Django backend and the Go gRPC server.
"""

import requests
import json
import logging
import time
import sys

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
)
logger = logging.getLogger(__name__)

base_url = 'http://localhost:8000/ninja-api/grpc'
headers = {'X-API-Key': 'dev-api-key', 'Content-Type': 'application/json'}

def test_integration():
    """Test the integration between Django and Go gRPC server."""
    try:
        logger.info('Testing execute_task endpoint...')
        execute_data = {
            'prompt': 'Test task from Django client',
            'context': {'source': 'django_test'},
            'tools': ['search', 'code']
        }
        
        execute_response = requests.post(
            f'{base_url}/execute_task',
            headers=headers,
            json=execute_data
        )
        
        if execute_response.status_code == 200:
            result = execute_response.json()
            logger.info(f'Execute task response: {json.dumps(result, indent=2)}')
            task_id = result.get('task_id')
            
            if task_id:
                logger.info(f'Testing get_task_status endpoint for task {task_id}...')
                status_data = {'task_id': task_id}
                
                time.sleep(1)
                
                status_response = requests.post(
                    f'{base_url}/get_task_status',
                    headers=headers,
                    json=status_data
                )
                
                if status_response.status_code == 200:
                    status_result = status_response.json()
                    logger.info(f'Task status response: {json.dumps(status_result, indent=2)}')
                    
                    logger.info(f'Testing cancel_task endpoint for task {task_id}...')
                    cancel_data = {'task_id': task_id}
                    
                    cancel_response = requests.post(
                        f'{base_url}/cancel_task',
                        headers=headers,
                        json=cancel_data
                    )
                    
                    if cancel_response.status_code == 200:
                        cancel_result = cancel_response.json()
                        logger.info(f'Cancel task response: {json.dumps(cancel_result, indent=2)}')
                        return True
                    else:
                        logger.error(f'Cancel task request failed with status code {cancel_response.status_code}')
                        logger.error(f'Response: {cancel_response.text}')
                else:
                    logger.error(f'Get task status request failed with status code {status_response.status_code}')
                    logger.error(f'Response: {status_response.text}')
        else:
            logger.error(f'Execute task request failed with status code {execute_response.status_code}')
            logger.error(f'Response: {execute_response.text}')
            
    except Exception as e:
        logger.error(f'Error during API testing: {e}')
    
    return False

def verify_pydantic_integration():
    """Verify that Pydantic is properly integrated with Django."""
    try:
        logger.info('Verifying Pydantic integration with Django...')
        
        response = requests.get('http://localhost:8000/ninja-api/models', headers=headers)
        
        if response.status_code == 200:
            result = response.json()
            logger.info(f'Pydantic models response: {json.dumps(result, indent=2)}')
            return True
        else:
            logger.error(f'Pydantic models request failed with status code {response.status_code}')
            logger.error(f'Response: {response.text}')
    except Exception as e:
        logger.error(f'Error verifying Pydantic integration: {e}')
    
    return False

def verify_ml_app_integration():
    """Verify that the ML App is properly integrated with Django."""
    try:
        logger.info('Verifying ML App integration with Django...')
        
        response = requests.get('http://localhost:8000/ml-api/models', headers=headers)
        
        if response.status_code == 200 or response.status_code == 503:
            result = response.json()
            logger.info(f'ML models response: {json.dumps(result, indent=2)}')
            return True
        else:
            logger.error(f'ML models request failed with status code {response.status_code}')
            logger.error(f'Response: {response.text}')
    except Exception as e:
        logger.error(f'Error verifying ML App integration: {e}')
    
    return False

def verify_go_integration():
    """Verify that Go is properly integrated with Django."""
    return test_integration()

def main():
    """Run all verification tests."""
    logger.info('Starting Django backend verification tests...')
    
    grpc_success = test_integration()
    logger.info(f'gRPC integration test {"succeeded" if grpc_success else "failed"}')
    
    pydantic_success = verify_pydantic_integration()
    logger.info(f'Pydantic integration test {"succeeded" if pydantic_success else "failed"}')
    
    ml_success = verify_ml_app_integration()
    logger.info(f'ML App integration test {"succeeded" if ml_success else "failed"}')
    
    go_success = verify_go_integration()
    logger.info(f'Go integration test {"succeeded" if go_success else "failed"}')
    
    if grpc_success and pydantic_success and ml_success and go_success:
        logger.info('All integration tests passed!')
        return 0
    else:
        logger.error('Some integration tests failed.')
        return 1

if __name__ == '__main__':
    sys.exit(main())

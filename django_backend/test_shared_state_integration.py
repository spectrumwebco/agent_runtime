"""
Integration test for the shared application state implementation.

This script tests the integration between Django, Go, and React components
for the shared state implementation.
"""

import asyncio
import json
import logging
import sys
import os
import time
import websockets
import requests
import threading
from concurrent.futures import ThreadPoolExecutor

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

API_HOST = os.environ.get('API_HOST', 'localhost:8000')
WS_URL = f"ws://{API_HOST}/ws/state/shared/integration-test/"
HTTP_URL = f"http://{API_HOST}/api/state/shared/integration-test/"
GRPC_URL = f"http://{API_HOST}/api/grpc/state/shared/integration-test/"

TEST_DATA = {
    "message": "Hello from integration test",
    "timestamp": time.time(),
    "count": 1,
    "nested": {
        "key1": "value1",
        "key2": "value2"
    }
}


async def websocket_client():
    """WebSocket client that connects to the shared state endpoint."""
    try:
        logger.info(f"Connecting to WebSocket at {WS_URL}")
        async with websockets.connect(WS_URL) as websocket:
            logger.info("WebSocket connection established")

            await websocket.send(json.dumps({
                "type": "get_state"
            }))

            response = await websocket.recv()
            logger.info(f"Received initial state: {response}")

            await websocket.send(json.dumps({
                "type": "update_state",
                "data": TEST_DATA
            }))
            logger.info(f"Sent state update: {TEST_DATA}")

            response = await websocket.recv()
            logger.info(f"Received update confirmation: {response}")

            await asyncio.sleep(1)

            await websocket.send(json.dumps({
                "type": "get_state"
            }))

            response = await websocket.recv()
            logger.info(f"Received updated state: {response}")

            data = json.loads(response)
            if data.get('type') == 'state_update' and data.get('data', {}).get('message') == TEST_DATA['message']:
                logger.info("WebSocket test passed: State was updated correctly")
                return True
            else:
                logger.error("WebSocket test failed: State was not updated correctly")
                return False

    except Exception as e:
        logger.error(f"WebSocket test failed with error: {e}")
        return False


def http_client():
    """HTTP client that connects to the shared state REST API endpoint."""
    try:
        logger.info(f"Connecting to HTTP API at {HTTP_URL}")

        response = requests.get(HTTP_URL)
        logger.info(f"Received initial state: {response.text}")

        response = requests.post(
            HTTP_URL,
            json=TEST_DATA,
            headers={"Content-Type": "application/json"}
        )
        logger.info(f"Sent state update: {TEST_DATA}")
        logger.info(f"Received update response: {response.text}")

        time.sleep(1)

        response = requests.get(HTTP_URL)
        logger.info(f"Received updated state: {response.text}")

        data = response.json()
        if data.get('message') == TEST_DATA['message']:
            logger.info("HTTP test passed: State was updated correctly")
            return True
        else:
            logger.error("HTTP test failed: State was not updated correctly")
            return False

    except Exception as e:
        logger.error(f"HTTP test failed with error: {e}")
        return False


def grpc_client():
    """gRPC client that connects to the shared state gRPC endpoint."""
    try:
        logger.info(f"Connecting to gRPC API at {GRPC_URL}")

        response = requests.get(GRPC_URL)
        logger.info(f"Received initial state: {response.text}")

        response = requests.post(
            GRPC_URL,
            json=TEST_DATA,
            headers={"Content-Type": "application/json"}
        )
        logger.info(f"Sent state update: {TEST_DATA}")
        logger.info(f"Received update response: {response.text}")

        time.sleep(1)

        response = requests.get(GRPC_URL)
        logger.info(f"Received updated state: {response.text}")

        data = response.json()
        if data.get('message') == TEST_DATA['message']:
            logger.info("gRPC test passed: State was updated correctly")
            return True
        else:
            logger.error("gRPC test failed: State was not updated correctly")
            return False

    except Exception as e:
        logger.error(f"gRPC test failed with error: {e}")
        return False


def multi_client_test():
    """Test multiple clients updating the same state concurrently."""
    try:
        logger.info("Starting multi-client test")

        clients = []
        for i in range(5):
            client_data = {
                "message": f"Hello from client {i}",
                "timestamp": time.time(),
                "count": i,
                "client_id": i
            }
            clients.append((i, client_data))

        def update_state(client_id, data):
            try:
                response = requests.post(
                    HTTP_URL,
                    json=data,
                    headers={"Content-Type": "application/json"}
                )
                logger.info(f"Client {client_id} sent state update: {data}")
                logger.info(f"Client {client_id} received update response: {response.text}")
                return True
            except Exception as e:
                logger.error(f"Client {client_id} failed to update state: {e}")
                return False

        threads = []
        for client_id, data in clients:
            thread = threading.Thread(target=update_state, args=(client_id, data))
            threads.append(thread)

        for thread in threads:
            thread.start()

        for thread in threads:
            thread.join()

        response = requests.get(HTTP_URL)
        logger.info(f"Received final state after multi-client test: {response.text}")

        data = response.json()
        if 'client_id' in data:
            logger.info("Multi-client test passed: State was updated correctly")
            return True
        else:
            logger.error("Multi-client test failed: State was not updated correctly")
            return False

    except Exception as e:
        logger.error(f"Multi-client test failed with error: {e}")
        return False


async def run_tests():
    """Run WebSocket, HTTP, and gRPC tests."""
    logger.info("Starting shared state integration tests")

    with ThreadPoolExecutor() as executor:
        http_future = executor.submit(http_client)
        grpc_future = executor.submit(grpc_client)
        
        ws_result = await websocket_client()
        
        http_result = http_future.result()
        grpc_result = grpc_future.result()

    multi_client_result = multi_client_test()

    if ws_result and http_result and grpc_result and multi_client_result:
        logger.info("All integration tests passed!")
        return 0
    else:
        logger.error("Some integration tests failed")
        return 1


if __name__ == "__main__":
    try:
        exit_code = asyncio.run(run_tests())
        sys.exit(exit_code)
    except KeyboardInterrupt:
        logger.info("Tests interrupted by user")
        sys.exit(130)
    except Exception as e:
        logger.error(f"Unhandled exception: {e}")
        sys.exit(1)

"""
Test script for the shared application state implementation.

This script tests the shared state implementation between Django and Go,
verifying that state updates are properly propagated between components.
"""

import asyncio
import json
import logging
import sys
import os
import time
import websockets
import requests
from concurrent.futures import ThreadPoolExecutor

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

API_HOST = os.environ.get('API_HOST', 'localhost:8000')
WS_URL = f"ws://{API_HOST}/ws/state/shared/test/"
HTTP_URL = f"http://{API_HOST}/api/state/shared/test/"

TEST_DATA = {
    "message": "Hello from test script",
    "timestamp": time.time(),
    "count": 1
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


async def run_tests():
    """Run WebSocket and HTTP tests concurrently."""
    logger.info("Starting shared state tests")

    with ThreadPoolExecutor() as executor:
        http_future = executor.submit(http_client)
        
        ws_result = await websocket_client()
        
        http_result = http_future.result()

    if ws_result and http_result:
        logger.info("All tests passed!")
        return 0
    else:
        logger.error("Some tests failed")
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

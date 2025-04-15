"""
gRPC server for the agent_runtime Django backend.
"""

import time
import signal
import grpc
import logging
from concurrent import futures
from django.conf import settings

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
)
logger = logging.getLogger(__name__)


class AgentRuntimeGrpcServer:
    """
    gRPC server for the agent_runtime Django backend.

    This server provides gRPC endpoints for the agent_runtime system,
    allowing the Go codebase to interact with the Django backend.
    """

    def __init__(self, host='0.0.0.0', port=50051, max_workers=10):
        self.host = host
        self.port = port
        self.server = grpc.server(
            futures.ThreadPoolExecutor(
                max_workers=max_workers))
        self.is_running = False

        signal.signal(signal.SIGTERM, self._handle_shutdown)
        signal.signal(signal.SIGINT, self._handle_shutdown)

    def _handle_shutdown(self, signum, frame):
        """Handle shutdown signals."""
        logger.info(f"Received signal {signum}, shutting down gRPC server...")
        self.stop()

    def start(self):
        """Start the gRPC server."""

        server_address = f"{self.host}:{self.port}"
        self.server.add_insecure_port(server_address)
        self.server.start()
        self.is_running = True
        logger.info(f"gRPC server started on {server_address}")

        return self

    def stop(self, grace=None):
        """Stop the gRPC server."""
        if self.is_running:
            self.server.stop(grace)
            self.is_running = False
            logger.info("gRPC server stopped")

    def wait_for_termination(self):
        """Wait for the server to terminate."""
        try:
            while self.is_running:
                time.sleep(1)
        except KeyboardInterrupt:
            self.stop()


def run_grpc_server():
    """Run the gRPC server."""
    host = getattr(settings, 'GRPC_SERVER_HOST', '0.0.0.0')
    port = getattr(settings, 'GRPC_SERVER_PORT', 50051)

    server = AgentRuntimeGrpcServer(host=host, port=port)
    server.start()

    try:
        server.wait_for_termination()
    except KeyboardInterrupt:
        server.stop()


if __name__ == '__main__':
    run_grpc_server()

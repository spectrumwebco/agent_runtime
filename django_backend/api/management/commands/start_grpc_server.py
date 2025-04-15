"""
Django management command to start the Go gRPC server.
"""

import os
import sys
import logging
import subprocess
import time
import signal
from django.core.management.base import BaseCommand
from django.conf import settings

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
)
logger = logging.getLogger(__name__)


class Command(BaseCommand):
    """
    Django management command to start the Go gRPC server.
    """

    help = 'Start the Go gRPC server for agent_runtime'

    def add_arguments(self, parser):
        """Add command arguments."""
        parser.add_argument(
            '--port',
            type=int,
            default=50051,
            help='Port for the gRPC server to listen on'
        )
        parser.add_argument(
            '--host',
            type=str,
            default='0.0.0.0',
            help='Host for the gRPC server to bind to'
        )

    def handle(self, *args, **options):
        """Handle the command."""
        port = options['port']
        host = options['host']
        
        self.stdout.write(self.style.SUCCESS(f'Starting Go gRPC server on {host}:{port}...'))
        
        base_dir = settings.BASE_DIR
        repo_root = os.path.dirname(os.path.dirname(os.path.dirname(base_dir)))
        server_script_path = os.path.join(repo_root, 'scripts', 'simple_grpc_server.go')
        
        if not os.path.exists(server_script_path):
            self.stderr.write(self.style.ERROR(f'Server script not found at {server_script_path}'))
            return
        
        try:
            os.chdir(os.path.dirname(server_script_path))
            
            self.stdout.write('Building Go gRPC server...')
            build_process = subprocess.run(
                ['go', 'build', '-o', 'simple_grpc_server', 'simple_grpc_server.go'],
                check=True,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                text=True
            )
            
            self.stdout.write('Starting Go gRPC server process...')
            server_process = subprocess.Popen(
                ['./simple_grpc_server'],
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                text=True
            )
            
            def signal_handler(sig, frame):
                self.stdout.write(self.style.WARNING('Received shutdown signal, stopping server...'))
                server_process.terminate()
                server_process.wait(timeout=5)
                sys.exit(0)
            
            signal.signal(signal.SIGINT, signal_handler)
            signal.signal(signal.SIGTERM, signal_handler)
            
            time.sleep(1)
            
            if server_process.poll() is not None:
                stdout, stderr = server_process.communicate()
                self.stderr.write(self.style.ERROR(f'Server process exited with code {server_process.returncode}'))
                self.stderr.write(f'STDOUT: {stdout}')
                self.stderr.write(f'STDERR: {stderr}')
                return
            
            self.stdout.write(self.style.SUCCESS(f'Go gRPC server running on {host}:{port}'))
            
            try:
                while True:
                    if server_process.poll() is not None:
                        stdout, stderr = server_process.communicate()
                        self.stderr.write(self.style.ERROR(f'Server process exited unexpectedly with code {server_process.returncode}'))
                        self.stderr.write(f'STDOUT: {stdout}')
                        self.stderr.write(f'STDERR: {stderr}')
                        break
                    
                    line = server_process.stdout.readline()
                    if line:
                        self.stdout.write(line.strip())
                    
                    time.sleep(0.1)
            except KeyboardInterrupt:
                self.stdout.write(self.style.WARNING('Keyboard interrupt received, stopping server...'))
                server_process.terminate()
                server_process.wait(timeout=5)
        
        except subprocess.CalledProcessError as e:
            self.stderr.write(self.style.ERROR(f'Failed to build Go gRPC server: {e}'))
            self.stderr.write(f'STDOUT: {e.stdout}')
            self.stderr.write(f'STDERR: {e.stderr}')
        except Exception as e:
            self.stderr.write(self.style.ERROR(f'Error starting Go gRPC server: {e}'))

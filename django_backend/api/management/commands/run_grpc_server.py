"""
Django management command to run the Go gRPC server.

This command starts the Go gRPC server that is integrated with the agent_runtime.
"""

import logging
import subprocess
import os
import signal
import sys
from django.core.management.base import BaseCommand
from django.conf import settings

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
)
logger = logging.getLogger(__name__)


class Command(BaseCommand):
    """
    Django management command to run the Go gRPC server.
    """

    help = 'Run the Go gRPC server for agent_runtime'
    proc = None

    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.proc = None

    def handle(self, *args, **options):
        """Handle the command."""
        self.stdout.write(self.style.SUCCESS('Starting Go gRPC server...'))

        agent_runtime_path = getattr(
            settings, 'AGENT_RUNTIME_PATH', 
            os.path.join(settings.BASE_DIR, '..', 'bin', 'agent_runtime'))
        
        host = getattr(settings, 'GRPC_SERVER_HOST', '0.0.0.0')
        port = getattr(settings, 'GRPC_SERVER_PORT', 50051)
        
        try:
            cmd = [agent_runtime_path, 'serve', '--grpc-only', 
                   f'--grpc-host={host}', f'--grpc-port={port}']
            
            self.stdout.write(f"Running command: {' '.join(cmd)}")
            
            self.proc = subprocess.Popen(
                cmd,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                universal_newlines=True,
            )
            
            signal.signal(signal.SIGINT, self.handle_shutdown)
            signal.signal(signal.SIGTERM, self.handle_shutdown)
            
            self.stdout.write(self.style.SUCCESS(
                f'Go gRPC server started on {host}:{port}'))
            
            while self.proc and self.proc.poll() is None:
                if self.proc.stdout:
                    line = self.proc.stdout.readline()
                    if line:
                        self.stdout.write(line.strip())
            
            if self.proc:
                return_code = self.proc.poll()
                if return_code != 0 and self.proc.stderr:
                    stderr = self.proc.stderr.read()
                    self.stdout.write(self.style.ERROR(
                        f'Go gRPC server exited with code {return_code}: {stderr}'))
            else:
                self.stdout.write(self.style.SUCCESS('Go gRPC server stopped'))
                
        except FileNotFoundError:
            self.stdout.write(self.style.ERROR(
                f'Agent runtime binary not found at {agent_runtime_path}'))
            self.stdout.write(self.style.WARNING(
                'Falling back to Python gRPC server...'))
            
            from django_backend.grpc_server import run_grpc_server
            run_grpc_server()
            
        except Exception as e:
            self.stdout.write(self.style.ERROR(f'Error starting Go gRPC server: {e}'))
            if self.proc:
                self.proc.terminate()
    
    def handle_shutdown(self, signum, frame):
        """Handle shutdown signals."""
        self.stdout.write(self.style.WARNING(
            f'Received signal {signum}, shutting down Go gRPC server...'))
        if self.proc:
            self.proc.terminate()
            try:
                self.proc.wait(timeout=5)
            except subprocess.TimeoutExpired:
                self.proc.kill()
            self.stdout.write(self.style.SUCCESS('Go gRPC server stopped'))
        sys.exit(0)

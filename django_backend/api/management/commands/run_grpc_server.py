"""
Django management command to run the gRPC server.
"""

import os
import sys
import time
import signal
import logging
from django.core.management.base import BaseCommand
from django.conf import settings

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
)
logger = logging.getLogger(__name__)


class Command(BaseCommand):
    """
    Django management command to run the gRPC server.
    """
    
    help = 'Run the gRPC server for agent_runtime'
    
    def handle(self, *args, **options):
        """Handle the command."""
        self.stdout.write(self.style.SUCCESS('Starting gRPC server...'))
        
        from django_backend.grpc_server import run_grpc_server
        
        try:
            run_grpc_server()
        except KeyboardInterrupt:
            self.stdout.write(self.style.SUCCESS('gRPC server stopped'))

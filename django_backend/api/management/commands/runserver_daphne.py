"""
Django management command to run the server with Daphne.
"""

import os
import sys
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
    Django management command to run the server with Daphne.
    
    This command runs the Django server using Daphne, which supports
    both HTTP and WebSocket protocols.
    """
    
    help = 'Run the server with Daphne (HTTP + WebSocket)'
    
    def add_arguments(self, parser):
        """Add command arguments."""
        parser.add_argument(
            '--host',
            default='0.0.0.0',
            help='Host to bind to'
        )
        parser.add_argument(
            '--port',
            default='8000',
            help='Port to bind to'
        )
    
    def handle(self, *args, **options):
        """Handle the command."""
        host = options['host']
        port = options['port']
        
        self.stdout.write(self.style.SUCCESS(f'Starting Daphne server on {host}:{port}...'))
        
        from daphne.cli import CommandLineInterface
        
        sys.argv = [
            'daphne',
            '-b', host,
            '-p', port,
            'agent_api.asgi:application'
        ]
        
        CommandLineInterface().run()

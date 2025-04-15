"""
WSGI config for agent_runtime Django backend.
"""

import os
from django.core.wsgi import get_wsgi_application

os.environ.setdefault('DJANGO_SETTINGS_MODULE', 'agent_api.settings')

application = get_wsgi_application()

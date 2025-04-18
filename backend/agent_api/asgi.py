"""
ASGI config for agent_api project.

It exposes the ASGI callable as a module-level variable named ``application``.

For more information on this file, see
https://docs.djangoproject.com/en/5.2/howto/deployment/asgi/
"""

from django.urls import path
import os
import django
from django.core.asgi import get_asgi_application
from channels.routing import ProtocolTypeRouter, URLRouter
from channels.auth import AuthMiddlewareStack

os.environ.setdefault('DJANGO_SETTINGS_MODULE', 'agent_api.settings')
django.setup()

from api.websocket import AgentConsumer  # noqa
from api.websocket_state import SharedStateConsumer  # noqa
from api.socketio_consumer import OpenHandsSocketIOConsumer  # noqa

# Get the Django ASGI application
django_asgi_app = get_asgi_application()

websocket_urlpatterns = [
    path('ws/agent/<str:client_id>/', AgentConsumer.as_asgi()),
    path('ws/agent/<str:client_id>/<str:task_id>/', AgentConsumer.as_asgi()),
    path('ws/state/<str:state_type>/<str:state_id>/', SharedStateConsumer.as_asgi()),
    path('ws/state/', SharedStateConsumer.as_asgi()),
    path('socket.io/', OpenHandsSocketIOConsumer.as_asgi()),
]

# Configure the ASGI application
application = ProtocolTypeRouter({
    'http': django_asgi_app,
    'websocket': AuthMiddlewareStack(
        URLRouter(
            websocket_urlpatterns
        )
    ),
})

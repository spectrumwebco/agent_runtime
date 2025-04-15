"""
WebSocket routing configuration for the agent_api project.
"""

from django.urls import path
from channels.routing import ProtocolTypeRouter, URLRouter
from channels.auth import AuthMiddlewareStack
from api.websocket import AgentConsumer

websocket_urlpatterns = [
    path('ws/agent/<str:client_id>/', AgentConsumer.as_asgi()),
    path('ws/agent/<str:client_id>/<str:task_id>/', AgentConsumer.as_asgi()),
]

application = ProtocolTypeRouter({
    'websocket': AuthMiddlewareStack(
        URLRouter(
            websocket_urlpatterns
        )
    ),
})

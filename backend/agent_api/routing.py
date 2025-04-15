"""
ASGI routing configuration for the agent_api project.
"""

from django.urls import path
from django.core.asgi import get_asgi_application
from channels.routing import ProtocolTypeRouter, URLRouter
from channels.auth import AuthMiddlewareStack
from channels.security.websocket import AllowedHostsOriginValidator
from api.websocket import AgentConsumer
from api.websocket_state import SharedStateConsumer
from api.socketio_consumer import OpenHandsSocketIOConsumer
from apps.python_agent.agent.django_views.consumers import AgentConsumer as PythonAgentConsumer

django_asgi_app = get_asgi_application()

websocket_urlpatterns = [
    path('ws/agent/<str:client_id>/', AgentConsumer.as_asgi()),
    path('ws/agent/<str:client_id>/<str:task_id>/', AgentConsumer.as_asgi()),
    path('ws/state/<str:state_type>/<str:state_id>/', SharedStateConsumer.as_asgi()),
    path('ws/state/', SharedStateConsumer.as_asgi()),
    path('socket.io/', OpenHandsSocketIOConsumer.as_asgi()),
    path('ws/python_agent/<str:thread_id>/', PythonAgentConsumer.as_asgi()),
]

application = ProtocolTypeRouter({
    'http': django_asgi_app,
    'websocket': AllowedHostsOriginValidator(
        AuthMiddlewareStack(
            URLRouter(
                websocket_urlpatterns
            )
        )
    ),
})

"""
URL configuration for the API app.
"""

from django.urls import path, include
from rest_framework.routers import DefaultRouter
from . import views
from .views.state_views import SharedStateViewSet
from .views.conversation_views import ConversationViewSet
from .views.auth_views import (
    login_view, register_view, logout_view, 
    github_auth, github_callback, 
    gitee_auth, gitee_callback,
    user_settings
)
from .views.events_views import send_event, forward_to_agent, get_events, create_event
from .views.options_views import get_models, get_agents, get_security_analyzers, get_config
from .views.billing_views import get_credits, add_credits, get_transactions, get_subscription

router = DefaultRouter()
router.register(r'users', views.UserViewSet)
router.register(r'state/shared', SharedStateViewSet, basename='shared-state')
router.register(r'conversations', ConversationViewSet, basename='conversation')

urlpatterns = [
    path('', views.api_root, name='api-root'),
    path('', include(router.urls)),
    path('tasks/', views.execute_agent_task, name='execute-agent-task'),
    
    path('auth/login/', login_view, name='login'),
    path('auth/register/', register_view, name='register'),
    path('auth/logout/', logout_view, name='logout'),
    path('auth/github/', github_auth, name='github-auth'),
    path('auth/github/callback/', github_callback, name='github-callback'),
    path('auth/gitee/', gitee_auth, name='gitee-auth'),
    path('auth/gitee/callback/', gitee_callback, name='gitee-callback'),
    path('auth/settings/', user_settings, name='user-settings'),
    
    path('events/send/', send_event, name='send-event'),
    path('events/forward/', forward_to_agent, name='forward-to-agent'),
    path('events/<str:conversation_id>/', get_events, name='get-events'),
    path('events/<str:conversation_id>/create/', create_event, name='create-event'),
    
    path('options/models/', get_models, name='get-models'),
    path('options/agents/', get_agents, name='get-agents'),
    path('options/security-analyzers/', get_security_analyzers, name='get-security-analyzers'),
    path('options/config/', get_config, name='get-config'),
    
    path('billing/credits/', get_credits, name='get-credits'),
    path('billing/credits/add/', add_credits, name='add-credits'),
    path('billing/transactions/', get_transactions, name='get-transactions'),
    path('billing/subscription/', get_subscription, name='get-subscription'),
    
    path('health/', views.health_check, name='health-check'),
    path('health/readiness/', views.readiness_check, name='readiness-check'),
    path('health/liveness/', views.liveness_check, name='liveness-check'),
]

"""
URL configuration for the API app.
"""

from django.urls import path, include
from rest_framework.routers import DefaultRouter
from . import views
from .views.state_views import SharedStateViewSet

router = DefaultRouter()
router.register(r'users', views.UserViewSet)
router.register(r'state/shared', SharedStateViewSet, basename='shared-state')

urlpatterns = [
    path('', views.api_root, name='api-root'),
    path('', include(router.urls)),
    path('tasks/', views.execute_agent_task, name='execute-agent-task'),
]

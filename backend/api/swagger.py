"""
Swagger API documentation for the agent_runtime Django backend.
"""

from django.urls import path
from django.views.generic import TemplateView
from rest_framework.schemas import get_schema_view
from rest_framework import permissions

schema_view = get_schema_view(
    title="Agent Runtime API",
    description="API for the Agent Runtime system",
    version="1.0.0",
    public=True,
    permission_classes=[permissions.AllowAny],
)

urlpatterns = [
    path('openapi/', schema_view, name='openapi-schema'),

    path('swagger-ui/', TemplateView.as_view(
        template_name='swagger-ui.html',
        extra_context={'schema_url': 'openapi-schema'}
    ), name='swagger-ui'),

    path('redoc/', TemplateView.as_view(
        template_name='redoc.html',
        extra_context={'schema_url': 'openapi-schema'}
    ), name='redoc'),
]

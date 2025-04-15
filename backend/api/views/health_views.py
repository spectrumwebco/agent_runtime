"""
Health check views for the API.

This module provides health check endpoints for the API.
"""

import logging
from django.http import JsonResponse
from django.views.decorators.http import require_http_methods
from django.views.decorators.csrf import csrf_exempt
from rest_framework.decorators import api_view, permission_classes
from rest_framework.permissions import AllowAny

logger = logging.getLogger(__name__)


@api_view(['GET'])
@permission_classes([AllowAny])
def health_check(request):
    """
    Health check endpoint for the API.
    
    This endpoint is used by Kubernetes liveness and readiness probes.
    """
    return JsonResponse({
        'status': 'ok',
        'message': 'API is healthy'
    })


@api_view(['GET'])
@permission_classes([AllowAny])
def readiness_check(request):
    """
    Readiness check endpoint for the API.
    
    This endpoint is used by Kubernetes readiness probes.
    """
    from django.db import connections
    from django.db.utils import OperationalError
    
    try:
        db_conn = connections['default']
        db_conn.cursor()
    except OperationalError:
        logger.error("Database is not available")
        return JsonResponse({
            'status': 'error',
            'message': 'Database is not available'
        }, status=503)
    
    from django_redis import get_redis_connection
    
    try:
        redis_conn = get_redis_connection("default")
        redis_conn.ping()
    except Exception as e:
        logger.error(f"Redis is not available: {e}")
        return JsonResponse({
            'status': 'error',
            'message': 'Redis is not available'
        }, status=503)
    
    return JsonResponse({
        'status': 'ok',
        'message': 'API is ready'
    })


@api_view(['GET'])
@permission_classes([AllowAny])
def liveness_check(request):
    """
    Liveness check endpoint for the API.
    
    This endpoint is used by Kubernetes liveness probes.
    """
    return JsonResponse({
        'status': 'ok',
        'message': 'API is alive'
    })

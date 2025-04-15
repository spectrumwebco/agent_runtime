"""
Views for the API app.
"""

from rest_framework import viewsets, permissions, status
from rest_framework.decorators import api_view, permission_classes
from rest_framework.response import Response
from django.contrib.auth.models import User
from .serializers import UserSerializer

import sys
from django.conf import settings

sys.path.append(str(settings.SRC_DIR))


class UserViewSet(viewsets.ReadOnlyModelViewSet):
    """
    API endpoint that allows users to be viewed.
    """
    queryset = User.objects.all().order_by('-date_joined')
    serializer_class = UserSerializer
    permission_classes = [permissions.IsAuthenticated]


@api_view(['GET'])
@permission_classes([permissions.AllowAny])
def api_root(request):
    """
    API root endpoint.
    """
    return Response({
        'status': 'online',
        'version': '1.0.0',
        'message': 'Agent Runtime API is running'
    })


@api_view(['POST'])
@permission_classes([permissions.IsAuthenticated])
def execute_agent_task(request):
    """
    Execute a task using the agent runtime.
    """
    try:
        return Response({
            'status': 'accepted',
            'task_id': 'placeholder-task-id',
            'message': 'Task submitted for execution'
        }, status=status.HTTP_202_ACCEPTED)
    except Exception as e:
        return Response({
            'status': 'error',
            'message': str(e)
        }, status=status.HTTP_500_INTERNAL_SERVER_ERROR)

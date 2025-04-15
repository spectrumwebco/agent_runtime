"""
Views for conversation management.

This module provides API views for managing conversations, including
creating, retrieving, updating, and deleting conversations, as well as
managing conversation events and trajectories.
"""

import os
import json
import logging
import zipfile
import tempfile
from typing import Dict, Any, List, Optional
from django.http import JsonResponse, HttpResponse, FileResponse
from django.views.decorators.http import require_http_methods
from django.views.decorators.csrf import csrf_exempt
from django.contrib.auth.decorators import login_required
from django.conf import settings
from rest_framework import viewsets, permissions
from rest_framework.decorators import action
from rest_framework.response import Response
from rest_framework import status

from api.models import Conversation, ConversationEvent
from api.serializers import ConversationSerializer
from api.utils.workspace import create_workspace, get_workspace_path, list_workspace_files, read_file_content

logger = logging.getLogger(__name__)


class ConversationViewSet(viewsets.ModelViewSet):
    """ViewSet for managing conversations."""
    
    serializer_class = ConversationSerializer
    permission_classes = [permissions.IsAuthenticated]
    
    def get_queryset(self):
        """Get the queryset for the current user."""
        return Conversation.objects.filter(user=self.request.user)
    
    def perform_create(self, serializer):
        """Create a new conversation."""
        workspace_path = create_workspace()
        
        serializer.save(
            user=self.request.user,
            workspace_path=workspace_path
        )
    
    @action(detail=True, methods=['get'])
    def events(self, request, pk=None):
        """Get events for a conversation."""
        conversation = self.get_object()
        
        events = ConversationEvent.objects.filter(
            conversation_id=str(conversation.id)
        ).order_by('event_id')
        
        event_dicts = [event.to_dict() for event in events]
        
        return Response(event_dicts)
    
    @action(detail=True, methods=['post'])
    def submit_feedback(self, request, pk=None):
        """Submit feedback for a conversation."""
        conversation = self.get_object()
        
        feedback = request.data.get('feedback')
        rating = request.data.get('rating')
        
        
        return Response({'status': 'success'})
    
    @action(detail=True, methods=['get'])
    def zip_directory(self, request, pk=None):
        """Get a zip file of the conversation workspace."""
        conversation = self.get_object()
        
        if not conversation.workspace_path:
            return Response(
                {'error': 'No workspace found for this conversation'},
                status=status.HTTP_404_NOT_FOUND
            )
        
        with tempfile.NamedTemporaryFile(delete=False, suffix='.zip') as temp_file:
            temp_path = temp_file.name
        
        workspace_path = get_workspace_path(conversation.workspace_path)
        with zipfile.ZipFile(temp_path, 'w', zipfile.ZIP_DEFLATED) as zipf:
            for root, _, files in os.walk(workspace_path):
                for file in files:
                    file_path = os.path.join(root, file)
                    arcname = os.path.relpath(file_path, workspace_path)
                    zipf.write(file_path, arcname)
        
        return FileResponse(
            open(temp_path, 'rb'),
            as_attachment=True,
            filename=f'conversation_{conversation.id}.zip'
        )
    
    @action(detail=True, methods=['get'])
    def vscode_url(self, request, pk=None):
        """Get a VSCode URL for the conversation workspace."""
        conversation = self.get_object()
        
        if not conversation.workspace_path:
            return Response(
                {'error': 'No workspace found for this conversation'},
                status=status.HTTP_404_NOT_FOUND
            )
        
        workspace_path = get_workspace_path(conversation.workspace_path)
        
        vscode_url = f"vscode://file{workspace_path}"
        
        return Response({'vscode_url': vscode_url})
    
    @action(detail=True, methods=['get'])
    def config(self, request, pk=None):
        """Get the runtime configuration for a conversation."""
        conversation = self.get_object()
        
        return Response({'runtime_id': str(conversation.id)})
    
    @action(detail=True, methods=['get'])
    def trajectory(self, request, pk=None):
        """Get the trajectory data for a conversation."""
        conversation = self.get_object()
        
        events = ConversationEvent.objects.filter(
            conversation_id=str(conversation.id)
        ).order_by('event_id')
        
        trajectory = {
            'id': str(conversation.id),
            'events': [event.to_dict() for event in events],
            'metadata': {
                'model': conversation.model,
                'agent': conversation.agent,
                'created_at': conversation.created_at.isoformat(),
                'updated_at': conversation.updated_at.isoformat(),
            }
        }
        
        return Response(trajectory)
    
    @action(detail=True, methods=['get'])
    def list_files(self, request, pk=None):
        """List files in the conversation workspace."""
        conversation = self.get_object()
        
        if not conversation.workspace_path:
            return Response(
                {'error': 'No workspace found for this conversation'},
                status=status.HTTP_404_NOT_FOUND
            )
        
        workspace_path = get_workspace_path(conversation.workspace_path)
        
        files = list_workspace_files(workspace_path)
        
        return Response({'files': files})
    
    @action(detail=True, methods=['post'])
    def select_file(self, request, pk=None):
        """Get the content of a file in the conversation workspace."""
        conversation = self.get_object()
        
        if not conversation.workspace_path:
            return Response(
                {'error': 'No workspace found for this conversation'},
                status=status.HTTP_404_NOT_FOUND
            )
        
        file_path = request.data.get('path')
        if not file_path:
            return Response(
                {'error': 'No file path provided'},
                status=status.HTTP_400_BAD_REQUEST
            )
        
        workspace_path = get_workspace_path(conversation.workspace_path)
        
        full_path = os.path.join(workspace_path, file_path)
        if not os.path.exists(full_path) or not os.path.isfile(full_path):
            return Response(
                {'error': f'File not found: {file_path}'},
                status=status.HTTP_404_NOT_FOUND
            )
        
        content = read_file_content(full_path)
        
        return Response({
            'path': file_path,
            'content': content,
            'size': os.path.getsize(full_path),
            'last_modified': os.path.getmtime(full_path),
        })

"""
Views for event handling.

This module provides API views for handling events, including
WebSocket events, agent events, and user events.
"""

import json
import logging
import asyncio
from datetime import datetime
from typing import Dict, Any, Optional
from django.http import JsonResponse, HttpResponse
from django.views.decorators.http import require_http_methods
from django.views.decorators.csrf import csrf_exempt
from django.contrib.auth.decorators import login_required
from django.conf import settings
from rest_framework.decorators import api_view, permission_classes
from rest_framework.permissions import IsAuthenticated
from rest_framework.response import Response

from api.models import Conversation, ConversationEvent
from api.socketio_consumer import (
    send_agent_message, send_agent_observation,
    send_agent_error, send_agent_state_update
)

logger = logging.getLogger(__name__)


@api_view(['POST'])
@permission_classes([IsAuthenticated])
def send_event(request):
    """Send an event to a conversation."""
    try:
        conversation_id = request.data.get('conversation_id')
        event_type = request.data.get('event_type')
        data = request.data.get('data', {})
        
        if not conversation_id:
            return Response({
                'status': 'error',
                'message': 'Missing conversation_id'
            }, status=400)
        
        if not event_type:
            return Response({
                'status': 'error',
                'message': 'Missing event_type'
            }, status=400)
        
        if event_type == 'agent_message':
            message = data.get('message')
            extras = data.get('extras')
            
            if not message:
                return Response({
                    'status': 'error',
                    'message': 'Missing message'
                }, status=400)
            
            asyncio.create_task(send_agent_message(conversation_id, message, extras))
            
        elif event_type == 'agent_observation':
            observation = data.get('observation')
            observation_type = data.get('observation_type', 'info')
            extras = data.get('extras')
            
            if not observation:
                return Response({
                    'status': 'error',
                    'message': 'Missing observation'
                }, status=400)
            
            asyncio.create_task(send_agent_observation(
                conversation_id, observation, observation_type, extras
            ))
            
        elif event_type == 'agent_error':
            message = data.get('message')
            error_id = data.get('error_id')
            
            if not message:
                return Response({
                    'status': 'error',
                    'message': 'Missing message'
                }, status=400)
            
            asyncio.create_task(send_agent_error(conversation_id, message, error_id))
            
        elif event_type == 'agent_state_update':
            state = data.get('state')
            
            if not state:
                return Response({
                    'status': 'error',
                    'message': 'Missing state'
                }, status=400)
            
            asyncio.create_task(send_agent_state_update(conversation_id, state))
            
        else:
            return Response({
                'status': 'error',
                'message': f'Unknown event type: {event_type}'
            }, status=400)
        
        return Response({
            'status': 'success',
            'message': f'Event {event_type} sent to conversation {conversation_id}'
        })
    
    except Exception as e:
        logger.error(f"Error sending event: {e}")
        return Response({
            'status': 'error',
            'message': str(e)
        }, status=500)


@api_view(['POST'])
@permission_classes([IsAuthenticated])
def forward_to_agent(request):
    """Forward a message to the agent."""
    try:
        conversation_id = request.data.get('conversation_id')
        message = request.data.get('message')
        
        if not conversation_id:
            return Response({
                'status': 'error',
                'message': 'Missing conversation_id'
            }, status=400)
        
        if not message:
            return Response({
                'status': 'error',
                'message': 'Missing message'
            }, status=400)
        
        event = ConversationEvent.objects.create(
            conversation_id=conversation_id,
            event_id=ConversationEvent.objects.filter(conversation_id=conversation_id).count() + 1,
            event_type='user_message',
            source='user',
            content=json.dumps({
                'message': message,
                'timestamp': int(datetime.now().timestamp() * 1000)
            })
        )
        
        
        return Response({
            'status': 'success',
            'message': 'Message forwarded to agent',
            'event_id': str(event.id)
        })
    
    except Exception as e:
        logger.error(f"Error forwarding message to agent: {e}")
        return Response({
            'status': 'error',
            'message': str(e)
        }, status=500)


@api_view(['GET'])
@permission_classes([IsAuthenticated])
def get_events(request, conversation_id):
    """Get events for a conversation."""
    try:
        events = ConversationEvent.objects.filter(
            conversation_id=conversation_id
        ).order_by('event_id')
        
        event_dicts = [event.to_dict() for event in events]
        
        return Response(event_dicts)
    
    except Exception as e:
        logger.error(f"Error getting events: {e}")
        return Response({
            'status': 'error',
            'message': str(e)
        }, status=500)


@api_view(['POST'])
@permission_classes([IsAuthenticated])
def create_event(request, conversation_id):
    """Create an event for a conversation."""
    try:
        event_type = request.data.get('event_type')
        source = request.data.get('source')
        content = request.data.get('content')
        
        if not event_type:
            return Response({
                'status': 'error',
                'message': 'Missing event_type'
            }, status=400)
        
        if not source:
            return Response({
                'status': 'error',
                'message': 'Missing source'
            }, status=400)
        
        if not content:
            return Response({
                'status': 'error',
                'message': 'Missing content'
            }, status=400)
        
        event = ConversationEvent.objects.create(
            conversation_id=conversation_id,
            event_id=ConversationEvent.objects.filter(conversation_id=conversation_id).count() + 1,
            event_type=event_type,
            source=source,
            content=json.dumps(content)
        )
        
        return Response({
            'status': 'success',
            'message': 'Event created',
            'event_id': str(event.id)
        })
    
    except Exception as e:
        logger.error(f"Error creating event: {e}")
        return Response({
            'status': 'error',
            'message': str(e)
        }, status=500)

"""
REST API views for shared application state.

This module provides REST API endpoints for interacting with the shared
application state system, enabling HTTP-based access to state data.
"""

from django.http import JsonResponse
from rest_framework.decorators import api_view, permission_classes
from rest_framework.permissions import IsAuthenticated
from rest_framework.response import Response
from rest_framework import status
import json
import logging

from api.grpc_bridge import grpc_bridge

logger = logging.getLogger(__name__)


@api_view(['GET', 'POST'])
def shared_state_view(request, state_id='default'):
    """
    REST API endpoint for shared state.
    
    GET: Retrieve the current state
    POST: Update the state
    
    Args:
        request: The HTTP request
        state_id: The ID of the state to retrieve or update
        
    Returns:
        Response: The HTTP response containing the state data or status
    """
    state_type = 'shared'
    
    if request.method == 'GET':
        try:
            response = grpc_bridge.get_state(
                state_type=state_type,
                state_id=state_id
            )
            
            if response and 'data' in response:
                return Response(response['data'])
            
            return Response(
                {"error": "Failed to retrieve state"},
                status=status.HTTP_404_NOT_FOUND
            )
        
        except Exception as e:
            logger.error(f"Error retrieving state: {e}")
            return Response(
                {"error": f"Failed to retrieve state: {str(e)}"},
                status=status.HTTP_500_INTERNAL_SERVER_ERROR
            )
    
    elif request.method == 'POST':
        try:
            data = request.data
            
            response = grpc_bridge.update_state(
                state_type=state_type,
                state_id=state_id,
                data=data
            )
            
            if response.get('status') == 'success':
                return Response(
                    {"status": "success", "message": "State updated successfully"},
                    status=status.HTTP_200_OK
                )
            
            return Response(
                {"error": "Failed to update state"},
                status=status.HTTP_500_INTERNAL_SERVER_ERROR
            )
        
        except Exception as e:
            logger.error(f"Error updating state: {e}")
            return Response(
                {"error": f"Failed to update state: {str(e)}"},
                status=status.HTTP_500_INTERNAL_SERVER_ERROR
            )

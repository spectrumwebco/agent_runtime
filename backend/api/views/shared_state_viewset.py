"""
ViewSet for shared application state.

This module provides a ViewSet for interacting with the shared
application state system, enabling REST API access to state data.
"""

from rest_framework import viewsets, status
from rest_framework.response import Response
from rest_framework.decorators import action
import logging

from api.grpc_bridge import grpc_bridge

logger = logging.getLogger(__name__)


class SharedStateViewSet(viewsets.ViewSet):
    """
    ViewSet for shared state operations.
    
    Provides endpoints for retrieving and updating shared state data.
    """
    
    def retrieve(self, request, pk=None):
        """
        Retrieve the current state.
        
        Args:
            request: The HTTP request
            pk: The primary key (state_id) of the state to retrieve
            
        Returns:
            Response: The HTTP response containing the state data
        """
        state_type = 'shared'
        state_id = pk or 'default'
        
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
    
    def update(self, request, pk=None):
        """
        Update the state.
        
        Args:
            request: The HTTP request
            pk: The primary key (state_id) of the state to update
            
        Returns:
            Response: The HTTP response containing the status
        """
        state_type = 'shared'
        state_id = pk or 'default'
        
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
    
    @action(detail=False, methods=['GET'])
    def list_states(self, request):
        """
        List all available states.
        
        Args:
            request: The HTTP request
            
        Returns:
            Response: The HTTP response containing the list of states
        """
        state_type = 'shared'
        
        try:
            response = grpc_bridge.list_states(state_type=state_type)
            
            if response and 'states' in response:
                return Response(response['states'])
            
            return Response(
                {"error": "Failed to list states"},
                status=status.HTTP_500_INTERNAL_SERVER_ERROR
            )
        
        except Exception as e:
            logger.error(f"Error listing states: {e}")
            return Response(
                {"error": f"Failed to list states: {str(e)}"},
                status=status.HTTP_500_INTERNAL_SERVER_ERROR
            )

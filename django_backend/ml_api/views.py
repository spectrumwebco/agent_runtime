"""
Views for the ML API app.
"""

import sys
import os
from django.conf import settings
from rest_framework import viewsets, status
from rest_framework.decorators import action, api_view, permission_classes
from rest_framework.response import Response
from rest_framework.permissions import IsAuthenticated, AllowAny

sys.path.append(str(settings.SRC_DIR))

try:
    from ml_infrastructure.api.client import MLInfrastructureAPIClient
    from models.api.ml_infrastructure_api_models import (
        ModelDetail, ModelList, FineTuningJobDetail, FineTuningJobCreate,
        InferenceServiceDetail, InferenceServiceCreate
    )
    from models.data_validation.validation_models import (
        RawDataModel, ChatFormatModel, CompletionFormatModel,
        ValidationResult, QualityMetrics
    )
    ML_CLIENT_AVAILABLE = True
except ImportError:
    ML_CLIENT_AVAILABLE = False

from .serializers import (
    ModelDetailSerializer, ModelListSerializer,
    FineTuningJobDetailSerializer, FineTuningJobCreateSerializer,
    InferenceServiceDetailSerializer, InferenceServiceCreateSerializer,
    ValidationResultSerializer, QualityMetricsSerializer
)


@api_view(['GET'])
@permission_classes([AllowAny])
def ml_api_root(request):
    """
    ML API root endpoint.
    """
    return Response({
        'status': 'online',
        'version': '1.0.0',
        'message': 'ML Infrastructure API is running',
        'ml_client_available': ML_CLIENT_AVAILABLE
    })


class MLModelViewSet(viewsets.ViewSet):
    """
    ViewSet for ML models.
    """
    permission_classes = [IsAuthenticated]
    
    def __init__(self, **kwargs):
        super().__init__(**kwargs)
        if ML_CLIENT_AVAILABLE:
            self.client = MLInfrastructureAPIClient()
    
    def list(self, request):
        """List all available models."""
        if not ML_CLIENT_AVAILABLE:
            return Response({"error": "ML Infrastructure client not available"}, 
                           status=status.HTTP_503_SERVICE_UNAVAILABLE)
        
        try:
            models = self.client.list_models()
            serializer = ModelListSerializer.from_pydantic(models)
            return Response(serializer.data)
        except Exception as e:
            return Response({"error": str(e)}, status=status.HTTP_500_INTERNAL_SERVER_ERROR)
    
    def retrieve(self, request, pk=None):
        """Retrieve a specific model by ID."""
        if not ML_CLIENT_AVAILABLE:
            return Response({"error": "ML Infrastructure client not available"}, 
                           status=status.HTTP_503_SERVICE_UNAVAILABLE)
        
        try:
            model = self.client.get_model(pk)
            serializer = ModelDetailSerializer.from_pydantic(model)
            return Response(serializer.data)
        except Exception as e:
            return Response({"error": str(e)}, status=status.HTTP_404_NOT_FOUND)


class FineTuningJobViewSet(viewsets.ViewSet):
    """
    ViewSet for fine-tuning jobs.
    """
    permission_classes = [IsAuthenticated]
    
    def __init__(self, **kwargs):
        super().__init__(**kwargs)
        if ML_CLIENT_AVAILABLE:
            self.client = MLInfrastructureAPIClient()
    
    def list(self, request):
        """List all fine-tuning jobs."""
        if not ML_CLIENT_AVAILABLE:
            return Response({"error": "ML Infrastructure client not available"}, 
                           status=status.HTTP_503_SERVICE_UNAVAILABLE)
        
        try:
            jobs = self.client.list_fine_tuning_jobs()
            serialized_jobs = [FineTuningJobDetailSerializer.from_pydantic(job).data for job in jobs]
            return Response(serialized_jobs)
        except Exception as e:
            return Response({"error": str(e)}, status=status.HTTP_500_INTERNAL_SERVER_ERROR)
    
    def retrieve(self, request, pk=None):
        """Retrieve a specific fine-tuning job by ID."""
        if not ML_CLIENT_AVAILABLE:
            return Response({"error": "ML Infrastructure client not available"}, 
                           status=status.HTTP_503_SERVICE_UNAVAILABLE)
        
        try:
            job = self.client.get_fine_tuning_job(pk)
            serializer = FineTuningJobDetailSerializer.from_pydantic(job)
            return Response(serializer.data)
        except Exception as e:
            return Response({"error": str(e)}, status=status.HTTP_404_NOT_FOUND)
    
    def create(self, request):
        """Create a new fine-tuning job."""
        if not ML_CLIENT_AVAILABLE:
            return Response({"error": "ML Infrastructure client not available"}, 
                           status=status.HTTP_503_SERVICE_UNAVAILABLE)
        
        try:
            job_create = FineTuningJobCreate(**request.data)
            
            job = self.client.create_fine_tuning_job(job_create)
            serializer = FineTuningJobDetailSerializer.from_pydantic(job)
            return Response(serializer.data, status=status.HTTP_201_CREATED)
        except Exception as e:
            return Response({"error": str(e)}, status=status.HTTP_400_BAD_REQUEST)
    
    @action(detail=True, methods=['post'])
    def cancel(self, request, pk=None):
        """Cancel a fine-tuning job."""
        if not ML_CLIENT_AVAILABLE:
            return Response({"error": "ML Infrastructure client not available"}, 
                           status=status.HTTP_503_SERVICE_UNAVAILABLE)
        
        try:
            result = self.client.cancel_fine_tuning_job(pk)
            return Response({"status": "success", "message": f"Job {pk} cancelled"})
        except Exception as e:
            return Response({"error": str(e)}, status=status.HTTP_400_BAD_REQUEST)

"""
Django Ninja API configuration for the API app.
"""

import sys
from typing import List, Dict, Any, Optional
from django.conf import settings
from ninja import NinjaAPI, Schema
from ninja.security import APIKeyHeader

sys.path.append(str(settings.SRC_DIR))

try:
    from models.api.ml_infrastructure_api_models import (
        ModelDetail
    )
    PYDANTIC_MODELS_AVAILABLE = True
except ImportError:
    PYDANTIC_MODELS_AVAILABLE = False


class ApiKey(APIKeyHeader):
    param_name = "X-API-Key"

    def authenticate(self, request, key):
        if key == settings.API_KEY:
            return key
        return None


api = NinjaAPI(
    title="Agent Runtime API",
    version="1.0.0",
    description="API for the Agent Runtime system",
    auth=ApiKey(),
    csrf=False,  # Disable CSRF for API testing
    urls_namespace="agent_api"
)


class TaskInput(Schema):
    """Schema for task input."""
    prompt: str
    context: Optional[Dict[str, Any]] = None
    tools: Optional[List[str]] = None


class TaskOutput(Schema):
    """Schema for task output."""
    task_id: str
    status: str
    message: str


@api.get("/", auth=None)
def api_root(request):
    """API root endpoint."""
    return {
        "status": "online",
        "version": "1.0.0",
        "message": "Agent Runtime API is running",
        "pydantic_models_available": PYDANTIC_MODELS_AVAILABLE
    }


@api.post("/tasks", response=TaskOutput)
def execute_task(request, task_input: TaskInput):
    """Execute a task using the agent runtime."""
    return {
        "task_id": "placeholder-task-id",
        "status": "accepted",
        "message": "Task submitted for execution"
    }


if PYDANTIC_MODELS_AVAILABLE:
    @api.get("/models", response=List[ModelDetail])
    def list_models(request):
        """List all available models."""
        return []

    @api.get("/models/{model_id}", response=ModelDetail)
    def get_model(request, model_id: str):
        """Get a specific model by ID."""
        return ModelDetail(
            id=model_id,
            name="Placeholder Model",
            version="1.0.0",
            description="Placeholder model for API testing",
            parameters={},
            created_at="2025-04-15T00:00:00Z",
            updated_at="2025-04-15T00:00:00Z"
        )

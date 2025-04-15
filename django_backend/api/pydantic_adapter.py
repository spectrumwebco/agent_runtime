"""
Adapter for Pydantic models to work with Django.
"""

from typing import Type, TypeVar, Dict, Any, Optional, List, Union
from pydantic import BaseModel

T = TypeVar('T', bound=BaseModel)


class PydanticModelAdapter:
    """
    Adapter for Pydantic models to work with Django.
    
    This adapter provides methods to convert between Pydantic models and
    Django-compatible dictionaries, enabling seamless integration between
    Django and Pydantic.
    """
    
    @staticmethod
    def to_dict(model: BaseModel) -> Dict[str, Any]:
        """Convert a Pydantic model to a dictionary."""
        return model.model_dump()
    
    @staticmethod
    def from_dict(model_class: Type[T], data: Dict[str, Any]) -> T:
        """Convert a dictionary to a Pydantic model."""
        return model_class(**data)
    
    @staticmethod
    def to_list_of_dicts(models: List[BaseModel]) -> List[Dict[str, Any]]:
        """Convert a list of Pydantic models to a list of dictionaries."""
        return [PydanticModelAdapter.to_dict(model) for model in models]
    
    @staticmethod
    def from_list_of_dicts(model_class: Type[T], data_list: List[Dict[str, Any]]) -> List[T]:
        """Convert a list of dictionaries to a list of Pydantic models."""
        return [PydanticModelAdapter.from_dict(model_class, data) for data in data_list]


class DjangoToPydanticMiddleware:
    """
    Middleware to convert Django request data to Pydantic models.
    
    This middleware can be used to automatically convert Django request data
    to Pydantic models before it reaches the view, and to convert Pydantic
    models to Django-compatible dictionaries before returning the response.
    """
    
    def __init__(self, get_response):
        self.get_response = get_response
    
    def __call__(self, request):
        
        response = self.get_response(request)
        
        
        return response


def pydantic_to_django(model: BaseModel) -> Dict[str, Any]:
    """Convert a Pydantic model to a Django-compatible dictionary."""
    return PydanticModelAdapter.to_dict(model)


def django_to_pydantic(model_class: Type[T], data: Dict[str, Any]) -> T:
    """Convert a Django-compatible dictionary to a Pydantic model."""
    return PydanticModelAdapter.from_dict(model_class, data)

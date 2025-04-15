"""
Views for options and configuration.

This module provides API views for retrieving options and configuration
data, such as available models, agents, and security analyzers.
"""

import logging
from django.http import JsonResponse
from django.views.decorators.http import require_http_methods
from django.views.decorators.csrf import csrf_exempt
from django.contrib.auth.decorators import login_required
from django.conf import settings
from rest_framework.decorators import api_view, permission_classes
from rest_framework.permissions import IsAuthenticated
from rest_framework.response import Response

logger = logging.getLogger(__name__)


@api_view(['GET'])
@permission_classes([IsAuthenticated])
def get_models(request):
    """Get available models."""
    models = [
        {
            'id': 'gpt-4',
            'name': 'GPT-4',
            'provider': 'openai',
            'description': 'OpenAI GPT-4 model',
            'context_length': 8192,
        },
        {
            'id': 'gpt-3.5-turbo',
            'name': 'GPT-3.5 Turbo',
            'provider': 'openai',
            'description': 'OpenAI GPT-3.5 Turbo model',
            'context_length': 4096,
        },
        {
            'id': 'claude-3-opus',
            'name': 'Claude 3 Opus',
            'provider': 'anthropic',
            'description': 'Anthropic Claude 3 Opus model',
            'context_length': 200000,
        },
        {
            'id': 'claude-3-sonnet',
            'name': 'Claude 3 Sonnet',
            'provider': 'anthropic',
            'description': 'Anthropic Claude 3 Sonnet model',
            'context_length': 180000,
        },
        {
            'id': 'llama-3-70b',
            'name': 'Llama 3 70B',
            'provider': 'meta',
            'description': 'Meta Llama 3 70B model',
            'context_length': 8192,
        },
    ]
    
    return Response(models)


@api_view(['GET'])
@permission_classes([IsAuthenticated])
def get_agents(request):
    """Get available agents."""
    agents = [
        {
            'id': 'default',
            'name': 'Default Agent',
            'description': 'Default agent with standard capabilities',
            'model': 'gpt-4',
        },
        {
            'id': 'code-assistant',
            'name': 'Code Assistant',
            'description': 'Specialized agent for code assistance',
            'model': 'gpt-4',
        },
        {
            'id': 'security-analyst',
            'name': 'Security Analyst',
            'description': 'Agent specialized in security analysis',
            'model': 'claude-3-opus',
        },
        {
            'id': 'data-scientist',
            'name': 'Data Scientist',
            'description': 'Agent specialized in data analysis and visualization',
            'model': 'claude-3-sonnet',
        },
    ]
    
    return Response(agents)


@api_view(['GET'])
@permission_classes([IsAuthenticated])
def get_security_analyzers(request):
    """Get available security analyzers."""
    analyzers = [
        {
            'id': 'semgrep',
            'name': 'Semgrep',
            'description': 'Static analysis tool for finding bugs and enforcing code standards',
            'languages': ['python', 'javascript', 'typescript', 'java', 'go'],
        },
        {
            'id': 'bandit',
            'name': 'Bandit',
            'description': 'Security linter for Python code',
            'languages': ['python'],
        },
        {
            'id': 'eslint',
            'name': 'ESLint',
            'description': 'Linter for JavaScript and TypeScript',
            'languages': ['javascript', 'typescript'],
        },
        {
            'id': 'gosec',
            'name': 'Gosec',
            'description': 'Security scanner for Go code',
            'languages': ['go'],
        },
    ]
    
    return Response(analyzers)


@api_view(['GET'])
@permission_classes([IsAuthenticated])
def get_config(request):
    """Get configuration options."""
    config = {
        'max_iterations': 100,
        'max_tokens': 4096,
        'temperature': 0.7,
        'top_p': 0.95,
        'frequency_penalty': 0.0,
        'presence_penalty': 0.0,
        'stop_sequences': [],
        'timeout': 300,
        'features': {
            'code_execution': True,
            'file_upload': True,
            'web_search': True,
            'web_browsing': True,
            'image_generation': True,
            'voice_input': False,
            'voice_output': False,
        },
    }
    
    return Response(config)

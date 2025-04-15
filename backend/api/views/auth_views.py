"""
Views for authentication.

This module provides API views for authentication, including
login, logout, registration, and OAuth authentication with
GitHub and Gitee.
"""

import json
import logging
import requests
from typing import Dict, Any, Optional
from django.http import JsonResponse, HttpResponse
from django.views.decorators.http import require_http_methods
from django.views.decorators.csrf import csrf_exempt
from django.contrib.auth import authenticate, login, logout
from django.contrib.auth.models import User
from django.contrib.auth.decorators import login_required
from django.conf import settings
from django.shortcuts import redirect
from rest_framework.decorators import api_view, permission_classes
from rest_framework.permissions import IsAuthenticated, AllowAny
from rest_framework.response import Response
from rest_framework_simplejwt.tokens import RefreshToken

from api.models import GitHubToken, GiteeToken, UserSettings

logger = logging.getLogger(__name__)


@api_view(['POST'])
@permission_classes([AllowAny])
def login_view(request):
    """Login a user."""
    try:
        username = request.data.get('username')
        password = request.data.get('password')
        
        if not username or not password:
            return Response({
                'status': 'error',
                'message': 'Missing username or password'
            }, status=400)
        
        user = authenticate(username=username, password=password)
        
        if not user:
            return Response({
                'status': 'error',
                'message': 'Invalid username or password'
            }, status=401)
        
        login(request, user)
        
        refresh = RefreshToken.for_user(user)
        
        return Response({
            'status': 'success',
            'message': 'Login successful',
            'user': {
                'id': user.id,
                'username': user.username,
                'email': user.email,
                'first_name': user.first_name,
                'last_name': user.last_name,
            },
            'tokens': {
                'refresh': str(refresh),
                'access': str(refresh.access_token),
            }
        })
    
    except Exception as e:
        logger.error(f"Error logging in: {e}")
        return Response({
            'status': 'error',
            'message': str(e)
        }, status=500)


@api_view(['POST'])
@permission_classes([IsAuthenticated])
def logout_view(request):
    """Logout a user."""
    try:
        logout(request)
        
        return Response({
            'status': 'success',
            'message': 'Logout successful'
        })
    
    except Exception as e:
        logger.error(f"Error logging out: {e}")
        return Response({
            'status': 'error',
            'message': str(e)
        }, status=500)


@api_view(['POST'])
@permission_classes([AllowAny])
def register_view(request):
    """Register a new user."""
    try:
        username = request.data.get('username')
        password = request.data.get('password')
        email = request.data.get('email')
        first_name = request.data.get('first_name', '')
        last_name = request.data.get('last_name', '')
        
        if not username or not password or not email:
            return Response({
                'status': 'error',
                'message': 'Missing required fields'
            }, status=400)
        
        if User.objects.filter(username=username).exists():
            return Response({
                'status': 'error',
                'message': 'Username already exists'
            }, status=400)
        
        if User.objects.filter(email=email).exists():
            return Response({
                'status': 'error',
                'message': 'Email already exists'
            }, status=400)
        
        user = User.objects.create_user(
            username=username,
            password=password,
            email=email,
            first_name=first_name,
            last_name=last_name
        )
        
        UserSettings.objects.create(user=user)
        
        refresh = RefreshToken.for_user(user)
        
        return Response({
            'status': 'success',
            'message': 'Registration successful',
            'user': {
                'id': user.id,
                'username': user.username,
                'email': user.email,
                'first_name': user.first_name,
                'last_name': user.last_name,
            },
            'tokens': {
                'refresh': str(refresh),
                'access': str(refresh.access_token),
            }
        })
    
    except Exception as e:
        logger.error(f"Error registering user: {e}")
        return Response({
            'status': 'error',
            'message': str(e)
        }, status=500)


@api_view(['GET'])
@permission_classes([IsAuthenticated])
def user_view(request):
    """Get the current user."""
    try:
        user = request.user
        
        return Response({
            'status': 'success',
            'user': {
                'id': user.id,
                'username': user.username,
                'email': user.email,
                'first_name': user.first_name,
                'last_name': user.last_name,
            }
        })
    
    except Exception as e:
        logger.error(f"Error getting user: {e}")
        return Response({
            'status': 'error',
            'message': str(e)
        }, status=500)


@api_view(['GET'])
@permission_classes([IsAuthenticated])
def user_settings_view(request):
    """Get the current user's settings."""
    try:
        settings, created = UserSettings.objects.get_or_create(user=request.user)
        
        return Response({
            'status': 'success',
            'settings': settings.settings
        })
    
    except Exception as e:
        logger.error(f"Error getting user settings: {e}")
        return Response({
            'status': 'error',
            'message': str(e)
        }, status=500)


@api_view(['PUT'])
@permission_classes([IsAuthenticated])
def update_user_settings_view(request):
    """Update the current user's settings."""
    try:
        settings, created = UserSettings.objects.get_or_create(user=request.user)
        
        settings.settings = request.data.get('settings', {})
        settings.save()
        
        return Response({
            'status': 'success',
            'message': 'Settings updated',
            'settings': settings.settings
        })
    
    except Exception as e:
        logger.error(f"Error updating user settings: {e}")
        return Response({
            'status': 'error',
            'message': str(e)
        }, status=500)


@api_view(['GET'])
@permission_classes([AllowAny])
def github_login_view(request):
    """Redirect to GitHub OAuth login."""
    try:
        github_oauth_url = f"https://github.com/login/oauth/authorize?client_id={settings.GITHUB_CLIENT_ID}&redirect_uri={settings.GITHUB_REDIRECT_URI}&scope=user,repo"
        
        return Response({
            'status': 'success',
            'oauth_url': github_oauth_url
        })
    
    except Exception as e:
        logger.error(f"Error redirecting to GitHub OAuth: {e}")
        return Response({
            'status': 'error',
            'message': str(e)
        }, status=500)


@api_view(['GET'])
@permission_classes([AllowAny])
def github_callback_view(request):
    """Handle GitHub OAuth callback."""
    try:
        code = request.query_params.get('code')
        
        if not code:
            return Response({
                'status': 'error',
                'message': 'Missing code'
            }, status=400)
        
        response = requests.post(
            'https://github.com/login/oauth/access_token',
            headers={
                'Accept': 'application/json'
            },
            data={
                'client_id': settings.GITHUB_CLIENT_ID,
                'client_secret': settings.GITHUB_CLIENT_SECRET,
                'code': code,
                'redirect_uri': settings.GITHUB_REDIRECT_URI
            }
        )
        
        if response.status_code != 200:
            return Response({
                'status': 'error',
                'message': 'Error exchanging code for access token'
            }, status=500)
        
        token_data = response.json()
        access_token = token_data.get('access_token')
        
        if not access_token:
            return Response({
                'status': 'error',
                'message': 'Missing access token'
            }, status=500)
        
        user_response = requests.get(
            'https://api.github.com/user',
            headers={
                'Authorization': f"token {access_token}"
            }
        )
        
        if user_response.status_code != 200:
            return Response({
                'status': 'error',
                'message': 'Error getting user data from GitHub'
            }, status=500)
        
        user_data = user_response.json()
        github_id = user_data.get('id')
        github_username = user_data.get('login')
        github_email = user_data.get('email')
        
        if not github_id or not github_username:
            return Response({
                'status': 'error',
                'message': 'Missing GitHub user data'
            }, status=500)
        
        user = None
        
        if github_email:
            user = User.objects.filter(email=github_email).first()
        
        if not user:
            username = github_username
            i = 1
            while User.objects.filter(username=username).exists():
                username = f"{github_username}{i}"
                i += 1
            
            user = User.objects.create_user(
                username=username,
                email=github_email or f"{username}@github.com",
                password=None
            )
            
            UserSettings.objects.create(user=user)
        
        github_token, created = GitHubToken.objects.get_or_create(
            user=user,
            defaults={
                'access_token': access_token,
                'refresh_token': token_data.get('refresh_token'),
                'expires_at': None
            }
        )
        
        if not created:
            github_token.access_token = access_token
            github_token.refresh_token = token_data.get('refresh_token')
            github_token.save()
        
        login(request, user)
        
        refresh = RefreshToken.for_user(user)
        
        redirect_url = f"{settings.FRONTEND_URL}/auth/callback?access_token={refresh.access_token}&refresh_token={refresh}"
        
        return redirect(redirect_url)
    
    except Exception as e:
        logger.error(f"Error handling GitHub callback: {e}")
        return Response({
            'status': 'error',
            'message': str(e)
        }, status=500)


@api_view(['GET'])
@permission_classes([AllowAny])
def gitee_login_view(request):
    """Redirect to Gitee OAuth login."""
    try:
        gitee_oauth_url = f"https://gitee.com/oauth/authorize?client_id={settings.GITEE_CLIENT_ID}&redirect_uri={settings.GITEE_REDIRECT_URI}&response_type=code"
        
        return Response({
            'status': 'success',
            'oauth_url': gitee_oauth_url
        })
    
    except Exception as e:
        logger.error(f"Error redirecting to Gitee OAuth: {e}")
        return Response({
            'status': 'error',
            'message': str(e)
        }, status=500)


@api_view(['GET'])
@permission_classes([AllowAny])
def gitee_callback_view(request):
    """Handle Gitee OAuth callback."""
    try:
        code = request.query_params.get('code')
        
        if not code:
            return Response({
                'status': 'error',
                'message': 'Missing code'
            }, status=400)
        
        response = requests.post(
            'https://gitee.com/oauth/token',
            data={
                'grant_type': 'authorization_code',
                'code': code,
                'client_id': settings.GITEE_CLIENT_ID,
                'client_secret': settings.GITEE_CLIENT_SECRET,
                'redirect_uri': settings.GITEE_REDIRECT_URI
            }
        )
        
        if response.status_code != 200:
            return Response({
                'status': 'error',
                'message': 'Error exchanging code for access token'
            }, status=500)
        
        token_data = response.json()
        access_token = token_data.get('access_token')
        
        if not access_token:
            return Response({
                'status': 'error',
                'message': 'Missing access token'
            }, status=500)
        
        user_response = requests.get(
            'https://gitee.com/api/v5/user',
            params={
                'access_token': access_token
            }
        )
        
        if user_response.status_code != 200:
            return Response({
                'status': 'error',
                'message': 'Error getting user data from Gitee'
            }, status=500)
        
        user_data = user_response.json()
        gitee_id = user_data.get('id')
        gitee_username = user_data.get('login')
        gitee_email = user_data.get('email')
        
        if not gitee_id or not gitee_username:
            return Response({
                'status': 'error',
                'message': 'Missing Gitee user data'
            }, status=500)
        
        user = None
        
        if gitee_email:
            user = User.objects.filter(email=gitee_email).first()
        
        if not user:
            username = gitee_username
            i = 1
            while User.objects.filter(username=username).exists():
                username = f"{gitee_username}{i}"
                i += 1
            
            user = User.objects.create_user(
                username=username,
                email=gitee_email or f"{username}@gitee.com",
                password=None
            )
            
            UserSettings.objects.create(user=user)
        
        gitee_token, created = GiteeToken.objects.get_or_create(
            user=user,
            defaults={
                'access_token': access_token,
                'refresh_token': token_data.get('refresh_token'),
                'expires_at': None
            }
        )
        
        if not created:
            gitee_token.access_token = access_token
            gitee_token.refresh_token = token_data.get('refresh_token')
            gitee_token.save()
        
        login(request, user)
        
        refresh = RefreshToken.for_user(user)
        
        redirect_url = f"{settings.FRONTEND_URL}/auth/callback?access_token={refresh.access_token}&refresh_token={refresh}"
        
        return redirect(redirect_url)
    
    except Exception as e:
        logger.error(f"Error handling Gitee callback: {e}")
        return Response({
            'status': 'error',
            'message': str(e)
        }, status=500)

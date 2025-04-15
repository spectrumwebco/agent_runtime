"""
Serializers for the API app.
"""

from rest_framework import serializers
from django.contrib.auth.models import User
from api.models import (
    Conversation, ConversationEvent, UserSettings,
    GitHubToken, GiteeToken, UserCredits, CreditTransaction,
    PolarSubscription
)


class UserSerializer(serializers.ModelSerializer):
    """Serializer for the User model."""

    class Meta:
        model = User
        fields = ['id', 'username', 'email', 'first_name', 'last_name']
        read_only_fields = ['id']


class UserSettingsSerializer(serializers.ModelSerializer):
    """Serializer for UserSettings model."""
    
    class Meta:
        model = UserSettings
        fields = ['id', 'user', 'settings', 'created_at', 'updated_at']
        read_only_fields = ['id', 'user', 'created_at', 'updated_at']


class GitHubTokenSerializer(serializers.ModelSerializer):
    """Serializer for GitHubToken model."""
    
    class Meta:
        model = GitHubToken
        fields = ['id', 'user', 'access_token', 'refresh_token', 'expires_at', 'created_at', 'updated_at']
        read_only_fields = ['id', 'user', 'created_at', 'updated_at']
        extra_kwargs = {
            'access_token': {'write_only': True},
            'refresh_token': {'write_only': True}
        }


class GiteeTokenSerializer(serializers.ModelSerializer):
    """Serializer for GiteeToken model."""
    
    class Meta:
        model = GiteeToken
        fields = ['id', 'user', 'access_token', 'refresh_token', 'expires_at', 'created_at', 'updated_at']
        read_only_fields = ['id', 'user', 'created_at', 'updated_at']
        extra_kwargs = {
            'access_token': {'write_only': True},
            'refresh_token': {'write_only': True}
        }


class ConversationSerializer(serializers.ModelSerializer):
    """Serializer for Conversation model."""
    
    class Meta:
        model = Conversation
        fields = [
            'id', 'user', 'title', 'created_at', 'updated_at',
            'is_active', 'model', 'agent', 'workspace_path'
        ]
        read_only_fields = ['id', 'user', 'created_at', 'updated_at']


class ConversationEventSerializer(serializers.ModelSerializer):
    """Serializer for ConversationEvent model."""
    
    content_json = serializers.SerializerMethodField()
    
    class Meta:
        model = ConversationEvent
        fields = [
            'id', 'conversation_id', 'event_id', 'event_type',
            'source', 'content', 'content_json', 'timestamp', 'created_at'
        ]
        read_only_fields = ['id', 'created_at']
    
    def get_content_json(self, obj):
        """Get the content as JSON."""
        import json
        try:
            return json.loads(obj.content)
        except:
            return None


class UserCreditsSerializer(serializers.ModelSerializer):
    """Serializer for UserCredits model."""
    
    class Meta:
        model = UserCredits
        fields = ['id', 'user', 'balance', 'created_at', 'updated_at']
        read_only_fields = ['id', 'user', 'created_at', 'updated_at']


class CreditTransactionSerializer(serializers.ModelSerializer):
    """Serializer for CreditTransaction model."""
    
    class Meta:
        model = CreditTransaction
        fields = ['id', 'user', 'amount', 'transaction_type', 'description', 'created_at']
        read_only_fields = ['id', 'user', 'created_at']


class PolarSubscriptionSerializer(serializers.ModelSerializer):
    """Serializer for PolarSubscription model."""
    
    class Meta:
        model = PolarSubscription
        fields = [
            'id', 'user', 'subscription_id', 'plan_id', 'status',
            'current_period_start', 'current_period_end',
            'cancel_at_period_end', 'created_at', 'updated_at'
        ]
        read_only_fields = ['id', 'user', 'created_at', 'updated_at']

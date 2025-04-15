"""
Models for the API app.

This module defines the database models for the API app, including
conversation data, events, and user settings.
"""

from django.db import models
from django.contrib.auth.models import User
import uuid
import json


class Conversation(models.Model):
    """Model for storing conversation data."""
    
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    user = models.ForeignKey(User, on_delete=models.CASCADE, null=True, blank=True)
    title = models.CharField(max_length=255, default="New Conversation")
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)
    is_active = models.BooleanField(default=True)
    
    model = models.CharField(max_length=100, null=True, blank=True)
    agent = models.CharField(max_length=100, null=True, blank=True)
    workspace_path = models.CharField(max_length=255, null=True, blank=True)
    
    class Meta:
        ordering = ['-updated_at']
    
    def __str__(self):
        return f"{self.title} ({self.id})"
    
    def to_dict(self):
        """Convert the conversation to a dictionary."""
        return {
            'id': str(self.id),
            'title': self.title,
            'created_at': self.created_at.isoformat(),
            'updated_at': self.updated_at.isoformat(),
            'is_active': self.is_active,
            'model': self.model,
            'agent': self.agent,
            'workspace_path': self.workspace_path,
        }


class ConversationEvent(models.Model):
    """Model for storing conversation events."""
    
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    conversation_id = models.CharField(max_length=255)
    event_id = models.IntegerField()
    event_type = models.CharField(max_length=100)
    source = models.CharField(max_length=100)
    content = models.TextField()
    timestamp = models.BigIntegerField(null=True, blank=True)
    created_at = models.DateTimeField(auto_now_add=True)
    
    class Meta:
        ordering = ['event_id']
        unique_together = ['conversation_id', 'event_id']
    
    def __str__(self):
        return f"Event {self.event_id} in {self.conversation_id}"
    
    def to_dict(self):
        """Convert the event to a dictionary."""
        return json.loads(self.content)


class UserSettings(models.Model):
    """Model for storing user settings."""
    
    user = models.OneToOneField(User, on_delete=models.CASCADE, related_name='settings')
    settings = models.JSONField(default=dict)
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)
    
    def __str__(self):
        return f"Settings for {self.user.username}"
    
    def get_setting(self, key, default=None):
        """Get a specific setting value."""
        return self.settings.get(key, default)
    
    def set_setting(self, key, value):
        """Set a specific setting value."""
        self.settings[key] = value
        self.save()
    
    def delete_setting(self, key):
        """Delete a specific setting."""
        if key in self.settings:
            del self.settings[key]
            self.save()


class GitHubToken(models.Model):
    """Model for storing GitHub access tokens."""
    
    user = models.OneToOneField(User, on_delete=models.CASCADE, related_name='github_token')
    access_token = models.CharField(max_length=255)
    refresh_token = models.CharField(max_length=255, null=True, blank=True)
    expires_at = models.DateTimeField(null=True, blank=True)
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)
    
    def __str__(self):
        return f"GitHub token for {self.user.username}"


class GiteeToken(models.Model):
    """Model for storing Gitee access tokens."""
    
    user = models.OneToOneField(User, on_delete=models.CASCADE, related_name='gitee_token')
    access_token = models.CharField(max_length=255)
    refresh_token = models.CharField(max_length=255, null=True, blank=True)
    expires_at = models.DateTimeField(null=True, blank=True)
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)
    
    def __str__(self):
        return f"Gitee token for {self.user.username}"


class UserCredits(models.Model):
    """Model for storing user credits for billing."""
    
    user = models.OneToOneField(User, on_delete=models.CASCADE, related_name='credits')
    balance = models.DecimalField(max_digits=10, decimal_places=2, default=0)
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)
    
    def __str__(self):
        return f"Credits for {self.user.username}: {self.balance}"


class CreditTransaction(models.Model):
    """Model for storing credit transactions."""
    
    TRANSACTION_TYPES = (
        ('purchase', 'Purchase'),
        ('usage', 'Usage'),
        ('refund', 'Refund'),
        ('bonus', 'Bonus'),
    )
    
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    user = models.ForeignKey(User, on_delete=models.CASCADE, related_name='transactions')
    amount = models.DecimalField(max_digits=10, decimal_places=2)
    transaction_type = models.CharField(max_length=20, choices=TRANSACTION_TYPES)
    description = models.CharField(max_length=255)
    created_at = models.DateTimeField(auto_now_add=True)
    
    class Meta:
        ordering = ['-created_at']
    
    def __str__(self):
        return f"{self.transaction_type} of {self.amount} for {self.user.username}"


class PolarSubscription(models.Model):
    """Model for storing Polar.sh subscription data."""
    
    user = models.OneToOneField(User, on_delete=models.CASCADE, related_name='polar_subscription')
    subscription_id = models.CharField(max_length=255)
    plan_id = models.CharField(max_length=255)
    status = models.CharField(max_length=50)
    current_period_start = models.DateTimeField()
    current_period_end = models.DateTimeField()
    cancel_at_period_end = models.BooleanField(default=False)
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)
    
    def __str__(self):
        return f"Polar subscription for {self.user.username}"

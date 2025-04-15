from django.db import models


class AgentConfig(models.Model):
    """
    Model for storing agent configuration settings.
    """

    name = models.CharField(max_length=100, unique=True)
    description = models.TextField(blank=True)
    config_json = models.JSONField()
    is_active = models.BooleanField(default=True)
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        verbose_name = "Agent Configuration"
        verbose_name_plural = "Agent Configurations"

    def __str__(self):
        return self.name


class AgentSession(models.Model):
    """
    Model for tracking agent execution sessions.
    """

    session_id = models.CharField(max_length=100, unique=True)
    agent_config = models.ForeignKey(
        AgentConfig, on_delete=models.CASCADE, related_name="sessions"
    )
    status = models.CharField(
        max_length=20,
        choices=[
            ("pending", "Pending"),
            ("running", "Running"),
            ("completed", "Completed"),
            ("failed", "Failed"),
        ],
        default="pending",
    )
    start_time = models.DateTimeField(auto_now_add=True)
    end_time = models.DateTimeField(null=True, blank=True)
    result_json = models.JSONField(null=True, blank=True)
    error_message = models.TextField(blank=True)

    class Meta:
        verbose_name = "Agent Session"
        verbose_name_plural = "Agent Sessions"

    def __str__(self):
        return f"Session {self.session_id} ({self.status})"


class AgentEvent(models.Model):
    """
    Model for tracking events during agent execution.
    """

    session = models.ForeignKey(
        AgentSession, on_delete=models.CASCADE, related_name="events"
    )
    event_type = models.CharField(max_length=50)
    event_data = models.JSONField()
    timestamp = models.DateTimeField(auto_now_add=True)

    class Meta:
        verbose_name = "Agent Event"
        verbose_name_plural = "Agent Events"
        ordering = ["timestamp"]

    def __str__(self):
        return f"{self.event_type} at {self.timestamp}"

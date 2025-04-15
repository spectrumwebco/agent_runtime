from django.urls import path, include
from . import views

app_name = "python_agent"

urlpatterns = [
    path("run/", views.run_agent_view, name="run_agent"),
    path("agent_framework/", include("apps.python_agent.agent_framework.django_views.urls")),
]

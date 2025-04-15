from django.urls import path
from . import views

app_name = "python_agent"

urlpatterns = [
    path("run/", views.run_agent_view, name="run_agent"),
]

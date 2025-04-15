from django.urls import path
from . import views

app_name = "python_ml"

urlpatterns = [
    path("inference/", views.run_inference_view, name="run_inference"),
    path("fine-tune/", views.fine_tune_model_view, name="fine_tune_model"),
]

"""
URL configuration for the ML API app.
"""

from django.urls import path, include
from rest_framework.routers import DefaultRouter
from . import views

router = DefaultRouter()
router.register(r'models', views.MLModelViewSet, basename='model')
router.register(
    r'fine-tuning-jobs',
    views.FineTuningJobViewSet,
    basename='fine-tuning-job')

urlpatterns = [
    path('', views.ml_api_root, name='ml-api-root'),
    path('', include(router.urls)),
]

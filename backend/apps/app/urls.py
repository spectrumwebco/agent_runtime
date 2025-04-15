from django.urls import path
from . import views

app_name = 'app'

urlpatterns = [
    path('status/', views.app_status_view, name='status'),
]

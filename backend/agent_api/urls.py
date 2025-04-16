"""
URL configuration for agent_api project.

The `urlpatterns` list routes URLs to views. For more information please see:
    https://docs.djangoproject.com/en/5.2/topics/http/urls/
Examples:
Function views
    1. Add an import:  from my_app import views
    2. Add a URL to urlpatterns:  path('', views.home, name='home')
Class-based views
    1. Add an import:  from other_app.views import Home
    2. Add a URL to urlpatterns:  path('', Home.as_view(), name='home')
Including another URLconf
    1. Import the include() function: from django.urls import include, path
    2. Add a URL to urlpatterns:  path('blog/', include('blog.urls'))
"""
from django.contrib import admin
from django.urls import path, include
from django.views.generic import RedirectView
from api.ninja_api import api as ninja_api
from api.grpc_service import router as grpc_router
from api.swagger import urlpatterns as swagger_urls

# Add router to ninja_api only if it hasn't been added already
try:
    ninja_api.add_router("/grpc", grpc_router)
except Exception as e:
    pass

urlpatterns = [
    path('admin/', admin.site.urls),
    path('api/', include('api.urls')),
    path('ml-api/', include('ml_api.urls')),
    path('ninja-api/', ninja_api.urls, name='ninja-api'),
    path('docs/', include(swagger_urls)),
    path('agent/', include('apps.python_agent.urls')),
    path('ml/', include('apps.python_ml.urls')),
    path('tools/', include('apps.python_agent.tools.urls')),
    path('app/', include('apps.app.urls')),
    path('', RedirectView.as_view(url='/api/', permanent=False)),
]

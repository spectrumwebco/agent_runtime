"""
URL configuration for agent_api project.
"""
from django.contrib import admin
from django.urls import path, include
from django.conf import settings
from django.conf.urls.static import static
from api.swagger import urlpatterns as swagger_urls

urlpatterns = [
    path('admin/', admin.site.urls),
    path('api/', include('api.urls')),
    path('api/app/', include('apps.app.urls')),
    path('api/agent/', include('apps.python_agent.urls')),
    path('api/ml/', include('apps.python_ml.urls')),
    path('api/tools/', include('apps.tools.urls')),
] + swagger_urls

if settings.DEBUG:
    urlpatterns += static(settings.STATIC_URL, document_root=settings.STATIC_ROOT)

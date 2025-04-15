from django.http import JsonResponse
from django.views.decorators.csrf import csrf_exempt
from django.views.decorators.http import require_http_methods
import logging

logger = logging.getLogger(__name__)


@csrf_exempt
@require_http_methods(["GET"])
def app_status_view(request):
    """
    Simple view to check the status of the main app.
    """
    return JsonResponse(
        {"status": "success", "message": "Main app is properly integrated with Django"}
    )

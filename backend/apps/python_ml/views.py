from django.http import JsonResponse
from django.views.decorators.csrf import csrf_exempt
from django.views.decorators.http import require_http_methods
import json
import logging
import asyncio

from apps.python_ml.api.client import MLInfrastructureAPIClient

logger = logging.getLogger(__name__)
client = MLInfrastructureAPIClient()


@csrf_exempt
@require_http_methods(["POST"])
def run_inference_view(request):
    """
    Django view to run ML inference with the provided data.
    """
    try:
        data = json.loads(request.body)

        service_id = data.get("service_id")
        input_text = data.get("input_text")
        parameters = data.get("parameters")

        if not service_id or not input_text:
            return JsonResponse(
                {
                    "status": "error",
                    "message": "Missing required parameters: service_id and input_text",
                },
                status=400,
            )

        loop = asyncio.new_event_loop()
        asyncio.set_event_loop(loop)

        result = loop.run_until_complete(
            client.predict(service_id, input_text, parameters)
        )

        if hasattr(result, "model_dump"):
            result_dict = result.model_dump()
        else:
            result_dict = result.dict()

        return JsonResponse({"status": "success", "result": result_dict})
    except Exception as e:
        logger.exception("Error running ML inference")
        return JsonResponse({"status": "error", "message": str(e)}, status=500)


@csrf_exempt
@require_http_methods(["POST"])
def fine_tune_model_view(request):
    """
    Django view to fine-tune a model with the provided data.
    """
    try:
        data = json.loads(request.body)

        model_type = data.get("model_type")
        training_data_path = data.get("training_data_path")
        validation_data_path = data.get("validation_data_path")
        hyperparameters = data.get("hyperparameters")

        if not model_type or not training_data_path:
            return JsonResponse(
                {
                    "status": "error",
                    "message": "Missing required parameters: model_type and training_data_path",
                },
                status=400,
            )

        loop = asyncio.new_event_loop()
        asyncio.set_event_loop(loop)

        result = loop.run_until_complete(
            client.create_fine_tuning_job(
                model_type, training_data_path, validation_data_path, hyperparameters
            )
        )

        if hasattr(result, "model_dump"):
            result_dict = result.model_dump()
        else:
            result_dict = result.dict()

        return JsonResponse({"status": "success", "result": result_dict})
    except Exception as e:
        logger.exception("Error fine-tuning model")
        return JsonResponse({"status": "error", "message": str(e)}, status=500)

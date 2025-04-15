from django.http import JsonResponse
from django.views.decorators.csrf import csrf_exempt
from django.views.decorators.http import require_http_methods
import json
import logging
import asyncio
import os

from apps.python_ml.api.client import MLInfrastructureAPIClient
from .integration.eventstream_integration import event_stream, Event, EventType, EventSource
from .integration.k8s_integration import k8s_client

logger = logging.getLogger(__name__)
client = MLInfrastructureAPIClient()


@csrf_exempt
@require_http_methods(["POST"])
def run_inference_view(request):
    """
    Run ML inference with the provided data.
    """
    try:
        data = json.loads(request.body)

        service_id = data.get("service_id")
        input_text = data.get("input_text")
        parameters = data.get("parameters", {})

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
        
        event_data = {
            'service_id': service_id,
            'input_length': len(input_text)
        }
        loop.run_until_complete(event_stream.publish(Event.new(
            EventType.ACTION,
            EventSource.ML,
            event_data,
            {'action': 'inference_request'}
        )))

        result = loop.run_until_complete(
            client.predict(service_id, input_text, parameters)
        )

        if hasattr(result, "model_dump"):
            result_dict = result.model_dump()
        else:
            result_dict = result.dict()
            
        event_data = {
            'service_id': service_id,
            'latency_ms': result_dict.get('latency_ms', 0)
        }
        loop.run_until_complete(event_stream.publish(Event.new(
            EventType.OBSERVATION,
            EventSource.ML,
            event_data,
            {'action': 'inference_response'}
        )))

        return JsonResponse({"status": "success", "result": result_dict})
    except Exception as e:
        logger.exception("Error running ML inference")
        return JsonResponse({"status": "error", "message": str(e)}, status=500)


@csrf_exempt
@require_http_methods(["POST"])
def fine_tune_model_view(request):
    """
    Fine-tune a model with the provided data.
    """
    try:
        data = json.loads(request.body)

        model_type = data.get("model_type")
        training_data_path = data.get("training_data_path")
        validation_data_path = data.get("validation_data_path")
        hyperparameters = data.get("hyperparameters", {})

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
        
        event_data = {
            'model_type': model_type,
            'training_data_path': training_data_path
        }
        loop.run_until_complete(event_stream.publish(Event.new(
            EventType.ACTION,
            EventSource.ML,
            event_data,
            {'action': 'fine_tuning_request'}
        )))

        result = loop.run_until_complete(
            client.create_fine_tuning_job(
                model_type, training_data_path, validation_data_path, hyperparameters
            )
        )

        if hasattr(result, "model_dump"):
            result_dict = result.model_dump()
        else:
            result_dict = result.dict()
            
        event_data = {
            'job_id': result_dict.get('id', ''),
            'model_type': model_type,
            'status': result_dict.get('status', 'created')
        }
        loop.run_until_complete(event_stream.publish(Event.new(
            EventType.STATE_UPDATE,
            EventSource.ML,
            event_data,
            {'action': 'fine_tuning_job_created'}
        )))

        return JsonResponse({"status": "success", "result": result_dict})
    except Exception as e:
        logger.exception("Error fine-tuning model")
        return JsonResponse({"status": "error", "message": str(e)}, status=500)


@csrf_exempt
@require_http_methods(["POST"])
def deploy_ml_infrastructure_view(request):
    """
    Deploy ML infrastructure to Kubernetes.
    """
    try:
        data = json.loads(request.body)
        
        namespace = data.get("namespace", "ml-infrastructure")
        components = data.get("components", ["mlflow", "kubeflow", "kserve"])
        
        loop = asyncio.new_event_loop()
        asyncio.set_event_loop(loop)
        
        event_data = {
            'namespace': namespace,
            'components': components
        }
        loop.run_until_complete(event_stream.publish(Event.new(
            EventType.ACTION,
            EventSource.ML,
            event_data,
            {'action': 'deploy_ml_infrastructure'}
        )))
        
        results = {}
        if "mlflow" in components:
            results["mlflow"] = loop.run_until_complete(k8s_client.deploy_mlflow(namespace))
            
        if "kubeflow" in components:
            results["kubeflow"] = loop.run_until_complete(k8s_client.deploy_kubeflow(namespace))
            
        if "kserve" in components:
            results["kserve"] = loop.run_until_complete(k8s_client.deploy_kserve(namespace))
        
        return JsonResponse({
            "status": "success",
            "results": results
        })
    except Exception as e:
        logger.exception("Error deploying ML infrastructure")
        return JsonResponse({
            "status": "error", 
            "message": str(e)
        }, status=500)

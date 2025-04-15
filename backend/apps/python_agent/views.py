from django.http import JsonResponse
from django.views.decorators.csrf import csrf_exempt
from django.views.decorators.http import require_http_methods
import json
import logging
import asyncio
import os

from apps.python_agent.agent_framework.runtime.abstract import (
    Action,
    Command,
    CreateSessionRequest,
)
from apps.python_agent.agent_framework.runtime.local import LocalRuntime

logger = logging.getLogger(__name__)
runtime = LocalRuntime()


@csrf_exempt
@require_http_methods(["POST"])
def run_agent_view(request):
    """
    Django view to run the agent with the provided configuration.
    """
    try:
        data = json.loads(request.body)

        loop = asyncio.new_event_loop()
        asyncio.set_event_loop(loop)

        session_request = CreateSessionRequest(
            working_directory=data.get("working_directory", os.getcwd()),
            session_id=data.get("session_id", None),
        )
        session_response = loop.run_until_complete(
            runtime.create_session(session_request)
        )
        session_id = session_response.session_id

        action = Action(
            session_id=session_id,
            command=Command(command=data.get("command", ""), args=data.get("args", [])),
        )
        result = loop.run_until_complete(runtime.run_in_session(action))

        if hasattr(result, "model_dump"):
            result_dict = result.model_dump()
        else:
            result_dict = result.dict()

        return JsonResponse(
            {"status": "success", "session_id": session_id, "result": result_dict}
        )
    except Exception as e:
        logger.exception("Error running agent")
        return JsonResponse({"status": "error", "message": str(e)}, status=500)

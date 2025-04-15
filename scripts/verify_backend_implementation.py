"""
Script to verify the backend implementation.

This script checks that all necessary components for the backend implementation
are in place and properly integrated.
"""

import os
import sys
import json
from pathlib import Path


def check_file_exists(path):
    """Check if a file exists."""
    return os.path.exists(path)


def main():
    """Verify the backend implementation."""
    root_dir = Path(__file__).resolve().parent.parent
    
    components = {
        "Django Backend": {
            "files": [
                "backend/agent_api/settings.py",
                "backend/agent_api/urls.py",
                "backend/agent_api/asgi.py",
                "backend/api/urls.py",
                "backend/api/views/health_views.py",
                "backend/api/views/events_views.py",
                "backend/api/views/auth_views.py",
                "backend/api/views/conversation_views.py",
                "backend/api/views/options_views.py",
                "backend/api/views/billing_views.py",
                "backend/api/websocket.py",
                "backend/api/websocket_state.py",
                "backend/api/socketio_consumer.py",
                "backend/api/grpc_client.py",
            ],
            "status": "Not checked"
        },
        "Go gRPC Bridge": {
            "files": [
                "protos/agent_bridge.proto",
                "internal/kled/grpc/bridge.go",
                "internal/server/socketio_routes.go",
                "internal/kled/socketio/server.go",
                "internal/kled/socketio/middleware.go",
                "internal/kled/routes/socketio.go",
            ],
            "status": "Not checked"
        },
        "Kubernetes Configuration": {
            "files": [
                "k8s/django/django-deployment.yaml",
                "k8s/django/django-service.yaml",
                "k8s/django/django-config.yaml",
                "k8s/django/django-ingress.yaml",
                "k8s/django/django-pvc.yaml",
            ],
            "status": "Not checked"
        },
        "Terraform Configuration": {
            "files": [
                "terraform/modules/django/main.tf",
                "terraform/modules/django/variables.tf",
                "terraform/modules/django/outputs.tf",
                "terraform/modules/django/files/settings.py",
                "terraform/modules/django/files/urls.py",
            ],
            "status": "Not checked"
        },
        "Python ML App": {
            "files": [
                "backend/apps/python_ml/apps.py",
                "backend/apps/python_ml/urls.py",
                "backend/apps/python_ml/views.py",
                "backend/apps/python_ml/integration/eventstream_integration.py",
                "backend/apps/python_ml/integration/k8s_integration.py",
                "backend/apps/python_ml/scrapers/github_scraper.py",
                "backend/apps/python_ml/scrapers/gitee_scraper.py",
                "backend/apps/python_ml/trajectories/generator.py",
                "backend/apps/python_ml/benchmarking/historical_benchmark.py",
            ],
            "status": "Not checked"
        },
        "Python Agent App": {
            "files": [
                "backend/apps/python_agent/apps.py",
                "backend/apps/python_agent/urls.py",
                "backend/apps/python_agent/views.py",
            ],
            "status": "Not checked"
        },
        "Tools App": {
            "files": [
                "backend/apps/tools/apps.py",
                "backend/apps/tools/urls.py",
                "backend/apps/tools/views.py",
            ],
            "status": "Not checked"
        },
        "App App": {
            "files": [
                "backend/apps/app/apps.py",
                "backend/apps/app/urls.py",
                "backend/apps/app/views.py",
                "backend/apps/app/models.py",
            ],
            "status": "Not checked"
        },
        "Protobuf Generation": {
            "files": [
                "scripts/generate_protos.py",
                "scripts/install_proto_deps.py",
                "protos/gen/python/__init__.py",
                "protos/gen/go/__init__.go",
            ],
            "status": "Not checked"
        },
        "Kled.io Framework": {
            "files": [
                "cmd/kled/main.go",
                "internal/kled/framework.go",
                "internal/kled/django.go",
                "internal/kled/pytorch.go",
                "internal/kled/framework_test.go",
                "internal/kled/django_test.go",
                "internal/kled/pytorch_test.go",
                "internal/kled/scripts/pytorch_tools.py",
            ],
            "status": "Not checked"
        }
    }
    
    for component, data in components.items():
        missing_files = []
        for file in data["files"]:
            file_path = os.path.join(root_dir, file)
            if not check_file_exists(file_path):
                missing_files.append(file)
        
        if missing_files:
            data["status"] = "Incomplete"
            data["missing_files"] = missing_files
        else:
            data["status"] = "Complete"
    
    print("Backend Implementation Verification Report")
    print("==========================================")
    print()
    
    all_complete = True
    for component, data in components.items():
        print(f"{component}: {data['status']}")
        if data["status"] == "Incomplete":
            all_complete = False
            print("  Missing files:")
            for file in data["missing_files"]:
                print(f"    - {file}")
        print()
    
    if all_complete:
        print("All components are complete!")
        return 0
    else:
        print("Some components are incomplete. Please check the missing files.")
        return 1


if __name__ == "__main__":
    sys.exit(main())

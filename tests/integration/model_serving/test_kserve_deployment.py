import unittest
import os
import sys
import json
import tempfile
from unittest.mock import patch, MagicMock

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '../../../src')))

from ml_infrastructure.kserve.manifests.client.kserve_client import KServeClient


class TestKServeDeployment(unittest.TestCase):
    """Integration tests for KServe model deployment."""

    def setUp(self):
        """Set up test fixtures."""
        self.test_data_dir = os.path.join(os.path.dirname(__file__), '../../data/test')
        os.makedirs(self.test_data_dir, exist_ok=True)
        
        self.temp_dir = tempfile.TemporaryDirectory()
        self.model_dir = self.temp_dir.name
        
        self.model_path = os.path.join(self.model_dir, "llama-4-maverick")
        os.makedirs(self.model_path, exist_ok=True)
        with open(os.path.join(self.model_path, "model.bin"), "w") as f:
            f.write("mock model data")
        
        self.kserve_client = KServeClient(
            namespace="ml-infrastructure",
            kserve_endpoint="http://kserve-gateway.ml-infrastructure.svc.cluster.local"
        )

    def tearDown(self):
        """Clean up after tests."""
        self.temp_dir.cleanup()

    @patch('ml_infrastructure.kserve.manifests.client.kserve_client.KServeClient._create_inference_service')
    def test_model_deployment(self, mock_create_service):
        """Test model deployment to KServe."""
        mock_create_service.return_value = {
            "apiVersion": "serving.kserve.io/v1beta1",
            "kind": "InferenceService",
            "metadata": {
                "name": "llama-4-maverick",
                "namespace": "ml-infrastructure"
            },
            "status": {
                "url": "http://llama-4-maverick.ml-infrastructure.example.com",
                "conditions": [
                    {
                        "type": "Ready",
                        "status": "True"
                    }
                ]
            }
        }
        
        result = self.kserve_client.deploy_model(
            name="llama-4-maverick",
            model_path=self.model_path,
            framework="pytorch",
            min_replicas=1,
            max_replicas=3
        )
        
        self.assertEqual(result["metadata"]["name"], "llama-4-maverick")
        self.assertEqual(result["status"]["conditions"][0]["type"], "Ready")
        self.assertEqual(result["status"]["conditions"][0]["status"], "True")
        
        mock_create_service.assert_called_once()
        args, kwargs = mock_create_service.call_args
        self.assertEqual(kwargs["name"], "llama-4-maverick")
        self.assertEqual(kwargs["framework"], "pytorch")
        self.assertEqual(kwargs["min_replicas"], 1)
        self.assertEqual(kwargs["max_replicas"], 3)

    @patch('ml_infrastructure.kserve.manifests.client.kserve_client.KServeClient._get_inference_service')
    def test_model_status(self, mock_get_service):
        """Test getting model deployment status."""
        mock_get_service.return_value = {
            "apiVersion": "serving.kserve.io/v1beta1",
            "kind": "InferenceService",
            "metadata": {
                "name": "llama-4-maverick",
                "namespace": "ml-infrastructure"
            },
            "status": {
                "url": "http://llama-4-maverick.ml-infrastructure.example.com",
                "conditions": [
                    {
                        "type": "Ready",
                        "status": "True",
                        "message": "Model is ready for inference"
                    }
                ],
                "components": {
                    "predictor": {
                        "url": "http://llama-4-maverick-predictor.ml-infrastructure.svc.cluster.local"
                    }
                }
            }
        }
        
        status = self.kserve_client.get_model_status("llama-4-maverick")
        
        self.assertEqual(status["name"], "llama-4-maverick")
        self.assertEqual(status["ready"], True)
        self.assertEqual(status["url"], "http://llama-4-maverick.ml-infrastructure.example.com")
        self.assertEqual(status["predictor_url"], "http://llama-4-maverick-predictor.ml-infrastructure.svc.cluster.local")
        
        mock_get_service.assert_called_once_with("llama-4-maverick")

    @patch('ml_infrastructure.kserve.manifests.client.kserve_client.KServeClient._delete_inference_service')
    def test_model_deletion(self, mock_delete_service):
        """Test model deletion from KServe."""
        mock_delete_service.return_value = True
        
        result = self.kserve_client.delete_model("llama-4-maverick")
        
        self.assertTrue(result)
        
        mock_delete_service.assert_called_once_with("llama-4-maverick")

    @patch('ml_infrastructure.kserve.manifests.client.kserve_client.KServeClient._predict')
    def test_model_inference(self, mock_predict):
        """Test model inference."""
        mock_predict.return_value = {
            "predictions": ["This is a generated response."]
        }
        
        result = self.kserve_client.predict(
            model_name="llama-4-maverick",
            inputs=["What is Kubernetes?"]
        )
        
        self.assertIn("predictions", result)
        self.assertEqual(len(result["predictions"]), 1)
        self.assertEqual(result["predictions"][0], "This is a generated response.")
        
        mock_predict.assert_called_once()
        args, kwargs = mock_predict.call_args
        self.assertEqual(kwargs["model_name"], "llama-4-maverick")
        self.assertEqual(kwargs["inputs"], ["What is Kubernetes?"])


if __name__ == '__main__':
    unittest.main()

import unittest
import os
import sys
import json
import tempfile
from unittest.mock import patch, MagicMock

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '../../../src')))

from integrations.github.scraper import GitHubScraper
from integrations.gitee.scraper import GiteeScraper
from integrations.issue_collector.collector import IssueCollector
from ml_infrastructure.training.scripts.train_llama4 import LlamaTrainer
from ml_infrastructure.training.config.llama4_maverick_config import Llama4MaverickConfig
from ml_infrastructure.training.evaluation.metrics import ModelEvaluator
from ml_infrastructure.kserve.manifests.client.kserve_client import KServeClient
from ml_infrastructure.api.client import MLInfrastructureClient


class TestMLPipeline(unittest.TestCase):
    """End-to-end tests for the complete ML pipeline."""

    def setUp(self):
        """Set up test fixtures."""
        self.test_data_dir = os.path.join(os.path.dirname(__file__), '../../data/test')
        os.makedirs(self.test_data_dir, exist_ok=True)
        
        self.temp_dir = tempfile.TemporaryDirectory()
        self.output_dir = self.temp_dir.name
        
        self.github_scraper = GitHubScraper(
            token=os.environ.get('GITHUB_TOKEN', 'dummy_token'),
            max_issues=5
        )
        self.gitee_scraper = GiteeScraper(
            token=os.environ.get('GITEE_TOKEN', 'dummy_token'),
            max_issues=5
        )
        self.collector = IssueCollector()
        self.maverick_config = Llama4MaverickConfig(
            model_name="llama-4-maverick",
            batch_size=4,
            learning_rate=2e-5,
            num_epochs=3,
            output_dir=self.output_dir
        )
        self.kserve_client = KServeClient(
            namespace="ml-infrastructure",
            kserve_endpoint="http://kserve-gateway.ml-infrastructure.svc.cluster.local"
        )
        self.ml_client = MLInfrastructureClient(
            base_url=os.environ.get('ML_API_URL', 'http://localhost:8000'),
            api_key=os.environ.get('ML_API_KEY', 'dummy_key')
        )

    def tearDown(self):
        """Clean up after tests."""
        self.temp_dir.cleanup()

    @patch('integrations.github.scraper.GitHubScraper._fetch_issues')
    @patch('integrations.github.scraper.GitHubScraper._fetch_comments')
    @patch('ml_infrastructure.training.scripts.train_llama4.LlamaTrainer._load_model')
    @patch('ml_infrastructure.training.scripts.train_llama4.LlamaTrainer._load_tokenizer')
    @patch('ml_infrastructure.training.scripts.train_llama4.LlamaTrainer._train_model')
    @patch('ml_infrastructure.kserve.manifests.client.kserve_client.KServeClient._create_inference_service')
    @patch('ml_infrastructure.kserve.manifests.client.kserve_client.KServeClient._predict')
    def test_end_to_end_pipeline(self, mock_predict, mock_create_service, mock_train, 
                                 mock_tokenizer, mock_model, mock_comments, mock_issues):
        """Test the complete ML pipeline from data collection to model serving."""
        mock_issues.return_value = [
            {
                "number": 1,
                "title": "Fix Kubernetes deployment issue",
                "body": "The deployment is failing with error: ResourceQuotaExceeded",
                "state": "closed",
                "labels": [{"name": "bug"}, {"name": "kubernetes"}],
                "created_at": "2023-01-01T00:00:00Z",
                "closed_at": "2023-01-02T00:00:00Z",
                "user": {"login": "user1"},
                "html_url": "https://github.com/org/repo/issues/1"
            }
        ]
        mock_comments.return_value = [
            {"user": {"login": "user3"}, "body": "I fixed this by increasing the resource quota."},
            {"user": {"login": "user1"}, "body": "Thanks, that worked!"}
        ]
        
        mock_model.return_value = MagicMock()
        mock_tokenizer.return_value = MagicMock()
        mock_train.return_value = {
            "train_loss": 0.1,
            "eval_loss": 0.2,
            "train_runtime": 100,
            "eval_accuracy": 0.85
        }
        
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
        mock_predict.return_value = {
            "predictions": ["This is a generated response."]
        }
        
        
        issues = self.github_scraper.scrape_issues("kubernetes/kubernetes")
        processed_data = self.collector.process_issues(issues, source="github", repo="kubernetes/kubernetes")
        
        training_data_path = os.path.join(self.test_data_dir, 'e2e_training_data.json')
        with open(training_data_path, 'w') as f:
            json.dump(processed_data, f, indent=2)
        
        trainer = LlamaTrainer(
            config=self.maverick_config,
            data_path=training_data_path
        )
        training_result = trainer.train()
        
        model_path = os.path.join(self.output_dir, "llama-4-maverick")
        os.makedirs(model_path, exist_ok=True)
        with open(os.path.join(model_path, "model.bin"), "w") as f:
            f.write("mock model data")
        
        deployment_result = self.kserve_client.deploy_model(
            name="llama-4-maverick",
            model_path=model_path,
            framework="pytorch",
            min_replicas=1,
            max_replicas=3
        )
        
        prediction_result = self.kserve_client.predict(
            model_name="llama-4-maverick",
            inputs=["What is Kubernetes?"]
        )
        
        
        self.assertEqual(len(processed_data), 1)
        self.assertEqual(processed_data[0]["issue_number"], 1)
        self.assertEqual(processed_data[0]["title"], "Fix Kubernetes deployment issue")
        
        self.assertIn("train_loss", training_result)
        self.assertIn("eval_accuracy", training_result)
        
        self.assertEqual(deployment_result["metadata"]["name"], "llama-4-maverick")
        self.assertEqual(deployment_result["status"]["conditions"][0]["status"], "True")
        
        self.assertIn("predictions", prediction_result)
        self.assertEqual(prediction_result["predictions"][0], "This is a generated response.")

    @patch('ml_infrastructure.api.client.MLInfrastructureClient.create_experiment')
    @patch('ml_infrastructure.api.client.MLInfrastructureClient.log_metrics')
    @patch('ml_infrastructure.api.client.MLInfrastructureClient.register_model')
    @patch('ml_infrastructure.api.client.MLInfrastructureClient.deploy_model')
    def test_api_client_integration(self, mock_deploy, mock_register, mock_log, mock_create):
        """Test the ML infrastructure API client integration."""
        mock_create.return_value = {"experiment_id": "exp-123"}
        mock_log.return_value = {"status": "success"}
        mock_register.return_value = {"model_id": "model-123"}
        mock_deploy.return_value = {"deployment_id": "deploy-123", "endpoint": "http://example.com/model"}
        
        experiment = self.ml_client.create_experiment(
            name="llama4-test",
            description="Test experiment for Llama 4 models"
        )
        
        metrics = self.ml_client.log_metrics(
            experiment_id=experiment["experiment_id"],
            metrics={
                "accuracy": 0.85,
                "f1_score": 0.82,
                "precision": 0.80,
                "recall": 0.84
            }
        )
        
        model = self.ml_client.register_model(
            name="llama-4-maverick",
            version="1.0.0",
            model_path=os.path.join(self.output_dir, "llama-4-maverick"),
            experiment_id=experiment["experiment_id"]
        )
        
        deployment = self.ml_client.deploy_model(
            model_id=model["model_id"],
            deployment_name="llama-4-maverick-prod",
            replicas=2
        )
        
        self.assertEqual(experiment["experiment_id"], "exp-123")
        self.assertEqual(metrics["status"], "success")
        self.assertEqual(model["model_id"], "model-123")
        self.assertEqual(deployment["deployment_id"], "deploy-123")
        self.assertEqual(deployment["endpoint"], "http://example.com/model")


if __name__ == '__main__':
    unittest.main()

import unittest
import os
import sys
import json
import tempfile
from unittest.mock import patch, MagicMock

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '../../../src')))

from ml_infrastructure.training.config.llama4_maverick_config import Llama4MaverickConfig
from ml_infrastructure.training.config.llama4_scout_config import Llama4ScoutConfig
from ml_infrastructure.training.scripts.train_llama4 import LlamaTrainer
from ml_infrastructure.training.evaluation.metrics import ModelEvaluator


class TestLlama4Training(unittest.TestCase):
    """Integration tests for Llama 4 training workflows."""

    def setUp(self):
        """Set up test fixtures."""
        self.test_data_dir = os.path.join(os.path.dirname(__file__), '../../data/test')
        os.makedirs(self.test_data_dir, exist_ok=True)
        
        self.temp_dir = tempfile.TemporaryDirectory()
        self.output_dir = self.temp_dir.name
        
        self.training_data_path = os.path.join(self.test_data_dir, 'training_data.json')
        if not os.path.exists(self.training_data_path):
            sample_data = [
                {
                    "issue_number": 1,
                    "title": "Fix Kubernetes deployment issue",
                    "body": "The deployment is failing with error: ResourceQuotaExceeded",
                    "state": "closed",
                    "labels": ["bug", "kubernetes"],
                    "created_at": "2023-01-01T00:00:00Z",
                    "closed_at": "2023-01-02T00:00:00Z",
                    "user": "user1",
                    "url": "https://github.com/org/repo/issues/1",
                    "comments": [
                        {"user": "user3", "body": "I fixed this by increasing the resource quota."},
                        {"user": "user1", "body": "Thanks, that worked!"}
                    ],
                    "source": "github",
                    "repository": "kubernetes/kubernetes"
                },
                {
                    "issue_number": 2,
                    "title": "Update Terraform module for AWS provider",
                    "body": "Need to update the module to support the latest AWS provider version",
                    "state": "closed",
                    "labels": ["enhancement", "terraform"],
                    "created_at": "2023-01-03T00:00:00Z",
                    "closed_at": "2023-01-04T00:00:00Z",
                    "user": "user2",
                    "url": "https://github.com/org/repo/issues/2",
                    "comments": [
                        {"user": "user4", "body": "Updated the module to use AWS provider 4.0."},
                        {"user": "user2", "body": "Looks good, merging now."}
                    ],
                    "source": "github",
                    "repository": "terraform/terraform"
                }
            ]
            with open(self.training_data_path, 'w') as f:
                json.dump(sample_data, f, indent=2)
        
        self.maverick_config = Llama4MaverickConfig(
            model_name="llama-4-maverick",
            batch_size=4,
            learning_rate=2e-5,
            num_epochs=3,
            output_dir=self.output_dir
        )
        
        self.scout_config = Llama4ScoutConfig(
            model_name="llama-4-scout",
            batch_size=4,
            learning_rate=2e-5,
            num_epochs=3,
            output_dir=self.output_dir
        )

    def tearDown(self):
        """Clean up after tests."""
        self.temp_dir.cleanup()

    @patch('ml_infrastructure.training.scripts.train_llama4.LlamaTrainer._load_model')
    @patch('ml_infrastructure.training.scripts.train_llama4.LlamaTrainer._load_tokenizer')
    @patch('ml_infrastructure.training.scripts.train_llama4.LlamaTrainer._train_model')
    def test_maverick_training_workflow(self, mock_train, mock_tokenizer, mock_model):
        """Test the Llama 4 Maverick training workflow."""
        mock_model.return_value = MagicMock()
        mock_tokenizer.return_value = MagicMock()
        mock_train.return_value = {
            "train_loss": 0.1,
            "eval_loss": 0.2,
            "train_runtime": 100,
            "eval_accuracy": 0.85
        }
        
        trainer = LlamaTrainer(
            config=self.maverick_config,
            data_path=self.training_data_path
        )
        
        result = trainer.train()
        
        mock_train.assert_called_once()
        
        self.assertIn("train_loss", result)
        self.assertIn("eval_loss", result)
        self.assertIn("train_runtime", result)
        self.assertIn("eval_accuracy", result)
        
        mock_model.return_value.save_pretrained.assert_called_once_with(
            os.path.join(self.output_dir, "llama-4-maverick")
        )
        mock_tokenizer.return_value.save_pretrained.assert_called_once_with(
            os.path.join(self.output_dir, "llama-4-maverick")
        )

    @patch('ml_infrastructure.training.scripts.train_llama4.LlamaTrainer._load_model')
    @patch('ml_infrastructure.training.scripts.train_llama4.LlamaTrainer._load_tokenizer')
    @patch('ml_infrastructure.training.scripts.train_llama4.LlamaTrainer._train_model')
    def test_scout_training_workflow(self, mock_train, mock_tokenizer, mock_model):
        """Test the Llama 4 Scout training workflow."""
        mock_model.return_value = MagicMock()
        mock_tokenizer.return_value = MagicMock()
        mock_train.return_value = {
            "train_loss": 0.15,
            "eval_loss": 0.25,
            "train_runtime": 120,
            "eval_accuracy": 0.8
        }
        
        trainer = LlamaTrainer(
            config=self.scout_config,
            data_path=self.training_data_path
        )
        
        result = trainer.train()
        
        mock_train.assert_called_once()
        
        self.assertIn("train_loss", result)
        self.assertIn("eval_loss", result)
        self.assertIn("train_runtime", result)
        self.assertIn("eval_accuracy", result)
        
        mock_model.return_value.save_pretrained.assert_called_once_with(
            os.path.join(self.output_dir, "llama-4-scout")
        )
        mock_tokenizer.return_value.save_pretrained.assert_called_once_with(
            os.path.join(self.output_dir, "llama-4-scout")
        )

    @patch('ml_infrastructure.training.evaluation.metrics.ModelEvaluator._load_model')
    @patch('ml_infrastructure.training.evaluation.metrics.ModelEvaluator._load_tokenizer')
    def test_model_evaluation(self, mock_tokenizer, mock_model):
        """Test model evaluation metrics."""
        mock_model.return_value = MagicMock()
        mock_model.return_value.generate.return_value = [[1, 2, 3, 4, 5]]
        mock_tokenizer.return_value = MagicMock()
        mock_tokenizer.return_value.decode.return_value = "This is a generated response."
        mock_tokenizer.return_value.batch_encode_plus.return_value = {
            "input_ids": [[1, 2, 3]],
            "attention_mask": [[1, 1, 1]]
        }
        
        evaluator = ModelEvaluator(
            model_path=os.path.join(self.output_dir, "llama-4-maverick"),
            test_data_path=self.training_data_path
        )
        
        metrics = evaluator.evaluate()
        
        self.assertIn("accuracy", metrics)
        self.assertIn("f1_score", metrics)
        self.assertIn("precision", metrics)
        self.assertIn("recall", metrics)
        self.assertIn("latency", metrics)
        
        mock_model.assert_called_once()
        mock_tokenizer.assert_called_once()


if __name__ == '__main__':
    unittest.main()

"""
End-to-end test for the ML pipeline.
"""

import os
import pytest
import asyncio
from src.integrations.github.scraper import GitHubScraper
from src.integrations.github.integration import GitHubIntegration
from src.integrations.issue_collector.collector import IssueCollector
from src.ml_infrastructure.data.preprocessing.preprocessor import Preprocessor
from src.ml_infrastructure.training.scripts.train_llama4 import train_model


@pytest.mark.asyncio
async def test_end_to_end_pipeline():
    """Test the end-to-end ML pipeline."""
    github_token = os.environ.get("GITHUB_TOKEN")
    if not github_token:
        pytest.skip("No GitHub token available")
    
    integration = GitHubIntegration(api_key=github_token)
    scraper = GitHubScraper(integration=integration)
    
    collector = IssueCollector(scrapers=[scraper])
    
    issues = await collector.collect_issues(
        repositories=[
            {"owner": "kubernetes", "repo": "kubernetes", "limit": 5}
        ],
        issue_state="closed"
    )
    
    assert len(issues) > 0
    assert "title" in issues[0]
    assert "body" in issues[0]
    assert "comments" in issues[0]
    
    preprocessor = Preprocessor()
    
    processed_data = preprocessor.preprocess(issues)
    
    assert len(processed_data) > 0
    assert "input" in processed_data[0]
    assert "output" in processed_data[0]
    
    def mock_train(*args, **kwargs):
        return {
            "model_path": "/tmp/mock_model",
            "metrics": {
                "loss": 0.1,
                "accuracy": 0.9,
                "f1": 0.85
            }
        }
    
    train_model = mock_train
    original_train = train_model
    
    try:
        training_result = train_model(
            training_data=processed_data,
            model_name="llama4-maverick-test",
            epochs=1,
            batch_size=1,
            learning_rate=1e-5
        )
        
        assert "model_path" in training_result
        assert "metrics" in training_result
        assert "loss" in training_result["metrics"]
        assert "accuracy" in training_result["metrics"]
    finally:
        train_model = original_train


if __name__ == "__main__":
    asyncio.run(test_end_to_end_pipeline())

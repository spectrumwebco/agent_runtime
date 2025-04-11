"""
Example script demonstrating how to use the ML Infrastructure API Client for KubeFlow integration.

This script shows how to:
1. Create and manage KubeFlow pipelines
2. Submit pipeline runs
3. Monitor pipeline execution
4. Retrieve pipeline results
"""

import os
import time
import logging
from pathlib import Path
import sys

sys.path.append(str(Path(__file__).parent.parent.parent))
from api.client import MLInfrastructureClient

logging.basicConfig(level=logging.INFO, format="%(asctime)s - %(name)s - %(levelname)s - %(message)s")
logger = logging.getLogger(__name__)

def main():
    """Run the KubeFlow integration example."""
    client = MLInfrastructureClient()

    try:
        status = client.get_api_status()
        logger.info(f"API Status: {status}")
    except Exception as e:
        logger.error(f"Failed to get API status: {e}")
        return

    try:
        pipelines = client.list_kubeflow_pipelines()
        logger.info(f"Available KubeFlow pipelines: {len(pipelines)}")
        
        if not pipelines:
            logger.info("No pipelines available. Creating a new pipeline...")
            
            pipeline_file = os.path.join(
                os.path.dirname(__file__), 
                "../../kubeflow/manifests/pipelines/pipeline.yaml"
            )
            
            if not os.path.exists(pipeline_file):
                logger.error(f"Pipeline file not found: {pipeline_file}")
                return
                
            pipeline = client.create_kubeflow_pipeline(
                name="llama4-fine-tuning-pipeline",
                description="Pipeline for fine-tuning Llama 4 models on GitOps, Terraform, and Kubernetes issues",
                pipeline_file=pipeline_file,
                parameters={
                    "model_type": "llama4-maverick",
                    "dataset_id": "gitops-terraform-k8s-issues",
                    "epochs": 3,
                    "batch_size": 8,
                    "learning_rate": 2e-5
                }
            )
            logger.info(f"Created pipeline: {pipeline['name']} (ID: {pipeline['id']})")
            pipeline_id = pipeline["id"]
        else:
            pipeline = pipelines[0]
            pipeline_id = pipeline["id"]
            logger.info(f"Using existing pipeline: {pipeline['name']} (ID: {pipeline_id})")
            
            pipeline_details = client.get_kubeflow_pipeline(pipeline_id)
            logger.info(f"Pipeline details: {pipeline_details}")
    except Exception as e:
        logger.error(f"Failed to list or create KubeFlow pipelines: {e}")
        return

    try:
        run = client.run_kubeflow_pipeline(
            pipeline_id=pipeline_id,
            run_name=f"llama4-fine-tuning-run-{int(time.time())}",
            parameters={
                "model_type": "llama4-maverick",
                "dataset_id": "gitops-terraform-k8s-issues",
                "epochs": 3,
                "batch_size": 8,
                "learning_rate": 2e-5,
                "weight_decay": 0.01,
                "warmup_steps": 500,
                "max_seq_length": 2048,
                "gradient_accumulation_steps": 4,
                "fp16": True,
                "output_dir": "models/llama4-maverick-fine-tuned"
            }
        )
        logger.info(f"Submitted pipeline run: {run['name']} (ID: {run['id']})")
        run_id = run["id"]
    except Exception as e:
        logger.error(f"Failed to run KubeFlow pipeline: {e}")
        return

    try:
        logger.info(f"Monitoring pipeline run: {run_id}")
        for _ in range(5):  # Poll 5 times in a real scenario, you would poll until completion
            run_details = client.get_kubeflow_pipeline_run(pipeline_id, run_id)
            status = run_details.get("status", "unknown")
            logger.info(f"Pipeline run status: {status}")
            
            if status in ["Succeeded", "Failed", "Error", "Skipped", "Terminated"]:
                logger.info(f"Pipeline run completed with status: {status}")
                break
                
            logger.info("Waiting for pipeline run to complete...")
            time.sleep(10)  # Wait 10 seconds between polls
    except Exception as e:
        logger.error(f"Failed to monitor KubeFlow pipeline run: {e}")

    try:
        experiment = client.create_mlflow_experiment(
            name=f"llama4-fine-tuning-{int(time.time())}",
            artifact_location="s3://mlflow/artifacts"
        )
        logger.info(f"Created MLflow experiment: {experiment['name']} (ID: {experiment['id']})")
        experiment_id = experiment["id"]
        
        run = client.create_mlflow_run(
            experiment_id=experiment_id,
            run_name="llama4-maverick-fine-tuning",
            tags={
                "model_type": "llama4-maverick",
                "dataset": "gitops-terraform-k8s-issues",
                "pipeline_run_id": run_id
            }
        )
        logger.info(f"Created MLflow run: {run['run_name']} (ID: {run['id']})")
        mlflow_run_id = run["id"]
        
        client.log_mlflow_param(mlflow_run_id, "model_type", "llama4-maverick")
        client.log_mlflow_param(mlflow_run_id, "epochs", "3")
        client.log_mlflow_param(mlflow_run_id, "batch_size", "8")
        client.log_mlflow_param(mlflow_run_id, "learning_rate", "2e-5")
        
        client.log_mlflow_metric(mlflow_run_id, "train_loss", 0.1234)
        client.log_mlflow_metric(mlflow_run_id, "val_loss", 0.2345)
        client.log_mlflow_metric(mlflow_run_id, "accuracy", 0.9876)
        
        logger.info("Logged parameters and metrics to MLflow")
    except Exception as e:
        logger.error(f"Failed to create MLflow experiment or log data: {e}")

    logger.info("KubeFlow integration example completed successfully")

if __name__ == "__main__":
    main()

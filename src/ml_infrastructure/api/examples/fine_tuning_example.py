"""
Example script demonstrating how to use the ML Infrastructure API Client for fine-tuning.

This script shows how to:
1. Create a dataset from GitHub and Gitee issues
2. Submit a fine-tuning job for Llama 4 models
3. Monitor training progress
4. Evaluate the fine-tuned model
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
    """Run the fine-tuning example."""
    client = MLInfrastructureClient()

    try:
        status = client.get_api_status()
        logger.info(f"API Status: {status}")
    except Exception as e:
        logger.error(f"Failed to get API status: {e}")
        return

    try:
        dataset = client.create_dataset(
            name="gitops-terraform-k8s-issues",
            description="Solved issues from GitOps, Terraform, and Kubernetes repositories",
            source="github-gitee",
            version="1.0.0",
            metadata={
                "topics": ["gitops", "terraform", "kubernetes"],
                "issue_type": "solved",
                "repositories": [
                    "kubernetes/kubernetes",
                    "hashicorp/terraform",
                    "fluxcd/flux2",
                    "argoproj/argo-cd"
                ]
            }
        )
        logger.info(f"Created dataset: {dataset['name']} (ID: {dataset['id']})")
        dataset_id = dataset["id"]
    except Exception as e:
        logger.error(f"Failed to create dataset: {e}")
        return

    dataset_file = os.path.join(os.path.dirname(__file__), "../../data/issues_dataset.json")
    if os.path.exists(dataset_file):
        try:
            upload_result = client.upload_dataset_file(dataset_id, dataset_file)
            logger.info(f"Uploaded dataset file: {upload_result}")
        except Exception as e:
            logger.error(f"Failed to upload dataset file: {e}")

    try:
        experiment = client.create_experiment(
            name="llama4-fine-tuning-gitops-terraform-k8s",
            description="Fine-tuning Llama 4 models on GitOps, Terraform, and Kubernetes issues",
            tags=["llama4", "fine-tuning", "gitops", "terraform", "kubernetes"]
        )
        logger.info(f"Created experiment: {experiment['name']} (ID: {experiment['id']})")
        experiment_id = experiment["id"]
    except Exception as e:
        logger.error(f"Failed to create experiment: {e}")
        return

    try:
        maverick_job = client.submit_training_job(
            name="llama4-maverick-fine-tuning",
            model_type="llama4-maverick",
            dataset_id=dataset_id,
            config={
                "training_type": "fine-tuning",
                "epochs": 3,
                "batch_size": 8,
                "learning_rate": 2e-5,
                "weight_decay": 0.01,
                "warmup_steps": 500,
                "max_seq_length": 2048,
                "gradient_accumulation_steps": 4,
                "fp16": True,
                "output_dir": "models/llama4-maverick-fine-tuned"
            },
            description="Fine-tuning Llama 4 Maverick on GitOps, Terraform, and Kubernetes issues",
            hyperparameters={
                "learning_rate": 2e-5,
                "epochs": 3,
                "batch_size": 8
            }
        )
        logger.info(f"Submitted Llama 4 Maverick fine-tuning job: {maverick_job['name']} (ID: {maverick_job['id']})")
        maverick_job_id = maverick_job["id"]
    except Exception as e:
        logger.error(f"Failed to submit Llama 4 Maverick fine-tuning job: {e}")
        return

    try:
        scout_job = client.submit_training_job(
            name="llama4-scout-fine-tuning",
            model_type="llama4-scout",
            dataset_id=dataset_id,
            config={
                "training_type": "fine-tuning",
                "epochs": 3,
                "batch_size": 16,
                "learning_rate": 3e-5,
                "weight_decay": 0.01,
                "warmup_steps": 500,
                "max_seq_length": 2048,
                "gradient_accumulation_steps": 2,
                "fp16": True,
                "output_dir": "models/llama4-scout-fine-tuned"
            },
            description="Fine-tuning Llama 4 Scout on GitOps, Terraform, and Kubernetes issues",
            hyperparameters={
                "learning_rate": 3e-5,
                "epochs": 3,
                "batch_size": 16
            }
        )
        logger.info(f"Submitted Llama 4 Scout fine-tuning job: {scout_job['name']} (ID: {scout_job['id']})")
        scout_job_id = scout_job["id"]
    except Exception as e:
        logger.error(f"Failed to submit Llama 4 Scout fine-tuning job: {e}")
        return

    logger.info("Monitoring Llama 4 Maverick training progress...")
    monitor_training_job(client, maverick_job_id)

    logger.info("Monitoring Llama 4 Scout training progress...")
    monitor_training_job(client, scout_job_id)

    try:
        maverick_eval = client.create_evaluation(
            name="llama4-maverick-evaluation",
            model_id=f"llama4-maverick-fine-tuned-{maverick_job_id}",
            dataset_id=dataset_id,
            metrics=["accuracy", "f1", "precision", "recall", "trajectory_similarity"],
            description="Evaluation of fine-tuned Llama 4 Maverick model",
            config={
                "split": "test",
                "num_samples": 100,
                "trajectory_metrics": True
            }
        )
        logger.info(f"Created evaluation for Llama 4 Maverick: {maverick_eval['name']} (ID: {maverick_eval['id']})")
        maverick_eval_id = maverick_eval["id"]
    except Exception as e:
        logger.error(f"Failed to create evaluation for Llama 4 Maverick: {e}")
        return

    try:
        scout_eval = client.create_evaluation(
            name="llama4-scout-evaluation",
            model_id=f"llama4-scout-fine-tuned-{scout_job_id}",
            dataset_id=dataset_id,
            metrics=["accuracy", "f1", "precision", "recall", "trajectory_similarity"],
            description="Evaluation of fine-tuned Llama 4 Scout model",
            config={
                "split": "test",
                "num_samples": 100,
                "trajectory_metrics": True
            }
        )
        logger.info(f"Created evaluation for Llama 4 Scout: {scout_eval['name']} (ID: {scout_eval['id']})")
        scout_eval_id = scout_eval["id"]
    except Exception as e:
        logger.error(f"Failed to create evaluation for Llama 4 Scout: {e}")
        return

    logger.info("Waiting for evaluations to complete...")
    time.sleep(30)  # In a real scenario, you would poll until completion

    try:
        maverick_results = client.get_evaluation_results(maverick_eval_id)
        logger.info(f"Llama 4 Maverick evaluation results: {maverick_results}")
    except Exception as e:
        logger.error(f"Failed to get evaluation results for Llama 4 Maverick: {e}")

    try:
        scout_results = client.get_evaluation_results(scout_eval_id)
        logger.info(f"Llama 4 Scout evaluation results: {scout_results}")
    except Exception as e:
        logger.error(f"Failed to get evaluation results for Llama 4 Scout: {e}")

    try:
        comparison = client.compare_evaluations([maverick_eval_id, scout_eval_id])
        logger.info(f"Evaluation comparison: {comparison}")
    except Exception as e:
        logger.error(f"Failed to compare evaluations: {e}")

    try:
        maverick_deployment = client.create_deployment(
            name="llama4-maverick-deployment",
            model_id=f"llama4-maverick-fine-tuned-{maverick_job_id}",
            replicas=2,
            description="Deployment of fine-tuned Llama 4 Maverick model",
            config={
                "resources": {
                    "limits": {
                        "cpu": "4",
                        "memory": "16Gi",
                        "nvidia.com/gpu": "1"
                    },
                    "requests": {
                        "cpu": "2",
                        "memory": "8Gi"
                    }
                },
                "scaling": {
                    "min_replicas": 1,
                    "max_replicas": 5,
                    "target_cpu_utilization": 80
                }
            }
        )
        logger.info(f"Created deployment for Llama 4 Maverick: {maverick_deployment['name']} (ID: {maverick_deployment['id']})")
    except Exception as e:
        logger.error(f"Failed to create deployment for Llama 4 Maverick: {e}")

    try:
        scout_deployment = client.create_deployment(
            name="llama4-scout-deployment",
            model_id=f"llama4-scout-fine-tuned-{scout_job_id}",
            replicas=2,
            description="Deployment of fine-tuned Llama 4 Scout model",
            config={
                "resources": {
                    "limits": {
                        "cpu": "4",
                        "memory": "16Gi",
                        "nvidia.com/gpu": "1"
                    },
                    "requests": {
                        "cpu": "2",
                        "memory": "8Gi"
                    }
                },
                "scaling": {
                    "min_replicas": 1,
                    "max_replicas": 5,
                    "target_cpu_utilization": 80
                }
            }
        )
        logger.info(f"Created deployment for Llama 4 Scout: {scout_deployment['name']} (ID: {scout_deployment['id']})")
    except Exception as e:
        logger.error(f"Failed to create deployment for Llama 4 Scout: {e}")

    logger.info("Fine-tuning example completed successfully")

def monitor_training_job(client, job_id, poll_interval=10, max_polls=12):
    """Monitor a training job until completion or max polls reached."""
    polls = 0
    while polls < max_polls:
        try:
            job = client.get_training_job(job_id)
            status = job.get("status", "unknown")
            logger.info(f"Job {job_id} status: {status}")
            
            metrics = client.get_training_job_metrics(job_id)
            if metrics:
                logger.info(f"Job {job_id} metrics: {metrics}")
            
            if status in ["completed", "failed", "stopped"]:
                logger.info(f"Job {job_id} {status}")
                
                logs = client.get_training_job_logs(job_id)
                if logs:
                    logger.info(f"Job {job_id} logs (last 5 lines): {logs[-5:]}")
                
                return status
            
            polls += 1
            time.sleep(poll_interval)
        except Exception as e:
            logger.error(f"Error monitoring job {job_id}: {e}")
            polls += 1
            time.sleep(poll_interval)
    
    logger.warning(f"Max polls reached for job {job_id}")
    return "unknown"

if __name__ == "__main__":
    main()

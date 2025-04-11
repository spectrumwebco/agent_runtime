"""
Example script demonstrating how to use the ML Infrastructure API Client for model inference.

This script shows how to:
1. Connect to a deployed model
2. Make predictions using the model
3. Perform batch predictions
4. Handle prediction errors
"""

import os
import json
import logging
from pathlib import Path
import sys

sys.path.append(str(Path(__file__).parent.parent.parent))
from api.client import MLInfrastructureClient

logging.basicConfig(level=logging.INFO, format="%(asctime)s - %(name)s - %(levelname)s - %(message)s")
logger = logging.getLogger(__name__)

def main():
    """Run the inference example."""
    client = MLInfrastructureClient()

    try:
        status = client.get_api_status()
        logger.info(f"API Status: {status}")
    except Exception as e:
        logger.error(f"Failed to get API status: {e}")
        return

    try:
        deployments = client.list_deployments()
        logger.info(f"Available deployments: {len(deployments)}")
        
        if not deployments:
            logger.warning("No deployments available. Please deploy a model first.")
            return
            
        deployment = deployments[0]
        deployment_id = deployment["id"]
        logger.info(f"Using deployment: {deployment['name']} (ID: {deployment_id})")
    except Exception as e:
        logger.error(f"Failed to list deployments: {e}")
        return

    try:
        status = client.get_deployment_status(deployment_id)
        logger.info(f"Deployment status: {status}")
        
        if status.get("status") != "running":
            logger.warning(f"Deployment is not running. Current status: {status.get('status')}")
            return
    except Exception as e:
        logger.error(f"Failed to get deployment status: {e}")
        return

    gitops_issue = {
        "repository": "fluxcd/flux2",
        "issue_title": "Helm chart fails to deploy with custom values",
        "issue_description": """
        I'm trying to deploy a Helm chart using Flux CD, but it fails when I provide custom values.
        
        Steps to reproduce:
        1. Create a HelmRelease with custom values
        2. Apply the HelmRelease to the cluster
        3. Observe the error in the logs
        
        Error message:
        ```
        failed to build helm release: failed to render chart: failed to parse values
        ```
        
        My HelmRelease looks like:
        ```yaml
        apiVersion: helm.toolkit.fluxcd.io/v2beta1
        kind: HelmRelease
        metadata:
          name: my-app
          namespace: default
        spec:
          interval: 5m
          chart:
            spec:
              chart: my-chart
              version: 1.0.0
              sourceRef:
                kind: HelmRepository
                name: my-repo
          values:
            custom:
              value: "test"
        ```
        """
    }

    try:
        logger.info("Making a single prediction...")
        prediction = client.predict(
            deployment_id=deployment_id,
            inputs=gitops_issue,
            parameters={
                "max_length": 1024,
                "temperature": 0.7,
                "top_p": 0.9
            }
        )
        logger.info("Prediction result:")
        logger.info(json.dumps(prediction, indent=2))
    except Exception as e:
        logger.error(f"Failed to make prediction: {e}")

    batch_inputs = [
        {
            "repository": "kubernetes/kubernetes",
            "issue_title": "Pod fails to start with CrashLoopBackOff",
            "issue_description": "My pod keeps restarting with CrashLoopBackOff status. The container exits with code 1."
        },
        {
            "repository": "hashicorp/terraform",
            "issue_title": "terraform apply fails with provider error",
            "issue_description": "When running terraform apply, I get an error from the AWS provider about invalid credentials."
        },
        {
            "repository": "argoproj/argo-cd",
            "issue_title": "Application sync fails with permission error",
            "issue_description": "My Argo CD application fails to sync with a permission error when trying to create resources."
        }
    ]

    try:
        logger.info("Making batch predictions...")
        batch_predictions = client.batch_predict(
            deployment_id=deployment_id,
            inputs=batch_inputs,
            parameters={
                "max_length": 512,
                "temperature": 0.5,
                "top_p": 0.95
            }
        )
        logger.info(f"Received {len(batch_predictions.get('predictions', []))} batch predictions")
        
        if batch_predictions.get("predictions"):
            logger.info("First batch prediction result:")
            logger.info(json.dumps(batch_predictions["predictions"][0], indent=2))
    except Exception as e:
        logger.error(f"Failed to make batch predictions: {e}")

    try:
        logger.info("Testing error handling with invalid input...")
        invalid_prediction = client.predict(
            deployment_id=deployment_id,
            inputs={"invalid": "input format"},
            parameters={"temperature": 0.5}
        )
        logger.info("Prediction with invalid input result:")
        logger.info(json.dumps(invalid_prediction, indent=2))
    except Exception as e:
        logger.info(f"Expected error occurred with invalid input: {e}")

    logger.info("Inference example completed successfully")

if __name__ == "__main__":
    main()

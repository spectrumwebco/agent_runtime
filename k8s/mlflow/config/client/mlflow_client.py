"""
MLFlow client for interacting with the MLFlow server.
"""

import os
import json
import logging
from typing import Dict, List, Any, Optional, Union
import mlflow
from mlflow.tracking import MlflowClient


class MLFlowClientConfig:
    """
    Configuration for MLFlow client.
    """

    def __init__(
        self,
        tracking_uri: str = "http://mlflow-server.mlflow.svc.cluster.local:5000",
        experiment_name: str = "llama4-maverick-fine-tuning",
        artifact_location: Optional[str] = None,
    ):
        """
        Initialize MLFlow client configuration.

        Args:
            tracking_uri: MLFlow tracking server URI
            experiment_name: Name of the experiment
            artifact_location: Location for artifacts
        """
        self.tracking_uri = tracking_uri
        self.experiment_name = experiment_name
        self.artifact_location = artifact_location


class MLFlowClient:
    """
    Client for interacting with MLFlow.
    """

    def __init__(self, config: MLFlowClientConfig):
        """
        Initialize MLFlow client.

        Args:
            config: MLFlow client configuration
        """
        self.config = config
        self.client = None
        self.experiment_id = None
        self.run_id = None

    def initialize(self) -> None:
        """
        Initialize MLFlow client and set up experiment.
        """
        mlflow.set_tracking_uri(self.config.tracking_uri)
        self.client = MlflowClient()

        experiment = mlflow.get_experiment_by_name(self.config.experiment_name)
        if experiment is None:
            logging.info(f"Creating experiment: {self.config.experiment_name}")
            self.experiment_id = mlflow.create_experiment(
                name=self.config.experiment_name,
                artifact_location=self.config.artifact_location,
            )
        else:
            self.experiment_id = experiment.experiment_id
            logging.info(f"Using existing experiment: {self.config.experiment_name} (ID: {self.experiment_id})")

    def start_run(
        self,
        run_name: Optional[str] = None,
        tags: Optional[Dict[str, str]] = None,
    ) -> str:
        """
        Start a new MLFlow run.

        Args:
            run_name: Name of the run
            tags: Tags for the run

        Returns:
            Run ID
        """
        if self.experiment_id is None:
            self.initialize()

        active_run = mlflow.start_run(
            experiment_id=self.experiment_id,
            run_name=run_name,
            tags=tags,
        )
        self.run_id = active_run.info.run_id
        logging.info(f"Started MLFlow run: {self.run_id}")
        return self.run_id

    def log_params(self, params: Dict[str, Any]) -> None:
        """
        Log parameters to the current run.

        Args:
            params: Parameters to log
        """
        for key, value in params.items():
            mlflow.log_param(key, value)

    def log_metrics(self, metrics: Dict[str, float], step: Optional[int] = None) -> None:
        """
        Log metrics to the current run.

        Args:
            metrics: Metrics to log
            step: Step number
        """
        for key, value in metrics.items():
            mlflow.log_metric(key, value, step=step)

    def log_artifact(self, local_path: str, artifact_path: Optional[str] = None) -> None:
        """
        Log an artifact to the current run.

        Args:
            local_path: Local path to the artifact
            artifact_path: Path within the artifact directory
        """
        mlflow.log_artifact(local_path, artifact_path)

    def log_model(
        self,
        model_path: str,
        model_name: str,
        flavor: str = "transformers",
        **kwargs,
    ) -> None:
        """
        Log a model to the current run.

        Args:
            model_path: Path to the model
            model_name: Name of the model
            flavor: Model flavor (e.g., "transformers", "pytorch")
            **kwargs: Additional arguments for the specific flavor
        """
        if flavor == "transformers":
            mlflow.transformers.log_model(
                transformers_model=model_path,
                artifact_path=model_name,
                **kwargs,
            )
        elif flavor == "pytorch":
            mlflow.pytorch.log_model(
                pytorch_model=model_path,
                artifact_path=model_name,
                **kwargs,
            )
        else:
            raise ValueError(f"Unsupported model flavor: {flavor}")

    def end_run(self, status: str = "FINISHED") -> None:
        """
        End the current run.

        Args:
            status: Run status
        """
        mlflow.end_run(status=status)
        logging.info(f"Ended MLFlow run: {self.run_id} with status: {status}")
        self.run_id = None

    def get_run(self, run_id: Optional[str] = None) -> Dict[str, Any]:
        """
        Get run information.

        Args:
            run_id: Run ID (uses current run if None)

        Returns:
            Run information
        """
        run_id = run_id or self.run_id
        if run_id is None:
            raise ValueError("No run ID specified and no current run")

        run = self.client.get_run(run_id)
        return {
            "run_id": run.info.run_id,
            "experiment_id": run.info.experiment_id,
            "status": run.info.status,
            "start_time": run.info.start_time,
            "end_time": run.info.end_time,
            "artifact_uri": run.info.artifact_uri,
            "metrics": run.data.metrics,
            "params": run.data.params,
            "tags": run.data.tags,
        }

    def search_runs(
        self,
        experiment_ids: Optional[List[str]] = None,
        filter_string: Optional[str] = None,
        order_by: Optional[List[str]] = None,
        max_results: int = 100,
    ) -> List[Dict[str, Any]]:
        """
        Search for runs.

        Args:
            experiment_ids: List of experiment IDs
            filter_string: Filter string
            order_by: Order by columns
            max_results: Maximum number of results

        Returns:
            List of runs
        """
        experiment_ids = experiment_ids or [self.experiment_id]
        if experiment_ids[0] is None:
            self.initialize()
            experiment_ids = [self.experiment_id]

        runs = self.client.search_runs(
            experiment_ids=experiment_ids,
            filter_string=filter_string,
            order_by=order_by,
            max_results=max_results,
        )

        return [
            {
                "run_id": run.info.run_id,
                "experiment_id": run.info.experiment_id,
                "status": run.info.status,
                "start_time": run.info.start_time,
                "end_time": run.info.end_time,
                "artifact_uri": run.info.artifact_uri,
                "metrics": run.data.metrics,
                "params": run.data.params,
                "tags": run.data.tags,
            }
            for run in runs
        ]

    def get_best_run(
        self,
        experiment_ids: Optional[List[str]] = None,
        metric_name: str = "validation_loss",
        ascending: bool = False,
    ) -> Dict[str, Any]:
        """
        Get the best run based on a metric.

        Args:
            experiment_ids: List of experiment IDs
            metric_name: Metric name to sort by
            ascending: Whether to sort in ascending order

        Returns:
            Best run
        """
        order_by = [f"metrics.{metric_name} {'ASC' if ascending else 'DESC'}"]
        runs = self.search_runs(
            experiment_ids=experiment_ids,
            order_by=order_by,
            max_results=1,
        )
        if not runs:
            raise ValueError(f"No runs found for metric: {metric_name}")
        return runs[0]

    def register_model(
        self,
        run_id: str,
        model_name: str,
        model_version: Optional[str] = None,
        description: Optional[str] = None,
        tags: Optional[Dict[str, str]] = None,
    ) -> Dict[str, Any]:
        """
        Register a model in the MLFlow model registry.

        Args:
            run_id: Run ID
            model_name: Model name
            model_version: Model version
            description: Model description
            tags: Model tags

        Returns:
            Registered model information
        """
        run = self.client.get_run(run_id)
        model_uri = f"runs:/{run_id}/artifacts/{model_name}"

        model_details = mlflow.register_model(
            model_uri=model_uri,
            name=model_name,
        )

        if description:
            self.client.update_registered_model(
                name=model_name,
                description=description,
            )

        if tags:
            for key, value in tags.items():
                self.client.set_registered_model_tag(
                    name=model_name,
                    key=key,
                    value=value,
                )

        return {
            "name": model_details.name,
            "version": model_details.version,
            "creation_timestamp": model_details.creation_timestamp,
            "status": "REGISTERED",
        }

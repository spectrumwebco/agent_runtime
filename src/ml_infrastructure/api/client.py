"""
API client for ML infrastructure.
"""

import os
import json
import logging
import requests
from typing import Dict, List, Any, Optional, Union
import asyncio
import aiohttp


class MLInfrastructureAPIClient:
    """
    API client for ML infrastructure.
    """

    def __init__(
        self,
        base_url: str = "http://ml-infrastructure-api.example.com",
        api_key: Optional[str] = None,
    ):
        """
        Initialize ML infrastructure API client.

        Args:
            base_url: Base URL for the API
            api_key: API key for authentication
        """
        self.base_url = base_url
        self.api_key = api_key
        self.headers = {
            "Content-Type": "application/json",
        }
        
        if api_key:
            self.headers["Authorization"] = f"Bearer {api_key}"
        
        logging.basicConfig(
            level=logging.INFO,
            format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
        )
        self.logger = logging.getLogger("MLInfrastructureAPIClient")

    async def _make_async_request(
        self,
        method: str,
        endpoint: str,
        params: Optional[Dict[str, Any]] = None,
        data: Optional[Dict[str, Any]] = None,
        files: Optional[Dict[str, Any]] = None,
    ) -> Dict[str, Any]:
        """
        Make an asynchronous request to the API.

        Args:
            method: HTTP method
            endpoint: API endpoint
            params: Query parameters
            data: Request data
            files: Files to upload

        Returns:
            Response data
        """
        url = f"{self.base_url}/{endpoint}"
        
        try:
            async with aiohttp.ClientSession() as session:
                if method == "GET":
                    async with session.get(url, headers=self.headers, params=params) as response:
                        response.raise_for_status()
                        return await response.json()
                elif method == "POST":
                    if files:
                        form_data = aiohttp.FormData()
                        for key, value in data.items():
                            form_data.add_field(key, str(value))
                        for key, file_obj in files.items():
                            form_data.add_field(key, file_obj)
                        
                        headers = {k: v for k, v in self.headers.items() if k != "Content-Type"}
                        
                        async with session.post(url, headers=headers, params=params, data=form_data) as response:
                            response.raise_for_status()
                            return await response.json()
                    else:
                        async with session.post(url, headers=self.headers, params=params, json=data) as response:
                            response.raise_for_status()
                            return await response.json()
                elif method == "PUT":
                    async with session.put(url, headers=self.headers, params=params, json=data) as response:
                        response.raise_for_status()
                        return await response.json()
                elif method == "DELETE":
                    async with session.delete(url, headers=self.headers, params=params, json=data) as response:
                        response.raise_for_status()
                        return await response.json()
                else:
                    raise ValueError(f"Unsupported HTTP method: {method}")
        
        except aiohttp.ClientError as e:
            self.logger.error(f"Error making request: {str(e)}")
            raise

    def _make_request(
        self,
        method: str,
        endpoint: str,
        params: Optional[Dict[str, Any]] = None,
        data: Optional[Dict[str, Any]] = None,
        files: Optional[Dict[str, Any]] = None,
    ) -> Dict[str, Any]:
        """
        Make a synchronous request to the API.

        Args:
            method: HTTP method
            endpoint: API endpoint
            params: Query parameters
            data: Request data
            files: Files to upload

        Returns:
            Response data
        """
        url = f"{self.base_url}/{endpoint}"
        
        try:
            if method == "GET":
                response = requests.get(url, headers=self.headers, params=params)
            elif method == "POST":
                if files:
                    headers = {k: v for k, v in self.headers.items() if k != "Content-Type"}
                    response = requests.post(url, headers=headers, params=params, data=data, files=files)
                else:
                    response = requests.post(url, headers=self.headers, params=params, json=data)
            elif method == "PUT":
                response = requests.put(url, headers=self.headers, params=params, json=data)
            elif method == "DELETE":
                response = requests.delete(url, headers=self.headers, params=params, json=data)
            else:
                raise ValueError(f"Unsupported HTTP method: {method}")
            
            response.raise_for_status()
            return response.json()
        
        except requests.exceptions.RequestException as e:
            self.logger.error(f"Error making request: {str(e)}")
            raise


    async def get_models(self, model_type: Optional[str] = None) -> List[Dict[str, Any]]:
        """
        Get available models.

        Args:
            model_type: Filter by model type

        Returns:
            List of models
        """
        params = {}
        if model_type:
            params["model_type"] = model_type
        
        return await self._make_async_request("GET", "models", params=params)

    async def get_model(self, model_id: str) -> Dict[str, Any]:
        """
        Get model details.

        Args:
            model_id: Model ID

        Returns:
            Model details
        """
        return await self._make_async_request("GET", f"models/{model_id}")


    async def create_fine_tuning_job(
        self,
        model_type: str,
        training_data_path: str,
        validation_data_path: Optional[str] = None,
        hyperparameters: Optional[Dict[str, Any]] = None,
    ) -> Dict[str, Any]:
        """
        Create a fine-tuning job.

        Args:
            model_type: Model type
            training_data_path: Path to training data
            validation_data_path: Path to validation data
            hyperparameters: Hyperparameters for fine-tuning

        Returns:
            Job details
        """
        data = {
            "model_type": model_type,
            "training_data_path": training_data_path,
        }
        
        if validation_data_path:
            data["validation_data_path"] = validation_data_path
        
        if hyperparameters:
            data["hyperparameters"] = hyperparameters
        
        return await self._make_async_request("POST", "fine-tuning/jobs", data=data)

    async def get_fine_tuning_job(self, job_id: str) -> Dict[str, Any]:
        """
        Get fine-tuning job details.

        Args:
            job_id: Job ID

        Returns:
            Job details
        """
        return await self._make_async_request("GET", f"fine-tuning/jobs/{job_id}")

    async def list_fine_tuning_jobs(
        self,
        model_type: Optional[str] = None,
        status: Optional[str] = None,
        limit: int = 10,
    ) -> List[Dict[str, Any]]:
        """
        List fine-tuning jobs.

        Args:
            model_type: Filter by model type
            status: Filter by status
            limit: Maximum number of jobs to return

        Returns:
            List of jobs
        """
        params = {"limit": limit}
        
        if model_type:
            params["model_type"] = model_type
        
        if status:
            params["status"] = status
        
        return await self._make_async_request("GET", "fine-tuning/jobs", params=params)

    async def cancel_fine_tuning_job(self, job_id: str) -> Dict[str, Any]:
        """
        Cancel a fine-tuning job.

        Args:
            job_id: Job ID

        Returns:
            Job details
        """
        return await self._make_async_request("POST", f"fine-tuning/jobs/{job_id}/cancel")


    async def upload_training_data(
        self,
        file_path: str,
        dataset_name: Optional[str] = None,
    ) -> Dict[str, Any]:
        """
        Upload training data.

        Args:
            file_path: Path to training data file
            dataset_name: Name of the dataset

        Returns:
            Upload details
        """
        with open(file_path, "rb") as f:
            files = {"file": f}
            data = {}
            
            if dataset_name:
                data["dataset_name"] = dataset_name
            
            return await self._make_async_request("POST", "datasets/upload", data=data, files=files)

    async def list_datasets(self, limit: int = 10) -> List[Dict[str, Any]]:
        """
        List available datasets.

        Args:
            limit: Maximum number of datasets to return

        Returns:
            List of datasets
        """
        params = {"limit": limit}
        return await self._make_async_request("GET", "datasets", params=params)

    async def get_dataset(self, dataset_id: str) -> Dict[str, Any]:
        """
        Get dataset details.

        Args:
            dataset_id: Dataset ID

        Returns:
            Dataset details
        """
        return await self._make_async_request("GET", f"datasets/{dataset_id}")

    async def delete_dataset(self, dataset_id: str) -> Dict[str, Any]:
        """
        Delete a dataset.

        Args:
            dataset_id: Dataset ID

        Returns:
            Deletion details
        """
        return await self._make_async_request("DELETE", f"datasets/{dataset_id}")


    async def create_inference_service(
        self,
        model_id: str,
        service_name: str,
        replicas: int = 1,
        resources: Optional[Dict[str, Any]] = None,
    ) -> Dict[str, Any]:
        """
        Create an inference service.

        Args:
            model_id: Model ID
            service_name: Service name
            replicas: Number of replicas
            resources: Resource requirements

        Returns:
            Service details
        """
        data = {
            "model_id": model_id,
            "service_name": service_name,
            "replicas": replicas,
        }
        
        if resources:
            data["resources"] = resources
        
        return await self._make_async_request("POST", "inference/services", data=data)

    async def get_inference_service(self, service_id: str) -> Dict[str, Any]:
        """
        Get inference service details.

        Args:
            service_id: Service ID

        Returns:
            Service details
        """
        return await self._make_async_request("GET", f"inference/services/{service_id}")

    async def list_inference_services(
        self,
        model_id: Optional[str] = None,
        limit: int = 10,
    ) -> List[Dict[str, Any]]:
        """
        List inference services.

        Args:
            model_id: Filter by model ID
            limit: Maximum number of services to return

        Returns:
            List of services
        """
        params = {"limit": limit}
        
        if model_id:
            params["model_id"] = model_id
        
        return await self._make_async_request("GET", "inference/services", params=params)

    async def delete_inference_service(self, service_id: str) -> Dict[str, Any]:
        """
        Delete an inference service.

        Args:
            service_id: Service ID

        Returns:
            Deletion details
        """
        return await self._make_async_request("DELETE", f"inference/services/{service_id}")

    async def predict(
        self,
        service_id: str,
        input_text: str,
        parameters: Optional[Dict[str, Any]] = None,
    ) -> Dict[str, Any]:
        """
        Make a prediction using an inference service.

        Args:
            service_id: Service ID
            input_text: Input text
            parameters: Prediction parameters

        Returns:
            Prediction result
        """
        data = {
            "input_text": input_text,
        }
        
        if parameters:
            data["parameters"] = parameters
        
        return await self._make_async_request("POST", f"inference/services/{service_id}/predict", data=data)


    async def get_experiments(self, limit: int = 10) -> List[Dict[str, Any]]:
        """
        Get MLFlow experiments.

        Args:
            limit: Maximum number of experiments to return

        Returns:
            List of experiments
        """
        params = {"limit": limit}
        return await self._make_async_request("GET", "mlflow/experiments", params=params)

    async def get_experiment(self, experiment_id: str) -> Dict[str, Any]:
        """
        Get MLFlow experiment details.

        Args:
            experiment_id: Experiment ID

        Returns:
            Experiment details
        """
        return await self._make_async_request("GET", f"mlflow/experiments/{experiment_id}")

    async def create_experiment(
        self,
        name: str,
        artifact_location: Optional[str] = None,
        tags: Optional[Dict[str, str]] = None,
    ) -> Dict[str, Any]:
        """
        Create MLFlow experiment.

        Args:
            name: Experiment name
            artifact_location: Artifact location
            tags: Experiment tags

        Returns:
            Experiment details
        """
        data = {
            "name": name,
        }
        
        if artifact_location:
            data["artifact_location"] = artifact_location
        
        if tags:
            data["tags"] = tags
        
        return await self._make_async_request("POST", "mlflow/experiments", data=data)

    async def get_runs(
        self,
        experiment_id: str,
        status: Optional[str] = None,
        limit: int = 10,
    ) -> List[Dict[str, Any]]:
        """
        Get MLFlow runs.

        Args:
            experiment_id: Experiment ID
            status: Filter by status
            limit: Maximum number of runs to return

        Returns:
            List of runs
        """
        params = {
            "experiment_id": experiment_id,
            "limit": limit,
        }
        
        if status:
            params["status"] = status
        
        return await self._make_async_request("GET", "mlflow/runs", params=params)

    async def get_run(self, run_id: str) -> Dict[str, Any]:
        """
        Get MLFlow run details.

        Args:
            run_id: Run ID

        Returns:
            Run details
        """
        return await self._make_async_request("GET", f"mlflow/runs/{run_id}")

    async def create_run(
        self,
        experiment_id: str,
        run_name: Optional[str] = None,
        tags: Optional[Dict[str, str]] = None,
    ) -> Dict[str, Any]:
        """
        Create MLFlow run.

        Args:
            experiment_id: Experiment ID
            run_name: Run name
            tags: Run tags

        Returns:
            Run details
        """
        data = {
            "experiment_id": experiment_id,
        }
        
        if run_name:
            data["run_name"] = run_name
        
        if tags:
            data["tags"] = tags
        
        return await self._make_async_request("POST", "mlflow/runs", data=data)

    async def log_metrics(
        self,
        run_id: str,
        metrics: Dict[str, float],
        step: Optional[int] = None,
    ) -> Dict[str, Any]:
        """
        Log metrics to MLFlow run.

        Args:
            run_id: Run ID
            metrics: Metrics to log
            step: Step number

        Returns:
            Log details
        """
        data = {
            "metrics": metrics,
        }
        
        if step is not None:
            data["step"] = step
        
        return await self._make_async_request("POST", f"mlflow/runs/{run_id}/metrics", data=data)

    async def log_params(
        self,
        run_id: str,
        params: Dict[str, str],
    ) -> Dict[str, Any]:
        """
        Log parameters to MLFlow run.

        Args:
            run_id: Run ID
            params: Parameters to log

        Returns:
            Log details
        """
        data = {
            "params": params,
        }
        
        return await self._make_async_request("POST", f"mlflow/runs/{run_id}/params", data=data)


    async def create_pipeline_run(
        self,
        pipeline_id: str,
        run_name: str,
        parameters: Optional[Dict[str, Any]] = None,
    ) -> Dict[str, Any]:
        """
        Create a KubeFlow pipeline run.

        Args:
            pipeline_id: Pipeline ID
            run_name: Run name
            parameters: Pipeline parameters

        Returns:
            Run details
        """
        data = {
            "pipeline_id": pipeline_id,
            "run_name": run_name,
        }
        
        if parameters:
            data["parameters"] = parameters
        
        return await self._make_async_request("POST", "kubeflow/pipelines/runs", data=data)

    async def get_pipeline_run(self, run_id: str) -> Dict[str, Any]:
        """
        Get KubeFlow pipeline run details.

        Args:
            run_id: Run ID

        Returns:
            Run details
        """
        return await self._make_async_request("GET", f"kubeflow/pipelines/runs/{run_id}")

    async def list_pipeline_runs(
        self,
        pipeline_id: Optional[str] = None,
        status: Optional[str] = None,
        limit: int = 10,
    ) -> List[Dict[str, Any]]:
        """
        List KubeFlow pipeline runs.

        Args:
            pipeline_id: Filter by pipeline ID
            status: Filter by status
            limit: Maximum number of runs to return

        Returns:
            List of runs
        """
        params = {"limit": limit}
        
        if pipeline_id:
            params["pipeline_id"] = pipeline_id
        
        if status:
            params["status"] = status
        
        return await self._make_async_request("GET", "kubeflow/pipelines/runs", params=params)


    async def create_kserve_model(
        self,
        name: str,
        model_uri: str,
        model_format: str = "pytorch",
        resources: Optional[Dict[str, Any]] = None,
        env: Optional[List[Dict[str, Any]]] = None,
    ) -> Dict[str, Any]:
        """
        Create a KServe model.

        Args:
            name: Model name
            model_uri: Model URI
            model_format: Model format
            resources: Resource requirements
            env: Environment variables

        Returns:
            Model details
        """
        data = {
            "name": name,
            "model_uri": model_uri,
            "model_format": model_format,
        }
        
        if resources:
            data["resources"] = resources
        
        if env:
            data["env"] = env
        
        return await self._make_async_request("POST", "kserve/models", data=data)

    async def get_kserve_model(self, name: str) -> Dict[str, Any]:
        """
        Get KServe model details.

        Args:
            name: Model name

        Returns:
            Model details
        """
        return await self._make_async_request("GET", f"kserve/models/{name}")

    async def list_kserve_models(self, limit: int = 10) -> List[Dict[str, Any]]:
        """
        List KServe models.

        Args:
            limit: Maximum number of models to return

        Returns:
            List of models
        """
        params = {"limit": limit}
        return await self._make_async_request("GET", "kserve/models", params=params)

    async def delete_kserve_model(self, name: str) -> Dict[str, Any]:
        """
        Delete a KServe model.

        Args:
            name: Model name

        Returns:
            Deletion details
        """
        return await self._make_async_request("DELETE", f"kserve/models/{name}")


    async def get_feature_values(
        self,
        entity_name: str,
        entity_ids: List[str],
        feature_names: List[str],
    ) -> Dict[str, Any]:
        """
        Get feature values from Feast.

        Args:
            entity_name: Entity name
            entity_ids: Entity IDs
            feature_names: Feature names

        Returns:
            Feature values
        """
        data = {
            "entity_name": entity_name,
            "entity_ids": entity_ids,
            "feature_names": feature_names,
        }
        
        return await self._make_async_request("POST", "feast/features", data=data)


    def run_async(self, coroutine):
        """
        Run an asynchronous coroutine.

        Args:
            coroutine: Asynchronous coroutine

        Returns:
            Coroutine result
        """
        return asyncio.run(coroutine)

    def get_models_sync(self, model_type: Optional[str] = None) -> List[Dict[str, Any]]:
        """
        Get available models (synchronous version).

        Args:
            model_type: Filter by model type

        Returns:
            List of models
        """
        params = {}
        if model_type:
            params["model_type"] = model_type
        
        return self._make_request("GET", "models", params=params)

    def predict_sync(
        self,
        service_id: str,
        input_text: str,
        parameters: Optional[Dict[str, Any]] = None,
    ) -> Dict[str, Any]:
        """
        Make a prediction using an inference service (synchronous version).

        Args:
            service_id: Service ID
            input_text: Input text
            parameters: Prediction parameters

        Returns:
            Prediction result
        """
        data = {
            "input_text": input_text,
        }
        
        if parameters:
            data["parameters"] = parameters
        
        return self._make_request("POST", f"inference/services/{service_id}/predict", data=data)

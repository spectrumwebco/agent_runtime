"""
CrunchyData PostgreSQL Operator integration for Django.

This module provides integration with CrunchyData's PostgreSQL Operator,
which manages PostgreSQL clusters on Kubernetes.
"""

import logging
import os
import subprocess
import json
from typing import Dict, Any, List, Optional, Union
from django.conf import settings
from pydantic import BaseModel, Field

logger = logging.getLogger(__name__)

class PostgresClusterConfig(BaseModel):
    """Postgres cluster configuration."""
    
    name: str = Field(default="agent-postgres-cluster")
    namespace: str = Field(default="default")
    version: str = Field(default="14")
    instances: int = Field(default=2)
    size: str = Field(default="10Gi")
    use_mock: bool = Field(default=False)

class CrunchyDataClient:
    """Client for CrunchyData PostgreSQL Operator."""
    
    def __init__(self, config: Optional[PostgresClusterConfig] = None):
        """Initialize the CrunchyData PostgreSQL Operator client."""
        self.config = config or self._get_default_config()
        self._use_mock = self.config.use_mock
        
        self._in_kubernetes = os.path.exists('/var/run/secrets/kubernetes.io/serviceaccount/token')
        
        if not self._in_kubernetes:
            logger.warning("Not running in Kubernetes, using mock mode")
            self._use_mock = True
        
        if self._use_mock:
            logger.warning("CrunchyData PostgreSQL Operator client running in mock mode")
    
    def _get_default_config(self) -> PostgresClusterConfig:
        """Get the default configuration from settings."""
        if os.environ.get('CI') == 'true':
            use_mock = True
        else:
            try:
                result = subprocess.run(
                    ["kubectl", "version", "--client"],
                    capture_output=True,
                    text=True,
                    check=False
                )
                use_mock = result.returncode != 0
            except (subprocess.SubprocessError, FileNotFoundError):
                use_mock = True
        
        return PostgresClusterConfig(
            name=getattr(settings, 'POSTGRES_CLUSTER_NAME', 'agent-postgres-cluster'),
            namespace=getattr(settings, 'POSTGRES_CLUSTER_NAMESPACE', 'default'),
            version=getattr(settings, 'POSTGRES_VERSION', '14'),
            instances=getattr(settings, 'POSTGRES_INSTANCES', 2),
            size=getattr(settings, 'POSTGRES_SIZE', '10Gi'),
            use_mock=use_mock
        )
    
    def get_cluster_status(self) -> Dict[str, Any]:
        """Get the status of the Postgres cluster."""
        if self._use_mock:
            return {
                "name": self.config.name,
                "namespace": self.config.namespace,
                "status": "Running",
                "instances": self.config.instances,
                "version": self.config.version,
                "mocked": True
            }
        
        try:
            result = subprocess.run(
                [
                    "kubectl", "get", "postgrescluster",
                    self.config.name,
                    "-n", self.config.namespace,
                    "-o", "json"
                ],
                capture_output=True,
                text=True,
                check=True
            )
            
            cluster_info = json.loads(result.stdout)
            
            return {
                "name": cluster_info.get("metadata", {}).get("name"),
                "namespace": cluster_info.get("metadata", {}).get("namespace"),
                "status": cluster_info.get("status", {}).get("conditions", [{}])[0].get("status"),
                "instances": len(cluster_info.get("status", {}).get("instances", [])),
                "version": cluster_info.get("spec", {}).get("postgresVersion"),
                "mocked": False
            }
        except Exception as e:
            logger.error(f"Error getting Postgres cluster status: {e}")
            return {
                "name": self.config.name,
                "namespace": self.config.namespace,
                "status": "Unknown",
                "error": str(e),
                "mocked": False
            }
    
    def get_connection_info(self) -> Dict[str, Any]:
        """Get connection information for the Postgres cluster."""
        if self._use_mock:
            return {
                "host": f"{self.config.name}-primary.{self.config.namespace}.svc.cluster.local",
                "port": 5432,
                "database": "agent_runtime",
                "user": "agent_user",
                "password": "mock_password",
                "mocked": True
            }
        
        try:
            result = subprocess.run(
                [
                    "kubectl", "get", "postgrescluster",
                    self.config.name,
                    "-n", self.config.namespace,
                    "-o", "jsonpath={.status.userInterface.secretName}"
                ],
                capture_output=True,
                text=True,
                check=True
            )
            
            secret_name = result.stdout.strip()
            
            result = subprocess.run(
                [
                    "kubectl", "get", "secret",
                    secret_name,
                    "-n", self.config.namespace,
                    "-o", "jsonpath={.data.password}"
                ],
                capture_output=True,
                text=True,
                check=True
            )
            
            import base64
            password = base64.b64decode(result.stdout.strip()).decode('utf-8')
            
            result = subprocess.run(
                [
                    "kubectl", "get", "secret",
                    secret_name,
                    "-n", self.config.namespace,
                    "-o", "jsonpath={.data.user}"
                ],
                capture_output=True,
                text=True,
                check=True
            )
            
            username = base64.b64decode(result.stdout.strip()).decode('utf-8')
            
            result = subprocess.run(
                [
                    "kubectl", "get", "service",
                    f"{self.config.name}-primary",
                    "-n", self.config.namespace,
                    "-o", "jsonpath={.spec.clusterIP}"
                ],
                capture_output=True,
                text=True,
                check=True
            )
            
            host = f"{self.config.name}-primary.{self.config.namespace}.svc.cluster.local"
            
            return {
                "host": host,
                "port": 5432,
                "database": "agent_runtime",
                "user": username,
                "password": password,
                "mocked": False
            }
        except Exception as e:
            logger.error(f"Error getting Postgres connection info: {e}")
            return {
                "host": f"{self.config.name}-primary.{self.config.namespace}.svc.cluster.local",
                "port": 5432,
                "database": "agent_runtime",
                "user": "agent_user",
                "password": "",
                "error": str(e),
                "mocked": False
            }
    
    def check_connection(self) -> Dict[str, Any]:
        """Check the connection to the CrunchyData PostgreSQL Operator."""
        if self._use_mock:
            return {
                "connected": False,
                "mocked": True,
                "message": "Running in mock mode"
            }
        
        try:
            result = subprocess.run(
                [
                    "kubectl", "get", "crd",
                    "postgresclusters.postgres-operator.crunchydata.com"
                ],
                capture_output=True,
                text=True,
                check=False
            )
            
            if result.returncode != 0:
                return {
                    "connected": False,
                    "mocked": False,
                    "message": "CrunchyData PostgreSQL Operator CRD not installed"
                }
            
            result = subprocess.run(
                [
                    "kubectl", "get", "deployment",
                    "postgres-operator",
                    "-n", "pgo",
                    "-o", "jsonpath={.status.readyReplicas}"
                ],
                capture_output=True,
                text=True,
                check=False
            )
            
            if result.returncode != 0 or not result.stdout.strip():
                return {
                    "connected": False,
                    "mocked": False,
                    "message": "CrunchyData PostgreSQL Operator not running"
                }
            
            return {
                "connected": True,
                "mocked": False,
                "message": "Successfully connected to CrunchyData PostgreSQL Operator"
            }
        except Exception as e:
            logger.error(f"Error checking CrunchyData PostgreSQL Operator connection: {e}")
            return {
                "connected": False,
                "mocked": False,
                "message": f"Connection error: {str(e)}"
            }
    
    def create_cluster(self, config: Optional[Dict[str, Any]] = None) -> bool:
        """Create a Postgres cluster."""
        if self._use_mock:
            logger.info(f"Mock creating Postgres cluster: {self.config.name}")
            return True
        
        try:
            import tempfile
            import yaml
            
            cluster_config = config or {
                "apiVersion": "postgres-operator.crunchydata.com/v1beta1",
                "kind": "PostgresCluster",
                "metadata": {
                    "name": self.config.name,
                    "namespace": self.config.namespace
                },
                "spec": {
                    "image": f"registry.developers.crunchydata.com/crunchydata/crunchy-postgres:ubi8-{self.config.version}.5-0",
                    "postgresVersion": self.config.version,
                    "instances": [
                        {
                            "name": "instance1",
                            "replicas": self.config.instances,
                            "dataVolumeClaimSpec": {
                                "accessModes": ["ReadWriteOnce"],
                                "resources": {
                                    "requests": {
                                        "storage": self.config.size
                                    }
                                }
                            }
                        }
                    ],
                    "backups": {
                        "pgbackrest": {
                            "image": "registry.developers.crunchydata.com/crunchydata/crunchy-pgbackrest:ubi8-2.38-0",
                            "repos": [
                                {
                                    "name": "repo1",
                                    "volume": {
                                        "volumeClaimSpec": {
                                            "accessModes": ["ReadWriteOnce"],
                                            "resources": {
                                                "requests": {
                                                    "storage": "5Gi"
                                                }
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "users": [
                        {
                            "name": "agent_user",
                            "databases": ["agent_runtime", "agent_db", "trajectory_db", "ml_db"],
                            "options": "SUPERUSER CREATEDB"
                        },
                        {
                            "name": "app_user",
                            "databases": ["agent_runtime"]
                        }
                    ]
                }
            }
            
            with tempfile.NamedTemporaryFile(mode='w', suffix='.yaml', delete=False) as temp:
                yaml.dump(cluster_config, temp)
                temp_path = temp.name
            
            result = subprocess.run(
                [
                    "kubectl", "apply", "-f", temp_path
                ],
                capture_output=True,
                text=True,
                check=True
            )
            
            os.unlink(temp_path)
            
            logger.info(f"Created Postgres cluster: {self.config.name}")
            return True
        except Exception as e:
            logger.error(f"Error creating Postgres cluster: {e}")
            return False
    
    def delete_cluster(self) -> bool:
        """Delete a Postgres cluster."""
        if self._use_mock:
            logger.info(f"Mock deleting Postgres cluster: {self.config.name}")
            return True
        
        try:
            result = subprocess.run(
                [
                    "kubectl", "delete", "postgrescluster",
                    self.config.name,
                    "-n", self.config.namespace
                ],
                capture_output=True,
                text=True,
                check=True
            )
            
            logger.info(f"Deleted Postgres cluster: {self.config.name}")
            return True
        except Exception as e:
            logger.error(f"Error deleting Postgres cluster: {e}")
            return False

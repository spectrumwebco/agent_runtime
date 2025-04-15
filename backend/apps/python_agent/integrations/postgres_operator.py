"""
Zalando Postgres Operator integration for Django.

This module provides integration with Zalando's Postgres Operator,
which manages PostgreSQL clusters on Kubernetes.
"""

import logging
import os
import subprocess
from typing import Dict, Any, List, Optional, Union
from django.conf import settings
from pydantic import BaseModel, Field

logger = logging.getLogger(__name__)

class PostgresClusterConfig(BaseModel):
    """Postgres cluster configuration."""
    
    name: str = Field(default="agent-postgres-cluster")
    namespace: str = Field(default="default")
    team_id: str = Field(default="agent")
    version: str = Field(default="15")
    instances: int = Field(default=2)
    size: str = Field(default="10Gi")
    use_mock: bool = Field(default=False)

class PostgresOperatorClient:
    """Client for Zalando Postgres Operator."""
    
    def __init__(self, config: Optional[PostgresClusterConfig] = None):
        """Initialize the Postgres Operator client."""
        self.config = config or self._get_default_config()
        self._use_mock = self.config.use_mock
        
        self._in_kubernetes = os.path.exists('/var/run/secrets/kubernetes.io/serviceaccount/token')
        
        if not self._in_kubernetes:
            logger.warning("Not running in Kubernetes, using mock mode")
            self._use_mock = True
        
        if self._use_mock:
            logger.warning("Postgres Operator client running in mock mode")
    
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
            team_id=getattr(settings, 'POSTGRES_TEAM_ID', 'agent'),
            version=getattr(settings, 'POSTGRES_VERSION', '15'),
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
                    "kubectl", "get", "postgresql",
                    self.config.name,
                    "-n", self.config.namespace,
                    "-o", "json"
                ],
                capture_output=True,
                text=True,
                check=True
            )
            
            import json
            cluster_info = json.loads(result.stdout)
            
            return {
                "name": cluster_info.get("metadata", {}).get("name"),
                "namespace": cluster_info.get("metadata", {}).get("namespace"),
                "status": cluster_info.get("status", {}).get("clusterStatus"),
                "instances": cluster_info.get("spec", {}).get("numberOfInstances"),
                "version": cluster_info.get("status", {}).get("postgresqlVersion"),
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
                "host": f"{self.config.name}.{self.config.namespace}.svc.cluster.local",
                "port": 5432,
                "database": "agent_runtime",
                "user": "agent_user",
                "password": "mock_password",
                "mocked": True
            }
        
        try:
            result = subprocess.run(
                [
                    "kubectl", "get", "postgresql",
                    self.config.name,
                    "-n", self.config.namespace,
                    "-o", "jsonpath={.spec.teamId}"
                ],
                capture_output=True,
                text=True,
                check=True
            )
            
            team_id = result.stdout.strip()
            secret_name = f"{self.config.name}.{team_id}.credentials.postgresql.acid.zalan.do"
            
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
                    "-o", "jsonpath={.data.username}"
                ],
                capture_output=True,
                text=True,
                check=True
            )
            
            username = base64.b64decode(result.stdout.strip()).decode('utf-8')
            
            return {
                "host": f"{self.config.name}.{self.config.namespace}.svc.cluster.local",
                "port": 5432,
                "database": "agent_runtime",
                "user": username,
                "password": password,
                "mocked": False
            }
        except Exception as e:
            logger.error(f"Error getting Postgres connection info: {e}")
            return {
                "host": f"{self.config.name}.{self.config.namespace}.svc.cluster.local",
                "port": 5432,
                "database": "agent_runtime",
                "user": "agent_user",
                "password": "",
                "error": str(e),
                "mocked": False
            }
    
    def check_connection(self) -> Dict[str, Any]:
        """Check the connection to the Postgres Operator."""
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
                    "postgresqls.acid.zalan.do"
                ],
                capture_output=True,
                text=True,
                check=False
            )
            
            if result.returncode != 0:
                return {
                    "connected": False,
                    "mocked": False,
                    "message": "Postgres Operator CRD not installed"
                }
            
            result = subprocess.run(
                [
                    "kubectl", "get", "deployment",
                    "postgres-operator",
                    "-n", "postgres-operator",
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
                    "message": "Postgres Operator not running"
                }
            
            return {
                "connected": True,
                "mocked": False,
                "message": "Successfully connected to Postgres Operator"
            }
        except Exception as e:
            logger.error(f"Error checking Postgres Operator connection: {e}")
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
                "apiVersion": "acid.zalan.do/v1",
                "kind": "postgresql",
                "metadata": {
                    "name": self.config.name,
                    "namespace": self.config.namespace
                },
                "spec": {
                    "teamId": self.config.team_id,
                    "volume": {
                        "size": self.config.size
                    },
                    "numberOfInstances": self.config.instances,
                    "users": {
                        "agent_user": ["superuser", "createdb"],
                        "app_user": []
                    },
                    "databases": {
                        "agent_runtime": "agent_user",
                        "agent_db": "agent_user",
                        "trajectory_db": "agent_user",
                        "ml_db": "agent_user"
                    },
                    "postgresql": {
                        "version": self.config.version,
                        "parameters": {
                            "shared_buffers": "256MB",
                            "max_connections": "200",
                            "log_statement": "all"
                        }
                    }
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
                    "kubectl", "delete", "postgresql",
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

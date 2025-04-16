"""
Django integration with PostgreSQL Operator.

This module provides integration between Django and PostgreSQL Operator,
implementing database operations and cluster management.
"""

import logging
import json
import subprocess
from typing import Dict, Any, Optional, List, Union
from django.conf import settings

logger = logging.getLogger(__name__)

class PostgresOperatorManager:
    """
    Manager for PostgreSQL Operator.
    
    This manager handles cluster management operations
    for PostgreSQL databases managed by the PostgreSQL Operator.
    """
    
    def __init__(self, namespace=None):
        """
        Initialize the PostgreSQL Operator manager.
        
        Args:
            namespace: Kubernetes namespace for the PostgreSQL Operator
        """
        self.namespace = namespace or getattr(settings, 'POSTGRES_OPERATOR_NAMESPACE', 'default')
    
    def apply_cluster_config(self, config_file):
        """
        Apply a PostgreSQL cluster configuration.
        
        Args:
            config_file: Path to the configuration file
            
        Returns:
            True if successful, False otherwise
        """
        try:
            result = subprocess.run(
                ['kubectl', 'apply', '-f', config_file, '-n', self.namespace],
                capture_output=True,
                text=True,
                check=True
            )
            logger.info(f"Applied PostgreSQL cluster configuration: {result.stdout}")
            return True
        except subprocess.CalledProcessError as e:
            logger.error(f"Error applying PostgreSQL cluster configuration: {e.stderr}")
            return False
    
    def delete_cluster(self, cluster_name):
        """
        Delete a PostgreSQL cluster.
        
        Args:
            cluster_name: Name of the cluster to delete
            
        Returns:
            True if successful, False otherwise
        """
        try:
            result = subprocess.run(
                ['kubectl', 'delete', 'postgrescluster', cluster_name, '-n', self.namespace],
                capture_output=True,
                text=True,
                check=True
            )
            logger.info(f"Deleted PostgreSQL cluster {cluster_name}: {result.stdout}")
            return True
        except subprocess.CalledProcessError as e:
            logger.error(f"Error deleting PostgreSQL cluster {cluster_name}: {e.stderr}")
            return False
    
    def get_clusters(self):
        """
        Get a list of PostgreSQL clusters managed by the operator.
        
        Returns:
            List of cluster information
        """
        try:
            result = subprocess.run(
                ['kubectl', 'get', 'postgresclusters', '-n', self.namespace, '-o', 'json'],
                capture_output=True,
                text=True,
                check=True
            )
            clusters = json.loads(result.stdout)
            return clusters.get('items', [])
        except subprocess.CalledProcessError as e:
            logger.error(f"Error getting PostgreSQL clusters: {e.stderr}")
            return []
        except json.JSONDecodeError as e:
            logger.error(f"Error parsing PostgreSQL clusters: {str(e)}")
            return []
    
    def get_cluster(self, cluster_name):
        """
        Get information about a specific PostgreSQL cluster.
        
        Args:
            cluster_name: Name of the cluster
            
        Returns:
            Cluster information
        """
        try:
            result = subprocess.run(
                ['kubectl', 'get', 'postgrescluster', cluster_name, '-n', self.namespace, '-o', 'json'],
                capture_output=True,
                text=True,
                check=True
            )
            return json.loads(result.stdout)
        except subprocess.CalledProcessError as e:
            logger.error(f"Error getting PostgreSQL cluster {cluster_name}: {e.stderr}")
            return {}
        except json.JSONDecodeError as e:
            logger.error(f"Error parsing PostgreSQL cluster {cluster_name}: {str(e)}")
            return {}
    
    def get_cluster_status(self, cluster_name):
        """
        Get the status of a PostgreSQL cluster.
        
        Args:
            cluster_name: Name of the cluster
            
        Returns:
            Cluster status
        """
        cluster = self.get_cluster(cluster_name)
        return cluster.get('status', {})
    
    def get_cluster_connection_info(self, cluster_name):
        """
        Get connection information for a PostgreSQL cluster.
        
        Args:
            cluster_name: Name of the cluster
            
        Returns:
            Connection information
        """
        status = self.get_cluster_status(cluster_name)
        return status.get('pgbouncer', {}).get('service', {})
    
    def get_cluster_pods(self, cluster_name):
        """
        Get pods for a PostgreSQL cluster.
        
        Args:
            cluster_name: Name of the cluster
            
        Returns:
            List of pod information
        """
        try:
            result = subprocess.run(
                ['kubectl', 'get', 'pods', '-l', f'postgres-operator.crunchydata.com/cluster={cluster_name}', '-n', self.namespace, '-o', 'json'],
                capture_output=True,
                text=True,
                check=True
            )
            pods = json.loads(result.stdout)
            return pods.get('items', [])
        except subprocess.CalledProcessError as e:
            logger.error(f"Error getting pods for PostgreSQL cluster {cluster_name}: {e.stderr}")
            return []
        except json.JSONDecodeError as e:
            logger.error(f"Error parsing pods for PostgreSQL cluster {cluster_name}: {str(e)}")
            return []
    
    def get_cluster_services(self, cluster_name):
        """
        Get services for a PostgreSQL cluster.
        
        Args:
            cluster_name: Name of the cluster
            
        Returns:
            List of service information
        """
        try:
            result = subprocess.run(
                ['kubectl', 'get', 'services', '-l', f'postgres-operator.crunchydata.com/cluster={cluster_name}', '-n', self.namespace, '-o', 'json'],
                capture_output=True,
                text=True,
                check=True
            )
            services = json.loads(result.stdout)
            return services.get('items', [])
        except subprocess.CalledProcessError as e:
            logger.error(f"Error getting services for PostgreSQL cluster {cluster_name}: {e.stderr}")
            return []
        except json.JSONDecodeError as e:
            logger.error(f"Error parsing services for PostgreSQL cluster {cluster_name}: {str(e)}")
            return []
    
    def get_cluster_secrets(self, cluster_name):
        """
        Get secrets for a PostgreSQL cluster.
        
        Args:
            cluster_name: Name of the cluster
            
        Returns:
            List of secret information
        """
        try:
            result = subprocess.run(
                ['kubectl', 'get', 'secrets', '-l', f'postgres-operator.crunchydata.com/cluster={cluster_name}', '-n', self.namespace, '-o', 'json'],
                capture_output=True,
                text=True,
                check=True
            )
            secrets = json.loads(result.stdout)
            return secrets.get('items', [])
        except subprocess.CalledProcessError as e:
            logger.error(f"Error getting secrets for PostgreSQL cluster {cluster_name}: {e.stderr}")
            return []
        except json.JSONDecodeError as e:
            logger.error(f"Error parsing secrets for PostgreSQL cluster {cluster_name}: {str(e)}")
            return []
    
    def get_operator_status(self):
        """
        Get the status of the PostgreSQL Operator.
        
        Returns:
            Operator status
        """
        try:
            result = subprocess.run(
                ['kubectl', 'get', 'deployment', 'postgres-operator', '-n', self.namespace, '-o', 'json'],
                capture_output=True,
                text=True,
                check=True
            )
            deployment = json.loads(result.stdout)
            return deployment.get('status', {})
        except subprocess.CalledProcessError as e:
            logger.error(f"Error getting PostgreSQL Operator status: {e.stderr}")
            return {}
        except json.JSONDecodeError as e:
            logger.error(f"Error parsing PostgreSQL Operator status: {str(e)}")
            return {}


def get_postgres_operator_manager(namespace=None):
    """
    Get a PostgreSQL Operator manager instance.
    
    Args:
        namespace: Kubernetes namespace for the PostgreSQL Operator
    
    Returns:
        PostgresOperatorManager instance
    """
    return PostgresOperatorManager(namespace)

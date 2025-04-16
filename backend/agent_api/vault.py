"""
Hashicorp Vault integration for the agent_api project.

This module provides integration with Hashicorp Vault for securely
storing and retrieving database credentials and other secrets.
"""

import os
import logging
import hvac
from typing import Dict, Any, Optional
from django.conf import settings

logger = logging.getLogger(__name__)


class VaultClient:
    """
    Client for interacting with Hashicorp Vault.
    
    This client provides methods for authenticating with Vault,
    reading and writing secrets, and managing policies.
    """
    
    def __init__(self, url=None, token=None, role_id=None, secret_id=None):
        """
        Initialize the Vault client.
        
        Args:
            url: Vault URL
            token: Vault token
            role_id: AppRole role ID
            secret_id: AppRole secret ID
        """
        self.url = url or os.environ.get('VAULT_ADDR', 'http://vault.default.svc.cluster.local:8200')
        self.token = token or os.environ.get('VAULT_TOKEN', '')
        self.role_id = role_id or os.environ.get('VAULT_ROLE_ID', '')
        self.secret_id = secret_id or os.environ.get('VAULT_SECRET_ID', '')
        
        self.client = None
        self.authenticated = False
        
        self.initialize()
    
    def initialize(self):
        """Initialize the Vault client and authenticate."""
        try:
            is_local = not os.path.exists('/var/run/secrets/kubernetes.io/serviceaccount/token')
            
            if is_local and 'vault.default.svc.cluster.local' in self.url:
                self.url = 'http://localhost:8200'
                logger.info(f"Local development detected, using Vault URL: {self.url}")
            
            self.client = hvac.Client(url=self.url)
            
            if self.token:
                self.client.token = self.token
                self.authenticated = self.client.is_authenticated()
            elif self.role_id and self.secret_id:
                self._authenticate_approle()
            else:
                self._authenticate_kubernetes()
            
            if not self.authenticated:
                if is_local:
                    logger.warning("Failed to authenticate with Vault in local development mode")
                else:
                    logger.warning("Failed to authenticate with Vault")
        
        except Exception as e:
            logger.error(f"Error initializing Vault client: {e}")
    
    def _authenticate_approle(self):
        """Authenticate with Vault using AppRole."""
        try:
            response = self.client.auth.approle.login(
                role_id=self.role_id,
                secret_id=self.secret_id
            )
            
            self.client.token = response['auth']['client_token']
            self.authenticated = self.client.is_authenticated()
        
        except Exception as e:
            logger.error(f"Error authenticating with Vault using AppRole: {e}")
    
    def _authenticate_kubernetes(self):
        """Authenticate with Vault using Kubernetes."""
        if not os.path.exists('/var/run/secrets/kubernetes.io/serviceaccount/token'):
            logger.warning("Not running in Kubernetes, skipping Kubernetes authentication")
            return
            
        try:
            with open('/var/run/secrets/kubernetes.io/serviceaccount/token', 'r') as f:
                jwt = f.read()
            
            response = self.client.auth.kubernetes.login(
                role='agent-api',
                jwt=jwt
            )
            
            self.client.token = response['auth']['client_token']
            self.authenticated = self.client.is_authenticated()
        
        except Exception as e:
            logger.error(f"Error authenticating with Vault using Kubernetes: {e}")
    
    def read_secret(self, path: str) -> Optional[Dict[str, Any]]:
        """
        Read a secret from Vault.
        
        Args:
            path: Path to the secret
            
        Returns:
            Dict[str, Any]: Secret data, or None if an error occurred
        """
        is_local = not os.path.exists('/var/run/secrets/kubernetes.io/serviceaccount/token')
        
        if not self.authenticated:
            self.initialize()
            
            if not self.authenticated:
                if is_local:
                    logger.warning("Not authenticated with Vault in local development mode")
                else:
                    logger.error("Not authenticated with Vault")
                return None
        
        try:
            response = self.client.secrets.kv.v2.read_secret_version(path=path)
            
            return response['data']['data']
        
        except Exception as e:
            if is_local:
                logger.warning(f"Error reading secret from Vault in local development mode: {e}")
            else:
                logger.error(f"Error reading secret from Vault: {e}")
            return None
    
    def write_secret(self, path: str, data: Dict[str, Any]) -> bool:
        """
        Write a secret to Vault.
        
        Args:
            path: Path to the secret
            data: Secret data
            
        Returns:
            bool: True if the write was successful, False otherwise
        """
        if not self.authenticated:
            self.initialize()
            
            if not self.authenticated:
                logger.error("Not authenticated with Vault")
                return False
        
        try:
            self.client.secrets.kv.v2.create_or_update_secret(
                path=path,
                secret=data
            )
            
            return True
        
        except Exception as e:
            logger.error(f"Error writing secret to Vault: {e}")
            return False
    
    def delete_secret(self, path: str) -> bool:
        """
        Delete a secret from Vault.
        
        Args:
            path: Path to the secret
            
        Returns:
            bool: True if the delete was successful, False otherwise
        """
        if not self.authenticated:
            self.initialize()
            
            if not self.authenticated:
                logger.error("Not authenticated with Vault")
                return False
        
        try:
            self.client.secrets.kv.v2.delete_metadata_and_all_versions(path=path)
            
            return True
        
        except Exception as e:
            logger.error(f"Error deleting secret from Vault: {e}")
            return False


class DatabaseSecrets:
    """
    Manager for database secrets stored in Vault.
    
    This class provides methods for retrieving database credentials
    from Vault and configuring Django database settings.
    """
    
    def __init__(self, vault_client: VaultClient = None):
        """
        Initialize the database secrets manager.
        
        Args:
            vault_client: Vault client
        """
        self.vault_client = vault_client or VaultClient()
    
    def get_database_credentials(self, database: str) -> Optional[Dict[str, Any]]:
        """
        Get database credentials from Vault.
        
        Args:
            database: Database name
            
        Returns:
            Dict[str, Any]: Database credentials, or None if an error occurred
        """
        path = f"database/{database}"
        
        return self.vault_client.read_secret(path)
    
    def store_database_credentials(self, database: str, credentials: Dict[str, Any]) -> bool:
        """
        Store database credentials in Vault.
        
        Args:
            database: Database name
            credentials: Database credentials
            
        Returns:
            bool: True if the store was successful, False otherwise
        """
        path = f"database/{database}"
        
        return self.vault_client.write_secret(path, credentials)
    
    def configure_django_databases(self) -> Dict[str, Dict[str, Any]]:
        """
        Configure Django database settings using credentials from Vault.
        
        Returns:
            Dict[str, Dict[str, Any]]: Django database settings
        """
        databases = {}
        
        default_credentials = self.get_database_credentials('default')
        
        if default_credentials:
            databases['default'] = {
                'ENGINE': default_credentials.get('engine', 'django.db.backends.postgresql'),
                'NAME': default_credentials.get('name', 'postgres'),
                'USER': default_credentials.get('user', 'postgres'),
                'PASSWORD': default_credentials.get('password', 'postgres'),
                'HOST': default_credentials.get('host', 'supabase-db.default.svc.cluster.local'),
                'PORT': default_credentials.get('port', 5432),
                'OPTIONS': {
                    'sslmode': 'require',
                },
            }
        
        agent_credentials = self.get_database_credentials('agent')
        
        if agent_credentials:
            databases['agent'] = {
                'ENGINE': agent_credentials.get('engine', 'django.db.backends.postgresql'),
                'NAME': agent_credentials.get('name', 'agent_db'),
                'USER': agent_credentials.get('user', 'postgres'),
                'PASSWORD': agent_credentials.get('password', 'postgres'),
                'HOST': agent_credentials.get('host', 'supabase-db.default.svc.cluster.local'),
                'PORT': agent_credentials.get('port', 5432),
                'OPTIONS': {
                    'sslmode': 'require',
                },
            }
        
        trajectory_credentials = self.get_database_credentials('trajectory')
        
        if trajectory_credentials:
            databases['trajectory'] = {
                'ENGINE': trajectory_credentials.get('engine', 'django.db.backends.postgresql'),
                'NAME': trajectory_credentials.get('name', 'trajectory_db'),
                'USER': trajectory_credentials.get('user', 'postgres'),
                'PASSWORD': trajectory_credentials.get('password', 'postgres'),
                'HOST': trajectory_credentials.get('host', 'supabase-db.default.svc.cluster.local'),
                'PORT': trajectory_credentials.get('port', 5432),
                'OPTIONS': {
                    'sslmode': 'require',
                },
            }
        
        ml_credentials = self.get_database_credentials('ml')
        
        if ml_credentials:
            databases['ml'] = {
                'ENGINE': ml_credentials.get('engine', 'django.db.backends.postgresql'),
                'NAME': ml_credentials.get('name', 'ml_db'),
                'USER': ml_credentials.get('user', 'postgres'),
                'PASSWORD': ml_credentials.get('password', 'postgres'),
                'HOST': ml_credentials.get('host', 'supabase-db.default.svc.cluster.local'),
                'PORT': ml_credentials.get('port', 5432),
                'OPTIONS': {
                    'sslmode': 'require',
                },
            }
        
        return databases


vault_client = VaultClient()

database_secrets = DatabaseSecrets(vault_client)

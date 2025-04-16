"""
Django integration with CrunchyData PostgreSQL Operator.

This module provides integration between Django and CrunchyData PostgreSQL Operator,
implementing database operations and cluster management.
"""

import logging
import json
import subprocess
from typing import Dict, Any, Optional, List, Union
from django.conf import settings
import psycopg2
from psycopg2.extras import RealDictCursor

logger = logging.getLogger(__name__)

class PostgresOperatorClient:
    """
    Client for CrunchyData PostgreSQL Operator.
    
    This client handles database operations and cluster management
    for PostgreSQL databases managed by the CrunchyData PostgreSQL Operator.
    """
    
    def __init__(self, connection_name='postgres'):
        """
        Initialize the PostgreSQL Operator client.
        
        Args:
            connection_name: Name of the Django database connection to use.
                            Defaults to 'postgres'.
        """
        self.connection_name = connection_name
        self.db_settings = settings.DATABASES.get(connection_name, {})
        self.namespace = getattr(settings, 'POSTGRES_OPERATOR_NAMESPACE', 'default')
    
    def get_connection(self):
        """
        Get a PostgreSQL database connection.
        
        Returns:
            PostgreSQL database connection
        """
        return psycopg2.connect(
            host=self.db_settings.get('HOST', 'localhost'),
            port=self.db_settings.get('PORT', 5432),
            user=self.db_settings.get('USER', 'postgres'),
            password=self.db_settings.get('PASSWORD', ''),
            dbname=self.db_settings.get('NAME', 'postgres'),
        )
    
    def execute_query(self, query, params=None):
        """
        Execute a SQL query on the PostgreSQL database.
        
        Args:
            query: SQL query to execute
            params: Parameters for the query
            
        Returns:
            Query results as a list of dictionaries
        """
        with self.get_connection() as conn:
            with conn.cursor(cursor_factory=RealDictCursor) as cursor:
                cursor.execute(query, params or ())
                return cursor.fetchall()
    
    def execute_update(self, query, params=None):
        """
        Execute an update query on the PostgreSQL database.
        
        Args:
            query: SQL query to execute
            params: Parameters for the query
            
        Returns:
            Number of rows affected
        """
        with self.get_connection() as conn:
            with conn.cursor() as cursor:
                cursor.execute(query, params or ())
                conn.commit()
                return cursor.rowcount
    
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
    
    def create_database(self, database_name):
        """
        Create a new database in the PostgreSQL cluster.
        
        Args:
            database_name: Name of the database to create
            
        Returns:
            True if successful, False otherwise
        """
        try:
            with self.get_connection() as conn:
                conn.autocommit = True
                with conn.cursor() as cursor:
                    cursor.execute(f"CREATE DATABASE {database_name}")
            return True
        except Exception as e:
            logger.error(f"Error creating database {database_name}: {str(e)}")
            return False
    
    def drop_database(self, database_name):
        """
        Drop a database from the PostgreSQL cluster.
        
        Args:
            database_name: Name of the database to drop
            
        Returns:
            True if successful, False otherwise
        """
        try:
            with self.get_connection() as conn:
                conn.autocommit = True
                with conn.cursor() as cursor:
                    cursor.execute(f"DROP DATABASE IF EXISTS {database_name}")
            return True
        except Exception as e:
            logger.error(f"Error dropping database {database_name}: {str(e)}")
            return False
    
    def create_user(self, username, password):
        """
        Create a new user in the PostgreSQL cluster.
        
        Args:
            username: Username for the new user
            password: Password for the new user
            
        Returns:
            True if successful, False otherwise
        """
        try:
            with self.get_connection() as conn:
                conn.autocommit = True
                with conn.cursor() as cursor:
                    cursor.execute(f"CREATE USER {username} WITH PASSWORD %s", (password,))
            return True
        except Exception as e:
            logger.error(f"Error creating user {username}: {str(e)}")
            return False
    
    def drop_user(self, username):
        """
        Drop a user from the PostgreSQL cluster.
        
        Args:
            username: Username to drop
            
        Returns:
            True if successful, False otherwise
        """
        try:
            with self.get_connection() as conn:
                conn.autocommit = True
                with conn.cursor() as cursor:
                    cursor.execute(f"DROP USER IF EXISTS {username}")
            return True
        except Exception as e:
            logger.error(f"Error dropping user {username}: {str(e)}")
            return False
    
    def grant_privileges(self, username, database_name, privileges='ALL'):
        """
        Grant privileges to a user on a database.
        
        Args:
            username: Username to grant privileges to
            database_name: Name of the database
            privileges: Privileges to grant
            
        Returns:
            True if successful, False otherwise
        """
        try:
            with self.get_connection() as conn:
                conn.autocommit = True
                with conn.cursor() as cursor:
                    cursor.execute(f"GRANT {privileges} ON DATABASE {database_name} TO {username}")
            return True
        except Exception as e:
            logger.error(f"Error granting privileges to {username} on {database_name}: {str(e)}")
            return False
    
    def get_connection_info(self):
        """
        Get information about the current database connection.
        
        Returns:
            Dictionary with connection information
        """
        return {
            'engine': self.db_settings.get('ENGINE', ''),
            'name': self.db_settings.get('NAME', ''),
            'host': self.db_settings.get('HOST', ''),
            'port': self.db_settings.get('PORT', ''),
            'user': self.db_settings.get('USER', ''),
        }


def get_postgres_operator_client(connection_name='postgres'):
    """
    Get a PostgreSQL Operator client instance.
    
    Args:
        connection_name: Name of the Django database connection to use.
                        Defaults to 'postgres'.
    
    Returns:
        PostgresOperatorClient instance
    """
    return PostgresOperatorClient(connection_name)

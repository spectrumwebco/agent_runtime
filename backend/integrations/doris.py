"""
Django integration with Apache Doris.

This module provides integration between Django and Apache Doris,
implementing database operations and query functionality.
"""

import logging
import pymysql
from django.conf import settings
from django.db import connections

logger = logging.getLogger(__name__)

class DorisClient:
    """
    Client for Apache Doris database operations.
    
    This client handles database operations and query execution
    for Apache Doris using the MySQL protocol.
    """
    
    def __init__(self, connection_name='default'):
        """
        Initialize the Doris client.
        
        Args:
            connection_name: Name of the Django database connection to use.
                            Defaults to 'default'.
        """
        self.connection_name = connection_name
    
    def execute_query(self, query, params=None):
        """
        Execute a SQL query on the Doris database.
        
        Args:
            query: SQL query to execute
            params: Parameters for the query
            
        Returns:
            Query results as a list of dictionaries
        """
        with connections[self.connection_name].cursor() as cursor:
            cursor.execute(query, params or ())
            columns = [col[0] for col in cursor.description]
            return [dict(zip(columns, row)) for row in cursor.fetchall()]
    
    def execute_update(self, query, params=None):
        """
        Execute an update query on the Doris database.
        
        Args:
            query: SQL query to execute
            params: Parameters for the query
            
        Returns:
            Number of rows affected
        """
        with connections[self.connection_name].cursor() as cursor:
            cursor.execute(query, params or ())
            return cursor.rowcount
    
    def create_table(self, table_name, columns, partition_by=None, distributed_by=None):
        """
        Create a table in the Doris database.
        
        Args:
            table_name: Name of the table to create
            columns: List of column definitions
            partition_by: Partition clause
            distributed_by: Distribution clause
            
        Returns:
            True if successful, False otherwise
        """
        column_defs = ", ".join(columns)
        query = f"CREATE TABLE IF NOT EXISTS {table_name} ({column_defs})"
        
        if partition_by:
            query += f" PARTITION BY {partition_by}"
        
        if distributed_by:
            query += f" DISTRIBUTED BY {distributed_by}"
        
        try:
            with connections[self.connection_name].cursor() as cursor:
                cursor.execute(query)
            return True
        except Exception as e:
            logger.error(f"Error creating table {table_name}: {str(e)}")
            return False
    
    def drop_table(self, table_name):
        """
        Drop a table from the Doris database.
        
        Args:
            table_name: Name of the table to drop
            
        Returns:
            True if successful, False otherwise
        """
        query = f"DROP TABLE IF EXISTS {table_name}"
        
        try:
            with connections[self.connection_name].cursor() as cursor:
                cursor.execute(query)
            return True
        except Exception as e:
            logger.error(f"Error dropping table {table_name}: {str(e)}")
            return False
    
    def get_table_schema(self, table_name):
        """
        Get the schema of a table in the Doris database.
        
        Args:
            table_name: Name of the table
            
        Returns:
            Table schema as a list of dictionaries
        """
        query = f"DESCRIBE {table_name}"
        
        try:
            return self.execute_query(query)
        except Exception as e:
            logger.error(f"Error getting schema for table {table_name}: {str(e)}")
            return []
    
    def get_tables(self):
        """
        Get a list of tables in the Doris database.
        
        Returns:
            List of table names
        """
        query = "SHOW TABLES"
        
        try:
            result = self.execute_query(query)
            return [list(row.values())[0] for row in result]
        except Exception as e:
            logger.error(f"Error getting tables: {str(e)}")
            return []
    
    def bulk_load(self, table_name, data, columns=None):
        """
        Bulk load data into a Doris table.
        
        Args:
            table_name: Name of the table
            data: List of tuples containing the data to load
            columns: List of column names
            
        Returns:
            Number of rows loaded
        """
        if not data:
            return 0
        
        column_clause = f"({', '.join(columns)})" if columns else ""
        placeholders = ", ".join(["%s"] * len(data[0]))
        query = f"INSERT INTO {table_name} {column_clause} VALUES ({placeholders})"
        
        try:
            with connections[self.connection_name].cursor() as cursor:
                cursor.executemany(query, data)
                return cursor.rowcount
        except Exception as e:
            logger.error(f"Error bulk loading data into {table_name}: {str(e)}")
            return 0
    
    def get_connection_info(self):
        """
        Get information about the current database connection.
        
        Returns:
            Dictionary with connection information
        """
        db_settings = settings.DATABASES.get(self.connection_name, {})
        return {
            'engine': db_settings.get('ENGINE', ''),
            'name': db_settings.get('NAME', ''),
            'host': db_settings.get('HOST', ''),
            'port': db_settings.get('PORT', ''),
            'user': db_settings.get('USER', ''),
        }


def get_doris_client(connection_name='default'):
    """
    Get a Doris client instance.
    
    Args:
        connection_name: Name of the Django database connection to use.
                        Defaults to 'default'.
    
    Returns:
        DorisClient instance
    """
    return DorisClient(connection_name)

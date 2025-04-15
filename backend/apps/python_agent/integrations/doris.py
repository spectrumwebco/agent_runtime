"""
Apache Doris integration for Django.

This module provides integration with Apache Doris, an enterprise-grade
MPP (Massively Parallel Processing) analytical database system.
"""

import logging
import os
import time
from typing import Dict, Any, List, Optional, Union
from django.conf import settings
from pydantic import BaseModel, Field

logger = logging.getLogger(__name__)

class DorisConfig(BaseModel):
    """Apache Doris configuration."""
    
    host: str = Field(default="localhost")
    http_port: int = Field(default=8030)
    query_port: int = Field(default=9030)
    username: str = Field(default="root")
    password: str = Field(default="")
    database: str = Field(default="agent_runtime")
    use_mock: bool = Field(default=False)

class DorisQueryResult(BaseModel):
    """Result of a Doris query."""
    
    column_names: List[str] = Field(default_factory=list)
    rows: List[List[Any]] = Field(default_factory=list)
    row_count: int = Field(default=0)
    execution_time_ms: int = Field(default=0)

class DorisClient:
    """Client for Apache Doris."""
    
    def __init__(self, config: Optional[DorisConfig] = None):
        """Initialize the Doris client."""
        self.config = config or self._get_default_config()
        self._connection = None
        self._cursor = None
        
        try:
            import pymysql
            self._pymysql = pymysql
            self._use_mock = self.config.use_mock
        except ImportError:
            logger.warning("pymysql not installed. Install with: pip install pymysql")
            self._pymysql = None
            self._use_mock = True
        
        if self._use_mock:
            logger.warning("Doris client running in mock mode")
    
    def _get_default_config(self) -> DorisConfig:
        """Get the default configuration from settings."""
        is_kubernetes = os.path.exists('/var/run/secrets/kubernetes.io/serviceaccount/token')
        
        if is_kubernetes:
            host = "doris-fe.default.svc.cluster.local"
        else:
            host = "localhost"
        
        if os.environ.get('CI') == 'true':
            use_mock = True
        else:
            import socket
            try:
                socket.create_connection((host, 9030), timeout=1)
                use_mock = False
            except (socket.timeout, socket.error):
                use_mock = True
        
        return DorisConfig(
            host=host,
            http_port=8030,
            query_port=9030,
            username=getattr(settings, 'DORIS_USERNAME', 'root'),
            password=getattr(settings, 'DORIS_PASSWORD', ''),
            database=getattr(settings, 'DORIS_DATABASE', 'agent_runtime'),
            use_mock=use_mock
        )
    
    def connect(self) -> bool:
        """Connect to Doris."""
        if self._use_mock:
            logger.info("Using mock connection to Doris")
            return True
        
        try:
            self._connection = self._pymysql.connect(
                host=self.config.host,
                port=self.config.query_port,
                user=self.config.username,
                password=self.config.password,
                database=self.config.database,
                connect_timeout=5
            )
            self._cursor = self._connection.cursor()
            logger.info(f"Connected to Doris at {self.config.host}:{self.config.query_port}")
            return True
        except Exception as e:
            logger.error(f"Failed to connect to Doris: {e}")
            return False
    
    def disconnect(self) -> None:
        """Disconnect from Doris."""
        if self._use_mock:
            return
        
        if self._cursor:
            self._cursor.close()
            self._cursor = None
        
        if self._connection:
            self._connection.close()
            self._connection = None
    
    def execute_query(self, query: str, params: Optional[Dict[str, Any]] = None) -> DorisQueryResult:
        """Execute a query on Doris."""
        if self._use_mock:
            return self._mock_execute_query(query, params)
        
        if not self._connection:
            self.connect()
        
        start_time = time.time()
        try:
            self._cursor.execute(query, params or {})
            rows = self._cursor.fetchall()
            column_names = [desc[0] for desc in self._cursor.description] if self._cursor.description else []
            
            result = DorisQueryResult(
                column_names=column_names,
                rows=rows,
                row_count=len(rows),
                execution_time_ms=int((time.time() - start_time) * 1000)
            )
            return result
        except Exception as e:
            logger.error(f"Error executing query: {e}")
            logger.error(f"Query: {query}")
            logger.error(f"Params: {params}")
            raise
    
    def _mock_execute_query(self, query: str, params: Optional[Dict[str, Any]] = None) -> DorisQueryResult:
        """Mock query execution for testing."""
        logger.info(f"Mock executing query: {query}")
        logger.info(f"Params: {params}")
        
        if "SELECT" in query.upper() and "VERSION()" in query.upper():
            return DorisQueryResult(
                column_names=["version"],
                rows=[["Apache Doris 2.0.2 (Mock)"]],
                row_count=1,
                execution_time_ms=5
            )
        elif "SHOW DATABASES" in query.upper():
            return DorisQueryResult(
                column_names=["Database"],
                rows=[["agent_runtime"], ["agent_db"], ["trajectory_db"], ["ml_db"]],
                row_count=4,
                execution_time_ms=10
            )
        elif "SHOW TABLES" in query.upper():
            return DorisQueryResult(
                column_names=["Tables"],
                rows=[["agents"], ["configs"], ["trajectories"], ["ml_models"]],
                row_count=4,
                execution_time_ms=15
            )
        else:
            return DorisQueryResult(
                column_names=["id", "name", "value"],
                rows=[[1, "mock_data", "mock_value"]],
                row_count=1,
                execution_time_ms=20
            )
    
    def check_connection(self) -> Dict[str, Any]:
        """Check the connection to Doris."""
        if self._use_mock:
            return {
                "connected": False,
                "mocked": True,
                "version": "Apache Doris 2.0.2 (Mock)",
                "message": "Running in mock mode"
            }
        
        try:
            if not self._connection:
                self.connect()
            
            self._cursor.execute("SELECT VERSION()")
            version = self._cursor.fetchone()[0]
            
            return {
                "connected": True,
                "mocked": False,
                "version": version,
                "message": "Successfully connected to Doris"
            }
        except Exception as e:
            logger.error(f"Error checking Doris connection: {e}")
            return {
                "connected": False,
                "mocked": False,
                "version": None,
                "message": f"Connection error: {str(e)}"
            }
    
    def create_database(self, database_name: str) -> bool:
        """Create a database in Doris."""
        if self._use_mock:
            logger.info(f"Mock creating database: {database_name}")
            return True
        
        try:
            if not self._connection:
                self.connect()
            
            self._cursor.execute(f"CREATE DATABASE IF NOT EXISTS `{database_name}`")
            self._connection.commit()
            logger.info(f"Created database: {database_name}")
            return True
        except Exception as e:
            logger.error(f"Error creating database {database_name}: {e}")
            return False
    
    def create_table(self, table_name: str, schema: Dict[str, str], 
                     partition_by: Optional[str] = None, 
                     distributed_by: Optional[List[str]] = None) -> bool:
        """Create a table in Doris."""
        if self._use_mock:
            logger.info(f"Mock creating table: {table_name}")
            logger.info(f"Schema: {schema}")
            logger.info(f"Partition by: {partition_by}")
            logger.info(f"Distributed by: {distributed_by}")
            return True
        
        try:
            if not self._connection:
                self.connect()
            
            columns = []
            for column_name, column_type in schema.items():
                columns.append(f"`{column_name}` {column_type}")
            
            create_table_sql = f"CREATE TABLE IF NOT EXISTS `{self.config.database}`.`{table_name}` (\n"
            create_table_sql += ",\n".join(columns)
            create_table_sql += "\n)"
            
            if partition_by:
                create_table_sql += f"\nPARTITION BY {partition_by}"
            
            if distributed_by:
                create_table_sql += f"\nDISTRIBUTED BY HASH({', '.join(distributed_by)}) BUCKETS 10"
            
            self._cursor.execute(create_table_sql)
            self._connection.commit()
            logger.info(f"Created table: {table_name}")
            return True
        except Exception as e:
            logger.error(f"Error creating table {table_name}: {e}")
            logger.error(f"SQL: {create_table_sql}")
            return False
    
    def insert_data(self, table_name: str, data: List[Dict[str, Any]]) -> bool:
        """Insert data into a Doris table."""
        if not data:
            logger.warning(f"No data to insert into {table_name}")
            return True
        
        if self._use_mock:
            logger.info(f"Mock inserting {len(data)} rows into {table_name}")
            return True
        
        try:
            if not self._connection:
                self.connect()
            
            columns = list(data[0].keys())
            
            insert_sql = f"INSERT INTO `{self.config.database}`.`{table_name}` ("
            insert_sql += ", ".join([f"`{col}`" for col in columns])
            insert_sql += ") VALUES "
            
            values = []
            for row in data:
                row_values = []
                for col in columns:
                    if col in row:
                        row_values.append(f"%({col})s")
                    else:
                        row_values.append("NULL")
                values.append(f"({', '.join(row_values)})")
            
            insert_sql += ", ".join(values)
            
            for row in data:
                self._cursor.execute(insert_sql, row)
            
            self._connection.commit()
            logger.info(f"Inserted {len(data)} rows into {table_name}")
            return True
        except Exception as e:
            logger.error(f"Error inserting data into {table_name}: {e}")
            return False
    
    def __enter__(self):
        """Context manager entry."""
        self.connect()
        return self
    
    def __exit__(self, exc_type, exc_val, exc_tb):
        """Context manager exit."""
        self.disconnect()

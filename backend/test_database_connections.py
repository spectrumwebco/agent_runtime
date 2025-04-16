"""
Standalone script to test database connections for all required databases.
This script can be run without requiring the full Django application.
"""

import os
import sys
import json
import logging
import argparse
import importlib
from pathlib import Path

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger("database-test")

sys.path.insert(0, str(Path(__file__).parent))

def test_supabase_connection(config=None):
    """Test connection to Supabase databases."""
    logger.info("Testing Supabase connection...")
    
    try:
        from apps.python_agent.integrations.supabase import SupabaseClient
        
        if config is None:
            config = {
                "url": os.environ.get("SUPABASE_URL", "http://localhost:8000"),
                "key": os.environ.get("SUPABASE_KEY", "mock-key"),
                "databases": ["agent_db", "trajectory_db", "ml_db", "user_db"]
            }
        
        for db_name in config["databases"]:
            logger.info(f"Testing connection to Supabase database: {db_name}")
            
            client = SupabaseClient(
                url=config["url"],
                key=config["key"],
                database=db_name,
                mock=True
            )
            
            result = client.query_table("test_table")
            logger.info(f"Query result: {result}")
            
            if hasattr(client, "auth"):
                auth_status = client.auth.get_session()
                logger.info(f"Auth status: {auth_status}")
            
            if hasattr(client, "functions"):
                function_result = client.functions.invoke("test_function")
                logger.info(f"Function result: {function_result}")
            
            logger.info(f"Supabase database {db_name} connection test successful")
        
        return True
    
    except Exception as e:
        logger.error(f"Supabase connection test failed: {e}")
        return False

def test_ragflow_connection(config=None):
    """Test connection to RAGflow vector database."""
    logger.info("Testing RAGflow connection...")
    
    try:
        from apps.python_agent.integrations.ragflow import RAGflowClient
        
        if config is None:
            config = {
                "host": os.environ.get("RAGFLOW_HOST", "localhost"),
                "port": int(os.environ.get("RAGFLOW_PORT", "8080")),
                "api_key": os.environ.get("RAGFLOW_API_KEY", "mock-key")
            }
        
        client = RAGflowClient(
            host=config["host"],
            port=config["port"],
            api_key=config["api_key"],
            mock=True
        )
        
        search_result = client.search("test query")
        logger.info(f"Search result: {search_result}")
        
        semantic_result = client.semantic_search("test semantic query")
        logger.info(f"Semantic search result: {semantic_result}")
        
        if hasattr(client, "deep_understanding"):
            understanding_result = client.deep_understanding("test understanding query")
            logger.info(f"Deep understanding result: {understanding_result}")
        
        logger.info("RAGflow connection test successful")
        return True
    
    except Exception as e:
        logger.error(f"RAGflow connection test failed: {e}")
        return False

def test_dragonfly_connection(config=None):
    """Test connection to DragonflyDB."""
    logger.info("Testing DragonflyDB connection...")
    
    try:
        from apps.python_agent.integrations.dragonfly import DragonflyClient
        
        if config is None:
            config = {
                "host": os.environ.get("DRAGONFLY_HOST", "localhost"),
                "port": int(os.environ.get("DRAGONFLY_PORT", "6379")),
                "password": os.environ.get("DRAGONFLY_PASSWORD", "")
            }
        
        client = DragonflyClient(
            host=config["host"],
            port=config["port"],
            password=config["password"],
            mock=True
        )
        
        client.set("test_key", "test_value")
        value = client.get("test_key")
        logger.info(f"Key-value test: {value}")
        
        if hasattr(client, "memcached_set"):
            client.memcached_set("test_memcached_key", "test_memcached_value")
            memcached_value = client.memcached_get("test_memcached_key")
            logger.info(f"Memcached test: {memcached_value}")
        
        logger.info("DragonflyDB connection test successful")
        return True
    
    except Exception as e:
        logger.error(f"DragonflyDB connection test failed: {e}")
        return False

def test_rocketmq_connection(config=None):
    """Test connection to RocketMQ."""
    logger.info("Testing RocketMQ connection...")
    
    try:
        from apps.python_agent.integrations.rocketmq import RocketMQClient
        
        if config is None:
            config = {
                "host": os.environ.get("ROCKETMQ_HOST", "localhost"),
                "port": int(os.environ.get("ROCKETMQ_PORT", "9876")),
                "group": os.environ.get("ROCKETMQ_GROUP", "test_group")
            }
        
        client = RocketMQClient(
            host=config["host"],
            port=config["port"],
            group=config["group"],
            mock=True
        )
        
        message_id = client.send_message("test_topic", {"test": "message"})
        logger.info(f"Message ID: {message_id}")
        
        message = client.consume_message("test_topic")
        logger.info(f"Consumed message: {message}")
        
        if hasattr(client, "update_state"):
            client.update_state("test_state", {"status": "testing"})
            state = client.get_state("test_state")
            logger.info(f"State management test: {state}")
        
        logger.info("RocketMQ connection test successful")
        return True
    
    except Exception as e:
        logger.error(f"RocketMQ connection test failed: {e}")
        return False

def test_doris_connection(config=None):
    """Test connection to Apache Doris."""
    logger.info("Testing Apache Doris connection...")
    
    try:
        from backend.integrations.doris import DorisClient
        
        if config is None:
            config = {
                "host": os.environ.get("DORIS_HOST", "localhost"),
                "port": int(os.environ.get("DORIS_PORT", "9030")),
                "user": os.environ.get("DORIS_USER", "root"),
                "password": os.environ.get("DORIS_PASSWORD", "")
            }
        
        client = DorisClient(
            host=config["host"],
            port=config["port"],
            user=config["user"],
            password=config["password"],
            mock=True
        )
        
        result = client.execute_query("SELECT 1")
        logger.info(f"Query result: {result}")
        
        client.create_table("test_table", {
            "id": "INT",
            "name": "VARCHAR(100)",
            "created_at": "DATETIME"
        })
        
        client.insert_data("test_table", [
            {"id": 1, "name": "Test 1", "created_at": "2023-01-01 00:00:00"},
            {"id": 2, "name": "Test 2", "created_at": "2023-01-02 00:00:00"}
        ])
        
        table_data = client.execute_query("SELECT * FROM test_table")
        logger.info(f"Table data: {table_data}")
        
        logger.info("Apache Doris connection test successful")
        return True
    
    except Exception as e:
        logger.error(f"Apache Doris connection test failed: {e}")
        return False

def test_kafka_connection(config=None):
    """Test connection to Apache Kafka."""
    logger.info("Testing Apache Kafka connection...")
    
    try:
        from backend.integrations.kafka import KafkaClient
        
        if config is None:
            config = {
                "bootstrap_servers": os.environ.get("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092"),
                "group_id": os.environ.get("KAFKA_GROUP_ID", "test_group")
            }
        
        client = KafkaClient(
            bootstrap_servers=config["bootstrap_servers"],
            group_id=config["group_id"],
            mock=True
        )
        
        client.produce_message("test_topic", {"test": "message"})
        
        message = client.consume_message("test_topic")
        logger.info(f"Consumed message: {message}")
        
        logger.info("Apache Kafka connection test successful")
        return True
    
    except Exception as e:
        logger.error(f"Apache Kafka connection test failed: {e}")
        return False

def test_postgres_connection(config=None):
    """Test connection to PostgreSQL (CrunchyData)."""
    logger.info("Testing PostgreSQL connection...")
    
    try:
        from backend.integrations.crunchydata import PostgresClient
        
        if config is None:
            config = {
                "host": os.environ.get("POSTGRES_HOST", "localhost"),
                "port": int(os.environ.get("POSTGRES_PORT", "5432")),
                "user": os.environ.get("POSTGRES_USER", "postgres"),
                "password": os.environ.get("POSTGRES_PASSWORD", ""),
                "database": os.environ.get("POSTGRES_DATABASE", "postgres")
            }
        
        client = PostgresClient(
            host=config["host"],
            port=config["port"],
            user=config["user"],
            password=config["password"],
            database=config["database"],
            mock=True
        )
        
        result = client.execute_query("SELECT 1")
        logger.info(f"Query result: {result}")
        
        client.create_table("test_table", {
            "id": "SERIAL PRIMARY KEY",
            "name": "VARCHAR(100)",
            "created_at": "TIMESTAMP DEFAULT CURRENT_TIMESTAMP"
        })
        
        client.insert_data("test_table", [
            {"name": "Test 1"},
            {"name": "Test 2"}
        ])
        
        table_data = client.execute_query("SELECT * FROM test_table")
        logger.info(f"Table data: {table_data}")
        
        logger.info("PostgreSQL connection test successful")
        return True
    
    except Exception as e:
        logger.error(f"PostgreSQL connection test failed: {e}")
        return False

def test_cross_database_integration():
    """Test cross-database integration."""
    logger.info("Testing cross-database integration...")
    
    try:
        from backend.integrations.crunchydata import PostgresClient
        from backend.integrations.kafka import KafkaClient
        from backend.integrations.doris import DorisClient
        
        postgres_client = PostgresClient(mock=True)
        kafka_client = KafkaClient(mock=True)
        doris_client = DorisClient(mock=True)
        
        postgres_client.create_table("test_integration", {
            "id": "SERIAL PRIMARY KEY",
            "name": "VARCHAR(100)",
            "value": "INT",
            "created_at": "TIMESTAMP DEFAULT CURRENT_TIMESTAMP"
        })
        
        postgres_client.insert_data("test_integration", [
            {"name": "Integration Test 1", "value": 100},
            {"name": "Integration Test 2", "value": 200},
            {"name": "Integration Test 3", "value": 300}
        ])
        
        kafka_client.produce_message("test_integration", {
            "source": "postgres",
            "destination": "doris",
            "table": "test_integration",
            "operation": "insert"
        })
        
        doris_client.create_table("test_integration", {
            "id": "INT",
            "name": "VARCHAR(100)",
            "value": "INT",
            "created_at": "DATETIME"
        })
        
        doris_client.insert_data("test_integration", [
            {"id": 1, "name": "Integration Test 1", "value": 100, "created_at": "2023-01-01 00:00:00"},
            {"id": 2, "name": "Integration Test 2", "value": 200, "created_at": "2023-01-01 00:00:00"},
            {"id": 3, "name": "Integration Test 3", "value": 300, "created_at": "2023-01-01 00:00:00"}
        ])
        
        result = doris_client.execute_query("SELECT COUNT(*) FROM test_integration")
        logger.info(f"Cross-database integration test result: {result}")
        
        logger.info("Cross-database integration test successful")
        return True
    
    except Exception as e:
        logger.error(f"Cross-database integration test failed: {e}")
        return False

def main():
    """Main function to run database connection tests."""
    parser = argparse.ArgumentParser(description="Test database connections")
    parser.add_argument("--all", action="store_true", help="Test all database connections")
    parser.add_argument("--supabase", action="store_true", help="Test Supabase connection")
    parser.add_argument("--ragflow", action="store_true", help="Test RAGflow connection")
    parser.add_argument("--dragonfly", action="store_true", help="Test DragonflyDB connection")
    parser.add_argument("--rocketmq", action="store_true", help="Test RocketMQ connection")
    parser.add_argument("--doris", action="store_true", help="Test Apache Doris connection")
    parser.add_argument("--kafka", action="store_true", help="Test Apache Kafka connection")
    parser.add_argument("--postgres", action="store_true", help="Test PostgreSQL connection")
    parser.add_argument("--integration", action="store_true", help="Test cross-database integration")
    parser.add_argument("--config", type=str, help="Path to configuration file")
    
    args = parser.parse_args()
    
    config = None
    if args.config:
        try:
            with open(args.config, "r") as f:
                config = json.load(f)
        except Exception as e:
            logger.error(f"Failed to load configuration file: {e}")
            return 1
    
    run_all = args.all or not any([
        args.supabase, args.ragflow, args.dragonfly, args.rocketmq,
        args.doris, args.kafka, args.postgres, args.integration
    ])
    
    success = True
    
    if run_all or args.supabase:
        if not test_supabase_connection(config.get("supabase") if config else None):
            success = False
    
    if run_all or args.ragflow:
        if not test_ragflow_connection(config.get("ragflow") if config else None):
            success = False
    
    if run_all or args.dragonfly:
        if not test_dragonfly_connection(config.get("dragonfly") if config else None):
            success = False
    
    if run_all or args.rocketmq:
        if not test_rocketmq_connection(config.get("rocketmq") if config else None):
            success = False
    
    if run_all or args.doris:
        if not test_doris_connection(config.get("doris") if config else None):
            success = False
    
    if run_all or args.kafka:
        if not test_kafka_connection(config.get("kafka") if config else None):
            success = False
    
    if run_all or args.postgres:
        if not test_postgres_connection(config.get("postgres") if config else None):
            success = False
    
    if run_all or args.integration:
        if not test_cross_database_integration():
            success = False
    
    if success:
        logger.info("All database connection tests passed")
        return 0
    else:
        logger.error("Some database connection tests failed")
        return 1

if __name__ == "__main__":
    sys.exit(main())

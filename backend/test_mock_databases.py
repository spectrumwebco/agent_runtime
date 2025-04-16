"""
Standalone script to test database integrations with mock databases.
This script can be run outside of Django to verify basic functionality.
"""

import os
import sys
import logging
import json
from pathlib import Path

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger("database-test")

backend_dir = Path(__file__).resolve().parent
if str(backend_dir) not in sys.path:
    sys.path.insert(0, str(backend_dir))

try:
    from apps.python_agent.integrations.mock_db import (
        MockSupabaseClient,
        MockRAGflowClient,
        MockDragonflyClient,
        MockRocketMQClient,
        MockDorisClient,
        MockKafkaClient,
        MockPostgresClient
    )
    logger.info("Successfully imported mock database clients")
except ImportError as e:
    logger.error(f"Failed to import mock database clients: {e}")
    sys.exit(1)

def test_supabase():
    """Test Supabase integration with mock database."""
    logger.info("Testing Supabase integration...")
    
    try:
        client = MockSupabaseClient(url="mock://supabase", key="mock-key")
        
        logger.info("Testing query operations...")
        records = client.query_table('test_table')
        logger.info(f"Query result: {records}")
        
        logger.info("Testing insert operations...")
        record = client.insert_record('test_table', {'name': 'Test Record'})
        logger.info(f"Insert result: {record}")
        
        logger.info("Supabase integration test passed")
        return True
        
    except Exception as e:
        logger.error(f"Supabase integration test failed: {e}")
        return False

def test_ragflow():
    """Test RAGflow integration with mock database."""
    logger.info("Testing RAGflow integration...")
    
    try:
        client = MockRAGflowClient(host="localhost", port=8000, api_key="mock-key")
        
        logger.info("Testing search functionality...")
        results = client.search("test query")
        logger.info(f"Search result: {results}")
        
        logger.info("Testing semantic search functionality...")
        results = client.semantic_search("test query")
        logger.info(f"Semantic search result: {results}")
        
        logger.info("RAGflow integration test passed")
        return True
        
    except Exception as e:
        logger.error(f"RAGflow integration test failed: {e}")
        return False

def test_dragonfly():
    """Test DragonflyDB integration with mock database."""
    logger.info("Testing DragonflyDB integration...")
    
    try:
        client = MockDragonflyClient(host="localhost", port=6379, mock=True)
        
        logger.info("Testing key-value operations...")
        client.set("test_key", "test_value")
        value = client.get("test_key")
        logger.info(f"Get result: {value}")
        
        logger.info("Testing memcached operations...")
        client.memcached_set("test_memcached_key", "test_memcached_value")
        value = client.memcached_get("test_memcached_key")
        logger.info(f"Memcached get result: {value}")
        
        logger.info("DragonflyDB integration test passed")
        return True
        
    except Exception as e:
        logger.error(f"DragonflyDB integration test failed: {e}")
        return False

def test_rocketmq():
    """Test RocketMQ integration with mock database."""
    logger.info("Testing RocketMQ integration...")
    
    try:
        client = MockRocketMQClient(host="localhost", port=9876, mock=True)
        
        logger.info("Testing message production and consumption...")
        message_id = client.send_message("test_topic", {"test": "message"})
        logger.info(f"Message ID: {message_id}")
        
        message = client.consume_message("test_topic")
        logger.info(f"Consumed message: {message}")
        
        logger.info("Testing state management...")
        client.update_state("test_state", {"status": "testing"})
        state = client.get_state("test_state")
        logger.info(f"State: {state}")
        
        logger.info("RocketMQ integration test passed")
        return True
        
    except Exception as e:
        logger.error(f"RocketMQ integration test failed: {e}")
        return False

def test_doris():
    """Test Apache Doris integration with mock database."""
    logger.info("Testing Apache Doris integration...")
    
    try:
        client = MockDorisClient(connection_params={
            "host": "localhost",
            "port": 9030,
            "user": "root",
            "password": ""
        }, mock=True)
        
        logger.info("Testing query execution...")
        results = client.execute_query("SELECT 1")
        logger.info(f"Query result: {results}")
        
        logger.info("Testing table operations...")
        client.create_table("test_table", {
            "id": "INT",
            "name": "VARCHAR(100)",
            "created_at": "DATETIME"
        })
        
        client.insert_data("test_table", [
            {"id": 1, "name": "Test 1", "created_at": "2023-01-01 00:00:00"},
            {"id": 2, "name": "Test 2", "created_at": "2023-01-02 00:00:00"}
        ])
        
        results = client.execute_query("SELECT * FROM test_table")
        logger.info(f"Table data: {results}")
        
        logger.info("Apache Doris integration test passed")
        return True
        
    except Exception as e:
        logger.error(f"Apache Doris integration test failed: {e}")
        return False

def test_postgres():
    """Test PostgreSQL integration with mock database."""
    logger.info("Testing PostgreSQL integration...")
    
    try:
        client = MockPostgresClient(host="localhost", port=5432, user="postgres", password="", mock=True)
        
        logger.info("Testing query execution...")
        results = client.execute_query("SELECT 1")
        logger.info(f"Query result: {results}")
        
        logger.info("Testing table operations...")
        client.create_table("test_table", {
            "id": "SERIAL PRIMARY KEY",
            "name": "VARCHAR(100)",
            "created_at": "TIMESTAMP DEFAULT CURRENT_TIMESTAMP"
        })
        
        client.insert_data("test_table", [
            {"name": "Test 1"},
            {"name": "Test 2"}
        ])
        
        results = client.execute_query("SELECT * FROM test_table")
        logger.info(f"Table data: {results}")
        
        logger.info("PostgreSQL integration test passed")
        return True
        
    except Exception as e:
        logger.error(f"PostgreSQL integration test failed: {e}")
        return False

def test_kafka():
    """Test Kafka integration with mock database."""
    logger.info("Testing Kafka integration...")
    
    try:
        client = MockKafkaClient(bootstrap_servers="localhost:9092", mock=True)
        
        logger.info("Testing message production and consumption...")
        client.produce_message("test_topic", {"test": "message"})
        
        message = client.consume_message("test_topic")
        logger.info(f"Consumed message: {message}")
        
        logger.info("Kafka integration test passed")
        return True
        
    except Exception as e:
        logger.error(f"Kafka integration test failed: {e}")
        return False

def main():
    """Run all database integration tests."""
    logger.info("Starting database integration tests")
    
    success = True
    
    if not test_supabase():
        success = False
    
    if not test_ragflow():
        success = False
    
    if not test_dragonfly():
        success = False
    
    if not test_rocketmq():
        success = False
    
    if not test_doris():
        success = False
    
    if not test_postgres():
        success = False
    
    if not test_kafka():
        success = False
    
    if success:
        logger.info("All database integration tests passed")
        return 0
    else:
        logger.error("Some database integration tests failed")
        return 1

if __name__ == "__main__":
    sys.exit(main())

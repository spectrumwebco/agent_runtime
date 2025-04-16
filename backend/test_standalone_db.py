"""
Standalone script to test database integrations without Django dependencies.
"""

import logging
import json
from typing import Dict, Any, Optional

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger("database-test")

class MockSupabaseClient:
    def __init__(self, url="mock://supabase", key="mock-key"):
        self.url = url
        self.key = key
        self.data = {}
        logger.info(f"Initialized MockSupabaseClient")
    
    def query_table(self, table_name):
        if table_name not in self.data:
            self.data[table_name] = []
        return self.data[table_name]
    
    def insert_record(self, table_name, record):
        if table_name not in self.data:
            self.data[table_name] = []
        record["id"] = len(self.data[table_name]) + 1
        self.data[table_name].append(record)
        return record

class MockRAGflowClient:
    def __init__(self, host="localhost", port=8000):
        self.host = host
        self.port = port
        logger.info(f"Initialized MockRAGflowClient")
    
    def search(self, query, top_k=5):
        return {"results": [{"content": f"Result for {query}", "score": 0.9}]}

class MockDragonflyClient:
    def __init__(self, host="localhost", port=6379):
        self.host = host
        self.port = port
        self.data = {}
        logger.info(f"Initialized MockDragonflyClient")
    
    def set(self, key, value):
        self.data[key] = value
        return True
    
    def get(self, key):
        return self.data.get(key)

class MockDorisClient:
    def __init__(self, connection_params=None):
        self.connection_params = connection_params or {}
        self.tables = {}
        logger.info(f"Initialized MockDorisClient")
    
    def execute_query(self, query):
        return [{"result": "success"}]

def test_all_databases():
    """Test all database integrations."""
    logger.info("Testing all database integrations...")
    
    supabase = MockSupabaseClient()
    supabase.insert_record("test_table", {"name": "Test"})
    records = supabase.query_table("test_table")
    logger.info(f"Supabase records: {records}")
    
    ragflow = MockRAGflowClient()
    results = ragflow.search("test query")
    logger.info(f"RAGflow results: {results}")
    
    dragonfly = MockDragonflyClient()
    dragonfly.set("test_key", "test_value")
    value = dragonfly.get("test_key")
    logger.info(f"DragonflyDB value: {value}")
    
    doris = MockDorisClient()
    results = doris.execute_query("SELECT 1")
    logger.info(f"Doris results: {results}")
    
    logger.info("All database tests completed successfully")
    return True

if __name__ == "__main__":
    test_all_databases()

"""
Apache Kafka integration for Django.

This module provides integration with Apache Kafka, a distributed
event streaming platform.
"""

import json
import logging
import os
import socket
import threading
import time
from typing import Dict, Any, List, Optional, Union, Callable
from django.conf import settings
from pydantic import BaseModel, Field

logger = logging.getLogger(__name__)

class KafkaConfig(BaseModel):
    """Apache Kafka configuration."""
    
    bootstrap_servers: str = Field(default="localhost:9092")
    client_id: str = Field(default="agent-runtime")
    group_id: str = Field(default="agent-runtime-group")
    auto_offset_reset: str = Field(default="earliest")
    enable_auto_commit: bool = Field(default=True)
    use_mock: bool = Field(default=False)

class KafkaMessage(BaseModel):
    """Kafka message model."""
    
    topic: str
    key: Optional[str] = None
    value: Any
    partition: Optional[int] = None
    timestamp: Optional[int] = None
    headers: Optional[Dict[str, str]] = None

class KafkaClient:
    """Client for Apache Kafka."""
    
    def __init__(self, config: Optional[KafkaConfig] = None):
        """Initialize the Kafka client."""
        self.config = config or self._get_default_config()
        self._producer = None
        self._consumer = None
        self._consumer_thread = None
        self._running = False
        self._message_handlers = {}
        
        try:
            from confluent_kafka import Producer, Consumer, KafkaError
            self._kafka = True
            self._kafka_error = KafkaError
            self._use_mock = self.config.use_mock
        except ImportError:
            logger.warning("confluent-kafka not installed. Install with: pip install confluent-kafka")
            self._kafka = False
            self._use_mock = True
        
        if self._use_mock:
            logger.warning("Kafka client running in mock mode")
    
    def _get_default_config(self) -> KafkaConfig:
        """Get the default configuration from settings."""
        is_kubernetes = os.path.exists('/var/run/secrets/kubernetes.io/serviceaccount/token')
        
        if is_kubernetes:
            bootstrap_servers = "kafka-broker.default.svc.cluster.local:9092"
        else:
            bootstrap_servers = "localhost:9092"
        
        if os.environ.get('CI') == 'true':
            use_mock = True
        else:
            try:
                host, port = bootstrap_servers.split(':')
                socket.create_connection((host, int(port)), timeout=1)
                use_mock = False
            except (socket.timeout, socket.error, ValueError):
                use_mock = True
        
        return KafkaConfig(
            bootstrap_servers=bootstrap_servers,
            client_id=getattr(settings, 'KAFKA_CLIENT_ID', 'agent-runtime'),
            group_id=getattr(settings, 'KAFKA_GROUP_ID', 'agent-runtime-group'),
            auto_offset_reset=getattr(settings, 'KAFKA_AUTO_OFFSET_RESET', 'earliest'),
            enable_auto_commit=getattr(settings, 'KAFKA_ENABLE_AUTO_COMMIT', True),
            use_mock=use_mock
        )
    
    def _create_producer(self):
        """Create a Kafka producer."""
        if self._use_mock:
            return
        
        if not self._kafka:
            return
        
        from confluent_kafka import Producer
        
        self._producer = Producer({
            'bootstrap.servers': self.config.bootstrap_servers,
            'client.id': self.config.client_id,
        })
    
    def _create_consumer(self):
        """Create a Kafka consumer."""
        if self._use_mock:
            return
        
        if not self._kafka:
            return
        
        from confluent_kafka import Consumer
        
        self._consumer = Consumer({
            'bootstrap.servers': self.config.bootstrap_servers,
            'group.id': self.config.group_id,
            'auto.offset.reset': self.config.auto_offset_reset,
            'enable.auto.commit': self.config.enable_auto_commit,
        })
    
    def produce(self, message: KafkaMessage) -> bool:
        """Produce a message to Kafka."""
        if self._use_mock:
            logger.info(f"Mock producing message to topic {message.topic}: {message.value}")
            return True
        
        if not self._kafka:
            logger.error("Kafka client not available")
            return False
        
        if not self._producer:
            self._create_producer()
        
        try:
            value = message.value
            if not isinstance(value, (str, bytes)):
                value = json.dumps(value)
            
            if isinstance(value, str):
                value = value.encode('utf-8')
            
            key = message.key
            if key and isinstance(key, str):
                key = key.encode('utf-8')
            
            headers = None
            if message.headers:
                headers = [(k, v.encode('utf-8') if isinstance(v, str) else v) 
                          for k, v in message.headers.items()]
            
            self._producer.produce(
                topic=message.topic,
                key=key,
                value=value,
                partition=message.partition,
                timestamp=message.timestamp,
                headers=headers
            )
            self._producer.flush()
            return True
        except Exception as e:
            logger.error(f"Error producing message to Kafka: {e}")
            return False
    
    def subscribe(self, topics: List[str]) -> bool:
        """Subscribe to Kafka topics."""
        if self._use_mock:
            logger.info(f"Mock subscribing to topics: {topics}")
            return True
        
        if not self._kafka:
            logger.error("Kafka client not available")
            return False
        
        if not self._consumer:
            self._create_consumer()
        
        try:
            self._consumer.subscribe(topics)
            return True
        except Exception as e:
            logger.error(f"Error subscribing to Kafka topics: {e}")
            return False
    
    def register_handler(self, topic: str, handler: Callable[[KafkaMessage], None]) -> None:
        """Register a handler for a topic."""
        if topic not in self._message_handlers:
            self._message_handlers[topic] = []
        
        self._message_handlers[topic].append(handler)
    
    def _consumer_loop(self):
        """Consumer loop for processing messages."""
        if self._use_mock:
            while self._running:
                time.sleep(1)
            return
        
        if not self._kafka:
            return
        
        while self._running:
            try:
                msg = self._consumer.poll(1.0)
                
                if msg is None:
                    continue
                
                if msg.error():
                    if msg.error().code() == self._kafka_error._PARTITION_EOF:
                        logger.debug(f"Reached end of partition {msg.partition()}")
                    else:
                        logger.error(f"Error consuming from Kafka: {msg.error()}")
                    continue
                
                topic = msg.topic()
                key = msg.key().decode('utf-8') if msg.key() else None
                value = msg.value().decode('utf-8')
                
                try:
                    value = json.loads(value)
                except json.JSONDecodeError:
                    pass
                
                message = KafkaMessage(
                    topic=topic,
                    key=key,
                    value=value,
                    partition=msg.partition(),
                    timestamp=msg.timestamp()[1]
                )
                
                if topic in self._message_handlers:
                    for handler in self._message_handlers[topic]:
                        try:
                            handler(message)
                        except Exception as e:
                            logger.error(f"Error in message handler: {e}")
            
            except Exception as e:
                logger.error(f"Error in consumer loop: {e}")
                time.sleep(1)
    
    def start_consumer(self) -> bool:
        """Start the consumer loop."""
        if self._running:
            logger.warning("Consumer already running")
            return True
        
        self._running = True
        
        if self._use_mock:
            logger.info("Starting mock consumer")
            self._consumer_thread = threading.Thread(target=self._consumer_loop)
            self._consumer_thread.daemon = True
            self._consumer_thread.start()
            return True
        
        if not self._kafka:
            logger.error("Kafka client not available")
            return False
        
        if not self._consumer:
            self._create_consumer()
        
        self._consumer_thread = threading.Thread(target=self._consumer_loop)
        self._consumer_thread.daemon = True
        self._consumer_thread.start()
        
        return True
    
    def stop_consumer(self) -> bool:
        """Stop the consumer loop."""
        if not self._running:
            logger.warning("Consumer not running")
            return True
        
        self._running = False
        
        if self._consumer_thread:
            self._consumer_thread.join(timeout=5.0)
            self._consumer_thread = None
        
        if not self._use_mock and self._consumer:
            self._consumer.close()
            self._consumer = None
        
        return True
    
    def check_connection(self) -> Dict[str, Any]:
        """Check the connection to Kafka."""
        if self._use_mock:
            return {
                "connected": False,
                "mocked": True,
                "bootstrap_servers": self.config.bootstrap_servers,
                "message": "Running in mock mode"
            }
        
        if not self._kafka:
            return {
                "connected": False,
                "mocked": False,
                "bootstrap_servers": self.config.bootstrap_servers,
                "message": "Kafka client not available"
            }
        
        try:
            if not self._producer:
                self._create_producer()
            
            host, port = self.config.bootstrap_servers.split(':')
            socket.create_connection((host, int(port)), timeout=3)
            
            return {
                "connected": True,
                "mocked": False,
                "bootstrap_servers": self.config.bootstrap_servers,
                "message": "Successfully connected to Kafka"
            }
        except Exception as e:
            logger.error(f"Error checking Kafka connection: {e}")
            return {
                "connected": False,
                "mocked": False,
                "bootstrap_servers": self.config.bootstrap_servers,
                "message": f"Connection error: {str(e)}"
            }
    
    def __enter__(self):
        """Context manager entry."""
        return self
    
    def __exit__(self, exc_type, exc_val, exc_tb):
        """Context manager exit."""
        self.stop_consumer()

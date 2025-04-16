"""
Django integration with Apache Kafka.

This module provides integration between Django and Apache Kafka,
implementing event streaming and message processing.
"""

import json
import logging
from typing import Dict, Any, Optional, List, Callable, Union
from django.conf import settings
from confluent_kafka import Producer, Consumer, KafkaError, KafkaException

logger = logging.getLogger(__name__)

class KafkaClient:
    """
    Client for Apache Kafka operations.
    
    This client handles event streaming and message processing
    for Apache Kafka.
    """
    
    def __init__(self, bootstrap_servers=None, client_id=None, group_id=None):
        """
        Initialize the Kafka client.
        
        Args:
            bootstrap_servers: Kafka bootstrap servers
            client_id: Client ID for Kafka
            group_id: Consumer group ID for Kafka
        """
        kafka_config = getattr(settings, 'KAFKA_CONFIG', {})
        self.bootstrap_servers = bootstrap_servers or kafka_config.get('bootstrap_servers', 'localhost:9092')
        self.client_id = client_id or f"django-kafka-{id(self)}"
        self.group_id = group_id or f"django-kafka-group-{id(self)}"
        self.topic_prefix = kafka_config.get('topic_prefix', 'django')
        self._producer = None
        self._consumer = None
    
    @property
    def producer(self):
        """Get or create a Kafka producer."""
        if self._producer is None:
            self._producer = Producer({
                'bootstrap.servers': self.bootstrap_servers,
                'client.id': self.client_id,
            })
        return self._producer
    
    def get_consumer(self, topics=None, group_id=None, auto_offset_reset='earliest'):
        """
        Get a Kafka consumer.
        
        Args:
            topics: List of topics to subscribe to
            group_id: Consumer group ID
            auto_offset_reset: Auto offset reset strategy
            
        Returns:
            Kafka consumer instance
        """
        consumer = Consumer({
            'bootstrap.servers': self.bootstrap_servers,
            'group.id': group_id or self.group_id,
            'auto.offset.reset': auto_offset_reset,
        })
        
        if topics:
            consumer.subscribe(topics)
        
        return consumer
    
    def get_full_topic_name(self, topic):
        """
        Get the full topic name with prefix.
        
        Args:
            topic: Base topic name
            
        Returns:
            Full topic name with prefix
        """
        return f"{self.topic_prefix}.{topic}" if self.topic_prefix else topic
    
    def produce(self, topic, value, key=None, headers=None, callback=None):
        """
        Produce a message to a Kafka topic.
        
        Args:
            topic: Topic to produce to
            value: Message value
            key: Message key
            headers: Message headers
            callback: Delivery callback function
            
        Returns:
            None
        """
        full_topic = self.get_full_topic_name(topic)
        
        if isinstance(value, dict) or isinstance(value, list):
            value = json.dumps(value).encode('utf-8')
        elif not isinstance(value, bytes):
            value = str(value).encode('utf-8')
        
        if key is not None and not isinstance(key, bytes):
            key = str(key).encode('utf-8')
        
        self.producer.produce(
            full_topic,
            value=value,
            key=key,
            headers=headers,
            callback=callback
        )
        self.producer.poll(0)
    
    def flush(self, timeout=10):
        """
        Flush the producer.
        
        Args:
            timeout: Flush timeout in seconds
            
        Returns:
            Number of messages still in queue
        """
        return self.producer.flush(timeout)
    
    def consume(self, topics, timeout=1.0, num_messages=1, group_id=None):
        """
        Consume messages from Kafka topics.
        
        Args:
            topics: List of topics to consume from
            timeout: Consume timeout in seconds
            num_messages: Maximum number of messages to consume
            group_id: Consumer group ID
            
        Returns:
            List of consumed messages
        """
        full_topics = [self.get_full_topic_name(topic) for topic in topics]
        consumer = self.get_consumer(full_topics, group_id)
        
        messages = []
        try:
            for _ in range(num_messages):
                msg = consumer.poll(timeout)
                if msg is None:
                    break
                
                if msg.error():
                    if msg.error().code() == KafkaError._PARTITION_EOF:
                        logger.debug(f"Reached end of partition for topic {msg.topic()}")
                    else:
                        logger.error(f"Error consuming from Kafka: {msg.error()}")
                else:
                    value = msg.value()
                    try:
                        value = json.loads(value.decode('utf-8'))
                    except (json.JSONDecodeError, UnicodeDecodeError):
                        pass
                    
                    messages.append({
                        'topic': msg.topic(),
                        'partition': msg.partition(),
                        'offset': msg.offset(),
                        'key': msg.key(),
                        'value': value,
                        'headers': msg.headers(),
                        'timestamp': msg.timestamp(),
                    })
        finally:
            consumer.close()
        
        return messages
    
    def consume_loop(self, topics, callback, group_id=None, timeout=1.0, exit_condition=None):
        """
        Consume messages in a loop.
        
        Args:
            topics: List of topics to consume from
            callback: Callback function for consumed messages
            group_id: Consumer group ID
            timeout: Consume timeout in seconds
            exit_condition: Function that returns True when the loop should exit
            
        Returns:
            None
        """
        full_topics = [self.get_full_topic_name(topic) for topic in topics]
        consumer = self.get_consumer(full_topics, group_id)
        
        try:
            while True:
                if exit_condition and exit_condition():
                    break
                
                msg = consumer.poll(timeout)
                if msg is None:
                    continue
                
                if msg.error():
                    if msg.error().code() == KafkaError._PARTITION_EOF:
                        logger.debug(f"Reached end of partition for topic {msg.topic()}")
                    else:
                        logger.error(f"Error consuming from Kafka: {msg.error()}")
                else:
                    value = msg.value()
                    try:
                        value = json.loads(value.decode('utf-8'))
                    except (json.JSONDecodeError, UnicodeDecodeError):
                        pass
                    
                    message = {
                        'topic': msg.topic(),
                        'partition': msg.partition(),
                        'offset': msg.offset(),
                        'key': msg.key(),
                        'value': value,
                        'headers': msg.headers(),
                        'timestamp': msg.timestamp(),
                    }
                    
                    callback(message)
                    consumer.commit(msg)
        finally:
            consumer.close()
    
    def create_topic(self, topic, num_partitions=1, replication_factor=1):
        """
        Create a Kafka topic.
        
        Args:
            topic: Topic to create
            num_partitions: Number of partitions
            replication_factor: Replication factor
            
        Returns:
            True if successful, False otherwise
        """
        from confluent_kafka.admin import AdminClient, NewTopic
        
        full_topic = self.get_full_topic_name(topic)
        admin_client = AdminClient({'bootstrap.servers': self.bootstrap_servers})
        
        new_topic = NewTopic(
            full_topic,
            num_partitions=num_partitions,
            replication_factor=replication_factor
        )
        
        try:
            admin_client.create_topics([new_topic])
            return True
        except KafkaException as e:
            logger.error(f"Error creating topic {full_topic}: {str(e)}")
            return False
    
    def delete_topic(self, topic):
        """
        Delete a Kafka topic.
        
        Args:
            topic: Topic to delete
            
        Returns:
            True if successful, False otherwise
        """
        from confluent_kafka.admin import AdminClient
        
        full_topic = self.get_full_topic_name(topic)
        admin_client = AdminClient({'bootstrap.servers': self.bootstrap_servers})
        
        try:
            admin_client.delete_topics([full_topic])
            return True
        except KafkaException as e:
            logger.error(f"Error deleting topic {full_topic}: {str(e)}")
            return False
    
    def list_topics(self):
        """
        List Kafka topics.
        
        Returns:
            Dictionary of topics and their metadata
        """
        from confluent_kafka.admin import AdminClient
        
        admin_client = AdminClient({'bootstrap.servers': self.bootstrap_servers})
        
        try:
            topics = admin_client.list_topics().topics
            return {topic: metadata.to_dict() for topic, metadata in topics.items()}
        except KafkaException as e:
            logger.error(f"Error listing topics: {str(e)}")
            return {}
    
    def close(self):
        """
        Close the Kafka client.
        
        Returns:
            None
        """
        if self._producer:
            self._producer.flush()
            self._producer = None
        
        if self._consumer:
            self._consumer.close()
            self._consumer = None


def get_kafka_client(bootstrap_servers=None, client_id=None, group_id=None):
    """
    Get a Kafka client instance.
    
    Args:
        bootstrap_servers: Kafka bootstrap servers
        client_id: Client ID for Kafka
        group_id: Consumer group ID for Kafka
    
    Returns:
        KafkaClient instance
    """
    return KafkaClient(bootstrap_servers, client_id, group_id)

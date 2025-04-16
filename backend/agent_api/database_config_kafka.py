"""
Kafka configuration for Django.
"""

import os
import logging
import socket
from pathlib import Path

logger = logging.getLogger(__name__)

def is_running_in_kubernetes():
    """Check if we're running in a Kubernetes environment."""
    return os.path.exists('/var/run/secrets/kubernetes.io/serviceaccount/token')

IN_KUBERNETES = is_running_in_kubernetes()

if IN_KUBERNETES:
    KAFKA_CONFIG = {
        'bootstrap_servers': 'kafka-broker.default.svc.cluster.local:9092',
        'client_id': 'agent-runtime',
        'group_id': 'agent-runtime-group',
        'auto_offset_reset': 'earliest',
        'enable_auto_commit': True,
    }
else:
    try:
        socket.create_connection(('localhost', 9092), timeout=1)
        kafka_available = True
    except (socket.timeout, socket.error):
        kafka_available = False
        logger.warning("Apache Kafka not available locally, using mock Kafka for development")
    
    if kafka_available:
        KAFKA_CONFIG = {
            'bootstrap_servers': 'localhost:9092',
            'client_id': 'agent-runtime',
            'group_id': 'agent-runtime-group',
            'auto_offset_reset': 'earliest',
            'enable_auto_commit': True,
        }
    else:
        KAFKA_CONFIG = {
            'bootstrap_servers': 'localhost:9092',
            'client_id': 'agent-runtime',
            'group_id': 'agent-runtime-group',
            'auto_offset_reset': 'earliest',
            'enable_auto_commit': True,
            'use_mock': True,
        }

KAFKA_TOPICS = {
    'agent_events': 'agent-events',
    'agent_commands': 'agent-commands',
    'agent_responses': 'agent-responses',
    'agent_logs': 'agent-logs',
    'trajectory_events': 'trajectory-events',
    'ml_events': 'ml-events',
    'ml_commands': 'ml-commands',
    'ml_responses': 'ml-responses',
    'ml_logs': 'ml-logs',
    'shared_state': 'shared-state',
}

KAFKA_CONSUMER_CONFIG = {
    'bootstrap_servers': KAFKA_CONFIG['bootstrap_servers'],
    'group_id': KAFKA_CONFIG['group_id'],
    'auto_offset_reset': KAFKA_CONFIG['auto_offset_reset'],
    'enable_auto_commit': KAFKA_CONFIG['enable_auto_commit'],
}

KAFKA_PRODUCER_CONFIG = {
    'bootstrap_servers': KAFKA_CONFIG['bootstrap_servers'],
    'client_id': KAFKA_CONFIG['client_id'],
}

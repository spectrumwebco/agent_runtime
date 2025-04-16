"""
Database configuration for Django with Apache Doris.
"""

import os
import logging
from pathlib import Path

logger = logging.getLogger(__name__)

def is_running_in_kubernetes():
    """Check if we're running in a Kubernetes environment."""
    return os.path.exists('/var/run/secrets/kubernetes.io/serviceaccount/token')

IN_KUBERNETES = is_running_in_kubernetes()

if IN_KUBERNETES:
    DATABASES = {
        'default': {
            'ENGINE': 'django.db.backends.mysql',
            'NAME': 'agent_runtime',
            'USER': 'root',
            'PASSWORD': '',  # Will be replaced by Vault
            'HOST': 'doris-fe.default.svc.cluster.local',
            'PORT': '9030',
            'OPTIONS': {
                'charset': 'utf8mb4',
                'use_unicode': True,
            },
        },
        'agent_db': {
            'ENGINE': 'django.db.backends.mysql',
            'NAME': 'agent_db',
            'USER': 'root',
            'PASSWORD': '',  # Will be replaced by Vault
            'HOST': 'doris-fe.default.svc.cluster.local',
            'PORT': '9030',
            'OPTIONS': {
                'charset': 'utf8mb4',
                'use_unicode': True,
            },
        },
        'trajectory_db': {
            'ENGINE': 'django.db.backends.mysql',
            'NAME': 'trajectory_db',
            'USER': 'root',
            'PASSWORD': '',  # Will be replaced by Vault
            'HOST': 'doris-fe.default.svc.cluster.local',
            'PORT': '9030',
            'OPTIONS': {
                'charset': 'utf8mb4',
                'use_unicode': True,
            },
        },
        'ml_db': {
            'ENGINE': 'django.db.backends.mysql',
            'NAME': 'ml_db',
            'USER': 'root',
            'PASSWORD': '',  # Will be replaced by Vault
            'HOST': 'doris-fe.default.svc.cluster.local',
            'PORT': '9030',
            'OPTIONS': {
                'charset': 'utf8mb4',
                'use_unicode': True,
            },
        },
    }
else:
    import socket
    try:
        socket.create_connection(('localhost', 9030), timeout=1)
        doris_available = True
    except (socket.timeout, socket.error):
        doris_available = False
        logger.warning("Apache Doris not available locally, using MariaDB for development")
    
    if doris_available:
        DATABASES = {
            'default': {
                'ENGINE': 'django.db.backends.mysql',
                'NAME': 'agent_runtime',
                'USER': 'root',
                'PASSWORD': '',
                'HOST': 'localhost',
                'PORT': '9030',
                'OPTIONS': {
                    'charset': 'utf8mb4',
                    'use_unicode': True,
                },
            },
            'agent_db': {
                'ENGINE': 'django.db.backends.mysql',
                'NAME': 'agent_db',
                'USER': 'root',
                'PASSWORD': '',
                'HOST': 'localhost',
                'PORT': '9030',
                'OPTIONS': {
                    'charset': 'utf8mb4',
                    'use_unicode': True,
                },
            },
            'trajectory_db': {
                'ENGINE': 'django.db.backends.mysql',
                'NAME': 'trajectory_db',
                'USER': 'root',
                'PASSWORD': '',
                'HOST': 'localhost',
                'PORT': '9030',
                'OPTIONS': {
                    'charset': 'utf8mb4',
                    'use_unicode': True,
                },
            },
            'ml_db': {
                'ENGINE': 'django.db.backends.mysql',
                'NAME': 'ml_db',
                'USER': 'root',
                'PASSWORD': '',
                'HOST': 'localhost',
                'PORT': '9030',
                'OPTIONS': {
                    'charset': 'utf8mb4',
                    'use_unicode': True,
                },
            },
        }
    else:
        DATABASES = {
            'default': {
                'ENGINE': 'django.db.backends.mysql',
                'NAME': 'agent_runtime',
                'USER': 'agent_user',
                'PASSWORD': 'agent_password',
                'HOST': 'localhost',
                'PORT': '3306',
                'OPTIONS': {
                    'charset': 'utf8mb4',
                    'use_unicode': True,
                },
            },
            'agent_db': {
                'ENGINE': 'django.db.backends.mysql',
                'NAME': 'agent_db',
                'USER': 'agent_user',
                'PASSWORD': 'agent_password',
                'HOST': 'localhost',
                'PORT': '3306',
                'OPTIONS': {
                    'charset': 'utf8mb4',
                    'use_unicode': True,
                },
            },
            'trajectory_db': {
                'ENGINE': 'django.db.backends.mysql',
                'NAME': 'trajectory_db',
                'USER': 'agent_user',
                'PASSWORD': 'agent_password',
                'HOST': 'localhost',
                'PORT': '3306',
                'OPTIONS': {
                    'charset': 'utf8mb4',
                    'use_unicode': True,
                },
            },
            'ml_db': {
                'ENGINE': 'django.db.backends.mysql',
                'NAME': 'ml_db',
                'USER': 'agent_user',
                'PASSWORD': 'agent_password',
                'HOST': 'localhost',
                'PORT': '3306',
                'OPTIONS': {
                    'charset': 'utf8mb4',
                    'use_unicode': True,
                },
            },
        }

DATABASE_ROUTERS = ['agent_api.database_routers.AgentDatabaseRouter']

if IN_KUBERNETES:
    REDIS_CONFIG = {
        'host': 'dragonfly-db.default.svc.cluster.local',
        'port': 6379,
        'db': 0,
        'password': '',  # Will be replaced by Vault
    }
else:
    try:
        socket.create_connection(('localhost', 6379), timeout=1)
        redis_available = True
    except (socket.timeout, socket.error):
        redis_available = False
        logger.warning("Redis not available locally, using in-memory cache for development")
    
    if redis_available:
        REDIS_CONFIG = {
            'host': 'localhost',
            'port': 6379,
            'db': 0,
            'password': '',
        }
    else:
        REDIS_CONFIG = {
            'host': 'localhost',
            'port': 6379,
            'db': 0,
            'password': '',
            'local_only': True,
        }

if IN_KUBERNETES:
    VECTOR_DB_CONFIG = {
        'host': 'ragflow.default.svc.cluster.local',
        'port': 8000,
        'api_key': '',  # Will be replaced by Vault
    }
else:
    try:
        socket.create_connection(('localhost', 8000), timeout=1)
        ragflow_available = True
    except (socket.timeout, socket.error):
        ragflow_available = False
        logger.warning("RAGflow not available locally, using mock vector database for development")
    
    if ragflow_available:
        VECTOR_DB_CONFIG = {
            'host': 'localhost',
            'port': 8000,
            'api_key': '',
        }
    else:
        VECTOR_DB_CONFIG = {
            'host': 'localhost',
            'port': 8000,
            'api_key': '',
            'local_only': True,
        }

if IN_KUBERNETES:
    ROCKETMQ_CONFIG = {
        'host': 'rocketmq.default.svc.cluster.local',
        'port': 9876,
        'access_key': '',  # Will be replaced by Vault
        'secret_key': '',  # Will be replaced by Vault
    }
else:
    try:
        socket.create_connection(('localhost', 9876), timeout=1)
        rocketmq_available = True
    except (socket.timeout, socket.error):
        rocketmq_available = False
        logger.warning("RocketMQ not available locally, using mock messaging for development")
    
    if rocketmq_available:
        ROCKETMQ_CONFIG = {
            'host': 'localhost',
            'port': 9876,
            'access_key': '',
            'secret_key': '',
        }
    else:
        ROCKETMQ_CONFIG = {
            'host': 'localhost',
            'port': 9876,
            'access_key': '',
            'secret_key': '',
            'local_only': True,
        }

DORIS_CONFIG = {
    'host': 'doris-fe.default.svc.cluster.local' if IN_KUBERNETES else 'localhost',
    'http_port': 8030,
    'query_port': 9030,
    'username': 'root',
    'password': '',  # Will be replaced by Vault
    'database': 'agent_runtime',
}

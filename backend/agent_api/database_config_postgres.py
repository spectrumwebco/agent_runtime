"""
CrunchyData PostgreSQL Operator configuration for Django.
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
    POSTGRES_DATABASES = {
        'agent_runtime': {
            'ENGINE': 'django.db.backends.postgresql',
            'NAME': 'agent_runtime',
            'USER': 'agent_user',
            'PASSWORD': os.environ.get('POSTGRES_PASSWORD', ''),
            'HOST': 'agent-postgres-cluster-primary.default.svc.cluster.local',
            'PORT': '5432',
            'CONN_MAX_AGE': 600,
        },
        'agent_db': {
            'ENGINE': 'django.db.backends.postgresql',
            'NAME': 'agent_db',
            'USER': 'agent_user',
            'PASSWORD': os.environ.get('POSTGRES_PASSWORD', ''),
            'HOST': 'agent-postgres-cluster-primary.default.svc.cluster.local',
            'PORT': '5432',
            'CONN_MAX_AGE': 600,
        },
        'trajectory_db': {
            'ENGINE': 'django.db.backends.postgresql',
            'NAME': 'trajectory_db',
            'USER': 'agent_user',
            'PASSWORD': os.environ.get('POSTGRES_PASSWORD', ''),
            'HOST': 'agent-postgres-cluster-primary.default.svc.cluster.local',
            'PORT': '5432',
            'CONN_MAX_AGE': 600,
        },
        'ml_db': {
            'ENGINE': 'django.db.backends.postgresql',
            'NAME': 'ml_db',
            'USER': 'agent_user',
            'PASSWORD': os.environ.get('POSTGRES_PASSWORD', ''),
            'HOST': 'agent-postgres-cluster-primary.default.svc.cluster.local',
            'PORT': '5432',
            'CONN_MAX_AGE': 600,
        },
    }
else:
    try:
        socket.create_connection(('localhost', 5432), timeout=1)
        postgres_available = True
    except (socket.timeout, socket.error):
        postgres_available = False
        logger.warning("PostgreSQL not available locally, using SQLite for development")
    
    if postgres_available:
        POSTGRES_DATABASES = {
            'agent_runtime': {
                'ENGINE': 'django.db.backends.postgresql',
                'NAME': 'agent_runtime',
                'USER': 'postgres',
                'PASSWORD': 'postgres',
                'HOST': 'localhost',
                'PORT': '5432',
                'CONN_MAX_AGE': 600,
            },
            'agent_db': {
                'ENGINE': 'django.db.backends.postgresql',
                'NAME': 'agent_db',
                'USER': 'postgres',
                'PASSWORD': 'postgres',
                'HOST': 'localhost',
                'PORT': '5432',
                'CONN_MAX_AGE': 600,
            },
            'trajectory_db': {
                'ENGINE': 'django.db.backends.postgresql',
                'NAME': 'trajectory_db',
                'USER': 'postgres',
                'PASSWORD': 'postgres',
                'HOST': 'localhost',
                'PORT': '5432',
                'CONN_MAX_AGE': 600,
            },
            'ml_db': {
                'ENGINE': 'django.db.backends.postgresql',
                'NAME': 'ml_db',
                'USER': 'postgres',
                'PASSWORD': 'postgres',
                'HOST': 'localhost',
                'PORT': '5432',
                'CONN_MAX_AGE': 600,
            },
        }
    else:
        POSTGRES_DATABASES = {
            'agent_runtime': {
                'ENGINE': 'django.db.backends.sqlite3',
                'NAME': os.path.join(os.path.dirname(os.path.dirname(os.path.abspath(__file__))), 'agent_runtime.sqlite3'),
            },
            'agent_db': {
                'ENGINE': 'django.db.backends.sqlite3',
                'NAME': os.path.join(os.path.dirname(os.path.dirname(os.path.abspath(__file__))), 'agent_db.sqlite3'),
            },
            'trajectory_db': {
                'ENGINE': 'django.db.backends.sqlite3',
                'NAME': os.path.join(os.path.dirname(os.path.dirname(os.path.abspath(__file__))), 'trajectory_db.sqlite3'),
            },
            'ml_db': {
                'ENGINE': 'django.db.backends.sqlite3',
                'NAME': os.path.join(os.path.dirname(os.path.dirname(os.path.abspath(__file__))), 'ml_db.sqlite3'),
            },
        }

DATABASE_ROUTERS = ['agent_api.database_routers.AgentRuntimeRouter']

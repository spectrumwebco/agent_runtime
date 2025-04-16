"""
Database configuration for the agent_api project.
This file contains the configuration for Supabase, Redis, DragonflyDB, and RAGflow.
"""

import os
import socket
import logging
from pydantic_settings import BaseSettings, SettingsConfigDict

logger = logging.getLogger(__name__)

def is_running_in_kubernetes():
    """Check if we're running in a Kubernetes environment."""
    return os.path.exists('/var/run/secrets/kubernetes.io/serviceaccount/token')

IN_KUBERNETES = is_running_in_kubernetes()
ENV = "kubernetes" if IN_KUBERNETES else "local"
logger.info(f"Detected environment: {ENV}")


class DatabaseSettings(BaseSettings):
    """Pydantic settings for database configuration."""

    environment: str = ENV

    db_engine: str = "django.db.backends.postgresql"
    db_name: str = "postgres"
    db_user: str = "postgres"
    db_password: str = "postgres"
    db_host: str = "supabase-db.default.svc.cluster.local" if IN_KUBERNETES else "localhost"
    db_port: int = 5432
    
    agent_db_name: str = "agent_db"
    agent_db_user: str = "postgres"
    agent_db_password: str = "postgres"
    agent_db_host: str = "supabase-db.default.svc.cluster.local" if IN_KUBERNETES else "localhost"
    agent_db_port: int = 5432
    
    trajectory_db_name: str = "trajectory_db"
    trajectory_db_user: str = "postgres"
    trajectory_db_password: str = "postgres"
    trajectory_db_host: str = "supabase-db.default.svc.cluster.local" if IN_KUBERNETES else "localhost"
    trajectory_db_port: int = 5432
    
    ml_db_name: str = "ml_db"
    ml_db_user: str = "postgres"
    ml_db_password: str = "postgres"
    ml_db_host: str = "supabase-db.default.svc.cluster.local" if IN_KUBERNETES else "localhost"
    ml_db_port: int = 5432

    redis_host: str = "dragonfly-db.default.svc.cluster.local" if IN_KUBERNETES else "localhost"
    redis_port: int = 6379
    redis_db: int = 0
    redis_password: str = ""
    redis_use_ssl: bool = False

    vector_db_api_key: str = ""
    vector_db_environment: str = "default"
    vector_db_index_name: str = "agent-docs"
    vector_db_host: str = "ragflow.default.svc.cluster.local" if IN_KUBERNETES else "localhost"
    vector_db_port: int = 8000

    rocketmq_host: str = "rocketmq.default.svc.cluster.local" if IN_KUBERNETES else "localhost"
    rocketmq_port: int = 9876

    model_config = SettingsConfigDict(
        env_file=".env",
        env_file_encoding="utf-8",
        env_prefix="AGENT_",
        extra="ignore",
    )


db_settings = DatabaseSettings()

def is_postgres_available():
    """Check if PostgreSQL is available on localhost."""
    if IN_KUBERNETES:
        return True
    
    try:
        s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        s.settimeout(1)
        s.connect(('localhost', 5432))
        s.close()
        return True
    except:
        return False

if ENV == "local" and not is_postgres_available():
    logger.warning("PostgreSQL not available locally, using MariaDB for development")
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
                'init_command': "SET sql_mode='STRICT_TRANS_TABLES'",
            },
        },
        'agent': {
            'ENGINE': 'django.db.backends.mysql',
            'NAME': 'agent_db',
            'USER': 'agent_user',
            'PASSWORD': 'agent_password',
            'HOST': 'localhost',
            'PORT': '3306',
            'OPTIONS': {
                'charset': 'utf8mb4',
                'init_command': "SET sql_mode='STRICT_TRANS_TABLES'",
            },
        },
        'trajectory': {
            'ENGINE': 'django.db.backends.mysql',
            'NAME': 'trajectory_db',
            'USER': 'agent_user',
            'PASSWORD': 'agent_password',
            'HOST': 'localhost',
            'PORT': '3306',
            'OPTIONS': {
                'charset': 'utf8mb4',
                'init_command': "SET sql_mode='STRICT_TRANS_TABLES'",
            },
        },
        'ml': {
            'ENGINE': 'django.db.backends.mysql',
            'NAME': 'ml_db',
            'USER': 'agent_user',
            'PASSWORD': 'agent_password',
            'HOST': 'localhost',
            'PORT': '3306',
            'OPTIONS': {
                'charset': 'utf8mb4',
                'init_command': "SET sql_mode='STRICT_TRANS_TABLES'",
            },
        },
    }
else:
    DATABASES = {
        'default': {
            'ENGINE': db_settings.db_engine,
            'NAME': db_settings.db_name,
            'USER': db_settings.db_user,
            'PASSWORD': db_settings.db_password,
            'HOST': db_settings.db_host,
            'PORT': db_settings.db_port,
            'OPTIONS': {
                'sslmode': 'prefer',  # Use 'require' in production
            },
        },
        'agent': {
            'ENGINE': db_settings.db_engine,
            'NAME': db_settings.agent_db_name,
            'USER': db_settings.agent_db_user,
            'PASSWORD': db_settings.agent_db_password,
            'HOST': db_settings.agent_db_host,
            'PORT': db_settings.agent_db_port,
            'OPTIONS': {
                'sslmode': 'prefer',  # Use 'require' in production
            },
        },
        'trajectory': {
            'ENGINE': db_settings.db_engine,
            'NAME': db_settings.trajectory_db_name,
            'USER': db_settings.trajectory_db_user,
            'PASSWORD': db_settings.trajectory_db_password,
            'HOST': db_settings.trajectory_db_host,
            'PORT': db_settings.trajectory_db_port,
            'OPTIONS': {
                'sslmode': 'prefer',  # Use 'require' in production
            },
        },
        'ml': {
            'ENGINE': db_settings.db_engine,
            'NAME': db_settings.ml_db_name,
            'USER': db_settings.ml_db_user,
            'PASSWORD': db_settings.ml_db_password,
            'HOST': db_settings.ml_db_host,
            'PORT': db_settings.ml_db_port,
            'OPTIONS': {
                'sslmode': 'prefer',  # Use 'require' in production
            },
        },
    }

DATABASE_ROUTERS = [
    'agent_api.database_routers.AgentRouter',
    'agent_api.database_routers.TrajectoryRouter',
    'agent_api.database_routers.MLRouter',
]

REDIS_CONFIG = {
    'host': db_settings.redis_host,
    'port': db_settings.redis_port,
    'db': db_settings.redis_db,
    'password': db_settings.redis_password,
    'ssl': db_settings.redis_use_ssl,
}

VECTOR_DB_CONFIG = {
    'api_key': db_settings.vector_db_api_key,
    'environment': db_settings.vector_db_environment,
    'index_name': db_settings.vector_db_index_name,
    'host': db_settings.vector_db_host,
    'port': db_settings.vector_db_port,
}

ROCKETMQ_CONFIG = {
    'host': db_settings.rocketmq_host,
    'port': db_settings.rocketmq_port,
}

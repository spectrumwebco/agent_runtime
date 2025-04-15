"""
Database configuration for the agent_api project.
This file contains the configuration for Supabase, Redis, DragonflyDB, and RAGflow.
"""

import os
from pydantic_settings import BaseSettings, SettingsConfigDict


class DatabaseSettings(BaseSettings):
    """Pydantic settings for database configuration."""

    db_engine: str = "django.db.backends.postgresql"
    db_name: str = "postgres"
    db_user: str = "postgres"
    db_password: str = "postgres"
    db_host: str = "supabase-db.default.svc.cluster.local"
    db_port: int = 5432
    
    agent_db_name: str = "agent_db"
    agent_db_user: str = "postgres"
    agent_db_password: str = "postgres"
    agent_db_host: str = "supabase-db.default.svc.cluster.local"
    agent_db_port: int = 5432
    
    trajectory_db_name: str = "trajectory_db"
    trajectory_db_user: str = "postgres"
    trajectory_db_password: str = "postgres"
    trajectory_db_host: str = "supabase-db.default.svc.cluster.local"
    trajectory_db_port: int = 5432
    
    ml_db_name: str = "ml_db"
    ml_db_user: str = "postgres"
    ml_db_password: str = "postgres"
    ml_db_host: str = "supabase-db.default.svc.cluster.local"
    ml_db_port: int = 5432

    redis_host: str = "dragonfly-db.default.svc.cluster.local"
    redis_port: int = 6379
    redis_db: int = 0
    redis_password: str = ""
    redis_use_ssl: bool = False

    vector_db_api_key: str = ""
    vector_db_environment: str = "default"
    vector_db_index_name: str = "agent-docs"
    vector_db_host: str = "ragflow.default.svc.cluster.local"
    vector_db_port: int = 8000

    rocketmq_host: str = "rocketmq.default.svc.cluster.local"
    rocketmq_port: int = 9876

    model_config = SettingsConfigDict(
        env_file=".env",
        env_file_encoding="utf-8",
        env_prefix="AGENT_",
        extra="ignore",
    )


db_settings = DatabaseSettings()

DATABASES = {
    'default': {
        'ENGINE': db_settings.db_engine,
        'NAME': db_settings.db_name,
        'USER': db_settings.db_user,
        'PASSWORD': db_settings.db_password,
        'HOST': db_settings.db_host,
        'PORT': db_settings.db_port,
        'OPTIONS': {
            'sslmode': 'require',
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
            'sslmode': 'require',
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
            'sslmode': 'require',
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
            'sslmode': 'require',
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

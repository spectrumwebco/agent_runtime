"""
Database configuration for the agent_api project.
This file contains the configuration for MariaDB, Redis, and Vector Database.
"""

import os
from pydantic_settings import BaseSettings, SettingsConfigDict


class DatabaseSettings(BaseSettings):
    """Pydantic settings for database configuration."""

    db_engine: str = "django.db.backends.mysql"
    db_name: str = "agent_runtime"
    db_user: str = "agent_user"
    db_password: str = ""
    db_host: str = "localhost"
    db_port: int = 3306

    redis_host: str = "localhost"
    redis_port: int = 6379
    redis_db: int = 0
    redis_password: str = ""
    redis_use_ssl: bool = False

    vector_db_api_key: str = ""
    vector_db_environment: str = "us-west1-gcp"
    vector_db_index_name: str = "agent-docs"

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
            'charset': 'utf8mb4',
            'init_command': "SET sql_mode='STRICT_TRANS_TABLES'",
        },
    }
}

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
}

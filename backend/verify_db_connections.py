"""
Standalone script to verify database connections.
This script doesn't rely on Django and can be run independently.
"""

import os
import sys
import socket
import logging
import time
from contextlib import contextmanager

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

def is_running_in_kubernetes():
    """Check if we're running in a Kubernetes environment."""
    return os.path.exists('/var/run/secrets/kubernetes.io/serviceaccount/token')

IN_KUBERNETES = is_running_in_kubernetes()
ENV = "kubernetes" if IN_KUBERNETES else "local"

class DatabaseConfig:
    """Database configuration."""
    
    def __init__(self):
        """Initialize database configuration."""
        self.environment = ENV
        
        self.db_host = "supabase-db.default.svc.cluster.local" if IN_KUBERNETES else "localhost"
        self.db_port = 5432
        self.db_user = "postgres"
        self.db_password = "postgres"
        self.db_name = "postgres"
        
        self.mariadb_host = "localhost"
        self.mariadb_port = 3306
        self.mariadb_user = "agent_user"
        self.mariadb_password = "agent_password"
        self.mariadb_name = "agent_runtime"
        
        self.redis_host = "dragonfly-db.default.svc.cluster.local" if IN_KUBERNETES else "localhost"
        self.redis_port = 6379
        
        self.ragflow_host = "ragflow.default.svc.cluster.local" if IN_KUBERNETES else "localhost"
        self.ragflow_port = 8000
        
        self.rocketmq_host = "rocketmq.default.svc.cluster.local" if IN_KUBERNETES else "localhost"
        self.rocketmq_port = 9876

db_config = DatabaseConfig()

def check_port_open(host, port, timeout=1):
    """Check if a port is open on a host."""
    try:
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.settimeout(timeout)
        result = sock.connect_ex((host, port))
        sock.close()
        return result == 0
    except Exception as e:
        logger.error(f"Error checking port {port} on {host}: {e}")
        return False

def check_postgres_connection():
    """Check PostgreSQL connection."""
    try:
        import psycopg2
        
        logger.info("Checking PostgreSQL connection...")
        
        conn = psycopg2.connect(
            host=db_config.db_host,
            port=db_config.db_port,
            user=db_config.db_user,
            password=db_config.db_password,
            dbname=db_config.db_name,
            connect_timeout=3
        )
        
        with conn.cursor() as cursor:
            cursor.execute("SELECT version();")
            version = cursor.fetchone()
            logger.info(f"Connected to PostgreSQL: {version[0]}")
        
        conn.close()
        return True
    except ImportError:
        logger.warning("psycopg2 not installed. Install with: pip install psycopg2-binary")
        return False
    except Exception as e:
        logger.error(f"PostgreSQL connection failed: {e}")
        return False

def check_mariadb_connection():
    """Check MariaDB connection."""
    try:
        import MySQLdb
        
        logger.info("Checking MariaDB connection...")
        
        conn = MySQLdb.connect(
            host=db_config.mariadb_host,
            port=db_config.mariadb_port,
            user=db_config.mariadb_user,
            passwd=db_config.mariadb_password,
            db=db_config.mariadb_name,
            connect_timeout=3
        )
        
        with conn.cursor() as cursor:
            cursor.execute("SELECT VERSION();")
            version = cursor.fetchone()
            logger.info(f"Connected to MariaDB: {version[0]}")
        
        conn.close()
        return True
    except ImportError:
        logger.warning("MySQLdb not installed. Install with: pip install mysqlclient")
        return False
    except Exception as e:
        logger.error(f"MariaDB connection failed: {e}")
        
        if "Unknown database" in str(e):
            try:
                logger.info("Attempting to create database and user...")
                
                conn = MySQLdb.connect(
                    host=db_config.mariadb_host,
                    port=db_config.mariadb_port,
                    user="root",
                    passwd="",
                    connect_timeout=3
                )
                
                with conn.cursor() as cursor:
                    cursor.execute(f"CREATE DATABASE IF NOT EXISTS {db_config.mariadb_name};")
                    
                    cursor.execute(f"CREATE USER IF NOT EXISTS '{db_config.mariadb_user}'@'localhost' IDENTIFIED BY '{db_config.mariadb_password}';")
                    
                    cursor.execute(f"GRANT ALL PRIVILEGES ON {db_config.mariadb_name}.* TO '{db_config.mariadb_user}'@'localhost';")
                    cursor.execute("FLUSH PRIVILEGES;")
                
                conn.close()
                
                logger.info("Database and user created successfully. Retrying connection...")
                
                conn = MySQLdb.connect(
                    host=db_config.mariadb_host,
                    port=db_config.mariadb_port,
                    user=db_config.mariadb_user,
                    passwd=db_config.mariadb_password,
                    db=db_config.mariadb_name,
                    connect_timeout=3
                )
                
                with conn.cursor() as cursor:
                    cursor.execute("SELECT VERSION();")
                    version = cursor.fetchone()
                    logger.info(f"Connected to MariaDB: {version[0]}")
                
                conn.close()
                return True
            except Exception as setup_error:
                logger.error(f"Failed to set up MariaDB: {setup_error}")
                return False
        
        return False

def check_redis_connection():
    """Check Redis connection."""
    try:
        import redis
        
        logger.info("Checking Redis (DragonflyDB) connection...")
        
        r = redis.Redis(
            host=db_config.redis_host,
            port=db_config.redis_port,
            socket_timeout=3
        )
        
        pong = r.ping()
        logger.info(f"Connected to Redis: {pong}")
        
        return True
    except ImportError:
        logger.warning("redis not installed. Install with: pip install redis")
        return False
    except Exception as e:
        logger.error(f"Redis connection failed: {e}")
        return False

def check_ragflow_connection():
    """Check RAGflow connection."""
    try:
        import requests
        
        logger.info("Checking RAGflow connection...")
        
        response = requests.get(
            f"http://{db_config.ragflow_host}:{db_config.ragflow_port}/health",
            timeout=3
        )
        
        if response.status_code == 200:
            logger.info(f"Connected to RAGflow: {response.json()}")
            return True
        else:
            logger.error(f"RAGflow returned status code: {response.status_code}")
            return False
    except ImportError:
        logger.warning("requests not installed. Install with: pip install requests")
        return False
    except Exception as e:
        logger.error(f"RAGflow connection failed: {e}")
        return False

def check_rocketmq_connection():
    """Check RocketMQ connection."""
    logger.info("Checking RocketMQ connection...")
    
    if check_port_open(db_config.rocketmq_host, db_config.rocketmq_port):
        logger.info("RocketMQ port is open")
        return True
    else:
        logger.error("RocketMQ port is closed")
        return False

def main():
    """Main function."""
    print("\n=== Database Connection Verification ===\n")
    
    print(f"Environment: {ENV}")
    print(f"Running in Kubernetes: {IN_KUBERNETES}\n")
    
    print("\n=== PostgreSQL (Supabase) Connection ===")
    postgres_ok = check_postgres_connection()
    print(f"PostgreSQL connection: {'✅ OK' if postgres_ok else '❌ Failed'}")
    
    if ENV == "local":
        print("\n=== MariaDB Connection (Local Development) ===")
        mariadb_ok = check_mariadb_connection()
        print(f"MariaDB connection: {'✅ OK' if mariadb_ok else '❌ Failed'}")
    
    print("\n=== Redis (DragonflyDB) Connection ===")
    redis_ok = check_redis_connection()
    print(f"Redis connection: {'✅ OK' if redis_ok else '❌ Failed'}")
    
    print("\n=== RAGflow Connection ===")
    ragflow_ok = check_ragflow_connection()
    print(f"RAGflow connection: {'✅ OK' if ragflow_ok else '❌ Failed'}")
    
    print("\n=== RocketMQ Connection ===")
    rocketmq_ok = check_rocketmq_connection()
    print(f"RocketMQ connection: {'✅ OK' if rocketmq_ok else '❌ Failed'}")
    
    print("\n=== Connection Summary ===")
    if ENV == "local":
        print(f"PostgreSQL: {'✅' if postgres_ok else '❌'}")
        print(f"MariaDB: {'✅' if mariadb_ok else '❌'}")
        print(f"Redis: {'✅' if redis_ok else '❌'}")
        print(f"RAGflow: {'✅' if ragflow_ok else '❌'}")
        print(f"RocketMQ: {'✅' if rocketmq_ok else '❌'}")
        
        if not postgres_ok and not mariadb_ok:
            print("\n⚠️ No database connections available. Please set up at least one database.")
        elif not postgres_ok and mariadb_ok:
            print("\n✅ MariaDB is available for local development.")
        elif postgres_ok and not mariadb_ok:
            print("\n✅ PostgreSQL is available for local development.")
        else:
            print("\n✅ Both PostgreSQL and MariaDB are available.")
    else:
        print(f"PostgreSQL: {'✅' if postgres_ok else '❌'}")
        print(f"Redis: {'✅' if redis_ok else '❌'}")
        print(f"RAGflow: {'✅' if ragflow_ok else '❌'}")
        print(f"RocketMQ: {'✅' if rocketmq_ok else '❌'}")
        
        if not postgres_ok:
            print("\n⚠️ PostgreSQL connection failed. Please check your Kubernetes configuration.")

if __name__ == "__main__":
    main()

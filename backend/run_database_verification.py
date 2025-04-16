"""
Script to run database verification tests for Apache Doris, Kafka, and PostgreSQL.
"""

import os
import sys
import subprocess
import argparse
import time
import logging

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    handlers=[
        logging.StreamHandler(),
        logging.FileHandler('database_verification.log')
    ]
)

logger = logging.getLogger(__name__)

def run_command(command, description):
    """Run a shell command and log the output."""
    logger.info(f"Running {description}...")
    try:
        result = subprocess.run(
            command,
            shell=True,
            check=True,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True
        )
        logger.info(f"{description} completed successfully")
        logger.info(f"Output: {result.stdout}")
        return True, result.stdout
    except subprocess.CalledProcessError as e:
        logger.error(f"{description} failed with error code {e.returncode}")
        logger.error(f"Error output: {e.stderr}")
        return False, e.stderr

def verify_database_connections():
    """Verify database connections for all configured databases."""
    logger.info("Verifying database connections...")
    
    success, output = run_command(
        "python manage.py check --database=default",
        "Django database configuration check for Apache Doris"
    )
    if not success:
        logger.error("Database configuration check failed for Apache Doris")
        return False
    
    success, output = run_command(
        "python manage.py check --database=agent_db",
        "Django database configuration check for PostgreSQL"
    )
    if not success:
        logger.error("Database configuration check failed for PostgreSQL")
        return False
    
    success, output = run_command(
        "python manage.py verify_database_integration --all",
        "Database integration verification"
    )
    if not success:
        logger.error("Database integration verification failed")
        return False
    
    logger.info("All database connections verified successfully")
    return True

def run_database_tests():
    """Run database integration tests."""
    logger.info("Running database integration tests...")
    
    success, output = run_command(
        "python manage.py test apps.python_agent.tests.test_database_models",
        "Database model tests"
    )
    if not success:
        logger.error("Database model tests failed")
        return False
    
    success, output = run_command(
        "python manage.py test apps.python_agent.tests.test_database_integration",
        "Database integration tests"
    )
    if not success:
        logger.error("Database integration tests failed")
        return False
    
    logger.info("All database tests passed successfully")
    return True

def verify_kafka_integration():
    """Verify Kafka integration."""
    logger.info("Verifying Kafka integration...")
    
    success, output = run_command(
        "python manage.py verify_database_integration --kafka",
        "Kafka integration verification"
    )
    if not success:
        logger.error("Kafka integration verification failed")
        return False
    
    success, output = run_command(
        "python -c \"from backend.integrations.kafka import KafkaClient; "
        "client = KafkaClient(); "
        "client.produce_message('test-topic', {'message': 'test'}); "
        "print('Message produced successfully'); "
        "message = client.consume_message('test-topic', timeout=10); "
        "print(f'Message consumed: {message}')\"",
        "Kafka message production and consumption test"
    )
    if not success:
        logger.error("Kafka message production and consumption test failed")
        return False
    
    logger.info("Kafka integration verified successfully")
    return True

def verify_postgres_integration():
    """Verify PostgreSQL integration."""
    logger.info("Verifying PostgreSQL integration...")
    
    success, output = run_command(
        "python manage.py verify_database_integration --postgres",
        "PostgreSQL integration verification"
    )
    if not success:
        logger.error("PostgreSQL integration verification failed")
        return False
    
    success, output = run_command(
        "python manage.py setup_postgres --check",
        "PostgreSQL cluster management test"
    )
    if not success:
        logger.error("PostgreSQL cluster management test failed")
        return False
    
    logger.info("PostgreSQL integration verified successfully")
    return True

def verify_doris_integration():
    """Verify Apache Doris integration."""
    logger.info("Verifying Apache Doris integration...")
    
    success, output = run_command(
        "python manage.py verify_database_integration --doris",
        "Apache Doris integration verification"
    )
    if not success:
        logger.error("Apache Doris integration verification failed")
        return False
    
    success, output = run_command(
        "python -c \"from django.db import connections; "
        "cursor = connections['default'].cursor(); "
        "cursor.execute('SELECT VERSION()'); "
        "version = cursor.fetchone()[0]; "
        "print(f'Apache Doris version: {version}'); "
        "cursor.execute('SHOW DATABASES'); "
        "databases = cursor.fetchall(); "
        "print(f'Databases: {databases}')\"",
        "Apache Doris query capabilities test"
    )
    if not success:
        logger.error("Apache Doris query capabilities test failed")
        return False
    
    logger.info("Apache Doris integration verified successfully")
    return True

def verify_cross_database_integration():
    """Verify integration between all database systems."""
    logger.info("Verifying cross-database integration...")
    
    success, output = run_command(
        "python manage.py verify_database_integration --integration",
        "Cross-database integration verification"
    )
    if not success:
        logger.error("Cross-database integration verification failed")
        return False
    
    logger.info("Cross-database integration verified successfully")
    return True

def main():
    """Main function to run database verification."""
    parser = argparse.ArgumentParser(description='Verify database integration')
    parser.add_argument('--all', action='store_true', help='Run all verification tests')
    parser.add_argument('--connections', action='store_true', help='Verify database connections')
    parser.add_argument('--tests', action='store_true', help='Run database tests')
    parser.add_argument('--kafka', action='store_true', help='Verify Kafka integration')
    parser.add_argument('--postgres', action='store_true', help='Verify PostgreSQL integration')
    parser.add_argument('--doris', action='store_true', help='Verify Apache Doris integration')
    parser.add_argument('--cross', action='store_true', help='Verify cross-database integration')
    
    args = parser.parse_args()
    
    os.chdir(os.path.dirname(os.path.abspath(__file__)))
    
    run_all = args.all or not any([
        args.connections, args.tests, args.kafka, 
        args.postgres, args.doris, args.cross
    ])
    
    success = True
    
    if run_all or args.connections:
        if not verify_database_connections():
            success = False
    
    if run_all or args.tests:
        if not run_database_tests():
            success = False
    
    if run_all or args.kafka:
        if not verify_kafka_integration():
            success = False
    
    if run_all or args.postgres:
        if not verify_postgres_integration():
            success = False
    
    if run_all or args.doris:
        if not verify_doris_integration():
            success = False
    
    if run_all or args.cross:
        if not verify_cross_database_integration():
            success = False
    
    if success:
        logger.info("All database verification tests passed successfully")
        return 0
    else:
        logger.error("Some database verification tests failed")
        return 1

if __name__ == '__main__':
    sys.exit(main())

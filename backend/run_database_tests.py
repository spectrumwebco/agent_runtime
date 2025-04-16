"""
Script to run comprehensive database integration tests for all database systems.
"""

import os
import sys
import subprocess
import argparse
import logging
from datetime import datetime

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    handlers=[
        logging.StreamHandler(),
        logging.FileHandler(f'database_tests_{datetime.now().strftime("%Y%m%d_%H%M%S")}.log')
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

def test_all_databases():
    """Run tests for all database systems."""
    logger.info("Testing all database systems...")
    
    success, output = run_command(
        "python manage.py test_all_databases --all",
        "All database systems test"
    )
    
    if not success:
        logger.error("All database systems test failed")
        return False
    
    logger.info("All database systems test completed successfully")
    return True

def test_supabase():
    """Test Supabase integration."""
    logger.info("Testing Supabase integration...")
    
    success, output = run_command(
        "python manage.py test_all_databases --supabase",
        "Supabase integration test"
    )
    
    if not success:
        logger.error("Supabase integration test failed")
        return False
    
    logger.info("Supabase integration test completed successfully")
    return True

def test_ragflow():
    """Test RAGflow integration."""
    logger.info("Testing RAGflow integration...")
    
    success, output = run_command(
        "python manage.py test_all_databases --ragflow",
        "RAGflow integration test"
    )
    
    if not success:
        logger.error("RAGflow integration test failed")
        return False
    
    logger.info("RAGflow integration test completed successfully")
    return True

def test_dragonfly():
    """Test DragonflyDB integration."""
    logger.info("Testing DragonflyDB integration...")
    
    success, output = run_command(
        "python manage.py test_all_databases --dragonfly",
        "DragonflyDB integration test"
    )
    
    if not success:
        logger.error("DragonflyDB integration test failed")
        return False
    
    logger.info("DragonflyDB integration test completed successfully")
    return True

def test_rocketmq():
    """Test RocketMQ integration."""
    logger.info("Testing RocketMQ integration...")
    
    success, output = run_command(
        "python manage.py test_all_databases --rocketmq",
        "RocketMQ integration test"
    )
    
    if not success:
        logger.error("RocketMQ integration test failed")
        return False
    
    logger.info("RocketMQ integration test completed successfully")
    return True

def test_doris():
    """Test Apache Doris integration."""
    logger.info("Testing Apache Doris integration...")
    
    success, output = run_command(
        "python manage.py test_all_databases --doris",
        "Apache Doris integration test"
    )
    
    if not success:
        logger.error("Apache Doris integration test failed")
        return False
    
    logger.info("Apache Doris integration test completed successfully")
    return True

def test_postgres():
    """Test PostgreSQL integration."""
    logger.info("Testing PostgreSQL integration...")
    
    success, output = run_command(
        "python manage.py test_all_databases --postgres",
        "PostgreSQL integration test"
    )
    
    if not success:
        logger.error("PostgreSQL integration test failed")
        return False
    
    logger.info("PostgreSQL integration test completed successfully")
    return True

def test_kafka():
    """Test Kafka integration."""
    logger.info("Testing Kafka integration...")
    
    success, output = run_command(
        "python manage.py test_all_databases --kafka",
        "Kafka integration test"
    )
    
    if not success:
        logger.error("Kafka integration test failed")
        return False
    
    logger.info("Kafka integration test completed successfully")
    return True

def run_django_tests():
    """Run Django test suite for database integration."""
    logger.info("Running Django test suite for database integration...")
    
    success, output = run_command(
        "python manage.py test apps.python_agent.tests.test_database_integration",
        "Django database integration tests"
    )
    
    if not success:
        logger.error("Django database integration tests failed")
        return False
    
    success, output = run_command(
        "python manage.py test apps.python_agent.tests.test_database_models",
        "Django database models tests"
    )
    
    if not success:
        logger.error("Django database models tests failed")
        return False
    
    logger.info("Django test suite for database integration completed successfully")
    return True

def verify_database_integration():
    """Verify database integration."""
    logger.info("Verifying database integration...")
    
    success, output = run_command(
        "python manage.py verify_database_integration --all",
        "Database integration verification"
    )
    
    if not success:
        logger.error("Database integration verification failed")
        return False
    
    logger.info("Database integration verification completed successfully")
    return True

def main():
    """Main function."""
    parser = argparse.ArgumentParser(description='Run database integration tests')
    parser.add_argument('--all', action='store_true', help='Run all tests')
    parser.add_argument('--supabase', action='store_true', help='Test Supabase integration')
    parser.add_argument('--ragflow', action='store_true', help='Test RAGflow integration')
    parser.add_argument('--dragonfly', action='store_true', help='Test DragonflyDB integration')
    parser.add_argument('--rocketmq', action='store_true', help='Test RocketMQ integration')
    parser.add_argument('--doris', action='store_true', help='Test Apache Doris integration')
    parser.add_argument('--postgres', action='store_true', help='Test PostgreSQL integration')
    parser.add_argument('--kafka', action='store_true', help='Test Kafka integration')
    parser.add_argument('--django', action='store_true', help='Run Django test suite')
    parser.add_argument('--verify', action='store_true', help='Verify database integration')
    
    args = parser.parse_args()
    
    os.chdir(os.path.dirname(os.path.abspath(__file__)))
    
    os.environ.setdefault('DJANGO_SETTINGS_MODULE', 'agent_api.settings')
    
    run_all = args.all or not any([
        args.supabase, args.ragflow, args.dragonfly, args.rocketmq,
        args.doris, args.postgres, args.kafka, args.django, args.verify
    ])
    
    success = True
    
    if run_all or args.supabase:
        if not test_supabase():
            success = False
    
    if run_all or args.ragflow:
        if not test_ragflow():
            success = False
    
    if run_all or args.dragonfly:
        if not test_dragonfly():
            success = False
    
    if run_all or args.rocketmq:
        if not test_rocketmq():
            success = False
    
    if run_all or args.doris:
        if not test_doris():
            success = False
    
    if run_all or args.postgres:
        if not test_postgres():
            success = False
    
    if run_all or args.kafka:
        if not test_kafka():
            success = False
    
    if run_all or args.django:
        if not run_django_tests():
            success = False
    
    if run_all or args.verify:
        if not verify_database_integration():
            success = False
    
    if success:
        logger.info("All tests completed successfully")
        return 0
    else:
        logger.error("Some tests failed")
        return 1

if __name__ == '__main__':
    sys.exit(main())

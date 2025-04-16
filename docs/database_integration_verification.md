# Database Integration Verification

This document outlines the verification process for the integration of Apache Doris, Apache Kafka, and CrunchyData PostgreSQL Operator with Django in the agent_runtime system.

## Overview

The agent_runtime system uses multiple database technologies:

1. **Apache Doris** - Enterprise-grade OLAP database (replacing MariaDB)
2. **Apache Kafka** - Event streaming platform for real-time data pipelines
3. **CrunchyData PostgreSQL** - Enterprise PostgreSQL database managed by Kubernetes Operator

## Verification Tools

The following tools have been created to verify database integration:

### 1. Management Commands

- `verify_database_integration.py` - Django management command to verify database connections and operations
  - Usage: `python manage.py verify_database_integration [--all|--doris|--kafka|--postgres|--integration]`

### 2. Test Suite

- `test_database_models.py` - Django test suite for database models and ORM integration
  - Usage: `python manage.py test apps.python_agent.tests.test_database_models`

### 3. Verification Scripts

- `run_database_verification.py` - Python script to run all verification tests
  - Usage: `python backend/run_database_verification.py [--all|--connections|--tests|--kafka|--postgres|--doris|--cross]`

- `run_database_verification.sh` - Shell script to run all verification tests
  - Usage: `./backend/run_database_verification.sh [--all|--connections|--tests|--kafka|--postgres|--doris|--cross]`

## Verification Process

### 1. Database Connections

Verify that Django can connect to all database systems:

```bash
python manage.py verify_database_integration --all
```

### 2. Database Models

Verify that Django ORM can interact with all database systems:

```bash
python manage.py test apps.python_agent.tests.test_database_models
```

### 3. Cross-Database Integration

Verify that data can flow between database systems:

```bash
python manage.py verify_database_integration --integration
```

## Kubernetes Integration

The database systems are deployed in Kubernetes using the following configurations:

1. **Apache Doris**
   - Deployment: `kubernetes/doris-deployment.yaml`
   - Service Discovery: `kubernetes/database-service-discovery.yaml`

2. **Apache Kafka**
   - Deployment: `kubernetes/kafka-deployment.yaml`
   - Service Discovery: `kubernetes/database-service-discovery.yaml`

3. **CrunchyData PostgreSQL**
   - Operator: `kubernetes/crunchydata-postgres-operator-deployment.yaml`
   - Cluster: `kubernetes/postgres-cluster-crunchy.yaml`
   - Service Discovery: `kubernetes/database-service-discovery.yaml`

## Terraform Integration

The database systems are provisioned using Terraform modules:

1. **Apache Doris**
   - Module: `terraform/modules/doris`

2. **Apache Kafka**
   - Module: `terraform/modules/kafka`

3. **CrunchyData PostgreSQL**
   - Module: `terraform/modules/postgres`

## Django Integration

The database systems are integrated with Django using the following configurations:

1. **Database Settings**
   - Configuration: `backend/agent_api/settings.py`
   - Doris Config: `backend/agent_api/database_config_doris.py`
   - Kafka Config: `backend/agent_api/database_config_kafka.py`
   - PostgreSQL Config: `backend/agent_api/database_config_postgres.py`

2. **Database Routers**
   - Router: `backend/agent_api/database_routers.py`

## Verification Results

The verification process checks the following aspects of database integration:

1. **Connection Verification**
   - Verify that Django can connect to all database systems
   - Verify that database credentials are properly managed using Hashicorp Vault

2. **CRUD Operations**
   - Verify that Django ORM can perform CRUD operations on all database systems
   - Verify that database models are properly mapped to database tables

3. **Cross-Database Integration**
   - Verify that data can flow between PostgreSQL and Apache Doris through Kafka
   - Verify that database routers correctly route queries to the appropriate database

4. **Performance Verification**
   - Verify that database operations meet performance requirements
   - Verify that database connections are properly pooled and managed

## Troubleshooting

If verification tests fail, check the following:

1. **Connection Issues**
   - Verify that database services are running in Kubernetes
   - Verify that database credentials are correctly configured in Vault
   - Verify that service discovery is properly configured

2. **ORM Issues**
   - Verify that database models are properly defined
   - Verify that database migrations have been applied
   - Verify that database routers are correctly configured

3. **Integration Issues**
   - Verify that Kafka topics are properly configured
   - Verify that data formats are compatible between systems
   - Verify that database schemas are properly aligned

## Conclusion

The database integration verification process ensures that Apache Doris, Apache Kafka, and CrunchyData PostgreSQL are properly integrated with Django in the agent_runtime system. The verification tools provide a comprehensive suite of tests to validate the integration and identify any issues.

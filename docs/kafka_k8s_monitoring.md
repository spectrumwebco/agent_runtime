# Kafka and Kubernetes Monitoring Integration

This document describes the integration between Apache Kafka and Kubernetes monitoring in the agent_runtime project.

## Overview

The agent_runtime project uses Apache Kafka to collect and process events from Kubernetes clusters. This integration enables real-time monitoring of Kubernetes resources and provides a stream of events that can be consumed by various components of the system.

## Architecture

The integration consists of the following components:

1. **Apache Kafka**: A distributed event streaming platform that serves as the central hub for all Kubernetes events.
2. **Kubernetes Monitor**: A Python service that watches Kubernetes resources and publishes events to Kafka.
3. **Event Consumers**: Various components that consume events from Kafka for processing, visualization, or storage.

## Kafka Configuration

Kafka is configured with the following topics:

- `k8s-events`: General Kubernetes events
- `k8s-pods`: Pod-specific events (creation, deletion, updates)
- `k8s-deployments`: Deployment-specific events
- `k8s-services`: Service-specific events
- `k8s-configmaps`: ConfigMap-specific events
- `k8s-secrets`: Secret-specific events (without sensitive data)
- `shared-state`: Events related to shared state changes
- `event-stream`: General event stream for the application

## Kubernetes Monitor

The Kubernetes Monitor is a Django management command that watches Kubernetes resources and publishes events to Kafka. It can be started with:

```bash
python manage.py start_k8s_monitor [--namespace NAMESPACE] [--poll-interval SECONDS] [--resources RESOURCES] [--daemon]
```

Options:
- `--namespace`: Kubernetes namespace to monitor (default: default)
- `--poll-interval`: Poll interval in seconds (default: 30)
- `--resources`: Comma-separated list of resources to monitor (default: pods,services,deployments,statefulsets,configmaps,secrets)
- `--daemon`: Run as a daemon (default: False)

To stop the monitor:

```bash
python manage.py stop_k8s_monitor
```

## Deployment

### Kubernetes

The Kafka and Kubernetes Monitor components are deployed using Kubernetes manifests:

```bash
kubectl apply -f kubernetes/kafka-deployment.yaml
kubectl apply -f kubernetes/k8s-monitor-deployment.yaml
```

### Terraform

For production deployments, Terraform modules are provided:

```hcl
module "kafka_k8s_monitor" {
  source = "./modules/kafka_k8s_monitor"
  
  namespace           = "default"
  kafka_replicas      = 3
  monitor_namespace   = "default"
  poll_interval       = 30
  resources_to_monitor = "pods,services,deployments,statefulsets,configmaps,secrets"
}
```

## Integration with Django

The Django application integrates with Kafka using the `kafka-python` library. The integration is configured in:

- `backend/agent_api/database_config_kafka.py`
- `backend/apps/python_agent/integrations/kafka.py`

### Local Development

For local development, you can set up Kafka using the provided Django management command:

```bash
python manage.py setup_kafka
```

This will:
- Check if Kafka is running
- Create the necessary topics
- Configure the Kafka client

## Event Consumers

The following components consume events from Kafka:

1. **Shared State Manager**: Consumes events from the `shared-state` topic to update the shared state.
2. **Event Stream Processor**: Consumes events from the `event-stream` topic for processing.
3. **Monitoring Dashboard**: Consumes events from various topics for visualization.

## Monitoring

The Kafka and Kubernetes Monitor components are themselves monitored using:

- Prometheus for metrics collection
- Grafana for visualization
- Kubernetes events fed back into Kafka

## Troubleshooting

If you encounter issues with the Kafka and Kubernetes Monitor integration, you can:

1. Check the Kafka logs:
   ```bash
   kubectl logs -n default deployment/kafka -c kafka
   ```

2. Check the Kubernetes Monitor logs:
   ```bash
   kubectl logs -n default deployment/k8s-monitor
   ```

3. Check the Kafka topics:
   ```bash
   kubectl exec -it -n default deployment/kafka -c kafka -- /opt/kafka/bin/kafka-topics.sh --list --zookeeper localhost:2181
   ```

4. Check the Kafka consumer groups:
   ```bash
   kubectl exec -it -n default deployment/kafka -c kafka -- /opt/kafka/bin/kafka-consumer-groups.sh --bootstrap-server localhost:9092 --list
   ```

5. Use the Kafka client in Django:
   ```python
   from apps.python_agent.integrations.kafka import KafkaClient
   client = KafkaClient()
   status = client.check_connection()
   print(status)
   ```

output "kafka_service_name" {
  description = "Service name for Kafka"
  value       = kubernetes_service.kafka.metadata[0].name
}

output "kafka_service_namespace" {
  description = "Namespace of the Kafka service"
  value       = kubernetes_service.kafka.metadata[0].namespace
}

output "zookeeper_service_name" {
  description = "Service name for ZooKeeper"
  value       = kubernetes_service.zookeeper.metadata[0].name
}

output "zookeeper_service_namespace" {
  description = "Namespace of the ZooKeeper service"
  value       = kubernetes_service.zookeeper.metadata[0].namespace
}

output "kafka_bootstrap_servers" {
  description = "Kafka bootstrap servers"
  value       = "${kubernetes_service.kafka.metadata[0].name}.${kubernetes_service.kafka.metadata[0].namespace}.svc.cluster.local:9092"
}

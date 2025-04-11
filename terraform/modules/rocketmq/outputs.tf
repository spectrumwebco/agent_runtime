output "namespace" {
  description = "Namespace where RocketMQ is deployed"
  value       = local.namespace
}

output "name_server_endpoint" {
  description = "Endpoint for RocketMQ name server"
  value       = "${var.name}-nameserver.${local.namespace}.svc.cluster.local:9876"
}

output "broker_endpoint" {
  description = "Endpoint for RocketMQ broker"
  value       = "${var.name}-broker.${local.namespace}.svc.cluster.local:10911"
}

output "dashboard_endpoint" {
  description = "Endpoint for RocketMQ dashboard"
  value       = var.dashboard_enabled ? "${var.name}-dashboard.${local.namespace}.svc.cluster.local:8080" : null
}

output "dashboard_url" {
  description = "URL for RocketMQ dashboard (if ingress is enabled)"
  value       = var.dashboard_enabled && var.ingress_enabled ? "https://rocketmq-dashboard.${var.ingress_domain}" : null
}

output "name_server_service_name" {
  description = "Name of the RocketMQ name server service"
  value       = kubernetes_service.nameserver.metadata[0].name
}

output "broker_service_name" {
  description = "Name of the RocketMQ broker service"
  value       = kubernetes_service.broker.metadata[0].name
}

output "dashboard_service_name" {
  description = "Name of the RocketMQ dashboard service"
  value       = var.dashboard_enabled ? kubernetes_service.dashboard[0].metadata[0].name : null
}

output "name_server_replicas" {
  description = "Number of RocketMQ name server replicas"
  value       = var.high_availability ? var.name_server_replicas : 1
}

output "broker_replicas" {
  description = "Number of RocketMQ broker replicas"
  value       = var.high_availability ? var.broker_replicas : 1
}

output "high_availability_enabled" {
  description = "Whether high availability is enabled for RocketMQ"
  value       = var.high_availability
}

output "topics" {
  description = "List of RocketMQ topics created"
  value       = [for topic in var.topic_configs : topic.name]
}

output "topic_configs" {
  description = "Configuration details for RocketMQ topics"
  value       = var.topic_configs
}

output "acl_enabled" {
  description = "Whether ACL is enabled for RocketMQ"
  value       = var.acl_enabled
}

output "prometheus_integration_enabled" {
  description = "Whether Prometheus integration is enabled for RocketMQ"
  value       = var.prometheus_integration
}

output "kata_container_integration_enabled" {
  description = "Whether Kata Containers integration is enabled for RocketMQ"
  value       = var.kata_container_integration
}

output "connection_string" {
  description = "Connection string for RocketMQ clients"
  value       = "${var.name}-nameserver.${local.namespace}.svc.cluster.local:9876;${var.name}-broker.${local.namespace}.svc.cluster.local:10911"
}

output "go_client_example" {
  description = "Example code for connecting to RocketMQ using the Go client"
  value       = <<-EOT
    // Example code for connecting to RocketMQ using the Go client
    package main
    
    import (
      "context"
      "fmt"
      "os"
      "time"
    
      "github.com/apache/rocketmq-client-go/v2"
      "github.com/apache/rocketmq-client-go/v2/primitive"
      "github.com/apache/rocketmq-client-go/v2/producer"
    )
    
    func main() {
      // Create a new producer
      p, err := rocketmq.NewProducer(
        producer.WithNameServer([]string{"${var.name}-nameserver.${local.namespace}.svc.cluster.local:9876"}),
        producer.WithRetry(2),
        producer.WithGroupName("agent-producer"),
      )
      if err != nil {
        fmt.Printf("Failed to create producer: %s\n", err)
        os.Exit(1)
      }
    
      // Start the producer
      err = p.Start()
      if err != nil {
        fmt.Printf("Failed to start producer: %s\n", err)
        os.Exit(1)
      }
      defer p.Shutdown()
    
      // Create a message
      msg := primitive.NewMessage(
        "k8s-lifecycle", // topic
        []byte("Hello RocketMQ!"), // body
      )
      
      // Add properties to the message
      msg.WithProperty("key1", "value1")
      msg.WithProperty("key2", "value2")
    
      // Send the message
      res, err := p.SendSync(context.Background(), msg)
      if err != nil {
        fmt.Printf("Failed to send message: %s\n", err)
        os.Exit(1)
      }
    
      fmt.Printf("Send message success: result=%s\n", res.String())
    }
  EOT
}

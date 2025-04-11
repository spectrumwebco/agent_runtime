output "service_name" {
  description = "Name of the DragonflyDB service"
  value       = kubernetes_service.dragonfly.metadata[0].name
}

output "service_namespace" {
  description = "Namespace of the DragonflyDB service"
  value       = kubernetes_service.dragonfly.metadata[0].namespace
}

output "service_port" {
  description = "Port of the DragonflyDB service"
  value       = kubernetes_service.dragonfly.spec[0].port[0].port
}

output "connection_string" {
  description = "Connection string for DragonflyDB"
  value       = "${kubernetes_service.dragonfly.metadata[0].name}.${kubernetes_service.dragonfly.metadata[0].namespace}.svc.cluster.local:${kubernetes_service.dragonfly.spec[0].port[0].port}"
}

output "headless_service_name" {
  description = "Name of the DragonflyDB headless service for replication"
  value       = kubernetes_service.dragonfly_headless.metadata[0].name
}

output "high_availability_enabled" {
  description = "Whether high availability is enabled for DragonflyDB"
  value       = var.high_availability
}

output "replicas" {
  description = "Number of DragonflyDB replicas"
  value       = var.high_availability ? var.replicas : 1
}

output "persistence_enabled" {
  description = "Whether persistence is enabled for DragonflyDB"
  value       = var.persistence_enabled
}

output "persistence_path" {
  description = "Path for DragonflyDB persistence files"
  value       = var.persistence_enabled ? var.persistence_path : null
}

output "rollback_enabled" {
  description = "Whether rollback mechanism is enabled for DragonflyDB"
  value       = var.rollback_enabled
}

output "rollback_scripts" {
  description = "List of available rollback scripts"
  value       = var.rollback_enabled ? [
    "rollback.sh",
    "create-snapshot.sh",
    "list-snapshots.sh"
  ] : []
}

output "snapshot_enabled" {
  description = "Whether snapshots are enabled for DragonflyDB"
  value       = var.snapshot_enabled
}

output "snapshot_schedule" {
  description = "Schedule for DragonflyDB snapshots"
  value       = var.snapshot_enabled ? var.snapshot_schedule : null
}

output "snapshot_retention" {
  description = "Number of snapshots to retain"
  value       = var.snapshot_enabled ? var.snapshot_retention : null
}

output "event_stream_integration_enabled" {
  description = "Whether Event Stream integration is enabled for DragonflyDB"
  value       = var.event_stream_integration
}

output "cache_invalidation_enabled" {
  description = "Whether cache invalidation is enabled for DragonflyDB"
  value       = var.cache_invalidation_enabled
}

output "cache_ttl" {
  description = "Default TTL for cached items in seconds"
  value       = var.cache_ttl
}

output "context_cache_ttl" {
  description = "TTL for context cache items in seconds"
  value       = var.context_cache_ttl
}

output "prometheus_integration_enabled" {
  description = "Whether Prometheus integration is enabled for DragonflyDB"
  value       = var.prometheus_integration
}

output "kata_container_integration_enabled" {
  description = "Whether Kata Containers integration is enabled for DragonflyDB"
  value       = var.kata_container_integration
}

output "context_rebuild_registry_enabled" {
  description = "Whether context rebuild registry is enabled for DragonflyDB"
  value       = var.event_stream_integration && var.cache_invalidation_enabled
}

output "context_rebuild_patterns" {
  description = "Patterns for context rebuild registry"
  value       = var.event_stream_integration && var.cache_invalidation_enabled ? [
    "agent:context:*",
    "tool:state:*",
    "execution:state:*"
  ] : []
}

output "go_client_example" {
  description = "Example code for connecting to DragonflyDB using Go"
  value       = <<-EOT
    // Example code for connecting to DragonflyDB using Go
    package main
    
    import (
      "context"
      "fmt"
      "time"
    
      "github.com/go-redis/redis/v8"
    )
    
    func main() {
      // Create a new Redis client
      client := redis.NewClient(&redis.Options{
        Addr:     "${kubernetes_service.dragonfly.metadata[0].name}.${kubernetes_service.dragonfly.metadata[0].namespace}.svc.cluster.local:${kubernetes_service.dragonfly.spec[0].port[0].port}",
        Password: "${var.password_enabled ? var.dragonfly_password : ""}",
        DB:       0,
      })
      
      ctx := context.Background()
      
      // Test connection
      pong, err := client.Ping(ctx).Result()
      if err != nil {
        fmt.Printf("Failed to connect to DragonflyDB: %s\n", err)
        return
      }
      fmt.Printf("Connected to DragonflyDB: %s\n", pong)
      
      // Set a key with TTL
      err = client.Set(ctx, "agent:context:example", "example-value", time.Duration(var.context_cache_ttl)*time.Second).Err()
      if err != nil {
        fmt.Printf("Failed to set key: %s\n", err)
        return
      }
      
      // Get a key
      val, err := client.Get(ctx, "agent:context:example").Result()
      if err != nil {
        fmt.Printf("Failed to get key: %s\n", err)
        return
      }
      fmt.Printf("Value: %s\n", val)
      
      // Subscribe to keyspace events (for cache invalidation)
      pubsub := client.PSubscribe(ctx, "__keyspace@0__:*")
      defer pubsub.Close()
      
      // Start a goroutine to handle messages
      go func() {
        for msg := range pubsub.Channel() {
          fmt.Printf("Channel: %s, Message: %s\n", msg.Channel, msg.Payload)
          // Handle cache invalidation here
        }
      }()
      
      // Wait for events
      time.Sleep(time.Second * 10)
    }
  EOT
}

output "event_stream_integration_example" {
  description = "Example code for integrating with Event Stream"
  value       = <<-EOT
    // Example code for integrating with Event Stream
    package main
    
    import (
      "context"
      "fmt"
      "time"
    
      "github.com/go-redis/redis/v8"
    )
    
    // ContextRebuildRegistry defines functions to rebuild context on cache misses
    type ContextRebuildRegistry struct {
      rebuildFunctions map[string]func(ctx context.Context, key string) (interface{}, error)
    }
    
    // NewContextRebuildRegistry creates a new registry
    func NewContextRebuildRegistry() *ContextRebuildRegistry {
      return &ContextRebuildRegistry{
        rebuildFunctions: make(map[string]func(ctx context.Context, key string) (interface{}, error)),
      }
    }
    
    // RegisterRebuildFunction registers a function to rebuild context for a pattern
    func (r *ContextRebuildRegistry) RegisterRebuildFunction(pattern string, fn func(ctx context.Context, key string) (interface{}, error)) {
      r.rebuildFunctions[pattern] = fn
    }
    
    // RebuildContext rebuilds context for a key
    func (r *ContextRebuildRegistry) RebuildContext(ctx context.Context, key string) (interface{}, error) {
      // Find the appropriate rebuild function
      for pattern, fn := range r.rebuildFunctions {
        // Check if the key matches the pattern
        // In a real implementation, you would use a proper pattern matching algorithm
        if pattern == "agent:context:*" && key[:14] == "agent:context:" {
          return fn(ctx, key)
        }
      }
      
      return nil, fmt.Errorf("no rebuild function found for key: %s", key)
    }
    
    // Example rebuild function for agent context
    func rebuildAgentContext(ctx context.Context, key string) (interface{}, error) {
      // In a real implementation, this would rebuild the context from source data
      return fmt.Sprintf("Rebuilt context for %s", key), nil
    }
    
    func main() {
      // Create a new Redis client
      client := redis.NewClient(&redis.Options{
        Addr:     "${kubernetes_service.dragonfly.metadata[0].name}.${kubernetes_service.dragonfly.metadata[0].namespace}.svc.cluster.local:${kubernetes_service.dragonfly.spec[0].port[0].port}",
        Password: "${var.password_enabled ? var.dragonfly_password : ""}",
        DB:       0,
      })
      
      ctx := context.Background()
      
      // Create a context rebuild registry
      registry := NewContextRebuildRegistry()
      
      // Register rebuild functions
      registry.RegisterRebuildFunction("agent:context:*", rebuildAgentContext)
      
      // Example of getting a value with cache miss handling
      key := "agent:context:example"
      val, err := client.Get(ctx, key).Result()
      if err == redis.Nil {
        // Cache miss, rebuild the context
        fmt.Printf("Cache miss for key: %s\n", key)
        
        // Rebuild the context
        newVal, err := registry.RebuildContext(ctx, key)
        if err != nil {
          fmt.Printf("Failed to rebuild context: %s\n", err)
          return
        }
        
        // Store the rebuilt context in the cache
        err = client.Set(ctx, key, newVal, time.Duration(var.context_cache_ttl)*time.Second).Err()
        if err != nil {
          fmt.Printf("Failed to set key: %s\n", err)
          return
        }
        
        fmt.Printf("Rebuilt and cached context for key: %s\n", key)
      } else if err != nil {
        fmt.Printf("Failed to get key: %s\n", err)
        return
      } else {
        fmt.Printf("Cache hit for key: %s, value: %s\n", key, val)
      }
    }
  EOT
}

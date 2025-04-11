
resource "helm_release" "jspolicy" {
  name       = var.name
  repository = "https://charts.loft.sh"
  chart      = "jspolicy"
  version    = var.chart_version
  namespace  = var.namespace
  
  set {
    name  = "image.repository"
    value = "loftsh/jspolicy"
  }
  
  set {
    name  = "image.tag"
    value = var.image_tag
  }
  
  set {
    name  = "replicas"
    value = var.high_availability ? 3 : 1
  }
  
  set {
    name  = "resources.limits.cpu"
    value = var.resources.limits.cpu
  }
  
  set {
    name  = "resources.limits.memory"
    value = var.resources.limits.memory
  }
  
  set {
    name  = "resources.requests.cpu"
    value = var.resources.requests.cpu
  }
  
  set {
    name  = "resources.requests.memory"
    value = var.resources.requests.memory
  }
  
  set {
    name  = "podSecurityContext.enabled"
    value = "true"
  }
  
  set {
    name  = "podSecurityContext.fsGroup"
    value = "1001"
  }
  
  set {
    name  = "containerSecurityContext.enabled"
    value = "true"
  }
  
  set {
    name  = "containerSecurityContext.runAsNonRoot"
    value = "true"
  }
  
  set {
    name  = "containerSecurityContext.runAsUser"
    value = "1001"
  }
  
  set {
    name  = "serviceAccount.create"
    value = "true"
  }
  
  set {
    name  = "rbac.create"
    value = "true"
  }
  
  set {
    name  = "webhook.enabled"
    value = "true"
  }
  
  set {
    name  = "webhook.failurePolicy"
    value = var.strict_mode ? "Fail" : "Ignore"
  }
  
  set {
    name  = "webhook.timeoutSeconds"
    value = "10"
  }
  
  set {
    name  = "webhook.namespaceSelector.matchExpressions[0].key"
    value = "kubernetes.io/metadata.name"
  }
  
  set {
    name  = "webhook.namespaceSelector.matchExpressions[0].operator"
    value = "NotIn"
  }
  
  set {
    name  = "webhook.namespaceSelector.matchExpressions[0].values[0]"
    value = "kube-system"
  }
  
  set {
    name  = "webhook.namespaceSelector.matchExpressions[0].values[1]"
    value = "jspolicy-system"
  }
  
  set {
    name  = "metrics.enabled"
    value = "true"
  }
  
  set {
    name  = "metrics.serviceMonitor.enabled"
    value = var.prometheus_enabled ? "true" : "false"
  }
}

resource "kubernetes_manifest" "resource_quota_policy" {
  manifest = {
    apiVersion = "policy.jspolicy.com/v1beta1"
    kind       = "Policy"
    metadata = {
      name = "resource-quota-policy"
    }
    spec = {
      operations = ["CREATE", "UPDATE"]
      resources  = ["Deployment", "StatefulSet", "DaemonSet"]
      javascript = <<-EOT
        function validate(request) {
          const object = request.object;
          
          // Skip system namespaces
          if (object.metadata.namespace === "kube-system" || 
              object.metadata.namespace === "jspolicy-system" ||
              object.metadata.namespace === "monitoring") {
            return { valid: true };
          }
          
          // Check for resource requests and limits
          const containers = object.spec.template.spec.containers || [];
          for (const container of containers) {
            if (!container.resources || !container.resources.requests || !container.resources.limits) {
              return {
                valid: false,
                message: `Container ${container.name} must specify resource requests and limits`
              };
            }
            
            // Check for CPU and memory requests and limits
            const requests = container.resources.requests || {};
            const limits = container.resources.limits || {};
            
            if (!requests.cpu || !requests.memory || !limits.cpu || !limits.memory) {
              return {
                valid: false,
                message: `Container ${container.name} must specify CPU and memory requests and limits`
              };
            }
          }
          
          return { valid: true };
        }
      EOT
    }
  }
  
  depends_on = [helm_release.jspolicy]
}

resource "kubernetes_manifest" "high_availability_policy" {
  manifest = {
    apiVersion = "policy.jspolicy.com/v1beta1"
    kind       = "Policy"
    metadata = {
      name = "high-availability-policy"
    }
    spec = {
      operations = ["CREATE", "UPDATE"]
      resources  = ["Deployment", "StatefulSet"]
      javascript = <<-EOT
        function validate(request) {
          const object = request.object;
          
          // Skip system namespaces
          if (object.metadata.namespace === "kube-system" || 
              object.metadata.namespace === "jspolicy-system" ||
              object.metadata.namespace === "monitoring") {
            return { valid: true };
          }
          
          // Check for production deployments
          if (object.metadata.labels && object.metadata.labels.environment === "production") {
            // Ensure high availability with multiple replicas
            if (object.kind === "Deployment" && (!object.spec.replicas || object.spec.replicas < 2)) {
              return {
                valid: false,
                message: "Production deployments must have at least 2 replicas for high availability"
              };
            }
            
            // Check for pod disruption budget
            if (!object.metadata.annotations || !object.metadata.annotations["pdb.kubernetes.io/configured"]) {
              return {
                valid: ${var.strict_mode ? "false" : "true"},
                message: "Warning: Production deployments should have a PodDisruptionBudget configured"
              };
            }
            
            // Check for pod anti-affinity
            const podSpec = object.spec.template.spec;
            if (!podSpec.affinity || 
                !podSpec.affinity.podAntiAffinity || 
                (!podSpec.affinity.podAntiAffinity.requiredDuringSchedulingIgnoredDuringExecution && 
                 !podSpec.affinity.podAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution)) {
              return {
                valid: ${var.strict_mode ? "false" : "true"},
                message: "Warning: Production deployments should use pod anti-affinity for high availability"
              };
            }
          }
          
          return { valid: true };
        }
      EOT
    }
  }
  
  depends_on = [helm_release.jspolicy]
}

resource "kubernetes_manifest" "security_policy" {
  manifest = {
    apiVersion = "policy.jspolicy.com/v1beta1"
    kind       = "Policy"
    metadata = {
      name = "security-policy"
    }
    spec = {
      operations = ["CREATE", "UPDATE"]
      resources  = ["Pod", "Deployment", "StatefulSet", "DaemonSet"]
      javascript = <<-EOT
        function validate(request) {
          const object = request.object;
          
          // Skip system namespaces
          if (object.metadata.namespace === "kube-system" || 
              object.metadata.namespace === "jspolicy-system") {
            return { valid: true };
          }
          
          let podSpec;
          if (object.kind === "Pod") {
            podSpec = object.spec;
          } else {
            podSpec = object.spec.template.spec;
          }
          
          // Check for privileged containers
          const containers = podSpec.containers || [];
          for (const container of containers) {
            if (container.securityContext && 
                container.securityContext.privileged === true) {
              return {
                valid: false,
                message: `Container ${container.name} must not be privileged`
              };
            }
            
            // Check for host network, PID, and IPC
            if (podSpec.hostNetwork === true || 
                podSpec.hostPID === true || 
                podSpec.hostIPC === true) {
              return {
                valid: false,
                message: "Pod must not use host network, PID, or IPC"
              };
            }
            
            // Check for capabilities
            if (container.securityContext && 
                container.securityContext.capabilities && 
                container.securityContext.capabilities.add) {
              const addedCaps = container.securityContext.capabilities.add;
              const dangerousCaps = ["ALL", "SYS_ADMIN", "NET_ADMIN"];
              
              for (const cap of dangerousCaps) {
                if (addedCaps.includes(cap)) {
                  return {
                    valid: false,
                    message: `Container ${container.name} must not add ${cap} capability`
                  };
                }
              }
            }
          }
          
          return { valid: true };
        }
      EOT
    }
  }
  
  depends_on = [helm_release.jspolicy]
}

resource "kubernetes_manifest" "network_policy" {
  manifest = {
    apiVersion = "policy.jspolicy.com/v1beta1"
    kind       = "Policy"
    metadata = {
      name = "network-policy-enforcer"
    }
    spec = {
      operations = ["CREATE"]
      resources  = ["Namespace"]
      javascript = <<-EOT
        function validate(request) {
          const object = request.object;
          
          // Skip system namespaces
          if (object.metadata.name === "kube-system" || 
              object.metadata.name === "jspolicy-system" ||
              object.metadata.name === "monitoring") {
            return { valid: true };
          }
          
          // Check if namespace has network policy label
          if (!object.metadata.labels || 
              !object.metadata.labels["network-policy"] || 
              object.metadata.labels["network-policy"] !== "enabled") {
            return {
              valid: ${var.strict_mode ? "false" : "true"},
              message: "Warning: Namespaces should have a network-policy=enabled label to indicate network policies are applied"
            };
          }
          
          return { valid: true };
        }
      EOT
    }
  }
  
  depends_on = [helm_release.jspolicy]
}

resource "kubernetes_manifest" "disaster_recovery_policy" {
  manifest = {
    apiVersion = "policy.jspolicy.com/v1beta1"
    kind       = "Policy"
    metadata = {
      name = "disaster-recovery-policy"
    }
    spec = {
      operations = ["CREATE", "UPDATE"]
      resources  = ["Deployment", "StatefulSet"]
      javascript = <<-EOT
        function validate(request) {
          const object = request.object;
          
          // Skip system namespaces
          if (object.metadata.namespace === "kube-system" || 
              object.metadata.namespace === "jspolicy-system") {
            return { valid: true };
          }
          
          // Check for production deployments
          if (object.metadata.labels && object.metadata.labels.environment === "production") {
            // Check for backup annotations
            if (!object.metadata.annotations || 
                !object.metadata.annotations["backup.velero.io/backup-volumes"]) {
              return {
                valid: ${var.strict_mode ? "false" : "true"},
                message: "Warning: Production deployments should have backup annotations for disaster recovery"
              };
            }
            
            // Check for StatefulSets with persistent volumes
            if (object.kind === "StatefulSet" && 
                object.spec.volumeClaimTemplates && 
                object.spec.volumeClaimTemplates.length > 0) {
              
              // Check for backup strategy
              if (!object.metadata.annotations || 
                  !object.metadata.annotations["backup.velero.io/backup-volumes"]) {
                return {
                  valid: ${var.strict_mode ? "false" : "true"},
                  message: "Warning: StatefulSets with persistent volumes should have backup annotations"
                };
              }
            }
          }
          
          return { valid: true };
        }
      EOT
    }
  }
  
  depends_on = [helm_release.jspolicy]
}

resource "kubernetes_manifest" "rollback_policy" {
  manifest = {
    apiVersion = "policy.jspolicy.com/v1beta1"
    kind       = "Policy"
    metadata = {
      name = "rollback-policy"
    }
    spec = {
      operations = ["CREATE", "UPDATE"]
      resources  = ["Deployment"]
      javascript = <<-EOT
        function validate(request) {
          const object = request.object;
          
          // Skip system namespaces
          if (object.metadata.namespace === "kube-system" || 
              object.metadata.namespace === "jspolicy-system") {
            return { valid: true };
          }
          
          // Check for production deployments
          if (object.metadata.labels && object.metadata.labels.environment === "production") {
            // Check for rollback strategy
            if (!object.spec.strategy || 
                !object.spec.strategy.rollingUpdate || 
                !object.spec.strategy.rollingUpdate.maxUnavailable || 
                !object.spec.strategy.rollingUpdate.maxSurge) {
              return {
                valid: ${var.strict_mode ? "false" : "true"},
                message: "Warning: Production deployments should have a rolling update strategy with maxUnavailable and maxSurge defined"
              };
            }
            
            // Check for revision history limit
            if (!object.spec.revisionHistoryLimit || object.spec.revisionHistoryLimit < 3) {
              return {
                valid: ${var.strict_mode ? "false" : "true"},
                message: "Warning: Production deployments should have a revisionHistoryLimit of at least 3 for rollback capability"
              };
            }
          }
          
          return { valid: true };
        }
      EOT
    }
  }
  
  depends_on = [helm_release.jspolicy]
}

resource "kubernetes_manifest" "resource_cleanup_policy" {
  manifest = {
    apiVersion = "policy.jspolicy.com/v1beta1"
    kind       = "Policy"
    metadata = {
      name = "resource-cleanup-policy"
    }
    spec = {
      operations = ["CREATE", "UPDATE"]
      resources  = ["Pod", "Job", "CronJob"]
      javascript = <<-EOT
        function validate(request) {
          const object = request.object;
          
          // Skip system namespaces
          if (object.metadata.namespace === "kube-system" || 
              object.metadata.namespace === "jspolicy-system") {
            return { valid: true };
          }
          
          if (object.kind === "Pod") {
            // Check for restart policy
            if (!object.spec.restartPolicy || object.spec.restartPolicy === "Always") {
              // Only warn for standalone pods, not those created by controllers
              if (!object.metadata.ownerReferences || object.metadata.ownerReferences.length === 0) {
                return {
                  valid: ${var.strict_mode ? "false" : "true"},
                  message: "Warning: Standalone pods should have a restart policy other than Always"
                };
              }
            }
          } else if (object.kind === "Job") {
            // Check for TTL after completion
            if (!object.spec.ttlSecondsAfterFinished) {
              return {
                valid: ${var.strict_mode ? "false" : "true"},
                message: "Warning: Jobs should have ttlSecondsAfterFinished set to automatically clean up completed jobs"
              };
            }
          } else if (object.kind === "CronJob") {
            // Check for history limits
            if (!object.spec.successfulJobsHistoryLimit || !object.spec.failedJobsHistoryLimit) {
              return {
                valid: ${var.strict_mode ? "false" : "true"},
                message: "Warning: CronJobs should have successfulJobsHistoryLimit and failedJobsHistoryLimit set"
              };
            }
          }
          
          return { valid: true };
        }
      EOT
    }
  }
  
  depends_on = [helm_release.jspolicy]
}

resource "kubernetes_manifest" "policy_bundle" {
  manifest = {
    apiVersion = "policy.jspolicy.com/v1beta1"
    kind       = "PolicyBundle"
    metadata = {
      name = "agent-runtime-policies"
    }
    spec = {
      policies = [
        "resource-quota-policy",
        "high-availability-policy",
        "security-policy",
        "network-policy-enforcer",
        "disaster-recovery-policy",
        "rollback-policy",
        "resource-cleanup-policy"
      ]
    }
  }
  
  depends_on = [
    kubernetes_manifest.resource_quota_policy,
    kubernetes_manifest.high_availability_policy,
    kubernetes_manifest.security_policy,
    kubernetes_manifest.network_policy,
    kubernetes_manifest.disaster_recovery_policy,
    kubernetes_manifest.rollback_policy,
    kubernetes_manifest.resource_cleanup_policy
  ]
}

resource "kubernetes_deployment" "recovery_controller" {
  count = var.enable_recovery_controller ? 1 : 0
  
  metadata {
    name      = "jspolicy-recovery-controller"
    namespace = var.namespace
    
    labels = {
      app = "jspolicy-recovery-controller"
    }
  }
  
  spec {
    replicas = var.high_availability ? 2 : 1
    
    selector {
      match_labels = {
        app = "jspolicy-recovery-controller"
      }
    }
    
    template {
      metadata {
        labels = {
          app = "jspolicy-recovery-controller"
        }
      }
      
      spec {
        service_account_name = kubernetes_service_account.recovery_controller[0].metadata[0].name
        
        container {
          name  = "controller"
          image = "loftsh/jspolicy-recovery-controller:${var.recovery_controller_version}"
          
          args = [
            "--v=2",
            "--namespace=${var.namespace}",
            "--recovery-interval=30s"
          ]
          
          resources {
            limits = {
              cpu    = "200m"
              memory = "256Mi"
            }
            
            requests = {
              cpu    = "100m"
              memory = "128Mi"
            }
          }
          
          liveness_probe {
            http_get {
              path = "/healthz"
              port = 8080
            }
            
            initial_delay_seconds = 30
            period_seconds        = 10
          }
          
          readiness_probe {
            http_get {
              path = "/readyz"
              port = 8080
            }
            
            initial_delay_seconds = 5
            period_seconds        = 10
          }
        }
      }
    }
  }
}

resource "kubernetes_service_account" "recovery_controller" {
  count = var.enable_recovery_controller ? 1 : 0
  
  metadata {
    name      = "jspolicy-recovery-controller"
    namespace = var.namespace
  }
}

resource "kubernetes_cluster_role" "recovery_controller" {
  count = var.enable_recovery_controller ? 1 : 0
  
  metadata {
    name = "jspolicy-recovery-controller"
  }
  
  rule {
    api_groups = [""]
    resources  = ["pods", "services", "configmaps", "secrets"]
    verbs      = ["get", "list", "watch"]
  }
  
  rule {
    api_groups = ["apps"]
    resources  = ["deployments", "statefulsets", "daemonsets"]
    verbs      = ["get", "list", "watch", "update", "patch"]
  }
  
  rule {
    api_groups = ["policy.jspolicy.com"]
    resources  = ["policies", "policybundles", "policyreports"]
    verbs      = ["get", "list", "watch"]
  }
  
  rule {
    api_groups = [""]
    resources  = ["events"]
    verbs      = ["create", "patch", "update"]
  }
}

resource "kubernetes_cluster_role_binding" "recovery_controller" {
  count = var.enable_recovery_controller ? 1 : 0
  
  metadata {
    name = "jspolicy-recovery-controller"
  }
  
  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "ClusterRole"
    name      = kubernetes_cluster_role.recovery_controller[0].metadata[0].name
  }
  
  subject {
    kind      = "ServiceAccount"
    name      = kubernetes_service_account.recovery_controller[0].metadata[0].name
    namespace = var.namespace
  }
}

# Kubernetes Module for Agent Runtime

locals {
  aws_node_sizes = {
    "small"  = "t3.medium"
    "medium" = "m5.large"
    "large"  = "m5.2xlarge"
    "xlarge" = "m5.4xlarge"
  }
  
  azure_node_sizes = {
    "small"  = "Standard_D2s_v3"
    "medium" = "Standard_D4s_v3"
    "large"  = "Standard_D8s_v3"
    "xlarge" = "Standard_D16s_v3"
  }
  
  ovh_node_sizes = {
    "small"  = "b2-7"
    "medium" = "b2-15"
    "large"  = "b2-30"
    "xlarge" = "b2-60"
  }
  
  fly_node_sizes = {
    "small"  = "dedicated-cpu-2x"
    "medium" = "dedicated-cpu-4x"
    "large"  = "dedicated-cpu-8x"
    "xlarge" = "dedicated-cpu-16x"
  }
}

resource "kubernetes_namespace" "agent_runtime" {
  metadata {
    name = var.namespace
    
    labels = {
      "app.kubernetes.io/name" = "agent-runtime"
      "app.kubernetes.io/part-of" = "agent-runtime"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }
}

resource "aws_eks_cluster" "this" {
  count = var.cloud_provider == "aws" ? 1 : 0
  
  name     = var.cluster_name
  role_arn = aws_iam_role.cluster[0].arn
  version  = var.kubernetes_version

  vpc_config {
    subnet_ids = aws_subnet.this[*].id
    security_group_ids = [aws_security_group.cluster[0].id]
    endpoint_private_access = true
    endpoint_public_access = true
  }

  encryption_config {
    provider {
      key_arn = aws_kms_key.eks[0].arn
    }
    resources = ["secrets"]
  }

  enabled_cluster_log_types = ["api", "audit", "authenticator", "controllerManager", "scheduler"]

  depends_on = [
    aws_iam_role_policy_attachment.cluster_policy[0],
    aws_iam_role_policy_attachment.service_policy[0],
    aws_cloudwatch_log_group.eks[0]
  ]

  lifecycle {
    prevent_destroy = true
  }
}

resource "aws_cloudwatch_log_group" "eks" {
  count = var.cloud_provider == "aws" ? 1 : 0
  
  name              = "/aws/eks/${var.cluster_name}/cluster"
  retention_in_days = 30
}

resource "aws_kms_key" "eks" {
  count = var.cloud_provider == "aws" ? 1 : 0
  
  description             = "EKS Secret Encryption Key"
  deletion_window_in_days = 7
  enable_key_rotation     = true
}

resource "aws_iam_role" "cluster" {
  count = var.cloud_provider == "aws" ? 1 : 0
  
  name = "${var.cluster_name}-cluster-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "eks.amazonaws.com"
        }
      },
    ]
  })
}

resource "aws_iam_role_policy_attachment" "cluster_policy" {
  count = var.cloud_provider == "aws" ? 1 : 0
  
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
  role       = aws_iam_role.cluster[0].name
}

resource "aws_iam_role_policy_attachment" "service_policy" {
  count = var.cloud_provider == "aws" ? 1 : 0
  
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSServicePolicy"
  role       = aws_iam_role.cluster[0].name
}

resource "aws_vpc" "this" {
  count = var.cloud_provider == "aws" ? 1 : 0
  
  cidr_block = "10.0.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true
  
  tags = {
    Name = "${var.cluster_name}-vpc"
  }
}

resource "aws_subnet" "this" {
  count = var.cloud_provider == "aws" ? 3 : 0
  
  vpc_id            = aws_vpc.this[0].id
  cidr_block        = "10.0.${count.index}.0/24"
  availability_zone = data.aws_availability_zones.available[0].names[count.index]
  map_public_ip_on_launch = true
  
  tags = {
    Name = "${var.cluster_name}-subnet-${count.index}"
    "kubernetes.io/cluster/${var.cluster_name}" = "shared"
  }
}

resource "aws_security_group" "cluster" {
  count = var.cloud_provider == "aws" ? 1 : 0
  
  name        = "${var.cluster_name}-cluster-sg"
  description = "Cluster security group"
  vpc_id      = aws_vpc.this[0].id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "${var.cluster_name}-cluster-sg"
  }
}

data "aws_availability_zones" "available" {
  count = var.cloud_provider == "aws" ? 1 : 0
  
  state = "available"
}

resource "aws_eks_node_group" "this" {
  count = var.cloud_provider == "aws" ? 1 : 0
  
  cluster_name    = aws_eks_cluster.this[0].name
  node_group_name = "${var.cluster_name}-node-group"
  node_role_arn   = aws_iam_role.node[0].arn
  subnet_ids      = aws_subnet.this[*].id
  instance_types  = [lookup(local.aws_node_sizes, var.node_size, "m5.large")]
  
  scaling_config {
    desired_size = var.node_count
    max_size     = var.node_count * 2
    min_size     = var.node_count >= 3 ? 3 : var.node_count
  }

  update_config {
    max_unavailable = 1
  }

  depends_on = [
    aws_iam_role_policy_attachment.node_policy[0],
    aws_iam_role_policy_attachment.cni_policy[0],
    aws_iam_role_policy_attachment.registry_policy[0],
  ]

  lifecycle {
    ignore_changes = [scaling_config[0].desired_size]
  }
}

resource "aws_iam_role" "node" {
  count = var.cloud_provider == "aws" ? 1 : 0
  
  name = "${var.cluster_name}-node-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ec2.amazonaws.com"
        }
      },
    ]
  })
}

resource "aws_iam_role_policy_attachment" "node_policy" {
  count = var.cloud_provider == "aws" ? 1 : 0
  
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy"
  role       = aws_iam_role.node[0].name
}

resource "aws_iam_role_policy_attachment" "cni_policy" {
  count = var.cloud_provider == "aws" ? 1 : 0
  
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKS-CNI-Policy"
  role       = aws_iam_role.node[0].name
}

resource "aws_iam_role_policy_attachment" "registry_policy" {
  count = var.cloud_provider == "aws" ? 1 : 0
  
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
  role       = aws_iam_role.node[0].name
}

resource "azurerm_resource_group" "this" {
  count = var.cloud_provider == "azure" ? 1 : 0
  
  name     = "${var.cluster_name}-rg"
  location = var.region
}

resource "azurerm_kubernetes_cluster" "this" {
  count = var.cloud_provider == "azure" ? 1 : 0
  
  name                = var.cluster_name
  location            = azurerm_resource_group.this[0].location
  resource_group_name = azurerm_resource_group.this[0].name
  dns_prefix          = var.cluster_name
  kubernetes_version  = var.kubernetes_version
  node_resource_group = "${var.cluster_name}-node-rg"

  default_node_pool {
    name                = "default"
    node_count          = var.node_count
    vm_size             = lookup(local.azure_node_sizes, var.node_size, "Standard_D4s_v3")
    availability_zones  = ["1", "2", "3"]
    enable_auto_scaling = true
    min_count           = var.node_count >= 3 ? 3 : var.node_count
    max_count           = var.node_count * 2
    os_disk_size_gb     = 100
    os_disk_type        = "Managed"
    vnet_subnet_id      = azurerm_subnet.this[0].id
  }

  identity {
    type = "SystemAssigned"
  }

  network_profile {
    network_plugin     = "azure"
    load_balancer_sku  = "standard"
    network_policy     = "calico"
    service_cidr       = "10.0.0.0/16"
    dns_service_ip     = "10.0.0.10"
    docker_bridge_cidr = "172.17.0.1/16"
  }

  role_based_access_control_enabled = true

  azure_active_directory_role_based_access_control {
    managed                = true
    admin_group_object_ids = var.admin_group_object_ids
  }

  addon_profile {
    oms_agent {
      enabled                    = true
      log_analytics_workspace_id = azurerm_log_analytics_workspace.this[0].id
    }
    azure_policy {
      enabled = true
    }
  }

  lifecycle {
    prevent_destroy = true
    ignore_changes  = [default_node_pool[0].node_count]
  }
}

resource "azurerm_virtual_network" "this" {
  count = var.cloud_provider == "azure" ? 1 : 0
  
  name                = "${var.cluster_name}-vnet"
  location            = azurerm_resource_group.this[0].location
  resource_group_name = azurerm_resource_group.this[0].name
  address_space       = ["10.1.0.0/16"]
}

resource "azurerm_subnet" "this" {
  count = var.cloud_provider == "azure" ? 1 : 0
  
  name                 = "${var.cluster_name}-subnet"
  resource_group_name  = azurerm_resource_group.this[0].name
  virtual_network_name = azurerm_virtual_network.this[0].name
  address_prefixes     = ["10.1.0.0/24"]
}

resource "azurerm_log_analytics_workspace" "this" {
  count = var.cloud_provider == "azure" ? 1 : 0
  
  name                = "${var.cluster_name}-workspace"
  location            = azurerm_resource_group.this[0].location
  resource_group_name = azurerm_resource_group.this[0].name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "ovh_cloud_project_kube" "this" {
  count = var.cloud_provider == "ovh" ? 1 : 0
  
  service_name = var.ovh_service_name
  name         = var.cluster_name
  region       = var.region
  version      = var.kubernetes_version
  
  private_network_id = ovh_cloud_project_network_private.this[0].id
}

resource "ovh_cloud_project_network_private" "this" {
  count = var.cloud_provider == "ovh" ? 1 : 0
  
  service_name = var.ovh_service_name
  name         = "${var.cluster_name}-network"
  regions      = [var.region]
  vlan_id      = 0
}

resource "ovh_cloud_project_network_private_subnet" "this" {
  count = var.cloud_provider == "ovh" ? 1 : 0
  
  service_name = var.ovh_service_name
  network_id   = ovh_cloud_project_network_private.this[0].id
  region       = var.region
  start        = "192.168.0.2"
  end          = "192.168.0.254"
  network      = "192.168.0.0/24"
  dhcp         = true
  no_gateway   = false
}

resource "ovh_cloud_project_kube_nodepool" "this" {
  count = var.cloud_provider == "ovh" ? 1 : 0
  
  service_name  = var.ovh_service_name
  kube_id       = ovh_cloud_project_kube.this[0].id
  name          = "${var.cluster_name}-pool"
  flavor_name   = lookup(local.ovh_node_sizes, var.node_size, "b2-7")
  desired_nodes = var.node_count
  max_nodes     = var.node_count * 2
  min_nodes     = var.node_count >= 3 ? 3 : var.node_count
  autoscale     = true
  
  monthly_billed = var.monthly_billed
}

resource "fly_app" "k3s_server" {
  count = var.cloud_provider == "fly" ? 1 : 0
  
  name = "${var.cluster_name}-server"
  
  deploy {
    strategy = "immediate"
  }
}

resource "fly_volume" "k3s_server_data" {
  count = var.cloud_provider == "fly" ? 1 : 0
  
  name   = "${var.cluster_name}-server-data"
  app    = fly_app.k3s_server[0].name
  size   = 50
  region = var.region
}

resource "fly_machine" "k3s_server" {
  count = var.cloud_provider == "fly" ? 3 : 0
  
  app    = fly_app.k3s_server[0].name
  region = var.region
  name   = "${var.cluster_name}-server-${count.index}"
  
  image  = "rancher/k3s:v${var.kubernetes_version}-k3s1"
  
  services = [
    {
      ports = [
        {
          port     = 6443
          handlers = ["tls"]
        }
      ]
      protocol = "tcp"
      internal_port = 6443
    }
  ]
  
  mounts = [
    {
      path   = "/var/lib/rancher/k3s"
      volume = fly_volume.k3s_server_data[0].id
    }
  ]
  
  cpus     = 4
  memory   = 8192
  
  cmd = [
    "server",
    "--tls-san", "${fly_app.k3s_server[0].name}.fly.dev",
    "--cluster-init",
    "--disable", "traefik",
    "--disable", "servicelb",
    "--disable", "local-storage",
    "--flannel-backend", "wireguard",
    "--node-taint", "CriticalAddonsOnly=true:NoExecute"
  ]
}

resource "fly_app" "k3s_agent" {
  count = var.cloud_provider == "fly" ? 1 : 0
  
  name = "${var.cluster_name}-agent"
  
  deploy {
    strategy = "immediate"
  }
}

resource "fly_machine" "k3s_agent" {
  count = var.cloud_provider == "fly" ? var.node_count : 0
  
  app    = fly_app.k3s_agent[0].name
  region = var.region
  name   = "${var.cluster_name}-agent-${count.index}"
  
  image  = "rancher/k3s:v${var.kubernetes_version}-k3s1"
  
  cpus   = lookup({
    "small"  = 2,
    "medium" = 4,
    "large"  = 8,
    "xlarge" = 16
  }, var.node_size, 4)
  
  memory = lookup({
    "small"  = 4096,
    "medium" = 8192,
    "large"  = 16384,
    "xlarge" = 32768
  }, var.node_size, 8192)
  
  cmd = [
    "agent",
    "--server", "https://${fly_app.k3s_server[0].name}.fly.dev:6443",
    "--token", var.k3s_token
  ]
  
  depends_on = [fly_machine.k3s_server]
}

resource "helm_release" "vcluster" {
  count = var.vcluster_enabled ? 1 : 0
  
  name       = "vcluster"
  repository = "https://charts.loft.sh"
  chart      = "vcluster"
  version    = var.vcluster_version
  namespace  = kubernetes_namespace.agent_runtime.metadata[0].name
  
  values = [
    file("${path.module}/../../k8s/vcluster/values.yaml")
  ]
  
  set {
    name  = "syncer.extraArgs"
    value = "{--tls-san=vcluster.${kubernetes_namespace.agent_runtime.metadata[0].name}}"
  }

  set {
    name  = "persistent"
    value = "true"
  }

  set {
    name  = "storage.persistence.enabled"
    value = "true"
  }

  set {
    name  = "storage.persistence.size"
    value = "10Gi"
  }

  set {
    name  = "sync.nodes.enabled"
    value = "true"
  }

  set {
    name  = "sync.ingresses.enabled"
    value = "true"
  }

  set {
    name  = "isolation.enabled"
    value = "true"
  }

  set {
    name  = "replicas"
    value = "3"
  }
}

resource "kubernetes_deployment" "vnode_runtime" {
  count = var.vnode_enabled ? 1 : 0
  
  metadata {
    name      = "vnode-runtime"
    namespace = kubernetes_namespace.agent_runtime.metadata[0].name
    labels = {
      app = "vnode-runtime"
    }
  }

  spec {
    replicas = 2

    selector {
      match_labels = {
        app = "vnode-runtime"
      }
    }

    template {
      metadata {
        labels = {
          app = "vnode-runtime"
        }
      }

      spec {
        container {
          image = "ghcr.io/loft-sh/vnode-runtime:0.0.1-alpha.1"
          name  = "vnode-runtime"

          env {
            name  = "VCLUSTER_NAME"
            value = "vcluster"
          }

          env {
            name  = "VCLUSTER_NAMESPACE"
            value = kubernetes_namespace.agent_runtime.metadata[0].name
          }

          resources {
            limits = {
              cpu    = "500m"
              memory = "512Mi"
            }
            requests = {
              cpu    = "250m"
              memory = "256Mi"
            }
          }
        }
      }
    }
  }
  
  depends_on = [helm_release.vcluster]
}

resource "kubernetes_manifest" "jspolicy" {
  count = var.jspolicy_enabled ? 1 : 0
  
  manifest = yamldecode(file("${path.module}/../../k8s/jspolicy/policies.yaml"))
  
  depends_on = [
    helm_release.vcluster
  ]
}

output "kubeconfig" {
  value = var.cloud_provider == "aws" ? aws_eks_cluster.this[0].kubeconfig : (
    var.cloud_provider == "azure" ? azurerm_kubernetes_cluster.this[0].kube_config_raw : (
      var.cloud_provider == "ovh" ? ovh_cloud_project_kube.this[0].kubeconfig : null
    )
  )
  sensitive = true
}

output "cluster_endpoint" {
  value = var.cloud_provider == "aws" ? aws_eks_cluster.this[0].endpoint : (
    var.cloud_provider == "azure" ? azurerm_kubernetes_cluster.this[0].kube_config.0.host : (
      var.cloud_provider == "ovh" ? ovh_cloud_project_kube.this[0].endpoint : (
        var.cloud_provider == "fly" ? "https://${fly_app.k3s_server[0].name}.fly.dev:6443" : null
      )
    )
  )
}

output "cluster_name" {
  value = var.cluster_name
}

output "cluster_region" {
  value = var.region
}

output "cluster_version" {
  value = var.kubernetes_version
}

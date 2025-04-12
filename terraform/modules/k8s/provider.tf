
terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 2.20.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = ">= 2.9.0"
    }
    aws = {
      source  = "hashicorp/aws"
      version = ">= 4.67.0"
    }
    azurerm = {
      source  = "hashicorp/azurerm"
      version = ">= 3.58.0"
    }
    ovh = {
      source  = "ovh/ovh"
      version = ">= 0.30.0"
    }
    fly = {
      source  = "fly-apps/fly"
      version = ">= 0.0.21"
    }
  }
  required_version = ">= 1.3.0"
}

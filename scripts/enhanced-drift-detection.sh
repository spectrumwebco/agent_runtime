#!/bin/bash

# Enhanced Drift Detection Script
# This script checks for drift between Kubernetes and Terraform configurations

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Starting enhanced drift detection between Kubernetes and Terraform...${NC}"

# Check if kubectl is available
if ! command -v kubectl &> /dev/null; then
    echo -e "${RED}Error: kubectl is not installed or not in PATH${NC}"
    exit 1
fi

# Check if terraform is available
if ! command -v terraform &> /dev/null; then
    echo -e "${RED}Error: terraform is not installed or not in PATH${NC}"
    exit 1
fi

# Get all namespaces from Kubernetes
k8s_namespaces=$(kubectl get namespaces -o jsonpath='{.items[*].metadata.name}')

# Get all modules from Terraform
cd terraform
terraform_modules=$(find modules -mindepth 1 -maxdepth 1 -type d | sed 's|modules/||')

# Check if each module is referenced in main.tf
echo -e "\n${YELLOW}Checking if all modules are referenced in main.tf...${NC}"
for module in $terraform_modules; do
    if grep -q "module \"$module\"" main.tf; then
        echo -e "${GREEN}✓ Module $module is referenced in main.tf${NC}"
    else
        echo -e "${RED}✗ Module $module is not referenced in main.tf${NC}"
    fi
done

# Check for each service in both Kubernetes and Terraform
echo -e "\n${YELLOW}Checking service existence in both Kubernetes and Terraform...${NC}"
for namespace in $k8s_namespaces; do
    # Skip system namespaces
    if [[ "$namespace" == "kube-system" || "$namespace" == "kube-public" || "$namespace" == "kube-node-lease" || "$namespace" == "default" ]]; then
        continue
    fi
    
    # Extract service name from namespace (remove -system suffix if present)
    service_name=${namespace%-system}
    
    # Check if service exists in Terraform modules
    if [[ " $terraform_modules " =~ " $service_name " ]]; then
        echo -e "${GREEN}✓ Service $service_name exists in both Kubernetes and Terraform${NC}"
        
        # Count resources in Kubernetes
        k8s_deployments=$(kubectl get deployments -n $namespace --no-headers 2>/dev/null | wc -l)
        k8s_services=$(kubectl get services -n $namespace --no-headers 2>/dev/null | wc -l)
        k8s_configmaps=$(kubectl get configmaps -n $namespace --no-headers 2>/dev/null | wc -l)
        
        # Count resources in Terraform
        tf_deployments=$(grep -c "kubernetes_deployment" modules/$service_name/main.tf 2>/dev/null || echo 0)
        tf_services=$(grep -c "kubernetes_service" modules/$service_name/main.tf 2>/dev/null || echo 0)
        tf_configmaps=$(grep -c "kubernetes_config_map" modules/$service_name/main.tf 2>/dev/null || echo 0)
        
        # Compare resource counts
        if [ "$k8s_deployments" -eq "$tf_deployments" ]; then
            echo -e "    ${GREEN}✓ Deployments: K8s=$k8s_deployments, TF=$tf_deployments${NC}"
        else
            echo -e "    ${RED}✗ Deployments: K8s=$k8s_deployments, TF=$tf_deployments${NC}"
        fi
        
        if [ "$k8s_services" -eq "$tf_services" ]; then
            echo -e "    ${GREEN}✓ Services: K8s=$k8s_services, TF=$tf_services${NC}"
        else
            echo -e "    ${RED}✗ Services: K8s=$k8s_services, TF=$tf_services${NC}"
        fi
        
        if [ "$k8s_configmaps" -eq "$tf_configmaps" ]; then
            echo -e "    ${GREEN}✓ ConfigMaps: K8s=$k8s_configmaps, TF=$tf_configmaps${NC}"
        else
            echo -e "    ${RED}✗ ConfigMaps: K8s=$k8s_configmaps, TF=$tf_configmaps${NC}"
        fi
        
        # Check if service is referenced in main.tf
        if grep -q "module \"$service_name\"" ../terraform/main.tf; then
            echo -e "    ${GREEN}✓ Service $service_name is referenced in main.tf${NC}"
        else
            echo -e "    ${RED}✗ Service $service_name is not referenced in main.tf${NC}"
        fi
    else
        echo -e "${RED}✗ Service $namespace exists in Kubernetes but not in Terraform${NC}"
    fi
done

# Check for Terraform modules without corresponding Kubernetes namespaces
echo -e "\n${YELLOW}Checking for Terraform modules without corresponding Kubernetes namespaces...${NC}"
for module in $terraform_modules; do
    # Check if namespace exists in Kubernetes
    if [[ " $k8s_namespaces " =~ " $module " || " $k8s_namespaces " =~ " $module-system " ]]; then
        continue
    else
        echo -e "${RED}✗ Module $module exists in Terraform but not in Kubernetes${NC}"
    fi
done

echo -e "\n${GREEN}Drift detection completed.${NC}"

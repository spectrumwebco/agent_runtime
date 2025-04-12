
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;36m'
NC='\033[0m' # No Color

extract_k8s_resources() {
  local dir=$1
  local output_file=$2
  
  echo "Extracting Kubernetes resources from $dir..."
  
  find "$dir" -name "*.yaml" -o -name "*.yml" | while read -r file; do
    grep -E "kind:|name:" "$file" | awk '{$1=$1};1' | paste -d " " - - | \
    sed 's/kind: //g' | sed 's/name: //g' | sed 's/  / /g' >> "$output_file"
  done
  
  sort "$output_file" -o "$output_file"
}

extract_tf_resources() {
  local dir=$1
  local output_file=$2
  
  echo "Extracting Terraform resources from $dir..."
  
  find "$dir" -name "*.tf" | while read -r file; do
    grep -n "resource \"" "$file" | while read -r line; do
      line_num=$(echo "$line" | cut -d':' -f1)
      resource_type=$(echo "$line" | sed 's/.*resource "\([^"]*\)".*/\1/')
      
      resource_name=$(sed -n "$((line_num+1))p" "$file" | grep -o '"[^"]*"' | head -1 | sed 's/"//g')
      
      case "$resource_type" in
        kubernetes_deployment)
          echo "Deployment $resource_name" >> "$output_file"
          ;;
        kubernetes_stateful_set)
          echo "StatefulSet $resource_name" >> "$output_file"
          ;;
        kubernetes_daemon_set)
          echo "DaemonSet $resource_name" >> "$output_file"
          ;;
        kubernetes_service)
          echo "Service $resource_name" >> "$output_file"
          ;;
        kubernetes_config_map)
          echo "ConfigMap $resource_name" >> "$output_file"
          ;;
        kubernetes_secret)
          echo "Secret $resource_name" >> "$output_file"
          ;;
        kubernetes_namespace)
          echo "Namespace $resource_name" >> "$output_file"
          ;;
        kubernetes_persistent_volume_claim)
          echo "PersistentVolumeClaim $resource_name" >> "$output_file"
          ;;
        kubernetes_persistent_volume)
          echo "PersistentVolume $resource_name" >> "$output_file"
          ;;
        kubernetes_role)
          echo "Role $resource_name" >> "$output_file"
          ;;
        kubernetes_role_binding)
          echo "RoleBinding $resource_name" >> "$output_file"
          ;;
        kubernetes_cluster_role)
          echo "ClusterRole $resource_name" >> "$output_file"
          ;;
        kubernetes_cluster_role_binding)
          echo "ClusterRoleBinding $resource_name" >> "$output_file"
          ;;
        kubernetes_service_account)
          echo "ServiceAccount $resource_name" >> "$output_file"
          ;;
        kubernetes_manifest)
          manifest_line=$(grep -A 5 -B 5 "resource \"kubernetes_manifest\"" "$file" | grep -o "yamldecode(file(.*yaml\")" | head -1)
          if [ -n "$manifest_line" ]; then
            yaml_file=$(echo "$manifest_line" | sed 's/yamldecode(file("\(.*\)")/\1/')
            if [[ "$yaml_file" == *"path.module"* ]]; then
              module_dir=$(dirname "$file")
              yaml_file=$(echo "$yaml_file" | sed "s|\${path.module}|$module_dir|g")
            else
              yaml_file="$REPO_ROOT/$yaml_file"
            fi
            
            if [ -f "$yaml_file" ]; then
              kind=$(grep "kind:" "$yaml_file" | head -1 | awk '{print $2}')
              name=$(grep "name:" "$yaml_file" | head -1 | awk '{print $2}')
              echo "$kind $name" >> "$output_file"
            fi
          fi
          ;;
        *)
          ;;
      esac
    done
    
    grep -n "helm_release" "$file" | while read -r line; do
      line_num=$(echo "$line" | cut -d':' -f1)
      
      release_name=$(sed -n "$((line_num+1))p" "$file" | grep -o '"[^"]*"' | head -1 | sed 's/"//g')
      
      chart_name=$(grep -A 10 "helm_release" "$file" | grep "chart" | head -1 | grep -o '"[^"]*"' | head -1 | sed 's/"//g')
      
      echo "HelmRelease $release_name ($chart_name)" >> "$output_file"
    done
  done
  
  sort "$output_file" -o "$output_file"
}

compare_resources() {
  local k8s_file=$1
  local tf_file=$2
  local component=$3
  
  echo -e "\n${BLUE}Comparing resources for $component...${NC}"
  
  local k8s_count=$(wc -l < "$k8s_file")
  local tf_count=$(wc -l < "$tf_file")
  
  echo -e "Kubernetes resources: $k8s_count"
  echo -e "Terraform resources: $tf_count"
  
  echo -e "\n${YELLOW}Resources in Kubernetes but not in Terraform:${NC}"
  comm -23 "$k8s_file" "$tf_file" | while read -r line; do
    echo -e "  $line"
  done
  
  echo -e "\n${YELLOW}Resources in Terraform but not in Kubernetes:${NC}"
  comm -13 "$k8s_file" "$tf_file" | while read -r line; do
    echo -e "  $line"
  done
  
  echo -e "\n${GREEN}Common resources:${NC}"
  comm -12 "$k8s_file" "$tf_file" | while read -r line; do
    echo -e "  $line"
  done
}

check_drift() {
  local k8s_dir=$1
  local tf_dir=$2
  local component=$3
  
  echo -e "\n${YELLOW}Checking drift for $component...${NC}"
  
  local k8s_resources=$(mktemp)
  local tf_resources=$(mktemp)
  
  extract_k8s_resources "$k8s_dir" "$k8s_resources"
  extract_tf_resources "$tf_dir" "$tf_resources"
  
  compare_resources "$k8s_resources" "$tf_resources" "$component"
  
  rm -f "$k8s_resources" "$tf_resources"
}

check_services() {
  echo -e "\n${YELLOW}Checking for services that exist in one but not the other...${NC}"
  
  local k8s_services=$(find "$REPO_ROOT/k8s" -mindepth 1 -maxdepth 1 -type d | sort)
  
  local tf_modules=$(find "$REPO_ROOT/terraform/modules" -mindepth 1 -maxdepth 1 -type d | sort)
  
  local k8s_service_names=$(mktemp)
  local tf_module_names=$(mktemp)
  
  for dir in $k8s_services; do
    basename "$dir" >> "$k8s_service_names"
  done
  
  for dir in $tf_modules; do
    basename "$dir" >> "$tf_module_names"
  done
  
  echo -e "\n${YELLOW}Services in Kubernetes but not in Terraform:${NC}"
  comm -23 <(sort "$k8s_service_names") <(sort "$tf_module_names") | while read -r line; do
    echo -e "  $line"
  done
  
  echo -e "\n${YELLOW}Services in Terraform but not in Kubernetes:${NC}"
  comm -13 <(sort "$k8s_service_names") <(sort "$tf_module_names") | while read -r line; do
    echo -e "  $line"
  done
  
  rm -f "$k8s_service_names" "$tf_module_names"
}

main() {
  echo -e "${YELLOW}Starting enhanced drift detection for Agent Runtime...${NC}"
  
  check_services
  
  check_drift "$REPO_ROOT/k8s/dragonfly" "$REPO_ROOT/terraform/modules/dragonfly" "DragonflyDB"
  
  check_drift "$REPO_ROOT/k8s/supabase" "$REPO_ROOT/terraform/modules/supabase" "Supabase"
  
  check_drift "$REPO_ROOT/k8s/monitoring/vector" "$REPO_ROOT/terraform/modules/monitoring" "Vector"
  
  check_drift "$REPO_ROOT/k8s/mcp" "$REPO_ROOT/terraform/modules/mcp" "MCP"
  
  check_drift "$REPO_ROOT/k8s/kata-containers" "$REPO_ROOT/terraform/modules/kata" "Kata Containers"
  
  echo -e "\n${GREEN}Enhanced drift detection completed${NC}"
}

main "$@"

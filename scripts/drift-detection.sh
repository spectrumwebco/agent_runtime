

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

validate_yaml() {
  local file=$1
  
  if [ ! -f "$file" ]; then
    echo -e "${RED}Error: File $file does not exist${NC}"
    return 1
  fi
  
  if command -v python3 &> /dev/null; then
    python3 -c "import yaml; yaml.safe_load(open('$file'))" 2>/dev/null
    if [ $? -ne 0 ]; then
      echo -e "${RED}Error: Invalid YAML in $file${NC}"
      return 1
    fi
  elif command -v python &> /dev/null; then
    python -c "import yaml; yaml.safe_load(open('$file'))" 2>/dev/null
    if [ $? -ne 0 ]; then
      echo -e "${RED}Error: Invalid YAML in $file${NC}"
      return 1
    fi
  else
    echo -e "${YELLOW}Warning: Python not found, skipping YAML validation for $file${NC}"
  fi
  
  return 0
}

check_drift() {
  local k8s_dir=$1
  local tf_dir=$2
  local component=$3
  
  echo -e "\n${YELLOW}Checking drift for $component...${NC}"
  
  for yaml_file in $(find "$k8s_dir" -name "*.yaml" -o -name "*.yml"); do
    echo -e "Validating $yaml_file..."
    validate_yaml "$yaml_file"
    if [ $? -ne 0 ]; then
      echo -e "${RED}Validation failed for $yaml_file${NC}"
      return 1
    fi
  done
  
  if [ ! -d "$tf_dir" ]; then
    echo -e "${RED}Error: Terraform directory $tf_dir does not exist${NC}"
    return 1
  fi
  
  local k8s_resources=$(grep -c "kind:" $(find "$k8s_dir" -name "*.yaml" -o -name "*.yml") 2>/dev/null || echo 0)
  local tf_resources=$(grep -c "resource" $(find "$tf_dir" -name "*.tf") 2>/dev/null || echo 0)
  
  echo -e "Kubernetes resources: $k8s_resources"
  echo -e "Terraform resources: $tf_resources"
  
  if [ $k8s_resources -lt $tf_resources ]; then
    echo -e "${RED}Drift detected: Kubernetes has fewer resources than Terraform${NC}"
    return 1
  elif [ $k8s_resources -gt $tf_resources ]; then
    echo -e "${YELLOW}Warning: Kubernetes has more resources than Terraform${NC}"
  else
    echo -e "${GREEN}No drift detected based on resource count${NC}"
  fi
  
  return 0
}

main() {
  echo -e "${YELLOW}Starting drift detection for Agent Runtime...${NC}"
  
  check_drift "$REPO_ROOT/k8s/dragonfly" "$REPO_ROOT/terraform/modules/dragonfly" "DragonflyDB"
  
  check_drift "$REPO_ROOT/k8s/supabase" "$REPO_ROOT/terraform/modules/supabase" "Supabase"
  
  check_drift "$REPO_ROOT/k8s/monitoring/vector" "$REPO_ROOT/terraform/modules/monitoring" "Vector"
  
  check_drift "$REPO_ROOT/k8s/mcp" "$REPO_ROOT/terraform/modules/mcp" "MCP"
  
  echo -e "\n${GREEN}Drift detection completed${NC}"
}

main "$@"

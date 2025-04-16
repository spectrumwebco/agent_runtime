
set -e

NAMESPACE="langsmith"
CONTEXT=""
KUBECONFIG=""
API_KEY=""
PROJECT_NAME="agent-runtime-verification"
SELF_HOSTED=true
VERBOSE=false

while [[ $# -gt 0 ]]; do
  key="$1"
  case $key in
    --namespace)
      NAMESPACE="$2"
      shift
      shift
      ;;
    --context)
      CONTEXT="$2"
      shift
      shift
      ;;
    --kubeconfig)
      KUBECONFIG="$2"
      shift
      shift
      ;;
    --api-key)
      API_KEY="$2"
      shift
      shift
      ;;
    --project-name)
      PROJECT_NAME="$2"
      shift
      shift
      ;;
    --hosted)
      SELF_HOSTED=false
      shift
      ;;
    --verbose)
      VERBOSE=true
      shift
      ;;
    *)
      echo "Unknown option: $1"
      exit 1
      ;;
  esac
done

KUBECTL_CMD="kubectl"
if [ -n "$CONTEXT" ]; then
  KUBECTL_CMD="$KUBECTL_CMD --context $CONTEXT"
fi
if [ -n "$KUBECONFIG" ]; then
  KUBECTL_CMD="$KUBECTL_CMD --kubeconfig $KUBECONFIG"
fi

log() {
  if [ "$VERBOSE" = true ]; then
    echo "$@"
  fi
}

command_exists() {
  command -v "$1" >/dev/null 2>&1
}

check_pod_running() {
  local pod_prefix=$1
  local namespace=$2
  
  log "Checking if pod with prefix '$pod_prefix' is running in namespace '$namespace'..."
  
  local pod_name=$($KUBECTL_CMD -n "$namespace" get pods -o jsonpath="{.items[?(@.metadata.name=~'$pod_prefix.*')].metadata.name}" 2>/dev/null)
  
  if [ -z "$pod_name" ]; then
    echo "Error: No pod with prefix '$pod_prefix' found in namespace '$namespace'"
    return 1
  fi
  
  local pod_status=$($KUBECTL_CMD -n "$namespace" get pod "$pod_name" -o jsonpath="{.status.phase}" 2>/dev/null)
  
  if [ "$pod_status" != "Running" ]; then
    echo "Error: Pod '$pod_name' is not running. Current status: $pod_status"
    return 1
  fi
  
  log "Pod '$pod_name' is running"
  return 0
}

check_service_exists() {
  local service_name=$1
  local namespace=$2
  
  log "Checking if service '$service_name' exists in namespace '$namespace'..."
  
  if ! $KUBECTL_CMD -n "$namespace" get service "$service_name" >/dev/null 2>&1; then
    echo "Error: Service '$service_name' not found in namespace '$namespace'"
    return 1
  fi
  
  log "Service '$service_name' exists"
  return 0
}

check_ingress_exists() {
  local ingress_name=$1
  local namespace=$2
  
  log "Checking if ingress '$ingress_name' exists in namespace '$namespace'..."
  
  if ! $KUBECTL_CMD -n "$namespace" get ingress "$ingress_name" >/dev/null 2>&1; then
    echo "Error: Ingress '$ingress_name' not found in namespace '$namespace'"
    return 1
  fi
  
  log "Ingress '$ingress_name' exists"
  return 0
}

check_langsmith_api() {
  local api_url=$1
  local api_key=$2
  
  log "Checking if LangSmith API is accessible at '$api_url'..."
  
  local response
  response=$(curl -s -o /dev/null -w "%{http_code}" -H "Authorization: Bearer $api_key" "$api_url/api/projects")
  
  if [ "$response" != "200" ]; then
    echo "Error: LangSmith API is not accessible. HTTP status code: $response"
    return 1
  fi
  
  log "LangSmith API is accessible"
  return 0
}

run_integration_test() {
  local api_url=$1
  local api_key=$2
  local project_name=$3
  
  log "Running LangSmith integration test..."
  
  export LANGCHAIN_TRACING_V2=true
  export LANGCHAIN_ENDPOINT="$api_url"
  export LANGCHAIN_API_KEY="$api_key"
  export LANGCHAIN_PROJECT="$project_name"
  export TEST_LANGSMITH=true
  
  SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
  TEST_SCRIPT="$SCRIPT_DIR/test_langsmith_integration.sh"
  
  if [ ! -f "$TEST_SCRIPT" ]; then
    echo "Error: Test script not found at '$TEST_SCRIPT'"
    return 1
  fi
  
  bash "$TEST_SCRIPT" --api-key "$api_key" --api-url "$api_url" --project-name "$project_name" ${SELF_HOSTED:+--self-hosted}
  
  local exit_code=$?
  if [ $exit_code -ne 0 ]; then
    echo "Error: LangSmith integration test failed with exit code $exit_code"
    return 1
  fi
  
  log "LangSmith integration test passed"
  return 0
}

echo "Verifying LangSmith deployment and integration with LangGraph..."

if [ "$SELF_HOSTED" = true ]; then
  echo "Verifying self-hosted LangSmith deployment in namespace '$NAMESPACE'..."
  
  if ! command_exists kubectl; then
    echo "Error: kubectl is not installed"
    exit 1
  fi
  
  if ! $KUBECTL_CMD get namespace "$NAMESPACE" >/dev/null 2>&1; then
    echo "Error: Namespace '$NAMESPACE' does not exist"
    exit 1
  fi
  
  check_pod_running "langsmith-api" "$NAMESPACE" || exit 1
  check_pod_running "langsmith-frontend" "$NAMESPACE" || exit 1
  
  check_service_exists "langsmith-api" "$NAMESPACE" || exit 1
  check_service_exists "langsmith-frontend" "$NAMESPACE" || exit 1
  
  check_ingress_exists "langsmith" "$NAMESPACE" || exit 1
  
  INGRESS_IP=$($KUBECTL_CMD -n "$NAMESPACE" get ingress langsmith -o jsonpath='{.status.loadBalancer.ingress[0].ip}' 2>/dev/null || echo "")
  INGRESS_HOSTNAME=$($KUBECTL_CMD -n "$NAMESPACE" get ingress langsmith -o jsonpath='{.status.loadBalancer.ingress[0].hostname}' 2>/dev/null || echo "")
  
  if [ -n "$INGRESS_IP" ]; then
    API_URL="http://$INGRESS_IP/api"
  elif [ -n "$INGRESS_HOSTNAME" ]; then
    API_URL="http://$INGRESS_HOSTNAME/api"
  else
    echo "No ingress IP or hostname found. Setting up port-forwarding for local access..."
    
    pkill -f "kubectl.*port-forward.*langsmith-api" || true
    
    $KUBECTL_CMD -n "$NAMESPACE" port-forward svc/langsmith-api 8000:8000 >/dev/null 2>&1 &
    PORT_FORWARD_PID=$!
    
    sleep 3
    
    API_URL="http://localhost:8000"
    
    trap 'kill $PORT_FORWARD_PID 2>/dev/null || true' EXIT
  fi
  
  if [ -z "$API_KEY" ]; then
    echo "No API key provided. Trying to get it from the secret..."
    
    API_KEY=$($KUBECTL_CMD -n "$NAMESPACE" get secret langsmith-secrets -o jsonpath='{.data.license-key}' 2>/dev/null | base64 --decode)
    
    if [ -z "$API_KEY" ]; then
      echo "Error: Could not get API key from secret"
      exit 1
    fi
  fi
else
  echo "Verifying hosted LangSmith integration..."
  
  API_URL="https://api.smith.langchain.com"
  
  if [ -z "$API_KEY" ]; then
    echo "Error: API key is required for hosted LangSmith"
    exit 1
  fi
fi

check_langsmith_api "$API_URL" "$API_KEY" || exit 1

run_integration_test "$API_URL" "$API_KEY" "$PROJECT_NAME" || exit 1

echo "LangSmith deployment and integration verification completed successfully!"
echo "LangSmith API URL: $API_URL"
echo "Project Name: $PROJECT_NAME"
echo "Self-hosted: $SELF_HOSTED"

exit 0

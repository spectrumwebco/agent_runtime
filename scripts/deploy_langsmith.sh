
set -e

NAMESPACE="langsmith"
CONTEXT=""
KUBECONFIG=""
LICENSE_KEY=""
DOMAIN="langsmith.example.com"
STORAGE_CLASS=""
DEPLOY_CLIENT=true
DEPLOY_SERVER=true

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
    --license-key)
      LICENSE_KEY="$2"
      shift
      shift
      ;;
    --domain)
      DOMAIN="$2"
      shift
      shift
      ;;
    --storage-class)
      STORAGE_CLASS="$2"
      shift
      shift
      ;;
    --client-only)
      DEPLOY_SERVER=false
      shift
      ;;
    --server-only)
      DEPLOY_CLIENT=false
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

if [ "$DEPLOY_SERVER" = true ]; then
  echo "Deploying LangSmith server components..."
  
  if [ -z "$LICENSE_KEY" ]; then
    echo "Warning: No LangSmith Enterprise license key provided."
    echo "You will need to provide a license key to use the self-hosted LangSmith server."
    echo "Continuing with deployment, but you will need to update the license key later."
  fi
  
  SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
  SETUP_SCRIPT="$SCRIPT_DIR/setup_langsmith_k8s.sh"
  
  SETUP_ARGS=""
  if [ -n "$NAMESPACE" ]; then
    SETUP_ARGS="$SETUP_ARGS --namespace $NAMESPACE"
  fi
  if [ -n "$STORAGE_CLASS" ]; then
    SETUP_ARGS="$SETUP_ARGS --storage-class $STORAGE_CLASS"
  fi
  if [ -n "$LICENSE_KEY" ]; then
    SETUP_ARGS="$SETUP_ARGS --license-key $LICENSE_KEY"
  fi
  
  POSTGRES_PASSWORD="$(openssl rand -base64 12)"
  SECRET_KEY="$(openssl rand -base64 32)"
  ADMIN_PASSWORD="$(openssl rand -base64 12)"
  
  SETUP_ARGS="$SETUP_ARGS --postgres-password $POSTGRES_PASSWORD --secret-key $SECRET_KEY --admin-password $ADMIN_PASSWORD"
  
  bash "$SETUP_SCRIPT" $SETUP_ARGS
  
  if [ -n "$DOMAIN" ] && [ "$DOMAIN" != "langsmith.example.com" ]; then
    echo "Updating ingress with domain: $DOMAIN"
    $KUBECTL_CMD -n "$NAMESPACE" get ingress langsmith -o yaml | \
      sed "s/host: langsmith.example.com/host: $DOMAIN/g" | \
      sed "s/- langsmith.example.com/- $DOMAIN/g" | \
      $KUBECTL_CMD apply -f -
  fi
  
  echo "LangSmith server components deployed successfully!"
  echo "Admin password: $ADMIN_PASSWORD"
  echo "PostgreSQL password: $POSTGRES_PASSWORD"
  
  echo "Waiting for all pods to be ready..."
  $KUBECTL_CMD -n "$NAMESPACE" wait --for=condition=ready pod --all --timeout=300s || true
  
  INGRESS_IP=$($KUBECTL_CMD -n "$NAMESPACE" get ingress langsmith -o jsonpath='{.status.loadBalancer.ingress[0].ip}' 2>/dev/null || echo "")
  INGRESS_HOSTNAME=$($KUBECTL_CMD -n "$NAMESPACE" get ingress langsmith -o jsonpath='{.status.loadBalancer.ingress[0].hostname}' 2>/dev/null || echo "")
  
  if [ -n "$INGRESS_IP" ]; then
    echo "LangSmith server is available at: http://$INGRESS_IP"
    echo "You may need to add an entry to your /etc/hosts file:"
    echo "$INGRESS_IP $DOMAIN"
  elif [ -n "$INGRESS_HOSTNAME" ]; then
    echo "LangSmith server is available at: http://$INGRESS_HOSTNAME"
  else
    echo "LangSmith server ingress is not yet available."
    echo "You can check the status with: $KUBECTL_CMD -n $NAMESPACE get ingress langsmith"
  fi
  
  echo "Setting up port forwarding for local access..."
  echo "Run the following command in a separate terminal to access LangSmith locally:"
  echo "$KUBECTL_CMD -n $NAMESPACE port-forward svc/langsmith-frontend 8080:80"
  echo "Then visit: http://localhost:8080"
fi

if [ "$DEPLOY_CLIENT" = true ]; then
  echo "Deploying LangSmith client components..."
  
  echo "Creating ConfigMap for LangSmith client configuration..."
  
  if [ "$DEPLOY_SERVER" = true ]; then
    API_URL="http://langsmith-api.$NAMESPACE.svc.cluster.local:8000"
  else
    API_URL="https://api.smith.langchain.com"
  fi
  
  cat <<EOF | $KUBECTL_CMD apply -f -
apiVersion: v1
kind: ConfigMap
metadata:
  name: langsmith-client-config
  namespace: $NAMESPACE
data:
  LANGCHAIN_TRACING_V2: "true"
  LANGCHAIN_ENDPOINT: "$API_URL"
  LANGCHAIN_PROJECT: "agent-runtime"
EOF
  
  if [ -n "$LICENSE_KEY" ]; then
    echo "Creating Secret for LangSmith API key..."
    cat <<EOF | $KUBECTL_CMD apply -f -
apiVersion: v1
kind: Secret
metadata:
  name: langsmith-client-secret
  namespace: $NAMESPACE
type: Opaque
stringData:
  LANGCHAIN_API_KEY: "$LICENSE_KEY"
EOF
  else
    echo "Warning: No LangSmith API key provided."
    echo "You will need to provide an API key to use the LangSmith client."
    echo "You can create a Secret manually with the following command:"
    echo "$KUBECTL_CMD -n $NAMESPACE create secret generic langsmith-client-secret --from-literal=LANGCHAIN_API_KEY=your-api-key"
  fi
  
  echo "LangSmith client components deployed successfully!"
  echo ""
  echo "To use LangSmith in your applications, add the following environment variables:"
  echo "- From ConfigMap 'langsmith-client-config':"
  echo "  - LANGCHAIN_TRACING_V2=true"
  echo "  - LANGCHAIN_ENDPOINT=$API_URL"
  echo "  - LANGCHAIN_PROJECT=agent-runtime"
  echo "- From Secret 'langsmith-client-secret':"
  echo "  - LANGCHAIN_API_KEY=<your-api-key>"
  echo ""
  echo "Example Kubernetes Deployment snippet:"
  echo "---"
  echo "apiVersion: apps/v1"
  echo "kind: Deployment"
  echo "metadata:"
  echo "  name: your-app"
  echo "  namespace: your-namespace"
  echo "spec:"
  echo "  template:"
  echo "    spec:"
  echo "      containers:"
  echo "      - name: your-container"
  echo "        image: your-image"
  echo "        envFrom:"
  echo "        - configMapRef:"
  echo "            name: langsmith-client-config"
  echo "        - secretRef:"
  echo "            name: langsmith-client-secret"
  echo "---"
fi

echo "LangSmith deployment complete!"


set -e

NAMESPACE="langsmith"
STORAGE_CLASS=""
LICENSE_KEY=""
POSTGRES_PASSWORD="$(openssl rand -base64 12)"
SECRET_KEY="$(openssl rand -base64 32)"
ADMIN_PASSWORD="$(openssl rand -base64 12)"

while [[ $# -gt 0 ]]; do
  key="$1"
  case $key in
    --namespace)
      NAMESPACE="$2"
      shift
      shift
      ;;
    --storage-class)
      STORAGE_CLASS="$2"
      shift
      shift
      ;;
    --license-key)
      LICENSE_KEY="$2"
      shift
      shift
      ;;
    --postgres-password)
      POSTGRES_PASSWORD="$2"
      shift
      shift
      ;;
    --secret-key)
      SECRET_KEY="$2"
      shift
      shift
      ;;
    --admin-password)
      ADMIN_PASSWORD="$2"
      shift
      shift
      ;;
    *)
      echo "Unknown option: $1"
      exit 1
      ;;
  esac
done

if [ -z "$LICENSE_KEY" ]; then
  echo "Error: LangSmith Enterprise license key is required."
  echo "Usage: $0 --license-key YOUR_LICENSE_KEY [options]"
  echo "Options:"
  echo "  --namespace NAMESPACE          Kubernetes namespace (default: langsmith)"
  echo "  --storage-class STORAGE_CLASS  Storage class for persistent volumes"
  echo "  --postgres-password PASSWORD   Password for PostgreSQL (default: random)"
  echo "  --secret-key SECRET_KEY        Secret key for LangSmith (default: random)"
  echo "  --admin-password PASSWORD      Admin password for LangSmith (default: random)"
  exit 1
fi

echo "Setting up Kubernetes environment for LangSmith self-hosting..."
echo "Namespace: $NAMESPACE"
if [ -n "$STORAGE_CLASS" ]; then
  echo "Storage Class: $STORAGE_CLASS"
fi

if ! kubectl get namespace "$NAMESPACE" &> /dev/null; then
  echo "Creating namespace $NAMESPACE..."
  kubectl create namespace "$NAMESPACE"
else
  echo "Namespace $NAMESPACE already exists."
fi

echo "Creating secrets..."
kubectl create secret generic langsmith-secrets \
  --namespace "$NAMESPACE" \
  --from-literal=database-url="postgresql://langsmith:$POSTGRES_PASSWORD@langsmith-postgres:5432/langsmith" \
  --from-literal=redis-url="redis://langsmith-redis:6379/0" \
  --from-literal=secret-key="$SECRET_KEY" \
  --from-literal=license-key="$LICENSE_KEY" \
  --from-literal=postgres-password="$POSTGRES_PASSWORD" \
  --from-literal=admin-password="$ADMIN_PASSWORD" \
  --dry-run=client -o yaml | kubectl apply -f -

echo "Creating persistent volume claims..."
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: langsmith-redis-data
  namespace: $NAMESPACE
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 5Gi
  $([ -n "$STORAGE_CLASS" ] && echo "storageClassName: $STORAGE_CLASS")
EOF

echo "Applying Kubernetes manifests..."
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(dirname "$SCRIPT_DIR")"

for file in "$REPO_ROOT"/kubernetes/langsmith/*.yaml; do
  sed "s/namespace: langsmith/namespace: $NAMESPACE/g" "$file" | kubectl apply -f -
done

echo "Waiting for pods to be ready..."
kubectl wait --for=condition=ready pod -l app=langsmith -n "$NAMESPACE" --timeout=300s || true

echo "LangSmith Kubernetes environment setup complete!"
echo "PostgreSQL password: $POSTGRES_PASSWORD"
echo "Admin password: $ADMIN_PASSWORD"
echo ""
echo "To access LangSmith, set up an ingress or port-forward:"
echo "kubectl port-forward svc/langsmith-frontend -n $NAMESPACE 8080:80"
echo ""
echo "Then visit: http://localhost:8080"
echo ""
echo "To use LangSmith with the Agent Runtime system, set the following environment variables:"
echo "export LANGCHAIN_TRACING_V2=true"
echo "export LANGCHAIN_ENDPOINT=http://langsmith-api.$NAMESPACE.svc.cluster.local:8000"
echo "export LANGCHAIN_API_KEY=$LICENSE_KEY"
echo "export LANGCHAIN_PROJECT=agent-runtime"
echo ""
echo "For external access, update the ingress configuration in kubernetes/langsmith/ingress.yaml"
echo "with your domain and TLS configuration, then reapply the manifest."

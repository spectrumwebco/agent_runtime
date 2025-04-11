# jsPolicy Configuration

This directory contains the configuration for jsPolicy, a policy engine for Kubernetes from Loft.

## Policies

The following policies are implemented:

1. **Resource Validation**: Ensures all pods and deployments have resource limits defined
2. **Security Context Validation**: Ensures pods run as non-root users
3. **Kata Containers Validation**: Ensures sandbox pods use the kata-containers runtime class
4. **High Availability Validation**: Ensures critical deployments have multiple replicas

## Usage

To deploy jsPolicy:

```bash
kubectl apply -f policies.yaml
```

## Integration with Agent Runtime

jsPolicy provides policy enforcement for the Agent Runtime Kubernetes deployment, ensuring:

1. Resource constraints are properly defined
2. Security best practices are followed
3. Sandbox environments use Kata Containers
4. Critical components are highly available

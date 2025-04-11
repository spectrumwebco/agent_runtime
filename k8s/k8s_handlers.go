package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func RegisterKubernetesTools(mcpServer *server.MCPServer, clientset *kubernetes.Clientset) {
	mcpServer.RegisterTool(mcp.Tool{
		Name:        "list_pods",
		Description: "List pods in a namespace",
		Parameters: []mcp.Parameter{
			{
				Name:        "namespace",
				Description: "Kubernetes namespace",
				Type:        "string",
				Required:    true,
			},
		},
		Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			return listPods(ctx, clientset, params)
		},
	})

	mcpServer.RegisterTool(mcp.Tool{
		Name:        "get_pod",
		Description: "Get details of a specific pod",
		Parameters: []mcp.Parameter{
			{
				Name:        "namespace",
				Description: "Kubernetes namespace",
				Type:        "string",
				Required:    true,
			},
			{
				Name:        "name",
				Description: "Pod name",
				Type:        "string",
				Required:    true,
			},
		},
		Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			return getPod(ctx, clientset, params)
		},
	})

	mcpServer.RegisterTool(mcp.Tool{
		Name:        "delete_pod",
		Description: "Delete a pod",
		Parameters: []mcp.Parameter{
			{
				Name:        "namespace",
				Description: "Kubernetes namespace",
				Type:        "string",
				Required:    true,
			},
			{
				Name:        "name",
				Description: "Pod name",
				Type:        "string",
				Required:    true,
			},
		},
		Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			return deletePod(ctx, clientset, params)
		},
	})

	mcpServer.RegisterTool(mcp.Tool{
		Name:        "list_deployments",
		Description: "List deployments in a namespace",
		Parameters: []mcp.Parameter{
			{
				Name:        "namespace",
				Description: "Kubernetes namespace",
				Type:        "string",
				Required:    true,
			},
		},
		Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			return listDeployments(ctx, clientset, params)
		},
	})

	mcpServer.RegisterTool(mcp.Tool{
		Name:        "get_deployment",
		Description: "Get details of a specific deployment",
		Parameters: []mcp.Parameter{
			{
				Name:        "namespace",
				Description: "Kubernetes namespace",
				Type:        "string",
				Required:    true,
			},
			{
				Name:        "name",
				Description: "Deployment name",
				Type:        "string",
				Required:    true,
			},
		},
		Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			return getDeployment(ctx, clientset, params)
		},
	})

	mcpServer.RegisterTool(mcp.Tool{
		Name:        "create_deployment",
		Description: "Create a new deployment",
		Parameters: []mcp.Parameter{
			{
				Name:        "namespace",
				Description: "Kubernetes namespace",
				Type:        "string",
				Required:    true,
			},
			{
				Name:        "name",
				Description: "Deployment name",
				Type:        "string",
				Required:    true,
			},
			{
				Name:        "image",
				Description: "Container image",
				Type:        "string",
				Required:    true,
			},
			{
				Name:        "replicas",
				Description: "Number of replicas",
				Type:        "number",
				Required:    true,
			},
			{
				Name:        "labels",
				Description: "Labels for the deployment",
				Type:        "object",
				Required:    false,
			},
		},
		Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			return createDeployment(ctx, clientset, params)
		},
	})

	mcpServer.RegisterTool(mcp.Tool{
		Name:        "update_deployment",
		Description: "Update an existing deployment",
		Parameters: []mcp.Parameter{
			{
				Name:        "namespace",
				Description: "Kubernetes namespace",
				Type:        "string",
				Required:    true,
			},
			{
				Name:        "name",
				Description: "Deployment name",
				Type:        "string",
				Required:    true,
			},
			{
				Name:        "image",
				Description: "Container image",
				Type:        "string",
				Required:    false,
			},
			{
				Name:        "replicas",
				Description: "Number of replicas",
				Type:        "number",
				Required:    false,
			},
		},
		Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			return updateDeployment(ctx, clientset, params)
		},
	})

	mcpServer.RegisterTool(mcp.Tool{
		Name:        "delete_deployment",
		Description: "Delete a deployment",
		Parameters: []mcp.Parameter{
			{
				Name:        "namespace",
				Description: "Kubernetes namespace",
				Type:        "string",
				Required:    true,
			},
			{
				Name:        "name",
				Description: "Deployment name",
				Type:        "string",
				Required:    true,
			},
		},
		Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			return deleteDeployment(ctx, clientset, params)
		},
	})

	mcpServer.RegisterTool(mcp.Tool{
		Name:        "scale_deployment",
		Description: "Scale a deployment",
		Parameters: []mcp.Parameter{
			{
				Name:        "namespace",
				Description: "Kubernetes namespace",
				Type:        "string",
				Required:    true,
			},
			{
				Name:        "name",
				Description: "Deployment name",
				Type:        "string",
				Required:    true,
			},
			{
				Name:        "replicas",
				Description: "Number of replicas",
				Type:        "number",
				Required:    true,
			},
		},
		Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			return scaleDeployment(ctx, clientset, params)
		},
	})

	mcpServer.RegisterTool(mcp.Tool{
		Name:        "list_services",
		Description: "List services in a namespace",
		Parameters: []mcp.Parameter{
			{
				Name:        "namespace",
				Description: "Kubernetes namespace",
				Type:        "string",
				Required:    true,
			},
		},
		Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			return listServices(ctx, clientset, params)
		},
	})

	mcpServer.RegisterTool(mcp.Tool{
		Name:        "create_service",
		Description: "Create a new service",
		Parameters: []mcp.Parameter{
			{
				Name:        "namespace",
				Description: "Kubernetes namespace",
				Type:        "string",
				Required:    true,
			},
			{
				Name:        "name",
				Description: "Service name",
				Type:        "string",
				Required:    true,
			},
			{
				Name:        "selector",
				Description: "Label selector for pods",
				Type:        "object",
				Required:    true,
			},
			{
				Name:        "ports",
				Description: "Port mappings",
				Type:        "array",
				Required:    true,
			},
			{
				Name:        "type",
				Description: "Service type (ClusterIP, NodePort, LoadBalancer)",
				Type:        "string",
				Required:    false,
			},
		},
		Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			return createService(ctx, clientset, params)
		},
	})

	mcpServer.RegisterTool(mcp.Tool{
		Name:        "delete_service",
		Description: "Delete a service",
		Parameters: []mcp.Parameter{
			{
				Name:        "namespace",
				Description: "Kubernetes namespace",
				Type:        "string",
				Required:    true,
			},
			{
				Name:        "name",
				Description: "Service name",
				Type:        "string",
				Required:    true,
			},
		},
		Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			return deleteService(ctx, clientset, params)
		},
	})

	mcpServer.RegisterTool(mcp.Tool{
		Name:        "list_namespaces",
		Description: "List all namespaces",
		Parameters:  []mcp.Parameter{},
		Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			return listNamespaces(ctx, clientset, params)
		},
	})

	mcpServer.RegisterTool(mcp.Tool{
		Name:        "create_namespace",
		Description: "Create a new namespace",
		Parameters: []mcp.Parameter{
			{
				Name:        "name",
				Description: "Namespace name",
				Type:        "string",
				Required:    true,
			},
			{
				Name:        "labels",
				Description: "Labels for the namespace",
				Type:        "object",
				Required:    false,
			},
		},
		Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			return createNamespace(ctx, clientset, params)
		},
	})

	mcpServer.RegisterTool(mcp.Tool{
		Name:        "delete_namespace",
		Description: "Delete a namespace",
		Parameters: []mcp.Parameter{
			{
				Name:        "name",
				Description: "Namespace name",
				Type:        "string",
				Required:    true,
			},
		},
		Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			return deleteNamespace(ctx, clientset, params)
		},
	})
}

func listPods(ctx context.Context, clientset *kubernetes.Clientset, params map[string]interface{}) (interface{}, error) {
	namespace := params["namespace"].(string)
	pods, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}

	var result []map[string]interface{}
	for _, pod := range pods.Items {
		result = append(result, map[string]interface{}{
			"name":      pod.Name,
			"namespace": pod.Namespace,
			"status":    pod.Status.Phase,
			"created":   pod.CreationTimestamp.Time,
		})
	}

	return result, nil
}

func getPod(ctx context.Context, clientset *kubernetes.Clientset, params map[string]interface{}) (interface{}, error) {
	namespace := params["namespace"].(string)
	name := params["name"].(string)

	pod, err := clientset.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get pod: %w", err)
	}

	return map[string]interface{}{
		"name":       pod.Name,
		"namespace":  pod.Namespace,
		"status":     pod.Status.Phase,
		"created":    pod.CreationTimestamp.Time,
		"containers": pod.Spec.Containers,
		"labels":     pod.Labels,
	}, nil
}

func deletePod(ctx context.Context, clientset *kubernetes.Clientset, params map[string]interface{}) (interface{}, error) {
	namespace := params["namespace"].(string)
	name := params["name"].(string)

	err := clientset.CoreV1().Pods(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to delete pod: %w", err)
	}

	return map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Pod %s in namespace %s deleted", name, namespace),
	}, nil
}

func listDeployments(ctx context.Context, clientset *kubernetes.Clientset, params map[string]interface{}) (interface{}, error) {
	namespace := params["namespace"].(string)
	deployments, err := clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments: %w", err)
	}

	var result []map[string]interface{}
	for _, deployment := range deployments.Items {
		result = append(result, map[string]interface{}{
			"name":      deployment.Name,
			"namespace": deployment.Namespace,
			"replicas":  deployment.Status.Replicas,
			"available": deployment.Status.AvailableReplicas,
			"created":   deployment.CreationTimestamp.Time,
		})
	}

	return result, nil
}

func getDeployment(ctx context.Context, clientset *kubernetes.Clientset, params map[string]interface{}) (interface{}, error) {
	namespace := params["namespace"].(string)
	name := params["name"].(string)

	deployment, err := clientset.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment: %w", err)
	}

	return map[string]interface{}{
		"name":       deployment.Name,
		"namespace":  deployment.Namespace,
		"replicas":   deployment.Status.Replicas,
		"available":  deployment.Status.AvailableReplicas,
		"created":    deployment.CreationTimestamp.Time,
		"containers": deployment.Spec.Template.Spec.Containers,
		"labels":     deployment.Labels,
	}, nil
}

func createDeployment(ctx context.Context, clientset *kubernetes.Clientset, params map[string]interface{}) (interface{}, error) {
	namespace := params["namespace"].(string)
	name := params["name"].(string)
	image := params["image"].(string)
	replicas := int32(params["replicas"].(float64))

	labels := map[string]string{
		"app": name,
	}
	if params["labels"] != nil {
		labelsMap, ok := params["labels"].(map[string]interface{})
		if ok {
			for k, v := range labelsMap {
				labels[k] = v.(string)
			}
		}
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  name,
							Image: image,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 8080,
								},
							},
						},
					},
				},
			},
		},
	}

	result, err := clientset.AppsV1().Deployments(namespace).Create(ctx, deployment, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create deployment: %w", err)
	}

	return map[string]interface{}{
		"name":      result.Name,
		"namespace": result.Namespace,
		"replicas":  result.Spec.Replicas,
		"created":   result.CreationTimestamp.Time,
	}, nil
}

func updateDeployment(ctx context.Context, clientset *kubernetes.Clientset, params map[string]interface{}) (interface{}, error) {
	namespace := params["namespace"].(string)
	name := params["name"].(string)

	deployment, err := clientset.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment: %w", err)
	}

	if params["image"] != nil {
		image := params["image"].(string)
		for i := range deployment.Spec.Template.Spec.Containers {
			deployment.Spec.Template.Spec.Containers[i].Image = image
		}
	}

	if params["replicas"] != nil {
		replicas := int32(params["replicas"].(float64))
		deployment.Spec.Replicas = &replicas
	}

	result, err := clientset.AppsV1().Deployments(namespace).Update(ctx, deployment, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to update deployment: %w", err)
	}

	return map[string]interface{}{
		"name":      result.Name,
		"namespace": result.Namespace,
		"replicas":  result.Spec.Replicas,
		"updated":   result.CreationTimestamp.Time,
	}, nil
}

func deleteDeployment(ctx context.Context, clientset *kubernetes.Clientset, params map[string]interface{}) (interface{}, error) {
	namespace := params["namespace"].(string)
	name := params["name"].(string)

	err := clientset.AppsV1().Deployments(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to delete deployment: %w", err)
	}

	return map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Deployment %s in namespace %s deleted", name, namespace),
	}, nil
}

func scaleDeployment(ctx context.Context, clientset *kubernetes.Clientset, params map[string]interface{}) (interface{}, error) {
	namespace := params["namespace"].(string)
	name := params["name"].(string)
	replicas := int32(params["replicas"].(float64))

	deployment, err := clientset.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment: %w", err)
	}

	deployment.Spec.Replicas = &replicas

	result, err := clientset.AppsV1().Deployments(namespace).Update(ctx, deployment, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to scale deployment: %w", err)
	}

	return map[string]interface{}{
		"name":      result.Name,
		"namespace": result.Namespace,
		"replicas":  result.Spec.Replicas,
		"scaled":    true,
	}, nil
}

func listServices(ctx context.Context, clientset *kubernetes.Clientset, params map[string]interface{}) (interface{}, error) {
	namespace := params["namespace"].(string)
	services, err := clientset.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	var result []map[string]interface{}
	for _, service := range services.Items {
		result = append(result, map[string]interface{}{
			"name":      service.Name,
			"namespace": service.Namespace,
			"type":      service.Spec.Type,
			"clusterIP": service.Spec.ClusterIP,
			"created":   service.CreationTimestamp.Time,
		})
	}

	return result, nil
}

func createService(ctx context.Context, clientset *kubernetes.Clientset, params map[string]interface{}) (interface{}, error) {
	namespace := params["namespace"].(string)
	name := params["name"].(string)
	
	selectorMap, ok := params["selector"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid selector format")
	}
	
	selector := make(map[string]string)
	for k, v := range selectorMap {
		selector[k] = v.(string)
	}
	
	portsArray, ok := params["ports"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid ports format")
	}
	
	var ports []corev1.ServicePort
	for _, portObj := range portsArray {
		portMap, ok := portObj.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid port format")
		}
		
		port := corev1.ServicePort{
			Name:     portMap["name"].(string),
			Protocol: corev1.ProtocolTCP,
		}
		
		if portVal, ok := portMap["port"]; ok {
			switch v := portVal.(type) {
			case float64:
				port.Port = int32(v)
			case string:
				portInt, err := strconv.Atoi(v)
				if err != nil {
					return nil, fmt.Errorf("invalid port number: %w", err)
				}
				port.Port = int32(portInt)
			default:
				return nil, fmt.Errorf("invalid port type")
			}
		}
		
		if targetPortVal, ok := portMap["targetPort"]; ok {
			switch v := targetPortVal.(type) {
			case float64:
				port.TargetPort.IntVal = int32(v)
			case string:
				port.TargetPort.StrVal = v
			default:
				return nil, fmt.Errorf("invalid targetPort type")
			}
		}
		
		ports = append(ports, port)
	}
	
	serviceType := corev1.ServiceTypeClusterIP
	if params["type"] != nil {
		typeStr := params["type"].(string)
		switch typeStr {
		case "ClusterIP":
			serviceType = corev1.ServiceTypeClusterIP
		case "NodePort":
			serviceType = corev1.ServiceTypeNodePort
		case "LoadBalancer":
			serviceType = corev1.ServiceTypeLoadBalancer
		default:
			return nil, fmt.Errorf("invalid service type: %s", typeStr)
		}
	}
	
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: selector,
			Ports:    ports,
			Type:     serviceType,
		},
	}
	
	result, err := clientset.CoreV1().Services(namespace).Create(ctx, service, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create service: %w", err)
	}
	
	return map[string]interface{}{
		"name":      result.Name,
		"namespace": result.Namespace,
		"type":      result.Spec.Type,
		"clusterIP": result.Spec.ClusterIP,
		"created":   result.CreationTimestamp.Time,
	}, nil
}

func deleteService(ctx context.Context, clientset *kubernetes.Clientset, params map[string]interface{}) (interface{}, error) {
	namespace := params["namespace"].(string)
	name := params["name"].(string)
	
	err := clientset.CoreV1().Services(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to delete service: %w", err)
	}
	
	return map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Service %s in namespace %s deleted", name, namespace),
	}, nil
}

func listNamespaces(ctx context.Context, clientset *kubernetes.Clientset, params map[string]interface{}) (interface{}, error) {
	namespaces, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %w", err)
	}
	
	var result []map[string]interface{}
	for _, namespace := range namespaces.Items {
		result = append(result, map[string]interface{}{
			"name":    namespace.Name,
			"status":  namespace.Status.Phase,
			"created": namespace.CreationTimestamp.Time,
		})
	}
	
	return result, nil
}

func createNamespace(ctx context.Context, clientset *kubernetes.Clientset, params map[string]interface{}) (interface{}, error) {
	name := params["name"].(string)
	
	labels := map[string]string{}
	if params["labels"] != nil {
		labelsMap, ok := params["labels"].(map[string]interface{})
		if ok {
			for k, v := range labelsMap {
				labels[k] = v.(string)
			}
		}
	}
	
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: labels,
		},
	}
	
	result, err := clientset.CoreV1().Namespaces().Create(ctx, namespace, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create namespace: %w", err)
	}
	
	return map[string]interface{}{
		"name":    result.Name,
		"status":  result.Status.Phase,
		"created": result.CreationTimestamp.Time,
	}, nil
}

func deleteNamespace(ctx context.Context, clientset *kubernetes.Clientset, params map[string]interface{}) (interface{}, error) {
	name := params["name"].(string)
	
	err := clientset.CoreV1().Namespaces().Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to delete namespace: %w", err)
	}
	
	return map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Namespace %s deleted", name),
	}, nil
}

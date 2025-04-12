
output "jupyterhub_namespace" {
  description = "Namespace where JupyterHub is deployed"
  value       = kubernetes_namespace.jupyterhub.metadata[0].name
}

output "jupyterhub_service" {
  description = "Name of JupyterHub service"
  value       = "jupyterhub.${kubernetes_namespace.jupyterhub.metadata[0].name}.svc.cluster.local"
}

output "jupyter_notebooks_config_map" {
  description = "Name of ConfigMap for Jupyter notebooks"
  value       = kubernetes_config_map.jupyter_notebooks.metadata[0].name
}

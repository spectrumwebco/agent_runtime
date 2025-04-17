apiVersion: pipecd.dev/v1beta1
kind: Piped
spec:
  projectID: "${project_id}"
  pipedID: "${piped_id}"
  pipedKeyFile: "/etc/pipecd-secret/piped-key"
  apiAddress: "pipecd-control-plane.pipecd.svc.cluster.local:9083"
  webAddress: "http://pipecd-control-plane.pipecd.svc.cluster.local"
  git:
    sshKeyFile: "/etc/pipecd-secret/ssh-key"
  repositories:
    %{ for repo in repositories ~}
    - repoId: "${repo.repo_id}"
      remote: "${repo.remote}"
      branch: "${repo.branch}"
    %{ endfor ~}
  platforms:
    - name: kubernetes
      type: KUBERNETES
      config:
        kubeconfigPath: "${kubernetes_config.kubeconfig_path}"
    - name: terraform
      type: TERRAFORM
      config:
        vars:
          %{ for key, value in terraform_config.vars ~}
          ${key}: "${value}"
          %{ endfor ~}

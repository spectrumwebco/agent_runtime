apiVersion: kubefirst.io/v1
kind: KubefirstConfig
spec:
  gitProvider: ${git_provider}
  gitAuth:
    type: basic
    username: "${git_username}"
    passwordSecretRef: "gitea-password"
  cloudProvider: ${cloud_provider}
  clusterName: ${cluster_name}
  gitopsTemplateURL: "${gitops_template_url}"
  gitopsTemplateBranch: "${gitops_template_branch}"

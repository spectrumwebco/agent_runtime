apiVersion: pipecd.dev/v1beta1
kind: ControlPlane
spec:
  address: ":9082"
  stateKey: "${state_key}"
  datastore:
    type: "filestore"
    config:
      path: "/data/filestore"
  filestore:
    type: "minio"
    config:
      endpoint: "${minio_endpoint}"
      bucket: "${minio_bucket}"
      accessKeyFile: "/etc/pipecd-secret/minio-access-key"
      secretKeyFile: "/etc/pipecd-secret/minio-secret-key"
      autoCreateBucket: true
  git:
    sshKeyFile: "/etc/pipecd-secret/ssh-key"
  repositories:
    %{ for repo in repositories ~}
    - repoId: "${repo.repo_id}"
      remote: "${repo.remote}"
      branch: "${repo.branch}"
    %{ endfor ~}

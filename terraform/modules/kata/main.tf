# Kata Containers Module for Agent Runtime
resource "kubernetes_manifest" "kata_runtime_class" {
  manifest = yamldecode(file("${path.module}/../../k8s/kata-containers/config.yaml"))
}

resource "kubernetes_manifest" "kata_sandbox" {
  manifest = yamldecode(file("${path.module}/../../k8s/kata-containers/sandbox.yaml"))
  
  depends_on = [
    kubernetes_manifest.kata_runtime_class
  ]
}

resource "kubernetes_manifest" "kata_service" {
  manifest = yamldecode(file("${path.module}/../../k8s/kata-containers/service.yaml"))
  
  depends_on = [
    kubernetes_manifest.kata_sandbox
  ]
}

resource "kubernetes_manifest" "kata_pvc" {
  manifest = yamldecode(file("${path.module}/../../k8s/kata-containers/pvc.yaml"))
}

resource "kubernetes_config_map" "kata_config" {
  metadata {
    name      = "kata-containers-config"
    namespace = var.namespace
  }
  
  data = {
    "configuration.toml" = file("${path.module}/../../config/kata/configuration.toml")
    
    "e2b-surf-config.json" = jsonencode({
      browser = {
        headless       = false
        defaultViewport = null
        args = [
          "--no-sandbox",
          "--disable-setuid-sandbox",
          "--disable-dev-shm-usage",
          "--disable-accelerated-2d-canvas",
          "--disable-gpu",
          "--window-size=1920,1080"
        ]
      }
      server = {
        port = 3000
        host = "0.0.0.0"
      }
    })
    
    "setup-desktop.sh" = <<-EOT
      set -e
      
      apt-get update
      DEBIAN_FRONTEND=noninteractive apt-get install -y \
        ubuntu-desktop-minimal \
        xrdp \
        firefox \
        curl \
        wget \
        git \
        build-essential \
        nodejs \
        npm \
        python3 \
        python3-pip
      
      systemctl enable xrdp
      
      wget -q https://packages.microsoft.com/keys/microsoft.asc -O- | apt-key add -
      add-apt-repository "deb [arch=amd64] https://packages.microsoft.com/repos/vscode stable main"
      apt-get update
      apt-get install -y code
      
      TOOLBOX_VERSION=$(curl -s "https://data.services.jetbrains.com/products/releases?code=TBA&latest=true&type=release" | grep -Po '"version":"\K[^"]+')
      wget -q "https://download.jetbrains.com/toolbox/jetbrains-toolbox-$TOOLBOX_VERSION.tar.gz"
      tar -xzf "jetbrains-toolbox-$TOOLBOX_VERSION.tar.gz"
      mv jetbrains-toolbox-*/jetbrains-toolbox /usr/local/bin/
      rm -rf jetbrains-toolbox-*
      
      git clone https://github.com/e2b-dev/surf.git /opt/e2b-surf
      
      cd /opt/e2b-surf
      npm install
      
      cat > /etc/systemd/system/e2b-surf.service <<EOF
      [Unit]
      Description=E2B Surf Browser Agent
      After=network.target
      
      [Service]
      Type=simple
      User=root
      WorkingDirectory=/opt/e2b-surf
      ExecStart=/usr/bin/node /opt/e2b-surf/src/index.js
      Restart=always
      
      [Install]
      WantedBy=multi-user.target
      EOF
      
      systemctl enable e2b-surf
    EOT
  }
}

resource "kubernetes_secret" "kata_certs" {
  metadata {
    name      = "kata-containers-certs"
    namespace = var.namespace
  }
  
  data = {
    "ca.crt"     = var.kata_ca_cert
    "server.crt" = var.kata_server_cert
    "server.key" = var.kata_server_key
  }
  
  type = "Opaque"
}

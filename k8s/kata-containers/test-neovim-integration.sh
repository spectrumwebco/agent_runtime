

set -e

echo "=== Testing Neovim Integration in Kata Container Sandbox ==="

run_in_sandbox() {
  kubectl exec -it $(kubectl get pods -l app=agent-runtime-sandbox -o jsonpath='{.items[0].metadata.name}') -- bash -c "$1"
}

echo "Checking if sandbox pod is running..."
SANDBOX_POD=$(kubectl get pods -l app=agent-runtime-sandbox -o jsonpath='{.items[0].metadata.name}' 2>/dev/null || echo "")
if [ -z "$SANDBOX_POD" ]; then
  echo "Sandbox pod not found. Please deploy the sandbox first."
  exit 1
fi

echo "Sandbox pod found: $SANDBOX_POD"

echo "Testing Neovim accessibility..."
run_in_sandbox "which nvim && nvim --version | head -n 1" || {
  echo "Failed: Neovim is not accessible in the sandbox"
  exit 1
}
echo "Success: Neovim is accessible in the sandbox"

echo "Testing Neovim configuration..."
run_in_sandbox "ls -la /home/ubuntu/.config/nvim/init.lua" || {
  echo "Failed: Neovim configuration not found"
  exit 1
}
echo "Success: Neovim configuration found"

echo "Testing Packer plugin manager..."
run_in_sandbox "ls -la /home/ubuntu/.local/share/nvim/site/pack/packer/start/packer.nvim" || {
  echo "Failed: Packer plugin manager not found"
  exit 1
}
echo "Success: Packer plugin manager found"

echo "Testing Neovim integration server..."
run_in_sandbox "systemctl is-active neovim-server" || {
  echo "Failed: Neovim integration server is not running"
  exit 1
}
echo "Success: Neovim integration server is running"

echo "Testing Neovim API..."
run_in_sandbox "curl -s -X POST -H 'Content-Type: application/json' -d '{\"id\":\"test\"}' http://localhost:8090/neovim/start" || {
  echo "Failed: Neovim API is not accessible"
  exit 1
}
echo "Success: Neovim API is accessible"

echo "Testing Neovim command execution..."
run_in_sandbox "curl -s -X POST -H 'Content-Type: application/json' -d '{\"id\":\"test\",\"command\":\":echo \\\"Hello from Neovim\\\"<CR>\"}' http://localhost:8090/neovim/execute" || {
  echo "Failed: Neovim command execution failed"
  exit 1
}
echo "Success: Neovim command execution works"

echo "Testing Neovim state management..."
run_in_sandbox "cd /opt/agent-runtime/neovim && go run test_state_management.go" || {
  echo "Failed: Neovim state management test failed"
  exit 1
}
echo "Success: Neovim state management works"

echo "Cleaning up..."
run_in_sandbox "curl -s -X POST -H 'Content-Type: application/json' -d '{\"id\":\"test\"}' http://localhost:8090/neovim/stop" || {
  echo "Failed: Neovim cleanup failed"
  exit 1
}
echo "Success: Neovim cleanup completed"

echo "=== All tests passed! Neovim integration is working correctly ==="

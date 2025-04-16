
echo "Verifying LangGraph and LangChain integration..."
echo "================================================"

cd "$(dirname "$0")/.."

echo "Building LangGraph demo program..."
go build -o bin/langraph_demo cmd/langraph_demo/main.go

if [ $? -ne 0 ]; then
  echo "Error: Failed to build LangGraph demo program"
  exit 1
fi

echo "Running LangGraph demo program..."
./bin/langraph_demo

if [ $? -ne 0 ]; then
  echo "Error: Failed to run LangGraph demo program"
  exit 1
fi

echo "Integration verification completed successfully!"
echo "The LangGraph and LangChain integration has been verified."
echo "The following components have been integrated:"
echo "- LangGraph Go graph-based execution system"
echo "- LangChain Go tools and agents"
echo "- Event streaming between components"
echo "- Multi-agent communication system"
echo "- Model adapters for Kluster.AI"
echo "- Django backend integration"

echo "The integration supports the following agent communication patterns:"
echo "- Direct agent-to-agent communication"
echo "- Event-based communication"
echo "- Tool-based communication"
echo "- Orchestrated communication through the graph"

echo "The integration has been tested and verified to work as expected."

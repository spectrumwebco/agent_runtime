
set -e

API_KEY=""
API_URL=""
PROJECT_NAME="agent-runtime-test"
SELF_HOSTED=false

while [[ $# -gt 0 ]]; do
  key="$1"
  case $key in
    --api-key)
      API_KEY="$2"
      shift
      shift
      ;;
    --api-url)
      API_URL="$2"
      shift
      shift
      ;;
    --project-name)
      PROJECT_NAME="$2"
      shift
      shift
      ;;
    --self-hosted)
      SELF_HOSTED=true
      shift
      ;;
    *)
      echo "Unknown option: $1"
      exit 1
      ;;
  esac
done

if [ -z "$API_KEY" ]; then
  echo "Error: LangSmith API key is required."
  echo "Usage: $0 --api-key YOUR_API_KEY [options]"
  echo "Options:"
  echo "  --api-url URL               LangSmith API URL (default: https://api.smith.langchain.com)"
  echo "  --project-name NAME         LangSmith project name (default: agent-runtime-test)"
  echo "  --self-hosted               Use self-hosted LangSmith instance"
  exit 1
fi

if [ -z "$API_URL" ]; then
  if [ "$SELF_HOSTED" = true ]; then
    API_URL="http://langsmith-api.langsmith.svc.cluster.local:8000"
  else
    API_URL="https://api.smith.langchain.com"
  fi
fi

echo "Testing LangSmith integration with LangGraph Go..."
echo "API URL: $API_URL"
echo "Project Name: $PROJECT_NAME"
echo "Self-hosted: $SELF_HOSTED"

export LANGCHAIN_TRACING_V2=true
export LANGCHAIN_ENDPOINT="$API_URL"
export LANGCHAIN_API_KEY="$API_KEY"
export LANGCHAIN_PROJECT="$PROJECT_NAME"

TEMP_DIR=$(mktemp -d)
trap 'rm -rf "$TEMP_DIR"' EXIT

cat <<EOF > "$TEMP_DIR/test_langsmith.go"
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spectrumwebco/agent_runtime/pkg/langraph"
	"github.com/spectrumwebco/agent_runtime/pkg/langsmith"
)

func main() {
	// Create a LangSmith client
	config := langsmith.DefaultConfig()
	client := langsmith.NewClient(config)

	// Create a LangGraph tracer
	tracer := langsmith.NewLangGraphTracer(client, os.Getenv("LANGCHAIN_PROJECT"))

	// Create a simple graph
	graph := langraph.NewGraph("Test Graph")
	
	// Add nodes to the graph
	nodeA := langraph.NewNode("NodeA", "test", func(ctx context.Context, inputs map[string]interface{}) (map[string]interface{}, error) {
		fmt.Println("Executing NodeA")
		return map[string]interface{}{"output": "Hello from NodeA"}, nil
	})
	
	nodeB := langraph.NewNode("NodeB", "test", func(ctx context.Context, inputs map[string]interface{}) (map[string]interface{}, error) {
		fmt.Println("Executing NodeB")
		inputValue, _ := inputs["input"].(string)
		return map[string]interface{}{"output": inputValue + " and NodeB"}, nil
	})
	
	// Add nodes to the graph
	graph.AddNode(nodeA)
	graph.AddNode(nodeB)
	
	// Add an edge from NodeA to NodeB
	graph.AddEdge(nodeA, nodeB, func(outputs map[string]interface{}) map[string]interface{} {
		return map[string]interface{}{"input": outputs["output"]}
	})
	
	// Create an executor with the tracer
	executor := langraph.NewExecutor(langraph.WithTracer(tracer))
	
	// Execute the graph
	ctx := context.Background()
	inputs := map[string]interface{}{"input": "Hello, LangSmith!"}
	
	fmt.Println("Executing graph...")
	outputs, err := executor.ExecuteGraph(ctx, graph, inputs)
	if err != nil {
		fmt.Printf("Error executing graph: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Println("Graph execution completed successfully!")
	fmt.Printf("Outputs: %v\n", outputs)
	
	// Create a run directly using the client
	fmt.Println("Creating a run directly using the client...")
	run, err := client.CreateRun(ctx, &langsmith.CreateRunRequest{
		Name:        "Test Run",
		RunType:     "test",
		StartTime:   time.Now(),
		Inputs:      map[string]interface{}{"input": "Direct client test"},
		ProjectName: os.Getenv("LANGCHAIN_PROJECT"),
		Status:      "running",
	})
	if err != nil {
		fmt.Printf("Error creating run: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("Run created with ID: %s\n", run.ID)
	
	// Update the run with outputs
	_, err = client.UpdateRun(ctx, run.ID, map[string]interface{}{
		"outputs": map[string]interface{}{"output": "Test completed successfully"},
		"end_time": time.Now(),
		"status": "completed",
	})
	if err != nil {
		fmt.Printf("Error updating run: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Println("Run updated successfully!")
	fmt.Println("LangSmith integration test completed successfully!")
}
EOF

echo "Building and running the test program..."
cd "$TEMP_DIR"
go mod init langsmith-test
go mod edit -replace github.com/spectrumwebco/agent_runtime=../../
go mod tidy
go build -o test_langsmith test_langsmith.go

./test_langsmith

echo "LangSmith integration test completed!"
echo "Check the LangSmith UI at $API_URL to see the test runs in project '$PROJECT_NAME'."

# Python FFI Integration for MCP Host

## Overview

This document explores approaches for integrating Python with our Go-based MCP Host application. The goal is to enable calling Python scripts from Go code, allowing us to leverage existing Python tools while benefiting from Go's performance and concurrency features.

## Requirements

1. **Bidirectional Communication**: Call Python from Go and vice versa
2. **Data Marshaling**: Convert data types between Go and Python
3. **Error Handling**: Properly propagate errors between languages
4. **Resource Management**: Manage Python interpreter lifecycle
5. **Performance**: Minimize overhead for cross-language calls
6. **Compatibility**: Support various Python versions (3.8+)
7. **Concurrency**: Handle concurrent Python calls from Go

## Integration Approaches

### 1. go-python/gopy

[go-python](https://github.com/sbinet/go-python) and [gopy](https://github.com/go-python/gopy) provide direct bindings to the Python C API.

#### Advantages:
- Direct integration with Python C API
- Low overhead for function calls
- Full access to Python functionality

#### Disadvantages:
- Requires CGO
- Tied to specific Python versions
- Complex memory management
- Potential for segmentation faults

#### Example:

```go
package main

import (
    "fmt"
    "github.com/sbinet/go-python"
)

func init() {
    err := python.Initialize()
    if err != nil {
        panic(err)
    }
}

func main() {
    defer python.Finalize()

    // Import a Python module
    pyModule := python.PyImport_ImportString("my_python_module")
    if pyModule == nil {
        panic("Could not import Python module")
    }
    defer pyModule.DecRef()

    // Get a Python function
    pyFunc := pyModule.GetAttrString("my_python_function")
    if pyFunc == nil {
        panic("Could not get Python function")
    }
    defer pyFunc.DecRef()

    // Create Python arguments
    args := python.PyTuple_New(1)
    defer args.DecRef()
    
    python.PyTuple_SetItem(args, 0, python.PyString_FromString("hello from Go"))

    // Call the Python function
    result := pyFunc.Call(args, python.PyDict_New())
    defer result.DecRef()

    // Convert result to Go type
    goResult := python.PyString_AsString(result)
    fmt.Println(goResult)
}
```

### 2. py4go

[py4go](https://github.com/DataDog/py4go) is a Go package that provides a higher-level API for Python integration.

#### Advantages:
- Simpler API than go-python
- Better memory management
- More idiomatic Go code

#### Disadvantages:
- Still requires CGO
- Less mature than go-python
- Limited documentation

#### Example:

```go
package main

import (
    "fmt"
    "github.com/DataDog/py4go"
)

func main() {
    // Initialize Python interpreter
    py, err := py4go.NewInterpreter()
    if err != nil {
        panic(err)
    }
    defer py.Close()

    // Import a Python module
    mod, err := py.Import("my_python_module")
    if err != nil {
        panic(err)
    }
    defer mod.DecRef()

    // Call a Python function
    result, err := mod.CallFunction("my_python_function", "hello from Go")
    if err != nil {
        panic(err)
    }
    defer result.DecRef()

    // Convert result to Go type
    var goResult string
    err = result.Parse(&goResult)
    if err != nil {
        panic(err)
    }

    fmt.Println(goResult)
}
```

### 3. gRPC/HTTP API

Using gRPC or HTTP APIs to communicate between Go and Python services.

#### Advantages:
- Language-agnostic
- No direct dependencies between Go and Python
- Can run in separate processes or containers
- Easier to scale and deploy

#### Disadvantages:
- Higher latency
- More complex setup
- Requires serialization/deserialization

#### Example (Go client):

```go
package main

import (
    "context"
    "fmt"
    "log"

    "google.golang.org/grpc"
    pb "github.com/spectrumwebco/agent_runtime/proto"
)

func main() {
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer conn.Close()

    client := pb.NewPythonServiceClient(conn)

    request := &pb.RunScriptRequest{
        Script: "my_python_module",
        Function: "my_python_function",
        Args: []string{"hello from Go"},
    }

    response, err := client.RunScript(context.Background(), request)
    if err != nil {
        log.Fatalf("Failed to run script: %v", err)
    }

    fmt.Println(response.Result)
}
```

#### Example (Python server):

```python
import grpc
from concurrent import futures
import time
import importlib

import python_service_pb2
import python_service_pb2_grpc

class PythonServiceServicer(python_service_pb2_grpc.PythonServiceServicer):
    def RunScript(self, request, context):
        try:
            module = importlib.import_module(request.script)
            function = getattr(module, request.function)
            result = function(*request.args)
            return python_service_pb2.RunScriptResponse(result=str(result))
        except Exception as e:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return python_service_pb2.RunScriptResponse()

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    python_service_pb2_grpc.add_PythonServiceServicer_to_server(
        PythonServiceServicer(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    try:
        while True:
            time.sleep(86400)
    except KeyboardInterrupt:
        server.stop(0)

if __name__ == '__main__':
    serve()
```

### 4. Embedded Python Interpreter

Running a Python interpreter embedded within the Go application.

#### Advantages:
- Single process
- Shared memory space
- Lower latency than network-based approaches

#### Disadvantages:
- Requires CGO
- Complex setup
- Potential for memory leaks

#### Example:

```go
package python

import (
    "fmt"
    "os"
    "path/filepath"
    "sync"
    "unsafe"

    "github.com/spectrumwebco/agent_runtime/internal/config"
)

// #cgo pkg-config: python3
// #include <Python.h>
// int PyRun_SimpleStringWithGIL(const char* command) {
//     int ret;
//     PyGILState_STATE gstate = PyGILState_Ensure();
//     ret = PyRun_SimpleString(command);
//     PyGILState_Release(gstate);
//     return ret;
// }
import "C"

var (
    initialized bool
    mutex       sync.Mutex
)

// Initialize initializes the Python interpreter
func Initialize() error {
    mutex.Lock()
    defer mutex.Unlock()

    if initialized {
        return nil
    }

    // Set Python path
    pythonPath := config.Get().PythonPath
    os.Setenv("PYTHONPATH", pythonPath)

    // Initialize Python interpreter
    C.Py_Initialize()
    if C.Py_IsInitialized() == 0 {
        return fmt.Errorf("failed to initialize Python interpreter")
    }

    // Initialize threads
    C.PyEval_InitThreads()
    C.PyEval_SaveThread() // Release GIL

    initialized = true
    return nil
}

// Finalize finalizes the Python interpreter
func Finalize() {
    mutex.Lock()
    defer mutex.Unlock()

    if !initialized {
        return
    }

    C.Py_Finalize()
    initialized = false
}

// RunScript runs a Python script
func RunScript(script string) error {
    if !initialized {
        return fmt.Errorf("Python interpreter not initialized")
    }

    cScript := C.CString(script)
    defer C.free(unsafe.Pointer(cScript))

    ret := C.PyRun_SimpleStringWithGIL(cScript)
    if ret != 0 {
        return fmt.Errorf("failed to run Python script: %d", ret)
    }

    return nil
}

// RunFunction runs a Python function
func RunFunction(module, function string, args ...interface{}) (interface{}, error) {
    // Implementation omitted for brevity
    return nil, nil
}
```

### 5. Python Subprocess

Executing Python scripts as subprocesses from Go.

#### Advantages:
- Simple implementation
- No CGO required
- Process isolation

#### Disadvantages:
- Higher latency
- Limited interaction
- Process overhead

#### Example:

```go
package python

import (
    "bytes"
    "encoding/json"
    "fmt"
    "os/exec"
    "path/filepath"

    "github.com/spectrumwebco/agent_runtime/internal/config"
)

// RunScript runs a Python script as a subprocess
func RunScript(scriptPath string, args ...string) (string, error) {
    cmd := exec.Command("python3", append([]string{scriptPath}, args...)...)
    var stdout, stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr

    err := cmd.Run()
    if err != nil {
        return "", fmt.Errorf("failed to run Python script: %v\n%s", err, stderr.String())
    }

    return stdout.String(), nil
}

// RunFunction runs a Python function as a subprocess
func RunFunction(module, function string, args ...interface{}) (interface{}, error) {
    // Create a temporary script to call the function
    scriptContent := fmt.Sprintf(`
import json
import sys
import %s

result = %s.%s(*json.loads(sys.argv[1]))
print(json.dumps(result))
`, module, module, function)

    tempDir := config.Get().TempDir
    scriptPath := filepath.Join(tempDir, "run_function.py")

    err := os.WriteFile(scriptPath, []byte(scriptContent), 0644)
    if err != nil {
        return nil, fmt.Errorf("failed to create temporary script: %v", err)
    }
    defer os.Remove(scriptPath)

    // Marshal arguments to JSON
    argsJSON, err := json.Marshal(args)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal arguments: %v", err)
    }

    // Run the script
    output, err := RunScript(scriptPath, string(argsJSON))
    if err != nil {
        return nil, err
    }

    // Unmarshal result
    var result interface{}
    err = json.Unmarshal([]byte(output), &result)
    if err != nil {
        return nil, fmt.Errorf("failed to unmarshal result: %v", err)
    }

    return result, nil
}
```

## Recommended Approach for Agent Runtime

For the Agent Runtime, we recommend a hybrid approach:

1. **Primary: gRPC/HTTP API**
   - Use gRPC for high-performance, structured communication
   - Implement a Python service that exposes SWE-Agent functionality
   - Create a Go client that communicates with the Python service

2. **Secondary: Python Subprocess**
   - Use for simpler, one-off script execution
   - Implement a utility package for running Python scripts
   - Use JSON for data exchange

### Implementation Plan

1. **Create Python Service**:
   - Implement a gRPC server in Python that exposes SWE-Agent functionality
   - Define protocol buffers for service methods
   - Implement service methods that call SWE-Agent functions

2. **Create Go Client**:
   - Generate Go client code from protocol buffers
   - Implement a client package that communicates with the Python service
   - Add error handling and retry logic

3. **Create Python Subprocess Utility**:
   - Implement a utility package for running Python scripts
   - Add support for passing arguments and receiving results
   - Implement error handling and logging

4. **Integrate with MCP Host**:
   - Create MCP tools that use the Python integration
   - Add configuration options for Python paths and settings
   - Implement resource management for Python processes

## Example Implementation

### Protocol Buffers

```protobuf
syntax = "proto3";

package python;

option go_package = "github.com/spectrumwebco/agent_runtime/proto/python";

service PythonService {
  rpc RunScript(RunScriptRequest) returns (RunScriptResponse);
  rpc RunFunction(RunFunctionRequest) returns (RunFunctionResponse);
  rpc RunTool(RunToolRequest) returns (RunToolResponse);
}

message RunScriptRequest {
  string script = 1;
  repeated string args = 2;
}

message RunScriptResponse {
  string result = 1;
  string error = 2;
}

message RunFunctionRequest {
  string module = 1;
  string function = 2;
  string args_json = 3;
}

message RunFunctionResponse {
  string result_json = 1;
  string error = 2;
}

message RunToolRequest {
  string tool_name = 1;
  string args_json = 2;
}

message RunToolResponse {
  string result_json = 1;
  string error = 2;
}
```

### Go Client

```go
package python

import (
    "context"
    "encoding/json"
    "fmt"

    "google.golang.org/grpc"
    pb "github.com/spectrumwebco/agent_runtime/proto/python"
)

// Client is a client for the Python service
type Client struct {
    conn   *grpc.ClientConn
    client pb.PythonServiceClient
}

// NewClient creates a new Python service client
func NewClient(address string) (*Client, error) {
    conn, err := grpc.Dial(address, grpc.WithInsecure())
    if err != nil {
        return nil, fmt.Errorf("failed to connect to Python service: %v", err)
    }

    client := pb.NewPythonServiceClient(conn)
    return &Client{
        conn:   conn,
        client: client,
    }, nil
}

// Close closes the client connection
func (c *Client) Close() error {
    return c.conn.Close()
}

// RunScript runs a Python script
func (c *Client) RunScript(ctx context.Context, script string, args ...string) (string, error) {
    request := &pb.RunScriptRequest{
        Script: script,
        Args:   args,
    }

    response, err := c.client.RunScript(ctx, request)
    if err != nil {
        return "", fmt.Errorf("failed to run script: %v", err)
    }

    if response.Error != "" {
        return "", fmt.Errorf("script error: %s", response.Error)
    }

    return response.Result, nil
}

// RunFunction runs a Python function
func (c *Client) RunFunction(ctx context.Context, module, function string, args ...interface{}) (interface{}, error) {
    argsJSON, err := json.Marshal(args)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal arguments: %v", err)
    }

    request := &pb.RunFunctionRequest{
        Module:   module,
        Function: function,
        ArgsJson: string(argsJSON),
    }

    response, err := c.client.RunFunction(ctx, request)
    if err != nil {
        return nil, fmt.Errorf("failed to run function: %v", err)
    }

    if response.Error != "" {
        return nil, fmt.Errorf("function error: %s", response.Error)
    }

    var result interface{}
    err = json.Unmarshal([]byte(response.ResultJson), &result)
    if err != nil {
        return nil, fmt.Errorf("failed to unmarshal result: %v", err)
    }

    return result, nil
}

// RunTool runs a Python tool
func (c *Client) RunTool(ctx context.Context, toolName string, args map[string]interface{}) (interface{}, error) {
    argsJSON, err := json.Marshal(args)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal arguments: %v", err)
    }

    request := &pb.RunToolRequest{
        ToolName: toolName,
        ArgsJson: string(argsJSON),
    }

    response, err := c.client.RunTool(ctx, request)
    if err != nil {
        return nil, fmt.Errorf("failed to run tool: %v", err)
    }

    if response.Error != "" {
        return nil, fmt.Errorf("tool error: %s", response.Error)
    }

    var result interface{}
    err = json.Unmarshal([]byte(response.ResultJson), &result)
    if err != nil {
        return nil, fmt.Errorf("failed to unmarshal result: %v", err)
    }

    return result, nil
}
```

### Python Service

```python
import grpc
import json
import importlib
from concurrent import futures
import time
import traceback

import python_service_pb2
import python_service_pb2_grpc

from sweagent.tools import tools

class PythonServiceServicer(python_service_pb2_grpc.PythonServiceServicer):
    def RunScript(self, request, context):
        try:
            # Execute the script with the provided arguments
            # This is a simplified implementation
            module = importlib.import_module(request.script)
            result = module.main(*request.args)
            return python_service_pb2.RunScriptResponse(result=str(result))
        except Exception as e:
            traceback.print_exc()
            return python_service_pb2.RunScriptResponse(error=str(e))

    def RunFunction(self, request, context):
        try:
            # Import the module and call the function
            module = importlib.import_module(request.module)
            function = getattr(module, request.function)
            
            # Parse arguments
            args = json.loads(request.args_json)
            
            # Call the function
            result = function(*args)
            
            # Serialize result
            result_json = json.dumps(result)
            
            return python_service_pb2.RunFunctionResponse(result_json=result_json)
        except Exception as e:
            traceback.print_exc()
            return python_service_pb2.RunFunctionResponse(error=str(e))

    def RunTool(self, request, context):
        try:
            # Parse arguments
            args = json.loads(request.args_json)
            
            # Get the tool
            tool = tools.get_tool(request.tool_name)
            if tool is None:
                return python_service_pb2.RunToolResponse(
                    error=f"Tool not found: {request.tool_name}"
                )
            
            # Run the tool
            result = tool.run(**args)
            
            # Serialize result
            result_json = json.dumps(result)
            
            return python_service_pb2.RunToolResponse(result_json=result_json)
        except Exception as e:
            traceback.print_exc()
            return python_service_pb2.RunToolResponse(error=str(e))

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    python_service_pb2_grpc.add_PythonServiceServicer_to_server(
        PythonServiceServicer(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    print("Python service started on port 50051")
    try:
        while True:
            time.sleep(86400)
    except KeyboardInterrupt:
        server.stop(0)

if __name__ == '__main__':
    serve()
```

## Conclusion

For the Agent Runtime, we recommend using a combination of gRPC/HTTP API and Python subprocess approaches. This hybrid approach provides the flexibility to handle both complex interactions and simple script execution while maintaining good performance and reliability.

The gRPC/HTTP API approach is particularly well-suited for integrating with the SWE-Agent and SWE-ReX frameworks, as it allows for structured communication and can be easily scaled. The Python subprocess approach provides a simpler alternative for one-off script execution and can be used as a fallback mechanism.

By implementing both approaches, we can leverage the strengths of both Go and Python while minimizing their weaknesses, resulting in a robust and flexible Agent Runtime.

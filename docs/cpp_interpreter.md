# C++ Interpreter for Agent Runtime

## Overview

The C++ Interpreter is a key component of the Agent Runtime framework that enables the execution of C++ code directly from Go. This capability allows the agent to implement performance-critical components in C++ and solve complex engineering tasks that benefit from C++'s performance characteristics.

## Architecture

The C++ Interpreter is implemented in the `internal/cpp` package and provides a bridge between Go and C++. It consists of the following components:

### Interpreter

The core of the C++ integration is the `Interpreter` struct, which manages the compilation and execution of C++ code. It provides methods for:

- Compiling and executing C++ code
- Executing C++ code with input
- Compiling C++ code into shared libraries
- Managing include paths and compiler flags

### Configuration

The C++ Interpreter can be configured using the `InterpreterConfig` struct, which allows customization of:

- Compiler flags
- Include directories
- Library paths
- Optimization levels

### Integration with MCP

The C++ Interpreter is integrated with the Model Context Protocol (MCP) through the `cpp_tool.go` implementation, which allows the agent to execute C++ code as part of its toolchain.

## Usage

### Basic Usage

```go
import (
    "context"
    "github.com/spectrumwebco/agent_runtime/internal/cpp"
)

func main() {
    // Initialize C++ interpreter with configuration
    config := cpp.InterpreterConfig{
        Flags: []string{"-std=c++17", "-O2"},
    }
    interpreter, err := cpp.NewInterpreter(config)
    if err != nil {
        log.Fatal(err)
    }
    defer interpreter.Close()
    
    // Execute C++ code
    ctx := context.Background()
    code := `
    #include <iostream>
    
    int main() {
        std::cout << "Hello from C++!" << std::endl;
        return 0;
    }
    `
    
    output, err := interpreter.Exec(ctx, code)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(output)
}
```

### Executing with Input

```go
// Execute C++ code with input
code := `
#include <iostream>
#include <string>

int main() {
    std::string name;
    std::getline(std::cin, name);
    std::cout << "Hello, " << name << "!" << std::endl;
    return 0;
}
`

output, err := interpreter.ExecWithInput(ctx, code, "Sam Sepiol")
if err != nil {
    log.Fatal(err)
}

fmt.Println(output) // Output: Hello, Sam Sepiol!
```

### Compiling Libraries

```go
// Compile C++ code into a shared library
code := `
#include <iostream>

extern "C" {
    int add(int a, int b) {
        return a + b;
    }
}
`

libPath, err := interpreter.CompileLibrary(ctx, code, "mathlib")
if err != nil {
    log.Fatal(err)
}

fmt.Println("Library compiled to:", libPath)
```

### Using Custom Include Directories

```go
// Initialize C++ interpreter with custom include directories
config := cpp.InterpreterConfig{
    Flags:        []string{"-std=c++17", "-O2"},
    IncludeDirs:  []string{"/usr/local/include", "/home/user/include"},
}
interpreter, err := cpp.NewInterpreter(config)
if err != nil {
    log.Fatal(err)
}

// Execute C++ code with custom headers
code := `
#include <iostream>
#include "custom_header.h"

int main() {
    std::cout << getCustomValue() << std::endl;
    return 0;
}
`

output, err := interpreter.Exec(ctx, code)
```

## Integration with Tools

The C++ Interpreter is integrated with the Agent Runtime tool system through the `CppTool` implementation in `pkg/tools/cpp_tool.go`. This allows the agent to execute C++ code as part of its toolchain.

### Tool Registration

```go
import "github.com/spectrumwebco/agent_runtime/pkg/tools"

func main() {
    // Create and register C++ tool
    cppTool, err := tools.NewCppTool("cpp", "Executes C++ code")
    if err != nil {
        log.Fatal(err)
    }
    
    // Register tool with the registry
    registry := tools.NewRegistry()
    registry.Register(cppTool)
}
```

### Tool Usage

```go
// Execute C++ code through the tool system
params := map[string]interface{}{
    "code": `
#include <iostream>

int main() {
    std::cout << "Hello from C++ tool!" << std::endl;
    return 0;
}
`,
}

result, err := cppTool.Execute(ctx, params)
if err != nil {
    log.Fatal(err)
}

fmt.Println(result)
```

## Benefits of C++ Integration

The C++ Interpreter provides several benefits to the Agent Runtime framework:

1. **Performance**: C++ offers superior performance for computationally intensive tasks compared to Python or Go.

2. **Memory Efficiency**: C++ provides fine-grained control over memory management, allowing for more efficient resource utilization.

3. **Ecosystem Access**: The C++ ecosystem includes numerous high-performance libraries for scientific computing, graphics, and system-level programming.

4. **Interoperability**: C++ can easily interface with C libraries, providing access to a vast ecosystem of system-level libraries.

5. **Problem-Solving Flexibility**: Having both Python and C++ interpreters allows the agent to choose the most appropriate language for each task.

## Use Cases

The C++ Interpreter is particularly useful for:

1. **High-Performance Computing**: Tasks that require intensive computation, such as numerical simulations or data processing.

2. **System-Level Programming**: Tasks that require low-level access to system resources or hardware.

3. **Real-Time Applications**: Applications with strict timing requirements that benefit from C++'s deterministic performance.

4. **Memory-Constrained Environments**: Scenarios where memory usage must be carefully controlled.

5. **Legacy Code Integration**: Incorporating existing C++ codebases into the agent's capabilities.

## Comparison with Python Interpreter

| Feature | C++ Interpreter | Python Interpreter |
|---------|----------------|-------------------|
| Performance | High | Moderate |
| Memory Efficiency | High | Moderate |
| Ease of Use | Moderate | High |
| Ecosystem | System-level, Scientific | General-purpose, Data Science |
| Compilation | Required | Not required |
| Debugging | More complex | Simpler |
| Use Cases | Performance-critical, System-level | Rapid prototyping, Data analysis |

## Conclusion

The C++ Interpreter is a powerful addition to the Agent Runtime framework, enabling the agent to execute C++ code for performance-critical tasks. By combining the strengths of both Python and C++, the agent can tackle a wider range of software engineering challenges with the most appropriate tool for each task.

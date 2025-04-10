package tools

import (
	"fmt"
	"sync"

	"github.com/spectrumwebco/agent_runtime/pkg/config"
)

type Registry struct {
	config *config.Config
	tools  map[string]Tool
	mutex  sync.RWMutex
}

func NewRegistry(cfg *config.Config) (*Registry, error) {
	registry := &Registry{
		config: cfg,
		tools:  make(map[string]Tool),
	}
	
	registry.registerBuiltinTools()
	
	return registry, nil
}

func (r *Registry) registerBuiltinTools() {
	r.Register(&ShellTool{
		name:        "shell",
		description: "Executes a shell command",
		sandbox:     r.config.Runtime.Sandbox,
	})
	
	r.Register(&FileTool{
		name:        "file",
		description: "Performs file operations",
		sandbox:     r.config.Runtime.Sandbox,
		allowedPaths: r.config.Runtime.AllowedPaths,
	})
	
	r.Register(&HTTPTool{
		name:        "http",
		description: "Makes HTTP requests",
	})
	
	pythonTool, err := NewPythonTool("python", "Executes Python code")
	if err == nil {
		r.Register(pythonTool)
	}
	
	cppTool, err := NewCppTool("cpp", "Executes C++ code")
	if err == nil {
		r.Register(cppTool)
	}
}

func (r *Registry) Register(tool Tool) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	if _, exists := r.tools[tool.Name()]; exists {
		return fmt.Errorf("tool already registered: %s", tool.Name())
	}
	
	r.tools[tool.Name()] = tool
	
	return nil
}

func (r *Registry) Get(name string) (Tool, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	tool, exists := r.tools[name]
	if !exists {
		return nil, fmt.Errorf("tool not found: %s", name)
	}
	
	return tool, nil
}

func (r *Registry) List() []Tool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	tools := make([]Tool, 0, len(r.tools))
	for _, tool := range r.tools {
		tools = append(tools, tool)
	}
	
	return tools
}

func (r *Registry) ListNames() []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	names := make([]string, 0, len(r.tools))
	for name := range r.tools {
		names = append(names, name)
	}
	
	return names
}

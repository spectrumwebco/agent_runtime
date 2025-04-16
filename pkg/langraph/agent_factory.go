package langraph

import (
	"context"
	"fmt"
	"sync"

	"github.com/spectrumwebco/agent_runtime/pkg/eventstream/models"
)

type AgentFactory struct {
	agents      map[string]*Agent
	agentsMutex sync.RWMutex
	eventStream EventStream
}

func NewAgentFactory(eventStream EventStream) *AgentFactory {
	return &AgentFactory{
		agents:      make(map[string]*Agent),
		eventStream: eventStream,
	}
}

func (f *AgentFactory) CreateFrontendAgent(name, description string) (*Agent, error) {
	config := AgentConfig{
		Name:        name,
		Description: description,
		Role:        AgentRoleFrontend,
		Capabilities: []AgentCapability{
			AgentCapabilityUIDesign,
			AgentCapabilityCodeGeneration,
		},
		Metadata: map[string]interface{}{
			"specialization": "frontend",
		},
	}

	handler := func(ctx context.Context, agent *Agent, input map[string]interface{}) (map[string]interface{}, error) {
		return map[string]interface{}{
			"status":  "processed",
			"message": "Frontend agent processed input",
			"input":   input,
		}, nil
	}

	agent := NewAgent(config, handler, f.eventStream)

	agent.AddTool("generate_ui_component", "Generate a UI component", func(ctx context.Context, agent *Agent, params map[string]interface{}) (map[string]interface{}, error) {
		return map[string]interface{}{
			"status":  "generated",
			"message": "UI component generated",
			"params":  params,
		}, nil
	}, map[string]interface{}{
		"component_type": "string",
		"props":          "object",
		"style":          "object",
	})

	agent.AddTool("analyze_ui_design", "Analyze UI design", func(ctx context.Context, agent *Agent, params map[string]interface{}) (map[string]interface{}, error) {
		return map[string]interface{}{
			"status":  "analyzed",
			"message": "UI design analyzed",
			"params":  params,
		}, nil
	}, map[string]interface{}{
		"design_file": "string",
		"criteria":    "array",
	})

	f.agentsMutex.Lock()
	defer f.agentsMutex.Unlock()
	f.agents[agent.Config.ID] = agent

	return agent, nil
}

func (f *AgentFactory) CreateAppBuilderAgent(name, description string) (*Agent, error) {
	config := AgentConfig{
		Name:        name,
		Description: description,
		Role:        AgentRoleAppBuilder,
		Capabilities: []AgentCapability{
			AgentCapabilityAPIDesign,
			AgentCapabilityDatabaseDesign,
			AgentCapabilityCodeGeneration,
		},
		Metadata: map[string]interface{}{
			"specialization": "app_builder",
		},
	}

	handler := func(ctx context.Context, agent *Agent, input map[string]interface{}) (map[string]interface{}, error) {
		return map[string]interface{}{
			"status":  "processed",
			"message": "App builder agent processed input",
			"input":   input,
		}, nil
	}

	agent := NewAgent(config, handler, f.eventStream)

	agent.AddTool("generate_api_endpoint", "Generate an API endpoint", func(ctx context.Context, agent *Agent, params map[string]interface{}) (map[string]interface{}, error) {
		return map[string]interface{}{
			"status":  "generated",
			"message": "API endpoint generated",
			"params":  params,
		}, nil
	}, map[string]interface{}{
		"endpoint_path": "string",
		"method":        "string",
		"request_body":  "object",
		"response_body": "object",
	})

	agent.AddTool("design_database_schema", "Design a database schema", func(ctx context.Context, agent *Agent, params map[string]interface{}) (map[string]interface{}, error) {
		return map[string]interface{}{
			"status":  "designed",
			"message": "Database schema designed",
			"params":  params,
		}, nil
	}, map[string]interface{}{
		"entities":    "array",
		"relationships": "array",
	})

	f.agentsMutex.Lock()
	defer f.agentsMutex.Unlock()
	f.agents[agent.Config.ID] = agent

	return agent, nil
}

func (f *AgentFactory) CreateCodegenAgent(name, description string) (*Agent, error) {
	config := AgentConfig{
		Name:        name,
		Description: description,
		Role:        AgentRoleCodegen,
		Capabilities: []AgentCapability{
			AgentCapabilityCodeGeneration,
			AgentCapabilityCodeReview,
			AgentCapabilityTesting,
		},
		Metadata: map[string]interface{}{
			"specialization": "codegen",
		},
	}

	handler := func(ctx context.Context, agent *Agent, input map[string]interface{}) (map[string]interface{}, error) {
		return map[string]interface{}{
			"status":  "processed",
			"message": "Codegen agent processed input",
			"input":   input,
		}, nil
	}

	agent := NewAgent(config, handler, f.eventStream)

	agent.AddTool("generate_code", "Generate code", func(ctx context.Context, agent *Agent, params map[string]interface{}) (map[string]interface{}, error) {
		return map[string]interface{}{
			"status":  "generated",
			"message": "Code generated",
			"params":  params,
		}, nil
	}, map[string]interface{}{
		"language":      "string",
		"requirements":  "string",
		"file_path":     "string",
		"code_template": "string",
	})

	agent.AddTool("review_code", "Review code", func(ctx context.Context, agent *Agent, params map[string]interface{}) (map[string]interface{}, error) {
		return map[string]interface{}{
			"status":  "reviewed",
			"message": "Code reviewed",
			"params":  params,
		}, nil
	}, map[string]interface{}{
		"code":       "string",
		"language":   "string",
		"criteria":   "array",
	})

	agent.AddTool("generate_tests", "Generate tests", func(ctx context.Context, agent *Agent, params map[string]interface{}) (map[string]interface{}, error) {
		return map[string]interface{}{
			"status":  "generated",
			"message": "Tests generated",
			"params":  params,
		}, nil
	}, map[string]interface{}{
		"code":       "string",
		"language":   "string",
		"test_framework": "string",
	})

	f.agentsMutex.Lock()
	defer f.agentsMutex.Unlock()
	f.agents[agent.Config.ID] = agent

	return agent, nil
}

func (f *AgentFactory) CreateEngineeringAgent(name, description string) (*Agent, error) {
	config := AgentConfig{
		Name:        name,
		Description: description,
		Role:        AgentRoleEngineering,
		Capabilities: []AgentCapability{
			AgentCapabilityCodeGeneration,
			AgentCapabilityCodeReview,
			AgentCapabilityAPIDesign,
			AgentCapabilityDatabaseDesign,
			AgentCapabilityTesting,
			AgentCapabilityDeployment,
			AgentCapabilityDocumentation,
			AgentCapabilityPlanning,
		},
		Metadata: map[string]interface{}{
			"specialization": "engineering",
		},
	}

	handler := func(ctx context.Context, agent *Agent, input map[string]interface{}) (map[string]interface{}, error) {
		return map[string]interface{}{
			"status":  "processed",
			"message": "Engineering agent processed input",
			"input":   input,
		}, nil
	}

	agent := NewAgent(config, handler, f.eventStream)

	agent.AddTool("analyze_requirements", "Analyze requirements", func(ctx context.Context, agent *Agent, params map[string]interface{}) (map[string]interface{}, error) {
		return map[string]interface{}{
			"status":  "analyzed",
			"message": "Requirements analyzed",
			"params":  params,
		}, nil
	}, map[string]interface{}{
		"requirements": "string",
		"context":      "string",
	})

	agent.AddTool("design_architecture", "Design architecture", func(ctx context.Context, agent *Agent, params map[string]interface{}) (map[string]interface{}, error) {
		return map[string]interface{}{
			"status":  "designed",
			"message": "Architecture designed",
			"params":  params,
		}, nil
	}, map[string]interface{}{
		"requirements": "string",
		"constraints":  "array",
		"technologies": "array",
	})

	agent.AddTool("create_technical_documentation", "Create technical documentation", func(ctx context.Context, agent *Agent, params map[string]interface{}) (map[string]interface{}, error) {
		return map[string]interface{}{
			"status":  "created",
			"message": "Technical documentation created",
			"params":  params,
		}, nil
	}, map[string]interface{}{
		"project":     "string",
		"components":  "array",
		"format":      "string",
	})

	f.agentsMutex.Lock()
	defer f.agentsMutex.Unlock()
	f.agents[agent.Config.ID] = agent

	return agent, nil
}

func (f *AgentFactory) CreateOrchestratorAgent(name, description string) (*Agent, error) {
	config := AgentConfig{
		Name:        name,
		Description: description,
		Role:        AgentRoleOrchestrator,
		Capabilities: []AgentCapability{
			AgentCapabilityOrchestration,
			AgentCapabilityPlanning,
		},
		Metadata: map[string]interface{}{
			"specialization": "orchestrator",
		},
	}

	handler := func(ctx context.Context, agent *Agent, input map[string]interface{}) (map[string]interface{}, error) {
		return map[string]interface{}{
			"status":  "processed",
			"message": "Orchestrator agent processed input",
			"input":   input,
		}, nil
	}

	agent := NewAgent(config, handler, f.eventStream)

	agent.AddTool("create_execution_plan", "Create execution plan", func(ctx context.Context, agent *Agent, params map[string]interface{}) (map[string]interface{}, error) {
		return map[string]interface{}{
			"status":  "created",
			"message": "Execution plan created",
			"params":  params,
		}, nil
	}, map[string]interface{}{
		"task":        "string",
		"constraints": "array",
		"agents":      "array",
	})

	agent.AddTool("assign_task", "Assign task to agent", func(ctx context.Context, agent *Agent, params map[string]interface{}) (map[string]interface{}, error) {
		return map[string]interface{}{
			"status":  "assigned",
			"message": "Task assigned",
			"params":  params,
		}, nil
	}, map[string]interface{}{
		"task":     "string",
		"agent_id": "string",
		"priority": "number",
	})

	agent.AddTool("monitor_progress", "Monitor task progress", func(ctx context.Context, agent *Agent, params map[string]interface{}) (map[string]interface{}, error) {
		return map[string]interface{}{
			"status":  "monitored",
			"message": "Progress monitored",
			"params":  params,
		}, nil
	}, map[string]interface{}{
		"task_id": "string",
	})

	f.agentsMutex.Lock()
	defer f.agentsMutex.Unlock()
	f.agents[agent.Config.ID] = agent

	return agent, nil
}

func (f *AgentFactory) GetAgent(id string) (*Agent, error) {
	f.agentsMutex.RLock()
	defer f.agentsMutex.RUnlock()

	agent, exists := f.agents[id]
	if !exists {
		return nil, fmt.Errorf("agent %s does not exist", id)
	}

	return agent, nil
}

func (f *AgentFactory) ListAgents() []*Agent {
	f.agentsMutex.RLock()
	defer f.agentsMutex.RUnlock()

	var agents []*Agent
	for _, agent := range f.agents {
		agents = append(agents, agent)
	}

	return agents
}

func (f *AgentFactory) ListAgentsByRole(role AgentRole) []*Agent {
	f.agentsMutex.RLock()
	defer f.agentsMutex.RUnlock()

	var agents []*Agent
	for _, agent := range f.agents {
		if agent.Config.Role == role {
			agents = append(agents, agent)
		}
	}

	return agents
}

func (f *AgentFactory) ListAgentsByCapability(capability AgentCapability) []*Agent {
	f.agentsMutex.RLock()
	defer f.agentsMutex.RUnlock()

	var agents []*Agent
	for _, agent := range f.agents {
		if agent.HasCapability(capability) {
			agents = append(agents, agent)
		}
	}

	return agents
}

func (f *AgentFactory) RemoveAgent(id string) error {
	f.agentsMutex.Lock()
	defer f.agentsMutex.Unlock()

	if _, exists := f.agents[id]; !exists {
		return fmt.Errorf("agent %s does not exist", id)
	}

	delete(f.agents, id)
	return nil
}

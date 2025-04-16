package types

import (
	"time"
)

type AgentType string

const (
	AgentTypeFrontend AgentType = "frontend"

	AgentTypeAppBuilder AgentType = "app_builder"

	AgentTypeCodegen AgentType = "codegen"

	AgentTypeEngineering AgentType = "engineering"

	AgentTypeOrchestrator AgentType = "orchestrator"
)

type AgentStatus string

const (
	AgentStatusInitialized AgentStatus = "initialized"

	AgentStatusRunning AgentStatus = "running"

	AgentStatusPaused AgentStatus = "paused"

	AgentStatusStopped AgentStatus = "stopped"

	AgentStatusError AgentStatus = "error"
)

type AgentConfig struct {
	ID string `json:"id" yaml:"id"`

	Name string `json:"name" yaml:"name"`

	Type AgentType `json:"type" yaml:"type"`

	Description string `json:"description" yaml:"description"`

	Capabilities []string `json:"capabilities" yaml:"capabilities"`

	Tools []string `json:"tools" yaml:"tools"`

	Dependencies []string `json:"dependencies" yaml:"dependencies"`

	Config map[string]interface{} `json:"config" yaml:"config"`
}

type AgentState struct {
	ID string `json:"id" yaml:"id"`

	Status AgentStatus `json:"status" yaml:"status"`

	LastUpdated time.Time `json:"last_updated" yaml:"last_updated"`

	CurrentTask string `json:"current_task" yaml:"current_task"`

	Progress int `json:"progress" yaml:"progress"`

	Error string `json:"error" yaml:"error"`

	Data map[string]interface{} `json:"data" yaml:"data"`
}

type AgentMessage struct {
	ID string `json:"id" yaml:"id"`

	SenderID string `json:"sender_id" yaml:"sender_id"`

	ReceiverID string `json:"receiver_id" yaml:"receiver_id"`

	Type string `json:"type" yaml:"type"`

	Content string `json:"content" yaml:"content"`

	Timestamp time.Time `json:"timestamp" yaml:"timestamp"`

	Metadata map[string]interface{} `json:"metadata" yaml:"metadata"`
}

type AgentTask struct {
	ID string `json:"id" yaml:"id"`

	AgentID string `json:"agent_id" yaml:"agent_id"`

	Type string `json:"type" yaml:"type"`

	Description string `json:"description" yaml:"description"`

	Status string `json:"status" yaml:"status"`

	Priority int `json:"priority" yaml:"priority"`

	CreatedAt time.Time `json:"created_at" yaml:"created_at"`

	StartedAt *time.Time `json:"started_at" yaml:"started_at"`

	CompletedAt *time.Time `json:"completed_at" yaml:"completed_at"`

	Parameters map[string]interface{} `json:"parameters" yaml:"parameters"`

	Result map[string]interface{} `json:"result" yaml:"result"`
}

type AgentCapability struct {
	Name string `json:"name" yaml:"name"`

	Description string `json:"description" yaml:"description"`

	Parameters map[string]interface{} `json:"parameters" yaml:"parameters"`
}

type AgentTool struct {
	Name string `json:"name" yaml:"name"`

	Description string `json:"description" yaml:"description"`

	Parameters map[string]interface{} `json:"parameters" yaml:"parameters"`
}

type AgentRegistry struct {
	Agents map[string]*AgentConfig `json:"agents" yaml:"agents"`
}

type AgentEvent struct {
	ID string `json:"id" yaml:"id"`

	AgentID string `json:"agent_id" yaml:"agent_id"`

	Type string `json:"type" yaml:"type"`

	Timestamp time.Time `json:"timestamp" yaml:"timestamp"`

	Data map[string]interface{} `json:"data" yaml:"data"`
}

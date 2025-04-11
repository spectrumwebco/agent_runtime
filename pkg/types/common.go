package types
import "time"


type Message struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"` // Can be string or list[dict]
}

type StepOutput struct {
	Thought      string                 `json:"thought"`
	Action       string                 `json:"action"`
	Output       string                 `json:"output"` // Model's raw output
	Observation  string                 `json:"observation"` // Result of executing action
	ExecutionTime float64                `json:"execution_time"`
	Done         bool                   `json:"done"`
	ExitStatus   interface{}            `json:"exit_status"` // Can be int or string
	Submission   *string                `json:"submission"` // Use pointer for optional field
	State        map[string]string      `json:"state"`
	ToolCalls    []map[string]interface{} `json:"tool_calls,omitempty"`
	ToolCallIDs  []string               `json:"tool_call_ids,omitempty"`
	ExtraInfo    map[string]interface{} `json:"extra_info,omitempty"`
}

type TrajectoryStep struct {
	Action       string                 `json:"action"`
	Observation  string                 `json:"observation"`
	Response     string                 `json:"response"` // Model's raw output for this step
	State        map[string]string      `json:"state"`
	Thought      string                 `json:"thought"`
	ExecutionTime float64                `json:"execution_time"`
	Messages     []Message              `json:"messages"` // History leading up to this step
	ExtraInfo    map[string]interface{} `json:"extra_info,omitempty"`
}

type AgentInfo struct {
	ModelStats      map[string]float64 `json:"model_stats,omitempty"`
	ExitStatus      *string            `json:"exit_status,omitempty"`
	Submission      *string            `json:"submission,omitempty"`
	Review          map[string]interface{} `json:"review,omitempty"`
	EditedFiles30   string             `json:"edited_files30,omitempty"`
	EditedFiles50   string             `json:"edited_files50,omitempty"`
	EditedFiles70   string             `json:"edited_files70,omitempty"`
	Summarizer      map[string]interface{} `json:"summarizer,omitempty"`
	SweAgentHash    string             `json:"swe_agent_hash,omitempty"`
	SweAgentVersion string             `json:"swe_agent_version,omitempty"`
	SweRexVersion   string             `json:"swe_rex_version,omitempty"`
	SweRexHash      string             `json:"swe_rex_hash,omitempty"`
}

type AgentRunResult struct {
	Info       AgentInfo        `json:"info"`
	Trajectory []TrajectoryStep `json:"trajectory"`
}



type ExampleType struct {
	ID   string
	Name string
}

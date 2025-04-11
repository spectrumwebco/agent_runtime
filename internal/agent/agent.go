package agent

import (
	"context"
	"fmt"
	"os" // Added for file system operations
	"path/filepath" // Added for path manipulation
	"time"

	"github.com/spectrumwebco/agent_runtime/internal/config" // Assuming config is internal
	"github.com/spectrumwebco/agent_runtime/internal/env"   // Assuming env is internal
	"github.com/spectrumwebco/agent_runtime/pkg/tools"
)

type Message struct {
	Role        string                 `json:"role"`         // system, assistant, user
	Content     string                 `json:"content"`      // The actual message content
	Agent       string                 `json:"agent"`        // Name of the agent that generated this message
	MessageType string                 `json:"message_type"` // Type of message (e.g., system_prompt, observation)
	Thought     string                 `json:"thought,omitempty"`
	Action      string                 `json:"action,omitempty"`
	RawOutput   string                 `json:"raw_output,omitempty"`
	ExtraInfo   map[string]interface{} `json:"extra_info,omitempty"`
}

type AgentInfo struct {
	Submission      string                 `json:"submission"`
	ExitStatus      string                 `json:"exit_status"`
	ModelStats      map[string]interface{} `json:"model_stats"` // Placeholder for model statistics
	EditedFiles     map[string]string      `json:"edited_files"`  // Placeholder for edited files context
	SweAgentHash    string                 `json:"swe_agent_hash"`    // Go equivalent placeholder
	SweAgentVersion string                 `json:"swe_agent_version"` // Go equivalent placeholder
	SweRexVersion   string                 `json:"swe_rex_version"`   // Go equivalent placeholder
	SweRexHash      string                 `json:"swe_rex_hash"`      // Go equivalent placeholder
}

type StepOutput struct {
	Thought       string                 `json:"thought"`
	Action        string                 `json:"action"`
	Observation   string                 `json:"observation"`
	Output        string                 `json:"output"` // Raw model output
	State         map[string]interface{} `json:"state"`  // Environment state snapshot
	Submission    string                 `json:"submission"`
	ExitStatus    string                 `json:"exit_status"`
	Done          bool                   `json:"done"`
	ExecutionTime time.Duration          `json:"execution_time"`
	ExtraInfo     map[string]interface{} `json:"extra_info"`
	ToolCalls     []interface{}          `json:"tool_calls"`    // Placeholder for tool calls
	ToolCallIDs   []string               `json:"tool_call_ids"` // Placeholder for tool call IDs
}

type TrajectoryStep struct {
	Action        string                 `json:"action"`
	Observation   string                 `json:"observation"`
	Response      string                 `json:"response"`
	Thought       string                 `json:"thought"`
	ExecutionTime time.Duration          `json:"execution_time"`
	State         map[string]interface{} `json:"state"`
	Messages      []map[string]string    `json:"messages"` // History messages at this step
	ExtraInfo     map[string]interface{} `json:"extra_info"`
}

type ProblemStatement struct {
	ID               string
	ProblemStatement string
}

func (ps *ProblemStatement) GetProblemStatement() string {
	if ps == nil {
		return ""
	}
	return ps.ProblemStatement
}

func (ps *ProblemStatement) GetID() string {
	if ps == nil {
		return ""
	}
	return ps.ID
}

type DefaultAgent struct {
	Name                 string
	Config               *config.Config // Placeholder for agent config (needs definition)
	Templates            interface{}    // Placeholder for template config (e.g., TemplateConfig struct)
	Tools                *tools.Handler // Placeholder for tool handler (needs definition)
	HistoryProcessors    []interface{}  // Placeholder for history processors
	Model                interface{}    // Placeholder for the language model interface
	MaxRequeries         int
	History              []Message           // Stores the conversation history as Message objects
	trajectory           []TrajectoryStep    // Stores the execution trajectory
	Info                 AgentInfo
	Env                  *env.SWEEnv       // Execution environment (exported for loop.go)
	problemStatement     *ProblemStatement // Task definition
	outputDir            string            // Directory for output files
	trajPath             string            // Path to save trajectory file
	alwaysRequireZeroExit bool
	nConsecutiveTimeouts int
	totalExecutionTime   time.Duration
}

func NewDefaultAgent(/* TODO: Add parameters based on SWE-Agent __init__ */) (*DefaultAgent, error) {
	fmt.Println("Placeholder: Initializing DefaultAgent...")
	agentInfo := AgentInfo{
		ModelStats:  make(map[string]interface{}),
		EditedFiles: make(map[string]string),
	}
	return &DefaultAgent{
		Name:              "default-go-agent",
		MaxRequeries:      3, // Default from Python
		History:           make([]Message, 0),
		trajectory:        make([]TrajectoryStep, 0),
		Info:              agentInfo, // Use initialized AgentInfo
	}, nil
}

func (a *DefaultAgent) Setup(environment *env.SWEEnv, problemStmt *ProblemStatement, outputDirPath string) error {
	if environment == nil {
		return fmt.Errorf("environment cannot be nil")
	}
	if problemStmt == nil {
		return fmt.Errorf("problem statement cannot be nil")
	}
	if outputDirPath == "" {
		outputDirPath = "." // Default to current directory
	}

	err := os.MkdirAll(outputDirPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create output directory %s: %w", outputDirPath, err)
	}

	a.outputDir = outputDirPath
	a.problemStatement = problemStmt
	a.Env = environment
	instanceID := a.problemStatement.GetID()

	fmt.Printf("Setting up agent for instance %s\n", instanceID)

	a.trajPath = filepath.Join(a.outputDir, instanceID+".traj")
	fmt.Printf("Trajectory will be saved to %s\n", a.trajPath)

	fmt.Println("Placeholder: Installing tools...")
	if a.Tools == nil {
		toolRegistry, _ := tools.NewRegistry(a.Config) // Assuming a.Config is initialized
		toolDefs := []tools.ToolDefinition{} // Load from config
		toolConf := tools.ToolConfig{} // Load from config
		a.Tools, _ = tools.NewHandler(toolConf, toolDefs, toolRegistry)
		fmt.Println("Warning: Tools handler was nil, initialized with placeholder.")
	}

	a.Info = AgentInfo{
		SweAgentHash:    "go-dev-hash",
		SweAgentVersion: "go-dev-version",
		SweRexVersion:   "go-rex-dev-version",
		SweRexHash:      "go-rex-dev-hash",
		ModelStats:      make(map[string]interface{}), // Initialize maps
		EditedFiles:     make(map[string]string),
	}
	a.History = make([]map[string]string, 0) // Reset history
	a.trajectory = make([]TrajectoryStep, 0) // Reset trajectory
	a.nConsecutiveTimeouts = 0
	a.totalExecutionTime = 0

	err = a.Env.SetEnvVariables(map[string]string{
		"PROBLEM_STATEMENT": a.problemStatement.GetProblemStatement(),
	})
	if err != nil {
		return fmt.Errorf("failed to set PROBLEM_STATEMENT env var: %w", err)
	}

	err = a.addSystemMessageToHistory()
	if err != nil {
		return fmt.Errorf("failed to add system message: %w", err)
	}
	err = a.addDemonstrationsToHistory()
	if err != nil {
		return fmt.Errorf("failed to add demonstrations: %w", err)
	}

	initialState, err := a.Tools.GetState(a.Env)
	if err != nil {
		fmt.Printf("Warning: Failed to get initial state: %v\n", err)
		initialState = make(map[string]interface{}) // Use empty state
	}

	err = a.addInstanceTemplateToHistory(initialState)
	if err != nil {
		return fmt.Errorf("failed to add instance template: %w", err)
	}

	fmt.Println("Agent setup complete.")
	return nil
}

func (a *DefaultAgent) QueryModel(ctx context.Context, history []Message) (map[string]interface{}, error) {
	fmt.Println("Placeholder: Querying model...")
	
	return map[string]interface{}{
		"message": fmt.Sprintf("Placeholder model output for step %d with thought and action", len(history)),
	}, nil
}

func (a *DefaultAgent) AddStepToHistory(thought, action string, modelOutput map[string]interface{}, observation string) {
	assistantContent := ""
	if thought != "" {
		assistantContent += thought + "\n"
	}
	assistantContent += fmt.Sprintf("```tool\n%s\n```", action) // Basic tool format

	assistantMsg := Message{
		Role:        "assistant",
		Content:     assistantContent,
		Thought:     thought,
		Action:      action,
		RawOutput:   modelOutput["message"].(string), // Store raw model output
		Agent:       a.Name,
		MessageType: "assistant_response",
	}
	a._appendHistory(assistantMsg)

	observationMsg := Message{
		Role:        "user", // Observations are treated as user input for the next turn
		Content:     fmt.Sprintf("Observation:\n%s", observation),
		Agent:       a.Name,
		MessageType: "observation",
	}
	a._appendHistory(observationMsg)
}


func (a *DefaultAgent) Step() (*StepOutput, error) {
	fmt.Println("Placeholder: Executing agent step...")
	stepOutput := &StepOutput{Done: true} // Placeholder to stop loop for now
	a.addStepToTrajectory(stepOutput)     // Add placeholder step
	return stepOutput, nil
}

func (a *DefaultAgent) Run(environment *env.SWEEnv, problemStmt *ProblemStatement, outputDirPath string) (*AgentRunResult, error) {
	fmt.Println("Placeholder: Starting agent run...")

	err := a.Setup(environment, problemStmt, outputDirPath)
	if err != nil {
		return nil, fmt.Errorf("agent setup failed: %w", err)
	}

	var stepOutput *StepOutput
	for i := 0; i < 1; i++ { // Limit steps for now
		stepOutput, err = a.Step()
		if err != nil {
			return nil, fmt.Errorf("agent step failed: %w", err)
		}
		if stepOutput.Done {
			break
		}
	}

	fmt.Println("Placeholder: Agent run finished.")
	return &AgentRunResult{
		Info:       a.Info,
		Trajectory: a.trajectory,
	}, nil
}

func (a *DefaultAgent) addStepToTrajectory(step *StepOutput) {
	histCopy := make([]map[string]string, len(a.History))
	for i, msg := range a.History {
		msgCopy := map[string]string{
			"role":         msg.Role,
			"content":      msg.Content,
			"agent":        msg.Agent,
			"message_type": msg.MessageType,
		}
		if msg.Thought != "" {
			msgCopy["thought"] = msg.Thought
		}
		if msg.Action != "" {
			msgCopy["action"] = msg.Action
		}
		if msg.RawOutput != "" {
			msgCopy["raw_output"] = msg.RawOutput
		}
		histCopy[i] = msgCopy
	}
	stateCopy := make(map[string]interface{})
	if step.State != nil {
		for k, v := range step.State {
			stateCopy[k] = v // Assuming state values are immutable or copying is sufficient
		}
	}

	trajectoryStep := TrajectoryStep{
		Action:        step.Action,
		Observation:   step.Observation,
		Response:      step.Output,
		Thought:       step.Thought,
		ExecutionTime: step.ExecutionTime,
		State:         stateCopy,
		Messages:      histCopy, // Use the deep copied history
		ExtraInfo:     step.ExtraInfo, // Assuming ExtraInfo is okay with shallow copy
	}
	a.trajectory = append(a.trajectory, trajectoryStep)
}

type AgentRunResult struct {
	Info       AgentInfo        `json:"info"`
	Trajectory []TrajectoryStep `json:"trajectory"`
}

func (a *DefaultAgent) _appendHistory(item Message) {
	msgCopy := Message{
		Role:        item.Role,
		Content:     item.Content,
		Agent:       item.Agent,
		MessageType: item.MessageType,
		Thought:     item.Thought,
		Action:      item.Action,
		RawOutput:   item.RawOutput,
	}
	
	if item.ExtraInfo != nil {
		msgCopy.ExtraInfo = make(map[string]interface{})
		for k, v := range item.ExtraInfo {
			msgCopy.ExtraInfo[k] = v
		}
	}
	
	a.History = append(a.History, msgCopy)
}

func (a *DefaultAgent) addSystemMessageToHistory() error {
	fmt.Println("Adding system message to history...")
	systemTemplate := "You are SWE-Agent, an autonomous software development agent. Follow the instructions carefully and use the provided tools to solve the task." // Placeholder template
	
	systemMsg := systemTemplate // Use raw template for now

	if systemMsg != "" {
		fmt.Printf("SYSTEM (%s)\n%s\n", a.Name, systemMsg)
	a._appendHistory(Message{
		Role:        "system",
		Content:     systemMsg,
		Agent:       a.Name,
		MessageType: "system_prompt",
	})
	}
	return nil
}

func (a *DefaultAgent) addDemonstrationsToHistory() error {
	fmt.Println("Adding demonstrations to history...")

	// 		}
	// 	}
	// }

	return nil
}

func (a *DefaultAgent) addInstanceTemplateToHistory(initialState map[string]interface{}) error {
	fmt.Println("Adding instance template to history...")
	instanceTemplate := "Solve the following task:\n{{problem_statement}}\n\nAvailable tools:\n{{command_docs}}\n\nInitial State:\n{{initial_state}}" // Placeholder template

	instanceMsg := fmt.Sprintf("Solve the following task:\n%s\n\nInitial State: %+v", a.problemStatement.GetProblemStatement(), initialState) // Basic placeholder rendering

	if instanceMsg != "" {
		fmt.Printf("USER (%s)\n%s\n", a.Name, instanceMsg)
	a._appendHistory(Message{
		Role:        "user",
		Content:     instanceMsg,
		Agent:       a.Name,
		MessageType: "instance_prompt",
	})
	}
	return nil
}

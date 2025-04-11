package agent

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spectrumwebco/agent_runtime/internal/config"
	"github.com/spectrumwebco/agent_runtime/internal/env"
	"github.com/spectrumwebco/agent_runtime/internal/ffi/python"
	"github.com/spectrumwebco/agent_runtime/internal/parser"
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
	Tools                *tools.ToolRegistry // Use ToolRegistry directly for now
	HistoryProcessors    []interface{}  // Placeholder for history processors
	Model                interface{}    // Placeholder for the language model interface
	Parser               parser.Parser  // Parser for model output
	ModelName            string         // Name of the model to use
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

type AgentOption func(*DefaultAgent) error

func WithTools(toolRegistry *tools.ToolRegistry) AgentOption {
	return func(a *DefaultAgent) error {
		a.Tools = toolRegistry
		return nil
	}
}

func WithParser(p parser.Parser) AgentOption {
	return func(a *DefaultAgent) error {
		a.Parser = p
		return nil
	}
}

func WithModelName(modelName string) AgentOption {
	return func(a *DefaultAgent) error {
		a.ModelName = modelName
		return nil
	}
}

func WithMaxRequeries(maxRequeries int) AgentOption {
	return func(a *DefaultAgent) error {
		a.MaxRequeries = maxRequeries
		return nil
	}
}

func WithConfig(cfg *config.Config) AgentOption {
	return func(a *DefaultAgent) error {
		a.Config = cfg
		return nil
	}
}

func NewDefaultAgent(options ...AgentOption) (*DefaultAgent, error) {
	fmt.Println("Initializing DefaultAgent...")
	agentInfo := AgentInfo{
		ModelStats:  make(map[string]interface{}),
		EditedFiles: make(map[string]string),
	}
	
	agent := &DefaultAgent{
		Name:              "default-go-agent",
		MaxRequeries:      3, // Default from Python
		History:           make([]Message, 0),
		trajectory:        make([]TrajectoryStep, 0),
		Info:              agentInfo,
		ModelName:         "gpt-4", // Default model
		Parser:            parser.NewThoughtActionParser(), // Default parser
	}
	
	for _, option := range options {
		if err := option(agent); err != nil {
			return nil, fmt.Errorf("failed to apply agent option: %w", err)
		}
	}
	
	return agent, nil
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
		toolConf := &tools.ToolConfig{
			ToolsDir: "tools", // Default path, should be configurable
		}
		var err error
		projectRoot, _ := os.Getwd() // Or determine project root differently
		toolsFullPath := filepath.Join(projectRoot, toolConf.ToolsDir)
		scriptRunner, runnerErr := python.NewScriptRunner(toolsFullPath)
		if runnerErr != nil {
			fmt.Printf("Warning: Failed to initialize Python script runner: %v. Python tools will not be available.\n", runnerErr)
			scriptRunner = nil // Explicitly set to nil if initialization fails
		}

		a.Tools, err = tools.NewToolRegistry(toolConf, scriptRunner) // Pass scriptRunner
		if err != nil {
			return fmt.Errorf("failed to initialize tool registry: %w", err)
		}
		// if err != nil {
		// }
		fmt.Println("Warning: Tools handler was nil, initialized with placeholder ToolRegistry.")
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
	a.History = make([]Message, 0) // Reset history
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


func (a *DefaultAgent) Step(ctx context.Context) (*StepOutput, error) {
	startTime := time.Now()
	fmt.Printf("\n--- Step %d ---\n", len(a.trajectory)+1)

	modelOutput, err := a.QueryModel(ctx, a.History)
	if err != nil {
		return nil, fmt.Errorf("failed to query model: %w", err)
	}
	fmt.Printf("Model Output: %+v\n", modelOutput)

	var parser parser.Parser
	if a.Parser != nil {
		parser = a.Parser
	} else {
		parser = parser.NewThoughtActionParser()
		fmt.Println("Warning: Using default ThoughtActionParser")
	}

	thought, action, err := parser.Parse(modelOutput, a.Tools.ListCommands())
	if err != nil {
		fmt.Printf("Warning: Failed to parse model output: %v. Attempting recovery...\n", err)
		thought = "Parsing failed, attempting raw output as action."
		action = modelOutput["message"].(string)
	}
	fmt.Printf("Thought: %s\nAction: %s\n", thought, action)

	observation, execErr := a.Tools.ExecuteAction(action, a.Env)
	executionTime := time.Since(startTime)
	a.totalExecutionTime += executionTime

	if execErr != nil {
		fmt.Printf("Execution Error: %v\n", execErr)
		observation = fmt.Sprintf("Error executing action: %v", execErr)
		a.nConsecutiveTimeouts++
		if a.nConsecutiveTimeouts >= a.Tools.GetConfig().MaxConsecutiveTimeouts {
			return nil, fmt.Errorf("too many consecutive timeouts (%d)", a.nConsecutiveTimeouts)
		}
	} else {
		a.nConsecutiveTimeouts = 0
		fmt.Printf("Observation: %s\n", observation)
	}

	done := false
	exitStatus := "incomplete"
	submission := ""
	if strings.HasPrefix(action, "submit") {
		done = true
		exitStatus = "submitted"
		submission = strings.TrimSpace(strings.TrimPrefix(action, "submit"))
		fmt.Println("Submit action detected.")
	}

	currentState, stateErr := a.Tools.GetState(a.Env)
	if stateErr != nil {
		fmt.Printf("Warning: Failed to get environment state: %v\n", stateErr)
		currentState = make(map[string]interface{}) // Use empty state on error
	}

	stepOutput := &StepOutput{
		Thought:       thought,
		Action:        action,
		Observation:   observation,
		Output:        modelOutput["message"].(string), // Raw model output
		State:         currentState,
		Submission:    submission,
		ExitStatus:    exitStatus,
		Done:          done,
		ExecutionTime: executionTime,
		ExtraInfo:     make(map[string]interface{}), // Initialize ExtraInfo
	}

	a.AddStepToHistory(thought, action, modelOutput, observation)
	a.addStepToTrajectory(stepOutput)

	if a.trajPath != "" {
		if err := a.saveTrajectory(); err != nil {
			fmt.Printf("Warning: Failed to save trajectory: %v\n", err)
		}
	}

	fmt.Printf("Step execution time: %s\n", executionTime)
	fmt.Printf("Total execution time: %s\n", a.totalExecutionTime)

	return stepOutput, nil
}

func (a *DefaultAgent) Run(environment *env.SWEEnv, problemStmt *ProblemStatement, outputDirPath string) (*AgentRunResult, error) {
	fmt.Println("Starting agent run...")

	err := a.Setup(environment, problemStmt, outputDirPath)
	if err != nil {
		return nil, fmt.Errorf("agent setup failed: %w", err)
	}

	totalTimeout := time.Duration(a.Tools.GetConfig().TotalExecutionTimeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), totalTimeout)
	defer cancel()

	var stepOutput *StepOutput
	maxSteps := 50 // Reasonable default, could be configurable
	for i := 0; i < maxSteps; i++ {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("agent run timed out after %s", totalTimeout)
		default:
			stepOutput, err = a.Step(ctx)
			if err != nil {
				return nil, fmt.Errorf("agent step failed: %w", err)
			}
			
			if stepOutput.Done {
				a.Info.ExitStatus = stepOutput.ExitStatus
				a.Info.Submission = stepOutput.Submission
				break
			}
		}
	}

	if err := a.saveTrajectory(); err != nil {
		fmt.Printf("Warning: Failed to save final trajectory: %v\n", err)
	}

	fmt.Println("Agent run finished.")
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

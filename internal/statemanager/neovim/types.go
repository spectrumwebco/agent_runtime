package neovim

import (
	"time"

	"github.com/google/uuid"
)

type StateType string

const (
	ExecutionState StateType = "neovim_execution_state"

	ToolState StateType = "neovim_tool_state"

	ContextState StateType = "neovim_context_state"

	BufferState StateType = "neovim_buffer_state"

	UIState StateType = "neovim_ui_state"
)

var TTLMap = map[StateType]int{
	ExecutionState: 3600,   // 1 hour
	ToolState:      7200,   // 2 hours
	ContextState:   86400,  // 24 hours
	BufferState:    1800,   // 30 minutes
	UIState:        1800,   // 30 minutes
}

type Status string

const (
	Active Status = "active"

	Inactive Status = "inactive"

	Transitioning Status = "transitioning"
)

type Mode string

const (
	NormalMode Mode = "normal"

	InsertMode Mode = "insert"

	VisualMode Mode = "visual"

	CommandMode Mode = "command"
)

type CursorPosition struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

type ExecutionStateData struct {
	Status         Status         `json:"status"`
	CurrentMode    Mode           `json:"current_mode"`
	CurrentFile    string         `json:"current_file"`
	CursorPosition CursorPosition `json:"cursor_position"`
	LastCommand    string         `json:"last_command"`
	CommandHistory []string       `json:"command_history"`
	LastUpdated    time.Time      `json:"last_updated"`
}

type Plugin struct {
	Name   string                 `json:"name"`
	Status string                 `json:"status"`
	Config map[string]interface{} `json:"config"`
}

type Terminal struct {
	ID       string `json:"id"`
	Command  string `json:"command"`
	Status   string `json:"status"`
	BufferID int    `json:"buffer_id"`
}

type ToolStateData struct {
	ActivePlugins      []Plugin   `json:"active_plugins"`
	TerminalIntegration struct {
		ActiveTerminals []Terminal `json:"active_terminals"`
	} `json:"terminal_integration"`
	LastUpdated time.Time `json:"last_updated"`
}

type Workspace struct {
	Root       string   `json:"root"`
	OpenFiles  []string `json:"open_files"`
	FileHistory []string `json:"file_history"`
}

type Session struct {
	ID         string    `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	LastActive time.Time `json:"last_active"`
}

type ContextStateData struct {
	Workspace   Workspace `json:"workspace"`
	Session     Session   `json:"session"`
	LastUpdated time.Time `json:"last_updated"`
}

type Diagnostic struct {
	Line     int    `json:"line"`
	Column   int    `json:"column"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

type BufferStateData struct {
	BufferID     int          `json:"buffer_id"`
	FilePath     string       `json:"file_path"`
	Language     string       `json:"language"`
	Modified     bool         `json:"modified"`
	ContentHash  string       `json:"content_hash"`
	Diagnostics  []Diagnostic `json:"diagnostics"`
	LastUpdated  time.Time    `json:"last_updated"`
}

type Window struct {
	ID       int `json:"id"`
	BufferID int `json:"buffer_id"`
	Position struct {
		Row int `json:"row"`
		Col int `json:"col"`
	} `json:"position"`
	Size struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"size"`
}

type Font struct {
	Name string `json:"name"`
	Size int    `json:"size"`
}

type UIStateData struct {
	Layout struct {
		Windows       []Window `json:"windows"`
		CurrentWindow int      `json:"current_window"`
	} `json:"layout"`
	Colorscheme string    `json:"colorscheme"`
	Font        Font      `json:"font"`
	LastUpdated time.Time `json:"last_updated"`
}

type State struct {
	Type      StateType              `json:"state_type"`
	ID        string                 `json:"state_id"`
	Data      map[string]interface{} `json:"data"`
	TaskID    string                 `json:"task_id"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

func NewState(stateType StateType, data map[string]interface{}, taskID string) *State {
	now := time.Now().UTC()
	return &State{
		Type:      stateType,
		ID:        generateStateID(stateType),
		Data:      data,
		TaskID:    taskID,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func generateStateID(stateType StateType) string {
	return string(stateType) + "-" + uuid.New().String()
}

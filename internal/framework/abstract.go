package rex

import (
	"context"
	"time"
)

type IsAliveResponse struct {
	IsAlive bool   `json:"is_alive"`
	Message string `json:"message,omitempty"`
}

type CreateBashSessionRequest struct {
	StartupSource   []string `json:"startup_source,omitempty"`
	Session         string   `json:"session,omitempty"`
	SessionType     string   `json:"session_type"`
	StartupTimeout  float64  `json:"startup_timeout,omitempty"`
}

type CreateBashSessionResponse struct {
	Output      string `json:"output,omitempty"`
	SessionType string `json:"session_type"`
}

type BashAction struct {
	Command              string   `json:"command"`
	Session              string   `json:"session,omitempty"`
	Timeout              *float64 `json:"timeout,omitempty"`
	IsInteractiveCommand bool     `json:"is_interactive_command,omitempty"`
	IsInteractiveQuit    bool     `json:"is_interactive_quit,omitempty"`
	Check                string   `json:"check,omitempty"`
	ErrorMsg             string   `json:"error_msg,omitempty"`
	Expect               []string `json:"expect,omitempty"`
	ActionType           string   `json:"action_type"`
}

type BashInterruptAction struct {
	Session    string   `json:"session,omitempty"`
	Timeout    float64  `json:"timeout,omitempty"`
	NRetry     int      `json:"n_retry,omitempty"`
	Expect     []string `json:"expect,omitempty"`
	ActionType string   `json:"action_type"`
}

type BashObservation struct {
	Output        string `json:"output,omitempty"`
	ExitCode      *int   `json:"exit_code,omitempty"`
	FailureReason string `json:"failure_reason,omitempty"`
	ExpectString  string `json:"expect_string,omitempty"`
	SessionType   string `json:"session_type"`
}

type CloseBashSessionRequest struct {
	Session     string `json:"session,omitempty"`
	SessionType string `json:"session_type"`
}

type CloseBashSessionResponse struct {
	SessionType string `json:"session_type"`
}

type Command struct {
	Command   interface{}       `json:"command"`
	Timeout   *float64          `json:"timeout,omitempty"`
	Shell     bool              `json:"shell,omitempty"`
	Check     bool              `json:"check,omitempty"`
	ErrorMsg  string            `json:"error_msg,omitempty"`
	Env       map[string]string `json:"env,omitempty"`
	Cwd       string            `json:"cwd,omitempty"`
}

type CommandResponse struct {
	Stdout   string `json:"stdout,omitempty"`
	Stderr   string `json:"stderr,omitempty"`
	ExitCode *int   `json:"exit_code,omitempty"`
}

type ReadFileRequest struct {
	Path     string `json:"path"`
	Encoding string `json:"encoding,omitempty"`
	Errors   string `json:"errors,omitempty"`
}

type ReadFileResponse struct {
	Content string `json:"content,omitempty"`
}

type WriteFileRequest struct {
	Content string `json:"content"`
	Path    string `json:"path"`
}

type WriteFileResponse struct {}

type UploadRequest struct {
	SourcePath string `json:"source_path"`
	TargetPath string `json:"target_path"`
}

type UploadResponse struct {}

type CloseResponse struct {}

type Runtime interface {
	IsAlive(ctx context.Context, timeout *time.Duration) (*IsAliveResponse, error)
	
	CreateSession(ctx context.Context, request *CreateBashSessionRequest) (*CreateBashSessionResponse, error)
	
	RunInSession(ctx context.Context, action interface{}) (*BashObservation, error)
	
	CloseSession(ctx context.Context, request *CloseBashSessionRequest) (*CloseBashSessionResponse, error)
	
	Execute(ctx context.Context, command *Command) (*CommandResponse, error)
	
	ReadFile(ctx context.Context, request *ReadFileRequest) (*ReadFileResponse, error)
	
	WriteFile(ctx context.Context, request *WriteFileRequest) (*WriteFileResponse, error)
	
	Upload(ctx context.Context, request *UploadRequest) (*UploadResponse, error)
	
	Close(ctx context.Context) (*CloseResponse, error)
}

package langchain

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/tools"
)

func NewShellTool() tools.Tool {
	return tools.NewTool(
		"shell",
		"Run shell commands and get the output",
		func(ctx context.Context, input string) (string, error) {
			cmd := exec.CommandContext(ctx, "bash", "-c", input)

			output, err := cmd.CombinedOutput()
			if err != nil {
				return "", fmt.Errorf("failed to run command: %w\nOutput: %s", err, string(output))
			}

			return string(output), nil
		},
	)
}

func NewSearchTool() tools.Tool {
	return tools.NewTool(
		"search",
		"Search the web for information",
		func(ctx context.Context, input string) (string, error) {
			return fmt.Sprintf("Search results for: %s", input), nil
		},
	)
}

func NewCalculatorTool() tools.Tool {
	return tools.NewTool(
		"calculator",
		"Perform mathematical calculations",
		func(ctx context.Context, input string) (string, error) {
			cmd := exec.CommandContext(ctx, "bc", "-l")
			cmd.Stdin = strings.NewReader(input)

			output, err := cmd.CombinedOutput()
			if err != nil {
				return "", fmt.Errorf("failed to calculate: %w\nOutput: %s", err, string(output))
			}

			return strings.TrimSpace(string(output)), nil
		},
	)
}

func NewFileTool() tools.Tool {
	return tools.NewTool(
		"file",
		"Read and write files",
		func(ctx context.Context, input string) (string, error) {
			var params struct {
				Action string `json:"action"` // read, write, append, list
				Path   string `json:"path"`
				Content string `json:"content,omitempty"`
			}
			if err := json.Unmarshal([]byte(input), &params); err != nil {
				return "", fmt.Errorf("failed to parse input: %w", err)
			}

			switch params.Action {
			case "read":
				cmd := exec.CommandContext(ctx, "cat", params.Path)

				output, err := cmd.CombinedOutput()
				if err != nil {
					return "", fmt.Errorf("failed to read file: %w\nOutput: %s", err, string(output))
				}

				return string(output), nil
			case "write":
				cmd := exec.CommandContext(ctx, "bash", "-c", fmt.Sprintf("echo '%s' > %s", params.Content, params.Path))

				output, err := cmd.CombinedOutput()
				if err != nil {
					return "", fmt.Errorf("failed to write file: %w\nOutput: %s", err, string(output))
				}

				return fmt.Sprintf("File %s written successfully", params.Path), nil
			case "append":
				cmd := exec.CommandContext(ctx, "bash", "-c", fmt.Sprintf("echo '%s' >> %s", params.Content, params.Path))

				output, err := cmd.CombinedOutput()
				if err != nil {
					return "", fmt.Errorf("failed to append to file: %w\nOutput: %s", err, string(output))
				}

				return fmt.Sprintf("Content appended to %s successfully", params.Path), nil
			case "list":
				cmd := exec.CommandContext(ctx, "ls", "-la", params.Path)

				output, err := cmd.CombinedOutput()
				if err != nil {
					return "", fmt.Errorf("failed to list directory: %w\nOutput: %s", err, string(output))
				}

				return string(output), nil
			default:
				return "", fmt.Errorf("unsupported action: %s", params.Action)
			}
		},
	)
}

func NewHTTPTool() tools.Tool {
	return tools.NewTool(
		"http",
		"Make HTTP requests",
		func(ctx context.Context, input string) (string, error) {
			var params struct {
				Method  string            `json:"method"`
				URL     string            `json:"url"`
				Headers map[string]string `json:"headers,omitempty"`
				Body    string            `json:"body,omitempty"`
			}
			if err := json.Unmarshal([]byte(input), &params); err != nil {
				return "", fmt.Errorf("failed to parse input: %w", err)
			}

			cmd := exec.CommandContext(ctx, "curl", "-s", "-X", params.Method)

			for k, v := range params.Headers {
				cmd.Args = append(cmd.Args, "-H", fmt.Sprintf("%s: %s", k, v))
			}

			if params.Body != "" {
				cmd.Args = append(cmd.Args, "-d", params.Body)
			}

			cmd.Args = append(cmd.Args, params.URL)

			output, err := cmd.CombinedOutput()
			if err != nil {
				return "", fmt.Errorf("failed to make HTTP request: %w\nOutput: %s", err, string(output))
			}

			return string(output), nil
		},
	)
}

func NewCodeTool() tools.Tool {
	return tools.NewTool(
		"code",
		"Execute code in various languages",
		func(ctx context.Context, input string) (string, error) {
			var params struct {
				Language string `json:"language"`
				Code     string `json:"code"`
			}
			if err := json.Unmarshal([]byte(input), &params); err != nil {
				return "", fmt.Errorf("failed to parse input: %w", err)
			}

			switch strings.ToLower(params.Language) {
			case "python":
				cmd := exec.CommandContext(ctx, "bash", "-c", fmt.Sprintf("echo '%s' > /tmp/code.py", params.Code))
				if err := cmd.Run(); err != nil {
					return "", fmt.Errorf("failed to create temporary file: %w", err)
				}

				cmd = exec.CommandContext(ctx, "python", "/tmp/code.py")
				output, err := cmd.CombinedOutput()
				if err != nil {
					return "", fmt.Errorf("failed to execute Python code: %w\nOutput: %s", err, string(output))
				}

				return string(output), nil
			case "javascript", "js":
				cmd := exec.CommandContext(ctx, "bash", "-c", fmt.Sprintf("echo '%s' > /tmp/code.js", params.Code))
				if err := cmd.Run(); err != nil {
					return "", fmt.Errorf("failed to create temporary file: %w", err)
				}

				cmd = exec.CommandContext(ctx, "node", "/tmp/code.js")
				output, err := cmd.CombinedOutput()
				if err != nil {
					return "", fmt.Errorf("failed to execute JavaScript code: %w\nOutput: %s", err, string(output))
				}

				return string(output), nil
			case "bash", "shell":
				cmd := exec.CommandContext(ctx, "bash", "-c", fmt.Sprintf("echo '%s' > /tmp/code.sh", params.Code))
				if err := cmd.Run(); err != nil {
					return "", fmt.Errorf("failed to create temporary file: %w", err)
				}

				cmd = exec.CommandContext(ctx, "chmod", "+x", "/tmp/code.sh")
				if err := cmd.Run(); err != nil {
					return "", fmt.Errorf("failed to make file executable: %w", err)
				}

				cmd = exec.CommandContext(ctx, "bash", "/tmp/code.sh")
				output, err := cmd.CombinedOutput()
				if err != nil {
					return "", fmt.Errorf("failed to execute Bash code: %w\nOutput: %s", err, string(output))
				}

				return string(output), nil
			default:
				return "", fmt.Errorf("unsupported language: %s", params.Language)
			}
		},
	)
}

func ToolToFunctionDefinition(tool tools.Tool) schema.FunctionDefinition {
	return schema.FunctionDefinition{
		Name:        tool.Name(),
		Description: tool.Description(),
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"input": map[string]any{
					"type":        "string",
					"description": "The input to the tool",
				},
			},
			"required": []string{"input"},
		},
	}
}

func ToolToFunctionCallable(tool tools.Tool) schema.FunctionCallable {
	return func(ctx context.Context, args string) (string, error) {
		var params struct {
			Input string `json:"input"`
		}
		if err := json.Unmarshal([]byte(args), &params); err != nil {
			return "", fmt.Errorf("failed to parse input: %w", err)
		}

		return tool.Call(ctx, params.Input)
	}
}

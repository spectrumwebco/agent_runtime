package parser

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/spectrumwebco/agent_runtime/pkg/tools"
)

type Parser interface {
	Parse(modelOutput map[string]interface{}, availableCommands []*tools.Command) (thought string, action string, err error)
}

type ActionParser struct{}

func NewActionParser() *ActionParser {
	return &ActionParser{}
}

func (p *ActionParser) Parse(modelOutput map[string]interface{}, availableCommands []*tools.Command) (string, string, error) {
	msg, ok := modelOutput["message"].(string)
	if !ok {
		return "", "", fmt.Errorf("model output message is not a string")
	}

	return "", msg, nil
}

type ThoughtActionParser struct{}

func NewThoughtActionParser() *ThoughtActionParser {
	return &ThoughtActionParser{}
}

func (p *ThoughtActionParser) Parse(modelOutput map[string]interface{}, availableCommands []*tools.Command) (string, string, error) {
	msg, ok := modelOutput["message"].(string)
	if !ok {
		return "", "", fmt.Errorf("model output message is not a string")
	}

	re := regexp.MustCompile("```([\\s\\S]*?)```")
	matches := re.FindStringSubmatch(msg)
	if len(matches) < 2 {
		return "", "", fmt.Errorf("no code block found in message")
	}

	thought := strings.TrimSpace(strings.Split(msg, "```")[0])
	action := strings.TrimSpace(matches[1])
	return thought, action, nil
}

type XMLThoughtActionParser struct{}

func NewXMLThoughtActionParser() *XMLThoughtActionParser {
	return &XMLThoughtActionParser{}
}

func (p *XMLThoughtActionParser) Parse(modelOutput map[string]interface{}, availableCommands []*tools.Command) (string, string, error) {
	msg, ok := modelOutput["message"].(string)
	if !ok {
		return "", "", fmt.Errorf("model output message is not a string")
	}

	re := regexp.MustCompile("<command>([\\s\\S]*?)</command>")
	matches := re.FindStringSubmatch(msg)
	if len(matches) < 2 {
		return "", "", fmt.Errorf("no command tags found in message")
	}

	thought := strings.TrimSpace(strings.Split(msg, "<command>")[0])
	action := strings.TrimSpace(matches[1])
	return thought, action, nil
}

type FunctionCallingParser struct{}

func NewFunctionCallingParser() *FunctionCallingParser {
	return &FunctionCallingParser{}
}

type functionCall struct {
	Thought string                 `json:"thought"`
	Command map[string]interface{} `json:"command"`
}

func (p *FunctionCallingParser) Parse(modelOutput map[string]interface{}, availableCommands []*tools.Command) (string, string, error) {
	msg, ok := modelOutput["message"].(string)
	if !ok {
		return "", "", fmt.Errorf("model output message is not a string")
	}

	var call functionCall
	if err := json.Unmarshal([]byte(msg), &call); err != nil {
		return "", "", fmt.Errorf("failed to parse function call: %v", err)
	}

	if call.Command == nil {
		return "", "", fmt.Errorf("no command found in function call")
	}

	cmdBytes, err := json.Marshal(call.Command)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal command: %v", err)
	}

	return call.Thought, string(cmdBytes), nil
}

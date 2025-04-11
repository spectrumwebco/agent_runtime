package parser

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/spectrumwebco/agent_runtime/pkg/tools"
)

func TestActionParser(t *testing.T) {
	parser := NewActionParser()
	command := &tools.Command{Name: "ls", Description: ""}
	
	thought, action, err := parser.Parse(map[string]interface{}{
		"message": "ls -l",
	}, []*tools.Command{command})
	
	assert.NoError(t, err)
	assert.Equal(t, "", thought)
	assert.Equal(t, "ls -l", action)
	
	_, _, err = parser.Parse(map[string]interface{}{
		"message": "invalid command",
	}, []*tools.Command{command})
	assert.Error(t, err)
}

func TestThoughtActionParser(t *testing.T) {
	parser := NewThoughtActionParser()
	modelResponse := "Let's look at the files in the current directory.\n```\nls -l\n```"
	
	thought, action, err := parser.Parse(map[string]interface{}{
		"message": modelResponse,
	}, nil)
	
	assert.NoError(t, err)
	assert.Equal(t, "Let's look at the files in the current directory.", thought)
	assert.Equal(t, "ls -l", action)
	
	_, _, err = parser.Parse(map[string]interface{}{
		"message": "No code block",
	}, nil)
	assert.Error(t, err)
}

func TestXMLThoughtActionParser(t *testing.T) {
	parser := NewXMLThoughtActionParser()
	modelResponse := "Let's look at the files in the current directory.\n<command>\nls -l\n</command>"
	
	thought, action, err := parser.Parse(map[string]interface{}{
		"message": modelResponse,
	}, nil)
	
	assert.NoError(t, err)
	assert.Equal(t, "Let's look at the files in the current directory.", thought)
	assert.Equal(t, "ls -l", action)
	
	_, _, err = parser.Parse(map[string]interface{}{
		"message": "No command tags",
	}, nil)
	assert.Error(t, err)
}

func TestFunctionCallingParser(t *testing.T) {
	parser := NewFunctionCallingParser()
	command := &tools.Command{
		Name:        "ls",
		Description: "",
		Arguments:   []tools.Argument{},
	}
	
	validResponse := map[string]interface{}{
		"thought": "List files",
		"command": map[string]interface{}{
			"name": "ls",
			"arguments": map[string]interface{}{
				"path": ".",
			},
		},
	}
	responseJSON, _ := json.Marshal(validResponse)
	
	thought, action, err := parser.Parse(map[string]interface{}{
		"message": string(responseJSON),
	}, []*tools.Command{command})
	
	assert.NoError(t, err)
	assert.Equal(t, "List files", thought)
	assert.Contains(t, action, "ls")
	
	_, _, err = parser.Parse(map[string]interface{}{
		"message": "invalid json",
	}, []*tools.Command{command})
	assert.Error(t, err)
}

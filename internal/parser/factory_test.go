package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParserFactory(t *testing.T) {
	tests := []struct {
		name       string
		parserType string
		wantType   string
	}{
		{
			name:       "action parser",
			parserType: "action",
			wantType:   "*parser.ActionParser",
		},
		{
			name:       "thought_action parser",
			parserType: "thought_action",
			wantType:   "*parser.ThoughtActionParser",
		},
		{
			name:       "xml_thought_action parser",
			parserType: "xml_thought_action",
			wantType:   "*parser.XMLThoughtActionParser",
		},
		{
			name:       "function_calling parser",
			parserType: "function_calling",
			wantType:   "*parser.FunctionCallingParser",
		},
		{
			name:       "default to thought_action parser",
			parserType: "unknown",
			wantType:   "*parser.ThoughtActionParser",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := ParserFactory(tt.parserType)
			assert.NotNil(t, parser)
			assert.Equal(t, tt.wantType, assert.TypeOf(t, parser).String())
		})
	}
}

func TestGetAvailableParsers(t *testing.T) {
	parsers := GetAvailableParsers()
	assert.NotEmpty(t, parsers)
	assert.Contains(t, parsers, "action")
	assert.Contains(t, parsers, "thought_action")
	assert.Contains(t, parsers, "xml_thought_action")
	assert.Contains(t, parsers, "function_calling")
}

func TestGetParserDescription(t *testing.T) {
	tests := []struct {
		name       string
		parserType string
		wantDesc   string
	}{
		{
			name:       "action parser description",
			parserType: "action",
			wantDesc:   "Simple action parser that treats the entire model output as an action",
		},
		{
			name:       "thought_action parser description",
			parserType: "thought_action",
			wantDesc:   "Parser that extracts thought and action from code blocks (```action```)",
		},
		{
			name:       "xml_thought_action parser description",
			parserType: "xml_thought_action",
			wantDesc:   "Parser that extracts thought and action from XML tags (<command>action</command>)",
		},
		{
			name:       "function_calling parser description",
			parserType: "function_calling",
			wantDesc:   "Parser that extracts thought and action from function calling format (JSON)",
		},
		{
			name:       "unknown parser description",
			parserType: "unknown",
			wantDesc:   "Unknown parser type: unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			desc := GetParserDescription(tt.parserType)
			assert.Equal(t, tt.wantDesc, desc)
		})
	}
}

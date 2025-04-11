package parser

import (
	"fmt"
)

func ParserFactory(parserType string) Parser {
	switch parserType {
	case "action":
		return NewActionParser()
	case "thought_action":
		return NewThoughtActionParser()
	case "xml_thought_action":
		return NewXMLThoughtActionParser()
	case "function_calling":
		return NewFunctionCallingParser()
	default:
		return NewThoughtActionParser()
	}
}

func GetAvailableParsers() []string {
	return []string{
		"action",
		"thought_action",
		"xml_thought_action",
		"function_calling",
	}
}

func GetParserDescription(parserType string) string {
	switch parserType {
	case "action":
		return "Simple action parser that treats the entire model output as an action"
	case "thought_action":
		return "Parser that extracts thought and action from code blocks (```action```)"
	case "xml_thought_action":
		return "Parser that extracts thought and action from XML tags (<command>action</command>)"
	case "function_calling":
		return "Parser that extracts thought and action from function calling format (JSON)"
	default:
		return fmt.Sprintf("Unknown parser type: %s", parserType)
	}
}

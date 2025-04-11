package parsing

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type ParsedCommand struct {
	Name string
	Args map[string]interface{}
}

type CmdParser struct {
	cmdPattern *regexp.Regexp
}

func NewCmdParser() *CmdParser {
	pattern := `^([\w_]+)\((.*)\)$`
	return &CmdParser{
		cmdPattern: regexp.MustCompile(pattern),
	}
}

func (p *CmdParser) Parse(cmdStr string) (*ParsedCommand, error) {
	cmdStr = strings.TrimSpace(cmdStr)

	if strings.HasPrefix(cmdStr, "{") && strings.HasSuffix(cmdStr, "}") {
		var jsonCmd struct {
			Name string                 `json:"name"`
			Args map[string]interface{} `json:"args"`
		}
		if err := json.Unmarshal([]byte(cmdStr), &jsonCmd); err == nil {
			return &ParsedCommand{
				Name: jsonCmd.Name,
				Args: jsonCmd.Args,
			}, nil
		}
	}

	matches := p.cmdPattern.FindStringSubmatch(cmdStr)
	if len(matches) == 3 {
		toolName := matches[1]
		argsStr := strings.TrimSpace(matches[2])
		args, err := p.parseArgs(argsStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse arguments: %w", err)
		}
		return &ParsedCommand{Name: toolName, Args: args}, nil
	}

	if !strings.Contains(cmdStr, "(") && !strings.Contains(cmdStr, ")") {
		words := strings.Fields(cmdStr)
		if len(words) > 0 {
			toolName := words[0]
			args := make(map[string]interface{})
			if len(words) > 1 {
				args["args"] = strings.Join(words[1:], " ")
			}
			return &ParsedCommand{Name: toolName, Args: args}, nil
		}
	}

	return nil, fmt.Errorf("invalid command format: '%s'. Expected 'tool_name(args)', 'tool_name', or JSON", cmdStr)
}

func (p *CmdParser) parseArgs(argsStr string) (map[string]interface{}, error) {
	args := make(map[string]interface{})
	if argsStr == "" {
		return args, nil
	}

	var argPairs []string
	var currentPair strings.Builder
	inQuotes := false
	quoteChar := rune(0)
	inBrackets := 0

	for _, char := range argsStr {
		switch {
		case char == '"' || char == '\'':
			if !inQuotes {
				inQuotes = true
				quoteChar = char
			} else if char == quoteChar {
				inQuotes = false
				quoteChar = rune(0)
			}
			currentPair.WriteRune(char)
		case char == '[' || char == '{':
			inBrackets++
			currentPair.WriteRune(char)
		case char == ']' || char == '}':
			inBrackets--
			currentPair.WriteRune(char)
		case char == ',' && !inQuotes && inBrackets == 0:
			argPairs = append(argPairs, currentPair.String())
			currentPair.Reset()
		default:
			currentPair.WriteRune(char)
		}
	}
	if currentPair.Len() > 0 {
		argPairs = append(argPairs, currentPair.String())
	}

	for _, pair := range argPairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}

		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			parsedValue, err := p.parseValue(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse value for key '%s': %w", key, err)
			}
			args[key] = parsedValue
		} else {
			key := strings.TrimSpace(parts[0])
			if key != "" {
				args[key] = true
			}
		}
	}

	return args, nil
}

func (p *CmdParser) parseValue(value string) (interface{}, error) {
	if (strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]")) ||
		(strings.HasPrefix(value, "{") && strings.HasSuffix(value, "}")) {
		var jsonValue interface{}
		if err := json.Unmarshal([]byte(value), &jsonValue); err == nil {
			return jsonValue, nil
		}
	}

	if (strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`)) ||
		(strings.HasPrefix(value, `'`) && strings.HasSuffix(value, `'`)) {
		return value[1 : len(value)-1], nil
	}

	switch value {
	case "true":
		return true, nil
	case "false":
		return false, nil
	}

	if num, err := p.parseNumber(value); err == nil {
		return num, nil
	}

	return value, nil
}

func (p *CmdParser) parseNumber(s string) (interface{}, error) {
	if i, err := fmt.Sscan(s, new(int)); err == nil && i == 1 {
		var val int
		fmt.Sscan(s, &val)
		return val, nil
	}
	if f, err := fmt.Sscan(s, new(float64)); err == nil && f == 1 {
		var val float64
		fmt.Sscan(s, &val)
		return val, nil
	}
	return nil, fmt.Errorf("not a number")
}

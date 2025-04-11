package parsing

import (
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
	matches := p.cmdPattern.FindStringSubmatch(cmdStr)

	if len(matches) != 3 {
		if !strings.Contains(cmdStr, "(") && !strings.Contains(cmdStr, ")") {
			if regexp.MustCompile(`^[\w_]+$`).MatchString(cmdStr) {
				return &ParsedCommand{Name: cmdStr, Args: make(map[string]interface{})}, nil
			}
		}
		return nil, fmt.Errorf("invalid command format: '%s'. Expected 'tool_name(args)' or 'tool_name'", cmdStr)
	}

	toolName := matches[1]
	argsStr := strings.TrimSpace(matches[2])
	args := make(map[string]interface{})

	if argsStr == "" {
		return &ParsedCommand{Name: toolName, Args: args}, nil
	}

	argPairs := strings.Split(argsStr, ",") // TODO: Handle commas within quotes/brackets
	for _, pair := range argPairs {
		parts := strings.SplitN(strings.TrimSpace(pair), "=", 2)
		if len(parts) != 2 {
			argName := strings.TrimSpace(parts[0])
			if argName != "" {
				args[argName] = true // Treat valueless args as boolean true
			}
			continue
		}
		argName := strings.TrimSpace(parts[0])
		argValueStr := strings.TrimSpace(parts[1])

		if (strings.HasPrefix(argValueStr, `"`) && strings.HasSuffix(argValueStr, `"`)) ||
			(strings.HasPrefix(argValueStr, `'`) && strings.HasSuffix(argValueStr, `'`)) {
			args[argName] = argValueStr[1 : len(argValueStr)-1] // Store as string without quotes
		} else if argValueStr == "true" {
			args[argName] = true
		} else if argValueStr == "false" {
			args[argName] = false
		} else if num, err := p.parseNumber(argValueStr); err == nil {
			args[argName] = num // Store as float64 or int
		} else {
			args[argName] = argValueStr // Default to string if no other type matches
		}
	}

	return &ParsedCommand{Name: toolName, Args: args}, nil
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

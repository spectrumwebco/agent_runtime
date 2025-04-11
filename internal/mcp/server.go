package mcp

import (
	"fmt"
)

type Host struct {
}

func NewHost(/* TODO: Add config parameter */) *Host {
	return &Host{}
}

func (h *Host) Start() {
	fmt.Println("Starting MCP Host and Servers...")
}

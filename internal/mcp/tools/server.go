package tools

import "fmt"

type Server struct {
}

func NewServer(/* config */) *Server {
	return &Server{}
}

func (s *Server) Run() {
	fmt.Println("Tool MCP Server placeholder running...")
	select {} // Keep running
}

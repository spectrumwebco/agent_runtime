package filesystem

import "fmt"

type Server struct {
}

func NewServer(/* config */) *Server {
	return &Server{}
}

func (s *Server) Run() {
	fmt.Println("Filesystem MCP Server placeholder running...")
	select {} // Keep running
}

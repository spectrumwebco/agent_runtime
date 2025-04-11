package mcp

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mark3labs/mcp-go/server"
)

type Host struct {
	server      *server.MCPServer
	clientURL   string
	serverURL   string
	apiKey      string
	httpClient  *http.Client
}

func NewHost(clientURL, serverURL string) *Host {
	apiKey := os.Getenv("LIBRECHAT_CODE_API_KEY")
	if apiKey == "" {
		log.Fatal("LIBRECHAT_CODE_API_KEY environment variable not set")
	}

	mcpServer := server.NewMCPServer(
		"agent-runtime/host",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	return &Host{
		server:     mcpServer,
		clientURL:  clientURL,
		serverURL:  serverURL,
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

func (h *Host) Start() {
	fmt.Println("Starting MCP Host and Servers...")

	h.server.AddTool(server.NewTool("proxy_request",
		server.WithDescription("Proxies requests between client and server"),
		server.WithString("target",
			server.Description("Target URL for the request"),
			server.Required(),
		),
		server.WithObject("payload",
			server.Description("Request payload"),
			server.Required(),
		),
		h.handleProxyRequest,
	))

	if err := h.server.Start(); err != nil {
		log.Fatalf("Failed to start MCP host: %v", err)
	}
}

func (h *Host) handleProxyRequest(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	target := args["target"].(string)
	payload := args["payload"].(map[string]interface{})

	var url string
	switch target {
	case "client":
		url = h.clientURL
	case "server":
		url = h.serverURL
	default:
		return nil, fmt.Errorf("invalid target: %s", target)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", h.apiKey)

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return result, nil
}

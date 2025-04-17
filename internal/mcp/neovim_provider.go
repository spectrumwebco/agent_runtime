package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/spectrumwebco/agent_runtime/pkg/interfaces"
)

type NeovimMCPProvider struct {
	terminalManager interfaces.TerminalManager
	apiURL          string
}

func NewNeovimMCPProvider(terminalManager interfaces.TerminalManager, apiURL string) *NeovimMCPProvider {
	return &NeovimMCPProvider{
		terminalManager: terminalManager,
		apiURL:          apiURL,
	}
}

func (p *NeovimMCPProvider) RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/mcp/neovim/create", p.handleCreate)
	mux.HandleFunc("/mcp/neovim/execute", p.handleExecute)
	mux.HandleFunc("/mcp/neovim/stop", p.handleStop)
	mux.HandleFunc("/mcp/neovim/list", p.handleList)
	mux.HandleFunc("/mcp/neovim/create_bulk", p.handleCreateBulk)
	mux.HandleFunc("/mcp/neovim/execute_all", p.handleExecuteAll)
}

func (p *NeovimMCPProvider) handleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var params struct {
		ID      string                 `json:"id"`
		Options map[string]interface{} `json:"options"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	if params.Options == nil {
		params.Options = make(map[string]interface{})
	}
	params.Options["api_url"] = p.apiURL

	terminal, err := p.terminalManager.CreateTerminal(r.Context(), "neovim", params.ID, params.Options)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create terminal: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"id":      terminal.ID(),
		"type":    terminal.GetType(),
		"running": terminal.IsRunning(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (p *NeovimMCPProvider) handleExecute(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var params struct {
		ID      string `json:"id"`
		Command string `json:"command"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	terminal, err := p.terminalManager.GetTerminal(params.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Terminal not found: %v", err), http.StatusNotFound)
		return
	}

	output, err := terminal.Execute(r.Context(), params.Command)
	if err != nil {
		http.Error(w, fmt.Sprintf("Command execution failed: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"id":      terminal.ID(),
		"output":  output,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (p *NeovimMCPProvider) handleStop(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var params struct {
		ID string `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	err := p.terminalManager.RemoveTerminal(r.Context(), params.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to stop terminal: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"id":      params.ID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (p *NeovimMCPProvider) handleList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	terminals := p.terminalManager.ListTerminals()
	result := make([]map[string]interface{}, 0, len(terminals))

	for _, terminal := range terminals {
		result = append(result, map[string]interface{}{
			"id":      terminal.ID(),
			"type":    terminal.GetType(),
			"running": terminal.IsRunning(),
		})
	}

	response := map[string]interface{}{
		"success":   true,
		"terminals": result,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (p *NeovimMCPProvider) handleCreateBulk(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var params struct {
		Count   int                    `json:"count"`
		Options map[string]interface{} `json:"options"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	if params.Options == nil {
		params.Options = make(map[string]interface{})
	}
	params.Options["api_url"] = p.apiURL

	terminals, err := p.terminalManager.CreateBulkTerminals(r.Context(), "neovim", params.Count, params.Options)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create bulk terminals: %v", err), http.StatusInternalServerError)
		return
	}

	result := make([]map[string]interface{}, 0, len(terminals))
	for _, terminal := range terminals {
		result = append(result, map[string]interface{}{
			"id":      terminal.ID(),
			"type":    terminal.GetType(),
			"running": terminal.IsRunning(),
		})
	}

	response := map[string]interface{}{
		"success":   true,
		"terminals": result,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (p *NeovimMCPProvider) handleExecuteAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var params struct {
		Command string `json:"command"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	results := make(map[string]string)
	terminals := p.terminalManager.ListTerminals()

	for _, terminal := range terminals {
		output, err := terminal.Execute(r.Context(), params.Command)
		if err != nil {
			results[terminal.ID()] = fmt.Sprintf("ERROR: %s", err.Error())
		} else {
			results[terminal.ID()] = output
		}
	}

	response := map[string]interface{}{
		"success": true,
		"results": results,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	log.Println("Testing Neovim state management...")

	log.Println("Test 1: Creating Neovim instance...")
	neovimID := "state-test"
	err := startNeovimInstance(neovimID)
	if err != nil {
		log.Fatalf("Failed to create Neovim instance: %v", err)
	}
	log.Println("Neovim instance created successfully")

	log.Println("Test 2: Creating a test file in Neovim...")
	testFilePath := "/tmp/neovim-state-test.txt"
	err = executeNeovimCommand(neovimID, fmt.Sprintf(":e %s<CR>iThis is a test file created by Neovim<Esc>:w<CR>", testFilePath))
	if err != nil {
		log.Fatalf("Failed to create test file: %v", err)
	}
	
	if _, err := os.Stat(testFilePath); os.IsNotExist(err) {
		log.Fatalf("Test file was not created: %v", err)
	}
	log.Println("Test file created successfully")

	log.Println("Test 3: Testing execution state...")
	executionState := map[string]interface{}{
		"status":         "active",
		"current_mode":   "normal",
		"current_file":   testFilePath,
		"cursor_position": map[string]interface{}{
			"line":   1,
			"column": 0,
		},
		"last_command":    ":w",
		"command_history": []string{":e", "i", "<Esc>", ":w"},
	}
	
	err = storeNeovimState("neovim_execution_state", "execution", executionState, "state-test-task")
	if err != nil {
		log.Fatalf("Failed to store execution state: %v", err)
	}
	
	retrievedState, exists := retrieveNeovimState("neovim_execution_state", "execution", "state-test-task")
	if !exists {
		log.Fatalf("Failed to retrieve execution state")
	}
	
	if retrievedState["status"] != "active" {
		log.Fatalf("Execution state status mismatch: expected 'active', got '%v'", retrievedState["status"])
	}
	log.Println("Execution state test passed")

	log.Println("Test 4: Testing buffer state...")
	bufferState := map[string]interface{}{
		"buffer_id":    1,
		"file_path":    testFilePath,
		"language":     "text",
		"modified":     false,
		"content_hash": "abc123",
		"diagnostics":  []interface{}{},
	}
	
	err = storeNeovimState("neovim_buffer_state", "buffer-1", bufferState, "state-test-task")
	if err != nil {
		log.Fatalf("Failed to store buffer state: %v", err)
	}
	
	retrievedState, exists = retrieveNeovimState("neovim_buffer_state", "buffer-1", "state-test-task")
	if !exists {
		log.Fatalf("Failed to retrieve buffer state")
	}
	
	if retrievedState["file_path"] != testFilePath {
		log.Fatalf("Buffer state file_path mismatch: expected '%s', got '%v'", testFilePath, retrievedState["file_path"])
	}
	log.Println("Buffer state test passed")

	log.Println("Test 5: Testing state update...")
	updates := map[string]interface{}{
		"modified": true,
	}
	
	err = updateNeovimState("neovim_buffer_state", "buffer-1", updates, "state-test-task")
	if err != nil {
		log.Fatalf("Failed to update buffer state: %v", err)
	}
	
	retrievedState, exists = retrieveNeovimState("neovim_buffer_state", "buffer-1", "state-test-task")
	if !exists {
		log.Fatalf("Failed to retrieve updated buffer state")
	}
	
	if retrievedState["modified"] != true {
		log.Fatalf("Buffer state modified mismatch: expected 'true', got '%v'", retrievedState["modified"])
	}
	log.Println("State update test passed")

	log.Println("Test 6: Cleaning up...")
	err = stopNeovimInstance(neovimID)
	if err != nil {
		log.Fatalf("Failed to stop Neovim instance: %v", err)
	}
	
	err = os.Remove(testFilePath)
	if err != nil {
		log.Printf("Warning: Failed to remove test file: %v", err)
	}
	
	log.Println("All Neovim state management tests passed!")
}


func startNeovimInstance(id string) error {
	url := "http://localhost:8090/neovim/start"
	payload := fmt.Sprintf(`{"id":"%s"}`, id)
	
	resp, err := http.Post(url, "application/json", strings.NewReader(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to start Neovim instance, status: %d", resp.StatusCode)
	}
	
	return nil
}

func stopNeovimInstance(id string) error {
	url := "http://localhost:8090/neovim/stop"
	payload := fmt.Sprintf(`{"id":"%s"}`, id)
	
	resp, err := http.Post(url, "application/json", strings.NewReader(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to stop Neovim instance, status: %d", resp.StatusCode)
	}
	
	return nil
}

func executeNeovimCommand(id string, command string) error {
	url := "http://localhost:8090/neovim/execute"
	payload := fmt.Sprintf(`{"id":"%s","command":"%s"}`, id, command)
	
	resp, err := http.Post(url, "application/json", strings.NewReader(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to execute Neovim command, status: %d", resp.StatusCode)
	}
	
	return nil
}


func storeNeovimState(stateType string, stateID string, stateData map[string]interface{}, taskID string) error {
	stateData["last_updated"] = time.Now().UTC().Format(time.RFC3339)
	
	jsonData, err := json.Marshal(stateData)
	if err != nil {
		return err
	}
	
	filename := fmt.Sprintf("/tmp/neovim-state-%s-%s-%s.json", taskID, stateType, stateID)
	return os.WriteFile(filename, jsonData, 0644)
}

func retrieveNeovimState(stateType string, stateID string, taskID string) (map[string]interface{}, bool) {
	filename := fmt.Sprintf("/tmp/neovim-state-%s-%s-%s.json", taskID, stateType, stateID)
	
	jsonData, err := os.ReadFile(filename)
	if err != nil {
		return nil, false
	}
	
	var stateData map[string]interface{}
	err = json.Unmarshal(jsonData, &stateData)
	if err != nil {
		return nil, false
	}
	
	return stateData, true
}

func updateNeovimState(stateType string, stateID string, updates map[string]interface{}, taskID string) error {
	stateData, exists := retrieveNeovimState(stateType, stateID, taskID)
	if !exists {
		return fmt.Errorf("state not found")
	}
	
	for k, v := range updates {
		stateData[k] = v
	}
	
	stateData["last_updated"] = time.Now().UTC().Format(time.RFC3339)
	
	return storeNeovimState(stateType, stateID, stateData, taskID)
}

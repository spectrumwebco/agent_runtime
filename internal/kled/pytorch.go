package kled

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/spectrumwebco/agent_runtime/internal/ffi/python"
)

type PyTorchIntegration struct {
	Framework *Framework

	ScriptRunner *python.ScriptRunner

	mutex sync.RWMutex

	Context map[string]interface{}

	ModelsDir string
	
	ScriptsDir string
}

var _ Integration = (*PyTorchIntegration)(nil)

type PyTorchIntegrationConfig struct {
	ModelsDir string `json:"models_dir"`

	ScriptsDir string `json:"scripts_dir"`

	RocketMQAddrs []string `json:"rocket_mq_addrs"`

	EventTopic string `json:"event_topic"`
}

func NewPyTorchIntegration(framework *Framework, config PyTorchIntegrationConfig) (*PyTorchIntegration, error) {
	if framework == nil {
		return nil, fmt.Errorf("framework cannot be nil")
	}

	if err := os.MkdirAll(config.ModelsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create models directory: %w", err)
	}

	if err := os.MkdirAll(config.ScriptsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create scripts directory: %w", err)
	}

	scriptRunner, err := python.NewScriptRunner(
		config.ScriptsDir,
		config.RocketMQAddrs,
		config.EventTopic,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create script runner: %w", err)
	}

	integration := &PyTorchIntegration{
		Framework:    framework,
		ScriptRunner: scriptRunner,
		Context:      make(map[string]interface{}),
		ModelsDir:    config.ModelsDir,
		ScriptsDir:   config.ScriptsDir,
	}

	if err := framework.RegisterIntegration(integration); err != nil {
		return nil, fmt.Errorf("failed to register PyTorch integration: %w", err)
	}

	if framework.Config.Debug {
		log.Printf("PyTorch integration created with models directory: %s\n", config.ModelsDir)
	}

	return integration, nil
}

func (p *PyTorchIntegration) Name() string {
	return "pytorch"
}

func (p *PyTorchIntegration) Start() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if !IsPyTorchInstalled() {
		if p.Framework.Config.Debug {
			log.Println("PyTorch not found, attempting to install...")
		}
		if err := InstallPyTorch(); err != nil {
			return fmt.Errorf("failed to install PyTorch: %w", err)
		}
	}

	if p.Framework.Config.Debug {
		log.Println("PyTorch integration started")
	}

	return nil
}

func (p *PyTorchIntegration) Stop() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if err := p.ScriptRunner.Close(); err != nil {
		return fmt.Errorf("failed to close script runner: %w", err)
	}

	if p.Framework.Config.Debug {
		log.Println("PyTorch integration stopped")
	}

	return nil
}

func (p *PyTorchIntegration) RunInference(ctx context.Context, modelName string, input map[string]interface{}) (interface{}, error) {
	args := map[string]interface{}{
		"model_name": modelName,
		"input":      input,
	}

	result, err := p.ScriptRunner.RunToolScript("run_inference", args)
	if err != nil {
		return nil, fmt.Errorf("failed to run inference: %w", err)
	}

	return result, nil
}

func (p *PyTorchIntegration) TrainModel(ctx context.Context, modelName string, config map[string]interface{}) (interface{}, error) {
	args := map[string]interface{}{
		"model_name": modelName,
		"config":     config,
	}

	result, err := p.ScriptRunner.RunToolScript("train_model", args)
	if err != nil {
		return nil, fmt.Errorf("failed to train model: %w", err)
	}

	return result, nil
}

func (p *PyTorchIntegration) SaveModel(ctx context.Context, modelName string, path string) (interface{}, error) {
	if !filepath.IsAbs(path) {
		path = filepath.Join(p.ModelsDir, path)
	}

	args := map[string]interface{}{
		"model_name": modelName,
		"path":       path,
	}

	result, err := p.ScriptRunner.RunToolScript("save_model", args)
	if err != nil {
		return nil, fmt.Errorf("failed to save model: %w", err)
	}

	return result, nil
}

func (p *PyTorchIntegration) LoadModel(ctx context.Context, modelName string, path string) (interface{}, error) {
	if !filepath.IsAbs(path) {
		path = filepath.Join(p.ModelsDir, path)
	}

	args := map[string]interface{}{
		"model_name": modelName,
		"path":       path,
	}

	result, err := p.ScriptRunner.RunToolScript("load_model", args)
	if err != nil {
		return nil, fmt.Errorf("failed to load model: %w", err)
	}

	return result, nil
}

func (p *PyTorchIntegration) EvaluateModel(ctx context.Context, modelName string, data map[string]interface{}) (interface{}, error) {
	args := map[string]interface{}{
		"model_name": modelName,
		"data":       data,
	}

	result, err := p.ScriptRunner.RunToolScript("evaluate_model", args)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate model: %w", err)
	}

	return result, nil
}

func (p *PyTorchIntegration) GetContext() map[string]interface{} {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	contextCopy := make(map[string]interface{})
	for k, v := range p.Context {
		contextCopy[k] = v
	}

	return contextCopy
}

func (p *PyTorchIntegration) SetContext(ctx map[string]interface{}) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.Context = ctx
}

func (p *PyTorchIntegration) UpdateContext(updates map[string]interface{}) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for k, v := range updates {
		p.Context[k] = v
	}
}

func (p *PyTorchIntegration) SaveState() ([]byte, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	state := map[string]interface{}{
		"context": p.Context,
	}

	return json.Marshal(state)
}

func (p *PyTorchIntegration) LoadState(data []byte) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	var state map[string]interface{}
	if err := json.Unmarshal(data, &state); err != nil {
		return fmt.Errorf("failed to unmarshal state: %w", err)
	}

	if context, ok := state["context"].(map[string]interface{}); ok {
		p.Context = context
	}

	return nil
}

func IsPyTorchInstalled() bool {
	cmd := exec.Command("python", "-c", "import torch; print(torch.__version__)")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func InstallPyTorch() error {
	cmd := exec.Command("pip", "install", "torch", "torchvision", "torchaudio")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install PyTorch: %w", err)
	}
	return nil
}

func (p *PyTorchIntegration) CreatePyTorchScript(name string, content string) error {
	scriptPath := filepath.Join(p.ScriptsDir, name+".py")
	
	if err := os.WriteFile(scriptPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to create PyTorch script: %w", err)
	}
	
	if p.Framework.Config.Debug {
		log.Printf("Created PyTorch script: %s\n", scriptPath)
	}
	
	return nil
}

func (p *PyTorchIntegration) ListModels() ([]string, error) {
	files, err := os.ReadDir(p.ModelsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to list models: %w", err)
	}
	
	var models []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".pt" {
			models = append(models, file.Name())
		}
	}
	
	return models, nil
}

func (p *PyTorchIntegration) DeleteModel(modelName string) error {
	modelPath := filepath.Join(p.ModelsDir, modelName)
	
	if err := os.Remove(modelPath); err != nil {
		return fmt.Errorf("failed to delete model: %w", err)
	}
	
	if p.Framework.Config.Debug {
		log.Printf("Deleted model: %s\n", modelPath)
	}
	
	return nil
}

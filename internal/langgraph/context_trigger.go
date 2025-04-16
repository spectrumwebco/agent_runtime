package langgraph

import (
	"context"
	"fmt"
	"regexp"
	"sync"
	"time"
)

type ContextTrigger struct {
	ID              string
	Name            string
	Description     string
	ContextPatterns []string
	ToolName        string
	Priority        int
	Enabled         bool
	LastTriggered   time.Time
	mu              sync.RWMutex
	compiledRegexps []*regexp.Regexp
}

func NewContextTrigger(id, name, description string, contextPatterns []string, toolName string, priority int) (*ContextTrigger, error) {
	compiledRegexps := make([]*regexp.Regexp, 0, len(contextPatterns))
	
	for _, pattern := range contextPatterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("failed to compile regex pattern %s: %w", pattern, err)
		}
		compiledRegexps = append(compiledRegexps, re)
	}
	
	return &ContextTrigger{
		ID:              id,
		Name:            name,
		Description:     description,
		ContextPatterns: contextPatterns,
		ToolName:        toolName,
		Priority:        priority,
		Enabled:         true,
		compiledRegexps: compiledRegexps,
	}, nil
}

func (ct *ContextTrigger) ShouldTrigger(contextData map[string]interface{}) bool {
	ct.mu.RLock()
	defer ct.mu.RUnlock()
	
	if !ct.Enabled {
		return false
	}
	
	for key, value := range contextData {
		strValue, ok := value.(string)
		if !ok {
			continue
		}
		
		for _, re := range ct.compiledRegexps {
			if re.MatchString(strValue) {
				return true
			}
		}
	}
	
	return false
}

func (ct *ContextTrigger) Execute(ctx context.Context, ai *AgentIntegration, contextData map[string]interface{}) (map[string]interface{}, error) {
	ct.mu.Lock()
	ct.LastTriggered = time.Now()
	ct.mu.Unlock()
	
	payload := map[string]interface{}{
		"tool":  ct.ToolName,
		"input": contextData,
		"state": map[string]interface{}{
			"trigger_id":   ct.ID,
			"trigger_name": ct.Name,
			"triggered_at": time.Now().Unix(),
		},
	}
	
	return ai.callDjangoAgent("execute_tool", payload)
}

type ContextTriggerManager struct {
	triggers     []*ContextTrigger
	integration  *AgentIntegration
	mu           sync.RWMutex
}

func NewContextTriggerManager(integration *AgentIntegration) *ContextTriggerManager {
	return &ContextTriggerManager{
		triggers:    make([]*ContextTrigger, 0),
		integration: integration,
	}
}

func (ctm *ContextTriggerManager) AddTrigger(trigger *ContextTrigger) {
	ctm.mu.Lock()
	defer ctm.mu.Unlock()
	
	ctm.triggers = append(ctm.triggers, trigger)
}

func (ctm *ContextTriggerManager) RemoveTrigger(triggerID string) bool {
	ctm.mu.Lock()
	defer ctm.mu.Unlock()
	
	for i, trigger := range ctm.triggers {
		if trigger.ID == triggerID {
			ctm.triggers = append(ctm.triggers[:i], ctm.triggers[i+1:]...)
			return true
		}
	}
	
	return false
}

func (ctm *ContextTriggerManager) GetTrigger(triggerID string) *ContextTrigger {
	ctm.mu.RLock()
	defer ctm.mu.RUnlock()
	
	for _, trigger := range ctm.triggers {
		if trigger.ID == triggerID {
			return trigger
		}
	}
	
	return nil
}

func (ctm *ContextTriggerManager) ProcessContext(ctx context.Context, contextData map[string]interface{}) ([]map[string]interface{}, error) {
	ctm.mu.RLock()
	triggersToExecute := make([]*ContextTrigger, 0)
	
	for _, trigger := range ctm.triggers {
		if trigger.ShouldTrigger(contextData) {
			triggersToExecute = append(triggersToExecute, trigger)
		}
	}
	ctm.mu.RUnlock()
	
	for i := 0; i < len(triggersToExecute)-1; i++ {
		for j := i + 1; j < len(triggersToExecute); j++ {
			if triggersToExecute[i].Priority < triggersToExecute[j].Priority {
				triggersToExecute[i], triggersToExecute[j] = triggersToExecute[j], triggersToExecute[i]
			}
		}
	}
	
	results := make([]map[string]interface{}, 0, len(triggersToExecute))
	
	for _, trigger := range triggersToExecute {
		result, err := trigger.Execute(ctx, ctm.integration, contextData)
		if err != nil {
			return results, err
		}
		
		results = append(results, result)
	}
	
	return results, nil
}

func (ctm *ContextTriggerManager) CreateKnowledgeToolTrigger() (*ContextTrigger, error) {
	patterns := []string{
		"(?i)knowledge base",
		"(?i)documentation",
		"(?i)learn about",
		"(?i)information on",
		"(?i)how (to|do) .+\\?",
		"(?i)what is .+\\?",
	}
	
	return NewContextTrigger(
		"knowledge-tool-trigger",
		"Knowledge Tool Trigger",
		"Triggers the knowledge tool when context indicates a need for information",
		patterns,
		"knowledge",
		100,
	)
}

func (ctm *ContextTriggerManager) CreateCodeAnalysisTrigger() (*ContextTrigger, error) {
	patterns := []string{
		"(?i)analyze (the )?code",
		"(?i)code review",
		"(?i)check (the )?code",
		"(?i)examine (the )?code",
		"(?i)code quality",
		"(?i)code metrics",
	}
	
	return NewContextTrigger(
		"code-analysis-trigger",
		"Code Analysis Trigger",
		"Triggers the code analysis tool when context indicates a need for code review",
		patterns,
		"code_analysis",
		90,
	)
}

func (ctm *ContextTriggerManager) CreateSearchTrigger() (*ContextTrigger, error) {
	patterns := []string{
		"(?i)search for",
		"(?i)find information",
		"(?i)look up",
		"(?i)research",
		"(?i)explore",
	}
	
	return NewContextTrigger(
		"search-trigger",
		"Search Trigger",
		"Triggers the search tool when context indicates a need for information search",
		patterns,
		"search",
		80,
	)
}

func (ctm *ContextTriggerManager) CreatePlanningTrigger() (*ContextTrigger, error) {
	patterns := []string{
		"(?i)create (a )?plan",
		"(?i)plan (for|to)",
		"(?i)strategy",
		"(?i)roadmap",
		"(?i)outline steps",
	}
	
	return NewContextTrigger(
		"planning-trigger",
		"Planning Trigger",
		"Triggers the planning tool when context indicates a need for planning",
		patterns,
		"planning",
		95,
	)
}

func (ctm *ContextTriggerManager) InitializeDefaultTriggers() error {
	knowledgeTrigger, err := ctm.CreateKnowledgeToolTrigger()
	if err != nil {
		return err
	}
	ctm.AddTrigger(knowledgeTrigger)
	
	codeAnalysisTrigger, err := ctm.CreateCodeAnalysisTrigger()
	if err != nil {
		return err
	}
	ctm.AddTrigger(codeAnalysisTrigger)
	
	searchTrigger, err := ctm.CreateSearchTrigger()
	if err != nil {
		return err
	}
	ctm.AddTrigger(searchTrigger)
	
	planningTrigger, err := ctm.CreatePlanningTrigger()
	if err != nil {
		return err
	}
	ctm.AddTrigger(planningTrigger)
	
	return nil
}

type ContextMonitor struct {
	triggerManager *ContextTriggerManager
	contextHistory []map[string]interface{}
	maxHistorySize int
	running        bool
	stopCh         chan struct{}
	mu             sync.RWMutex
}

func NewContextMonitor(triggerManager *ContextTriggerManager, maxHistorySize int) *ContextMonitor {
	if maxHistorySize <= 0 {
		maxHistorySize = 100
	}
	
	return &ContextMonitor{
		triggerManager: triggerManager,
		contextHistory: make([]map[string]interface{}, 0, maxHistorySize),
		maxHistorySize: maxHistorySize,
		running:        false,
		stopCh:         make(chan struct{}),
	}
}

func (cm *ContextMonitor) Start(ctx context.Context) {
	cm.mu.Lock()
	if cm.running {
		cm.mu.Unlock()
		return
	}
	
	cm.running = true
	cm.stopCh = make(chan struct{})
	cm.mu.Unlock()
	
	go cm.monitorLoop(ctx)
}

func (cm *ContextMonitor) Stop() {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	if !cm.running {
		return
	}
	
	close(cm.stopCh)
	cm.running = false
}

func (cm *ContextMonitor) AddContext(contextData map[string]interface{}) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	cm.contextHistory = append(cm.contextHistory, contextData)
	
	if len(cm.contextHistory) > cm.maxHistorySize {
		cm.contextHistory = cm.contextHistory[len(cm.contextHistory)-cm.maxHistorySize:]
	}
}

func (cm *ContextMonitor) GetContextHistory() []map[string]interface{} {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	
	history := make([]map[string]interface{}, len(cm.contextHistory))
	copy(history, cm.contextHistory)
	
	return history
}

func (cm *ContextMonitor) monitorLoop(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			cm.processLatestContext(ctx)
		case <-cm.stopCh:
			return
		case <-ctx.Done():
			return
		}
	}
}

func (cm *ContextMonitor) processLatestContext(ctx context.Context) {
	cm.mu.RLock()
	if len(cm.contextHistory) == 0 {
		cm.mu.RUnlock()
		return
	}
	
	latestContext := cm.contextHistory[len(cm.contextHistory)-1]
	cm.mu.RUnlock()
	
	_, _ = cm.triggerManager.ProcessContext(ctx, latestContext)
}

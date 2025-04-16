package vercel

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/spectrumwebco/agent_runtime/pkg/eventstream"
	"github.com/spectrumwebco/agent_runtime/pkg/sharedstate"
)

type ComponentType string

const (
	ButtonComponent      ComponentType = "button"
	CardComponent        ComponentType = "card"
	FormComponent        ComponentType = "form"
	TableComponent       ComponentType = "table"
	ChartComponent       ComponentType = "chart"
	ProgressComponent    ComponentType = "progress"
	AlertComponent       ComponentType = "alert"
	ModalComponent       ComponentType = "modal"
	TabsComponent        ComponentType = "tabs"
	AccordionComponent   ComponentType = "accordion"
	TimelineComponent    ComponentType = "timeline"
	CodeBlockComponent   ComponentType = "codeblock"
	TerminalComponent    ComponentType = "terminal"
	MarkdownComponent    ComponentType = "markdown"
	ContainerComponent   ComponentType = "container"
	ToolOutputComponent  ComponentType = "tooloutput"
	AgentActionComponent ComponentType = "agentaction"
)

type RSCAdapter struct {
	sharedStateManager *sharedstate.SharedStateManager
	eventStream        *eventstream.Stream
	aiClient           *VercelAIClient
	componentRegistry  map[ComponentType]ComponentGenerator
	mu                 sync.RWMutex
	componentHistory   []GeneratedComponent
}

type ComponentGenerator func(ctx context.Context, props map[string]interface{}) (GeneratedComponent, error)

type GeneratedComponent struct {
	ID           string                 `json:"id"`
	Type         ComponentType          `json:"type"`
	Props        map[string]interface{} `json:"props"`
	Children     []GeneratedComponent   `json:"children,omitempty"`
	StateID      string                 `json:"state_id,omitempty"`
	AgentID      string                 `json:"agent_id,omitempty"`
	ToolID       string                 `json:"tool_id,omitempty"`
	ActionID     string                 `json:"action_id,omitempty"`
	CreatedAt    int64                  `json:"created_at"`
	UpdatedAt    int64                  `json:"updated_at"`
	Dependencies []string               `json:"dependencies,omitempty"`
}

func NewRSCAdapter(sharedStateManager *sharedstate.SharedStateManager, eventStream *eventstream.Stream, aiClient *VercelAIClient) *RSCAdapter {
	adapter := &RSCAdapter{
		sharedStateManager: sharedStateManager,
		eventStream:        eventStream,
		aiClient:           aiClient,
		componentRegistry:  make(map[ComponentType]ComponentGenerator),
		componentHistory:   make([]GeneratedComponent, 0),
	}

	adapter.registerDefaultComponentGenerators()

	return adapter
}

func (a *RSCAdapter) registerDefaultComponentGenerators() {
	a.RegisterComponentGenerator(ButtonComponent, a.generateButtonComponent)
	a.RegisterComponentGenerator(CardComponent, a.generateCardComponent)
	a.RegisterComponentGenerator(FormComponent, a.generateFormComponent)
	a.RegisterComponentGenerator(TableComponent, a.generateTableComponent)
	a.RegisterComponentGenerator(ChartComponent, a.generateChartComponent)
	a.RegisterComponentGenerator(ProgressComponent, a.generateProgressComponent)
	a.RegisterComponentGenerator(AlertComponent, a.generateAlertComponent)
	a.RegisterComponentGenerator(ModalComponent, a.generateModalComponent)
	a.RegisterComponentGenerator(TabsComponent, a.generateTabsComponent)
	a.RegisterComponentGenerator(AccordionComponent, a.generateAccordionComponent)
	a.RegisterComponentGenerator(TimelineComponent, a.generateTimelineComponent)
	a.RegisterComponentGenerator(CodeBlockComponent, a.generateCodeBlockComponent)
	a.RegisterComponentGenerator(TerminalComponent, a.generateTerminalComponent)
	a.RegisterComponentGenerator(MarkdownComponent, a.generateMarkdownComponent)
	a.RegisterComponentGenerator(ContainerComponent, a.generateContainerComponent)
	a.RegisterComponentGenerator(ToolOutputComponent, a.generateToolOutputComponent)
	a.RegisterComponentGenerator(AgentActionComponent, a.generateAgentActionComponent)
}

func (a *RSCAdapter) RegisterComponentGenerator(componentType ComponentType, generator ComponentGenerator) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.componentRegistry[componentType] = generator
}

func (a *RSCAdapter) GenerateComponent(ctx context.Context, componentType ComponentType, props map[string]interface{}) (GeneratedComponent, error) {
	a.mu.RLock()
	generator, exists := a.componentRegistry[componentType]
	a.mu.RUnlock()

	if !exists {
		return GeneratedComponent{}, fmt.Errorf("component generator for type %s not found", componentType)
	}

	component, err := generator(ctx, props)
	if err != nil {
		return GeneratedComponent{}, fmt.Errorf("error generating component: %w", err)
	}

	a.mu.Lock()
	a.componentHistory = append(a.componentHistory, component)
	a.mu.Unlock()

	a.notifyComponentGenerated(component)

	return component, nil
}

func (a *RSCAdapter) GenerateComponentFromAgentAction(ctx context.Context, agentID string, actionID string, actionType string, actionData map[string]interface{}) (GeneratedComponent, error) {
	componentType, props := a.mapActionToComponent(actionType, actionData)

	props["agentID"] = agentID
	props["actionID"] = actionID
	props["actionType"] = actionType

	return a.GenerateComponent(ctx, componentType, props)
}

func (a *RSCAdapter) GenerateComponentFromToolUsage(ctx context.Context, agentID string, toolID string, toolName string, toolInput map[string]interface{}, toolOutput map[string]interface{}) (GeneratedComponent, error) {
	props := map[string]interface{}{
		"agentID":    agentID,
		"toolID":     toolID,
		"toolName":   toolName,
		"toolInput":  toolInput,
		"toolOutput": toolOutput,
	}

	return a.GenerateComponent(ctx, ToolOutputComponent, props)
}

func (a *RSCAdapter) mapActionToComponent(actionType string, actionData map[string]interface{}) (ComponentType, map[string]interface{}) {
	props := make(map[string]interface{})
	for k, v := range actionData {
		props[k] = v
	}

	switch actionType {
	case "thinking":
		return MarkdownComponent, props
	case "search":
		return TableComponent, props
	case "code_generation":
		return CodeBlockComponent, props
	case "tool_execution":
		return ToolOutputComponent, props
	case "file_operation":
		return CardComponent, props
	case "database_query":
		return TableComponent, props
	case "api_call":
		return CardComponent, props
	case "progress_update":
		return ProgressComponent, props
	case "error":
		return AlertComponent, props
	case "success":
		return AlertComponent, props
	default:
		return CardComponent, props
	}
}

func (a *RSCAdapter) GetComponentHistory() []GeneratedComponent {
	a.mu.RLock()
	defer a.mu.RUnlock()

	history := make([]GeneratedComponent, len(a.componentHistory))
	copy(history, a.componentHistory)

	return history
}

func (a *RSCAdapter) ClearComponentHistory() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.componentHistory = make([]GeneratedComponent, 0)
}

func (a *RSCAdapter) notifyComponentGenerated(component GeneratedComponent) {
	if a.aiClient != nil {
		a.aiClient.SendStreamingEvent(map[string]interface{}{
			"type":      "component_generated",
			"component": component,
		})
	}

	if a.eventStream != nil {
		componentJSON, err := json.Marshal(component)
		if err != nil {
			log.Printf("Error marshaling component: %v", err)
			return
		}

		event := eventstream.Event{
			Type:    "component_generated",
			Payload: string(componentJSON),
		}

		err = a.eventStream.Publish(event)
		if err != nil {
			log.Printf("Error publishing component generated event: %v", err)
		}
	}
}


func (a *RSCAdapter) generateButtonComponent(ctx context.Context, props map[string]interface{}) (GeneratedComponent, error) {
	return GeneratedComponent{
		Type:  ButtonComponent,
		Props: props,
	}, nil
}

func (a *RSCAdapter) generateCardComponent(ctx context.Context, props map[string]interface{}) (GeneratedComponent, error) {
	return GeneratedComponent{
		Type:  CardComponent,
		Props: props,
	}, nil
}

func (a *RSCAdapter) generateFormComponent(ctx context.Context, props map[string]interface{}) (GeneratedComponent, error) {
	return GeneratedComponent{
		Type:  FormComponent,
		Props: props,
	}, nil
}

func (a *RSCAdapter) generateTableComponent(ctx context.Context, props map[string]interface{}) (GeneratedComponent, error) {
	return GeneratedComponent{
		Type:  TableComponent,
		Props: props,
	}, nil
}

func (a *RSCAdapter) generateChartComponent(ctx context.Context, props map[string]interface{}) (GeneratedComponent, error) {
	return GeneratedComponent{
		Type:  ChartComponent,
		Props: props,
	}, nil
}

func (a *RSCAdapter) generateProgressComponent(ctx context.Context, props map[string]interface{}) (GeneratedComponent, error) {
	return GeneratedComponent{
		Type:  ProgressComponent,
		Props: props,
	}, nil
}

func (a *RSCAdapter) generateAlertComponent(ctx context.Context, props map[string]interface{}) (GeneratedComponent, error) {
	return GeneratedComponent{
		Type:  AlertComponent,
		Props: props,
	}, nil
}

func (a *RSCAdapter) generateModalComponent(ctx context.Context, props map[string]interface{}) (GeneratedComponent, error) {
	return GeneratedComponent{
		Type:  ModalComponent,
		Props: props,
	}, nil
}

func (a *RSCAdapter) generateTabsComponent(ctx context.Context, props map[string]interface{}) (GeneratedComponent, error) {
	return GeneratedComponent{
		Type:  TabsComponent,
		Props: props,
	}, nil
}

func (a *RSCAdapter) generateAccordionComponent(ctx context.Context, props map[string]interface{}) (GeneratedComponent, error) {
	return GeneratedComponent{
		Type:  AccordionComponent,
		Props: props,
	}, nil
}

func (a *RSCAdapter) generateTimelineComponent(ctx context.Context, props map[string]interface{}) (GeneratedComponent, error) {
	return GeneratedComponent{
		Type:  TimelineComponent,
		Props: props,
	}, nil
}

func (a *RSCAdapter) generateCodeBlockComponent(ctx context.Context, props map[string]interface{}) (GeneratedComponent, error) {
	return GeneratedComponent{
		Type:  CodeBlockComponent,
		Props: props,
	}, nil
}

func (a *RSCAdapter) generateTerminalComponent(ctx context.Context, props map[string]interface{}) (GeneratedComponent, error) {
	return GeneratedComponent{
		Type:  TerminalComponent,
		Props: props,
	}, nil
}

func (a *RSCAdapter) generateMarkdownComponent(ctx context.Context, props map[string]interface{}) (GeneratedComponent, error) {
	return GeneratedComponent{
		Type:  MarkdownComponent,
		Props: props,
	}, nil
}

func (a *RSCAdapter) generateContainerComponent(ctx context.Context, props map[string]interface{}) (GeneratedComponent, error) {
	return GeneratedComponent{
		Type:  ContainerComponent,
		Props: props,
	}, nil
}

func (a *RSCAdapter) generateToolOutputComponent(ctx context.Context, props map[string]interface{}) (GeneratedComponent, error) {
	return GeneratedComponent{
		Type:  ToolOutputComponent,
		Props: props,
	}, nil
}

func (a *RSCAdapter) generateAgentActionComponent(ctx context.Context, props map[string]interface{}) (GeneratedComponent, error) {
	return GeneratedComponent{
		Type:  AgentActionComponent,
		Props: props,
	}, nil
}

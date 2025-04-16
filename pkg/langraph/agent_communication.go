package langraph

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/spectrumwebco/agent_runtime/pkg/eventstream/models"
)

type CommunicationChannel struct {
	ID            string                 `json:"id"`
	SourceAgentID string                 `json:"source_agent_id"`
	TargetAgentID string                 `json:"target_agent_id"`
	Pattern       string                 `json:"pattern"`
	State         map[string]interface{} `json:"state,omitempty"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

type Message struct {
	ID              string                 `json:"id"`
	ChannelID       string                 `json:"channel_id"`
	SourceAgentID   string                 `json:"source_agent_id"`
	TargetAgentID   string                 `json:"target_agent_id"`
	Type            string                 `json:"type"`
	Content         map[string]interface{} `json:"content"`
	Metadata        map[string]string      `json:"metadata,omitempty"`
	RequiresReply   bool                   `json:"requires_reply"`
	ReplyToID       string                 `json:"reply_to_id,omitempty"`
	Status          string                 `json:"status"`
	ProcessedAt     time.Time              `json:"processed_at,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
}

type CommunicationManager struct {
	Channels       map[string]*CommunicationChannel
	Messages       map[string]*Message
	EventStream    EventStream
	AgentRegistry  map[string]*Agent
	mutex          sync.RWMutex
}

func NewCommunicationManager(eventStream EventStream) *CommunicationManager {
	return &CommunicationManager{
		Channels:      make(map[string]*CommunicationChannel),
		Messages:      make(map[string]*Message),
		EventStream:   eventStream,
		AgentRegistry: make(map[string]*Agent),
	}
}

func (cm *CommunicationManager) RegisterAgent(agent *Agent) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	
	cm.AgentRegistry[agent.Config.ID] = agent
}

func (cm *CommunicationManager) CreateChannel(sourceAgentID, targetAgentID, pattern string) (*CommunicationChannel, error) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	
	if _, exists := cm.AgentRegistry[sourceAgentID]; !exists {
		return nil, fmt.Errorf("source agent %s not found", sourceAgentID)
	}
	
	if _, exists := cm.AgentRegistry[targetAgentID]; !exists {
		return nil, fmt.Errorf("target agent %s not found", targetAgentID)
	}
	
	now := time.Now().UTC()
	channel := &CommunicationChannel{
		ID:            uuid.New().String(),
		SourceAgentID: sourceAgentID,
		TargetAgentID: targetAgentID,
		Pattern:       pattern,
		State:         make(map[string]interface{}),
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	
	cm.Channels[channel.ID] = channel
	
	event := models.NewEvent(
		models.EventTypeSystem,
		models.EventSourceSystem,
		map[string]interface{}{
			"action":          "channel_created",
			"channel_id":      channel.ID,
			"source_agent_id": sourceAgentID,
			"target_agent_id": targetAgentID,
			"pattern":         pattern,
		},
		map[string]string{
			"channel_id": channel.ID,
		},
	)
	
	if cm.EventStream != nil {
		if err := cm.EventStream.AddEvent(event); err != nil {
			fmt.Printf("Failed to add event to stream: %v\n", err)
		}
	}
	
	return channel, nil
}

func (cm *CommunicationManager) SendMessage(ctx context.Context, sourceAgentID, targetAgentID, channelID, messageType string, content map[string]interface{}, metadata map[string]string, requiresReply bool, replyToID string) (*Message, error) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	
	channel, exists := cm.Channels[channelID]
	if !exists {
		return nil, fmt.Errorf("channel %s not found", channelID)
	}
	
	if channel.SourceAgentID != sourceAgentID && channel.TargetAgentID != sourceAgentID {
		return nil, fmt.Errorf("agent %s is not authorized to use channel %s", sourceAgentID, channelID)
	}
	
	message := &Message{
		ID:            uuid.New().String(),
		ChannelID:     channelID,
		SourceAgentID: sourceAgentID,
		TargetAgentID: targetAgentID,
		Type:          messageType,
		Content:       content,
		Metadata:      metadata,
		RequiresReply: requiresReply,
		ReplyToID:     replyToID,
		Status:        "sent",
		CreatedAt:     time.Now().UTC(),
	}
	
	cm.Messages[message.ID] = message
	
	event := models.NewEvent(
		models.EventTypeMessage,
		models.EventSourceAgent,
		map[string]interface{}{
			"message_id":      message.ID,
			"channel_id":      channelID,
			"source_agent_id": sourceAgentID,
			"target_agent_id": targetAgentID,
			"type":            messageType,
			"content":         content,
			"requires_reply":  requiresReply,
			"reply_to_id":     replyToID,
		},
		map[string]string{
			"message_id": message.ID,
			"channel_id": channelID,
		},
	)
	
	if cm.EventStream != nil {
		if err := cm.EventStream.AddEvent(event); err != nil {
			fmt.Printf("Failed to add event to stream: %v\n", err)
		}
	}
	
	go cm.processMessage(ctx, message)
	
	return message, nil
}

func (cm *CommunicationManager) processMessage(ctx context.Context, message *Message) {
	cm.mutex.RLock()
	targetAgent, exists := cm.AgentRegistry[message.TargetAgentID]
	cm.mutex.RUnlock()
	
	if !exists {
		fmt.Printf("Target agent %s not found for message %s\n", message.TargetAgentID, message.ID)
		return
	}
	
	result, err := targetAgent.Process(ctx, map[string]interface{}{
		"message_id":      message.ID,
		"channel_id":      message.ChannelID,
		"source_agent_id": message.SourceAgentID,
		"type":            message.Type,
		"content":         message.Content,
		"requires_reply":  message.RequiresReply,
		"reply_to_id":     message.ReplyToID,
	})
	
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	
	message.Status = "processed"
	message.ProcessedAt = time.Now().UTC()
	
	event := models.NewEvent(
		models.EventTypeResponse,
		models.EventSourceAgent,
		map[string]interface{}{
			"message_id":      message.ID,
			"channel_id":      message.ChannelID,
			"target_agent_id": message.TargetAgentID,
			"result":          result,
			"error":           err != nil,
			"error_msg":       err,
		},
		map[string]string{
			"message_id": message.ID,
			"channel_id": message.ChannelID,
		},
	)
	
	if cm.EventStream != nil {
		if err := cm.EventStream.AddEvent(event); err != nil {
			fmt.Printf("Failed to add event to stream: %v\n", err)
		}
	}
	
	if message.RequiresReply && err == nil {
		replyContent := result
		if replyContent == nil {
			replyContent = make(map[string]interface{})
		}
		
		replyMessage := &Message{
			ID:            uuid.New().String(),
			ChannelID:     message.ChannelID,
			SourceAgentID: message.TargetAgentID,
			TargetAgentID: message.SourceAgentID,
			Type:          "reply",
			Content:       replyContent,
			Metadata:      message.Metadata,
			RequiresReply: false,
			ReplyToID:     message.ID,
			Status:        "sent",
			CreatedAt:     time.Now().UTC(),
		}
		
		cm.Messages[replyMessage.ID] = replyMessage
		
		replyEvent := models.NewEvent(
			models.EventTypeMessage,
			models.EventSourceAgent,
			map[string]interface{}{
				"message_id":      replyMessage.ID,
				"channel_id":      replyMessage.ChannelID,
				"source_agent_id": replyMessage.SourceAgentID,
				"target_agent_id": replyMessage.TargetAgentID,
				"type":            replyMessage.Type,
				"content":         replyMessage.Content,
				"requires_reply":  replyMessage.RequiresReply,
				"reply_to_id":     replyMessage.ReplyToID,
			},
			map[string]string{
				"message_id": replyMessage.ID,
				"channel_id": replyMessage.ChannelID,
			},
		)
		
		if cm.EventStream != nil {
			if err := cm.EventStream.AddEvent(replyEvent); err != nil {
				fmt.Printf("Failed to add event to stream: %v\n", err)
			}
		}
		
		go cm.processMessage(ctx, replyMessage)
	}
}

func (cm *CommunicationManager) GetMessage(messageID string) (*Message, error) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	
	message, exists := cm.Messages[messageID]
	if !exists {
		return nil, fmt.Errorf("message %s not found", messageID)
	}
	
	return message, nil
}

func (cm *CommunicationManager) GetChannel(channelID string) (*CommunicationChannel, error) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	
	channel, exists := cm.Channels[channelID]
	if !exists {
		return nil, fmt.Errorf("channel %s not found", channelID)
	}
	
	return channel, nil
}

func (cm *CommunicationManager) GetChannelsBetweenAgents(agentID1, agentID2 string) ([]*CommunicationChannel, error) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	
	var channels []*CommunicationChannel
	
	for _, channel := range cm.Channels {
		if (channel.SourceAgentID == agentID1 && channel.TargetAgentID == agentID2) ||
		   (channel.SourceAgentID == agentID2 && channel.TargetAgentID == agentID1) {
			channels = append(channels, channel)
		}
	}
	
	return channels, nil
}

func (cm *CommunicationManager) GetMessagesByChannel(channelID string) ([]*Message, error) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	
	var messages []*Message
	
	for _, message := range cm.Messages {
		if message.ChannelID == channelID {
			messages = append(messages, message)
		}
	}
	
	return messages, nil
}

func (cm *CommunicationManager) UpdateChannelState(channelID string, state map[string]interface{}) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	
	channel, exists := cm.Channels[channelID]
	if !exists {
		return fmt.Errorf("channel %s not found", channelID)
	}
	
	for k, v := range state {
		channel.State[k] = v
	}
	
	channel.UpdatedAt = time.Now().UTC()
	
	return nil
}

func (cm *CommunicationManager) CloseChannel(channelID string) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	
	if _, exists := cm.Channels[channelID]; !exists {
		return fmt.Errorf("channel %s not found", channelID)
	}
	
	delete(cm.Channels, channelID)
	
	event := models.NewEvent(
		models.EventTypeSystem,
		models.EventSourceSystem,
		map[string]interface{}{
			"action":     "channel_closed",
			"channel_id": channelID,
		},
		map[string]string{
			"channel_id": channelID,
		},
	)
	
	if cm.EventStream != nil {
		if err := cm.EventStream.AddEvent(event); err != nil {
			fmt.Printf("Failed to add event to stream: %v\n", err)
		}
	}
	
	return nil
}

func (cm *CommunicationManager) CreateStandardCommunicationChannels(ctx context.Context, agents []*Agent) error {
	agentsByRole := make(map[AgentRole]*Agent)
	for _, agent := range agents {
		agentsByRole[agent.Config.Role] = agent
		cm.RegisterAgent(agent)
	}
	
	standardPatterns := []struct {
		SourceRole AgentRole
		TargetRole AgentRole
		Pattern    string
	}{
		{AgentRoleOrchestrator, AgentRoleFrontend, "task_assignment"},
		{AgentRoleOrchestrator, AgentRoleAppBuilder, "task_assignment"},
		{AgentRoleOrchestrator, AgentRoleCodegen, "task_assignment"},
		{AgentRoleOrchestrator, AgentRoleEngineering, "task_assignment"},
		{AgentRoleFrontend, AgentRoleAppBuilder, "api_specification_request"},
		{AgentRoleAppBuilder, AgentRoleFrontend, "api_specification_response"},
		{AgentRoleFrontend, AgentRoleCodegen, "code_generation_request"},
		{AgentRoleCodegen, AgentRoleFrontend, "code_generation_response"},
		{AgentRoleFrontend, AgentRoleEngineering, "integration_support_request"},
		{AgentRoleEngineering, AgentRoleFrontend, "integration_support_response"},
		{AgentRoleAppBuilder, AgentRoleCodegen, "code_generation_request"},
		{AgentRoleCodegen, AgentRoleAppBuilder, "code_generation_response"},
		{AgentRoleAppBuilder, AgentRoleEngineering, "integration_support_request"},
		{AgentRoleEngineering, AgentRoleAppBuilder, "integration_support_response"},
		{AgentRoleCodegen, AgentRoleEngineering, "code_review_request"},
		{AgentRoleEngineering, AgentRoleCodegen, "code_review_response"},
	}
	
	for _, pattern := range standardPatterns {
		sourceAgent, sourceExists := agentsByRole[pattern.SourceRole]
		targetAgent, targetExists := agentsByRole[pattern.TargetRole]
		
		if !sourceExists || !targetExists {
			continue
		}
		
		_, err := cm.CreateChannel(sourceAgent.Config.ID, targetAgent.Config.ID, pattern.Pattern)
		if err != nil {
			return fmt.Errorf("failed to create channel from %s to %s: %v", pattern.SourceRole, pattern.TargetRole, err)
		}
	}
	
	return nil
}

func (cm *CommunicationManager) InitiateTaskExecution(ctx context.Context, task map[string]interface{}) error {
	cm.mutex.RLock()
	var orchestratorAgent *Agent
	for _, agent := range cm.AgentRegistry {
		if agent.Config.Role == AgentRoleOrchestrator {
			orchestratorAgent = agent
			break
		}
	}
	cm.mutex.RUnlock()
	
	if orchestratorAgent == nil {
		return fmt.Errorf("orchestrator agent not found")
	}
	
	_, err := orchestratorAgent.Process(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to process task with orchestrator agent: %v", err)
	}
	
	return nil
}

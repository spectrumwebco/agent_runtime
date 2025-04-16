package models

// Tool represents a tool in the system
type Tool struct {
	Base
	Name        string `gorm:"not null" json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Command     string `json:"command"`
	AgentID     uint   `json:"agent_id"`
	Agent       Agent  `gorm:"foreignKey:AgentID" json:"agent"`
}

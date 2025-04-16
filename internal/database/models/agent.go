package models

// Agent represents an agent in the system
type Agent struct {
	Base
	Name        string `gorm:"not null" json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Status      string `json:"status"`
	UserID      uint   `json:"user_id"`
	User        User   `gorm:"foreignKey:UserID" json:"user"`
}

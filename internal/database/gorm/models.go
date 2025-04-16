package gorm

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

type User struct {
	BaseModel
	Username     string    `gorm:"size:255;not null;uniqueIndex" json:"username"`
	Email        string    `gorm:"size:255;not null;uniqueIndex" json:"email"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	FirstName    string    `gorm:"size:255" json:"first_name"`
	LastName     string    `gorm:"size:255" json:"last_name"`
	LastLogin    time.Time `json:"last_login"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	IsAdmin      bool      `gorm:"default:false" json:"is_admin"`
	Roles        []Role    `gorm:"many2many:user_roles;" json:"roles"`
}

type Role struct {
	BaseModel
	Name        string `gorm:"size:255;not null;uniqueIndex" json:"name"`
	Description string `gorm:"size:255" json:"description"`
	Users       []User `gorm:"many2many:user_roles;" json:"-"`
}

type Permission struct {
	BaseModel
	Name        string `gorm:"size:255;not null;uniqueIndex" json:"name"`
	Description string `gorm:"size:255" json:"description"`
	Roles       []Role `gorm:"many2many:role_permissions;" json:"-"`
}

type Agent struct {
	BaseModel
	Name        string `gorm:"size:255;not null" json:"name"`
	Description string `gorm:"size:255" json:"description"`
	Type        string `gorm:"size:50;not null" json:"type"`
	Config      string `gorm:"type:text" json:"config"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
	UserID      uint   `json:"user_id"`
	User        User   `gorm:"foreignKey:UserID" json:"-"`
	Tasks       []Task `json:"tasks"`
}

type Task struct {
	BaseModel
	Title       string    `gorm:"size:255;not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	Status      string    `gorm:"size:50;not null" json:"status"`
	Priority    string    `gorm:"size:50;not null" json:"priority"`
	DueDate     time.Time `json:"due_date"`
	CompletedAt time.Time `json:"completed_at"`
	AgentID     uint      `json:"agent_id"`
	Agent       Agent     `gorm:"foreignKey:AgentID" json:"-"`
	UserID      uint      `json:"user_id"`
	User        User      `gorm:"foreignKey:UserID" json:"-"`
}

type Tool struct {
	BaseModel
	Name        string `gorm:"size:255;not null;uniqueIndex" json:"name"`
	Description string `gorm:"size:255" json:"description"`
	Type        string `gorm:"size:50;not null" json:"type"`
	Config      string `gorm:"type:text" json:"config"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
}

type AgentTool struct {
	AgentID   uint      `gorm:"primaryKey" json:"agent_id"`
	ToolID    uint      `gorm:"primaryKey" json:"tool_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Workspace struct {
	BaseModel
	Name        string `gorm:"size:255;not null" json:"name"`
	Description string `gorm:"size:255" json:"description"`
	Path        string `gorm:"size:255;not null" json:"path"`
	UserID      uint   `json:"user_id"`
	User        User   `gorm:"foreignKey:UserID" json:"-"`
}

type Session struct {
	BaseModel
	Token      string    `gorm:"size:255;not null;uniqueIndex" json:"token"`
	ExpiresAt  time.Time `json:"expires_at"`
	UserID     uint      `json:"user_id"`
	User       User      `gorm:"foreignKey:UserID" json:"-"`
	UserAgent  string    `gorm:"size:255" json:"user_agent"`
	IPAddress  string    `gorm:"size:50" json:"ip_address"`
	LastActive time.Time `json:"last_active"`
}

type AuditLog struct {
	BaseModel
	Action     string `gorm:"size:50;not null" json:"action"`
	EntityType string `gorm:"size:50;not null" json:"entity_type"`
	EntityID   uint   `json:"entity_id"`
	UserID     uint   `json:"user_id"`
	User       User   `gorm:"foreignKey:UserID" json:"-"`
	IPAddress  string `gorm:"size:50" json:"ip_address"`
	Details    string `gorm:"type:text" json:"details"`
}

type Setting struct {
	BaseModel
	Key         string `gorm:"size:255;not null;uniqueIndex" json:"key"`
	Value       string `gorm:"type:text" json:"value"`
	Description string `gorm:"size:255" json:"description"`
	IsPublic    bool   `gorm:"default:false" json:"is_public"`
}

type APIKey struct {
	BaseModel
	Key         string    `gorm:"size:255;not null;uniqueIndex" json:"key"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	Description string    `gorm:"size:255" json:"description"`
	ExpiresAt   time.Time `json:"expires_at"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	UserID      uint      `json:"user_id"`
	User        User      `gorm:"foreignKey:UserID" json:"-"`
	LastUsed    time.Time `json:"last_used"`
}

type Notification struct {
	BaseModel
	Title     string `gorm:"size:255;not null" json:"title"`
	Message   string `gorm:"type:text" json:"message"`
	Type      string `gorm:"size:50;not null" json:"type"`
	IsRead    bool   `gorm:"default:false" json:"is_read"`
	UserID    uint   `json:"user_id"`
	User      User   `gorm:"foreignKey:UserID" json:"-"`
	ReadAt    time.Time `json:"read_at"`
	EntityType string `gorm:"size:50" json:"entity_type"`
	EntityID   uint   `json:"entity_id"`
}

func GetModels() []interface{} {
	return []interface{}{
		&User{},
		&Role{},
		&Permission{},
		&Agent{},
		&Task{},
		&Tool{},
		&AgentTool{},
		&Workspace{},
		&Session{},
		&AuditLog{},
		&Setting{},
		&APIKey{},
		&Notification{},
	}
}

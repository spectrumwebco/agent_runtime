package models

// User represents a user in the system
type User struct {
	Base
	Username  string `gorm:"uniqueIndex;not null" json:"username"`
	Email     string `gorm:"uniqueIndex;not null" json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
}

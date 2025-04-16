package repositories

import (
	"context"

	"gorm.io/gorm"
)

// Repository is the interface that all repositories should implement
type Repository interface {
	Find(ctx context.Context, id uint) (interface{}, error)
	FindAll(ctx context.Context) (interface{}, error)
	Create(ctx context.Context, entity interface{}) error
	Update(ctx context.Context, entity interface{}) error
	Delete(ctx context.Context, id uint) error
}

// BaseRepository is the base implementation of Repository
type BaseRepository struct {
	DB *gorm.DB
}

// NewBaseRepository creates a new BaseRepository
func NewBaseRepository(db *gorm.DB) *BaseRepository {
	return &BaseRepository{
		DB: db,
	}
}

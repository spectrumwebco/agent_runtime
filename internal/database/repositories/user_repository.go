package repositories

import (
	"errors"

	"github.com/spectrumwebco/agent_runtime/internal/database"
	"github.com/spectrumwebco/agent_runtime/internal/database/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *database.DB
}

func NewUserRepository(db *database.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) List(offset, limit int) ([]models.User, error) {
	var users []models.User
	err := r.db.Offset(offset).Limit(limit).Find(&users).Error
	return users, err
}

func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

func (r *UserRepository) GetUserWorkspaces(userID uint) ([]models.Workspace, error) {
	var user models.User
	err := r.db.Preload("Workspaces").First(&user, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return user.Workspaces, nil
}

func (r *UserRepository) AddUserToWorkspace(userID, workspaceID uint) error {
	return r.db.Exec("INSERT INTO user_workspaces (user_id, workspace_id) VALUES (?, ?)", userID, workspaceID).Error
}

func (r *UserRepository) RemoveUserFromWorkspace(userID, workspaceID uint) error {
	return r.db.Exec("DELETE FROM user_workspaces WHERE user_id = ? AND workspace_id = ?", userID, workspaceID).Error
}

func (r *UserRepository) CreateApiKey(apiKey *models.ApiKey) error {
	return r.db.Create(apiKey).Error
}

func (r *UserRepository) GetApiKeyByKey(key string) (*models.ApiKey, error) {
	var apiKey models.ApiKey
	err := r.db.Where("key = ? AND is_active = ?", key, true).First(&apiKey).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &apiKey, nil
}

func (r *UserRepository) GetUserApiKeys(userID uint) ([]models.ApiKey, error) {
	var apiKeys []models.ApiKey
	err := r.db.Where("user_id = ?", userID).Find(&apiKeys).Error
	return apiKeys, err
}

func (r *UserRepository) DeleteApiKey(id uint) error {
	return r.db.Delete(&models.ApiKey{}, id).Error
}

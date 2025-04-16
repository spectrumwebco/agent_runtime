package jwk

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"time"

	"github.com/go-jose/go-jose/v3"
)

type Manager interface {
	GetSet(ctx context.Context, setID string) (*jose.JSONWebKeySet, error)
	
	CreateSet(ctx context.Context, setID, keyType, use, algorithm string) (*jose.JSONWebKeySet, error)
	
	UpdateSet(ctx context.Context, setID string, keys []jose.JSONWebKey) (*jose.JSONWebKeySet, error)
	
	DeleteSet(ctx context.Context, setID string) error
	
	ListSets(ctx context.Context) ([]string, error)
	
	GetKey(ctx context.Context, setID, keyID string) (*jose.JSONWebKey, error)
	
	UpdateKey(ctx context.Context, setID, keyID string, key jose.JSONWebKey) (*jose.JSONWebKey, error)
	
	DeleteKey(ctx context.Context, setID, keyID string) error
	
	GenerateKey(ctx context.Context, setID, keyType, use, algorithm string) (*jose.JSONWebKey, error)
}

type DefaultManager struct {
	storage Storage
}

type Storage interface {
	GetSet(ctx context.Context, setID string) (*jose.JSONWebKeySet, error)
	
	CreateSet(ctx context.Context, setID string, set *jose.JSONWebKeySet) error
	
	UpdateSet(ctx context.Context, setID string, set *jose.JSONWebKeySet) error
	
	DeleteSet(ctx context.Context, setID string) error
	
	ListSets(ctx context.Context) ([]string, error)
}

func NewManager(storage Storage) Manager {
	return &DefaultManager{
		storage: storage,
	}
}

func (m *DefaultManager) GetSet(ctx context.Context, setID string) (*jose.JSONWebKeySet, error) {
	return m.storage.GetSet(ctx, setID)
}

func (m *DefaultManager) CreateSet(ctx context.Context, setID, keyType, use, algorithm string) (*jose.JSONWebKeySet, error) {
	key, err := m.GenerateKey(ctx, setID, keyType, use, algorithm)
	if err != nil {
		return nil, err
	}
	
	set := &jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{*key},
	}
	
	if err := m.storage.CreateSet(ctx, setID, set); err != nil {
		return nil, err
	}
	
	return set, nil
}

func (m *DefaultManager) UpdateSet(ctx context.Context, setID string, keys []jose.JSONWebKey) (*jose.JSONWebKeySet, error) {
	set := &jose.JSONWebKeySet{
		Keys: keys,
	}
	
	if err := m.storage.UpdateSet(ctx, setID, set); err != nil {
		return nil, err
	}
	
	return set, nil
}

func (m *DefaultManager) DeleteSet(ctx context.Context, setID string) error {
	return m.storage.DeleteSet(ctx, setID)
}

func (m *DefaultManager) ListSets(ctx context.Context) ([]string, error) {
	return m.storage.ListSets(ctx)
}

func (m *DefaultManager) GetKey(ctx context.Context, setID, keyID string) (*jose.JSONWebKey, error) {
	set, err := m.storage.GetSet(ctx, setID)
	if err != nil {
		return nil, err
	}
	
	for _, key := range set.Keys {
		if key.KeyID == keyID {
			return &key, nil
		}
	}
	
	return nil, errors.New("key not found")
}

func (m *DefaultManager) UpdateKey(ctx context.Context, setID, keyID string, key jose.JSONWebKey) (*jose.JSONWebKey, error) {
	set, err := m.storage.GetSet(ctx, setID)
	if err != nil {
		return nil, err
	}
	
	found := false
	for i, k := range set.Keys {
		if k.KeyID == keyID {
			set.Keys[i] = key
			found = true
			break
		}
	}
	
	if !found {
		return nil, errors.New("key not found")
	}
	
	if err := m.storage.UpdateSet(ctx, setID, set); err != nil {
		return nil, err
	}
	
	return &key, nil
}

func (m *DefaultManager) DeleteKey(ctx context.Context, setID, keyID string) error {
	set, err := m.storage.GetSet(ctx, setID)
	if err != nil {
		return err
	}
	
	found := false
	newKeys := make([]jose.JSONWebKey, 0, len(set.Keys))
	for _, k := range set.Keys {
		if k.KeyID == keyID {
			found = true
			continue
		}
		newKeys = append(newKeys, k)
	}
	
	if !found {
		return errors.New("key not found")
	}
	
	set.Keys = newKeys
	
	return m.storage.UpdateSet(ctx, setID, set)
}

func (m *DefaultManager) GenerateKey(ctx context.Context, setID, keyType, use, algorithm string) (*jose.JSONWebKey, error) {
	var key jose.JSONWebKey
	
	switch keyType {
	case "RSA":
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return nil, err
		}
		key = jose.JSONWebKey{
			Key:       privateKey,
			KeyID:     setID,
			Algorithm: algorithm,
			Use:       use,
		}
	default:
		return nil, errors.New("unsupported key type")
	}
	
	return &key, nil
}

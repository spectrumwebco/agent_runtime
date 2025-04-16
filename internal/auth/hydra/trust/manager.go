package trust

import (
	"context"
	"errors"
	"time"
)

type TrustedIssuer struct {
	ID          string    `json:"id"`
	Issuer      string    `json:"issuer"`
	PublicKey   string    `json:"public_key"`
	AllowedAud  []string  `json:"allowed_audience"`
	Scope       string    `json:"scope"`
	ExpiresAt   time.Time `json:"expires_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Manager interface {
	GetTrustedIssuer(ctx context.Context, id string) (*TrustedIssuer, error)
	
	CreateTrustedIssuer(ctx context.Context, issuer *TrustedIssuer) error
	
	UpdateTrustedIssuer(ctx context.Context, issuer *TrustedIssuer) error
	
	DeleteTrustedIssuer(ctx context.Context, id string) error
	
	ListTrustedIssuers(ctx context.Context) ([]*TrustedIssuer, error)
	
	ValidateJWT(ctx context.Context, token string) (*JWTValidationResult, error)
}

type JWTValidationResult struct {
	Valid       bool      `json:"valid"`
	Issuer      string    `json:"issuer"`
	Subject     string    `json:"subject"`
	Audience    []string  `json:"audience"`
	Scope       string    `json:"scope"`
	ExpiresAt   time.Time `json:"expires_at"`
	IssuedAt    time.Time `json:"issued_at"`
	NotBefore   time.Time `json:"not_before"`
	Error       string    `json:"error,omitempty"`
}

type DefaultManager struct {
	storage Storage
}

type Storage interface {
	GetTrustedIssuer(ctx context.Context, id string) (*TrustedIssuer, error)
	
	CreateTrustedIssuer(ctx context.Context, issuer *TrustedIssuer) error
	
	UpdateTrustedIssuer(ctx context.Context, issuer *TrustedIssuer) error
	
	DeleteTrustedIssuer(ctx context.Context, id string) error
	
	ListTrustedIssuers(ctx context.Context) ([]*TrustedIssuer, error)
}

func NewManager(storage Storage) Manager {
	return &DefaultManager{
		storage: storage,
	}
}

func (m *DefaultManager) GetTrustedIssuer(ctx context.Context, id string) (*TrustedIssuer, error) {
	return m.storage.GetTrustedIssuer(ctx, id)
}

func (m *DefaultManager) CreateTrustedIssuer(ctx context.Context, issuer *TrustedIssuer) error {
	issuer.CreatedAt = time.Now()
	issuer.UpdatedAt = time.Now()
	
	return m.storage.CreateTrustedIssuer(ctx, issuer)
}

func (m *DefaultManager) UpdateTrustedIssuer(ctx context.Context, issuer *TrustedIssuer) error {
	issuer.UpdatedAt = time.Now()
	
	return m.storage.UpdateTrustedIssuer(ctx, issuer)
}

func (m *DefaultManager) DeleteTrustedIssuer(ctx context.Context, id string) error {
	return m.storage.DeleteTrustedIssuer(ctx, id)
}

func (m *DefaultManager) ListTrustedIssuers(ctx context.Context) ([]*TrustedIssuer, error) {
	return m.storage.ListTrustedIssuers(ctx)
}

func (m *DefaultManager) ValidateJWT(ctx context.Context, token string) (*JWTValidationResult, error) {
	
	return nil, errors.New("not implemented")
}

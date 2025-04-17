package interfaces

import (
	"context"
)

type Terminal interface {
	ID() string

	Execute(ctx context.Context, command string) (string, error)

	Start(ctx context.Context) error

	Stop(ctx context.Context) error

	IsRunning() bool

	GetType() string
}

type TerminalManager interface {
	CreateTerminal(ctx context.Context, terminalType string, id string, options map[string]interface{}) (Terminal, error)

	GetTerminal(id string) (Terminal, error)

	ListTerminals() []Terminal

	RemoveTerminal(ctx context.Context, id string) error

	CreateBulkTerminals(ctx context.Context, terminalType string, count int, options map[string]interface{}) ([]Terminal, error)
}

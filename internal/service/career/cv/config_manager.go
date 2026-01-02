package cv

import (
	"context"
	career "github.com/baphled/kariya/internal/domain/career"
)

// ConfigManager defines the interface for managing CV configurations.
type ConfigManager interface {
	// LoadConfig loads a CV configuration by name.
	// Returns ErrConfigNotFound if the configuration doesn't exist.
	LoadConfig(ctx context.Context, name string) (*career.CVConfig, error)

	// SaveConfig saves or updates a CV configuration.
	// Creates the directory if it doesn't exist.
	SaveConfig(ctx context.Context, config *career.CVConfig) error

	// DeleteConfig deletes a CV configuration by name.
	// Returns ErrConfigNotFound if the configuration doesn't exist.
	DeleteConfig(ctx context.Context, name string) error

	// ListConfigs returns all available CV configurations.
	ListConfigs(ctx context.Context) ([]*career.CVConfig, error)

	// GetConfigPath returns the file path for a configuration name.
	GetConfigPath(name string) string

	// ConfigExists checks if a configuration exists.
	ConfigExists(ctx context.Context, name string) (bool, error)
}


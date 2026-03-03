package cv

import (
	"context"
	"errors"
	"sync"
	"time"

	career "github.com/baphled/kariya/internal/domain/career"
)

var (
	// ErrConfigNotFound is returned when a configuration is not found.
	ErrConfigNotFound = errors.New("configuration not found")

	// ErrInvalidConfigName is returned when configuration name is invalid.
	ErrInvalidConfigName = errors.New("invalid configuration name")
)

// ConfigManager defines the interface for managing CV configurations.
//
//nolint:interfacebloat // Config management requires load, save, delete, list, and validation methods
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

// MemoryConfigManager provides an in-memory implementation of ConfigManager for testing.
type MemoryConfigManager struct {
	configs map[string]*career.CVConfig
	mu      sync.RWMutex
}

// NewMemoryConfigManager creates a new in-memory configuration manager.
//
// Returns:
//   - A fully initialized MemoryConfigManager ready for use.
//
// Side effects:
//   - None.
func NewMemoryConfigManager() *MemoryConfigManager {
	return &MemoryConfigManager{
		configs: make(map[string]*career.CVConfig),
	}
}

// LoadConfig loads a configuration by name.
//
// Expected:
//   - ctx: valid context for cancellation
//   - name: configuration name to load
//
// Returns:
//   - *career.CVConfig: loaded configuration, nil if not found
//   - error: ErrConfigNotFound if config does not exist, nil on success
//
// Side effects:
//   - None.
func (m *MemoryConfigManager) LoadConfig(_ context.Context, name string) (*career.CVConfig, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	config, exists := m.configs[name]
	if !exists {
		return nil, ErrConfigNotFound
	}

	// Return a copy to prevent external modification
	configCopy := *config
	return &configCopy, nil
}

// SaveConfig saves or updates a configuration.
//
// Expected:
//   - config must be a valid configuration object.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (m *MemoryConfigManager) SaveConfig(_ context.Context, config *career.CVConfig) error {
	if config == nil {
		return ErrInvalidConfigName
	}

	if config.Name == "" {
		return ErrInvalidConfigName
	}

	if err := config.Validate(); err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Update timestamps
	now := time.Now()
	if config.CreatedAt.IsZero() {
		config.CreatedAt = now
	}
	config.UpdatedAt = now

	m.configs[config.Name] = config
	return nil
}

// DeleteConfig deletes a configuration by name.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (m *MemoryConfigManager) DeleteConfig(_ context.Context, name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.configs[name]; !exists {
		return ErrConfigNotFound
	}

	delete(m.configs, name)
	return nil
}

// ListConfigs returns all configurations.
//
// Expected:
//   - context: A context for cancellation and timeouts.
//
// Returns:
//   - []*career.CVConfig: A slice of all stored configurations.
//   - error: An error if the operation fails.
//
// Side effects:
//   - None.
func (m *MemoryConfigManager) ListConfigs(_ context.Context) ([]*career.CVConfig, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	configs := make([]*career.CVConfig, 0, len(m.configs))
	for _, config := range m.configs {
		configCopy := *config
		configs = append(configs, &configCopy)
	}

	return configs, nil
}

// GetConfigPath returns the file path for a configuration name.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *MemoryConfigManager) GetConfigPath(name string) string {
	return "memory://" + name
}

// ConfigExists checks if a configuration exists.
//
// Expected:
//   - context: A context for cancellation and timeouts.
//   - name: The name of the configuration to check.
//
// Returns:
//   - bool: True if the configuration exists, false otherwise.
//   - error: An error if the operation fails.
//
// Side effects:
//   - None.
func (m *MemoryConfigManager) ConfigExists(_ context.Context, name string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, exists := m.configs[name]
	return exists, nil
}

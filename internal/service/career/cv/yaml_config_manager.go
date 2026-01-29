package cv

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	career "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	"gopkg.in/yaml.v3"
)

var (
	// ErrInvalidConfigPath is returned when the config path is invalid.
	ErrInvalidConfigPath = errors.New("invalid configuration path")

	// ErrConfigDirectoryFailed is returned when directory operations fail.
	ErrConfigDirectoryFailed = errors.New("failed to create or access configuration directory")
)

// YAMLConfigManager manages CV configurations stored as YAML files.
type YAMLConfigManager struct {
	configDir string
	logger    *logger.Logger
}

// NewYAMLConfigManager creates a new YAMLConfigManager instance.
// It initializes the configuration directory if it doesn't exist.
func NewYAMLConfigManager(log *logger.Logger) (*YAMLConfigManager, error) {
	// Get home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".kariya", "cv_configs")

	// Create directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0o750); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	return &YAMLConfigManager{
		configDir: configDir,
		logger:    log,
	}, nil
}

// LoadConfig loads a CV configuration from a YAML file.
func (m *YAMLConfigManager) LoadConfig(ctx context.Context, name string) (*career.CVConfig, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if name = strings.TrimSpace(name); name == "" {
		return nil, fmt.Errorf("config name cannot be empty: %w", ErrInvalidConfigPath)
	}

	configPath := m.GetConfigPath(name)

	// Verify path is within configDir (security check)
	if !isPathWithinDirectory(configPath, m.configDir) {
		return nil, fmt.Errorf("config path outside allowed directory: %w", ErrInvalidConfigPath)
	}

	// Read file
	// #nosec G304 -- configPath is validated above to be within configDir
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			m.logger.Info("Config not found: %s at path %s", name, configPath)
			return nil, ErrConfigNotFound
		}
		m.logger.Error("Failed to read config file: %v at path %s", err, configPath)
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Unmarshal YAML
	config := &career.CVConfig{}
	if err := yaml.Unmarshal(data, config); err != nil {
		m.logger.Error("Failed to unmarshal config YAML: %v for config %s", err, name)
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate config
	if err := config.Validate(); err != nil {
		m.logger.Warn("Config validation failed: %v for config %s", err, name)
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	m.logger.Info("Config loaded successfully: %s", name)
	return config, nil
}

// SaveConfig saves a CV configuration to a YAML file.
// Creates the directory if needed and uses atomic write (write to temp, then rename).
func (m *YAMLConfigManager) SaveConfig(ctx context.Context, config *career.CVConfig) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if config == nil {
		return errors.New("config cannot be nil")
	}

	// Validate config before saving
	if err := config.Validate(); err != nil {
		m.logger.Warn("Config validation failed before save: %v for config %s", err, config.Name)
		return fmt.Errorf("config validation failed: %w", err)
	}

	configPath := m.GetConfigPath(config.Name)

	// Verify path is within configDir (security check)
	if !isPathWithinDirectory(configPath, m.configDir) {
		return fmt.Errorf("config path outside allowed directory: %w", ErrInvalidConfigPath)
	}

	// Ensure directory exists
	if err := os.MkdirAll(m.configDir, 0o750); err != nil {
		m.logger.Error("Failed to create config directory: %v at path %s", err, m.configDir)
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Update timestamps
	now := time.Now()
	config.UpdatedAt = now
	if config.CreatedAt.IsZero() {
		config.CreatedAt = now
	}

	// Marshal YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		m.logger.Error("Failed to marshal config to YAML: %v for config %s", err, config.Name)
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to temporary file
	tmpFile, err := os.CreateTemp(m.configDir, ".tmp-*.yaml")
	if err != nil {
		m.logger.Error("Failed to create temp file: %v", err)
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()

	// Write data to temp file
	if _, err := tmpFile.Write(data); err != nil {
		if closeErr := tmpFile.Close(); closeErr != nil {
			m.logger.Error("Failed to close temp file: %v", closeErr)
		}
		if removeErr := os.Remove(tmpPath); removeErr != nil {
			m.logger.Error("Failed to remove temp file: %v", removeErr)
		}
		m.logger.Error("Failed to write to temp file: %v", err)
		return fmt.Errorf("failed to write to temp file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		m.logger.Error("Failed to close temp file: %v", err)
	}

	// Atomic rename
	if err := os.Rename(tmpPath, configPath); err != nil {
		if removeErr := os.Remove(tmpPath); removeErr != nil {
			m.logger.Error("Failed to remove temp file: %v", removeErr)
		}
		m.logger.Error("Failed to rename temp file: %v from %s to %s", err, tmpPath, configPath)
		return fmt.Errorf("failed to save config file: %w", err)
	}

	m.logger.Info("Config saved successfully: %s at path %s", config.Name, configPath)
	return nil
}

// DeleteConfig deletes a CV configuration file.
func (m *YAMLConfigManager) DeleteConfig(ctx context.Context, name string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if name = strings.TrimSpace(name); name == "" {
		return fmt.Errorf("config name cannot be empty: %w", ErrInvalidConfigPath)
	}

	configPath := m.GetConfigPath(name)

	// Verify path is within configDir (security check)
	if !isPathWithinDirectory(configPath, m.configDir) {
		return fmt.Errorf("config path outside allowed directory: %w", ErrInvalidConfigPath)
	}

	// Check if file exists
	if _, err := os.Stat(configPath); err != nil {
		if os.IsNotExist(err) {
			m.logger.Info("Config not found for deletion: %s", name)
			return ErrConfigNotFound
		}
		m.logger.Error("Failed to stat config file: %v", err)
		return err
	}

	// Delete file
	if err := os.Remove(configPath); err != nil {
		m.logger.Error("Failed to delete config file: %v at path %s", err, configPath)
		return fmt.Errorf("failed to delete config file: %w", err)
	}

	m.logger.Info("Config deleted successfully: %s", name)
	return nil
}

// ListConfigs returns all available CV configurations.
func (m *YAMLConfigManager) ListConfigs(ctx context.Context) ([]*career.CVConfig, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// Read directory
	entries, err := os.ReadDir(m.configDir)
	if err != nil {
		if os.IsNotExist(err) {
			m.logger.Info("Config directory not found")
			return []*career.CVConfig{}, nil
		}
		m.logger.Error("Failed to read config directory: %v", err)
		return nil, fmt.Errorf("failed to read config directory: %w", err)
	}

	var configs []*career.CVConfig

	// Iterate through files
	for _, entry := range entries {
		// Skip directories and non-YAML files
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") && !strings.HasSuffix(entry.Name(), ".yml") {
			continue
		}

		// Extract config name (remove extension)
		configName := strings.TrimSuffix(entry.Name(), ".yaml")
		configName = strings.TrimSuffix(configName, ".yml")

		// Load config
		config, err := m.LoadConfig(ctx, configName)
		if err != nil {
			m.logger.Warn("Failed to load config during list: %v for config %s", err, configName)
			continue
		}

		configs = append(configs, config)
	}

	return configs, nil
}

// GetConfigPath returns the file path for a configuration name.
func (m *YAMLConfigManager) GetConfigPath(name string) string {
	// Sanitize name to prevent path traversal
	name = strings.TrimSpace(name)
	name = sanitizeFileName(name)
	return filepath.Join(m.configDir, name+".yaml")
}

// ConfigExists checks if a configuration exists.
func (m *YAMLConfigManager) ConfigExists(ctx context.Context, name string) (bool, error) {
	if ctx.Err() != nil {
		return false, ctx.Err()
	}

	if name = strings.TrimSpace(name); name == "" {
		return false, fmt.Errorf("config name cannot be empty: %w", ErrInvalidConfigPath)
	}

	configPath := m.GetConfigPath(name)

	_, err := os.Stat(configPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// isPathWithinDirectory checks if a path is within a directory (security check).
func isPathWithinDirectory(filePath, dirPath string) bool {
	absFile, err := filepath.Abs(filePath)
	if err != nil {
		return false
	}

	absDir, err := filepath.Abs(dirPath)
	if err != nil {
		return false
	}

	// Ensure directory path ends with separator for proper prefix check
	if !strings.HasSuffix(absDir, string(filepath.Separator)) {
		absDir += string(filepath.Separator)
	}

	return strings.HasPrefix(absFile, absDir)
}

// sanitizeFileName removes potentially dangerous characters from a filename.
func sanitizeFileName(name string) string {
	// Replace path separators and other dangerous characters
	replacer := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
		"\x00", "",
	)
	return replacer.Replace(name)
}

// GetConfigDirectory returns the configuration directory path.
func (m *YAMLConfigManager) GetConfigDirectory() string {
	return m.configDir
}

// VerifyDirectory checks if the config directory exists and is writable.
func (m *YAMLConfigManager) VerifyDirectory() error {
	// Check if directory exists
	info, err := os.Stat(m.configDir)
	if err != nil {
		if os.IsNotExist(err) {
			if m.logger != nil {
				m.logger.Warn("Config directory does not exist: %s, attempting to create", m.configDir)
			}
			// Try to create it
			if err := os.MkdirAll(m.configDir, 0o750); err != nil {
				if m.logger != nil {
					m.logger.Error("Failed to create config directory: %v", err)
				}
				return fmt.Errorf("failed to create config directory: %w", err)
			}
			return nil
		}
		if m.logger != nil {
			m.logger.Error("Failed to stat config directory: %v", err)
		}
		return fmt.Errorf("failed to stat config directory: %w", err)
	}

	// Check if it's a directory
	if !info.IsDir() {
		if m.logger != nil {
			m.logger.Error("Config path exists but is not a directory: %s", m.configDir)
		}
		return fmt.Errorf("config path exists but is not a directory: %s", m.configDir)
	}

	// Try to write a test file to verify permissions
	testFile := filepath.Join(m.configDir, ".write-test")
	if err := os.WriteFile(testFile, []byte("test"), 0o600); err != nil {
		if m.logger != nil {
			m.logger.Error("Config directory is not writable: %v", err)
		}
		return fmt.Errorf("config directory is not writable: %w", err)
	}
	// Clean up test file
	if err := os.Remove(testFile); err != nil && m.logger != nil {
		m.logger.Error("Failed to remove write test file: %v", err)
	}

	if m.logger != nil {
		m.logger.Info("Config directory verified: %s", m.configDir)
	}
	return nil
}

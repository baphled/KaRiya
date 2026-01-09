package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	System  SystemConfig  `yaml:"system"`
	Profile ProfileConfig `yaml:"profile"`
	CV      CVConfig      `yaml:"cv"`
	Export  ExportConfig  `yaml:"export"`
	Display DisplayConfig `yaml:"display"`
}

// SystemConfig contains system-level configuration
type SystemConfig struct {
	DataDir     string `yaml:"data_dir"`
	LogLevel    string `yaml:"log_level"`
	AutoBackup  bool   `yaml:"auto_backup"`
	BackupCount int    `yaml:"backup_count"`
}

// ProfileConfig contains user profile configuration
type ProfileConfig struct {
	Name            string `yaml:"name"`
	Email           string `yaml:"email"`
	DefaultRole     string `yaml:"default_role"`
	DefaultAudience string `yaml:"default_audience"`
}

// CVConfig contains CV generation configuration
type CVConfig struct {
	DefaultFormat string `yaml:"default_format"`
	MaxBullets    int    `yaml:"max_bullets"`
}

// ExportConfig contains export configuration
type ExportConfig struct {
	DefaultDestination string `yaml:"default_destination"`
	AutoOpen           bool   `yaml:"auto_open"`
}

// DisplayConfig contains display/UI configuration
type DisplayConfig struct {
	Theme      string `yaml:"theme"`
	Animations bool   `yaml:"animations"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	dataDir := filepath.Join(homeDir, ".kariya")

	return &Config{
		System: SystemConfig{
			DataDir:     dataDir,
			LogLevel:    "info",
			AutoBackup:  true,
			BackupCount: 5,
		},
		Profile: ProfileConfig{
			Name:            "",
			Email:           "",
			DefaultRole:     "senior_ic",
			DefaultAudience: "technical",
		},
		CV: CVConfig{
			DefaultFormat: "markdown",
			MaxBullets:    50,
		},
		Export: ExportConfig{
			DefaultDestination: "file",
			AutoOpen:           false,
		},
		Display: DisplayConfig{
			Theme:      "dark",
			Animations: true,
		},
	}
}

// GetConfigPath returns the path to the config file
func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	return filepath.Join(homeDir, ".kariya", "config.yaml"), nil
}

// LoadConfig loads configuration from the default location
func LoadConfig() (*Config, error) {
	path, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	return LoadConfigFromPath(path)
}

// LoadConfigFromPath loads configuration from a specific file path
func LoadConfigFromPath(path string) (*Config, error) {
	// Check if config file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Return default config if file doesn't exist
		return DefaultConfig(), nil
	}

	data, err := os.ReadFile(path) // #nosec G304 -- path from GetConfigPath (user home directory)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}

// SaveConfig saves configuration to the default location
func SaveConfig(cfg *Config) error {
	path, err := GetConfigPath()
	if err != nil {
		return err
	}

	return SaveConfigToPath(cfg, path)
}

// SaveConfigToPath saves configuration to a specific file path
func SaveConfigToPath(cfg *Config, path string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

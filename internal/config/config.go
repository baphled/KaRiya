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

	// Narrative CV profile fields
	Title         string   `yaml:"title"`          // e.g., "Senior Software Engineer / Technical Consultant"
	Location      string   `yaml:"location"`       // e.g., "Remote (UK)"
	GitHub        string   `yaml:"github"`         // e.g., "https://github.com/username"
	Portfolio     string   `yaml:"portfolio"`      // e.g., "http://portfolio.example.com"
	CoreStrengths []string `yaml:"core_strengths"` // List of core strengths
	Languages     string   `yaml:"languages"`      // e.g., "Ruby, Go, PHP, C/C++, JavaScript, Shell"
	Frontend      string   `yaml:"frontend"`       // e.g., "Vue.js, React"
	Systems       string   `yaml:"systems"`        // e.g., "Linux, SQL, APIs, CI/CD, automation"
	WhatIBring    []string `yaml:"what_i_bring"`   // List of value propositions
}

// CVConfig contains CV generation configuration
type CVConfig struct {
	DefaultFormat   string                `yaml:"default_format"`
	MaxBullets      int                   `yaml:"max_bullets"`
	AudienceBullets AudienceBulletsConfig `yaml:"audience_bullets"`
}

// AudienceBulletsConfig defines bullets per company for each audience type
type AudienceBulletsConfig struct {
	Recruiter     int `yaml:"recruiter"`      // Bullets per company for recruiters (default: 4)
	HiringManager int `yaml:"hiring_manager"` // Bullets per company for hiring managers (default: 6)
	Peer          int `yaml:"peer"`           // Bullets per company for peers (default: 8)
	Default       int `yaml:"default"`        // Default bullets per company (default: 5)
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
			Title:           "",
			Location:        "",
			GitHub:          "",
			Portfolio:       "",
			CoreStrengths:   []string{},
			Languages:       "",
			Frontend:        "",
			Systems:         "",
			WhatIBring:      []string{},
		},
		CV: CVConfig{
			DefaultFormat: "markdown",
			MaxBullets:    50,
			AudienceBullets: AudienceBulletsConfig{
				Recruiter:     4,
				HiringManager: 6,
				Peer:          8,
				Default:       5,
			},
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
	// Clean path to prevent path traversal attacks
	cleanPath := filepath.Clean(path)

	// Check if config file exists
	if _, err := os.Stat(cleanPath); os.IsNotExist(err) {
		// Return default config if file doesn't exist
		return DefaultConfig(), nil
	}

	data, err := os.ReadFile(cleanPath) // #nosec G304 - path is cleaned above
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Apply defaults for any missing values
	applyDefaults(&cfg)

	return &cfg, nil
}

// applyDefaults fills in default values for any zero-value fields in the config.
// This ensures that when loading a config file with missing fields, sensible
// defaults are applied rather than leaving them as zero values.
func applyDefaults(cfg *Config) {
	defaults := DefaultConfig()

	// System defaults
	if cfg.System.DataDir == "" {
		cfg.System.DataDir = defaults.System.DataDir
	}
	if cfg.System.LogLevel == "" {
		cfg.System.LogLevel = defaults.System.LogLevel
	}
	if cfg.System.BackupCount == 0 {
		cfg.System.BackupCount = defaults.System.BackupCount
	}
	// Note: AutoBackup is bool, can't distinguish false from unset

	// Profile defaults
	if cfg.Profile.DefaultRole == "" {
		cfg.Profile.DefaultRole = defaults.Profile.DefaultRole
	}
	if cfg.Profile.DefaultAudience == "" {
		cfg.Profile.DefaultAudience = defaults.Profile.DefaultAudience
	}

	// CV defaults
	if cfg.CV.DefaultFormat == "" {
		cfg.CV.DefaultFormat = defaults.CV.DefaultFormat
	}
	if cfg.CV.MaxBullets == 0 {
		cfg.CV.MaxBullets = defaults.CV.MaxBullets
	}

	// AudienceBullets defaults
	if cfg.CV.AudienceBullets.Recruiter == 0 {
		cfg.CV.AudienceBullets.Recruiter = defaults.CV.AudienceBullets.Recruiter
	}
	if cfg.CV.AudienceBullets.HiringManager == 0 {
		cfg.CV.AudienceBullets.HiringManager = defaults.CV.AudienceBullets.HiringManager
	}
	if cfg.CV.AudienceBullets.Peer == 0 {
		cfg.CV.AudienceBullets.Peer = defaults.CV.AudienceBullets.Peer
	}
	if cfg.CV.AudienceBullets.Default == 0 {
		cfg.CV.AudienceBullets.Default = defaults.CV.AudienceBullets.Default
	}

	// Export defaults
	if cfg.Export.DefaultDestination == "" {
		cfg.Export.DefaultDestination = defaults.Export.DefaultDestination
	}
	// Note: AutoOpen is bool, can't distinguish false from unset

	// Display defaults
	if cfg.Display.Theme == "" {
		cfg.Display.Theme = defaults.Display.Theme
	}
	// Note: Animations is bool, can't distinguish false from unset
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

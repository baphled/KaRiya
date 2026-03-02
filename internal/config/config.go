package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

// configPathOverride holds an optional override for the config path.
// This is used by tests to isolate config file writes.
var (
	configPathOverride string
	configPathMu       sync.RWMutex
)

// isTestEnvironment checks if we're running in a test environment.
// This detects both `go test` and test binaries.
func isTestEnvironment() bool {
	// Check if running under go test (the binary name ends with .test)
	executable, err := os.Executable()
	if err == nil && strings.HasSuffix(executable, ".test") {
		return true
	}

	// Check for test flags in os.Args
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-test.") {
			return true
		}
	}

	return false
}

// requireTestIsolation panics if we're in a test environment but the config
// path override hasn't been set. This prevents tests from accidentally
// writing to the user's real config file (BUG-007 prevention).
func requireTestIsolation(operation string) {
	if !isTestEnvironment() {
		return
	}

	configPathMu.RLock()
	hasOverride := configPathOverride != ""
	configPathMu.RUnlock()

	if !hasOverride {
		panic(fmt.Sprintf(
			"BUG-007 PROTECTION: %s called in test without config isolation!\n\n"+
				"Tests must isolate config writes to prevent polluting ~/.kariya/config.yaml.\n\n"+
				"Fix: Call config.SetConfigPathForTesting(path) before using %s,\n"+
				"     or use harness.Setup()/harness.SetupWithOnboarding() which handle isolation.\n\n"+
				"Example:\n"+
				"    tempDir := t.TempDir()\n"+
				"    config.SetConfigPathForTesting(filepath.Join(tempDir, \"config.yaml\"))\n"+
				"    defer config.ResetConfigPath()\n",
			operation, operation,
		))
	}
}

// Config represents the application configuration.
type Config struct {
	System  SystemConfig  `yaml:"system"`
	Profile ProfileConfig `yaml:"profile"`
	CV      CVConfig      `yaml:"cv"`
	Export  ExportConfig  `yaml:"export"`
	Display DisplayConfig `yaml:"display"`
	Scoring ScoringConfig `yaml:"scoring"`
}

// SystemConfig contains system-level configuration.
type SystemConfig struct {
	DataDir     string `yaml:"data_dir"`
	LogLevel    string `yaml:"log_level"`
	AutoBackup  bool   `yaml:"auto_backup"`
	BackupCount int    `yaml:"backup_count"`
}

// ProfileConfig contains user profile configuration.
type ProfileConfig struct {
	Name            string `yaml:"name"`
	Email           string `yaml:"email"`
	DefaultRole     string `yaml:"default_role"`
	DefaultAudience string `yaml:"default_audience"`

	// Narrative CV profile fields
	Title         string   `yaml:"title"`
	Location      string   `yaml:"location"`
	GitHub        string   `yaml:"github"`
	Portfolio     string   `yaml:"portfolio"`
	CoreStrengths []string `yaml:"core_strengths"`
	Languages     []string `yaml:"languages"`
	Frontend      []string `yaml:"frontend"`
	Systems       []string `yaml:"systems"`
	WhatIBring    []string `yaml:"what_i_bring"`
}

// CVConfig contains CV generation configuration.
type CVConfig struct {
	DefaultFormat   string                `yaml:"default_format"`
	MaxBullets      int                   `yaml:"max_bullets"`
	AudienceBullets AudienceBulletsConfig `yaml:"audience_bullets"`
}

// AudienceBulletsConfig defines bullets per company for each audience type.
type AudienceBulletsConfig struct {
	Recruiter     int `yaml:"recruiter"`
	HiringManager int `yaml:"hiring_manager"`
	Peer          int `yaml:"peer"`
	Default       int `yaml:"default"`
}

// ExportConfig contains export configuration.
type ExportConfig struct {
	DefaultDestination string `yaml:"default_destination"`
	AutoOpen           bool   `yaml:"auto_open"`
}

// DisplayConfig contains display/UI configuration.
type DisplayConfig struct {
	Theme      string `yaml:"theme"`
	Animations bool   `yaml:"animations"`
}

// ScoringConfig contains bullet scoring configuration for CV generation.
// This allows customizing how bullets are scored and filtered.
type ScoringConfig struct {
	Weights      ScoringWeights            `yaml:"weights"`
	Thresholds   ScoringThresholds         `yaml:"thresholds"`
	RoleSettings map[string]RoleScoringCfg `yaml:"role_settings"`
}

// ScoringWeights defines the weights for each scoring component.
// All weights must sum to 1.0.
type ScoringWeights struct {
	RoleScore     float64 `yaml:"role_score"`
	AudienceScore float64 `yaml:"audience_score"`
	MetricScore   float64 `yaml:"metric_score"`
	ImpactScore   float64 `yaml:"impact_score"`
	Confidence    float64 `yaml:"confidence"`
}

// ScoringThresholds defines confidence and scoring thresholds.
type ScoringThresholds struct {
	FactDefaultConfidence  float64 `yaml:"fact_default_confidence"`
	EventDefaultConfidence float64 `yaml:"event_default_confidence"`
	HighConfidence         float64 `yaml:"high_confidence"`
	HighImpactConfidence   float64 `yaml:"high_impact_confidence"`
}

// RoleScoringCfg defines scoring settings for a specific role.
type RoleScoringCfg struct {
	MinConfidence        float64 `yaml:"min_confidence"`
	MaxBulletsPerCompany int     `yaml:"max_bullets_per_company"`
}

// ValidateWeights checks that scoring weights sum to 1.0 within tolerance.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (s *ScoringConfig) ValidateWeights() error {
	sum := s.Weights.RoleScore + s.Weights.AudienceScore + s.Weights.MetricScore +
		s.Weights.ImpactScore + s.Weights.Confidence

	const tolerance = 0.01
	if sum < 1.0-tolerance || sum > 1.0+tolerance {
		return fmt.Errorf("scoring weights must sum to 1.0, got %.4f", sum)
	}
	return nil
}

// DefaultConfig returns the default configuration.
//
// Returns:
//   - A fully initialized Config ready for use.
//
// Side effects:
//   - None.
func DefaultConfig() *Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory if home dir unavailable
		homeDir = "."
	}
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
			Languages:       []string{},
			Frontend:        []string{},
			Systems:         []string{},
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
		Scoring: ScoringConfig{
			Weights: ScoringWeights{
				RoleScore:     0.25,
				AudienceScore: 0.20,
				MetricScore:   0.20,
				ImpactScore:   0.20,
				Confidence:    0.15,
			},
			Thresholds: ScoringThresholds{
				FactDefaultConfidence:  0.85,
				EventDefaultConfidence: 0.80,
				HighConfidence:         0.80,
				HighImpactConfidence:   0.85,
			},
			RoleSettings: map[string]RoleScoringCfg{
				"principal": {MinConfidence: 0.80, MaxBulletsPerCompany: 4},
				"staff":     {MinConfidence: 0.75, MaxBulletsPerCompany: 5},
				"em":        {MinConfidence: 0.75, MaxBulletsPerCompany: 4},
				"senior_ic": {MinConfidence: 0.75, MaxBulletsPerCompany: 5},
			},
		},
	}
}

// SetConfigPathForTesting overrides the config path for testing purposes.
//
// Expected:
//   - Must be a valid string.
//
// Side effects:
//   - None.
func SetConfigPathForTesting(path string) {
	configPathMu.Lock()
	defer configPathMu.Unlock()
	configPathOverride = path
}

// SwapConfigPathForTesting sets a new config path and returns the previous one.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func SwapConfigPathForTesting(path string) string {
	configPathMu.Lock()
	defer configPathMu.Unlock()
	prev := configPathOverride
	configPathOverride = path
	return prev
}

// ResetConfigPath clears the config path override and restores default behavior.
//
// Side effects:
//   - None.
func ResetConfigPath() {
	configPathMu.Lock()
	defer configPathMu.Unlock()
	configPathOverride = ""
}

// GetConfigPath returns the path to the config file.
// If SetConfigPathForTesting was called, returns the overridden path.
// Otherwise, returns the default path: ~/.kariya/config.yaml.
func GetConfigPath() (string, error) {
	configPathMu.RLock()
	override := configPathOverride
	configPathMu.RUnlock()

	if override != "" {
		return override, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	return filepath.Join(homeDir, ".kariya", "config.yaml"), nil
}

// LoadConfig loads configuration from the default location.
// In test environments, this will panic if SetConfigPathForTesting hasn't been called.
func LoadConfig() (*Config, error) {
	requireTestIsolation("config.LoadConfig()")

	path, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	return LoadConfigFromPath(path)
}

// LoadConfigFromPath loads configuration from a specific file path.
func LoadConfigFromPath(path string) (*Config, error) {
	// Clean path to prevent path traversal attacks
	cleanPath := filepath.Clean(path)

	// Check if config file exists
	if _, err := os.Stat(cleanPath); os.IsNotExist(err) {
		// Return default config if file doesn't exist
		return DefaultConfig(), nil
	}

	data, err := os.ReadFile(cleanPath)
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
	// AutoBackup is bool, can't distinguish false from unset

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
	// AutoOpen is bool, can't distinguish false from unset

	// Display defaults
	if cfg.Display.Theme == "" {
		cfg.Display.Theme = defaults.Display.Theme
	}
	// Animations is bool, can't distinguish false from unset

	// Scoring defaults (auto-migration for configs without scoring section)
	applyScoringDefaults(cfg, defaults)
}

// applyScoringDefaults applies default values for the scoring configuration.
// This enables auto-migration when loading configs without a scoring section.
//
// To avoid the issue where a user-provided value of 0.0 gets overwritten,
// we check if the entire scoring section appears uninitialized. If ALL weights
// are zero AND ALL thresholds are zero AND role settings are empty, we assume
// the section is missing and apply all defaults. This means users who want to
// set some values to zero must set at least one non-zero value in the section.
func applyScoringDefaults(cfg, defaults *Config) {
	// Check if the entire scoring section appears uninitialized:
	// - All weights are zero
	// - All thresholds are zero
	// - No role settings defined
	allWeightsZero := cfg.Scoring.Weights.RoleScore == 0 &&
		cfg.Scoring.Weights.AudienceScore == 0 &&
		cfg.Scoring.Weights.MetricScore == 0 &&
		cfg.Scoring.Weights.ImpactScore == 0 &&
		cfg.Scoring.Weights.Confidence == 0

	allThresholdsZero := cfg.Scoring.Thresholds.FactDefaultConfidence == 0 &&
		cfg.Scoring.Thresholds.EventDefaultConfidence == 0 &&
		cfg.Scoring.Thresholds.HighConfidence == 0 &&
		cfg.Scoring.Thresholds.HighImpactConfidence == 0

	noRoleSettings := len(cfg.Scoring.RoleSettings) == 0

	// If the entire section is uninitialized, apply all defaults
	if allWeightsZero && allThresholdsZero && noRoleSettings {
		cfg.Scoring.Weights = defaults.Scoring.Weights
		cfg.Scoring.Thresholds = defaults.Scoring.Thresholds
		cfg.Scoring.RoleSettings = defaults.Scoring.RoleSettings
	}
	// Otherwise, the user has customized at least part of the scoring section,
	// so we respect their configuration (including any explicit zeros).
}

// SaveConfig saves configuration to the default location.
//
// Expected:
//   - config must be a valid configuration object.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func SaveConfig(cfg *Config) error {
	requireTestIsolation("config.SaveConfig()")

	path, err := GetConfigPath()
	if err != nil {
		return err
	}

	return SaveConfigToPath(cfg, path)
}

// SaveConfigToPath saves configuration to a specific file path.
//
// Expected:
//   - config must be a valid configuration object.
//   - Must be a valid string.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func SaveConfigToPath(cfg *Config, path string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

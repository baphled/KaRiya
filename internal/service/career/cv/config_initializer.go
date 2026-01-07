package cv

import (
	"context"
	"fmt"
	"time"

	career "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
)

// ConfigInitializer handles initialization of CV configuration system
type ConfigInitializer struct {
	configManager ConfigManager
	logger        *logger.Logger
}

// NewConfigInitializer creates a new configuration initializer
func NewConfigInitializer(configManager ConfigManager, log *logger.Logger) *ConfigInitializer {
	return &ConfigInitializer{
		configManager: configManager,
		logger:        log,
	}
}

// Initialize sets up the configuration system and creates default configs if needed
func (ci *ConfigInitializer) Initialize(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if ci.logger != nil {
		ci.logger.Info("Initializing CV configuration system")
	}

	// Verify directory exists and is accessible (for YAML config managers)
	if yamlMgr, ok := ci.configManager.(*YAMLConfigManager); ok {
		if err := yamlMgr.VerifyDirectory(); err != nil {
			if ci.logger != nil {
				ci.logger.Error("Failed to verify config directory: %v", err)
			}
			return fmt.Errorf("failed to verify config directory: %w", err)
		}
	}

	// Check if any configs exist
	configs, err := ci.configManager.ListConfigs(ctx)
	if err != nil {
		if ci.logger != nil {
			ci.logger.Error("Failed to list existing configs: %v", err)
		}
		return fmt.Errorf("failed to list existing configs: %w", err)
	}

	// If no configs exist, create default ones
	if len(configs) == 0 {
		if ci.logger != nil {
			ci.logger.Info("No existing configs found, creating defaults")
		}
		if err := ci.createDefaultConfigs(ctx); err != nil {
			if ci.logger != nil {
				ci.logger.Error("Failed to create default configs: %v", err)
			}
			return fmt.Errorf("failed to create default configs: %w", err)
		}
		if ci.logger != nil {
			ci.logger.Info("Default configs created successfully")
		}
	} else {
		if ci.logger != nil {
			ci.logger.Info("Found %d existing configs", len(configs))
		}
	}

	return nil
}

// createDefaultConfigs creates a set of default CV configurations
func (ci *ConfigInitializer) createDefaultConfigs(ctx context.Context) error {
	defaultConfigs := []*career.CVConfig{
		{
			Name:           "Principal Engineer",
			TargetRole:     "principal",
			TargetAudience: "hiring_manager",
			EventFilters: map[string]interface{}{
				"categories": []string{"technical", "leadership", "product"},
			},
		},
		{
			Name:           "Staff Engineer",
			TargetRole:     "staff",
			TargetAudience: "hiring_manager",
			EventFilters: map[string]interface{}{
				"categories": []string{"technical", "product"},
			},
		},
		{
			Name:           "Engineering Manager",
			TargetRole:     "em",
			TargetAudience: "hiring_manager",
			EventFilters: map[string]interface{}{
				"categories": []string{"leadership", "mentoring"},
			},
		},
		{
			Name:           "Senior IC",
			TargetRole:     "senior_ic",
			TargetAudience: "recruiter",
			EventFilters: map[string]interface{}{
				"categories": []string{"technical", "achievement"},
			},
		},
	}

	for _, config := range defaultConfigs {
		// Set timestamps
		now := time.Now()
		config.CreatedAt = now
		config.UpdatedAt = now

		// Validate config
		if err := config.Validate(); err != nil {
			if ci.logger != nil {
				ci.logger.Warn("Default config validation failed for %s: %v", config.Name, err)
			}
			return fmt.Errorf("default config validation failed for %s: %w", config.Name, err)
		}

		// Save config
		if err := ci.configManager.SaveConfig(ctx, config); err != nil {
			if ci.logger != nil {
				ci.logger.Error("Failed to save default config %s: %v", config.Name, err)
			}
			return fmt.Errorf("failed to save default config %s: %w", config.Name, err)
		}

		if ci.logger != nil {
			ci.logger.Info("Created default config: %s", config.Name)
		}
	}

	return nil
}

// ValidateSetup validates that the configuration system is properly set up
func (ci *ConfigInitializer) ValidateSetup(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if ci.logger != nil {
		ci.logger.Info("Validating configuration setup")
	}

	// Try to list configs - this validates the system is accessible
	configs, err := ci.configManager.ListConfigs(ctx)
	if err != nil {
		if ci.logger != nil {
			ci.logger.Error("Configuration validation failed: %v", err)
		}
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	if ci.logger != nil {
		ci.logger.Info("Configuration validation passed, found %d configs", len(configs))
	}

	return nil
}

// EnsureConfigExists ensures that at least one config exists, creating defaults if needed
func (ci *ConfigInitializer) EnsureConfigExists(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	configs, err := ci.configManager.ListConfigs(ctx)
	if err != nil {
		return fmt.Errorf("failed to list configs: %w", err)
	}

	// If no configs exist, create defaults
	if len(configs) == 0 {
		if ci.logger != nil {
			ci.logger.Info("No configs found, creating defaults")
		}
		return ci.createDefaultConfigs(ctx)
	}

	return nil
}

// CheckAndCreateConfig checks if a config exists and creates it if it doesn't
func (ci *ConfigInitializer) CheckAndCreateConfig(ctx context.Context, configName string) (*career.CVConfig, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if configName == "" {
		return nil, fmt.Errorf("config name cannot be empty")
	}

	// Try to load the config
	config, err := ci.configManager.LoadConfig(ctx, configName)
	if err == nil {
		// Config exists
		if ci.logger != nil {
			ci.logger.Info("Config already exists: %s", configName)
		}
		return config, nil
	}

	// Config doesn't exist, check if it's a not found error
	if err != ErrConfigNotFound {
		if ci.logger != nil {
			ci.logger.Error("Error loading config %s: %v", configName, err)
		}
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	// Config not found, create it
	if ci.logger != nil {
		ci.logger.Info("Config not found, creating new config: %s", configName)
	}

	// Create a new config with sensible defaults
	config = &career.CVConfig{
		Name:           configName,
		TargetRole:     "staff",
		TargetAudience: "hiring_manager",
		EventFilters: map[string]interface{}{
			"categories": []string{"technical"},
		},
	}

	// Validate the config
	if err := config.Validate(); err != nil {
		if ci.logger != nil {
			ci.logger.Error("Config validation failed: %v", err)
		}
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	// Save the config
	if err := ci.configManager.SaveConfig(ctx, config); err != nil {
		if ci.logger != nil {
			ci.logger.Error("Failed to save new config: %v", err)
		}
		return nil, fmt.Errorf("failed to save new config: %w", err)
	}

	if ci.logger != nil {
		ci.logger.Info("New config created successfully: %s", configName)
	}

	return config, nil
}

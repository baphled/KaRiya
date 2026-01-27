// Package bootstrap handles pre-application setup including first-run onboarding
// and service initialization. This runs before the main app loop starts.
package bootstrap

import (
	"errors"

	"github.com/baphled/kariya/internal/config"
	"github.com/baphled/kariya/internal/logger"
	careerservice "github.com/baphled/kariya/internal/service/career"
)

// ErrUserAborted is returned when the user aborts onboarding (e.g., Ctrl+C).
var ErrUserAborted = errors.New("user aborted onboarding")

// Result contains everything needed to start the main application.
type Result struct {
	Config   *config.Config
	Services *Services
}

// Run performs pre-application bootstrap:
// 1. Loads configuration
// 2. Checks if profile is complete
// 3. Runs onboarding wizard if needed (blocking)
// 4. Initializes services
// 5. Returns everything the app needs
//
// Returns nil if the user aborted onboarding (Ctrl+C).
func Run(careerService *careerservice.Service, log *logger.Logger) (*Result, error) {
	// Load existing configuration.
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error("Failed to load config: %v, using defaults", err)
		cfg = config.DefaultConfig()
	}

	// Check if onboarding is needed.
	if !IsProfileComplete(cfg) {
		log.Info("Profile incomplete, starting onboarding wizard")

		profile, err := runOnboarding(&cfg.Profile)
		if err != nil {
			return nil, err
		}

		cfg.Profile = *profile

		if err := config.SaveConfig(cfg); err != nil {
			log.Error("Failed to save config after onboarding: %v", err)
		} else {
			log.Info("Profile saved successfully")
		}
	}

	// Initialize services.
	services := InitServices(careerService, cfg, log)

	return &Result{
		Config:   cfg,
		Services: services,
	}, nil
}

// RunWithConfig performs bootstrap with a pre-loaded config.
// This is useful for testing or when config is already loaded.
func RunWithConfig(cfg *config.Config, careerService *careerservice.Service, log *logger.Logger) (*Result, error) {
	// Check if onboarding is needed.
	if !IsProfileComplete(cfg) {
		log.Info("Profile incomplete, starting onboarding wizard")

		profile, err := runOnboarding(&cfg.Profile)
		if err != nil {
			return nil, err
		}

		cfg.Profile = *profile

		if err := config.SaveConfig(cfg); err != nil {
			log.Error("Failed to save config after onboarding: %v", err)
		}
	}

	services := InitServices(careerService, cfg, log)

	return &Result{
		Config:   cfg,
		Services: services,
	}, nil
}

// SkipOnboarding initializes services without checking profile completeness.
// This is primarily for testing scenarios where onboarding should be bypassed.
func SkipOnboarding(cfg *config.Config, careerService *careerservice.Service, log *logger.Logger) *Result {
	services := InitServices(careerService, cfg, log)

	return &Result{
		Config:   cfg,
		Services: services,
	}
}

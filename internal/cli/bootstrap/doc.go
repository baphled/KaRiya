// Package bootstrap provides application initialization and dependency injection.
//
// # Overview
//
// The bootstrap package handles the creation and wiring of all application
// components including services, repositories, and configuration. It provides
// a centralized location for dependency management during application startup.
//
// # Responsibilities
//
//   - Service initialization and configuration
//   - Repository setup (database connections)
//   - Configuration loading and validation
//   - Onboarding flow coordination
//   - Test environment setup
//
// # Usage
//
// Bootstrap the application:
//
//	cfg := bootstrap.LoadConfig()
//	services := bootstrap.InitializeServices(cfg)
//
// For testing:
//
//	env := bootstrap.NewTestEnvironment()
//	defer env.Cleanup()
package bootstrap

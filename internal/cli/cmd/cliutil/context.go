// Package cliutil provides shared utilities and interfaces for CLI commands.
package cliutil

import careerservice "github.com/baphled/kariya/internal/service/career"

// ServiceContext defines the interface for CLI context that subcommands need.
type ServiceContext interface {
	// Service returns the initialized career service.
	// Returns nil if InitService has not been called yet.
	Service() *careerservice.Service
}

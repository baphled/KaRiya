// Package service provides CLI-specific service implementations.
//
// # Overview
//
// The service package contains service implementations that bridge the domain
// layer with the CLI layer. These services wrap domain services and add
// CLI-specific functionality like progress reporting and TUI integration.
//
// # Responsibilities
//
//   - Event service with progress callbacks
//   - Career data management
//   - Import/export operations
//   - CV generation coordination
//
// # Architecture
//
// Services in this package wrap domain services from internal/service:
//
//	CLIService -> DomainService
//
// This allows the CLI to add TUI-specific behavior while maintaining
// clean domain boundaries.
package service

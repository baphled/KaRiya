// Package configtypes provides configuration domain types and enums.
//
// # Overview
//
// The configtypes package defines the core types used throughout the
// configuration system. These types are shared between the intents package
// and the configuration implementation to avoid circular dependencies.
//
// # Types
//
//   - ConfigurationDomain: Enum for configuration sections (system, profile, export, ui)
//   - ConfigurationSetting: Individual configuration entry with metadata
//   - ConfigurationState: Workflow state for configuration editing
//   - ConfigurationChanges: Tracks pending configuration modifications
//
// # Usage
//
// These types are primarily used by the configure intent and related screens
// to provide type-safe configuration management.
package configtypes

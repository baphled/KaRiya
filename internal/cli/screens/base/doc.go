// Package base provides base screen implementations.
//
// # Overview
//
// The base package contains base screen types and utilities that
// provide common functionality for all screens in the application.
//
// # Types
//
//   - BaseScreen: Common screen foundation
//   - DetailScreen: Base for detail views
//   - ScreenResult: Result types for screen communication
//
// # Usage
//
// Embed in a custom screen:
//
//	type MyScreen struct {
//	    *base.BaseScreen
//	    // ... custom fields
//	}
package base

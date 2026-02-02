// Package models provides form wrapper types for the TUI application.
//
// # Overview
//
// The models package contains wrapper types that encapsulate form state and
// behavior. These wrappers provide a cleaner interface than using huh forms
// directly in intents.
//
// # Deprecation Notice
//
// This package is DEPRECATED. New code should use screens/*FormScreen instead
// of models.*Form wrappers. The forms package should be the only place that
// imports huh directly.
//
// # Migration Path
//
// Old: models.CaptureForm
// New: screens/capture/form_screen.go with embedded forms.Form
//
// See docs/rules/FORMS_WORKFLOW_GUIDE.md for the new pattern.
package models

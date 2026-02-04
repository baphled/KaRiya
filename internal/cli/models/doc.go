// Package models provides shared message types and base model utilities for the TUI application.
//
// # Overview
//
// The models package contains message types (e.g. SubmitMsg) and base model
// utilities (BaseStandardModel) used across the application. Legacy form
// wrappers (CaptureForm, SkillForm) have been migrated to screens/*FormScreen.
//
// # Migration Status
//
// CaptureForm → screens/capture.EventFormScreen (completed)
// SkillForm → screens/skills.SkillFormScreen (completed)
//
// Remaining types in this package:
//   - BaseStandardModel: shared model base
//   - SubmitMsg and other message types
//   - Modal wrappers (MetadataEditorModelNew, BurstSuggestionModelNew, FactEditorModelNew)
//
// See docs/rules/FORMS_WORKFLOW_GUIDE.md for the new pattern.
package models

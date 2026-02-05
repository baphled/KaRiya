// Package forms provides form configuration and builders for the TUI application.
//
// # Overview
//
// The forms package wraps the charmbracelet/huh library with KaRiya-specific
// configuration and styling. It provides a consistent interface for creating
// interactive forms across the application.
//
// # Architecture
//
// The forms package is the ONLY package that should import huh directly.
// All other packages must use the form builders and wrappers provided here.
//
// # Usage
//
// Create a form using entity-specific constructors:
//
//	data := &forms.SkillFormData{Name: "Go"}
//	form := forms.NewSkillForm(data, width, height)
//
// Check form state:
//
//	if forms.IsCompleted(form) {
//	    data := forms.GetData(form)
//	}
//
// For more details, see the forms documentation in docs/FORMS_GUIDE.md.
package forms

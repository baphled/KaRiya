package models

import "errors"

// Custom errors for model operations
var (
	ErrInvalidBreadcrumbIndex = errors.New("invalid breadcrumb index")
	ErrEmptyNavigationHistory = errors.New("navigation history is empty")
	ErrInvalidModelState      = errors.New("invalid model state")
	ErrValidationFailed       = errors.New("validation failed")
	ErrContextNotAvailable    = errors.New("context not available")
	ErrShortcutNotRegistered  = errors.New("shortcut not registered")
)

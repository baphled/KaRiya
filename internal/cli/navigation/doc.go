// Package navigation provides keyboard navigation handlers and utilities.
//
// # Overview
//
// The navigation package implements consistent keyboard navigation patterns
// across the TUI application. It provides reusable handlers for common
// navigation scenarios like lists, forms, and menus.
//
// # Available Handlers
//
//   - ListNavigationHandler: Up/down navigation with vim keys
//   - FormNavigationHandler: Tab/shift-tab between fields
//   - ModalNavigationHandler: Modal-specific key handling
//
// # Usage
//
// Use with TableBehavior:
//
//	table := behaviors.NewTableBehavior(...)
//	table.SetNavigationHandler(navigation.NewListHandler())
package navigation

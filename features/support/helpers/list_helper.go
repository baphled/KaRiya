// Package helpers provides test helper utilities for BDD scenarios.
package helpers

import (
	"context"
	"fmt"
	"strings"

	"github.com/baphled/kariya/features/support"
	"github.com/baphled/kariya/internal/testutil/e2e"
	"github.com/cucumber/godog"
)

// ListHelper provides list navigation abstraction.
type ListHelper struct {
	env *e2e.TestEnv
	ctx context.Context
}

// NewListHelper creates a list helper for the given context.
//
// Expected:
//   - ctx must contain a valid TestEnv set by the BDD environment.
//
// Returns:
//   - A ListHelper bound to the test environment, or error if env is nil.
//
// Side effects:
//   - None.
func NewListHelper(ctx context.Context) (*ListHelper, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return nil, godog.ErrPending
	}
	return &ListHelper{env: env, ctx: ctx}, nil
}

// NavigateToItem navigates through a list until the item is visible/focused.
//
// Expected:
//   - itemName must be a non-empty string matching a visible list item.
//
// Returns:
//   - error if item not found after max attempts.
//
// Side effects:
//   - Sends navigation key events to the test environment.
func (l *ListHelper) NavigateToItem(itemName string) error {
	const maxAttempts = 20

	for range maxAttempts {
		view := l.env.GetView()
		if strings.Contains(view, itemName) {
			return nil
		}
		l.env.NavigateDown()
	}

	return fmt.Errorf("item '%s' not found after %d attempts", itemName, maxAttempts)
}

// SelectCurrentItem confirms selection of the currently focused item.
//
// Returns:
//   - nil always (error interface for consistency).
//
// Side effects:
//   - Sends confirm key event to the test environment.
func (l *ListHelper) SelectCurrentItem() error {
	l.env.Confirm()
	return nil
}

// SelectItemByName navigates to an item and selects it.
//
// Expected:
//   - itemName must match a visible list item.
//
// Returns:
//   - error if the item cannot be found.
//
// Side effects:
//   - Sends navigation and confirm key events to the test environment.
func (l *ListHelper) SelectItemByName(itemName string) error {
	if err := l.NavigateToItem(itemName); err != nil {
		return err
	}
	return l.SelectCurrentItem()
}

// NavigateUp moves up one item in the list.
//
// Side effects:
//   - Sends up navigation key event to the test environment.
func (l *ListHelper) NavigateUp() {
	l.env.NavigateUp()
}

// NavigateDown moves down one item in the list.
//
// Side effects:
//   - Sends down navigation key event to the test environment.
func (l *ListHelper) NavigateDown() {
	l.env.NavigateDown()
}

// IsItemVisible checks if an item is currently visible.
//
// Expected:
//   - itemName must be a non-empty string.
//
// Returns:
//   - true if the item name appears in the current view.
//
// Side effects:
//   - None.
func (l *ListHelper) IsItemVisible(itemName string) bool {
	view := l.env.GetView()
	return strings.Contains(view, itemName)
}

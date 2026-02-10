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
func NewListHelper(ctx context.Context) (*ListHelper, error) {
	env := support.GetAppEnv(ctx)
	if env == nil {
		return nil, godog.ErrPending
	}
	return &ListHelper{env: env, ctx: ctx}, nil
}

// NavigateToItem navigates through a list until the item is visible/focused.
func (l *ListHelper) NavigateToItem(itemName string) error {
	const maxAttempts = 20

	for i := 0; i < maxAttempts; i++ {
		view := l.env.GetView()
		if strings.Contains(view, itemName) {
			return nil
		}
		l.env.NavigateDown()
	}

	return fmt.Errorf("item '%s' not found after %d attempts", itemName, maxAttempts)
}

// SelectCurrentItem confirms selection of the currently focused item.
func (l *ListHelper) SelectCurrentItem() error {
	l.env.Confirm()
	return nil
}

// SelectItemByName navigates to an item and selects it.
func (l *ListHelper) SelectItemByName(itemName string) error {
	if err := l.NavigateToItem(itemName); err != nil {
		return err
	}
	return l.SelectCurrentItem()
}

// NavigateUp moves up one item in the list.
func (l *ListHelper) NavigateUp() {
	l.env.NavigateUp()
}

// NavigateDown moves down one item in the list.
func (l *ListHelper) NavigateDown() {
	l.env.NavigateDown()
}

// IsItemVisible checks if an item is currently visible.
func (l *ListHelper) IsItemVisible(itemName string) bool {
	view := l.env.GetView()
	return strings.Contains(view, itemName)
}

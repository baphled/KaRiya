package components

import (
	"strings"

	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/charmbracelet/lipgloss"
)

// ListContainer displays a list of items with optional pagination and empty state handling.
// It provides consistent list layout with support for pagination information and
// custom empty state messages.
// ListContainer is a stateless rendering component.
type ListContainer struct {
	items           []string
	emptyMessage    string
	paginationInfo  string
	hasPagination   bool
	hasEmptyMessage bool
}

// NewListContainer creates a new ListContainer.
func NewListContainer() *ListContainer {
	return &ListContainer{
		items:           []string{},
		emptyMessage:    "No items to display",
		paginationInfo:  "",
		hasPagination:   false,
		hasEmptyMessage: false,
	}
}

// SetItems sets the list items to display.
// This method uses the builder pattern to allow method chaining.
func (lc *ListContainer) SetItems(items []string) *ListContainer {
	lc.items = items
	return lc
}

// SetEmptyStateMessage sets the message to display when the list is empty.
// This method uses the builder pattern to allow method chaining.
func (lc *ListContainer) SetEmptyStateMessage(message string) *ListContainer {
	lc.emptyMessage = message
	lc.hasEmptyMessage = true
	return lc
}

// SetPaginationInfo sets the pagination information to display.
// This method uses the builder pattern to allow method chaining.
func (lc *ListContainer) SetPaginationInfo(info string) *ListContainer {
	lc.paginationInfo = info
	lc.hasPagination = true
	return lc
}

// Render returns the styled list container as a string.
// It displays items or an empty state message, with optional pagination information.
func (lc *ListContainer) Render() string {
	var content string

	// Handle empty list
	if len(lc.items) == 0 {
		emptyStyle := lipgloss.NewStyle().
			Foreground(styles.ColorTextMuted).
			Italic(true)
		content = emptyStyle.Render(lc.emptyMessage)
	} else {
		var parts []string

		// Render items
		itemStyle := lipgloss.NewStyle().
			Foreground(styles.ColorTextPrimary)

		for _, item := range lc.items {
			parts = append(parts, itemStyle.Render(item))
		}

		content = strings.Join(parts, "\n")
	}

	// Add pagination info if present (even for empty lists)
	if lc.hasPagination && lc.paginationInfo != "" {
		paginationStyle := lipgloss.NewStyle().
			Foreground(styles.ColorTextSecondary).
			MarginTop(1)
		content += "\n" + paginationStyle.Render(lc.paginationInfo)
	}

	return content
}

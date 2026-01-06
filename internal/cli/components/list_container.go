package components

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

// ListContainer displays a table-based list and ensures compatibility with Bubble Tea components.
type ListContainer struct {
	table           table.Model
	emptyMessage    string
	items           []string
	hasPagination   bool
	hasEmptyMessage bool
	paginationInfo  string
}

// NewListContainer initializes a ListContainer instance with default configurations.
func NewListContainer() *ListContainer {
	columns := []table.Column{
		{Title: "Column 1", Width: 20},
		{Title: "Column 2", Width: 30},
	}
	rows := []table.Row{
		{"Data 1", "Data 2"},
		{"Data 3", "Data 4"},
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)
	return &ListContainer{
		table:        t,
		emptyMessage: "No items to display",
	}
}

// SetItems sets the list of string items (and also updates the table, if needed).
func (lc *ListContainer) SetItems(items []string) *ListContainer {
	lc.items = items
	// Reset table rows as a fallback
	rows := make([]table.Row, len(items))
	for i, item := range items {
		rows[i] = table.Row{item}
	}
	lc.table.SetRows(rows)
	return lc
}

// SetEmptyStateMessage sets the message to display when the list is empty.
func (lc *ListContainer) SetEmptyStateMessage(message string) *ListContainer {
	lc.emptyMessage = message
	lc.hasEmptyMessage = message != ""
	return lc
}

// SetPaginationInfo sets pagination info to be displayed if desired.
func (lc *ListContainer) SetPaginationInfo(info string) *ListContainer {
	lc.paginationInfo = info
	lc.hasPagination = info != ""
	return lc
}

// Render combines list, empty state, and pagination into the rendered output expected by tests.
func (lc *ListContainer) Render() string {
	style := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(1, 2)
	if len(lc.items) == 0 {
		msg := lc.emptyMessage
		if msg == "" {
			msg = "No items to display"
		}
		return style.Render(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true).Render(msg))
	}
	// Combine all items, separated by newlines
	out := ""
	for _, item := range lc.items {
		out += item + "\n"
	}
	if lc.hasPagination && lc.paginationInfo != "" {
		out += lc.paginationInfo + "\n"
	}
	return style.Render(out)
}

// View renders the ListContainer using table.Model (backward compatibility).
func (lc *ListContainer) View() string {
	style := lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(1, 2)
	if len(lc.table.Rows()) == 0 && lc.emptyMessage != "" {
		return style.Render(lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true).Render(lc.emptyMessage))
	}
	return style.Render(lc.table.View())
}

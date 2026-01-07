package components

import (
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/table"
)

// TestTableListContainer_NoHeader verifies header is not rendered after refactor
func TestTableListContainer_NoHeader(t *testing.T) {
	columns := []table.Column{
		{Title: "Name", Width: 20},
		{Title: "Value", Width: 20},
	}

	rows := []table.Row{
		{"Item 1", "Value 1"},
		{"Item 2", "Value 2"},
	}

	tableModel := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithHeight(5),
	)

	container := NewTableListContainer(tableModel, "Test Container", 80)
	output := container.Render()

	// Should NOT contain breadcrumb separator
	if strings.Contains(output, "▸") {
		t.Error("Output should not contain breadcrumb separator (▸) - header should be removed")
	}

	// Should NOT contain "KaRiya" branding (from header)
	if strings.Contains(output, "KaRiya") {
		t.Error("Output should not contain 'KaRiya' branding - header should be removed")
	}

	// Should still contain table content
	if !strings.Contains(output, "Item 1") {
		t.Error("Expected table content to be present")
	}
}

// TestTableListContainer_NoHelpFooter verifies help footer is not rendered
func TestTableListContainer_NoHelpFooter(t *testing.T) {
	columns := []table.Column{
		{Title: "Name", Width: 20},
	}

	rows := []table.Row{
		{"Item 1"},
	}

	tableModel := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithHeight(5),
	)

	container := NewTableListContainer(tableModel, "Test Container", 80)
	output := container.Render()

	// Should NOT contain navigation help patterns
	helpPatterns := []string{
		"↑/k",
		"Navigate up",
		"Enter: Confirm",
	}

	for _, pattern := range helpPatterns {
		if strings.Contains(output, pattern) {
			t.Errorf("Output should not contain help pattern '%s' - help footer should be removed", pattern)
		}
	}
}

// TestTableListContainer_KeepsPagination verifies pagination still works
func TestTableListContainer_KeepsPagination(t *testing.T) {
	columns := []table.Column{
		{Title: "Name", Width: 20},
	}

	rows := []table.Row{
		{"Item 1"},
		{"Item 2"},
		{"Item 3"},
	}

	tableModel := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithHeight(5),
	)

	container := NewTableListContainer(tableModel, "Test Container", 80)
	container.SetPaginationInfo("Page 1 of 1 | Items: 3")
	output := container.Render()

	if !strings.Contains(output, "Page 1 of 1") {
		t.Error("Expected pagination info to be present")
	}

	if !strings.Contains(output, "Items: 3") {
		t.Error("Expected item count in pagination to be present")
	}
}

// TestTableListContainer_KeepsTable verifies table content still renders
func TestTableListContainer_KeepsTable(t *testing.T) {
	columns := []table.Column{
		{Title: "Date", Width: 15},
		{Title: "Event", Width: 40},
	}

	rows := []table.Row{
		{"2025-01-01", "First event"},
		{"2025-01-02", "Second event"},
	}

	tableModel := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithHeight(5),
	)

	container := NewTableListContainer(tableModel, "Events", 80)
	output := container.Render()

	// Table content should be present
	if !strings.Contains(output, "First event") {
		t.Error("Expected first event in output")
	}

	if !strings.Contains(output, "Second event") {
		t.Error("Expected second event in output")
	}

	// Column headers should be present (from table itself, not HeaderModel)
	if !strings.Contains(output, "Date") || !strings.Contains(output, "Event") {
		t.Error("Expected column headers from table")
	}
}

// TestTableListContainer_EmptyState verifies empty state message works
func TestTableListContainer_EmptyState(t *testing.T) {
	columns := []table.Column{
		{Title: "Name", Width: 20},
	}

	tableModel := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
		table.WithHeight(5),
	)

	container := NewTableListContainer(tableModel, "Empty Test", 80)
	container.SetEmptyStateMessage("No items found")
	output := container.Render()

	if !strings.Contains(output, "No items found") {
		t.Error("Expected empty state message")
	}
}

// TestTableListContainer_ErrorState verifies error message display
func TestTableListContainer_ErrorState(t *testing.T) {
	columns := []table.Column{
		{Title: "Name", Width: 20},
	}

	tableModel := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
		table.WithHeight(5),
	)

	container := NewTableListContainer(tableModel, "Error Test", 80)
	container.SetErrorMessage("Failed to load data")
	output := container.Render()

	if !strings.Contains(output, "Failed to load data") {
		t.Error("Expected error message")
	}
}

package models

import (
	"github.com/baphled/kariya/internal/domain/career"
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/cli/importer"
	"github.com/baphled/kariya/internal/cli/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ImportReviewMsg signals that import review is complete
type ImportReviewMsg struct {
	SelectedRows []int  // Row numbers to import
	Action       string // "import", "cancel"
}

// ImportResultMsg signals that import is complete
type ImportResultMsg struct {
	Result *importer.ImportResult
	Error  error
}

// ImportReviewModel handles the import review screen
type ImportReviewModel struct {
	ParsedRows      []*importer.ParsedRow
	selectedRows    map[int]bool // Track selected rows by row number
	focusedRowIndex int          // Index in the display (0-based)
	summary         map[string]int
	width           int
	height          int
	err             error
}

// NewImportReviewModel creates a new import review model
func NewImportReviewModel(parsedRows []*importer.ParsedRow) *ImportReviewModel {
	model := &ImportReviewModel{
		ParsedRows:      parsedRows,
		selectedRows:    make(map[int]bool),
		focusedRowIndex: 0,
		width:           80,
		height:          24,
		summary:         calculateSummary(parsedRows),
	}

	// Pre-select all valid, non-duplicate rows
	for _, row := range parsedRows {
		if row.IsValid && !row.IsDuplicate {
			model.selectedRows[row.RowNumber] = true
		}
	}

	return model
}

// Update handles user input
func (m *ImportReviewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.focusedRowIndex > 0 {
				m.focusedRowIndex--
			}
			return m, nil

		case "down", "j":
			if m.focusedRowIndex < len(m.ParsedRows)-1 {
				m.focusedRowIndex++
			}
			return m, nil

		case "space":
			// Toggle selection for current row
			if m.focusedRowIndex < len(m.ParsedRows) {
				row := m.ParsedRows[m.focusedRowIndex]
				// Only allow selection of valid, non-duplicate rows
				if row.IsValid && !row.IsDuplicate {
					m.selectedRows[row.RowNumber] = !m.selectedRows[row.RowNumber]
				}
			}
			return m, nil

		case "a":
			// Select all valid rows
			for _, row := range m.ParsedRows {
				if row.IsValid && !row.IsDuplicate {
					m.selectedRows[row.RowNumber] = true
				}
			}
			return m, nil

		case "d":
			// Deselect all rows
			m.selectedRows = make(map[int]bool)
			return m, nil

		case "enter":
			// Start import
			selectedRows := m.getSelectedRowNumbers()
			return m, func() tea.Msg {
				return ImportReviewMsg{
					SelectedRows: selectedRows,
					Action:       "import",
				}
			}

		case "escape", "q":
			// Cancel import
			return m, func() tea.Msg {
				return ImportReviewMsg{
					Action: "cancel",
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

// View renders the import review screen
func (m *ImportReviewModel) View() string {
	title := styles.HeaderMain.Render("Import Career Events Review")

	// Render summary
	summaryText := fmt.Sprintf(
		"Total: %d | Valid: %d | Invalid: %d | Duplicates: %d",
		m.summary["total"],
		m.summary["valid"],
		m.summary["invalid"],
		m.summary["duplicate"],
	)
	summaryBox := styles.InfoBox.Width(styles.MaxWidth(m.width) - 4).Render(summaryText)

	// Render rows
	rowsContent := m.renderRows()

	// Render help text
	helpText := m.renderHelp()

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		summaryBox,
		"",
		rowsContent,
		"",
		helpText,
	)

	card := styles.ResponsiveCard(m.width).Render(content)
	return styles.Center(card, m.width, m.height)
}

// renderRows renders the list of parsed rows
func (m *ImportReviewModel) renderRows() string {
	var rows []string

	for i, parsedRow := range m.ParsedRows {
		isSelected := m.selectedRows[parsedRow.RowNumber]
		isFocused := i == m.focusedRowIndex

		// Build row indicator
		var indicator string
		if parsedRow.IsDuplicate {
			indicator = "⚠ DUP"
		} else if !parsedRow.IsValid {
			indicator = "✗ ERR"
		} else if isSelected {
			indicator = "✓ SEL"
		} else {
			indicator = "  "
		}

		// Truncate text
		text := parsedRow.RawData["Text"]
		if len(text) > 40 {
			text = text[:37] + "..."
		}

		// Build row string
		rowStr := fmt.Sprintf("%s [%d] %s", indicator, parsedRow.RowNumber, text)

		// Apply styling
		if isFocused {
			rowStr = styles.ListItemFocused.Render(rowStr)
		} else if parsedRow.IsDuplicate || !parsedRow.IsValid {
			rowStr = styles.ErrorText.Render(rowStr)
		} else {
			rowStr = styles.ListItem.Render(rowStr)
		}

		rows = append(rows, rowStr)

		// Show validation errors if any
		if len(parsedRow.ValidationErrors) > 0 && isFocused {
			errorText := strings.Join(parsedRow.ValidationErrors, "; ")
			errorLine := styles.ErrorBox.Render("Error: " + errorText)
			rows = append(rows, errorLine)
		}

		// Show duplicate info if row is duplicate and focused
		if parsedRow.IsDuplicate && isFocused {
			dupInfo := m.GetDuplicateInfo(i)
			dupLine := styles.InfoBox.Render("Duplicate: " + dupInfo)
			rows = append(rows, dupLine)
		}
	}

	return strings.Join(rows, "\n")
}
// renderHelp renders the help text
func (m *ImportReviewModel) renderHelp() string {
	helpLines := []string{
		styles.InfoText.Render("↑/k") + " - Move up | " + styles.InfoText.Render("↓/j") + " - Move down",
		styles.InfoText.Render("space") + " - Toggle selection | " + styles.InfoText.Render("a") + " - Select all | " + styles.InfoText.Render("d") + " - Deselect all",
		styles.InfoText.Render("Enter") + " - Import selected | " + styles.InfoText.Render("Esc/q") + " - Cancel",
	}
	return strings.Join(helpLines, "\n")
}

// getSelectedRowNumbers returns the list of selected row numbers
func (m *ImportReviewModel) getSelectedRowNumbers() []int {
	var selected []int
	for rowNum := range m.selectedRows {
		selected = append(selected, rowNum)
	}
	return selected
}

// calculateSummary calculates the import summary
func calculateSummary(parsedRows []*importer.ParsedRow) map[string]int {
	summary := map[string]int{
		"total":     len(parsedRows),
		"valid":     0,
		"invalid":   0,
		"duplicate": 0,
	}

	for _, row := range parsedRows {
		if row.IsDuplicate {
			summary["duplicate"]++
		} else if row.IsValid {
			summary["valid"]++
		} else {
			summary["invalid"]++
		}
	}

	return summary
}

// ImportProgressModel shows the import progress
type ImportProgressModel struct {
	totalRows  int
	currentRow int
	Completed  bool
	result     *importer.ImportResult
	err        error
	width      int
	height     int
}

// NewImportProgressModel creates a new import progress model
func NewImportProgressModel(totalRows int) *ImportProgressModel {
	return &ImportProgressModel{
		totalRows:  totalRows,
		currentRow: 0,
		Completed:  false,
		width:      80,
		height:     24,
	}
}

// Update handles messages
func (m *ImportProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ImportResultMsg:
		m.result = msg.Result
		m.err = msg.Error
		m.Completed = true
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

// View renders the progress screen
func (m *ImportProgressModel) View() string {
	title := styles.HeaderMain.Render("Importing Career Events")

	if m.Completed {
		return m.renderResult()
	}

	// Render progress bar
	progress := float64(m.currentRow) / float64(m.totalRows)
	barWidth := 40
	filledWidth := int(float64(barWidth) * progress)
	bar := "[" + strings.Repeat("=", filledWidth) + strings.Repeat(" ", barWidth-filledWidth) + "]"

	statusText := fmt.Sprintf("Importing: %d/%d events", m.currentRow, m.totalRows)
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		statusText,
		bar,
		"",
		"Processing...",
	)

	card := styles.ResponsiveCard(m.width).Render(content)
	return styles.Center(card, m.width, m.height)
}

// renderResult renders the import result
func (m *ImportProgressModel) renderResult() string {
	title := styles.HeaderMain.Render("Import Complete")

	if m.err != nil {
		errorBox := styles.ErrorBox.Render(fmt.Sprintf("Error: %v", m.err))
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			"",
			errorBox,
			"",
			styles.InfoText.Render("Press any key to continue..."),
		)
		card := styles.ResponsiveCard(m.width).Render(content)
		return styles.Center(card, m.width, m.height)
	}

	resultText := fmt.Sprintf(
		"✓ Success: %d | ⚠ Skipped: %d | ✗ Failed: %d",
		m.result.SuccessCount,
		m.result.SkippedCount,
		m.result.FailedCount,
	)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		resultText,
		"",
		styles.SuccessBox.Render(fmt.Sprintf("Imported %d events successfully", m.result.SuccessCount)),
		"",
		styles.InfoText.Render("Press any key to continue..."),
	)

	card := styles.ResponsiveCard(m.width).Render(content)
	return styles.Center(card, m.width, m.height)
}

// GetFieldOrigins returns which fields came from CSV vs defaults for a row
func (m *ImportReviewModel) GetFieldOrigins(rowIndex int) map[string]bool {
	if rowIndex < 0 || rowIndex >= len(m.ParsedRows) {
		return make(map[string]bool)
	}

	row := m.ParsedRows[rowIndex]
	origins := make(map[string]bool) // true = from CSV, false = default

	// Check which fields are in RawData (came from CSV)
	for _, field := range []string{"text", "date", "company", "project", "tags", "categories"} {
		_, inRaw := row.RawData[field]
		origins[field] = inRaw
	}

	return origins
}

// HasDefaultFields checks if a row has any default (non-CSV) fields
func (m *ImportReviewModel) HasDefaultFields(rowIndex int) bool {
	origins := m.GetFieldOrigins(rowIndex)
	for _, fromCSV := range origins {
		if !fromCSV {
			return true
		}
	}
	return false
}

// GetDefaultFieldsList returns a list of fields that are defaults for a row
func (m *ImportReviewModel) GetDefaultFieldsList(rowIndex int) []string {
	origins := m.GetFieldOrigins(rowIndex)
	var defaults []string
	for field, fromCSV := range origins {
		if !fromCSV {
			defaults = append(defaults, field)
		}
	}
	return defaults
}

// GetParsingWarnings returns parsing-related warnings for a row
func (m *ImportReviewModel) GetParsingWarnings(rowIndex int) []string {
	if rowIndex < 0 || rowIndex >= len(m.ParsedRows) {
		return []string{}
	}

	row := m.ParsedRows[rowIndex]
	var warnings []string

	// Check for missing required fields
	if row.Event.Text == "" {
		warnings = append(warnings, "missing text")
	}
	if row.Event.Date.IsZero() {
		warnings = append(warnings, "missing date")
	}

	// Check for fields that are defaults (not in CSV)
	defaults := m.GetDefaultFieldsList(rowIndex)
	if len(defaults) > 0 {
		warnings = append(warnings, fmt.Sprintf("using defaults for: %s", strings.Join(defaults, ", ")))
	}

	return warnings
}

// GetDuplicateInfo returns duplicate detection information for a row
func (m *ImportReviewModel) GetDuplicateInfo(rowIndex int) string {
	if rowIndex < 0 || rowIndex >= len(m.ParsedRows) {
		return ""
	}

	row := m.ParsedRows[rowIndex]
	if row.IsDuplicate {
		return fmt.Sprintf("duplicate of event #%s", row.DuplicateOf)
	}
	return ""
}

// HasWarnings checks if a row has parsing warnings
func (m *ImportReviewModel) HasWarnings(rowIndex int) bool {
	return len(m.GetParsingWarnings(rowIndex)) > 0
}

// GetSelectedEvents returns the selected events for bulk operations
func (m *ImportReviewModel) GetSelectedEvents() []*career.CareerEvent {
	var selectedEvents []*career.CareerEvent
	for _, row := range m.ParsedRows {
		if m.selectedRows[row.RowNumber] {
			selectedEvents = append(selectedEvents, row.Event)
		}
	}
	return selectedEvents
}

// GetSelectedEventIDs returns the IDs of selected events
func (m *ImportReviewModel) GetSelectedEventIDs() []string {
	var ids []string
	for _, event := range m.GetSelectedEvents() {
		if event.ID != "" {
			ids = append(ids, event.ID)
		}
	}
	return ids
}

// CanBulkEdit checks if bulk operations are available
func (m *ImportReviewModel) CanBulkEdit() bool {
	return len(m.selectedRows) > 1
}

// ApplyBulkMetadataUpdate applies bulk metadata changes to selected events
func (m *ImportReviewModel) ApplyBulkMetadataUpdate(company, project string, tags, categories []string) int {
	count := 0
	for _, row := range m.ParsedRows {
		if m.selectedRows[row.RowNumber] {
			if company != "" {
				row.Event.Company = company
			}
			if project != "" {
				row.Event.Project = project
			}
			if len(tags) > 0 {
				row.Event.Tags = tags
			}
			if len(categories) > 0 {
				row.Event.Categories = categories
			}
			count++
		}
	}
	return count
}

// Init implements tea.Model
func (m *ImportReviewModel) Init() tea.Cmd {
	return nil
}

// Init implements tea.Model for ImportProgressModel
func (m *ImportProgressModel) Init() tea.Cmd {
	return nil
}

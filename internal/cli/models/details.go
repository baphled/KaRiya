package models

import (
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// DetailsModel represents the event details screen
type DetailsModel struct {
	*BaseStandardModel
	event      *career.CareerEvent
	width      int
	height     int
	header     components.HeaderModel
	helpFooter components.HelpFooterModel
}

// NewDetailsModel creates a new details model
func NewDetailsModel(event *career.CareerEvent) *DetailsModel {
	return &DetailsModel{
		BaseStandardModel: NewBaseStandardModel(),
		event:             event,
		width:             80,
		height:            24,
		header:            components.NewHeader("Event Details", 80),
		helpFooter:        components.NewHelpFooter("details", 80),
	}
}

// Init initializes the model
func (m *DetailsModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m *DetailsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.header.SetWidth(msg.Width)
		m.helpFooter.SetWidth(msg.Width)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "backspace":
			// Signal back navigation to parent
			return m, func() tea.Msg { return BackMsg{} }
		case "ctrl+c", "q":
			// Signal quit to parent
			return m, func() tea.Msg { return QuitMsg{} }
		}
	}
	return m, nil
}

// renderDescriptionSection renders the event description using SectionContainer
func (m *DetailsModel) renderDescriptionSection() string {
	content := lipgloss.NewStyle().Foreground(styles.ColorTextPrimary).Render(m.event.Text)
	return components.NewSectionContainer(content).
		SetTitle("Description").
		Render()
}

// renderMetadataSection renders event metadata (date, company, project) using SectionContainer
func (m *DetailsModel) renderMetadataSection() string {
	var metaLines []string

	metaLines = append(metaLines,
		styles.InputLabel.Render("Date:")+" "+lipgloss.NewStyle().Foreground(styles.ColorTextPrimary).Render(m.event.Date.Format("2006-01-02")),
	)

	if m.event.Company != "" {
		metaLines = append(metaLines,
			styles.InputLabel.Render("Company:")+" "+lipgloss.NewStyle().Foreground(styles.ColorTextPrimary).Render(m.event.Company),
		)
	}

	if m.event.Project != "" {
		metaLines = append(metaLines,
			styles.InputLabel.Render("Project:")+" "+lipgloss.NewStyle().Foreground(styles.ColorTextPrimary).Render(m.event.Project),
		)
	}

	content := strings.Join(metaLines, "\n")
	return components.NewSectionContainer(content).
		SetTitle("Metadata").
		Render()
}

// renderTagsSection renders tags using SectionContainer
func (m *DetailsModel) renderTagsSection() string {
	if len(m.event.Tags) == 0 {
		return ""
	}

	var tagContent strings.Builder
	for _, tag := range m.event.Tags {
		tagContent.WriteString(styles.TagBase.Render(tag))
		tagContent.WriteString(" ")
	}

	return components.NewSectionContainer(tagContent.String()).
		SetTitle("Tags").
		Render()
}

// renderCategoriesSection renders categories using SectionContainer
func (m *DetailsModel) renderCategoriesSection() string {
	if len(m.event.Categories) == 0 {
		return ""
	}

	var categoryContent strings.Builder
	for _, category := range m.event.Categories {
		categoryContent.WriteString(styles.TagBase.Render(category))
		categoryContent.WriteString(" ")
	}

	return components.NewSectionContainer(categoryContent.String()).
		SetTitle("Categories").
		Render()
}

// renderSystemInfoSection renders system information using SectionContainer
func (m *DetailsModel) renderSystemInfoSection() string {
	var sysLines []string

	sysLines = append(sysLines,
		styles.InputLabel.Render("Event ID:")+" "+lipgloss.NewStyle().Foreground(styles.ColorTextMuted).Render(m.event.ID),
	)

	sysLines = append(sysLines,
		styles.InputLabel.Render("Created:")+" "+lipgloss.NewStyle().Foreground(styles.ColorTextMuted).Render(m.formatTime(m.event.CreatedAt)),
	)

	sysLines = append(sysLines,
		styles.InputLabel.Render("Last Updated:")+" "+lipgloss.NewStyle().Foreground(styles.ColorTextMuted).Render(m.formatTime(m.event.UpdatedAt)),
	)

	content := strings.Join(sysLines, "\n")
	return components.NewSectionContainer(content).
		SetTitle("System Information").
		Render()
}

// renderDetailsWithContainers renders all details sections using SectionContainers
func (m *DetailsModel) renderDetailsWithContainers() string {
	var sections []string

	// Add description section
	sections = append(sections, m.renderDescriptionSection())

	// Add metadata section
	sections = append(sections, m.renderMetadataSection())

	// Add tags section if present
	if tagsSection := m.renderTagsSection(); tagsSection != "" {
		sections = append(sections, tagsSection)
	}

	// Add categories section if present
	if categoriesSection := m.renderCategoriesSection(); categoriesSection != "" {
		sections = append(sections, categoriesSection)
	}

	// Add system info section
	sections = append(sections, m.renderSystemInfoSection())

	return strings.Join(sections, "\n\n")
}

// View renders the details view
func (m *DetailsModel) View() string {
	if m.event == nil {
		return styles.ErrorBox.Render("No event selected\n\nPress 'esc' to return")
	}

	// Header
	headerContent := m.header.View()

	// Render details using SectionContainers
	detailsContent := m.renderDetailsWithContainers()

	// Help footer
	m.helpFooter.SetWidth(styles.MaxWidth(m.width))
	helpFooterContent := m.helpFooter.View()

	// Combine all content
	fullContent := lipgloss.JoinVertical(
		lipgloss.Left,
		headerContent,
		"",
		detailsContent,
		"",
		helpFooterContent,
	)

	// Wrap in card
	card := styles.CardBase.
		Width(styles.MaxWidth(m.width) - 4).
		Render(fullContent)

	return card
}

// formatTime formats a timestamp for display
func (m *DetailsModel) formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// Event returns the event being displayed
func (m *DetailsModel) Event() *career.CareerEvent {
	return m.event
}

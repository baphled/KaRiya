package models

import (
	"context"
	"errors"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/styles"
	careerdom "github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	cvservice "github.com/baphled/kariya/internal/service/career/cv"
)

type SourceEventTracerModel struct {
	*BaseStandardModel
	bullet              *careerdom.CVBullet
	traceabilityService *cvservice.TraceabilityService
	eventRepository     careerrepo.Repository
	factRepository      careerrepo.FactRepository
	sourceEvents        []*careerdom.CareerEvent
	sourceFacts         []*careerdom.Fact
	loading             bool
	width               int
	height              int
	scrollOffset        int
	header              components.HeaderModel
	helpFooter          components.HelpFooterModel
	onBack              func()
}

// NewSourceEventTracerModel creates a new source event tracer model.
func NewSourceEventTracerModel(
	baseModel *BaseStandardModel,
	bullet *careerdom.CVBullet,
	traceabilityService *cvservice.TraceabilityService,
	eventRepository careerrepo.Repository,
	factRepository careerrepo.FactRepository,
) *SourceEventTracerModel {
	return &SourceEventTracerModel{
		BaseStandardModel:   baseModel,
		bullet:              bullet,
		traceabilityService: traceabilityService,
		eventRepository:     eventRepository,
		factRepository:      factRepository,
		sourceEvents:        make([]*careerdom.CareerEvent, 0),
		sourceFacts:         make([]*careerdom.Fact, 0),
		loading:             true,
		width:               80,
		height:              20,
		scrollOffset:        0,
		header:              components.NewHeader("Bullet Source Traceability", 80),
		helpFooter:          components.NewHelpFooter("source_event_tracer", 80),
		onBack:              nil,
	}
}

// Init initializes the source event tracer model.
func (m *SourceEventTracerModel) Init() tea.Cmd {
	return tea.Batch(
		m.BaseStandardModel.Init(),
		m.loadSources(),
	)
}

// loadSources loads the source events and facts for the bullet.
func (m *SourceEventTracerModel) loadSources() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Load source events
		events := make([]*careerdom.CareerEvent, 0)
		if m.eventRepository != nil {
			for _, eventID := range m.bullet.SourceEventIDs {
				event, err := m.eventRepository.GetByID(ctx, eventID)
				if err != nil {
					return SourcesLoadErrorMsg{Err: fmt.Sprintf("failed to load event: %v", err)}
				}
				if event != nil {
					events = append(events, event)
				}
			}
		}

		// Load source facts
		facts := make([]*careerdom.Fact, 0)
		if m.factRepository != nil {
			for _, factID := range m.bullet.SourceFactIDs {
				fact, err := m.factRepository.GetByID(ctx, factID)
				if err != nil {
					return SourcesLoadErrorMsg{Err: fmt.Sprintf("failed to load fact: %v", err)}
				}
				if fact != nil {
					facts = append(facts, fact)
				}
			}
		}

		return SourcesLoadedMsg{Events: events, Facts: facts}
	}
}

// Update handles messages and updates the model state.
func (m *SourceEventTracerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case SourcesLoadedMsg:
		m.loading = false
		m.sourceEvents = msg.Events
		m.sourceFacts = msg.Facts
		return m, nil

	case SourcesLoadErrorMsg:
		m.loading = false
		m.SetError(errors.New(msg.Err))
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.header.SetWidth(msg.Width)
		m.helpFooter.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEscape:
			if m.onBack != nil {
				m.onBack()
			}
			return m, func() tea.Msg {
				return BackMsg{}
			}
		}

		if m.loading {
			return m, nil
		}

		maxScroll := len(m.sourceEvents) + len(m.sourceFacts) - (m.height - 6)
		if maxScroll < 0 {
			maxScroll = 0
		}

		switch msg.Type {
		case tea.KeyUp:
			if m.scrollOffset > 0 {
				m.scrollOffset--
			}
			return m, nil

		case tea.KeyDown:
			if m.scrollOffset < maxScroll {
				m.scrollOffset++
			}
			return m, nil

		case tea.KeyRunes:
			// Handle letter keys for navigation
			if len(msg.Runes) > 0 {
				switch msg.Runes[0] {
				case 'j':
					if m.scrollOffset < maxScroll {
						m.scrollOffset++
					}
					return m, nil
				case 'k':
					if m.scrollOffset > 0 {
						m.scrollOffset--
					}
					return m, nil
				}
			}
		}
	}

	// Return model unchanged for unhandled messages
	return m, nil
}

// View renders the source event tracer screen.
func (m *SourceEventTracerModel) View() string {
	headerContent := m.header.View()

	if m.loading {
		content := lipgloss.NewStyle().
			Padding(1, 2).
			Render("Loading source events and facts...")
		return lipgloss.JoinVertical(
			lipgloss.Top,
			headerContent,
			"",
			content,
		)
	}

	var lines []string

	// Add bullet information
	bulletSection := m.renderBulletSection()
	lines = append(lines, bulletSection)
	lines = append(lines, "")

	// Add source events
	if len(m.sourceEvents) > 0 {
		eventsSection := m.renderEventsSection()
		lines = append(lines, eventsSection)
		lines = append(lines, "")
	} else {
		lines = append(lines, styles.ListItem.Render("No source events"))
		lines = append(lines, "")
	}

	// Add source facts
	if len(m.sourceFacts) > 0 {
		factsSection := m.renderFactsSection()
		lines = append(lines, factsSection)
	} else {
		lines = append(lines, styles.ListItem.Render("No source facts"))
	}

	// Apply scroll offset
	content := lipgloss.JoinVertical(lipgloss.Left, lines...)

	// Add padding
	contentBox := lipgloss.NewStyle().
		Padding(1, 2).
		Render(content)

	footerContent := m.helpFooter.View()

	return lipgloss.JoinVertical(
		lipgloss.Top,
		headerContent,
		"",
		contentBox,
		"",
		footerContent,
	)
}

// renderBulletSection renders the bullet information.
func (m *SourceEventTracerModel) renderBulletSection() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("6"))

	title := titleStyle.Render("Bullet:")
	content := fmt.Sprintf("%s\n%s", title, m.bullet.Text)

	// Add metadata
	metadata := fmt.Sprintf(
		"Rank: %.2f | Confidence: %.2f | Reason: %s",
		m.bullet.Rank,
		m.bullet.Confidence,
		m.bullet.InclusionReason,
	)
	metadataStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	return lipgloss.JoinVertical(
		lipgloss.Left,
		content,
		metadataStyle.Render(metadata),
	)
}

// renderEventsSection renders the source events.
func (m *SourceEventTracerModel) renderEventsSection() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("2"))

	title := titleStyle.Render(fmt.Sprintf("Source Events (%d):", len(m.sourceEvents)))

	var eventLines []string
	for i, event := range m.sourceEvents {
		eventLine := fmt.Sprintf(
			"[%d] %s (%s) - %s",
			i+1,
			event.Date.Format("2006-01-02"),
			event.Company,
			event.Text,
		)

		if len(eventLine) > m.width-4 {
			eventLine = eventLine[:m.width-7] + "..."
		}

		eventLines = append(eventLines, eventLine)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		append([]string{title}, eventLines...)...,
	)
}

// renderFactsSection renders the source facts.
func (m *SourceEventTracerModel) renderFactsSection() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("3"))

	title := titleStyle.Render(fmt.Sprintf("Source Facts (%d):", len(m.sourceFacts)))

	var factLines []string
	for i, fact := range m.sourceFacts {
		factLine := fmt.Sprintf(
			"[%d] %s - %s",
			i+1,
			fact.StrengthSignal,
			fact.Text,
		)

		if len(factLine) > m.width-4 {
			factLine = factLine[:m.width-7] + "..."
		}

		factLines = append(factLines, factLine)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		append([]string{title}, factLines...)...,
	)
}

// GetBullet returns the bullet.
func (m *SourceEventTracerModel) GetBullet() *careerdom.CVBullet {
	return m.bullet
}

// GetSourceEvents returns the loaded source events.
func (m *SourceEventTracerModel) GetSourceEvents() []*careerdom.CareerEvent {
	return m.sourceEvents
}

// GetSourceFacts returns the loaded source facts.
func (m *SourceEventTracerModel) GetSourceFacts() []*careerdom.Fact {
	return m.sourceFacts
}

// GetWidth returns the model width.
func (m *SourceEventTracerModel) GetWidth() int {
	return m.width
}

// GetHeight returns the model height.
func (m *SourceEventTracerModel) GetHeight() int {
	return m.height
}

// GetScrollOffset returns the current scroll offset.
func (m *SourceEventTracerModel) GetScrollOffset() int {
	return m.scrollOffset
}

// SetLoading sets the loading state.
func (m *SourceEventTracerModel) SetLoading(loading bool) {
	m.loading = loading
}

// SetOnBack sets the callback for back navigation.
func (m *SourceEventTracerModel) SetOnBack(fn func()) {
	m.onBack = fn
}

// SourcesLoadedMsg indicates sources have been loaded.
type SourcesLoadedMsg struct {
	Events []*careerdom.CareerEvent
	Facts  []*careerdom.Fact
}

// SourcesLoadErrorMsg indicates an error loading sources.
type SourcesLoadErrorMsg struct {
	Err string
}

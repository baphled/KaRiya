package event

import (
	"strings"

	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
	tea "github.com/charmbracelet/bubbletea"
)

// BurstSuggestionActionResult holds data for a single suggestion action (confirm/reject).
type BurstSuggestionActionResult struct {
	Action     string
	Suggestion display.BurstSuggestion
	Index      int
}

// BurstSuggestionCompleteResult holds data when all suggestions are processed.
type BurstSuggestionCompleteResult struct {
	Action         string
	ConfirmedCount int
	RejectedCount  int
}

// BurstSuggestion is a presentational-only View for reviewing burst suggestions.
type BurstSuggestion struct {
	widgets.BaseView
	suggestions   []display.BurstSuggestion
	currentIdx    int
	confirmed     []display.BurstSuggestion
	rejected      []display.BurstSuggestion
	editing       bool
	relatedEvents map[int][]display.Event
	editedNames   map[int]string
	editedDescs   map[int]string
	width         int
	height        int
}

// NewBurstSuggestion creates a BurstSuggestion view for the given suggestions.
//
// Expected:
//   - burstsuggestion must be valid.
//
// Returns:
//   - A fully initialized BurstSuggestion ready for use.
//
// Side effects:
//   - None.
func NewBurstSuggestion(suggestions []display.BurstSuggestion) *BurstSuggestion {
	return &BurstSuggestion{
		BaseView:      widgets.BaseView{},
		suggestions:   suggestions,
		currentIdx:    0,
		confirmed:     []display.BurstSuggestion{},
		rejected:      []display.BurstSuggestion{},
		relatedEvents: map[int][]display.Event{},
		editedNames:   map[int]string{},
		editedDescs:   map[int]string{},
		width:         80,
		height:        24,
	}
}

// Init returns nil (no async initialisation needed).
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (m *BurstSuggestion) Init() tea.Cmd {
	return nil
}

// Update handles messages and returns a ViewResult on user action.
//
// Expected:
//   - msg is a valid tea.Msg (key press, window resize, or custom message).
//
// Returns:
//   - A tea.Cmd and widgets.ViewResult pair.
//
// Side effects:
//   - Updates internal view state based on message type.
func (m *BurstSuggestion) Update(msg tea.Msg) (tea.Cmd, widgets.ViewResult) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return nil, nil
}

func (m *BurstSuggestion) handleKeyMsg(msg tea.KeyMsg) (tea.Cmd, widgets.ViewResult) {
	if len(m.suggestions) == 0 {
		return nil, nil
	}
	key := msg.String()
	if m.editing {
		if key == "e" {
			m.editing = false
		}
		return nil, nil
	}

	switch key {
	case "down", "j":
		if m.currentIdx < len(m.suggestions)-1 {
			m.currentIdx++
		}
	case "up", "k":
		if m.currentIdx > 0 {
			m.currentIdx--
		}
	case "y":
		m.confirmed = append(m.confirmed, m.suggestions[m.currentIdx])
		return nil, m.processSuggestionAction("confirm", m.currentIdx)
	case "n":
		m.rejected = append(m.rejected, m.suggestions[m.currentIdx])
		return nil, m.processSuggestionAction("reject", m.currentIdx)
	case "esc":
		return nil, &widgets.CancelViewResult{}
	case "e":
		m.editing = true
	}
	return nil, nil
}

func (m *BurstSuggestion) processSuggestionAction(
	action string,
	idx int,
) widgets.ViewResult {
	if idx == len(m.suggestions)-1 {
		return m.completeResult()
	}
	m.currentIdx++
	if action == "confirm" {
		return &widgets.SubmitViewResult{
			FormData: BurstSuggestionActionResult{
				Action:     action,
				Suggestion: m.suggestions[idx],
				Index:      idx,
			},
		}
	}
	return &widgets.NavigateViewResult{
		ResultData: BurstSuggestionActionResult{
			Action:     action,
			Suggestion: m.suggestions[idx],
			Index:      idx,
		},
	}
}

func (m *BurstSuggestion) completeResult() *widgets.SubmitViewResult {
	return &widgets.SubmitViewResult{
		FormData: BurstSuggestionCompleteResult{
			Action:         "complete",
			ConfirmedCount: len(m.confirmed),
			RejectedCount:  len(m.rejected),
		},
	}
}

// RenderContent returns the rendered suggestion content.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *BurstSuggestion) RenderContent() string {
	if len(m.suggestions) == 0 {
		return "No burst suggestions available"
	}
	if m.editing {
		return m.renderEditView()
	}
	return m.renderReviewView()
}

// HelpText returns key binding instructions.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *BurstSuggestion) HelpText() string {
	return "y: confirm | n: reject | e: edit | esc: cancel | j/k/up/down: navigate"
}

// SetRelatedEvents associates related events with a suggestion index.
//
// Expected:
//   - int must be valid.
//   - event must be valid.
//
// Side effects:
//   - None.
func (m *BurstSuggestion) SetRelatedEvents(idx int, events []display.Event) {
	m.relatedEvents[idx] = events
}

// GetConfirmed returns the confirmed suggestions.
//
// Returns:
//   - A []display.BurstSuggestion value.
//
// Side effects:
//   - None.
func (m *BurstSuggestion) GetConfirmed() []display.BurstSuggestion {
	return m.confirmed
}

// GetRejected returns the rejected suggestions.
//
// Returns:
//   - A []display.BurstSuggestion value.
//
// Side effects:
//   - None.
func (m *BurstSuggestion) GetRejected() []display.BurstSuggestion {
	return m.rejected
}

// IsDone returns true when all suggestions have been processed.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *BurstSuggestion) IsDone() bool {
	return len(m.confirmed)+len(m.rejected) >= len(m.suggestions)
}

// IsEditing returns whether the view is in edit mode.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *BurstSuggestion) IsEditing() bool {
	return m.editing
}

// CurrentIdx returns the current suggestion index.
//
// Returns:
//   - A int value.
//
// Side effects:
//   - None.
func (m *BurstSuggestion) CurrentIdx() int {
	return m.currentIdx
}

// Width returns the view width.
//
// Returns:
//   - A int value.
//
// Side effects:
//   - None.
func (m *BurstSuggestion) Width() int {
	return m.width
}

// Height returns the view height.
//
// Returns:
//   - A int value.
//
// Side effects:
//   - None.
func (m *BurstSuggestion) Height() int {
	return m.height
}

// SetEditedName sets an edited name for the suggestion at the given index.
//
// Expected:
//   - int must be valid.
//   - Must be a valid string.
//
// Side effects:
//   - None.
func (m *BurstSuggestion) SetEditedName(idx int, name string) {
	m.editedNames[idx] = name
}

// SetEditedDescription sets an edited description for the suggestion at the given index.
//
// Expected:
//   - int must be valid.
//   - Must be a valid string.
//
// Side effects:
//   - None.
func (m *BurstSuggestion) SetEditedDescription(idx int, desc string) {
	m.editedDescs[idx] = desc
}

// ExitEditMode exits the editing state.
//
// Side effects:
//   - None.
func (m *BurstSuggestion) ExitEditMode() {
	m.editing = false
}

// Rendering methods.
func (m *BurstSuggestion) renderReviewView() string {
	out := "Review: " + m.suggestions[m.currentIdx].Name
	related := m.renderRelatedEvents()
	if related != "" {
		out += "\n" + related
	}
	return out
}

func (m *BurstSuggestion) renderEditView() string {
	return "Edit: " + m.suggestions[m.currentIdx].Name
}

func (m *BurstSuggestion) renderRelatedEvents() string {
	events := m.relatedEvents[m.currentIdx]
	if len(events) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("Related Events: ")
	for idx := range events {
		ev := events[idx]
		b.WriteString(ev.Text)
		b.WriteByte(' ')
	}
	return b.String()
}

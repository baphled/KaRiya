package event

import (
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/ui/types"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
	tea "github.com/charmbracelet/bubbletea"
)

// StrategySelect allows users to choose between Quick and Manual capture strategies.
type StrategySelect struct {
	widgets.BaseView
	items         []types.CaptureStrategy
	renderer      func(types.CaptureStrategy) string
	title         string
	selectedIndex int
	scrollOffset  int
	visibleItems  int
}

// NewStrategySelect creates a new StrategySelect.
//
// Returns:
//   - A fully initialized StrategySelect ready for use.
//
// Side effects:
//   - None.
func NewStrategySelect() *StrategySelect {
	strategies := []types.CaptureStrategy{
		types.StrategyQuick,
		types.StrategyManual,
	}

	renderer := func(strategy types.CaptureStrategy) string {
		var label, description string
		switch strategy {
		case types.StrategyQuick:
			label = "Quick"
			description = "Capture with minimal fields (event text only)"
		case types.StrategyManual:
			label = "Manual"
			description = "Full form with optional fields (date, company, project, tags)"
		default:
			label = string(strategy)
			description = "Unknown strategy"
		}
		return fmt.Sprintf("%s - %s", label, description)
	}

	view := &StrategySelect{
		items:         strategies,
		renderer:      renderer,
		title:         "Select Capture Strategy",
		selectedIndex: 0,
		scrollOffset:  0,
		visibleItems:  10,
	}
	view.SetTerminalInfo(120, 40)
	view.updateVisibleItems()
	return view
}

// WithInitialSelection sets the initial selection index.
//
// Expected:
//   - int must be valid.
//
// Returns:
//   - A fully initialized StrategySelect ready for use.
//
// Side effects:
//   - None.
func (v *StrategySelect) WithInitialSelection(index int) *StrategySelect {
	if index < 0 {
		index = 0
	}
	if index >= len(v.items) {
		index = len(v.items) - 1
	}
	if index < 0 {
		index = 0
	}
	v.selectedIndex = index
	v.updateScrollOffset()
	return v
}

// Init returns nil (no async loading needed).
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (v *StrategySelect) Init() tea.Cmd {
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
func (v *StrategySelect) Update(msg tea.Msg) (tea.Cmd, widgets.ViewResult) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		v.SetTerminalInfo(msg.Width, msg.Height)
		v.updateVisibleItems()
		return nil, nil
	case tea.KeyMsg:
		return v.handleKeyMessage(msg)
	}
	return nil, nil
}

// RenderContent returns the list content string — no chrome.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (v *StrategySelect) RenderContent() string {
	if len(v.items) == 0 {
		return "\n  No items available\n"
	}

	var b strings.Builder

	if v.title != "" {
		b.WriteString("\n  ")
		b.WriteString(v.title)
		b.WriteString("\n\n")
	}

	start := v.scrollOffset
	end := v.scrollOffset + v.visibleItems
	if end > len(v.items) {
		end = len(v.items)
	}

	for i := start; i < end; i++ {
		prefix := "  "
		if i == v.selectedIndex {
			prefix = "▶ "
		}

		itemText := v.renderer(v.items[i])
		b.WriteString(prefix)
		b.WriteString(itemText)
		b.WriteString("\n")
	}

	if v.scrollOffset > 0 {
		b.WriteString("\n  ↑ More items above")
	}
	if end < len(v.items) {
		b.WriteString("\n  ↓ More items below")
	}

	return b.String()
}

// HelpText returns the footer help text for this view.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (v *StrategySelect) HelpText() string {
	if len(v.items) == 0 {
		return "Esc: Back  q: Quit"
	}
	return "↑/↓/j/k: Navigate  g/G: Jump  Enter: Select  Esc: Back  q: Quit"
}

func (v *StrategySelect) handleKeyMessage(msg tea.KeyMsg) (tea.Cmd, widgets.ViewResult) {
	switch msg.String() {
	case "down", "j":
		v.navigateDown()
		return nil, nil
	case "up", "k":
		v.navigateUp()
		return nil, nil
	case "g":
		v.jumpToTop()
		return nil, nil
	case "G":
		v.jumpToBottom()
		return nil, nil
	case "enter":
		return nil, v.handleSelection()
	case "esc":
		return nil, v.handleCancellation()
	}

	return nil, nil
}

func (v *StrategySelect) navigateDown() {
	if len(v.items) == 0 {
		return
	}
	if v.selectedIndex < len(v.items)-1 {
		v.selectedIndex++
		v.updateScrollOffset()
	}
}

func (v *StrategySelect) navigateUp() {
	if len(v.items) == 0 {
		return
	}
	if v.selectedIndex > 0 {
		v.selectedIndex--
		v.updateScrollOffset()
	}
}

func (v *StrategySelect) jumpToTop() {
	if len(v.items) == 0 {
		return
	}
	v.selectedIndex = 0
	v.scrollOffset = 0
}

func (v *StrategySelect) jumpToBottom() {
	if len(v.items) == 0 {
		return
	}
	v.selectedIndex = len(v.items) - 1
	v.updateScrollOffset()
}

func (v *StrategySelect) updateScrollOffset() {
	if v.selectedIndex < v.scrollOffset {
		v.scrollOffset = v.selectedIndex
	}
	if v.selectedIndex >= v.scrollOffset+v.visibleItems {
		v.scrollOffset = v.selectedIndex - v.visibleItems + 1
	}
}

func (v *StrategySelect) updateVisibleItems() {
	availableHeight := v.GetTerminalHeight() - 15
	if availableHeight < 5 {
		availableHeight = 5
	}
	v.visibleItems = availableHeight
}

func (v *StrategySelect) handleSelection() widgets.ViewResult {
	if len(v.items) == 0 {
		return nil
	}
	selectedItem := v.items[v.selectedIndex]
	result := &widgets.NavigateViewResult{ResultData: selectedItem}
	result.WithMetadata("selected_index", v.selectedIndex)
	return result
}

func (v *StrategySelect) handleCancellation() widgets.ViewResult {
	result := &widgets.CancelViewResult{}
	result.WithMetadata("selected_index", v.selectedIndex)
	return result
}

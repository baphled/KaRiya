package capture_event

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/models"
	captureScreens "github.com/baphled/kariya/internal/cli/screens/capture"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/domain/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// setCompleted marks the intent as completed with a result.
func (i *Intent) setCompleted(result *Result) {
	i.result = intents.NewCompletedResult(result)
	i.active = false
}

// setCancelled marks the intent as cancelled by the user.
func (i *Intent) setCancelled() {
	i.result = intents.NewCancelledResult[*Result]()
	i.active = false
}

// setFailed marks the intent as failed with an error.
func (i *Intent) setFailed(code, message string, cause error) {
	i.result = intents.NewFailedResult[*Result](code, message, cause)
	i.active = false
}

// setFailedCmd is a helper to set failed result and return nil command.
func (i *Intent) setFailedCmd(code, message string, cause error) tea.Cmd {
	i.setFailed(code, message, cause)
	return nil
}

// showSubmitModal creates and shows the loading modal for event submission.
func (i *Intent) showSubmitModal() tea.Cmd {
	i.state.submitModal = feedback.NewLoadingModal("Saving event...", false)
	return tea.Batch(i.performSubmit(), i.state.submitModal.Init())
}

// performSubmit performs the actual submission of the event.
// It calls the domain service to persist the event to the database and optionally enriches it.
func (i *Intent) performSubmit() tea.Cmd {
	event := i.state.reviewState.Event
	acceptedFacts := i.state.reviewState.AcceptedFacts
	strategy := i.state.strategy
	careerService := i.state.context.CareerService
	eventService := i.eventService

	return func() tea.Msg {
		if event == nil {
			return SubmitErrorMsg{
				Code:    "MISSING_EVENT",
				Message: "No event data to submit",
				Cause:   nil,
			}
		}

		if err := event.Validate(); err != nil {
			return SubmitErrorMsg{
				Code:    "VALIDATION_ERROR",
				Message: fmt.Sprintf("Event validation failed: %v", err),
				Cause:   err,
			}
		}

		if eventService == nil {
			return SubmitErrorMsg{
				Code:    "SERVICE_ERROR",
				Message: "Event service not initialized",
				Cause:   nil,
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if strategy == StrategyQuick && event.Date.IsZero() {
			event.Date = time.Now()
		}

		mode := careerservice.ManualEntry

		if careerService == nil {
			return SubmitErrorMsg{
				Code:    "SERVICE_ERROR",
				Message: "Career service not initialized",
				Cause:   nil,
			}
		}

		err := careerService.CaptureEvent(ctx, event, mode)

		if err != nil {
			return SubmitErrorMsg{
				Code:    "PERSISTENCE_ERROR",
				Message: fmt.Sprintf("Failed to save event: %v", err),
				Cause:   err,
			}
		}

		if careerService != nil && len(acceptedFacts) > 0 {
			for _, fact := range acceptedFacts {
				if fact.ID == "" {
					if fact.SourceEventID == "" {
						fact.SourceEventID = event.ID
					}

					if err := careerService.SaveFact(ctx, fact); err != nil {
						continue
					}
				}
			}
		}

		return SubmitCompleteMsg{}
	}
}

// initializeFormForNew initializes the form for capturing a new event.
func (i *Intent) initializeFormForNew() tea.Cmd {
	i.state.reviewState.Event = &career.Event{
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Tags:       make([]string, 0),
		Categories: make([]string, 0),
	}
	return func() tea.Msg { return nil }
}

// transitionToStrategyScreen transitions to the strategy selection screen.
func (i *Intent) transitionToStrategyScreen() tea.Cmd {
	i.state.currentState = StateChooseStrategy
	breadcrumbs := []string{"Main Menu", "Capture Event"}
	i.activeScreen = captureScreens.NewStrategySelectScreen(breadcrumbs)

	termInfo := i.GetTerminalInfo()
	width, height := 120, 40
	if termInfo != nil {
		width = termInfo.Width
		height = termInfo.Height
	}

	i.activeScreen.SetTerminalInfo(width, height)
	i.activeScreen.SetTheme(i.Theme())
	i.activeScreen.SetLogo(i.GetLogo(), i.GetLogoSpacing())

	return nil
}

// transitionToFormScreen transitions to the event form screen.
func (i *Intent) transitionToFormScreen(strategy CaptureStrategy) tea.Cmd {
	i.state.currentState = StateForm
	i.state.strategy = strategy
	breadcrumbs := []string{"Main Menu", "Capture Event", "Form"}
	i.activeScreen = captureScreens.NewEventFormScreen(
		i.eventService,
		breadcrumbs,
		strategy,
	)

	termInfo := i.GetTerminalInfo()
	width, height := 120, 40
	if termInfo != nil {
		width = termInfo.Width
		height = termInfo.Height
	}

	i.activeScreen.SetTerminalInfo(width, height)
	i.activeScreen.SetTheme(i.Theme())
	i.activeScreen.SetLogo(i.GetLogo(), i.GetLogoSpacing())

	if initable, ok := i.activeScreen.(interface{ Init() tea.Cmd }); ok {
		return initable.Init()
	}
	return nil
}

// navigateReviewItems moves the selection cursor through inferred bursts and facts
// by the given delta (+1 for down, -1 for up). It wraps around at both ends and
// automatically updates the SelectedItemType based on the new index position.
func (i *Intent) navigateReviewItems(delta int) {
	totalItems := len(i.state.reviewState.InferredBursts) + len(i.state.reviewState.InferredFacts)
	if totalItems == 0 {
		return
	}

	i.state.reviewState.SelectedIndex += delta
	if i.state.reviewState.SelectedIndex >= totalItems {
		i.state.reviewState.SelectedIndex = 0
	}
	if i.state.reviewState.SelectedIndex < 0 {
		i.state.reviewState.SelectedIndex = totalItems - 1
	}

	if i.state.reviewState.SelectedIndex < len(i.state.reviewState.InferredBursts) {
		i.state.reviewState.SelectedItemType = "burst"
	} else {
		i.state.reviewState.SelectedItemType = "fact"
		i.state.reviewState.SelectedIndex -= len(i.state.reviewState.InferredBursts)
	}
}

// acceptCurrentItem accepts the currently selected burst or fact.
func (i *Intent) acceptCurrentItem() {
	if i.state.reviewState.SelectedItemType == "burst" {
		idx := i.state.reviewState.SelectedIndex
		if idx >= 0 && idx < len(i.state.reviewState.InferredBursts) {
			burst := i.state.reviewState.InferredBursts[idx]
			i.state.reviewState.AcceptedBursts = append(i.state.reviewState.AcceptedBursts, burst)
			i.state.reviewState.InferredBursts = append(
				i.state.reviewState.InferredBursts[:idx],
				i.state.reviewState.InferredBursts[idx+1:]...)
			if i.state.reviewState.SelectedIndex >= len(i.state.reviewState.InferredBursts) && len(i.state.reviewState.InferredBursts) > 0 {
				i.state.reviewState.SelectedIndex = len(i.state.reviewState.InferredBursts) - 1
			}
		}
	} else if i.state.reviewState.SelectedItemType == "fact" {
		idx := i.state.reviewState.SelectedIndex
		if idx >= 0 && idx < len(i.state.reviewState.InferredFacts) {
			fact := i.state.reviewState.InferredFacts[idx]
			i.state.reviewState.AcceptedFacts = append(i.state.reviewState.AcceptedFacts, fact)
			i.state.reviewState.InferredFacts = append(
				i.state.reviewState.InferredFacts[:idx],
				i.state.reviewState.InferredFacts[idx+1:]...)
			if i.state.reviewState.SelectedIndex >= len(i.state.reviewState.InferredFacts) && len(i.state.reviewState.InferredFacts) > 0 {
				i.state.reviewState.SelectedIndex = len(i.state.reviewState.InferredFacts) - 1
			}
		}
	}
}

// rejectCurrentItem rejects the currently selected burst or fact.
func (i *Intent) rejectCurrentItem() {
	if i.state.reviewState.SelectedItemType == "burst" {
		idx := i.state.reviewState.SelectedIndex
		if idx >= 0 && idx < len(i.state.reviewState.InferredBursts) {
			burst := i.state.reviewState.InferredBursts[idx]
			if i.state.reviewState.RejectedItems == nil {
				i.state.reviewState.RejectedItems = make(map[string]string)
			}
			i.state.reviewState.RejectedItems[burst.ID] = "user_rejected"
			i.state.reviewState.InferredBursts = append(
				i.state.reviewState.InferredBursts[:idx],
				i.state.reviewState.InferredBursts[idx+1:]...)
			if i.state.reviewState.SelectedIndex >= len(i.state.reviewState.InferredBursts) && len(i.state.reviewState.InferredBursts) > 0 {
				i.state.reviewState.SelectedIndex = len(i.state.reviewState.InferredBursts) - 1
			}
		}
	} else if i.state.reviewState.SelectedItemType == "fact" {
		idx := i.state.reviewState.SelectedIndex
		if idx >= 0 && idx < len(i.state.reviewState.InferredFacts) {
			fact := i.state.reviewState.InferredFacts[idx]
			if i.state.reviewState.RejectedItems == nil {
				i.state.reviewState.RejectedItems = make(map[string]string)
			}
			i.state.reviewState.RejectedItems[fact.ID] = "user_rejected"
			i.state.reviewState.InferredFacts = append(
				i.state.reviewState.InferredFacts[:idx],
				i.state.reviewState.InferredFacts[idx+1:]...)
			if i.state.reviewState.SelectedIndex >= len(i.state.reviewState.InferredFacts) && len(i.state.reviewState.InferredFacts) > 0 {
				i.state.reviewState.SelectedIndex = len(i.state.reviewState.InferredFacts) - 1
			}
		}
	}
}

// getTheme returns the theme or a default.
func (i *Intent) getTheme() themes.Theme {
	if theme := i.Theme(); theme != nil {
		return theme
	}
	return themes.NewDefaultTheme()
}

// getCardStyle returns the card base style from the current theme for content containers.
func (i *Intent) getCardStyle() lipgloss.Style {
	return i.getTheme().Styles().CardBase
}

// getPrimaryColor returns the primary text color from theme.
func (i *Intent) getPrimaryColor() lipgloss.Color {
	return i.getTheme().ForegroundColor()
}

// getAccentColor returns the accent color from theme.
func (i *Intent) getAccentColor() lipgloss.Color {
	return i.getTheme().PrimaryColor()
}

// getStateName returns a human-readable name for the current state.
func (i *Intent) getStateName() string {
	switch i.state.currentState {
	case StateChooseStrategy:
		return "Choose Strategy"
	case StateForm:
		return "Enter Details"
	case StateReview:
		return "Review"
	case StateSubmit:
		return "Submit"
	default:
		return string(i.state.currentState)
	}
}

// getStateContent returns the content for the current state.
func (i *Intent) getStateContent() string {
	if i.state.error != nil && !i.HasError() {
		return i.viewError()
	}

	switch i.state.currentState {
	case StateChooseStrategy:
		return i.viewChooseStrategy()
	case StateForm:
		return i.viewCaptureForm()
	case StateReview:
		return i.viewReviewInferredEvent()
	case StateSubmit:
		return i.viewSubmit()
	default:
		return fmt.Sprintf("Unknown state: %s", i.state.currentState)
	}
}

// getContextHelp returns context-aware help text for the current state.
func (i *Intent) getContextHelp() string {
	theme := i.Theme()

	switch i.state.currentState {
	case StateChooseStrategy:
		return intents.CombineThemedFooters(
			intents.ThemedNavigationFooter(theme),
			intents.ThemedGlobalBadges(theme),
		)
	case StateForm:
		if i.state.strategy == StrategyManual {
			return intents.CombineThemedFooters(
				intents.ThemedFormFooter(theme),
				intents.ThemedCustomFooter(theme,
					primitives.HelpKeyBadge("Ctrl+O", "Toggle fields", theme),
				),
				intents.ThemedGlobalBadges(theme),
			)
		}
		return intents.CombineThemedFooters(
			intents.ThemedFormFooter(theme),
			intents.ThemedGlobalBadges(theme),
		)
	case StateReview:
		if i.state.reviewState.EditingMode != EditingModeNone {
			return intents.ThemedCustomFooter(theme,
				primitives.HelpKeyBadge("Editing", "...", theme),
				primitives.CancelBadge(theme),
				primitives.SaveBadge(theme),
			)
		}
		return intents.CombineThemedFooters(
			intents.ThemedCustomFooter(theme,
				primitives.NavigateBadge(theme),
				primitives.EditBadge(theme),
				primitives.HelpKeyBadge("b", "Bursts", theme),
				primitives.HelpKeyBadge("f", "Facts", theme),
				primitives.HelpKeyBadge("a", "Accept", theme),
				primitives.HelpKeyBadge("r", "Reject", theme),
				primitives.BackBadge(theme),
			),
			intents.ThemedGlobalBadges(theme),
		)
	case StateSubmit:
		return intents.CombineThemedFooters(
			intents.ThemedCustomFooter(theme,
				primitives.HelpKeyBadge("Enter", "Continue", theme),
				primitives.BackBadge(theme),
			),
			intents.ThemedGlobalBadges(theme),
		)
	default:
		return intents.ThemedGlobalBadges(theme)
	}
}

// modalContentData holds the title, body content, and footer text used to render
// an editing modal overlay on top of the review screen.
type modalContentData struct {
	// title is the modal header text (e.g. "Edit Metadata").
	title string

	// content is the modal body rendered by the underlying form model.
	content string

	// footer is the help/shortcut text shown at the bottom of the modal.
	footer string
}

// getMetadataModalContent returns the modal content for metadata editing.
func (i *Intent) getMetadataModalContent() *modalContentData {
	if i.state.reviewState.metadataModal == nil {
		i.state.reviewState.metadataModal = models.NewMetadataEditorModelNew(
			i.state.reviewState.Event,
			i.context.CareerService,
			i.context.CLIEventService,
			context.Background(),
		)
	}
	return &modalContentData{
		title:   i.state.reviewState.metadataModal.GetTitle(),
		content: i.state.reviewState.metadataModal.GetContent(),
		footer:  i.state.reviewState.metadataModal.GetFooter(),
	}
}

// getBurstModalContent returns the modal content for burst editing.
func (i *Intent) getBurstModalContent() *modalContentData {
	if i.state.reviewState.burstModal == nil {
		return nil
	}
	return &modalContentData{
		title:   i.state.reviewState.burstModal.GetTitle(),
		content: i.state.reviewState.burstModal.GetContent(),
		footer:  i.state.reviewState.burstModal.GetFooter(),
	}
}

// getFactModalContent returns the modal content for fact editing.
func (i *Intent) getFactModalContent() *modalContentData {
	if i.state.reviewState.factModal == nil {
		return nil
	}
	return &modalContentData{
		title:   i.state.reviewState.factModal.GetTitle(),
		content: i.state.reviewState.factModal.GetContent(),
		footer:  i.state.reviewState.factModal.GetFooter(),
	}
}

// getEditingModalContent returns the modal content based on current editing mode.
func (i *Intent) getEditingModalContent() *modalContentData {
	if i.state.reviewState == nil {
		return nil
	}
	switch i.state.reviewState.EditingMode {
	case EditingModeMetadata:
		return i.getMetadataModalContent()
	case EditingModeBursts:
		return i.getBurstModalContent()
	case EditingModeFacts:
		return i.getFactModalContent()
	default:
		return nil
	}
}

// renderModalOverlay renders a modal overlay on top of the background content.
func (i *Intent) renderModalOverlay(background string, modalContent *modalContentData) string {
	if modalContent == nil {
		return background
	}

	info := i.GetTerminalInfo()
	width := 80
	height := 24
	if info != nil {
		width = info.Width
		height = info.Height
	}

	overlay := feedback.NewOverlayModal(modalContent.title, modalContent.content)
	overlay.SetFooter(modalContent.footer)
	overlay.SetWidth(80)

	return overlay.RenderCentered(background, width, height)
}

// overlayModal overlays modal content on top of background content (centered).
func (i *Intent) overlayModal(background, modal string, width, _ int) string {
	bgLines := strings.Split(background, "\n")
	modalLines := strings.Split(modal, "\n")

	bgHeight := len(bgLines)
	modalHeight := len(modalLines)
	startLine := (bgHeight - modalHeight) / 2
	if startLine < 0 {
		startLine = 0
	}

	result := make([]string, len(bgLines))
	copy(result, bgLines)

	for idx, modalLine := range modalLines {
		lineIndex := startLine + idx
		if lineIndex >= 0 && lineIndex < len(result) {
			centeredModalLine := lipgloss.PlaceHorizontal(width, lipgloss.Center, modalLine)
			result[lineIndex] = centeredModalLine
		}
	}

	return strings.Join(result, "\n")
}

// viewChooseStrategy renders the strategy selection UI.
func (i *Intent) viewChooseStrategy() string {
	var content strings.Builder
	content.WriteString("\n📝 Select Capture Strategy\n\n")

	strategies := []struct {
		value       CaptureStrategy
		label       string
		description string
	}{
		{StrategyQuick, "Quick", "Capture with minimal fields (event text only)"},
		{StrategyManual, "Manual", "Full form with optional fields (date, company, project, tags)"},
	}

	for idx, s := range strategies {
		prefix := "  "
		if idx == i.state.selectedStrategyIndex {
			prefix = "▶ "
		}

		optStyle := lipgloss.NewStyle().Foreground(i.getPrimaryColor())
		if idx == i.state.selectedStrategyIndex {
			optStyle = optStyle.Foreground(i.getAccentColor()).Bold(true)
		}

		line := fmt.Sprintf("%s%s - %s", prefix, s.label, s.description)
		content.WriteString(optStyle.Render(line) + "\n")
	}

	return i.getCardStyle().Render(content.String())
}

// viewCaptureForm renders the form for capturing event details.
func (i *Intent) viewCaptureForm() string {
	if i.state.captureForm == nil {
		return "Error: Form not initialized"
	}
	return i.state.captureForm.View()
}

// viewReviewInferredEvent renders the review UI for inferred bursts and facts.
func (i *Intent) viewReviewInferredEvent() string {
	baseView := i.buildReviewBaseView()

	switch i.state.reviewState.EditingMode {
	case EditingModeMetadata:
		return i.renderModalOverlay(baseView, i.getMetadataModalContent())
	case EditingModeBursts:
		return i.renderModalOverlay(baseView, i.getBurstModalContent())
	case EditingModeFacts:
		return i.renderModalOverlay(baseView, i.getFactModalContent())
	}

	return baseView
}

// buildReviewBaseView builds the base review view content.
func (i *Intent) buildReviewBaseView() string {
	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString("┌─ Review Inferred Event ────────────────────────┐\n")
	sb.WriteString("│                                                │\n")

	if i.state.reviewState.Event != nil {
		title := i.state.reviewState.Event.Text
		if len(title) > 40 {
			title = title[:37] + "..."
		}
		sb.WriteString(fmt.Sprintf("│ Event: %s                    │\n", title))
		sb.WriteString("│                                                │\n")
	}

	sb.WriteString("│ Inferred Bursts:                               │\n")
	bursts := i.state.reviewState.InferredBursts
	if len(bursts) == 0 {
		bursts = i.state.reviewState.AcceptedBursts
	}
	if len(bursts) > 0 {
		for idx, burst := range bursts {
			burstTitle := burst.Name
			if len(burstTitle) > 35 {
				burstTitle = burstTitle[:32] + "..."
			}
			sb.WriteString(fmt.Sprintf("│   [✓] Burst %d: %s              │\n", idx+1, burstTitle))
		}
	} else {
		sb.WriteString("│   (No bursts detected)                         │\n")
	}
	sb.WriteString("│                                                │\n")

	sb.WriteString("│ Inferred Facts:                                │\n")
	facts := i.state.reviewState.InferredFacts
	if len(facts) == 0 {
		facts = i.state.reviewState.AcceptedFacts
	}
	if len(facts) > 0 {
		for idx, fact := range facts {
			desc := fact.Text
			if len(desc) > 35 {
				desc = desc[:32] + "..."
			}
			sb.WriteString(fmt.Sprintf("│   [✓] Fact %d: %s              │\n", idx+1, desc))
		}
	} else {
		sb.WriteString("│   (No facts detected)                          │\n")
	}
	sb.WriteString("│                                                │\n")
	sb.WriteString("└────────────────────────────────────────────────┘\n")
	return sb.String()
}

// viewSubmit renders the submit confirmation.
func (i *Intent) viewSubmit() string {
	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString("┌─ Confirm Submission ───────────────────────────┐\n")
	sb.WriteString("│                                                │\n")

	if i.state.reviewState.Event != nil {
		title := i.state.reviewState.Event.Text
		if len(title) > 40 {
			title = title[:37] + "..."
		}
		sb.WriteString(fmt.Sprintf("│ Event: %s                    │\n", title))
		sb.WriteString(fmt.Sprintf("│ Date: %s                      │\n", i.state.reviewState.Event.Date))
		sb.WriteString("│                                                │\n")
		sb.WriteString(fmt.Sprintf("│ Bursts: %d                                    │\n", len(i.state.reviewState.AcceptedBursts)))
		sb.WriteString(fmt.Sprintf("│ Facts: %d                                     │\n", len(i.state.reviewState.AcceptedFacts)))
		sb.WriteString("│                                                │\n")
	}

	sb.WriteString("│ Ready to submit? Press Enter to confirm.       │\n")
	sb.WriteString("│                                                │\n")
	sb.WriteString("└────────────────────────────────────────────────┘\n")
	return sb.String()
}

// viewError renders an error state.
func (i *Intent) viewError() string {
	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString("┌─ Error ─────────────────────────────────────────┐\n")
	sb.WriteString("│                                                │\n")

	if i.state.error != nil {
		code := i.state.error.Code
		if len(code) > 40 {
			code = code[:37] + "..."
		}
		sb.WriteString(fmt.Sprintf("│ Code: %s                        │\n", code))

		msg := i.state.error.Message
		if len(msg) > 40 {
			msg = msg[:37] + "..."
		}
		sb.WriteString(fmt.Sprintf("│ Message: %s                 │\n", msg))
		sb.WriteString("│                                                │\n")
	}

	sb.WriteString("│ Press 'r' to retry or Esc to cancel.           │\n")
	sb.WriteString("│                                                │\n")
	sb.WriteString("└────────────────────────────────────────────────┘\n")

	return sb.String()
}

package capture_event

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/models"
	captureScreens "github.com/baphled/kariya/internal/cli/screens/capture"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/domain/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// setCompleted marks the intent as successfully completed.
//
// Expected:
//   - result must be non-nil.
//
// Side effects:
//   - Sets i.result to a completed IntentResult wrapping the given data.
//   - Sets i.active to false, stopping further Update processing.
func (i *Intent) setCompleted(result *Result) {
	i.result = intents.NewCompletedResult(result)
	i.active = false
}

// setCancelled marks the intent as cancelled by the user.
//
// Side effects:
//   - Sets i.result to a cancelled IntentResult.
//   - Sets i.active to false, stopping further Update processing.
func (i *Intent) setCancelled() {
	i.result = intents.NewCancelledResult[*Result]()
	i.active = false
}

// setFailed marks the intent as failed with a structured error.
//
// Expected:
//   - code is a non-empty error code string (e.g. "VALIDATION_ERROR").
//   - message describes the failure in human-readable form.
//   - cause may be nil if no underlying error exists.
//
// Side effects:
//   - Sets i.result to a failed IntentResult.
//   - Sets i.active to false, stopping further Update processing.
func (i *Intent) setFailed(code, message string, cause error) {
	i.result = intents.NewFailedResult[*Result](code, message, cause)
	i.active = false
}

// setFailedCmd marks the intent as failed and returns nil.
//
// Expected:
//   - code is a non-empty error code string.
//   - message describes the failure in human-readable form.
//   - cause may be nil if no underlying error exists.
//
// Returns:
//   - nil (always). Intended as a terminal return in Update branches.
//
// Side effects:
//   - Delegates to setFailed to mark the intent as failed and inactive.
func (i *Intent) setFailedCmd(code, message string, cause error) tea.Cmd {
	i.setFailed(code, message, cause)
	return nil
}

// showSubmitModal creates the loading modal and starts event submission.
//
// Returns:
//   - A batched tea.Cmd that runs both performSubmit and the modal's Init.
//
// Side effects:
//   - Creates a new LoadingModal on i.submitModal.
func (i *Intent) showSubmitModal() tea.Cmd {
	i.submitModal = feedback.NewLoadingModal("Saving event...", false)
	return tea.Batch(i.performSubmit(), i.submitModal.Init())
}

// performSubmit persists the captured event and accepted facts via domain services.
//
// Expected:
//   - i.reviewState.Event is non-nil and valid.
//   - i.eventService and i.context.CareerService are non-nil.
//
// Returns:
//   - A tea.Cmd that runs asynchronously and produces a SubmitCompleteMsg
//     on success or a SubmitErrorMsg on failure.
//
// Side effects:
//   - Calls CareerService.CaptureEvent to persist the event.
//   - Calls CareerService.SaveFact for each accepted fact without an ID.
//   - Sets a default date for quick-strategy events with a zero date.
func (i *Intent) performSubmit() tea.Cmd {
	event := i.reviewState.Event
	acceptedFacts := i.reviewState.AcceptedFacts
	strategy := i.strategy
	careerService := i.context.CareerService
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

// initializeFormForNew prepares an empty event for a new capture session.
//
// Returns:
//   - A no-op tea.Cmd (returns nil message).
//
// Side effects:
//   - Sets i.reviewState.Event to a fresh career.Event with current
//     timestamps and empty slices for tags and categories.
func (i *Intent) initializeFormForNew() tea.Cmd {
	i.reviewState.Event = &career.Event{
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Tags:       make([]string, 0),
		Categories: make([]string, 0),
	}
	return func() tea.Msg { return nil }
}

// transitionToStrategyScreen creates the strategy selection screen and activates it.
//
// Returns:
//   - nil (always). The screen is rendered on the next View call.
//
// Side effects:
//   - Sets i.currentState to StateChooseStrategy.
//   - Creates a new StrategySelectScreen and assigns it to i.activeScreen.
//   - Configures the screen with terminal dimensions, theme, and logo.
func (i *Intent) transitionToStrategyScreen() tea.Cmd {
	i.currentState = StateChooseStrategy
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

// transitionToFormScreen creates the event form screen for the given strategy.
//
// Expected:
//   - strategy is a valid CaptureStrategy (StrategyQuick or StrategyManual).
//
// Returns:
//   - A tea.Cmd from the screen's Init method if it implements Init, nil otherwise.
//
// Side effects:
//   - Sets i.currentState to StateForm and i.strategy.
//   - Creates a new EventFormScreen and assigns it to i.activeScreen.
//   - Configures the screen with terminal dimensions, theme, and logo.
func (i *Intent) transitionToFormScreen(strategy CaptureStrategy) tea.Cmd {
	i.currentState = StateForm
	i.strategy = strategy
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

// navigateReviewItems moves the selection cursor through inferred review items.
//
// Expected:
//   - delta is +1 (down) or -1 (up).
//   - i.reviewState is non-nil.
//
// Side effects:
//   - Updates SelectedIndex with wraparound at both ends.
//   - Updates SelectedItemType to "burst" or "fact" based on the new position
//     relative to the burst/fact boundary.
func (i *Intent) navigateReviewItems(delta int) {
	totalItems := len(i.reviewState.InferredBursts) + len(i.reviewState.InferredFacts)
	if totalItems == 0 {
		return
	}

	i.reviewState.SelectedIndex += delta
	if i.reviewState.SelectedIndex >= totalItems {
		i.reviewState.SelectedIndex = 0
	}
	if i.reviewState.SelectedIndex < 0 {
		i.reviewState.SelectedIndex = totalItems - 1
	}

	if i.reviewState.SelectedIndex < len(i.reviewState.InferredBursts) {
		i.reviewState.SelectedItemType = "burst"
	} else {
		i.reviewState.SelectedItemType = "fact"
		i.reviewState.SelectedIndex -= len(i.reviewState.InferredBursts)
	}
}

// acceptCurrentItem moves the currently selected item to the accepted list.
//
// Expected:
//   - i.reviewState is non-nil with a valid SelectedItemType and index.
//
// Side effects:
//   - Appends the selected burst or fact to AcceptedBursts/AcceptedFacts.
//   - Removes it from InferredBursts/InferredFacts.
//   - Clamps SelectedIndex to remain within bounds after removal.
func (i *Intent) acceptCurrentItem() {
	if i.reviewState.SelectedItemType == "burst" {
		idx := i.reviewState.SelectedIndex
		if idx >= 0 && idx < len(i.reviewState.InferredBursts) {
			burst := i.reviewState.InferredBursts[idx]
			i.reviewState.AcceptedBursts = append(i.reviewState.AcceptedBursts, burst)
			i.reviewState.InferredBursts = append(
				i.reviewState.InferredBursts[:idx],
				i.reviewState.InferredBursts[idx+1:]...)
			if i.reviewState.SelectedIndex >= len(i.reviewState.InferredBursts) && len(i.reviewState.InferredBursts) > 0 {
				i.reviewState.SelectedIndex = len(i.reviewState.InferredBursts) - 1
			}
		}
	} else if i.reviewState.SelectedItemType == "fact" {
		idx := i.reviewState.SelectedIndex
		if idx >= 0 && idx < len(i.reviewState.InferredFacts) {
			fact := i.reviewState.InferredFacts[idx]
			i.reviewState.AcceptedFacts = append(i.reviewState.AcceptedFacts, fact)
			i.reviewState.InferredFacts = append(
				i.reviewState.InferredFacts[:idx],
				i.reviewState.InferredFacts[idx+1:]...)
			if i.reviewState.SelectedIndex >= len(i.reviewState.InferredFacts) && len(i.reviewState.InferredFacts) > 0 {
				i.reviewState.SelectedIndex = len(i.reviewState.InferredFacts) - 1
			}
		}
	}
}

// rejectCurrentItem marks the currently selected item as rejected.
//
// Expected:
//   - i.reviewState is non-nil with a valid SelectedItemType and index.
//
// Side effects:
//   - Adds the item's ID to RejectedItems with reason "user_rejected".
//   - Removes the item from InferredBursts/InferredFacts.
//   - Clamps SelectedIndex to remain within bounds after removal.
func (i *Intent) rejectCurrentItem() {
	if i.reviewState.SelectedItemType == "burst" {
		idx := i.reviewState.SelectedIndex
		if idx >= 0 && idx < len(i.reviewState.InferredBursts) {
			burst := i.reviewState.InferredBursts[idx]
			if i.reviewState.RejectedItems == nil {
				i.reviewState.RejectedItems = make(map[string]string)
			}
			i.reviewState.RejectedItems[burst.ID] = "user_rejected"
			i.reviewState.InferredBursts = append(
				i.reviewState.InferredBursts[:idx],
				i.reviewState.InferredBursts[idx+1:]...)
			if i.reviewState.SelectedIndex >= len(i.reviewState.InferredBursts) && len(i.reviewState.InferredBursts) > 0 {
				i.reviewState.SelectedIndex = len(i.reviewState.InferredBursts) - 1
			}
		}
	} else if i.reviewState.SelectedItemType == "fact" {
		idx := i.reviewState.SelectedIndex
		if idx >= 0 && idx < len(i.reviewState.InferredFacts) {
			fact := i.reviewState.InferredFacts[idx]
			if i.reviewState.RejectedItems == nil {
				i.reviewState.RejectedItems = make(map[string]string)
			}
			i.reviewState.RejectedItems[fact.ID] = "user_rejected"
			i.reviewState.InferredFacts = append(
				i.reviewState.InferredFacts[:idx],
				i.reviewState.InferredFacts[idx+1:]...)
			if i.reviewState.SelectedIndex >= len(i.reviewState.InferredFacts) && len(i.reviewState.InferredFacts) > 0 {
				i.reviewState.SelectedIndex = len(i.reviewState.InferredFacts) - 1
			}
		}
	}
}

// modalContentData holds the rendered parts of an editing modal overlay.
//
// Used by renderModalOverlay to compose a centred overlay on top of the
// review screen's base view.
type modalContentData struct {
	// title is the modal header text (e.g. "Edit Metadata").
	title string

	// content is the modal body rendered by the underlying form model.
	content string

	// footer is the help/shortcut text shown at the bottom of the modal.
	footer string
}

// getMetadataModalContent returns the rendered modal parts for metadata editing.
//
// Returns:
//   - A modalContentData with title, content, and footer from the metadata modal.
//
// Side effects:
//   - Lazily creates the metadataModal if it is nil.
func (i *Intent) getMetadataModalContent() *modalContentData {
	if i.reviewState.metadataModal == nil {
		i.reviewState.metadataModal = models.NewMetadataEditorModelNew(
			i.reviewState.Event,
			i.context.CareerService,
			i.context.CLIEventService,
			context.Background(),
		)
	}
	return &modalContentData{
		title:   i.reviewState.metadataModal.GetTitle(),
		content: i.reviewState.metadataModal.GetContent(),
		footer:  i.reviewState.metadataModal.GetFooter(),
	}
}

// getBurstModalContent returns the rendered modal parts for burst editing.
//
// Returns:
//   - A modalContentData with title, content, and footer from the burst modal.
//   - nil if the burstModal has not been created.
func (i *Intent) getBurstModalContent() *modalContentData {
	if i.reviewState.burstModal == nil {
		return nil
	}
	return &modalContentData{
		title:   i.reviewState.burstModal.GetTitle(),
		content: i.reviewState.burstModal.GetContent(),
		footer:  i.reviewState.burstModal.GetFooter(),
	}
}

// getFactModalContent returns the rendered modal parts for fact editing.
//
// Returns:
//   - A modalContentData with title, content, and footer from the fact modal.
//   - nil if the factModal has not been created.
func (i *Intent) getFactModalContent() *modalContentData {
	if i.reviewState.factModal == nil {
		return nil
	}
	return &modalContentData{
		title:   i.reviewState.factModal.GetTitle(),
		content: i.reviewState.factModal.GetContent(),
		footer:  i.reviewState.factModal.GetFooter(),
	}
}

// getEditingModalContent returns the modal content for the active editing mode.
//
// Returns:
//   - A modalContentData for metadata, burst, or fact editing.
//   - nil if reviewState is nil or EditingMode is EditingModeNone.
func (i *Intent) getEditingModalContent() *modalContentData {
	if i.reviewState == nil {
		return nil
	}
	switch i.reviewState.EditingMode {
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

// renderModalOverlay renders a centred modal overlay on top of background content.
//
// Expected:
//   - background is the fully rendered base view string.
//   - modalContent may be nil, in which case background is returned unchanged.
//
// Returns:
//   - The composited string with the modal centred over the background.
//   - The unmodified background if modalContent is nil.
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

// overlayModal composites modal lines over background lines, centred horizontally.
//
// Expected:
//   - background and modal are newline-separated rendered strings.
//   - width is the terminal width for horizontal centring.
//
// Returns:
//   - A newline-joined string with modal lines replacing background lines
//     at the vertical centre.
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

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
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
)

// terminalDimensions returns the current terminal dimensions for modal sizing.
// Falls back to nil if terminal info is unavailable, letting the model
// apply its own defaults.
func (i *Intent) terminalDimensions() *models.MetadataEditorDimensions {
	info := i.GetTerminalInfo()
	if info == nil {
		return nil
	}
	return &models.MetadataEditorDimensions{
		TerminalWidth:  info.Width,
		TerminalHeight: info.Height,
	}
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
			i.terminalDimensions(),
		)
	}
	return &modalContentData{
		title:   i.reviewState.metadataModal.GetTitle(),
		content: i.reviewState.metadataModal.GetContent(),
		footer:  renderFormModalFooter(),
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
	// Detect editing state from the model's footer text.
	// The deprecated models package does not export an IsEditing() method.
	editing := strings.Contains(i.reviewState.burstModal.GetFooter(), "Save")
	return &modalContentData{
		title:   i.reviewState.burstModal.GetTitle(),
		content: i.reviewState.burstModal.GetContent(),
		footer:  renderBurstModalFooter(editing),
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
		footer:  renderFormModalFooter(),
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
//
// The footer is appended to the content rather than passed via SetFooter()
// to avoid OverlayModal.buildContent() applying MutedColor styling that
// would conflict with pre-styled UIKit badge text.
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

	// Include the footer as part of the content body to preserve
	// badge styling. OverlayModal.SetFooter() wraps text in MutedColor
	// which strips pre-styled badge colours.
	content := modalContent.content
	if modalContent.footer != "" {
		content = content + "\n\n" + modalContent.footer
	}

	overlay := feedback.NewOverlayModal(modalContent.title, content)
	overlay.SetWidth(80)

	return overlay.RenderCentered(background, width, height)
}

// renderFormModalFooter returns a UIKit badge-styled footer for form modals.
//
// Used by metadata and fact editing modals that share the same key bindings.
func renderFormModalFooter() string {
	th := themes.NewDefaultTheme()
	return primitives.RenderHelpFooter(th,
		primitives.ConfirmBadge(th),
		primitives.NextFieldBadge(th),
		primitives.PrevBadge(th),
		primitives.CancelBadge(th),
	)
}

// renderBurstModalFooter returns a UIKit badge-styled footer for the burst modal.
//
// The burst modal has two modes: editing (form fields) and navigating (suggestion list).
// Each mode shows different key bindings.
func renderBurstModalFooter(editing bool) string {
	th := themes.NewDefaultTheme()
	if editing {
		return primitives.RenderHelpFooter(th,
			primitives.NextFieldBadge(th),
			primitives.ConfirmBadge(th),
			primitives.CancelBadge(th),
		)
	}
	return primitives.RenderHelpFooter(th,
		primitives.NavigateBadge(th),
		primitives.HelpKeyBadge("y", "Confirm", th),
		primitives.HelpKeyBadge("n", "Reject", th),
		primitives.EditBadge(th),
		primitives.BackBadge(th),
	)
}

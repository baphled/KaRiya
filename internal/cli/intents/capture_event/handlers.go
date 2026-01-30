package capture_event

import (
	"context"
	"fmt"

	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/models"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/domain/career"
	burstfact "github.com/baphled/kariya/internal/service/career/burstfact"
	tea "github.com/charmbracelet/bubbletea"
)

// handleScreenResult dispatches a screen result to the appropriate handler method.
//
// Returns:
//   - A tea.Cmd from the matched handler, or nil.
func (i *Intent) handleScreenResult(result screens.ScreenResult) tea.Cmd {
	return behaviors.NewScreenResultDispatcher(i).Dispatch(result)
}

// HandleNavigate processes navigation results from screens.
//
// Expected:
//   - result.Data() is either a string action ("edit_metadata", "edit_bursts",
//     "edit_facts") or a CaptureStrategy value.
//
// Returns:
//   - A tea.Cmd to initialise the appropriate modal or screen.
//   - nil and marks intent as failed if data type is unrecognised.
//
// Side effects:
//   - Sets editing mode and creates modal instances for edit actions.
//   - Transitions to the form screen for strategy selection.
//
// Implements behaviors.ScreenResultHandler.
func (i *Intent) HandleNavigate(result *screens.NavigateResult) tea.Cmd {
	data := result.Data()

	if action, ok := data.(string); ok {
		switch action {
		case "edit_metadata":
			if i.reviewState == nil || i.reviewState.Event == nil {
				return i.setFailedCmd("NO_EVENT", "No event to edit", nil)
			}
			i.reviewState.EditingMode = EditingModeMetadata
			i.reviewState.metadataModal = models.NewMetadataEditorModelNew(
				i.reviewState.Event,
				i.context.CareerService,
				i.context.CLIEventService,
				context.Background(),
			)
			return i.reviewState.metadataModal.Init()

		case "edit_bursts":
			i.reviewState.EditingMode = EditingModeBursts
			var suggestions []burstfact.BurstSuggestion
			for _, b := range i.reviewState.InferredBursts {
				suggestions = append(suggestions, burstfact.BurstSuggestion{
					Name:        b.Name,
					Description: b.Description,
				})
			}
			i.reviewState.burstModal = models.NewBurstSuggestionModelNew(
				i.context.CareerService,
				suggestions,
				context.Background(),
			)
			return i.reviewState.burstModal.Init()

		case "edit_facts":
			i.reviewState.EditingMode = EditingModeFacts
			var fact *career.Fact
			if len(i.reviewState.InferredFacts) > 0 {
				fact = i.reviewState.InferredFacts[0]
			} else {
				fact = &career.Fact{Text: ""}
			}
			i.reviewState.factModal = models.NewFactEditorModelNew(
				fact,
				i.context.CareerService,
				context.Background(),
			)
			return i.reviewState.factModal.Init()

		default:
			return i.setFailedCmd("INVALID_NAVIGATION", fmt.Sprintf("Unknown navigation action: %s", action), nil)
		}
	}

	if strategy, ok := data.(CaptureStrategy); ok {
		i.strategy = strategy
		return i.transitionToFormScreen(strategy)
	}

	return i.setFailedCmd("INVALID_NAVIGATION_DATA", fmt.Sprintf("Invalid navigation data type: %T", data), nil)
}

// HandleCancel processes cancellation results from screens.
//
// Returns:
//   - nil after marking the intent cancelled (from StateChooseStrategy or
//     StateForm with a previous event).
//   - A tea.Cmd to transition back to the strategy or form screen.
//
// Side effects:
//   - Marks the intent as cancelled when no back-navigation is possible.
//   - Transitions to a prior state when cancellation acts as "go back".
//
// Implements behaviors.ScreenResultHandler.
func (i *Intent) HandleCancel(_ *screens.CancelResult) tea.Cmd {
	switch i.currentState {
	case StateChooseStrategy:
		i.setCancelled()
		return nil

	case StateForm:
		if i.context.PreviousEvent != nil {
			i.setCancelled()
			return nil
		}
		return i.transitionToStrategyScreen()

	case StateReview:
		return i.transitionToFormScreen(i.strategy)

	case StateSubmit:
		return nil

	default:
		return i.setFailedCmd("INVALID_CANCEL_STATE", fmt.Sprintf("Cannot cancel from state: %s", i.currentState), nil)
	}
}

// HandleSubmit processes form and data submission results from screens.
//
// Expected:
//   - result.Data() is a *career.Event (StateForm), map[string]interface{} with
//     "event", "bursts", "facts" keys (StateReview/StateSubmit).
//
// Returns:
//   - A tea.Cmd to show the submit modal or nil on completion.
//   - A failure command if the data type is unrecognised or validation fails.
//
// Side effects:
//   - Initialises review state from form data.
//   - Updates accepted bursts/facts from review data.
//   - Completes the intent with final result on submit confirmation.
//
// Implements behaviors.ScreenResultHandler.
func (i *Intent) HandleSubmit(result *screens.SubmitResult) tea.Cmd {
	data := result.Data()

	switch i.currentState {
	case StateForm:
		if event, ok := data.(*career.Event); ok {
			if err := event.Validate(); err != nil {
				return i.setFailedCmd("VALIDATION_ERROR", fmt.Sprintf("Event validation failed: %v", err), err)
			}

			i.reviewState = &ReviewInferredEventState{
				Event:          event,
				InferredBursts: make([]*career.Burst, 0),
				InferredFacts:  make([]*career.Fact, 0),
				AcceptedBursts: make([]*career.Burst, 0),
				AcceptedFacts:  make([]*career.Fact, 0),
				RejectedItems:  make(map[string]string),
			}

			return i.showSubmitModal()
		}
		return i.setFailedCmd("INVALID_FORM_DATA", fmt.Sprintf("Invalid form data type: %T", data), nil)

	case StateReview:
		if reviewData, ok := data.(map[string]interface{}); ok {
			//nolint:errcheck // Type assertions are safe for map data extraction.
			event, _ := reviewData["event"].(*career.Event)
			//nolint:errcheck // Type assertions are safe for map data extraction.
			bursts, _ := reviewData["bursts"].([]*career.Burst)
			//nolint:errcheck // Type assertions are safe for map data extraction.
			facts, _ := reviewData["facts"].([]*career.Fact)

			i.reviewState.Event = event
			i.reviewState.AcceptedBursts = bursts
			i.reviewState.AcceptedFacts = facts

			return i.showSubmitModal()
		}
		return i.setFailedCmd("INVALID_REVIEW_DATA", fmt.Sprintf("Invalid review data type: %T", data), nil)

	case StateSubmit:
		if submitData, ok := data.(map[string]interface{}); ok {
			//nolint:errcheck // Type assertion is safe for map data extraction.
			event, _ := submitData["event"].(*career.Event)
			//nolint:errcheck // Type assertion is safe for map data extraction.
			bursts, _ := submitData["bursts"].([]*career.Burst)
			//nolint:errcheck // Type assertion is safe for map data extraction.
			facts, _ := submitData["facts"].([]*career.Fact)

			i.result = &intents.IntentResult[*Result]{
				Status: intents.Completed,
				Data: &Result{
					Event:  event,
					Bursts: bursts,
					Facts:  facts,
				},
			}
			i.active = false
			return nil
		}
		return i.setFailedCmd("INVALID_SUBMIT_DATA", fmt.Sprintf("Invalid submit data type: %T", data), nil)

	default:
		return i.setFailedCmd("INVALID_SUBMIT_STATE", fmt.Sprintf("Cannot submit from state: %s", i.currentState), nil)
	}
}

// HandleError processes error results from screens.
//
// Expected:
//   - result.Data() is a map[string]interface{} with "error" (error) and
//     "message" (string) keys, or any value convertible to a string.
//
// Returns:
//   - nil after marking the intent as failed.
//
// Side effects:
//   - Marks the intent as failed with the extracted error details.
//
// Implements behaviors.ScreenResultHandler.
func (i *Intent) HandleError(result *screens.ErrorResult) tea.Cmd {
	data := result.Data()
	if errorData, ok := data.(map[string]interface{}); ok {
		//nolint:errcheck // Type assertion is safe for map data extraction.
		err, _ := errorData["error"].(error)
		//nolint:errcheck // Type assertion is safe for map data extraction.
		msg, _ := errorData["message"].(string)
		return i.setFailedCmd("SCREEN_ERROR", msg, err)
	}
	return i.setFailedCmd("SCREEN_ERROR", "Unknown screen error", fmt.Errorf("%v", data))
}

// handleReviewSubmit processes submit/confirm actions from the review state.
//
// Returns:
//   - nil after completing the intent (when postSaveReview is true).
//   - A tea.Cmd to perform event persistence (when postSaveReview is false).
//
// Side effects:
//   - Completes the intent with review data if post-save review is active.
//   - Transitions to StateSubmit and triggers persistence otherwise.
func (i *Intent) handleReviewSubmit() tea.Cmd {
	if i.postSaveReview {
		result := &Result{
			Event:          i.reviewState.Event,
			Bursts:         i.reviewState.InferredBursts,
			Facts:          i.reviewState.InferredFacts,
			AcceptedFields: make(map[string]bool),
			RejectedFields: i.reviewState.RejectedItems,
		}
		i.setCompleted(result)
		return nil
	}
	i.currentState = StateSubmit
	return i.performSubmit()
}

// updateChooseStrategy handles messages in the strategy selection state.
//
// Expected:
//   - msg is a tea.KeyMsg or StrategySelectedMsg.
//
// Returns:
//   - A tea.Cmd to initialise the form (on enter/strategy selection).
//   - tea.Quit on quit key, nil for navigation and help.
//
// Side effects:
//   - Updates selectedStrategyIndex on arrow/vim key navigation.
//   - Transitions to StateForm and initialises the capture form on selection.
//   - Toggles help or cancels the intent on global keys.
func (i *Intent) updateChooseStrategy(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			if i.selectedStrategyIndex > 0 {
				i.selectedStrategyIndex--
			}
			return nil

		case tea.KeyDown:
			if i.selectedStrategyIndex < 1 {
				i.selectedStrategyIndex++
			}
			return nil

		case tea.KeyEnter:
			strategies := []CaptureStrategy{StrategyQuick, StrategyManual}
			i.strategy = strategies[i.selectedStrategyIndex]
			i.captureForm.SetStrategy(string(i.strategy))
			i.currentState = StateForm
			return i.captureForm.Init()
		}

		switch msg.String() {
		case "k":
			if i.selectedStrategyIndex > 0 {
				i.selectedStrategyIndex--
			}
			return nil

		case "j":
			if i.selectedStrategyIndex < 1 {
				i.selectedStrategyIndex++
			}
			return nil
		}

		switch intents.HandleGlobalKeys(msg) {
		case intents.KeyQuit:
			return tea.Quit
		case intents.KeyHelp:
			i.ToggleHelp()
			return nil
		case intents.KeyBack:
			i.setCancelled()
			return nil
		}

	case StrategySelectedMsg:
		i.currentState = StateForm
		return nil
	}

	return nil
}

// updateCaptureForm handles messages in the form entry state.
//
// Expected:
//   - msg is a tea.KeyMsg, models.SubmitMsg, or FormSubmittedMsg.
//
// Returns:
//   - A tea.Cmd from the form's Update method for input handling.
//   - tea.Quit on quit key, nil for navigation and help.
//
// Side effects:
//   - Delegates non-global key messages to the capture form model.
//   - Transitions to StateReview on successful form submission.
//   - Sets an error on validation failure or nil event submission.
//   - Transitions back to StateChooseStrategy on back key (unless editing).
func (i *Intent) updateCaptureForm(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch intents.HandleGlobalKeys(msg) {
		case intents.KeyQuit:
			return tea.Quit
		case intents.KeyHelp:
			i.ToggleHelp()
			return nil
		case intents.KeyBack:
			if i.context.PreviousEvent != nil {
				i.setCancelled()
				return nil
			}
			i.currentState = StateChooseStrategy
			return nil
		}

		switch msg.String() {
		case "ctrl+s":
			return i.captureForm.SubmitForm()
		}

	case models.SubmitMsg:
		if msg.Err != nil {
			i.intentError = &intents.IntentError{
				Code:    "FORM_SUBMISSION_ERROR",
				Message: msg.Err.Error(),
				Cause:   msg.Err,
			}
			return nil
		}

		if msg.Event == nil {
			i.setFailed("INVALID_FORM", "Form submission with nil event", nil)
			return nil
		}

		if err := msg.Event.Validate(); err != nil {
			i.setFailed("VALIDATION_ERROR", fmt.Sprintf("Form validation failed: %v", err), err)
			return nil
		}

		i.reviewState.Event = msg.Event
		i.currentState = StateReview
		return nil

	case FormSubmittedMsg:
		if msg.Event == nil {
			i.setFailed("INVALID_FORM", "Form submission with nil event", nil)
			return nil
		}

		if err := msg.Event.Validate(); err != nil {
			i.setFailed("VALIDATION_ERROR", fmt.Sprintf("Form validation failed: %v", err), err)
			return nil
		}

		i.reviewState.Event = msg.Event
		i.currentState = StateReview
		return nil
	}

	_, formCmd := i.captureForm.Update(msg)

	return formCmd
}

// updateReviewInferredEvent handles messages in the review state.
//
// Expected:
//   - msg is a tea.KeyMsg, ReviewConfirmedMsg, ReviewCancelledMsg, or
//     ReviewBackMsg.
//
// Returns:
//   - A tea.Cmd from modal Update when an editing modal is active.
//   - A tea.Cmd from handleReviewSubmit on enter/ctrl+s.
//   - tea.Quit on quit key, nil for navigation and item accept/reject.
//
// Side effects:
//   - Delegates to the active editing modal (metadata, bursts, facts).
//   - Clears editing mode when a modal is submitted or cancelled.
//   - Navigates review items, accepts/rejects the selected item.
//   - Opens editing modals on e/b/f keys.
//   - Transitions to StateSubmit on ReviewConfirmedMsg.
func (i *Intent) updateReviewInferredEvent(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch intents.HandleGlobalKeys(msg) {
		case intents.KeyQuit:
			return tea.Quit
		case intents.KeyHelp:
			i.ToggleHelp()
			return nil
		case intents.KeyBack:
			if i.reviewState.EditingMode != EditingModeNone {
				i.reviewState.metadataModal = nil
				i.reviewState.burstModal = nil
				i.reviewState.factModal = nil
				i.reviewState.EditingMode = EditingModeNone
				return nil
			}
			i.currentState = StateForm
			return nil
		}
	}

	switch i.reviewState.EditingMode {
	case EditingModeMetadata:
		if i.reviewState.metadataModal != nil {
			modal, cmd := i.reviewState.metadataModal.Update(msg)
			//nolint:errcheck // Type assertion is safe - Update always returns same type.
			i.reviewState.metadataModal = modal.(*models.MetadataEditorModelNew)

			if i.reviewState.metadataModal.IsSubmitted() {
				i.reviewState.Event = i.reviewState.metadataModal.GetEvent()
				i.reviewState.metadataModal = nil
				i.reviewState.EditingMode = EditingModeNone
			} else if i.reviewState.metadataModal.IsCancelled() {
				i.reviewState.metadataModal = nil
				i.reviewState.EditingMode = EditingModeNone
			}
			return cmd
		}

	case EditingModeBursts:
		if i.reviewState.burstModal != nil {
			modal, cmd := i.reviewState.burstModal.Update(msg)
			//nolint:errcheck // Type assertion is safe - Update always returns same type.
			i.reviewState.burstModal = modal.(*models.BurstSuggestionModelNew)

			return cmd
		}

	case EditingModeFacts:
		if i.reviewState.factModal != nil {
			modal, cmd := i.reviewState.factModal.Update(msg)
			//nolint:errcheck // Type assertion is safe - Update always returns same type.
			i.reviewState.factModal = modal.(*models.FactEditorModelNew)

			if i.reviewState.factModal.IsSubmitted() {
				i.reviewState.factModal = nil
				i.reviewState.EditingMode = EditingModeNone
			} else if i.reviewState.factModal.IsCancelled() {
				i.reviewState.factModal = nil
				i.reviewState.EditingMode = EditingModeNone
			}
			return cmd
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle special keys via tea.Key* constants.
		switch msg.Type {
		case tea.KeyEnter:
			return i.handleReviewSubmit()

		case tea.KeyUp:
			i.navigateReviewItems(-1)
			return nil

		case tea.KeyDown:
			i.navigateReviewItems(1)
			return nil
		}

		// Handle rune/vim keys and chords via string comparison.
		switch msg.String() {
		case "ctrl+s":
			return i.handleReviewSubmit()

		case "e":
			i.reviewState.EditingMode = EditingModeMetadata
			if i.reviewState.metadataModal == nil && i.reviewState.Event != nil {
				i.reviewState.metadataModal = models.NewMetadataEditorModelNew(
					i.reviewState.Event,
					i.context.CareerService,
					i.context.CLIEventService,
					context.Background(),
				)
				return i.reviewState.metadataModal.Init()
			}
			return nil

		case "b":
			i.reviewState.EditingMode = EditingModeBursts
			if i.reviewState.burstModal == nil {
				var suggestions []burstfact.BurstSuggestion
				i.reviewState.burstModal = models.NewBurstSuggestionModelNew(
					i.context.CareerService,
					suggestions,
					context.Background(),
				)
				return i.reviewState.burstModal.Init()
			}
			return nil

		case "f":
			i.reviewState.EditingMode = EditingModeFacts
			if i.reviewState.factModal == nil {
				i.reviewState.factModal = models.NewFactEditorModelNew(
					&career.Fact{},
					i.context.CareerService,
					context.Background(),
				)
				return i.reviewState.factModal.Init()
			}
			return nil

		case "a":
			i.acceptCurrentItem()
			return nil

		case "r":
			i.rejectCurrentItem()
			return nil

		case "j":
			i.navigateReviewItems(1)
			return nil

		case "k":
			i.navigateReviewItems(-1)
			return nil
		}

	case ReviewConfirmedMsg:
		i.reviewState.AcceptedBursts = msg.AcceptedBursts
		i.reviewState.AcceptedFacts = msg.AcceptedFacts
		i.reviewState.RejectedItems = msg.RejectedItems
		i.currentState = StateSubmit
		return i.performSubmit()

	case ReviewCancelledMsg:
		i.setCancelled()
		return nil

	case ReviewBackMsg:
		i.currentState = StateForm
		return nil
	}

	return nil
}

// updateSubmit handles messages in the submission state.
//
// Expected:
//   - msg is a SubmitCompleteMsg, SubmitErrorMsg, or tea.KeyMsg.
//
// Returns:
//   - nil after completing or failing the intent on submit messages.
//   - A tea.Cmd to retry submission on "r" key.
//   - tea.Quit on quit key, nil for help and back navigation.
//
// Side effects:
//   - Completes the intent with the final result on SubmitCompleteMsg.
//   - Sets an error on SubmitErrorMsg.
//   - Transitions back to StateReview on back key.
func (i *Intent) updateSubmit(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case SubmitCompleteMsg:
		result := &Result{
			Event:          i.reviewState.Event,
			Bursts:         i.reviewState.AcceptedBursts,
			Facts:          i.reviewState.AcceptedFacts,
			AcceptedFields: make(map[string]bool),
			RejectedFields: i.reviewState.RejectedItems,
		}
		i.setCompleted(result)
		return nil

	case SubmitErrorMsg:
		i.intentError = &intents.IntentError{
			Code:    msg.Code,
			Message: msg.Message,
			Cause:   msg.Cause,
		}
		return nil

	case tea.KeyMsg:
		switch intents.HandleGlobalKeys(msg) {
		case intents.KeyQuit:
			return tea.Quit
		case intents.KeyHelp:
			i.ToggleHelp()
			return nil
		case intents.KeyBack:
			i.currentState = StateReview
			return nil
		}

		switch msg.String() {
		case "r":
			return i.performSubmit()
		}
	}

	return nil
}

// updateEditingModal handles messages when an editing modal is active in the
// screens architecture path.
//
// Expected:
//   - reviewState.EditingMode is not EditingModeNone.
//   - The corresponding modal (metadataModal, burstModal, factModal) is non-nil.
//
// Returns:
//   - nil on escape key (closes the modal).
//   - A tea.Cmd from the active modal's Update method.
//
// Side effects:
//   - Clears all modals and resets EditingMode to EditingModeNone on escape.
//   - Updates the event from metadata modal on submission.
//   - Resets editing state when a modal is submitted or cancelled.
func (i *Intent) updateEditingModal(msg tea.Msg) tea.Cmd {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.Type == tea.KeyEsc {
			i.reviewState.metadataModal = nil
			i.reviewState.burstModal = nil
			i.reviewState.factModal = nil
			i.reviewState.EditingMode = EditingModeNone
			return nil
		}
	}

	switch i.reviewState.EditingMode {
	case EditingModeMetadata:
		if i.reviewState.metadataModal != nil {
			modal, cmd := i.reviewState.metadataModal.Update(msg)
			//nolint:errcheck // Type assertion is safe - Update always returns same type.
			i.reviewState.metadataModal = modal.(*models.MetadataEditorModelNew)

			if i.reviewState.metadataModal.IsSubmitted() {
				i.reviewState.Event = i.reviewState.metadataModal.GetEvent()
				i.reviewState.metadataModal = nil
				i.reviewState.EditingMode = EditingModeNone
			} else if i.reviewState.metadataModal.IsCancelled() {
				i.reviewState.metadataModal = nil
				i.reviewState.EditingMode = EditingModeNone
			}
			return cmd
		}

	case EditingModeBursts:
		if i.reviewState.burstModal != nil {
			modal, cmd := i.reviewState.burstModal.Update(msg)
			//nolint:errcheck // Type assertion is safe - Update always returns same type.
			i.reviewState.burstModal = modal.(*models.BurstSuggestionModelNew)
			return cmd
		}

	case EditingModeFacts:
		if i.reviewState.factModal != nil {
			modal, cmd := i.reviewState.factModal.Update(msg)
			//nolint:errcheck // Type assertion is safe - Update always returns same type.
			i.reviewState.factModal = modal.(*models.FactEditorModelNew)

			if i.reviewState.factModal.IsSubmitted() {
				i.reviewState.factModal = nil
				i.reviewState.EditingMode = EditingModeNone
			} else if i.reviewState.factModal.IsCancelled() {
				i.reviewState.factModal = nil
				i.reviewState.EditingMode = EditingModeNone
			}
			return cmd
		}
	}

	return nil
}

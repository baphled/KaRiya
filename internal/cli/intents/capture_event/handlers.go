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

// handleScreenResult processes results from screen updates.
// This is the central hub for all screen-to-intent communication.
func (i *Intent) handleScreenResult(result screens.ScreenResult) tea.Cmd {
	return behaviors.NewScreenResultDispatcher(i).Dispatch(result)
}

// HandleNavigate handles navigation actions from screens.
//
// Implements ScreenResultHandler interface.
func (i *Intent) HandleNavigate(result *screens.NavigateResult) tea.Cmd {
	data := result.Data()

	if action, ok := data.(string); ok {
		switch action {
		case "edit_metadata":
			if i.state.reviewState == nil || i.state.reviewState.Event == nil {
				return i.setFailedCmd("NO_EVENT", "No event to edit", nil)
			}
			i.state.reviewState.EditingMode = EditingModeMetadata
			i.state.reviewState.metadataModal = models.NewMetadataEditorModelNew(
				i.state.reviewState.Event,
				i.context.CareerService,
				i.context.CLIEventService,
				context.Background(),
			)
			return i.state.reviewState.metadataModal.Init()

		case "edit_bursts":
			i.state.reviewState.EditingMode = EditingModeBursts
			var suggestions []burstfact.BurstSuggestion
			for _, b := range i.state.reviewState.InferredBursts {
				suggestions = append(suggestions, burstfact.BurstSuggestion{
					Name:        b.Name,
					Description: b.Description,
				})
			}
			i.state.reviewState.burstModal = models.NewBurstSuggestionModelNew(
				i.context.CareerService,
				suggestions,
				context.Background(),
			)
			return i.state.reviewState.burstModal.Init()

		case "edit_facts":
			i.state.reviewState.EditingMode = EditingModeFacts
			var fact *career.Fact
			if len(i.state.reviewState.InferredFacts) > 0 {
				fact = i.state.reviewState.InferredFacts[0]
			} else {
				fact = &career.Fact{Text: ""}
			}
			i.state.reviewState.factModal = models.NewFactEditorModelNew(
				fact,
				i.context.CareerService,
				context.Background(),
			)
			return i.state.reviewState.factModal.Init()

		default:
			return i.setFailedCmd("INVALID_NAVIGATION", fmt.Sprintf("Unknown navigation action: %s", action), nil)
		}
	}

	if strategy, ok := data.(CaptureStrategy); ok {
		i.state.strategy = strategy
		return i.transitionToFormScreen(strategy)
	}

	return i.setFailedCmd("INVALID_NAVIGATION_DATA", fmt.Sprintf("Invalid navigation data type: %T", data), nil)
}

// HandleCancel handles cancellation from screens.
//
// Implements ScreenResultHandler interface.
func (i *Intent) HandleCancel(_ *screens.CancelResult) tea.Cmd {
	switch i.state.currentState {
	case StateChooseStrategy:
		i.setCancelled()
		return nil

	case StateForm:
		if i.state.context.PreviousEvent != nil {
			i.setCancelled()
			return nil
		}
		return i.transitionToStrategyScreen()

	case StateReview:
		return i.transitionToFormScreen(i.state.strategy)

	case StateSubmit:
		return nil

	default:
		return i.setFailedCmd("INVALID_CANCEL_STATE", fmt.Sprintf("Cannot cancel from state: %s", i.state.currentState), nil)
	}
}

// HandleSubmit handles form/data submission from screens.
//
// Implements ScreenResultHandler interface.
func (i *Intent) HandleSubmit(result *screens.SubmitResult) tea.Cmd {
	data := result.Data()

	switch i.state.currentState {
	case StateForm:
		if event, ok := data.(*career.Event); ok {
			if err := event.Validate(); err != nil {
				return i.setFailedCmd("VALIDATION_ERROR", fmt.Sprintf("Event validation failed: %v", err), err)
			}

			i.state.reviewState = &ReviewInferredEventState{
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

			i.state.reviewState.Event = event
			i.state.reviewState.AcceptedBursts = bursts
			i.state.reviewState.AcceptedFacts = facts

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
		return i.setFailedCmd("INVALID_SUBMIT_STATE", fmt.Sprintf("Cannot submit from state: %s", i.state.currentState), nil)
	}
}

// HandleError handles errors from screens.
//
// Implements ScreenResultHandler interface.
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

// updateChooseStrategy handles messages while choosing capture strategy.
func (i *Intent) updateChooseStrategy(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if i.state.selectedStrategyIndex > 0 {
				i.state.selectedStrategyIndex--
			}
			return nil

		case "down", "j":
			if i.state.selectedStrategyIndex < 1 {
				i.state.selectedStrategyIndex++
			}
			return nil

		case "enter":
			strategies := []CaptureStrategy{StrategyQuick, StrategyManual}
			i.state.strategy = strategies[i.state.selectedStrategyIndex]
			i.state.captureForm.SetStrategy(string(i.state.strategy))
			i.state.currentState = StateForm
			return i.state.captureForm.Init()
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
		i.state.currentState = StateForm
		return nil
	}

	return nil
}

// updateCaptureForm handles messages while capturing event details.
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
			i.state.currentState = StateChooseStrategy
			return nil
		}

		switch msg.String() {
		case "ctrl+s":
			return i.state.captureForm.SubmitForm()
		}

	case models.SubmitMsg:
		if msg.Err != nil {
			i.state.error = &intents.IntentError{
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

		i.state.reviewState.Event = msg.Event
		i.state.currentState = StateReview
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

		i.state.reviewState.Event = msg.Event
		i.state.currentState = StateReview
		return nil
	}

	_, formCmd := i.state.captureForm.Update(msg)

	return formCmd
}

// updateReviewInferredEvent handles messages while reviewing inferred bursts and facts.
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
			if i.state.reviewState.EditingMode != EditingModeNone {
				i.state.reviewState.metadataModal = nil
				i.state.reviewState.burstModal = nil
				i.state.reviewState.factModal = nil
				i.state.reviewState.EditingMode = EditingModeNone
				return nil
			}
			i.state.currentState = StateForm
			return nil
		}
	}

	switch i.state.reviewState.EditingMode {
	case EditingModeMetadata:
		if i.state.reviewState.metadataModal != nil {
			modal, cmd := i.state.reviewState.metadataModal.Update(msg)
			//nolint:errcheck // Type assertion is safe - Update always returns same type.
			i.state.reviewState.metadataModal = modal.(*models.MetadataEditorModelNew)

			if i.state.reviewState.metadataModal.IsSubmitted() {
				i.state.reviewState.Event = i.state.reviewState.metadataModal.GetEvent()
				i.state.reviewState.metadataModal = nil
				i.state.reviewState.EditingMode = EditingModeNone
			} else if i.state.reviewState.metadataModal.IsCancelled() {
				i.state.reviewState.metadataModal = nil
				i.state.reviewState.EditingMode = EditingModeNone
			}
			return cmd
		}

	case EditingModeBursts:
		if i.state.reviewState.burstModal != nil {
			modal, cmd := i.state.reviewState.burstModal.Update(msg)
			//nolint:errcheck // Type assertion is safe - Update always returns same type.
			i.state.reviewState.burstModal = modal.(*models.BurstSuggestionModelNew)

			return cmd
		}

	case EditingModeFacts:
		if i.state.reviewState.factModal != nil {
			modal, cmd := i.state.reviewState.factModal.Update(msg)
			//nolint:errcheck // Type assertion is safe - Update always returns same type.
			i.state.reviewState.factModal = modal.(*models.FactEditorModelNew)

			if i.state.reviewState.factModal.IsSubmitted() {
				i.state.reviewState.factModal = nil
				i.state.reviewState.EditingMode = EditingModeNone
			} else if i.state.reviewState.factModal.IsCancelled() {
				i.state.reviewState.factModal = nil
				i.state.reviewState.EditingMode = EditingModeNone
			}
			return cmd
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+s", "enter":
			if i.state.postSaveReview {
				result := &Result{
					Event:          i.state.reviewState.Event,
					Bursts:         i.state.reviewState.InferredBursts,
					Facts:          i.state.reviewState.InferredFacts,
					AcceptedFields: make(map[string]bool),
					RejectedFields: i.state.reviewState.RejectedItems,
				}
				i.setCompleted(result)
				return nil
			}
			i.state.currentState = StateSubmit
			return i.performSubmit()

		case "e":
			i.state.reviewState.EditingMode = EditingModeMetadata
			if i.state.reviewState.metadataModal == nil && i.state.reviewState.Event != nil {
				i.state.reviewState.metadataModal = models.NewMetadataEditorModelNew(
					i.state.reviewState.Event,
					i.context.CareerService,
					i.context.CLIEventService,
					context.Background(),
				)
				return i.state.reviewState.metadataModal.Init()
			}
			return nil

		case "b":
			i.state.reviewState.EditingMode = EditingModeBursts
			if i.state.reviewState.burstModal == nil {
				var suggestions []burstfact.BurstSuggestion
				i.state.reviewState.burstModal = models.NewBurstSuggestionModelNew(
					i.context.CareerService,
					suggestions,
					context.Background(),
				)
				return i.state.reviewState.burstModal.Init()
			}
			return nil

		case "f":
			i.state.reviewState.EditingMode = EditingModeFacts
			if i.state.reviewState.factModal == nil {
				i.state.reviewState.factModal = models.NewFactEditorModelNew(
					&career.Fact{},
					i.context.CareerService,
					context.Background(),
				)
				return i.state.reviewState.factModal.Init()
			}
			return nil

		case "a":
			i.acceptCurrentItem()
			return nil

		case "r":
			i.rejectCurrentItem()
			return nil

		case "j", "down":
			totalItems := len(i.state.reviewState.InferredBursts) + len(i.state.reviewState.InferredFacts)
			if totalItems > 0 {
				i.state.reviewState.SelectedIndex++
				if i.state.reviewState.SelectedIndex >= totalItems {
					i.state.reviewState.SelectedIndex = 0
				}
				if i.state.reviewState.SelectedIndex < len(i.state.reviewState.InferredBursts) {
					i.state.reviewState.SelectedItemType = "burst"
				} else {
					i.state.reviewState.SelectedItemType = "fact"
					i.state.reviewState.SelectedIndex -= len(i.state.reviewState.InferredBursts)
				}
			}
			return nil

		case "k", "up":
			totalItems := len(i.state.reviewState.InferredBursts) + len(i.state.reviewState.InferredFacts)
			if totalItems > 0 {
				i.state.reviewState.SelectedIndex--
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
			return nil
		}

	case ReviewConfirmedMsg:
		i.state.reviewState.AcceptedBursts = msg.AcceptedBursts
		i.state.reviewState.AcceptedFacts = msg.AcceptedFacts
		i.state.reviewState.RejectedItems = msg.RejectedItems
		i.state.currentState = StateSubmit
		return i.performSubmit()

	case ReviewCancelledMsg:
		i.setCancelled()
		return nil

	case ReviewBackMsg:
		i.state.currentState = StateForm
		return nil
	}

	return nil
}

// updateSubmit handles messages while submitting the event.
func (i *Intent) updateSubmit(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case SubmitCompleteMsg:
		result := &Result{
			Event:          i.state.reviewState.Event,
			Bursts:         i.state.reviewState.AcceptedBursts,
			Facts:          i.state.reviewState.AcceptedFacts,
			AcceptedFields: make(map[string]bool),
			RejectedFields: i.state.reviewState.RejectedItems,
		}
		i.setCompleted(result)
		return nil

	case SubmitErrorMsg:
		i.state.error = &intents.IntentError{
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
			i.state.currentState = StateReview
			return nil
		}

		switch msg.String() {
		case "r":
			return i.performSubmit()
		}
	}

	return nil
}

// updateEditingModal handles modal updates when using screens architecture.
// This is called when an editing modal (metadata, bursts, facts) is active.
func (i *Intent) updateEditingModal(msg tea.Msg) tea.Cmd {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.String() == "esc" {
			i.state.reviewState.metadataModal = nil
			i.state.reviewState.burstModal = nil
			i.state.reviewState.factModal = nil
			i.state.reviewState.EditingMode = EditingModeNone
			return nil
		}
	}

	switch i.state.reviewState.EditingMode {
	case EditingModeMetadata:
		if i.state.reviewState.metadataModal != nil {
			modal, cmd := i.state.reviewState.metadataModal.Update(msg)
			//nolint:errcheck // Type assertion is safe - Update always returns same type.
			i.state.reviewState.metadataModal = modal.(*models.MetadataEditorModelNew)

			if i.state.reviewState.metadataModal.IsSubmitted() {
				i.state.reviewState.Event = i.state.reviewState.metadataModal.GetEvent()
				i.state.reviewState.metadataModal = nil
				i.state.reviewState.EditingMode = EditingModeNone
			} else if i.state.reviewState.metadataModal.IsCancelled() {
				i.state.reviewState.metadataModal = nil
				i.state.reviewState.EditingMode = EditingModeNone
			}
			return cmd
		}

	case EditingModeBursts:
		if i.state.reviewState.burstModal != nil {
			modal, cmd := i.state.reviewState.burstModal.Update(msg)
			//nolint:errcheck // Type assertion is safe - Update always returns same type.
			i.state.reviewState.burstModal = modal.(*models.BurstSuggestionModelNew)
			return cmd
		}

	case EditingModeFacts:
		if i.state.reviewState.factModal != nil {
			modal, cmd := i.state.reviewState.factModal.Update(msg)
			//nolint:errcheck // Type assertion is safe - Update always returns same type.
			i.state.reviewState.factModal = modal.(*models.FactEditorModelNew)

			if i.state.reviewState.factModal.IsSubmitted() {
				i.state.reviewState.factModal = nil
				i.state.reviewState.EditingMode = EditingModeNone
			} else if i.state.reviewState.factModal.IsCancelled() {
				i.state.reviewState.factModal = nil
				i.state.reviewState.EditingMode = EditingModeNone
			}
			return cmd
		}
	}

	return nil
}

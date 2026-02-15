package captureevent

import (
	"context"
	"fmt"
	"time"

	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/burst_management/modals"
	captureScreens "github.com/baphled/kariya/internal/cli/screens/capture"
	"github.com/baphled/kariya/internal/domain/career"
	burstfact "github.com/baphled/kariya/internal/service/career/burstfact"
	"github.com/baphled/kariya/internal/service/career/skillinference"
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
			if i.context.CareerService == nil {
				return i.setFailedCmd("NO_SERVICE", "Career service not available for metadata editing", nil)
			}
			i.reviewState.EditingMode = EditingModeMetadata
			i.reviewState.metadataModal = NewMetadataEditorModelNew(
				context.Background(),
				i.reviewState.Event,
				i.context.CareerService,
				i.context.CLIEventService,
				i.terminalDimensions(),
			)
			return i.reviewState.metadataModal.Init()

		case "edit_bursts":
			if i.context.CareerService == nil {
				return i.setFailedCmd("NO_SERVICE", "Career service not available for burst editing", nil)
			}
			i.reviewState.EditingMode = EditingModeBursts
			var suggestions []burstfact.BurstSuggestion
			for _, b := range i.reviewState.InferredBursts {
				suggestions = append(suggestions, burstfact.BurstSuggestion{
					Name:        b.Name,
					Description: b.Description,
				})
			}
			i.reviewState.burstModal = NewBurstSuggestionModelNew(
				context.Background(),
				i.context.CareerService,
				suggestions,
			)
			return i.reviewState.burstModal.Init()

		case "edit_facts":
			if i.context.CareerService == nil {
				return i.setFailedCmd("NO_SERVICE", "Career service not available for fact editing", nil)
			}
			i.reviewState.EditingMode = EditingModeFacts
			var fact *career.Fact
			if len(i.reviewState.InferredFacts) > 0 {
				fact = i.reviewState.InferredFacts[0]
			} else {
				fact = &career.Fact{Text: ""}
			}
			i.reviewState.factModal = NewFactEditorModelNew(
				context.Background(),
				fact,
				i.context.CareerService,
			)
			return i.reviewState.factModal.Init()

		case "edit_skills":
			i.reviewState.EditingMode = EditingModeSkills
			i.reviewState.skillModal = modals.NewSkillSuggestionModal(
				i.reviewState.InferredSkills,
				i.Theme(),
			)

			dims := i.terminalDimensions()
			if dims != nil {
				i.reviewState.skillModal.SetDimensions(dims.TerminalWidth, dims.TerminalHeight)
			}

			return i.reviewState.skillModal.Init()

		default:
			return i.setFailedCmd("INVALID_NAVIGATION", "Unknown navigation action: "+action, nil)
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
// Expected: the cancel result originates from the currently active screen.
//
// Returns:
//   - nil after marking the intent cancelled (from StateChooseStrategy or
//     StateForm with a previous event).
//   - A tea.Cmd to transition back to the strategy or form screen.
//
// Side effects:
//   - Marks the intent as cancelled when no back-navigation is possible.
//   - Transitions to a prior state when cancellation acts as "go back".
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
		if formData, ok := data.(*forms.CaptureEventFormData); ok {
			if !formData.SubmitConfirmed {
				return nil
			}

			event, err := eventFromFormData(formData)
			if err != nil {
				return i.showValidationErrorModal(err.Error())
			}

			if err := event.Validate(); err != nil {
				// Show validation error modal instead of silently failing
				return i.showValidationErrorModal(err.Error())
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
			event, eventOK := reviewData["event"].(*career.Event)
			if !eventOK {
				return i.setFailedCmd("INVALID_REVIEW_DATA", "Review data missing valid event", nil)
			}
			bursts, _ := reviewData["bursts"].([]*career.Burst)
			facts, _ := reviewData["facts"].([]*career.Fact)
			convertedSkills := extractSkillsFromReviewData(reviewData)

			i.reviewState.Event = event
			i.reviewState.AcceptedBursts = bursts
			i.reviewState.AcceptedFacts = facts
			i.reviewState.AcceptedSkills = convertedSkills

			if i.postSaveReview {
				// Persist accepted skills that were edited after initial save.
				if len(convertedSkills) > 0 && i.context.CareerService != nil {
					ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
					defer cancel()
					skillRepo := i.context.CareerService.GetSkillRepository()
					eventRepo := i.context.CareerService.GetEventRepository()
					for _, skill := range convertedSkills {
						if skill.ID == "" {
							if err := skillRepo.Create(ctx, skill); err != nil {
								return i.setFailedCmd("SKILL_SAVE_ERROR", fmt.Sprintf("Failed to save skill: %v", err), err)
							}
						}
						if err := eventRepo.LinkSkill(ctx, event.ID, skill.ID); err != nil {
							return i.setFailedCmd("SKILL_LINK_ERROR", fmt.Sprintf("Failed to link skill: %v", err), err)
						}
					}
				}

				i.result = &intents.IntentResult[*Result]{
					Status: intents.Completed,
					Data: &Result{
						Event:  event,
						Bursts: bursts,
						Facts:  facts,
						Skills: convertedSkills,
					},
				}
				i.active = false
				return nil
			}

			return i.showSubmitModal()
		}
		return i.setFailedCmd("INVALID_REVIEW_DATA", fmt.Sprintf("Invalid review data type: %T", data), nil)

	case StateSubmit:
		if submitData, ok := data.(map[string]interface{}); ok {
			event, eventOK := submitData["event"].(*career.Event)
			if !eventOK {
				return i.setFailedCmd("INVALID_SUBMIT_DATA", "Submit data missing valid event", nil)
			}
			bursts, _ := submitData["bursts"].([]*career.Burst)
			facts, _ := submitData["facts"].([]*career.Fact)

			i.result = &intents.IntentResult[*Result]{
				Status: intents.Completed,
				Data: &Result{
					Event:  event,
					Bursts: bursts,
					Facts:  facts,
					Skills: extractSkillsFromReviewData(submitData),
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
		err, _ := errorData["error"].(error)
		msg, _ := errorData["message"].(string)
		if msg == "" {
			msg = "Unknown screen error"
		}
		return i.setFailedCmd("SCREEN_ERROR", msg, err)
	}
	return i.setFailedCmd("SCREEN_ERROR", "Unknown screen error", fmt.Errorf("%v", data))
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
			if typed, ok := modal.(*MetadataEditorModelNew); ok {
				i.reviewState.metadataModal = typed
			}

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
			if typed, ok := modal.(*BurstSuggestionModelNew); ok {
				i.reviewState.burstModal = typed
			}
			return cmd
		}

	case EditingModeFacts:
		if i.reviewState.factModal != nil {
			modal, cmd := i.reviewState.factModal.Update(msg)
			if typed, ok := modal.(*FactEditorModelNew); ok {
				i.reviewState.factModal = typed
			}

			if i.reviewState.factModal.IsSubmitted() {
				i.reviewState.factModal = nil
				i.reviewState.EditingMode = EditingModeNone
			} else if i.reviewState.factModal.IsCancelled() {
				i.reviewState.factModal = nil
				i.reviewState.EditingMode = EditingModeNone
			}
			return cmd
		}

	case EditingModeSkills:
		if i.reviewState.skillModal != nil {
			model, cmd := i.reviewState.skillModal.Update(msg)
			if typed, ok := model.(*modals.SuggestionReviewModal); ok {
				i.reviewState.skillModal = typed
			}

			if !i.reviewState.skillModal.IsVisible() {
				accepted := i.reviewState.skillModal.GetAcceptedSkills()
				if len(accepted) > 0 {
					for _, s := range accepted {
						if s.Name == "" {
							continue
						}
						i.reviewState.AcceptedSkills = append(i.reviewState.AcceptedSkills, &career.Skill{
							Name:     s.Name,
							Category: s.Category,
						})
					}

					if screen, ok := i.activeScreen.(*captureScreens.EventReviewScreen); ok {
						var screenSkills []skillinference.SkillSuggestion
						for _, sk := range i.reviewState.AcceptedSkills {
							if sk.Name == "" {
								continue
							}
							screenSkills = append(screenSkills, skillinference.SkillSuggestion{
								Name:       sk.Name,
								Category:   sk.Category,
								Confidence: 1.0,
							})
						}
						screen.SetAcceptedSkills(screenSkills)
					}
				}
				i.reviewState.skillModal = nil
				i.reviewState.EditingMode = EditingModeNone
			}
			return cmd
		}
	}

	return nil
}

// eventFromFormData converts CaptureEventFormData to a career.Event.
//
// Expected:
//   - data must not be nil.
//
// Returns:
//   - A fully initialized career.Event ready for use.
//   - An error if date parsing fails.
//
// Side effects:
//   - None.
func eventFromFormData(data *forms.CaptureEventFormData) (*career.Event, error) {
	var eventDate time.Time
	var err error

	if data.Date == "" {
		eventDate = time.Now()
	} else {
		eventDate, err = forms.ParseDateString(data.Date)
		if err != nil {
			return nil, err
		}
	}

	event := &career.Event{
		Text:       data.Text,
		Date:       eventDate,
		Company:    data.Company,
		Project:    data.Project,
		Tags:       data.Tags,
		Categories: data.Categories,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	return event, nil
}

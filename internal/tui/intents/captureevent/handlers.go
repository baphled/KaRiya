package captureevent

import (
	"fmt"

	domcapture "github.com/baphled/kariya/internal/domain/capture"
	"github.com/baphled/kariya/internal/domain/career"
	burstfact "github.com/baphled/kariya/internal/service/career/burstfact"
	"github.com/baphled/kariya/internal/tui/forms"
	"github.com/baphled/kariya/internal/tui/intents"
	burstviews "github.com/baphled/kariya/internal/tui/views/burst"
	"github.com/baphled/kariya/internal/tui/views/event"
	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
	tea "github.com/charmbracelet/bubbletea"
)

// reviewUpdater defines the interface for views that can accept review updates
// from modal handlers. This decouples modal handlers from concrete view types.
type reviewUpdater interface {
	SetAcceptedBursts([]display.Burst)
	SetAcceptedFacts([]display.Fact)
	SetAcceptedSkills([]display.SkillSuggestion)
}

func (i *Intent) handleViewResult(result widgets.ViewResult) tea.Cmd {
	switch result.Type() {
	case widgets.ResultCancel:
		return i.HandleCancel(result)
	case widgets.ResultNavigate:
		return i.HandleNavigate(result)
	case widgets.ResultSubmit:
		return i.HandleSubmit(result)
	case widgets.ResultError:
		return i.HandleError(result)
	default:
		return nil
	}
}

// HandleNavigate processes navigation results from views.
//
// Expected:
//   - result.Data() is either a string action ("edit_metadata", "suggest_bursts",
//     "suggest_facts") or a CaptureStrategy value.
//
// Returns:
//   - A tea.Cmd to initialise the appropriate modal or screen.
//   - nil and marks intent as failed if data type is unrecognised.
//
// Side effects:
//   - Sets editing mode and creates modal instances for edit actions.
//   - Transitions to the form screen for strategy selection.
func (i *Intent) HandleNavigate(result widgets.ViewResult) tea.Cmd {
	data := result.Data()

	if action, ok := data.(string); ok {
		return i.handleNavigateAction(action)
	}

	if strategy, ok := data.(CaptureStrategy); ok {
		i.strategy = strategy
		return i.transitionToFormScreen(strategy)
	}

	return i.setFailedCmd("INVALID_NAVIGATION_DATA", fmt.Sprintf("Invalid navigation data type: %T", data), nil)
}

// handleNavigateAction dispatches string-typed navigation actions to specific handlers.
//
// Returns:
//   - A tea.Cmd from the appropriate action handler, or a failure command.
func (i *Intent) handleNavigateAction(action string) tea.Cmd {
	switch action {
	case "edit_metadata":
		return i.handleEditMetadata()
	case "suggest_bursts":
		return i.handleSuggestBursts()
	case "suggest_facts":
		return i.handleSuggestFacts()
	case "suggest_skills":
		return i.handleSuggestSkills()
	default:
		return i.setFailedCmd("INVALID_NAVIGATION", "Unknown navigation action: "+action, nil)
	}
}

// handleEditMetadata opens the metadata editing modal for the current event.
//
// Returns:
//   - A tea.Cmd from the modal's Init method, or a failure command.
//
// Side effects:
//   - Sets editing mode and creates the metadata modal instance.
func (i *Intent) handleEditMetadata() tea.Cmd {
	if i.reviewState == nil || i.reviewState.Event == nil {
		return i.setFailedCmd("NO_EVENT", "No event to edit", nil)
	}
	if i.context.CareerService == nil {
		return i.setFailedCmd("NO_SERVICE", "Career service not available for metadata editing", nil)
	}
	i.reviewState.EditingMode = EditingModeMetadata
	i.buildMetadataModal()
	return i.reviewState.metadataModal.Init()
}

// handleSuggestBursts opens the burst suggestion review modal.
//
// Returns:
//   - A tea.Cmd from the modal's Init method, or a failure command.
//
// Side effects:
//   - Sets editing mode and creates the burst suggestion modal.
func (i *Intent) handleSuggestBursts() tea.Cmd {
	if i.context.CareerService == nil {
		return i.setFailedCmd("NO_SERVICE", "Career service not available for burst editing", nil)
	}
	i.reviewState.EditingMode = EditingModeBursts
	var suggestions []burstfact.BurstSuggestion
	for _, b := range i.reviewState.InferredBursts {
		suggestions = append(suggestions, burstfact.BurstSuggestion{
			Name:        b.Name,
			Description: b.Description,
			EventIDs:    b.EventIDs,
		})
	}
	i.reviewState.burstModal = burstviews.NewSuggestionReview(display.BurstSuggestionsFromDomain(suggestions), i.Theme())
	dims := i.terminalDimensions()
	if dims != nil {
		i.reviewState.burstModal.SetDimensions(dims.TerminalWidth, dims.TerminalHeight)
	}
	return i.reviewState.burstModal.Init()
}

// handleSuggestFacts opens the fact suggestion review modal.
//
// Returns:
//   - A tea.Cmd from the modal's Init method.
//
// Side effects:
//   - Sets editing mode and creates the fact suggestion modal.
func (i *Intent) handleSuggestFacts() tea.Cmd {
	i.reviewState.EditingMode = EditingModeFacts

	var facts []career.Fact
	for _, f := range i.reviewState.InferredFacts {
		if f != nil {
			facts = append(facts, *f)
		}
	}

	displayFacts := make([]display.Fact, len(facts))
	for idx := range facts {
		displayFacts[idx] = display.FactFromDomain(&facts[idx])
	}
	i.reviewState.InferredDisplayFacts = displayFacts

	i.reviewState.factSuggestionModal = burstviews.NewFactSuggestion(
		displayFacts,
		i.Theme(),
	)

	dims := i.terminalDimensions()
	if dims != nil {
		i.reviewState.factSuggestionModal.SetDimensions(dims.TerminalWidth, dims.TerminalHeight)
	}

	return i.reviewState.factSuggestionModal.Init()
}

// handleSuggestSkills opens the skill suggestion review modal.
//
// Returns:
//   - A tea.Cmd from the modal's Init method.
//
// Side effects:
//   - Sets editing mode and creates the skill suggestion modal.
func (i *Intent) handleSuggestSkills() tea.Cmd {
	i.reviewState.EditingMode = EditingModeSkills
	i.reviewState.skillModal = burstviews.NewSkillSuggestion(
		display.SkillSuggestionsFromDomain(i.reviewState.InferredSkills),
		i.Theme(),
	)

	dims := i.terminalDimensions()
	if dims != nil {
		i.reviewState.skillModal.SetDimensions(dims.TerminalWidth, dims.TerminalHeight)
	}

	return i.reviewState.skillModal.Init()
}

// HandleCancel processes cancellation results from views.
//
// Expected: the cancel result originates from the currently active view.
//
// Returns:
//   - nil after marking the intent cancelled (from StateChooseStrategy or
//     StateForm with a previous event).
//   - A tea.Cmd to transition back to the strategy or form screen.
//
// Side effects:
//   - Marks the intent as cancelled when no back-navigation is possible.
//   - Transitions to a prior state when cancellation acts as "go back".
func (i *Intent) HandleCancel(_ widgets.ViewResult) tea.Cmd {
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

// HandleSubmit processes form and data submission results from views.
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
func (i *Intent) HandleSubmit(result widgets.ViewResult) tea.Cmd {
	data := result.Data()

	switch i.currentState {
	case StateForm:
		return i.handleFormSubmitData(data)
	case StateReview:
		return i.handleReviewSubmitData(data)
	case StateSubmit:
		return i.handleSubmitConfirmData(data)
	default:
		return i.setFailedCmd("INVALID_SUBMIT_STATE", fmt.Sprintf("Cannot submit from state: %s", i.currentState), nil)
	}
}

// handleFormSubmitData processes form submission data from the capture event form.
//
// Returns:
//   - A tea.Cmd to show the submit modal, a validation error modal, or a failure command.
//
// Side effects:
//   - Creates review state from validated form data.
func (i *Intent) handleFormSubmitData(data interface{}) tea.Cmd {
	formData, ok := data.(*forms.CaptureEventFormData)
	if !ok {
		return i.setFailedCmd("INVALID_FORM_DATA", fmt.Sprintf("Invalid form data type: %T", data), nil)
	}
	if !formData.SubmitConfirmed {
		return nil
	}

	capturedEvent, err := eventFromFormData(formData)
	if err != nil {
		return i.showValidationErrorModal(err.Error())
	}

	if err := capturedEvent.Validate(); err != nil {
		return i.showValidationErrorModal(err.Error())
	}

	i.reviewState = &ReviewInferredEventState{
		Event:          capturedEvent,
		InferredBursts: make([]*career.Burst, 0),
		InferredFacts:  make([]*career.Fact, 0),
		AcceptedBursts: make([]*career.Burst, 0),
		AcceptedFacts:  make([]*career.Fact, 0),
		AcceptedSkills: make([]*career.Skill, 0),
		RejectedItems:  make(map[string]string),
	}

	return i.showSubmitModal()
}

// handleReviewSubmitData processes review submission data with event, bursts, and facts.
//
// Returns:
//   - A tea.Cmd to perform post-save persistence, or a failure command.
//
// Side effects:
//   - Updates review state with accepted bursts, facts, and skills.
//   - Transitions to StateSubmit.
func (i *Intent) handleReviewSubmitData(data interface{}) tea.Cmd {
	reviewResult, ok := data.(event.ReviewResult)
	if !ok {
		return i.setFailedCmd("INVALID_REVIEW_DATA", fmt.Sprintf("Invalid review data type: %T", data), nil)
	}
	bursts := displayBurstsToPointers(reviewResult.Bursts, i.reviewState.InferredBursts)
	facts := displayFactsToPointers(reviewResult.Facts, i.reviewState.InferredFacts)
	convertedSkills := skillsFromDisplaySuggestions(reviewResult.Skills)

	i.reviewState.Event = domainEventFromDisplayEvent(reviewResult.Event)
	i.reviewState.AcceptedBursts = bursts
	i.reviewState.AcceptedFacts = facts
	i.reviewState.AcceptedSkills = convertedSkills

	i.currentState = StateSubmit
	return i.performPostSavePersistence(i.reviewState.Event, facts, convertedSkills, bursts)
}

// handleSubmitConfirmData processes the final submit confirmation data.
//
// Returns:
//   - nil after completing the intent, or a failure command.
//
// Side effects:
//   - Sets the intent result and marks it inactive on success.
func (i *Intent) handleSubmitConfirmData(data interface{}) tea.Cmd {
	reviewResult, ok := data.(event.ReviewResult)
	if !ok {
		return i.setFailedCmd("INVALID_SUBMIT_DATA", fmt.Sprintf("Invalid submit data type: %T", data), nil)
	}

	submittedEvent := domainEventFromDisplayEvent(reviewResult.Event)
	bursts := displayBurstsToPointers(reviewResult.Bursts, i.reviewState.InferredBursts)
	facts := displayFactsToPointers(reviewResult.Facts, i.reviewState.InferredFacts)
	convertedSkills := skillsFromDisplaySuggestions(reviewResult.Skills)

	i.result = &intents.IntentResult[*Result]{
		Status: intents.Completed,
		Data: &Result{
			Event:  submittedEvent,
			Bursts: bursts,
			Facts:  facts,
			Skills: convertedSkills,
		},
	}
	i.active = false
	return nil
}

// HandleError processes error results from views.
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
func (i *Intent) HandleError(result widgets.ViewResult) tea.Cmd {
	data := result.Data()
	if errorData, ok := data.(map[string]interface{}); ok {
		err, _ := errorData["error"].(error)
		msg, _ := errorData["message"].(string)
		if msg == "" {
			msg = "Unknown view error"
		}
		return i.setFailedCmd("VIEW_ERROR", msg, err)
	}
	return i.setFailedCmd("VIEW_ERROR", "Unknown view error", fmt.Errorf("%v", data))
}

// updateEditingModal handles messages when an editing modal is active in the
// view architecture path.
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
			i.reviewState.factSuggestionModal = nil
			i.reviewState.EditingMode = EditingModeNone
			return nil
		}
	}

	switch i.reviewState.EditingMode {
	case EditingModeMetadata:
		return i.updateMetadataModal(msg)
	case EditingModeBursts:
		return i.updateBurstReviewModal(msg)
	case EditingModeFacts:
		return i.updateFactReviewModal(msg)
	case EditingModeSkills:
		return i.updateSkillReviewModal(msg)
	}

	return nil
}

//nolint:dupl // Similar pattern for burst, fact, and skill modals
func (i *Intent) updateBurstReviewModal(msg tea.Msg) tea.Cmd {
	if i.reviewState.burstModal == nil {
		return nil
	}
	if !i.reviewState.burstModal.HasSuggestions() {
		i.reviewState.burstModal = nil
		i.reviewState.EditingMode = EditingModeNone
		return nil
	}
	model, cmd := i.reviewState.burstModal.Update(msg)
	if typed, ok := model.(*burstviews.SuggestionReview); ok {
		i.reviewState.burstModal = typed
	}
	if i.reviewState.burstModal.IsVisible() {
		return cmd
	}
	i.processAcceptedBursts()
	i.reviewState.burstModal = nil
	i.reviewState.EditingMode = EditingModeNone
	return cmd
}

// processAcceptedBursts transfers accepted burst suggestions to the review state.
//
// Side effects:
//   - Appends matched inferred bursts to AcceptedBursts.
//   - Updates the review view with accepted bursts if available.
func (i *Intent) processAcceptedBursts() {
	accepted := i.reviewState.burstModal.GetAcceptedSuggestions()
	for _, s := range accepted {
		name := s.Name
		if name == "" {
			name = fmt.Sprintf("Burst of %d events", len(s.EventIDs))
		}
		original := findInferredBurst(i.reviewState.InferredBursts, name)
		if original == nil {
			continue
		}
		i.reviewState.AcceptedBursts = append(i.reviewState.AcceptedBursts, original)
	}
	if reviewView, ok := i.activeView.(reviewUpdater); ok {
		reviewView.SetAcceptedBursts(display.BurstsFromDomain(i.reviewState.AcceptedBursts))
	}
}

func (i *Intent) updateMetadataModal(msg tea.Msg) tea.Cmd {
	if i.reviewState.metadataModal != nil {
		cmd, _ := i.reviewState.metadataModal.Update(msg)
		if i.reviewState.metadataModal.IsSubmitted() {
			i.reviewState.Event = domainEventFromDisplayEvent(i.reviewState.metadataModal.GetEvent())
			i.reviewState.metadataModal = nil
			i.reviewState.EditingMode = EditingModeNone
		} else if i.reviewState.metadataModal.IsCancelled() {
			i.reviewState.metadataModal = nil
			i.reviewState.EditingMode = EditingModeNone
		}
		return cmd
	}
	return nil
}

// updateFactReviewModal handles the fact suggestion review modal updates.
//
// Expected:
//   - msg must be a valid tea.Msg.
//
// Returns:
//   - A tea.Cmd if the modal processes the message.
//
// Side effects:
//   - May update reviewState.AcceptedFacts and the review view.
//
//nolint:dupl // Similar pattern for burst, fact, and skill modals
func (i *Intent) updateFactReviewModal(msg tea.Msg) tea.Cmd {
	if i.reviewState.factSuggestionModal == nil {
		return nil
	}
	if !i.reviewState.factSuggestionModal.HasSuggestions() {
		i.reviewState.factSuggestionModal = nil
		i.reviewState.EditingMode = EditingModeNone
		return nil
	}
	model, cmd := i.reviewState.factSuggestionModal.Update(msg)
	if typed, ok := model.(*burstviews.SuggestionReview); ok {
		i.reviewState.factSuggestionModal = typed
	}
	if i.reviewState.factSuggestionModal.IsVisible() {
		return cmd
	}
	i.processAcceptedFacts()
	i.reviewState.factSuggestionModal = nil
	i.reviewState.EditingMode = EditingModeNone
	return cmd
}

// processAcceptedFacts transfers accepted fact suggestions to the review state.
//
// Side effects:
//   - Appends matched facts to AcceptedFacts.
//   - Updates the review view with accepted facts if available.
func (i *Intent) processAcceptedFacts() {
	accepted := i.reviewState.factSuggestionModal.GetAcceptedFacts()
	if len(accepted) > 0 {
		i.reviewState.AcceptedFacts = append(
			i.reviewState.AcceptedFacts,
			displayFactsToPointers(accepted, i.reviewState.InferredFacts)...,
		)
		if reviewView, ok := i.activeView.(reviewUpdater); ok {
			reviewView.SetAcceptedFacts(display.FactsFromDomain(i.reviewState.AcceptedFacts))
		}
	}
}

// updateSkillReviewModal handles the skill suggestion review modal updates.
//
// Expected:
//   - msg must be a valid tea.Msg.
//
// Returns:
//   - A tea.Cmd if the modal processes the message.
//
// Side effects:
//   - May update reviewState.AcceptedSkills and the review view.
//
//nolint:dupl // Similar pattern for burst, fact, and skill modals
func (i *Intent) updateSkillReviewModal(msg tea.Msg) tea.Cmd {
	if i.reviewState.skillModal == nil {
		return nil
	}
	if !i.reviewState.skillModal.HasSuggestions() {
		i.reviewState.skillModal = nil
		i.reviewState.EditingMode = EditingModeNone
		return nil
	}
	model, cmd := i.reviewState.skillModal.Update(msg)
	if typed, ok := model.(*burstviews.SuggestionReview); ok {
		i.reviewState.skillModal = typed
	}
	if i.reviewState.skillModal.IsVisible() {
		return cmd
	}
	i.processAcceptedSkills()
	i.reviewState.skillModal = nil
	i.reviewState.EditingMode = EditingModeNone
	return cmd
}

// processAcceptedSkills transfers accepted skill suggestions to the review state.
//
// Side effects:
//   - Appends career.Skill entries to AcceptedSkills.
//   - Updates the review view with accepted skills if available.
func (i *Intent) processAcceptedSkills() {
	accepted := i.reviewState.skillModal.GetAcceptedSkills()
	for _, s := range accepted {
		if s.Name == "" {
			continue
		}
		i.reviewState.AcceptedSkills = append(i.reviewState.AcceptedSkills, &career.Skill{
			Name:     s.Name,
			Category: s.Category,
		})
	}
	if reviewView, ok := i.activeView.(reviewUpdater); ok {
		var acceptedSkills []display.SkillSuggestion
		for _, sk := range i.reviewState.AcceptedSkills {
			if sk.Name == "" {
				continue
			}
			acceptedSkills = append(acceptedSkills, display.SkillSuggestion{
				Name:       sk.Name,
				Category:   sk.Category,
				Confidence: 1.0,
			})
		}
		reviewView.SetAcceptedSkills(acceptedSkills)
	}
}

// eventFromFormData converts CaptureEventFormData to a career.Event.
// Delegates to the pure domain function capture.NewEventFromInput.
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
	return domcapture.NewEventFromInput(domcapture.EventInput{
		Text:       data.Text,
		Date:       data.Date,
		Company:    data.Company,
		Project:    data.Project,
		Tags:       data.Tags,
		Categories: data.Categories,
	})
}

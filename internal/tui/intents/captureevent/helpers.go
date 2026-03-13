package captureevent

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/constants"
	"github.com/baphled/kariya/internal/domain/career"
	repo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	burstfact "github.com/baphled/kariya/internal/service/career/burstfact"
	"github.com/baphled/kariya/internal/service/career/skillinference"
	"github.com/baphled/kariya/internal/tui/intents"
	"github.com/baphled/kariya/internal/tui/views/event"
	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/themes"
	"github.com/baphled/kariya/internal/ui/uikit/feedback"
	"github.com/baphled/kariya/internal/ui/uikit/primitives"
	tea "github.com/charmbracelet/bubbletea"
)

// terminalDimensions returns the current terminal dimensions for modal sizing.
// Falls back to nil if terminal info is unavailable, letting the model
// apply its own defaults.
func (i *Intent) terminalDimensions() *ReviewEnrichmentDimensions {
	info := i.GetTerminalInfo()
	if info == nil {
		return nil
	}
	return &ReviewEnrichmentDimensions{
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

// showValidationErrorModal displays a validation error modal without exiting the intent.
// This allows the user to see what went wrong and correct the input.
//
// Expected:
//   - message describes the validation failure in human-readable form.
//
// Returns:
//   - A tea.Cmd from the modal's Init method.
//
// Side effects:
//   - Creates an error modal on i.submitModal.
func (i *Intent) showValidationErrorModal(message string) tea.Cmd {
	i.submitModal = feedback.NewErrorModal("Validation Error", message)
	return i.submitModal.Init()
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
//   - i.context.CareerService is non-nil.
//
// Returns:
//   - A tea.Cmd that runs asynchronously and produces a SubmitCompleteMsg
//     on success or a SubmitErrorMsg on failure.
//
// Side effects:
//   - Calls CareerService.CaptureEvent to persist the event.
//   - Calls CareerService.SaveFact for each accepted fact without an ID.
//   - Sets a default date for quick-strategy events with a zero date.
//   - Persists accepted skills via careerService.SaveSkill and careerService.LinkSkillToEvent.
func (i *Intent) performSubmit() tea.Cmd {
	capturedEvent := i.reviewState.Event
	acceptedFacts := i.reviewState.AcceptedFacts
	acceptedSkills := i.reviewState.AcceptedSkills
	strategy := i.strategy
	careerService := i.context.CareerService

	return func() tea.Msg {
		if errMsg := submitValidateEvent(capturedEvent, careerService); errMsg != nil {
			return *errMsg
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if strategy == StrategyQuick && capturedEvent.Date.IsZero() {
			capturedEvent.Date = time.Now()
		}

		err := careerService.CaptureEvent(ctx, capturedEvent, careerservice.ManualEntry)
		if err != nil {
			return SubmitErrorMsg{
				Code:    "PERSISTENCE_ERROR",
				Message: fmt.Sprintf("Failed to save event: %v", err),
				Cause:   err,
			}
		}

		if errMsg := submitPersistFacts(ctx, acceptedFacts, capturedEvent.ID, careerService); errMsg != nil {
			return *errMsg
		}

		if errMsg := submitPersistSkills(ctx, acceptedSkills, capturedEvent.ID, careerService); errMsg != nil {
			return *errMsg
		}

		return SubmitCompleteMsg{}
	}
}

// submitValidateEvent checks the event and career service are valid for submission.
//
// Returns:
//   - A SubmitErrorMsg pointer if validation fails, or nil on success.
func submitValidateEvent(capturedEvent *career.Event, svc *careerservice.Service) *SubmitErrorMsg {
	if capturedEvent == nil {
		return &SubmitErrorMsg{Code: "MISSING_EVENT", Message: "No event data to submit"}
	}
	if err := capturedEvent.Validate(); err != nil {
		return &SubmitErrorMsg{Code: "VALIDATION_ERROR", Message: fmt.Sprintf("Event validation failed: %v", err), Cause: err}
	}
	if svc == nil {
		return &SubmitErrorMsg{Code: "SERVICE_ERROR", Message: "Career service not initialized"}
	}
	return nil
}

// submitPersistFacts saves accepted facts that lack an ID and links them to the event.
//
// Returns:
//   - A SubmitErrorMsg pointer if any fact save failed, or nil on success.
func submitPersistFacts(ctx context.Context, facts []*career.Fact, eventID string, svc *careerservice.Service) *SubmitErrorMsg {
	var factErrors []string
	for _, fact := range facts {
		if fact.ID != "" {
			continue
		}
		if fact.SourceEventID == "" {
			fact.SourceEventID = eventID
		}
		if err := svc.SaveFact(ctx, fact); err != nil {
			factErrors = append(factErrors, err.Error())
		}
	}
	if len(factErrors) > 0 {
		return &SubmitErrorMsg{
			Code:    "PARTIAL_SAVE",
			Message: fmt.Sprintf("Event saved but %d fact(s) failed: %s", len(factErrors), strings.Join(factErrors, "; ")),
			Cause:   fmt.Errorf("fact save failures: %s", strings.Join(factErrors, "; ")),
		}
	}
	return nil
}

// submitPersistSkills saves accepted skills and links them to the event.
//
// Returns:
//   - A SubmitErrorMsg pointer if any skill operation failed, or nil on success.
func submitPersistSkills(ctx context.Context, skills []*career.Skill, eventID string, svc *careerservice.Service) *SubmitErrorMsg {
	var skillErrors []string
	for _, skill := range skills {
		if err := svc.SaveSkill(ctx, skill); err != nil {
			skillErrors = append(skillErrors, err.Error())
			continue
		}
		if err := svc.LinkSkillToEvent(ctx, eventID, skill.ID); err != nil {
			skillErrors = append(skillErrors, err.Error())
		}
	}
	if len(skillErrors) > 0 {
		return &SubmitErrorMsg{
			Code:    "SKILL_SAVE_ERROR",
			Message: fmt.Sprintf("Event saved but %d skill(s) failed: %s", len(skillErrors), strings.Join(skillErrors, "; ")),
			Cause:   fmt.Errorf("skill save failures: %s", strings.Join(skillErrors, "; ")),
		}
	}
	return nil
}

// performPostSavePersistence persists skills, facts, and confirms bursts after
// the initial event save during post-save review.
//
// Expected:
//   - event is non-nil and has a valid ID.
//   - facts, skills, bursts may be empty or nil.
//
// Returns:
//   - A tea.Cmd that runs asynchronously and produces a
//     PostSavePersistenceCompleteMsg on success or a SubmitErrorMsg on failure.
//
// Side effects:
//   - Saves skills via careerService.SaveSkill and links them via careerService.LinkSkillToEvent.
//   - Saves facts via careerService.SaveFact (sets SourceEventID if empty).
//   - Confirms bursts via careerService.ConfirmBurst.
func (i *Intent) performPostSavePersistence(
	savedEvent *career.Event,
	facts []*career.Fact,
	skills []*career.Skill,
	bursts []*career.Burst,
) tea.Cmd {
	careerService := i.context.CareerService

	return func() tea.Msg {
		if careerService == nil {
			return PostSavePersistenceCompleteMsg{
				Event:  savedEvent,
				Bursts: bursts,
				Facts:  facts,
				Skills: skills,
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if errMsg, ok := persistPostSaveSkills(ctx, careerService, savedEvent.ID, skills); ok {
			return errMsg
		}
		if errMsg, ok := persistPostSaveFacts(ctx, careerService, savedEvent.ID, facts); ok {
			return errMsg
		}
		if errMsg, ok := confirmPostSaveBursts(ctx, careerService, bursts); ok {
			return errMsg
		}

		return PostSavePersistenceCompleteMsg{
			Event:  savedEvent,
			Bursts: bursts,
			Facts:  facts,
			Skills: skills,
		}
	}
}

// persistPostSaveSkills persists skills and links them to the event.
//
// Expected:
//   - ctx is non-nil.
//   - svc is non-nil.
//   - eventID is non-empty for successful linking.
//
// Returns:
//   - SubmitErrorMsg and true when a save/link fails.
//   - Zero value and false when all skills are processed.
//
// Side effects:
//   - Saves skills and links them to the event.
func persistPostSaveSkills(
	ctx context.Context,
	svc *careerservice.Service,
	eventID string,
	skills []*career.Skill,
) (SubmitErrorMsg, bool) {
	for _, skill := range skills {
		if err := svc.SaveSkill(ctx, skill); err != nil {
			return SubmitErrorMsg{
				Code:    "SKILL_SAVE_ERROR",
				Message: fmt.Sprintf("Failed to save skill: %v", err),
				Cause:   err,
			}, true
		}
		if err := svc.LinkSkillToEvent(ctx, eventID, skill.ID); err != nil {
			return SubmitErrorMsg{
				Code:    "SKILL_LINK_ERROR",
				Message: fmt.Sprintf("Failed to link skill: %v", err),
				Cause:   err,
			}, true
		}
	}
	return SubmitErrorMsg{}, false
}

// persistPostSaveFacts saves newly created facts and sets their source event.
//
// Expected:
//   - ctx is non-nil.
//   - svc is non-nil.
//   - eventID is non-empty for new fact linkage.
//
// Returns:
//   - SubmitErrorMsg and true when a save fails.
//   - Zero value and false when all facts are processed.
//
// Side effects:
//   - Sets SourceEventID for new facts and persists them.
func persistPostSaveFacts(
	ctx context.Context,
	svc *careerservice.Service,
	eventID string,
	facts []*career.Fact,
) (SubmitErrorMsg, bool) {
	for _, fact := range facts {
		if fact.ID != "" {
			continue
		}
		fact.SourceEventID = eventID
		if err := svc.SaveFact(ctx, fact); err != nil {
			return SubmitErrorMsg{
				Code:    "FACT_SAVE_ERROR",
				Message: fmt.Sprintf("Failed to save fact: %v", err),
				Cause:   err,
			}, true
		}
	}
	return SubmitErrorMsg{}, false
}

// confirmPostSaveBursts confirms bursts after post-save review.
//
// Expected:
//   - ctx is non-nil.
//   - svc is non-nil.
//
// Returns:
//   - SubmitErrorMsg and true when a confirm fails.
//   - Zero value and false when all bursts are confirmed.
//
// Side effects:
//   - Calls ConfirmBurst for each burst.
func confirmPostSaveBursts(
	ctx context.Context,
	svc *careerservice.Service,
	bursts []*career.Burst,
) (SubmitErrorMsg, bool) {
	for _, burst := range bursts {
		if err := svc.ConfirmBurst(ctx, burst); err != nil {
			return SubmitErrorMsg{
				Code:    "BURST_CONFIRM_ERROR",
				Message: fmt.Sprintf("Failed to confirm burst: %v", err),
				Cause:   err,
			}, true
		}
	}
	return SubmitErrorMsg{}, false
}

// performInference runs LLM inference calls (skill inference, burst suggestion,
// fact extraction) in a background tea.Cmd with its own timeout. This is
// decoupled from performSubmit so that event persistence completes promptly
// without waiting for potentially slow LLM calls.
//
// Expected:
//   - i.reviewState.Event is non-nil (event was already saved).
//   - i.context.CareerService may be nil (inference is skipped).
//
// Returns:
//   - A tea.Cmd that runs asynchronously and produces an InferenceCompleteMsg.
//
// Side effects:
//   - Calls SkillInferenceService.InferSkillsFromEvents.
//   - Calls CareerService.SuggestBursts and ExtractFactsFromBurst/Event.
//   - May persist burst suggestions via SaveBurstSuggestions.
func (i *Intent) performInference() tea.Cmd {
	capturedEvent := i.reviewState.Event
	careerService := i.context.CareerService
	skillService := i.context.SkillInferenceService

	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		inferredSkills := inferSkillsFromEvent(ctx, capturedEvent, skillService)
		result := inferBurstsAndFacts(ctx, capturedEvent, careerService)

		return InferenceCompleteMsg{
			InferredSkills:           inferredSkills,
			InferredFacts:            result.facts,
			InferredBursts:           result.bursts,
			InferredBurstSuggestions: result.burstSuggestions,
		}
	}
}

type inferenceResult struct {
	facts            []*career.Fact
	bursts           []*career.Burst
	burstSuggestions []burstfact.BurstSuggestion
}

// inferSkillsFromEvent runs skill inference for a single event.
//
// Returns:
//   - Skill suggestions from the inference service, or nil if unavailable.
func inferSkillsFromEvent(
	ctx context.Context,
	evt *career.Event,
	svc skillinference.SkillInferenceService,
) []skillinference.SkillSuggestion {
	if svc == nil {
		return nil
	}
	inferResult, inferErr := svc.InferSkillsFromEvents(ctx, []*career.Event{evt})
	if inferErr != nil || inferResult == nil {
		return nil
	}
	return inferResult.Suggestions
}

// inferBurstsAndFacts runs burst suggestion and fact extraction for an event.
//
// Returns:
//   - An inferenceResult with discovered bursts, facts, and suggestions.
func inferBurstsAndFacts(ctx context.Context, evt *career.Event, svc *careerservice.Service) inferenceResult {
	if svc == nil {
		return inferenceResult{}
	}
	allEvents, listErr := svc.ListEvents(ctx, repo.EventListFilters{Limit: -1})
	if listErr != nil || len(allEvents) < 2 {
		return inferenceResult{facts: extractEventFacts(ctx, evt, svc)}
	}
	var allEventIDs []string
	for _, e := range allEvents {
		allEventIDs = append(allEventIDs, e.ID)
	}
	suggestions, suggestErr := svc.SuggestBursts(ctx, allEventIDs)
	if suggestErr != nil || len(suggestions) == 0 {
		return inferenceResult{facts: extractEventFacts(ctx, evt, svc)}
	}
	savedBursts, saveErr := svc.SaveBurstSuggestions(ctx, suggestions)
	if saveErr != nil || len(savedBursts) == 0 {
		return inferenceResult{facts: extractEventFacts(ctx, evt, svc)}
	}
	var allFacts []career.Fact
	for _, burst := range savedBursts {
		facts, extractErr := svc.ExtractFactsFromBurst(ctx, burst)
		if extractErr == nil {
			allFacts = append(allFacts, facts...)
		}
	}
	result := inferenceResult{bursts: savedBursts, burstSuggestions: suggestions, facts: factsToPointers(allFacts)}
	if len(result.facts) == 0 {
		result.facts = extractEventFacts(ctx, evt, svc)
	}
	return result
}

// extractEventFacts extracts facts directly from a single event.
//
// Returns:
//   - Fact pointers from event extraction, or nil on error.
func extractEventFacts(ctx context.Context, evt *career.Event, svc *careerservice.Service) []*career.Fact {
	eventFacts, extractErr := svc.ExtractFactsFromEvent(ctx, evt)
	if extractErr != nil {
		return nil
	}
	return factsToPointers(eventFacts)
}

// transitionToStrategyScreen creates the strategy selection screen and activates it.
//
// Returns:
//   - nil (always). The screen is rendered on the next View call.
//
// Side effects:
//   - Sets i.currentState to StateChooseStrategy.
//   - Creates a new StrategySelectScreen and assigns it to i.activeView.
//   - Configures the screen with terminal dimensions, theme, and logo.
func (i *Intent) transitionToStrategyScreen() tea.Cmd {
	i.currentState = StateChooseStrategy
	i.activeView = event.NewStrategySelect()
	return nil
}

// transitionToFormScreen creates the form screen and activates it for the given strategy.
//
// Expected:
//   - strategy is a valid CaptureStrategy (StrategyQuick or StrategyManual).
//
// Returns:
//   - nil (always). The screen is rendered on the next View call.
//
// Side effects:
//   - Sets i.currentState to StateForm and i.strategy.
//   - Creates a new EventFormScreen and assigns it to i.activeView.
//   - Configures the screen with terminal dimensions, theme, and logo.
func (i *Intent) transitionToFormScreen(strategy CaptureStrategy) tea.Cmd {
	i.currentState = StateForm
	i.strategy = strategy

	var prevEvent *career.Event
	if i.context.PreviousEvent != nil {
		prevEvent = i.context.PreviousEvent
	}

	i.activeView = event.NewForm(display.EventFromDomain(prevEvent), strategy)
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

// buildMetadataModal creates and assigns the metadata modal to reviewState.
//
// Side effects:
//   - Sets i.reviewState.metadataModal.
func (i *Intent) buildMetadataModal() {
	dims := i.terminalDimensions()
	var termWidth, termHeight int
	if dims != nil {
		termWidth = dims.TerminalWidth
		termHeight = dims.TerminalHeight
	}
	allTags := constants.AllEventTags()
	tags := make([]string, len(allTags))
	for idx, t := range allTags {
		tags[idx] = string(t)
	}
	allCats := constants.AllCompetencyCategories()
	cats := make([]string, len(allCats))
	for idx, c := range allCats {
		cats[idx] = string(c)
	}
	i.reviewState.metadataModal = event.NewReviewEnrichment(
		display.EventFromDomain(i.reviewState.Event),
		event.ReviewEnrichmentConfig{
			AvailableTags:       tags,
			AvailableCategories: cats,
			AvailableSkills:     constants.SkillCategoryStrings(),
			TerminalWidth:       termWidth,
			TerminalHeight:      termHeight,
		},
	)
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
		i.buildMetadataModal()
	}
	return &modalContentData{
		title:   i.reviewState.metadataModal.GetTitle(),
		content: i.reviewState.metadataModal.GetContent(),
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

	var width, height int
	info := i.GetTerminalInfo()
	if info != nil && info.Width > 0 && info.Height > 0 {
		width = info.Width
		height = info.Height
	} else {
		width = 80
		height = 24
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

// renderFormModalFooter returns a UIKit badge-styled footer for form views.
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

// extractSkillsFromReviewData converts the "skills" field from review data
// into []*career.Skill, supporting both []skillinference.SkillSuggestion
// and []*career.Skill input types.
//
// Expected:
//   - reviewData contains a "skills" key with either type.
//
// Returns:
//   - Converted []*career.Skill slice, empty-named entries excluded.
//   - nil if "skills" key is missing or unrecognised type.
//
// Side effects: None.
func extractSkillsFromReviewData(reviewData map[string]interface{}) []*career.Skill {
	if suggestions, ok := reviewData["skills"].([]display.SkillSuggestion); ok {
		return skillsFromDisplaySuggestions(suggestions)
	}
	if suggestions, ok := reviewData["skills"].([]skillinference.SkillSuggestion); ok {
		return skillsFromInferenceSuggestions(suggestions)
	}
	if skills, ok := reviewData["skills"].([]*career.Skill); ok {
		return skillsFromDomainSkills(skills)
	}
	return nil
}

// skillsFromDisplaySuggestions converts display skill suggestions to domain skills.
//
// Expected:
//   - suggestions may be nil or empty.
//
// Returns:
//   - Converted skills with empty names excluded.
//
// Side effects: None.
func skillsFromDisplaySuggestions(suggestions []display.SkillSuggestion) []*career.Skill {
	var result []*career.Skill
	for _, s := range suggestions {
		if s.Name == "" {
			continue
		}
		result = append(result, &career.Skill{
			Name:     s.Name,
			Category: s.Category,
		})
	}
	return result
}

// skillsFromInferenceSuggestions converts inference suggestions to domain skills.
//
// Expected:
//   - suggestions may be nil or empty.
//
// Returns:
//   - Converted skills with empty names excluded.
//
// Side effects: None.
func skillsFromInferenceSuggestions(suggestions []skillinference.SkillSuggestion) []*career.Skill {
	var result []*career.Skill
	for _, s := range suggestions {
		if s.Name == "" {
			continue
		}
		result = append(result, &career.Skill{
			Name:     s.Name,
			Category: s.Category,
		})
	}
	return result
}

// skillsFromDomainSkills filters domain skills with valid names.
//
// Expected:
//   - skills may be nil or empty.
//
// Returns:
//   - Filtered skill pointers with empty names excluded.
//
// Side effects: None.
func skillsFromDomainSkills(skills []*career.Skill) []*career.Skill {
	var result []*career.Skill
	for _, s := range skills {
		if s == nil || s.Name == "" {
			continue
		}
		result = append(result, s)
	}
	return result
}

// extractBurstsFromReviewData extracts bursts from review data, converting display.Burst to domain pointers.
func extractBurstsFromReviewData(reviewData map[string]interface{}, inferredBursts []*career.Burst) []*career.Burst {
	if bursts, ok := reviewData["bursts"].([]display.Burst); ok {
		return displayBurstsToPointers(bursts, inferredBursts)
	}

	if bursts, ok := reviewData["bursts"].([]*career.Burst); ok {
		return bursts
	}

	return nil
}

// extractFactsFromReviewData extracts facts from review data, converting display.Fact to domain pointers.
func extractFactsFromReviewData(reviewData map[string]interface{}, inferredFacts []*career.Fact) []*career.Fact {
	if facts, ok := reviewData["facts"].([]display.Fact); ok {
		return displayFactsToPointers(facts, inferredFacts)
	}

	if facts, ok := reviewData["facts"].([]*career.Fact); ok {
		return facts
	}

	return nil
}

// findInferredBurst returns the inferred burst whose Name matches the given name, or nil if not found.
func findInferredBurst(inferred []*career.Burst, name string) *career.Burst {
	for _, b := range inferred {
		if b.Name == name {
			return b
		}
	}
	return nil
}

// factsToPointers converts a slice of career.Fact values to a slice of pointers.
//
// Expected:
//   - facts may be nil or empty.
//
// Returns:
//   - A slice of pointers to each fact in the input.
//   - An empty slice (not nil) if facts is nil or empty.
//
// Side effects: None.
func factsToPointers(facts []career.Fact) []*career.Fact {
	if len(facts) == 0 {
		return []*career.Fact{}
	}

	result := make([]*career.Fact, len(facts))
	for i := range facts {
		result[i] = &facts[i]
	}
	return result
}

// displayFactsToPointers converts display.Fact slice to domain career.Fact pointers, matching by ID.
func displayFactsToPointers(displayFacts []display.Fact, inferredFacts []*career.Fact) []*career.Fact {
	if len(displayFacts) == 0 {
		return []*career.Fact{}
	}

	byID := make(map[string]*career.Fact, len(inferredFacts))
	for idx := range inferredFacts {
		fact := inferredFacts[idx]
		if fact != nil && fact.ID != "" {
			byID[fact.ID] = fact
		}
	}

	result := make([]*career.Fact, 0, len(displayFacts))
	for idx := range displayFacts {
		displayFact := displayFacts[idx]
		if fact, ok := byID[displayFact.ID]; ok {
			result = append(result, fact)
			continue
		}

		result = append(result, &career.Fact{
			ID:                   displayFact.ID,
			Text:                 displayFact.Text,
			CompetencyCategories: append([]string(nil), displayFact.CompetencyCategories...),
			RoleFit:              career.RoleFit(displayFact.RoleFit),
			AudienceRelevance:    append([]string(nil), displayFact.AudienceRelevance...),
			StrengthSignal:       displayFact.StrengthSignal,
			SourceEventID:        displayFact.SourceEventID,
			SourceBurstID:        displayFact.SourceBurstID,
			CreatedAt:            displayFact.CreatedAt,
			UpdatedAt:            displayFact.UpdatedAt,
		})
	}

	return result
}

// displayBurstsToPointers converts display.Burst slice to domain career.Burst pointers, matching by ID.
func displayBurstsToPointers(displayBursts []display.Burst, inferredBursts []*career.Burst) []*career.Burst {
	if len(displayBursts) == 0 {
		return []*career.Burst{}
	}

	byID := make(map[string]*career.Burst, len(inferredBursts))
	for idx := range inferredBursts {
		burst := inferredBursts[idx]
		if burst != nil && burst.ID != "" {
			byID[burst.ID] = burst
		}
	}

	result := make([]*career.Burst, 0, len(displayBursts))
	for idx := range displayBursts {
		displayBurst := displayBursts[idx]
		if burst, ok := byID[displayBurst.ID]; ok {
			result = append(result, burst)
			continue
		}

		result = append(result, &career.Burst{
			ID:          displayBurst.ID,
			Name:        displayBurst.Name,
			Description: displayBurst.Description,
			EventIDs:    append([]string(nil), displayBurst.EventIDs...),
			CreatedAt:   displayBurst.CreatedAt,
			UpdatedAt:   displayBurst.UpdatedAt,
		})
	}

	return result
}

// domainEventFromDisplayEvent converts a display.Event to a domain career.Event.
func domainEventFromDisplayEvent(displayEvent display.Event) *career.Event {
	return &career.Event{
		ID:         displayEvent.ID,
		Text:       displayEvent.Text,
		Date:       displayEvent.Date,
		Company:    displayEvent.Company,
		Project:    displayEvent.Project,
		Tags:       append([]string(nil), displayEvent.Tags...),
		Categories: append([]string(nil), displayEvent.Categories...),
		Skills:     append([]string(nil), displayEvent.Skills...),
		CreatedAt:  displayEvent.CreatedAt,
		UpdatedAt:  displayEvent.UpdatedAt,
	}
}

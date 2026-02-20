package captureevent

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/intents"
	captureScreens "github.com/baphled/kariya/internal/cli/screens/capture"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/domain/career"
	repo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	burstfact "github.com/baphled/kariya/internal/service/career/burstfact"
	"github.com/baphled/kariya/internal/service/career/skillinference"
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
//   - Persists accepted skills using the skill repository.
func (i *Intent) performSubmit() tea.Cmd {
	event := i.reviewState.Event
	acceptedFacts := i.reviewState.AcceptedFacts
	acceptedSkills := i.reviewState.AcceptedSkills
	strategy := i.strategy
	careerService := i.context.CareerService

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

		if careerService == nil {
			return SubmitErrorMsg{
				Code:    "SERVICE_ERROR",
				Message: "Career service not initialized",
				Cause:   nil,
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if strategy == StrategyQuick && event.Date.IsZero() {
			event.Date = time.Now()
		}

		mode := careerservice.ManualEntry

		err := careerService.CaptureEvent(ctx, event, mode)

		if err != nil {
			return SubmitErrorMsg{
				Code:    "PERSISTENCE_ERROR",
				Message: fmt.Sprintf("Failed to save event: %v", err),
				Cause:   err,
			}
		}

		var factErrors []string
		if len(acceptedFacts) > 0 {
			for _, fact := range acceptedFacts {
				if fact.ID == "" {
					if fact.SourceEventID == "" {
						fact.SourceEventID = event.ID
					}

					if err := careerService.SaveFact(ctx, fact); err != nil {
						factErrors = append(factErrors, err.Error())
					}
				}
			}
		}

		if len(factErrors) > 0 {
			return SubmitErrorMsg{
				Code:    "PARTIAL_SAVE",
				Message: fmt.Sprintf("Event saved but %d fact(s) failed: %s", len(factErrors), strings.Join(factErrors, "; ")),
				Cause:   fmt.Errorf("fact save failures: %s", strings.Join(factErrors, "; ")),
			}
		}

		// Persist accepted skills using the skill repository.
		if len(acceptedSkills) > 0 && careerService != nil {
			skillRepo := careerService.GetSkillRepository()
			eventRepo := careerService.GetEventRepository()
			for _, skill := range acceptedSkills {
				if skill.ID == "" {
					if err := skillRepo.Create(ctx, skill); err != nil {
						return SubmitErrorMsg{
							Code:    "SKILL_SAVE_ERROR",
							Message: fmt.Sprintf("Failed to save skill %s: %v", skill.Name, err),
							Cause:   err,
						}
					}
				}
				if err := eventRepo.LinkSkill(ctx, event.ID, skill.ID); err != nil {
					return SubmitErrorMsg{
						Code:    "SKILL_LINK_ERROR",
						Message: fmt.Sprintf("Failed to link skill %s to event: %v", skill.Name, err),
						Cause:   err,
					}
				}
			}
		}

		return SubmitCompleteMsg{}
	}
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
//   - Creates skills via skillRepo.Create and links them via eventRepo.LinkSkill.
//   - Saves facts via careerService.SaveFact (sets SourceEventID if empty).
//   - Confirms bursts via careerService.ConfirmBurst.
func (i *Intent) performPostSavePersistence(
	event *career.Event,
	facts []*career.Fact,
	skills []*career.Skill,
	bursts []*career.Burst,
) tea.Cmd {
	careerService := i.context.CareerService

	return func() tea.Msg {
		if careerService == nil {
			return PostSavePersistenceCompleteMsg{
				Event:  event,
				Bursts: bursts,
				Facts:  facts,
				Skills: skills,
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		skillRepo := careerService.GetSkillRepository()
		eventRepo := careerService.GetEventRepository()

		for _, skill := range skills {
			if skill.ID == "" {
				if err := skillRepo.Create(ctx, skill); err != nil {
					return SubmitErrorMsg{
						Code:    "SKILL_SAVE_ERROR",
						Message: fmt.Sprintf("Failed to save skill: %v", err),
						Cause:   err,
					}
				}
			}
			if err := eventRepo.LinkSkill(ctx, event.ID, skill.ID); err != nil {
				return SubmitErrorMsg{
					Code:    "SKILL_LINK_ERROR",
					Message: fmt.Sprintf("Failed to link skill: %v", err),
					Cause:   err,
				}
			}
		}

		for _, fact := range facts {
			if fact.ID == "" {
				fact.SourceEventID = event.ID
				if err := careerService.SaveFact(ctx, fact); err != nil {
					return SubmitErrorMsg{
						Code:    "FACT_SAVE_ERROR",
						Message: fmt.Sprintf("Failed to save fact: %v", err),
						Cause:   err,
					}
				}
			}
		}

		for _, burst := range bursts {
			if err := careerService.ConfirmBurst(ctx, burst); err != nil {
				return SubmitErrorMsg{
					Code:    "BURST_CONFIRM_ERROR",
					Message: fmt.Sprintf("Failed to confirm burst: %v", err),
					Cause:   err,
				}
			}
		}

		return PostSavePersistenceCompleteMsg{
			Event:  event,
			Bursts: bursts,
			Facts:  facts,
			Skills: skills,
		}
	}
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
	event := i.reviewState.Event
	careerService := i.context.CareerService
	skillService := i.context.SkillInferenceService

	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var inferredSkills []skillinference.SkillSuggestion
		if skillService != nil {
			inferResult, inferErr := skillService.InferSkillsFromEvents(ctx, []*career.Event{event})
			if inferErr == nil && inferResult != nil {
				inferredSkills = inferResult.Suggestions
			}
		}

		var inferredFacts []*career.Fact
		var inferredBursts []*career.Burst
		var inferredBurstSuggestions []burstfact.BurstSuggestion

		if careerService != nil {
			allEvents, listErr := careerService.ListEvents(ctx, repo.EventListFilters{Limit: -1})
			if listErr == nil && len(allEvents) >= 2 {
				var allEventIDs []string
				for _, e := range allEvents {
					allEventIDs = append(allEventIDs, e.ID)
				}

				suggestions, suggestErr := careerService.SuggestBursts(ctx, allEventIDs)
				if suggestErr == nil && len(suggestions) > 0 {
					savedBursts, saveErr := careerService.SaveBurstSuggestions(ctx, suggestions)
					if saveErr == nil && len(savedBursts) > 0 {
						inferredBursts = savedBursts
						inferredBurstSuggestions = suggestions
						var allFacts []career.Fact
						for _, burst := range savedBursts {
							facts, extractErr := careerService.ExtractFactsFromBurst(ctx, burst)
							if extractErr == nil {
								allFacts = append(allFacts, facts...)
							}
						}
						inferredFacts = factsToPointers(allFacts)
					}
				}
			}

			if len(inferredFacts) == 0 {
				eventFacts, extractErr := careerService.ExtractFactsFromEvent(ctx, event)
				if extractErr == nil {
					inferredFacts = factsToPointers(eventFacts)
				}
			}
		}

		return InferenceCompleteMsg{
			InferredSkills:           inferredSkills,
			InferredFacts:            inferredFacts,
			InferredBursts:           inferredBursts,
			InferredBurstSuggestions: inferredBurstSuggestions,
		}
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
	if termInfo != nil && termInfo.Width > 0 && termInfo.Height > 0 {
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

	var event *career.Event
	if i.context.PreviousEvent != nil {
		event = i.context.PreviousEvent
	}

	formScreen := captureScreens.NewEventFormScreen(event, breadcrumbs, strategy)
	i.captureFormScreen = formScreen
	i.activeScreen = formScreen

	termInfo := i.GetTerminalInfo()
	width, height := 120, 40
	if termInfo != nil && termInfo.Width > 0 && termInfo.Height > 0 {
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
		i.reviewState.metadataModal = NewReviewEnrichmentModel(
			context.Background(),
			i.reviewState.Event,
			i.context.CareerService,
			i.context.CLIEventService,
			i.terminalDimensions(),
		)
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

	info := i.GetTerminalInfo()
	width := 80
	height := 40
	if info != nil && info.Width > 0 && info.Height > 0 {
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
	if suggestions, ok := reviewData["skills"].([]skillinference.SkillSuggestion); ok {
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

	if skills, ok := reviewData["skills"].([]*career.Skill); ok {
		var result []*career.Skill
		for _, s := range skills {
			if s == nil || s.Name == "" {
				continue
			}
			result = append(result, s)
		}
		return result
	}

	return nil
}

// findInferredBurst returns the inferred burst whose Name matches name, or nil if not found.
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

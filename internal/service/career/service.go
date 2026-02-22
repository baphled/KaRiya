package career

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"

	domain "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	repo "github.com/baphled/kariya/internal/repository/career"
	burst_fact "github.com/baphled/kariya/internal/service/career/burstfact"
)

// EventCaptureMode defines the different ways events can be captured.
type EventCaptureMode string

const (
	// TimelineJournaling is for logging events in real-time.
	TimelineJournaling EventCaptureMode = "timeline"

	// CVBackfill is for importing events from existing CVs.
	CVBackfill EventCaptureMode = "cv_backfill"

	// ManualEntry is for manually adding individual events.
	ManualEntry EventCaptureMode = "manual"
)

// Service provides business logic for career event management.
type Service struct {
	repo      repo.EventRepository
	factRepo  repo.FactRepository
	burstRepo repo.BurstRepository
	skillRepo repo.SkillRepository
	logger    *logger.Logger
}

// NewService creates a new career event service.
//
// Expected:
//   - eventrepository must be valid.
//
// Returns:
//   - A fully initialized Service ready for use.
//
// Side effects:
//   - None.
func NewService(repository repo.EventRepository) *Service {
	return &Service{
		repo:   repository,
		logger: logger.DefaultLogger(),
	}
}

// SetFactRepository sets the fact repository (optional, for fact extraction features).
//
// Expected:
//   - factrepository must be valid.
//
// Side effects:
//   - None.
func (s *Service) SetFactRepository(factRepo repo.FactRepository) {
	s.factRepo = factRepo
}

// SetBurstRepository sets the burst repository (optional, for burst detection features).
//
// Expected:
//   - burstrepository must be valid.
//
// Side effects:
//   - None.
func (s *Service) SetBurstRepository(burstRepo repo.BurstRepository) {
	s.burstRepo = burstRepo
}

// SetSkillRepository sets the skill repository (optional, for user-defined skills features).
//
// Expected:
//   - skillrepository must be valid.
//
// Side effects:
//   - None.
func (s *Service) SetSkillRepository(skillRepo repo.SkillRepository) {
	s.skillRepo = skillRepo
}

// CaptureEvent adds a new career event with specified capture mode.
//
// Expected:
//   - event must be valid.
//   - eventcapturemode must be valid.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (s *Service) CaptureEvent(ctx context.Context, event *domain.Event, mode EventCaptureMode) error {
	// Validate the event
	if err := event.Validate(); err != nil {
		s.logger.
			WithFields(map[string]string{
				"event_text":       event.Text,
				"capture_mode":     string(mode),
				"validation_error": err.Error(),
			}).
			Error("Event validation failed")
		return err
	}

	// Apply mode-specific validation or transformations
	switch mode {
	case TimelineJournaling:
		// Ensure date is close to current time for timeline entries
		if time.Since(event.Date) > 30*24*time.Hour {
			s.logger.
				WithFields(map[string]string{
					"event_date": event.Date.String(),
					"mode":       string(mode),
				}).
				Warn("Timeline event outside 30-day window")
			return errors.New("timeline events must be recent (within 30 days)")
		}
	case CVBackfill:
		// For imported events, allow older dates
		s.logger.
			WithFields(map[string]string{
				"event_text": event.Text,
				"event_date": event.Date.String(),
			}).
			Info("Importing event from CV")
	case ManualEntry:
		// No additional constraints for manual entry
		s.logger.
			WithFields(map[string]string{
				"event_text": event.Text,
			}).
			Info("Manually entered event")
	default:
		s.logger.
			WithFields(map[string]string{
				"mode": string(mode),
			}).
			Error("Invalid event capture mode")
		return errors.New("invalid event capture mode")
	}

	// Generate UUID if not provided
	if event.ID == "" {
		event.ID = uuid.New().String()
		s.logger.
			WithFields(map[string]string{
				"event_id":   event.ID,
				"event_text": event.Text,
			}).
			Debug("Generated new UUID for event")
	}

	// Set timestamps
	now := time.Now()
	event.CreatedAt = now
	event.UpdatedAt = now

	// Persist the event
	if err := s.repo.Create(ctx, event); err != nil {
		s.logger.
			WithFields(map[string]string{
				"event_id":   event.ID,
				"event_text": event.Text,
				"error":      err.Error(),
			}).
			Error("Failed to persist event")
		return err
	}

	s.logger.
		WithFields(map[string]string{
			"event_id":     event.ID,
			"capture_mode": string(mode),
			"event_text":   event.Text,
		}).
		Info("Event captured successfully")

	return nil
}

// UpdateEvent modifies an existing career event.
//
// Expected:
//   - event must be valid.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (s *Service) UpdateEvent(ctx context.Context, event *domain.Event) error {
	// Validate the updated event
	if err := event.Validate(); err != nil {
		s.logger.
			WithFields(map[string]string{
				"event_id":         event.ID,
				"validation_error": err.Error(),
			}).
			Error("Event validation failed during update")
		return err
	}

	// Ensure the event exists before updating
	existingEvent, err := s.repo.GetByID(ctx, event.ID)
	if err != nil {
		s.logger.
			WithFields(map[string]string{
				"event_id": event.ID,
				"error":    err.Error(),
			}).
			Error("Failed to retrieve existing event for update")
		return err
	}

	// Preserve creation timestamp
	event.CreatedAt = existingEvent.CreatedAt
	// Update modification timestamp
	event.UpdatedAt = time.Now()

	updateErr := s.repo.Update(ctx, event)
	if updateErr != nil {
		s.logger.
			WithFields(map[string]string{
				"event_id": event.ID,
				"error":    updateErr.Error(),
			}).
			Error("Failed to update event")
	}

	return updateErr
}

// DeleteEvent removes a career event.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (s *Service) DeleteEvent(ctx context.Context, eventID string) error {
	deleteErr := s.repo.Delete(ctx, eventID)
	if deleteErr != nil {
		s.logger.
			WithFields(map[string]string{
				"event_id": eventID,
				"error":    deleteErr.Error(),
			}).
			Error("Failed to delete event")
	}
	return deleteErr
}

// ListEvents retrieves career events with optional filtering.
func (s *Service) ListEvents(ctx context.Context, filters repo.EventListFilters) ([]*domain.Event, error) {
	events, err := s.repo.List(ctx, filters)
	if err != nil {
		s.logger.
			WithFields(map[string]string{
				"filters": fmt.Sprintf("%+v", filters),
				"error":   err.Error(),
			}).
			Warn("Failed to list events")
	}
	return events, err
}

// CountEvents returns the total number of events matching filters.
func (s *Service) CountEvents(ctx context.Context, filters repo.EventListFilters) (int, error) {
	count, err := s.repo.Count(ctx, filters)
	if err != nil {
		s.logger.
			WithFields(map[string]string{
				"filters": fmt.Sprintf("%+v", filters),
				"error":   err.Error(),
			}).
			Warn("Failed to count events")
	}
	return count, err
}

// GetEventByID retrieves a specific event.
func (s *Service) GetEventByID(ctx context.Context, eventID string) (*domain.Event, error) {
	event, err := s.repo.GetByID(ctx, eventID)
	if err != nil {
		s.logger.
			WithFields(map[string]string{
				"event_id": eventID,
				"error":    err.Error(),
			}).
			Warn("Failed to retrieve event")
	}
	return event, err
}

// SuggestBursts detects and suggests bursts for provided event IDs.
func (s *Service) SuggestBursts(ctx context.Context, eventIDs []string) ([]burst_fact.BurstSuggestion, error) {
	opts := &burst_fact.DetectionOptions{
		MinConfidence:       0.6,
		TemporalWindow:      6 * 30 * 24 * time.Hour,
		MinEventCount:       2,
		MaxSuggestionsCount: 0,
	}
	return s.SuggestBurstsWithOptions(ctx, eventIDs, opts)
}

// SuggestBurstsWithOptions detects bursts with custom detection options.
func (s *Service) SuggestBurstsWithOptions(
	ctx context.Context,
	eventIDs []string,
	opts *burst_fact.DetectionOptions,
) ([]burst_fact.BurstSuggestion, error) {
	if len(eventIDs) < 2 {
		return []burst_fact.BurstSuggestion{}, nil
	}

	// Retrieve events from repository
	var events []domain.Event
	for _, id := range eventIDs {
		event, err := s.repo.GetByID(ctx, id)
		if err != nil {
			s.logger.Debug("Event not found for burst detection")
			continue
		}
		events = append(events, *event)
	}

	if len(events) < 2 {
		return []burst_fact.BurstSuggestion{}, nil
	}

	// Detect bursts
	detector := burst_fact.NewBurstDetector()
	suggestions, err := detector.DetectBursts(ctx, events, opts)
	if err != nil {
		s.logger.Warn("Burst detection failed")
		return nil, err
	}

	s.logger.Info("Burst suggestions generated")
	return suggestions, nil
}

// SaveBurst validates and persists a burst without setting confirmation state.
//
// Expected:
//   - burst must be valid.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (s *Service) SaveBurst(ctx context.Context, burst *domain.Burst) error {
	if burst == nil {
		return errors.New("burst cannot be nil")
	}

	if err := burst.Validate(); err != nil {
		s.logger.WithFields(map[string]string{"error": err.Error()}).Warn("Burst validation failed")
		return err
	}

	if s.burstRepo != nil {
		if err := s.burstRepo.Create(ctx, burst); err != nil {
			s.logger.WithFields(map[string]string{
				"burst_id": burst.ID,
				"error":    err.Error(),
			}).Warn("Failed to save burst")
			return fmt.Errorf("failed to save burst: %w", err)
		}

		s.logger.WithFields(map[string]string{
			"burst_id": burst.ID,
			"name":     burst.Name,
		}).Info("Burst saved")
	} else {
		s.logger.Info("Burst saved (no repository configured)")
	}

	return nil
}

// ConfirmBurst marks an existing burst as confirmed and updates the repository.
//
// Expected:
//   - burst must be valid.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (s *Service) ConfirmBurst(ctx context.Context, burst *domain.Burst) error {
	if burst == nil {
		return errors.New("burst cannot be nil")
	}

	if err := burst.Validate(); err != nil {
		s.logger.WithFields(map[string]string{"error": err.Error()}).Warn("Burst validation failed")
		return err
	}

	// Capture original state for rollback on failure.
	origConfirmed := burst.Confirmed
	origConfirmedAt := burst.ConfirmedAt
	origUpdatedAt := burst.UpdatedAt

	now := time.Now()
	burst.Confirmed = true
	burst.ConfirmedAt = &now
	burst.UpdatedAt = now

	if s.burstRepo != nil {
		if err := s.burstRepo.Update(ctx, burst); err != nil {
			// Rollback in-memory changes on failure.
			burst.Confirmed = origConfirmed
			burst.ConfirmedAt = origConfirmedAt
			burst.UpdatedAt = origUpdatedAt
			s.logger.WithFields(map[string]string{
				"burst_id": burst.ID,
				"error":    err.Error(),
			}).Warn("Failed to confirm burst")
			return fmt.Errorf("failed to confirm burst: %w", err)
		}

		s.logger.WithFields(map[string]string{
			"burst_id": burst.ID,
			"name":     burst.Name,
		}).Info("Burst confirmed")
	}

	return nil
}

// DeleteBurst removes a burst from the repository.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (s *Service) DeleteBurst(ctx context.Context, burstID string) error {
	if s.burstRepo == nil {
		// Burst repository is optional - return sentinel error without logging
		return ErrBurstRepositoryNotConfigured
	}

	if burstID == "" {
		s.logger.Warn("Cannot delete burst with empty ID")
		return errors.New("burst ID cannot be empty")
	}

	if err := s.burstRepo.Delete(ctx, burstID); err != nil {
		s.logger.
			WithFields(map[string]string{
				"burst_id": burstID,
				"error":    err.Error(),
			}).
			Error("Failed to delete burst")
		return fmt.Errorf("failed to delete burst: %w", err)
	}

	s.logger.
		WithFields(map[string]string{
			"burst_id": burstID,
		}).
		Info("Burst deleted successfully")

	return nil
}

// RejectBurstSuggestion records rejection of burst suggestion.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (s *Service) RejectBurstSuggestion(_ context.Context, eventIDs []string) error {
	if len(eventIDs) == 0 {
		return nil
	}
	s.logger.Debug("Burst suggestion rejected")
	return nil
}

// SaveBurstSuggestions converts burst suggestions to actual bursts and saves them.
func (s *Service) SaveBurstSuggestions(ctx context.Context, suggestions []burst_fact.BurstSuggestion) ([]*domain.Burst, error) {
	if s.burstRepo == nil {
		s.logger.Debug("No burst repository configured - suggestions not saved")
		return nil, nil
	}

	var savedBursts []*domain.Burst

	for _, suggestion := range suggestions {
		// Convert suggestion to burst domain object
		burst := &domain.Burst{
			ID:          uuid.New().String(),
			Name:        suggestion.Name,
			Description: suggestion.Description,
			EventIDs:    suggestion.EventIDs,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		// Save to repository (unconfirmed - user must review and confirm).
		if err := s.SaveBurst(ctx, burst); err != nil {
			s.logger.WithFields(map[string]string{
				"suggestion_name": suggestion.Name,
				"error":           err.Error(),
			}).Warn("Failed to save burst suggestion")
			continue
		}

		savedBursts = append(savedBursts, burst)
	}

	if len(savedBursts) > 0 {
		s.logger.WithFields(map[string]string{
			"saved_count":     strconv.Itoa(len(savedBursts)),
			"suggested_count": strconv.Itoa(len(suggestions)),
		}).Info("Burst suggestions saved as persistent bursts")
	}

	return savedBursts, nil
}

// ExtractFactsFromEvent extracts facts from a single career event.
func (s *Service) ExtractFactsFromEvent(ctx context.Context, event *domain.Event) ([]domain.Fact, error) {
	if event == nil {
		s.logger.Warn("Cannot extract facts from nil event")
		return nil, errors.New("event cannot be nil")
	}

	if event.ID == "" {
		s.logger.Warn("Cannot extract facts from event with empty ID")
		return nil, errors.New("event ID cannot be empty")
	}

	// Create extractor and classifier
	classifier := burst_fact.NewClassifier()
	extractor := burst_fact.NewExtractor(classifier)

	// Extract facts
	facts := extractor.ExtractFromEvent(ctx, event)

	if len(facts) == 0 {
		s.logger.Debug("No facts extracted from event")
		return []domain.Fact{}, nil
	}

	s.logger.
		WithFields(map[string]string{
			"event_id":   event.ID,
			"fact_count": strconv.Itoa(len(facts)),
			"event_text": event.Text,
		}).
		Info("Facts extracted from event")

	return facts, nil
}

// ExtractFactsFromBurst extracts facts from a burst (multiple related events).
func (s *Service) ExtractFactsFromBurst(ctx context.Context, burst *domain.Burst) ([]domain.Fact, error) {
	if burst == nil {
		s.logger.Warn("Cannot extract facts from nil burst")
		return nil, errors.New("burst cannot be nil")
	}

	if burst.ID == "" {
		s.logger.Warn("Cannot extract facts from burst with empty ID")
		return nil, errors.New("burst ID cannot be empty")
	}

	if len(burst.EventIDs) == 0 {
		s.logger.Debug("Burst has no events")
		return []domain.Fact{}, nil
	}

	// Retrieve all events in burst
	var events []*domain.Event
	for _, eventID := range burst.EventIDs {
		event, err := s.repo.GetByID(ctx, eventID)
		if err != nil {
			s.logger.
				WithFields(map[string]string{
					"event_id": eventID,
					"burst_id": burst.ID,
					"error":    err.Error(),
				}).
				Debug("Event not found for burst fact extraction")
			continue
		}
		events = append(events, event)
	}

	if len(events) == 0 {
		s.logger.
			WithFields(map[string]string{
				"burst_id":    burst.ID,
				"event_count": strconv.Itoa(len(burst.EventIDs)),
			}).
			Warn("No events found for burst fact extraction")
		return []domain.Fact{}, nil
	}

	// Create extractor and classifier
	classifier := burst_fact.NewClassifier()
	extractor := burst_fact.NewExtractor(classifier)

	// Extract facts
	facts := extractor.ExtractFromBurst(ctx, burst, events)

	if len(facts) == 0 {
		s.logger.Debug("No facts extracted from burst")
		return []domain.Fact{}, nil
	}

	s.logger.
		WithFields(map[string]string{
			"burst_id":    burst.ID,
			"burst_name":  burst.Name,
			"fact_count":  strconv.Itoa(len(facts)),
			"event_count": strconv.Itoa(len(events)),
		}).
		Info("Facts extracted from burst")

	return facts, nil
}

// ValidateFact checks if a fact meets all validation criteria.
//
// Expected:
//   - fact must be valid.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (s *Service) ValidateFact(_ context.Context, fact *domain.Fact) error {
	if fact == nil {
		s.logger.Warn("Cannot validate nil fact")
		return errors.New("fact cannot be nil")
	}

	// Perform domain validation
	if err := fact.Validate(); err != nil {
		s.logger.
			WithFields(map[string]string{
				"fact_id":          fact.ID,
				"fact_text":        fact.Text,
				"validation_error": err.Error(),
			}).
			Warn("Fact validation failed")
		return err
	}

	s.logger.
		WithFields(map[string]string{
			"fact_id":   fact.ID,
			"fact_text": fact.Text,
			"role_fit":  string(fact.RoleFit),
		}).
		Debug("Fact validated successfully")

	return nil
}

// GetFactsBySourceEventID retrieves all facts extracted from a specific event.
func (s *Service) GetFactsBySourceEventID(ctx context.Context, eventID string) ([]*domain.Fact, error) {
	if s.factRepo == nil {
		s.logger.Debug("Fact repository not configured")
		return []*domain.Fact{}, nil
	}

	if eventID == "" {
		s.logger.Warn("Cannot get facts for empty event ID")
		return nil, errors.New("event ID cannot be empty")
	}

	facts, err := s.factRepo.GetBySourceEventID(ctx, eventID)
	if err != nil {
		s.logger.
			WithFields(map[string]string{
				"event_id": eventID,
				"error":    err.Error(),
			}).
			Error("Failed to retrieve facts for event")
		return nil, fmt.Errorf("failed to retrieve facts: %w", err)
	}

	s.logger.
		WithFields(map[string]string{
			"event_id":   eventID,
			"fact_count": strconv.Itoa(len(facts)),
		}).
		Debug("Facts retrieved for event")

	return facts, nil
}

// GetFactsBySourceBurstID retrieves all facts extracted from a specific burst.
func (s *Service) GetFactsBySourceBurstID(ctx context.Context, burstID string) ([]*domain.Fact, error) {
	if s.factRepo == nil {
		s.logger.Debug("Fact repository not configured")
		return []*domain.Fact{}, nil
	}

	if burstID == "" {
		s.logger.Warn("Cannot get facts for empty burst ID")
		return nil, errors.New("burst ID cannot be empty")
	}

	facts, err := s.factRepo.GetBySourceBurstID(ctx, burstID)
	if err != nil {
		s.logger.
			WithFields(map[string]string{
				"burst_id": burstID,
				"error":    err.Error(),
			}).
			Error("Failed to retrieve facts for burst")
		return nil, fmt.Errorf("failed to retrieve facts: %w", err)
	}

	s.logger.
		WithFields(map[string]string{
			"burst_id":   burstID,
			"fact_count": strconv.Itoa(len(facts)),
		}).
		Debug("Facts retrieved for burst")

	return facts, nil
}

// SaveFact persists a fact to the repository.
//
// Expected:
//   - fact must be valid.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (s *Service) SaveFact(ctx context.Context, fact *domain.Fact) error {
	if s.factRepo == nil {
		// Fact repository is optional - return sentinel error without logging
		return ErrFactRepositoryNotConfigured
	}

	if fact == nil {
		s.logger.Warn("Cannot save nil fact")
		return errors.New("fact cannot be nil")
	}

	// Generate UUID if not provided (must be done before validation)
	if fact.ID == "" {
		fact.ID = uuid.New().String()
	}

	// Set timestamps (must be done before validation)
	now := time.Now()
	if fact.CreatedAt.IsZero() {
		fact.CreatedAt = now
	}
	fact.UpdatedAt = now

	// Validate fact before saving (after ID and timestamps are set)
	if err := s.ValidateFact(ctx, fact); err != nil {
		return err
	}

	// Check if fact already exists
	existing, err := s.factRepo.GetByID(ctx, fact.ID)
	if err == nil && existing != nil {
		// Update existing fact
		if err := s.factRepo.Update(ctx, fact); err != nil {
			s.logger.
				WithFields(map[string]string{
					"fact_id": fact.ID,
					"error":   err.Error(),
				}).
				Error("Failed to update fact")
			return fmt.Errorf("failed to update fact: %w", err)
		}

		s.logger.
			WithFields(map[string]string{
				"fact_id": fact.ID,
			}).
			Info("Fact updated successfully")

		return nil
	}

	// Create new fact
	if err := s.factRepo.Create(ctx, fact); err != nil {
		s.logger.
			WithFields(map[string]string{
				"fact_id": fact.ID,
				"error":   err.Error(),
			}).
			Error("Failed to create fact")
		return fmt.Errorf("failed to create fact: %w", err)
	}

	s.logger.
		WithFields(map[string]string{
			"fact_id": fact.ID,
		}).
		Info("Fact created successfully")

	return nil
}

// DeleteFact removes a fact from the repository.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (s *Service) DeleteFact(ctx context.Context, factID string) error {
	if s.factRepo == nil {
		// Fact repository is optional - return sentinel error without logging
		return ErrFactRepositoryNotConfigured
	}

	if factID == "" {
		s.logger.Warn("Cannot delete fact with empty ID")
		return errors.New("fact ID cannot be empty")
	}

	if err := s.factRepo.Delete(ctx, factID); err != nil {
		s.logger.
			WithFields(map[string]string{
				"fact_id": factID,
				"error":   err.Error(),
			}).
			Error("Failed to delete fact")
		return fmt.Errorf("failed to delete fact: %w", err)
	}

	s.logger.
		WithFields(map[string]string{
			"fact_id": factID,
		}).
		Info("Fact deleted successfully")

	return nil
}

// SaveSkill persists a new skill if it does not yet have an ID.
//
// Expected:
//   - skill is non-nil.
//
// Returns:
//   - An error value.
//
// Side effects:
//   - Creates a new skill record via the skill repository.
func (s *Service) SaveSkill(ctx context.Context, skill *domain.Skill) error {
	if s.skillRepo == nil {
		return ErrSkillRepositoryNotConfigured
	}

	if skill == nil {
		return errors.New("skill cannot be nil")
	}

	if skill.ID != "" {
		return nil
	}

	skill.ID = uuid.New().String()

	now := time.Now()
	if skill.CreatedAt.IsZero() {
		skill.CreatedAt = now
	}
	skill.UpdatedAt = now

	if err := s.skillRepo.Create(ctx, skill); err != nil {
		s.logger.
			WithFields(map[string]string{
				"skill_name": skill.Name,
				"error":      err.Error(),
			}).
			Error("Failed to create skill")
		return fmt.Errorf("failed to create skill: %w", err)
	}

	s.logger.
		WithFields(map[string]string{
			"skill_id":   skill.ID,
			"skill_name": skill.Name,
		}).
		Info("Skill created successfully")

	return nil
}

// LinkSkillToEvent creates an association between a skill and an event.
//
// Expected:
//   - eventID is a non-empty string.
//   - skillID is a non-empty string.
//
// Returns:
//   - An error value.
//
// Side effects:
//   - Creates an event-skill association via the event repository.
func (s *Service) LinkSkillToEvent(ctx context.Context, eventID, skillID string) error {
	if eventID == "" {
		return errors.New("event ID cannot be empty")
	}

	if skillID == "" {
		return errors.New("skill ID cannot be empty")
	}

	if err := s.repo.LinkSkill(ctx, eventID, skillID); err != nil {
		s.logger.
			WithFields(map[string]string{
				"event_id": eventID,
				"skill_id": skillID,
				"error":    err.Error(),
			}).
			Error("Failed to link skill to event")
		return fmt.Errorf("failed to link skill to event: %w", err)
	}

	s.logger.
		WithFields(map[string]string{
			"event_id": eventID,
			"skill_id": skillID,
		}).
		Info("Skill linked to event successfully")

	return nil
}

// GetBurstRepository returns the burst repository.
//
// Returns:
//   - A repo.BurstRepository value.
//
// Side effects:
//   - None.
func (s *Service) GetBurstRepository() repo.BurstRepository {
	return s.burstRepo
}

// GetFactRepository returns the fact repository.
//
// Returns:
//   - A repo.FactRepository value.
//
// Side effects:
//   - None.
func (s *Service) GetFactRepository() repo.FactRepository {
	return s.factRepo
}

// GetEventRepository returns the event repository.
//
// Returns:
//   - A repo.EventRepository value.
//
// Side effects:
//   - None.
func (s *Service) GetEventRepository() repo.EventRepository {
	return s.repo
}

// GetSkillRepository returns the skill repository.
//
// Returns:
//   - A repo.SkillRepository value.
//
// Side effects:
//   - None.
func (s *Service) GetSkillRepository() repo.SkillRepository {
	return s.skillRepo
}

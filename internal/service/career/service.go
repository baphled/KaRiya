package career

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	domain "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	repo "github.com/baphled/kariya/internal/repository/career"
	burst_fact "github.com/baphled/kariya/internal/service/career/burst_fact"
)

// EventCaptureMode defines the different ways events can be captured
type EventCaptureMode string

const (
	// TimelineJournaling is for logging events in real-time
	TimelineJournaling EventCaptureMode = "timeline"

	// CVBackfill is for importing events from existing CVs
	CVBackfill EventCaptureMode = "cv_backfill"

	// ManualEntry is for manually adding individual events
	ManualEntry EventCaptureMode = "manual"
)

// Service provides business logic for career event management
type Service struct {
	repo   repo.Repository
	logger *logger.Logger
}

// NewService creates a new career event service
func NewService(repository repo.Repository) *Service {
	return &Service{
		repo:   repository,
		logger: logger.DefaultLogger(),
	}
}

// CaptureEvent adds a new career event with specified capture mode
func (s *Service) CaptureEvent(ctx context.Context, event *domain.CareerEvent, mode EventCaptureMode) error {
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

// UpdateEvent modifies an existing career event
func (s *Service) UpdateEvent(ctx context.Context, event *domain.CareerEvent) error {
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

// DeleteEvent removes a career event
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

// ListEvents retrieves career events with optional filtering
func (s *Service) ListEvents(ctx context.Context, filters repo.ListFilters) ([]*domain.CareerEvent, error) {
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

// CountEvents returns the total number of events matching filters
func (s *Service) CountEvents(ctx context.Context, filters repo.ListFilters) (int, error) {
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

// GetEventByID retrieves a specific event
func (s *Service) GetEventByID(ctx context.Context, eventID string) (*domain.CareerEvent, error) {
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
// SuggestBursts detects and suggests bursts for provided event IDs
func (s *Service) SuggestBursts(ctx context.Context, eventIDs []string) ([]burst_fact.BurstSuggestion, error) {
	opts := &burst_fact.DetectionOptions{
		MinConfidence:       0.6,
		TemporalWindow:      6 * 30 * 24 * time.Hour,
		MinEventCount:       2,
		MaxSuggestionsCount: 10,
	}
	return s.SuggestBurstsWithOptions(ctx, eventIDs, opts)
}

// SuggestBurstsWithOptions detects bursts with custom detection options
func (s *Service) SuggestBurstsWithOptions(
	ctx context.Context,
	eventIDs []string,
	opts *burst_fact.DetectionOptions,
) ([]burst_fact.BurstSuggestion, error) {
	if len(eventIDs) < 2 {
		return []burst_fact.BurstSuggestion{}, nil
	}

	// Retrieve events from repository
	var events []domain.CareerEvent
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

// ConfirmBurst validates a burst suggestion
func (s *Service) ConfirmBurst(ctx context.Context, burst *domain.Burst) error {
	if burst == nil {
		return fmt.Errorf("burst cannot be nil")
	}

	// Validate burst
	if err := burst.Validate(); err != nil {
		s.logger.Warn("Burst validation failed")
		return err
	}

	s.logger.Info("Burst confirmed")
	return nil
}

// RejectBurstSuggestion records rejection of burst suggestion
func (s *Service) RejectBurstSuggestion(ctx context.Context, eventIDs []string) error {
	if len(eventIDs) == 0 {
		return nil
	}
	s.logger.Debug("Burst suggestion rejected")
	return nil
}

// ExtractFactsFromEvent extracts facts from a single career event
func (s *Service) ExtractFactsFromEvent(ctx context.Context, event *domain.CareerEvent) ([]domain.Fact, error) {
	if event == nil {
		s.logger.Warn("Cannot extract facts from nil event")
		return nil, fmt.Errorf("event cannot be nil")
	}

	if event.ID == "" {
		s.logger.Warn("Cannot extract facts from event with empty ID")
		return nil, fmt.Errorf("event ID cannot be empty")
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
			"event_id":    event.ID,
			"fact_count":  fmt.Sprintf("%d", len(facts)),
			"event_text":  event.Text,
		}).
		Info("Facts extracted from event")

	return facts, nil
}

// ExtractFactsFromBurst extracts facts from a burst (multiple related events)
func (s *Service) ExtractFactsFromBurst(ctx context.Context, burst *domain.Burst) ([]domain.Fact, error) {
	if burst == nil {
		s.logger.Warn("Cannot extract facts from nil burst")
		return nil, fmt.Errorf("burst cannot be nil")
	}

	if burst.ID == "" {
		s.logger.Warn("Cannot extract facts from burst with empty ID")
		return nil, fmt.Errorf("burst ID cannot be empty")
	}

	if len(burst.EventIDs) == 0 {
		s.logger.Debug("Burst has no events")
		return []domain.Fact{}, nil
	}

	// Retrieve all events in burst
	var events []*domain.CareerEvent
	for _, eventID := range burst.EventIDs {
		event, err := s.repo.GetByID(ctx, eventID)
		if err != nil {
			s.logger.
				WithFields(map[string]string{
					"event_id":   eventID,
					"burst_id":   burst.ID,
					"error":      err.Error(),
				}).
				Debug("Event not found for burst fact extraction")
			continue
		}
		events = append(events, event)
	}

	if len(events) == 0 {
		s.logger.
			WithFields(map[string]string{
				"burst_id":     burst.ID,
				"event_count":  fmt.Sprintf("%d", len(burst.EventIDs)),
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
			"fact_count":  fmt.Sprintf("%d", len(facts)),
			"event_count": fmt.Sprintf("%d", len(events)),
		}).
		Info("Facts extracted from burst")

	return facts, nil
}

// ValidateFact checks if a fact meets all validation criteria
func (s *Service) ValidateFact(ctx context.Context, fact *domain.Fact) error {
	if fact == nil {
		s.logger.Warn("Cannot validate nil fact")
		return fmt.Errorf("fact cannot be nil")
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

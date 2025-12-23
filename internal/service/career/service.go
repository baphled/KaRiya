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

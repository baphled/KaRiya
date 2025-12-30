package service

import (
	"context"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
)

// CLIEventService wraps the existing career service to provide CLI-specific event operations
type CLIEventService struct {
	service *careerservice.Service
}

// NewCLIEventService creates a new CLI event service
func NewCLIEventService(service *careerservice.Service) *CLIEventService {
	return &CLIEventService{
		service: service,
	}
}

// CaptureEvent captures a new career event using the existing service
func (c *CLIEventService) CaptureEvent(ctx context.Context, text string, date time.Time, mode careerservice.EventCaptureMode, opts ...Option) error {
	// Apply default and optional configurations
	config := defaultConfig()
	for _, opt := range opts {
		opt(config)
	}

	// Create event
	event := &career.CareerEvent{
		Text:       text,
		Date:       date,
		Company:    config.Company,
		Project:    config.Project,
		Tags:       config.Tags,
		Categories: config.Categories,
	}

	// Capture event using existing service
	return c.service.CaptureEvent(ctx, event, mode)
}

// ListEvents retrieves a list of events with optional filtering
func (c *CLIEventService) ListEvents(ctx context.Context, filters *careerrepo.ListFilters) ([]*career.CareerEvent, error) {
	if filters == nil {
		filters = &careerrepo.ListFilters{}
	}
	return c.service.ListEvents(ctx, *filters)
}

// GetEventByID retrieves a specific event by its ID
func (c *CLIEventService) GetEventByID(ctx context.Context, eventID string) (*career.CareerEvent, error) {
	return c.service.GetEventByID(ctx, eventID)
}

// UpdateEvent updates an existing career event
func (c *CLIEventService) UpdateEvent(ctx context.Context, eventID string, text string, date time.Time, opts ...Option) error {
	// Apply default and optional configurations
	config := defaultConfig()
	for _, opt := range opts {
		opt(config)
	}

	// Create updated event
	event := &career.CareerEvent{
		ID:         eventID,
		Text:       text,
		Date:       date,
		Company:    config.Company,
		Project:    config.Project,
		Tags:       config.Tags,
		Categories: config.Categories,
	}

	// Update event using existing service
	return c.service.UpdateEvent(ctx, event)
}

// DeleteEvent deletes a career event by ID
func (c *CLIEventService) DeleteEvent(ctx context.Context, eventID string) error {
	return c.service.DeleteEvent(ctx, eventID)
}

// Configuration options for event capture
type Option func(*EventConfig)

// EventConfig holds optional configuration for event capture
type EventConfig struct {
	Company    string
	Project    string
	Tags       []string
	Categories []string
}

// defaultConfig provides default event configuration
func defaultConfig() *EventConfig {
	return &EventConfig{
		Tags:       []string{},
		Categories: []string{},
	}
}

// WithCompany sets the company for the event
func WithCompany(company string) Option {
	return func(ec *EventConfig) {
		ec.Company = company
	}
}

// WithProject sets the project for the event
func WithProject(project string) Option {
	return func(ec *EventConfig) {
		ec.Project = project
	}
}

// WithTags sets tags for the event
func WithTags(tags []string) Option {
	return func(ec *EventConfig) {
		ec.Tags = tags
	}
}

// WithCategories sets categories for the event
func WithCategories(categories []string) Option {
	return func(ec *EventConfig) {
		ec.Categories = categories
	}
}

// UpdateEventMetadata updates only the metadata fields of an event (company, project, tags, categories)
// This is used by the metadata editor to update event metadata without changing the text or date
func (c *CLIEventService) UpdateEventMetadata(ctx context.Context, event *career.CareerEvent) error {
	if event == nil {
		return ErrNilEvent
	}

	if event.ID == "" {
		return ErrEmptyEventID
	}

	// Get the existing event first
	existingEvent, err := c.service.GetEventByID(ctx, event.ID)
	if err != nil {
		return err
	}

	// Preserve original text and date, only update metadata fields
	event.Text = existingEvent.Text
	event.Date = existingEvent.Date
	event.CreatedAt = existingEvent.CreatedAt

	// Update the event with the new metadata
	return c.service.UpdateEvent(ctx, event)
}

// Error definitions for metadata operations
var (
	ErrNilEvent       = NewMetadataError("event cannot be nil")
	ErrEmptyEventID   = NewMetadataError("event ID cannot be empty")
	ErrEmptyEventList = NewMetadataError("event list cannot be empty")
)

// MetadataError represents an error during metadata operations
type MetadataError struct {
	message string
}

// NewMetadataError creates a new metadata error
func NewMetadataError(message string) *MetadataError {
	return &MetadataError{message: message}
}

// Error implements the error interface
func (me *MetadataError) Error() string {
	return me.message
}

// BulkUpdateSummary contains the results of a bulk update operation
type BulkUpdateSummary struct {
	UpdatedCount int
	FailedCount  int
	TotalCount   int
}

// BulkUpdateMetadata updates metadata for multiple events with transaction-like behavior
// If any event fails validation, no events are updated
func (c *CLIEventService) BulkUpdateMetadata(ctx context.Context, eventIDs []string, company, project string, tags, categories []string) (*BulkUpdateSummary, error) {
	// Validate input
	if len(eventIDs) == 0 {
		return nil, ErrEmptyEventList
	}

	// Load all events and validate before updating any
	events := make([]*career.CareerEvent, 0, len(eventIDs))
	for _, id := range eventIDs {
		event, err := c.service.GetEventByID(ctx, id)
		if err != nil {
			// Event not found or error retrieving - fail the entire operation
			return nil, err
		}
		events = append(events, event)
	}

	// Validate all events before updating any (transaction-like behavior)
	for _, event := range events {
		if err := event.Validate(); err != nil {
			return nil, err
		}
	}

	// All validations passed - update all events
	updatedCount := 0
	for _, event := range events {
		// Update metadata fields
		if company != "" {
			event.Company = company
		}
		if project != "" {
			event.Project = project
		}
		if len(tags) > 0 {
			event.Tags = tags
		}
		if len(categories) > 0 {
			event.Categories = categories
		}
		// Update timestamp
		event.UpdatedAt = time.Now()

		// Persist the update
		if err := c.service.UpdateEvent(ctx, event); err != nil {
			// In case of error, return error (operation is not truly atomic, but we validate first)
			return nil, err
		}
		updatedCount++
	}

	return &BulkUpdateSummary{
		UpdatedCount: updatedCount,
		FailedCount:  0,
		TotalCount:   len(eventIDs),
	}, nil
}

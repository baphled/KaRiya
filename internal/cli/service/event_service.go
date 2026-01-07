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

// BulkMetadataUpdate represents metadata fields to update in bulk
type BulkMetadataUpdate struct {
	Company             string
	ApplyIfEmptyCompany bool
	Project             string
	ApplyIfEmptyProject bool
	Tags                []string
	Categories          []string
}

// BulkOperationsSummary contains the results of a bulk update operation
type BulkOperationsSummary struct {
	EventsAffected int
	FieldsUpdated  []string
	Errors         []string
	AppliedCount   int
	SkippedCount   int
}

// BulkUpdateMetadata updates metadata for multiple events with transaction-like behavior
// If any event fails validation, it's reported in the summary but operation continues for other events
func (c *CLIEventService) BulkUpdateMetadata(ctx context.Context, eventIDs []string, update *BulkMetadataUpdate) (*BulkOperationsSummary, error) {
	// Validate input
	if len(eventIDs) == 0 {
		return nil, ErrEmptyEventList
	}

	if update == nil {
		return nil, NewMetadataError("update cannot be nil")
	}

	summary := &BulkOperationsSummary{
		EventsAffected: len(eventIDs),
		FieldsUpdated:  []string{},
		Errors:         []string{},
		AppliedCount:   0,
		SkippedCount:   0,
	}

	// Track which fields are being updated
	fieldsMap := make(map[string]bool)

	// Load all events first
	events := make(map[string]*career.CareerEvent)
	for _, id := range eventIDs {
		event, err := c.service.GetEventByID(ctx, id)
		if err != nil {
			summary.Errors = append(summary.Errors, "event "+id+" not found")
			summary.SkippedCount++
			continue
		}
		events[id] = event
	}

	// Validate all events before updating any (transaction-like behavior)
	for _, event := range events {
		if err := event.Validate(); err != nil {
			summary.Errors = append(summary.Errors, err.Error())
			summary.SkippedCount++
			continue
		}
	}

	// Apply updates to all valid events
	for _, id := range eventIDs {
		event, exists := events[id]
		if !exists {
			continue // Already reported in errors
		}

		updated := false

		// Update company field
		if update.Company != "" {
			if update.ApplyIfEmptyCompany {
				if event.Company == "" {
					event.Company = update.Company
					updated = true
					fieldsMap["company"] = true
				}
			} else {
				event.Company = update.Company
				updated = true
				fieldsMap["company"] = true
			}
		}

		// Update project field
		if update.Project != "" {
			if update.ApplyIfEmptyProject {
				if event.Project == "" {
					event.Project = update.Project
					updated = true
					fieldsMap["project"] = true
				}
			} else {
				event.Project = update.Project
				updated = true
				fieldsMap["project"] = true
			}
		}

		// Update tags field
		if len(update.Tags) > 0 {
			// Validate tags before applying
			validTags := true
			for _, tag := range update.Tags {
				if !career.AllowedTags[tag] {
					summary.Errors = append(summary.Errors, "invalid tag: "+tag)
					validTags = false
					break
				}
			}
			if validTags {
				event.Tags = update.Tags
				updated = true
				fieldsMap["tags"] = true
			}
		}

		// Update categories field
		if len(update.Categories) > 0 {
			event.Categories = update.Categories
			updated = true
			fieldsMap["categories"] = true
		}

		// Persist the update if any changes were made
		if updated {
			event.UpdatedAt = time.Now()
			if err := c.service.UpdateEvent(ctx, event); err != nil {
				summary.Errors = append(summary.Errors, "failed to update "+id+": "+err.Error())
				summary.SkippedCount++
				continue
			}
			summary.AppliedCount++
		}
	}

	// Build fields updated list
	for field := range fieldsMap {
		summary.FieldsUpdated = append(summary.FieldsUpdated, field)
	}

	return summary, nil
}

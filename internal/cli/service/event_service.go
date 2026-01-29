package service

import (
	"context"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
)

// CLIEventService wraps the existing career service to provide CLI-specific event operations.
type CLIEventService struct {
	service *careerservice.Service
}

// NewCLIEventService creates a new CLI event service.
func NewCLIEventService(service *careerservice.Service) *CLIEventService {
	return &CLIEventService{
		service: service,
	}
}

// CaptureEvent captures a new career event using the existing service.
func (c *CLIEventService) CaptureEvent(
	ctx context.Context, text string, date time.Time,
	mode careerservice.EventCaptureMode, opts ...Option,
) error {
	config := defaultConfig()
	for _, opt := range opts {
		opt(config)
	}

	event := &career.Event{
		Text:       text,
		Date:       date,
		Company:    config.Company,
		Project:    config.Project,
		Tags:       config.Tags,
		Categories: config.Categories,
	}

	return c.service.CaptureEvent(ctx, event, mode)
}

// ListEvents retrieves a list of events with optional filtering.
func (c *CLIEventService) ListEvents(ctx context.Context, filters *careerrepo.EventListFilters) ([]*career.Event, error) {
	if filters == nil {
		filters = &careerrepo.EventListFilters{}
	}
	return c.service.ListEvents(ctx, *filters)
}

// GetEventByID retrieves a specific event by its ID.
func (c *CLIEventService) GetEventByID(ctx context.Context, eventID string) (*career.Event, error) {
	return c.service.GetEventByID(ctx, eventID)
}

// UpdateEvent updates an existing career event.
func (c *CLIEventService) UpdateEvent(ctx context.Context, eventID string, text string, date time.Time, opts ...Option) error {
	config := defaultConfig()
	for _, opt := range opts {
		opt(config)
	}

	event := &career.Event{
		ID:         eventID,
		Text:       text,
		Date:       date,
		Company:    config.Company,
		Project:    config.Project,
		Tags:       config.Tags,
		Categories: config.Categories,
	}

	return c.service.UpdateEvent(ctx, event)
}

// DeleteEvent deletes a career event by ID.
func (c *CLIEventService) DeleteEvent(ctx context.Context, eventID string) error {
	return c.service.DeleteEvent(ctx, eventID)
}

// Option is a functional option that configures optional metadata for event
// capture and update operations.
//
// Callers pass zero or more Option values to CaptureEvent or UpdateEvent to
// attach secondary fields without requiring all parameters up-front. Each
// Option receives a mutable *EventConfig and sets one field on it.
//
// Available options:
//
//   - WithCompany sets the EventConfig.Company field to the given string.
//   - WithProject sets the EventConfig.Project field to the given string.
//   - WithTags replaces the EventConfig.Tags slice with the given []string.
//   - WithCategories replaces the EventConfig.Categories slice with the
//     given []string.
//
// When no options are supplied the service falls back to the defaults
// returned by defaultConfig, which initialises Tags and Categories to
// empty slices and leaves Company and Project as zero-value strings.
type Option func(*EventConfig)

// EventConfig holds optional configuration for event capture.
type EventConfig struct {
	Company    string
	Project    string
	Tags       []string
	Categories []string
}

// defaultConfig provides default event configuration.
func defaultConfig() *EventConfig {
	return &EventConfig{
		Tags:       []string{},
		Categories: []string{},
	}
}

// WithCompany sets the company for the event.
func WithCompany(company string) Option {
	return func(ec *EventConfig) {
		ec.Company = company
	}
}

// WithProject sets the project for the event.
func WithProject(project string) Option {
	return func(ec *EventConfig) {
		ec.Project = project
	}
}

// WithTags sets tags for the event.
func WithTags(tags []string) Option {
	return func(ec *EventConfig) {
		ec.Tags = tags
	}
}

// WithCategories sets categories for the event.
func WithCategories(categories []string) Option {
	return func(ec *EventConfig) {
		ec.Categories = categories
	}
}

// UpdateEventMetadata updates only the metadata fields of an event (company, project, tags, categories).
// This is used by the metadata editor to update event metadata without changing the text or date.
func (c *CLIEventService) UpdateEventMetadata(ctx context.Context, event *career.Event) error {
	if event == nil {
		return ErrNilEvent
	}

	if event.ID == "" {
		return ErrEmptyEventID
	}

	existingEvent, err := c.service.GetEventByID(ctx, event.ID)
	if err != nil {
		return err
	}

	event.Text = existingEvent.Text
	event.Date = existingEvent.Date
	event.CreatedAt = existingEvent.CreatedAt

	return c.service.UpdateEvent(ctx, event)
}

// Error definitions for metadata operations.
var (
	ErrNilEvent     = NewMetadataError("event cannot be nil")
	ErrEmptyEventID = NewMetadataError("event ID cannot be empty")
)

// MetadataError represents an error during metadata operations.
type MetadataError struct {
	message string
}

// NewMetadataError creates a new metadata error.
func NewMetadataError(message string) *MetadataError {
	return &MetadataError{message: message}
}

// Error implements the error interface.
func (me *MetadataError) Error() string {
	return me.message
}

// GetSkillsForEvent retrieves all skills associated with an event.
func (c *CLIEventService) GetSkillsForEvent(ctx context.Context, eventID string) ([]*career.Skill, error) {
	if c.service == nil {
		return []*career.Skill{}, nil
	}
	skillRepo := c.service.GetSkillRepository()
	if skillRepo == nil {
		return []*career.Skill{}, nil
	}
	return skillRepo.GetSkillsForEvent(ctx, eventID)
}

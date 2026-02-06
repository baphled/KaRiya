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
//
// Expected:
//   - service must be valid.
//
// Returns:
//   - A fully initialized CLIEventService ready for use.
//
// Side effects:
//   - None.
func NewCLIEventService(service *careerservice.Service) *CLIEventService {
	return &CLIEventService{
		service: service,
	}
}

// CaptureEvent captures a new career event using the existing service.
//
// Expected:
//   - ctx must be a valid context.Context.
//   - text must be a valid string.
//   - date must be a valid time.Time.
//   - mode must be a valid EventCaptureMode.
//   - opts must be valid Option values (can be empty).
//
// Returns:
//   - A error value if capture failed.
//
// Side effects:
//   - Creates a new event in the database.
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
//
// Expected:
//   - ctx must be a valid context.Context.
//   - filters must be a valid EventListFilters pointer (can be nil).
//
// Returns:
//   - A []*career.Event value containing matching events.
//   - An error value if retrieval failed.
//
// Side effects:
//   - None.
func (c *CLIEventService) ListEvents(ctx context.Context, filters *careerrepo.EventListFilters) ([]*career.Event, error) {
	if filters == nil {
		filters = &careerrepo.EventListFilters{}
	}
	return c.service.ListEvents(ctx, *filters)
}

// GetEventByID retrieves a specific event by its ID.
//
// Expected:
//   - ctx must be a valid context.Context.
//   - eventid must be a valid string.
//
// Returns:
//   - A *career.Event value if found.
//   - An error value if retrieval failed.
//
// Side effects:
//   - None.
func (c *CLIEventService) GetEventByID(ctx context.Context, eventID string) (*career.Event, error) {
	return c.service.GetEventByID(ctx, eventID)
}

// UpdateEvent updates an existing career event.
//
// Expected:
//   - Must be a valid string.
//   - Must be a valid string.
//   - time must be valid.
//   - option must be valid.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
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
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
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
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A Option value.
//
// Side effects:
//   - None.
func WithCompany(company string) Option {
	return func(ec *EventConfig) {
		ec.Company = company
	}
}

// WithProject sets the project for the event.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A Option value.
//
// Side effects:
//   - None.
func WithProject(project string) Option {
	return func(ec *EventConfig) {
		ec.Project = project
	}
}

// WithTags sets tags for the event.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A Option value.
//
// Side effects:
//   - None.
func WithTags(tags []string) Option {
	return func(ec *EventConfig) {
		ec.Tags = tags
	}
}

// WithCategories sets categories for the event.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A Option value.
//
// Side effects:
//   - None.
func WithCategories(categories []string) Option {
	return func(ec *EventConfig) {
		ec.Categories = categories
	}
}

// UpdateEventMetadata updates only the metadata fields of an event (company, project, tags, categories).
//
// Expected:
//   - event must be valid.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
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
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A fully initialized MetadataError ready for use.
//
// Side effects:
//   - None.
func NewMetadataError(message string) *MetadataError {
	return &MetadataError{message: message}
}

// Error implements the error interface.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (me *MetadataError) Error() string {
	return me.message
}

// GetSkillsForEvent retrieves all skills associated with an event.
//
// Expected:
//   - ctx must be a valid context.Context.
//   - eventid must be a valid string.
//
// Returns:
//   - A []*career.Skill value containing associated skills.
//   - An error value if retrieval failed.
//
// Side effects:
//   - None.
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

// LinkSkillToEvent creates an association between an event and a skill.
//
// Expected:
//   - ctx must be a valid context.Context.
//   - eventID must be a valid string identifier for an existing event.
//   - skillID must be a valid string identifier for an existing skill.
//
// Returns:
//   - An error value if linking failed.
//
// Side effects:
//   - Creates a link in the database between the event and skill.
func (c *CLIEventService) LinkSkillToEvent(ctx context.Context, eventID string, skillID string) error {
	eventRepo := c.service.GetEventRepository()
	return eventRepo.LinkSkill(ctx, eventID, skillID)
}

// UnlinkSkillFromEvent removes an association between an event and a skill.
//
// Expected:
//   - ctx must be a valid context.Context.
//   - eventID must be a valid string identifier for an existing event.
//   - skillID must be a valid string identifier for an existing skill.
//
// Returns:
//   - An error value if unlinking failed.
//
// Side effects:
//   - Removes the link in the database between the event and skill.
func (c *CLIEventService) UnlinkSkillFromEvent(ctx context.Context, eventID string, skillID string) error {
	eventRepo := c.service.GetEventRepository()
	return eventRepo.UnlinkSkill(ctx, eventID, skillID)
}

// ListAllSkills retrieves all skills without any filters.
//
// Expected:
//   - ctx must be a valid context.Context.
//
// Returns:
//   - A []*career.Skill value containing all skills.
//   - An error value if retrieval failed.
//
// Side effects:
//   - None.
func (c *CLIEventService) ListAllSkills(ctx context.Context) ([]*career.Skill, error) {
	skillRepo := c.service.GetSkillRepository()
	if skillRepo == nil {
		return []*career.Skill{}, nil
	}
	return skillRepo.List(ctx, nil)
}

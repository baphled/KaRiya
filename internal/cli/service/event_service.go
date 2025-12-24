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

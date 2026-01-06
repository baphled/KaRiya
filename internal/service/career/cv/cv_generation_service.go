package cv

import (
	"context"
	"fmt"
	"time"

	career "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	"github.com/google/uuid"
)

// CVGenerationService orchestrates the CV generation workflow
type CVGenerationService interface {
	// GenerateCV generates a CV from a saved configuration by name
	GenerateCV(ctx context.Context, configName string) (*career.CVView, error)

	// GenerateCVFromConfig generates a CV from a configuration object
	GenerateCVFromConfig(ctx context.Context, config *career.CVConfig) (*career.CVView, error)
}

// DefaultCVGenerationService is the default implementation of CVGenerationService
type DefaultCVGenerationService struct {
	eventRepo       careerrepo.Repository
	factRepo        careerrepo.FactRepository
	configManager   ConfigManager
	bulletGenerator BulletGenerator
	sectionBuilder  SectionBuilder
	logger          *logger.Logger
}

// NewCVGenerationService creates a new CVGenerationService instance
func NewCVGenerationService(
	eventRepo careerrepo.Repository,
	factRepo careerrepo.FactRepository,
	configManager ConfigManager,
	bulletGenerator BulletGenerator,
	sectionBuilder SectionBuilder,
	log *logger.Logger,
) *DefaultCVGenerationService {
	return &DefaultCVGenerationService{
		eventRepo:       eventRepo,
		factRepo:        factRepo,
		configManager:   configManager,
		bulletGenerator: bulletGenerator,
		sectionBuilder:  sectionBuilder,
		logger:          log,
	}
}

// GenerateCV generates a CV from a saved configuration by name
func (svc *DefaultCVGenerationService) GenerateCV(ctx context.Context, configName string) (*career.CVView, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// Load configuration
	config, err := svc.configManager.LoadConfig(ctx, configName)
	if err != nil {
		svc.logger.Error("Failed to load config %s: %v", configName, err)
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	return svc.GenerateCVFromConfig(ctx, config)
}

// GenerateCVFromConfig generates a CV from a configuration object
func (svc *DefaultCVGenerationService) GenerateCVFromConfig(ctx context.Context, config *career.CVConfig) (*career.CVView, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if config == nil {
		return nil, fmt.Errorf("configuration cannot be nil")
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		svc.logger.Error("Invalid configuration: %v", err)
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	svc.logger.Info("Starting CV generation for role %s with audience %s", config.TargetRole, config.TargetAudience)

	// Retrieve events with filters
	events, err := svc.retrieveEventsWithFilters(ctx, config.EventFilters)
	if err != nil {
		svc.logger.Error("Failed to retrieve events: %v", err)
		return nil, fmt.Errorf("failed to retrieve events: %w", err)
	}

	if len(events) == 0 {
		svc.logger.Info("No events found matching filters")
		return &career.CVView{
			ID:               uuid.New().String(),
			Name:             config.Name,
			TargetRole:       config.TargetRole,
			TargetAudience:   config.TargetAudience,
			EventFilters:     config.EventFilters,
			GeneratedAt:      time.Now(),
			SourceEventCount: 0,
			SourceFactCount:  0,
		}, nil
	}

	// Retrieve facts
	facts, err := svc.retrieveFacts(ctx)
	if err != nil {
		svc.logger.Error("Failed to retrieve facts: %v", err)
		// Continue anyway - facts are optional
		facts = []*career.Fact{}
	}

	// Generate bullets using BulletGenerator
	bullets, err := svc.bulletGenerator.GenerateBullets(ctx, events, facts, config.TargetRole, config.TargetAudience)
	if err != nil {
		svc.logger.Error("Failed to generate bullets: %v", err)
		return nil, fmt.Errorf("failed to generate bullets: %w", err)
	}

	svc.logger.Info("Generated %d bullets from %d events and %d facts", len(bullets), len(events), len(facts))

	// Build sections using SectionBuilder
	sections, err := svc.sectionBuilder.BuildSections(ctx, bullets, events, facts, config.TargetRole)
	if err != nil {
		svc.logger.Error("Failed to build sections: %v", err)
		return nil, fmt.Errorf("failed to build sections: %w", err)
	}

	// Create CVView with metadata and sections
	cvView := &career.CVView{
		ID:               uuid.New().String(),
		Name:             config.Name,
		TargetRole:       config.TargetRole,
		TargetAudience:   config.TargetAudience,
		EventFilters:     config.EventFilters,
		GeneratedAt:      time.Now(),
		SourceEventCount: len(events),
		SourceFactCount:  len(facts),
		Sections:         sections, // Include generated sections in the CV view
	}

	svc.logger.Info("CV generated successfully: %s (role: %s, sections: %d)", config.Name, config.TargetRole, len(sections))

	return cvView, nil
}

// retrieveEventsWithFilters retrieves events based on filter criteria
func (svc *DefaultCVGenerationService) retrieveEventsWithFilters(ctx context.Context, filters map[string]interface{}) ([]*career.CareerEvent, error) {
	// Return all events for CV generation
	// Use a high limit to ensure we get all events (repository defaults to 100)
	events, err := svc.eventRepo.List(ctx, careerrepo.ListFilters{
		Limit: 10000, // High enough to get all events
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}

	// Apply filters if they exist
	if len(filters) > 0 {
		events = svc.applyEventFilters(events, filters)
	}

	return events, nil
}

// applyEventFilters applies filter criteria to events
func (svc *DefaultCVGenerationService) applyEventFilters(events []*career.CareerEvent, filters map[string]interface{}) []*career.CareerEvent {
	var filtered []*career.CareerEvent

	for _, event := range events {
		if svc.eventMatchesFilters(event, filters) {
			filtered = append(filtered, event)
		}
	}

	return filtered
}

// eventMatchesFilters checks if an event matches all filter criteria
func (svc *DefaultCVGenerationService) eventMatchesFilters(event *career.CareerEvent, filters map[string]interface{}) bool {
	// Check date range filters
	if minDate, ok := filters["minDate"].(time.Time); ok {
		if event.Date.Before(minDate) {
			return false
		}
	}

	if maxDate, ok := filters["maxDate"].(time.Time); ok {
		if event.Date.After(maxDate) {
			return false
		}
	}

	// Check company filter
	if companies, ok := filters["companies"].([]string); ok {
		found := false
		for _, company := range companies {
			if event.Company == company {
				found = true
				break
			}
		}
		if !found && len(companies) > 0 {
			return false
		}
	}

	// Check tags filter
	if tags, ok := filters["tags"].([]string); ok {
		found := false
		for _, tag := range tags {
			for _, eventTag := range event.Tags {
				if tag == eventTag {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if !found && len(tags) > 0 {
			return false
		}
	}

	// Check categories filter
	// If event has categories, match against them; otherwise fall back to tags
	if categories, ok := filters["categories"].([]string); ok && len(categories) > 0 {
		found := false

		// If event has categories, match against them
		if len(event.Categories) > 0 {
			for _, category := range categories {
				for _, eventCategory := range event.Categories {
					if category == eventCategory {
						found = true
						break
					}
				}
				if found {
					break
				}
			}
		} else {
			// Fall back to tags if no categories are set
			for _, category := range categories {
				for _, eventTag := range event.Tags {
					if category == eventTag {
						found = true
						break
					}
				}
				if found {
					break
				}
			}
		}

		if !found {
			return false
		}
	}

	return true
}

// retrieveFacts retrieves all facts from the repository
func (svc *DefaultCVGenerationService) retrieveFacts(ctx context.Context) ([]*career.Fact, error) {
	if svc.factRepo == nil {
		return []*career.Fact{}, nil
	}

	facts, err := svc.factRepo.List(ctx, careerrepo.FactListFilters{})
	if err != nil {
		return nil, fmt.Errorf("failed to list facts: %w", err)
	}

	return facts, nil
}

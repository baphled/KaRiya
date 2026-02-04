//go:generate mockgen -destination=../../../testutil/mocks/service/cv_generation_service_mock.go -package=mocksvc github.com/baphled/kariya/internal/service/career/cv CVGenerationService

package cv

import (
	"context"
	"errors"
	"fmt"
	"time"

	career "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	"github.com/google/uuid"
)

// CVGenerationService orchestrates the CV generation workflow.
type CVGenerationService interface {
	// GenerateCV generates a CV from a saved configuration by name
	GenerateCV(ctx context.Context, configName string) (*career.CVView, error)

	// GenerateCVFromConfig generates a CV from a configuration object
	GenerateCVFromConfig(ctx context.Context, config *career.CVConfig) (*career.CVView, error)
}

// DefaultCVGenerationService is the default implementation of CVGenerationService.
type DefaultCVGenerationService struct {
	eventRepo       careerrepo.EventRepository
	factRepo        careerrepo.FactRepository
	configManager   ConfigManager
	bulletGenerator BulletGenerator
	dataProcessor   DataProcessingService
	sectionBuilder  SectionBuilder
	logger          *logger.Logger
}

// NewCVGenerationService creates a new CVGenerationService instance.
func NewCVGenerationService(
	eventRepo careerrepo.EventRepository,
	factRepo careerrepo.FactRepository,
	configManager ConfigManager,
	bulletGenerator BulletGenerator,
	dataProcessor DataProcessingService,
	sectionBuilder SectionBuilder,
	log *logger.Logger,
) *DefaultCVGenerationService {
	if dataProcessor == nil {
		panic("dataProcessor cannot be nil")
	}
	return &DefaultCVGenerationService{
		eventRepo:       eventRepo,
		factRepo:        factRepo,
		configManager:   configManager,
		bulletGenerator: bulletGenerator,
		dataProcessor:   dataProcessor,
		sectionBuilder:  sectionBuilder,
		logger:          log,
	}
}

// GenerateCV generates a CV from a saved configuration by name.
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

// GenerateCVFromConfig generates a CV from a configuration object.
func (svc *DefaultCVGenerationService) GenerateCVFromConfig(ctx context.Context, config *career.CVConfig) (*career.CVView, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if config == nil {
		return nil, errors.New("configuration cannot be nil")
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

	// BUG-010: Apply length format constraints - filter events by date
	if config.LengthFormat != "" {
		lengthConfig := GetLengthFormatConfig(LengthFormat(config.LengthFormat))
		if lengthConfig.MaxYearsHistory != nil {
			originalCount := len(events)
			filteredEvents := make([]*career.Event, 0, len(events))
			for _, event := range events {
				if lengthConfig.ShouldIncludeEvent(event.Date) {
					filteredEvents = append(filteredEvents, event)
				}
			}
			events = filteredEvents
			svc.logger.Info("Applied length format %s: filtered events from %d to %d (max years: %d)",
				config.LengthFormat, originalCount, len(events), *lengthConfig.MaxYearsHistory)
		}
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

	// Extract achievements from events for metric detection
	var achievements []*Achievement
	for _, event := range events {
		relatedFacts := svc.filterFactsForEvent(facts, event.ID)
		eventAchievements, err := svc.dataProcessor.ExtractAchievements(ctx, event, relatedFacts)
		if err != nil {
			svc.logger.Warn("Failed to extract achievements for event %s: %v", event.ID, err)
			continue
		}
		achievements = append(achievements, eventAchievements...)
	}

	svc.logger.Info("Extracted %d achievements from %d events", len(achievements), len(events))

	// Generate bullets (BUG-008: role-based scoring) WITH achievements
	bullets, err := svc.bulletGenerator.GenerateBullets(ctx, events, facts, achievements, config.TargetRole, config.TargetAudience)
	if err != nil {
		svc.logger.Error("Failed to generate bullets: %v", err)
		return nil, fmt.Errorf("failed to generate bullets: %w", err)
	}

	svc.logger.Info("Generated %d bullets from %d events and %d facts", len(bullets), len(events), len(facts))

	// Apply technology-based filtering if not Language Agnostic (Phase 10 - Task 40)
	if config.TechnologyFocus != "" && config.TechnologyFocus != string(TechnologyFocusLanguageAgnostic) {
		techFocus := TechnologyFocus(config.TechnologyFocus)
		bullets = svc.bulletGenerator.FilterByTechnologies(bullets, events, techFocus, config.SelectedTechnologies)
		svc.logger.Info("Applied technology filtering (%s) with %d technologies, %d bullets after filtering",
			config.TechnologyFocus, len(config.SelectedTechnologies), len(bullets))
	}

	// BUG-010: Apply length format constraints - filter bullets by confidence
	if config.LengthFormat != "" {
		lengthConfig := GetLengthFormatConfig(LengthFormat(config.LengthFormat))
		originalCount := len(bullets)
		filteredBullets := make([]*Bullet, 0, len(bullets))
		for _, bullet := range bullets {
			if lengthConfig.MeetsConfidenceThreshold(bullet.Confidence) {
				filteredBullets = append(filteredBullets, bullet)
			}
		}
		bullets = filteredBullets
		svc.logger.Info("Applied length format %s confidence filter: %d bullets filtered to %d (min confidence: %.2f)",
			config.LengthFormat, originalCount, len(bullets), lengthConfig.MinConfidence)
	}

	// Convert to domain bullets for SectionBuilder
	cvBullets := ConvertBullets(bullets)

	// Build sections using SectionBuilder (Phase 11 - Task 40: pass skills format config)
	skillsConfig := &SkillsFormatConfig{
		Format:               config.SkillsFormat,
		Limit:                config.SkillsLimit,
		SelectedTechnologies: config.SelectedTechnologies,
	}
	sections, err := svc.sectionBuilder.BuildSections(ctx, cvBullets, events, facts, config.TargetRole, skillsConfig)
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
		Sections:         sections,
	}

	svc.logger.Info("CV generated successfully: %s (role: %s, sections: %d)", config.Name, config.TargetRole, len(sections))

	return cvView, nil
}

// retrieveEventsWithFilters retrieves events based on filter criteria.
func (svc *DefaultCVGenerationService) retrieveEventsWithFilters(
	ctx context.Context, filters map[string]interface{},
) ([]*career.Event, error) {
	// Return all events for CV generation
	// Use a high limit to ensure we get all events (repository defaults to 100)
	events, err := svc.eventRepo.List(ctx, careerrepo.EventListFilters{
		Limit: 10000,
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

// applyEventFilters applies filter criteria to events.
func (svc *DefaultCVGenerationService) applyEventFilters(
	events []*career.Event, filters map[string]interface{},
) []*career.Event {
	var filtered []*career.Event

	for _, event := range events {
		if svc.eventMatchesFilters(event, filters) {
			filtered = append(filtered, event)
		}
	}

	return filtered
}

// eventMatchesFilters checks if an event matches all filter criteria.
func (svc *DefaultCVGenerationService) eventMatchesFilters(event *career.Event, filters map[string]interface{}) bool {
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

// retrieveFacts retrieves all facts from the repository.
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

// filterFactsForEvent returns facts that originated from a specific event.
func (svc *DefaultCVGenerationService) filterFactsForEvent(facts []*career.Fact, eventID string) []*career.Fact {
	var result []*career.Fact
	for _, fact := range facts {
		if fact.SourceEventID == eventID {
			result = append(result, fact)
		}
	}
	return result
}

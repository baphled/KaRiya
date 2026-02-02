package cv

import (
	"context"
	"errors"
	"fmt"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
)

// TraceabilityService provides traceability functionality for CV bullets.
type TraceabilityService struct {
	eventRepo careerrepo.EventRepository
	factRepo  careerrepo.FactRepository
	logger    *logger.Logger
}

// NewTraceabilityService creates a new TraceabilityService.
func NewTraceabilityService(
	eventRepo careerrepo.EventRepository,
	factRepo careerrepo.FactRepository,
	log *logger.Logger,
) *TraceabilityService {
	return &TraceabilityService{
		eventRepo: eventRepo,
		factRepo:  factRepo,
		logger:    log,
	}
}

// GetBulletSources retrieves source events and facts for a bullet.
func (ts *TraceabilityService) GetBulletSources(
	ctx context.Context,
	bulletID string,
	bullet *career.CVBullet,
) ([]*career.Event, []*career.Fact, error) {
	if bullet == nil {
		return nil, nil, errors.New("bullet cannot be nil")
	}

	ts.logger.Info(
		"Getting sources for bullet %s with %d events and %d facts",
		bulletID, len(bullet.SourceEventIDs), len(bullet.SourceFactIDs),
	)

	sourceEvents := make([]*career.Event, 0, len(bullet.SourceEventIDs))
	for _, eventID := range bullet.SourceEventIDs {
		event, err := ts.eventRepo.GetByID(ctx, eventID)
		if err != nil {
			ts.logger.Error("Failed to retrieve event %s: %v", eventID, err)
			continue
		}
		if event != nil {
			sourceEvents = append(sourceEvents, event)
		}
	}

	sourceFacts := make([]*career.Fact, 0, len(bullet.SourceFactIDs))
	for _, factID := range bullet.SourceFactIDs {
		fact, err := ts.factRepo.GetByID(ctx, factID)
		if err != nil {
			ts.logger.Error("Failed to retrieve fact %s: %v", factID, err)
			continue
		}
		if fact != nil {
			sourceFacts = append(sourceFacts, fact)
		}
	}

	ts.logger.Info("Retrieved sources for bullet %s with %d events and %d facts", bulletID, len(sourceEvents), len(sourceFacts))
	return sourceEvents, sourceFacts, nil
}

// ValidationReport represents validation results.
type ValidationReport struct {
	OrphanedBullets      []*career.CVBullet
	MissingSourceBullets []*career.CVBullet
	TotalBullets         int
	ValidBullets         int
	InvalidBullets       int
	Issues               []string
}

// IsValid returns true if all bullets are valid.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (vr *ValidationReport) IsValid() bool {
	return vr.InvalidBullets == 0
}

// GetEventUsage finds all bullets using a specific event.
func (ts *TraceabilityService) GetEventUsage(
	_ context.Context,
	eventID string,
	allBullets []*career.CVBullet,
) []*career.CVBullet {
	ts.logger.Info("Finding bullets using event %s", eventID)

	usageBullets := make([]*career.CVBullet, 0)
	for _, bullet := range allBullets {
		for _, sourceEventID := range bullet.SourceEventIDs {
			if sourceEventID == eventID {
				usageBullets = append(usageBullets, bullet)
				break
			}
		}
	}

	ts.logger.Info("Found %d bullets using event %s", len(usageBullets), eventID)
	return usageBullets
}

// GetFactUsage finds all bullets using a specific fact.
func (ts *TraceabilityService) GetFactUsage(
	_ context.Context,
	factID string,
	allBullets []*career.CVBullet,
) []*career.CVBullet {
	ts.logger.Info("Finding bullets using fact %s", factID)

	usageBullets := make([]*career.CVBullet, 0)
	for _, bullet := range allBullets {
		for _, sourceFactID := range bullet.SourceFactIDs {
			if sourceFactID == factID {
				usageBullets = append(usageBullets, bullet)
				break
			}
		}
	}

	ts.logger.Info("Found %d bullets using fact %s", len(usageBullets), factID)
	return usageBullets
}

// ValidateTraceability validates all bullets have sources.
func (ts *TraceabilityService) ValidateTraceability(
	_ context.Context,
	bullets []*career.CVBullet,
) *ValidationReport {
	ts.logger.Info("Validating traceability for %d bullets", len(bullets))

	report := &ValidationReport{
		TotalBullets:   len(bullets),
		ValidBullets:   0,
		InvalidBullets: 0,
		Issues:         make([]string, 0),
	}

	for _, bullet := range bullets {
		hasValidSources := len(bullet.SourceEventIDs) > 0 || len(bullet.SourceFactIDs) > 0

		if !hasValidSources {
			report.InvalidBullets++
			report.Issues = append(report.Issues, fmt.Sprintf("Bullet %s has no sources", bullet.ID))
			continue
		}

		report.ValidBullets++
	}

	ts.logger.Info("Validation completed: %d valid, %d invalid", report.ValidBullets, report.InvalidBullets)
	return report
}

// GetEventBulletMapping returns a mapping of events to bullets.
func (ts *TraceabilityService) GetEventBulletMapping(
	_ context.Context,
	bullets []*career.CVBullet,
) map[string][]*career.CVBullet {
	ts.logger.Info("Creating event-to-bullet mapping for %d bullets", len(bullets))

	mapping := make(map[string][]*career.CVBullet)
	for _, bullet := range bullets {
		for _, eventID := range bullet.SourceEventIDs {
			mapping[eventID] = append(mapping[eventID], bullet)
		}
	}

	ts.logger.Info("Event mapping created for %d events", len(mapping))
	return mapping
}

// GetFactBulletMapping returns a mapping of facts to bullets.
func (ts *TraceabilityService) GetFactBulletMapping(
	_ context.Context,
	bullets []*career.CVBullet,
) map[string][]*career.CVBullet {
	ts.logger.Info("Creating fact-to-bullet mapping for %d bullets", len(bullets))

	mapping := make(map[string][]*career.CVBullet)
	for _, bullet := range bullets {
		for _, factID := range bullet.SourceFactIDs {
			mapping[factID] = append(mapping[factID], bullet)
		}
	}

	ts.logger.Info("Fact mapping created for %d facts", len(mapping))
	return mapping
}

// SummaryString returns a human-readable summary of validation results.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (vr *ValidationReport) SummaryString() string {
	if vr.IsValid() {
		return fmt.Sprintf("All %d bullets have valid sources", vr.ValidBullets)
	}
	return fmt.Sprintf(
		"Validation failed: %d valid, %d invalid out of %d total bullets",
		vr.ValidBullets, vr.InvalidBullets, vr.TotalBullets,
	)
}

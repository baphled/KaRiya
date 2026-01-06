package cv

import (
	"context"
	"math"
	"sort"
	"strings"
	"time"

	career "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
)

// BulletGenerator generates ranked CV bullets from events and facts
type BulletGenerator interface {
	// GenerateBullets generates a list of ranked CV bullets from events and facts
	// Applies inclusion criteria, ranking algorithm, and role/audience-specific filtering
	GenerateBullets(ctx context.Context, events []*career.CareerEvent, facts []*career.Fact, targetRole string, targetAudience string) ([]*career.CVBullet, error)
}

// DefaultBulletGenerator is the default implementation of BulletGenerator
type DefaultBulletGenerator struct {
	eventRepo careerrepo.Repository
	factRepo  careerrepo.FactRepository
	logger    *logger.Logger
}

// NewBulletGenerator creates a new BulletGenerator instance
func NewBulletGenerator(eventRepo careerrepo.Repository, factRepo careerrepo.FactRepository, log *logger.Logger) *DefaultBulletGenerator {
	return &DefaultBulletGenerator{
		eventRepo: eventRepo,
		factRepo:  factRepo,
		logger:    log,
	}
}

// GenerateBullets generates ranked CV bullets from events and facts
func (bg *DefaultBulletGenerator) GenerateBullets(ctx context.Context, events []*career.CareerEvent, facts []*career.Fact, targetRole string, targetAudience string) ([]*career.CVBullet, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if len(events) == 0 {
		bg.logger.Info("No events provided for bullet generation")
		return []*career.CVBullet{}, nil
	}

	// Generate initial bullets from events and facts
	bullets := bg.generateInitialBullets(events, facts, targetRole, targetAudience)

	// Apply inclusion criteria filtering
	filteredBullets := bg.filterByInclusionCriteria(bullets, targetRole)

	// Rank the bullets
	rankedBullets := bg.rankBullets(filteredBullets, events)

	// Apply role-specific compression
	compressedBullets := bg.compressByRole(rankedBullets, targetRole)

	bg.logger.Info("Generated %d bullets from %d events for role %s", len(compressedBullets), len(events), targetRole)
	return compressedBullets, nil
}

// generateInitialBullets creates initial bullets from events and facts
func (bg *DefaultBulletGenerator) generateInitialBullets(events []*career.CareerEvent, facts []*career.Fact, targetRole string, targetAudience string) []*career.CVBullet {
	var bullets []*career.CVBullet

	// Generate bullets from facts (higher quality source)
	for _, fact := range facts {
		if !bg.isFactRelevantToRole(fact, targetRole) {
			continue
		}

		if !bg.isFactRelevantToAudience(fact, targetAudience) {
			continue
		}

		// Populate SourceEventIDs from fact's source event
		sourceEventIDs := []string{}
		if fact.SourceEventID != "" {
			sourceEventIDs = []string{fact.SourceEventID}
		}

		bullet := &career.CVBullet{
			ID:              fact.ID,
			Text:            fact.Text,
			SourceFactIDs:   []string{fact.ID},
			SourceEventIDs:  sourceEventIDs,
			InclusionReason: "fact_extraction",
			Confidence:      0.8,
		}
		bullets = append(bullets, bullet)
	}

	// Generate bullets from events
	for _, event := range events {
		if !bg.isEventRelevantToRole(event, targetRole) {
			continue
		}

		if !bg.isEventRelevantToAudience(event, targetAudience) {
			continue
		}

		bullet := &career.CVBullet{
			ID:              event.ID,
			Text:            event.Text,
			SourceEventIDs:  []string{event.ID},
			SourceFactIDs:   []string{},
			InclusionReason: "event_direct",
			Confidence:      0.7, // Default confidence for direct events
		}
		bullets = append(bullets, bullet)
	}

	return bullets
}

// filterByInclusionCriteria filters bullets based on inclusion/exclusion criteria
func (bg *DefaultBulletGenerator) filterByInclusionCriteria(bullets []*career.CVBullet, targetRole string) []*career.CVBullet {
	var filtered []*career.CVBullet

	for _, bullet := range bullets {
		// Must have at least one source
		if len(bullet.SourceEventIDs) == 0 && len(bullet.SourceFactIDs) == 0 {
			continue
		}

		// Single claim bullets only
		if !career.IsSingleClaimBullet(bullet.Text) {
			bg.logger.Info("Filtered bullet: multiple claims in %s", bullet.Text[:min(50, len(bullet.Text))])
			continue
		}

		// No aspirational language
		if career.IsAspirationLanguage(bullet.Text) {
			bg.logger.Info("Filtered bullet: aspirational language in %s", bullet.Text[:min(50, len(bullet.Text))])
			continue
		}

		// No inferred metrics
		if career.HasInferredMetrics(bullet.Text) {
			bg.logger.Info("Filtered bullet: inferred metrics in %s", bullet.Text[:min(50, len(bullet.Text))])
			continue
		}

		// No role inflation
		if career.IsRoleInflation(bullet.Text, targetRole) {
			bg.logger.Info("Filtered bullet: role inflation in %s", bullet.Text[:min(50, len(bullet.Text))])
			continue
		}

		filtered = append(filtered, bullet)
	}

	return filtered
}

// rankBullets ranks bullets using priority-based scoring
// Sorts by rank (score) first, then by source event date descending (newest first) as tiebreaker
func (bg *DefaultBulletGenerator) rankBullets(bullets []*career.CVBullet, events []*career.CareerEvent) []*career.CVBullet {
	// Build event map for date lookup
	eventMap := make(map[string]*career.CareerEvent)
	for _, event := range events {
		eventMap[event.ID] = event
	}

	// Calculate scores
	for _, bullet := range bullets {
		bullet.Rank = bg.calculateBulletScore(bullet)
	}

	// Sort by rank descending, with date descending as tiebreaker
	sort.Slice(bullets, func(i, j int) bool {
		// First, compare by rank
		if bullets[i].Rank != bullets[j].Rank {
			return bullets[i].Rank > bullets[j].Rank
		}

		// Tiebreaker: sort by source event date (newest first)
		// Get dates for both bullets
		var dateI, dateJ time.Time
		if len(bullets[i].SourceEventIDs) > 0 {
			if event, exists := eventMap[bullets[i].SourceEventIDs[0]]; exists {
				dateI = event.Date
			}
		}
		if len(bullets[j].SourceEventIDs) > 0 {
			if event, exists := eventMap[bullets[j].SourceEventIDs[0]]; exists {
				dateJ = event.Date
			}
		}

		// Newer dates (later in time) come first
		return dateI.After(dateJ)
	})

	return bullets
}

// calculateBulletScore calculates a bullet's priority score
// Note: This is called BEFORE section building, so we don't have access to events here
// The scoring must be based solely on bullet properties
func (bg *DefaultBulletGenerator) calculateBulletScore(bullet *career.CVBullet) float64 {
	baseScore := 0.5

	// Apply inclusion reason bonus
	switch bullet.InclusionReason {
	case "fact_extraction":
		baseScore += 0.3 // Facts are higher quality
	case "event_direct":
		baseScore += 0.2
	}

	// Apply confidence multiplier
	baseScore *= (0.5 + bullet.Confidence*0.5) // Confidence range: 0.5 to 1.0

	// Apply source count bonus
	sourceCount := len(bullet.SourceEventIDs) + len(bullet.SourceFactIDs)
	if sourceCount > 1 {
		baseScore += float64(sourceCount-1) * 0.1 // Bonus for multiple sources
	}

	// IMPORTANT: Bullets with source events are more valuable because they can be
	// grouped by company/project in the experience section
	// Bullets without source events can only go in standalone sections
	if len(bullet.SourceEventIDs) > 0 {
		baseScore += 0.15 // Significant bonus for having source events
	}

	return math.Min(baseScore, 1.0) // Cap at 1.0
}

// compressByRole applies role-specific bullet caps
func (bg *DefaultBulletGenerator) compressByRole(bullets []*career.CVBullet, targetRole string) []*career.CVBullet {
	maxBullets := bg.getBulletCapForRole(targetRole)

	if len(bullets) <= maxBullets {
		return bullets
	}

	// Remove lower-ranked bullets
	compressed := bullets[:maxBullets]
	bg.logger.Info("Compressed bullets from %d to %d for role %s", len(bullets), len(compressed), targetRole)
	return compressed
}

// getBulletCapForRole returns the maximum number of bullets for a role
// These are total bullets across ALL companies/sections, not per company
// Set generously to ensure good facts aren't filtered out before section building
func (bg *DefaultBulletGenerator) getBulletCapForRole(targetRole string) int {
	switch strings.ToLower(targetRole) {
	case "principal":
		return 50 // Allow all high-quality principal facts through
	case "staff":
		return 40
	case "em":
		return 40
	case "senior_ic":
		return 40
	default:
		return 30 // Default cap
	}
}

// isEventRelevantToRole checks if an event is relevant to a role
func (bg *DefaultBulletGenerator) isEventRelevantToRole(event *career.CareerEvent, targetRole string) bool {
	// Filter by category if role-specific categories exist
	roleCategories := bg.getCategoriesForRole(targetRole)
	if len(roleCategories) == 0 {
		return true // All categories relevant if not specified
	}

	for _, category := range event.Categories {
		for _, roleCategory := range roleCategories {
			if category == roleCategory {
				return true
			}
		}
	}

	return false
}

// isEventRelevantToAudience checks if an event is relevant to an audience
func (bg *DefaultBulletGenerator) isEventRelevantToAudience(event *career.CareerEvent, audience string) bool {
	if audience == "" {
		return true // All audiences relevant if not specified
	}

	// For now, accept all events for all audiences
	// In future, could implement audience-specific filtering
	return true
}

// isFactRelevantToRole checks if a fact is relevant to a role
func (bg *DefaultBulletGenerator) isFactRelevantToRole(fact *career.Fact, targetRole string) bool {
	// Match facts to target role based on role_fit field
	// This ensures a principal CV only includes principal-level facts
	return string(fact.RoleFit) == targetRole
}

// isFactRelevantToAudience checks if a fact is relevant to an audience
func (bg *DefaultBulletGenerator) isFactRelevantToAudience(fact *career.Fact, audience string) bool {
	if audience == "" {
		return true
	}

	// For now, accept all facts for all audiences
	return true
}

// getCategoriesForRole returns the preferred categories for a role
func (bg *DefaultBulletGenerator) getCategoriesForRole(targetRole string) []string {
	switch strings.ToLower(targetRole) {
	case "principal":
		return []string{"leadership", "strategy", "technical"}
	case "staff":
		return []string{"technical", "leadership"}
	case "em":
		return []string{"leadership", "mentoring"}
	case "senior_ic":
		return []string{"technical", "architecture"}
	default:
		return []string{}
	}
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

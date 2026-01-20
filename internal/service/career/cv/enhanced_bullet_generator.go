package cv

import (
	"context"
	"math"
	"sort"
	"strings"

	career "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	"github.com/google/uuid"
)

// EnhancedBulletGenerator generates professionally enhanced, ranked CV bullets
type EnhancedBulletGenerator interface {
	// GenerateBullets generates enhanced bullets from events and facts
	GenerateBullets(ctx context.Context,
		events []*career.CareerEvent,
		facts []*career.Fact,
		achievements []*Achievement,
		targetRole string,
		targetAudience string) ([]*EnhancedBullet, error)

	// FilterByRole filters bullets based on role relevance
	FilterByRole(bullets []*EnhancedBullet, role string) []*EnhancedBullet

	// FilterByAudience filters bullets based on audience fit
	FilterByAudience(bullets []*EnhancedBullet, audience string) []*EnhancedBullet

	// RankByRelevance ranks bullets using multi-factor scoring
	RankByRelevance(bullets []*EnhancedBullet, role string, audience string) []*EnhancedBullet

	// EnhanceWording improves bullet text for professional CV use
	EnhanceWording(bullet *EnhancedBullet, role string) (*EnhancedBullet, error)
}

// EnhancedBullet represents a CV bullet with scoring metadata
type EnhancedBullet struct {
	ID                string
	Text              string
	EnhancedText      string
	SourceEventIDs    []string
	SourceFactIDs     []string
	AudienceRelevance []string // Audience types this bullet is relevant to
	Metrics           []*Metric
	Confidence        float64
	RoleScore         float64 // 0.0-1.0
	AudienceScore     float64 // 0.0-1.0
	MetricScore       float64 // 0.0-1.0
	ImpactScore       float64 // 0.0-1.0
	FinalScore        float64 // Weighted combination
	ImpactLevel       string  // "low", "medium", "high"
	KeywordMatches    []string
	InclusionReason   string
	Rank              float64
}

// ToCVBullet converts an EnhancedBullet to a CVBullet for the domain layer.
// Uses EnhancedText as the primary Text if available, otherwise falls back to original Text.
func (eb *EnhancedBullet) ToCVBullet() *career.CVBullet {
	text := eb.Text
	if eb.EnhancedText != "" {
		text = eb.EnhancedText
	}
	return &career.CVBullet{
		ID:              eb.ID,
		Text:            text,
		EnhancedText:    eb.EnhancedText,
		SourceEventIDs:  eb.SourceEventIDs,
		SourceFactIDs:   eb.SourceFactIDs,
		Rank:            eb.Rank,
		InclusionReason: eb.InclusionReason,
		Confidence:      eb.Confidence,
		RoleScore:       eb.RoleScore,
		AudienceScore:   eb.AudienceScore,
		MetricScore:     eb.MetricScore,
		ImpactScore:     eb.ImpactScore,
		ImpactLevel:     eb.ImpactLevel,
		KeywordMatches:  eb.KeywordMatches,
	}
}

// ConvertBullets converts a slice of EnhancedBullets to CVBullets.
func ConvertBullets(enhanced []*EnhancedBullet) []*career.CVBullet {
	if enhanced == nil {
		return nil
	}
	bullets := make([]*career.CVBullet, len(enhanced))
	for i, eb := range enhanced {
		bullets[i] = eb.ToCVBullet()
	}
	return bullets
}

// DefaultEnhancedBulletGenerator implements EnhancedBulletGenerator
type DefaultEnhancedBulletGenerator struct {
	logger *logger.Logger
}

// NewEnhancedBulletGenerator creates a new enhanced bullet generator
func NewEnhancedBulletGenerator(log *logger.Logger) EnhancedBulletGenerator {
	return &DefaultEnhancedBulletGenerator{
		logger: log,
	}
}

// RoleFilter defines role-specific filtering criteria
type RoleFilter struct {
	PrimaryCategories   []string
	SecondaryCategories []string
	MinConfidence       float64
	PreferredMetrics    []string
}

// AudienceFilter defines audience-specific filtering criteria
type AudienceFilter struct {
	Name             string
	FocusAreas       []string
	PreferredMetrics []string
	MinImpactLevel   string
}

// GenerateBullets generates enhanced bullets from events and facts
func (ebg *DefaultEnhancedBulletGenerator) GenerateBullets(ctx context.Context,
	events []*career.CareerEvent,
	facts []*career.Fact,
	achievements []*Achievement,
	targetRole string,
	targetAudience string) ([]*EnhancedBullet, error) {

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if len(events) == 0 && len(facts) == 0 && len(achievements) == 0 {
		return []*EnhancedBullet{}, nil
	}

	// Create initial bullets from achievements (highest quality)
	bullets := ebg.createBulletsFromAchievements(achievements)

	// Add bullets from facts (filter by audience during creation)
	bullets = append(bullets, ebg.createBulletsFromFacts(facts, targetAudience)...)

	// Add bullets from events
	bullets = append(bullets, ebg.createBulletsFromEvents(events)...)

	// Filter by role
	bullets = ebg.FilterByRole(bullets, targetRole)

	// Filter by audience
	bullets = ebg.FilterByAudience(bullets, targetAudience)

	// Deduplicate bullets with identical text (keeps highest confidence)
	bullets = ebg.deduplicateBullets(bullets)

	// Calculate scores
	bullets = ebg.calculateScores(bullets, targetRole, targetAudience)

	// Rank by relevance
	bullets = ebg.RankByRelevance(bullets, targetRole, targetAudience)

	// Enhance wording
	for i, bullet := range bullets {
		enhanced, err := ebg.EnhanceWording(bullet, targetRole)
		if err != nil {
			ebg.logger.Warn("Failed to enhance bullet: %v", err)
			continue
		}
		bullets[i] = enhanced
	}

	// Note: Per-company/project caps are applied by SectionBuilder based on audience
	ebg.logger.Info("Generated %d enhanced bullets for role %s with audience %s",
		len(bullets), targetRole, targetAudience)

	return bullets, nil
}

// FilterByRole filters bullets based on role relevance
func (ebg *DefaultEnhancedBulletGenerator) FilterByRole(bullets []*EnhancedBullet, role string) []*EnhancedBullet {
	filter := ebg.getRoleFilter(role)
	var filtered []*EnhancedBullet

	for _, bullet := range bullets {
		// Must meet minimum confidence
		if bullet.Confidence < filter.MinConfidence {
			continue
		}

		filtered = append(filtered, bullet)
	}

	return filtered
}

// FilterByAudience filters bullets based on audience fit
func (ebg *DefaultEnhancedBulletGenerator) FilterByAudience(bullets []*EnhancedBullet, audience string) []*EnhancedBullet {
	if audience == "" {
		return bullets
	}

	var filtered []*EnhancedBullet
	for _, bullet := range bullets {
		// Bullets without audience relevance (from events) pass through
		if len(bullet.AudienceRelevance) == 0 {
			filtered = append(filtered, bullet)
			continue
		}

		// Check if bullet is relevant to target audience
		for _, relevantAudience := range bullet.AudienceRelevance {
			if relevantAudience == audience {
				filtered = append(filtered, bullet)
				break
			}
		}
	}

	return filtered
}

// RankByRelevance ranks bullets using multi-factor scoring
func (ebg *DefaultEnhancedBulletGenerator) RankByRelevance(bullets []*EnhancedBullet, role string, audience string) []*EnhancedBullet {
	// Calculate final scores
	for _, bullet := range bullets {
		bullet.FinalScore = ebg.calculateFinalScore(bullet)
	}

	// Sort by final score descending
	sort.Slice(bullets, func(i, j int) bool {
		if bullets[i].FinalScore != bullets[j].FinalScore {
			return bullets[i].FinalScore > bullets[j].FinalScore
		}
		// Tie-breaker: higher confidence
		return bullets[i].Confidence > bullets[j].Confidence
	})

	// Assign rank
	for i, bullet := range bullets {
		bullet.Rank = float64(i + 1)
	}

	return bullets
}

// EnhanceWording improves bullet text for professional CV use
func (ebg *DefaultEnhancedBulletGenerator) EnhanceWording(bullet *EnhancedBullet, role string) (*EnhancedBullet, error) {
	enhanced := bullet.EnhancedText
	if enhanced == "" {
		enhanced = bullet.Text
	}

	// Apply action verb enhancement
	enhanced = ebg.enhanceActionVerb(enhanced, role)

	// Add metric context
	if len(bullet.Metrics) > 0 {
		enhanced = ebg.addMetricContext(enhanced, bullet.Metrics)
	}

	// Structure for impact
	enhanced = ebg.structureForImpact(enhanced)

	// Role-specific customization
	enhanced = ebg.customizeForRole(enhanced, role)

	// Capitalize first letter
	if len(enhanced) > 0 {
		enhanced = strings.ToUpper(enhanced[:1]) + enhanced[1:]
	}

	bullet.EnhancedText = enhanced
	return bullet, nil
}

// Helper functions

// createBulletsFromAchievements creates bullets from achievements (highest quality)
func (ebg *DefaultEnhancedBulletGenerator) createBulletsFromAchievements(achievements []*Achievement) []*EnhancedBullet {
	var bullets []*EnhancedBullet

	for _, achievement := range achievements {
		bullet := &EnhancedBullet{
			ID:              uuid.New().String(),
			Text:            achievement.Description,
			EnhancedText:    achievement.Description,
			SourceEventIDs:  []string{achievement.EventID},
			SourceFactIDs:   achievement.FactIDs,
			Metrics:         achievement.Metrics,
			Confidence:      achievement.Confidence,
			InclusionReason: "achievement_extraction",
			ImpactLevel:     ebg.determineImpactLevel(achievement),
		}
		bullets = append(bullets, bullet)
	}

	return bullets
}

// createBulletsFromFacts creates bullets from facts, filtering by audience if specified
func (ebg *DefaultEnhancedBulletGenerator) createBulletsFromFacts(facts []*career.Fact, targetAudience string) []*EnhancedBullet {
	var bullets []*EnhancedBullet

	for _, fact := range facts {
		// Filter by audience relevance
		if !ebg.isFactRelevantToAudience(fact, targetAudience) {
			continue
		}

		// Populate SourceEventIDs from fact's source event
		sourceEventIDs := []string{}
		if fact.SourceEventID != "" {
			sourceEventIDs = []string{fact.SourceEventID}
		}

		bullet := &EnhancedBullet{
			ID:                fact.ID,
			Text:              fact.Text,
			EnhancedText:      fact.Text,
			SourceFactIDs:     []string{fact.ID},
			SourceEventIDs:    sourceEventIDs,
			AudienceRelevance: fact.AudienceRelevance,
			Confidence:        0.85,
			InclusionReason:   "fact_extraction",
			ImpactLevel:       "medium",
		}
		bullets = append(bullets, bullet)
	}

	return bullets
}

// isFactRelevantToAudience checks if a fact is relevant to an audience
// Uses the Fact.AudienceRelevance field to filter facts based on target audience
func (ebg *DefaultEnhancedBulletGenerator) isFactRelevantToAudience(fact *career.Fact, audience string) bool {
	if audience == "" {
		return true // All audiences relevant if not specified
	}

	// Check if the requested audience is in the fact's relevance list
	for _, relevantAudience := range fact.AudienceRelevance {
		if relevantAudience == audience {
			return true
		}
	}

	return false
}

// createBulletsFromEvents creates bullets from events
func (ebg *DefaultEnhancedBulletGenerator) createBulletsFromEvents(events []*career.CareerEvent) []*EnhancedBullet {
	var bullets []*EnhancedBullet

	for _, event := range events {
		bullet := &EnhancedBullet{
			ID:              event.ID,
			Text:            event.Text,
			EnhancedText:    event.Text,
			SourceEventIDs:  []string{event.ID},
			Confidence:      0.80,
			InclusionReason: "event_direct",
			ImpactLevel:     "low",
		}
		bullets = append(bullets, bullet)
	}

	return bullets
}

// deduplicateBullets removes bullets with identical text, keeping the one with highest confidence
// When merging, it combines source IDs to preserve lineage information
func (ebg *DefaultEnhancedBulletGenerator) deduplicateBullets(bullets []*EnhancedBullet) []*EnhancedBullet {
	if len(bullets) <= 1 {
		return bullets
	}

	// Map to track unique bullets by normalized text
	seen := make(map[string]*EnhancedBullet)

	for _, bullet := range bullets {
		// Normalize text for comparison (lowercase, trim spaces)
		normalizedText := strings.ToLower(strings.TrimSpace(bullet.Text))

		if existing, exists := seen[normalizedText]; exists {
			// Merge: keep the one with higher confidence, combine source IDs
			if bullet.Confidence > existing.Confidence {
				// Keep new bullet but merge source IDs from existing
				bullet.SourceEventIDs = mergeUniqueStrings(bullet.SourceEventIDs, existing.SourceEventIDs)
				bullet.SourceFactIDs = mergeUniqueStrings(bullet.SourceFactIDs, existing.SourceFactIDs)
				seen[normalizedText] = bullet
			} else {
				// Keep existing but merge source IDs from new
				existing.SourceEventIDs = mergeUniqueStrings(existing.SourceEventIDs, bullet.SourceEventIDs)
				existing.SourceFactIDs = mergeUniqueStrings(existing.SourceFactIDs, bullet.SourceFactIDs)
			}
		} else {
			seen[normalizedText] = bullet
		}
	}

	// Convert back to slice
	result := make([]*EnhancedBullet, 0, len(seen))
	for _, bullet := range seen {
		result = append(result, bullet)
	}

	ebg.logger.Info("Deduplicated bullets: %d -> %d", len(bullets), len(result))
	return result
}

// mergeUniqueStrings merges two string slices, removing duplicates
func mergeUniqueStrings(a, b []string) []string {
	seen := make(map[string]bool)
	for _, s := range a {
		seen[s] = true
	}
	for _, s := range b {
		seen[s] = true
	}
	result := make([]string, 0, len(seen))
	for s := range seen {
		result = append(result, s)
	}
	return result
}

// calculateScores calculates all score components
func (ebg *DefaultEnhancedBulletGenerator) calculateScores(bullets []*EnhancedBullet, role string, audience string) []*EnhancedBullet {
	for _, bullet := range bullets {
		bullet.RoleScore = ebg.calculateRoleScore(bullet, role)
		bullet.AudienceScore = ebg.calculateAudienceScore(bullet, audience)
		bullet.MetricScore = ebg.calculateMetricScore(bullet)
		bullet.ImpactScore = ebg.calculateImpactScore(bullet)
	}

	return bullets
}

// calculateFinalScore calculates weighted final score
func (ebg *DefaultEnhancedBulletGenerator) calculateFinalScore(bullet *EnhancedBullet) float64 {
	score := (0.25 * bullet.RoleScore) +
		(0.20 * bullet.AudienceScore) +
		(0.20 * bullet.MetricScore) +
		(0.20 * bullet.ImpactScore) +
		(0.15 * bullet.Confidence)

	return math.Min(score, 1.0)
}

// calculateRoleScore calculates role relevance score
func (ebg *DefaultEnhancedBulletGenerator) calculateRoleScore(bullet *EnhancedBullet, role string) float64 {
	score := 0.5 // Base score

	// Bonus for achievement-based bullets
	if bullet.InclusionReason == "achievement_extraction" {
		score += 0.3
	} else if bullet.InclusionReason == "fact_extraction" {
		score += 0.2
	}

	// Bonus for high confidence
	if bullet.Confidence > 0.8 {
		score += 0.1
	}

	return math.Min(score, 1.0)
}

// calculateAudienceScore calculates audience fit score
func (ebg *DefaultEnhancedBulletGenerator) calculateAudienceScore(bullet *EnhancedBullet, audience string) float64 {
	if audience == "" {
		return 1.0
	}

	// Base score
	score := 0.5

	// Bonus for matching audience relevance
	for _, relevantAudience := range bullet.AudienceRelevance {
		if relevantAudience == audience {
			score += 0.3 // Strong match
			break
		}
	}

	// Additional bonus for relevant impact level
	if bullet.ImpactLevel == "high" {
		score += 0.15
	} else if bullet.ImpactLevel == "medium" {
		score += 0.05
	}

	return math.Min(score, 1.0)
}

// calculateMetricScore calculates metric presence and quality score
func (ebg *DefaultEnhancedBulletGenerator) calculateMetricScore(bullet *EnhancedBullet) float64 {
	if len(bullet.Metrics) == 0 {
		return 0.3 // Lower score for no metrics
	}

	score := 0.6 + float64(len(bullet.Metrics))*0.1
	return math.Min(score, 1.0)
}

// calculateImpactScore calculates impact score based on metrics and impact level
func (ebg *DefaultEnhancedBulletGenerator) calculateImpactScore(bullet *EnhancedBullet) float64 {
	score := 0.5

	switch bullet.ImpactLevel {
	case "high":
		score += 0.4
	case "medium":
		score += 0.2
	}

	// Bonus for multiple metrics
	if len(bullet.Metrics) > 1 {
		score += 0.1
	}

	return math.Min(score, 1.0)
}

// determineImpactLevel determines impact level from achievement
func (ebg *DefaultEnhancedBulletGenerator) determineImpactLevel(achievement *Achievement) string {
	// High impact: multiple metrics, high confidence
	if len(achievement.Metrics) > 1 && achievement.Confidence > 0.85 {
		return "high"
	}

	// Medium impact: at least one metric
	if len(achievement.Metrics) > 0 {
		return "medium"
	}

	// Low impact: no metrics
	return "low"
}

// enhanceActionVerb replaces weak verbs with strong action verbs
func (ebg *DefaultEnhancedBulletGenerator) enhanceActionVerb(text string, role string) string {
	weakVerbs := map[string]string{
		"worked on":    "led",
		"helped with":  "contributed to",
		"was involved": "drove",
		"participated": "executed",
		"did":          "accomplished",
		"made":         "delivered",
		"got":          "achieved",
	}

	result := text
	lowerResult := strings.ToLower(result)

	for weak, strong := range weakVerbs {
		if strings.Contains(lowerResult, weak) {
			// Find the actual position in the original text (case-insensitive)
			idx := strings.Index(lowerResult, weak)
			if idx != -1 {
				// Replace the original text at the found position
				result = result[:idx] + strong + result[idx+len(weak):]
			}
			break
		}
	}

	return result
}

// addMetricContext adds context around metrics
func (ebg *DefaultEnhancedBulletGenerator) addMetricContext(text string, metrics []*Metric) string {
	if len(metrics) == 0 {
		return text
	}

	// Add first metric to text if not already present
	metric := metrics[0]
	if !strings.Contains(text, metric.Value) {
		text = text + " (" + metric.Value + metric.Unit + ")"
	}

	return text
}

// structureForImpact applies "action + context + result" structure
func (ebg *DefaultEnhancedBulletGenerator) structureForImpact(text string) string {
	// Remove common weak starters
	starters := []string{"I ", "We ", "The team "}
	result := text
	for _, starter := range starters {
		if strings.HasPrefix(result, starter) {
			result = strings.TrimPrefix(result, starter)
			break
		}
	}

	return strings.TrimSpace(result)
}

// customizeForRole applies role-specific customization
func (ebg *DefaultEnhancedBulletGenerator) customizeForRole(text string, role string) string {
	// For now, return as-is
	// Future: apply role-specific wording preferences
	return text
}

// getRoleFilter returns the filter for a specific role
func (ebg *DefaultEnhancedBulletGenerator) getRoleFilter(role string) *RoleFilter {
	switch strings.ToLower(role) {
	case "principal":
		return &RoleFilter{
			PrimaryCategories:   []string{"leadership", "strategy", "architecture"},
			SecondaryCategories: []string{"technical", "mentoring"},
			MinConfidence:       0.8,
			PreferredMetrics:    []string{"percentage", "count", "currency"},
		}
	case "staff":
		return &RoleFilter{
			PrimaryCategories:   []string{"technical", "architecture"},
			SecondaryCategories: []string{"leadership", "mentoring"},
			MinConfidence:       0.75,
			PreferredMetrics:    []string{"percentage", "count"},
		}
	case "em":
		return &RoleFilter{
			PrimaryCategories:   []string{"leadership", "mentoring"},
			SecondaryCategories: []string{"strategy", "product"},
			MinConfidence:       0.75,
			PreferredMetrics:    []string{"count", "percentage"},
		}
	case "senior_ic":
		return &RoleFilter{
			PrimaryCategories:   []string{"technical", "architecture"},
			SecondaryCategories: []string{"leadership", "strategy"},
			MinConfidence:       0.75,
			PreferredMetrics:    []string{"percentage", "count"},
		}
	default:
		return &RoleFilter{
			MinConfidence: 0.7,
		}
	}
}

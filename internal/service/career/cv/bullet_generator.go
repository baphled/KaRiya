package cv

import (
	"context"
	"math"
	"sort"
	"strings"

	"github.com/baphled/kariya/internal/config"
	"github.com/baphled/kariya/internal/constants"
	career "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	"github.com/google/uuid"
)

// Role scoring constants for calculateRoleScore.
const (
	roleScoreBase                    = 0.50 // Base score for all bullets
	roleScorePrimaryCategoryBoost    = 0.30 // Boost for primary category match
	roleScoreSecondaryCategoryBoost  = 0.15 // Boost for secondary category match
	roleScoreAchievementBonus        = 0.10 // Bonus for achievement-based bullets
	roleScoreFactBonus               = 0.05 // Bonus for fact-based bullets
	roleScoreHighConfidenceBonus     = 0.05 // Bonus for high confidence bullets
	roleScoreHighConfidenceThreshold = 0.80 // Threshold for high confidence
	technologySkillMatchBonus        = 0.15 // Bonus for matching selected technologies
	impactLevelHighBonus             = 0.15 // Bonus for high impact level
	impactLevelMediumBonus           = 0.05 // Bonus for medium impact level
)

// BulletGenerator generates professionally ranked CV bullets.
type BulletGenerator interface {
	// GenerateBullets generates bullets from events and facts
	GenerateBullets(ctx context.Context,
		events []*career.CareerEvent,
		facts []*career.Fact,
		achievements []*Achievement,
		targetRole string,
		targetAudience string) ([]*Bullet, error)

	// FilterByRole filters bullets based on role relevance
	FilterByRole(bullets []*Bullet, role string) []*Bullet

	// FilterByAudience filters bullets based on audience fit
	FilterByAudience(bullets []*Bullet, audience string) []*Bullet

	// RankByRelevance ranks bullets using multi-factor scoring
	RankByRelevance(bullets []*Bullet, role string, audience string) []*Bullet

	// EnhanceWording improves bullet text for professional CV use
	EnhanceWording(bullet *Bullet, role string) (*Bullet, error)

	// FilterByTechnologies filters and boosts bullets based on selected technologies
	FilterByTechnologies(bullets []*Bullet, events []*career.CareerEvent,
		techFocus TechnologyFocus, technologies []string) []*Bullet
}

// Bullet represents a CV bullet with scoring metadata.
type Bullet struct {
	ID                string
	Text              string
	EnhancedText      string
	SourceEventIDs    []string
	SourceFactIDs     []string
	AudienceRelevance []string                     // Audience types this bullet is relevant to
	Category          constants.CompetencyCategory // Primary competency category (BUG-008)
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

// ToCVBullet converts a Bullet to a CVBullet for the domain layer.
// Uses EnhancedText if available, otherwise falls back to original Text.
func (eb *Bullet) ToCVBullet() *career.CVBullet {
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
		Category:        eb.Category, // BUG-008: propagate category
		RoleScore:       eb.RoleScore,
		AudienceScore:   eb.AudienceScore,
		MetricScore:     eb.MetricScore,
		ImpactScore:     eb.ImpactScore,
		ImpactLevel:     eb.ImpactLevel,
		KeywordMatches:  eb.KeywordMatches,
	}
}

// ConvertBullets converts bullets to domain CVBullets.
func ConvertBullets(bullets []*Bullet) []*career.CVBullet {
	if bullets == nil {
		return nil
	}
	result := make([]*career.CVBullet, len(bullets))
	for i, b := range bullets {
		result[i] = b.ToCVBullet()
	}
	return result
}

// DefaultBulletGenerator implements BulletGenerator
type DefaultBulletGenerator struct {
	logger        *logger.Logger
	scoringConfig *config.ScoringConfig
}

// NewBulletGenerator creates a new bullet generator
func NewBulletGenerator(log *logger.Logger, scoringConfig *config.ScoringConfig) BulletGenerator {
	return &DefaultBulletGenerator{
		logger:        log,
		scoringConfig: scoringConfig,
	}
}

// RoleFilter defines role-specific filtering criteria
type RoleFilter struct {
	PrimaryCategories   []constants.CompetencyCategory
	SecondaryCategories []constants.CompetencyCategory
	MinConfidence       float64
	PreferredMetrics    []string
}

// GenerateBullets generates bullets from events and facts
func (bg *DefaultBulletGenerator) GenerateBullets(ctx context.Context,
	events []*career.CareerEvent,
	facts []*career.Fact,
	achievements []*Achievement,
	targetRole string,
	targetAudience string,
) ([]*Bullet, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	if len(events) == 0 && len(facts) == 0 && len(achievements) == 0 {
		return []*Bullet{}, nil
	}

	// Create initial bullets from achievements (highest quality)
	bullets := bg.createBulletsFromAchievements(achievements)

	// Add bullets from facts (filter by audience during creation)
	bullets = append(bullets, bg.createBulletsFromFacts(facts, targetAudience)...)

	// Add bullets from events
	bullets = append(bullets, bg.createBulletsFromEvents(events)...)

	// Filter by role
	bullets = bg.FilterByRole(bullets, targetRole)

	// Filter by audience
	bullets = bg.FilterByAudience(bullets, targetAudience)

	// Build event map for company-aware deduplication (BUG-013).
	eventMap := make(map[string]*career.CareerEvent, len(events))
	for _, event := range events {
		eventMap[event.ID] = event
	}

	// Deduplicate bullets with identical text at the same company (keeps highest confidence).
	bullets = bg.deduplicateBullets(bullets, eventMap)

	// Calculate scores
	bullets = bg.calculateScores(bullets, targetRole, targetAudience)

	// Rank by relevance
	bullets = bg.RankByRelevance(bullets, targetRole, targetAudience)

	// Enhance wording
	for i, bullet := range bullets {
		enhanced, err := bg.EnhanceWording(bullet, targetRole)
		if err != nil {
			bg.logger.Warn("Failed to enhance bullet: %v", err)
			continue
		}
		bullets[i] = enhanced
	}

	// Note: Per-company/project caps are applied by SectionBuilder based on audience
	bg.logger.Info("Generated %d bullets for role %s with audience %s",
		len(bullets), targetRole, targetAudience)

	return bullets, nil
}

// FilterByRole filters bullets based on role relevance
func (bg *DefaultBulletGenerator) FilterByRole(bullets []*Bullet, role string) []*Bullet {
	filter := bg.getRoleFilter(role)
	var filtered []*Bullet

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
func (bg *DefaultBulletGenerator) FilterByAudience(bullets []*Bullet, audience string) []*Bullet {
	if audience == "" {
		return bullets
	}

	var filtered []*Bullet
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
func (bg *DefaultBulletGenerator) RankByRelevance(bullets []*Bullet, _ string, _ string) []*Bullet {
	// Calculate final scores
	for _, bullet := range bullets {
		bullet.FinalScore = bg.calculateFinalScore(bullet)
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
func (bg *DefaultBulletGenerator) EnhanceWording(bullet *Bullet, role string) (*Bullet, error) {
	enhanced := bullet.EnhancedText
	if enhanced == "" {
		enhanced = bullet.Text
	}

	// Apply action verb enhancement
	enhanced = bg.enhanceActionVerb(enhanced, role)

	// Add metric context
	if len(bullet.Metrics) > 0 {
		enhanced = bg.addMetricContext(enhanced, bullet.Metrics)
	}

	// Structure for impact
	enhanced = bg.structureForImpact(enhanced)

	// Capitalize first letter
	if enhanced != "" {
		enhanced = strings.ToUpper(enhanced[:1]) + enhanced[1:]
	}

	bullet.EnhancedText = enhanced
	return bullet, nil
}

// Helper functions

// createBulletsFromAchievements creates bullets from achievements (highest quality)
func (bg *DefaultBulletGenerator) createBulletsFromAchievements(achievements []*Achievement) []*Bullet {
	var bullets []*Bullet

	for _, achievement := range achievements {
		bullet := &Bullet{
			ID:              uuid.New().String(),
			Text:            achievement.Description,
			EnhancedText:    achievement.Description,
			SourceEventIDs:  []string{achievement.EventID},
			SourceFactIDs:   achievement.FactIDs,
			Metrics:         achievement.Metrics,
			Category:        achievement.Category, // BUG-008: propagate from achievement
			Confidence:      achievement.Confidence,
			InclusionReason: string(constants.InclusionReasonAchievementExtraction),
			ImpactLevel:     bg.determineImpactLevel(achievement),
		}
		bullets = append(bullets, bullet)
	}

	return bullets
}

// createBulletsFromFacts creates bullets from facts, filtering by audience if specified
func (bg *DefaultBulletGenerator) createBulletsFromFacts(facts []*career.Fact, targetAudience string) []*Bullet {
	var bullets []*Bullet

	for _, fact := range facts {
		// Filter by audience relevance
		if !bg.isFactRelevantToAudience(fact, targetAudience) {
			continue
		}

		// Populate SourceEventIDs from fact's source event
		sourceEventIDs := []string{}
		if fact.SourceEventID != "" {
			sourceEventIDs = []string{fact.SourceEventID}
		}

		bullet := &Bullet{
			ID:                fact.ID,
			Text:              fact.Text,
			EnhancedText:      fact.Text,
			SourceFactIDs:     []string{fact.ID},
			SourceEventIDs:    sourceEventIDs,
			AudienceRelevance: fact.AudienceRelevance,
			Category:          bg.extractPrimaryCategory(fact.CompetencyCategories), // BUG-008
			Confidence:        0.85,
			InclusionReason:   string(constants.InclusionReasonFactExtraction),
			ImpactLevel:       "medium",
		}
		bullets = append(bullets, bullet)
	}

	return bullets
}

// isFactRelevantToAudience checks if a fact is relevant to an audience
// Uses the Fact.AudienceRelevance field to filter facts based on target audience
func (bg *DefaultBulletGenerator) isFactRelevantToAudience(fact *career.Fact, audience string) bool {
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
func (bg *DefaultBulletGenerator) createBulletsFromEvents(events []*career.CareerEvent) []*Bullet {
	var bullets []*Bullet

	for _, event := range events {
		bullet := &Bullet{
			ID:              event.ID,
			Text:            event.Text,
			EnhancedText:    event.Text,
			SourceEventIDs:  []string{event.ID},
			Category:        bg.extractPrimaryCategory(event.Categories), // BUG-008
			Confidence:      0.80,
			InclusionReason: string(constants.InclusionReasonEventDirect),
			ImpactLevel:     "low",
		}
		bullets = append(bullets, bullet)
	}

	return bullets
}

// extractPrimaryCategory extracts the primary (first) category from a list of categories
// and converts it to the type-safe CompetencyCategory constant (BUG-008)
func (bg *DefaultBulletGenerator) extractPrimaryCategory(categories []string) constants.CompetencyCategory {
	if len(categories) == 0 {
		return ""
	}
	// Use first category as primary
	primary := strings.ToLower(categories[0])
	if constants.IsValidCompetencyCategory(primary) {
		return constants.CompetencyCategory(primary)
	}
	return ""
}

// dedupKey is a struct-based key for bullet deduplication, avoiding string
// delimiter collisions (e.g. company names containing "|").
type dedupKey struct {
	text    string
	company string
}

// deduplicateBullets removes bullets with identical text at the same company,
// keeping the one with highest confidence. Bullets at different companies with
// identical text are kept separate (BUG-013). When merging, it combines source
// IDs to preserve lineage information.
func (bg *DefaultBulletGenerator) deduplicateBullets(bullets []*Bullet, eventMap map[string]*career.CareerEvent) []*Bullet {
	if len(bullets) <= 1 {
		return bullets
	}

	// Map to track unique bullets by normalized text + company (BUG-013).
	seen := make(map[dedupKey]*Bullet)

	for _, bullet := range bullets {
		// Normalize text for comparison (lowercase, trim spaces).
		normalizedText := strings.ToLower(strings.TrimSpace(bullet.Text))

		// Resolve primary company from source events (BUG-013).
		company := resolvePrimaryCompany(bullet.SourceEventIDs, eventMap)

		// When company cannot be resolved (no company on source events),
		// fall back to the first SourceEventID to prevent unrelated bullets
		// from merging under an empty key (BUG-015 defence-in-depth).
		// If SourceEventIDs is also empty, use bullet ID as last resort.
		if company == "" {
			if len(bullet.SourceEventIDs) > 0 {
				company = bullet.SourceEventIDs[0]
			} else {
				company = bullet.ID
			}
		}

		key := dedupKey{text: normalizedText, company: company}

		if existing, exists := seen[key]; exists {
			// Merge: keep the one with higher confidence, combine source IDs.
			if bullet.Confidence > existing.Confidence {
				// Keep new bullet but merge source IDs from existing.
				bullet.SourceEventIDs = mergeUniqueStrings(bullet.SourceEventIDs, existing.SourceEventIDs)
				bullet.SourceFactIDs = mergeUniqueStrings(bullet.SourceFactIDs, existing.SourceFactIDs)
				seen[key] = bullet
			} else {
				// Keep existing but merge source IDs from new.
				existing.SourceEventIDs = mergeUniqueStrings(existing.SourceEventIDs, bullet.SourceEventIDs)
				existing.SourceFactIDs = mergeUniqueStrings(existing.SourceFactIDs, bullet.SourceFactIDs)
			}
		} else {
			seen[key] = bullet
		}
	}

	// Convert back to slice.
	result := make([]*Bullet, 0, len(seen))
	for _, bullet := range seen {
		result = append(result, bullet)
	}

	bg.logger.Info("Deduplicated bullets: %d -> %d", len(bullets), len(result))
	return result
}

// resolvePrimaryCompany determines the primary company for a bullet by counting
// which company appears most frequently across its source events (BUG-013).
// When multiple companies are tied, the lexicographically smallest name wins
// to ensure deterministic results across runs.
func resolvePrimaryCompany(sourceEventIDs []string, eventMap map[string]*career.CareerEvent) string {
	if len(sourceEventIDs) == 0 || len(eventMap) == 0 {
		return ""
	}

	companyCounts := make(map[string]int)
	for _, eid := range sourceEventIDs {
		if event, ok := eventMap[eid]; ok && event.Company != "" {
			companyCounts[event.Company]++
		}
	}

	// Sort company names for deterministic tie-breaking.
	companies := make([]string, 0, len(companyCounts))
	for company := range companyCounts {
		companies = append(companies, company)
	}
	sort.Strings(companies)

	primaryCompany := ""
	maxCount := 0
	for _, company := range companies {
		if companyCounts[company] > maxCount {
			maxCount = companyCounts[company]
			primaryCompany = company
		}
	}

	return primaryCompany
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
func (bg *DefaultBulletGenerator) calculateScores(bullets []*Bullet, role string, audience string) []*Bullet {
	for _, bullet := range bullets {
		bullet.RoleScore = bg.calculateRoleScore(bullet, role)
		bullet.AudienceScore = bg.calculateAudienceScore(bullet, audience)
		bullet.MetricScore = bg.calculateMetricScore(bullet)
		bullet.ImpactScore = bg.calculateImpactScore(bullet)
	}

	return bullets
}

// calculateFinalScore calculates weighted final score using config weights
func (bg *DefaultBulletGenerator) calculateFinalScore(bullet *Bullet) float64 {
	// Use config weights if available, otherwise use defaults
	weights := bg.getWeights()

	score := (weights.RoleScore * bullet.RoleScore) +
		(weights.AudienceScore * bullet.AudienceScore) +
		(weights.MetricScore * bullet.MetricScore) +
		(weights.ImpactScore * bullet.ImpactScore) +
		(weights.Confidence * bullet.Confidence)

	return math.Min(score, 1.0)
}

// getWeights returns the scoring weights from config or defaults
func (bg *DefaultBulletGenerator) getWeights() config.ScoringWeights {
	if bg.scoringConfig != nil {
		return bg.scoringConfig.Weights
	}
	// Default weights (same as config.DefaultConfig)
	return config.ScoringWeights{
		RoleScore:     0.25,
		AudienceScore: 0.20,
		MetricScore:   0.20,
		ImpactScore:   0.20,
		Confidence:    0.15,
	}
}

// calculateRoleScore calculates role relevance score based on category alignment
func (bg *DefaultBulletGenerator) calculateRoleScore(bullet *Bullet, role string) float64 {
	score := roleScoreBase

	// Get role filter for category matching
	filter := bg.getRoleFilter(role)

	// Primary category match: strong boost
	for _, primary := range filter.PrimaryCategories {
		if bullet.Category == primary {
			score += roleScorePrimaryCategoryBoost
			break
		}
	}

	// Secondary category match: medium boost - only if no primary match
	if score == roleScoreBase {
		for _, secondary := range filter.SecondaryCategories {
			if bullet.Category == secondary {
				score += roleScoreSecondaryCategoryBoost
				break
			}
		}
	}

	// Bonus for achievement-based bullets
	if bullet.InclusionReason == string(constants.InclusionReasonAchievementExtraction) {
		score += roleScoreAchievementBonus
	} else if bullet.InclusionReason == string(constants.InclusionReasonFactExtraction) {
		score += roleScoreFactBonus
	}

	// Bonus for high confidence
	if bullet.Confidence > roleScoreHighConfidenceThreshold {
		score += roleScoreHighConfidenceBonus
	}

	return math.Min(score, 1.0)
}

// calculateAudienceScore calculates audience fit score
func (bg *DefaultBulletGenerator) calculateAudienceScore(bullet *Bullet, audience string) float64 {
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

	// Additional bonus for relevant impact level.
	if bullet.ImpactLevel == "high" {
		score += impactLevelHighBonus
	} else if bullet.ImpactLevel == "medium" {
		score += impactLevelMediumBonus
	}

	return math.Min(score, 1.0)
}

// calculateMetricScore calculates metric presence and quality score
func (bg *DefaultBulletGenerator) calculateMetricScore(bullet *Bullet) float64 {
	if len(bullet.Metrics) == 0 {
		return 0.3 // Lower score for no metrics
	}

	score := 0.6 + float64(len(bullet.Metrics))*0.1
	return math.Min(score, 1.0)
}

// calculateImpactScore calculates impact score based on metrics and impact level
func (bg *DefaultBulletGenerator) calculateImpactScore(bullet *Bullet) float64 {
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
func (bg *DefaultBulletGenerator) determineImpactLevel(achievement *Achievement) string {
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
func (bg *DefaultBulletGenerator) enhanceActionVerb(text string, _ string) string {
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
func (bg *DefaultBulletGenerator) addMetricContext(text string, metrics []*Metric) string {
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
func (bg *DefaultBulletGenerator) structureForImpact(text string) string {
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

// getRoleFilter returns the filter for a specific role
func (bg *DefaultBulletGenerator) getRoleFilter(role string) *RoleFilter {
	// Get MinConfidence from config if available
	minConfidence := bg.getMinConfidenceForRole(role)

	switch strings.ToLower(role) {
	case "principal":
		return &RoleFilter{
			PrimaryCategories:   []constants.CompetencyCategory{constants.CompetencyLeadership},
			SecondaryCategories: []constants.CompetencyCategory{constants.CompetencyTechnical, constants.CompetencyMentoring},
			MinConfidence:       minConfidence,
			PreferredMetrics:    []string{"percentage", "count", "currency"},
		}
	case "staff":
		return &RoleFilter{
			PrimaryCategories:   []constants.CompetencyCategory{constants.CompetencyTechnical},
			SecondaryCategories: []constants.CompetencyCategory{constants.CompetencyLeadership, constants.CompetencyMentoring},
			MinConfidence:       minConfidence,
			PreferredMetrics:    []string{"percentage", "count"},
		}
	case "em":
		return &RoleFilter{
			PrimaryCategories:   []constants.CompetencyCategory{constants.CompetencyLeadership, constants.CompetencyMentoring},
			SecondaryCategories: []constants.CompetencyCategory{constants.CompetencyProduct},
			MinConfidence:       minConfidence,
			PreferredMetrics:    []string{"count", "percentage"},
		}
	case "senior_ic":
		return &RoleFilter{
			PrimaryCategories:   []constants.CompetencyCategory{constants.CompetencyTechnical},
			SecondaryCategories: []constants.CompetencyCategory{constants.CompetencyLeadership},
			MinConfidence:       minConfidence,
			PreferredMetrics:    []string{"percentage", "count"},
		}
	default:
		return &RoleFilter{
			MinConfidence: minConfidence,
		}
	}
}

// getMinConfidenceForRole returns the minimum confidence for a role from config or defaults
func (bg *DefaultBulletGenerator) getMinConfidenceForRole(role string) float64 {
	if bg.scoringConfig != nil && bg.scoringConfig.RoleSettings != nil {
		if settings, ok := bg.scoringConfig.RoleSettings[strings.ToLower(role)]; ok {
			return settings.MinConfidence
		}
	}
	// Default values (same as config.DefaultConfig)
	switch strings.ToLower(role) {
	case "principal":
		return 0.80
	case "staff", "em", "senior_ic":
		return 0.75
	default:
		return 0.70
	}
}

// FilterByTechnologies filters and boosts bullets based on selected technologies.
// Applies technology-based scoring adjustments and sorts by final rank.
func (bg *DefaultBulletGenerator) FilterByTechnologies(
	bullets []*Bullet,
	events []*career.CareerEvent,
	techFocus TechnologyFocus,
	technologies []string,
) []*Bullet {
	// Language Agnostic: No filtering or boosting
	if techFocus == TechnologyFocusLanguageAgnostic {
		return bullets
	}

	// Build event map for skill lookup
	eventMap := make(map[string]*career.CareerEvent)
	for _, event := range events {
		eventMap[event.ID] = event
	}

	// Build technology set for fast lookup
	techSet := make(map[string]bool)
	for _, tech := range technologies {
		techSet[tech] = true
	}

	// Apply skill match bonus to bullets
	for _, bullet := range bullets {
		// Skip bullets without source events
		if len(bullet.SourceEventIDs) == 0 {
			continue
		}

		// Get source event
		event, exists := eventMap[bullet.SourceEventIDs[0]]
		if !exists {
			continue
		}

		// Check if event has any selected technology
		hasSelectedTech := false
		for _, skill := range event.Skills {
			if techSet[skill] {
				hasSelectedTech = true
				break
			}
		}

		// Apply skill match bonus for matching selected technologies
		if hasSelectedTech {
			bullet.Rank = math.Min(bullet.Rank+technologySkillMatchBonus, 1.0)
		}
		// No penalty for events without skills (keep baseline score)
	}

	// Sort bullets by rank descending (highest first)
	sort.Slice(bullets, func(i, j int) bool {
		return bullets[i].Rank > bullets[j].Rank
	})

	return bullets
}

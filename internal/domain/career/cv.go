package career

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/constants"
)

// CVView represents an in-memory, ephemeral CV generated from career events and facts.
// CVViews are NOT stored in the database - they are generated on-demand.
type CVView struct {
	ID               string                 `json:"id"`
	Name             string                 `json:"name"`
	TargetRole       string                 `json:"target_role"`
	TargetAudience   string                 `json:"target_audience"`
	EventFilters     map[string]interface{} `json:"event_filters"`
	GeneratedAt      time.Time              `json:"generated_at"`
	SourceEventCount int                    `json:"source_event_count"`
	SourceFactCount  int                    `json:"source_fact_count"`
	Sections         []*CVSection           `json:"sections,omitempty"`
}

// Validate checks if the CVView meets all defined criteria.
func (cv *CVView) Validate() error {
	// Validate ID
	if err := cv.validateID(); err != nil {
		return err
	}

	// Validate name
	if err := cv.validateName(); err != nil {
		return err
	}

	// Validate target role
	if err := cv.validateTargetRole(); err != nil {
		return err
	}

	// Validate target audience
	if err := cv.validateTargetAudience(); err != nil {
		return err
	}

	// Validate source counts
	if err := cv.validateSourceCounts(); err != nil {
		return err
	}

	return nil
}

// validateID ensures ID is not empty.
func (cv *CVView) validateID() error {
	if strings.TrimSpace(cv.ID) == "" {
		return errors.New("CV ID cannot be empty")
	}
	return nil
}

// validateName ensures name is not empty and within length constraints.
func (cv *CVView) validateName() error {
	trimmedName := strings.TrimSpace(cv.Name)
	if trimmedName == "" {
		return errors.New("CV name cannot be empty")
	}
	if len(trimmedName) > 200 {
		return errors.New("CV name cannot exceed 200 characters")
	}
	return nil
}

// validateTargetRole ensures target role is from the allowed set.
func (cv *CVView) validateTargetRole() error {
	trimmedRole := strings.TrimSpace(strings.ToLower(cv.TargetRole))
	if trimmedRole == "" {
		return ErrInvalidCVRole
	}
	if !constants.IsValidRoleFit(trimmedRole) {
		return fmt.Errorf("invalid target role: %s", cv.TargetRole)
	}
	return nil
}

// validateTargetAudience ensures a valid audience is specified.
func (cv *CVView) validateTargetAudience() error {
	trimmedAudience := strings.TrimSpace(strings.ToLower(cv.TargetAudience))
	if trimmedAudience == "" {
		return ErrInvalidAudience
	}
	if !constants.IsValidAudience(trimmedAudience) {
		return fmt.Errorf("invalid target audience: %s", cv.TargetAudience)
	}
	return nil
}

// validateSourceCounts ensures source counts are non-negative.
func (cv *CVView) validateSourceCounts() error {
	if cv.SourceEventCount < 0 {
		return errors.New("source event count cannot be negative")
	}
	if cv.SourceFactCount < 0 {
		return errors.New("source fact count cannot be negative")
	}
	return nil
}

// SectionContentGroup represents a grouped section (company/project with date and bullets).
type SectionContentGroup struct {
	Header    string      `json:"header"`
	StartDate string      `json:"start_date,omitempty"`
	EndDate   string      `json:"end_date,omitempty"`
	Bullets   []*CVBullet `json:"bullets"`
}

// CVSection represents a section within a CV (e.g., Experience, Skills, Summary).
type CVSection struct {
	ID          string                 `json:"id"`
	CVViewID    string                 `json:"cv_view_id"`
	SectionType string                 `json:"section_type"`
	Title       string                 `json:"title"`
	Order       int                    `json:"order"`
	Content     []*SectionContentGroup `json:"content,omitempty"`
	Summary     string                 `json:"summary,omitempty"`
}

// Validate checks if the CVSection meets all defined criteria.
func (cs *CVSection) Validate() error {
	// Validate ID
	if err := cs.validateID(); err != nil {
		return err
	}

	// Validate CVViewID
	if err := cs.validateCVViewID(); err != nil {
		return err
	}

	// Validate section type
	if err := cs.validateSectionType(); err != nil {
		return err
	}

	// Validate title
	if err := cs.validateTitle(); err != nil {
		return err
	}

	// Validate order
	if err := cs.validateOrder(); err != nil {
		return err
	}

	return nil
}

// validateID ensures ID is not empty.
func (cs *CVSection) validateID() error {
	if strings.TrimSpace(cs.ID) == "" {
		return errors.New("section ID cannot be empty")
	}
	return nil
}

// validateCVViewID ensures CVViewID is not empty.
func (cs *CVSection) validateCVViewID() error {
	if strings.TrimSpace(cs.CVViewID) == "" {
		return errors.New("CV view ID cannot be empty")
	}
	return nil
}

// validateSectionType ensures section type is from the allowed set.
func (cs *CVSection) validateSectionType() error {
	trimmedType := strings.TrimSpace(strings.ToLower(cs.SectionType))
	if trimmedType == "" {
		return errors.New("section type cannot be empty")
	}
	if !constants.IsValidSectionType(trimmedType) {
		return fmt.Errorf("invalid section type: %s", cs.SectionType)
	}
	return nil
}

// validateTitle ensures title is not empty and within length constraints.
func (cs *CVSection) validateTitle() error {
	trimmedTitle := strings.TrimSpace(cs.Title)
	if trimmedTitle == "" {
		return errors.New("section title cannot be empty")
	}
	if len(trimmedTitle) > 200 {
		return errors.New("section title cannot exceed 200 characters")
	}
	return nil
}

// validateOrder ensures order is non-negative.
func (cs *CVSection) validateOrder() error {
	if cs.Order < 0 {
		return errors.New("section order cannot be negative")
	}
	return nil
}

// CVBullet represents a single bullet point within a CV section.
type CVBullet struct {
	ID              string   `json:"id"`
	SectionID       string   `json:"section_id"`
	Text            string   `json:"text"`
	SourceEventIDs  []string `json:"source_event_ids"`
	SourceFactIDs   []string `json:"source_fact_ids"`
	Rank            float64  `json:"rank"`
	InclusionReason string   `json:"inclusion_reason"`
	Confidence      float64  `json:"confidence"`

	// Enhanced fields from BulletGenerator (Task 44)
	EnhancedText   string                       `json:"enhanced_text,omitempty"`
	Category       constants.CompetencyCategory `json:"category,omitempty"`
	RoleScore      float64                      `json:"role_score,omitempty"`
	AudienceScore  float64                      `json:"audience_score,omitempty"`
	MetricScore    float64                      `json:"metric_score,omitempty"`
	ImpactScore    float64                      `json:"impact_score,omitempty"`
	ImpactLevel    string                       `json:"impact_level,omitempty"`
	KeywordMatches []string                     `json:"keyword_matches,omitempty"`
}

// Validate checks if the CVBullet meets all defined criteria.
func (cb *CVBullet) Validate() error {
	// Validate ID
	if err := cb.validateID(); err != nil {
		return err
	}

	// Validate SectionID
	if err := cb.validateSectionID(); err != nil {
		return err
	}

	// Validate text
	if err := cb.validateText(); err != nil {
		return err
	}

	// Validate source events
	if err := cb.validateSourceEventIDs(); err != nil {
		return err
	}

	// Validate rank and confidence
	if err := cb.validateScores(); err != nil {
		return err
	}

	// Validate inclusion reason
	if err := cb.validateInclusionReason(); err != nil {
		return err
	}

	// Validate enhanced scores (optional fields, Task 44)
	if err := cb.validateEnhancedScores(); err != nil {
		return err
	}

	// Validate impact level (optional field, Task 44)
	if err := cb.validateImpactLevel(); err != nil {
		return err
	}

	return nil
}

// validateID ensures ID is not empty.
func (cb *CVBullet) validateID() error {
	if strings.TrimSpace(cb.ID) == "" {
		return errors.New("bullet ID cannot be empty")
	}
	return nil
}

// validateSectionID ensures SectionID is not empty.
func (cb *CVBullet) validateSectionID() error {
	if strings.TrimSpace(cb.SectionID) == "" {
		return errors.New("section ID cannot be empty")
	}
	return nil
}

// validateText ensures text is not empty and within length constraints.
func (cb *CVBullet) validateText() error {
	trimmedText := strings.TrimSpace(cb.Text)
	if trimmedText == "" {
		return errors.New("bullet text cannot be empty")
	}
	if len(trimmedText) > 500 {
		return errors.New("bullet text cannot exceed 500 characters")
	}
	return nil
}

// validateSourceEventIDs ensures at least one source event.
func (cb *CVBullet) validateSourceEventIDs() error {
	if len(cb.SourceEventIDs) == 0 {
		return ErrNoSourceEvents
	}

	// Check for empty event IDs
	seen := make(map[string]bool)
	for _, eventID := range cb.SourceEventIDs {
		if strings.TrimSpace(eventID) == "" {
			return errors.New("event ID in source list cannot be empty")
		}
		if seen[eventID] {
			return errors.New("duplicate event ID in source list")
		}
		seen[eventID] = true
	}
	return nil
}

// validateScores ensures rank and confidence are between 0.0 and 1.0.
func (cb *CVBullet) validateScores() error {
	if cb.Rank < 0.0 || cb.Rank > 1.0 {
		return fmt.Errorf("rank must be between 0.0 and 1.0, got %f", cb.Rank)
	}
	if cb.Confidence < 0.0 || cb.Confidence > 1.0 {
		return fmt.Errorf("confidence must be between 0.0 and 1.0, got %f", cb.Confidence)
	}
	return nil
}

// validateInclusionReason ensures inclusion reason is from the allowed set.
func (cb *CVBullet) validateInclusionReason() error {
	trimmedReason := strings.TrimSpace(strings.ToLower(cb.InclusionReason))
	if trimmedReason == "" {
		return errors.New("inclusion reason cannot be empty")
	}
	if !constants.IsValidInclusionReason(trimmedReason) {
		return fmt.Errorf("invalid inclusion reason: %s", cb.InclusionReason)
	}
	return nil
}

// validateEnhancedScores ensures enhanced score fields are in valid range (0.0-1.0).
// These fields are optional, so 0.0 is valid (not set).
func (cb *CVBullet) validateEnhancedScores() error {
	if cb.RoleScore < 0.0 || cb.RoleScore > 1.0 {
		return fmt.Errorf("role score must be between 0.0 and 1.0, got %f", cb.RoleScore)
	}
	if cb.AudienceScore < 0.0 || cb.AudienceScore > 1.0 {
		return fmt.Errorf("audience score must be between 0.0 and 1.0, got %f", cb.AudienceScore)
	}
	if cb.MetricScore < 0.0 || cb.MetricScore > 1.0 {
		return fmt.Errorf("metric score must be between 0.0 and 1.0, got %f", cb.MetricScore)
	}
	if cb.ImpactScore < 0.0 || cb.ImpactScore > 1.0 {
		return fmt.Errorf("impact score must be between 0.0 and 1.0, got %f", cb.ImpactScore)
	}
	return nil
}

// validateImpactLevel ensures impact level is valid.
func (cb *CVBullet) validateImpactLevel() error {
	if !constants.IsValidImpactLevel(cb.ImpactLevel) {
		return fmt.Errorf("invalid impact level: %s (must be low, medium, high, or empty)", cb.ImpactLevel)
	}
	return nil
}

// CVConfig represents the configuration for CV generation.
// Stored as YAML files in $HOME/.kariya/cv_configs/
// NOT stored in database - file-based configuration only.
type CVConfig struct {
	Name           string                 `yaml:"name" json:"name"`
	TargetRole     string                 `yaml:"target_role" json:"target_role"`
	TargetAudience string                 `yaml:"target_audience" json:"target_audience"`
	EventFilters   map[string]interface{} `yaml:"event_filters,omitempty" json:"event_filters,omitempty"`

	// Technology selections (Phase 10 - Task 40)
	// TechnologyFocus: language_agnostic, generalist, specialist
	TechnologyFocus string `yaml:"technology_focus,omitempty" json:"technology_focus,omitempty"`
	// SelectedTechnologies: Skill IDs for specialist focus
	SelectedTechnologies []string `yaml:"selected_technologies,omitempty" json:"selected_technologies,omitempty"`
	// FocusArea: backend, frontend, fullstack, devops
	FocusArea string `yaml:"focus_area,omitempty" json:"focus_area,omitempty"`
	// LengthFormat: ultra_short, short, standard, full
	LengthFormat string `yaml:"length_format,omitempty" json:"length_format,omitempty"`

	// Skills section format (Phase 11 enhancement)
	SkillsFormat string `yaml:"skills_format,omitempty" json:"skills_format,omitempty"`
	SkillsLimit  int    `yaml:"skills_limit,omitempty" json:"skills_limit,omitempty"`

	CreatedAt time.Time `yaml:"created_at" json:"created_at"`
	UpdatedAt time.Time `yaml:"updated_at" json:"updated_at"`
}

// Validate checks if the CVConfig meets all defined criteria.
func (cc *CVConfig) Validate() error {
	// Validate name
	if err := cc.validateName(); err != nil {
		return err
	}

	// Validate target role
	if err := cc.validateTargetRole(); err != nil {
		return err
	}

	// Validate target audience
	if err := cc.validateTargetAudience(); err != nil {
		return err
	}

	return nil
}

// validateName ensures name is not empty and contains valid filename characters.
func (cc *CVConfig) validateName() error {
	trimmedName := strings.TrimSpace(cc.Name)
	if trimmedName == "" {
		return errors.New("config name cannot be empty")
	}
	if len(trimmedName) > 200 {
		return errors.New("config name cannot exceed 200 characters")
	}

	// Check for valid filename characters
	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range invalidChars {
		if strings.Contains(trimmedName, char) {
			return fmt.Errorf("config name contains invalid character: %s", char)
		}
	}

	return nil
}

// validateTargetRole ensures target role is from the allowed set.
func (cc *CVConfig) validateTargetRole() error {
	trimmedRole := strings.TrimSpace(strings.ToLower(cc.TargetRole))
	if trimmedRole == "" {
		return ErrInvalidCVRole
	}
	if !constants.IsValidRoleFit(trimmedRole) {
		return fmt.Errorf("invalid target role: %s", cc.TargetRole)
	}
	return nil
}

// validateTargetAudience ensures a valid audience is specified.
func (cc *CVConfig) validateTargetAudience() error {
	trimmedAudience := strings.TrimSpace(strings.ToLower(cc.TargetAudience))
	if trimmedAudience == "" {
		return ErrInvalidAudience
	}
	if !constants.IsValidAudience(trimmedAudience) {
		return fmt.Errorf("invalid target audience: %s", cc.TargetAudience)
	}
	return nil
}

// ToJSON converts CVConfig to JSON bytes.
func (cc *CVConfig) ToJSON() ([]byte, error) {
	return json.Marshal(cc)
}

// FromJSON populates CVConfig from JSON bytes.
func (cc *CVConfig) FromJSON(data []byte) error {
	return json.Unmarshal(data, cc)
}

// Validation Helper Functions

// IsAspirationLanguage checks if text contains aspirational language
// Aspirational language includes: "will", "hoping", "aiming", "trying", "working towards", "in progress", etc.
func IsAspirationLanguage(text string) bool {
	lowerText := strings.ToLower(text)

	aspirationalKeywords := []string{
		"will ",
		"hoping to",
		"aiming to",
		"trying to",
		"working towards",
		"in progress",
		"planning to",
		"going to",
		"intend to",
		"plan to",
	}

	for _, keyword := range aspirationalKeywords {
		if strings.Contains(lowerText, keyword) {
			return true
		}
	}

	return false
}

// IsSingleClaimBullet checks if a bullet contains a single, focused claim
// Multiple claims are indicated by: "and", "while", "also", "in addition", etc.
func IsSingleClaimBullet(text string) bool {
	lowerText := strings.ToLower(text)

	// Check for multiple claims indicators
	multipleClaimIndicators := []string{
		" and ",
		" while ",
		" also ",
		" in addition",
		" additionally",
		" furthermore",
		" moreover",
		"; ",
		" | ",
	}

	count := 0
	for _, indicator := range multipleClaimIndicators {
		count += strings.Count(lowerText, indicator)
	}

	// Allow up to one conjunction (single "and" is often acceptable)
	return count <= 1
}

// HasInferredMetrics checks if text contains inferred or assumed metrics.
// Inferred metrics are vague quantifiers without specific numbers or context
// Note: "improved", "increased", "decreased" are valid action verbs when used with specific metrics.
func HasInferredMetrics(text string) bool {
	lowerText := strings.ToLower(text)

	// Only flag truly vague/inferred quantifiers
	// These patterns indicate non-specific claims without measurable data
	inferredPatterns := []string{
		"thousands of",
		"millions of",
		"billions of",
		"significantly",
		"substantially",
		"dramatically",
		"much faster",
		"much better",
		"greatly improved",
		"x%",
		"n users",
		"many users",
		"numerous",
	}

	for _, pattern := range inferredPatterns {
		if strings.Contains(lowerText, pattern) {
			return true
		}
	}

	return false
}

// IsRoleInflation checks if text exaggerates the scope of work for a given role.
// This is a placeholder function - real implementation would need role context.
func IsRoleInflation(text string, targetRole string) bool {
	lowerText := strings.ToLower(text)
	lowerRole := strings.ToLower(targetRole)

	// Define role-specific inflation checks
	inflationPatterns := map[string][]string{
		"senior_ic": {
			"led the entire company",
			"managed the team",
			"owned the strategy",
		},
		"principal": {
			"single-handedly built",
			"sole architect",
		},
		"staff": {
			"designed the entire platform",
			"owned the architecture",
		},
	}

	patterns, exists := inflationPatterns[lowerRole]
	if !exists {
		return false
	}

	for _, pattern := range patterns {
		if strings.Contains(lowerText, pattern) {
			return true
		}
	}

	return false
}

// Error Types for CV Validation

var (
	// ErrInvalidCVRole is returned when the target role is invalid.
	ErrInvalidCVRole = errors.New("invalid target role for CV")

	// ErrInvalidAudience is returned when no valid target audience is specified.
	ErrInvalidAudience = errors.New("a valid target audience must be specified")

	// ErrBulletMultipleClaims is returned when a bullet contains multiple claims.
	ErrBulletMultipleClaims = errors.New("bullet contains multiple claims - only single-claim bullets allowed")

	// ErrBulletAspirationLanguage is returned when a bullet contains aspirational language.
	ErrBulletAspirationLanguage = errors.New("bullet contains aspirational language - use factual language only")

	// ErrBulletInferredMetrics is returned when a bullet contains inferred metrics.
	ErrBulletInferredMetrics = errors.New("bullet contains inferred metrics without clear attribution")

	// ErrNoSourceEvents is returned when a bullet has no source events.
	ErrNoSourceEvents = errors.New("bullet must have at least one source event")
)

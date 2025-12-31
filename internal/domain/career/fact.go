package career

import (
	"errors"
	"strings"
	"time"
)

// RoleFit represents the career level fit for a fact
type RoleFit string

const (
	RoleFitPrincipal RoleFit = "principal"
	RoleFitEM        RoleFit = "em"
	RoleFitStaff     RoleFit = "staff"
	RoleFitSeniorIC  RoleFit = "senior_ic"
)

// AllowedRoleFits defines the set of valid role fit values
var AllowedRoleFits = map[string]bool{
	"principal": true,
	"em":        true,
	"staff":     true,
	"senior_ic": true,
}

// AllowedAudienceRelevance defines the set of valid audience types
var AllowedAudienceRelevance = map[string]bool{
	"hiring_manager": true,
	"recruiter":      true,
	"peer":           true,
}

// AspirationKeywords contains words that indicate aspirational language
var AspirationKeywords = map[string]bool{
	"will":    true,
	"should":  true,
	"could":   true,
	"might":   true,
	"may":     true,
	"want":    true,
	"wish":    true,
	"hope":    true,
	"plan":    true,
	"intend":  true,
	"attempt": true,
	"try":     true,
	"would":   true,
}

// Fact represents an extracted or inferred competency fact from an event or burst
type Fact struct {
	ID                   string    `json:"id"`
	Text                 string    `json:"text"`
	CompetencyCategories []string  `json:"competency_categories"`
	RoleFit              RoleFit   `json:"role_fit"`
	AudienceRelevance    []string  `json:"audience_relevance"`
	StrengthSignal       string    `json:"strength_signal"`
	SourceEventID        string    `json:"source_event_id,omitempty"`
	SourceBurstID        string    `json:"source_burst_id,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// Validate checks if the Fact meets all defined criteria
func (f *Fact) Validate() error {
	// Validate ID
	if err := f.validateID(); err != nil {
		return err
	}

	// Validate text
	if err := f.validateText(); err != nil {
		return err
	}

	// Validate competency categories
	if err := f.validateCompetencyCategories(); err != nil {
		return err
	}

	// Validate role fit
	if err := f.validateRoleFit(); err != nil {
		return err
	}

	// Validate audience relevance
	if err := f.validateAudienceRelevance(); err != nil {
		return err
	}

	// Validate source references
	if err := f.validateSourceReferences(); err != nil {
		return err
	}

	// Validate aspirational language
	if err := f.validateAspirationLanguage(); err != nil {
		return err
	}

	return nil
}

// validateID ensures ID is not empty
func (f *Fact) validateID() error {
	if strings.TrimSpace(f.ID) == "" {
		return errors.New("fact ID cannot be empty")
	}
	return nil
}

// validateText ensures text is not empty and within length constraints
func (f *Fact) validateText() error {
	trimmedText := strings.TrimSpace(f.Text)
	if trimmedText == "" {
		return errors.New("fact text cannot be empty")
	}
	if len(trimmedText) > 2000 {
		return errors.New("fact text cannot exceed 2000 characters")
	}
	return nil
}

// validateCompetencyCategories ensures categories are valid and not empty
func (f *Fact) validateCompetencyCategories() error {
	if len(f.CompetencyCategories) == 0 {
		return errors.New("fact must have at least one competency category")
	}

	// Check for duplicates and validate each category
	seen := make(map[string]bool)
	for _, category := range f.CompetencyCategories {
		if seen[category] {
			return errors.New("fact contains duplicate competency categories")
		}
		if !AllowedCategories[category] {
			return errors.New("invalid competency category: " + category)
		}
		seen[category] = true
	}

	return nil
}

// validateRoleFit ensures role fit is valid
func (f *Fact) validateRoleFit() error {
	if string(f.RoleFit) == "" {
		return errors.New("fact role fit cannot be empty")
	}
	if !AllowedRoleFits[string(f.RoleFit)] {
		return errors.New("invalid role fit: must be one of principal, em, staff, or senior_ic")
	}
	return nil
}

// validateAudienceRelevance ensures audience relevance values are valid
func (f *Fact) validateAudienceRelevance() error {
	if len(f.AudienceRelevance) == 0 {
		return errors.New("fact must have at least one audience relevance type")
	}

	// Check for duplicates and validate each audience
	seen := make(map[string]bool)
	for _, audience := range f.AudienceRelevance {
		if seen[audience] {
			return errors.New("fact contains duplicate audience relevance types")
		}
		if !AllowedAudienceRelevance[audience] {
			return errors.New("invalid audience relevance: " + audience)
		}
		seen[audience] = true
	}

	return nil
}

// validateSourceReferences ensures at least one source is provided
func (f *Fact) validateSourceReferences() error {
	if strings.TrimSpace(f.SourceEventID) == "" && strings.TrimSpace(f.SourceBurstID) == "" {
		return errors.New("fact must have at least one source (event or burst)")
	}
	return nil
}

// validateAspirationLanguage checks for aspirational language (will, should, could, etc.)
func (f *Fact) validateAspirationLanguage() error {
	lowerText := strings.ToLower(f.Text)
	words := strings.Fields(lowerText)

	for _, word := range words {
		// Remove punctuation for matching
		cleanWord := strings.Trim(word, ".,!?;:")
		if AspirationKeywords[cleanWord] {
			return errors.New("fact text contains aspirational language: " + word)
		}
	}

	return nil
}


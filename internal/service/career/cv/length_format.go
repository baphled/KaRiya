package cv

import "time"

// LengthFormatConfig defines length-specific filtering and configuration.
type LengthFormatConfig struct {
	ID               LengthFormat
	Name             string
	Description      string
	MaxYearsHistory  *int
	MaxCompanies     *int
	MaxBulletsPerJob *int
	MinConfidence    float64
	TargetPages      string
	IncludeSummary   bool
	IncludeEducation bool
	IncludeSkills    bool
}

// GetLengthFormatConfig returns the configuration for a length format.
//
// Expected:
//   - lengthformat must be valid.
//
// Returns:
//   - A fully initialized LengthFormatConfig ready for use.
//
// Side effects:
//   - None.
func GetLengthFormatConfig(format LengthFormat) *LengthFormatConfig {
	configs := getLengthFormatConfigMap()
	if config, exists := configs[format]; exists {
		return config
	}
	// Return default config for unknown format
	return &LengthFormatConfig{
		ID:               format,
		Name:             string(format),
		Description:      "Standard CV format",
		MaxYearsHistory:  nil,
		MaxCompanies:     nil,
		MaxBulletsPerJob: nil,
		MinConfidence:    0.65,
		TargetPages:      "2-3",
		IncludeSummary:   true,
		IncludeEducation: true,
		IncludeSkills:    true,
	}
}

// ListLengthFormatConfigs returns all length format configurations.
//
// Returns:
//   - A []*LengthFormatConfig value.
//
// Side effects:
//   - None.
func ListLengthFormatConfigs() []*LengthFormatConfig {
	return []*LengthFormatConfig{
		GetLengthFormatConfig(LengthFull),
		GetLengthFormatConfig(LengthStandard),
		GetLengthFormatConfig(LengthShort),
		GetLengthFormatConfig(LengthUltraShort),
	}
}

// ShouldIncludeEvent returns true if an event date falls within the MaxYearsHistory limit.
//
// Expected:
//   - time must be valid.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (c *LengthFormatConfig) ShouldIncludeEvent(eventDate time.Time) bool {
	if c.MaxYearsHistory == nil {
		return true
	}

	cutoffDate := time.Now().AddDate(-*c.MaxYearsHistory, 0, 0)
	return !eventDate.Before(cutoffDate)
}

// ShouldIncludeEventByYear returns true if an event year falls within the MaxYearsHistory limit.
//
// Expected:
//   - int must be valid.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (c *LengthFormatConfig) ShouldIncludeEventByYear(year int) bool {
	if c.MaxYearsHistory == nil {
		return true
	}

	currentYear := time.Now().Year()
	cutoffYear := currentYear - *c.MaxYearsHistory
	return year >= cutoffYear
}

// FilterCompaniesByLimit returns companies limited to MaxCompanies.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A []string value.
//
// Side effects:
//   - None.
func (c *LengthFormatConfig) FilterCompaniesByLimit(companies []string) []string {
	if c.MaxCompanies == nil || len(companies) <= *c.MaxCompanies {
		return companies
	}
	return companies[:*c.MaxCompanies]
}

// GetEffectiveBulletLimit returns the bullet limit per job, using default if not specified.
//
// Expected:
//   - int must be valid.
//
// Returns:
//   - A int value.
//
// Side effects:
//   - None.
func (c *LengthFormatConfig) GetEffectiveBulletLimit(defaultLimit int) int {
	if c.MaxBulletsPerJob == nil {
		return defaultLimit
	}
	return *c.MaxBulletsPerJob
}

// MeetsConfidenceThreshold returns true if a bullet's confidence meets the minimum threshold.
//
// Expected:
//   - float64 must be valid.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (c *LengthFormatConfig) MeetsConfidenceThreshold(confidence float64) bool {
	return confidence >= c.MinConfidence
}

// intPtr returns a pointer to an int value.
// Used to set optional *int config fields.
func intPtr(i int) *int {
	return &i
}

// getLengthFormatConfigMap returns the map of all length format configurations.
func getLengthFormatConfigMap() map[LengthFormat]*LengthFormatConfig {
	return map[LengthFormat]*LengthFormatConfig{
		LengthFull: {
			ID:               LengthFull,
			Name:             "Full",
			Description:      "Complete career history with all relevant details",
			MaxYearsHistory:  nil,
			MaxCompanies:     nil,
			MaxBulletsPerJob: nil,
			MinConfidence:    0.50,
			TargetPages:      "3+",
			IncludeSummary:   true,
			IncludeEducation: true,
			IncludeSkills:    true,
		},
		LengthStandard: {
			ID:               LengthStandard,
			Name:             "Standard",
			Description:      "Balanced format with recent 10 years emphasis",
			MaxYearsHistory:  intPtr(10),
			MaxCompanies:     nil,
			MaxBulletsPerJob: intPtr(5),
			MinConfidence:    0.65,
			TargetPages:      "2-3",
			IncludeSummary:   true,
			IncludeEducation: true,
			IncludeSkills:    true,
		},
		LengthShort: {
			ID:               LengthShort,
			Name:             "Short",
			Description:      "Concise format focusing on recent 5 years",
			MaxYearsHistory:  intPtr(5),
			MaxCompanies:     intPtr(5),
			MaxBulletsPerJob: intPtr(3),
			MinConfidence:    0.75,
			TargetPages:      "1-2",
			IncludeSummary:   true,
			IncludeEducation: false,
			IncludeSkills:    true,
		},
		LengthUltraShort: {
			ID:               LengthUltraShort,
			Name:             "1-Page",
			Description:      "Key highlights only for quick review",
			MaxYearsHistory:  intPtr(3),
			MaxCompanies:     intPtr(3),
			MaxBulletsPerJob: intPtr(2),
			MinConfidence:    0.85,
			TargetPages:      "1",
			IncludeSummary:   true,
			IncludeEducation: false,
			IncludeSkills:    false,
		},
	}
}

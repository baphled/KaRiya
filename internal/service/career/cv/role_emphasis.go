package cv

// RoleEmphasisConfig defines role-specific filtering and scoring configuration
type RoleEmphasisConfig struct {
	ID                  RoleEmphasis
	Name                string
	Description         string
	PrimaryCategories   []string // Categories with highest score boost
	SecondaryCategories []string // Categories with medium score boost
	BulletFocus         string   // Description of what bullets to emphasize
	MinConfidenceBoost  float64  // Additional confidence threshold for this role
}

// GetRoleEmphasisConfig returns the configuration for a role emphasis
func GetRoleEmphasisConfig(emphasis RoleEmphasis) *RoleEmphasisConfig {
	configs := getRoleEmphasisConfigMap()
	if config, exists := configs[emphasis]; exists {
		return config
	}
	// Return default config for unknown emphasis
	return &RoleEmphasisConfig{
		ID:                  emphasis,
		Name:                string(emphasis),
		Description:         "General role emphasis",
		PrimaryCategories:   []string{"technical", "delivery"},
		SecondaryCategories: []string{"leadership", "product"},
		BulletFocus:         "General professional achievements",
		MinConfidenceBoost:  0.0,
	}
}

// ListRoleEmphasisConfigs returns all role emphasis configurations
func ListRoleEmphasisConfigs() []*RoleEmphasisConfig {
	return []*RoleEmphasisConfig{
		GetRoleEmphasisConfig(RoleEmphasisSeniorBackend),
		GetRoleEmphasisConfig(RoleEmphasisStaffPrincipal),
		GetRoleEmphasisConfig(RoleEmphasisConsulting),
		GetRoleEmphasisConfig(RoleEmphasisLanguageAgnostic),
	}
}

// ScoreBulletCategory returns a score (0.0-1.0) for a category based on this role emphasis
func (c *RoleEmphasisConfig) ScoreBulletCategory(category string) float64 {
	// Check primary categories (highest score)
	for _, primary := range c.PrimaryCategories {
		if primary == category {
			return 1.0
		}
	}

	// Check secondary categories (medium score)
	for _, secondary := range c.SecondaryCategories {
		if secondary == category {
			return 0.6
		}
	}

	// Unrelated category (low score)
	return 0.3
}

// getRoleEmphasisConfigMap returns the map of all role emphasis configurations
func getRoleEmphasisConfigMap() map[RoleEmphasis]*RoleEmphasisConfig {
	return map[RoleEmphasis]*RoleEmphasisConfig{
		RoleEmphasisSeniorBackend: {
			ID:                  RoleEmphasisSeniorBackend,
			Name:                "Senior Backend",
			Description:         "Emphasizes backend engineering skills, technical depth, and system design",
			PrimaryCategories:   []string{"technical", "architecture"},
			SecondaryCategories: []string{"product", "delivery"},
			BulletFocus:         "Technical depth, system design, and product impact",
			MinConfidenceBoost:  0.0,
		},
		RoleEmphasisStaffPrincipal: {
			ID:                  RoleEmphasisStaffPrincipal,
			Name:                "Staff/Principal",
			Description:         "Emphasizes technical leadership, architecture, and cross-team impact",
			PrimaryCategories:   []string{"leadership", "strategy", "architecture"},
			SecondaryCategories: []string{"technical", "mentoring"},
			BulletFocus:         "Architecture decisions, mentorship, and cross-team influence",
			MinConfidenceBoost:  0.05,
		},
		RoleEmphasisConsulting: {
			ID:                  RoleEmphasisConsulting,
			Name:                "Consulting",
			Description:         "Emphasizes client engagements, adaptability, and business impact",
			PrimaryCategories:   []string{"strategy", "delivery", "consulting"},
			SecondaryCategories: []string{"technical", "leadership"},
			BulletFocus:         "Client work, rapid technology assessment, and business outcomes",
			MinConfidenceBoost:  0.0,
		},
		RoleEmphasisLanguageAgnostic: {
			ID:                  RoleEmphasisLanguageAgnostic,
			Name:                "Language-Agnostic",
			Description:         "Emphasizes transferable skills without technology-specific focus",
			PrimaryCategories:   []string{"technical", "architecture"},
			SecondaryCategories: []string{"leadership", "delivery", "product", "mentoring"},
			BulletFocus:         "Multi-language evidence, adaptability, and transferable skills",
			MinConfidenceBoost:  0.0,
		},
	}
}

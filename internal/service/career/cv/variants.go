package cv

import "fmt"

// RoleEmphasis defines what aspect of experience to emphasize in a CV
type RoleEmphasis string

const (
	// RoleEmphasisSeniorBackend emphasizes backend engineering skills and technical depth
	RoleEmphasisSeniorBackend RoleEmphasis = "senior_backend"

	// RoleEmphasisStaffPrincipal emphasizes technical leadership, architecture, and cross-team impact
	RoleEmphasisStaffPrincipal RoleEmphasis = "staff_principal"

	// RoleEmphasisConsulting emphasizes client engagements, adaptability, and business impact
	RoleEmphasisConsulting RoleEmphasis = "consulting"

	// RoleEmphasisLanguageAgnostic emphasizes transferable skills without technology-specific focus
	RoleEmphasisLanguageAgnostic RoleEmphasis = "language_agnostic"
)

// TechnologyFocus represents the technology presentation style for a CV.
// Used by the wizard modal for technology-specific filtering.
type TechnologyFocus string

const (
	// TechnologyFocusLanguageAgnostic - Technology-agnostic narrative style
	TechnologyFocusLanguageAgnostic TechnologyFocus = "language_agnostic"

	// TechnologyFocusGeneralist - Generalist style highlighting 2-5 technologies
	TechnologyFocusGeneralist TechnologyFocus = "generalist"

	// TechnologyFocusSpecialist - Specialist style focused on 1 technology
	TechnologyFocusSpecialist TechnologyFocus = "specialist"
)

// FocusArea represents a career focus area for the wizard modal.
type FocusArea string

const (
	// FocusAreaBackend - Backend development focus
	FocusAreaBackend FocusArea = "backend"

	// FocusAreaFrontend - Frontend development focus
	FocusAreaFrontend FocusArea = "frontend"

	// FocusAreaFullstack - Fullstack development focus
	FocusAreaFullstack FocusArea = "fullstack"

	// FocusAreaDevOps - DevOps/infrastructure focus
	FocusAreaDevOps FocusArea = "devops"
)

// LengthFormat defines CV density/length
type LengthFormat string

const (
	// LengthFull includes complete career history with all relevant details
	LengthFull LengthFormat = "full"

	// LengthStandard is a balanced 2-3 page format with recent 10 years emphasis
	LengthStandard LengthFormat = "standard"

	// LengthShort is a concise 1-2 page format focusing on recent 5 years
	LengthShort LengthFormat = "short"

	// LengthUltraShort is a 1-page highlights format with only top achievements
	LengthUltraShort LengthFormat = "ultra_short"
)

// CVVariant combines role emphasis and length format with configuration
type CVVariant struct {
	ID              string           // Unique identifier (e.g., "senior_backend_standard")
	Name            string           // Friendly display name
	Description     string           // Brief description for UI
	RoleEmphasis    RoleEmphasis     // Which role aspect to emphasize
	LengthFormat    LengthFormat     // Target length/density
	BaseStructure   CVStructure      // Underlying structure (standard, narrative, consulting, highlights)
	Sections        []SectionConfig  // Section configurations
	BulletConfig    BulletConfig     // Bullet filtering configuration
	ProfileOverride *ProfileOverride // Optional profile customization
	IsBuiltIn       bool             // True for system-defined variants
}

// SectionConfig defines how a section should be rendered
type SectionConfig struct {
	SectionType   string  // Section identifier (experience, skills, core_strengths, etc.)
	Title         string  // Display title (can be customized)
	Order         int     // Display order (lower = earlier)
	Enabled       bool    // Whether to include this section
	MinConfidence float64 // Minimum confidence threshold for bullets in this section
}

// BulletConfig defines bullet filtering and limits
type BulletConfig struct {
	MaxBulletsPerCompany *int     // Max bullets per company (nil = use default)
	MinConfidence        *float64 // Minimum confidence threshold (nil = use default)
	MaxYearsHistory      *int     // Max years of history to include (nil = no limit)
	MaxCompanies         *int     // Max companies to include (nil = no limit)
}

// ProfileOverride allows variant-specific profile customization
type ProfileOverride struct {
	ProfessionalTitle     *string  // Override title (e.g., "Senior Consulting Engineer")
	CoreStrengths         []string // Override core strengths
	CareerDifferentiators []string // Unique value propositions
	CareerPositioning     *string  // Career positioning narrative
}

// RoleEmphasisInfo provides metadata about a role emphasis for UI display
type RoleEmphasisInfo struct {
	ID          RoleEmphasis
	Name        string
	Description string
}

// LengthFormatInfo provides metadata about a length format for UI display
type LengthFormatInfo struct {
	ID          LengthFormat
	Name        string
	Description string
}

// BuiltInVariants contains all 16 predefined CV variants
var BuiltInVariants = []*CVVariant{
	// Senior Backend variants (4)
	{
		ID:            "senior_backend_full",
		Name:          "Senior Backend (Full)",
		Description:   "Complete career history emphasizing backend engineering",
		RoleEmphasis:  RoleEmphasisSeniorBackend,
		LengthFormat:  LengthFull,
		BaseStructure: CVStructureStandard,
		BulletConfig: BulletConfig{
			MinConfidence: floatPtr(0.60),
		},
		IsBuiltIn: true,
	},
	{
		ID:            "senior_backend_standard",
		Name:          "Senior Backend (Standard)",
		Description:   "Balanced 2-3 page CV for senior backend roles",
		RoleEmphasis:  RoleEmphasisSeniorBackend,
		LengthFormat:  LengthStandard,
		BaseStructure: CVStructureStandard,
		BulletConfig: BulletConfig{
			MinConfidence:   floatPtr(0.65),
			MaxYearsHistory: intPtr(10),
		},
		IsBuiltIn: true,
	},
	{
		ID:            "senior_backend_short",
		Name:          "Senior Backend (Short)",
		Description:   "Concise 1-2 page CV for senior backend roles",
		RoleEmphasis:  RoleEmphasisSeniorBackend,
		LengthFormat:  LengthShort,
		BaseStructure: CVStructureStandard,
		BulletConfig: BulletConfig{
			MinConfidence:   floatPtr(0.75),
			MaxYearsHistory: intPtr(5),
		},
		IsBuiltIn: true,
	},
	{
		ID:            "senior_backend_ultra_short",
		Name:          "Senior Backend (1-Page)",
		Description:   "One-page highlights for senior backend roles",
		RoleEmphasis:  RoleEmphasisSeniorBackend,
		LengthFormat:  LengthUltraShort,
		BaseStructure: CVStructureHighlights,
		BulletConfig: BulletConfig{
			MinConfidence: floatPtr(0.85),
			MaxCompanies:  intPtr(3),
		},
		IsBuiltIn: true,
	},

	// Staff/Principal variants (4)
	{
		ID:            "staff_principal_full",
		Name:          "Staff/Principal (Full)",
		Description:   "Complete career history emphasizing technical leadership",
		RoleEmphasis:  RoleEmphasisStaffPrincipal,
		LengthFormat:  LengthFull,
		BaseStructure: CVStructureStandard,
		BulletConfig: BulletConfig{
			MinConfidence: floatPtr(0.70),
		},
		IsBuiltIn: true,
	},
	{
		ID:            "staff_principal_standard",
		Name:          "Staff/Principal (Standard)",
		Description:   "Balanced 2-3 page CV for staff/principal roles",
		RoleEmphasis:  RoleEmphasisStaffPrincipal,
		LengthFormat:  LengthStandard,
		BaseStructure: CVStructureStandard,
		BulletConfig: BulletConfig{
			MinConfidence:   floatPtr(0.75),
			MaxYearsHistory: intPtr(10),
		},
		IsBuiltIn: true,
	},
	{
		ID:            "staff_principal_short",
		Name:          "Staff/Principal (Short)",
		Description:   "Concise 1-2 page CV for staff/principal roles",
		RoleEmphasis:  RoleEmphasisStaffPrincipal,
		LengthFormat:  LengthShort,
		BaseStructure: CVStructureStandard,
		BulletConfig: BulletConfig{
			MinConfidence:   floatPtr(0.80),
			MaxYearsHistory: intPtr(5),
		},
		IsBuiltIn: true,
	},
	{
		ID:            "staff_principal_ultra_short",
		Name:          "Staff/Principal (1-Page)",
		Description:   "One-page highlights for staff/principal roles",
		RoleEmphasis:  RoleEmphasisStaffPrincipal,
		LengthFormat:  LengthUltraShort,
		BaseStructure: CVStructureHighlights,
		BulletConfig: BulletConfig{
			MinConfidence: floatPtr(0.90),
			MaxCompanies:  intPtr(3),
		},
		IsBuiltIn: true,
	},

	// Consulting variants (4)
	{
		ID:            "consulting_full",
		Name:          "Consulting (Full)",
		Description:   "Complete career history emphasizing client engagements",
		RoleEmphasis:  RoleEmphasisConsulting,
		LengthFormat:  LengthFull,
		BaseStructure: CVStructureConsulting,
		BulletConfig: BulletConfig{
			MinConfidence: floatPtr(0.50),
		},
		IsBuiltIn: true,
	},
	{
		ID:            "consulting_standard",
		Name:          "Consulting (Standard)",
		Description:   "Balanced 2-3 page CV for consulting roles",
		RoleEmphasis:  RoleEmphasisConsulting,
		LengthFormat:  LengthStandard,
		BaseStructure: CVStructureConsulting,
		BulletConfig: BulletConfig{
			MinConfidence:   floatPtr(0.60),
			MaxYearsHistory: intPtr(10),
		},
		IsBuiltIn: true,
	},
	{
		ID:            "consulting_short",
		Name:          "Consulting (Short)",
		Description:   "Concise 1-2 page CV for consulting roles",
		RoleEmphasis:  RoleEmphasisConsulting,
		LengthFormat:  LengthShort,
		BaseStructure: CVStructureConsulting,
		BulletConfig: BulletConfig{
			MinConfidence:   floatPtr(0.70),
			MaxYearsHistory: intPtr(5),
		},
		IsBuiltIn: true,
	},
	{
		ID:            "consulting_ultra_short",
		Name:          "Consulting (1-Page)",
		Description:   "One-page highlights for consulting roles",
		RoleEmphasis:  RoleEmphasisConsulting,
		LengthFormat:  LengthUltraShort,
		BaseStructure: CVStructureHighlights,
		BulletConfig: BulletConfig{
			MinConfidence: floatPtr(0.80),
			MaxCompanies:  intPtr(3),
		},
		IsBuiltIn: true,
	},

	// Language-Agnostic variants (4)
	{
		ID:            "language_agnostic_full",
		Name:          "Language-Agnostic (Full)",
		Description:   "Complete career history emphasizing transferable skills",
		RoleEmphasis:  RoleEmphasisLanguageAgnostic,
		LengthFormat:  LengthFull,
		BaseStructure: CVStructureNarrative,
		BulletConfig: BulletConfig{
			MinConfidence: floatPtr(0.65),
		},
		IsBuiltIn: true,
	},
	{
		ID:            "language_agnostic_standard",
		Name:          "Language-Agnostic (Standard)",
		Description:   "Balanced 2-3 page CV emphasizing transferable skills",
		RoleEmphasis:  RoleEmphasisLanguageAgnostic,
		LengthFormat:  LengthStandard,
		BaseStructure: CVStructureNarrative,
		BulletConfig: BulletConfig{
			MinConfidence:   floatPtr(0.70),
			MaxYearsHistory: intPtr(10),
		},
		IsBuiltIn: true,
	},
	{
		ID:            "language_agnostic_short",
		Name:          "Language-Agnostic (Short)",
		Description:   "Concise 1-2 page CV emphasizing transferable skills",
		RoleEmphasis:  RoleEmphasisLanguageAgnostic,
		LengthFormat:  LengthShort,
		BaseStructure: CVStructureNarrative,
		BulletConfig: BulletConfig{
			MinConfidence:   floatPtr(0.78),
			MaxYearsHistory: intPtr(5),
		},
		IsBuiltIn: true,
	},
	{
		ID:            "language_agnostic_ultra_short",
		Name:          "Language-Agnostic (1-Page)",
		Description:   "One-page highlights emphasizing transferable skills",
		RoleEmphasis:  RoleEmphasisLanguageAgnostic,
		LengthFormat:  LengthUltraShort,
		BaseStructure: CVStructureHighlights,
		BulletConfig: BulletConfig{
			MinConfidence: floatPtr(0.88),
			MaxCompanies:  intPtr(3),
		},
		IsBuiltIn: true,
	},
}

// VariantService provides access to CV variants
type VariantService interface {
	// ListVariants returns all available variants
	ListVariants() []*CVVariant

	// GetVariant returns a variant by ID
	GetVariant(id string) (*CVVariant, error)

	// GetVariantByDimensions finds a variant by role emphasis and length format
	GetVariantByDimensions(role RoleEmphasis, length LengthFormat) (*CVVariant, error)

	// ListRoleEmphases returns available role emphases with metadata
	ListRoleEmphases() []RoleEmphasisInfo

	// ListLengthFormats returns available length formats with metadata
	ListLengthFormats() []LengthFormatInfo
}

// DefaultVariantService is the default implementation of VariantService
type DefaultVariantService struct {
	variants     map[string]*CVVariant
	byDimensions map[string]*CVVariant // key: "role:length"
}

// NewVariantService creates a new VariantService
func NewVariantService() VariantService {
	service := &DefaultVariantService{
		variants:     make(map[string]*CVVariant),
		byDimensions: make(map[string]*CVVariant),
	}

	// Index built-in variants
	for _, v := range BuiltInVariants {
		service.variants[v.ID] = v
		key := fmt.Sprintf("%s:%s", v.RoleEmphasis, v.LengthFormat)
		service.byDimensions[key] = v
	}

	return service
}

// ListVariants returns all available variants
func (s *DefaultVariantService) ListVariants() []*CVVariant {
	return BuiltInVariants
}

// GetVariant returns a variant by ID
func (s *DefaultVariantService) GetVariant(id string) (*CVVariant, error) {
	v, exists := s.variants[id]
	if !exists {
		return nil, fmt.Errorf("variant not found: %s", id)
	}
	return v, nil
}

// GetVariantByDimensions finds a variant by role emphasis and length format
func (s *DefaultVariantService) GetVariantByDimensions(role RoleEmphasis, length LengthFormat) (*CVVariant, error) {
	key := fmt.Sprintf("%s:%s", role, length)
	v, exists := s.byDimensions[key]
	if !exists {
		return nil, fmt.Errorf("variant not found for role=%s, length=%s", role, length)
	}
	return v, nil
}

// ListRoleEmphases returns available role emphases with metadata
func (s *DefaultVariantService) ListRoleEmphases() []RoleEmphasisInfo {
	return []RoleEmphasisInfo{
		{
			ID:          RoleEmphasisSeniorBackend,
			Name:        "Senior Backend",
			Description: "Emphasizes backend engineering skills, technical depth, and system design",
		},
		{
			ID:          RoleEmphasisStaffPrincipal,
			Name:        "Staff/Principal",
			Description: "Emphasizes technical leadership, architecture, and cross-team impact",
		},
		{
			ID:          RoleEmphasisConsulting,
			Name:        "Consulting",
			Description: "Emphasizes client engagements, adaptability, and business impact",
		},
		{
			ID:          RoleEmphasisLanguageAgnostic,
			Name:        "Language-Agnostic",
			Description: "Emphasizes transferable skills without technology-specific focus",
		},
	}
}

// ListLengthFormats returns available length formats with metadata
func (s *DefaultVariantService) ListLengthFormats() []LengthFormatInfo {
	return []LengthFormatInfo{
		{
			ID:          LengthFull,
			Name:        "Full",
			Description: "Complete career history with all relevant details (3+ pages)",
		},
		{
			ID:          LengthStandard,
			Name:        "Standard",
			Description: "Balanced format with recent 10 years emphasis (2-3 pages)",
		},
		{
			ID:          LengthShort,
			Name:        "Short",
			Description: "Concise format focusing on recent 5 years (1-2 pages)",
		},
		{
			ID:          LengthUltraShort,
			Name:        "1-Page",
			Description: "Key highlights only for quick review (1 page)",
		},
	}
}

// Helper functions for pointer values
func intPtr(i int) *int {
	return &i
}

func floatPtr(f float64) *float64 {
	return &f
}

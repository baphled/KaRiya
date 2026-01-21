// Package constants provides centralized type-safe constants for the KaRiya application.
// This package is a leaf package with no internal dependencies, ensuring a single source
// of truth for all enumerations and constant values used throughout the application.
package constants

// Role represents a career role/level target for CVs and facts.
type Role string

const (
	RolePrincipal Role = "principal"
	RoleEM        Role = "em"
	RoleStaff     Role = "staff"
	RoleSeniorIC  Role = "senior_ic"
)

// AllRoles returns all defined roles.
func AllRoles() []Role {
	return []Role{
		RolePrincipal,
		RoleEM,
		RoleStaff,
		RoleSeniorIC,
	}
}

// IsValidRole checks if the given string is a valid role.
func IsValidRole(s string) bool {
	for _, r := range AllRoles() {
		if string(r) == s {
			return true
		}
	}
	return false
}

// CVTargetRoles returns roles suitable for CV targeting.
func CVTargetRoles() []Role {
	return AllRoles()
}

// CompetencyCategory represents broad professional competency categories.
type CompetencyCategory string

const (
	CompetencyTechnical  CompetencyCategory = "technical"
	CompetencyLeadership CompetencyCategory = "leadership"
	CompetencyProduct    CompetencyCategory = "product"
	CompetencyConsulting CompetencyCategory = "consulting"
	CompetencyResearch   CompetencyCategory = "research"
	CompetencyMentoring  CompetencyCategory = "mentoring"
)

// AllCompetencyCategories returns all defined competency categories.
func AllCompetencyCategories() []CompetencyCategory {
	return []CompetencyCategory{
		CompetencyTechnical,
		CompetencyLeadership,
		CompetencyProduct,
		CompetencyConsulting,
		CompetencyResearch,
		CompetencyMentoring,
	}
}

// IsValidCompetencyCategory checks if the given string is a valid competency category.
func IsValidCompetencyCategory(s string) bool {
	for _, c := range AllCompetencyCategories() {
		if string(c) == s {
			return true
		}
	}
	return false
}

// SkillLevel represents the proficiency level of a skill.
type SkillLevel string

const (
	SkillLevelBeginner     SkillLevel = "beginner"
	SkillLevelIntermediate SkillLevel = "intermediate"
	SkillLevelAdvanced     SkillLevel = "advanced"
	SkillLevelExpert       SkillLevel = "expert"
)

// AllSkillLevels returns all defined skill levels.
func AllSkillLevels() []SkillLevel {
	return []SkillLevel{
		SkillLevelBeginner,
		SkillLevelIntermediate,
		SkillLevelAdvanced,
		SkillLevelExpert,
	}
}

// IsValidSkillLevel checks if the given string is a valid skill level.
func IsValidSkillLevel(s string) bool {
	for _, l := range AllSkillLevels() {
		if string(l) == s {
			return true
		}
	}
	return false
}

// EventTag represents tags that can be applied to career events.
type EventTag string

const (
	EventTagProject     EventTag = "project"
	EventTagAchievement EventTag = "achievement"
	EventTagLeadership  EventTag = "leadership"
	EventTagTechnical   EventTag = "technical"
	EventTagConsulting  EventTag = "consulting"
	EventTagResearch    EventTag = "research"
	EventTagProduct     EventTag = "product"
	EventTagMentoring   EventTag = "mentoring"
)

// AllEventTags returns all defined event tags.
func AllEventTags() []EventTag {
	return []EventTag{
		EventTagProject,
		EventTagAchievement,
		EventTagLeadership,
		EventTagTechnical,
		EventTagConsulting,
		EventTagResearch,
		EventTagProduct,
		EventTagMentoring,
	}
}

// IsValidEventTag checks if the given string is a valid event tag.
func IsValidEventTag(s string) bool {
	for _, t := range AllEventTags() {
		if string(t) == s {
			return true
		}
	}
	return false
}

// Audience represents the target audience for CVs and facts.
type Audience string

const (
	AudienceHiringManager Audience = "hiring_manager"
	AudienceRecruiter     Audience = "recruiter"
	AudiencePeer          Audience = "peer"
)

// AllAudiences returns all defined audiences.
func AllAudiences() []Audience {
	return []Audience{
		AudienceHiringManager,
		AudienceRecruiter,
		AudiencePeer,
	}
}

// IsValidAudience checks if the given string is a valid audience.
func IsValidAudience(s string) bool {
	for _, a := range AllAudiences() {
		if string(a) == s {
			return true
		}
	}
	return false
}

// SkillCategory represents suggested categories for organizing skills.
// Note: These are suggestions, not constraints. Users can define custom categories.
type SkillCategory string

const (
	SkillCategoryBackend  SkillCategory = "backend"
	SkillCategoryFrontend SkillCategory = "frontend"
	SkillCategoryDevOps   SkillCategory = "devops"
	SkillCategoryDatabase SkillCategory = "database"
	SkillCategoryCloud    SkillCategory = "cloud"
	SkillCategoryMobile   SkillCategory = "mobile"
	SkillCategoryTooling  SkillCategory = "tooling"
	SkillCategoryOther    SkillCategory = "other"
)

// SuggestedSkillCategories returns all suggested skill categories.
// This includes "mobile" which was previously missing from forms.
func SuggestedSkillCategories() []SkillCategory {
	return []SkillCategory{
		SkillCategoryBackend,
		SkillCategoryFrontend,
		SkillCategoryDevOps,
		SkillCategoryDatabase,
		SkillCategoryCloud,
		SkillCategoryMobile,
		SkillCategoryTooling,
		SkillCategoryOther,
	}
}

// SectionType represents the types of sections in a CV.
type SectionType string

const (
	SectionTypeExperience SectionType = "experience"
	SectionTypeProjects   SectionType = "projects"
	SectionTypeSkills     SectionType = "skills"
	SectionTypeSummary    SectionType = "summary"
)

// AllSectionTypes returns all defined section types.
func AllSectionTypes() []SectionType {
	return []SectionType{
		SectionTypeExperience,
		SectionTypeProjects,
		SectionTypeSkills,
		SectionTypeSummary,
	}
}

// IsValidSectionType checks if the given string is a valid section type.
func IsValidSectionType(s string) bool {
	for _, t := range AllSectionTypes() {
		if string(t) == s {
			return true
		}
	}
	return false
}

// ImpactLevel represents the impact level of a CV bullet.
type ImpactLevel string

const (
	ImpactLevelLow    ImpactLevel = "low"
	ImpactLevelMedium ImpactLevel = "medium"
	ImpactLevelHigh   ImpactLevel = "high"
)

// AllImpactLevels returns all defined impact levels (excluding empty).
func AllImpactLevels() []ImpactLevel {
	return []ImpactLevel{
		ImpactLevelLow,
		ImpactLevelMedium,
		ImpactLevelHigh,
	}
}

// IsValidImpactLevel checks if the given string is a valid impact level.
// Empty string is valid (represents "not set").
func IsValidImpactLevel(s string) bool {
	if s == "" {
		return true
	}
	for _, l := range AllImpactLevels() {
		if string(l) == s {
			return true
		}
	}
	return false
}

// InclusionReason represents the reason a bullet is included in a CV.
type InclusionReason string

// Semantic reasons (manual categorization)
const (
	InclusionReasonOwnership    InclusionReason = "ownership"
	InclusionReasonContribution InclusionReason = "contribution"
	InclusionReasonStrategy     InclusionReason = "strategy"
	InclusionReasonExecution    InclusionReason = "execution"
	InclusionReasonOutcome      InclusionReason = "outcome"
	InclusionReasonActivity     InclusionReason = "activity"
)

// Source-based reasons (used by BulletGenerator and EnhancedBulletGenerator)
const (
	InclusionReasonFactExtraction        InclusionReason = "fact_extraction"
	InclusionReasonEventDirect           InclusionReason = "event_direct"
	InclusionReasonAchievementExtraction InclusionReason = "achievement_extraction"
)

// AllInclusionReasons returns all defined inclusion reasons.
func AllInclusionReasons() []InclusionReason {
	return []InclusionReason{
		// Semantic reasons
		InclusionReasonOwnership,
		InclusionReasonContribution,
		InclusionReasonStrategy,
		InclusionReasonExecution,
		InclusionReasonOutcome,
		InclusionReasonActivity,
		// Source-based reasons
		InclusionReasonFactExtraction,
		InclusionReasonEventDirect,
		InclusionReasonAchievementExtraction,
	}
}

// IsValidInclusionReason checks if the given string is a valid inclusion reason.
func IsValidInclusionReason(s string) bool {
	for _, r := range AllInclusionReasons() {
		if string(r) == s {
			return true
		}
	}
	return false
}

// FocusArea represents a career focus area based on skill categories.
type FocusArea string

const (
	FocusAreaBackend   FocusArea = "backend"
	FocusAreaFrontend  FocusArea = "frontend"
	FocusAreaFullstack FocusArea = "fullstack"
	FocusAreaDevOps    FocusArea = "devops"
)

// AllFocusAreas returns all defined focus areas.
func AllFocusAreas() []FocusArea {
	return []FocusArea{
		FocusAreaBackend,
		FocusAreaFrontend,
		FocusAreaFullstack,
		FocusAreaDevOps,
	}
}

// IsValidFocusArea checks if the given string is a valid focus area.
func IsValidFocusArea(s string) bool {
	for _, a := range AllFocusAreas() {
		if string(a) == s {
			return true
		}
	}
	return false
}

// TechnologyFocus represents how technology should be emphasized in a CV.
type TechnologyFocus string

const (
	TechnologyFocusLanguageAgnostic TechnologyFocus = "language_agnostic"
	TechnologyFocusGeneralist       TechnologyFocus = "generalist"
	TechnologyFocusSpecialist       TechnologyFocus = "specialist"
)

// AllTechnologyFocuses returns all defined technology focuses.
func AllTechnologyFocuses() []TechnologyFocus {
	return []TechnologyFocus{
		TechnologyFocusLanguageAgnostic,
		TechnologyFocusGeneralist,
		TechnologyFocusSpecialist,
	}
}

// IsValidTechnologyFocus checks if the given string is a valid technology focus.
func IsValidTechnologyFocus(s string) bool {
	for _, f := range AllTechnologyFocuses() {
		if string(f) == s {
			return true
		}
	}
	return false
}

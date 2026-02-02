// Package constants provides centralized type-safe enumerations for the KaRiya
// career management application. It is a leaf package with no internal dependencies,
// serving as the single source of truth for domain values such as role targets,
// competency areas, skill proficiencies, and CV structure metadata. Other packages
// import these constants to ensure consistent validation and categorization across
// the domain, service, and CLI layers.
package constants

// RoleFit classifies the seniority or leadership level a CV should target. It
// drives content selection and phrasing so that career events, facts, and bullet
// points are framed in terms appropriate for the intended position.
type RoleFit string

// RoleFit values control how CV content is tailored. Each level emphasizes
// different career accomplishments: principal highlights system-wide architectural
// impact, em highlights people management and org-level outcomes, staff highlights
// deep technical expertise, and senior_ic highlights sustained individual output.
const (
	// RoleFitPrincipal targets CV content toward principal engineer positions,
	// emphasizing system design, technical vision, and cross-team impact.
	RoleFitPrincipal RoleFit = "principal"
	// RoleFitEM targets CV content toward engineering manager positions,
	// emphasizing team building, process improvement, and delivery outcomes.
	RoleFitEM RoleFit = "em"
	// RoleFitStaff targets CV content toward staff engineer positions,
	// emphasizing deep technical expertise, mentorship, and complex problem solving.
	RoleFitStaff RoleFit = "staff"
	// RoleFitSeniorIC targets CV content toward senior individual contributor positions,
	// emphasizing hands-on execution, feature ownership, and measurable output.
	RoleFitSeniorIC RoleFit = "senior_ic"
)

// AllRoleFits returns every defined RoleFit value in declaration order.
//
// Returns:
//   - A []RoleFit value.
//
// Side effects:
//   - None.
func AllRoleFits() []RoleFit {
	return []RoleFit{
		RoleFitPrincipal,
		RoleFitEM,
		RoleFitStaff,
		RoleFitSeniorIC,
	}
}

// IsValidRoleFit reports whether s matches a recognised RoleFit value.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func IsValidRoleFit(s string) bool {
	for _, r := range AllRoleFits() {
		if string(r) == s {
			return true
		}
	}
	return false
}

// CVTargetRoleFits returns the RoleFit values that are valid targets when
//
// Returns:
//   - A []RoleFit value.
//
// Side effects:
//   - None.
func CVTargetRoleFits() []RoleFit {
	return AllRoleFits()
}

// CompetencyCategory models the broad professional discipline that a career fact
// or event demonstrates. It is used to classify accomplishments so the CV generator
// can balance content across different competency areas.
type CompetencyCategory string

// CompetencyCategory values partition professional capabilities into orthogonal
// disciplines. CV generation uses these to ensure well-rounded representation,
// and the timeline browser uses them to filter and group career events.
const (
	// CompetencyTechnical covers engineering skills such as architecture,
	// coding, debugging, and systems design.
	CompetencyTechnical CompetencyCategory = "technical"
	// CompetencyLeadership covers people and organisational skills such as
	// team management, hiring, performance reviews, and strategic direction.
	CompetencyLeadership CompetencyCategory = "leadership"
	// CompetencyProduct covers product thinking skills such as roadmap
	// prioritisation, user research, stakeholder communication, and feature scoping.
	CompetencyProduct CompetencyCategory = "product"
	// CompetencyConsulting covers client-facing advisory skills such as
	// requirements gathering, solution design, stakeholder workshops, and delivery.
	CompetencyConsulting CompetencyCategory = "consulting"
	// CompetencyResearch covers investigative skills such as prototyping,
	// experimentation, literature review, and technical spike execution.
	CompetencyResearch CompetencyCategory = "research"
	// CompetencyMentoring covers knowledge transfer skills such as coaching
	// junior engineers, conducting code reviews, and running workshops.
	CompetencyMentoring CompetencyCategory = "mentoring"

	// CompetencyCommunication covers presenting, documenting, and explaining
	// technical concepts to diverse audiences including stakeholders and teams.
	CompetencyCommunication CompetencyCategory = "communication"
	// CompetencyCollaboration covers cross-functional teamwork, partnering with
	// other teams, and coordinating efforts across organisational boundaries.
	CompetencyCollaboration CompetencyCategory = "collaboration"
	// CompetencyProblemSolving covers analytical debugging, troubleshooting
	// complex issues, and systematic root-cause investigation.
	CompetencyProblemSolving CompetencyCategory = "problem-solving"
	// CompetencyProjectManagement covers planning sprints, tracking milestones,
	// estimating effort, and delivering projects on schedule.
	CompetencyProjectManagement CompetencyCategory = "project-management"
	// CompetencyArchitecture covers system design, scalability decisions,
	// distributed systems thinking, and technical design authority.
	CompetencyArchitecture CompetencyCategory = "architecture"
)

// AllCompetencyCategories returns every defined CompetencyCategory value in
//
// Returns:
//   - A []CompetencyCategory value.
//
// Side effects:
//   - None.
func AllCompetencyCategories() []CompetencyCategory {
	return []CompetencyCategory{
		CompetencyTechnical,
		CompetencyLeadership,
		CompetencyProduct,
		CompetencyConsulting,
		CompetencyResearch,
		CompetencyMentoring,
		CompetencyCommunication,
		CompetencyCollaboration,
		CompetencyProblemSolving,
		CompetencyProjectManagement,
		CompetencyArchitecture,
	}
}

// IsValidCompetencyCategory reports whether s matches a recognised
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func IsValidCompetencyCategory(s string) bool {
	for _, c := range AllCompetencyCategories() {
		if string(c) == s {
			return true
		}
	}
	return false
}

// SkillLevel models a person's self-assessed proficiency with a particular
// technology or practice. It is stored alongside skills and used by the CV
// generator to decide which skills to highlight and how prominently to feature them.
type SkillLevel string

// SkillLevel values form an ordered progression from novice to authority-level
// mastery. The CV generator may use these to filter out low-proficiency skills
// when space is limited, or to sort skills by strength.
const (
	// SkillLevelBeginner indicates foundational awareness with limited
	// practical experience, typically less than one year of active use.
	SkillLevelBeginner SkillLevel = "beginner"
	// SkillLevelIntermediate indicates working proficiency with the ability
	// to complete tasks independently, typically one to three years of use.
	SkillLevelIntermediate SkillLevel = "intermediate"
	// SkillLevelAdvanced indicates strong command with the ability to mentor
	// others and tackle complex problems, typically three to six years of use.
	SkillLevelAdvanced SkillLevel = "advanced"
	// SkillLevelExpert indicates authoritative mastery with deep knowledge of
	// internals, edge cases, and best practices, often used as a go-to resource.
	SkillLevelExpert SkillLevel = "expert"
)

// AllSkillLevels returns every defined SkillLevel value in ascending order of
//
// Returns:
//   - A []SkillLevel value.
//
// Side effects:
//   - None.
func AllSkillLevels() []SkillLevel {
	return []SkillLevel{
		SkillLevelBeginner,
		SkillLevelIntermediate,
		SkillLevelAdvanced,
		SkillLevelExpert,
	}
}

// IsValidSkillLevel reports whether s matches a recognised SkillLevel value.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func IsValidSkillLevel(s string) bool {
	for _, l := range AllSkillLevels() {
		if string(l) == s {
			return true
		}
	}
	return false
}

// EventTag classifies career events by their nature so they can be filtered,
// grouped, and weighted during CV generation. A single event may carry multiple
// tags to reflect overlapping concerns (for example, a project that also
// demonstrated leadership).
type EventTag string

// EventTag values label the primary nature of a career event. The timeline
// browser uses these for filtering, and the CV generator uses them to ensure
// coverage of the categories that matter most for the target role.
const (
	// EventTagProject marks events that describe a discrete body of work
	// with a defined scope, timeline, and deliverable.
	EventTagProject EventTag = "project"
	// EventTagAchievement marks events that highlight a notable outcome,
	// award, or measurable result such as a performance improvement or cost saving.
	EventTagAchievement EventTag = "achievement"
	// EventTagLeadership marks events that demonstrate people management,
	// organisational influence, or decision-making authority.
	EventTagLeadership EventTag = "leadership"
	// EventTagTechnical marks events focused on engineering work such as
	// building systems, solving hard bugs, or introducing new tools.
	EventTagTechnical EventTag = "technical"
	// EventTagConsulting marks events involving client-facing advisory work,
	// requirements workshops, or external stakeholder delivery.
	EventTagConsulting EventTag = "consulting"
	// EventTagResearch marks events centered on investigation, prototyping,
	// or exploratory work that informs future engineering decisions.
	EventTagResearch EventTag = "research"
	// EventTagProduct marks events related to product definition, roadmap
	// influence, user research, or feature prioritisation.
	EventTagProduct EventTag = "product"
	// EventTagMentoring marks events involving coaching, onboarding, code
	// review, or knowledge sharing with other engineers.
	EventTagMentoring EventTag = "mentoring"
)

// AllEventTags returns every defined EventTag value in declaration order.
//
// Returns:
//   - A []EventTag value.
//
// Side effects:
//   - None.
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

// IsValidEventTag reports whether s matches a recognised EventTag value.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func IsValidEventTag(s string) bool {
	for _, t := range AllEventTags() {
		if string(t) == s {
			return true
		}
	}
	return false
}

// Audience identifies who will read the generated CV, allowing the generator
// to adjust tone, detail level, and emphasis. A hiring manager cares about
// impact and fit, a recruiter scans for keywords and seniority signals, and
// a peer evaluates technical depth.
type Audience string

// Audience values determine the voice and focus of generated CV content.
// Each audience expects different information density and framing, so the
// generator tailors bullet wording and section ordering accordingly.
const (
	// AudienceHiringManager targets someone evaluating culture fit, leadership
	// potential, and business impact, favouring outcome-driven language.
	AudienceHiringManager Audience = "hiring_manager"
	// AudienceRecruiter targets someone screening for role match and seniority
	// signals, favouring keyword-rich and concise bullet points.
	AudienceRecruiter Audience = "recruiter"
	// AudiencePeer targets a fellow engineer assessing technical credibility,
	// favouring concrete implementation details and architectural reasoning.
	AudiencePeer Audience = "peer"
)

// AllAudiences returns every defined Audience value in declaration order.
//
// Returns:
//   - A []Audience value.
//
// Side effects:
//   - None.
func AllAudiences() []Audience {
	return []Audience{
		AudienceHiringManager,
		AudienceRecruiter,
		AudiencePeer,
	}
}

// IsValidAudience reports whether s matches a recognised Audience value.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func IsValidAudience(s string) bool {
	for _, a := range AllAudiences() {
		if string(a) == s {
			return true
		}
	}
	return false
}

// SkillCategory provides a suggested grouping for technical skills on a CV.
// Categories are advisory rather than mandatory; users may define their own
// custom categories when the predefined ones do not fit.
type SkillCategory string

// SkillCategory values offer sensible defaults for organising skills into
// sections on a CV. The skills form and the CV renderer both use these to
// group related technologies together for readability.
const (
	// SkillCategoryBackend groups server-side languages, frameworks, and
	// runtime technologies such as Go, Java, Node.js, and REST API design.
	SkillCategoryBackend SkillCategory = "backend"
	// SkillCategoryFrontend groups client-side technologies such as
	// JavaScript frameworks, CSS, HTML, and browser APIs.
	SkillCategoryFrontend SkillCategory = "frontend"
	// SkillCategoryDevOps groups infrastructure and delivery pipeline
	// technologies such as Docker, Kubernetes, Terraform, and CI/CD tooling.
	SkillCategoryDevOps SkillCategory = "devops"
	// SkillCategoryDatabase groups data storage technologies such as
	// PostgreSQL, MongoDB, Redis, and query optimisation expertise.
	SkillCategoryDatabase SkillCategory = "database"
	// SkillCategoryCloud groups cloud platform services and architecture
	// patterns such as AWS, GCP, Azure, and serverless design.
	SkillCategoryCloud SkillCategory = "cloud"
	// SkillCategoryMobile groups mobile development technologies such as
	// Swift, Kotlin, React Native, and platform-specific APIs.
	SkillCategoryMobile SkillCategory = "mobile"
	// SkillCategoryTooling groups developer productivity tools such as
	// editors, linters, profilers, build systems, and version control.
	SkillCategoryTooling SkillCategory = "tooling"
	// SkillCategoryTesting groups quality assurance technologies such as
	// test frameworks, assertion libraries, and code coverage tools.
	SkillCategoryTesting SkillCategory = "testing"
	// SkillCategoryData groups data engineering technologies such as
	// ETL pipelines, analytics platforms, and business intelligence tools.
	SkillCategoryData SkillCategory = "data"
	// SkillCategoryML groups machine learning and artificial intelligence
	// technologies such as TensorFlow, PyTorch, and LLM tooling.
	SkillCategoryML SkillCategory = "ml"
	// SkillCategoryMonitoring groups observability technologies such as
	// Prometheus, Grafana, ELK stack, and alerting systems.
	SkillCategoryMonitoring SkillCategory = "monitoring"
	// SkillCategoryArchitecture groups system design patterns and
	// architectural concepts such as microservices, DDD, and CQRS.
	SkillCategoryArchitecture SkillCategory = "architecture"
	// SkillCategorySecurity groups security and compliance technologies
	// such as OAuth, TLS, encryption, and access control systems.
	SkillCategorySecurity SkillCategory = "security"
	// SkillCategoryPractices groups engineering methodologies and practices
	// such as Agile, TDD, code review, and incident management.
	SkillCategoryPractices SkillCategory = "practices"
	// SkillCategoryOther is a catch-all for skills that do not fit neatly
	// into the predefined categories above.
	SkillCategoryOther SkillCategory = "other"
)

// AllSkillCategories returns every defined SkillCategory value in declaration
//
// Returns:
//   - A []SkillCategory value.
//
// Side effects:
//   - None.
func AllSkillCategories() []SkillCategory {
	return []SkillCategory{
		SkillCategoryBackend,
		SkillCategoryFrontend,
		SkillCategoryDevOps,
		SkillCategoryDatabase,
		SkillCategoryCloud,
		SkillCategoryMobile,
		SkillCategoryTooling,
		SkillCategoryTesting,
		SkillCategoryData,
		SkillCategoryML,
		SkillCategoryMonitoring,
		SkillCategoryArchitecture,
		SkillCategorySecurity,
		SkillCategoryPractices,
		SkillCategoryOther,
	}
}

// SuggestedSkillCategories returns every predefined SkillCategory value in
//
// Returns:
//   - A []SkillCategory value.
//
// Side effects:
//   - None.
func SuggestedSkillCategories() []SkillCategory {
	return AllSkillCategories()
}

// IsValidSkillCategory reports whether s matches a recognised SkillCategory
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func IsValidSkillCategory(s string) bool {
	for _, c := range AllSkillCategories() {
		if string(c) == s {
			return true
		}
	}
	return false
}

// SkillCategoryStrings returns a string slice of all skill category values
//
// Returns:
//   - A []string value.
//
// Side effects:
//   - None.
func SkillCategoryStrings() []string {
	cats := AllSkillCategories()
	result := make([]string, len(cats))
	for i, c := range cats {
		result[i] = string(c)
	}
	return result
}

// SectionType identifies the kind of content block within a generated CV. The
// CV renderer uses it to apply the correct layout, ordering, and formatting
// rules for each block.
type SectionType string

// SectionType values enumerate the structural building blocks of a CV. Each
// type has distinct rendering rules: experience shows company-grouped bullets,
// projects shows standalone deliverables, skills shows a categorised grid, and
// summary shows a prose paragraph.
const (
	// SectionTypeExperience identifies a work history block that groups
	// bullet points under employer names and date ranges.
	SectionTypeExperience SectionType = "experience"
	// SectionTypeProjects identifies a block listing standalone projects,
	// open-source contributions, or side work outside of employment.
	SectionTypeProjects SectionType = "projects"
	// SectionTypeSkills identifies a block presenting technical and
	// professional skills, typically rendered as a categorised grid or list.
	SectionTypeSkills SectionType = "skills"
	// SectionTypeSummary identifies the introductory prose paragraph at the
	// top of a CV that frames the candidate's career narrative.
	SectionTypeSummary SectionType = "summary"
)

// CV grouping constants control how career events are organised into
// company-based sections. They handle edge cases such as missing employer
// names and multiple tenures at the same organisation.
const (
	// DefaultCompanyName is the fallback label applied when a career event
	// has no employer specified, ensuring every event belongs to a group.
	DefaultCompanyName = "Other"

	// TenureSeparator delimits the numeric suffix appended to a company name
	// when the user held multiple distinct tenures there. For example,
	// "Acme Corp#1" and "Acme Corp#2" represent two separate stints.
	TenureSeparator = "#"
)

// AllSectionTypes returns every defined SectionType value in declaration order.
//
// Returns:
//   - A []SectionType value.
//
// Side effects:
//   - None.
func AllSectionTypes() []SectionType {
	return []SectionType{
		SectionTypeExperience,
		SectionTypeProjects,
		SectionTypeSkills,
		SectionTypeSummary,
	}
}

// IsValidSectionType reports whether s matches a recognised SectionType value.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func IsValidSectionType(s string) bool {
	for _, t := range AllSectionTypes() {
		if string(t) == s {
			return true
		}
	}
	return false
}

// ImpactLevel quantifies the significance of a CV bullet point. The CV
// generator uses it to prioritise which bullets appear when space is limited,
// preferring high-impact items over low-impact ones.
type ImpactLevel string

// ImpactLevel values rank bullet points by their business or technical
// significance. When the CV must be condensed, the generator drops low-impact
// bullets first, then medium, keeping high-impact items visible.
const (
	// ImpactLevelLow marks a bullet as routine work that is useful context
	// but not essential, such as minor bug fixes or incremental improvements.
	ImpactLevelLow ImpactLevel = "low"
	// ImpactLevelMedium marks a bullet as meaningful work that demonstrates
	// competence, such as feature delivery or moderate process improvements.
	ImpactLevelMedium ImpactLevel = "medium"
	// ImpactLevelHigh marks a bullet as a standout accomplishment with
	// significant business or technical outcomes, such as major cost savings
	// or system redesigns.
	ImpactLevelHigh ImpactLevel = "high"
)

// AllImpactLevels returns every defined ImpactLevel value in ascending order
//
// Returns:
//   - A []ImpactLevel value.
//
// Side effects:
//   - None.
func AllImpactLevels() []ImpactLevel {
	return []ImpactLevel{
		ImpactLevelLow,
		ImpactLevelMedium,
		ImpactLevelHigh,
	}
}

// IsValidImpactLevel reports whether s matches a recognised ImpactLevel value.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
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

// InclusionReason records why a particular bullet was added to a CV. It enables
// auditing of generated content and helps users understand the provenance of
// each bullet point when reviewing or editing their CV.
type InclusionReason string

// Semantic inclusion reasons describe the narrative role a bullet plays in
// the CV. These are typically assigned during manual curation to explain what
// aspect of the candidate's work the bullet is meant to showcase.
const (
	// InclusionReasonOwnership indicates the bullet demonstrates end-to-end
	// responsibility for a system, product area, or initiative.
	InclusionReasonOwnership InclusionReason = "ownership"
	// InclusionReasonContribution indicates the bullet highlights a specific
	// hands-on technical or collaborative contribution to a larger effort.
	InclusionReasonContribution InclusionReason = "contribution"
	// InclusionReasonStrategy indicates the bullet showcases strategic
	// thinking such as roadmap planning, architectural vision, or long-term
	// decision making.
	InclusionReasonStrategy InclusionReason = "strategy"
	// InclusionReasonExecution indicates the bullet highlights the ability to
	// ship reliably, meeting deadlines and quality standards under pressure.
	InclusionReasonExecution InclusionReason = "execution"
	// InclusionReasonOutcome indicates the bullet emphasises a measurable
	// result such as revenue growth, latency reduction, or user adoption.
	InclusionReasonOutcome InclusionReason = "outcome"
	// InclusionReasonActivity indicates the bullet describes day-to-day work
	// that provides context without a specific standout result.
	InclusionReasonActivity InclusionReason = "activity"
)

// Source-based inclusion reasons record the automated pipeline stage that
// produced the bullet. These are assigned by the BulletGenerator during CV
// assembly and are not typically set by users.
const (
	// InclusionReasonFactExtraction indicates the bullet was synthesised from
	// one or more extracted career facts by the fact extraction pipeline.
	InclusionReasonFactExtraction InclusionReason = "fact_extraction"
	// InclusionReasonEventDirect indicates the bullet was derived directly
	// from a single career event without additional synthesis.
	InclusionReasonEventDirect InclusionReason = "event_direct"
	// InclusionReasonAchievementExtraction indicates the bullet was generated
	// from an achievement record identified during career event analysis.
	InclusionReasonAchievementExtraction InclusionReason = "achievement_extraction"
)

// AllInclusionReasons returns every defined InclusionReason value in
//
// Returns:
//   - A []InclusionReason value.
//
// Side effects:
//   - None.
func AllInclusionReasons() []InclusionReason {
	return []InclusionReason{
		InclusionReasonOwnership,
		InclusionReasonContribution,
		InclusionReasonStrategy,
		InclusionReasonExecution,
		InclusionReasonOutcome,
		InclusionReasonActivity,
		InclusionReasonFactExtraction,
		InclusionReasonEventDirect,
		InclusionReasonAchievementExtraction,
	}
}

// IsValidInclusionReason reports whether s matches a recognised
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func IsValidInclusionReason(s string) bool {
	for _, r := range AllInclusionReasons() {
		if string(r) == s {
			return true
		}
	}
	return false
}

// FocusArea identifies the engineering discipline a CV should emphasise. It
// influences which career events and skills are prioritised during generation,
// ensuring the resulting CV aligns with the type of role being targeted.
type FocusArea string

// FocusArea values steer the CV generator toward the engineering discipline
// that best matches the target position. Backend and frontend focus on their
// respective stacks, fullstack balances both, and devops highlights
// infrastructure and delivery pipeline expertise.
const (
	// FocusAreaBackend prioritises server-side, API, and systems engineering
	// content, de-emphasising frontend and UI work.
	FocusAreaBackend FocusArea = "backend"
	// FocusAreaFrontend prioritises client-side, UI, and user experience
	// content, de-emphasising infrastructure and backend work.
	FocusAreaFrontend FocusArea = "frontend"
	// FocusAreaFullstack balances backend and frontend content equally,
	// suitable for roles that require end-to-end feature delivery.
	FocusAreaFullstack FocusArea = "fullstack"
	// FocusAreaDevOps prioritises infrastructure, CI/CD, cloud architecture,
	// and site reliability content over application-level work.
	FocusAreaDevOps FocusArea = "devops"
)

// AllFocusAreas returns every defined FocusArea value in declaration order.
//
// Returns:
//   - A []FocusArea value.
//
// Side effects:
//   - None.
func AllFocusAreas() []FocusArea {
	return []FocusArea{
		FocusAreaBackend,
		FocusAreaFrontend,
		FocusAreaFullstack,
		FocusAreaDevOps,
	}
}

// IsValidFocusArea reports whether s matches a recognised FocusArea value.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func IsValidFocusArea(s string) bool {
	for _, a := range AllFocusAreas() {
		if string(a) == s {
			return true
		}
	}
	return false
}

// TechnologyFocus controls how prominently specific technologies are named in
// the generated CV. It lets users choose between emphasising transferable
// skills versus deep expertise in particular tools.
type TechnologyFocus string

// TechnologyFocus values tell the CV generator how to handle technology names.
// Language-agnostic omits specific tool names in favour of abstract concepts,
// generalist names tools but keeps the emphasis on breadth, and specialist
// foregrounds deep expertise in a narrow set of technologies.
const (
	// TechnologyFocusLanguageAgnostic omits specific language and framework
	// names, framing accomplishments in terms of design patterns, paradigms,
	// and transferable engineering principles.
	TechnologyFocusLanguageAgnostic TechnologyFocus = "language_agnostic"
	// TechnologyFocusGeneralist names technologies used but keeps the focus
	// broad, suitable for roles valuing versatility across multiple stacks.
	TechnologyFocusGeneralist TechnologyFocus = "generalist"
	// TechnologyFocusSpecialist foregrounds specific technologies and deep
	// expertise, suitable for roles requiring mastery of a particular stack.
	TechnologyFocusSpecialist TechnologyFocus = "specialist"
)

// AllTechnologyFocuses returns every defined TechnologyFocus value in
//
// Returns:
//   - A []TechnologyFocus value.
//
// Side effects:
//   - None.
func AllTechnologyFocuses() []TechnologyFocus {
	return []TechnologyFocus{
		TechnologyFocusLanguageAgnostic,
		TechnologyFocusGeneralist,
		TechnologyFocusSpecialist,
	}
}

// IsValidTechnologyFocus reports whether s matches a recognised
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func IsValidTechnologyFocus(s string) bool {
	for _, f := range AllTechnologyFocuses() {
		if string(f) == s {
			return true
		}
	}
	return false
}

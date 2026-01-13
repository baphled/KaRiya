# Task 40: Role Emphasis Redesign - Technology-Focused CV Generation

## Overview
- **Goal**: Replace current role emphasis with technology-focused system that uses user-defined skills for CV generation
- **Time Estimate**: 3-4 days
- **Prerequisites**: Task 39 (User-Defined Skills) must be complete

## Session Contract Acknowledgment
- [x] Ran `make session-start` and it passed (2026-01-13)
- [x] Acknowledge and commit to following all workflow rules
- [x] Token count: 46,176 (< 50k ✅)

## Pre-Task Checklist (MUST COMPLETE BEFORE STARTING)
- [x] `make check-compliance` passes ✅
- [x] Task 39 (User-Defined Skills) is complete ✅
- [x] Reviewed existing patterns in:
  - `internal/service/career/cv/` (no variants.go or role_emphasis.go - will create)
  - `internal/cli/intents/generate_cv_intent.go` (CV generation flow)
  - `internal/domain/career/skill.go` (Skill domain model with Category)
  - `internal/repository/career/skill_repository.go` (SkillRepository interface)
- [x] Confirmed this is ONE atomic task (CV generation redesign) ✅
- [x] Identified which test files will be created/modified (see Phase checklists)

## Context

Currently, CV generation uses "Role Emphasis" (Senior Backend, Staff/Principal, Consulting, Language-Agnostic) to determine presentation style. This task replaces that with a **Technology Focus** system that:

1. Derives technologies from user-defined skills (Task 39)
2. Allows selecting presentation style: Language Agnostic, Generalist (2-5 techs), or Specialist (1 tech)
3. Allows selecting Focus Area: Backend, Frontend, Fullstack, DevOps (derived from skill categories)
4. Populates CV Skills section with selected technologies

## Design Summary

### New Flow
```
[Select Profile] ← Contains TargetRole (Principal/Staff/EM/Senior IC)
    ↓
[Select Audience] ← Hiring Manager / Recruiter / Peer
    ↓
[Extracting Technologies...] (loading - aggregate user skills)
    ↓
[Select Technology Focus] ← REPLACES "Role Emphasis"
    ├─ Language Agnostic → [Select Focus Area] → [Select Length]
    ├─ Generalist (2-5) → [Select Technologies] → [Select Focus Area] → [Select Length]
    └─ Specialist (1) → [Select Technology] → [Select Focus Area] → [Select Length]
    ↓
[Select Length Format]
    ↓
[Generate CV]
```

### Variant ID Structure

**Language Agnostic:**
```
agnostic_{focus_area}_{length}
Example: agnostic_backend_standard
```

**Generalist:**
```
generalist_{focus_area}_{length}
Example: generalist_fullstack_short
```

**Specialist:**
```
specialist_{tech}_{focus_area}_{length}
Example: specialist_ruby_backend_full
```

**Note**: Target Role (Principal/Staff/EM/Senior IC) is NOT part of variant ID - it's a generation parameter from the Profile that affects bullet filtering.

### CV Structure Mapping

| Technology Focus | CV Structure |
|-----------------|--------------|
| Language Agnostic | Narrative |
| Generalist | Standard |
| Specialist | Standard |

Ultra-Short always uses Highlights structure regardless of Technology Focus.

## Files to Create

### Service Layer
- [ ] `internal/service/career/technology/extractor.go` - Technology extraction service
- [ ] `internal/service/career/technology/extractor_test.go` - Extraction tests
- [ ] `internal/service/career/technology/focus_area.go` - Focus area analyzer
- [ ] `internal/service/career/technology/focus_area_test.go` - Focus area tests

### Intent Tests
- [ ] `internal/cli/intents/generate_cv_technology_test.go` - Technology selection tests
- [ ] `internal/cli/intents/generate_cv_focus_area_test.go` - Focus area selection tests

## Files to Modify

### Domain/Service Layer
- [ ] `internal/service/career/cv/variants.go` - Remove old RoleEmphasis constants, add TechnologyFocus type, add FocusArea type
- [ ] `internal/service/career/cv/role_emphasis.go` - Remove old configs, add new technology focus configs
- [ ] `internal/service/career/cv/bullet_generator.go` - Add technology-based filtering

### Intent Layer
- [ ] `internal/cli/intents/generate_cv.go` - Add new states, state data fields, update types
- [ ] `internal/cli/intents/generate_cv_intent.go` - Add extraction, technology selection, focus area views/handlers
- [ ] `internal/cli/intents/generate_cv_structure_test.go` - Update role emphasis tests

### Documentation
- [ ] `docs/guides/CV_VARIANTS_GUIDE.md` - Complete rewrite
- [ ] `docs/guides/CV_GENERATION_GUIDE.md` - Significant updates

## Implementation Plan

### Phase 1: Technology Extraction Service

**Goal**: Create service to aggregate and filter user skills

#### Technology Extractor
```go
// internal/service/career/technology/extractor.go

type ExtractedTechnology struct {
    ID          string   // Skill ID
    Name        string   // "Ruby", "PostgreSQL"
    Category    string   // "backend", "database"
    EventCount  int      // How many events use this skill
    EventIDs    []string // Which events
}

type Extractor struct {
    skillRepo       career.SkillRepository
    eventRepository career.EventRepository
}

// ExtractFromUser aggregates user's defined skills with event associations
func (e *Extractor) ExtractFromUser(ctx context.Context) ([]*ExtractedTechnology, error)

// FilterByThreshold removes skills with < N events
func (e *Extractor) FilterByThreshold(techs []*ExtractedTechnology, minEvents int) []*ExtractedTechnology
```

#### Extraction Logic
1. Load all user-defined skills from SkillRepository
2. For each skill, count associated events from event_skills table
3. Filter out skills with < 3 events
4. Sort by event count (descending)
5. Return as ExtractedTechnology slice

#### Handling Events Without Skills

**Requirement from Task 39 completion**: Some events may not have associated skills (legacy events, or events where skills weren't applicable). The CV generation system must handle this gracefully.

**Approach**:
1. **Inclusion Strategy**: Events without skills are still included in CV generation
2. **Filtering Logic**: When technology focus is selected:
   - **Language Agnostic**: All events included (skills optional)
   - **Generalist (2-5 techs)**: Events with ANY of the selected technologies are prioritized
   - **Specialist (1 tech)**: Events with the selected technology are prioritized
3. **Prioritization**: Events WITH selected technologies score higher in bullet generation
4. **Fallback**: Events without skills can still appear if they're high-quality (strong bullets, recent dates)

**Implementation Details**:
- Bullet scoring includes skill match bonus (e.g., +0.15 if event has selected technology)
- Events without skills have baseline score (no bonus, no penalty)
- This ensures skills-based filtering is additive, not subtractive

**TDD Checklist - Phase 1:**
- [ ] Write failing test: ExtractFromUser loads user skills
- [ ] Test passes
- [ ] Write failing test: ExtractFromUser counts events per skill
- [ ] Test passes
- [ ] Write failing test: FilterByThreshold removes low-count skills
- [ ] Test passes
- [ ] Write failing test: Empty skills returns empty list
- [ ] Test passes
- [ ] Commit: `test(technology): add extractor tests`
- [ ] Commit: `feat(technology): implement technology extractor`

### Phase 2: Focus Area Analyzer

**Goal**: Suggest focus area based on skill categories

#### Focus Area Analyzer
```go
// internal/service/career/technology/focus_area.go

type FocusArea string

const (
    FocusAreaBackend   FocusArea = "backend"
    FocusAreaFrontend  FocusArea = "frontend"
    FocusAreaFullstack FocusArea = "fullstack"
    FocusAreaDevOps    FocusArea = "devops"
)

type FocusAreaSuggestion struct {
    Area       FocusArea
    Confidence float64       // 0.0-1.0
    Evidence   map[string]int // Category counts: {"backend": 12, "frontend": 3}
}

type Analyzer struct{}

// AnalyzeSkills suggests focus area from skill categories
func (a *Analyzer) AnalyzeSkills(techs []*ExtractedTechnology) *FocusAreaSuggestion
```

#### Analysis Logic
1. Count skills by category
2. Determine dominant focus area:
   - If backend > 70%: Backend
   - If frontend > 70%: Frontend
   - If devops > 70%: DevOps
   - If mix of backend + frontend: Fullstack
3. Calculate confidence based on distribution
4. Return suggestion with evidence

**TDD Checklist - Phase 2:**
- [ ] Write failing test: AnalyzeSkills suggests Backend (70%+ backend)
- [ ] Test passes
- [ ] Write failing test: AnalyzeSkills suggests Frontend (70%+ frontend)
- [ ] Test passes
- [ ] Write failing test: AnalyzeSkills suggests DevOps (70%+ devops)
- [ ] Test passes
- [ ] Write failing test: AnalyzeSkills suggests Fullstack (mixed)
- [ ] Test passes
- [ ] Write failing test: Confidence calculation
- [ ] Test passes
- [ ] Write failing test: Evidence map populated
- [ ] Test passes
- [ ] Commit: `test(technology): add focus area analyzer tests`
- [ ] Commit: `feat(technology): implement focus area analyzer`

### Phase 3: Update Domain Types

**Goal**: Replace old RoleEmphasis with TechnologyFocus

#### Remove Old Types
```go
// DELETE from variants.go:
RoleEmphasisSeniorBackend
RoleEmphasisStaffPrincipal
RoleEmphasisConsulting
```

#### Add New Types
```go
// internal/service/career/cv/variants.go

type TechnologyFocus string

const (
    TechnologyFocusLanguageAgnostic TechnologyFocus = "language_agnostic"
    TechnologyFocusGeneralist       TechnologyFocus = "generalist"
    TechnologyFocusSpecialist       TechnologyFocus = "specialist"
)

type FocusArea string

const (
    FocusAreaBackend   FocusArea = "backend"
    FocusAreaFrontend  FocusArea = "frontend"
    FocusAreaFullstack FocusArea = "fullstack"
    FocusAreaDevOps    FocusArea = "devops"
)
```

#### Update CVVariant
```go
type CVVariant struct {
    ID              string
    Name            string
    Description     string
    TechnologyFocus TechnologyFocus  // NEW (replaces RoleEmphasis)
    FocusArea       FocusArea        // NEW
    LengthFormat    LengthFormat
    Technologies    []string         // NEW - selected techs (for Generalist/Specialist)
    BaseStructure   CVStructure
    // ...
}
```

**TDD Checklist - Phase 3:**
- [ ] Write failing test: TechnologyFocus constants exist
- [ ] Test passes
- [ ] Write failing test: FocusArea constants exist
- [ ] Test passes
- [ ] Write failing test: CVVariant has TechnologyFocus field
- [ ] Test passes
- [ ] Write failing test: CVVariant has FocusArea field
- [ ] Test passes
- [ ] Write failing test: CVVariant has Technologies field
- [ ] Test passes
- [ ] Commit: `refactor(cv): replace RoleEmphasis with TechnologyFocus`

### Phase 4: Update State Machine

**Goal**: Add new states for technology and focus area selection

#### New States
```go
// internal/cli/intents/generate_cv.go

const (
    GenerateCVStateExtractingTechnologies GenerateCVState = "extracting_technologies"
    GenerateCVStateSelectTechnologyFocus  GenerateCVState = "select_technology_focus"  // Renamed from SelectRoleEmphasis
    GenerateCVStateSelectTechnologies     GenerateCVState = "select_technologies"      // NEW
    GenerateCVStateSelectFocusArea        GenerateCVState = "select_focus_area"        // NEW
    GenerateCVStateSelectLengthFormat     GenerateCVState = "select_length_format"
)
```

#### New State Data
```go
type GenerateCVModel struct {
    // ... existing fields ...
    
    // Technology extraction
    extractedTechnologies []*technology.ExtractedTechnology
    technologiesAvailable bool  // true if 3+ technologies found
    
    // Technology Focus selection
    selectedTechnologyFocus TechnologyFocus
    technologyFocusIndex    int
    
    // Technology selection (for Generalist/Specialist)
    selectedTechnologies []string      // Skill IDs
    technologyCursor     int
    technologySelected   map[int]bool  // Multi-select state
    
    // Focus area
    focusAreaSuggestion *technology.FocusAreaSuggestion
    selectedFocusArea   FocusArea
    focusAreaCursor     int
}
```

**TDD Checklist - Phase 4:**
- [ ] Write failing test: New states defined
- [ ] Test passes
- [ ] Write failing test: State data fields added
- [ ] Test passes
- [ ] Commit: `refactor(cv): add technology selection states`

### Phase 5: Technology Extraction Flow

**Goal**: Extract technologies after audience selection

#### Update Handler
```go
// Transition from SelectAudience to ExtractingTechnologies
func (i *GenerateCVIntent) updateSelectAudience(msg tea.Msg) tea.Cmd {
    case "enter":
        i.state.currentState = GenerateCVStateExtractingTechnologies
        return i.extractTechnologies()
}

// Extract technologies command
func (i *GenerateCVIntent) extractTechnologies() tea.Cmd {
    return func() tea.Msg {
        extractor := technology.NewExtractor(skillRepo, eventRepo)
        techs, err := extractor.ExtractFromUser(ctx)
        filtered := extractor.FilterByThreshold(techs, 3)
        
        analyzer := &technology.Analyzer{}
        suggestion := analyzer.AnalyzeSkills(filtered)
        
        return TechnologiesExtractedMsg{
            Technologies: filtered,
            Suggestion:   suggestion,
            Error:        err,
        }
    }
}
```

**TDD Checklist - Phase 5:**
- [ ] Write failing test: Audience selection triggers extraction
- [ ] Test passes
- [ ] Write failing test: Extraction command created
- [ ] Test passes
- [ ] Write failing test: TechnologiesExtractedMsg handled
- [ ] Test passes
- [ ] Write failing test: Extraction stores results in state
- [ ] Test passes
- [ ] Commit: `feat(cv): add technology extraction flow`

### Phase 6: Technology Focus Selection View

**Goal**: Replace Role Emphasis view with Technology Focus

#### View
```go
func (i *GenerateCVIntent) viewSelectTechnologyFocus() string {
    // Show 3 options: Language Agnostic, Generalist, Specialist
    // If < 3 technologies, disable Generalist and Specialist
    // Show tech count: "Found 8 technologies across your career events"
}
```

#### Update Handler
```go
func (i *GenerateCVIntent) updateSelectTechnologyFocus(msg tea.Msg) tea.Cmd {
    switch selected {
    case TechnologyFocusLanguageAgnostic:
        // Go directly to SelectFocusArea
    case TechnologyFocusGeneralist:
        // Go to SelectTechnologies (multi-select, 2-5)
    case TechnologyFocusSpecialist:
        // Go to SelectTechnologies (single-select)
    }
}
```

**TDD Checklist - Phase 6:**
- [ ] Write failing test: View shows 3 options
- [ ] Test passes
- [ ] Write failing test: View shows tech count
- [ ] Test passes
- [ ] Write failing test: Generalist/Specialist disabled if < 3 techs
- [ ] Test passes
- [ ] Write failing test: Language Agnostic transitions to FocusArea
- [ ] Test passes
- [ ] Write failing test: Generalist transitions to Technologies
- [ ] Test passes
- [ ] Write failing test: Specialist transitions to Technologies
- [ ] Test passes
- [ ] Commit: `feat(cv): implement technology focus selection`

### Phase 7: Technology Selection View

**Goal**: Allow selecting technologies (multi or single)

#### View (Generalist - Multi-Select)
```go
func (i *GenerateCVIntent) viewSelectTechnologies() string {
    // Show extracted technologies with checkboxes
    // For Generalist: "Select 2-5 technologies"
    // For Specialist: "Select 1 technology"
    // Show event count per technology
    // Space to toggle, Enter to confirm
}
```

#### Update Handler
```go
func (i *GenerateCVIntent) updateSelectTechnologies(msg tea.Msg) tea.Cmd {
    case "space":
        // Toggle selection
        // Validate count (2-5 for Generalist, exactly 1 for Specialist)
    case "enter":
        // Confirm and transition to SelectFocusArea
}
```

**TDD Checklist - Phase 7:**
- [ ] Write failing test: View shows technologies (Generalist)
- [ ] Test passes
- [ ] Write failing test: View shows technologies (Specialist)
- [ ] Test passes
- [ ] Write failing test: Space toggles selection
- [ ] Test passes
- [ ] Write failing test: Generalist requires 2-5 selections
- [ ] Test passes
- [ ] Write failing test: Specialist requires exactly 1
- [ ] Test passes
- [ ] Write failing test: Enter confirms and transitions
- [ ] Test passes
- [ ] Commit: `feat(cv): implement technology selection view`

### Phase 8: Focus Area Selection View

**Goal**: Allow selecting focus area with suggestions

#### View
```go
func (i *GenerateCVIntent) viewSelectFocusArea() string {
    // Show 4 options: Backend, Frontend, Fullstack, DevOps
    // Highlight suggested option
    // Show evidence: "backend: 12 skills, frontend: 3 skills"
}
```

#### Update Handler
```go
func (i *GenerateCVIntent) updateSelectFocusArea(msg tea.Msg) tea.Cmd {
    case "enter":
        // Save selection and transition to SelectLengthFormat
}
```

**TDD Checklist - Phase 8:**
- [ ] Write failing test: View shows 4 focus areas
- [ ] Test passes
- [ ] Write failing test: Suggestion highlighted
- [ ] Test passes
- [ ] Write failing test: Evidence displayed
- [ ] Test passes
- [ ] Write failing test: Selection saved
- [ ] Test passes
- [ ] Write failing test: Transition to length format
- [ ] Test passes
- [ ] Commit: `feat(cv): implement focus area selection view`

### Phase 9: Update Variant System

**Goal**: Create new variants, remove old ones

#### Remove Old Variants
- Delete all 16 old variants (senior_backend_*, staff_principal_*, consulting_*, language_agnostic_*)

#### Create New Variant Templates
```go
// Variant lookup now dynamic based on selections
func GetVariantBySelections(
    techFocus TechnologyFocus,
    focusArea FocusArea,
    length LengthFormat,
    technologies []string,
) (*CVVariant, error) {
    
    variantID := buildVariantID(techFocus, focusArea, length, technologies)
    
    // Create variant dynamically
    return &CVVariant{
        ID:              variantID,
        TechnologyFocus: techFocus,
        FocusArea:       focusArea,
        LengthFormat:    length,
        Technologies:    technologies,
        BaseStructure:   determineStructure(techFocus, length),
        BulletConfig:    buildBulletConfig(length),
    }, nil
}

func determineStructure(techFocus TechnologyFocus, length LengthFormat) CVStructure {
    if length == LengthUltraShort {
        return CVStructureHighlights
    }
    
    switch techFocus {
    case TechnologyFocusLanguageAgnostic:
        return CVStructureNarrative
    case TechnologyFocusGeneralist, TechnologyFocusSpecialist:
        return CVStructureStandard
    }
}
```

**TDD Checklist - Phase 9:**
- [ ] Write failing test: Old variants removed
- [ ] Test passes
- [ ] Write failing test: GetVariantBySelections creates Language Agnostic variant
- [ ] Test passes
- [ ] Write failing test: GetVariantBySelections creates Generalist variant
- [ ] Test passes
- [ ] Write failing test: GetVariantBySelections creates Specialist variant
- [ ] Test passes
- [ ] Write failing test: Variant ID generation
- [ ] Test passes
- [ ] Write failing test: Structure determination
- [ ] Test passes
- [ ] Commit: `refactor(cv): replace static variants with dynamic system`

### Phase 10: Update Bullet Generation

**Goal**: Filter/prioritize bullets based on selected technologies

#### Technology-Based Filtering
```go
// internal/service/career/cv/bullet_generator.go

func (g *BulletGenerator) FilterByTechnologies(
    bullets []*Bullet,
    techFocus TechnologyFocus,
    technologies []string,
) []*Bullet {
    
    switch techFocus {
    case TechnologyFocusLanguageAgnostic:
        // No filtering - show all (skills optional)
        return bullets
        
    case TechnologyFocusSpecialist:
        // Prioritize bullets with selected tech, keep high-quality bullets without skills
        return filterAndBoostWithTech(bullets, technologies[0])
        
    case TechnologyFocusGeneralist:
        // Boost bullets with selected techs, keep others with baseline score
        return boostBulletsWithTechs(bullets, technologies)
    }
}

// Skill Matching Bonus System
// +0.15 bonus if event has selected technology
// +0.00 baseline if event has no skills (not penalized)
// This makes skills additive, not subtractive
```

#### Handling Events Without Skills in Bullet Generation

**Scoring Logic**:
1. **Base Score**: All bullets start with base score from EnhancedBulletGenerator (0.0-1.0)
2. **Skill Match Bonus**: +0.15 if event has selected technology
3. **No Penalty**: Events without skills keep base score (no deduction)
4. **High-Quality Fallback**: Events without skills can still rank high if they have:
   - Strong action verbs (0.70+ base confidence)
   - Recent dates (within last 2 years)
   - Clear impact/results

**Example Scoring**:
```
Bullet A (with selected tech): 0.85 base + 0.15 skill bonus = 1.00 (capped)
Bullet B (no skills): 0.85 base + 0.00 = 0.85
Bullet C (with different tech): 0.75 base + 0.00 = 0.75
Bullet D (no skills, strong): 0.90 base + 0.00 = 0.90 (beats C!)
```

**Result**: Skills provide advantage, but quality bullets without skills still appear.

**TDD Checklist - Phase 10:**
- [ ] Write failing test: Language Agnostic shows all bullets (including events without skills)
- [ ] Test passes
- [ ] Write failing test: Events without skills are not penalized (baseline score)
- [ ] Test passes
- [ ] Write failing test: Events with selected technology get skill match bonus
- [ ] Test passes
- [ ] Write failing test: High-quality events without skills can rank higher than low-quality events with skills
- [ ] Test passes
- [ ] Write failing test: Specialist filters to selected tech
- [ ] Test passes
- [ ] Write failing test: Generalist boosts selected techs
- [ ] Test passes
- [ ] Commit: `feat(cv): add technology-based bullet filtering`

### Phase 11: Skills Section Population

**Goal**: Populate Skills section with selected technologies

#### Update Skills Section Builder
```go
// internal/service/career/cv/section_builder.go

func (b *SectionBuilder) buildSkillsSection(
    variant *CVVariant,
    userSkills []*career.Skill,
) *CVSection {
    
    // For Generalist/Specialist: Prioritize selected technologies
    // Show selected technologies first, then others
    // Group by category
    // Include event counts
}
```

**TDD Checklist - Phase 11:**
- [ ] Write failing test: Skills section shows selected technologies
- [ ] Test passes
- [ ] Write failing test: Selected techs appear first
- [ ] Test passes
- [ ] Write failing test: Skills grouped by category
- [ ] Test passes
- [ ] Write failing test: Event counts displayed
- [ ] Test passes
- [ ] Commit: `feat(cv): populate skills section with selected technologies`

### Phase 12: Update Documentation

**Goal**: Rewrite variant and generation guides

#### CV_VARIANTS_GUIDE.md
- Remove all references to old role emphasis
- Document new Technology Focus options
- Document Focus Area selection
- Document variant naming scheme
- Provide examples for each combination

#### CV_GENERATION_GUIDE.md
- Update selection flow diagrams
- Document technology extraction
- Document focus area suggestions
- Update examples

**TDD Checklist - Phase 12:**
- [ ] CV_VARIANTS_GUIDE.md rewritten
- [ ] CV_GENERATION_GUIDE.md updated
- [ ] Examples added for each variant type
- [ ] Commit: `docs(cv): update variant and generation guides`

## Pre-Commit Checklist (BEFORE EACH COMMIT)
- [ ] `make check-compliance` passes (REQUIRED before commit)
- [ ] Use `make ai-commit MSG="type(scope): description"` for AI-generated code
- [ ] Commit message explains **WHY**, not just WHAT
- [ ] Commit is atomic (ONE logical change)

## Post-Task Checklist (MUST COMPLETE BEFORE NEXT TASK)
- [ ] `make check-compliance` passes
- [ ] All checkboxes above completed
- [ ] Task marked complete `[x]` in task file
- [ ] Token count: _____ (< 100k to continue)

## Acceptance Criteria
- [ ] Users can select Technology Focus (Language Agnostic / Generalist / Specialist)
- [ ] Users can select technologies (multi-select for Generalist, single for Specialist)
- [ ] Users can select Focus Area (Backend/Frontend/Fullstack/DevOps) with suggestions
- [ ] Technology extraction works with 3+ event threshold
- [ ] Focus area is suggested based on skill categories
- [ ] Old role emphasis variants removed
- [ ] New variant system works with all combinations
- [ ] CV Skills section populated with selected technologies
- [ ] Bullet filtering works based on technology selection
- [ ] All tests pass (100% pass rate)
- [ ] Coverage maintained ≥ 80%
- [ ] Zero staticcheck warnings
- [ ] Zero race conditions
- [ ] Documentation updated

## Rollback Plan
- This is a breaking change - old variant IDs will not work
- Users will need to regenerate CVs with new flow
- Database is not affected (no schema changes)
- Can revert code changes via git if needed

## Breaking Changes

### For Users
- Old CV configs referencing `senior_backend`, `staff_principal`, `consulting` variants will not work
- Must regenerate CVs using new flow

### For Developers
- `RoleEmphasis` type removed
- `RoleEmphasisConfig` removed
- All 16 old variants removed
- New `TechnologyFocus` and `FocusArea` types added

## Notes

### Variant ID Examples

**Language Agnostic:**
- `agnostic_backend_full`
- `agnostic_frontend_standard`
- `agnostic_fullstack_short`
- `agnostic_devops_ultra_short`

**Generalist:**
- `generalist_backend_standard`
- `generalist_fullstack_short`

**Specialist:**
- `specialist_ruby_backend_full`
- `specialist_react_frontend_standard`
- `specialist_kubernetes_devops_short`

### Technology Threshold
- Minimum 3 events with a skill to appear in selection
- Prevents noise from rarely-used skills
- Users can always add more skills in Manage Skills (Task 39)

### Focus Area Suggestions
- Backend: 70%+ backend category skills
- Frontend: 70%+ frontend category skills
- DevOps: 70%+ devops category skills
- Fullstack: Mix of backend + frontend
- User always makes final choice

### Target Role Independence
- Target Role (Principal/Staff/EM/Senior IC) remains in Profile
- Affects bullet caps and filtering (per PRD)
- Not part of variant identity
- Orthogonal to Technology Focus

### CV Structure Mapping
- Language Agnostic → Narrative (emphasizes adaptability)
- Generalist → Standard (traditional format)
- Specialist → Standard (traditional format)
- Ultra-Short → Highlights (always, regardless of focus)

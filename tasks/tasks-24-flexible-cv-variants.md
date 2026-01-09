# Task 24: Flexible CV Variants

## Overview

- **Goal**: Implement a flexible CV variant system with two-dimensional selection (role emphasis x length), adding two new structures (Consulting, Highlights) and replacing the current Standard/Narrative structure selection
- **Time Estimate**: ~3-4 weeks (7 phases)
- **Prerequisites**: 
  - Task 23 complete (CV structure selection) ✅
  - Understanding of PRD sections 5-6 (fact selection, role/audience filtering)

**Related Documents**:
- `docs/PRD_MASTER.md` - Master PRD with role/audience filtering specs
- `docs/guides/CV_GENERATION_GUIDE.md` - CV generation workflow
- `docs/CUSTOM_CV_FORMAT_PROPOSAL.md` - Current structure implementation

---

## Key Design Decisions

| Decision | Value | Rationale |
|----------|-------|-----------|
| **Storage** | Built-in variants in Go code | User variants deferred to future |
| **UI Selection** | Two dropdowns: Role Emphasis → Length | Clear mental model |
| **Structure per variant** | Implicit (variant determines structure) | Simplifies flow |
| **Structures** | 4 total: standard, narrative, consulting, highlights | Covers all use cases |
| **Audience filtering** | Use existing `Fact.AudienceRelevance` field | Infrastructure exists |
| **Audience affects bullets** | Both number AND type | Per PRD documentation |
| **Date filtering** | Hard filter (events excluded) | Clean, predictable |
| **Density thresholds** | Simple: "Show if data exists, skip if not" | Pragmatic approach |
| **User notification** | Show in preview only (not in export) | Clean exports |
| **Positioning Statement** | Optional - use if set, omit if empty | Flexible |
| **Client Engagements** | Group by company (defer client/employer distinction) | Simplify initial scope |

---

## Structures (4 Total)

| Structure | Sections | Used By |
|-----------|----------|---------|
| **standard** | Summary, Experience, Projects, Skills | senior_backend, staff_principal |
| **narrative** | Profile Header, Positioning Statement*, Summary, Core Strengths, Technologies, Selected Experience, What I Bring | language_agnostic |
| **consulting** | Profile Header, Summary, Client Engagements (by company), Technical Capabilities*, What I Bring | consulting |
| **highlights** | Profile Header, Summary, Key Capabilities, Selected Highlights, Languages & Systems | all ultra_short variants |

*\* = Optional section, shown if data exists, omitted if not*

---

## Variants (16 Total)

### Variant Matrix

| Role Emphasis | Full | Standard | Short | Ultra-Short |
|---------------|------|----------|-------|-------------|
| **senior_backend** | standard | standard | standard | highlights |
| **staff_principal** | standard | standard | standard | highlights |
| **consulting** | consulting | consulting | consulting | highlights |
| **language_agnostic** | narrative | narrative | narrative | highlights |

### Variant Configuration Details

| ID | Name | Role Emphasis | Length | MinConfidence | MaxYears | MaxCompanies | Structure |
|----|------|---------------|--------|---------------|----------|--------------|-----------|
| `senior_backend_full` | Senior Backend (Full) | senior_backend | full | 0.60 | nil | nil | standard |
| `senior_backend_standard` | Senior Backend (Standard) | senior_backend | standard | 0.65 | 10 | nil | standard |
| `senior_backend_short` | Senior Backend (Short) | senior_backend | short | 0.75 | 5 | nil | standard |
| `senior_backend_ultra_short` | Senior Backend (1-Page) | senior_backend | ultra_short | 0.85 | nil | 3 | highlights |
| `staff_principal_full` | Staff/Principal (Full) | staff_principal | full | 0.70 | nil | nil | standard |
| `staff_principal_standard` | Staff/Principal (Standard) | staff_principal | standard | 0.75 | 10 | nil | standard |
| `staff_principal_short` | Staff/Principal (Short) | staff_principal | short | 0.80 | 5 | nil | standard |
| `staff_principal_ultra_short` | Staff/Principal (1-Page) | staff_principal | ultra_short | 0.90 | nil | 3 | highlights |
| `consulting_full` | Consulting (Full) | consulting | full | 0.50 | nil | nil | consulting |
| `consulting_standard` | Consulting (Standard) | consulting | standard | 0.60 | 10 | nil | consulting |
| `consulting_short` | Consulting (Short) | consulting | short | 0.70 | 5 | nil | consulting |
| `consulting_ultra_short` | Consulting (1-Page) | consulting | ultra_short | 0.80 | nil | 3 | highlights |
| `language_agnostic_full` | Language-Agnostic (Full) | language_agnostic | full | 0.65 | nil | nil | narrative |
| `language_agnostic_standard` | Language-Agnostic (Standard) | language_agnostic | standard | 0.70 | 10 | nil | narrative |
| `language_agnostic_short` | Language-Agnostic (Short) | language_agnostic | short | 0.78 | 5 | nil | narrative |
| `language_agnostic_ultra_short` | Language-Agnostic (1-Page) | language_agnostic | ultra_short | 0.88 | nil | 3 | highlights |

---

## Implementation Phases

### Phase 0: Implement Audience Filtering (Prerequisite) ✅ COMPLETE

**Objective**: Implement the audience filtering documented in PRD but currently stubbed out.

**Status**: Completed 2026-01-09

**Commits**:
- `c2b21e7` - feat(cv): implement audience filtering for bullet generation
- `eed48a8` - feat(cv): add audience filtering to EnhancedBulletGenerator

#### Step 0.1: Write Failing Tests ✅
**File**: `internal/service/career/cv/bullet_generator_test.go`

Tests added (9 new tests):
- [x] Filters facts by hiring_manager audience (returns facts with hiring_manager in AudienceRelevance)
- [x] Filters facts by recruiter audience
- [x] Filters facts by peer audience
- [x] Returns all facts when audience is empty string
- [x] Excludes facts not relevant to selected audience
- [x] Multiple audience relevance (facts with multiple audiences)
- [x] Legacy event-based generation tests (3 tests)

#### Step 0.2: Implement Fact Audience Filtering ✅
**File**: `internal/service/career/cv/bullet_generator.go`

- [x] Update `isFactRelevantToAudience()` to check `fact.AudienceRelevance` contains requested audience
- [x] Fix `GenerateBullets` to process facts even when events slice is empty
- [x] Update log message to include facts count

#### Step 0.3: Implement Enhanced Bullet Generator Filtering ✅
**File**: `internal/service/career/cv/enhanced_bullet_generator.go`

- [x] Add `AudienceRelevance` field to `EnhancedBullet` struct
- [x] Update `createBulletsFromFacts()` to filter by audience during creation
- [x] Implement `isFactRelevantToAudience()` helper
- [x] Update `FilterByAudience()` to filter bullets by stored audience relevance
- [x] Enhance `calculateAudienceScore()` to give bonus for matching audience
- [x] Add 6 comprehensive tests for audience filtering scenarios

#### Verification ✅
- [x] All existing tests pass (215/216 - 1 clipboard test requires display)
- [x] New audience filtering tests pass (15 new tests total)
- [x] `go test -race ./internal/service/career/cv/...` passes

---

### Phase 1: Domain Models and Built-In Variants ✅ COMPLETE

**Objective**: Define the type system, 4 structures, and 16 built-in variants in code.

**Status**: Completed 2026-01-09

**Commit**: `f2c0ffa` - feat(cv): add CV variant types and 16 built-in variants

#### Step 1.1: Create Variant Types ✅
**File**: `internal/service/career/cv/variants.go` (new)

```go
// RoleEmphasis defines what aspect of experience to emphasize
type RoleEmphasis string

const (
    RoleEmphasisSeniorBackend    RoleEmphasis = "senior_backend"
    RoleEmphasisStaffPrincipal   RoleEmphasis = "staff_principal"
    RoleEmphasisConsulting       RoleEmphasis = "consulting"
    RoleEmphasisLanguageAgnostic RoleEmphasis = "language_agnostic"
)

// LengthFormat defines CV density/length
type LengthFormat string

const (
    LengthFull       LengthFormat = "full"
    LengthStandard   LengthFormat = "standard"
    LengthShort      LengthFormat = "short"
    LengthUltraShort LengthFormat = "ultra_short"
)

// CVVariant combines dimensions with configuration
type CVVariant struct {
    ID              string
    Name            string           // Friendly display name
    Description     string
    RoleEmphasis    RoleEmphasis
    LengthFormat    LengthFormat
    BaseStructure   CVStructure      // standard, narrative, consulting, highlights
    Sections        []SectionConfig
    BulletConfig    BulletConfig
    ProfileOverride *ProfileOverride // Phase 6
    IsBuiltIn       bool
}

type SectionConfig struct {
    SectionType   string   // experience, skills, core_strengths, etc.
    Title         string   // Display title (customizable)
    Order         int      // Display order
    Enabled       bool     // Show/hide
    MinConfidence float64  // Filter threshold for this section
}

type BulletConfig struct {
    MaxBulletsPerCompany *int     // nil = use role default
    MinConfidence        *float64 // nil = use role default
    MaxYearsHistory      *int     // nil = no limit
    MaxCompanies         *int     // nil = no limit
}

type ProfileOverride struct {
    ProfessionalTitle     *string   // e.g., "Senior Consulting Engineer"
    CoreStrengths         []string  // 6 key competencies
    CareerDifferentiators []string  // Unique value propositions
    CareerPositioning     *string   // Career positioning narrative
}
```

#### Step 1.2: Add New Structure Constants ✅
**File**: `internal/service/career/cv/cv_helpers.go`

- [x] Add `CVStructureConsulting CVStructure = "consulting"`
- [x] Add `CVStructureHighlights CVStructure = "highlights"`

#### Step 1.3: Define Built-In Variants ✅
**File**: `internal/service/career/cv/variants.go`

- [x] Define all 16 variants as `var BuiltInVariants []*CVVariant`
- [x] Include bullet configurations per length format (MinConfidence, MaxYearsHistory, MaxCompanies)
- [x] All variants marked as built-in

#### Step 1.4: Create Variant Service ✅
**File**: `internal/service/career/cv/variants.go` (combined with types)

- [x] `VariantService` interface with all methods
- [x] `DefaultVariantService` implementation
- [x] `NewVariantService()` constructor
- [x] Index variants by ID and by dimensions (role:length)

#### Step 1.5: Write Tests ✅
**File**: `internal/service/career/cv/variants_test.go` (new)

- [x] Lists all 16 built-in variants
- [x] Gets variant by ID
- [x] Returns error for unknown variant ID
- [x] Gets variant by dimensions (role + length)
- [x] Lists 4 role emphases with metadata
- [x] Lists 4 length formats with metadata
- [x] All variants have valid configuration
- [x] All variants map to valid structure
- [x] Tests for each role emphasis (4 variants each)
- [x] 31 total test specs

#### Verification ✅
- [x] `go build ./...` compiles
- [x] All 31 variant tests pass
- [x] No regressions (246/247 CV tests pass - 1 clipboard test requires display)

---

### Phase 2: New Structures (Consulting, Highlights) ✅ COMPLETE

**Objective**: Implement the two new CV structures and their section rendering.

**Status**: Completed 2026-01-09

**Commit**: `b8074e4` - feat(cv): implement Consulting and Highlights export structures

#### Step 2.1: Consulting Structure Export ✅
**File**: `internal/service/career/cv/export_service.go`

Consulting structure sections:
1. Profile Header (role, location, contact)
2. Summary
3. Client Engagements (experience grouped by company with dates)
4. What I Bring (value propositions from profile)

- [x] Add `exportConsultingWithProfile()` routing method
- [x] Add `exportConsultingText()` for plain text
- [x] Add `exportConsultingMarkdown()` for markdown

#### Step 2.2: Highlights Structure Export ✅
**File**: `internal/service/career/cv/export_service.go`

Highlights structure sections:
1. Condensed Profile Header (one line)
2. Short Summary
3. Key Capabilities (4-6 from core_strengths or defaults)
4. Selected Highlights (top 5 bullets by confidence)
5. Technologies (languages and systems)

- [x] Add `exportHighlightsWithProfile()` routing method
- [x] Add `exportHighlightsText()` for plain text
- [x] Add `exportHighlightsMarkdown()` for markdown
- [x] Add `getTopBulletsByConfidence()` helper

#### Step 2.3: Update Export Routing ✅
**File**: `internal/service/career/cv/export_service.go`

- [x] Update `ExportWithProfile()` switch to handle all 4 structures
- [x] CVStructureConsulting routes to consulting export
- [x] CVStructureHighlights routes to highlights export

#### Step 2.4: Tests ✅
**File**: `internal/service/career/cv/structure_test.go` (new)

24 tests covering:
- [x] Consulting: profile header, summary, client engagements, bullets, what I bring
- [x] Consulting: text and markdown formats
- [x] Highlights: condensed header, summary, key capabilities, selected highlights
- [x] Highlights: top 5 by confidence filtering, technologies section
- [x] Highlights: text and markdown formats
- [x] Structure routing for all 4 types

#### Verification ✅
- [x] All 24 structure tests pass
- [x] No regressions (270/271 CV tests pass - 1 clipboard test requires display)

---

### Phase 3: Role Emphasis Dimension ✅ COMPLETE

**Objective**: Implement role emphasis-specific bullet filtering and section configuration.

**Status**: Completed 2026-01-09

**Commit**: `c5dea6b` - feat(cv): add role emphasis configuration and scoring

#### Step 3.1: Define Role Emphasis Behaviors ✅
**File**: `internal/service/career/cv/role_emphasis.go` (new)

| Role Emphasis | Primary Categories | Secondary Categories | Bullet Focus |
|---------------|-------------------|---------------------|--------------|
| senior_backend | technical, architecture | product, delivery | Technical depth, product impact |
| staff_principal | leadership, strategy, architecture | technical, mentoring | Architecture, mentorship, cross-team |
| consulting | strategy, delivery, consulting | technical, leadership | Client work, rapid assessment |
| language_agnostic | technical, architecture | all | Multi-language evidence, adaptability |

- [x] Define `RoleEmphasisConfig` struct
- [x] Define `GetRoleEmphasisConfig()` function
- [x] Define `ListRoleEmphasisConfigs()` for TUI
- [x] Add `ScoreBulletCategory()` for category-based scoring

#### Step 3.2: Write Tests ✅
**File**: `internal/service/career/cv/role_emphasis_test.go` (new)

- [x] Returns config for each role emphasis type
- [x] Returns default config for unknown emphasis
- [x] Lists all 4 role emphasis configs
- [x] ScoreBulletCategory returns 1.0 for primary categories
- [x] ScoreBulletCategory returns 0.6 for secondary categories
- [x] ScoreBulletCategory returns 0.3 for unrelated categories
- [x] 11 total test specs

#### Verification ✅
- [x] All role emphasis tests pass
- [x] No regressions in existing tests

---

### Phase 4: Length Dimension ✅ COMPLETE

**Objective**: Implement length-based compression with hard date filtering.

**Status**: Completed 2026-01-09

**Commit**: `729d265` - feat(cv): add length format configuration and filtering

#### Step 4.1: Define Length Behaviors ✅
**File**: `internal/service/career/cv/length_format.go` (new)

| Length | Max Years | Max Companies | Min Confidence | Target Pages |
|--------|-----------|---------------|----------------|--------------|
| full | unlimited | unlimited | 0.50 | 3+ |
| standard | 10 | unlimited | 0.65 | 2-3 |
| short | 5 | 5 | 0.75 | 1-2 |
| ultra_short | 3 | 3 | 0.85 | 1 |

- [x] Define `LengthFormatConfig` struct (MaxYearsHistory, MaxCompanies, MaxBulletsPerJob, MinConfidence, etc.)
- [x] Define `GetLengthFormatConfig()` function
- [x] Define `ListLengthFormatConfigs()` for TUI
- [x] Add `ShouldIncludeEvent()` for date-based filtering
- [x] Add `ShouldIncludeEventByYear()` for year-based filtering
- [x] Add `FilterCompaniesByLimit()` for company limiting
- [x] Add `GetEffectiveBulletLimit()` for bullet limits
- [x] Add `MeetsConfidenceThreshold()` for confidence filtering

#### Step 4.2: Write Tests ✅
**File**: `internal/service/career/cv/length_format_test.go` (new)

- [x] Returns config for each length format type
- [x] Returns default config for unknown format
- [x] Lists all 4 length format configs
- [x] Full format includes all events (no year limit)
- [x] Standard format excludes events older than 10 years
- [x] Short format excludes events older than 5 years
- [x] Ultra-short excludes events older than 3 years
- [x] FilterCompaniesByLimit limits to MaxCompanies
- [x] GetEffectiveBulletLimit returns configured or default limit
- [x] MeetsConfidenceThreshold filters by MinConfidence
- [x] Progressively stricter confidence thresholds
- [x] Progressively stricter year limits
- [x] Progressively stricter bullet limits
- [x] 31 total test specs

#### Verification ✅
- [x] All length tests pass
- [x] No regressions in existing tests

---

### Phase 5: UI Integration ✅ COMPLETE

**Objective**: Replace structure selection with two-dropdown variant selection.

**Status**: Completed 2026-01-09

**Commit**: `d43b7ca` - feat(tui): add variant-based CV generation workflow

#### Step 5.1: Update GenerateCV Types ✅
**File**: `internal/cli/intents/generate_cv.go`

- [x] Add `GenerateCVStateSelectRoleEmphasis GenerateCVState = "select_role_emphasis"`
- [x] Add `GenerateCVStateSelectLengthFormat GenerateCVState = "select_length_format"`
- [x] Keep `GenerateCVStateSelectStructure` (deprecated but backward compatible)
- [x] Add `selectedRoleEmphasis RoleEmphasis` to model
- [x] Add `selectedLengthFormat LengthFormat` to model
- [x] Add `selectedVariant *CVVariant` to model
- [x] Update `GenerateCVResult` to include `SelectedVariant`
- [x] Add type aliases `RoleEmphasis` and `LengthFormat` for convenience

#### Step 5.2: Implement Role Emphasis Selection ✅
**File**: `internal/cli/intents/generate_cv_intent.go`

- [x] Add `updateSelectRoleEmphasis()` handler
- [x] Add `viewSelectRoleEmphasis()` view
- [x] Show 4 role emphases with descriptions from RoleEmphasisConfig
- [x] Navigate with j/k or arrows, select with Enter
- [x] Esc goes back to audience selection
- [x] q/ctrl+c/m cancels

#### Step 5.3: Implement Length Selection ✅
**File**: `internal/cli/intents/generate_cv_intent.go`

- [x] Add `updateSelectLengthFormat()` handler
- [x] Add `viewSelectLengthFormat()` view
- [x] Show 4 length formats with page counts from LengthFormatConfig
- [x] Default to Standard (index 1)
- [x] Navigate with j/k or arrows, select with Enter
- [x] Esc goes back to role emphasis selection
- [x] On Enter, resolve variant and proceed to generation

#### Step 5.4: Update State Flow ✅
**File**: `internal/cli/intents/generate_cv_intent.go`

New flow:
```
SelectProfile → SelectAudience → SelectRoleEmphasis → SelectLengthFormat → Generating
```

- [x] Update `updateSelectAudience()` to transition to `SelectRoleEmphasis`
- [x] Update breadcrumbs: "Select Role Emphasis", "Select Length"
- [x] Update context help for new states
- [x] Update view routing in `getStateContent()`

#### Step 5.5: Wire Variant to Generation ✅
**File**: `internal/cli/intents/generate_cv_intent.go`

- [x] Lookup variant via `VariantService.GetVariantByDimensions()`
- [x] Set `selectedVariant` on model
- [x] Set `selectedCVStructure` from variant's `BaseStructure`
- [x] Include variant in result metadata (role_emphasis, length_format, variant_id)

#### Step 5.6: Update Tests ✅
**Files**: Multiple test files updated

- [x] `generate_cv_test.go`: audience → role emphasis transition
- [x] `generate_cv_structure_test.go`: renamed to variant-based tests
- [x] `generate_cv_workflow_test.go`: full workflow tests
- [x] `generate_cv_preview_test.go`: preview tests for variants

#### Verification ✅
- [x] All 675 intent tests pass
- [x] No regressions in existing tests

---

### Phase 6: ProfileOverride Support ✅ COMPLETE

**Objective**: Allow variant-specific profile customization via ProfileOverride.

**Status**: Completed 2026-01-09

**Commit**: `b518e33` - feat(cv): wire ProfileOverride from variants to export

#### Step 6.1: Implement Profile Override Functions ✅
**File**: `internal/service/career/cv/cv_helpers.go`

- [x] Add `ApplyProfileOverride(cfg *ProfileConfig, override *ProfileOverride) *ProfileConfig`
  - Creates copy of config to avoid mutation
  - Applies ProfessionalTitle override (maps to Title)
  - Applies CoreStrengths override
  - Applies CareerDifferentiators override (maps to WhatIBring)
- [x] Add `ApplyProfileOverrideToNarrative(profile *NarrativeProfileData, override *ProfileOverride) *NarrativeProfileData`
  - Direct override for narrative profile data
  - Creates copy to avoid mutation

#### Step 6.2: Wire ProfileOverride to Export ✅
**File**: `internal/cli/intents/generate_cv_intent.go`

- [x] Update `exportCVAsync()` to apply ProfileOverride from selected variant
- [x] Pass modified profileCfg to ExportWithProfile()

#### Step 6.3: Write Tests ✅
**File**: `internal/service/career/cv/cv_helpers_test.go` (new)

ApplyProfileOverride tests:
- [x] Returns original config when override is nil
- [x] Creates config with override values when config is nil
- [x] Overrides title when ProfessionalTitle is set
- [x] Overrides core strengths when CoreStrengths is set
- [x] Overrides WhatIBring when CareerDifferentiators is set
- [x] Applies all overrides when multiple fields set
- [x] Does not modify the original config

ApplyProfileOverrideToNarrative tests:
- [x] Returns original profile when override is nil
- [x] Returns nil when profile is nil
- [x] Overrides role when ProfessionalTitle is set
- [x] Overrides core strengths
- [x] Overrides value propositions
- [x] Does not modify the original profile
- [x] 13 total test specs

#### Verification ✅
- [x] All helper tests pass
- [x] 325/326 CV tests pass (1 clipboard test requires display)
- [x] No regressions

---

### Phase 7: Section Configuration (DEFERRED)

**Objective**: Enable section configuration per variant (order, enable/disable, custom titles).

**Status**: Deferred to future work - Phase 6 (ProfileOverride) completed instead.

**Note**: ProfileOverride functionality was prioritized as it provides immediate value for customizing CV exports based on role emphasis. Section configuration (order, enable/disable, custom titles) will be implemented in a future task when more complex variant customization is needed.

#### Future Work
- Update `BuildSections()` to accept `[]SectionConfig`
- Order sections by `SectionConfig.Order`
- Skip sections where `Enabled: false`
- Use `SectionConfig.Title` for custom titles
- Apply `SectionConfig.MinConfidence` for per-section filtering

---

## Files Summary

### Files to Create

| File | Purpose |
|------|---------|
| `internal/service/career/cv/variants.go` | Variant types, built-in definitions, VariantService |
| `internal/service/career/cv/variants_test.go` | Variant tests |
| `internal/service/career/cv/variant_service.go` | VariantService interface and implementation |
| `internal/service/career/cv/role_emphasis.go` | Role emphasis behaviors and config |
| `internal/service/career/cv/role_emphasis_test.go` | Role emphasis tests |
| `internal/service/career/cv/length_format.go` | Length format behaviors and config |
| `internal/service/career/cv/length_format_test.go` | Length format tests |
| `internal/service/career/cv/structure_test.go` | Consulting and Highlights structure tests |
| `internal/service/career/cv/section_config_test.go` | Section configuration tests |
| `internal/service/career/cv/profile_override_test.go` | Profile override tests |
| `internal/cli/intents/generate_cv_variant_test.go` | UI variant selection tests |

### Files to Modify

| File | Changes |
|------|---------|
| `internal/service/career/cv/cv_helpers.go` | Add new structure constants, profile merging |
| `internal/service/career/cv/bullet_generator.go` | Implement audience filtering |
| `internal/service/career/cv/enhanced_bullet_generator.go` | Role emphasis filtering, audience filtering |
| `internal/service/career/cv/cv_generation_service.go` | Date/company filtering, variant integration |
| `internal/service/career/cv/section_builder.go` | New structures, section config rendering |
| `internal/service/career/cv/export_service.go` | New structure export, section config, profile override |
| `internal/cli/intents/generate_cv.go` | New states, variant fields in model and result |
| `internal/cli/intents/generate_cv_intent.go` | New handlers, views, state flow changes |

---

## Test Count (Actual)

| Phase | Tests | Status |
|-------|-------|--------|
| Phase 0 (Audience) | 15 | ✅ |
| Phase 1 (Variants) | 31 | ✅ |
| Phase 2 (Structures) | 24 | ✅ |
| Phase 3 (Role Emphasis) | 11 | ✅ |
| Phase 4 (Length) | 31 | ✅ |
| Phase 5 (UI) | ~40 (updated existing tests) | ✅ |
| Phase 6 (ProfileOverride) | 13 | ✅ |
| Phase 7 (Section Config) | - | Deferred |
| **Total New Tests** | ~125 | ✅ |

---

## Migration Notes

### Breaking Changes

1. **State flow changed**: `SelectStructure` state removed, replaced by `SelectRoleEmphasis` + `SelectLength`
2. **Result changed**: `SelectedStructure CVStructure` replaced by `SelectedVariant *CVVariant`
3. **New structures**: `consulting` and `highlights` added to `CVStructure` type

### Backward Compatibility

- Existing CV configs continue to work (they specify role + audience, variant is additive)
- Standard and Narrative structures still exist (as base structures for variants)
- Existing exports continue to work

---

## Rollback Plan

Each phase can be rolled back independently:

1. **Phase 0**: Revert bullet_generator.go changes (return to stub)
2. **Phase 1**: Delete new variant files
3. **Phase 2**: Revert section_builder.go and export_service.go structure additions
4. **Phase 3-4**: Revert filter changes
5. **Phase 5**: Revert intent changes, restore SelectStructure state
6. **Phase 6-7**: Revert section/profile changes

---

## Deferred / Future Work

- **User-created custom variants** stored in `~/.kariya/cv_variants/`
- **Client/employer distinction** for Consulting structure (tag-based or relationship type)
- **Smart Positioning Statement generation** (derive from events/facts)
- **Variant sharing/export** between users
- **AI-assisted variant suggestion** based on job description

---

## Documentation Updates (After Implementation)

| File | Changes |
|------|---------|
| `docs/CUSTOM_CV_FORMAT_PROPOSAL.md` | Update to reflect variant system |
| `docs/guides/CV_GENERATION_GUIDE.md` | Add variant selection documentation |
| `docs/guides/NARRATIVE_CV_GUIDE.md` | Update to mention language_agnostic variants |
| `AGENTS.md` | Update to mention CV variants and new structures |

### New Documentation to Create

| File | Purpose |
|------|---------|
| `docs/guides/CV_VARIANTS_GUIDE.md` | Complete guide to variant selection |
| `docs/guides/CONSULTING_CV_GUIDE.md` | Guide for consulting structure |

---

## Acceptance Criteria

- [ ] 16 built-in variants available via two-dropdown selection
- [ ] 4 structures (standard, narrative, consulting, highlights) render correctly
- [ ] Audience filtering works (uses Fact.AudienceRelevance)
- [ ] Audience affects both bullet count and type
- [ ] Role emphasis affects bullet prioritization
- [ ] Length affects date/company filtering (hard filter)
- [ ] Optional sections omit gracefully when no data
- [ ] Preview shows which sections are omitted and why
- [ ] Export output respects variant configuration
- [ ] All tests pass with zero regressions
- [ ] Code coverage maintained >87%

---

**Status**: ✅ PHASES 0-6 COMPLETE (Phase 7 deferred)
**Next Step**: Documentation updates and future work (Section Configuration)
**Last Updated**: 2026-01-09

### Implementation Summary

**Completed**:
- Phase 0: Audience filtering (15 tests) ✅
- Phase 1: Domain models and 16 built-in variants (31 tests) ✅
- Phase 2: Consulting and Highlights structures (24 tests) ✅
- Phase 3: Role emphasis configuration (11 tests) ✅
- Phase 4: Length format configuration (31 tests) ✅
- Phase 5: TUI variant selection workflow (~40 tests updated) ✅
- Phase 6: ProfileOverride support (13 tests) ✅

**Total**: ~125 new tests, 325/326 CV tests passing (1 clipboard test requires display)

**Key Commits**:
- `c2b21e7` - Audience filtering in bullet generator
- `eed48a8` - Audience filtering in EnhancedBulletGenerator
- `f2c0ffa` - CV variant types and 16 built-in variants
- `b8074e4` - Consulting and Highlights export structures
- `c5dea6b` - Role emphasis configuration and scoring
- `729d265` - Length format configuration and filtering
- `d43b7ca` - TUI variant-based CV generation workflow
- `b518e33` - ProfileOverride wiring to export

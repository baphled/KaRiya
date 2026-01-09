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

### Phase 3: Role Emphasis Dimension

**Objective**: Implement role emphasis-specific bullet filtering and section configuration.

#### Step 3.1: Define Role Emphasis Behaviors
**File**: `internal/service/career/cv/role_emphasis.go` (new)

| Role Emphasis | Primary Categories | Secondary Categories | Bullet Focus |
|---------------|-------------------|---------------------|--------------|
| senior_backend | technical, architecture | product, delivery | Technical depth, product impact |
| staff_principal | leadership, strategy, architecture | technical, mentoring | Architecture, mentorship, cross-team |
| consulting | strategy, delivery, consulting | technical, leadership | Client work, rapid assessment |
| language_agnostic | technical, architecture | all | Multi-language evidence, adaptability |

- [ ] Define `RoleEmphasisConfig` struct
- [ ] Define `GetRoleEmphasisConfig()` function
- [ ] Map each role emphasis to preferred categories and bullet types

#### Step 3.2: Implement Role Emphasis Filtering
**File**: `internal/service/career/cv/enhanced_bullet_generator.go`

- [ ] Add `filterByRoleEmphasis()` method
- [ ] Update scoring to weight bullets by role emphasis match
- [ ] Update `getRoleFilter()` to consider role emphasis (not just target role)

#### Step 3.3: Write Tests
**File**: `internal/service/career/cv/role_emphasis_test.go` (new)

- [ ] Senior backend emphasizes technical bullets
- [ ] Staff/principal emphasizes leadership/architecture bullets
- [ ] Consulting emphasizes client/delivery bullets
- [ ] Language-agnostic emphasizes cross-stack evidence
- [ ] Role emphasis filtering integrates with bullet generator

#### Verification
- [ ] All role emphasis tests pass
- [ ] No regressions in existing tests

---

### Phase 4: Length Dimension

**Objective**: Implement length-based compression with hard date filtering.

#### Step 4.1: Define Length Behaviors
**File**: `internal/service/career/cv/length_format.go` (new)

| Length | Max Years | Max Companies | Confidence Boost | Notes |
|--------|-----------|---------------|------------------|-------|
| full | unlimited | unlimited | 0.0 | All sections |
| standard | 10 | unlimited | +0.05 | All sections |
| short | 5 | unlimited | +0.10 | Reduced sections |
| ultra_short | unlimited | 3 | +0.15 | Highlights only |

- [ ] Define `LengthFormatConfig` struct
- [ ] Define `GetLengthFormatConfig()` function

#### Step 4.2: Implement Date Filtering
**File**: `internal/service/career/cv/cv_generation_service.go`

- [ ] Add `filterEventsByDateRange()` method (hard filter - exclude events outside range)
- [ ] Add `filterEventsByCompanyLimit()` method (keep most recent N companies)
- [ ] Apply filters based on variant's `BulletConfig`

#### Step 4.3: Implement Section Limiting
**File**: `internal/service/career/cv/section_builder.go`

- [ ] Update section building to respect variant's section configuration
- [ ] Ultra-short variants only get highlights structure sections

#### Step 4.4: Write Tests
**File**: `internal/service/career/cv/length_format_test.go` (new)

- [ ] Full length includes all events
- [ ] Standard length excludes events older than 10 years
- [ ] Short length excludes events older than 5 years
- [ ] Ultra-short limits to 3 most recent companies
- [ ] Date filtering is hard (excluded events don't appear)
- [ ] Company limiting keeps most recent companies

#### Verification
- [ ] All length tests pass
- [ ] No regressions

---

### Phase 5: UI Integration

**Objective**: Replace structure selection with two-dropdown variant selection.

#### Step 5.1: Update GenerateCV Types
**File**: `internal/cli/intents/generate_cv.go`

- [ ] Add `GenerateCVStateSelectRoleEmphasis GenerateCVState = "select_role_emphasis"`
- [ ] Add `GenerateCVStateSelectLength GenerateCVState = "select_length"`
- [ ] Remove `GenerateCVStateSelectStructure` (replaced by new states)
- [ ] Add `selectedRoleEmphasis RoleEmphasis` to model
- [ ] Add `selectedLengthFormat LengthFormat` to model
- [ ] Add `selectedVariant *CVVariant` to model
- [ ] Update `GenerateCVResult` to include `SelectedVariant`

#### Step 5.2: Implement Role Emphasis Selection
**File**: `internal/cli/intents/generate_cv_intent.go`

- [ ] Add `updateSelectRoleEmphasis()` handler
- [ ] Add `viewSelectRoleEmphasis()` view
- [ ] Show 4 role emphases with descriptions:
  - Senior Backend (Product Teams) - "Product-focused backend roles"
  - Staff/Principal - "Technical leadership roles"
  - Consulting - "Advisory and fractional roles"
  - Language-Agnostic - "Cross-stack adaptability"
- [ ] Navigate with j/k or arrows, select with Enter
- [ ] Esc goes back to audience selection

#### Step 5.3: Implement Length Selection
**File**: `internal/cli/intents/generate_cv_intent.go`

- [ ] Add `updateSelectLength()` handler
- [ ] Add `viewSelectLength()` view
- [ ] Show 4 length formats with descriptions:
  - Full - "Complete history (3+ pages)"
  - Standard - "Balanced (2-3 pages)"
  - Short - "Condensed (2 pages)"
  - Ultra-Short - "Key highlights (1 page)"
- [ ] Navigate with j/k or arrows, select with Enter
- [ ] Esc goes back to role emphasis selection
- [ ] On Enter, resolve variant and proceed to generation

#### Step 5.4: Update State Flow
**File**: `internal/cli/intents/generate_cv_intent.go`

Old flow:
```
SelectProfile → SelectAudience → SelectStructure → Generating
```

New flow:
```
SelectProfile → SelectAudience → SelectRoleEmphasis → SelectLength → Generating
```

- [ ] Update `updateSelectAudience()` to transition to `SelectRoleEmphasis` (not `SelectStructure`)
- [ ] Update breadcrumbs to show role emphasis and length
- [ ] Update context help for new states
- [ ] Remove structure selection code (replaced)

#### Step 5.5: Wire Variant to Generation
**File**: `internal/cli/intents/generate_cv_intent.go`

- [ ] Update `generateCVAsync()` to use variant configuration
- [ ] Pass variant's confidence threshold to bullet generator
- [ ] Pass variant's date/company filters to generation service
- [ ] Use variant's base structure for section building and export

#### Step 5.6: Write Tests
**File**: `internal/cli/intents/generate_cv_variant_test.go` (new)

- [ ] Transitions from select_audience to select_role_emphasis on enter
- [ ] Shows 4 role emphasis options with descriptions
- [ ] Navigates role emphases with j/k
- [ ] Navigates role emphases with up/down arrows
- [ ] Transitions from select_role_emphasis to select_length on enter
- [ ] Goes back to select_audience on esc from role_emphasis
- [ ] Shows 4 length format options with descriptions
- [ ] Navigates lengths with j/k
- [ ] Navigates lengths with up/down arrows
- [ ] Transitions from select_length to generating on enter
- [ ] Goes back to select_role_emphasis on esc from length
- [ ] Resolves correct variant from dimensions
- [ ] Includes variant in result
- [ ] Cancels on q from any selection state

#### Verification
- [ ] All UI tests pass
- [ ] Manual testing confirms new flow works
- [ ] No regressions in existing tests

---

### Phase 6: Section Configuration

**Objective**: Enable section configuration per variant (order, enable/disable, custom titles).

#### Step 6.1: Implement Section Rendering with Config
**File**: `internal/service/career/cv/section_builder.go`

- [ ] Update `BuildSections()` to accept `[]SectionConfig`
- [ ] Order sections by `SectionConfig.Order`
- [ ] Skip sections where `Enabled: false`
- [ ] Use `SectionConfig.Title` for custom titles
- [ ] Apply `SectionConfig.MinConfidence` for per-section filtering

#### Step 6.2: Update Export Service
**File**: `internal/service/career/cv/export_service.go`

- [ ] Pass section config through export pipeline
- [ ] Respect section order in export output
- [ ] Respect custom titles in export output

#### Step 6.3: Implement Optional Section Logic
**File**: `internal/service/career/cv/section_builder.go`

For sections marked as optional (Positioning Statement, Technical Capabilities, Approach):
- [ ] Check if data exists to populate section
- [ ] If data exists → include section
- [ ] If no data → omit section (don't show empty section)

#### Step 6.4: Write Tests
**File**: `internal/service/career/cv/section_config_test.go` (new)

- [ ] Sections render in configured order
- [ ] Disabled sections are skipped
- [ ] Custom titles are used in output
- [ ] Default config works (all sections, default order, default titles)
- [ ] Optional sections omitted when no data
- [ ] Optional sections included when data exists
- [ ] Per-section confidence threshold applies

#### Verification
- [ ] All section config tests pass
- [ ] Exported CVs respect section configuration

---

### Phase 7: Profile Overrides

**Objective**: Allow variant-specific profile customization.

#### Step 7.1: Implement Profile Merging
**File**: `internal/service/career/cv/cv_helpers.go`

- [ ] Add `MergeProfileWithOverride(base *ProfileConfig, override *ProfileOverride) *NarrativeProfileData`
- [ ] Base profile provides default values
- [ ] Override replaces non-nil fields (ProfessionalTitle, CoreStrengths, CareerDifferentiators, CareerPositioning)
- [ ] Nil override fields fall back to base

#### Step 7.2: Wire Profile Override to Export
**File**: `internal/service/career/cv/export_service.go`

- [ ] Update export methods to accept optional `ProfileOverride`
- [ ] Merge override with base profile before rendering
- [ ] Pass merged profile to narrative/consulting/highlights renderers

#### Step 7.3: Add Profile Overrides to Built-In Variants
**File**: `internal/service/career/cv/variants.go`

| Role Emphasis | ProfessionalTitle | CoreStrengths | CareerDifferentiators |
|---------------|-------------------|---------------|----------------------|
| senior_backend | nil | Backend-focused | nil |
| staff_principal | "Staff Engineer / Technical Lead" | Leadership-focused | Leadership differentiators |
| consulting | "Senior Consulting Engineer" | Client-focused | Advisory differentiators |
| language_agnostic | nil | nil (use defaults) | nil (use defaults) |

- [ ] Add `ProfileOverride` to relevant variants
- [ ] Keep overrides minimal (only what differs from default)

#### Step 7.4: Write Tests
**File**: `internal/service/career/cv/profile_override_test.go` (new)

- [ ] Nil override uses base profile entirely
- [ ] Override replaces ProfessionalTitle when set
- [ ] Override replaces CoreStrengths when set
- [ ] Override replaces CareerDifferentiators when set
- [ ] Override replaces CareerPositioning when set
- [ ] Override preserves unspecified fields from base
- [ ] Export uses merged profile correctly

#### Verification
- [ ] All profile override tests pass
- [ ] Exported CVs show correct profile per variant

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

## Test Count Estimate

| Phase | Tests |
|-------|-------|
| Phase 0 (Audience) | ~10 |
| Phase 1 (Variants) | ~12 |
| Phase 2 (Structures) | ~14 |
| Phase 3 (Role Emphasis) | ~8 |
| Phase 4 (Length) | ~10 |
| Phase 5 (UI) | ~16 |
| Phase 6 (Section Config) | ~8 |
| Phase 7 (Profile Override) | ~7 |
| **Total** | ~85 |

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

**Status**: Ready for implementation
**Next Step**: Phase 0 - Implement audience filtering
**Last Updated**: 2026-01-09

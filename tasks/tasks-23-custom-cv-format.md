# Task 23: CV Structure Selection - Standard and Narrative Templates

## Overview

**Goal**: Add CV structure selection to KaRiya's CV generation system. Users select a structure (Standard or Narrative) *before* CV generation, which determines how the CV content is organized and presented. This is separate from export format (Text/Markdown/YAML) which determines the output file type.

**Time Estimate**: ~2 weeks (5 phases)

**Prerequisites**:
- Existing CV generation system functional (✅ Complete)
- Export service supports Text/Markdown/YAML (✅ Complete)
- Understanding of `docs/CUSTOM_CV_FORMAT_PROPOSAL.md` (✅ Complete)

**Related Docs**:
- `docs/CUSTOM_CV_FORMAT_PROPOSAL.md` - Original proposal (needs update after implementation)
- `docs/guides/CV_GENERATION_GUIDE.md` - CV generation workflow (needs update after implementation)
- `internal/service/career/cv/export_service.go` - Current export implementation

---

## Architecture Overview

### Key Concepts

| Concept | Description | When Selected |
|---------|-------------|---------------|
| **CV Structure** | How content is organized (sections, emphasis, filtering) | Before generation |
| **Export Format** | File output format (text, markdown, yaml) | After preview |

### CV Structures

| Structure | Description | Use Case |
|-----------|-------------|----------|
| `standard` | Traditional CV with Experience, Projects, Skills, Summary sections | Most job applications |
| `narrative` | Language-agnostic professional format with Core Strengths, Technologies, What I Bring sections | Emphasizing pragmatic, cross-domain experience |

### State Flow

```
Select Profile → Select Audience → Select Structure → Generate → Preview → Export Format
                                        ↑                           ↑
                                  Standard | Narrative         Text | Markdown | YAML
```

---

## Success Criteria

- [x] CV Structure selection appears before generation (Standard/Narrative)
- [x] Narrative structure renders CV with correct sections in preview
- [x] Narrative structure exports correctly to Text/Markdown formats
- [x] YAML always uses standard structure (it's data format)
- [x] Profile configuration available in Configure System → Profile
- [x] All tests pass with zero regressions
- [x] Code coverage maintained >87%
- [x] User guide complete with examples

---

## Implementation Phases

### Phase 1: Cleanup - Remove Incorrect Implementation

**Objective**: Remove the incorrectly implemented "Custom" export format

**Background**: Phase 1 was previously implemented incorrectly - "Custom" was added as an export format (alongside Text/Markdown/YAML) when it should be a CV structure selected before generation.

#### Files to Delete
- [x] `internal/service/career/cv/export_custom.go` - Wrong architecture
- [x] `internal/service/career/cv/export_custom_test.go` - Tests for deleted code

#### Files to Modify

##### 1.1: Remove ExportFormatCustom from export_service.go
- [x] Remove `ExportFormatCustom ExportFormat = "custom"` constant (N/A - never existed after reset)
- [x] Remove custom case from `getFileExtension()` (N/A - never existed after reset)

##### 1.2: Remove CVExportFormatCustom from generate_cv.go
- [x] Remove `CVExportFormatCustom CVExportFormat = "custom"` constant (N/A - never existed after reset)

##### 1.3: Revert generate_cv_intent.go to 3 export formats
- [x] Revert `updateExportSelectFormat()` to handle 3 formats (maxIndex = 2) (N/A - never existed after reset)
- [x] Revert `viewExportSelectFormat()` to show only Text, Markdown, YAML (N/A - never existed after reset)
- [x] Remove custom case from `exportCVAsync()` switch (N/A - never existed after reset)
- [x] Remove "Custom" from format name display in view methods (N/A - never existed after reset)

#### Verification
- [x] `go build ./...` compiles successfully
- [x] `go test ./...` passes

**Note**: Phase 1 cleanup was simplified by resetting the branch to `origin/next` which removed all incorrect commits.

---

### Phase 2: Add CV Structure Types and Selection State (TDD)

**Objective**: Add new state for CV structure selection between audience and generation

#### Step 2.1: Write Failing Tests First
**File**: `internal/cli/intents/generate_cv_structure_test.go` (new)

Tests to implement:
- [x] Transitions from select_audience to select_structure on enter
- [x] Transitions from select_structure to generating on enter
- [x] Goes back to select_audience on esc from select_structure
- [x] Cancels on q from select_structure
- [x] Shows Standard and Narrative structure options
- [x] Shows descriptions for each structure
- [x] Highlights currently selected structure
- [x] Navigates between structure options with j/k
- [x] Navigates between structure options with up/down
- [x] Stores selected structure in model
- [x] Defaults to Standard structure
- [x] Includes structure in breadcrumbs
- [x] Includes selected structure in result

#### Step 2.2: Add Types and Constants
**File**: `internal/cli/intents/generate_cv.go`

- [x] Add `CVStructure` type
- [x] Add `CVStructureStandard CVStructure = "standard"` constant
- [x] Add `CVStructureNarrative CVStructure = "narrative"` constant
- [x] Add `GenerateCVStateSelectStructure GenerateCVState = "select_structure"` constant
- [x] Add `selectedCVStructure CVStructure` to `GenerateCVModel`
- [x] Add `structureIndex int` to `GenerateCVModel`
- [x] Add `SelectedStructure CVStructure` to `GenerateCVResult`

#### Step 2.3: Implement State Handler
**File**: `internal/cli/intents/generate_cv_intent.go`

- [x] Add `case GenerateCVStateSelectStructure:` to `Update()` switch
- [x] Implement `updateSelectStructure(msg tea.Msg) tea.Cmd`
  - j/k or up/down: navigate between Standard (0) and Narrative (1)
  - enter: set structure, transition to generating, call generateCVAsync()
  - esc: go back to select_audience
  - q/ctrl+c: cancel intent
  - m: return to main menu
- [x] Implement `viewSelectStructure() string`

#### Step 2.4: Update State Transitions
- [x] Modify `updateSelectAudience()`: enter → `GenerateCVStateSelectStructure` (not generating)
- [x] Remove `generateCVAsync()` call from `updateSelectAudience()` (moved to structure selection)

#### Step 2.5: Update Helper Methods
- [x] Update `getStateContent()` - add case for select_structure
- [x] Update `getBreadcrumbs()` - add "Select Structure" crumb
- [x] Update `getContextHelp()` - add help text for structure selection

#### Step 2.6: Update Result
- [x] Modify `setCompleted()` to include `SelectedStructure` in result and metadata

---

### Phase 3: Add Structure-Aware Preview (TDD)

**Objective**: Preview displays CV using the selected structure

#### Step 3.1: Write Failing Tests First

Tests for Standard Structure Preview:
- [x] Renders CV name as header
- [x] Renders all sections in order
- [x] Renders experience with company headers and dates
- [x] Renders bullets for each group

Tests for Narrative Structure Preview:
- [x] Renders profile header with name, role, location, contact
- [x] Renders Summary section
- [x] Renders Core Strengths section
- [x] Renders Languages & Technologies section
- [x] Renders Selected Experience with confidence filtering
- [x] Renders What I Bring section
- [x] Uses default values when sections are empty
- [x] Filters bullets by confidence >= 0.75

#### Step 3.2: Create Helper File
**File**: `internal/cli/intents/generate_cv_helpers.go` (new)

- [x] Add `NarrativeProfileData` struct
- [x] Add `DefaultNarrativeProfile()` function (hardcoded for Phase 1)
- [x] Add `extractStrengthsFromSections()` helper
- [x] Add `extractTechnologiesFromSections()` helper
- [x] Add `extractValuePropositions()` helper

#### Step 3.3: Split Preview Method
**File**: `internal/cli/intents/generate_cv_intent.go`

- [x] Modify `viewPreview()` to switch on `selectedCVStructure`
- [x] Rename current preview logic to `viewPreviewStandard()`
- [x] Implement `viewPreviewNarrative()` with:
  - Profile header (using DefaultNarrativeProfile())
  - Summary section
  - Core Strengths section (using extractStrengthsFromSections or defaults)
  - Languages & Technologies section (using extractTechnologiesFromSections or defaults)
  - Selected Experience section (filtered by confidence >= 0.75)
  - What I Bring section (using extractValuePropositions or defaults)

---

### Phase 4: Add Structure-Aware Export (TDD)

**Objective**: Export methods use selected structure to render content

#### Step 4.1: Write Failing Tests First
**File**: `internal/service/career/cv/export_service_test.go`

- [x] Exports standard structure to text format
- [x] Exports standard structure to markdown format
- [x] Exports narrative structure to text format
- [x] Exports narrative structure to markdown format
- [x] Always uses standard structure for YAML format
- [x] Returns error for nil CV

**Commits**:
- `7ee650b` - test(service): add structure-aware CV export tests

#### Step 4.2: Create Export Helper File
**File**: `internal/service/career/cv/cv_helpers.go` (new)

- [x] Add `CVStructure` type and constants
- [x] Add `NarrativeProfileData` struct
- [x] Add `DefaultNarrativeProfile()` function
- [x] Add `extractStrengthsFromSections()` helper (not needed - using hardcoded profile)
- [x] Add `extractTechnologiesFromSections()` helper (not needed - using hardcoded profile)
- [x] Add `extractValuePropositions()` helper (not needed - using hardcoded profile)

#### Step 4.3: Add Export Method with Structure
**File**: `internal/service/career/cv/export_service.go`

- [x] Add `Export()` method with structure parameter
- [x] Add `exportStandard()` internal method (delegates to existing ExportToText/Markdown)
- [x] Add `exportNarrative()` internal method (narrative structure rendering)
- [x] Add `exportNarrativeText()` for plain text narrative output
- [x] Add `exportNarrativeMarkdown()` for markdown narrative output
- [x] Keep existing `ExportToText()`, `ExportToMarkdown()`, `ExportToYAML()` unchanged (backward compat)

**Commits**:
- `669501b` - feat(service): implement structure-aware CV export

#### Step 4.4: Wire Export in Intent
**File**: `internal/cli/intents/generate_cv_intent.go`

- [x] Update `exportCVAsync()` to use new `Export()` method with structure

**Commits**:
- `f13128f` - feat(intents): wire structure-aware export in GenerateCV intent

**Phase 4 Status**: ✅ COMPLETE

---

### Phase 5: Profile Configuration ✅ COMPLETE

**Objective**: Allow users to configure narrative profile in Configure System

#### Step 5.1: Extend ProfileConfig
**File**: `internal/config/config.go`

- [x] Add `Title` field (professional title)
- [x] Add `Location` field
- [x] Add `GitHub` field (URL)
- [x] Add `Portfolio` field (URL)
- [x] Add `Languages` field (comma-separated)
- [x] Add `Frontend` field (comma-separated)
- [x] Add `Systems` field (comma-separated)
- [x] Add `CoreStrengths` field (slice)
- [x] Add `WhatIBring` field (slice)

#### Step 5.2: Add Profile Conversion Helper
**File**: `internal/service/career/cv/cv_helpers.go`

- [x] Add `NarrativeProfileFromConfig()` function
- [x] Falls back to defaults for empty fields

#### Step 5.3: Add Profile-Aware Export
**File**: `internal/service/career/cv/export_service.go`

- [x] Add `ExportWithProfile()` method
- [x] Add `exportNarrativeWithProfile()` internal method
- [x] Add `exportNarrativeTextWithProfile()` for text output
- [x] Add `exportNarrativeMarkdownWithProfile()` for markdown output

#### Step 5.4: Update ConfigureSystem Intent
**File**: `internal/cli/intents/configure_system.go`

- [x] Add profile fields to `settingsFromConfig()` (Title, Location, GitHub, Portfolio, Languages, Frontend, Systems)
- [x] Update `applyProfileChange()` to handle new fields
- [x] Update domain description for Profile

#### Step 5.5: Wire Profile to GenerateCV
**File**: `internal/cli/intents/generate_cv.go`

- [x] Add `ProfileConfig` field to `GenerateCVContext`

**File**: `internal/cli/intents/generate_cv_intent.go`

- [x] Update `exportCVAsync()` to use `ExportWithProfile()`

**File**: `internal/cli/app/app.go`

- [x] Load profile config and pass to GenerateCVContext

#### Step 5.6: Add Tests
**File**: `internal/service/career/cv/export_service_test.go`

- [x] Test default profile when config is nil
- [x] Test custom profile in text format
- [x] Test custom profile in markdown format
- [x] Test fallback to defaults for empty fields
- [x] Test YAML ignores profile (uses standard structure)
- [x] Test NarrativeProfileFromConfig with nil config
- [x] Test NarrativeProfileFromConfig with full config
- [x] Test NarrativeProfileFromConfig with partial config

**Commits**:
- `fc59cf4` - feat(intents): add configurable profile for narrative CV exports

---

## File Summary

### Files to Delete
| File | Reason |
|------|--------|
| `internal/service/career/cv/export_custom.go` | Wrong architecture |
| `internal/service/career/cv/export_custom_test.go` | Tests for deleted code |

### Files to Create
| File | Purpose |
|------|---------|
| `internal/cli/intents/generate_cv_structure_test.go` | Tests for structure selection state |
| `internal/cli/intents/generate_cv_helpers.go` | Shared helpers for narrative preview |
| `internal/service/career/cv/cv_helpers.go` | Shared helpers for narrative export |

### Files to Modify
| File | Changes |
|------|---------|
| `internal/service/career/cv/export_service.go` | Remove custom format, add Export() with structure |
| `internal/service/career/cv/export_service_test.go` | Add structure-aware export tests |
| `internal/cli/intents/generate_cv.go` | Add CVStructure type, state, model fields |
| `internal/cli/intents/generate_cv_intent.go` | Add state handler, structure-aware preview, wire export |

---

## Test Count Estimate

| Category | Tests |
|----------|-------|
| Structure selection state | ~12 |
| Preview tests | ~10 |
| Export tests | ~6 |
| **Total new tests** | ~28 |

---

## Narrative CV Output Format

The narrative structure produces content with:

```markdown
# Yomi Colledge

**Senior Software Engineer / Technical Consultant**  
Remote (UK)  
Email: [yomi@boodah.net](mailto:yomi@boodah.net)  
GitHub: https://github.com/baphled  
Portfolio: http://boodah.net

---

## Summary
[Professional summary from CV sections or default]

---

## Core Strengths
- Language-agnostic backend and systems engineering  
- System design and architectural ownership  
[...]

---

## Languages & Technologies
**Languages:** Ruby, Go, PHP, C/C++, JavaScript, Shell  
**Frontend:** Vue.js  
**Systems:** Linux, SQL, APIs, CI/CD, automation

---

## Selected Experience

### Company Name
*Jan 2020 - Present*

- High confidence bullet (>= 0.75)  
- Another achievement  

---

## What I Bring
- Languages as tools, not identity  
- Calm handling of complexity  
[...]

---

**References available on request.**
```

---

## Rollback Plan

If issues arise:

1. Revert commits atomically
2. Structure selection is additive - no existing functionality broken
3. Export formats remain unchanged (Text/Markdown/YAML)

**Impact**: Zero - feature is additive

---

## Documentation Updates (After Implementation)

### Files to Update
| File | Changes |
|------|---------|
| `docs/CUSTOM_CV_FORMAT_PROPOSAL.md` | Rewrite to reflect correct architecture |
| `docs/guides/CV_GENERATION_GUIDE.md` | Add "CV Structure Selection" section |
| `docs/guides/CV_TROUBLESHOOTING.md` | Add narrative structure troubleshooting |
| `AGENTS.md` | Update to mention CV structures |

### Files to Create
| File | Purpose |
|------|---------|
| `docs/guides/NARRATIVE_CV_GUIDE.md` | Explain narrative structure, when to use it |

---

## Resolved Decisions

| Question | Decision | Rationale |
|----------|----------|-----------|
| Terminology | Standard/Narrative (not Default/Custom) | Clearer meaning |
| Selection point | Before generation | Structure affects presentation |
| Helper duplication | Keep separate (intents vs service) | Simpler for now |
| Export backward compat | Keep existing methods unchanged | No breaking changes |
| YAML + Narrative | YAML always standard | YAML is data, not presentation |
| Documentation timing | After implementation | Avoid updating twice |

---

### Phase 6: Role Emphasis Redesign - Technology-Focused CV Generation

**Prerequisites**: Task 39 (User-Defined Skills Management) must be complete

**Objective**: Replace current role emphasis with technology-focused system that uses user-defined skills for CV generation

**Time Estimate**: 3-4 days

#### Context

Currently, CV generation uses "Role Emphasis" (Senior Backend, Staff/Principal, Consulting, Language-Agnostic) to determine presentation style. Phase 6 replaces that with a **Technology Focus** system that:

1. Derives technologies from user-defined skills (Task 39)
2. Allows selecting presentation style: Language Agnostic, Generalist (2-5 techs), or Specialist (1 tech)
3. Allows selecting Focus Area: Backend, Frontend, Fullstack, DevOps (derived from skill categories)
4. Populates CV Skills section with selected technologies

#### New Flow

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

#### Variant ID Structure

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

#### CV Structure Mapping

| Technology Focus | CV Structure |
|-----------------|--------------|
| Language Agnostic | Narrative |
| Generalist | Standard |
| Specialist | Standard |

Ultra-Short always uses Highlights structure regardless of Technology Focus.

#### Implementation Steps

**Step 6.1: Technology Extraction Service**

Create service to aggregate and filter user skills.

Files to create:
- [ ] `internal/service/career/technology/extractor.go` - Technology extraction service
- [ ] `internal/service/career/technology/extractor_test.go` - Extraction tests

**Step 6.2: Focus Area Analyzer**

Suggest focus area based on skill categories.

Files to create:
- [ ] `internal/service/career/technology/focus_area.go` - Focus area analyzer
- [ ] `internal/service/career/technology/focus_area_test.go` - Focus area tests

Analysis Logic:
1. Count skills by category
2. Determine dominant focus area:
   - If backend > 70%: Backend
   - If frontend > 70%: Frontend
   - If devops > 70%: DevOps
   - If mix of backend + frontend: Fullstack
3. Calculate confidence based on distribution
4. Return suggestion with evidence

**Step 6.3: Update Domain Types**

Replace old RoleEmphasis with TechnologyFocus.

Files to modify:
- [ ] `internal/service/career/cv/variants.go` - Remove old RoleEmphasis constants, add TechnologyFocus type, add FocusArea type

Remove:
```go
RoleEmphasisSeniorBackend
RoleEmphasisStaffPrincipal
RoleEmphasisConsulting
RoleEmphasisLanguageAgnostic
```

Add:
```go
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

Update CVVariant:
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

**Step 6.4: Update State Machine**

Add new states for technology and focus area selection.

Files to modify:
- [ ] `internal/cli/intents/generate_cv.go` - Add new states, state data fields, update types

New states:
```go
const (
    GenerateCVStateExtractingTechnologies GenerateCVState = "extracting_technologies"
    GenerateCVStateSelectTechnologyFocus  GenerateCVState = "select_technology_focus"  // Renamed from SelectRoleEmphasis
    GenerateCVStateSelectTechnologies     GenerateCVState = "select_technologies"      // NEW
    GenerateCVStateSelectFocusArea        GenerateCVState = "select_focus_area"        // NEW
    GenerateCVStateSelectLengthFormat     GenerateCVState = "select_length_format"
)
```

New state data:
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

**Step 6.5: Technology Extraction Flow**

Extract technologies after audience selection.

Files to modify:
- [ ] `internal/cli/intents/generate_cv_intent.go` - Add extraction, technology selection, focus area views/handlers

**Step 6.6: Technology Focus Selection View**

Replace Role Emphasis view with Technology Focus.

Files to create:
- [ ] `internal/cli/intents/generate_cv_technology_test.go` - Technology selection tests

**Step 6.7: Technology Selection View**

Allow selecting technologies (multi or single).

View shows:
- For Generalist: "Select 2-5 technologies"
- For Specialist: "Select 1 technology"
- Event count per technology
- Space to toggle, Enter to confirm

**Step 6.8: Focus Area Selection View**

Allow selecting focus area with suggestions.

Files to create:
- [ ] `internal/cli/intents/generate_cv_focus_area_test.go` - Focus area selection tests

View shows:
- 4 options: Backend, Frontend, Fullstack, DevOps
- Highlight suggested option
- Show evidence: "backend: 12 skills, frontend: 3 skills"

**Step 6.9: Update Variant System**

Create new variants, remove old ones.

Files to modify:
- [ ] `internal/service/career/cv/role_emphasis.go` - Remove old configs, add new technology focus configs

Remove all 16 old variants (senior_backend_*, staff_principal_*, consulting_*, language_agnostic_*)

Create new variant lookup:
```go
func GetVariantBySelections(
    techFocus TechnologyFocus,
    focusArea FocusArea,
    length LengthFormat,
    technologies []string,
) (*CVVariant, error)
```

**Step 6.10: Update Bullet Generation**

Filter/prioritize bullets based on selected technologies.

Files to modify:
- [ ] `internal/service/career/cv/bullet_generator.go` - Add technology-based filtering

Technology-based filtering:
- Language Agnostic: No filtering - show all
- Specialist: Strongly filter - only bullets with selected tech
- Generalist: Boost bullets with selected techs, keep others

**Step 6.11: Skills Section Population**

Populate Skills section with selected technologies.

Update section builder to:
- Prioritize selected technologies
- Show selected technologies first, then others
- Group by category
- Include event counts

**Step 6.12: Update Documentation**

Rewrite variant and generation guides.

Files to modify:
- [ ] `docs/guides/CV_VARIANTS_GUIDE.md` - Complete rewrite
- [ ] `docs/guides/CV_GENERATION_GUIDE.md` - Significant updates

#### Phase 6 Acceptance Criteria

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

#### Phase 6 Breaking Changes

**For Users:**
- Old CV configs referencing `senior_backend`, `staff_principal`, `consulting` variants will not work
- Must regenerate CVs using new flow

**For Developers:**
- `RoleEmphasis` type removed
- `RoleEmphasisConfig` removed
- All 16 old variants removed
- New `TechnologyFocus` and `FocusArea` types added

#### Phase 6 Notes

**Variant ID Examples:**

Language Agnostic:
- `agnostic_backend_full`
- `agnostic_frontend_standard`
- `agnostic_fullstack_short`
- `agnostic_devops_ultra_short`

Generalist:
- `generalist_backend_standard`
- `generalist_fullstack_short`

Specialist:
- `specialist_ruby_backend_full`
- `specialist_react_frontend_standard`
- `specialist_kubernetes_devops_short`

**Technology Threshold:**
- Minimum 3 events with a skill to appear in selection
- Prevents noise from rarely-used skills
- Users can always add more skills in Manage Skills (Task 39)

**Focus Area Suggestions:**
- Backend: 70%+ backend category skills
- Frontend: 70%+ frontend category skills
- DevOps: 70%+ devops category skills
- Fullstack: Mix of backend + frontend
- User always makes final choice

**Target Role Independence:**
- Target Role (Principal/Staff/EM/Senior IC) remains in Profile
- Affects bullet caps and filtering (per PRD)
- Not part of variant identity
- Orthogonal to Technology Focus

**CV Structure Mapping:**
- Language Agnostic → Narrative (emphasizes adaptability)
- Generalist → Standard (traditional format)
- Specialist → Standard (traditional format)
- Ultra-Short → Highlights (always, regardless of focus)

---

## References

- [CV Generation Service](../internal/service/career/cv/cv_generation_service.go)
- [Export Service](../internal/service/career/cv/export_service.go)
- [GenerateCV Intent](../internal/cli/intents/generate_cv_intent.go)

---

**Last Updated**: 2026-01-11  
**Status**: ✅ PHASES 1-5 COMPLETE - Phase 6 requires Task 39 (User-Defined Skills)  
**Next Step**: Complete Task 39, then implement Phase 6 (Role Emphasis Redesign)

## Implementation Summary

### Completed Phases

| Phase | Description | Tests Added | Commits | Status |
|-------|-------------|-------------|---------|--------|
| Phase 1 | Cleanup - Reset branch | - | Branch reset | ✅ Complete |
| Phase 2 | CV Structure types and selection state | 22 | 3 commits | ✅ Complete |
| Phase 3 | Structure-aware preview | 42 | 2 commits | ✅ Complete |
| Phase 4 | Structure-aware export | 7 | 3 commits | ✅ Complete |
| Phase 5 | Profile configuration | 11 | 1 commit | ✅ Complete |
| Phase 6 | Role Emphasis Redesign (Technology Focus) | TBD | TBD | ⏳ Blocked by Task 39 |

**Total new tests (Phases 1-5)**: 82  
**Estimated tests (Phase 6)**: ~150

### Key Files Created

| File | Purpose |
|------|---------|
| `internal/cli/intents/generate_cv_structure_test.go` | Structure selection tests |
| `internal/cli/intents/generate_cv_preview_test.go` | Preview rendering tests |
| `internal/cli/intents/generate_cv_helpers.go` | Narrative preview helpers |
| `internal/service/career/cv/cv_helpers.go` | Export helpers, types, and profile conversion |

### Key Files Modified

| File | Changes |
|------|---------|
| `internal/cli/intents/generate_cv.go` | Added CVStructure type, state, model fields, ProfileConfig |
| `internal/cli/intents/generate_cv_intent.go` | Added state handler, structure-aware preview and export with profile |
| `internal/cli/intents/configure_system.go` | Added profile fields (Title, Location, GitHub, etc.) |
| `internal/cli/app/app.go` | Load and pass profile config to GenerateCV |
| `internal/config/config.go` | Added narrative profile fields to ProfileConfig |
| `internal/service/career/cv/export_service.go` | Added Export() and ExportWithProfile() methods |
| `internal/service/career/cv/export_service_test.go` | Added structure-aware and profile-aware export tests |

### Branch and PR

- **Branch**: `feature/custom-cv-export-format`
- **PR**: #44
- **Base**: `next`

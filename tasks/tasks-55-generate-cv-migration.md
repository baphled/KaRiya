# Task: Generate CV Intent Migration

## Summary

Migrate `generate_cv` from a 2,776-line monolith with 3 feature-flagged code paths
to subdirectory architecture. This is the largest intent (P8) and requires
comprehensive E2E coverage before any structural changes.

**Estimated Total Effort**: 32-42 hours across 7 phases
**Deadline**: 2026-03-01 (per `.legacy-intents`)
**Branch**: `refactor/generate-cv-subdirectory-migration`

---

## Current State

| Metric | Value |
|--------|-------|
| `generate_cv_intent.go` | 2,776 lines (4.6x over 600-line limit) |
| `generate_cv.go` | 405 lines (context, types, state) |
| Total intent code | 3,181 lines across 2 flat files |
| States | 22 (16 legacy + 6 wizard) |
| Functions | 73 in intent file alone |
| Feature flags | 2 (`useWizardFlow`, `useScreens`) |
| Active code paths | 1 (wizard flow only) |
| Dead code | ~2,500 lines (legacy + screens flow) |
| Architecture violations | `huh` import in `ExportOptionsModal`, `context.Background()` usage |
| Dormant screens | 3 (`ProfileSelect`, `AudienceSelect`, `Generating`) |
| Unused message types | 8 |
| Subdirectory | Does NOT exist |

### Production Code Path (Wizard Flow)

```
CVStateConfiguring → CVStateExtracting → CVStateGenerating →
CVStateReview → CVStatePreview → CVStateExporting →
GenerateCVStateExporting → GenerateCVStateExportComplete
```

Components used in production:
- `CVConfigWizardModal` (components/, 646 lines) -- ACTIVE
- `CVProgressModal` (components/, 394 lines) -- ACTIVE
- `ExportOptionsModal` (components/, 332 lines) -- ACTIVE, has `huh` violation
- `ReviewScreen` (screens/cv/, 201 lines) -- ACTIVE via factory
- `CVPreviewScreen` (screens/cv/, 458 lines) -- ACTIVE via factory

### Known Bug

`SkillRepository` and `EventRepository` are NOT set in `registrar.go` (lines 183-200),
so `extractTechnologiesAsync()` always returns empty results in production. The tech
extraction step is effectively a no-op.

---

## Acceptance Criteria

- [ ] Phase 1: Service test gaps filled (Extractor error paths)
- [ ] Phase 2: Mock infrastructure available in TestEnv
- [ ] Phase 3: Comprehensive E2E test suite passes (37 scenarios)
- [ ] Phase 4: Dead code removed (~2,500 lines)
- [ ] Phase 5: Subdirectory structure created (5 core + 3 optional files)
- [ ] Phase 6: Modals moved from `components/` to `screens/cv/modals/`
- [ ] Phase 7: Intent reduced to broker pattern (<400 lines)
- [ ] All existing tests still pass after migration
- [ ] `make check-intent-architecture` passes
- [ ] `make check-compliance` passes
- [ ] No `huh` imports outside `forms/`
- [ ] No `context.Background()` in intent code
- [ ] Committed with `make ai-commit`

---

## Phase 1: Service Test Hardening

**Purpose**: Fill coverage gaps in services we will mock during E2E testing.
**Effort**: 1-2 hours

### 1.1 Technology Extractor Error Paths

**File**: `internal/service/career/technology/extractor_test.go`

The `Extractor` has error return paths that are untested:

- [ ] Test: `ExtractFromUser` returns error when `skillRepo.List()` fails
- [ ] Test: `ExtractFromUser` returns error when `eventRepo.List()` fails
- [ ] Test: `ExtractFromUser` handles context cancellation

**Files to modify**:
- `internal/service/career/technology/extractor_test.go` (add ~60 lines)

### 1.2 Verify Existing Service Tests Pass

- [ ] Run `make test-suite SUITE=./internal/service/career/cv/...`
- [ ] Run `make test-suite SUITE=./internal/service/career/technology/...`
- [ ] All existing test cases pass

### 1.3 Service Coverage Baseline

Current coverage (verified by audit):

| Service | Unit Tests | Error Paths | Integration | Verdict |
|---------|------------|-------------|-------------|---------|
| `CVGenerationService` | 22 tests, 865 lines (main file); 34 total across 3 files | 3 strict error tests (9 total negative-path) | 1 integration file (11 tests) | Safe to mock |
| `TechnologyExtractor` | 8 tests, 200 lines (extractor_test.go only) | 0 error tests | None | Needs 1.1 first |
| `ExportService` | 42 tests, 654 lines | 7 error tests | Implicit (real I/O) | Safe to mock |

---

## Phase 2: E2E Mock Infrastructure

**Purpose**: Add service mocking and file verification capabilities to TestEnv.
**Effort**: 2-3 hours

### 2.1 Create Mock Services

**File**: `internal/testutil/e2e/mocks.go` (NEW, ~150 lines)

- [ ] `MockCVGenerationService` implementing `cv.CVGenerationService`
  - Configurable: return specific CV or error
  - Tracks calls for assertion
- [ ] `MockExportService` wrapping export behavior
  - Configurable: return specific path or error
  - Tracks exported CVs
- [ ] Compile-time interface compliance checks

### 2.2 Extend TestEnv

**File**: `internal/testutil/e2e/helpers.go` (modify, ~100 lines added)

Mock configuration methods:
- [ ] `WithMockCVGeneration(cv *career.CVView, err error) *TestEnv`
- [ ] `WithMockCVExport(path string, err error) *TestEnv`
- [ ] `WithMockTechExtraction(techs []*ExtractedTechnology, err error) *TestEnv`

File verification helpers:
- [ ] `AssertFileExists(path string) *TestEnv`
- [ ] `AssertFileExtension(path, expectedExt string) *TestEnv`
- [ ] `AssertFileContains(path, substr string) *TestEnv`
- [ ] `GetFileContent(path string) string`

View polling (for async operations):
- [ ] `WaitForViewContaining(substr string, timeout time.Duration) *TestEnv`

### 2.3 Mock Injection via Registrar

**File**: `internal/cli/app/registrar.go` (modify, ~20 lines)

- [ ] Add mock service detection in `registerGenerateCV()`
- [ ] When test context has mock services, use them instead of real services
- [ ] Fall back to real services when mocks not present

### 2.4 Verify Existing Tests Unaffected

- [ ] All existing E2E tests pass with no changes
- [ ] Write 1 smoke test using mocked CV generation
- [ ] Write 1 smoke test using mocked export failure

---

## Phase 3: Comprehensive E2E Test Suite

**Purpose**: Exhaustive coverage of every user journey through wizard flow.
**Effort**: 16-20 hours (including bug fixing)
**File**: `internal/testutil/e2e/generate_cv_comprehensive_e2e_test.go` (NEW, ~1,200 lines)

Tests simulate real user behavior via `e2e.TestEnv` (keyboard input, view assertions).
No private field access. No internal message injection.

### 3.1 Happy Paths (13 scenarios)

#### Complete Workflows
- [ ] A1: Complete wizard to done (no export)
- [ ] A2: Complete wizard to file export (Markdown)
- [ ] A3: Complete wizard to clipboard export
- [ ] A4: Skip wizard with Ctrl+S (use defaults)

#### Wizard Navigation
- [ ] A5: Tab through all wizard steps
- [ ] A6: Navigate back with Esc at each step (step 3 → 2 → 1)
- [ ] A7: Navigate between review and preview (Enter/Esc)

#### Configuration Variations
- [ ] A8: All profile options generate successfully
- [ ] A9: All audience options (hiring_manager, recruiter, peer)
- [ ] A10: Technology focus variations (language_agnostic, generalist, specialist)
- [ ] A11: CV length variations (1_page, 2_page, standard, detailed)
- [ ] A12: Skills format variations (flat, grouped, categorized)
- [ ] A13: All export formats (Text, Markdown, YAML) with file verification

### 3.2 Sad Paths (8 scenarios)

#### Empty State
- [ ] B1: Empty database shows warning modal, dismiss returns to menu

#### Cancellation
- [ ] B2: Cancel from wizard step 1 (Esc → main menu)
- [ ] B3: Cancel from mid-wizard (step 2/3 → back → cancel)
- [ ] B7: Cancel export modal (Esc → back to preview)

#### Error Handling (Mocked Services)
- [ ] B4: Tech extraction failure → error shown, can retry or cancel
- [ ] B5: CV generation failure → error shown, returns to wizard
- [ ] B6: Export failure → error shown, can retry

#### Stress
- [ ] B8: Rapid state transitions (Tab/Esc/Enter/j/k mashing)

### 3.3 Edge Cases (7 scenarios)

#### Minimal Data
- [ ] C1: Single event in database → CV generated
- [ ] C2: Single profile configured → wizard works
- [ ] C3: No technologies detected → tech step skipped

#### UI Robustness
- [ ] C4: Rapid key presses (20+ Tab presses)
- [ ] C5: Window resize during wizard (WindowSizeMsg injection)
- [ ] C6: Window resize during progress modal
- [ ] C7: Long CV preview scrolling (j, k, PgDn, PgUp, Home, End)

### 3.4 Data Integrity (5 scenarios)

- [ ] D1: Profile selection preserved through full workflow
- [ ] D2: Audience selection preserved through full workflow
- [ ] D3: Technology selections preserved through generation
- [ ] D4: State reset on re-entry after cancellation
- [ ] D5: Session persistence (data survives restart)

### 3.5 Navigation Consistency (4 scenarios)

- [ ] E1: Global 'q' key works from all reachable states
- [ ] E2: Global 'm' key returns to main menu from all states
- [ ] E3: Esc behavior is context-appropriate at every state
- [ ] E4: Help footer shows correct shortcuts at each state

### Bug Fixing Protocol

When a test fails due to a bug in the workflow (not a test bug):
1. Document the bug in this task doc (append to Bugs Discovered section)
2. Fix the bug immediately
3. Verify the fix passes
4. Continue writing tests

---

## Phase 4: Dead Code Removal

**Purpose**: Remove ~2,500 lines of dead code to simplify migration.
**Effort**: 2-3 hours
**Prerequisite**: Phase 3 E2E tests passing (safety net)

### 4.1 Remove Dormant Screens

- [ ] Delete `screens/cv/profile_select.go` + test (99 lines)
- [ ] Delete `screens/cv/audience_select.go` + test (94 lines)
- [ ] Delete `screens/cv/generating.go` + test (104 lines)
- [ ] Run tests, verify nothing breaks

### 4.2 Remove `useScreens` Code Path

- [ ] Remove `useScreens` field from intent struct
- [ ] Remove `activeScreen` field from intent struct
- [ ] Remove `EnableScreens()` method (line 2350)
- [ ] Remove `transitionToScreen()` (lines 178-204)
- [ ] Remove `NewCVProfileSelectScreenFromIntent()` (lines 217-244)
- [ ] Remove `handleScreenResult()` and its 4 sub-handlers (lines 729-842):
  - `handleNavigateResult()` (line 753)
  - `handleCancelResult()` (line 797)
  - `handleSubmitResult()` (line 827)
  - `handleErrorResult()` (line 834)
- [ ] Run tests, verify nothing breaks

### 4.3 Remove Legacy 17-State Code Path

- [ ] Remove 15 legacy `update*` methods (lines 875-2569, ~1,700 lines)
- [ ] Remove 14 legacy `view*` methods (lines 1946-2776, ~750 lines)
- [ ] Remove legacy `updateGenerating()` (starts line 1366)
- [ ] Remove legacy `updatePreview()` (starts line 1404)
- [ ] Remove legacy `updateReview()` (starts line 1438)
- [ ] Remove legacy `updateConfirm()` (starts line 1463)
- [ ] Remove legacy `updateExportSelectFormat()` (starts line 2413)
- [ ] Remove legacy `updateExportSelectLocation()` (starts line 2453)
- [ ] Remove legacy `updateExporting()` (starts line 2499)
- [ ] Remove legacy `updateExportComplete()` (starts line 2530)
- [ ] Remove `getStateContent()` (starts line 1500)
- [ ] Remove `getContextHelp()` (starts line 1540)
- [ ] Remove `getBreadcrumbs()` (starts line 1879)
- [ ] Remove 5 theme helper methods (lines 257-277: getCardStyle, getPrimaryColor, getAccentColor, getErrorColor, getBorderColor)
- [ ] Remove 16 legacy state constants
- [ ] Run tests, verify nothing breaks

### 4.4 Remove Unused Types

- [ ] Remove `CVGeneratedMsg` (lines 299-302)
- [ ] Remove `CVGenerationStartedMsg` (lines 305)
- [ ] Remove `TechnologyFocusSelectedMsg` (lines 323-325)
- [ ] Remove `TechnologiesSelectedMsg` (lines 328-330)
- [ ] Remove `FocusAreaSelectedMsg` (lines 333-335)
- [ ] Remove `LengthFormatSelectedMsg` (lines 338-340)
- [ ] Remove `CVExportFormatSelectedMsg` (lines 392-394)
- [ ] Remove `CVExportOptionSelectedMsg` (lines 397-399)
- [ ] Note: `GenerateCVStateSelectLengthFormat` (line 58) is already covered by Phase 4.3's legacy state constant removal
- [ ] Run tests, verify nothing breaks

### 4.5 Update Legacy Tests

Legacy intent-level tests (`generate_cv_test.go`, `*_navigation_test.go`, etc.)
test the deprecated 17-state flow. After removing legacy code:

- [ ] Audit each test file to determine if it tests legacy or wizard flow
- [ ] Delete tests that ONLY test removed legacy code
- [ ] Keep tests that test wizard flow or shared behavior
- [ ] Rename `generate_cv_wizard_e2e_test.go` → `generate_cv_wizard_integration_test.go`
- [ ] Verify all remaining tests pass

### 4.6 Verify

- [ ] Run `make test` -- all tests pass
- [ ] Run `make check-compliance`
- [ ] Commit dead code removal

---

## Phase 5: Subdirectory Migration

**Purpose**: Create required 5-core-file subdirectory structure.
**Effort**: 3-4 hours
**Prerequisite**: Phase 4 complete (dead code removed, file is ~800 lines)

### 5.1 Create Directory Structure

```
internal/cli/intents/generate_cv/     # NEW
├── context.go       # GenerateCVContext, Validate(), type aliases
├── result.go        # GenerateCVResult
├── constants.go     # GenerateCVState (wizard-only), CVExportFormat, CVExportOption
├── messages.go      # WizardCompleteMsg, CVGenerationCompleteMsg, TechnologiesExtractedMsg, CVExportCompleteMsg
├── intent.go        # GenerateCVIntent struct, New, Init, Update, View, Result (<400 lines)
├── types.go         # Intent struct definition (if intent.go > 300 lines)
├── handlers.go      # handleWizardComplete, handleTechExtracted, handleCVGenerated, handleReviewScreenResult, handlePreviewScreenResult, handleExportComplete
└── helpers.go       # generateCVAsync, extractTechnologiesAsync, exportCVAsync, initWizardFlow, setCompleted, setCancelled, wizard view helpers
```

### 5.2 Extract Core Files

- [ ] Create `generate_cv/constants.go` -- wizard-only state enum + export enums (~40 lines)
- [ ] Create `generate_cv/messages.go` -- 4 active message types only (~50 lines)
- [ ] Create `generate_cv/result.go` -- `GenerateCVResult` struct (~30 lines)
- [ ] Create `generate_cv/context.go` -- `GenerateCVContext`, `Validate()`, type aliases (~150 lines)
- [ ] Run tests after each extraction

### 5.3 Extract Optional Files

- [ ] Create `generate_cv/intent.go` -- lifecycle broker (<400 lines)
- [ ] Create `generate_cv/handlers.go` -- all `handle*` functions (~250 lines)
- [ ] Create `generate_cv/helpers.go` -- async operations, state helpers, view helpers (~300 lines)
- [ ] Create `generate_cv/types.go` -- if intent.go exceeds 300 lines
- [ ] Run tests after each extraction

### 5.4 Fix Architecture Violations

- [ ] Replace `context.Background()` at lines 1037, 2587 with `i.getContext()` or `i.GetContext()`
- [ ] Verify no rendering logic in `intent.go` (only `View()` + modal overlay helper)
- [ ] Verify `intent.go` is under 400 lines
- [ ] Run `make check-intent-architecture`

### 5.5 Delete Old Files

- [ ] Delete `intents/generate_cv.go`
- [ ] Delete `intents/generate_cv_intent.go`
- [ ] Move test files to `intents/generate_cv/` package
- [ ] Run tests after each move

### 5.6 Verify

- [ ] `make check-intent-architecture` passes
- [ ] `make test` -- all tests pass
- [ ] `make check-compliance` passes
- [ ] Commit subdirectory creation

---

## Phase 6: Modal Relocation

**Purpose**: Move modals from deprecated `components/` to correct locations.
**Effort**: 2-3 hours
**Prerequisite**: Phase 5 complete

### 6.1 Create Modals Directory

```
internal/cli/screens/cv/modals/       # NEW
```

### 6.2 Move Config Wizard Modal

- [ ] Move `components/cv_config_wizard_modal.go` → `screens/cv/modals/config_wizard_modal.go`
- [ ] Move test file
- [ ] Update package name to `modals`
- [ ] Update imports (intent, registrar)
- [ ] Run tests

### 6.3 Move Progress Modal

- [ ] Move `components/cv_progress_modal.go` → `screens/cv/modals/progress_modal.go`
- [ ] Move test file
- [ ] Update package name to `modals`
- [ ] Update imports (intent)
- [ ] Run tests

### 6.4 Move Export Modal + Fix huh Violation

- [ ] Move `components/export_options_modal.go` → `screens/cv/modals/export_modal.go`
- [ ] Move test file
- [ ] Update package name to `modals`
- [ ] **Fix**: Replace direct `huh` import with `forms/` package usage
- [ ] Update imports (intent)
- [ ] Run tests

### 6.5 Update Registrar

- [ ] Update `app/registrar.go` imports to reference new modal locations
- [ ] Verify factory wiring still works
- [ ] Run tests

### 6.6 Verify

- [ ] `make check-compliance` passes (no `huh` imports)
- [ ] `make test` -- all tests pass
- [ ] Commit modal relocation

---

## Phase 7: Final Validation

**Purpose**: Verify complete migration meets all architecture requirements.
**Effort**: 1-2 hours

### 7.1 Architecture Checks

- [ ] `make check-intent-architecture` passes
- [ ] `make check-compliance` passes
- [ ] `make check-patterns` passes
- [ ] `make check-patterns-strict` passes

### 7.2 Size Verification

| File | Target | Max |
|------|--------|-----|
| `intent.go` | <400 lines | 600 lines |
| `handlers.go` | <400 lines | 800 lines |
| `helpers.go` | <300 lines | 500 lines |
| `context.go` | <200 lines | N/A |
| `result.go` | N/A | N/A |
| `constants.go` | N/A | N/A |
| `messages.go` | N/A | N/A |
| Each screen | <300 lines | N/A |
| Each modal | <400 lines | N/A |

### 7.3 Dependency Verification

- [ ] No circular dependencies
- [ ] `screens/cv/` does NOT import `intents/`
- [ ] `screens/cv/modals/` does NOT import `intents/`
- [ ] No `huh` imports outside `forms/`
- [ ] No `context.Background()` in intent code
- [ ] No forbidden comment markers

### 7.4 Test Suite Health

- [ ] `make test` -- all tests pass
- [ ] `make ci-local` -- full CI passes
- [ ] E2E coverage: 37+ scenarios passing

### 7.5 Metrics

| Metric | Before | After | Change |
|--------|--------|--------|---------|
| Intent file lines | 2,777 | <400 | -85% |
| Total intent package lines | 3,183 | ~1,200 | -62% |
| Dead code | ~2,500 | 0 | -100% |
| Architecture violations | 3 | 0 | -100% |
| Dormant screens | 3 | 0 | -100% |
| Unused message types | 8 | 0 | -100% |
| E2E test scenarios | 34 active + 12 pending baseline | 37+ (comprehensive) | Replaced |
| States | 22 | 8 | -64% |

### 7.6 Fix Known Bug

- [ ] Wire `SkillRepository` and `EventRepository` in `registerGenerateCV()`
- [ ] Verify tech extraction works with real services
- [ ] Run tests, verify fix passes

---

## Bugs Discovered

_This section will be populated during Phase 3 as E2E tests uncover issues._

| Bug | Phase Found | Severity | Status |
|-----|------------|----------|--------|
| `SkillRepository`/`EventRepository` not set in registrar | Pre-task audit | Medium | Known |
| | | | |

---

## Dependencies

- Phase 1 depends on nothing
- Phase 2 depends on Phase 1 (service tests must pass before mocking)
- Phase 3 depends on Phase 2 (mock infrastructure required for E2E tests)
- Phase 4 depends on Phase 3 (E2E tests are safety net for dead code removal)
- Phase 5 depends on Phase 4 (smaller file easier to restructure)
- Phase 6 depends on Phase 5 (intent imports must be updated first)
- Phase 7 depends on Phase 6 (final validation)

---

## Risks

| Risk | Impact | Likelihood | Mitigation |
|-------|--------|------------|------------|
| Wizard form Tab sequence doesn't work in E2E | High | Medium | Start with simplest case, add logging |
| Async operations hang in tests | High | Medium | Use `WaitForViewContaining()` with generous timeouts |
| Mock injection breaks real service behavior | High | Low | Fall back to real services when mocks absent |
| Legacy test removal breaks coverage metrics | Medium | Medium | Audit each test file before deletion |
| Modal relocation creates import cycles | Medium | Low | Follow established pattern from browse_timeline |
| Export file verification flaky on CI | Low | Medium | Use TestEnv temp directory, clean in AfterEach |
| Discovery of fundamental bugs | Variable | High | Fix as we go (immediate protocol) |
| Window resize handling untested in code | Medium | Low | Skip if not feasible to test |
| Multiple export formats not all wired up correctly | Low | Medium | Test one format thoroughly, others lightly |

---

## Definition of Done

- [ ] All 7 phases complete
- [ ] 37+ E2E test scenarios passing
- [ ] `generate_cv/intent.go` under 400 lines
- [ ] Zero architecture violations
- [ ] Zero dead code
- [ ] Dormant screens removed
- [ ] Modals moved to correct location
- [ ] All existing tests still pass after migration
- [ ] `make check-intent-architecture` passes
- [ ] `make check-compliance` passes
- [ ] `make check-patterns` passes
- [ ] No `huh` imports outside `forms/`
- [ ] No `context.Background()` in intent code
- [ ] Committed with `make ai-commit`

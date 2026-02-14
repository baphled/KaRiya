# Cucumber Tests: Replace DB Assertions with View Assertions

## TL;DR

> **Quick Summary**: Refactor all cucumber/godog BDD tests to assert against the rendered TUI view (`env.GetView()`) instead of querying the database directly (`env.GetEvents/Facts/Skills/Bursts()`). Also fix "When" steps that bypass the UI entirely by writing directly to the DB, and update views that don't currently display data needed for assertions.
>
> **Deliverables**:
> - All 46 DB/model-checking "Then" step functions migrated to view assertions
> - All 9 UI-bypassing "When" step functions fixed to drive the actual UI
> - View components updated to display data needed for assertions
> - Feature file Gherkin text updated where it references "database"
> - Zero regressions — same scenarios pass
>
> **Estimated Effort**: XL
> **Parallel Execution**: YES - limited parallelism (domain-by-domain, but view audit + feature text can be parallel)
> **Critical Path**: View Audit → Browse → Facts → Skills → Bursts → Capture → Onboarding → Feature File Text

---

## Context

### Original Request
Cucumber tests don't test that the view changes when data is modified. They circumvent the process by checking the DB directly. This has led to views not updating correctly, but tests still passing because they only verify DB state.

### Interview Summary
**Key Discussions**:
- Three layers of DB-circumvention found: Then assertions (46 functions), When bypasses (9 functions), and convenient TestEnv DB helpers
- Skills/agents correctly say "test what the user sees" — the code doesn't follow the guidance
- Priority: Fix When and Then steps in parallel per feature domain
- Onboarding must also check view, not just model result
- Views must be updated to show data before tests can assert on it

**Research Findings**:
- capture_steps.go: 71% DB assertions (15 DB functions, 6 view functions)
- facts_steps.go: 25% DB assertions (5 DB, ~15 view)
- skills_steps.go: 25% DB assertions (5 DB, ~15 view)
- browse_steps.go: 5% DB assertions (1 DB, ~20 view)
- onboarding_steps.go: 70% model assertions (7 model, 3 view)
- bursts_steps.go: additional DB checkers (not fully audited — Metis found 2 more When bypasses)

### Metis Review
**Identified Gaps** (addressed):
- Missed 2 UI-bypassing When steps in bursts_steps.go (iSubmitTheBurstForm, iConfirmTheAction) — added
- "Given" steps that seed DB are ACCEPTABLE and out of scope — guardrail set
- View audit is MANDATORY before any refactoring — added as Task 0
- Feature file Gherkin text changes must be SEPARATE from Go code changes — separated
- Count assertions may be fragile via view — acknowledged, approach TBD per view audit
- huh form submission timing may need adjustment — noted as risk
- Pending steps must NOT be fixed during this refactoring — guardrail set

---

## Work Objectives

### Core Objective
Make cucumber BDD tests verify what the USER SEES on screen (the rendered Bubble Tea view), not what's in the database. This ensures that view rendering bugs are caught by tests instead of silently passing.

### Concrete Deliverables
- 46 "Then" step functions using `env.GetView()` instead of `env.GetEvents/Facts/Skills/Bursts()`
- 9 "When" step functions driving the actual UI instead of writing to DB directly
- View components updated where data isn't currently visible
- Feature files with "database"-referencing Gherkin text updated
- Baseline BDD test results preserved (zero regressions)

### Definition of Done
- [ ] `make bdd` passes with same pass/pending/fail counts as baseline
- [ ] `grep -rn 'env\.GetEvents\|env\.GetFacts\|env\.GetSkills\|env\.GetBursts\|env\.AssertEventCount\|env\.AssertBurstCount\|env\.AssertFactCount\|env\.AssertSkillCount' features/steps/ | grep -v 'func iHave\|func theDatabaseIsEmpty'` returns empty
- [ ] `grep -rn 'Repo()\.\(Create\|Update\|List\|Delete\)' features/steps/ | grep -v 'func iHave'` returns empty
- [ ] `make build` passes
- [ ] `make test` passes

### Must Have
- All "Then" steps assert via view, not DB
- All "When" steps drive the UI, not bypass it
- Views display all data that tests need to assert
- Feature file text reflects view-based assertions
- Zero regressions in BDD test outcomes

### Must NOT Have (Guardrails)
- MUST NOT change "Given" (iHave*) step implementations — DB seeding for test fixtures is acceptable BDD practice
- MUST NOT fix pending steps (godog.ErrPending) — those are intentionally unimplemented features
- MUST NOT restructure TestEnv or remove its DB helper methods — they're used by E2E tests too
- MUST NOT change the public API of step registrations (Register*Steps functions)
- MUST NOT modify feature file Gherkin text in the same commit as Go code changes
- MUST NOT start a new domain until the current domain is fully green (`make bdd` passes)
- MUST NOT add view components for data that NO existing passing scenario asserts — only add rendering if needed by an existing passing test

---

## Verification Strategy (MANDATORY)

> **UNIVERSAL RULE: ZERO HUMAN INTERVENTION**
>
> ALL tasks in this plan MUST be verifiable WITHOUT any human action.

### Test Decision
- **Infrastructure exists**: YES (godog, gomega, make bdd)
- **Automated tests**: TDD — write view-based assertions, verify they fail, fix implementation, verify they pass
- **Framework**: godog + gomega

### Agent-Executed QA Scenarios (MANDATORY — ALL tasks)

**Verification Tool by Deliverable Type:**

| Type | Tool | How Agent Verifies |
|------|------|-------------------|
| Step definitions | Bash (make bdd) | Run BDD tests, assert pass counts match baseline |
| Grep verification | Bash (grep) | Run grep commands, assert empty output for DB patterns |
| Build verification | Bash (make build) | Run build, assert exit code 0 |

---

## Execution Strategy

### Parallel Execution Waves

```
Wave 1 (Start Immediately):
├── Task 0: Baseline + View Audit (PREREQUISITE)
└── (nothing else — must complete first)

Wave 2 (After Task 0):
├── Task 1: Browse domain (simplest — 1 DB fix)
└── Task 7: Feature file Gherkin text (independent of Go changes)

Wave 3 (After Task 1):
├── Task 2: Facts domain
└── Task 8: Skill/agent guidance updates (independent)

Wave 4 (After Task 2):
└── Task 3: Skills domain

Wave 5 (After Task 3):
└── Task 4: Bursts domain

Wave 6 (After Task 4):
└── Task 5: Capture domain (hardest — most fixes)

Wave 7 (After Task 5):
└── Task 6: Onboarding domain (separate env)

Critical Path: Task 0 → Task 1 → Task 2 → Task 3 → Task 4 → Task 5 → Task 6
```

### Dependency Matrix

| Task | Depends On | Blocks | Can Parallelize With |
|------|------------|--------|---------------------|
| 0 | None | All others | None |
| 1 | 0 | 2 | 7, 8 |
| 2 | 1 | 3 | 7, 8 |
| 3 | 2 | 4 | 7, 8 |
| 4 | 3 | 5 | 7, 8 |
| 5 | 4 | 6 | 7, 8 |
| 6 | 5 | None | 7, 8 |
| 7 | 0 | None | 1-6, 8 |
| 8 | 0 | None | 1-7 |

### Agent Dispatch Summary

| Wave | Tasks | Recommended Agents |
|------|-------|-------------------|
| 1 | 0 | task(category="deep", load_skills=["cucumber", "e2e-testing", "bubble-tea-testing"], run_in_background=false) |
| 2-7 | 1-6 | task(category="deep", load_skills=["cucumber", "e2e-testing", "bdd-workflow", "tdd-workflow", "bubble-tea-testing"], run_in_background=false) |
| 2 | 7 | task(category="quick", load_skills=["cucumber"], run_in_background=true) |
| 3 | 8 | task(category="quick", load_skills=["cucumber", "e2e-testing"], run_in_background=true) |

---

## TODOs

- [x] 0. View Audit + Baseline Recording (PREREQUISITE) - COMPLETED

  **What to do**:
  - Run `make bdd` and record the baseline pass/pending/fail counts
  - For each entity type (event, fact, skill, burst), navigate to its view screen in TestEnv and dump `GetView()` output
  - Catalogue which fields render in each view and which don't
  - Create a mapping: step function → required data → visible in view? → blocked?
  - Identify which view components need updates (Task dependencies for later tasks)

  **Must NOT do**:
  - Don't modify any code — this is READ-ONLY analysis
  - Don't fix any pending steps

  **Recommended Agent Profile**:
  - **Category**: `deep`
    - Reason: Requires navigating multiple screens and building a comprehensive mapping
  - **Skills**: [`cucumber`, `e2e-testing`, `bubble-tea-testing`]
    - `cucumber`: Understands BDD step patterns and Gherkin mapping
    - `e2e-testing`: Knows TestEnv infrastructure and how to drive the app
    - `bubble-tea-testing`: Knows how to inspect Bubble Tea model views

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Sequential (Wave 1)
  - **Blocks**: Tasks 1-8
  - **Blocked By**: None

  **References**:

  **Pattern References**:
  - `internal/testutil/e2e/helpers.go:1308-1310` - GetView() method returns Model.View() output
  - `internal/testutil/e2e/helpers.go:1322-1331` - AssertViewContains() pattern to follow
  - `internal/testutil/e2e/helpers.go:573-626` - SelectIntentByName() for navigating to screens

  **API/Type References**:
  - `internal/testutil/e2e/helpers.go:54-93` - TestEnv struct with all fields
  - `features/support/env.go:304-306` - NewAppEnv() creates test environment

  **Test References**:
  - `features/steps/capture_steps.go:204-216` - Example view assertion (iShouldSeeTheSuccessMessage)
  - `features/steps/browse_steps.go:120-132` - Example view assertion (iShouldSeeAListOfEvents)

  **Acceptance Criteria**:

  **Agent-Executed QA Scenarios:**

  ```
  Scenario: Record BDD baseline
    Tool: Bash
    Preconditions: Project builds
    Steps:
      1. Run: make bdd 2>&1 | tail -20
      2. Capture pass/pending/fail counts
      3. Save to: .sisyphus/evidence/task-0-bdd-baseline.txt
    Expected Result: Baseline counts recorded
    Evidence: .sisyphus/evidence/task-0-bdd-baseline.txt

  Scenario: Audit event view
    Tool: Bash (go test with print)
    Steps:
      1. Write a small test function that creates TestEnv, adds an event with all fields populated (description, company, project, tags, categories, skills, date), navigates to browse_timeline, and prints GetView()
      2. Run it and capture output
      3. Document which fields appear in rendered output
    Expected Result: Mapping of event fields → visible/not visible
    Evidence: .sisyphus/evidence/task-0-event-view-audit.txt

  Scenario: Audit fact/skill/burst views
    Tool: Bash
    Steps:
      1. Repeat for fact_management, manage_skills, burst_management screens
      2. Document all visible fields
    Expected Result: Complete field visibility mapping
    Evidence: .sisyphus/evidence/task-0-view-audit-complete.txt
  ```

  **Evidence to Capture:**
  - [ ] .sisyphus/evidence/task-0-bdd-baseline.txt — BDD pass/pending/fail counts
  - [ ] .sisyphus/evidence/task-0-view-audit-complete.txt — Field visibility mapping per entity type
  - [ ] .sisyphus/evidence/task-0-blocked-steps.txt — Steps that can't be migrated without view updates

  **Commit**: NO (read-only task, no code changes)

---

- [x] 1. Browse Domain — Replace DB Assertions with View Assertions - PARTIALLY COMPLETE (1/2)

  **What to do**:
  - Fix `iShouldSee1Event()` in browse_steps.go: replace `env.AssertEventCount(1)` with view assertion checking rendered event list
  - Fix `theEventHasSkills()` in browse_steps.go: replace `env.GetEvents()` + `eventRepo.Update()` with UI-based skill assignment (this is a "Given" step that writes to DB for setup — BUT the `eventRepo.Update` call at end is the issue)
  - Verify all other browse_steps.go functions already use view assertions
  - Run `make bdd` — browse scenarios must pass, baseline must be preserved

  **Must NOT do**:
  - Don't modify feature files (.feature) in this task
  - Don't fix pending steps (iShouldSeeDifferentEvents, iShouldSeeTheOriginalEvents, etc.)

  **Recommended Agent Profile**:
  - **Category**: `deep`
    - Reason: Requires understanding the view rendering and adjusting assertions
  - **Skills**: [`cucumber`, `e2e-testing`, `bdd-workflow`, `tdd-workflow`, `bubble-tea-testing`]
    - `cucumber`: BDD step pattern expertise
    - `e2e-testing`: TestEnv and view assertion patterns
    - `bdd-workflow`: Red-Green-Refactor cycle for test refactoring
    - `tdd-workflow`: TDD discipline for each change
    - `bubble-tea-testing`: Bubble Tea view inspection

  **Parallelization**:
  - **Can Run In Parallel**: YES (with Task 7, Task 8)
  - **Parallel Group**: Wave 2
  - **Blocks**: Task 2
  - **Blocked By**: Task 0

  **References**:

  **Pattern References**:
  - `features/steps/browse_steps.go:120-132` - iShouldSeeAListOfEvents (correct view pattern)
  - `features/steps/browse_steps.go:654-662` - iShouldSee1Event (DB pattern to replace)
  - `features/steps/browse_steps.go:737-788` - theEventHasSkills (DB write to replace)

  **Test References**:
  - `features/browse_timeline.feature` - All browse scenarios
  - `features/chained_workflows.feature` - Browse scenarios in workflows

  **Acceptance Criteria**:

  **Agent-Executed QA Scenarios:**

  ```
  Scenario: Browse DB patterns removed
    Tool: Bash (grep)
    Steps:
      1. grep -n 'env\.AssertEventCount\|env\.GetEvents\|eventRepo' features/steps/browse_steps.go | grep -v 'func iHave'
      2. Assert: output is empty (no DB patterns in Then/When steps)
    Expected Result: Zero DB-checking patterns in browse steps
    Evidence: grep output captured

  Scenario: Browse BDD tests pass
    Tool: Bash
    Steps:
      1. Run: make bdd 2>&1 | tail -20
      2. Compare pass/pending/fail with baseline
    Expected Result: Same or better counts
    Evidence: .sisyphus/evidence/task-1-bdd-results.txt
  ```

  **Commit**: YES
  - Message: `refactor(tests): replace DB assertions with view assertions in browse steps`
  - Files: `features/steps/browse_steps.go`
  - Pre-commit: `make bdd`

  **Status Notes**:
  - ✅ FIXED: `iShouldSee1Event()` - Replaced `env.AssertEventCount(1)` with view assertion using `env.GetView()` and checking for table indicators
  - ❌ BLOCKED: `theEventHasSkills()` - Event edit form (CaptureEventFormData) does NOT include Skills field. UI cannot assign skills to events. Requires view update to add Skills to edit form.
  - See: `.sisyphus/evidence/task-0-blocked-steps.txt` for full blocker details

---

- [x] 2. Facts Domain — Replace DB Assertions + Fix UI Bypasses - COMPLETED (6/7 functions)

  **What to do**:
  - Fix 5 "Then" DB-checking functions: `thereShouldBeNFacts`, `theFactShouldHaveText`, `theFactShouldHaveCategories`, `theFactShouldHaveAudiences`, `thereShouldBeAFactWithText`
  - Fix hybrid function `iShouldSeeFactText` — remove DB dependency, use only view
  - Fix 2 "When" UI-bypassing functions: `iSubmitTheFactForm` and `iSaveTheFactEdit` — drive the actual huh form UI instead of repo calls
  - If view audit (Task 0) shows fact categories/audiences don't render, update the fact detail view to show them
  - Run `make bdd` after each function change

  **Must NOT do**:
  - Don't modify feature files
  - Don't change "Given" steps (iHaveNFactsInMyProfile, iHaveAFact, etc.)
  - Don't fix pending steps

  **Recommended Agent Profile**:
  - **Category**: `deep`
    - Reason: Requires form-driving logic for fact editor + view updates
  - **Skills**: [`cucumber`, `e2e-testing`, `bdd-workflow`, `tdd-workflow`, `bubble-tea-testing`, `huh`]
    - `cucumber`: BDD step patterns
    - `e2e-testing`: TestEnv view assertions
    - `bdd-workflow`: Red-Green-Refactor
    - `tdd-workflow`: TDD discipline
    - `bubble-tea-testing`: View inspection
    - `huh`: huh form interaction patterns for driving fact editor

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 3 (sequential after Task 1)
  - **Blocks**: Task 3
  - **Blocked By**: Task 1

  **References**:

  **Pattern References**:
  - `features/steps/facts_steps.go:263-270` - thereShouldBeNFacts (DB pattern to replace)
  - `features/steps/facts_steps.go:335-369` - iSubmitTheFactForm (UI bypass to fix)
  - `features/steps/facts_steps.go:634-665` - iSaveTheFactEdit (UI bypass to fix)
  - `features/steps/facts_steps.go:106-118` - iShouldSeeAListOfFacts (correct view pattern to follow)

  **API/Type References**:
  - `internal/domain/career/fact.go` - Fact domain struct (fields to check)
  - `internal/testutil/e2e/helpers.go:1174-1188` - SubmitFact (current bypass mechanism)

  **Test References**:
  - `features/fact_management.feature` - All fact scenarios

  **Acceptance Criteria**:

  **Agent-Executed QA Scenarios:**

  ```
  Scenario: Facts DB patterns removed
    Tool: Bash (grep)
    Steps:
      1. grep -n 'env\.GetFacts\|env\.AssertFactCount\|factRepo' features/steps/facts_steps.go | grep -v 'func iHave'
      2. Assert: output is empty
    Expected Result: Zero DB-checking patterns in facts steps
    Evidence: grep output captured

  Scenario: Facts BDD tests pass
    Tool: Bash
    Steps:
      1. Run: make bdd 2>&1 | tail -20
      2. Compare with baseline
    Expected Result: Same or better counts
    Evidence: .sisyphus/evidence/task-2-bdd-results.txt
  ```

  **Commit**: YES
  - Message: `refactor(tests): replace DB assertions with view assertions in facts steps`
  - Files: `features/steps/facts_steps.go`, possibly view components
  - Pre-commit: `make bdd`

  **Status Notes**:
  - ✅ FIXED: `iShouldSeeFactText()` - Removed `env.GetFacts()`, now uses view assertion only
  - ✅ FIXED: `thereShouldBeNFacts()` - Replaced `env.AssertFactCount()` with view-based counting (counts ✓ symbols)
  - ✅ FIXED: `theFactShouldHaveText()` - Replaced `env.GetFacts()` with `env.GetView()` assertion
  - ✅ FIXED: `theFactShouldHaveCategories()` - Replaced DB query with view assertion (handles 30-char truncation)
  - ✅ FIXED: `thereShouldBeAFactWithText()` - Replaced DB iteration with view assertion (handles 50-char truncation)
  - ✅ FIXED: `iSubmitTheFactForm()` - Replaced repo calls with UI form submission (Tab navigation + Ctrl+S)
  - ✅ FIXED: `iSaveTheFactEdit()` - Replaced repo calls with UI form submission (TypeText + Ctrl+S)
  - ❌ BLOCKED: `theFactShouldHaveAudiences()` - Audiences only visible in fact detail view, no navigation step exists
  - Removed unused import: `careerrepo`
  - Added helper function: `min(a, b int)` for truncation handling

---

- [ ] 3. Skills Domain — Replace DB Assertions + Fix UI Bypasses

  **What to do**:
  - Fix 5 "Then" DB-checking functions: `thereShouldBeNSkills`, `theSkillShouldHaveName`, `theSkillShouldHaveLevel`, `theSkillShouldHaveYears`, `eachSkillShouldHaveUniqueCategory`
  - Fix 1 "When" UI-bypassing function: `iSubmitTheSkillForm` — drive the actual huh form UI (handle add vs edit through the form, not by querying repo to decide)
  - If view audit shows skill level/years don't render, update skill detail view to show them
  - Run `make bdd` after each function change

  **Must NOT do**:
  - Don't modify feature files
  - Don't change "Given" steps
  - Don't fix pending steps (iShouldSeeTheSkillSuggestionsModal, etc.)

  **Recommended Agent Profile**:
  - **Category**: `deep`
    - Reason: Requires form-driving logic for skill editor + possible view updates
  - **Skills**: [`cucumber`, `e2e-testing`, `bdd-workflow`, `tdd-workflow`, `bubble-tea-testing`, `huh`]

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 4 (sequential after Task 2)
  - **Blocks**: Task 4
  - **Blocked By**: Task 2

  **References**:

  **Pattern References**:
  - `features/steps/skills_steps.go:794-801` - thereShouldBeNSkills (DB pattern)
  - `features/steps/skills_steps.go:622-659` - iSubmitTheSkillForm (UI bypass — note hardcoded "Python"/"TypeScript")
  - `features/steps/skills_steps.go:292-304` - iShouldSeeAListOfSkills (correct view pattern)

  **Test References**:
  - `features/skills_management.feature` - All skill scenarios

  **Acceptance Criteria**:

  **Agent-Executed QA Scenarios:**

  ```
  Scenario: Skills DB patterns removed
    Tool: Bash (grep)
    Steps:
      1. grep -n 'env\.GetSkills\|env\.AssertSkillCount\|skillRepo' features/steps/skills_steps.go | grep -v 'func iHave'
      2. Assert: output is empty
    Expected Result: Zero DB-checking patterns in skills steps
    Evidence: grep output captured

  Scenario: Skills BDD tests pass
    Tool: Bash
    Steps:
      1. Run: make bdd 2>&1 | tail -20
      2. Compare with baseline
    Expected Result: Same or better counts
    Evidence: .sisyphus/evidence/task-3-bdd-results.txt
  ```

  **Commit**: YES
  - Message: `refactor(tests): replace DB assertions with view assertions in skills steps`
  - Files: `features/steps/skills_steps.go`, possibly view components
  - Pre-commit: `make bdd`

---

- [ ] 4. Bursts Domain — Replace DB Assertions + Fix UI Bypasses

  **What to do**:
  - Read bursts_steps.go fully (not fully audited yet)
  - Fix all "Then" DB-checking functions: `theBurstShouldHaveName`, `theBurstShouldHaveDescription`, `theBurstShouldNotBeConfirmed`, `theBurstShouldBeConfirmed`, etc.
  - Fix 2 "When" UI-bypassing functions (identified by Metis): `iSubmitTheBurstForm` and `iConfirmTheAction`
  - If view audit shows burst name/description don't render, update burst detail view
  - Run `make bdd` after each function change

  **Must NOT do**:
  - Don't modify feature files
  - Don't change "Given" steps

  **Recommended Agent Profile**:
  - **Category**: `deep`
    - Reason: Requires form-driving logic + possible view updates for burst management
  - **Skills**: [`cucumber`, `e2e-testing`, `bdd-workflow`, `tdd-workflow`, `bubble-tea-testing`, `huh`]

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 5 (sequential after Task 3)
  - **Blocks**: Task 5
  - **Blocked By**: Task 3

  **References**:

  **Pattern References**:
  - `features/steps/bursts_steps.go` - Full file needs reading (Metis identified iSubmitTheBurstForm:410-443 and iConfirmTheAction:553-573)
  - `internal/testutil/e2e/helpers.go:1230-1241` - SubmitBurstUpdate (bypass mechanism)
  - `internal/testutil/e2e/helpers.go:1255-1277` - ConfirmBurst (bypass mechanism)

  **Test References**:
  - `features/burst_management.feature` - All burst scenarios

  **Acceptance Criteria**:

  **Agent-Executed QA Scenarios:**

  ```
  Scenario: Bursts DB patterns removed
    Tool: Bash (grep)
    Steps:
      1. grep -n 'env\.GetBursts\|env\.AssertBurstCount\|burstRepo' features/steps/bursts_steps.go | grep -v 'func iHave'
      2. Assert: output is empty
    Expected Result: Zero DB-checking patterns in bursts steps
    Evidence: grep output captured

  Scenario: Bursts BDD tests pass
    Tool: Bash
    Steps:
      1. Run: make bdd 2>&1 | tail -20
      2. Compare with baseline
    Expected Result: Same or better counts
    Evidence: .sisyphus/evidence/task-4-bdd-results.txt
  ```

  **Commit**: YES
  - Message: `refactor(tests): replace DB assertions with view assertions in bursts steps`
  - Files: `features/steps/bursts_steps.go`, possibly view components
  - Pre-commit: `make bdd`

---

- [ ] 5. Capture Domain — Replace DB Assertions + Fix UI Bypasses (HARDEST)

  **What to do**:
  - Fix 15 "Then" DB-checking functions (see full list in draft audit)
  - Fix 5 "When" UI-bypassing functions: `iSaveMetadataChanges`/`updateEventMetadata`, `iAcceptTheSuggestedBurst`/`createSuggestedBurstFromEvents`, `iSaveTheBurstEdit`/`createEditedBurstFromEvents`, `iAcceptAllInferredSkills`/`createInferredSkillsFromLatestEvent`
  - Handle success modal timing (DismissSuccessModal bypass is OK — it's a test speed optimisation, not a UI bypass)
  - Handle enrichment review screen navigation
  - If view audit shows event tags/categories/date/skills don't render in relevant screens, update those views
  - Run `make bdd` after each function change — this domain has the most changes

  **Must NOT do**:
  - Don't modify feature files
  - Don't change "Given" steps (iHaveAnEventAtCompany, etc.)
  - Don't change DismissSuccessModal (that's a legitimate test speed helper)

  **Recommended Agent Profile**:
  - **Category**: `deep`
    - Reason: Most complex domain — 20 function changes, multiple UI flows
  - **Skills**: [`cucumber`, `e2e-testing`, `bdd-workflow`, `tdd-workflow`, `bubble-tea-testing`, `huh`]

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 6 (sequential after Task 4)
  - **Blocks**: Task 6
  - **Blocked By**: Task 4

  **References**:

  **Pattern References**:
  - `features/steps/capture_steps.go:287-294` - thereShouldBeNEvents (DB pattern)
  - `features/steps/capture_steps.go:305-314` - theEventShouldHaveDescription (DB pattern)
  - `features/steps/capture_steps.go:633-665` - iSaveMetadataChanges + updateEventMetadata (UI bypass)
  - `features/steps/capture_steps.go:406-463` - iAcceptTheSuggestedBurst + createSuggestedBurstFromEvents (UI bypass)
  - `features/steps/capture_steps.go:465-511` - iAcceptAllInferredSkills + createInferredSkillsFromLatestEvent (UI bypass)
  - `features/steps/capture_steps.go:204-216` - iShouldSeeTheSuccessMessage (CORRECT view pattern to follow)

  **Test References**:
  - `features/capture_event.feature` - All capture scenarios
  - `features/skill_inference_review.feature` - Skill inference scenarios

  **Acceptance Criteria**:

  **Agent-Executed QA Scenarios:**

  ```
  Scenario: Capture DB patterns removed
    Tool: Bash (grep)
    Steps:
      1. grep -n 'env\.GetEvents\|env\.AssertEventCount\|env\.GetBursts\|env\.AssertBurstCount\|env\.GetSkills\|eventRepo\|burstRepo' features/steps/capture_steps.go | grep -v 'func iHave\|func theDatabaseIsEmpty\|func persistEventWithSkills'
      2. Assert: output is empty
    Expected Result: Zero DB-checking patterns in capture Then/When steps
    Evidence: grep output captured

  Scenario: Capture BDD tests pass
    Tool: Bash
    Steps:
      1. Run: make bdd 2>&1 | tail -20
      2. Compare with baseline
    Expected Result: Same or better counts
    Evidence: .sisyphus/evidence/task-5-bdd-results.txt
  ```

  **Commit**: YES (may be split into multiple commits for large changes)
  - Message: `refactor(tests): replace DB assertions with view assertions in capture steps`
  - Files: `features/steps/capture_steps.go`, possibly view components
  - Pre-commit: `make bdd`

---

- [ ] 6. Onboarding Domain — Replace Model Assertions with View Assertions

  **What to do**:
  - Fix 7 "Then" model-checking functions: `profileShouldHaveName`, `profileShouldHaveEmail`, `profileShouldHaveLocation`, `profileShouldHaveTitle`, `profileShouldHaveGitHubUsername`, `profileShouldHavePortfolio`, `profileShouldNotHaveTitle`
  - These use `env.Result()` (OnboardingEnv model result) — replace with `env.View()` (rendered wizard view)
  - NOTE: Onboarding uses OnboardingEnv (separate from TestEnv) — different API
  - The wizard view after completion should show the profile data. If it doesn't, this may need a view update to show a summary.

  **Must NOT do**:
  - Don't modify feature files
  - Don't change wizard navigation steps

  **Recommended Agent Profile**:
  - **Category**: `deep`
    - Reason: Different test env (OnboardingEnv) requires separate understanding
  - **Skills**: [`cucumber`, `e2e-testing`, `bubble-tea-testing`, `huh`]

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 7 (sequential after Task 5)
  - **Blocks**: None
  - **Blocked By**: Task 5

  **References**:

  **Pattern References**:
  - `features/steps/onboarding_steps.go:152-162` - profileShouldHaveName (model pattern to replace)
  - `features/steps/onboarding_steps.go:70-78` - iAmOnStep (CORRECT view pattern — checks View() for step indicator)
  - `features/support/env.go:50-52` - OnboardingEnv.View() method

  **Test References**:
  - `features/onboarding.feature` - All onboarding scenarios

  **Acceptance Criteria**:

  **Agent-Executed QA Scenarios:**

  ```
  Scenario: Onboarding model patterns removed
    Tool: Bash (grep)
    Steps:
      1. grep -n 'env\.Result()' features/steps/onboarding_steps.go
      2. Assert: output is empty (all replaced with view checks)
    Expected Result: Zero model-checking patterns in onboarding steps
    Evidence: grep output captured

  Scenario: Onboarding BDD tests pass
    Tool: Bash
    Steps:
      1. Run: make bdd 2>&1 | tail -20
      2. Compare with baseline
    Expected Result: Same or better counts
    Evidence: .sisyphus/evidence/task-6-bdd-results.txt
  ```

  **Commit**: YES
  - Message: `refactor(tests): replace model assertions with view assertions in onboarding steps`
  - Files: `features/steps/onboarding_steps.go`
  - Pre-commit: `make bdd`

---

- [ ] 7. Feature File Gherkin Text — Update Database References

  **What to do**:
  - Update Gherkin step text in ALL .feature files that references "database":
    - `the database is empty` → `I have no data` or similar domain-appropriate text
    - `the event should be saved` → `I should see the event` or similar
    - `the database should be created` → appropriate CLI output assertion
  - Update corresponding step registrations in Go files to match new step text
  - NOTE: This task only changes .feature file text and step regex patterns — NOT assertion logic

  **Must NOT do**:
  - Don't change assertion implementations (those are handled in Tasks 1-6)
  - Don't change scenario structure or add/remove scenarios

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: Text substitution across files — straightforward
  - **Skills**: [`cucumber`]
    - `cucumber`: Gherkin best practices for step text

  **Parallelization**:
  - **Can Run In Parallel**: YES (with Tasks 1-6, but see coordination note)
  - **Parallel Group**: Wave 2+ (after Task 0)
  - **Blocks**: None
  - **Blocked By**: Task 0
  - **Coordination Note**: If Tasks 1-6 change step regex patterns in Go code, Task 7's feature file text MUST match. Best executed AFTER Tasks 1-6 complete to avoid merge conflicts. Can start in Wave 2 but should be finalised last.

  **References**:

  **Pattern References**:
  - `features/capture_event.feature:8` - `Given the database is empty` (30+ occurrences across all feature files)
  - `features/cli_commands.feature:165-169` - Database CLI scenarios

  **Acceptance Criteria**:

  **Agent-Executed QA Scenarios:**

  ```
  Scenario: No database references in feature files
    Tool: Bash (grep)
    Steps:
      1. grep -rn 'database' features/*.feature | grep -v 'PostgreSQL\|database migration'
      2. Assert: output is empty or only contains legitimate data descriptions (not test setup/assertions)
    Expected Result: No "database is empty" or "database should be created" style steps remain
    Evidence: grep output captured

  Scenario: BDD tests still pass
    Tool: Bash
    Steps:
      1. Run: make bdd 2>&1 | tail -20
    Expected Result: Same counts as baseline
    Evidence: .sisyphus/evidence/task-7-bdd-results.txt
  ```

  **Commit**: YES
  - Message: `refactor(tests): update feature file Gherkin text to remove database references`
  - Files: `features/*.feature`, corresponding step registration changes
  - Pre-commit: `make bdd`

---

- [ ] 8. Update Skills & Agent Guidance — Add Anti-Pattern Documentation

  **What to do**:
  - The skills/agents already say the right thing. But add EXPLICIT anti-patterns to prevent regression:
    - In BDD/cucumber skill at `~/.config/opencode/skills/cucumber/`: add "NEVER use env.GetEvents/Facts/Skills/Bursts() in Then steps — always use env.GetView()"
    - In e2e-testing skill at `~/.config/opencode/skills/e2e-testing/`: add "NEVER bypass UI with direct repo calls in When steps"
    - Add examples of CORRECT vs INCORRECT assertion patterns
  - All agent/skill guidance changes go in `~/.config/opencode/` — NOT in the project repository
  - This is the prevention layer to stop future agents from recreating the problem

  **Must NOT do**:
  - Don't change any test code
  - Don't modify the project's `AGENTS.md` or any files inside the repository for this task
  - All changes MUST go in `~/.config/opencode/skills/` only

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: Documentation update — small scope
  - **Skills**: [`cucumber`, `e2e-testing`]

  **Parallelization**:
  - **Can Run In Parallel**: YES (with Tasks 1-7)
  - **Parallel Group**: Wave 2+ (after Task 0)
  - **Blocks**: None
  - **Blocked By**: Task 0

  **References**:

  **Documentation References**:
  - `~/.config/opencode/skills/cucumber/` — Cucumber skill files (MODIFY HERE)
  - `~/.config/opencode/skills/e2e-testing/` — E2E testing skill files (MODIFY HERE)

  **Acceptance Criteria**:

  **Agent-Executed QA Scenarios:**

  ```
  Scenario: Anti-patterns documented
    Tool: Bash (grep)
    Steps:
      1. grep -l 'GetEvents\|GetFacts\|GetSkills' in relevant skill/doc files
      2. Assert: anti-pattern documentation exists
    Expected Result: Clear "NEVER do this" guidance for DB assertions in BDD
    Evidence: Documentation content captured
  ```

  **Commit**: YES
  - Message: `docs(skills): add anti-patterns for DB assertions in cucumber tests`
  - Files: Skill/agent documentation files
  - Pre-commit: N/A (documentation only)

---

## Commit Strategy

| After Task | Message | Files | Verification |
|------------|---------|-------|--------------|
| 0 | No commit (read-only) | N/A | Evidence files saved |
| 1 | `refactor(tests): replace DB assertions with view assertions in browse steps` | browse_steps.go | make bdd |
| 2 | `refactor(tests): replace DB assertions with view assertions in facts steps` | facts_steps.go, possibly view components | make bdd |
| 3 | `refactor(tests): replace DB assertions with view assertions in skills steps` | skills_steps.go, possibly view components | make bdd |
| 4 | `refactor(tests): replace DB assertions with view assertions in bursts steps` | bursts_steps.go, possibly view components | make bdd |
| 5 | `refactor(tests): replace DB assertions with view assertions in capture steps` | capture_steps.go, possibly view components | make bdd |
| 6 | `refactor(tests): replace model assertions with view assertions in onboarding steps` | onboarding_steps.go | make bdd |
| 7 | `refactor(tests): update feature file Gherkin text to remove database references` | *.feature files, step registrations | make bdd |
| 8 | `docs(skills): add anti-patterns for DB assertions in cucumber tests` | `~/.config/opencode/skills/cucumber/`, `~/.config/opencode/skills/e2e-testing/` | N/A (not in repo) |

---

## Success Criteria

### Verification Commands
```bash
# 1. All BDD tests pass (same as baseline)
make bdd  # Expected: same pass/pending/fail as baseline

# 2. No DB assertions in Then steps
grep -rn 'env\.GetEvents\|env\.GetFacts\|env\.GetSkills\|env\.GetBursts\|env\.AssertEventCount\|env\.AssertBurstCount\|env\.AssertFactCount\|env\.AssertSkillCount' features/steps/ | grep -v 'func iHave\|func theDatabaseIsEmpty'
# Expected: empty

# 3. No direct repo calls in When steps
grep -rn 'Repo()\.\(Create\|Update\|List\|Delete\)' features/steps/ | grep -v 'func iHave'
# Expected: empty

# 4. No model-only assertions in onboarding
grep -n 'env\.Result()' features/steps/onboarding_steps.go
# Expected: empty

# 5. Build passes
make build  # Expected: exit 0

# 6. All tests pass
make test  # Expected: exit 0
```

### Final Checklist
- [ ] All "Must Have" present
- [ ] All "Must NOT Have" absent
- [ ] All BDD tests pass (same baseline)
- [ ] All unit tests pass
- [ ] Build succeeds
- [ ] No DB patterns in Then steps (grep verified)
- [ ] No repo calls in When steps (grep verified)
- [ ] Feature files reference views, not database
- [ ] Skill/agent guidance updated with anti-patterns

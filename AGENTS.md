# KaRiya AI Agent Instructions

**Quick reference for AI coding agents working on KaRiya - a Go TUI application using Bubble Tea framework.**

## 🚀 Getting Started

### Session Start (MANDATORY)
```bash
make session-start   # MUST run first - validates environment
```

### Essential Commands

| Task | Command |
|------|---------|
| **Run all tests** | `make test` |
| **Run single test** | `make individual-test TEST="test description"` |
| **Run test suite** | `make test-suite SUITE=./path/...` |
| **Build** | `make build` |
| **Format** | `make fmt` |
| **Lint** | `make vet && make staticcheck` |
| **Full compliance** | `make check-compliance` |
| **Commit** | `make ai-commit FILE=/tmp/commit.txt` |

## 📋 Development Workflow

### Top-Down BDD (MANDATORY - Start with E2E Tests)

**ALWAYS develop from the outside-in:**

1. **Write E2E test FIRST** - Start with full user journey (`*_e2e_test.go`)
2. **Watch it fail** - Verify it fails for the right reason
3. **Implement top-down** - Write only code needed to pass the E2E test
4. **Drop to lower levels only when needed** - Write unit tests for complex logic

```go
// Example: Feature - Edit event metadata
// File: internal/cli/intents/captureevent/capture_event_e2e_test.go

It("allows editing event metadata", func() {
    // Arrange - Full user journey from start to finish
    intent := NewIntent(ctx, event)
    
    // Act - Simulate user actions
    intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})  // Edit
    intent.Update(tea.KeyMsg{Type: tea.KeyTab})                         // Next field
    intent.Update(tea.KeyMsg{Type: tea.KeyEnter})                       // Submit
    
    // Assert - Verify final state
    Expect(intent.IsComplete()).To(BeTrue())
    Expect(savedEvent.Company).To(Equal("New Company"))
})
```

**Why E2E first?**
- ✅ Proves the feature works end-to-end
- ✅ No integration surprises
- ✅ Drives architecture naturally
- ✅ Documents expected behavior
- ❌ **NOT** bottom-up (build components then hope they integrate)

### TDD Cycle (MANDATORY)
```bash
make tdd-red         # 1. Write failing E2E test FIRST
make tdd-green       # 2. Implement minimal code to pass
make tdd-refactor    # 3. Improve code quality
make check-compliance # 4. Verify all checks pass
make ai-commit       # 5. Commit with AI attribution
```

**Test priority order:**
1. **E2E** (`*_e2e_test.go`) - Full user journey ← **START HERE**
2. **Navigation** (`*_navigation_test.go`) - Multi-screen flows
3. **Integration** - Component interaction
4. **Unit** - Complex logic only (when E2E test is hard to debug)

**When writing E2E tests:**
- Build self-documenting helpers that match user actions (like Cucumber)
- Test complete end-to-end flows, not isolated components
- Use `env.SelectIntentByName()`, `env.Confirm()`, `env.TypeText()` patterns
- See [E2E_TEST_HELPERS_GUIDE.md](docs/development/E2E_TEST_HELPERS_GUIDE.md)

**📖 Full workflow:** [DEVELOPMENT_WORKFLOW.md](docs/development/DEVELOPMENT_WORKFLOW.md)  
**📖 BDD details:** [BDD_WORKFLOW.md](docs/development/BDD_WORKFLOW.md)  
**📖 Session protocol:** [SESSION_PROTOCOL.md](docs/development/SESSION_PROTOCOL.md)

## 🏗️ Architecture

### Layer Hierarchy (MUST follow)
```
App → Intents → Screens/Modals → UIKit → Behaviors
```

**Dependency Rules:**
- `screens/` **NEVER** imports `intents/` ❌
- `uikit/` **NEVER** imports `screens/` or `intents/` ❌
- Use `ScreenResult` for screen → intent communication

**📖 Full architecture:** [ARCHITECTURE_OVERVIEW.md](docs/development/ARCHITECTURE_OVERVIEW.md)  
**📖 Intent patterns:** [INTENT_ARCHITECTURE_GUIDE.md](docs/INTENT_ARCHITECTURE_GUIDE.md)  
**📖 TUI architecture:** [TUI_ARCHITECTURE.md](docs/architecture/TUI_ARCHITECTURE.md)

## 📝 Code Style

### Imports
```go
import (
    // 1. Standard library
    "context"
    "fmt"
    
    // 2. External (alphabetical)
    tea "github.com/charmbracelet/bubbletea"
    
    // 3. Internal (alphabetical)
    "github.com/baphled/kariya/internal/cli/intents"
    "github.com/baphled/kariya/internal/domain/career"
)
```

### Naming
- **Files:** `snake_case.go`
- **Packages:** `lowercase` (single word)
- **Types:** `PascalCase`
- **Private:** `camelCase`
- **Public:** `PascalCase`

**📖 Conventions:** [NAMING_CONVENTIONS.md](docs/rules/NAMING_CONVENTIONS.md)  
**📖 Screen naming:** [SCREEN_NAMING.md](docs/conventions/SCREEN_NAMING.md)  
**📖 Modal naming:** [MODAL_NAMING.md](docs/conventions/MODAL_NAMING.md)

### Documentation (ENFORCED)
```go
// Package mypackage provides...
package mypackage

// MyFunction does something.
//
// Expected: param must not be nil
// Returns: result or error
// Side effects: None
func MyFunction(param string) error { ... }
```

**📖 Full rules:** [GO_DOCUMENTATION_RULES.md](docs/conventions/GO_DOCUMENTATION_RULES.md)

### Comments (FORBIDDEN)
❌ **NO comments inside function bodies** - extract to named methods instead  
❌ **NO inline comments** - use descriptive names  
❌ **NO markers:** `TODO`, `FIXME`, `HACK`, `XXX`, `NOTE`, `IMPORTANT`

## 🎯 Component Usage

| Need | Use | NOT |
|------|-----|-----|
| Table | `behaviors.TableBehavior[T]` | `table.New()` |
| Forms | `base.FormScreen[T]`, `forms.NewInput()`, `screens/*FormScreen` | Direct `huh.*`, `models/*Form` (legacy) |
| Modals | `feedback.Modal`, `behaviors.RenderModalOverlay()` | Custom code |
| Text | `primitives.Title()`, `primitives.Body()` | Raw lipgloss |
| Colors | `theme.Primary()` | `lipgloss.Color("#xxx")` |
| Layout | `layout.ScreenLayout` | Manual composition |

**📖 UIKit guide:** [UIKIT_GUIDE.md](docs/UIKIT_GUIDE.md)  
**📖 Forms guide:** [FORMS_GUIDE.md](docs/FORMS_GUIDE.md)  
**📖 Modal patterns:** [MODAL_PATTERNS.md](docs/MODAL_PATTERNS.md)  
**📖 Component rules:** [COMPONENT_USAGE_RULES.md](docs/rules/COMPONENT_USAGE_RULES.md)

## 🎬 VHS Demo Generation

When creating or modifying TUI workflows, generate visual documentation.

**IMPORTANT**: When a feature task is marked as "done", the [VHS Demo Generation Prompt](docs/prompts/VHS_DEMO_GENERATION_PROMPT.md) should be triggered to ensure visual documentation exists.

**MANDATORY**: Before marking any UI-related task as complete, you MUST:
1. Check if VHS tapes need to be created or updated
2. Follow the VHS Demo Generation Prompt
3. Include demo evidence in the PR description

**Trigger Conditions**:
| Change Type | Action | Location |
|-------------|--------|----------|
| **New Feature** | Create 3 new tapes | `demos/vhs/tapes/<feature>/` |
| **Bug Fix** | Find and update existing tape | Existing tape location |
| **Enhancement** | Find and update existing tape | Existing tape location |

**Required for**:
- New intents/workflows
- UI changes that affect user experience
- Bug fixes with visual impact

**Demo Scenarios** (all three required for new features):

| Scenario | File | Purpose |
|----------|------|---------|
| **Happy path** | `happy-path.tape` | Successful workflow completion |
| **Sad path** | `sad-path.tape` | Error handling, validation |
| **Edge cases** | `edge-cases.tape` | Cancel, back nav, empty states |

**Workflow**:
```bash
# 1. Create feature demo directory
mkdir -p demos/vhs/tapes/your-feature

# 2. Copy templates
cp demos/vhs/tapes/template/*.tape demos/vhs/tapes/your-feature/

# 3. Customize tapes for your feature

# 4. Generate demos
make vhs-feature FEATURE=your-feature

# 5. Include GIFs in PR description
```

**Quick Commands**:
```bash
make vhs-demos                    # Generate all workflow demos
make vhs-{workflow}               # Generate specific workflow (vhs-capture, vhs-browse, etc.)
make vhs-feature FEATURE=name     # Generate feature demo for PR
make vhs-golden-compare           # Visual regression test
```

**📖 VHS Demo Prompt:** [VHS_DEMO_GENERATION_PROMPT.md](docs/prompts/VHS_DEMO_GENERATION_PROMPT.md)  
**📖 Workflow Guide:** [WORKFLOW_DOCUMENTATION_GUIDE.md](docs/workflows/WORKFLOW_DOCUMENTATION_GUIDE.md)  
**📖 VHS README:** [demos/vhs/README.md](demos/vhs/README.md)

## 🧪 Testing

### E2E Test Location (CRITICAL)

**All E2E tests MUST live in ONE location:**
```
internal/testutil/e2e/{workflow}_e2e_test.go
```

**ONE file per workflow** - consolidate all scenarios for a workflow into a single file:
- `capture_event_e2e_test.go` - ALL capture event scenarios
- `browse_timeline_e2e_test.go` - ALL timeline browsing scenarios
- `generate_cv_e2e_test.go` - ALL CV generation scenarios

**WRONG locations** (must be moved):
- ❌ `internal/cli/intents/{feature}_e2e_test.go`
- ❌ `internal/cli/uikit/{component}_e2e_test.go`
- ❌ Any E2E test outside `internal/testutil/e2e/`

**If E2E tests exist in wrong location:**
1. **FIRST:** Move to `internal/testutil/e2e/`
2. **SECOND:** Consolidate duplicates
3. **THIRD:** Implement feature/fix

**📖 Location guide:** [E2E_TEST_LOCATION_GUIDE.md](docs/development/E2E_TEST_LOCATION_GUIDE.md)

### Test Hierarchy (ALWAYS start at top)
1. **E2E** (`*_e2e_test.go`) - Full user journey ← **START HERE**
2. **Navigation** (`*_navigation_test.go`) - Multi-screen flows
3. **Integration** - Component interaction
4. **Unit** - Complex logic only (when E2E is hard to debug)

**Critical Rule:** Always write the highest-level test possible. Only drop to lower levels when:
- E2E test is too slow (database, network)
- E2E test is hard to debug (complex calculation)
- Testing edge cases in isolation

### Self-Documenting E2E Helpers (Like Cucumber)

Build test helpers that read like user actions:

```go
// ✅ Self-documenting - reads like natural language
It("should allow editing event metadata", func() {
    env.SelectIntentByName("capture_event")
    env.Confirm()
    env.SubmitEvent(testEvent)
    env.PressKeyRune('e')  // Edit
    env.TypeText("Updated company")
    env.Confirm()
    
    env.AssertViewContains("Updated company")
})
```

**Available helpers:** `SelectIntentByName()`, `Confirm()`, `Cancel()`, `TypeText()`, `NavigateDown()`, `AssertViewContains()`, `PopulateTestData()`, and more.

**When adding features, build new helpers that match real user behavior.**

**📖 Testing guides:**
- [E2E_TEST_LOCATION_GUIDE.md](docs/development/E2E_TEST_LOCATION_GUIDE.md) - Where E2E tests live
- [E2E_TEST_HELPERS_GUIDE.md](docs/development/E2E_TEST_HELPERS_GUIDE.md) - Building self-documenting helpers
- [NAVIGATION_TESTING_GUIDE.md](docs/development/NAVIGATION_TESTING_GUIDE.md) - Navigation flows
- [integration-test-strategy.md](docs/integration-test-strategy.md) - E2E vs integration
- [BDD_WORKFLOW.md](docs/development/BDD_WORKFLOW.md) - AAA pattern, fixtures

## 🔧 Common Tasks

```bash
make what-to-use NEED="keyword"   # Lookup component usage
make check-patterns               # Check TUI pattern compliance
make new-feature TASK="..."       # Create feature task
make new-bug BUG="..."            # Create bug report
make coverage                     # Coverage report
```

**📖 Common tasks:** [COMMON_TASKS.md](docs/development/COMMON_TASKS.md)  
**📖 Troubleshooting:** [TROUBLESHOOTING.md](docs/TROUBLESHOOTING.md)

## 🚫 When to REFUSE

**Immediately refuse if asked to:**
- Skip `make session-start`
- Commit directly to `next`/`main` without feature branch
- Write implementation before E2E test (TDD violation)
- Build components bottom-up without E2E test
- Write unit tests before E2E tests (wrong order)
- Create E2E tests outside `internal/testutil/e2e/` directory
- Create duplicate E2E test files for same workflow
- Implement feature when E2E tests are in wrong location (move them first)
- Use `git commit` (must use `make ai-commit`)
- Add `TODO`/`FIXME`/`HACK` comments
- Put comments inside function bodies
- Import `intents/` from `screens/`
- Import `huh` outside `forms/` package
- Use hardcoded colors instead of theme
- Mark UI task complete without VHS demo (see VHS Demo Generation section)

**Refusal template:**
```
I cannot proceed - this violates project rules.

VIOLATION: [specific rule]
See: [relevant doc link]

Required approach: [correct way]
Run: make check-compliance
```

## 📚 Full Documentation Index

### Development Guides
- [DEVELOPMENT_WORKFLOW.md](docs/development/DEVELOPMENT_WORKFLOW.md) - Complete workflow
- [BDD_WORKFLOW.md](docs/development/BDD_WORKFLOW.md) - TDD/BDD specifics
- [E2E_TEST_HELPERS_GUIDE.md](docs/development/E2E_TEST_HELPERS_GUIDE.md) - Building self-documenting helpers
- [E2E_TEST_LOCATION_GUIDE.md](docs/development/E2E_TEST_LOCATION_GUIDE.md) - E2E test organization
- [SESSION_PROTOCOL.md](docs/development/SESSION_PROTOCOL.md) - Session management
- [INTENT_PATTERNS_LIBRARY.md](docs/development/INTENT_PATTERNS_LIBRARY.md) - Intent patterns
- [KEYBOARD_SYSTEM_GUIDE.md](docs/development/KEYBOARD_SYSTEM_GUIDE.md) - Keyboard handling
- [STATE_TRANSITION_PATTERNS.md](docs/development/STATE_TRANSITION_PATTERNS.md) - State machines

### Rules & Guidelines
- [senior-engineer-guidelines.md](docs/rules/senior-engineer-guidelines.md) - SOLID principles
- [go-guidelines.md](docs/rules/go-guidelines.md) - Go idioms
- [atomic-commits.md](docs/rules/atomic-commits.md) - Commit strategy
- [AI_COMMIT_ATTRIBUTION.md](docs/rules/AI_COMMIT_ATTRIBUTION.md) - AI commits
- [FORMS_WORKFLOW_GUIDE.md](docs/rules/FORMS_WORKFLOW_GUIDE.md) - Form handling

### Architecture & Patterns
- [ARCHITECTURE_OVERVIEW.md](docs/development/ARCHITECTURE_OVERVIEW.md) - System overview
- [INTENT_ARCHITECTURE_GUIDE.md](docs/INTENT_ARCHITECTURE_GUIDE.md) - Intent design
- [TUI_ARCHITECTURE.md](docs/architecture/TUI_ARCHITECTURE.md) - TUI structure
- [MODAL_OVERLAY_PATTERN.md](docs/development/MODAL_OVERLAY_PATTERN.md) - Modal overlays

### Migration & Setup
- [INTENT_MIGRATION_TO_SUBDIRECTORY.md](docs/guides/INTENT_MIGRATION_TO_SUBDIRECTORY.md)
- [SCREEN_MODAL_EXTRACTION_GUIDE.md](docs/guides/SCREEN_MODAL_EXTRACTION_GUIDE.md)

### CI/CD
- [BRANCHING_STRATEGY.md](docs/BRANCHING_STRATEGY.md) - Git workflow
- [CI_CD_PIPELINE.md](docs/CI_CD_PIPELINE.md) - Pipeline details
- [CI_LOCAL_GUIDE.md](docs/CI_LOCAL_GUIDE.md) - Run CI locally

### Quick References
- [TASK_QUICK_REF.md](docs/rules/TASK_QUICK_REF.md) - Task workflow
- [COMMIT_QUICK_REFERENCE.md](docs/rules/COMMIT_QUICK_REFERENCE.md) - Commit format
- [COMPLIANCE_QUICK_REF.md](docs/rules/COMPLIANCE_QUICK_REF.md) - Compliance checks
- [AI_COMMIT_CHECKLIST.md](docs/rules/AI_COMMIT_CHECKLIST.md) - Commit checklist

---

**For detailed rules, architectural requirements, and enforcement checks, see the documents above.**

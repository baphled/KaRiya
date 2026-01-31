# KaRiya AI Agent Instructions

## Session Start (MANDATORY)

```bash
make session-start   # MUST run first - validates environment, acknowledges rules
```

**If `session-start` fails, REFUSE to proceed until issues are fixed.**

---

## Critical Rules (Zero Tolerance)

1. **Feature branches** - NEVER commit directly to `next` or `main`. Always create a feature branch first.
2. **PRs target `next`** - Never `main`. Only `next->main` for releases.
3. **TDD** - Write test FIRST, then implementation (Red->Green->Refactor)
4. **Commits** - Use `make ai-commit FILE=<path>` only (not `git commit`)
5. **Compliance** - Run `make check-compliance` before AND after tasks
6. **One task** - One logical change per commit
7. **Senior Engineer Identity** - Apply SOLID, DRY, KISS, YAGNI principles
8. **Architecture Compliance** - Follow layer hierarchy, no shortcuts
9. **Comment Hygiene** - No TODO/FIXME/HACK/XXX in merged code, no inline comments

---

## Comment Rules (STRICTLY ENFORCED)

### Philosophy: Code Over Comments

**Write self-documenting code.** Use clear variable names, extract methods, and apply SOLID principles instead of explaining what code does.

**Comments should explain WHY, never WHAT.** If you need a comment to explain what code does, the code needs refactoring.

---

### Allowed Comment Locations (ONLY)

Comments are ONLY permitted in these locations:

| Location | Purpose | Example |
|----------|---------|---------|
| **Package documentation** | Describe package purpose | `// Package intents implements...` |
| **Type documentation** | Describe type/struct (godoc) | `// BrowseTimelineIntent manages...` |
| **Public function documentation** | Describe exported functions (godoc) | `// NewIntent creates a new...` |
| **Constant/variable groups** | Document const/var blocks | `// State constants for...` |
| **Complex algorithms** | Explain non-obvious "why" | `// Using binary search because...` |

**All other comments are FORBIDDEN.**

---

### Forbidden Comment Locations

#### 1. Inside Function Bodies (STRICTLY FORBIDDEN)

```go
// ❌ BAD - Comments inside function
func (i *Intent) Update(msg tea.Msg) tea.Cmd {
    // Handle screen updates
    cmd, result := i.activeScreen.Update(msg)
    
    // Process the result
    if result != nil {
        return i.handleScreenResult(result)
    }
    
    return cmd
}

// ✅ GOOD - Extract to well-named methods
func (i *Intent) Update(msg tea.Msg) tea.Cmd {
    cmd, result := i.delegateToActiveScreen(msg)
    
    if result != nil {
        return i.handleScreenResult(result)
    }
    
    return cmd
}

func (i *Intent) delegateToActiveScreen(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
    return i.activeScreen.Update(msg)
}
```

#### 2. Inline Comments at End of Lines (FORBIDDEN)

```go
// ❌ BAD - Inline comment
x := 42 // magic number for calculation

// ✅ GOOD - Use named constant instead
const defaultRetryCount = 42

x := defaultRetryCount
```

#### 3. Struct Field Comments (FORBIDDEN)

```go
// ❌ BAD - Field comments
type TableBehavior[T any] struct {
    allItems      []T // Original unfiltered items
    displayItems  []T // Filtered/sorted items
    selectedIndex int // Current selection
}

// ✅ GOOD - Godoc explains data model
// TableBehavior maintains three data representations:
// allItems (original data), displayItems (after filter/sort),
// and selectedIndex (current selection in displayItems).
type TableBehavior[T any] struct {
    allItems      []T
    displayItems  []T
    selectedIndex int
}
```

#### 4. Section Divider Comments (FORBIDDEN)

```go
// ❌ BAD - Section dividers
func (i *Intent) processModals(msg tea.Msg) tea.Cmd {
    // Error modal handling
    if i.errorModal != nil && i.errorModal.IsVisible() {
        return i.updateErrorModal(msg)
    }
    
    // Help modal handling
    if i.helpModal != nil && i.helpModal.IsVisible() {
        return i.updateHelpModal(msg)
    }
    
    return nil
}

// ✅ GOOD - Extract to named methods
func (i *Intent) processModals(msg tea.Msg) tea.Cmd {
    if cmd := i.tryUpdateErrorModal(msg); cmd != nil {
        return cmd
    }
    
    if cmd := i.tryUpdateHelpModal(msg); cmd != nil {
        return cmd
    }
    
    return nil
}
```

#### 5. Switch Case Explanations (FORBIDDEN)

```go
// ❌ BAD - Comments in switch cases
switch action {
case "edit":
    // Extract burst and show edit modal
    burst := actionData["burst"].(*career.Burst)
    i.showEditModal(burst)
    
case "delete":
    // Show delete confirmation
    i.showDeleteConfirmation()
}

// ✅ GOOD - Extract to named handlers
switch action {
case "edit":
    return i.handleEditAction(actionData)
case "delete":
    return i.handleDeleteAction()
}
```

---

### Forbidden Comment Markers (must resolve before merge)

| Marker | Issue | Action |
|--------|-------|--------|
| `TODO` | Incomplete work | Complete the work or create an issue |
| `FIXME` | Known bug | Fix it or create a bug report |
| `HACK` | Technical debt | Refactor properly |
| `XXX` | Attention needed | Resolve the issue |
| `NOTE` | Explanatory comment | Remove or move to godoc |
| `IMPORTANT` | Emphasis | Remove or move to godoc |
| `BUG` | Bug marker | Fix the bug or create a tracked issue |

**No exceptions.** All markers are forbidden.

---

### Exception: E2E Test Files Only

E2E test files (`*_e2e_test.go`) MAY use inline comments for readability:

```go
// ✅ ALLOWED in e2e tests only
env.NavigateDown() // Go to Manual
env.Confirm()      // Select Manual strategy
env.Cancel()       // Go back
```

Regular test files (`*_test.go`) MUST NOT use inline comments. Use well-named variables, helper functions, and descriptive test names instead.

**Rationale**: E2E tests describe complex multi-step user interactions where inline comments clarify intent.

---

### Comment Style (When Comments Are Used)

When writing allowed comments:

1. **Use complete sentences** - Comments should end with a period.
2. **Explain WHY, not WHAT** - The code shows what; comments explain why.
3. **Keep above code** - Never beside it (except in tests).
4. **Be concise** - If it takes a paragraph, refactor the code instead.

---

### Enforcement

Pre-commit hooks check for:
- Forbidden markers (TODO, FIXME, HACK, XXX, NOTE, IMPORTANT)
- Inline comments (except in e2e test files)
- Comments inside function bodies (linter warning)

**When in doubt, delete the comment and refactor the code.**

---

## Architecture (VITAL - Strictly Enforced)

### Layer Hierarchy (MUST follow)

```
App (Router)
    ↓
Intents (State machines, orchestration)
    ↓
Screens (Stateless views) ←→ Modals (Overlays)
    ↓
UIKit (Primitives, Containers, Layout)
    ↓
Behaviors (TableBehavior, CRUDBehavior)
```

### Dependency Rules (BLOCKING violations)

| Layer | Can Import | NEVER Import |
|-------|------------|--------------|
| `intents/` | screens, uikit, behaviors, components | - |
| `screens/` | uikit, behaviors | **intents** (FORBIDDEN) |
| `uikit/` | themes only | **screens, intents** (FORBIDDEN) |
| `behaviors/` | uikit, themes | **screens, intents** (FORBIDDEN) |

**Circular dependencies = immediate rejection.**

### File Structure Rules (STRICTLY ENFORCED)

**Each type in its own location:**

| Type | Contains | Location | Example |
|------|----------|----------|---------|
| **Context** | Input parameters (events, services, config) | Separate file in `intents/` | `browse_timeline.go` |
| **Model** | State wrapper (SHOULD BE FLATTENED) | Prefer: flattened into intent<br>Alternative: separate file | Flatten into `browse_timeline_intent.go`<br>OR `browse_timeline_model.go` |
| **Intent** | Implementation (Update, View, handlers) | `*_intent.go` in `intents/` | `browse_timeline_intent.go` |
| **Screen** | UI component (view + update) | `screens/feature/` package | `screens/browse/list_screen.go` |
| **Modal** | Overlay component | `components/` or `uikit/feedback/` | `components/delete_modal.go` |

**Example structure:**
```
internal/cli/intents/
├── browse_timeline.go          # Context struct
├── browse_timeline_intent.go   # Intent implementation
screens/browse/
├── list_screen.go              # ListScreen
├── detail_screen.go            # DetailScreen
components/
└── delete_modal.go             # DeleteModal
```

**Enforcement**: Check #16 (`check-intent-architecture.sh`) blocks commits with violations.

### Intent Requirements

All intents MUST:
```go
type MyIntent struct {
    *BaseIntent              // REQUIRED: embed BaseIntent
    state      MyState       // State machine enum
    active     bool          // Is intent active
    result     *IntentResult[*MyResult]
    
    // Screens (one per state)
    listScreen   *myfeature.ListScreen
    detailScreen *myfeature.DetailScreen
    
    // Modals (shared across states)
    deleteModal *components.DeleteConfirmModal
}
```

### Screen Requirements

All screens MUST:
```go
type MyScreen struct {
    *base.BaseScreen         // REQUIRED: embed BaseScreen
    table *behaviors.TableBehavior[*MyItem]  // Use behaviors
}

// Return ScreenResult, not mutate intent state
func (s *MyScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "esc" {
            return nil, screens.NewCancelResult("")  // Return result
        }
    }
    return nil, nil
}
```

### Keyboard Handling (MANDATORY)

**ALWAYS use `tea.Key*` constants for special keys** - never use string comparison for special keys.

```go
// GOOD - Use tea.Key* constants for special keys
switch msg.Type {
case tea.KeyEsc:
    return nil, &screens.CancelResult{}
case tea.KeyUp:
    s.table.HandleNavigation("up")
case tea.KeyDown:
    s.table.HandleNavigation("down")
case tea.KeyPgDown:
    s.table.HandleNavigation("pgdn")  // Note: pgdn not pgdown
case tea.KeyPgUp:
    s.table.HandleNavigation("pgup")
case tea.KeyHome:
    s.table.HandleNavigation("home")
case tea.KeyEnd:
    s.table.HandleNavigation("end")
case tea.KeyCtrlD:
    s.table.HandleNavigation("ctrl+d")
case tea.KeyCtrlU:
    s.table.HandleNavigation("ctrl+u")
case tea.KeyEnter:
    // Handle enter
case tea.KeyBackspace:
    // Handle backspace
}

// String comparison ONLY for vim keys and runes
switch msg.String() {
case "j":
    s.table.HandleNavigation("down")
case "k":
    s.table.HandleNavigation("up")
case "g":
    s.table.HandleNavigation("home")
case "G":
    s.table.HandleNavigation("end")
case "q":
    // Quit
}

// BAD - String comparison for special keys (will fail for some keys)
switch msg.String() {
case "esc":        // Use tea.KeyEsc instead
case "pgdown":     // Use tea.KeyPgDown instead (also: pgdown != pgdn)
case "up", "down": // Use tea.KeyUp, tea.KeyDown instead
}
```

**Why**: `msg.String()` returns inconsistent values for special keys across terminals. `tea.Key*` constants work reliably everywhere.

### Modal Requirements

All modals MUST:
```go
func (m *MyModal) View() string {
    // 1. Check visibility
    if !m.visible {
        return ""
    }
    
    // 2. Nil theme guard (REQUIRED)
    theme := m.theme
    if theme == nil {
        theme = themes.NewDefaultTheme()
    }
    
    // 3. Use UIKit with SOLID background (REQUIRED)
    return containers.NewBox(theme).
        Content(content).
        Background(theme.BackgroundColor()).  // REQUIRED for overlays
        Render()
}
```

### Rendering Modals Over Screens

```go
func (i *MyIntent) View() string {
    baseView := i.currentScreen.View()
    
    // Use behaviors.RenderModalOverlay (NOT custom overlay code)
    if i.modal != nil && i.modal.IsVisible() {
        return behaviors.RenderModalOverlay(i.modal, baseView)
    }
    return baseView
}
```

### Forms/Huh Architecture (STRICTLY ENFORCED)

`huh` (the form library) should ONLY be used at the lowest level.

**Import Rules**:

| Package | Can Import `huh`? | Use Instead |
|---------|-------------------|-------------|
| `forms/` | **YES** (only place) | - |
| `screens/` | **NO** | `forms/` package directly |
| `intents/` | **NEVER** | `screens/*FormScreen` |
| `models/` | **DEPRECATED** | Use `screens/` instead |
| `components/` | **DEPRECATED** | Use `screens/` instead |

**Form Primitives** (use these in `forms/` package):
```go
// forms/ package - wraps huh with KaRiya config
forms.NewInput(FieldConfig{...})     // Text input
forms.NewText(FieldConfig{...})       // Text area
forms.NewSelect(key, title, ...)      // Single select
forms.NewMultiSelect(key, title, ...) // Multi select
forms.NewConfirm(key, title, ...)     // Yes/No confirm
forms.NewForm(groups...)              // Form builder
```

**Form Screens** (use these in intents):
```go
// screens/{feature}/form_screen.go - Screen with embedded form
type FormScreen struct {
    *base.BaseScreen
    form     forms.Form    // Direct use of forms package
    formData *MyFormData
}

// Check form state via forms package (NOT huh)
func (s *FormScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
    if forms.IsCompleted(s.form) {
        return nil, screens.NewSubmitResult(s.formData)
    }
    if forms.IsAborted(s.form) {
        return nil, screens.NewCancelResult("form")
    }
    // ...
}
```

**NEW CODE MUST**:
- Create form builders in `forms/` package only
- Create form screens in `screens/{feature}/form_screen.go`
- Use `forms.IsCompleted()` / `forms.IsAborted()` for state checks
- NEVER import `github.com/charmbracelet/huh` outside `forms/`
- **NEVER use `models/` package for forms** (DEPRECATED)

**DEPRECATED PATTERNS** (DO NOT USE):
```go
// ❌ WRONG - models/ package for forms
import "github.com/baphled/kariya/internal/cli/models"
type MyIntent struct {
    form *models.CaptureForm  // DEPRECATED
}

// ✅ CORRECT - screens/ package
import "github.com/baphled/kariya/internal/cli/screens/myfeature"
type MyIntent struct {
    formScreen *myfeature.FormScreen  // Use screen
}
```

**Reference**: [Forms Guide](docs/FORMS_GUIDE.md), [Forms Workflow](docs/rules/FORMS_WORKFLOW_GUIDE.md)

---

### Intent File Organization (STRICT STANDARD)

ALL intents MUST follow this subdirectory structure:

#### File Structure (REQUIRED)

**Core Files (5 required):**
```
intents/{feature}/
├── context.go    # IntentContext struct + Validate() + domain types
├── result.go     # Result struct (20-50 lines)
├── constants.go  # State enum ONLY (15-50 lines)
├── messages.go   # ALL *Msg types (30-100 lines)
└── intent.go     # NewIntent, Init, Update, View, Result (200-400 lines)
```

**Optional Recommended Files (for larger intents):**
```
intents/{feature}/
├── types.go      # Intent struct definition (if intent.go > 300 lines)
├── handlers.go   # ALL handle* functions (see Handler Organization Rules below)
├── helpers.go    # Utilities, async operations, state management (NO handle* functions)
├── filters.go    # Domain-specific filter/sort logic
└── interfaces.go # Service interfaces for dependency injection
```

**Screens Structure:**
```
screens/{feature}/
├── list_screen.go      # List view (TableBehavior)
├── detail_screen.go    # Detail view
├── form_screen.go      # Form view (if applicable)
└── modals/             # Feature-specific modals (if >2 modals)
    ├── filter_modal.go
    ├── search_modal.go
    └── helpers.go      # Modal helpers
```

**Reference Implementation**: `intents/browse_timeline/` and `screens/timeline/`

#### File Responsibilities

| File | Contains | Max Lines | Enforcement |
|------|----------|-----------|-------------|
| **context.go** | IntentContext struct, Validate(), domain types (Filters) | 50-200 | Guideline |
| **result.go** | Result struct | 20-50 | Guideline |
| **constants.go** | State enum ONLY | 15-50 | Guideline |
| **messages.go** | ALL *Msg types | 30-100 | Guideline |
| **intent.go** | NewIntent, Init, Update, View, Result | 200-400 | **BLOCKED at 600** (Check #18) |
| **types.go** | Intent struct definition (optional) | 50-150 | Guideline |
| **handlers.go** | ALL `handle*` functions (optional) | 200-800 | Guideline |
| **helpers.go** | Utilities, async ops (NO `handle*`) (optional) | 100-500 | Guideline |

#### Intent Rules (ENFORCED)

**Intent MUST**:
- Have 5 core files (context.go, result.go, constants.go, messages.go, intent.go)
- Keep intent.go under 400 lines (warning) / 600 lines (hard block)
- Delegate business logic to Context
- Delegate rendering to Screens
- Use centralized modals (`uikit/feedback/`) or feature modals (`screens/{feature}/modals/`)
- Only orchestrate - NO implementation

**Intent MUST NOT**:
- Contain rendering logic (extract to `screens/{feature}/`)
- Contain business logic (move to Context in `context.go`)
- Define State enum in intent.go (must be in `constants.go`)
- Define *Msg types in intent.go (must be in `messages.go`)
- Exceed 600 lines (hard block by Check #18)
- Have >2 render methods (only `View()` + optional helper)

#### Type Location Rules

| Type | File | Rule |
|------|------|------|
| `*IntentContext struct` | `context.go` | Input params, Validate(), domain types |
| `*Result struct` | `result.go` | Output type |
| `*State string` | `constants.go` | State enum ONLY |
| `const (...)` | `constants.go` | State constants |
| `*Msg struct` | `messages.go` | **ALL message types** |
| `*Intent struct` | `intent.go` OR `types.go` | Struct definition |
| ALL `handle*` methods | `handlers.go` (optional) | **ALL handlers must be here** |
| Utility methods | `helpers.go` (optional) | NO `handle*` functions allowed |

#### Handler Organization Rules (MANDATORY)

**ALL functions named `handle*` MUST live in `handlers.go`**. This includes:

| Handler Type | Examples | Location |
|--------------|----------|----------|
| Screen Result Handlers | `HandleCancel`, `HandleNavigate`, `HandleSubmit`, `HandleError` | `handlers.go` |
| Modal Update Handlers | `handleModalUpdates`, `handleDeleteModalUpdate`, `handleEditModalUpdate` | `handlers.go` |
| Message Handlers | `handleBurstSuggestionsLoaded`, `handleFactExtractionComplete` | `handlers.go` |
| Keyboard Handlers | `handleKeyShortcuts`, `handleDetailModalKeypress` | `handlers.go` |
| Internal Dispatchers | `handleScreenResult`, `handleActionData` | `handlers.go` |

**helpers.go MUST NOT contain any `handle*` functions**. It should contain:
- Utility functions (e.g., `loadBurstEvents`, `showBurstDetailModal`)
- Async command builders (e.g., `startBurstDetection`, `createBurstDetectionCmd`)
- State management helpers (e.g., `clearSuggestionState`, `transitionToScreen`)
- View helpers (e.g., `getStateName`, `getContextHelp`, `rebuildModalRegistry`)

**Rationale**: Consolidating all handlers in one file makes it easy to:
1. Find how any message/event is handled
2. Understand the intent's interaction model
3. Ensure consistent handler patterns
4. Review handler logic in code reviews

**Reference Implementation**: `intents/burst_management/handlers.go`

#### Enforcement

Run before every commit:
```bash
make check-intent-architecture  # Checks #17-29 enforce structure
```

**Automated Checks**:
- **Check #17**: Subdirectory structure (blocks incomplete structures)
- **Check #18**: Intent file size (warns >400 lines, blocks >600 lines)
- **Check #19**: No rendering in intent (blocks >2 render methods)
- **Check #20**: Type location (blocks types in wrong files)
- **Check #21**: UIKit component usage
- **Check #22**: Deprecated models package usage
- **Check #23**: Screens directory structure
- **Check #24**: Modal structs in intents/ package (blocks any modal in intents)
- **Check #25**: huh import location (blocks huh import outside forms/)
- **Check #26**: Render methods in helpers.go (blocks >2 render methods)
- **Check #27**: Screens existence for multi-state intents (blocks 2+ states without screens)
- **Check #28**: Modal/Screen location validation (blocks structs in wrong packages)
- **Check #29**: Naming convention enforcement (blocks non-compliant names)

#### AI Agent Behavior

**When creating new intents**:
1. ✅ ALWAYS use subdirectory structure with 5 core files
2. ✅ ALWAYS extract views to `screens/{feature}/`
3. ✅ ALWAYS keep `intent.go` under 400 lines
4. ✅ USE `types.go`, `handlers.go`, `helpers.go` if intent.go > 300 lines
5. ✅ ALWAYS put `*Msg` types in `messages.go`
6. ✅ CREATE `screens/{feature}/modals/` if >2 modals
7. ✅ PUT ALL `handle*` functions in `handlers.go` (NEVER in helpers.go)

**When modifying existing intents**:
- If adding >50 lines → **REFUSE**, suggest migration first
- If intent exceeds 600 lines → **REFUSE** all changes except migration
- If intent not in subdirectory → **REFUSE**, require migration

**Refusal Template**:
```
I cannot modify this intent - it doesn't follow the subdirectory structure.

Current: intents/burst_management_intent.go (1,488 lines)
Required: intents/burst_management/ subdirectory with 5+ files

This intent must be migrated first:
1. Create intents/burst_management/ subdirectory
2. Create 5 core files: context.go, result.go, constants.go, messages.go, intent.go
3. (Optional) Create types.go, handlers.go, helpers.go to keep intent.go small
4. (Optional) Create handlers.go to keep intent.go small
5. Extract views to screens/burst_management/
5. Move modals to screens/burst_management/modals/
6. Reduce intent.go to 200-400 lines (broker only)

See: docs/guides/INTENT_MIGRATION_TO_SUBDIRECTORY.md
```

**Template Location**: `examples/intent_subdirectory_template/`

---

### Architecture Violations to REFUSE

The AI agent MUST refuse code that:

- Has screens importing from `intents/` package
- Has UIKit importing from `screens/` or `intents/`
- **Imports `huh` outside of `forms/` package**
- Uses raw `*huh.Form` in intents (must use `models.*Form` wrapper)
- Uses raw lipgloss styling (must use UIKit primitives)
- Uses `components.KeyBadge` (must use `primitives.HelpKeyBadge()`)
- Uses `components.StandardView` (must use `layout.NewScreenLayout()`)
- Creates modals without solid background
- Mutates intent state from screen code
- Has circular package dependencies
- **Defines modal structs in `intents/` package** (must be in `screens/*/modals/` or `uikit/feedback/`)
- **Defines screen structs outside `screens/` package**
- **Has >2 render methods in helpers.go** (must extract to screens)
- **Has 2+ states without screens directory** (must extract screens)
- **Uses wrong naming conventions** (Screen suffix, Modal suffix, underscore packages)

**Reference**: [Intent Architecture Guide](docs/INTENT_ARCHITECTURE_GUIDE.md)

---

### Naming Conventions (STRICTLY ENFORCED)

#### Screen Naming

| Aspect | Convention | Example |
|--------|------------|---------|
| Package | `screens/{feature_name}/` (underscore) | `screens/fact_management/` |
| File | `{type}.go` | `list.go`, `detail.go`, `form.go` |
| Struct | `{Entity}{Type}Screen` | `FactListScreen`, `FactDetailScreen` |
| Constructor | `New{StructName}()` | `NewFactListScreen()` |

**Screen Types**: `List`, `Detail`, `Form`, `Delete`, `Select`, `Confirm`, `Preview`, `Review`

#### Modal Naming

| Aspect | Convention | Example |
|--------|------------|---------|
| Package | `screens/{feature}/modals/` | `screens/fact_management/modals/` |
| File | `{action}_modal.go` | `edit_modal.go`, `filter_modal.go` |
| Struct | `{Action}Modal` | `EditModal`, `FilterModal` |
| Constructor | `New{StructName}()` | `NewEditModal()` |

**Modal Types**: `Edit`, `Filter`, `Sort`, `Search`, `Confirm`, `Delete`, `Detail`, `QuickAdd`

**Allowed Modal Locations**:
- `internal/cli/uikit/feedback/` (reusable modals)
- `internal/cli/screens/{feature}/modals/` (feature-specific)
- `internal/cli/components/` (legacy, deprecated)

**FORBIDDEN Modal Locations**:
- `internal/cli/intents/` - **NEVER**
- `internal/cli/models/` - **NEVER**

**Reference**: [Screen Naming](docs/conventions/SCREEN_NAMING.md), [Modal Naming](docs/conventions/MODAL_NAMING.md)

---

## Task Types & Workflows

### Development Tasks (New Features)

**Workflow**: `make pre-task` -> TDD Cycle -> `make check-compliance` -> Commit

| Step | Command | Purpose |
|------|---------|---------|
| 1 | `make pre-task` | Verify environment ready |
| 2 | `make tdd-red` | Write failing test FIRST |
| 3 | `make tdd-green` | Minimal code to pass |
| 4 | `make tdd-refactor` | Improve code quality |
| 5 | `make check-compliance` | Verify all checks pass |
| 6 | `make ai-commit FILE=...` | Commit with attribution |

**Reference Docs**:
- [Development Workflow](docs/development/DEVELOPMENT_WORKFLOW.md) - Complete workflow
- [BDD Workflow](docs/development/BDD_WORKFLOW.md) - TDD/BDD specifics
- [Master Task Prompt](docs/rules/master-task-prompt.md) - 5-phase workflow

---

### Testing Tasks

**Test Types**:
| Type | Pattern | Example |
|------|---------|---------|
| Unit tests | `*_test.go` | `service_test.go` |
| E2E tests | `*_e2e_test.go` | `capture_e2e_test.go` |
| Navigation tests | `*_navigation_test.go` | `browse_navigation_test.go` |
| Escape tests | `*_escape_test.go` | `intent_escape_test.go` |

**Commands**:
```bash
make test                          # Run all tests
make test-suite SUITE=./path/...   # Run specific suite
make individual-test TEST="name"   # Run single test
make coverage                      # Generate coverage report
go test -race ./...                # Check race conditions
go test -v ./... -run "TestName"   # Verbose single test
```

**Coverage Requirements**: >= 95% on changed code (enforced by pre-commit hook)

**Reference Docs**:
- [BDD Workflow](docs/development/BDD_WORKFLOW.md) - Test structure with Ginkgo
- [Navigation Testing Guide](docs/development/NAVIGATION_TESTING_GUIDE.md) - TUI navigation tests
- [Integration Test Strategy](docs/integration-test-strategy.md) - E2E patterns

---

### Debugging Tasks

**Debug Commands**:
```bash
go test -v ./... -run "TestName"   # Verbose test output
go test -race ./...                # Race condition detection
go test -cpuprofile=cpu.prof ./... # CPU profiling
go test -bench=. ./path/...        # Run benchmarks
go test -benchmem ./path/...       # Memory profiling
```

**Common Issues & Solutions**:

| Issue | Cause | Solution |
|-------|-------|----------|
| Race conditions | Shared state | Add mutex/channels, check `GlobalContext` |
| Memory leaks | Goroutine leaks | Profile with pprof, check cleanup |
| Slow tests | DB queries | Check N+1 issues, use mocks |
| Multiple Ginkgo entry points | Multiple test files | One Ginkgo suite per package |
| Form alignment issues | Raw `*huh.Form` in intent | Use wrapper model pattern |

**Reference Docs**:
- [Troubleshooting](docs/TROUBLESHOOTING.md) - User-facing issues
- [Error Handling Guide](docs/guides/ERROR_HANDLING_GUIDE.md) - Error patterns
- [Common Tasks](docs/development/COMMON_TASKS.md) - Development troubleshooting

---

### Bug Fixing Tasks

**Workflow**: Report -> Regression Test -> Fix -> Verify

```bash
# 1. Create bug report
make new-bug BUG="description"

# 2. Write regression test FIRST (TDD)
make tdd-red

# 3. Fix the bug
make tdd-green

# 4. Verify and commit
make check-compliance
make ai-commit FILE=/tmp/commit.txt
```

**Bug Test Pattern**:
```go
Describe("Bug Regressions", func() {
    It("BUG-XXX: prevents [bug behavior]", func() {
        // Test that verifies the bug is fixed
    })
})
```

**Reference Docs**:
- [Bug Template](docs/templates/bug-template.md) - Bug report format
- [Bug Task Template](docs/templates/bug-task-template.md) - Bug fix task format

---

### Code Review / Refactoring Tasks

**Pre-refactor Checklist**:
- [ ] Tests exist and pass
- [ ] Scope is limited to current task
- [ ] No behavior changes (tests still pass)

**Refactoring Commands**:
```bash
make fmt                    # Format code
make vet                    # Static analysis
make staticcheck            # Advanced static analysis
make check-patterns         # Check TUI pattern compliance
make check-patterns-strict  # Strict pattern enforcement
```

**Reference Docs**:
- [Senior Engineer Guidelines](docs/rules/senior-engineer-guidelines.md) - SOLID principles
- [Go Guidelines](docs/rules/go-guidelines.md) - Go idioms

---

### CI/CD & Deployment Tasks

**Pre-push Validation**:
```bash
make ci-local              # Run ALL CI checks locally (recommended)
make pre-pr                # Validate before PR (targets next branch)
```

**Individual CI Checks**:
```bash
make ci-install-tools      # Install all CI tools
make fmt                   # Code formatting
make vet                   # Static analysis
make staticcheck           # Staticcheck
make gosec                 # Security scanning
make test                  # All tests
make coverage              # Coverage report
```

**Branch Strategy**:
- `next` - Integration branch, all PRs target here
- `main` - Production releases only (from `next`)

**Reference Docs**:
- [Branching Strategy](docs/BRANCHING_STRATEGY.md) - Branch workflow
- [CI/CD Pipeline](docs/CI_CD_PIPELINE.md) - Pipeline details
- [CI Local Guide](docs/CI_LOCAL_GUIDE.md) - Running CI locally

---

### Documentation Tasks

**Documentation Commands**:
```bash
make generate-docs          # Generate all docs
make generate-diagrams      # Generate Mermaid diagrams
make generate-state-matrix  # Generate state matrix
```

**No TDD Required For**:
- `.md` files (documentation)
- `.yaml`, `.json` files (configuration)
- `.sh` files (scripts)

---

## Essential Make Commands

### Session Management
| Command | Purpose |
|---------|---------|
| `make session-start` | **MUST run first** - validates environment |
| `make session-end` | End session (cleanup) |
| `make session-reset` | Recovery after crash |
| `make pre-task` | Checklist before any task |

### Quality & Compliance
| Command | Purpose |
|---------|---------|
| `make check-compliance` | Full compliance check (before/after tasks) |
| `make check-patterns` | TUI pattern violations |
| `make check-patterns-strict` | Strict pattern check (blocking) |
| `make pre-commit` | Quick pre-commit checks |

### TDD Workflow
| Command | Purpose |
|---------|---------|
| `make tdd-red` | Start: write failing test |
| `make tdd-green` | Make test pass |
| `make tdd-refactor` | Improve code quality |
| `make tdd-document` | Finalize and commit |

### Testing
| Command | Purpose |
|---------|---------|
| `make test` | Run all tests |
| `make test-suite SUITE=...` | Run specific suite |
| `make individual-test TEST=...` | Run single test |
| `make coverage` | Generate coverage report |

### Code Quality
| Command | Purpose |
|---------|---------|
| `make fmt` | Format code |
| `make vet` | Static analysis |
| `make staticcheck` | Advanced static analysis |
| `make gosec` | Security scanning |
| `make ci-local` | Run ALL CI checks |

### Commits & PRs
| Command | Purpose |
|---------|---------|
| `make ai-commit FILE=...` | Create AI-attributed commit |
| `make review-commit` | Review staged changes |
| `make pre-pr` | Validate before PR |

### Task Management
| Command | Purpose |
|---------|---------|
| `make new-feature TASK="x"` | Create feature task |
| `make new-bug BUG="x"` | Create bug report |

### Component Lookup
| Command | Purpose |
|---------|---------|
| `make what-to-use NEED="x"` | Lookup component (table, form, modal, color...) |

---

## Component Patterns (Enforced)

| Need | Use | Not |
|------|-----|-----|
| Table | `behaviors.TableBehavior[T]` | `table.New()` |
| Form in intent | `screens/*FormScreen` | `models.*Form` (DEPRECATED) |
| Form primitives | `forms.NewInput()`, `forms.NewSelect()` | Direct `huh.NewInput()` |
| Form state check | `forms.IsCompleted(f)` | `f.State == huh.StateCompleted` |
| Text/titles | `primitives.Title()`, `primitives.Body()` | Raw lipgloss |
| Badges | `primitives.HelpKeyBadge()` | `components.KeyBadge` |
| Colors | `theme.Primary()` etc | `lipgloss.Color("#xxx")` |
| Layout | `layout.ScreenLayout` | Manual composition |
| Modals | `feedback.Modal`, `behaviors.RenderModalOverlay()` | Custom modal code |
| Intent | Embed `*BaseIntent` | Custom base |

Run `make what-to-use NEED="keyword"` for detailed usage and examples.

---

## When to Refuse

The AI agent MUST refuse if asked to:

### Workflow Violations
- Skip `make session-start`
- Commit directly to `next` or `main` (always use feature branches)
- Write implementation before test (TDD violation)
- Make multiple changes per commit
- Use `git commit` directly (must use `make ai-commit`)
- Skip compliance checks (`make check-compliance`, `make check-intent-architecture`)
- Create PR targeting `main` (must target `next`)

### Comment Violations
- Add comments inside function bodies (extract to named methods instead)
- Add inline comments at end of lines (except in e2e test files)
- Add section divider comments within functions
- Add field-level inline comments in structs
- Use forbidden markers (TODO, FIXME, HACK, XXX, NOTE, IMPORTANT)
- Explain WHAT code does (use better names instead)
- Add comments when explicit code would be clearer

**When user requests a comment**: Suggest refactoring instead (extract method, rename variable, add godoc).

### Architecture Violations (AUTOMATED ENFORCEMENT)

These violations are **automatically detected** and will **block commits**:

**Run `make check-intent-architecture` before every commit.**

#### 1. Untyped State Fields
```go
// ❌ REFUSE THIS
type MyIntent struct {
    state string  // WRONG: Raw string
}

// ✅ REQUIRE THIS
type MyIntentState string
const (
    StateList MyIntentState = "list"
)
type MyIntent struct {
    state MyIntentState  // CORRECT: Typed enum
}
```

#### 2. Wrapped State Models
```go
// ❌ REFUSE THIS
type MyIntent struct {
    state *MyIntentModel  // WRONG: Wrapped model
}

// ✅ REQUIRE THIS
type MyIntent struct {
    // State fields flattened directly
    state      MyIntentState
    active     bool
    items      []*Item
    selected   *Item
}
```

#### 3. context.Background() in Intents
```go
// ❌ REFUSE THIS
ctx := context.Background()

// ✅ REQUIRE THIS
ctx := i.getContext()
```

#### 4. Dead Code Markers
```go
// ❌ REFUSE THIS
case *SomeScreen:
    // This case should not be reached
    return screen.View()

// ✅ REQUIRE THIS
// Remove the unreachable code OR create cleanup task
```

#### 5. Missing *BaseIntent
```go
// ❌ REFUSE THIS
type MyIntent struct {
    // Missing BaseIntent
}

// ✅ REQUIRE THIS
type MyIntent struct {
    *BaseIntent  // REQUIRED
}
```

#### 6. Missing ScreenResultHandler
```go
// ❌ REFUSE THIS
type MyIntent struct {
    *BaseIntent
    activeScreen screens.Screen
}
// No ScreenResultHandler implementation

// ✅ REQUIRE THIS
var _ ScreenResultHandler = (*MyIntent)(nil)

func (i *MyIntent) HandleCancel(*screens.CancelResult) tea.Cmd { ... }
func (i *MyIntent) HandleNavigate(*screens.NavigateResult) tea.Cmd { ... }
func (i *MyIntent) HandleSubmit(*screens.SubmitResult) tea.Cmd { ... }
func (i *MyIntent) HandleError(*screens.ErrorResult) tea.Cmd { ... }
```

#### 7. Direct Modal Rendering (Not Using RenderModalOverlay)
```go
// ❌ REFUSE THIS
func (i *MyIntent) View() string {
    return i.modal.View()  // WRONG: Direct rendering
}

// ✅ REQUIRE THIS
func (i *MyIntent) View() string {
    baseView := i.screen.View()
    if i.modal != nil && i.modal.IsVisible() {
        return behaviors.RenderModalOverlay(i.modal, baseView)
    }
    return baseView
}
```

#### 8. Screens Importing Intents
```go
// ❌ REFUSE THIS (in screens/)
import "github.com/baphled/kariya/internal/cli/intents"

// ✅ REQUIRE THIS
// Screens NEVER import intents
// Use ScreenResult to communicate back to intent
```

#### 9. Direct huh Import Outside forms/
```go
// ❌ REFUSE THIS (in intents/, screens/, components/, models/)
import "github.com/charmbracelet/huh"

// ✅ REQUIRE THIS
import "github.com/baphled/kariya/internal/cli/forms"
import "github.com/baphled/kariya/internal/cli/models"
```

#### 10. Missing Context Field
```go
// ❌ REFUSE THIS
type MyIntent struct {
    *BaseIntent
    items []*Item  // Raw parameters - NO context struct
}

// ✅ REQUIRE THIS
type MyIntentContext struct {
    Items []*Item
}

type MyIntent struct {
    *BaseIntent
    context *MyIntentContext  // REQUIRED
}
```

#### 11. Missing State Field
```go
// ❌ REFUSE THIS
type MyIntent struct {
    *BaseIntent
    // Missing state field
}

// ✅ REQUIRE THIS
type MyIntentState string

const (
    StateList   MyIntentState = "list"
    StateDetail MyIntentState = "detail"
)

type MyIntent struct {
    *BaseIntent
    state MyIntentState  // REQUIRED
}
```

#### 12. Generic activeScreen Without Typed Fields
```go
// ❌ REFUSE THIS
type MyIntent struct {
    *BaseIntent
    activeScreen screens.Screen  // WRONG: No typed fields
}

// ✅ REQUIRE THIS
type MyIntent struct {
    *BaseIntent
    
    // Explicit typed fields (REQUIRED when using screens)
    listScreen   *myfeature.ListScreen
    detailScreen *myfeature.DetailScreen
    
    // Generic pointer
    activeScreen screens.Screen
}
```

#### 13. Business Logic in Update() Method
```go
// ❌ REFUSE THIS
func (i *MyIntent) Update(msg tea.Msg) tea.Cmd {
    // WRONG: Direct SQL queries
    rows, err := db.Query("SELECT * FROM...")
    
    // WRONG: Complex business logic
    for _, item := range items {
        // Complex processing...
    }
}

// ✅ REQUIRE THIS - Delegate to services
func (i *MyIntent) Update(msg tea.Msg) tea.Cmd {
    // Orchestrate, don't implement
    cmd, result := i.listScreen.Update(msg)
    
    if result != nil {
        return i.handleScreenResult(result)
    }
    
    return cmd
}
```

#### 14. All Types in One File (File Separation)

**CRITICAL**: Context (input params) and Models (state wrappers) MUST be in separate files from intent implementation.

```go
// ❌ REFUSE THIS - All in one file
// File: my_intent.go
package intents

type MyIntentContext struct { ... }  // WRONG: Context in intent file
type MyIntentModel struct { ... }    // WRONG: Model in intent file
type MyIntent struct { ... }

// ✅ REQUIRE THIS - Proper separation

// File: my.go (or my_context.go)
// Purpose: Input parameters (Events, Services, Config)
package intents

type MyIntentContext struct {
    Events  []*career.Event    // Input data
    Service *service.MyService // Dependencies
    Config  *MyConfig          // Configuration
}

// File: my_intent.go
// Purpose: Intent implementation ONLY
package intents

type MyIntent struct {
    *BaseIntent
    context *MyIntentContext  // Reference to context
    state   MyIntentState     // State fields (flattened, NOT wrapped)
    active  bool
    
    // Explicit screen fields
    listScreen   *myfeature.ListScreen
    detailScreen *myfeature.DetailScreen
}

// Screens: screens/myfeature/*.go
// Modals: components/*.go or uikit/feedback/*.go
```

**Key Points**:
- **Context** = Input parameters (events, services, config) → Separate file
- **Model** = State wrapper → Should be FLATTENED (preferred) or separate file
- **Intent** = Implementation (Update, View, handlers) → Main file
- **Screens** = UI components → `screens/` package
- **Modals** = Overlay components → `components/` or `uikit/feedback/`

#### 15. String-Based Key Handling (WARNING)
```go
// ⚠️ DISCOURAGED - String comparisons
if keyMsg.String() == "q" {
    return tea.Quit
}

// ✅ RECOMMENDED - Use HandleGlobalKeys
switch HandleGlobalKeys(keyMsg) {
case KeyQuit:
    return tea.Quit
case KeyHelp:
    i.helpModal.Toggle()
}
```

### Pattern Violations
- Use hardcoded colors/styles (must use theme system)
- Use raw `*huh.Form` in intents (must use wrapper models)
- Use `components.KeyBadge` (must use `primitives.HelpKeyBadge()`)
- Use `components.StandardView` (must use `layout.NewScreenLayout()`)
- Write code violating SOLID principles

### Enforcement Commands

Before committing, the AI agent **MUST** run:

```bash
# Architecture validation (REQUIRED)
make check-intent-architecture

# Full compliance check (REQUIRED)
make check-compliance

# Comprehensive linting (RECOMMENDED)
make golangci-lint
```

**These checks are automatically run by the pre-commit hook and will BLOCK commits with violations.**

---

**Refusal Template** (Use when violations detected):
```
I cannot proceed with this request.

ARCHITECTURE VIOLATION DETECTED

Violation: [Specific violation, e.g., "Untyped state field"]
Rule: [Rule description]
Automated Check: [Script that catches this, e.g., "check-intent-architecture.sh (Check #1)"]

Required correction:
1. [Specific fix needed]

This violation is AUTOMATICALLY ENFORCED by:
- Pre-commit hook: .git/hooks/pre-commit
- Linter: scripts/check-intent-architecture.sh
- CI/CD: GitHub Actions workflow

Run `make check-intent-architecture` to verify compliance.

See: docs/checklists/INTENT_DEVELOPMENT_CHECKLIST.md

This is non-negotiable for project compliance.
```

---

**When in doubt**: `make what-to-use NEED="keyword"` or `make check-patterns`

---

## Commit Message Format

```bash
# Create commit message file
cat > /tmp/commit.txt << 'EOF'
type(scope): description

Optional body explaining WHY

Optional footer (issue refs)
EOF

# Commit with AI attribution
make ai-commit FILE=/tmp/commit.txt
```

**Types**: `feat`, `fix`, `docs`, `refactor`, `test`, `chore`
**Scopes**: `domain`, `service`, `cli`, `intents`, `uikit`, `forms`

---

## Documentation Reference

### Development Guides (docs/development/)

| Topic | Document |
|-------|----------|
| Session protocol | [SESSION_PROTOCOL.md](docs/development/SESSION_PROTOCOL.md) |
| BDD workflow | [BDD_WORKFLOW.md](docs/development/BDD_WORKFLOW.md) |
| Architecture | [ARCHITECTURE_OVERVIEW.md](docs/development/ARCHITECTURE_OVERVIEW.md) |
| Development workflow | [DEVELOPMENT_WORKFLOW.md](docs/development/DEVELOPMENT_WORKFLOW.md) |
| Common tasks | [COMMON_TASKS.md](docs/development/COMMON_TASKS.md) |
| Intent patterns | [INTENT_PATTERNS_LIBRARY.md](docs/development/INTENT_PATTERNS_LIBRARY.md) |
| Navigation testing | [NAVIGATION_TESTING_GUIDE.md](docs/development/NAVIGATION_TESTING_GUIDE.md) |
| Keyboard system | [KEYBOARD_SYSTEM_GUIDE.md](docs/development/KEYBOARD_SYSTEM_GUIDE.md) |
| State transitions | [STATE_TRANSITION_PATTERNS.md](docs/development/STATE_TRANSITION_PATTERNS.md) |
| Modal overlays | [MODAL_OVERLAY_PATTERN.md](docs/development/MODAL_OVERLAY_PATTERN.md) |

### Rules (docs/rules/)

| Topic | Document |
|-------|----------|
| Master workflow | [master-task-prompt.md](docs/rules/master-task-prompt.md) |
| Senior engineer | [senior-engineer-guidelines.md](docs/rules/senior-engineer-guidelines.md) |
| Go guidelines | [go-guidelines.md](docs/rules/go-guidelines.md) |
| Atomic commits | [atomic-commits.md](docs/rules/atomic-commits.md) |
| AI attribution | [AI_COMMIT_ATTRIBUTION.md](docs/rules/AI_COMMIT_ATTRIBUTION.md) |
| Token efficiency | [token-efficiency.md](docs/rules/token-efficiency.md) |
| Compliance check | [rules-compliance-check.md](docs/rules/rules-compliance-check.md) |

### Quick References (docs/rules/)

| Topic | Document |
|-------|----------|
| Task workflow | [TASK_QUICK_REF.md](docs/rules/TASK_QUICK_REF.md) |
| Commit format | [COMMIT_QUICK_REFERENCE.md](docs/rules/COMMIT_QUICK_REFERENCE.md) |
| Compliance | [COMPLIANCE_QUICK_REF.md](docs/rules/COMPLIANCE_QUICK_REF.md) |
| AI commits | [AI_COMMIT_CHECKLIST.md](docs/rules/AI_COMMIT_CHECKLIST.md) |

### Component Guides

| Topic | Document |
|-------|----------|
| UIKit | [UIKIT_GUIDE.md](docs/UIKIT_GUIDE.md) |
| Forms | [FORMS_GUIDE.md](docs/FORMS_GUIDE.md) |
| Forms workflow | [FORMS_WORKFLOW_GUIDE.md](docs/rules/FORMS_WORKFLOW_GUIDE.md) |
| Modals | [MODAL_PATTERNS.md](docs/MODAL_PATTERNS.md) |
| Intent architecture | [INTENT_ARCHITECTURE_GUIDE.md](docs/INTENT_ARCHITECTURE_GUIDE.md) |
| Error handling | [ERROR_HANDLING_GUIDE.md](docs/guides/ERROR_HANDLING_GUIDE.md) |
| Screen naming | [SCREEN_NAMING.md](docs/conventions/SCREEN_NAMING.md) |
| Modal naming | [MODAL_NAMING.md](docs/conventions/MODAL_NAMING.md) |

### Troubleshooting

| Topic | Document |
|-------|----------|
| General | [TROUBLESHOOTING.md](docs/TROUBLESHOOTING.md) |
| CV issues | [CV_TROUBLESHOOTING.md](docs/guides/CV_TROUBLESHOOTING.md) |
| Error handling | [ERROR_HANDLING_GUIDE.md](docs/guides/ERROR_HANDLING_GUIDE.md) |

---

## AI Agent Configuration

The `ai-commit` script auto-detects the AI agent. Override with:

```bash
export AI_AGENT="Claude Code"
export AI_MODEL="Claude Sonnet 4"
```

Auto-detection checks (in order):
1. `AI_AGENT` environment variable
2. `CLAUDE_CODE` or `ANTHROPIC_API_KEY` -> Claude Code
3. `CURSOR_SESSION` or `CURSOR` -> Cursor
4. Parent process containing "claude" -> Claude Code
5. Default: Claude Code

---

## Project Structure

```
internal/cli/
├── intents/     # Workflows (state machines)
├── behaviors/   # Reusable behaviors (TableBehavior, CRUD)
├── components/  # LEGACY - migrate to uikit/
├── uikit/       # UIKit component library
│   ├── primitives/  # Text, Button, Badge, Input
│   ├── containers/  # Box, Overlay
│   ├── feedback/    # Modal, ModalContainer, HelpModal
│   ├── layout/      # ScreenLayout, Header, Footer
│   └── theme/       # Theme integration
├── models/      # Form wrappers
└── forms/       # Form configs
```

---

## Code Examples

See `examples/` directory:
- `intent_template.go.example` - Intent pattern
- `form_wrapper_template.go.example` - Form wrapper pattern
- `behavior_usage.go.example` - Behavior patterns

---

**When in doubt**: `make what-to-use NEED="keyword"` or `make check-patterns`

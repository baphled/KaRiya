# Intent Development Checklist

**Purpose**: Mandatory checklist for all intent development  
**Enforcement**: `make check-intent-architecture` (automated)  
**Status**: REQUIRED - Must pass before commit

---

## Pre-Commit Validation

Run this command **before every commit**:

```bash
make check-intent-architecture
```

This checklist is **automatically enforced** by the pre-commit hook. Violations will **block your commit**.

---

## Architecture Requirements (BLOCKING)

### 1. BaseIntent Embedding ✅

**Rule**: All intents MUST embed `*BaseIntent`

```go
type MyIntent struct {
    *BaseIntent  // REQUIRED
    
    // ... other fields
}
```

**Why**: Provides terminal info, theme, logo, and help modal infrastructure.

**Checked by**: `check-intent-architecture.sh` (Check #6)

---

### 2. Typed State Enum ✅

**Rule**: State fields MUST use typed enums, not raw strings

```go
// REQUIRED: Define typed state
type MyIntentState string

const (
    StateList   MyIntentState = "list"
    StateDetail MyIntentState = "detail"
)

type MyIntent struct {
    *BaseIntent
    state MyIntentState  // Typed, not string
}
```

**Why**: Type safety prevents typos and enables compile-time checks.

**Checked by**: `check-intent-architecture.sh` (Check #1)

---

### 3. Flattened State Model ✅

**Rule**: Intent state MUST be flattened directly into the intent struct, not wrapped in a separate model

```go
// ❌ BAD: Wrapped state
type MyIntent struct {
    *BaseIntent
    state *MyIntentModel  // WRONG
}

// ✅ GOOD: Flattened state
type MyIntent struct {
    *BaseIntent
    
    // State machine
    state      MyIntentState
    active     bool
    result     *IntentResult[*MyResult]
    
    // Screens
    listScreen   *myfeature.ListScreen
    detailScreen *myfeature.DetailScreen
    
    // Modals
    deleteModal *feedback.ConfirmModal
    
    // Data (flattened)
    items         []*Item
    selectedItem  *Item
    filters       *FilterConfig
}
```

**Why**: Reduces indirection, matches architecture guide patterns.

**Checked by**: `check-intent-architecture.sh` (Check #2)

---

### 4. Required State Field ✅

**Rule**: All intents MUST have a state field (typed enum)

```go
// ❌ BAD: No state field
type MyIntent struct {
    *BaseIntent
    // Missing state field
}

// ✅ GOOD: Typed state field
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

**Why**: 
- State machine pattern for orchestration
- Type-safe state transitions
- Clear intent lifecycle

**Checked by**: `check-intent-architecture.sh` (Check #10)

---

### 5. Required Context Field ✅

**Rule**: All intents MUST have a context field for input parameters

```go
// ❌ BAD: No context field
type MyIntent struct {
    *BaseIntent
    items []*Item  // Raw parameters
}

// ✅ GOOD: Context struct
type MyIntentContext struct {
    Items []*Item
    Mode  string
}

type MyIntent struct {
    *BaseIntent
    context *MyIntentContext
    
    // ... other fields
}
```

**Why**: 
- Centralizes input validation
- Makes intent reusable with different contexts
- Clear separation between input and state

**Checked by**: `check-intent-architecture.sh` (Check #14)

---

### 6. Explicit Screen Fields ✅

**Rule**: Intents using screens MUST declare explicit typed screen fields

```go
// ❌ BAD: Only generic activeScreen
type MyIntent struct {
    *BaseIntent
    activeScreen screens.Screen  // WRONG: No typed fields
}

// ✅ GOOD: Explicit typed fields
type MyIntent struct {
    *BaseIntent
    
    // Explicit screen fields (REQUIRED)
    listScreen   *myfeature.ListScreen
    detailScreen *myfeature.DetailScreen
    formScreen   *myfeature.FormScreen
    
    // Generic pointer (points to one of the above)
    activeScreen screens.Screen
}
```

**Why**:
- Makes intent structure clear
- Easy to see all possible screens
- Type-safe screen management
- Better code navigation

**Checked by**: `check-intent-architecture.sh` (Check #11)

---

### 7. Explicit Modal Fields ⚠️

**Rule**: Intents using modals SHOULD declare explicit typed modal fields

```go
// ⚠️ ACCEPTABLE but not ideal
type MyIntent struct {
    *BaseIntent
    // Modals used but not declared
}

// ✅ BETTER: Explicit modal fields
type MyIntent struct {
    *BaseIntent
    
    // Modals (explicit declaration)
    deleteModal   *feedback.ConfirmModal
    successModal  *feedback.SuccessModal
    errorModal    *feedback.ErrorModal
    filterModal   *components.FilterModal
}
```

**Why**:
- Makes intent structure clear
- Easy to see all possible modals
- Better code documentation

**Checked by**: `check-intent-architecture.sh` (Check #12)

---

### 8. Lean Intent Pattern ✅

**Rule**: Intents should orchestrate, NOT contain business logic

```go
// ❌ BAD: Business logic in Update
func (i *MyIntent) Update(msg tea.Msg) tea.Cmd {
    // WRONG: Direct SQL queries
    rows, err := db.Query("SELECT * FROM items WHERE...")
    
    // WRONG: Complex computations
    for _, item := range items {
        for _, tag := range item.Tags {
            // Complex processing...
        }
    }
    
    return nil
}

// ✅ GOOD: Delegate to services
func (i *MyIntent) Update(msg tea.Msg) tea.Cmd {
    // Delegate to screen
    cmd, result := i.listScreen.Update(msg)
    
    if result != nil {
        return i.handleScreenResult(result)
    }
    
    return cmd
}

// Business logic goes in service layer
func (i *MyIntent) handleScreenResult(result screens.ScreenResult) tea.Cmd {
    if submit, ok := result.(*screens.SubmitResult); ok {
        // Delegate to service
        ctx := i.getContext()
        err := i.service.SaveItem(ctx, submit.Data)
        // ...
    }
}
```

**Why**:
- Separation of concerns
- Testable business logic
- Reusable services
- Maintainable intents

**Checked by**: `check-intent-architecture.sh` (Check #9)

---

### 9. Context Usage ✅

**Rule**: Use `i.getContext()` instead of `context.Background()`

```go
// ❌ BAD
ctx := context.Background()

// ✅ GOOD
ctx := i.getContext()
```

**Why**: Enables proper cancellation and timeout support.

**Checked by**: `check-intent-architecture.sh` (Check #3)

---

### 10. No Dead Code ✅

**Rule**: Remove code marked "should not be reached" or create cleanup task

```go
// ❌ BAD
case *timeline.TimelineEventDetailScreen:
    // LEGACY: This case should not be reached
    return screen.View()
```

**Action**: Either remove the code or create a task: `make new-bug BUG="Remove legacy code"`

**Checked by**: `check-intent-architecture.sh` (Check #4)

---

### 11. ScreenResultHandler Implementation ✅

**Rule**: Intents using screens MUST implement `ScreenResultHandler`

```go
// REQUIRED: Interface compliance check
var _ ScreenResultHandler = (*MyIntent)(nil)

// Implement all methods:
func (i *MyIntent) HandleCancel(*screens.CancelResult) tea.Cmd
func (i *MyIntent) HandleNavigate(*screens.NavigateResult) tea.Cmd
func (i *MyIntent) HandleSubmit(*screens.SubmitResult) tea.Cmd
func (i *MyIntent) HandleError(*screens.ErrorResult) tea.Cmd
```

**Checked by**: `check-intent-architecture.sh` (Check #9)

---

### 12. Modal Overlay Pattern ✅

**Rule**: Use `behaviors.RenderModalOverlay()` for modal rendering

```go
// ✅ GOOD
func (i *MyIntent) View() string {
    baseView := i.listScreen.View()
    
    if i.deleteModal != nil && i.deleteModal.IsVisible() {
        return behaviors.RenderModalOverlay(i.deleteModal, baseView)
    }
    
    return baseView
}
```

**Checked by**: `check-intent-architecture.sh` (Check #8)

---

## Code Quality (WARNINGS)

### 13. File Separation ✅

**Rule**: Context, Model, Screen, and Modal structs MUST be in separate files

```go
// ❌ VIOLATION - All in one file: my_intent.go
package intents

type MyIntentContext struct { ... }  // WRONG: In intent file
type MyIntentModel struct { ... }    // WRONG: In intent file
type MyListScreen struct { ... }     // WRONG: In intent file
type MyModal struct { ... }          // WRONG: In intent file

type MyIntent struct { ... }

// ✅ REQUIRED - Proper file separation:

// File: my_context.go
type MyIntentContext struct { ... }

// File: my_model.go (if needed, but prefer flattening)
type MyIntentModel struct { ... }

// File: screens/myfeature/list_screen.go
type ListScreen struct { ... }

// File: components/my_modal.go or uikit/feedback/my_modal.go
type MyModal struct { ... }

// File: my_intent.go
type MyIntent struct {
    context *MyIntentContext
    // Flattened model fields here
}
```

**File Structure**:
```
internal/cli/
├── intents/
│   ├── my_intent.go          # Intent implementation only
│   ├── my_context.go          # Context struct
│   └── my_model.go            # Model struct (if not flattened)
├── screens/
│   └── myfeature/
│       ├── list_screen.go     # Screen implementations
│       └── detail_screen.go
└── components/
    └── my_modal.go            # Modal implementations
```

**Why**:
- Separation of concerns
- Single responsibility per file
- Easier to navigate and maintain
- Clear boundaries between types

**Checked by**: `check-intent-architecture.sh` (Check #16)

---

### 14. Key Handling Pattern ⚠️

**Rule**: Use `HandleGlobalKeys()` for common key handling, avoid string comparisons

```go
// ⚠️ ACCEPTABLE but not ideal
if keyMsg.String() == "q" {
    return tea.Quit
}

// ✅ BETTER - Use HandleGlobalKeys
switch HandleGlobalKeys(keyMsg) {
case KeyQuit:
    return tea.Quit
case KeyHelp:
    i.helpModal.Toggle()
    return nil
}
```

**Why**: 
- Centralized key handling logic
- Easier to maintain and test
- Consistent behavior across intents
- Type-safe (not string-based)

**Checked by**: `check-intent-architecture.sh` (Check #10, #12)

---

### 15. Godoc Completeness ⚠️

**Rule**: All exported functions should have godoc comments

```go
// HandleCancel handles screen cancellation (back/escape).
//
// Implements ScreenResultHandler interface.
func (i *MyIntent) HandleCancel(result *screens.CancelResult) tea.Cmd {
    // ...
}
```

**Why**: Documentation is code. Complete sentences required.

**Checked by**: `check-intent-architecture.sh` (Check #7), `golangci-lint`

---

### 16. Screen Management Pattern ⚠️

**Recommendation**: Declare typed screen fields for clarity

```go
// ⚠️ ACCEPTABLE but not ideal
type MyIntent struct {
    *BaseIntent
    activeScreen screens.Screen  // Generic
}

// ✅ BETTER
type MyIntent struct {
    *BaseIntent
    listScreen   *myfeature.ListScreen
    detailScreen *myfeature.DetailScreen
    activeScreen screens.Screen  // Points to one of the above
}
```

**Checked by**: `check-intent-architecture.sh` (Check #5)

---

## Testing Requirements

### E2E Tests

**Rule**: All intents MUST have E2E tests with Happy and Sad paths

```bash
# Check file exists
ls internal/cli/intents/my_intent_e2e_test.go

# Must contain:
grep "Happy Paths" internal/cli/intents/my_intent_e2e_test.go
grep "Sad Paths" internal/cli/intents/my_intent_e2e_test.go
```

**Checked by**: `scripts/check-patterns-strict.sh` (Check #7)

---

### Coverage Requirement

**Rule**: Package coverage ≥ 95%

```bash
make coverage
```

**Checked by**: `scripts/check-patterns-strict.sh` (Check #9)

---

## Architectural Constraints

### Layer Dependencies

**NEVER** violate these rules (enforced by `.golangci.yml`):

| Layer | Can Import | NEVER Import |
|-------|------------|--------------|
| `intents/` | screens, uikit, behaviors, components | - |
| `screens/` | uikit, behaviors | **intents** (FORBIDDEN) |
| `uikit/` | themes only | **screens, intents** (FORBIDDEN) |
| `behaviors/` | uikit, themes | **screens, intents** (FORBIDDEN) |

**Checked by**: `golangci-lint` (depguard rules)

---

### Forms Architecture

**Rule**: `huh` can ONLY be imported in `forms/` package

```go
// ❌ NEVER in intents
import "github.com/charmbracelet/huh"

// ✅ ALWAYS use wrappers
import "github.com/baphled/kariya/internal/cli/models"
import "github.com/baphled/kariya/internal/cli/forms"
```

**Checked by**: `golangci-lint` (depguard rules)

---

## Documentation Requirements

### New Intents

When creating a new intent, update:

1. **Workflow documentation** (`docs/workflows/`)
2. **Intent patterns** (`docs/development/INTENT_PATTERNS_LIBRARY.md`)
3. **Architecture guide** (`docs/INTENT_ARCHITECTURE_GUIDE.md` - add to reference list)

**Checked by**: `scripts/check-patterns-strict.sh` (Check #8)

---

## Quick Reference Commands

| Command | Purpose |
|---------|---------|
| `make check-intent-architecture` | Run all architecture checks |
| `make check-compliance` | Full compliance (includes architecture) |
| `make golangci-lint` | Comprehensive linting |
| `make check-patterns-strict` | Strict pattern validation |
| `make pre-commit` | Quick pre-commit checks |

---

## Example: Compliant Intent Structure

```go
package intents

import (
    "github.com/baphled/kariya/internal/cli/behaviors"
    "github.com/baphled/kariya/internal/cli/screens"
    "github.com/baphled/kariya/internal/cli/screens/myfeature"
    "github.com/baphled/kariya/internal/cli/uikit/feedback"
    tea "github.com/charmbracelet/bubbletea"
)

// MyIntentState defines the state machine states.
type MyIntentState string

const (
    StateList   MyIntentState = "list"
    StateDetail MyIntentState = "detail"
)

// Ensure interface compliance.
var _ ScreenResultHandler = (*MyIntent)(nil)

// MyIntent implements the Intent interface for [description].
type MyIntent struct {
    *BaseIntent
    
    // State machine
    state  MyIntentState
    active bool
    result *IntentResult[*MyResult]
    
    // Screens
    listScreen   *myfeature.ListScreen
    detailScreen *myfeature.DetailScreen
    
    // Modals
    deleteModal *feedback.ConfirmModal
    
    // Data
    items         []*Item
    selectedItem  *Item
}

// NewMyIntent creates a new MyIntent.
func NewMyIntent(ctx context.Context, items []*Item) (*MyIntent, error) {
    return &MyIntent{
        BaseIntent: NewBaseIntent(),
        state:      StateList,
        active:     true,
        items:      items,
    }, nil
}

// Init initializes the intent.
func (i *MyIntent) Init() tea.Cmd {
    i.listScreen = myfeature.NewListScreen(i.items)
    i.listScreen.SetTerminalInfo(i.GetTerminalInfo())
    i.listScreen.SetTheme(i.Theme())
    return nil
}

// Update processes messages.
func (i *MyIntent) Update(msg tea.Msg) tea.Cmd {
    if !i.active {
        return nil
    }
    
    // Modal handling
    if i.deleteModal != nil && i.deleteModal.IsVisible() {
        model, cmd := i.deleteModal.Update(msg)
        i.deleteModal = model.(*feedback.ConfirmModal)
        return cmd
    }
    
    // Screen delegation
    switch i.state {
    case StateList:
        cmd, result := i.listScreen.Update(msg)
        if result != nil {
            return i.handleScreenResult(result)
        }
        return cmd
    }
    
    return nil
}

// View renders the intent.
func (i *MyIntent) View() string {
    var baseView string
    
    switch i.state {
    case StateList:
        baseView = i.listScreen.View()
    case StateDetail:
        baseView = i.detailScreen.View()
    }
    
    // Overlay modals
    if i.deleteModal != nil && i.deleteModal.IsVisible() {
        return behaviors.RenderModalOverlay(i.deleteModal, baseView)
    }
    
    return baseView
}

// Result returns the intent result.
func (i *MyIntent) Result() *IntentResult[interface{}] {
    // Type erasure for interface compliance
    if i.result == nil {
        return nil
    }
    return &IntentResult[interface{}]{
        Status: i.result.Status,
        Data:   i.result.Data,
    }
}

// HandleCancel handles screen cancellation.
func (i *MyIntent) HandleCancel(result *screens.CancelResult) tea.Cmd {
    i.setCancelled()
    return nil
}

// HandleNavigate handles screen navigation.
func (i *MyIntent) HandleNavigate(result *screens.NavigateResult) tea.Cmd {
    // State transitions here
    return nil
}

// HandleSubmit handles form submission.
func (i *MyIntent) HandleSubmit(result *screens.SubmitResult) tea.Cmd {
    return nil
}

// HandleError handles errors from screens.
func (i *MyIntent) HandleError(result *screens.ErrorResult) tea.Cmd {
    return nil
}

// getContext returns context for service calls.
func (i *MyIntent) getContext() context.Context {
    // Use proper context with cancellation
    return i.BaseIntent.GetContext()
}

// setCancelled marks the intent as cancelled.
func (i *MyIntent) setCancelled() {
    i.result = &IntentResult[*MyResult]{
        Status: Cancelled,
    }
    i.active = false
}
```

---

## File Organization (STRICT STANDARD)

### 17. Subdirectory Structure ✅

**Rule**: New intents MUST use subdirectory structure with 5 files

**Required Structure**:
```
intents/{feature}/
├── context.go     # Business logic, data (100-200 lines)
├── result.go      # Output type (20-50 lines)
├── constants.go   # State enum, error constants (50-80 lines)
├── messages.go    # ALL *Msg types (50-100 lines)
└── intent.go      # Broker ONLY (200-400 lines, MAX 600)

screens/{feature}/
├── list_screen.go    # List view
├── detail_screen.go  # Detail view
└── form_screen.go    # Form view
```

**Legacy Intents**: 7 intents still using flat structure - marked for migration

**Why**:
- File organization standard
- Enforced separation of concerns
- Intent.go stays lean (broker only)
- Clear type boundaries

**Checked by**: `check-intent-architecture.sh` (Check #17)

---

### 18. Intent File Size Limit ✅

**Rule**: Intent.go MUST NOT exceed 600 lines (target 200-400 lines)

```bash
# Warnings at >400 lines
# BLOCKS commits at >600 lines
```

**Action for >600 lines**:
1. REFUSE all feature additions
2. Only allow critical bug fixes
3. Apply Boy Scout Rule - refactor touched areas
4. Migrate to subdirectory structure

**Why**:
- Intent is broker, not implementation
- Forces extraction of business logic → context.go
- Forces extraction of views → screens/
- Forces centralized modals → uikit/feedback/

**Checked by**: `check-intent-architecture.sh` (Check #18)

**Current Violators** (5 intents >1,000 lines):
- `generate_cv_intent.go` - 2,746 lines
- `manage_skills_intent.go` - 2,415 lines
- `capture_event_intent.go` - 1,884 lines
- `burst_management_intent.go` - 1,488 lines
- `browse_timeline_intent.go` - 1,195 lines

---

### 19. No Rendering in Intent ✅

**Rule**: Intent MUST only have View() that delegates to screens (max 2 render methods)

```go
// ❌ BAD: Rendering in intent
func (i *MyIntent) View() string {
    // WRONG: 200 lines of lipgloss styling
    title := lipgloss.NewStyle().Foreground(...)
    // ... massive rendering code
}

func (i *MyIntent) renderHeader() string { ... }
func (i *MyIntent) renderBody() string { ... }
func (i *MyIntent) renderFooter() string { ... }
func (i *MyIntent) renderModal() string { ... }
// 8+ render methods = VIOLATION

// ✅ GOOD: Delegate to screens
func (i *MyIntent) View() string {
    baseView := i.activeScreen.View()
    
    if i.modal != nil && i.modal.IsVisible() {
        return behaviors.RenderModalOverlay(i.modal, baseView)
    }
    
    return baseView
}
// Only 1 method, delegates rendering
```

**Why**:
- Screens handle rendering, not intents
- Reduces intent file size
- Reusable screen components
- Clear separation of concerns

**Checked by**: `check-intent-architecture.sh` (Check #19)

---

### 20. Type Location Enforcement ✅

**Rule**: Types MUST be in correct files

| Type | File | Violation |
|------|------|-----------|
| `*Context struct` | `context.go` | 🔴 In intent.go |
| `*Result struct` | `result.go` | 🔴 In intent.go |
| `*State string` | `constants.go` | 🔴 In intent.go |
| `const (...)` | `constants.go` | 🔴 In intent.go |
| `var Err...` | `constants.go` | 🔴 In intent.go |
| `*Msg struct` | `messages.go` | 🔴 In intent.go or constants.go |
| `*Intent struct` | `intent.go` | ✅ Correct location |

**Example Structure**:
```go
// File: intents/myfeature/context.go
type MyFeatureContext struct {
    Events  []*Event
    Service *MyService
}

// File: intents/myfeature/result.go
type MyFeatureResult struct {
    Item *Item
}

// File: intents/myfeature/constants.go
type MyFeatureState string

const (
    StateList   MyFeatureState = "list"
    StateDetail MyFeatureState = "detail"
)

var ErrNoItems = errors.New("no items found")

// File: intents/myfeature/messages.go
type ItemSelectedMsg struct {
    Item *Item
}

type FilterAppliedMsg struct {
    Filters *FilterConfig
}

// File: intents/myfeature/intent.go
type MyFeatureIntent struct {
    *BaseIntent
    context *MyFeatureContext
    state   MyFeatureState
    // Implementation only
}
```

**Why**:
- Clear file responsibilities
- Single responsibility per file
- Easy navigation
- Enforced separation

**Checked by**: `check-intent-architecture.sh` (Check #20)

---

### 21. UIKit Component Usage ✅

**Rule**: Use UIKit components, not deprecated `components/` package

**Violations (BLOCKING)**:
```go
// ❌ NEVER use these
import "github.com/baphled/kariya/internal/cli/components"

components.KeyBadge("key", "label")      // WRONG
components.StandardView{...}             // WRONG

// ✅ ALWAYS use these
import "github.com/baphled/kariya/internal/cli/uikit/primitives"
import "github.com/baphled/kariya/internal/cli/uikit/layout"

primitives.HelpKeyBadge("key", "label", theme)  // CORRECT
layout.NewScreenLayout(theme).WithContent(...)  // CORRECT
```

**Warnings (TECHNICAL DEBT)**:
```go
// ⚠️ DISCOURAGED: Excessive raw lipgloss
lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))

// ✅ RECOMMENDED: Use theme system
style := lipgloss.NewStyle().Foreground(theme.Error())

// ✅ OR UIKit primitives
primitives.ErrorText("message", theme)
```

**Why**:
- `components/` is LEGACY, being phased out
- UIKit provides consistent, reusable components
- Theme integration built-in
- Reduces duplicate styling code

**Checked by**: `check-intent-architecture.sh` (Check #21)

**See Also**: [UIKit Guide](../UIKIT_GUIDE.md)

---

### 22. Deprecated Models Package for Forms ✅

**Rule**: NEVER use `models/` package for forms - use `screens/*FormScreen` instead

**Violations (BLOCKING)**:
```go
// ❌ NEVER use these (DEPRECATED)
import "github.com/baphled/kariya/internal/cli/models"

type MyIntent struct {
    form *models.CaptureForm  // DEPRECATED
}

formModel := models.NewCaptureForm(service)  // DEPRECATED

// ✅ ALWAYS use these
import "github.com/baphled/kariya/internal/cli/screens/myfeature"

type MyIntent struct {
    formScreen *myfeature.FormScreen  // CORRECT
}

i.formScreen = myfeature.NewFormScreen()  // CORRECT
```

**Architecture Flow**:
```
❌ WRONG (DEPRECATED):
intents → models.*Form → forms package → huh

✅ CORRECT:
intents → screens/*FormScreen → forms package → huh
```

**Why**:
- `models/` package for forms is DEPRECATED
- Cleaner architecture with fewer layers
- Consistent with screen pattern (list, detail, form)
- Forms are just another type of screen
- Easier to test and maintain

**Migration Example**:
```go
// BEFORE (DEPRECATED)
type CaptureEventIntent struct {
    *BaseIntent
    form *models.CaptureForm  // WRONG
}

func (i *CaptureEventIntent) Init() tea.Cmd {
    i.form = models.NewCaptureForm(i.service)
    return i.form.Init()
}

func (i *CaptureEventIntent) Update(msg tea.Msg) tea.Cmd {
    var cmd tea.Cmd
    i.form, cmd = i.form.Update(msg)
    
    if i.form.Submitted() {
        // handle submission
    }
    return cmd
}

func (i *CaptureEventIntent) View() string {
    return i.form.View()
}

// AFTER (CORRECT)
type CaptureEventIntent struct {
    *BaseIntent
    formScreen *capture.FormScreen  // CORRECT
}

func (i *CaptureEventIntent) Init() tea.Cmd {
    i.formScreen = capture.NewFormScreen(i.service)
    i.formScreen.SetTerminalInfo(i.GetTerminalInfo())
    i.formScreen.SetTheme(i.Theme())
    i.activeScreen = i.formScreen
    return nil
}

func (i *CaptureEventIntent) Update(msg tea.Msg) tea.Cmd {
    cmd, result := i.formScreen.Update(msg)
    
    if result != nil {
        return i.handleScreenResult(result)
    }
    return cmd
}

func (i *CaptureEventIntent) View() string {
    return i.activeScreen.View()
}
```

**Checked by**: `check-intent-architecture.sh` (Check #22)

**See Also**:
- [Forms Guide](../FORMS_GUIDE.md)
- [Forms Workflow](../rules/FORMS_WORKFLOW_GUIDE.md)
- [Screen Extraction Guide](../guides/SCREEN_EXTRACTION_GUIDE.md)

---

## Enforcement Summary

| Check | Type | Severity | Script |
|-------|------|----------|--------|
| BaseIntent embedding | Automated | 🔴 BLOCKING | check-intent-architecture.sh (#5) |
| Typed state enum | Automated | 🔴 BLOCKING | check-intent-architecture.sh (#1) |
| Flattened state | Automated | 🔴 BLOCKING | check-intent-architecture.sh (#2) |
| **Required state field** | **Automated** | **🔴 BLOCKING** | **check-intent-architecture.sh (#10)** |
| **Required context field** | **Automated** | **🔴 BLOCKING** | **check-intent-architecture.sh (#14)** |
| Context usage (getContext) | Automated | 🔴 BLOCKING | check-intent-architecture.sh (#3) |
| Dead code | Automated | 🔴 BLOCKING | check-intent-architecture.sh (#4) |
| **Lean intent pattern** | **Automated** | **🔴 BLOCKING** | **check-intent-architecture.sh (#9)** |
| ScreenResultHandler | Automated | 🔴 BLOCKING | check-intent-architecture.sh (#8) |
| Modal overlay | Automated | 🔴 BLOCKING | check-intent-architecture.sh (#7) |
| **Explicit screen fields** | **Automated** | **🔴 BLOCKING** | **check-intent-architecture.sh (#11)** |
| **File separation** | **Automated** | **🔴 BLOCKING** | **check-intent-architecture.sh (#16)** |
| **Subdirectory structure** | **Automated** | **🟡 WARNING** | **check-intent-architecture.sh (#17)** |
| **Intent file size limit** | **Automated** | **🔴 BLOCKING (>600 lines)** | **check-intent-architecture.sh (#18)** |
| **No rendering in intent** | **Automated** | **🔴 BLOCKING** | **check-intent-architecture.sh (#19)** |
| **Type location** | **Automated** | **🔴 BLOCKING** | **check-intent-architecture.sh (#20)** |
| **UIKit component usage** | **Automated** | **🔴 BLOCKING + 🟡 WARNING** | **check-intent-architecture.sh (#21)** |
| **Deprecated models/ for forms** | **Automated** | **🔴 BLOCKING** | **check-intent-architecture.sh (#22)** |
| Layer dependencies | Automated | 🔴 BLOCKING | golangci-lint (depguard) |
| Forms architecture | Automated | 🔴 BLOCKING | golangci-lint (depguard) |
| Explicit modal fields | Automated | 🟡 WARNING | check-intent-architecture.sh (#12) |
| Key handling pattern | Automated | 🟡 WARNING | check-intent-architecture.sh (#13) |
| HandleGlobalKeys usage | Automated | 🟡 WARNING | check-intent-architecture.sh (#15) |
| Godoc completeness | Automated | 🟡 WARNING | check-intent-architecture.sh (#6) |
| E2E tests | Automated | 🔴 BLOCKING | check-patterns-strict.sh |
| Coverage ≥95% | Automated | 🔴 BLOCKING | check-patterns-strict.sh |

---

**See Also**:
- [Intent Architecture Guide](../INTENT_ARCHITECTURE_GUIDE.md) - Full architecture documentation
- [Intent Patterns Library](../development/INTENT_PATTERNS_LIBRARY.md) - Pattern examples
- [AGENTS.md](../../AGENTS.md) - AI agent refusal rules

---

**Last Updated**: 2026-01-22  
**Status**: Production - Enforced by pre-commit hook

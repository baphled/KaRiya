# Intent Migration to Subdirectory Structure

**Purpose**: Step-by-step guide for migrating flat-structure intents to subdirectory organization  
**Status**: REQUIRED for all legacy intents (7 identified, 5 critical >1,000 lines)  
**Enforcement**: Check #17-20 in `check-intent-architecture.sh`

---

## Overview

### Why Migrate?

Legacy intents are **massively bloated**:
- `generate_cv_intent.go` - **2,746 lines** (should be ~300)
- `manage_skills_intent.go` - **2,415 lines** (should be ~300)
- `capture_event_intent.go` - **1,884 lines** (should be ~300)
- `burst_management_intent.go` - **1,488 lines** (should be ~300)
- `browse_timeline_intent.go` - **1,195 lines** (should be ~300)

**Problems**:
1. Violates single responsibility principle
2. Mixes business logic with orchestration
3. Inlines rendering instead of using screens
4. Hard to navigate and maintain
5. Impossible to test in isolation
6. Blocks new feature development (>600 line limit enforced)

**Solution**: Migrate to **subdirectory structure** with 5 focused files.

---

## Migration Strategy

### Phase 1: Create Subdirectory Structure

**Target Structure**:
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
└── form_screen.go    # Form view (if applicable)
```

### Phase 2: Extract Types

**Order of extraction** (follow this sequence):
1. **constants.go** - State enum, error constants (easiest)
2. **messages.go** - All `*Msg` types (second easiest)
3. **result.go** - Output type (simple)
4. **context.go** - Input parameters, business logic (moderate)
5. **intent.go** - Broker implementation (last, references above)

### Phase 3: Extract Views

**Extract rendering** from intent to screens:
1. Identify all `render*()` methods in intent
2. Create screen structs in `screens/{feature}/`
3. Move rendering logic to screen `View()` methods
4. Update intent to delegate: `return i.activeScreen.View()`

### Phase 4: Centralize Modals

**Replace custom modals** with centralized ones:
1. Identify custom modal structs (delete confirmation, success, error)
2. Replace with `uikit/feedback/` modals
3. Use `behaviors.RenderModalOverlay()` for rendering

---

## Step-by-Step Migration

### Example: Migrating `burst_management_intent.go`

**Current State**: 1,488 lines in single file  
**Target**: 5 files totaling ~600 lines + 3 screen files

---

### Step 1: Create Subdirectory

```bash
# Create subdirectory
mkdir -p internal/cli/intents/burst_management

# Create screen directory
mkdir -p internal/cli/screens/burst_management
```

---

### Step 2: Extract Constants (constants.go)

**Identify constants** in original file:
```go
// From burst_management_intent.go (BEFORE)
type BurstManagementState string

const (
    StateInitial     BurstManagementState = "initial"
    StateSelectMode  BurstManagementState = "select_mode"
    StateList        BurstManagementState = "list"
    StateCreateForm  BurstManagementState = "create_form"
    StateEditForm    BurstManagementState = "edit_form"
    StateConfirmDelete BurstManagementState = "confirm_delete"
)

var (
    ErrNoBurstFound = errors.New("no burst found")
    ErrInvalidMode  = errors.New("invalid mode")
)
```

**Create new file**: `intents/burst_management/constants.go`
```go
package burst_management

import "errors"

// BurstManagementState defines the state machine states.
type BurstManagementState string

const (
    // StateInitial is the initial state.
    StateInitial BurstManagementState = "initial"
    
    // StateSelectMode is the mode selection state.
    StateSelectMode BurstManagementState = "select_mode"
    
    // StateList is the list view state.
    StateList BurstManagementState = "list"
    
    // StateCreateForm is the create form state.
    StateCreateForm BurstManagementState = "create_form"
    
    // StateEditForm is the edit form state.
    StateEditForm BurstManagementState = "edit_form"
    
    // StateConfirmDelete is the delete confirmation state.
    StateConfirmDelete BurstManagementState = "confirm_delete"
)

// Error constants.
var (
    ErrNoBurstFound = errors.New("no burst found")
    ErrInvalidMode  = errors.New("invalid mode")
)
```

**Result**: ~35 lines (constants.go)

---

### Step 3: Extract Messages (messages.go)

**Identify message types** in original file:
```go
// From burst_management_intent.go (BEFORE)
type BurstSelectedMsg struct {
    Burst *domain.Burst
}

type BurstCreatedMsg struct {
    Burst *domain.Burst
}

type BurstDeletedMsg struct {
    ID string
}

type FilterAppliedMsg struct {
    Filters *FilterConfig
}
```

**Create new file**: `intents/burst_management/messages.go`
```go
package burst_management

import "github.com/baphled/kariya/internal/domain"

// BurstSelectedMsg is sent when a burst is selected.
type BurstSelectedMsg struct {
    Burst *domain.Burst
}

// BurstCreatedMsg is sent when a burst is created.
type BurstCreatedMsg struct {
    Burst *domain.Burst
}

// BurstDeletedMsg is sent when a burst is deleted.
type BurstDeletedMsg struct {
    ID string
}

// FilterAppliedMsg is sent when filters are applied.
type FilterAppliedMsg struct {
    Filters *FilterConfig
}

// FilterConfig holds filter configuration.
type FilterConfig struct {
    StartDate string
    EndDate   string
    Tags      []string
}
```

**Result**: ~45 lines (messages.go)

---

### Step 4: Extract Result (result.go)

**Identify result type** in original file:
```go
// From burst_management_intent.go (BEFORE)
type BurstManagementResult struct {
    Bursts  []*domain.Burst
    Created *domain.Burst
    Updated *domain.Burst
    Deleted string
}
```

**Create new file**: `intents/burst_management/result.go`
```go
package burst_management

import "github.com/baphled/kariya/internal/domain"

// BurstManagementResult is the output of the burst management intent.
type BurstManagementResult struct {
    // Bursts is the list of bursts (for list mode).
    Bursts []*domain.Burst
    
    // Created is the newly created burst.
    Created *domain.Burst
    
    // Updated is the updated burst.
    Updated *domain.Burst
    
    // Deleted is the ID of the deleted burst.
    Deleted string
}
```

**Result**: ~25 lines (result.go)

---

### Step 5: Extract Context (context.go)

**Identify input parameters** and business logic:
```go
// From burst_management_intent.go (BEFORE)
type BurstManagementIntent struct {
    *BaseIntent
    
    // Input parameters (should be in Context)
    mode         string
    burstService *service.BurstService
    
    // Business logic (should be in Context)
    bursts       []*domain.Burst
    selectedBurst *domain.Burst
    
    // ... state machine fields
}

// Business logic methods (should be in Context)
func (i *BurstManagementIntent) loadBursts() error { ... }
func (i *BurstManagementIntent) filterBursts(filters *FilterConfig) []*domain.Burst { ... }
func (i *BurstManagementIntent) validateBurst(burst *domain.Burst) error { ... }
```

**Create new file**: `intents/burst_management/context.go`
```go
package burst_management

import (
    "context"
    "github.com/baphled/kariya/internal/domain"
    "github.com/baphled/kariya/internal/service"
)

// BurstManagementContext holds input parameters and business logic.
type BurstManagementContext struct {
    // Input parameters
    Mode         string
    BurstService *service.BurstService
    
    // Data
    bursts       []*domain.Burst
    selectedBurst *domain.Burst
}

// NewBurstManagementContext creates a new context.
func NewBurstManagementContext(mode string, service *service.BurstService) *BurstManagementContext {
    return &BurstManagementContext{
        Mode:         mode,
        BurstService: service,
    }
}

// LoadBursts loads all bursts from the service.
func (c *BurstManagementContext) LoadBursts(ctx context.Context) error {
    bursts, err := c.BurstService.GetAll(ctx)
    if err != nil {
        return err
    }
    c.bursts = bursts
    return nil
}

// FilterBursts filters bursts by the given criteria.
func (c *BurstManagementContext) FilterBursts(filters *FilterConfig) []*domain.Burst {
    var filtered []*domain.Burst
    for _, burst := range c.bursts {
        if c.matchesFilter(burst, filters) {
            filtered = append(filtered, burst)
        }
    }
    return filtered
}

// ValidateBurst validates a burst.
func (c *BurstManagementContext) ValidateBurst(burst *domain.Burst) error {
    if burst.Title == "" {
        return errors.New("title is required")
    }
    // ... validation logic
    return nil
}

// GetBursts returns all bursts.
func (c *BurstManagementContext) GetBursts() []*domain.Burst {
    return c.bursts
}

// GetSelectedBurst returns the selected burst.
func (c *BurstManagementContext) GetSelectedBurst() *domain.Burst {
    return c.selectedBurst
}

// SetSelectedBurst sets the selected burst.
func (c *BurstManagementContext) SetSelectedBurst(burst *domain.Burst) {
    c.selectedBurst = burst
}

// Private helper
func (c *BurstManagementContext) matchesFilter(burst *domain.Burst, filters *FilterConfig) bool {
    // ... filter logic
    return true
}
```

**Result**: ~120 lines (context.go)

---

### Step 6: Extract Screens

**Identify rendering** in original intent:
```go
// From burst_management_intent.go (BEFORE)
func (i *BurstManagementIntent) View() string {
    // 200+ lines of rendering code
    switch i.state {
    case StateList:
        return i.renderList()  // WRONG: Should be in screen
    case StateCreateForm:
        return i.renderCreateForm()  // WRONG: Should be in screen
    }
}

func (i *BurstManagementIntent) renderList() string {
    // 100+ lines of lipgloss styling
    title := lipgloss.NewStyle().Foreground(...)
    body := lipgloss.NewStyle().Padding(...)
    // ...
}
```

**Create new file**: `screens/burst_management/list_screen.go`
```go
package burst_management

import (
    "github.com/baphled/kariya/internal/cli/behaviors"
    "github.com/baphled/kariya/internal/cli/screens"
    "github.com/baphled/kariya/internal/cli/screens/base"
    "github.com/baphled/kariya/internal/cli/uikit/layout"
    "github.com/baphled/kariya/internal/cli/uikit/primitives"
    "github.com/baphled/kariya/internal/domain"
    tea "github.com/charmbracelet/bubbletea"
)

// ListScreen displays a list of bursts.
type ListScreen struct {
    *base.BaseScreen
    table *behaviors.TableBehavior[*domain.Burst]
}

// NewListScreen creates a new list screen.
func NewListScreen(bursts []*domain.Burst) *ListScreen {
    columns := []behaviors.Column{
        {Title: "Title", Width: 30},
        {Title: "Start", Width: 20},
        {Title: "End", Width: 20},
        {Title: "Tags", Width: 30},
    }
    
    table := behaviors.NewTableBehavior(
        bursts,
        columns,
        func(b *domain.Burst) []string {
            return []string{
                b.Title,
                b.StartDate.Format("2006-01-02"),
                b.EndDate.Format("2006-01-02"),
                strings.Join(b.Tags, ", "),
            }
        },
    )
    
    return &ListScreen{
        BaseScreen: base.NewBaseScreen(),
        table:      table,
    }
}

// Update handles screen updates.
func (s *ListScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "enter":
            if selected := s.table.SelectedItem(); selected != nil {
                return nil, screens.NewNavigateResult("detail", selected)
            }
        case "esc":
            return nil, screens.NewCancelResult("list")
        }
    }
    
    // Delegate to table behavior
    cmd := s.table.Update(msg)
    return cmd, nil
}

// View renders the screen.
func (s *ListScreen) View() string {
    theme := s.GetTheme()
    
    header := primitives.Title("Burst Management", theme)
    tableView := s.table.View()
    help := primitives.HelpKeyBadge("enter", "select", theme) + " " +
            primitives.HelpKeyBadge("esc", "back", theme)
    
    return layout.NewScreenLayout(theme).
        WithHeader(header).
        WithContent(tableView).
        WithFooter(help).
        Render()
}
```

**Result**: ~90 lines (list_screen.go)

---

### Step 7: Refactor Intent (intent.go)

**Original**: 1,488 lines with everything  
**Target**: ~300 lines broker only

**Create new file**: `intents/burst_management/intent.go`
```go
package burst_management

import (
    "github.com/baphled/kariya/internal/cli/behaviors"
    "github.com/baphled/kariya/internal/cli/intents"
    "github.com/baphled/kariya/internal/cli/screens"
    "github.com/baphled/kariya/internal/cli/screens/burst_management"
    "github.com/baphled/kariya/internal/cli/uikit/feedback"
    tea "github.com/charmbracelet/bubbletea"
)

// Ensure interface compliance.
var _ intents.ScreenResultHandler = (*BurstManagementIntent)(nil)

// BurstManagementIntent orchestrates burst management workflow.
type BurstManagementIntent struct {
    *intents.BaseIntent
    
    // Context
    context *BurstManagementContext
    
    // State machine
    state  BurstManagementState
    active bool
    result *intents.IntentResult[*BurstManagementResult]
    
    // Screens (explicit typed fields)
    listScreen   *burst_management.ListScreen
    detailScreen *burst_management.DetailScreen
    formScreen   *burst_management.FormScreen
    
    // Active screen (generic pointer)
    activeScreen screens.Screen
    
    // Modals (centralized from uikit/feedback)
    deleteModal  *feedback.ConfirmModal
    errorModal   *feedback.ErrorModal
    successModal *feedback.SuccessModal
}

// NewBurstManagementIntent creates a new burst management intent.
func NewBurstManagementIntent(ctx *BurstManagementContext) (*BurstManagementIntent, error) {
    return &BurstManagementIntent{
        BaseIntent: intents.NewBaseIntent(),
        context:    ctx,
        state:      StateInitial,
        active:     true,
    }, nil
}

// Init initializes the intent.
func (i *BurstManagementIntent) Init() tea.Cmd {
    // Load bursts
    ctx := i.GetContext()
    if err := i.context.LoadBursts(ctx); err != nil {
        i.errorModal = feedback.NewErrorModal("Failed to load bursts: " + err.Error())
        i.errorModal.Show()
        return nil
    }
    
    // Initialize list screen
    i.listScreen = burst_management.NewListScreen(i.context.GetBursts())
    i.listScreen.SetTerminalInfo(i.GetTerminalInfo())
    i.listScreen.SetTheme(i.Theme())
    i.activeScreen = i.listScreen
    i.state = StateList
    
    return nil
}

// Update processes messages.
func (i *BurstManagementIntent) Update(msg tea.Msg) tea.Cmd {
    if !i.active {
        return nil
    }
    
    // Modal handling
    if i.deleteModal != nil && i.deleteModal.IsVisible() {
        model, cmd := i.deleteModal.Update(msg)
        i.deleteModal = model.(*feedback.ConfirmModal)
        
        if i.deleteModal.Confirmed() {
            return i.handleDeleteConfirm()
        }
        if i.deleteModal.Cancelled() {
            i.deleteModal.Hide()
        }
        
        return cmd
    }
    
    // Screen delegation
    cmd, result := i.activeScreen.Update(msg)
    
    if result != nil {
        return i.handleScreenResult(result)
    }
    
    return cmd
}

// View renders the intent.
func (i *BurstManagementIntent) View() string {
    baseView := i.activeScreen.View()
    
    // Overlay modals
    if i.deleteModal != nil && i.deleteModal.IsVisible() {
        return behaviors.RenderModalOverlay(i.deleteModal, baseView)
    }
    if i.errorModal != nil && i.errorModal.IsVisible() {
        return behaviors.RenderModalOverlay(i.errorModal, baseView)
    }
    if i.successModal != nil && i.successModal.IsVisible() {
        return behaviors.RenderModalOverlay(i.successModal, baseView)
    }
    
    return baseView
}

// Result returns the intent result.
func (i *BurstManagementIntent) Result() *intents.IntentResult[interface{}] {
    if i.result == nil {
        return nil
    }
    return &intents.IntentResult[interface{}]{
        Status: i.result.Status,
        Data:   i.result.Data,
    }
}

// HandleCancel handles screen cancellation.
func (i *BurstManagementIntent) HandleCancel(result *screens.CancelResult) tea.Cmd {
    i.setCancelled()
    return nil
}

// HandleNavigate handles screen navigation.
func (i *BurstManagementIntent) HandleNavigate(result *screens.NavigateResult) tea.Cmd {
    switch result.Target {
    case "detail":
        burst := result.Data.(*domain.Burst)
        i.context.SetSelectedBurst(burst)
        i.detailScreen = burst_management.NewDetailScreen(burst)
        i.detailScreen.SetTerminalInfo(i.GetTerminalInfo())
        i.detailScreen.SetTheme(i.Theme())
        i.activeScreen = i.detailScreen
        i.state = StateDetail
    }
    return nil
}

// HandleSubmit handles form submission.
func (i *BurstManagementIntent) HandleSubmit(result *screens.SubmitResult) tea.Cmd {
    // Delegate to service
    ctx := i.GetContext()
    burst := result.Data.(*domain.Burst)
    
    if err := i.context.BurstService.Save(ctx, burst); err != nil {
        i.errorModal = feedback.NewErrorModal("Failed to save burst: " + err.Error())
        i.errorModal.Show()
        return nil
    }
    
    i.successModal = feedback.NewSuccessModal("Burst saved successfully")
    i.successModal.Show()
    return nil
}

// HandleError handles errors from screens.
func (i *BurstManagementIntent) HandleError(result *screens.ErrorResult) tea.Cmd {
    i.errorModal = feedback.NewErrorModal(result.Error.Error())
    i.errorModal.Show()
    return nil
}

// Private helpers
func (i *BurstManagementIntent) handleDeleteConfirm() tea.Cmd {
    ctx := i.GetContext()
    burst := i.context.GetSelectedBurst()
    
    if err := i.context.BurstService.Delete(ctx, burst.ID); err != nil {
        i.errorModal = feedback.NewErrorModal("Failed to delete burst: " + err.Error())
        i.errorModal.Show()
        return nil
    }
    
    i.deleteModal.Hide()
    i.successModal = feedback.NewSuccessModal("Burst deleted successfully")
    i.successModal.Show()
    
    // Reload bursts
    i.context.LoadBursts(ctx)
    i.listScreen = burst_management.NewListScreen(i.context.GetBursts())
    i.listScreen.SetTerminalInfo(i.GetTerminalInfo())
    i.listScreen.SetTheme(i.Theme())
    i.activeScreen = i.listScreen
    i.state = StateList
    
    return nil
}

func (i *BurstManagementIntent) setCancelled() {
    i.result = &intents.IntentResult[*BurstManagementResult]{
        Status: intents.Cancelled,
    }
    i.active = false
}

func (i *BurstManagementIntent) handleScreenResult(result screens.ScreenResult) tea.Cmd {
    switch r := result.(type) {
    case *screens.CancelResult:
        return i.HandleCancel(r)
    case *screens.NavigateResult:
        return i.HandleNavigate(r)
    case *screens.SubmitResult:
        return i.HandleSubmit(r)
    case *screens.ErrorResult:
        return i.HandleError(r)
    }
    return nil
}
```

**Result**: ~290 lines (intent.go) - WITHIN TARGET!

---

## Migration Checklist

Use this checklist for each intent migration:

### Phase 1: Preparation
- [ ] Read original intent file completely
- [ ] Identify all types (State, Context, Result, Msg types)
- [ ] Identify all business logic methods
- [ ] Identify all rendering methods
- [ ] Identify all modals
- [ ] Create subdirectory: `intents/{feature}/`
- [ ] Create screen directory: `screens/{feature}/`

### Phase 2: Extract Types
- [ ] Create `constants.go` (State enum + error constants)
- [ ] Create `messages.go` (ALL *Msg types)
- [ ] Create `result.go` (Output type)
- [ ] Create `context.go` (Input params + business logic)
- [ ] Verify types compile independently

### Phase 3: Extract Screens
- [ ] Create `list_screen.go` (if list view exists)
- [ ] Create `detail_screen.go` (if detail view exists)
- [ ] Create `form_screen.go` (if form exists)
- [ ] Move rendering logic from intent to screens
- [ ] Use UIKit components (not raw lipgloss)
- [ ] Verify screens compile independently

### Phase 4: Extract Modals
- [ ] Identify custom modals
- [ ] Replace with `uikit/feedback/` modals:
  - Delete confirmation → `feedback.ConfirmModal`
  - Success message → `feedback.SuccessModal`
  - Error message → `feedback.ErrorModal`
  - Loading → `feedback.LoadingModal`
- [ ] Remove custom modal code

### Phase 5: Refactor Intent
- [ ] Create `intent.go` (broker only)
- [ ] Embed `*BaseIntent`
- [ ] Add context field
- [ ] Add typed screen fields
- [ ] Add typed modal fields
- [ ] Implement `Init()` (setup screens)
- [ ] Implement `Update()` (delegate to screens)
- [ ] Implement `View()` (delegate to screens, overlay modals)
- [ ] Implement `ScreenResultHandler` interface
- [ ] Verify intent.go < 600 lines

### Phase 6: Testing
- [ ] Run `make check-intent-architecture`
- [ ] Verify Check #17 passes (subdirectory structure)
- [ ] Verify Check #18 passes (file size)
- [ ] Verify Check #19 passes (no rendering in intent)
- [ ] Verify Check #20 passes (type location)
- [ ] Verify Check #21 passes (UIKit usage)
- [ ] Run unit tests
- [ ] Run E2E tests
- [ ] Manual testing

### Phase 7: Cleanup
- [ ] Delete original flat file
- [ ] Update imports in other files
- [ ] Run `make fmt`
- [ ] Run `make vet`
- [ ] Run `make golangci-lint`
- [ ] Commit with proper message

---

## Common Pitfalls

### 1. Forgetting to Update Package Name

```go
// ❌ WRONG - Old package name
package intents

// ✅ CORRECT - Subdirectory package name
package burst_management
```

### 2. Circular Dependencies

```go
// ❌ WRONG - Context imports intent
// File: context.go
import "github.com/baphled/kariya/internal/cli/intents/burst_management"

// ✅ CORRECT - Context is standalone
// File: context.go
import "github.com/baphled/kariya/internal/domain"
```

### 3. Business Logic in Intent

```go
// ❌ WRONG - Business logic in intent
func (i *BurstManagementIntent) Update(msg tea.Msg) tea.Cmd {
    // Direct SQL query
    rows, err := db.Query("SELECT * FROM bursts")
}

// ✅ CORRECT - Delegate to context
func (i *BurstManagementIntent) Update(msg tea.Msg) tea.Cmd {
    ctx := i.GetContext()
    i.context.LoadBursts(ctx)
}
```

### 4. Rendering in Intent

```go
// ❌ WRONG - Rendering in intent
func (i *BurstManagementIntent) View() string {
    title := lipgloss.NewStyle().Foreground(...)
    // 100+ lines of styling
}

// ✅ CORRECT - Delegate to screen
func (i *BurstManagementIntent) View() string {
    return i.activeScreen.View()
}
```

---

## Validation

After migration, run these commands:

```bash
# Architecture checks
make check-intent-architecture

# Full compliance
make check-compliance

# Comprehensive linting
make golangci-lint

# Tests
make test
make coverage
```

**All checks MUST pass** before the migration is considered complete.

---

## Example Migrations

### Priority Order (Largest First)

1. **generate_cv_intent.go** (2,746 lines) - CRITICAL
2. **manage_skills_intent.go** (2,415 lines) - CRITICAL
3. **capture_event_intent.go** (1,884 lines) - HIGH
4. **burst_management_intent.go** (1,488 lines) - HIGH
5. **browse_timeline_intent.go** (1,195 lines) - MEDIUM

**Estimate**: ~2-4 hours per intent migration (depending on complexity)

---

## Resources

### Templates
- `examples/intent_subdirectory_template/` - Complete template system

### Documentation
- [Intent Architecture Guide](../INTENT_ARCHITECTURE_GUIDE.md) - Architecture patterns
- [Intent Development Checklist](../checklists/INTENT_DEVELOPMENT_CHECKLIST.md) - Compliance checklist
- [Screen Extraction Guide](SCREEN_EXTRACTION_GUIDE.md) - Extracting views to screens
- [Helper Extraction Guide](HELPER_EXTRACTION_GUIDE.md) - Where helpers should live
- [AGENTS.md](../../AGENTS.md) - AI agent behavior rules

### Tools
- `make check-intent-architecture` - Architecture validation
- `make check-compliance` - Full compliance check
- `make what-to-use NEED="keyword"` - Component lookup

---

**Last Updated**: 2026-01-22  
**Status**: Required for all legacy intents  
**Enforcement**: Automated by pre-commit hook (Check #17-21)

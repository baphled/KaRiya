# Intent Subdirectory Template

**Purpose**: Standardized structure for creating new intents in KaRiya.

**Status**: REQUIRED for all new intents (enforced by Check #17)

---

## Structure Overview

```
intents/{feature}/
├── context.go    # Business logic, data management (100-200 lines)
├── result.go     # Output type (20-50 lines)
├── constants.go  # State enum, error constants (50-80 lines)
├── messages.go   # ALL *Msg types (50-100 lines)
└── intent.go     # Broker/orchestration ONLY (200-400 lines)

screens/{feature}/
├── list_screen.go    # List view
├── detail_screen.go  # Detail view
└── form_screen.go    # Form view
```

---

## Quick Start

### Step 1: Copy Template Files

```bash
# Set your feature name
FEATURE="myfeature"  # lowercase

# Create intent directory
mkdir -p internal/cli/intents/$FEATURE

# Copy template files
cp examples/intent_subdirectory_template/context.go.template \
   internal/cli/intents/$FEATURE/context.go
cp examples/intent_subdirectory_template/result.go.template \
   internal/cli/intents/$FEATURE/result.go
cp examples/intent_subdirectory_template/constants.go.template \
   internal/cli/intents/$FEATURE/constants.go
cp examples/intent_subdirectory_template/messages.go.template \
   internal/cli/intents/$FEATURE/messages.go
cp examples/intent_subdirectory_template/intent.go.template \
   internal/cli/intents/$FEATURE/intent.go
```

### Step 2: Replace Placeholders

Replace these in ALL files:
- `{feature}` → your feature name (lowercase, e.g., `burst_management`)
- `{Feature}` → PascalCase (e.g., `BurstManagement`)
- `{Entity}` → domain entity (e.g., `Burst`, `Skill`, `Event`)

### Step 3: Create Screens

```bash
# Create screen directory
mkdir -p internal/cli/screens/$FEATURE

# Copy screen template
cp examples/intent_subdirectory_template/screens/list_screen.go.template \
   internal/cli/screens/$FEATURE/list_screen.go

# Replace placeholders (same as above)
```

### Step 4: Verify Compliance

```bash
# Run architecture checks (should pass all)
make check-intent-architecture

# Expected results:
# ✅ Check #17: Subdirectory structure complete
# ✅ Check #18: Intent file size within limits
# ✅ Check #19: No rendering in intent
# ✅ Check #20: Types in correct locations
```

---

## File Responsibilities

### context.go (Business Logic - 100-200 lines)

**Contains**:
- `Context` struct with all state/data
- Constructor (`NewContext`)
- Business methods (`Load`, `Create`, `Update`, `Delete`)
- Filter/sort logic
- Validation

**Rules**:
- ALL business logic HERE (not in intent.go)
- Data management and persistence
- Filtering and sorting
- Validation and error handling

**Example**:
```go
func (c *Context) LoadItems() error {
    items, err := c.Repository.List(c.Context)
    if err != nil {
        return err
    }
    c.Items = items
    c.TotalItems = len(items)
    return nil
}
```

---

### result.go (Output Type - 20-50 lines)

**Contains**:
- `Result` struct (what intent returns on completion)
- Typically: Action, Entity, Message

**Rules**:
- Keep minimal
- Define output contract
- Used by router to handle intent completion

**Example**:
```go
type Result struct {
    Action  string         // "created", "updated", "deleted"
    Item    *career.Burst  // Entity involved
    Message string         // User message
}
```

---

### constants.go (State & Errors - 50-80 lines)

**Contains**:
- `State` enum (typed string)
- State constants (const block)
- Error constants (var block)

**Rules**:
- NO Msg types here (use messages.go)
- Only state machine and errors
- Well-documented constants

**Example**:
```go
type State string

const (
    StateList   State = "list"
    StateDetail State = "detail"
)

var ErrNoItemSelected = errors.New("no item selected")
```

---

### messages.go (Custom Messages - 50-100 lines)

**Contains**:
- ALL `*Msg` types for this intent
- State transition messages
- Async operation results

**Rules**:
- CRITICAL: ALL *Msg types MUST be in this file
- NOT in constants.go
- NOT in intent.go
- Follow naming: `{Entity}{Action}Msg`

**Example**:
```go
type ItemsLoadedMsg struct {
    Items []*career.Item
    Error error
}

type ItemSelectedMsg struct {
    Item  *career.Item
    Index int
}
```

---

### intent.go (Broker - 200-400 lines, MAX 600)

**Contains**:
- `Intent` struct (embeds BaseIntent)
- Constructor (`NewIntent`)
- `Init/Update/View/Result` methods
- `ScreenResultHandler` implementation
- State transition helpers

**Rules**:
- **NO rendering logic** (extract to screens/{feature}/)
- **NO business logic** (move to context.go)
- Only orchestration and delegation
- Target: 200-400 lines
- **HARD BLOCK**: 600 lines (Check #18)

**Example**:
```go
func (i *Intent) Update(msg tea.Msg) tea.Cmd {
    // 1. Check modals FIRST
    if i.deleteModal != nil && i.deleteModal.IsVisible() {
        // Handle modal
    }

    // 2. Handle global keys
    // 3. Delegate to screen
    cmd, result := i.activeScreen.Update(msg)
    if result != nil {
        return i.handleScreenResult(result)
    }
    return cmd
}

func (i *Intent) View() string {
    // Delegate to screen
    return i.activeScreen.View()
}
```

---

## Target File Sizes

| File | Target | Maximum | Enforcement |
|------|--------|---------|-------------|
| context.go | 100-200 lines | 300 lines | Guideline |
| result.go | 20-50 lines | 100 lines | Guideline |
| constants.go | 50-80 lines | 150 lines | Guideline |
| messages.go | 50-100 lines | 200 lines | Guideline |
| **intent.go** | **200-400 lines** | **600 lines** | **BLOCKED (Check #18)** |

---

## Common Patterns

### Centralized Modals

Use modals from `uikit/feedback/` (NOT custom):

```go
// Delete confirmation
i.deleteModal = feedback.NewConfirmModal(
    "Delete Item",
    "Are you sure?",
).WithVariant(feedback.ConfirmDestructive)

return i.deleteModal.Init()  // CRITICAL: Must call Init()
```

### Screen Transitions

```go
func (i *Intent) transitionToDetail() tea.Cmd {
    i.context.CurrentState = StateDetail
    i.detailScreen = screens.NewDetailScreen(i.context.GetSelectedItem())
    i.detailScreen.SetTerminalInfo(i.GetTerminalInfo())
    i.detailScreen.SetTheme(i.Theme())
    i.activeScreen = i.detailScreen
    return nil
}
```

### Business Logic Delegation

```go
// ❌ BAD: Business logic in intent
func (i *Intent) handleSubmit(item *Item) tea.Cmd {
    // Direct database call in intent
    if err := i.repository.Create(item); err != nil {
        return i.setError(err)
    }
}

// ✅ GOOD: Delegate to context
func (i *Intent) handleSubmit(item *Item) tea.Cmd {
    // Context handles business logic
    if err := i.context.CreateItem(item); err != nil {
        return i.setError(err)
    }
}
```

---

## Testing Checklist

After creating your intent, verify:

- [ ] All 5 files exist (context, result, constants, messages, intent)
- [ ] Placeholders replaced ({feature}, {Feature}, {Entity})
- [ ] `make check-intent-architecture` passes
- [ ] intent.go is <400 lines (warn at 400, block at 600)
- [ ] All *Msg types in messages.go (NOT constants.go or intent.go)
- [ ] Business logic in context.go (NOT intent.go)
- [ ] Rendering in screens/{feature}/ (NOT intent.go)
- [ ] Unit tests exist for context business logic
- [ ] Navigation tests exist for intent orchestration

---

## Migration Guide

For existing intents that need migration to this structure:

See: `docs/guides/INTENT_MIGRATION_TO_SUBDIRECTORY.md`

---

## References

- **Architecture Checks**: `scripts/check-intent-architecture.sh`
- **Checklist**: `docs/checklists/INTENT_DEVELOPMENT_CHECKLIST.md`
- **Screen Extraction**: `docs/guides/SCREEN_EXTRACTION_GUIDE.md`
- **Migration Guide**: `docs/guides/INTENT_MIGRATION_TO_SUBDIRECTORY.md`
- **AGENTS.md**: Intent File Organization section

---

## Troubleshooting

### Check #17 Fails (Missing files)

**Error**: "Incomplete subdirectory structure"

**Fix**: Ensure all 5 files exist:
```bash
ls internal/cli/intents/{feature}/
# Should see: context.go result.go constants.go messages.go intent.go
```

### Check #18 Fails (File too large)

**Error**: "intent.go exceeds 600 lines"

**Fix**:
1. Extract rendering to `screens/{feature}/`
2. Move business logic to `context.go`
3. Extract helpers to utility packages

### Check #19 Fails (Rendering in intent)

**Error**: "Rendering methods in intent.go"

**Fix**:
1. Create screen files in `screens/{feature}/`
2. Move render logic to screens
3. Keep only `View()` in intent.go that delegates

### Check #20 Fails (Types in wrong file)

**Error**: "Msg types defined in intent.go"

**Fix**:
- Move ALL *Msg types to `messages.go`
- Move Context to `context.go`
- Move Result to `result.go`
- Move State enum to `constants.go`

---

## Example: Complete Intent Structure

```
internal/cli/intents/burst_management/
├── context.go (180 lines)
│   ├── Context struct
│   ├── NewContext()
│   ├── LoadBursts()
│   ├── CreateBurst()
│   ├── UpdateBurst()
│   ├── DeleteBurst()
│   └── Validate()
│
├── result.go (25 lines)
│   └── Result struct
│
├── constants.go (75 lines)
│   ├── State enum
│   ├── const block (6 states)
│   └── var block (3 error constants)
│
├── messages.go (85 lines)
│   ├── BurstsLoadedMsg
│   ├── BurstSelectedMsg
│   ├── BurstCreatedMsg
│   ├── BurstUpdatedMsg
│   └── BurstDeletedMsg
│
└── intent.go (320 lines)
    ├── Intent struct
    ├── NewIntent()
    ├── Init/Update/View/Result
    ├── ScreenResultHandler impl (4 methods)
    ├── State transitions (4 methods)
    └── Helpers (2 methods)

internal/cli/screens/burst_management/
├── list_screen.go (150 lines)
├── view_screen.go (100 lines)
└── editor_screen.go (120 lines)
```

**Total**: 1,055 lines  
**Intent**: 320 lines (✅ within limits)  
**Screens**: 370 lines (✅ separated)  
**Context**: 365 lines (✅ business logic isolated)

---

## Success Criteria

Your intent follows the standard when:

1. ✅ Subdirectory structure with 5 files
2. ✅ intent.go is 200-400 lines (max 600)
3. ✅ All *Msg types in messages.go
4. ✅ Business logic in context.go
5. ✅ Rendering in screens/{feature}/
6. ✅ Passes `make check-intent-architecture`
7. ✅ All tests pass
8. ✅ Clear separation of concerns

**Remember**: Intent is a BROKER. It orchestrates, it doesn't implement.

---
name: create-intent
description: Create a new intent with proper subdirectory structure following KaRiya architecture rules
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

## What I do

Create a new intent with the proper subdirectory structure and all required files.

## When to use me

Use this skill when creating a new workflow/feature that requires an intent (state machine).

## Quick Start

```bash
make new-intent NAME=feature_name
```

This creates the required structure in `internal/cli/intents/feature_name/`.

## Required File Structure

```
intents/{feature}/
├── context.go    # IntentContext struct + Validate() + domain types
├── result.go     # Result struct (20-50 lines)
├── constants.go  # State enum ONLY (15-50 lines)
├── messages.go   # ALL *Msg types (30-100 lines)
└── intent.go     # NewIntent, Init, Update, View, Result (200-400 lines)
```

**Optional files (for larger intents):**
```
├── types.go      # Intent struct definition (if intent.go > 300 lines)
├── handlers.go   # ALL handle* functions
├── helpers.go    # Utilities, async operations (NO handle* functions)
└── interfaces.go # Service interfaces for DI
```

## Intent Requirements

Every intent MUST:

```go
type MyIntent struct {
    *BaseIntent              // REQUIRED: embed BaseIntent
    context *MyIntentContext // REQUIRED: context struct
    state   MyIntentState    // REQUIRED: typed state enum
    active  bool
    result  *IntentResult[*MyResult]
    
    // Explicit typed screen fields
    listScreen   *myfeature.ListScreen
    detailScreen *myfeature.DetailScreen
    
    // Generic pointer
    activeScreen screens.Screen
}
```

## Screen Result Handler

All intents with screens MUST implement:

```go
var _ ScreenResultHandler = (*MyIntent)(nil)

func (i *MyIntent) HandleCancel(*screens.CancelResult) tea.Cmd { ... }
func (i *MyIntent) HandleNavigate(*screens.NavigateResult) tea.Cmd { ... }
func (i *MyIntent) HandleSubmit(*screens.SubmitResult) tea.Cmd { ... }
func (i *MyIntent) HandleError(*screens.ErrorResult) tea.Cmd { ... }
```

## Handler Organization

**ALL `handle*` functions MUST be in `handlers.go`**:
- Screen result handlers
- Modal update handlers
- Message handlers
- Keyboard handlers

**`helpers.go` MUST NOT contain any `handle*` functions**.

## File Size Limits

| File | Max Lines | Enforcement |
|------|-----------|-------------|
| intent.go | 400 (warning), 600 (BLOCK) | Check #18 |
| context.go | 50-200 | Guideline |
| constants.go | 15-50 | Guideline |

## Screens Structure

Create corresponding screens in:
```
screens/{feature}/
├── list_screen.go
├── detail_screen.go
├── form_screen.go
└── modals/
    ├── filter_modal.go
    └── search_modal.go
```

## Validation

```bash
make check-intent-architecture
```

## Related skills

- `create-screen` - Create screens for the intent
- `bubble-tea-expert` - Bubble Tea TUI patterns (Update signatures, ScreenResult, modals)
- `fix-architecture` - Fix architecture violations
- `check-compliance` - Validate after creation

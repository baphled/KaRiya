---
name: fix-architecture
description: Diagnose and fix KaRiya architecture violations detected by check-intent-architecture
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

## What I do

Help diagnose and fix architecture violations in KaRiya intents and screens.

## When to use me

Use this skill when `make check-intent-architecture` reports violations.

## Run Architecture Check

```bash
make check-intent-architecture
```

## Common Violations and Fixes

### 1. Untyped State Fields

```go
// WRONG
type MyIntent struct {
    state string
}

// CORRECT
type MyIntentState string
const (
    StateList MyIntentState = "list"
)
type MyIntent struct {
    state MyIntentState
}
```

### 2. Missing *BaseIntent

```go
// WRONG
type MyIntent struct {
    // No BaseIntent
}

// CORRECT
type MyIntent struct {
    *BaseIntent  // REQUIRED
}
```

### 3. Missing ScreenResultHandler

```go
// REQUIRED implementation
var _ ScreenResultHandler = (*MyIntent)(nil)

func (i *MyIntent) HandleCancel(*screens.CancelResult) tea.Cmd { ... }
func (i *MyIntent) HandleNavigate(*screens.NavigateResult) tea.Cmd { ... }
func (i *MyIntent) HandleSubmit(*screens.SubmitResult) tea.Cmd { ... }
func (i *MyIntent) HandleError(*screens.ErrorResult) tea.Cmd { ... }
```

### 4. context.Background() in Intents

```go
// WRONG
ctx := context.Background()

// CORRECT
ctx := i.getContext()
```

### 5. Direct huh Import Outside forms/

```go
// WRONG (in intents/, screens/, components/)
import "github.com/charmbracelet/huh"

// CORRECT
import "github.com/baphled/kariya/internal/cli/forms"
```

### 6. Modal Struct in Intents Package

```go
// WRONG - Modal defined in intents/
type MyModal struct { ... }

// CORRECT - Move to:
// - screens/{feature}/modals/my_modal.go
// - uikit/feedback/my_modal.go
```

### 7. Intent File Too Large (>600 lines)

```bash
# Check file size
wc -l internal/cli/intents/my_feature/intent.go
```

**Solution:** Extract to separate files:
- `handlers.go` - All `handle*` functions
- `helpers.go` - Utilities (NO `handle*` functions)
- `types.go` - Struct definitions

### 8. Screens Importing Intents

```go
// WRONG (in screens/)
import "github.com/baphled/kariya/internal/cli/intents"

// CORRECT
// Screens NEVER import intents
// Use ScreenResult to communicate back
```

### 9. Missing Context Struct

```go
// WRONG
type MyIntent struct {
    items []*Item  // Raw parameters
}

// CORRECT
type MyIntentContext struct {
    Items []*Item
}
type MyIntent struct {
    context *MyIntentContext
}
```

### 10. All Types in One File

**WRONG:** Everything in `my_intent.go`

**CORRECT:** Separate files:
- `context.go` - IntentContext
- `constants.go` - State enum
- `messages.go` - *Msg types
- `intent.go` - Implementation

## Layer Hierarchy

```
intents/ -> screens/, uikit/, behaviors/
screens/ -> uikit/, behaviors/ (NEVER intents/)
uikit/ -> themes only (NEVER screens/, intents/)
behaviors/ -> uikit/, themes (NEVER screens/, intents/)
```

## Refusal Template

When violations are detected, refuse to proceed:

```
I cannot proceed with this request.

ARCHITECTURE VIOLATION DETECTED

Violation: [Specific violation]
Rule: [Rule description]
Automated Check: check-intent-architecture.sh (Check #X)

Required correction:
1. [Specific fix needed]

Run `make check-intent-architecture` to verify compliance.
```

## Related skills

- `create-intent` - Create compliant intents
- `create-screen` - Create compliant screens
- `check-compliance` - Full compliance check

---
name: architecture
description: Enforce KaRiya architectural patterns and layer boundaries
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

## What I do

Enforce architectural patterns and layer boundaries in the KaRiya codebase.

## When to use me

Use this skill when:
- Creating new components
- Reviewing code for architectural compliance
- Deciding where code should live
- Resolving dependency questions

## Layer Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        App (Router)                          │
│                    internal/cli/app/                         │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                   Intents (State Machines)                   │
│                   internal/cli/intents/                      │
│         Orchestrate workflows, manage state transitions      │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Screens (Stateless Views)                 │
│                    internal/cli/screens/                     │
│              Render UI, return ScreenResults                 │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    UIKit (Primitives)                        │
│                    internal/cli/uikit/                       │
│           Text, Buttons, Badges, Modals, Layout              │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Behaviors (Composable)                    │
│                   internal/cli/behaviors/                    │
│              TableBehavior, CRUDBehavior, etc.               │
└─────────────────────────────────────────────────────────────┘
```

## Dependency Rules (STRICT)

| Layer | Can Import | NEVER Import |
|-------|------------|--------------|
| `app/` | intents, screens, uikit, behaviors | - |
| `intents/` | screens, uikit, behaviors, domain, service | - |
| `screens/` | uikit, behaviors, domain | **intents/** (FORBIDDEN) |
| `uikit/` | themes only | **screens/, intents/** |
| `behaviors/` | uikit, themes | **screens/, intents/** |
| `forms/` | huh (only place allowed) | **screens/, intents/** |
| `service/` | domain, repository | **cli/*** |
| `repository/` | domain, models | **service/, cli/*** |
| `domain/` | nothing | everything else |

## Package Structure

```
internal/
├── domain/career/           # Domain models (Event, Skill, Burst, Fact)
├── repository/career/       # Repository interfaces + implementations
├── service/                 # Business logic services
└── cli/
    ├── app/                 # Application router
    ├── intents/             # Workflow state machines
    │   └── {feature}/       # Subdirectory per intent
    ├── screens/             # UI screens
    │   └── {feature}/       # Package per feature
    │       └── modals/      # Feature-specific modals
    ├── uikit/               # UI component library
    │   ├── primitives/      # Text, Badge, Button
    │   ├── containers/      # Box, Overlay
    │   ├── feedback/        # Modal, HelpModal
    │   └── layout/          # ScreenLayout, Header
    ├── behaviors/           # Reusable behaviors
    ├── forms/               # Form primitives (huh wrapper)
    ├── themes/              # Theme definitions
    └── types/               # Shared types
```

## Intent Architecture

Every intent follows the subdirectory structure:

```
intents/{feature}/
├── context.go     # Input parameters, Validate()
├── result.go      # Output result type
├── constants.go   # State enum
├── messages.go    # All *Msg types
├── intent.go      # Implementation (NewIntent, Init, Update, View)
├── handlers.go    # All handle* functions (optional)
└── helpers.go     # Utilities, NO handle* (optional)
```

### Intent Requirements

```go
type MyIntent struct {
    *BaseIntent                    // REQUIRED
    context *MyIntentContext       // REQUIRED
    state   MyIntentState          // REQUIRED: typed enum
    
    // Explicit screen fields
    listScreen   *feature.ListScreen
    detailScreen *feature.DetailScreen
    activeScreen screens.Screen
}

// REQUIRED: Implement ScreenResultHandler
var _ ScreenResultHandler = (*MyIntent)(nil)
```

## Screen Architecture

Screens are stateless views that return results:

```go
type MyScreen struct {
    *base.BaseScreen             // REQUIRED
    table *behaviors.TableBehavior[*Item]
}

func (s *MyScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
    // Return results, don't mutate intent state
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.Type == tea.KeyEsc {
            return nil, screens.NewCancelResult("")
        }
    }
    return nil, nil
}
```

## Communication Patterns

### Intent → Screen: Direct method calls
```go
cmd, result := i.listScreen.Update(msg)
```

### Screen → Intent: ScreenResult
```go
return nil, screens.NewNavigateResult("detail", item)
```

### Intent → Service: Method calls
```go
events, err := i.context.Service.ListEvents(ctx)
```

### Service → Repository: Interface
```go
events, err := s.repo.List(ctx, filters)
```

## Form Architecture

`huh` ONLY in `forms/` package:

```go
// forms/inputs.go - Wraps huh
func NewInput(config FieldConfig) *huh.Input { ... }

// screens/{feature}/form_screen.go - Uses forms package
import "github.com/baphled/kariya/internal/cli/forms"
form := forms.NewInput(forms.FieldConfig{...})
```

## Validation Commands

```bash
make check-intent-architecture  # Full architecture check
make check-patterns             # Pattern compliance
make check-compliance           # All checks
```

## Architectural Violations to Refuse

| Violation | Why It's Wrong |
|-----------|----------------|
| Screen imports intent | Circular dependency, breaks layer isolation |
| UIKit imports screen | Lower layer depending on higher |
| huh outside forms/ | Leaky abstraction |
| Business logic in intent | Intent should orchestrate, not implement |
| State mutation in screen | Screens are stateless views |
| Modal in intents/ | Modals belong in screens/ or uikit/ |
| Raw lipgloss in intent | Use UIKit primitives |

## Related skills

- `create-intent` - Create architecturally compliant intents
- `create-screen` - Create compliant screens
- `fix-architecture` - Fix violations
- `code-reviewer` - Review for architecture

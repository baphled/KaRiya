# KaRiya Examples - Canonical Patterns

**Purpose**: These examples demonstrate the CORRECT patterns for extending KaRiya.
AI agents and developers should reference these before creating new components.

## Directory Contents

| Example | Purpose | When to Use |
|---------|---------|-------------|
| `intent_template.go` | Minimal intent implementation | Creating any new intent |
| `form_wrapper_template.go` | Form wrapper for intents | Adding forms to intents |
| `modal_template.go` | Modal with overlay pattern | Creating modals |
| `behavior_usage.go` | Using TableBehavior[T] | Adding table-based views |

## Quick Reference: Must-Use Components

Before writing ANY TUI code, check if a component exists:

### Behaviors (`internal/cli/behaviors/`)
- `TableBehavior[T]` - Data binding, pagination, navigation for tables
- `CRUDBehavior[T]` - Create/edit/delete with confirmation modals
- `FilterMenuBehavior[T]` - Sectioned filter menu with multi-select
- `SortMenuBehavior[T]` - Sort options menu with comparators

### UIKit (`internal/cli/uikit/`)
- `primitives/Text` - Semantic text (Title, Body, Muted, Error, Success)
- `primitives/Button` - Interactive buttons (Primary, Secondary, Danger, Ghost)
- `primitives/ButtonGroup` - Keyboard-navigable button groups
- `primitives/Input` - Themed text input wrapper
- `primitives/Badge` - Status badges (Key, Status, Tag)
- `containers/Box` - Bordered containers with variants
- `containers/Overlay` - Centered modal overlay

### Forms (`internal/cli/forms/`)
- `validators.go` - 20+ reusable validators (Required, ValidateDate, Email, URL)
- `burst_form.go` - Burst editing form config
- `fact_form.go` - Fact editing form config
- `metadata_form.go` - Metadata editing form config
- `skill_form.go` - Skill editing form config
- `capture_event_form.go` - Event capture form config

### Components (`internal/cli/components/`)
- `StandardView` - Standardized layout with logo and footer
- `ModalContent` - 5 modal types (Error, Loading, Progress, Success, Warning)
- `FilterModalModel` - Filter modal
- `HelpModal` - Context-aware help modal

## Anti-Patterns to Avoid

```go
// BAD: Direct *huh.Form in intent
type MyIntent struct {
    form *huh.Form  // WRONG - causes alignment issues
}

// GOOD: Use form wrapper model
type MyIntent struct {
    form *models.CaptureForm  // CORRECT - handles WindowSizeMsg
}
```

```go
// BAD: Manual table implementation
table := table.New(...)  // WRONG - duplicates behavior code

// GOOD: Use TableBehavior
table := behaviors.NewTableBehavior(theme, columns, formatter)
```

```go
// BAD: Hardcoded colors
style := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))

// GOOD: Use theme
style := lipgloss.NewStyle().Foreground(theme.Error())
```

## See Also

- [AGENTS.md](../AGENTS.md) - Complete project documentation
- [docs/FORMS_GUIDE.md](../docs/FORMS_GUIDE.md) - Forms implementation guide
- [docs/TUI_STANDARDS.md](../docs/TUI_STANDARDS.md) - TUI design standards

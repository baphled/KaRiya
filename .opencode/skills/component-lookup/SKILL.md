---
name: component-lookup
description: Find the correct KaRiya component or pattern for a specific need
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

## What I do

Help you find the correct component, pattern, or approach for common UI and architecture needs.

## When to use me

Use this skill when you're unsure which component or pattern to use for a specific need.

## Quick Lookup

```bash
make what-to-use NEED="keyword"
```

Keywords: table, form, modal, color, badge, layout, theme, screen

## Component Reference

### Tables

| Need | Use |
|------|-----|
| Table with selection | `behaviors.TableBehavior[T]` |
| CRUD operations | `behaviors.CRUDBehavior[T]` |

```go
table := behaviors.NewTableBehavior(items, &behaviors.TableConfig{
    Columns: []behaviors.Column{...},
    OnSelect: func(item *T) tea.Cmd { ... },
})
```

### Forms

| Need | Use | NOT |
|------|-----|-----|
| Form in screen | `forms.NewInput()`, `forms.NewSelect()` | Direct `huh.*` |
| Form state check | `forms.IsCompleted(f)` | `f.State == huh.StateCompleted` |
| Form in intent | Use screen with embedded form | Raw `*huh.Form` |

```go
// In forms/ package
input := forms.NewInput(forms.FieldConfig{
    Key:   "title",
    Title: "Title",
})

// Check completion
if forms.IsCompleted(form) { ... }
if forms.IsAborted(form) { ... }
```

### Colors and Themes

| Need | Use | NOT |
|------|-----|-----|
| Primary color | `theme.Primary()` | `lipgloss.Color("#xxx")` |
| Secondary color | `theme.Secondary()` | Hardcoded colors |
| Background | `theme.BackgroundColor()` | |
| Border color | `theme.BorderColor()` | |

### Layout

| Need | Use |
|------|-----|
| Screen layout | `layout.NewScreenLayout()` |
| Header | `layout.NewHeader()` |
| Footer | `layout.NewFooter()` |

```go
layout := layout.NewScreenLayout(theme).
    Header(header).
    Content(content).
    Footer(footer).
    Render()
```

### Modals

| Need | Use |
|------|-----|
| Create modal | `feedback.NewModal()` |
| Render over screen | `behaviors.RenderModalOverlay(modal, baseView)` |
| Help modal | `feedback.NewHelpModal()` |
| Confirm modal | `feedback.NewConfirmModal()` |

```go
// In View()
if i.modal != nil && i.modal.IsVisible() {
    return behaviors.RenderModalOverlay(i.modal, baseView)
}
```

### Badges and Help

| Need | Use | NOT |
|------|-----|-----|
| Key badge | `primitives.HelpKeyBadge("key", "action")` | `components.KeyBadge` |
| Status badge | `primitives.Badge()` | |

### Text Primitives

| Need | Use |
|------|-----|
| Title | `primitives.Title(text)` |
| Subtitle | `primitives.Subtitle(text)` |
| Body text | `primitives.Body(text)` |
| Muted text | `primitives.Muted(text)` |

### Containers

| Need | Use |
|------|-----|
| Box with border | `containers.NewBox(theme)` |
| Overlay | `containers.NewOverlay()` |

## Architecture Patterns

### Screen Result Communication

```go
// Screen returns results
func (s *MyScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
    return nil, screens.NewNavigateResult("detail", item)
}

// Intent handles results
func (i *MyIntent) HandleNavigate(r *screens.NavigateResult) tea.Cmd {
    i.state = StateDetail
    return nil
}
```

### Modal Rendering

```go
func (i *MyIntent) View() string {
    baseView := i.activeScreen.View()
    
    if i.modal != nil && i.modal.IsVisible() {
        return behaviors.RenderModalOverlay(i.modal, baseView)
    }
    
    return baseView
}
```

### State Machine Pattern

```go
type MyIntentState string

const (
    StateList   MyIntentState = "list"
    StateDetail MyIntentState = "detail"
    StateEdit   MyIntentState = "edit"
)

func (i *MyIntent) Update(msg tea.Msg) tea.Cmd {
    switch i.state {
    case StateList:
        return i.updateListState(msg)
    case StateDetail:
        return i.updateDetailState(msg)
    }
    return nil
}
```

## Documentation

- UIKit Guide: `docs/UIKIT_GUIDE.md`
- Forms Guide: `docs/FORMS_GUIDE.md`
- Modal Patterns: `docs/MODAL_PATTERNS.md`
- Intent Architecture: `docs/INTENT_ARCHITECTURE_GUIDE.md`
- Intent Patterns: `docs/development/INTENT_PATTERNS_LIBRARY.md`

## Related skills

- `create-screen` - Create screens using these components
- `create-intent` - Create intents with proper patterns
- `fix-architecture` - Fix component usage violations

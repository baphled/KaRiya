# Centralized Key Handling Developer Guide

> **Last Updated**: 2026-01-12  
> **Status**: ✅ Production Ready

## Overview

KaRiya uses a centralized key handling system to ensure consistent keyboard behavior across all TUI components. This guide explains how to use the centralized keymaps and patterns.

## Table of Contents

1. [Core Concepts](#core-concepts)
2. [Available Keymaps](#available-keymaps)
3. [Usage Patterns](#usage-patterns)
4. [MessageInterceptor Pattern](#messageinterceptor-pattern)
5. [Migration Guide](#migration-guide)
6. [Testing](#testing)

---

## Core Concepts

### Why Centralized Key Handling?

**Before** (String Literals - ❌ Error-Prone):
```go
case "esc":
    // Handle escape
case "q", "ctrl+c":
    // Handle quit
```

**Problems**:
- Typos cause silent failures
- Inconsistent key names across codebase
- No compile-time safety
- Difficult to refactor

**After** (Centralized Keymaps - ✅ Type-Safe):
```go
globalKeys := navigation.DefaultGlobalKeyMap()

switch {
case key.Matches(msg, globalKeys.Back):
    // Handle escape
case key.Matches(msg, globalKeys.Quit):
    // Handle quit
}
```

**Benefits**:
- Type-safe key definitions
- Compile-time checking
- Consistent across codebase
- Easy to refactor and extend

---

## Available Keymaps

All keymaps are defined in `internal/cli/navigation/keymaps.go`.

### 1. GlobalKeyMap

**Purpose**: Universal shortcuts that work everywhere

**Definition**:
```go
type GlobalKeyMap struct {
    Quit         key.Binding // q, ctrl+c
    Back         key.Binding // esc
    Help         key.Binding // ?, h
    MainMenu     key.Binding // m
    Confirm      key.Binding // y, Y, enter
    Cancel       key.Binding // n, N, esc
}
```

**Usage**:
```go
globalKeys := navigation.DefaultGlobalKeyMap()

switch {
case key.Matches(msg, globalKeys.Quit):
    return m, tea.Quit
case key.Matches(msg, globalKeys.Back):
    m.state = PreviousState
    return m, nil
case key.Matches(msg, globalKeys.Help):
    m.showHelp = !m.showHelp
    return m, nil
}
```

### 2. FormKeyMap

**Purpose**: Form input navigation

**Definition**:
```go
type FormKeyMap struct {
    NextField     key.Binding // tab, down, j
    PrevField     key.Binding // shift+tab, up, k
    Submit        key.Binding // enter
    Cancel        key.Binding // esc
}
```

**Usage**:
```go
formKeys := navigation.DefaultFormKeyMap()

switch {
case key.Matches(msg, formKeys.NextField):
    m.focusNextField()
case key.Matches(msg, formKeys.PrevField):
    m.focusPrevField()
case key.Matches(msg, formKeys.Submit):
    return m, m.submitForm()
}
```

### 3. ListKeyMap

**Purpose**: List/table navigation

**Definition**:
```go
type ListKeyMap struct {
    Up            key.Binding // up, k
    Down          key.Binding // down, j
    PageUp        key.Binding // pgup
    PageDown      key.Binding // pgdown
    Home          key.Binding // home, g
    End           key.Binding // end, G
    Select        key.Binding // enter, space
    Filter        key.Binding // /
}
```

**Usage**:
```go
listKeys := navigation.DefaultListKeyMap()

switch {
case key.Matches(msg, listKeys.Up):
    m.cursor--
case key.Matches(msg, listKeys.Down):
    m.cursor++
case key.Matches(msg, listKeys.Select):
    return m, m.selectItem()
}
```

---

## Usage Patterns

### Pattern 1: Simple Key Handling

For straightforward key handling without delegation:

```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    globalKeys := navigation.DefaultGlobalKeyMap()
    listKeys := navigation.DefaultListKeyMap()

    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch {
        case key.Matches(msg, globalKeys.Quit):
            return m, tea.Quit
        case key.Matches(msg, globalKeys.Back):
            m.state = PreviousState
            return m, nil
        case key.Matches(msg, listKeys.Up):
            m.cursor--
            return m, nil
        case key.Matches(msg, listKeys.Down):
            m.cursor++
            return m, nil
        }
    }
    return m, nil
}
```

### Pattern 2: HandleGlobalKeys Helper

For intents that need standard global key handling:

```go
import "github.com/baphled/kariya/internal/cli/navigation"

func (i *MyIntent) Update(msg tea.Msg) tea.Cmd {
    // Handle global keys first (quit, back, help)
    if cmd := navigation.HandleGlobalKeys(msg, i.BaseIntent); cmd != nil {
        return cmd
    }

    // Handle intent-specific keys
    switch msg := msg.(type) {
    case tea.KeyMsg:
        listKeys := navigation.DefaultListKeyMap()
        switch {
        case key.Matches(msg, listKeys.Select):
            return i.selectItem()
        }
    }
    return nil
}
```

**What HandleGlobalKeys does**:
- Handles `q`/`ctrl+c` → Quit
- Handles `esc` → Back/Cancel based on state type
- Handles `?`/`h` → Toggle help
- Handles `m` → Return to main menu

### Pattern 3: MessageInterceptor (Recommended for Forms/Modals)

For components that delegate to child models (forms, modals):

```go
import "github.com/baphled/kariya/internal/cli/navigation"

func (i *MyIntent) updateFormState(msg tea.Msg) tea.Cmd {
    return navigation.NewMessageInterceptor().
        OnQuit(navigation.StandardQuitHandler()).
        OnHelp(navigation.StandardHelpHandler(i.BaseIntent)).
        OnBack(func() tea.Cmd {
            // Custom back logic
            i.state = PreviousState
            return nil
        }).
        InterceptOr(msg, func() tea.Cmd {
            // Delegate to child model
            return i.form.Update(msg)
        })
}
```

**Advantages**:
- Prevents child models from consuming global keys
- Ensures consistent behavior
- Clean separation of concerns
- Chainable API

---

## MessageInterceptor Pattern

The `MessageInterceptor` is a fluent API for intercepting messages before delegation.

### API Methods

```go
type MessageInterceptor struct { ... }

// Create new interceptor
func NewMessageInterceptor() *MessageInterceptor

// Register handlers (chainable)
func (mi *MessageInterceptor) OnQuit(handler func() tea.Cmd) *MessageInterceptor
func (mi *MessageInterceptor) OnHelp(handler func() tea.Cmd) *MessageInterceptor
func (mi *MessageInterceptor) OnBack(handler func() tea.Cmd) *MessageInterceptor
func (mi *MessageInterceptor) OnMainMenu(handler func() tea.Cmd) *MessageInterceptor

// Intercept or delegate
func (mi *MessageInterceptor) InterceptOr(msg tea.Msg, delegate func() tea.Cmd) tea.Cmd
```

### Complete Example

```go
func (i *CaptureEventIntent) updateFormState(msg tea.Msg) tea.Cmd {
    return navigation.NewMessageInterceptor().
        // Handle quit (q, ctrl+c)
        OnQuit(navigation.StandardQuitHandler()).
        
        // Handle help (?, h)
        OnHelp(navigation.StandardHelpHandler(i.BaseIntent)).
        
        // Handle back (esc) - custom logic
        OnBack(func() tea.Cmd {
            if i.form.HasUnsavedChanges() {
                i.showConfirmDialog = true
                return nil
            }
            i.state = CaptureStateChooseStrategy
            return nil
        }).
        
        // Handle main menu (m)
        OnMainMenu(func() tea.Cmd {
            i.result = &IntentResult[*CaptureEventResult]{
                Status: StatusCancelled,
            }
            return nil
        }).
        
        // If no handlers match, delegate to form
        InterceptOr(msg, func() tea.Cmd {
            return i.form.Update(msg)
        })
}
```

### When to Use MessageInterceptor

✅ **Use MessageInterceptor when**:
- Delegating to child models (forms, lists, modals)
- Need to intercept global keys before child handles them
- Want clean separation between global and local key handling

❌ **Don't use MessageInterceptor when**:
- Simple key handling without delegation
- Performance-critical tight loops
- No child models involved

---

## Migration Guide

### Step 1: Import Required Packages

```go
import (
    "github.com/baphled/kariya/internal/cli/navigation"
    "github.com/charmbracelet/bubbles/key"
)
```

### Step 2: Replace String Literals

**Before**:
```go
case "esc":
    m.state = PreviousState
case "q", "ctrl+c":
    return m, tea.Quit
case "?", "h":
    m.showHelp = !m.showHelp
```

**After**:
```go
globalKeys := navigation.DefaultGlobalKeyMap()

switch {
case key.Matches(msg, globalKeys.Back):
    m.state = PreviousState
case key.Matches(msg, globalKeys.Quit):
    return m, tea.Quit
case key.Matches(msg, globalKeys.Help):
    m.showHelp = !m.showHelp
}
```

### Step 3: Use MessageInterceptor for Delegation

**Before**:
```go
// Form could consume global keys - BAD
return m, m.form.Update(msg)
```

**After**:
```go
return m, navigation.NewMessageInterceptor().
    OnQuit(navigation.StandardQuitHandler()).
    OnHelp(navigation.StandardHelpHandler(m.BaseIntent)).
    OnBack(func() tea.Cmd {
        m.state = PreviousState
        return nil
    }).
    InterceptOr(msg, func() tea.Cmd {
        return m.form.Update(msg)
    })
```

### Step 4: Update Tests

**Before**:
```go
m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
```

**After** (same - no test changes needed):
```go
m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
```

---

## Testing

### Testing Key Handling

```go
var _ = Describe("Key Handling", func() {
    It("should handle escape key", func() {
        intent := NewMyIntent(context)
        
        // Press escape
        intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
        
        // Verify state changed
        Expect(intent.state).To(Equal(PreviousState))
    })
    
    It("should quit on 'q'", func() {
        intent := NewMyIntent(context)
        
        cmd := intent.Update(tea.KeyMsg{
            Type: tea.KeyRunes,
            Runes: []rune{'q'},
        })
        
        Expect(cmd).To(Equal(tea.Quit))
    })
})
```

### Testing MessageInterceptor

```go
It("should intercept escape before form", func() {
    intent := NewMyIntent(context)
    intent.state = MyStateForm
    
    // Escape should change state, not delegate to form
    intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
    
    Expect(intent.state).To(Equal(PreviousState))
})
```

---

## Best Practices

### ✅ DO

1. **Always use centralized keymaps** for global keys (quit, back, help)
2. **Use MessageInterceptor** when delegating to child models
3. **Handle global keys first** before intent-specific keys
4. **Use `key.Matches()`** for all key comparisons
5. **Test escape behavior** for every state

### ❌ DON'T

1. **Don't use string literals** for key checking (`case "esc":`)
2. **Don't delegate without intercepting** global keys
3. **Don't create custom keymaps** for global keys (use DefaultGlobalKeyMap)
4. **Don't mix string literals and key.Matches** (choose one pattern)
5. **Don't skip testing** escape key behavior

---

## Related Documentation

- **[State Matrix](../STATE_MATRIX.md)** - Complete state and escape behavior documentation
- **[Keyboard System Guide](KEYBOARD_SYSTEM_GUIDE.md)** - Comprehensive keyboard system guide
- **[Navigation Testing Guide](NAVIGATION_TESTING_GUIDE.md)** - Testing navigation patterns
- **[TUI Standards](../TUI_STANDARDS.md)** - Universal keyboard shortcuts

---

## Migration Status

**Completed**:
- ✅ `internal/cli/screens/capture/event_form_screen.go` - Fully migrated
- ✅ `internal/cli/intents/browse_timeline_intent.go` - MessageInterceptor pattern
- ✅ `internal/cli/models/` - Package deleted (modals migrated to captureevent/)

**Pending**:
- ⏳ 8 remaining intents (capture_event, generate_cv, etc.)
- ⏳ Modal components (EditBurstModal, EditMetadataModal, EditFactModal)

---

*For questions or suggestions, please open an issue on GitHub.*

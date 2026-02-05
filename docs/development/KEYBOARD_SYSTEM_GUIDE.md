# KaRiya Keyboard Shortcut System - Developer Guide

> **Note**: Some file references in this guide point to `internal/cli/models/` which has been
> deleted. The keyboard patterns described here are still valid — they now live in
> `internal/cli/intents/` and `internal/cli/screens/` packages.

**Complete Guide for Implementing and Extending Keyboard Shortcuts**

**Last Updated**: 2026-02-04
**Version**: 2.1 (models/ package deleted — patterns migrated to intents/screens)
**Audience**: Developers implementing TUI features

---

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Component Reference](#component-reference)
3. [Integration Guide](#integration-guide)
4. [Implementation Patterns](#implementation-patterns)
5. [Testing Shortcuts](#testing-shortcuts)
6. [Best Practices](#best-practices)
7. [Performance Considerations](#performance-considerations)
8. [Troubleshooting](#troubleshooting)

---

## Architecture Overview

The KaRiya keyboard shortcut system provides a consistent, type-safe way to handle keyboard input across all TUI components. It integrates with the StandardView system and supports context-aware shortcut resolution.

### Current System (5 Files)

The shortcut system is currently spread across 5 files:

| File | Lines | Purpose | Status |
|------|-------|---------|--------|
| **shortcut_handler.go** | 125 | Basic shortcut registration | Active |
| **shortcut_mapper.go** | 255 | Global + context shortcuts | Active |
| **context_shortcut_handler.go** | 180 | Context stack management | Active |
| **shortcut_customizer.go** | 261 | UI for customization | Active |
| **shortcut_help_system.go** | 275 | Help display | Active |

**Total**: 1,096 lines across 5 files

### System Layers

```
┌──────────────────────────────────────────────────┐
│         Keyboard Shortcut System                 │
├──────────────────────────────────────────────────┤
│                                                  │
│  Intent Layer                                    │
│  ├─ Intent implements keyboard handlers         │
│  ├─ Context-specific shortcuts                  │
│  └─ Integration with StandardView               │
│                                                  │
│  Shortcut Mapper (shortcut_mapper.go)          │
│  ├─ Global shortcuts registry                   │
│  ├─ Context shortcuts registry                  │
│  └─ Shortcut lookup and merging                 │
│                                                  │
│  Context Handler (context_shortcut_handler.go)  │
│  ├─ Context stack management                    │
│  ├─ Context metadata storage                    │
│  └─ Active shortcut resolution                  │
│                                                  │
│  Help System (shortcut_help_system.go)         │
│  ├─ Help text generation                        │
│  ├─ Context-aware help                          │
│  └─ Search and discovery                        │
│                                                  │
│  Customizer (shortcut_customizer.go)           │
│  └─ UI for shortcut customization              │
│                                                  │
└──────────────────────────────────────────────────┘
```

### Key Design Principles

1. **Context-Aware**: Shortcuts change based on current intent state
2. **Type-Safe**: All shortcuts use BubbleTea's `key.Binding`
3. **Centralized Help**: Helper functions generate consistent footer text
4. **StandardView Integration**: Keyboard shortcuts work with StandardView layout
5. **Testable**: All keyboard handlers have comprehensive test coverage

---

## Component Reference

### 1. ShortcutMapper

**File**: `internal/cli/models/shortcut_mapper.go`

**Purpose**: Manages global and context-specific shortcuts with conflict detection

**Core API**:

```go
type ShortcutMapper struct {
    globalShortcuts  map[string]*ShortcutBinding
    contextShortcuts map[string]map[string]*ShortcutBinding
}

// Global shortcuts
func (sm *ShortcutMapper) RegisterGlobalShortcut(id string, binding key.Binding, action ShortcutAction) error
func (sm *ShortcutMapper) GetGlobalShortcut(id string) (*ShortcutBinding, bool)
func (sm *ShortcutMapper) GetAllGlobalShortcuts() map[string]*ShortcutBinding

// Context shortcuts
func (sm *ShortcutMapper) RegisterContextShortcut(context, id string, binding key.Binding, action ShortcutAction) error
func (sm *ShortcutMapper) GetContextShortcut(context, id string) (*ShortcutBinding, bool)
func (sm *ShortcutMapper) GetContextShortcuts(context string) map[string]*ShortcutBinding

// Utilities
func (sm *ShortcutMapper) CheckConflicts(id string, binding key.Binding) []string
func (sm *ShortcutMapper) GetMergedShortcuts(context string) map[string]*ShortcutBinding
```

**Usage Example**:

```go
// Create mapper
mapper := models.NewShortcutMapper()

// Register global shortcuts
quitBinding := key.NewBinding(
    key.WithKeys("q"),
    key.WithHelp("q", "quit"),
)
mapper.RegisterGlobalShortcut("quit", quitBinding, func() tea.Cmd {
    return tea.Quit
})

// Register context shortcuts
editBinding := key.NewBinding(
    key.WithKeys("e"),
    key.WithHelp("e", "edit"),
)
mapper.RegisterContextShortcut("list", "edit", editBinding, editAction)

// Get merged shortcuts (context overrides global)
shortcuts := mapper.GetMergedShortcuts("list")
```

---

### 2. ContextShortcutHandler

**File**: `internal/cli/models/context_shortcut_handler.go`

**Purpose**: Manages context stacks and context-aware shortcut resolution

**Core API**:

```go
type ContextShortcutHandler struct {
    mapper         *ShortcutMapper
    contextStack   []string
    currentContext string
}

// Context management
func (csh *ContextShortcutHandler) SetCurrentContext(context string)
func (csh *ContextShortcutHandler) GetCurrentContext() string
func (csh *ContextShortcutHandler) PushContext(context string)
func (csh *ContextShortcutHandler) PopContext() string
func (csh *ContextShortcutHandler) GetContextStack() []string

// Shortcut resolution
func (csh *ContextShortcutHandler) GetActiveShortcuts() map[string]*ShortcutBinding
func (csh *ContextShortcutHandler) HandleMsg(msg tea.Msg) (tea.Cmd, bool)
```

**Usage Example**:

```go
// Create handler
handler := models.NewContextShortcutHandler()
handler.SetMapper(mapper)

// Switch contexts
handler.SetCurrentContext("list")

// Get active shortcuts for current context
activeShortcuts := handler.GetActiveShortcuts()

// Use context stack for modal dialogs
handler.PushContext("modal")
// ... modal operations ...
handler.PopContext() // Back to previous context
```

---

### 3. ShortcutHelpSystem

**File**: `internal/cli/models/shortcut_help_system.go`

**Purpose**: Provides help information and discoverability for shortcuts

**Core API**:

```go
type ShortcutHelpSystem struct {
    shortcuts  map[string]*ShortcutInfo
    categories map[string][]string
}

// Registration
func (shs *ShortcutHelpSystem) RegisterShortcut(id string, binding key.Binding, description string, contexts []string)
func (shs *ShortcutHelpSystem) CategorizeShortcut(id string, category string)

// Help generation
func (shs *ShortcutHelpSystem) GenerateContextHelp(context string) string
func (shs *ShortcutHelpSystem) GenerateAllHelp() string
func (shs *ShortcutHelpSystem) SearchShortcuts(query string) []string
```

**Usage Example**:

```go
// Create help system
help := models.NewShortcutHelpSystem()

// Register shortcuts with help text
help.RegisterShortcut("save", binding, "Save the file", []string{"editor"})

// Categorize for organization
help.CategorizeShortcut("save", "File Operations")

// Get context-specific help
contextHelp := help.GenerateContextHelp("editor")

// Search for shortcuts
results := help.SearchShortcuts("save")
```

---

### 4. ShortcutCustomizer

**File**: `internal/cli/models/shortcut_customizer.go`

**Purpose**: Allows users to customize shortcuts and manage profiles

**Core API**:

```go
type ShortcutCustomizer struct {
    mapper   *ShortcutMapper
    profiles map[string]*ShortcutProfile
}

// Customization
func (sc *ShortcutCustomizer) OverrideShortcut(id string, newBinding key.Binding) error
func (sc *ShortcutCustomizer) ResetShortcut(id string) error
func (sc *ShortcutCustomizer) GetAllOverrides() map[string]key.Binding

// Profiles
func (sc *ShortcutCustomizer) CreateProfile(name string) *ShortcutProfile
func (sc *ShortcutCustomizer) ApplyProfile(name string) error
func (sc *ShortcutCustomizer) GetActiveProfile() string
```

**Usage Example**:

```go
// Create customizer
customizer := models.NewShortcutCustomizer()

// Override a shortcut
newBinding := key.NewBinding(key.WithKeys("ctrl+shift+s"))
customizer.OverrideShortcut("save", newBinding)

// Create profiles
profile := customizer.CreateProfile("work")
profile.OverrideShortcut("save", binding)

// Apply profile
customizer.ApplyProfile("work")
```

---

## Integration Guide

### Adding Shortcuts to an Intent

All intents in KaRiya use a consistent pattern for keyboard shortcuts:

#### Step 1: Define Keyboard Handler in Update()

```go
func (i *MyIntent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch i.state {
        case StateInitial:
            return i.updateInitial(msg)
        case StateWorking:
            return i.updateWorking(msg)
        case StateFinal:
            return i.updateFinal(msg)
        }
    }
    return i, nil
}
```

#### Step 2: Implement State-Specific Handlers

```go
func (i *MyIntent) updateInitial(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    switch msg.String() {
    case "esc":
        if i.isRootState() {
            i.setCancelled()  // Cancel intent
        } else {
            i.state = i.previousState  // Go back
        }
        return i, nil
    
    case "m":
        i.setCancelled()  // Always return to main menu
        return i, nil
    
    case "q", "ctrl+c":
        return i, tea.Quit  // Quit application
    
    case "up", "k":
        i.selectedIndex = max(0, i.selectedIndex-1)
        return i, nil
    
    case "down", "j":
        i.selectedIndex = min(len(i.items)-1, i.selectedIndex+1)
        return i, nil
    
    case "enter":
        // Handle selection
        return i, i.handleSelection()
    }
    
    return i, nil
}
```

#### Step 3: Provide Context-Aware Help

```go
func (i *MyIntent) getContextHelp() string {
    theme := i.Theme()
    if theme == nil {
        return DefaultHelp()
    }
    
    switch i.state {
    case StateInitial:
        return CombineThemedFooters(
            ThemedNavigationFooter(theme),  // ↑/k Up  ↓/j Down  Enter Select
            ThemedGlobalBadges(theme),       // Esc Back  m Main  q Quit
        )
    
    case StateWorking:
        return CombineThemedFooters(
            ThemedCustomFooter(theme,
                components.NewKeyBadge("Ctrl+S", "Save"),
                components.NewKeyBadge("Enter", "Continue"),
            ),
            ThemedGlobalBadges(theme),
        )
    
    case StateFinal:
        return CombineThemedFooters(
            ThemedCustomFooter(theme,
                components.NewKeyBadge("Enter", "Done"),
                components.NewKeyBadge("r", "Retry"),
            ),
            ThemedGlobalBadges(theme),
        )
    }
    
    return DefaultHelp()
}
```

#### Step 4: Integrate with StandardView

```go
func (i *MyIntent) View() string {
    // Build content based on state
    content := i.buildContent()
    
    // Get context-aware help
    help := i.getContextHelp()
    
    // Create StandardView with breadcrumbs
    view := i.CreateViewWithBreadcrumbs(
        "Main Menu",
        "My Intent",
        i.getStateName(),
    )
    
    view.WithContent(content)
    view.WithHelp(help).WithFooterSeparator(true)
    
    // Apply modals if needed (error, loading, etc.)
    return view.Render()
}
```

---

### Context Management Patterns

#### Modal Dialog Pattern

When showing a modal, use context stack:

```go
// Before showing modal
func (i *MyIntent) showModal() tea.Cmd {
    originalState := i.state
    i.state = StateModal
    i.previousState = originalState
    
    // Modal will handle Esc to return to previousState
    return nil
}

// In modal Update()
case "esc":
    i.state = i.previousState  // Return to previous state
    return i, nil
```

#### Async Operation Pattern

For async operations (generating, exporting, saving):

```go
case StateInProgress:
    switch msg.String() {
    case "esc":
        // Let operation complete in background
        i.state = i.previousState
        return i, nil
    
    case "m":
        // Cancel operation immediately
        i.setCancelled()
        return i, nil
    
    case "q", "ctrl+c":
        return i, tea.Quit
    }
```

#### Error State Pattern

When showing errors, preserve them when going back:

```go
case StateError:
    switch msg.String() {
    case "esc":
        // Keep error visible when going back
        i.state = i.previousState
        // Error field remains set
        return i, nil
    
    case "r":
        // Retry - clear error
        i.ClearError()
        return i, i.retry()
    
    case "m":
        i.setCancelled()
        return i, nil
    }
```

---

## Implementation Patterns

### Pattern 1: Global Shortcuts (All Intents)

**Universal shortcuts that work in every intent**:

```go
// In every Update() method, handle these first:
switch msg.String() {
case "q", "ctrl+c":
    return i, tea.Quit  // Quit application

case "?":
    i.ToggleHelp()  // Toggle help modal
    return i, nil

case "m":
    i.setCancelled()  // Return to main menu
    return i, nil

case "esc":
    // Context-dependent (see Escape Key Pattern)
    if i.isRootState() {
        i.setCancelled()
    } else {
        i.state = i.previousState
    }
    return i, nil
}
```

### Pattern 2: Navigation Shortcuts (Lists and Menus)

**Standard navigation pattern**:

```go
case "up", "k":
    i.selectedIndex = max(0, i.selectedIndex-1)
    return i, nil

case "down", "j":
    i.selectedIndex = min(len(i.items)-1, i.selectedIndex+1)
    return i, nil

case "pgup":
    i.selectedIndex = max(0, i.selectedIndex-10)
    return i, nil

case "pgdown":
    i.selectedIndex = min(len(i.items)-1, i.selectedIndex+10)
    return i, nil

case "home", "g":
    i.selectedIndex = 0
    return i, nil

case "end", "G":
    i.selectedIndex = len(i.items) - 1
    return i, nil

case "enter":
    return i, i.handleSelection()
```

### Pattern 3: Form Navigation (Huh Forms)

**Forms use Huh library for navigation**:

```go
// Huh form automatically handles:
// - Tab / Shift+Tab (field navigation)
// - Enter (submit)
// - Arrow keys (within fields)

// Intent only needs to handle:
case "ctrl+s":
    // Force submit
    return i, i.submitForm()

case "ctrl+o":
    // Toggle optional fields (Manual mode only)
    if i.captureMode == "manual" {
        i.showOptionalFields = !i.showOptionalFields
    }
    return i, nil

case "esc":
    // Cancel form
    i.state = i.previousState
    return i, nil
```

### Pattern 4: Help Footer Construction

**Build context-aware help using helper functions**:

```go
func (i *MyIntent) getContextHelp() string {
    theme := i.Theme()
    if theme == nil {
        return DefaultHelp()
    }
    
    // Navigation footer (↑/k Up  ↓/j Down  Enter Select)
    navFooter := ThemedNavigationFooter(theme)
    
    // Form footer (Tab Next  Shift+Tab Prev  Enter Submit)
    formFooter := ThemedFormFooter(theme)
    
    // Custom badges for specific actions
    customFooter := ThemedCustomFooter(theme,
        components.NewKeyBadge("e", "Edit"),
        components.NewKeyBadge("d", "Delete"),
    )
    
    // Global badges (Esc Back  m Main  q Quit)
    globalFooter := ThemedGlobalBadges(theme)
    
    // Combine footers
    return CombineThemedFooters(navFooter, customFooter, globalFooter)
}
```

**Available Helper Functions** (from `internal/cli/intents/view_helpers.go`):

```go
// Predefined footers
func ThemedNavigationFooter(theme *themes.Theme) string
func ThemedFormFooter(theme *themes.Theme) string
func ThemedGlobalBadges(theme *themes.Theme) string

// Custom footer
func ThemedCustomFooter(theme *themes.Theme, badges ...components.KeyBadge) string

// Combine multiple footers
func CombineThemedFooters(footers ...string) string
```

---

## Testing Shortcuts

### Unit Testing Pattern

**Test keyboard handlers for each state**:

```go
var _ = Describe("MyIntent Shortcuts", func() {
    var intent *MyIntent
    
    BeforeEach(func() {
        intent = NewMyIntent(mockContext)
    })
    
    Describe("State: Initial", func() {
        It("should navigate up with 'k'", func() {
            intent.selectedIndex = 5
            
            msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
            intent.Update(msg)
            
            Expect(intent.selectedIndex).To(Equal(4))
        })
        
        It("should navigate down with 'j'", func() {
            intent.selectedIndex = 5
            
            msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
            intent.Update(msg)
            
            Expect(intent.selectedIndex).To(Equal(6))
        })
        
        It("should cancel with 'esc' in root state", func() {
            msg := tea.KeyMsg{Type: tea.KeyEsc}
            intent.Update(msg)
            
            result := intent.Result()
            Expect(result.Status).To(Equal(StatusCancelled))
        })
        
        It("should go back with 'esc' in non-root state", func() {
            intent.state = StateWorking
            intent.previousState = StateInitial
            
            msg := tea.KeyMsg{Type: tea.KeyEsc}
            intent.Update(msg)
            
            Expect(intent.state).To(Equal(StateInitial))
        })
        
        It("should quit with 'q'", func() {
            msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
            _, cmd := intent.Update(msg)
            
            Expect(cmd).NotTo(BeNil())
            // cmd should be tea.Quit
        })
    })
})
```

### Integration Testing Pattern

**Test complete keyboard flows through multiple states**:

```go
var _ = Describe("MyIntent Navigation Flow", func() {
    It("should complete full workflow with keyboard only", func() {
        intent := NewMyIntent(mockContext)
        
        // Initial state - select item
        intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
        intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
        intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
        
        Expect(intent.state).To(Equal(StateWorking))
        
        // Working state - confirm
        intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
        
        Expect(intent.state).To(Equal(StateFinal))
        
        // Final state - complete
        intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
        
        result := intent.Result()
        Expect(result.Status).To(Equal(StatusCompleted))
    })
    
    It("should handle back navigation correctly", func() {
        intent := NewMyIntent(mockContext)
        intent.state = StateFinal
        intent.previousState = StateWorking
        
        // Go back
        intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
        Expect(intent.state).To(Equal(StateWorking))
        
        // Go back again
        intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
        Expect(intent.state).To(Equal(StateInitial))
        
        // Cancel from root
        intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
        
        result := intent.Result()
        Expect(result.Status).To(Equal(StatusCancelled))
    })
})
```

### View Testing Pattern

**Test that help footer displays correct shortcuts**:

```go
var _ = Describe("MyIntent Help Footer", func() {
    var intent *MyIntent
    
    BeforeEach(func() {
        intent = NewMyIntent(mockContext)
    })
    
    It("should show navigation shortcuts in Initial state", func() {
        intent.state = StateInitial
        
        view := intent.View()
        
        Expect(view).To(ContainSubstring("↑/k"))
        Expect(view).To(ContainSubstring("Up"))
        Expect(view).To(ContainSubstring("Enter"))
        Expect(view).To(ContainSubstring("Select"))
    })
    
    It("should show form shortcuts in Working state", func() {
        intent.state = StateWorking
        
        view := intent.View()
        
        Expect(view).To(ContainSubstring("Ctrl+S"))
        Expect(view).To(ContainSubstring("Save"))
    })
    
    It("should always show global shortcuts", func() {
        for _, state := range []string{StateInitial, StateWorking, StateFinal} {
            intent.state = state
            view := intent.View()
            
            Expect(view).To(ContainSubstring("Esc"))
            Expect(view).To(ContainSubstring("m"))
            Expect(view).To(ContainSubstring("q"))
        }
    })
})
```

---

## Best Practices

### 1. Use Consistent Patterns

✅ **DO**: Follow established keyboard handler patterns

```go
// Good - consistent with all intents
case "esc":
    if i.isRootState() {
        i.setCancelled()
    } else {
        i.state = i.previousState
    }
```

❌ **DON'T**: Create custom navigation that breaks user expectations

```go
// Bad - inconsistent with other intents
case "esc":
    i.state = StateRandom  // Unpredictable behavior
```

### 2. Always Provide Help Text

✅ **DO**: Implement `getContextHelp()` for every intent

```go
func (i *MyIntent) getContextHelp() string {
    theme := i.Theme()
    switch i.state {
    case StateInitial:
        return CombineThemedFooters(
            ThemedNavigationFooter(theme),
            ThemedGlobalBadges(theme),
        )
    // ... all states
    }
}
```

❌ **DON'T**: Return empty or static help text

```go
func (i *MyIntent) getContextHelp() string {
    return ""  // Bad - users can't discover shortcuts
}
```

### 3. Handle Global Shortcuts First

✅ **DO**: Check global shortcuts before state-specific ones

```go
// Global shortcuts handled first
switch msg.String() {
case "q", "ctrl+c":
    return i, tea.Quit
case "m":
    i.setCancelled()
    return i, nil
case "?":
    i.ToggleHelp()
    return i, nil
}

// Then state-specific shortcuts
switch i.state {
case StateInitial:
    return i.updateInitial(msg)
}
```

❌ **DON'T**: Let state handlers override global shortcuts

```go
// Bad - 'q' might not quit in some states
switch i.state {
case StateInitial:
    // State handler might not handle 'q'
}
```

### 4. Preserve Errors on Back Navigation

✅ **DO**: Keep errors visible when user goes back

```go
case "esc":
    // Error remains visible
    i.state = i.previousState
    return i, nil
```

❌ **DON'T**: Clear errors on back navigation

```go
case "esc":
    i.ClearError()  // Bad - user loses context
    i.state = i.previousState
```

### 5. Support Both Arrow and Vim Keys

✅ **DO**: Accept both navigation styles

```go
case "up", "k":
    i.selectedIndex--
case "down", "j":
    i.selectedIndex++
```

❌ **DON'T**: Support only one navigation style

```go
case "up":  // Bad - vim users can't use 'k'
    i.selectedIndex--
```

### 6. Test All Keyboard Paths

✅ **DO**: Write tests for every keyboard shortcut

```go
It("should navigate up with 'k'", func() { ... })
It("should navigate up with arrow key", func() { ... })
It("should handle esc from each state", func() { ... })
```

❌ **DON'T**: Test only happy path

```go
It("should complete workflow", func() {
    // Only tests Enter key
})
```

---

## Performance Considerations

### Shortcut Lookup Complexity

- **GetGlobalShortcut**: O(1) via map lookup
- **GetContextShortcut**: O(1) via nested map lookup
- **GetMergedShortcuts**: O(n + m) where n = global, m = context shortcuts
- **CheckConflicts**: O(n) where n = total shortcuts

### Optimization Tips

1. **Cache merged shortcuts** when context doesn't change frequently

```go
type MyIntent struct {
    mergedCache map[string]*ShortcutBinding
    cacheDirty  bool
}

func (i *MyIntent) GetActiveShortcuts() map[string]*ShortcutBinding {
    if !i.cacheDirty {
        return i.mergedCache
    }
    
    i.mergedCache = i.mapper.GetMergedShortcuts(i.currentContext)
    i.cacheDirty = false
    return i.mergedCache
}
```

2. **Avoid redundant help text generation**

```go
type MyIntent struct {
    helpCache   map[string]string
    lastState   string
}

func (i *MyIntent) getContextHelp() string {
    if cached, ok := i.helpCache[i.state]; ok && i.state == i.lastState {
        return cached
    }
    
    help := i.generateHelp()
    i.helpCache[i.state] = help
    i.lastState = i.state
    return help
}
```

3. **Lazy-load customizer profiles**

```go
// Only load profile when user accesses settings
func (sc *ShortcutCustomizer) GetProfile(name string) *ShortcutProfile {
    if profile, ok := sc.profiles[name]; ok {
        return profile
    }
    
    // Load from disk only when needed
    return sc.loadProfile(name)
}
```

---

## Troubleshooting

### Shortcuts Not Registering

**Problem**: Shortcuts don't work in new intent

**Checklist**:
1. ✅ Intent Update() method handles tea.KeyMsg
2. ✅ Global shortcuts (q, ?, m, Esc) handled first
3. ✅ State-specific shortcuts in correct case statement
4. ✅ Help footer shows registered shortcuts

**Solution**:

```go
func (i *MyIntent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Handle global shortcuts first
        switch msg.String() {
        case "q", "ctrl+c":
            return i, tea.Quit
        case "m":
            i.setCancelled()
            return i, nil
        }
        
        // Then state-specific shortcuts
        switch i.state {
        case StateInitial:
            return i.updateInitial(msg)
        }
    }
    return i, nil
}
```

### Conflicts Between Shortcuts

**Problem**: Two shortcuts use the same key

**Checklist**:
1. ✅ Check if shortcut is global vs context-specific
2. ✅ Verify context boundaries
3. ✅ Use ShortcutMapper.CheckConflicts()

**Solution**:

```go
conflicts := mapper.CheckConflicts("my-action", newBinding)
if len(conflicts) > 0 {
    log.Printf("Shortcut conflict detected: %v", conflicts)
    // Reassign to different key or resolve conflict
}
```

### Help Text Not Showing

**Problem**: Footer doesn't display shortcuts

**Checklist**:
1. ✅ Intent implements getContextHelp()
2. ✅ View() calls view.WithHelp(help)
3. ✅ Terminal size sufficient (> 40 width)
4. ✅ Footer separator enabled with WithFooterSeparator(true)

**Solution**:

```go
func (i *MyIntent) View() string {
    content := i.buildContent()
    help := i.getContextHelp()  // Must implement
    
    view := i.CreateViewWithBreadcrumbs("Main", "Intent", i.state)
    view.WithContent(content)
    view.WithHelp(help).WithFooterSeparator(true)  // Must call WithHelp
    
    return view.Render()
}
```

### Context Stack Issues

**Problem**: Context doesn't restore after modal

**Checklist**:
1. ✅ PushContext() called before modal
2. ✅ PopContext() called when modal closes
3. ✅ Context stack depth verified

**Solution**:

```go
// Before modal
handler.PushContext("modal")
originalDepth := handler.GetContextDepth()

// After modal closes
handler.PopContext()
newDepth := handler.GetContextDepth()

Expect(newDepth).To(Equal(originalDepth - 1))
```

---

## Related Documentation

- **[Keyboard Shortcuts Guide](../KEYBOARD_SHORTCUTS_GUIDE.md)** - User-facing keyboard reference
- **[TUI Standards](../TUI_STANDARDS.md)** - Complete TUI design standards
- **[TUI Developer Guide](../TUI_DEVELOPER_GUIDE.md)** - General TUI development guide
- **[StandardView Guide](../STANDARDVIEW_GUIDE.md)** - StandardView integration
- **[Intent Development Checklist](../INTENT_DEVELOPMENT_CHECKLIST.md)** - Checklist for new intents

---

## Future Enhancements

Planned improvements to the keyboard shortcut system:

1. **Unified ShortcutManager** - Consolidate 5 files into single manager (see UNIFIED_SHORTCUT_SYSTEM_DESIGN.md)
2. **Key Macro Recording** - Record and replay key sequences
3. **Shortcut Analytics** - Track which shortcuts users use most
4. **Smart Recommendations** - Suggest shortcuts based on usage patterns
5. **Accessibility** - Voice command shortcuts
6. **Internationalization** - Shortcuts for different keyboard layouts

**See**: [Unified Shortcut System Design](../UNIFIED_SHORTCUT_SYSTEM_DESIGN.md) for detailed consolidation plan

---

**Document Status**: ✅ Complete and Production-Ready
**Last Updated**: 2026-01-12
**Version**: 2.0 (Consolidated)
**Next Review**: 2026-04-01

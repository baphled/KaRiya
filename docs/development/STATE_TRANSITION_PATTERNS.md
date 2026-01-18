# State Transition Patterns Guide

**Last Updated**: 2026-01-14  
**Status**: Best Practices Guide  
**Audience**: Intent Developers

---

## Table of Contents

1. [Overview](#overview)
2. [Why Use Transition Helpers](#why-use-transition-helpers)
3. [When to Extract Helpers](#when-to-extract-helpers)
4. [Pattern: Extracted Transition Helpers](#pattern-extracted-transition-helpers)
5. [Pattern: Inline Transitions](#pattern-inline-transitions)
6. [Reference Implementation](#reference-implementation)
7. [Best Practices](#best-practices)
8. [Common Mistakes](#common-mistakes)

---

## Overview

State transitions are a core part of intent state machines. As intents grow in complexity, managing state transitions becomes critical for maintainability and correctness.

This guide documents two approaches:
1. **Extracted Transition Helpers** - Dedicated methods for each transition (recommended for complex intents)
2. **Inline Transitions** - Direct state changes in event handlers (acceptable for simple intents)

**Key Decision**: Use extracted helpers when an intent has **4+ states** or **complex initialization/cleanup** per state.

---

## Why Use Transition Helpers

### Benefits

✅ **Encapsulation**: All state change logic in one place  
✅ **Consistency**: Same transition always follows same steps  
✅ **Testability**: Easy to test transitions in isolation  
✅ **Readability**: Clear intent: `transitionToDetailScreen()` vs `currentState = StateDetail`  
✅ **Safety**: Ensures cleanup happens, prevents partial transitions  
✅ **Maintainability**: Single place to update when transition logic changes

### Tradeoffs

⚠️ **Overhead**: Extra methods to maintain  
⚠️ **Indirection**: One more hop to understand flow  
⚠️ **Not Always Needed**: Overkill for simple 2-3 state machines

---

## When to Extract Helpers

### Extract Helpers When:

✅ Intent has **4+ states**  
✅ Transitions require **initialization** (creating screens, loading data)  
✅ Transitions require **cleanup** (clearing modals, resetting state)  
✅ Multiple paths lead to the same target state  
✅ Transition logic is **15+ lines**  
✅ You find yourself copying transition code

### Use Inline When:

✅ Intent has **≤3 states**  
✅ Transitions are **simple** (just `currentState = NewState`)  
✅ Each transition is **unique** (no duplication)  
✅ Transition logic is **<10 lines**  
✅ Intent uses **different architecture** (not screen-based)

---

## Pattern: Extracted Transition Helpers

### Structure

```go
// transitionToXState transitions to state X with full initialization.
func (i *YourIntent) transitionToXState() tea.Cmd {
    // 1. Set the new state
    i.currentState = StateX
    
    // 2. Initialize resources for new state
    i.xResource = createXResource()
    
    // 3. Clean up old state (if needed)
    i.oldResource = nil
    
    // 4. Return initialization command (or nil)
    return i.xResource.Init()
}
```

### Naming Convention

Use clear, consistent naming:
- `transitionToListScreen()` - Screen-based architecture
- `transitionToDetailState()` - State machine
- `transitionToXFromY()` - When transition logic varies by source

### Example: Simple Transition

```go
// transitionToListScreen returns to the skills list.
func (i *ManageSkillsIntent) transitionToListScreen() tea.Cmd {
    i.currentState = SkillsStateList
    
    // Create list screen
    listScreen := NewSkillsListScreenFromIntent(i.skills, i.GetThemeManager())
    
    // Apply intent context (terminal, theme, logo)
    i.applyIntentContextToScreen(listScreen)
    
    i.activeScreen = listScreen
    return nil
}
```

### Example: Conditional Transition

```go
// transitionToFormScreen transitions to add or edit form based on skill.
func (i *ManageSkillsIntent) transitionToFormScreen(skill *domain.Skill) tea.Cmd {
    // Conditional state based on parameter
    if skill == nil {
        i.currentState = SkillsStateAdd
    } else {
        i.currentState = SkillsStateEdit
    }
    
    formScreen := NewSkillFormScreenFromIntent(skill, i.GetThemeManager())
    i.applyIntentContextToScreen(formScreen)
    i.activeScreen = formScreen
    return nil
}
```

### Example: Async Initialization

```go
// transitionToLoadingState starts async data load.
func (i *YourIntent) transitionToLoadingState() tea.Cmd {
    i.currentState = StateLoading
    i.isLoading = true
    
    // Return async command
    return func() tea.Msg {
        data, err := loadData()
        return DataLoadedMsg{Data: data, Error: err}
    }
}
```

---

## Pattern: Inline Transitions

### When Acceptable

Inline transitions are acceptable for:
- Simple state changes with no initialization
- One-off transitions that won't be reused
- Intents with ≤3 states

### Structure

```go
func (i *SimpleIntent) Update(msg tea.Msg) tea.Cmd {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "enter":
            // Inline transition - acceptable for simple cases
            i.currentState = StateNext
            return nil
        }
    }
    return nil
}
```

### Keep It Simple

```go
// ✅ Good: Simple inline transition
i.state.currentState = GenerateCVStateSelectAudience
i.activeScreen = nil
return nil

// ❌ Bad: Complex inline transition (extract this!)
i.currentState = SkillsStateList
listScreen := NewSkillsListScreenFromIntent(i.skills, i.GetThemeManager())
if skillScreen, ok := listScreen.(*skills_screens.SkillsListScreen); ok {
    if i.eventCounts != nil {
        skillScreen.SetEventCounts(i.eventCounts)
    }
}
i.applyIntentContextToScreen(listScreen)
i.activeScreen = listScreen
return nil
// This should be: return i.transitionToListScreen()
```

---

## Reference Implementation

**Intent**: `ManageSkillsIntent` (complex, screen-based)  
**File**: `internal/cli/intents/manage_skills_intent.go`  
**Lines**: 2189-2264

### Available Helpers

```go
// Screen transitions
func (i *ManageSkillsIntent) transitionToListScreen() tea.Cmd
func (i *ManageSkillsIntent) transitionToDetailScreen() tea.Cmd
func (i *ManageSkillsIntent) transitionToFormScreen(skill *domain.Skill) tea.Cmd
func (i *ManageSkillsIntent) transitionToDeleteScreen(skill *domain.Skill) tea.Cmd

// Error handling transition
func (i *ManageSkillsIntent) handleErrorInternal(err error) tea.Cmd
```

### Why This Works

1. **Consistent Structure**: All helpers follow same pattern
2. **Screen Architecture**: Each transition creates and configures screen
3. **Context Application**: Common setup (terminal, theme, logo) applied in helper
4. **Clear Naming**: Method names clearly indicate target state
5. **Reusable**: Multiple code paths can call same transition

---

## Best Practices

### DO

✅ **Name transitions clearly**: `transitionToXState()` or `transitionToXScreen()`  
✅ **Document purpose**: Add godoc explaining when/why transition occurs  
✅ **Set state first**: `i.currentState = NewState` before initialization  
✅ **Clean up**: Nil out old resources, clear flags  
✅ **Return commands**: Return initialization commands (or nil)  
✅ **Be consistent**: If you extract one, extract them all  
✅ **Test transitions**: Write unit tests for each helper

### DON'T

❌ **Mix patterns**: Don't use both extracted and inline in same intent  
❌ **Forget cleanup**: Always clear old state to prevent leaks  
❌ **Skip initialization**: Ensure new state is fully ready  
❌ **Hardcode state**: Use const values: `StateList` not `"list"`  
❌ **Return wrong state**: Ensure state matches what you set  
❌ **Ignore errors**: Handle initialization errors properly

---

## Common Mistakes

### Mistake 1: Partial Transition

```go
// ❌ BAD: State set but screen not initialized
func (i *Intent) transitionToDetail() tea.Cmd {
    i.currentState = StateDetail
    // FORGOT: i.activeScreen = NewDetailScreen()
    return nil
}
```

**Fix**: Always initialize all resources for new state.

### Mistake 2: Forgetting Cleanup

```go
// ❌ BAD: Old modal not cleared
func (i *Intent) transitionToList() tea.Cmd {
    i.currentState = StateList
    // FORGOT: i.editModal = nil
    return nil
}
```

**Fix**: Always clean up old state resources.

### Mistake 3: Inconsistent Pattern

```go
// ❌ BAD: Mix of extracted and inline
func (i *Intent) Update(msg tea.Msg) tea.Cmd {
    switch msg.String() {
    case "a":
        return i.transitionToAddScreen()  // Extracted
    case "l":
        i.currentState = StateList        // Inline
        return nil
    }
}
```

**Fix**: Pick one pattern and stick with it.

### Mistake 4: No Return Path

```go
// ❌ BAD: Sets state but doesn't return command
func (i *Intent) transitionToLoading() {  // Missing tea.Cmd return
    i.currentState = StateLoading
    // FORGOT: return i.loadData()
}
```

**Fix**: Always return `tea.Cmd` (or nil).

---

## Template: Creating New Transition Helpers

### Step 1: Define the method

```go
// transitionToXState transitions to state X with full initialization.
// Use this when [explain when this transition occurs].
func (i *YourIntent) transitionToXState() tea.Cmd {
    // Implementation
}
```

### Step 2: Set the state

```go
i.currentState = StateX
```

### Step 3: Initialize resources

```go
// Create any resources needed for new state
i.xResource = createXResource()

// Configure resources
i.xResource.SetTheme(i.Theme())
```

### Step 4: Clean up old state

```go
// Clear old resources
i.oldResource = nil
i.oldModal = nil
```

### Step 5: Return initialization command

```go
// Return initialization command or nil
return i.xResource.Init()
```

---

## Testing Transition Helpers

### Unit Test Template

```go
Describe("transitionToXState", func() {
    It("should set current state to X", func() {
        intent := setupIntent()
        
        intent.transitionToXState()
        
        Expect(intent.currentState).To(Equal(StateX))
    })
    
    It("should initialize X resources", func() {
        intent := setupIntent()
        
        intent.transitionToXState()
        
        Expect(intent.xResource).NotTo(BeNil())
    })
    
    It("should clean up old resources", func() {
        intent := setupIntent()
        intent.oldResource = &Resource{}
        
        intent.transitionToXState()
        
        Expect(intent.oldResource).To(BeNil())
    })
    
    It("should return initialization command", func() {
        intent := setupIntent()
        
        cmd := intent.transitionToXState()
        
        // Verify command behavior
        Expect(cmd).NotTo(BeNil())
    })
})
```

---

## Related Documentation

- **Intent Patterns Library**: `docs/development/INTENT_PATTERNS_LIBRARY.md`
- **Common Intent Patterns**: `docs/development/COMMON_INTENT_PATTERNS.md`
- **TUI Developer Guide**: `docs/TUI_DEVELOPER_GUIDE.md`
- **State Matrix**: `docs/STATE_MATRIX.md`

---

## Summary

**Use extracted transition helpers when**:
- Intent has 4+ states
- Transitions require initialization or cleanup
- Same transition is used from multiple places
- Transition logic is >15 lines

**Use inline transitions when**:
- Intent has ≤3 states
- Transitions are simple (<10 lines)
- Each transition is unique
- Intent uses different architecture

**Reference**: ManageSkillsIntent (`manage_skills_intent.go` lines 2189-2264) for best-practice implementation.

---

**Last Updated**: 2026-01-14  
**Next Review**: After next complex intent implementation  
**Owner**: KaRiya TUI Architecture Team

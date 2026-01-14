---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# KaRiya CLI Workflow Refinement Strategy

## Fundamental Architectural Principles

### 1. Intent-Driven Navigation
Instead of screen-to-screen transitions, we'll model user intents explicitly.

```go
type NavIntent int

const (
    IntentStartCapture NavIntent = iota
    IntentAdvanceStep
    IntentGoBack
    IntentViewEvent
    IntentExit
)

type NavigationState struct {
    CurrentScreen Screen
    AllowedIntents []NavIntent
}

func (ns NavigationState) Apply(intent NavIntent) (NavigationState, error) {
    // Pure transition with explicit intent validation
}
```

### 2. Immutable State Transitions
All state changes will be pure functions returning new states.

```go
// Instead of mutating methods
func (state FormState) WithFieldValue(field string, value interface{}) FormState {
    newState := state.Clone()
    newState.Fields[field] = value
    return newState
}
```

### 3. Explicit Input Dispatch
Clear ownership and consumption of input events.

```go
type InputDispatcher struct {
    modalLayer    InputHandler
    screenLayer   InputHandler
    componentLayer InputHandler
}

func (d *InputDispatcher) Dispatch(msg tea.Msg) tea.Cmd {
    // Hierarchical input handling
    // Modal > Screen > Component priority
}
```

### 4. Screen-Local Focus Management
Focus managed autonomously within each screen.

```go
type FocusPolicy interface {
    Next(current int, totalFields int) int
    Prev(current int, totalFields int) int
}

type DefaultCyclicFocusPolicy struct{}

func (p DefaultCyclicFocusPolicy) Next(current, total int) int {
    return (current + 1) % total
}
```

### 5. Constrained Error Management
Explicit error scoping and lifecycle.

```go
type ErrorScope int
type ErrorLifetime int

const (
    ScopeField ErrorScope = iota
    ScopeScreen
    ScopeFlow
)

const (
    LifetimeUntilFixed ErrorLifetime = iota
    LifetimeUntilNavigate
    LifetimeImmediate
)

type UserError struct {
    Message   string
    Scope     ErrorScope
    Lifetime  ErrorLifetime
}
```

### 6. Validation with Explicit Phases
Deterministic, immutable form validation.

```go
type ValidationPhase int

const (
    PhaseField ValidationPhase = iota
    PhaseCrossField
    PhaseSubmission
)

type ValidationRule struct {
    Phase     ValidationPhase
    DependsOn []string
    Validate  func(FormSnapshot) []UserError
}
```

## Migration Strategy

### Strict Migration Rules
1. **Hard Freeze**: Migrated screens cannot use legacy navigation
2. **Incremental Replacement**: One screen/component at a time
3. **Comprehensive Test Coverage**: 100% coverage before migration

### Recommended Migration Order
1. Input Dispatch System
2. Navigation Intent Model
3. Screen-Local Focus Management
4. Form Validation Framework
5. Error Management System

## Technical Debt Reduction Metrics
- Cyclomatic complexity reduction
- Fewer global state mutations
- More predictable state transitions
- Increased test coverage
- Reduced debugging complexity

## Risks and Mitigations
- **Over-abstraction**: Maintain clear, simple interfaces
- **Performance**: Benchmark and optimize pure functions
- **Developer Learning Curve**: Comprehensive documentation

## Success Criteria
- Predictable navigation flows
- Explicit user intent modeling
- Immutable, traceable state changes
- Reduced coupling between components

---

**Strategy Version**: 2.0
**Focused on**: Constraint-Driven Design
**Core Principle**: Make Illegal States Unrepresentable


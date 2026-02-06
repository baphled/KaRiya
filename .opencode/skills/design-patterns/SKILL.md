---
name: design-patterns
description: Apply appropriate design patterns during refactoring - know when and why to use each pattern
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

## What I do

Guide selection and implementation of design patterns. Patterns are tools - use them when they solve a problem, not to show off.

## When to use me

- During REFACTOR phase of TDD
- When code smells indicate a pattern would help
- When reviewing code for structural improvements
- When designing new components

## Golden Rule

**Don't use a pattern unless you have the problem it solves.**

Premature pattern application creates complexity. Wait until you feel the pain, then apply the pattern.

## Creational Patterns

### Factory Method
**Problem:** Need to create objects without specifying exact class
**When:** Object creation logic is complex or varies by context

```go
// Interface
type Screen interface {
    View() string
    Update(tea.Msg) (tea.Cmd, ScreenResult)
}

// Factory
func NewScreen(screenType string, ctx *ScreenContext) Screen {
    switch screenType {
    case "list":
        return NewListScreen(ctx)
    case "detail":
        return NewDetailScreen(ctx)
    case "form":
        return NewFormScreen(ctx)
    default:
        return NewErrorScreen("unknown screen type")
    }
}
```

**KaRiya usage:** Screen creation in intents

### Builder
**Problem:** Complex object construction with many optional parameters
**When:** Constructor has >4 parameters or many combinations

```go
// Builder pattern
type ModalBuilder struct {
    modal *Modal
}

func NewModalBuilder(theme *themes.Theme) *ModalBuilder {
    return &ModalBuilder{
        modal: &Modal{theme: theme},
    }
}

func (b *ModalBuilder) Title(t string) *ModalBuilder {
    b.modal.title = t
    return b
}

func (b *ModalBuilder) Content(c string) *ModalBuilder {
    b.modal.content = c
    return b
}

func (b *ModalBuilder) WithCancel() *ModalBuilder {
    b.modal.showCancel = true
    return b
}

func (b *ModalBuilder) Build() *Modal {
    return b.modal
}

// Usage
modal := NewModalBuilder(theme).
    Title("Confirm Delete").
    Content("Are you sure?").
    WithCancel().
    Build()
```

**KaRiya usage:** UIKit containers, complex form builders

### Functional Options
**Problem:** Same as Builder but more idiomatic Go
**When:** Configurable objects with sensible defaults

```go
type Option func(*Config)

func WithTimeout(d time.Duration) Option {
    return func(c *Config) {
        c.timeout = d
    }
}

func WithRetries(n int) Option {
    return func(c *Config) {
        c.retries = n
    }
}

func NewClient(opts ...Option) *Client {
    cfg := &Config{
        timeout: 30 * time.Second,  // default
        retries: 3,                  // default
    }
    for _, opt := range opts {
        opt(cfg)
    }
    return &Client{config: cfg}
}

// Usage
client := NewClient(
    WithTimeout(10 * time.Second),
    WithRetries(5),
)
```

**KaRiya usage:** Service configuration, behavior options

## Structural Patterns

### Adapter
**Problem:** Interface mismatch between components
**When:** Integrating external libraries or legacy code

```go
// External library uses different interface
type ExternalLogger interface {
    Log(level int, msg string)
}

// Our interface
type Logger interface {
    Info(msg string)
    Error(msg string)
}

// Adapter
type LoggerAdapter struct {
    external ExternalLogger
}

func (a *LoggerAdapter) Info(msg string) {
    a.external.Log(0, msg)
}

func (a *LoggerAdapter) Error(msg string) {
    a.external.Log(2, msg)
}
```

**KaRiya usage:** Wrapping huh forms, external service clients

### Composite
**Problem:** Treat individual objects and compositions uniformly
**When:** Tree structures, nested UI components

```go
type Component interface {
    Render() string
}

type Container struct {
    children []Component
}

func (c *Container) Render() string {
    var result strings.Builder
    for _, child := range c.children {
        result.WriteString(child.Render())
    }
    return result.String()
}

type Text struct {
    content string
}

func (t *Text) Render() string {
    return t.content
}

// Usage - both are Components
container := &Container{
    children: []Component{
        &Text{content: "Hello"},
        &Container{children: []Component{&Text{content: "Nested"}}},
    },
}
```

**KaRiya usage:** UIKit layout composition

### Decorator
**Problem:** Add behavior without modifying original
**When:** Cross-cutting concerns (logging, caching, auth)

```go
type Repository interface {
    Get(id string) (*Entity, error)
}

// Base implementation
type SQLRepository struct {
    db *sql.DB
}

func (r *SQLRepository) Get(id string) (*Entity, error) {
    // SQL query
}

// Logging decorator
type LoggingRepository struct {
    wrapped Repository
    logger  Logger
}

func (r *LoggingRepository) Get(id string) (*Entity, error) {
    r.logger.Info("Getting entity: " + id)
    result, err := r.wrapped.Get(id)
    if err != nil {
        r.logger.Error("Failed to get entity: " + err.Error())
    }
    return result, err
}

// Caching decorator
type CachingRepository struct {
    wrapped Repository
    cache   map[string]*Entity
}

func (r *CachingRepository) Get(id string) (*Entity, error) {
    if cached, ok := r.cache[id]; ok {
        return cached, nil
    }
    result, err := r.wrapped.Get(id)
    if err == nil {
        r.cache[id] = result
    }
    return result, err
}

// Usage - stack decorators
repo := &CachingRepository{
    wrapped: &LoggingRepository{
        wrapped: &SQLRepository{db: db},
        logger:  logger,
    },
    cache: make(map[string]*Entity),
}
```

**KaRiya usage:** Repository wrappers, middleware

## Behavioral Patterns

### Strategy
**Problem:** Need interchangeable algorithms
**When:** Multiple ways to do the same thing

```go
type SortStrategy interface {
    Sort(items []Item) []Item
}

type DateSortStrategy struct{}

func (s *DateSortStrategy) Sort(items []Item) []Item {
    sort.Slice(items, func(i, j int) bool {
        return items[i].Date.Before(items[j].Date)
    })
    return items
}

type NameSortStrategy struct{}

func (s *NameSortStrategy) Sort(items []Item) []Item {
    sort.Slice(items, func(i, j int) bool {
        return items[i].Name < items[j].Name
    })
    return items
}

type ItemList struct {
    items    []Item
    strategy SortStrategy
}

func (l *ItemList) SetStrategy(s SortStrategy) {
    l.strategy = s
}

func (l *ItemList) SortedItems() []Item {
    return l.strategy.Sort(l.items)
}
```

**KaRiya usage:** Capture strategies, filter strategies

### State
**Problem:** Object behavior changes based on internal state
**When:** State machines, workflow management

```go
type IntentState interface {
    HandleKey(key string) IntentState
    View() string
}

type ListState struct {
    intent *MyIntent
}

func (s *ListState) HandleKey(key string) IntentState {
    switch key {
    case "enter":
        return &DetailState{intent: s.intent}
    case "n":
        return &FormState{intent: s.intent}
    }
    return s
}

func (s *ListState) View() string {
    return s.intent.listScreen.View()
}

type DetailState struct {
    intent *MyIntent
}

func (s *DetailState) HandleKey(key string) IntentState {
    switch key {
    case "esc":
        return &ListState{intent: s.intent}
    case "e":
        return &EditState{intent: s.intent}
    }
    return s
}
```

**KaRiya usage:** Intent state management (though KaRiya uses enum-based states)

### Observer
**Problem:** Objects need to react to changes in other objects
**When:** Event-driven systems, UI updates

```go
type Observer interface {
    OnEvent(event Event)
}

type Subject struct {
    observers []Observer
}

func (s *Subject) Subscribe(o Observer) {
    s.observers = append(s.observers, o)
}

func (s *Subject) Notify(event Event) {
    for _, o := range s.observers {
        o.OnEvent(event)
    }
}
```

**KaRiya usage:** Bubble Tea's Cmd/Msg pattern is observer-like

### Command
**Problem:** Encapsulate requests as objects
**When:** Undo/redo, queuing operations, macro recording

```go
type Command interface {
    Execute() error
    Undo() error
}

type CreateEventCommand struct {
    service *EventService
    event   *Event
    created *Event  // for undo
}

func (c *CreateEventCommand) Execute() error {
    created, err := c.service.Create(c.event)
    if err != nil {
        return err
    }
    c.created = created
    return nil
}

func (c *CreateEventCommand) Undo() error {
    if c.created == nil {
        return errors.New("nothing to undo")
    }
    return c.service.Delete(c.created.ID)
}
```

**KaRiya usage:** Bubble Tea's `tea.Cmd` is a command pattern

## Pattern Selection Guide

### By Code Smell

| Smell | Consider Pattern |
|-------|------------------|
| Long parameter list | Builder, Functional Options |
| Switch on type | Strategy, State |
| Duplicate code | Template Method, Strategy |
| Complex conditionals | State, Strategy |
| Tight coupling | Adapter, Dependency Injection |
| God object | Facade, decomposition |
| Feature envy | Move method, Visitor |

### By Problem

| Problem | Pattern |
|---------|---------|
| Create objects flexibly | Factory, Builder |
| Add behavior dynamically | Decorator |
| Handle multiple algorithms | Strategy |
| Manage complex state | State |
| Decouple components | Observer, Mediator |
| Wrap incompatible interface | Adapter |
| Simplify complex subsystem | Facade |

## Anti-Patterns to Avoid

### Pattern Abuse
```go
// WRONG - Factory for single type
func NewUser() *User {
    return &User{}
}

// RIGHT - Just use constructor
user := &User{}
```

### Premature Abstraction
```go
// WRONG - Strategy for one algorithm
type OnlySortStrategy struct{}

// RIGHT - Wait until you need multiple strategies
func sortItems(items []Item) []Item {
    // direct implementation
}
```

### Over-Engineering
```go
// WRONG - 5 patterns for simple CRUD
type AbstractFactoryBuilderStrategyDecoratorFacade struct{}

// RIGHT - Simple, direct code
func SaveUser(u *User) error {
    return db.Save(u)
}
```

## Related Skills

- `clean-code` - When to refactor
- `tdd-workflow` - Apply patterns during REFACTOR phase
- `refactor` - Refactoring techniques
- `go-expert` - Go-specific patterns
- `architecture` - Architectural patterns
- `concurrency` - Concurrent patterns

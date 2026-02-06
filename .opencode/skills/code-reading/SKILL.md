---
name: code-reading
description: Understand unfamiliar codebases quickly - navigation strategies, building mental models, finding entry points
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

## What I do

Guide efficient reading and understanding of unfamiliar code. Reading code is a skill - learn to build mental models quickly and navigate effectively.

## When to use me

- Joining a new project
- Working in unfamiliar part of codebase
- Understanding dependencies
- Debugging unfamiliar code
- Code review of new areas

## Core Principles

1. **Top-down then bottom-up** - Understand structure before details
2. **Follow the data** - Trace how data flows through the system
3. **Tests tell the truth** - Tests show intended behaviour
4. **Don't read everything** - Read what's relevant to your task
5. **Build incrementally** - Understanding grows over time

## Reading Strategies

### The 5-Minute Overview

```markdown
## Quick Orientation Checklist

1. [ ] README - What does this do?
2. [ ] Directory structure - How is it organised?
3. [ ] Entry points - Where does execution start?
4. [ ] Dependencies - What does it use?
5. [ ] Tests - What's the expected behaviour?

TIME: 5-10 minutes max, then dive deeper
```

### Directory Structure Analysis

```bash
# Get overview
tree -L 2 -d

# KaRiya structure:
internal/
├── cli/           # TUI application
│   ├── intents/   # Workflows (start here for features)
│   ├── screens/   # UI components
│   ├── behaviors/ # Reusable behaviours
│   └── uikit/     # UI primitives
├── domain/        # Business entities (start here for data)
├── service/       # Business logic
└── repository/    # Data access

# Pattern: Find the layer relevant to your task
```

### Finding Entry Points

```bash
# Main function
grep -r "func main" --include="*.go"

# HTTP handlers
grep -r "http.Handle\|gin\.\|echo\." --include="*.go"

# CLI commands
grep -r "cobra\|flag\." --include="*.go"

# Test entry points
grep -r "func Test\|var _ = Describe" --include="*_test.go"
```

## Building Mental Models

### The Onion Model

```
Layer 1: WHAT (5 min)
├── What does this project do?
├── Who uses it?
└── What problem does it solve?

Layer 2: HOW - Structure (15 min)
├── How is code organised?
├── What are the main components?
└── How do they connect?

Layer 3: HOW - Flow (30 min)
├── How does data flow?
├── What's the happy path?
└── Where are the boundaries?

Layer 4: WHY - Details (as needed)
├── Why this design choice?
├── What are the edge cases?
└── What are the gotchas?
```

### Component Mapping

```markdown
## Component Map: KaRiya Timeline Feature

### User Action: View Timeline
1. App router → browsetimeline intent
2. Intent creates ListScreen
3. ListScreen uses TableBehavior
4. TableBehavior queries service
5. Service calls repository
6. Repository returns Events
7. Screen renders via View()

### Key Files
- Entry: intents/browsetimeline/intent.go
- Screen: screens/timeline/list_screen.go
- Data: domain/career/event.go
- Storage: repository/event_repository.go
```

### Tracing Data Flow

```go
// Start with the data structure
type Event struct {
    ID   string
    Name string
    Date time.Time
}

// Find where it's created
grep -r "Event{" --include="*.go"

// Find where it's used
grep -r "*Event" --include="*.go"

// Find where it's stored
grep -r "Save.*Event\|Insert.*Event" --include="*.go"

// Find where it's retrieved
grep -r "Find.*Event\|Get.*Event" --include="*.go"
```

## Navigation Techniques

### IDE Navigation

```markdown
## Essential Shortcuts (Go/VS Code)

Go to Definition     : F12 / Cmd+Click
Find References      : Shift+F12
Go to File           : Cmd+P
Go to Symbol         : Cmd+Shift+O
Find in Files        : Cmd+Shift+F
Peek Definition      : Alt+F12
Go Back              : Ctrl+- (after jumping)
```

### Command Line Navigation

```bash
# Find type definition
grep -rn "type Event struct" --include="*.go"

# Find implementations of interface
grep -rn "func.*EventRepository" --include="*.go"

# Find usages
grep -rn "\.CreateEvent\|\.GetEvent" --include="*.go"

# Find tests for a function
grep -rn "TestCreateEvent\|CreateEvent" --include="*_test.go"
```

### Call Hierarchy

```bash
# Who calls this function?
grep -rn "CreateEvent(" --include="*.go" | grep -v "func.*CreateEvent"

# What does this function call?
# Read the function, note calls, then find their definitions
```

## Reading Tests

### Tests as Documentation

```go
// Tests tell you:
// 1. What the function should do
// 2. What inputs are valid
// 3. What outputs to expect
// 4. What edge cases exist

var _ = Describe("EventService", func() {
    Describe("Create", func() {
        It("creates event with valid input", func() {
            // Shows: happy path, expected input format
        })
        
        It("returns error for empty title", func() {
            // Shows: validation rule
        })
        
        It("generates ID if not provided", func() {
            // Shows: automatic behaviour
        })
    })
})
```

### Finding Relevant Tests

```bash
# Tests for a specific file
find . -name "*_test.go" | xargs grep -l "EventService"

# E2E tests for a feature
find . -name "*_e2e_test.go" | xargs grep -l "timeline"

# All tests for a package
ls internal/service/*_test.go
```

## Understanding Patterns

### Recognising Common Patterns

```go
// Repository Pattern
type Repository interface {
    Find(id string) (*Entity, error)
    Save(entity *Entity) error
}

// Service Pattern
type Service struct {
    repo Repository  // Dependency injection
}

// Factory Pattern
func NewService(repo Repository) *Service {
    return &Service{repo: repo}
}

// Intent/Screen Pattern (KaRiya-specific)
type Intent struct {
    state  IntentState
    screen Screen
}

func (i *Intent) Update(msg tea.Msg) tea.Cmd {
    // State machine
}
```

### Spotting Architecture

```bash
# Find dependency injection
grep -rn "func New.*(" --include="*.go" | head -20

# Find interfaces
grep -rn "type.*interface" --include="*.go"

# Find configuration
grep -rn "type.*Config\|config\." --include="*.go"
```

## Reading Strategies by Goal

### Goal: Fix a Bug

```markdown
1. Reproduce the bug
2. Find the symptoms in code (error messages, UI)
3. Trace backwards to cause
4. Read tests to understand expected behaviour
5. Read minimal code to fix
```

### Goal: Add a Feature

```markdown
1. Find similar existing feature
2. Trace its implementation top-to-bottom
3. Identify patterns used
4. Map where your feature fits
5. Read tests for similar feature
```

### Goal: Code Review

```markdown
1. Understand the PR description
2. Look at file changes (what areas?)
3. Read tests first (what should it do?)
4. Read implementation (does it do that?)
5. Check edge cases
```

### Goal: Onboarding

```markdown
Week 1: Orientation
- README and docs
- Directory structure
- Main entry points
- One simple bug fix

Week 2: Shallow understanding
- Major components
- Data models
- Common patterns
- One small feature

Week 3+: Deeper understanding
- Edge cases
- Historical context (git blame)
- Architecture decisions
- Larger features
```

## Avoiding Common Mistakes

### DON'T

- Try to understand everything at once
- Read linearly like a book
- Ignore tests
- Skip the README
- Assume without verifying

### DO

- Start with a specific goal
- Follow the data/control flow
- Read tests alongside code
- Use IDE navigation
- Build understanding incrementally

## Questions to Ask

### When Reading a Function

```markdown
- What does it take as input?
- What does it return?
- What side effects does it have?
- What can go wrong?
- Why does it exist?
```

### When Reading a Class/Struct

```markdown
- What data does it hold?
- What operations does it support?
- What are its dependencies?
- How is it constructed?
- Who uses it?
```

### When Reading a Package

```markdown
- What is its purpose?
- What does it export?
- What are its dependencies?
- How does it fit in the architecture?
- What are the entry points?
```

## Related Skills

- `research` - Systematic investigation
- `debug-test` - Debugging techniques
- `architecture` - Understanding structure
- `question-resolver` - Answering questions about code

---
name: refactor
description: Systematic code refactoring with safety nets, incremental changes, and behavior preservation
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

## What I do

Guide systematic refactoring that improves code structure while preserving behavior.

## When to use me

Use this skill when:
- Improving code structure
- Reducing complexity
- Eliminating duplication
- Improving naming
- Extracting abstractions
- Consolidating patterns

## Golden Rule

**Refactoring changes structure, not behavior. Tests must pass before and after.**

## Refactoring Process

### 1. Ensure Safety Net

Before ANY refactoring:

```bash
# Verify tests exist and pass
make test

# Check coverage for the area
go test -cover ./path/to/refactor/...
```

**No tests = No refactoring.** Write tests first.

### 2. Identify the Smell

| Smell | Indicator | Refactoring |
|-------|-----------|-------------|
| Long function | > 30 lines | Extract Method |
| Long parameter list | > 3 params | Introduce Parameter Object |
| Duplicate code | Copy-paste | Extract Method/Function |
| Feature envy | Uses other object's data | Move Method |
| Data clump | Same fields together | Extract Class |
| Primitive obsession | Strings for everything | Value Object |
| Switch statements | Type checking | Polymorphism |
| Parallel inheritance | Matching hierarchies | Collapse Hierarchy |
| Lazy class | Does too little | Inline Class |
| Speculative generality | "Might need someday" | Remove |
| Temporary field | Only sometimes set | Extract Class |
| Message chain | a.b().c().d() | Hide Delegate |
| Middle man | Just delegates | Remove Middle Man |
| Inappropriate intimacy | Too coupled | Move Method |
| Comments | Explaining bad code | Rename/Extract |

### 3. Choose the Refactoring

#### Extract Method
```go
// Before
func (s *Service) Process(items []*Item) error {
    // validate
    for _, item := range items {
        if item.Name == "" {
            return errors.New("name required")
        }
    }
    // process
    for _, item := range items {
        item.Status = "processed"
        s.repo.Save(item)
    }
    return nil
}

// After
func (s *Service) Process(items []*Item) error {
    if err := s.validateItems(items); err != nil {
        return err
    }
    return s.processItems(items)
}

func (s *Service) validateItems(items []*Item) error {
    for _, item := range items {
        if item.Name == "" {
            return errors.New("name required")
        }
    }
    return nil
}

func (s *Service) processItems(items []*Item) error {
    for _, item := range items {
        item.Status = "processed"
        s.repo.Save(item)
    }
    return nil
}
```

#### Introduce Parameter Object
```go
// Before
func Query(startDate, endDate time.Time, status string, limit, offset int) {}

// After
type QueryParams struct {
    StartDate time.Time
    EndDate   time.Time
    Status    string
    Limit     int
    Offset    int
}

func Query(params QueryParams) {}
```

#### Replace Conditional with Polymorphism
```go
// Before
func (e *Event) Icon() string {
    switch e.Type {
    case "work":
        return "💼"
    case "education":
        return "🎓"
    default:
        return "📌"
    }
}

// After - each type implements EventType interface
type EventType interface {
    Icon() string
}

type WorkEvent struct{}
func (WorkEvent) Icon() string { return "💼" }

type EducationEvent struct{}
func (EducationEvent) Icon() string { return "🎓" }
```

### 4. Execute Incrementally

**Small steps, frequent verification:**

1. Make ONE small change
2. Run tests
3. Commit if green
4. Repeat

```bash
# After each small change
make test

# If tests pass, commit
git add -A && git commit -m "refactor: extract validateItems method"
```

### 5. Verify Behavior Preserved

```bash
# Full test suite
make test

# Coverage should be same or better
go test -cover ./path/...

# No new linter warnings
make check-patterns
```

## Refactoring Catalog

### Composing Methods
- Extract Method
- Inline Method
- Replace Temp with Query
- Split Temporary Variable
- Remove Assignments to Parameters

### Moving Features
- Move Method
- Move Field
- Extract Class
- Inline Class
- Hide Delegate
- Remove Middle Man

### Organizing Data
- Replace Magic Number with Constant
- Encapsulate Field
- Replace Data Value with Object
- Replace Array with Object

### Simplifying Conditionals
- Decompose Conditional
- Consolidate Conditional Expression
- Replace Nested Conditional with Guard Clauses
- Replace Conditional with Polymorphism

### Simplifying Method Calls
- Rename Method
- Add Parameter
- Remove Parameter
- Preserve Whole Object
- Replace Parameter with Method Call
- Introduce Parameter Object

## KaRiya-Specific Refactorings

### Intent Too Large
```
Problem: intent.go > 400 lines
Solution:
1. Extract handlers to handlers.go
2. Extract helpers to helpers.go
3. Extract types to types.go
4. Move screens to screens/{feature}/
```

### Screen in Intent
```
Problem: Screen struct defined in intents/
Solution:
1. Create screens/{feature}/ package
2. Move screen struct
3. Update imports
```

### Deprecated Pattern
```
Problem: Using components.KeyBadge
Solution:
1. Replace with primitives.HelpKeyBadge()
2. Update imports
3. Remove old import
```

## Anti-Patterns

- **Big Bang Refactoring** - Changing everything at once
- **Refactoring Without Tests** - No safety net
- **Refactoring While Adding Features** - Mixed concerns
- **Gold Plating** - Over-engineering during refactor
- **Ignoring Test Failures** - "I'll fix them later"

## Refactoring Checklist

Before:
- [ ] Tests exist and pass
- [ ] Understand current behavior
- [ ] Identified specific smell

During:
- [ ] Small incremental changes
- [ ] Run tests after each change
- [ ] Commit working states

After:
- [ ] All tests pass
- [ ] Coverage maintained
- [ ] No new warnings
- [ ] Code is cleaner

## Related skills

- `clean-code` - Clean code principles
- `architecture` - Structural patterns
- `code-reviewer` - Identify refactoring opportunities
- `tech-debt` - Track refactoring needs

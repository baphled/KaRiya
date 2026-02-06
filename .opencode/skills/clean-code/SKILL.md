---
name: clean-code
description: Write clean, maintainable code following Boy Scout Rule and SOLID principles
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

## What I do

Enforce clean code practices and the Boy Scout Rule: "Leave the code cleaner than you found it."

## When to use me

Use this skill when writing or modifying any code. This should be your default mindset.

## The Boy Scout Rule

**Always leave the code better than you found it.**

When you touch a file:
1. Fix any obvious issues you see (naming, formatting, small refactors)
2. Remove dead code
3. Improve unclear variable/function names
4. Extract magic numbers to named constants
5. Add missing documentation for exports

## Clean Code Principles

### 1. Meaningful Names

```go
// BAD
func p(e *E) error { ... }
d := 86400

// GOOD
func processEvent(event *Event) error { ... }
const secondsPerDay = 86400
```

### 2. Functions Should Do One Thing

```go
// BAD - Does too much
func handleRequest(r *Request) {
    validate(r)
    save(r)
    notify(r)
    log(r)
}

// GOOD - Single responsibility
func handleRequest(r *Request) error {
    if err := r.Validate(); err != nil {
        return err
    }
    return r.Save()
}
```

### 3. No Comments Inside Functions

In this project, comments inside function bodies are **FORBIDDEN**.

```go
// BAD
func process(items []Item) {
    // Filter active items
    active := filterActive(items)
    // Sort by date
    sorted := sortByDate(active)
}

// GOOD - Self-documenting code
func process(items []Item) {
    activeItems := filterActiveItems(items)
    sortedItems := sortItemsByDate(activeItems)
}
```

### 4. DRY - Don't Repeat Yourself

```go
// BAD - Duplication
if user.Role == "admin" {
    // 20 lines of permission logic
}
if user.Role == "superadmin" {
    // Same 20 lines
}

// GOOD - Extract common logic
func hasAdminPermissions(role string) bool {
    return role == "admin" || role == "superadmin"
}
```

### 5. KISS - Keep It Simple

```go
// BAD - Over-engineered
type EventProcessorFactoryBuilder struct { ... }

// GOOD - Simple and direct
func NewEventProcessor(db *DB) *EventProcessor { ... }
```

### 6. YAGNI - You Aren't Gonna Need It

Don't add features "just in case." Only implement what's needed now.

## SOLID Principles

| Principle | Meaning | Example |
|-----------|---------|---------|
| **S**ingle Responsibility | One reason to change | Intent orchestrates, Screen renders |
| **O**pen/Closed | Open for extension, closed for modification | Use interfaces |
| **L**iskov Substitution | Subtypes must be substitutable | Implement full interface contract |
| **I**nterface Segregation | Many specific interfaces | `EventRepository` not `Repository` |
| **D**ependency Inversion | Depend on abstractions | Accept interfaces, return structs |

## Code Smells to Fix

When you encounter these, fix them:

| Smell | Fix |
|-------|-----|
| Long function (>30 lines) | Extract methods |
| Long parameter list (>3) | Use config struct |
| Nested conditionals (>2 deep) | Early returns, extract methods |
| Magic numbers | Named constants |
| Commented-out code | Delete it (git has history) |
| TODO/FIXME comments | Complete the work or create issue |
| Duplicate code | Extract to shared function |
| God object | Split into focused types |

## Refactoring Checklist

Before finishing any change:

- [ ] Names are meaningful and intention-revealing
- [ ] Functions do one thing
- [ ] No duplicate code
- [ ] No magic numbers
- [ ] No comments inside functions
- [ ] No dead code
- [ ] Error handling is explicit
- [ ] Tests cover the change

## Go-Specific Clean Code

### Error Handling
```go
// BAD
result, _ := doSomething()

// GOOD
result, err := doSomething()
if err != nil {
    return fmt.Errorf("doing something: %w", err)
}
```

### Avoid Naked Returns
```go
// BAD
func get() (result string, err error) {
    result = "value"
    return  // What's being returned?
}

// GOOD
func get() (string, error) {
    return "value", nil
}
```

### Use Early Returns
```go
// BAD
func process(x int) error {
    if x > 0 {
        if x < 100 {
            // actual logic
        }
    }
    return nil
}

// GOOD
func process(x int) error {
    if x <= 0 {
        return errors.New("x must be positive")
    }
    if x >= 100 {
        return errors.New("x must be less than 100")
    }
    // actual logic
    return nil
}
```

## Related skills

- `code-reviewer` - Review code for clean code violations
- `architecture` - Ensure architectural cleanliness
- `tdd-workflow` - Clean code through TDD

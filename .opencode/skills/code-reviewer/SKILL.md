---
name: code-reviewer
description: Comprehensive code review covering clean code, architecture, security, and KaRiya patterns
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

## What I do

Perform comprehensive code reviews covering all aspects of code quality.

## When to use me

Use this skill when:
- Reviewing code before commit
- Self-reviewing your changes
- Preparing for PR review
- Auditing existing code

## Review Dimensions

### 1. Correctness
- Does it do what it's supposed to?
- Are edge cases handled?
- Are errors handled properly?

### 2. Clean Code
- Meaningful names?
- Functions do one thing?
- No duplication?
- No magic numbers?

### 3. Architecture
- Correct layer placement?
- Proper dependencies?
- Following patterns?

### 4. Security
- Input validated?
- No SQL injection?
- No hardcoded secrets?
- Errors don't leak info?

### 5. Testing
- Tests exist?
- Coverage adequate (>=95%)?
- Edge cases covered?
- Tests are readable?

### 6. Documentation
- Exports documented?
- Expected/Returns/Side effects?
- No comments inside functions?

## Review Checklist

### Quick Check (Every Change)

```
[ ] Code compiles without warnings
[ ] Tests pass
[ ] No new linter warnings
[ ] Changes match the intent
```

### Clean Code Check

```
[ ] Names are meaningful and intention-revealing
[ ] Functions are small (<30 lines ideal)
[ ] Functions do one thing
[ ] No duplicate code
[ ] No magic numbers (use named constants)
[ ] No commented-out code
[ ] No TODO/FIXME (complete or create issue)
[ ] No comments inside function bodies
[ ] Early returns used (no deep nesting)
```

### Architecture Check

```
[ ] Code is in the correct layer
[ ] Dependencies flow downward only
[ ] No circular dependencies
[ ] Intent uses subdirectory structure
[ ] Screens don't import intents
[ ] huh only used in forms/
[ ] Modals in screens/ or uikit/, not intents/
[ ] State enum in constants.go
[ ] Messages in messages.go
```

### Security Check

```
[ ] User input is validated
[ ] No SQL string concatenation
[ ] File paths are sanitized
[ ] No hardcoded credentials
[ ] Errors don't expose internal details
[ ] Timeouts set for external calls
```

### Testing Check

```
[ ] New code has tests
[ ] Tests cover happy path
[ ] Tests cover error cases
[ ] Tests cover edge cases
[ ] Coverage >= 95% for modified packages
[ ] Tests are readable (Arrange-Act-Assert)
[ ] No test pollution (tests are independent)
```

### Documentation Check

```
[ ] Package has doc.go
[ ] Exported types have godoc comments
[ ] Exported functions have Expected/Returns/Side effects
[ ] No inline comments in function bodies
[ ] Type godoc starts with type name
```

## Review Commands

```bash
# Run all checks
make check-compliance

# Architecture check
make check-intent-architecture

# Pattern check
make check-patterns

# Security scan
make gosec

# Coverage
go test -cover ./path/to/modified/...

# Full CI locally
make ci-local
```

## Common Issues

### Issue: Long Functions

```go
// BAD: 80-line function
func processEvents(events []*Event) error {
    // ... 80 lines of logic
}

// GOOD: Extract to focused functions
func processEvents(events []*Event) error {
    validated := validateEvents(events)
    grouped := groupByCategory(validated)
    return saveGroups(grouped)
}
```

### Issue: Deep Nesting

```go
// BAD: 4 levels deep
if x != nil {
    if x.Value > 0 {
        if x.IsActive {
            if x.HasPermission {
                // do work
            }
        }
    }
}

// GOOD: Early returns
if x == nil {
    return nil
}
if x.Value <= 0 {
    return nil
}
if !x.IsActive || !x.HasPermission {
    return nil
}
// do work
```

### Issue: Magic Numbers

```go
// BAD
if len(items) > 100 {
    items = items[:100]
}

// GOOD
const maxDisplayItems = 100

if len(items) > maxDisplayItems {
    items = items[:maxDisplayItems]
}
```

### Issue: Error Swallowing

```go
// BAD
result, _ := doSomething()

// GOOD
result, err := doSomething()
if err != nil {
    return fmt.Errorf("doing something: %w", err)
}
```

### Issue: Leaky Abstraction

```go
// BAD: huh in intent
import "github.com/charmbracelet/huh"
form := huh.NewForm(...)

// GOOD: Use forms package
import "github.com/baphled/kariya/internal/cli/forms"
form := forms.NewForm(...)
```

## Review Feedback Templates

### Blocking Issue
```
[BLOCKING] {Description}

This needs to be fixed before merge because {reason}.

Suggested fix: {solution}
```

### Suggestion
```
[SUGGESTION] {Description}

Consider {alternative approach} because {benefit}.

Current code works, but this would {improvement}.
```

### Question
```
[QUESTION] {Question}

I'm not sure about {aspect}. Could you clarify {specific thing}?
```

### Praise
```
[NICE] {Description}

Good job on {specific thing}. This {benefit}.
```

## Boy Scout Rule in Review

When reviewing, also note opportunities to leave code cleaner:

```
[CLEANUP] While you're here, consider:
- [ ] Renaming `x` to `eventCount`
- [ ] Extracting the validation to a method
- [ ] Removing the commented-out code on line 45
```

## Related skills

- `clean-code` - Clean code principles
- `architecture` - Architectural patterns
- `security` - Security review
- `check-compliance` - Automated checks

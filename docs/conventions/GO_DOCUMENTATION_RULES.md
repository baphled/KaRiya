# Go Documentation & Comment Rules

This document defines the **mandatory** documentation, comment, and naming standards for Go code in this project. These rules are enforced by `golangci-lint` (via `stylecheck`, `revive`, `godot`) and manual code review.

> **Guiding Principle**: Documentation explains *why and what* at boundaries. Code explains *how* through structure and naming.

---

## 0. Go Idioms & Naming Conventions (MANDATORY)

All code must follow standard Go idioms from [Effective Go](https://go.dev/doc/effective_go) and [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments).

### Package Names

- **Lowercase, single-word names** - no underscores, no mixedCaps
- Package name should describe what the package provides
- Avoid generic names like `util`, `common`, `misc`

```go
// GOOD
package timeline
package career
package repository

// BAD - underscores
package browse_timeline  // Use: browsetimeline or timeline
package fact_management  // Use: factmanagement or facts

// BAD - mixedCaps
package careerService    // Use: career or service
```

### Exported Names

- **MixedCaps** (CamelCase), not underscores
- Acronyms should be all caps: `ID`, `URL`, `HTTP`, `API`, `JSON`, `SQL`

```go
// GOOD
type UserID string
func GetHTTPClient() *http.Client
var DefaultAPIURL = "..."

// BAD
type User_ID string      // Use: UserID
func GetHttpClient()     // Use: GetHTTPClient
var Default_Api_Url      // Use: DefaultAPIURL
```

### Enforcement

- **Linter**: `stylecheck` (ST1003 - package naming)
- **Linter**: `revive` (var-naming)

---

## 1. Package-Level Documentation (MANDATORY)

Every package **must** have a package comment.

### Requirements

- Must start with `Package <name> ...`
- Must describe:
  - The purpose of the package
  - Its responsibilities
  - What it explicitly does *not* handle (if relevant)

### Example

```go
// Package career provides domain models and business logic for career event
// management. It handles event creation, validation, and lifecycle management.
//
// This package does NOT handle:
//   - Persistence (see repository package)
//   - UI rendering (see screens package)
//   - External API communication (see service package)
package career
```

### Violations

- Missing package comment
- Vague or restating the package name only (e.g., `// Package career is for career stuff`)

### Enforcement

- **Linter**: `stylecheck` (ST1000)
- **Linter**: `revive` (package-comments)

---

## 2. Exported Identifiers (STRICT)

All exported functions, methods, structs, interfaces, constants, and variables **must** have doc comments.

### Requirements

- **Start with the identifier name** (e.g., `// ProcessOrder validates...`)
- Describe *intent and contract*, not implementation details
- Use structured sections where applicable (see below)

### Structured Sections

Use these sections when applicable:

| Section | When to Use |
|---------|-------------|
| `Expected:` or `Parameters:` | Document preconditions, parameter constraints |
| `Returns:` | Document return values and their meaning |
| `Side effects:` | Document state changes, I/O operations, events emitted |
| `Errors:` | Document non-trivial error conditions |

### Example

```go
// ProcessOrder validates and executes an order.
//
// Expected:
//   - order must be non-nil
//   - order.Items must not be empty
//
// Returns:
//   - Receipt on success
//   - error if validation or execution fails
//
// Side effects:
//   - Persists order to database
//   - Emits OrderProcessed event
func ProcessOrder(order *Order) (*Receipt, error) {
    // ...
}
```

### Interface Example

```go
// EventRepository defines the contract for event persistence.
//
// Implementations must:
//   - Be safe for concurrent access
//   - Return ErrNotFound when an event does not exist
//   - Validate events before persistence
//
// Implementations must NOT:
//   - Perform business logic validation
//   - Emit events (caller's responsibility)
type EventRepository interface {
    // Save persists an event to storage.
    //
    // Returns:
    //   - nil on success
    //   - ErrDuplicateID if event with same ID exists
    //   - ErrValidation if event fails persistence validation
    Save(ctx context.Context, event *Event) error

    // FindByID retrieves an event by its unique identifier.
    //
    // Returns:
    //   - Event and nil error on success
    //   - nil and ErrNotFound if event does not exist
    FindByID(ctx context.Context, id string) (*Event, error)
}
```

### Violations

- Missing doc comment on exported identifier
- Comment that only restates the name (e.g., `// Add adds two numbers`)
- Missing required sections when applicable

### Enforcement

- **Linter**: `stylecheck` (ST1020, ST1021, ST1022)
- **Linter**: `revive` (exported)

---

## 3. Inline Comments Are PROHIBITED

**Inline comments are NOT allowed.** This includes:

- Trailing comments on the same line as code
- Comments inside function bodies
- Comments explaining control flow, conditionals, or implementation steps

### Disallowed

```go
// BAD: Trailing comment
x := foo() // explain foo

// BAD: Comment inside function body explaining implementation
func Process(items []Item) {
    // check if valid
    if valid {
        doThing()
    }
    
    // iterate through items
    for _, item := range items {
        // process each item
        handle(item)
    }
}
```

### Allowed

- Package comments
- File-level comments (e.g., copyright headers)
- Doc comments attached to declarations
- Linter directives (`//nolint`) when unavoidable

### Why?

Inline comments often indicate:
1. **Code that should be extracted** - If you need to explain what code does, extract it into a well-named function.
2. **Poor naming** - If variable/function names need explanation, rename them.
3. **Excessive complexity** - If logic requires a comment, simplify it.

### Enforcement

- **Manual Review**: Inline comments are flagged during code review.
- **Note**: No automated linter catches all inline comments; reviewers must enforce this rule.

---

## 4. Code Must Be Self-Explanatory

If a comment appears necessary, refactor the code instead:

| Instead of... | Do this |
|---------------|---------|
| Comment explaining a condition | Extract to a well-named function |
| Comment explaining a variable | Rename the variable |
| Comment explaining a block of code | Extract to a function with a descriptive name |
| Comment explaining complex logic | Introduce a type or simplify the logic |

### Example

**Not acceptable:**
```go
if attempts > 3 { // max retries
    return ErrMaxRetriesExceeded
}
```

**Preferred:**
```go
if exceedsRetryLimit(attempts) {
    return ErrMaxRetriesExceeded
}

// In a helper or as a constant:
const maxRetries = 3

func exceedsRetryLimit(attempts int) bool {
    return attempts > maxRetries
}
```

---

## 5. Private Code Rules

- Private (unexported) functions do **not** require doc comments.
- Private code must remain readable **without comments**.
- Inline comments are **still prohibited** in private code.

If private code is unclear:
1. **Refactor** - Break down complex logic.
2. **Rename** - Use descriptive names.
3. **Extract** - Create well-named helper functions.

---

## 6. Interfaces Define Behaviour Contracts

Interfaces must be documented with:

- Behavioural guarantees (what implementations must do)
- Constraints on implementations (thread safety, error handling)
- Error semantics (what errors mean)

**Comments on interfaces are preferred over comments on implementations.** The interface doc comment is the source of truth; implementations should not duplicate the documentation.

### Example

```go
// Validator checks whether a value meets specific criteria.
//
// Implementations must:
//   - Be stateless and safe for concurrent use
//   - Return nil for valid values
//   - Return a descriptive error for invalid values
//
// Implementations must NOT:
//   - Modify the input value
//   - Perform I/O operations
type Validator interface {
    Validate(value interface{}) error
}
```

---

## 7. Comment Style

### Rules

- Doc comments must end with a period.
- Use complete sentences.
- Place comments above the code they describe, not beside it.

### Enforcement

- **Linter**: `godot` (checks comments end with a period)

---

## 8. Enforcement Summary

| Rule | Linter | Check |
|------|--------|-------|
| Package comments present | `stylecheck` | ST1000 |
| Package comments present | `revive` | package-comments |
| Exported function docs | `stylecheck` | ST1020 |
| Exported type docs | `stylecheck` | ST1021 |
| Exported var/const docs | `stylecheck` | ST1022 |
| Exported identifier docs | `revive` | exported |
| Comments end with period | `godot` | - |
| No TODO/FIXME/HACK/XXX | `godox` | - |
| Inline comments | **Manual Review** | - |

### Running Enforcement

```bash
# Run all linters including documentation checks
make golangci-lint

# Or directly
golangci-lint run ./...
```

---

## 9. Legacy Code Exceptions

Legacy code in certain directories has relaxed enforcement. This is **temporary** to allow gradual migration.

**Legacy directories** (stylecheck/godot relaxed):
- `internal/cli/components/`
- `internal/cli/terminal/`
- `internal/cli/screens/`
- `internal/cli/intents/`
- `internal/cli/models/`
- And others (see `.golangci.yml` exclusions)

**New code** in these directories **should** follow documentation rules. The exclusions exist only for pre-existing code.

---

## 10. AI Agent Behaviour

When violations are found, agents must:

1. **Fail validation explicitly**, or
2. **Refactor the code** to comply, without adding inline comments

Agents must **not**:
- Justify violations
- Add explanatory inline comments
- Relax rules for "clarity" or "brevity"

### Refusal Template

```
I cannot proceed with this code.

DOCUMENTATION VIOLATION DETECTED

Violation: [Specific violation, e.g., "Missing doc comment on exported function"]
Rule: [Rule number from this document]
Enforcement: [Linter and check, e.g., "stylecheck ST1020"]

Required correction:
1. [Specific fix needed]

See: docs/conventions/GO_DOCUMENTATION_RULES.md
```

---

## References

- [Effective Go - Commentary](https://go.dev/doc/effective_go#commentary)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [stylecheck documentation](https://staticcheck.dev/docs/checks/)
- [revive rules](https://revive.run/r)

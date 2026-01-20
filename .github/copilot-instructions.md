# Copilot Code Review Instructions

## Project Context

KaRiya is a Go-based TUI (Terminal User Interface) application using the Bubble Tea framework. It follows strict architectural patterns and coding standards.

## Review Focus Areas

### 1. Architecture Compliance

- **Intent Pattern**: All workflows must embed `*BaseIntent` and follow the state machine pattern
- **UIKit Components**: Use `uikit/` components, not legacy `components/`
- **Behaviors**: Use `behaviors.TableBehavior[T]` for tables, not `table.New()` directly
- **Forms**: Use `models.*Form` wrappers in intents, never `*huh.Form` directly

### 2. Theme & Styling

Flag violations of:
- Hardcoded colors (e.g., `lipgloss.Color("#xxx")`) - should use `theme.Primary()`, `theme.Secondary()`, etc.
- Raw lipgloss styling - should use `primitives.Title()`, `primitives.Body()`, etc.
- Custom modal code - should use `feedback.Modal` and `behaviors.RenderModalOverlay()`

### 3. Code Quality

- **TDD**: Tests should exist for new functionality
- **Error Handling**: All errors must be handled, not ignored
- **Naming**: Follow Go naming conventions (camelCase for private, PascalCase for public)
- **Comments**: Exported functions should have doc comments

### 4. Bubble Tea Patterns

- `Init()` should return appropriate initial commands
- `Update()` must handle all relevant messages
- `View()` should be pure (no side effects)
- Use `tea.Batch()` for multiple commands

### 5. Testing Standards

- Table-driven tests preferred
- Test edge cases and error conditions
- Use testify assertions consistently
- Mock external dependencies

## Things to Flag

1. **Breaking Changes**: API modifications, removed functions, changed signatures
2. **Security Issues**: Hardcoded credentials, unsafe input handling
3. **Performance**: Unnecessary allocations, N+1 queries, blocking operations
4. **Accessibility**: Missing keyboard navigation, unclear UI states
5. **Documentation**: Missing or outdated comments on public APIs

## Things to Ignore

- Generated files (`*_gen.go`, `*.pb.go`)
- Vendor directory
- Test fixtures and mock data
- Markdown formatting preferences

## Code Patterns to Enforce

### Correct Intent Structure
```go
type MyIntent struct {
    *intents.BaseIntent
    // ... fields
}

func NewMyIntent() *MyIntent {
    return &MyIntent{
        BaseIntent: intents.NewBaseIntent("my-intent"),
    }
}
```

### Correct Theme Usage
```go
// Good
style := lipgloss.NewStyle().Foreground(theme.Primary())

// Bad - flag this
style := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))
```

### Correct Component Usage
```go
// Good
text := primitives.Title("Hello")
badge := primitives.HelpKeyBadge("?", "Help")

// Bad - flag this
text := lipgloss.NewStyle().Bold(true).Render("Hello")
```

## Review Tone

- Be constructive and specific
- Explain why something should change, not just what
- Suggest concrete alternatives when flagging issues
- Acknowledge good patterns when you see them

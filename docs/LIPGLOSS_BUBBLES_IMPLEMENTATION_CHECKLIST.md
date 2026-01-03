# Lipgloss & Bubbles Implementation Checklist

This checklist helps ensure consistent, professional terminal UI implementation across KaRiya using lipgloss for styling and bubbles for interactive components.

## Pre-Implementation

- [ ] Review the [Lipgloss & Bubbles Guide](./LIPGLOSS_BUBBLES_GUIDE.md)
- [ ] Review the [Terminal UI Styling Reference](./TERMINAL_UI_STYLING_REFERENCE.md)
- [ ] Review the [Enhanced Capture Example](../internal/cli/intents/enhanced_capture_example.go)
- [ ] Understand the existing [Styles Package](../internal/cli/styles/styles.go)
- [ ] Understand the existing [Components Package](../internal/cli/components/)

## Component Development Checklist

### 1. Style Setup

- [ ] All styles defined in `internal/cli/styles/styles.go`
- [ ] No inline style definitions in component code
- [ ] Use existing color constants from styles package
- [ ] Document any new colors added to palette
- [ ] Test color contrast for accessibility
- [ ] Verify colors work in both light and dark terminals

### 2. Component Structure

- [ ] Component implements proper Bubble Tea model interface
- [ ] Component has clear Init() method
- [ ] Component has Update(msg tea.Msg) tea.Cmd method
- [ ] Component has View() string method
- [ ] Component manages its own state cleanly
- [ ] Component uses composition over duplication

### 3. Interactive Components

For models using bubbles components (textinput, list, etc.):

- [ ] Bubbles component properly initialized
- [ ] Bubbles component styled with styles package colors
- [ ] Focus management implemented clearly
- [ ] Tab navigation between fields works
- [ ] Shift+Tab navigation between fields works
- [ ] Enter key handling implemented
- [ ] Escape key handling implemented
- [ ] Keyboard shortcuts documented

### 4. Form Implementation

If implementing a form:

- [ ] All form fields use bubbles textinput or similar
- [ ] Form fields properly focused/blurred
- [ ] Form validation implemented
- [ ] Validation errors displayed clearly
- [ ] Error colors use ColorError
- [ ] Focus indicator shows which field is active
- [ ] Tab order is logical
- [ ] Required fields marked clearly

### 5. Layout & Spacing

- [ ] Uses CardContainer for cards
- [ ] Uses FormContainer for forms
- [ ] Uses ListContainer for lists
- [ ] Uses ModalContainer for modals
- [ ] Consistent padding/margin throughout
- [ ] Responsive to terminal size changes
- [ ] Handles minimum terminal width (40 chars)
- [ ] Handles maximum useful width (120 chars)
- [ ] Uses lipgloss.JoinVertical/JoinHorizontal appropriately

### 6. Visual Hierarchy

- [ ] Headers use CardHeader style
- [ ] Body text uses CardContent style
- [ ] Footers use CardFooter style
- [ ] Primary actions use ButtonPrimary style
- [ ] Secondary actions use ButtonSecondary style
- [ ] Text hierarchy clear (bold, italic, normal)
- [ ] Color used to indicate status/state
- [ ] Emoji/icons used appropriately for clarity

### 7. User Feedback

- [ ] Success messages use ColorSuccess
- [ ] Error messages use ColorError
- [ ] Warning messages use ColorWarning
- [ ] Info messages use ColorInfo
- [ ] Loading states shown with spinner
- [ ] Progress shown with progress bar
- [ ] Disabled state uses Faint() style
- [ ] Focus state clearly indicated

### 8. Keyboard Navigation

- [ ] All interactive elements keyboard accessible
- [ ] Tab moves to next element
- [ ] Shift+Tab moves to previous element
- [ ] Enter confirms/submits
- [ ] Escape cancels/goes back
- [ ] Ctrl+C quits application
- [ ] Keyboard shortcuts documented in footer
- [ ] Help text available (?)

### 9. Accessibility

- [ ] No color-only differentiation
- [ ] Text labels for all controls
- [ ] High contrast text and backgrounds
- [ ] Focus indicators clearly visible
- [ ] Error messages in text, not just color
- [ ] Keyboard navigation complete
- [ ] No flashing or rapid animations
- [ ] Terminal resize handled gracefully

### 10. Testing

- [ ] Unit tests for component logic
- [ ] Tests for state transitions
- [ ] Tests for keyboard input handling
- [ ] Tests for rendering at different sizes
- [ ] Tests for focus management
- [ ] Tests for validation logic
- [ ] Manual testing with actual terminal
- [ ] >90% code coverage for component

### 11. Documentation

- [ ] Component documented with comments
- [ ] State machine documented
- [ ] Keyboard shortcuts documented in code
- [ ] Usage example provided
- [ ] Style choices documented
- [ ] Integration guide provided
- [ ] Dependencies documented
- [ ] Known limitations documented

### 12. Code Quality

- [ ] Code formatted with gofmt
- [ ] No golangci-lint warnings
- [ ] No go vet warnings
- [ ] Proper error handling
- [ ] No panic() calls
- [ ] Proper resource cleanup
- [ ] No race conditions
- [ ] Follows project conventions

## Integration Checklist

### 1. Intent Integration

If implementing as an Intent:

- [ ] Implements Intent interface properly
- [ ] Init() returns appropriate command
- [ ] Update() handles all message types
- [ ] View() returns styled output
- [ ] Result() returns IntentResult[T]
- [ ] State machine clearly defined
- [ ] Back navigation preserves context
- [ ] Metadata stored for restoration

### 2. Router Integration

- [ ] Intent registered with IntentRouter
- [ ] Factory function creates new instances
- [ ] Result handler implemented
- [ ] Navigation to intent works
- [ ] Navigation from intent works
- [ ] Back navigation works
- [ ] Result callbacks handled
- [ ] Global shortcuts still work

### 3. Styling Consistency

- [ ] Uses existing color palette
- [ ] Follows existing layout patterns
- [ ] Consistent with other intents
- [ ] Button styles match elsewhere
- [ ] Card styles match elsewhere
- [ ] Font styles match elsewhere
- [ ] Spacing matches elsewhere
- [ ] Emoji usage consistent

### 4. Performance

- [ ] Rendering completes in <100ms
- [ ] State transitions in <10ms
- [ ] No memory leaks
- [ ] No CPU spikes
- [ ] Efficient re-rendering
- [ ] Caches styles appropriately
- [ ] Handles large datasets
- [ ] Terminal resize responsive

## Final Validation

### Pre-Merge Checklist

- [ ] All tests passing
- [ ] Coverage >90%
- [ ] No lint warnings
- [ ] No format issues
- [ ] Manual testing complete
- [ ] Keyboard navigation verified
- [ ] Responsive sizing tested
- [ ] Accessibility verified
- [ ] Documentation complete
- [ ] Code review approved

### Testing Scenarios

- [ ] Test with 40-char wide terminal
- [ ] Test with 80-char wide terminal
- [ ] Test with 120-char wide terminal
- [ ] Test with 200-char wide terminal
- [ ] Test with 10-line high terminal
- [ ] Test with 30-line high terminal
- [ ] Test with 50-line high terminal
- [ ] Test rapid window resizing
- [ ] Test all keyboard shortcuts
- [ ] Test all validation paths
- [ ] Test error conditions
- [ ] Test with mouse input disabled

## Common Issues & Solutions

### Issue: Inline Style Definitions

**Problem**: Styles defined in component view methods
```go
// ❌ Bad
func (m Model) View() string {
    style := lipgloss.NewStyle().Foreground(lipgloss.Color("#fff"))
    return style.Render("text")
}
```

**Solution**: Move to styles package
```go
// ✅ Good - in styles/styles.go
var MyStyle = lipgloss.NewStyle().Foreground(ColorTextPrimary)

// ✅ Good - in component
func (m Model) View() string {
    return styles.MyStyle.Render("text")
}
```

### Issue: No Focus Management

**Problem**: Multiple inputs but no clear focus handling
```go
// ❌ Bad
func (m Model) Update(msg tea.Msg) tea.Cmd {
    m.input1, _ = m.input1.Update(msg)
    m.input2, _ = m.input2.Update(msg)
}
```

**Solution**: Implement proper focus management
```go
// ✅ Good
func (m *Model) updateFocus() {
    m.input1.Blur()
    m.input2.Blur()
    if m.focused == 0 {
        m.input1.Focus()
    } else {
        m.input2.Focus()
    }
}
```

### Issue: Fixed Terminal Size

**Problem**: Component assumes fixed terminal dimensions
```go
// ❌ Bad
const Width = 80
const Height = 24
```

**Solution**: Respond to window size messages
```go
// ✅ Good
case tea.WindowSizeMsg:
    m.width = msg.Width
    m.height = msg.Height
```

### Issue: Color-Only Differentiation

**Problem**: Using only color to indicate state
```go
// ❌ Bad
if hasError {
    return ColorError.Render(text)
}
```

**Solution**: Use text labels + color
```go
// ✅ Good
if hasError {
    return "✗ " + ColorError.Render(text)
}
```

### Issue: Inconsistent Component Usage

**Problem**: Duplicating card rendering logic
```go
// ❌ Bad
func (m Model) viewCard1() string {
    return lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        Render(content1)
}

func (m Model) viewCard2() string {
    return lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        Render(content2)
}
```

**Solution**: Use CardContainer component
```go
// ✅ Good
func (m Model) viewCard1() string {
    return components.NewCardContainer().
        SetBody(content1).Render()
}

func (m Model) viewCard2() string {
    return components.NewCardContainer().
        SetBody(content2).Render()
}
```

## Review Checklist for Code Reviewers

- [ ] All styles from styles package (not inline)
- [ ] Using existing component containers
- [ ] Proper focus management implemented
- [ ] Responsive to terminal size
- [ ] Keyboard navigation complete
- [ ] Error handling clear
- [ ] Tests >90% coverage
- [ ] Documentation complete
- [ ] No accessibility issues
- [ ] Performance acceptable
- [ ] Follows project conventions
- [ ] Code quality high (lint, format, vet)

## Resources

- [Lipgloss Documentation](https://github.com/charmbracelet/lipgloss)
- [Bubbles Components](https://github.com/charmbracelet/bubbles)
- [Bubble Tea Framework](https://github.com/charmbracelet/bubbletea)
- [KaRiya Styles Package](../internal/cli/styles/styles.go)
- [KaRiya Components](../internal/cli/components/)
- [Enhanced Example](../internal/cli/intents/enhanced_capture_example.go)
- [Terminal UI Guide](./LIPGLOSS_BUBBLES_GUIDE.md)
- [Styling Reference](./TERMINAL_UI_STYLING_REFERENCE.md)

## Quick Links

- **Style Definition**: `internal/cli/styles/styles.go`
- **Reusable Components**: `internal/cli/components/`
- **Color Palette**: `internal/cli/styles/constants_export.go`
- **Example Implementation**: `examples/enhanced_capture_example.go`
- **Intent Template**: `internal/cli/intents/capture_event.go`

---

**Version**: 1.0
**Last Updated**: 2026-01-03
**Maintained By**: KaRiya Development Team


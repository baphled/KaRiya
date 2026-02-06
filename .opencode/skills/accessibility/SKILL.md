# Accessibility Skill

## Identity

You are an accessibility expert ensuring terminal applications are usable by everyone, including users with visual impairments, motor disabilities, and cognitive differences. You understand WCAG principles adapted for TUI contexts.

## Core Principles (POUR)

### Perceivable
- Information must be presentable in ways users can perceive
- Don't rely solely on color to convey meaning
- Ensure sufficient contrast
- Support screen readers where possible

### Operable
- All functionality available via keyboard
- No time limits that can't be extended
- Clear navigation and focus management
- Predictable, consistent interactions

### Understandable
- Clear, simple language
- Predictable behavior
- Help users avoid and correct mistakes
- Consistent navigation patterns

### Robust
- Work across different terminal emulators
- Degrade gracefully with limited capabilities
- Don't require specific terminal features

## TUI Accessibility Guidelines

### Keyboard Navigation
```go
// REQUIRED: Full keyboard support
// - Tab/Shift+Tab for form fields
// - Arrow keys for lists
// - Enter for selection
// - Esc for cancel/back
// - Shortcuts announced in help

// WRONG: Mouse-only interactions (impossible in TUI anyway)
```

### Color and Contrast
```go
// CORRECT - Don't rely solely on color
// Use symbols + color for status
"[x] Completed"  // Checkmark + green
"[ ] Pending"    // Empty + default
"[!] Error"      // Exclamation + red

// WRONG - Color only
// Green text for success, red for error (colorblind users can't distinguish)
```

### Text and Readability
```go
// CORRECT - Clear, concise labels
"Save changes"
"Delete event"
"Cancel"

// WRONG - Ambiguous labels
"OK"
"Yes"
"Submit"
```

### Focus Management
```go
// CORRECT - Clear focus indicators
// Selected item has visible highlight
// Focus moves logically through interface
// Focus trapped in modals

// WRONG - Invisible focus
// Focus jumps unpredictably
// No indication of current position
```

## KaRiya Accessibility Patterns

### Status Indicators
```go
// Use text + color + symbol
type Status int

const (
    StatusPending Status = iota
    StatusInProgress
    StatusCompleted
    StatusError
)

func (s Status) String() string {
    switch s {
    case StatusPending:
        return "[ ] Pending"
    case StatusInProgress:
        return "[~] In Progress"
    case StatusCompleted:
        return "[x] Completed"
    case StatusError:
        return "[!] Error"
    }
}
```

### Help Text
```go
// Always show available actions
footer := []string{
    primitives.HelpKeyBadge("enter", "Select"),
    primitives.HelpKeyBadge("esc", "Back"),
    primitives.HelpKeyBadge("?", "Help"),
}
```

### Error Messages
```go
// CORRECT - Clear, actionable errors
"Event title is required. Press Tab to return to the title field."

// WRONG - Vague errors
"Invalid input"
"Error occurred"
```

### Confirmation Dialogs
```go
// CORRECT - Clear consequences
"Delete event 'Team Meeting'? This cannot be undone."
"[Delete] [Cancel]"

// WRONG - Unclear consequences
"Are you sure?"
"[Yes] [No]"
```

## Accessibility Checklist

When reviewing for accessibility:
- [ ] All interactive elements reachable via keyboard
- [ ] Focus order is logical (top-to-bottom, left-to-right)
- [ ] Focus indicator is clearly visible
- [ ] Color is not the only means of conveying information
- [ ] Text has sufficient contrast (use theme system)
- [ ] Error messages are clear and actionable
- [ ] Help is available (? key)
- [ ] Destructive actions require confirmation
- [ ] No flashing/blinking content
- [ ] Language is clear and concise

## Testing Accessibility

### Manual Testing
1. **Keyboard only** - Navigate entire app without mouse
2. **High contrast** - Test with high contrast terminal theme
3. **Monochrome** - Test with colors disabled
4. **Screen magnification** - Test at larger font sizes
5. **Screen reader** - Test with terminal screen reader if available

### Questions to Ask
- Can a keyboard-only user complete all tasks?
- Can a colorblind user understand all information?
- Can a user with low vision read all text?
- Are error states clearly communicated?
- Is the current focus always visible?

## Common Accessibility Issues

| Issue | Solution |
|-------|----------|
| Color-only status | Add symbols (x, !, ~, etc.) |
| Hidden shortcuts | Show in help footer |
| Unclear errors | Provide specific, actionable messages |
| Lost focus | Manage focus explicitly after state changes |
| Ambiguous buttons | Use descriptive labels ("Delete event" not "OK") |
| No confirmation | Add modal for destructive actions |

## Related Skills
- `ui-design` - Visual design principles
- `ux-design` - Interaction patterns
- `bubble-tea-expert` - Implementation details

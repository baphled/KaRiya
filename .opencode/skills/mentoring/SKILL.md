---
name: mentoring
description: Teaching and guiding junior engineers, code review coaching, knowledge transfer
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
  always_active: false
---

## What I do

Guide effective mentoring practices for teaching junior engineers, providing constructive code review feedback, and transferring knowledge within the team.

## When to use me

Use this skill when:
- Reviewing code from junior engineers
- Explaining architectural decisions
- Onboarding new team members
- Writing documentation for learning
- Providing feedback that teaches

## Mentoring Principles

### 1. Guide, Don't Gatekeep

```
❌ "This is wrong, fix it."
✅ "This approach has a potential issue with X. 
    Consider Y because [reason]. 
    Here's an example: [code]"
```

### 2. Explain the Why

```
❌ "Use TableBehavior here."
✅ "Use TableBehavior here because it:
    - Handles keyboard navigation consistently
    - Provides filtering/sorting out of the box
    - Follows the pattern we use across all list views
    
    See: internal/cli/screens/timeline/list_screen.go for an example"
```

### 3. Ask Questions Before Prescribing

```
❌ "You should use the repository pattern."
✅ "I notice you're calling the database directly here. 
    What made you choose this approach? 
    Have you considered using the repository pattern?
    
    It might help because [reason]."
```

## Code Review for Learning

### Feedback Levels

| Level | When | Example |
|-------|------|---------|
| **Praise** | Good patterns | "Great use of TableBehavior here - exactly the right pattern" |
| **Suggest** | Could be better | "Consider extracting this to a helper function for reusability" |
| **Request** | Should change | "This needs to use the theme system instead of hardcoded colours" |
| **Block** | Must change | "This violates our architecture - screens cannot import intents" |

### Constructive Feedback Template

```markdown
## Observation
[What you noticed]

## Impact
[Why it matters - performance, maintainability, consistency]

## Suggestion
[What to do instead]

## Example
[Code snippet showing the better approach]

## Learning Resource
[Link to docs, guide, or example in codebase]
```

### Example Reviews

**Teaching Patterns:**

```markdown
### Feedback: Form Handling

I noticed you're using `*huh.Form` directly in the intent.

**Why this matters:** Our architecture keeps the `huh` library encapsulated 
in the `forms/` package. This helps us:
- Maintain consistent form styling
- Swap form libraries if needed
- Keep intents focused on orchestration

**Suggestion:** Use a FormScreen instead:

```go
// Instead of:
type MyIntent struct {
    form *huh.Form
}

// Use:
type MyIntent struct {
    formScreen *myfeature.FormScreen
}
```

**Reference:** See `screens/timeline/capture_screen.go` for the pattern.
```

**Teaching Architecture:**

```markdown
### Feedback: Layer Violation

This screen imports from `intents/`:

```go
import "github.com/baphled/kariya/internal/cli/intents"
```

**Why this matters:** Our layer hierarchy requires:
- Intents → import → Screens (allowed)
- Screens → import → Intents (forbidden)

This prevents circular dependencies and keeps screens reusable.

**Solution:** Use ScreenResult to communicate back:

```go
// In screen
return nil, screens.NewNavigateResult("detail", selectedItem)

// Intent handles the result and owns the navigation logic
```

**Docs:** See [INTENT_ARCHITECTURE_GUIDE.md](docs/INTENT_ARCHITECTURE_GUIDE.md)
```

## Onboarding Guide

### First Week Tasks

| Day | Focus | Activities |
|-----|-------|------------|
| 1 | Environment | Setup, run tests, explore codebase |
| 2 | Architecture | Read AGENTS.md, architecture docs |
| 3 | Small Fix | Pick up a `good-first-issue` bug |
| 4 | Code Review | Review a small PR, ask questions |
| 5 | First Feature | Start small feature with guidance |

### Recommended Reading Order

1. `AGENTS.md` - Project rules and conventions
2. `docs/development/ARCHITECTURE_OVERVIEW.md` - System design
3. `docs/INTENT_ARCHITECTURE_GUIDE.md` - TUI patterns
4. `docs/UIKIT_GUIDE.md` - Component library
5. Example code in `examples/` directory

### First Contribution Guide

```markdown
## Your First KaRiya Contribution

### 1. Pick an Issue

Look for issues labelled `good-first-issue`:
```bash
gh issue list --label "good-first-issue"
```

### 2. Understand the Context

Before coding:
- Read related code
- Understand the intent/screen involved
- Check for existing patterns

Ask questions! Use: `/ask-codebase [your question]`

### 3. Follow TDD

```bash
make tdd-red    # Write failing test first
make tdd-green  # Make it pass
make tdd-refactor  # Clean up
```

### 4. Create PR

```bash
make ai-commit FILE=/tmp/commit.txt
git push -u origin feature/your-feature
gh pr create --base next
```

### 5. Request Review

Tag a mentor for review. Expect feedback - it's how we learn!
```

## Knowledge Transfer Techniques

### Pair Programming Sessions

```markdown
## Pairing Session: [Topic]

**Goal:** [What we'll learn]
**Duration:** 1-2 hours
**Prerequisites:** [What to read first]

### Session Structure

1. **Context** (10 min) - Explain the problem
2. **Navigate** (10 min) - Junior drives, explores codebase
3. **Design** (15 min) - Discuss approach together
4. **Implement** (45 min) - Trade driver/navigator
5. **Review** (10 min) - What did we learn?

### Notes

[Capture key learnings, patterns discovered, questions for follow-up]
```

### Documentation for Learning

When documenting, include:

```go
// ProcessEvents handles batch event processing.
//
// This function is the preferred way to process events because
// it handles batching efficiently - processing 100 events is
// nearly as fast as processing 1 event due to [reason].
//
// Example usage:
//
//     results := ProcessEvents(ctx, events)
//     for _, r := range results {
//         if r.Error != nil {
//             log.Error("failed", "event_id", r.ID, "error", r.Error)
//         }
//     }
//
// Common mistakes:
//   - Don't call ProcessEvent in a loop; use this instead
//   - Always check results for partial failures
//
// See also: ProcessEvent (deprecated), EventProcessor interface
func ProcessEvents(ctx context.Context, events []*Event) []Result
```

## Feedback Techniques

### The COIN Model

| Element | Question | Example |
|---------|----------|---------|
| **C**ontext | When/where? | "In the ProcessEvents function..." |
| **O**bservation | What did you see? | "I noticed the error handling ignores..." |
| **I**mpact | Why does it matter? | "This could cause silent failures..." |
| **N**ext | What to do? | "Consider wrapping errors with context..." |

### Growth-Focused Feedback

```markdown
### What's Going Well

[Specific praise for good patterns, improvements, or learnings]

### Growth Opportunities

[Areas for improvement, framed as learning opportunities]

### Suggested Next Steps

[Concrete actions for skill development]
```

## Common Teaching Moments

| Situation | Teaching Approach |
|-----------|------------------|
| Hardcoded values | Explain constants, configuration, theme system |
| God function | Demonstrate single responsibility, extraction |
| Missing tests | Walk through TDD workflow together |
| Copy-paste code | Show DRY principle, extraction, generics |
| Architecture violation | Explain layers, show correct pattern |
| Missing error handling | Demonstrate error wrapping, context |

## Mentoring Anti-Patterns

**DON'T:**
- Rewrite their code without explanation
- Say "just do it this way" without why
- Compare to other engineers
- Assume knowledge without checking
- Rush through explanations
- Take over instead of guiding

**DO:**
- Ask what they've tried first
- Explain reasoning behind patterns
- Celebrate progress
- Create safe space for questions
- Follow up on previous feedback
- Share your own learning journey

## Related Skills

- `code-reviewer` - Detailed review practices
- `documentation-writing` - Writing for understanding
- `pair-programming` - Collaborative learning
- `clean-code` - Principles to teach

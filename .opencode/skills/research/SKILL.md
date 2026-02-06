---
name: research
description: Systematic research and investigation for understanding codebases, technologies, and solutions
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

## What I do

Conduct systematic research to understand codebases, evaluate technologies, and find solutions to problems.

## When to use me

Use this skill when:
- Exploring an unfamiliar codebase area
- Evaluating technical approaches
- Investigating how something works
- Finding existing patterns to follow
- Understanding dependencies
- Researching best practices

## Research Principles

1. **Gather before concluding** - Collect evidence first
2. **Primary sources** - Prefer code over documentation
3. **Multiple perspectives** - Look at related implementations
4. **Document findings** - Share what you learn
5. **Question assumptions** - Verify, don't assume

## Research Process

### 1. Define the Question

Be specific about what you need to know:

```
# Vague
"How does the timeline work?"

# Specific
"How does the timeline intent handle filtering events by date range?"
```

### 2. Identify Sources

| Source | When to Use |
|--------|-------------|
| Code | Primary source of truth |
| Tests | Reveals expected behavior |
| Existing patterns | Similar implementations |
| Git history | Why changes were made |
| Documentation | High-level context |

### 3. Systematic Exploration

#### Codebase Exploration

```bash
# Find relevant files
find internal/cli -name "*timeline*" -type f

# Search for patterns
grep -rn "FilterByDate" internal/cli/

# Find usages
grep -rn "TimelineIntent" internal/cli/

# Check git history for context
git log --oneline -20 -- internal/cli/intents/browsetimeline/
```

#### Understanding Flow

1. Start at the entry point
2. Trace the execution path
3. Note key decision points
4. Identify dependencies
5. Map the data flow

#### Pattern Recognition

```bash
# Find similar implementations
grep -rn "TableBehavior" internal/cli/intents/

# See how others solved similar problems
grep -rn "HandleNavigation" internal/cli/screens/
```

### 4. Synthesize Findings

Organize what you learned:

```markdown
## Research: [Topic]

### Question
What I was trying to understand.

### Findings

#### How It Works
- Step 1: ...
- Step 2: ...

#### Key Components
- `ComponentA` - Does X
- `ComponentB` - Does Y

#### Patterns Used
- Pattern 1: Used for...
- Pattern 2: Used for...

### Implications
What this means for the current task.

### Open Questions
Things still unclear.
```

## Research Techniques

### Following the Thread

Start from a known point and trace connections:

```
Entry Point → Handler → Service → Repository → Database
```

### Comparing Implementations

Find 2-3 similar implementations and compare:

```markdown
| Aspect | Implementation A | Implementation B |
|--------|------------------|------------------|
| Pattern used | X | Y |
| Error handling | Returns error | Shows modal |
| State management | Local | Context |
```

### Reverse Engineering

When documentation is lacking:

1. Read the tests - they show expected behavior
2. Add logging temporarily
3. Step through with debugger
4. Check git blame for context

### Dependency Mapping

```markdown
## Component: TimelineIntent

### Depends On
- EventService (for data)
- ListScreen (for display)
- FilterModal (for filtering)

### Used By
- Router (navigation target)
- MainMenu (menu option)
```

## Research Questions by Context

### Understanding Existing Code

- What does this component do?
- Why was it implemented this way?
- What are the inputs and outputs?
- What are the edge cases?
- How is it tested?

### Evaluating Approaches

- What are the options?
- What are the trade-offs?
- What does the codebase already do?
- What's the simplest solution?
- What are the risks?

### Debugging

- When does the problem occur?
- What changed recently?
- What are the preconditions?
- What does the error mean?
- Where else does this pattern exist?

### Planning Implementation

- Where should this code live?
- What patterns should I follow?
- What components can I reuse?
- What are the dependencies?
- How will this be tested?

## Documenting Research

### For Yourself

Keep notes as you go:

```markdown
## Research Log: [Date]

### Goal
What I'm trying to find out.

### Steps Taken
1. Looked at X, found Y
2. Searched for Z, discovered...

### Key Insights
- Important finding 1
- Important finding 2

### Next Steps
- Still need to understand...
```

### For Others

When research reveals something important:

1. Update relevant documentation
2. Add code comments for non-obvious things
3. Create a task if action needed
4. Share in commit message

## Red Flags During Research

| Red Flag | Implication |
|----------|-------------|
| No tests | Behavior undefined, add tests first |
| Conflicting patterns | Tech debt, follow newer pattern |
| No similar examples | Novel problem, research more |
| Circular dependencies | Architecture issue, flag it |
| Magic values | Document or refactor |

## Research Anti-Patterns

- **Assuming** without verifying
- **Skimming** when deep dive needed
- **Stopping** at first answer
- **Ignoring** contradictory evidence
- **Not documenting** findings

## Related skills

- `critical-thinking` - Evaluate findings
- `architecture` - Understand structure
- `debug-test` - Investigate issues
- `tech-debt` - Document findings as debt

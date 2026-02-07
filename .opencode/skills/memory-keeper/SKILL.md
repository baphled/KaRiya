---
name: memory-keeper
description: Store findings, fixes, and common problems in memory for learning and quick resolution
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
  always_active: true
---

## What I do

Build institutional memory by storing findings, fixes, and common problems so we can learn from mistakes and quickly resolve recurring issues.

## When to use me

**Always active.** This skill triggers automatically at two key moments:

### 1. DISCOVERY - When Learning Something New
Triggered when you:
- Discover a validation rule or constraint
- Learn why code exists (from git history, docs, or investigation)
- Find a bug pattern or gotcha
- Understand a business rule or domain logic
- Identify a UI/UX pattern or decision

### 2. CHANGE - When Modifying Behaviour  
Triggered when you:
- Add or change validation rules
- Modify business logic
- Change UI/UX behaviour
- Update architectural patterns
- Make any decision that affects how the system works

**Also use this skill:**
- **After fixing a bug** - Record the problem and solution
- **After debugging** - Store the diagnostic process and findings
- **Before debugging** - Recall similar issues we've solved before
- **During retrospectives** - Store lessons learned
- **When patterns emerge** - Document recurring problems
- **At session start** - Surface relevant memories for context

## Mandatory Memory Triggers

### On Discovery
```
DISCOVERED: [what you learned]
CONTEXT: [where/how you found it]
IMPLICATION: [why it matters]

→ Store as memory entity with observations
```

### On Change
```
CHANGED: [what changed]
FROM: [old behaviour]
TO: [new behaviour]  
REASON: [why the change was made]
IMPACT: [what this affects]

→ Store as memory entity with observations
→ Update any existing related entities
```

## Memory Storage via MCP

This skill uses the MCP memory server to store and retrieve memories. The MCP server provides:
- Persistent storage across sessions
- Semantic search for finding relevant memories
- Categorisation by scope (project/domain/tool/language)
- Cross-project learning

## Memory Categories

### Project-Specific (KaRiya)

**MUST capture:**
- **Validation rules** - e.g., "Event text minimum 10 characters"
- **Business rules** - e.g., "Skills must have valid category from constants"
- **UI/UX decisions** - e.g., "Confirm dialogs default to Cancel (right arrow for Submit)"
- **Architectural decisions** - e.g., "Screens never import intents"
- **Component patterns** - e.g., "Use SubmitSkill() helper to bypass form complexity in tests"
- **Domain constraints** - e.g., "Valid skill categories: backend, frontend, devops..."

**Entity naming convention:**
- `KaRiya [Component] [Type]` - e.g., "KaRiya Event Validation Rules"
- `KaRiya [Feature] Pattern` - e.g., "KaRiya Skills BDD Testing"
- `KaRiya [Decision] Rationale` - e.g., "KaRiya 10-char Minimum Rationale"

### Domain-Specific
- CLI/TUI patterns and anti-patterns
- Terminal application quirks
- User interaction patterns
- Testing strategies for CLI apps

### Tool-Specific
- GORM quirks and gotchas
- Ginkgo/Gomega patterns
- Bubble Tea common issues
- Git workflow issues
- GitHub Actions problems

### Language-Specific (Go)
- Go idioms and anti-patterns
- Concurrency pitfalls
- Common compilation errors
- Performance patterns

## Workflow

### Recording a Memory

After fixing a bug or solving a problem:

1. **Identify the memory scope:**
   - Is it KaRiya-specific? → `scope:project`
   - Is it a CLI/TUI pattern? → `scope:domain`
   - Is it a GORM/Ginkgo issue? → `scope:tool`
   - Is it a Go language issue? → `scope:language`

2. **Structure the memory:**
   ```
   Problem: [Clear description of the issue]
   Context: [When/where it occurs]
   Symptoms: [How it manifested]
   Root Cause: [Why it happened]
   Solution: [How we fixed it]
   Prevention: [How to avoid it in future]
   Tags: [searchable keywords]
   ```

3. **Store via MCP memory tools:**
   Use the MCP memory server's storage tool to persist the memory.

### Recalling a Memory

Before starting debugging or when encountering a familiar issue:

1. **Query the memory store:**
   Use the MCP memory server's retrieval tool with relevant keywords.

2. **Review similar issues:**
   Check if we've solved this before or something similar.

3. **Apply or adapt the solution:**
   Use the stored solution as a starting point.

4. **Update if needed:**
   If the situation differs, record the variation.

## Memory Templates

### Bug Fix Memory
```markdown
**Type:** Bug Fix
**Scope:** [project|domain|tool|language]
**Component:** [specific area affected]

**Problem:**
[Clear description of the bug]

**Symptoms:**
- [Symptom 1]
- [Symptom 2]

**Root Cause:**
[Why the bug occurred - technical explanation]

**Solution:**
[How we fixed it - specific code or approach]

**Files Changed:**
- [file:line]

**Prevention:**
- [How to avoid this in future]
- [Linter rule / test pattern / documentation]

**Tags:**
[keyword1, keyword2, keyword3]

**Related Issues:**
- [Link to PR/issue if applicable]
```

### Pattern/Gotcha Memory
```markdown
**Type:** Pattern
**Scope:** [project|domain|tool|language]
**Tool/Area:** [specific tool or area]

**Pattern Name:**
[Concise name for the pattern]

**Context:**
[When does this pattern apply?]

**Problem:**
[What goes wrong if you don't follow this pattern?]

**Correct Approach:**
```go
// Good example
```

**Incorrect Approach:**
```go
// Bad example - avoid this
```

**Rationale:**
[Why the correct approach is better]

**Tags:**
[keyword1, keyword2, keyword3]
```

### Diagnostic Process Memory
```markdown
**Type:** Diagnostic
**Scope:** [project|domain|tool|language]
**Issue Type:** [test failure|runtime error|compilation error|performance]

**Initial Symptoms:**
[What we observed first]

**Diagnostic Steps:**
1. [First thing we checked]
   - Result: [what we found]
2. [Second thing we checked]
   - Result: [what we found]
3. [Root cause identified]

**Tools Used:**
- [debugging tool 1]
- [debugging tool 2]

**Key Insight:**
[What clue led us to the solution?]

**Solution:**
[How we resolved it]

**Time Spent:**
[How long it took - helps prioritise prevention]

**Tags:**
[keyword1, keyword2, keyword3]
```

## Integration with Other Skills

### With `debug-test`
When debugging failing tests:
1. **Recall** memories tagged with test failure symptoms
2. Check if we've seen this pattern before
3. After fixing, **record** the solution for next time

### With `retrospective`
During retrospectives:
1. **Extract** key learnings from the retro
2. **Store** as memories with appropriate scope
3. **Link** to action items and decisions

### With `session-start`
At session start:
1. **Surface** recent memories relevant to current work
2. **Review** common issues for the area you're working in
3. Provide context for the session

### With `incident-response`
During/after incidents:
1. **Record** the incident details
2. **Store** diagnostic process and solution
3. **Link** to post-mortem document

## Common Memory Queries

### "Have we seen this test failure before?"
```
Query: test failure, [test name], [error message keywords]
Scope: project, tool:ginkgo
```

### "How do we handle this GORM pattern?"
```
Query: GORM, [pattern name], repository
Scope: tool:gorm
```

### "What's the correct way to structure this screen?"
```
Query: screen structure, bubble tea, [screen type]
Scope: domain:tui, project
```

### "Why does this Go pattern fail?"
```
Query: [pattern name], concurrency, [error type]
Scope: language:go
```

## Memory Hygiene

### Regular Reviews
- **Weekly:** Review memories added this week
- **Monthly:** Look for patterns across multiple memories
- **Quarterly:** Archive outdated memories, update stale ones

### Quality Checks
- Are solutions still valid?
- Are tags accurate and helpful?
- Is the scope correct?
- Are there duplicates?

### Consolidation
When similar memories accumulate:
1. Identify the common pattern
2. Create a consolidated memory
3. Link individual cases as examples
4. Consider creating a lint rule or test pattern

## Anti-Patterns

**DON'T:**
- Store vague or incomplete information
- Skip recording because "I'll remember"
- Over-categorise (use tags liberally instead)
- Store sensitive information (credentials, keys, PII)
- Record without context (always include why/how)

**DO:**
- Be specific and detailed
- Record immediately while fresh
- Include code examples
- Tag generously for discoverability
- Link to related resources (PRs, docs, issues)

## Metrics

Track memory effectiveness:
- **Recall rate:** How often do we find relevant memories?
- **Resolution time:** Does memory reduce debugging time?
- **Recurrence:** Are issues recurring less?
- **Coverage:** Do we have memories for common areas?

## Related Skills

- `debug-test` - Debug failures using stored memories
- `retrospective` - Extract and store learnings
- `incident-response` - Record incident solutions
- `session-start` - Surface relevant context
- `code-reviewer` - Reference common pitfalls from memory
- `prove-correctness` - Use memories to guide testing

## Examples

### Example 1: Recording a Bug Fix

```
After fixing a race condition in the event list screen:

**Type:** Bug Fix
**Scope:** project
**Component:** UI/EventList

**Problem:**
Event list screen crashes when events are updated while rendering.

**Symptoms:**
- Panic: concurrent map read and write
- Only happens under load
- Tests pass with -race flag failure

**Root Cause:**
EventList model reads from shared state without mutex protection.
Multiple goroutines updating the model concurrently.

**Solution:**
Added sync.RWMutex to EventListModel:
- internal/cli/intents/career/ui/event_list_model.go:23
- Wrap reads with RLock()
- Wrap writes with Lock()

**Prevention:**
- Always use mutex for shared state in Bubble Tea models
- Run `make test-race` before committing
- Added test case with concurrent updates

**Tags:**
race condition, bubble tea, concurrent access, event list, mutex

**Related Issues:**
- PR #234
```

### Example 2: Recalling Before Debugging

```
Symptom: Test fails with "Found more than one test suite file"

Query memory:
- "ginkgo multiple suite files"
- "test suite file error"

Retrieved memory:
"GORM Repository Test Suite Conflict"
- Each package needs exactly ONE *_suite_test.go
- Check for duplicate RunSpecs calls
- Solution: internal/domain/repositories/career/career_suite_test.go:15

Applied solution → Fixed in 2 minutes vs 30 minutes of investigation
```

## Notes

- Memories complement but don't replace documentation
- Documentation = general knowledge, Memory = specific experiences
- Documentation = what should happen, Memory = what actually happened
- Use memories to identify documentation gaps

# Memory Keeper Skill

Build institutional memory by storing findings, fixes, and common problems for learning and quick resolution.

## Overview

The `memory-keeper` skill helps you:
- 🧠 **Remember** solutions to problems you've solved before
- 🔍 **Recall** relevant knowledge when encountering familiar issues
- 📚 **Learn** from past mistakes and successes
- ⚡ **Resolve** recurring issues faster
- 🎯 **Prevent** repeating the same mistakes

## Quick Start

### Recording a Memory After Fixing a Bug

```
1. Load the skill: mcp_skill(name="memory-keeper")
2. Use Bug Fix template to structure the memory
3. Store via MCP memory server
4. Tag appropriately for future discovery
```

### Recalling a Memory Before Debugging

```
1. Load the skill: mcp_skill(name="memory-keeper")
2. Query memory store with relevant keywords
3. Review similar past issues
4. Apply or adapt the solution
```

## Memory Scopes

| Scope | Description | Examples |
|-------|-------------|----------|
| **project** | KaRiya-specific | Architecture decisions, component bugs |
| **domain** | CLI/TUI patterns | Terminal interaction, screen patterns |
| **tool** | Tool-specific | GORM, Ginkgo, Bubble Tea gotchas |
| **language** | Go-specific | Concurrency, idioms, performance |

## Memory Types

| Type | Template | Use Case |
|------|----------|----------|
| **bug_fix** | Bug Fix Memory | After fixing any bug |
| **pattern** | Pattern/Gotcha Memory | Documenting best practices |
| **diagnostic** | Diagnostic Process Memory | Complex debugging sessions |

## Integration Points

### With Other Skills

- **debug-test**: Query memories before debugging test failures
- **retrospective**: Store learnings from retrospectives
- **session-start**: Surface relevant context at session start
- **incident-response**: Record incident solutions
- **code-reviewer**: Reference common pitfalls

### With Development Workflow

```
session-start → check memories for context
     ↓
  develop
     ↓
encounter issue → query memories
     ↓
   debug
     ↓
fix issue → record memory
     ↓
code review → reference memories
     ↓
retrospective → extract learnings → store memories
```

## Files

- **SKILL.md**: Complete skill specification
- **USAGE.md**: Detailed usage guide with examples
- **README.md**: This file

## MCP Memory Server

This skill requires an MCP memory server with these capabilities:
- **Storage**: Persist structured memories
- **Retrieval**: Semantic search across memories
- **Tagging**: Categorise and filter memories
- **Metadata**: Store context (scope, type, timestamps)

## Example Memories

### Bug Fix: Race Condition
```
Stored a race condition fix in EventListModel
Tags: race-condition, bubble-tea, concurrency
Retrieval: Query "race condition event list"
Time saved on recurrence: ~30 minutes
```

### Pattern: GORM Preload
```
Stored correct GORM preload pattern
Tags: gorm, preload, relationships
Retrieval: Query "gorm relationships not loading"
Time saved on recurrence: ~15 minutes
```

### Diagnostic: Test Flakiness
```
Stored diagnostic process for flaky tests
Tags: test-flakiness, ginkgo, timing
Retrieval: Query "test fails intermittently"
Time saved on recurrence: ~45 minutes
```

## Benefits

### For Developers
- Reduce time debugging familiar issues
- Build expertise faster
- Share knowledge automatically
- Prevent regression of old bugs

### For Projects
- Institutional knowledge survives team changes
- Patterns become actionable documentation
- Common issues get resolved faster
- Quality improves through learning

### For Teams
- Collective learning across sessions
- Consistent solutions to common problems
- Better onboarding for new team members
- Reduced duplicate debugging effort

## Best Practices

### DO
✅ Record immediately while details are fresh
✅ Be specific with problem descriptions
✅ Include code examples
✅ Tag generously for discoverability
✅ Link to related resources (PRs, docs)

### DON'T
❌ Store vague or incomplete information
❌ Skip recording because "you'll remember"
❌ Store sensitive information (credentials, PII)
❌ Record without context
❌ Over-categorise (use tags instead)

## Metrics to Track

Monitor memory effectiveness:
- **Recall rate**: How often do we find relevant memories?
- **Resolution time**: Reduction in debugging time
- **Recurrence**: Decrease in repeated issues
- **Coverage**: Memory breadth across codebase areas

## Future Enhancements

Planned improvements:
- Auto-categorisation with ML
- Pattern detection across memories
- Proactive suggestions based on context
- Memory analytics dashboard
- Cross-project memory search
- Auto-generate docs from frequently-used memories

## Contributing

To improve this skill:
1. Use it regularly and note what works/doesn't
2. Suggest new memory types or templates
3. Share effective tagging strategies
4. Report integration issues
5. Propose new search patterns

## License

MIT

## Related Skills

- `debug-test` - Debug failing tests
- `retrospective` - Post-mortems and learning
- `incident-response` - Handle production issues
- `session-start` - Session initialisation
- `code-reviewer` - Code review patterns
- `prove-correctness` - Testing strategies

## Support

For issues or questions:
1. Check USAGE.md for detailed examples
2. Review memory templates in SKILL.md
3. Verify MCP memory server is configured
4. Check memory query syntax

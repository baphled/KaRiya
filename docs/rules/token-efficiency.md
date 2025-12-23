# Token Efficiency Guidelines

## Purpose

Keep conversations focused, token usage low, and work efficient. This guide helps you avoid token bloat and maintain productive AI-assisted development sessions.

---

## Token Budget Awareness

### Token Thresholds

- **0-20k tokens**: ✅ Healthy - Continue normally
- **20k-50k tokens**: ⚠️ Warning - Start being more concise
- **50k-100k tokens**: 🔶 High - Consider fresh conversation soon
- **>100k tokens**: 🔴 Critical - Start fresh conversation NOW

### Current Token Count

Always visible in your interface. Check regularly:
- After long conversations
- Before starting new major task
- When responses feel slow

---

## Token Conservation Strategies

### 1. Use Tools Over Text ⭐

**❌ Token Wasteful:**
```
Can you show me the contents of internal/service/career/service.go?
[AI pastes entire 200-line file]
```

**✅ Token Efficient:**
```
view:internal/service/career/service.go
[AI uses tool, references specific lines]
```

**Best Practices:**
- Use `view` to read files
- Use `grep` to search code
- Use `ls` to explore structure
- Use `glob` to find patterns
- Don't ask AI to paste code you can view

### 2. Be Concise

**❌ Verbose:**
```
I was wondering if you could please help me understand how the
classification system works. I'd like to know about the different
categories and how they are determined. Could you explain this to
me in detail with examples?
```

**✅ Concise:**
```
Explain classification system: categories and determination logic.
```

**Tips:**
- Skip pleasantries in follow-ups
- Use bullet points
- State requirements directly
- Avoid redundant context

### 3. Batch Operations

**❌ Sequential (Multiple Roundtrips):**
```
User: Check if file X exists
AI: [checks] Yes
User: Now read its contents
AI: [reads] Here it is
User: Now check file Y
AI: [checks] Here it is
```

**✅ Batched (Single Roundtrip):**
```
User: Check files X and Y, read both if they exist
AI: [checks both, reads both]
```

**Apply to:**
- File operations
- Test runs
- Code reviews
- Multiple related checks

### 4. Reference, Don't Repeat

**❌ Repetitive:**
```
As I mentioned earlier, we need to follow atomic commit guidelines.
These guidelines state that each commit should have one logical change.
This means we need to ensure each commit...
```

**✅ Referenced:**
```
Following atomic commit guidelines (see earlier message).
```

**Tips:**
- Reference previous context
- Don't re-explain established patterns
- Assume understanding of repeated concepts

### 5. Focus on Deltas

**❌ Full State:**
```
The entire CareerEvent struct has these fields:
ID, Text, Date, Company, Project, Tags, Categories, CreatedAt, UpdatedAt.
Now I'm adding a new field Duration...
```

**✅ Delta Only:**
```
Adding Duration field to CareerEvent struct.
```

**Apply to:**
- Code changes
- Status updates
- Test results
- Documentation updates

---

## Anti-Patterns (Avoid These)

### 🚫 Asking AI to Repeat Code

**Bad:**
```
Can you show me the updated version of the entire file?
```

**Good:**
```
Apply these changes to service.go [specific changes]
```

### 🚫 Verbose Explanations for Simple Changes

**Bad:**
```
I need to add a new method. This method will take a parameter
and return a value. It's important because...
[5 paragraphs of explanation]
```

**Good:**
```
Add GetEventsByTag(tag string) method to service.
```

### 🚫 Re-explaining Context Each Time

**Bad:**
```
As we discussed, the project uses DDD architecture with domain,
service, and repository layers. The domain layer contains...
Now I need to add a new field...
```

**Good:**
```
Add field to domain layer following established DDD pattern.
```

### 🚫 Copy-Pasting Large Code Blocks

**Bad:**
```
Here's the entire 300-line file I need you to review:
[paste 300 lines]
```

**Good:**
```
Review service.go focusing on validation logic (lines 45-78).
```

### 🚫 Multiple Small Questions in Sequence

**Bad:**
```
Q1: Does file X exist?
[wait for response]
Q2: What's in it?
[wait for response]
Q3: Does it have method Y?
[wait for response]
```

**Good:**
```
Check file X: exists? contents? has method Y?
```

---

## Efficient Communication Patterns

### Pattern 1: Tool-First Approach

```
# Instead of asking, use tools directly
✅ view:internal/domain/career/event.go
✅ grep:"validation" internal/service/
✅ ls:internal/cli/
```

### Pattern 2: Specific Line References

```
# Don't paste whole files
✅ "See lines 45-67 in service.go for context"
❌ [Paste entire file]
```

### Pattern 3: Bullet Points

```
✅ Need to:
- Add validation method
- Update tests
- Document changes

❌ I need to add a validation method to the service layer.
After that, I'll need to update the tests to cover the new
validation. Then I should update the documentation...
```

### Pattern 4: Command Chains

```
✅ "Run: fmt, vet, test, coverage"
❌ "First run fmt. Then run vet. After that run test..."
```

### Pattern 5: Assume Context

```
# After discussing DDD multiple times
✅ "Add to domain layer"
❌ "Add to the domain layer, which as we know is responsible
for business logic and validation according to DDD principles..."
```

---

## When to Start Fresh

### Indicators You Should Start New Conversation

- ✅ Token count > 100k
- ✅ Switching to completely different task
- ✅ AI responses become slow/confused
- ✅ Repeating same information multiple times
- ✅ Conversation no longer focused

### How to Start Fresh Effectively

1. **Document current state** (if mid-task)
   ```bash
   # Save progress
   git stash -u -m "WIP: current task state"
   ```

2. **Summarize what's done** (1-2 sentences)
   ```
   Completed: Event validation and tests
   Next: Add filtering to repository layer
   ```

3. **Start new conversation with minimal context**
   ```
   Task: Add date range filtering to repository
   Context: Repository interface in internal/repository/career/
   Pattern: Follow existing tag filtering pattern
   ```

---

## Token Efficiency Checklist

Before each interaction, ask:

- [ ] Can I use a tool instead of asking?
- [ ] Am I being concise?
- [ ] Am I batching multiple operations?
- [ ] Am I referencing previous context instead of repeating?
- [ ] Am I focusing on deltas, not full state?
- [ ] Is my token count reasonable?

---

## Measuring Efficiency

### Good Session Example

```
Tokens used: 15,000
Tasks completed: 3
Lines of code: ~500
Tests written: 15
Time: 1 hour

Efficiency: ✅ Excellent
- Used tools effectively
- Concise communication
- Focused on specific tasks
```

### Poor Session Example

```
Tokens used: 95,000
Tasks completed: 2
Lines of code: ~300
Tests written: 8
Time: 2 hours

Efficiency: ❌ Poor
- Too much back-and-forth
- Repeated explanations
- Verbose responses
- Unfocused conversations
```

---

## Makefile Integration

Add to your `Makefile`:

```makefile
.PHONY: token-check

token-check:
	@echo "Token Efficiency Reminders:"
	@echo "  ✅ Use tools (view, grep, ls) over text"
	@echo "  ✅ Be concise and specific"
	@echo "  ✅ Batch multiple operations"
	@echo "  ✅ Reference context, don't repeat"
	@echo "  ✅ Focus on deltas, not full state"
	@echo ""
	@echo "Current best practices:"
	@echo "  - Token count < 50k: ✅ Good"
	@echo "  - Token count 50-100k: ⚠️ High (be more concise)"
	@echo "  - Token count >100k: 🔴 Critical (start fresh)"
```

---

## AI Prompt Optimization

### When Prompting AI

**❌ Inefficient Prompt:**
```
I have a service that handles career events. The service has methods
for creating, reading, updating, and deleting events. It uses a
repository for persistence. The repository has an interface that
defines methods. I want to add a new method for filtering events
by tags. The tags are strings in a slice. Can you help me implement
this? It should follow the existing patterns we have.
```

**✅ Efficient Prompt:**
```
Add FilterByTags(tags []string) method to CareerService.
Follow pattern from existing methods in internal/service/career/service.go.
```

### When AI Responds Verbosely

**Your response:**
```
✅ "Understood. Proceed with implementation."
❌ "Thank you for the detailed explanation. I understand
that we need to follow these patterns and..." [5 paragraphs]
```

---

## Summary

**Golden Rules:**
1. **Tools > Text**: Use view/grep/ls instead of asking
2. **Concise > Verbose**: Say less, communicate more
3. **Batch > Sequential**: Combine operations
4. **Reference > Repeat**: Point to context, don't re-explain
5. **Delta > Full State**: Focus on changes only

**Token Thresholds:**
- < 20k: ✅ Healthy
- 20-50k: ⚠️ Be concise
- 50-100k: 🔶 Consider fresh start
- \> 100k: 🔴 Start fresh NOW

**Check Regularly:**
```bash
make token-check
```

---

**Related Documentation:**
- [Rules Compliance Check](./rules-compliance-check.md)
- [Process Task List](./process-task-list.md)
- [Senior Engineer Guidelines](./senior-engineer-guidelines.md)

---

*Last Updated: 2025-12-23*
*Version: 1.0*


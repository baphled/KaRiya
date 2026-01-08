# Master Task Execution System Setup

## What Was Created

A comprehensive master prompt system that integrates ALL project rules into a single, actionable workflow for executing tasks.

---

## 📚 Documentation Created

### 1. Master Task Execution Prompt
**File:** `docs/rules/master-task-prompt.md` (1000+ lines)

**Complete workflow in 5 phases:**

1. **Preparation Phase**
   - Token awareness check
   - Full compliance check
   - Task review and breakdown
   - Pattern exploration

2. **TDD Phase (Red-Green-Refactor)**
   - RED: Write failing test first
   - GREEN: Implement minimal code
   - REFACTOR: Clean up (if needed)
   - Atomic commits at each step

3. **Compliance Verification Phase**
   - Code quality checks
   - Architecture compliance
   - Documentation updates

4. **Final Verification Phase**
   - Review all commits
   - Full compliance check
   - Token count verification

5. **Task Completion Phase**
   - Summary of work
   - Handoff notes

**Includes:**
- Detailed step-by-step instructions
- Troubleshooting guide
- Before/after checklists
- Example workflows
- Token efficiency integration
- All rules referenced

### 2. Task Quick Reference
**File:** `docs/rules/TASK_QUICK_REF.md** (One-page reference)

**Quick checklist version:**
- 5-phase workflow condensed
- Essential commands
- Commit template
- Token thresholds
- Common fixes

---

## 🎯 What This Integrates

### All Project Rules ✅

**Code Quality:**
- Go Guidelines (formatting, idioms, error handling)
- Senior Engineer Guidelines (SOLID, TDD)
- Testing Standards (Ginkgo/Gomega, 80% coverage)

**Commit Quality:**
- Atomic Commits (one logical change)
- Conventional Commits (type, scope, message format)
- Review process (no generated files, debug code)

**Task Processing:**
- Process Task List (one task at a time)
- Task Instructions (checklist, existing patterns)
- Tool-driven verification

**Token Efficiency:**
- Token awareness and thresholds
- Efficiency strategies
- When to start fresh

**Compliance:**
- Automated checking
- Violation detection
- Fix guidance

---

## 🚀 How to Use

### For Every Task

```bash
# 1. Before starting ANY task
make check-compliance

# 2. Open quick reference
cat docs/rules/TASK_QUICK_REF.md

# 3. Follow the 5-phase workflow:
#    - Prepare
#    - TDD (Red-Green-Refactor)
#    - Verify Compliance
#    - Final Check
#    - Complete

# 4. Before finishing
make check-compliance
```

### Daily Workflow

**Morning:**
```bash
make check-compliance    # Verify clean state
make token-check         # Review efficiency tips
```

**For each task:**
```bash
# Follow TASK_QUICK_REF.md
# Or full guide: master-task-prompt.md
```

**Before each commit:**
```bash
make review-commit
```

**End of session:**
```bash
make check-compliance
# Note token count for next session
```

---

## 📋 The 5-Phase Workflow

### Phase 1: Preparation (5 minutes)
- Check token count (< 50k?)
- Run compliance check
- Understand the single task
- Review existing patterns
- Verify task is atomic

### Phase 2: TDD (Main work)
**RED:** Write failing test → Commit
```bash
git commit -m "test(scope): add failing test for X"
```

**GREEN:** Implement → Commit
```bash
git commit -m "feat(scope): implement X"
```

**REFACTOR:** Clean up → Commit (if needed)
```bash
git commit -m "refactor(scope): improve X"
```

### Phase 3: Compliance Verification
- Format code (`make fmt`)
- Static analysis (`make vet`)
- Run tests (`make test`)
- Check races (`go test -race ./...`)
- Verify coverage (`make coverage`)

### Phase 4: Final Verification
- Review commits (`git log`)
- Full compliance (`make check-compliance`)
- Token count check

### Phase 5: Task Completion
- Verify all checks pass
- Note what was done
- Handoff for next task/session

---

## 💡 Key Features

### Integrated Compliance
Every phase includes compliance checks:
- Before starting (clean state)
- During work (incremental verification)
- Before finishing (final validation)

### Token Awareness
Built-in token monitoring:
- Check at start of each phase
- Guidance when approaching limits
- Clear thresholds (< 20k, 50k, 100k)

### TDD Enforcement
Strict Red-Green-Refactor:
- Tests ALWAYS written first
- Implementation only after failing test
- Refactoring is separate step
- Each step gets its own commit

### Atomic Commits
Natural atomic commits:
- One commit per TDD phase
- Clear separation of concerns
- Easy to review and revert
- Conventional format enforced

### Tool-Driven
Emphasizes tool usage:
- `view` for reading code
- `grep` for searching
- `ls` for exploration
- Automated checks via `make`

---

## 📖 Example: Complete Task Execution

### Task: Add Duration field to CareerEvent

```bash
# PHASE 1: PREPARATION
make check-compliance        # ✅ Pass
# Token: 15k ✅
ls internal/domain/career/   # Explore
view internal/domain/career/event.go

# PHASE 2: TDD

## RED
view internal/domain/career/event_test.go
# Add test for Duration field
git add internal/domain/career/event_test.go
make review-commit
git commit -m "test(domain): add test for Duration field validation

Duration should be non-negative integer.
Test validates field constraints."

## GREEN
view internal/domain/career/event.go
# Add Duration field and validation
make test                    # ✅ Pass
git add internal/domain/career/event.go
make review-commit
git commit -m "feat(domain): add Duration field to CareerEvent

Duration tracks hours spent on event.
Validated as non-negative integer."

# PHASE 3: COMPLIANCE VERIFICATION
make fmt                     # ✅
make vet                     # ✅
make test                    # ✅
go test -race ./...          # ✅
make coverage                # ✅ 82%

# PHASE 4: FINAL VERIFICATION
git log --oneline -2
# Verify both commits are atomic ✅
make check-compliance        # ✅ All pass
# Token: 22k ✅

# PHASE 5: COMPLETE
# Summary: Added Duration field with validation
# Commits: 2 atomic commits
# Coverage: 82% (maintained)
# Token: 22k (healthy)
# Ready for: Next task
```

---

## 🎓 Best Practices

### Always Start Clean
```bash
make check-compliance
```
Don't start new work with violations.

### One Task at a Time
Focus on single, atomic task. If it touches multiple layers, break it down.

### Test First, Always
Never write production code before the test. This is non-negotiable.

### Commit Frequently
After each TDD phase (Red, Green, Refactor). Small, atomic commits.

### Check Regularly
Run `make test` frequently. Don't accumulate broken code.

### Monitor Tokens
Check token count at each phase transition. Start fresh if > 100k.

### Use Tools
View files with `view`, search with `grep`, explore with `ls`.

### Be Concise
Especially when token count > 50k. Bullet points, no fluff.

---

## 🔧 Troubleshooting

### "Task is too large"
**Solution:** Break into subtasks
- One subtask per layer (domain, service, repo)
- Each subtask follows full 5-phase workflow
- Each subtask should take < 30 minutes

### "Don't know what test to write"
**Solution:** Find similar pattern
```bash
grep -r "similar_functionality" internal/
view path/to/similar_test.go
# Copy pattern, modify for your case
```

### "Token count too high"
**Solution:** Optimize communication
- Use ONLY tools (view/grep/ls)
- One-word confirmations
- No explanations
- Batch all operations
- If > 100k: STOP, start fresh

### "Compliance check fails"
**Solution:** Fix violations immediately
```bash
make check-compliance
# Read specific violation
# Apply suggested fix
# Re-run until clean
```

### "Can't make test pass"
**Solution:** Debug systematically
```bash
go test -v ./... -run "TestName"
# Add logging (remove before commit)
# Check intermediate values
# Verify test expectations
# Ask for help if stuck > 30 min
```

---

## 📊 Success Metrics

### Good Task Execution
- ✅ Completed in one session
- ✅ All compliance checks pass
- ✅ 2-5 atomic commits
- ✅ Tests all passing
- ✅ Coverage maintained/improved
- ✅ Token count < 50k
- ✅ Clear commit history

### Needs Improvement
- ❌ Multiple sessions needed
- ❌ Compliance violations
- ❌ Large commits
- ❌ Tests failing
- ❌ Coverage decreased
- ❌ Token count > 100k
- ❌ Unclear commits

---

## 🔗 Related Documentation

### Full Guides
- **Master Task Prompt:** `docs/rules/master-task-prompt.md`
- **Task Quick Reference:** `docs/rules/TASK_QUICK_REF.md`
- **Atomic Commits:** `docs/rules/atomic-commits.md`
- **Token Efficiency:** `docs/rules/token-efficiency.md`
- **Rules Compliance:** `docs/rules/rules-compliance-check.md`
- **Senior Engineer:** `docs/rules/senior-engineer-guidelines.md`
- **Go Guidelines:** `docs/rules/go-guidelines.md`
- **Task Processing:** Integrated into `docs/rules/master-task-prompt.md`

### Quick References
- **Task Quick Ref:** `docs/rules/TASK_QUICK_REF.md`
- **Compliance Quick Ref:** `docs/rules/COMPLIANCE_QUICK_REF.md`
- **Commit Quick Ref:** `docs/rules/COMMIT_QUICK_REFERENCE.md`

### Scripts
- `scripts/check-compliance.sh` - Full compliance check
- `scripts/review-commit.sh` - Commit review

### Makefile Targets
```bash
make check-compliance    # Full compliance
make review-commit       # Commit review
make token-check         # Token tips
make test                # Run tests
make coverage            # Coverage report
make fmt                 # Format code
make vet                 # Static analysis
make staticcheck         # Advanced static analysis
```

---

## 🎯 Quick Commands

```bash
# Check compliance
make check-compliance

# View quick reference
cat docs/rules/TASK_QUICK_REF.md

# View full guide
cat docs/rules/master-task-prompt.md

# Token efficiency tips
make token-check

# Before commit
make review-commit

# Code quality
make fmt && make vet && make test
```

---

## 📝 Summary

**What you have:**
- ✅ Comprehensive master prompt (1000+ lines)
- ✅ Quick reference (1-page)
- ✅ Integration of ALL project rules
- ✅ 5-phase workflow (Prepare, TDD, Verify, Check, Complete)
- ✅ Token efficiency built-in
- ✅ Compliance checks at each phase
- ✅ TDD enforcement (Red-Green-Refactor)
- ✅ Atomic commit guidance
- ✅ Troubleshooting guide
- ✅ Example workflows
- ✅ Success metrics

**How to use:**
1. Open `docs/rules/TASK_QUICK_REF.md`
2. Follow 5-phase workflow for EVERY task
3. Run `make check-compliance` before and after
4. Monitor token count throughout
5. Create atomic commits at each TDD step

**This ensures:**
- All rules followed automatically
- High code quality maintained
- Atomic commits created naturally
- Token usage stays efficient
- Tasks completed successfully

---

**Start your next task:**
```bash
make check-compliance
cat docs/rules/TASK_QUICK_REF.md
# Follow the workflow
```

---

*Created: 2025-12-23*
*Version: 1.0*


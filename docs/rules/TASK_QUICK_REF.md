# Task Execution Quick Reference

## Quick Start

```bash
make check-compliance  # Before starting
# Follow workflow below
make check-compliance  # Before finishing
```

---

## The Workflow (5 Phases)

### 1. Prepare (5 min)
- [ ] Check token count (< 50k?)
- [ ] `make check-compliance` passes
- [ ] Understand the ONE task
- [ ] Review existing patterns

### 2. TDD (Red-Green-Refactor)

**RED: Write Test**
```bash
view *_test.go          # See existing tests
# Write failing test
git add *_test.go
make review-commit
git commit -m "test(scope): add failing test for X"
```

**GREEN: Implement**
```bash
view *.go               # See existing code
# Write minimal implementation
make test               # Verify passes
git add *.go
make review-commit
git commit -m "feat(scope): implement X"
```

**REFACTOR: Clean Up (if needed)**
```bash
# Improve without changing behavior
make test               # Still passes
git add *.go
make review-commit
git commit -m "refactor(scope): improve X"
```

### 3. Verify Compliance
```bash
make fmt                # Format
make vet                # Static analysis
make test               # All tests
go test -race ./...     # Race check
make coverage           # Coverage check
```

### 4. Final Check
```bash
git log --oneline -5    # Review commits
make check-compliance   # Full check
# Token count: _____ (< 100k?)
```

### 5. Complete
- [ ] All tests pass
- [ ] Coverage ≥ 80%
- [ ] Atomic commits
- [ ] Compliance passes
- [ ] Token count OK

---

## Commit Template

```
<type>(<scope>): <subject>

<why this change was needed>

<issue references>
```

**Types:** feat, fix, docs, refactor, test, chore
**Scopes:** domain, service, repo, cli, logger

---

## Before Each Commit

```bash
make review-commit
```

Check:
- [ ] ONE logical change
- [ ] Clear message (type, scope, subject)
- [ ] Explains WHY
- [ ] No generated files
- [ ] No debug code
- [ ] Tests included

---

## Commands

```bash
make check-compliance    # Full check
make review-commit       # Commit check
make token-check         # Token tips
make fmt                 # Format
make vet                 # Analysis
make test                # Tests
make coverage            # Coverage
```

---

## Token Efficiency

**Thresholds:**
- < 20k: ✅ Healthy
- 20-50k: ⚠️ Be concise
- 50-100k: 🔶 Consider fresh start
- \> 100k: 🔴 Start fresh NOW

**Best Practices:**
- Use tools (view, grep, ls)
- Be concise (bullet points)
- Batch operations
- Reference context
- Focus on deltas

---

## Troubleshooting

**Task too large?** → Break into subtasks
**Test unclear?** → Find similar test pattern
**Token high?** → Use tools only, be concise
**Compliance fails?** → Follow fix suggestions

---

## Example

```bash
# Prep
make check-compliance

# RED
view internal/service/career/service_test.go
# Add test
git commit -m "test: add test for X"

# GREEN
view internal/service/career/service.go
# Implement
make test
git commit -m "feat: implement X"

# Verify
make check-compliance

# Done ✅
```

---

**Full Guide:** `docs/rules/master-task-prompt.md`


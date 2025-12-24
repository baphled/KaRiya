# Master Task Execution Prompt

## Purpose

This is the **definitive prompt** for executing tasks in the KaRiya project. It integrates ALL project rules, guidelines, and best practices into a single, step-by-step workflow.

Use this prompt for EVERY task to ensure:
- ✅ All project rules are followed
- ✅ Token efficiency is maintained
- ✅ Code quality standards are met
- ✅ Atomic commits are created
- ✅ TDD principles are applied
- ✅ Tasks are completed successfully

---

## Quick Start

```bash
# 1. Check compliance before starting
make check-compliance

# 2. Follow the workflow below

# 3. Check compliance before finishing
make check-compliance
```

---

## Master Workflow

### Phase 1: Preparation (5 minutes)

#### 1.1 Token Awareness
```bash
# Check token efficiency reminders
make token-check
```

**Current token count:** _____ (fill in from interface)
- < 20k: ✅ Healthy
- 20-50k: ⚠️ Be concise
- 50-100k: 🔶 Consider fresh start soon
- > 100k: 🔴 Start fresh NOW

**Efficiency mindset:**
- Use tools (view, grep, ls) over text
- Be concise and specific
- Batch multiple operations
- Reference context, don't repeat
- Focus on deltas, not full state

#### 1.2 Rules Check
```bash
# Run full compliance check
make check-compliance
```

**Fix any violations before proceeding.**

#### 1.3 Task Review

**Read the task carefully:**
- What is the SINGLE goal?
- What files will be affected?
- What tests are needed?
- What documentation needs updating?

**Check existing patterns:**
```bash
# Explore relevant areas
ls internal/domain/
ls internal/service/
ls internal/repository/

# Search for similar functionality
grep -r "similar_function" internal/
```

**Decision: Is this ONE atomic task?**
- ✅ Yes → Proceed
- ❌ No → Break down into smaller tasks

---

### Phase 2: Test-Driven Development (Red-Green-Refactor)

#### 2.1 RED: Write Failing Test

**Before writing ANY production code, write the test.**

```bash
# View existing test file
view internal/service/career/service_test.go

# Or create new test file
# Follow existing patterns
```

**Test Requirements:**
- Use Ginkgo/Gomega
- One expectation per `It` block
- Descriptive test names
- Proper `Describe` and `Context` structure
- Test should FAIL initially

**Example:**
```go
var _ = Describe("ServiceName", func() {
    Context("when adding new feature", func() {
        It("should return expected result", func() {
            result := service.NewFeature()
            Expect(result).To(Equal(expectedValue))
        })
    })
})
```

**Commit the failing test:**
```bash
git add *_test.go
make review-commit
git commit -m "test(scope): add failing test for new feature

Describes expected behavior for [feature name].
Test currently fails as implementation doesn't exist yet."
```

**Token check:** Still < 50k? ✅

#### 2.2 GREEN: Implement Minimal Code

**Now implement the minimal code to make the test pass.**

```bash
# View relevant files
view internal/service/career/service.go

# Make minimal changes
```

**Implementation Requirements:**
- Minimal code to pass test (no more)
- Follow Go idioms and guidelines
- Proper error handling
- Follow existing patterns
- Add comments for complex logic

**Verify:**
```bash
# Run the specific test
make individual-test TEST="your test name"

# Run all tests
make test

# Check for race conditions
go test -race ./...

# Check off the corresponding task in your checklist
sed -i '' '/- \[ \] Implement the feature/a\- [x] Implement the feature' your-checklist-file.md
```

**Commit the implementation:**
```bash
git add *.go
make review-commit
git commit -m "feat(scope): implement new feature

[Brief description of what was implemented and why]

Implementation follows existing patterns in [relevant file].
All tests now pass."
```

**Token check:** Still < 50k? ✅

#### 2.3 REFACTOR: Clean Up Code (if needed)

**If code has duplication or can be improved, refactor now.**

**Refactoring Requirements:**
- No behavior changes (tests still pass)
- Extract methods for clarity
- Remove duplication (DRY)
- Improve naming
- Simplify complex logic

**Verify after refactoring:**
```bash
make test
go test -race ./...
```

**Commit refactoring separately:**
```bash
git add *.go
make review-commit
git commit -m "refactor(scope): extract method for clarity

[Description of what was refactored and why]

No behavior changes - all tests still pass."
```

**Token check:** Still < 50k? ✅

---

### Phase 3: Compliance Verification

#### 3.1 Code Quality Check

```bash
# Format code
make fmt

# Run static analysis
make vet

# Run all tests
make test

# Check race conditions
go test -race ./...

# Check coverage
make coverage
```

**Requirements:**
- ✅ All tests pass
- ✅ No race conditions
- ✅ Coverage ≥ 80%
- ✅ No vet warnings
- ✅ Code formatted

#### 3.2 Architecture Check

**Verify layer boundaries:**
- Domain layer has NO dependencies on service/repository
- Service layer depends on domain + repository
- Repository layer depends only on domain
- CLI layer depends on service

**Check imports:**
```bash
# Domain should have no service/repo imports
grep -r "service\|repository" internal/domain/

# Should return empty
```

#### 3.3 Documentation Check

**Update documentation if needed:**
- [ ] Add/update function comments
- [ ] Update README if user-facing changes
- [ ] Update AGENTS.md if patterns changed
- [ ] Update feature docs if applicable

**Commit documentation separately:**
```bash
git add *.md
make review-commit
git commit -m "docs: update documentation for new feature

[What was documented and why]"
```

**Token check:** Still < 50k? ✅

---

### Phase 4: Final Verification

#### 4.1 Review All Commits

```bash
# View commits made
git log --oneline -10

# Review each commit
git show <commit-hash>
```

**Verify each commit:**
- [ ] Represents ONE logical change
- [ ] Has clear commit message (type, scope, subject)
- [ ] Message explains WHY, not just WHAT
- [ ] No generated files
- [ ] No debug code
- [ ] Tests included

**If any commit is not atomic:**
```bash
# Use interactive rebase to fix
git rebase -i HEAD~N

# Split or reorder commits as needed
```

#### 4.2 Full Compliance Check

```bash
# Run comprehensive compliance check
make check-compliance
```

**ALL checks must pass before finishing:**
- ✅ Code quality (fmt, vet, build, test, race)
- ✅ Test coverage (≥ 80%)
- ✅ No violations in staged area
- ✅ Architecture compliance
- ✅ Documentation complete
- ✅ Testing standards met
- ✅ Dependencies healthy
- ✅ File organization correct
- ✅ Git health good

**If any checks fail:**
- Fix the issues
- Re-run compliance check
- Update affected commits if needed

#### 4.3 Final Token Check

**Current token count:** _____ (fill in from interface)

**Actions:**
- < 50k: ✅ Good, can continue with more tasks
- 50-100k: ⚠️ High, consider finishing here
- > 100k: 🔴 Must stop, start fresh next time

---

### Phase 5: Task Completion

#### 5.1 Summary

**What was accomplished:**
- [ ] Test written and passing
- [ ] Implementation complete
- [ ] Code refactored (if needed)
- [ ] Documentation updated
- [ ] All tests passing
- [ ] Coverage maintained/improved
- [ ] Atomic commits created
- [ ] Compliance checks passed

#### 5.2 Handoff Notes

**If continuing with next task, note:**
```
Completed: [Brief task description]
Commits: [List commit hashes]
Coverage: [Current coverage %]
Token count: [Current count]
Ready for: [Next task]
```

**If stopping here, note:**
```
Completed: [Brief task description]
Status: Ready to push
Branch: [Branch name]
Next session: [What to work on next]
```

---

## Troubleshooting

### Issue: Task Too Large

**Symptom:** Task affects multiple layers or features

**Solution:**
```
1. Break task into smaller subtasks:
   - Subtask 1: Domain layer changes
   - Subtask 2: Service layer changes
   - Subtask 3: Repository layer changes
   - Subtask 4: CLI layer changes

2. Execute each subtask separately using this workflow

3. Each subtask should take < 30 minutes
```

### Issue: Test Unclear

**Symptom:** Don't know what test to write

**Solution:**
```
1. Review task requirement carefully
2. Find similar functionality:
   grep -r "similar_feature" internal/
3. View existing test as template:
   view internal/service/career/service_test.go
4. Write test for expected behavior (can be simple)
5. Implement to make it pass
```

### Issue: Token Count High

**Symptom:** Token count > 50k

**Solution:**
```
1. Check if task is complete
   - If yes: Stop here, start fresh next session
   - If no: Continue but be VERY concise

2. Apply token efficiency:
   - Use tools exclusively
   - One-line responses
   - No explanations unless asked
   - Batch all operations

3. If > 100k: MUST stop and start fresh
```

### Issue: Compliance Check Fails

**Symptom:** `make check-compliance` returns violations

**Solution:**
```
1. Read the specific violation
2. Apply suggested fix
3. Re-run compliance check
4. Repeat until all checks pass

Common fixes:
- go fmt ./...          # Formatting
- go vet ./...          # Vet warnings
- make test             # Fix failing tests
- git reset <file>      # Remove generated files
```

### Issue: Can't Make Test Pass

**Symptom:** Implementation not making test pass

**Solution:**
```
1. Run test with verbose output:
   go test -v ./... -run "TestName"

2. Check test expectations:
   - Is expected value correct?
   - Is setup/teardown working?
   - Are mocks configured correctly?

3. Debug implementation:
   - Add logging (remove before commit)
   - Check intermediate values
   - Verify assumptions

4. Ask for help if stuck > 30 minutes
```

---

## Checklists

### Start of Task Checklist

Before writing any code:

- [ ] Token count checked (< 50k?)
- [ ] `make check-compliance` passed
- [ ] Task is clearly understood
- [ ] Task is atomic (ONE thing)
- [ ] Existing patterns reviewed
- [ ] Test plan identified

### Before Each Commit Checklist

Before running `git commit`:

- [ ] `make review-commit` passed
- [ ] One logical change only
- [ ] Commit message written (type, scope, subject)
- [ ] Message explains WHY
- [ ] No generated files
- [ ] No debug code
- [ ] Tests included (if applicable)

### End of Task Checklist

Before marking task complete:

- [ ] All tests pass (`make test`)
- [ ] No race conditions (`go test -race ./...`)
- [ ] Coverage ≥ 80% (`make coverage`)
- [ ] Code formatted (`make fmt`)
- [ ] No vet warnings (`make vet`)
- [ ] `make check-compliance` passed
- [ ] All commits are atomic
- [ ] Documentation updated
- [ ] Token count reasonable (< 100k)

---

## Quick Reference Commands

```bash
# Compliance and checks
make check-compliance    # Full rules compliance check
make review-commit       # Commit-specific review
make token-check         # Token efficiency reminder

# Code quality
make fmt                 # Format code
make vet                 # Static analysis
make test                # Run all tests
make coverage            # Generate coverage

# Development
make build               # Build application
make pre-commit          # Quick pre-commit checks

# Git workflow
git add <files>          # Stage specific files
git add -p <file>        # Interactive staging
git commit               # Commit with message
git log --oneline -10    # Recent commits
git show <hash>          # Review specific commit
```

---

## Rules Reference

This workflow implements ALL project rules:

### Code Quality
- **Go Guidelines** (`docs/rules/go-guidelines.md`)
  - Formatting with gofmt
  - Idiomatic Go patterns
  - Proper error handling
  - Interface-based design
  - Structured logging

- **Senior Engineer Guidelines** (`docs/rules/senior-engineer-guidelines.md`)
  - SOLID principles
  - Red-Green-Refactor (TDD)
  - Clean code practices
  - Architectural decisions
  - Observability

### Commit Quality
- **Atomic Commits** (`docs/rules/atomic-commits.md`)
  - One logical change per commit
  - Self-contained commits
  - Clear commit messages
  - No generated files

- **Conventional Commits**
  - Type: feat, fix, docs, refactor, test, chore
  - Scope: domain, service, repo, cli, logger
  - Subject: Clear, concise (< 50 chars)
  - Body: Explains WHY
  - Footer: Issue references

### Task Processing
- **Process Task List** (`docs/rules/process-task-list.md`)
  - One task at a time
  - Tool-driven verification
  - Authority order (tools > checklist > guidelines)
  - Deterministic execution

- **Task Instructions** (`docs/rules/task-instructions.md`)
  - Strict checklist following
  - Reference codebase first
  - Follow existing patterns
  - One expectation per test

### Token Efficiency
- **Token Efficiency** (`docs/rules/token-efficiency.md`)
  - Use tools over text
  - Be concise and specific
  - Batch operations
  - Reference context
  - Focus on deltas

---

## Examples

### Example 1: Adding New Domain Field

**Task:** Add `Duration` field to CareerEvent

**Workflow:**

```bash
# Phase 1: Preparation
make check-compliance
# Token count: 15k ✅

# Phase 2: TDD

# RED: Write failing test
view internal/domain/career/event_test.go
# Add test for Duration field validation

git add internal/domain/career/event_test.go
make review-commit
git commit -m "test(domain): add test for Duration field

Duration should be a positive integer representing hours.
Test validates field presence and constraints."

# GREEN: Implement
view internal/domain/career/event.go
# Add Duration field to struct
# Add validation for Duration

make test
# Tests pass ✅

git add internal/domain/career/event.go
make review-commit
git commit -m "feat(domain): add Duration field to CareerEvent

Adds Duration field to track time spent on events.
Duration is validated to be non-negative integer."

# Phase 3: Compliance
make check-compliance
# All pass ✅

# Phase 4: Final verification
git log --oneline -2
# Verify atomic commits ✅

# Token count: 18k ✅

# Phase 5: Complete
```

### Example 2: Fixing a Bug

**Task:** Fix nil pointer in service layer

**Workflow:**

```bash
# Phase 1: Preparation
make check-compliance

# Phase 2: TDD

# RED: Write failing test that reproduces bug
view internal/service/career/service_test.go
# Add test with nil event

git add internal/service/career/service_test.go
make review-commit
git commit -m "test(service): add test for nil event handling

Test verifies that passing nil event returns appropriate error.
Currently fails - reproduces reported bug."

# GREEN: Fix the bug
view internal/service/career/service.go
# Add nil check

make test
# Tests pass ✅

git add internal/service/career/service.go
make review-commit
git commit -m "fix(service): prevent nil pointer in CaptureEvent

Add nil check before accessing event properties.
Returns descriptive error when nil event provided.

Fixes #87"

# Phase 3: Compliance
make check-compliance

# Phase 4: Final verification
git log --oneline -2

# Phase 5: Complete
```

---

## Token Efficiency in Practice

### Good Session Example

```
Token count: 15k → 22k → 28k → 35k
Tasks completed: 3
Time: 1.5 hours
Quality: All checks passed

Why efficient:
- Used view/grep/ls for exploration
- Concise responses
- Batched operations
- Followed workflow strictly
```

### Poor Session Example

```
Token count: 20k → 45k → 78k → 105k
Tasks completed: 2
Time: 2 hours
Quality: Some checks failed

Why inefficient:
- Asked for code instead of using tools
- Verbose explanations
- Many back-and-forth exchanges
- Didn't follow workflow
```

---

## Summary

**This workflow ensures:**
- ✅ TDD is followed (Red-Green-Refactor)
- ✅ Atomic commits are created
- ✅ All code quality rules met
- ✅ Architecture boundaries respected
- ✅ Tests comprehensive and passing
- ✅ Documentation up to date
- ✅ Token usage efficient
- ✅ Tasks completed successfully

**Use for EVERY task:**
1. Check compliance (`make check-compliance`)
2. Follow TDD (Red → Green → Refactor)
3. Verify compliance (`make check-compliance`)
4. Create atomic commits (`make review-commit`)
5. Monitor token usage (`make token-check`)

**Success metrics:**
- All compliance checks pass
- All tests pass with ≥ 80% coverage
- All commits are atomic
- Token count stays reasonable
- Task is complete and working

---

**Related Documentation:**
- [Atomic Commits](./atomic-commits.md)
- [Token Efficiency](./token-efficiency.md)
- [Rules Compliance Check](./rules-compliance-check.md)
- [Senior Engineer Guidelines](./senior-engineer-guidelines.md)
- [Go Guidelines](./go-guidelines.md)
- [Process Task List](./process-task-list.md)

---

*Last Updated: 2025-12-23*
*Version: 1.0*


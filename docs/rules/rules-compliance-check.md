# Rules Compliance Check

## Purpose

This prompt helps ensure ALL project rules are being followed during development work. Use this regularly (especially before major commits or when starting new tasks) to verify compliance across all documented guidelines.

---

## Quick Compliance Check

Run through this checklist periodically during work:

### ✅ Code Quality Rules

**Go Guidelines** (`docs/rules/go-guidelines.md`)
- [ ] Code is formatted with `gofmt`
- [ ] Following idiomatic Go patterns
- [ ] Using proper error handling (explicit checks, wrapping with context)
- [ ] Interfaces used for testability
- [ ] Context passed for cancellation/timeouts
- [ ] Structured logging (not print statements)
- [ ] Code organized in proper package structure

**Senior Engineer Guidelines** (`docs/rules/senior-engineer-guidelines.md`)
- [ ] Following SOLID principles
- [ ] Applying Red-Green-Refactor (TDD)
- [ ] Writing tests BEFORE implementation
- [ ] Code is properly documented
- [ ] Error handling is fail-fast with descriptive messages
- [ ] Observability included (logging, tracing where applicable)

**Testing Standards**
- [ ] Tests pass: `go test ./...`
- [ ] No race conditions: `go test -race ./...`
- [ ] Coverage meets threshold (80%+)
- [ ] Using Ginkgo/Gomega patterns correctly
- [ ] One expectation per `It` block
- [ ] Test names are descriptive

### ✅ Commit Rules

**Atomic Commits** (`docs/rules/atomic-commits.md`)
- [ ] Each commit represents ONE logical change
- [ ] Commit message follows conventional format: `<type>(<scope>): <subject>`
- [ ] Commit message explains WHY, not just WHAT
- [ ] No generated files in staging (coverage.out, *.exe, etc.)
- [ ] No debug statements left in code
- [ ] No secrets or credentials
- [ ] Related tests included with code changes
- [ ] Ran `make review-commit` before committing

### ✅ Task Processing Rules

**Process Task List** (`docs/rules/process-task-list.md`)
- [ ] Working on exactly ONE task at a time
- [ ] Task checklist is locked and followed sequentially
- [ ] Using tools (not assumptions) to verify state
- [ ] Marking tasks complete after verification
- [ ] Not skipping ahead or batch-processing tasks

**Task Processing Guidelines**
- [ ] Strictly following task checklist order (see process-task-list.md)
- [ ] Referencing codebase before making changes
- [ ] Following existing patterns in codebase
- [ ] Avoiding side effects and scope creep
- [ ] One expectation per test case

### ✅ Documentation Rules

- [ ] README updated if user-facing changes
- [ ] API documentation updated if interfaces changed
- [ ] Comments explain WHY, not WHAT
- [ ] Complex logic has explanatory comments
- [ ] ADRs created for architectural decisions

---

## Detailed Compliance Review

### 1. Code Quality Check

```bash
# Run these commands to verify code quality
go fmt ./...                    # Format code
go vet ./...                   # Static analysis
go test ./...                  # All tests pass
go test -race ./...            # No race conditions
make coverage                  # Check coverage
```

**Review Checklist:**
- [ ] All commands pass without errors
- [ ] Coverage is 80% or higher
- [ ] No warnings from go vet
- [ ] No race conditions detected

### 2. Commit Compliance Check

```bash
# Review staged changes
git diff --cached --stat
git diff --cached --name-only

# Run automated review
make review-commit
```

**Review Checklist:**
- [ ] Fewer than 10 files changed (unless initial setup)
- [ ] Fewer than 500 lines changed (unless initial setup)
- [ ] All changed files relate to ONE logical change
- [ ] No generated files (*.out, *.test, coverage.*)
- [ ] No debug code (TODO, FIXME, console.log, fmt.Println("DEBUG"))
- [ ] Commit message draft follows format

**Commit Message Validation:**
```
<type>(<scope>): <subject line - max 50 chars>

<body - explain WHY, wrap at 72 chars>

<footer - issue references>
```

Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`, `perf`
Scopes: `domain`, `service`, `repo`, `cli`, `logger`

### 3. Test Compliance Check

**Ginkgo/Gomega Best Practices:**
- [ ] Test files have `_test.go` suffix
- [ ] Using `Describe` for grouping related tests
- [ ] Using `Context` for different scenarios
- [ ] Using `It` for individual test cases
- [ ] Using `BeforeEach`/`AfterEach` for setup/teardown
- [ ] Gomega matchers used correctly (Expect, To, Equal, HaveOccurred, etc.)
- [ ] Test names are descriptive and clear
- [ ] One logical assertion per `It` block

**Example Pattern:**
```go
var _ = Describe("ServiceName", func() {
    var (
        service *Service
        repo    *Repository
    )

    BeforeEach(func() {
        repo = NewRepository()
        service = NewService(repo)
    })

    Context("when condition X", func() {
        It("should do Y", func() {
            result := service.DoSomething()
            Expect(result).To(Equal(expectedValue))
        })
    })
})
```

### 4. Architecture Compliance Check

**Domain-Driven Design (DDD) Principles:**
- [ ] Domain layer has no dependencies on other layers
- [ ] Domain models enforce business rules through validation
- [ ] Service layer orchestrates domain objects
- [ ] Repository layer handles only persistence
- [ ] Clear separation of concerns maintained

**Layer Responsibilities:**
```
Domain Layer (internal/domain/):
  - Pure domain logic
  - Business rule validation
  - Domain models (structs, enums)
  - No external dependencies

Service Layer (internal/service/):
  - Business logic orchestration
  - Transaction management
  - Logging and observability
  - Depends on: Domain, Repository

Repository Layer (internal/repository/):
  - Data persistence
  - Query building
  - Database operations
  - Depends on: Domain

CLI Layer (internal/cli/):
  - User interface
  - Input parsing
  - Output formatting
  - Depends on: Service
```

### 5. Dependency Check

```bash
# Check for outdated dependencies
go list -u -m all

# Check for security vulnerabilities (if using)
go mod tidy
```

**Review:**
- [ ] No unused dependencies
- [ ] Dependencies are pinned to specific versions
- [ ] No known security vulnerabilities
- [ ] go.mod and go.sum are in sync

### 6. File Organization Check

**Project Structure Compliance:**
```
✅ Correct:
internal/domain/career/event.go          # Domain model
internal/domain/career/event_test.go     # Co-located test
internal/service/career/service.go       # Service layer
internal/repository/career/repository.go # Repository interface

❌ Incorrect:
src/event.go                             # Wrong location
domain/event.go                          # Missing internal/
internal/domain/event_test.go            # Wrong package
```

**Naming Conventions:**
- [ ] Package names are lowercase, no underscores
- [ ] File names are lowercase with underscores
- [ ] Test files end with `_test.go`
- [ ] Exported names start with capital letter
- [ ] Unexported names start with lowercase

### 7. Documentation Compliance

**Required Documentation:**
- [ ] Public functions have doc comments
- [ ] Doc comments start with function/type name
- [ ] Complex algorithms explained
- [ ] Non-obvious business logic documented
- [ ] README reflects current state
- [ ] AGENTS.md updated with new patterns (if applicable)

**Doc Comment Format:**
```go
// MethodName does X and returns Y.
// It handles Z edge case by doing A.
func MethodName() {}
```

---

## Token Efficiency Check

### Current Conversation Health

**Quick Indicators:**
- Token count > 50,000? → Consider starting fresh conversation
- Repeated information? → Remove redundancy
- Large code blocks? → Use tools to read files instead
- Verbose explanations? → Be more concise

### Token Conservation Strategies

1. **Use Tools Over Text**
   - ✅ Use `view` tool to read files
   - ✅ Use `grep` to search code
   - ✅ Use `ls` to explore structure
   - ❌ Paste entire file contents
   - ❌ Repeat information already provided
   - ❌ Explain obvious code

2. **Be Concise**
   - State intent clearly and briefly
   - Skip pleasantries in follow-ups
   - Don't repeat context unnecessarily
   - Focus on what's needed NOW

3. **Batch Operations**
   - Group related file operations
   - Make multiple tool calls in one message
   - Avoid back-and-forth for simple checks

4. **Context Management**
   - Reference previous messages instead of repeating
   - Assume understanding of established patterns
   - Focus on deltas (what changed) not full state

---

## Compliance Violations - Common Fixes

### Violation: Large Commit

**Symptom:** More than 10 files or 500 lines changed

**Fix:**
```bash
git reset
# Commit in logical groups (see atomic-commits.md)
```

### Violation: Multiple Changes in One Commit

**Symptom:** Commit touches multiple layers or features

**Fix:**
```bash
git reset
git add <domain-files>
git commit -m "feat(domain): add domain change"
git add <service-files>
git commit -m "feat(service): add service change"
```

### Violation: No Tests with Code

**Symptom:** Production code changed but no test changes

**Fix:**
```bash
# Don't commit yet! Write tests first
vim *_test.go
# Add tests, then commit together
```

### Violation: Generated Files Staged

**Symptom:** coverage.out, *.exe, etc. in staging

**Fix:**
```bash
git reset coverage.out *.exe
# Update .gitignore if needed
echo "coverage.out" >> .gitignore
```

### Violation: Poor Test Coverage

**Symptom:** Coverage below 80%

**Fix:**
```bash
# Identify gaps
make coverage
# Open coverage/index.html
# Add tests for uncovered code
```

### Violation: Race Conditions

**Symptom:** `go test -race` fails

**Fix:**
```bash
# Identify race
go test -race ./...
# Add proper synchronization (mutex, channels)
# Re-test until clean
```

### Violation: Not Following TDD

**Symptom:** Implementing code before tests

**Fix:**
```bash
# Back up if already committed
git reset --soft HEAD~1
# Recommit with proper order:
# 1. Write failing test
git add *_test.go
git commit -m "test: add failing test for X"
# 2. Implement
git add *.go
git commit -m "feat: implement X"
```

---

## Automated Compliance Check Script

Create `scripts/check-compliance.sh`:

```bash
#!/bin/bash

echo "================================================"
echo "🔍 RULES COMPLIANCE CHECK"
echo "================================================"
echo ""

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

VIOLATIONS=0

# 1. Code Quality
echo "📋 CODE QUALITY"
echo "------------------------------------------------"

echo -n "Formatting: "
if [ -z "$(gofmt -l . | grep -v vendor)" ]; then
    echo -e "${GREEN}✅ Pass${NC}"
else
    echo -e "${RED}❌ Fail - Run: go fmt ./...${NC}"
    VIOLATIONS=$((VIOLATIONS+1))
fi

echo -n "Build: "
if go build ./... 2>/dev/null; then
    echo -e "${GREEN}✅ Pass${NC}"
else
    echo -e "${RED}❌ Fail${NC}"
    VIOLATIONS=$((VIOLATIONS+1))
fi

echo -n "Tests: "
if go test ./... 2>/dev/null; then
    echo -e "${GREEN}✅ Pass${NC}"
else
    echo -e "${RED}❌ Fail${NC}"
    VIOLATIONS=$((VIOLATIONS+1))
fi

echo -n "Race Detection: "
if go test -race ./... 2>/dev/null; then
    echo -e "${GREEN}✅ Pass${NC}"
else
    echo -e "${RED}❌ Fail${NC}"
    VIOLATIONS=$((VIOLATIONS+1))
fi

echo ""

# 2. Test Coverage
echo "📊 TEST COVERAGE"
echo "------------------------------------------------"
COVERAGE=$(go test -cover ./... 2>/dev/null | grep -oP 'coverage: \K[0-9.]+' | awk '{sum+=$1; count++} END {print sum/count}')
if (( $(echo "$COVERAGE >= 80" | bc -l) )); then
    echo -e "Coverage: ${GREEN}${COVERAGE}% ✅${NC}"
else
    echo -e "Coverage: ${RED}${COVERAGE}% ❌ (Minimum: 80%)${NC}"
    VIOLATIONS=$((VIOLATIONS+1))
fi

echo ""

# 3. Staged Changes (if any)
echo "📝 STAGED CHANGES"
echo "------------------------------------------------"
if git diff --cached --quiet; then
    echo -e "${GREEN}No staged changes${NC}"
else
    FILE_COUNT=$(git diff --cached --name-only | wc -l)
    echo "Files staged: $FILE_COUNT"

    if [ "$FILE_COUNT" -gt 10 ]; then
        echo -e "${YELLOW}⚠️  Warning: >10 files (consider splitting)${NC}"
    fi

    if git diff --cached --name-only | grep -E '\.(out|exe|test)$'; then
        echo -e "${RED}❌ Generated files in staging!${NC}"
        VIOLATIONS=$((VIOLATIONS+1))
    fi
fi

echo ""

# 4. Documentation
echo "📚 DOCUMENTATION"
echo "------------------------------------------------"
echo -n "README exists: "
if [ -f "README.md" ]; then
    echo -e "${GREEN}✅${NC}"
else
    echo -e "${RED}❌${NC}"
    VIOLATIONS=$((VIOLATIONS+1))
fi

echo -n "AGENTS.md exists: "
if [ -f "AGENTS.md" ]; then
    echo -e "${GREEN}✅${NC}"
else
    echo -e "${YELLOW}⚠️  Consider creating${NC}"
fi

echo ""

# Summary
echo "================================================"
if [ $VIOLATIONS -eq 0 ]; then
    echo -e "${GREEN}✅ ALL CHECKS PASSED${NC}"
    exit 0
else
    echo -e "${RED}❌ $VIOLATIONS VIOLATION(S) FOUND${NC}"
    exit 1
fi
```

Make executable:
```bash
chmod +x scripts/check-compliance.sh
```

Add to Makefile:
```makefile
.PHONY: check-compliance

check-compliance:
	@bash scripts/check-compliance.sh
```

---

## Integration into Workflow

### Daily Workflow

```bash
# Morning: Check overall compliance
make check-compliance

# Before each commit: Check commit compliance
make review-commit

# After significant work: Re-check compliance
make check-compliance
```

### Pre-PR Checklist

Before creating a pull request:

```bash
# 1. Full compliance check
make check-compliance

# 2. Review all commits
git log origin/main..HEAD --oneline

# 3. Verify each commit is atomic
git log origin/main..HEAD --stat

# 4. Check coverage
make coverage

# 5. Final test run
make test
```

---

## Troubleshooting

### "Too many violations, where do I start?"

Priority order:
1. Fix failing tests (blocks everything else)
2. Fix race conditions (critical safety issue)
3. Fix formatting (easy quick win)
4. Add missing tests (improves stability)
5. Fix commit issues (improves maintainability)
6. Update documentation (helps team)

### "Compliance check fails but I think it's wrong"

1. Verify the check is correct
2. If check is outdated, update the check script
3. Document exceptions in project-specific rules
4. Don't ignore real violations

### "Takes too long to check compliance"

Run targeted checks:
```bash
# Just code quality
go fmt ./... && go vet ./...

# Just tests
make test

# Just commit check
make review-commit

# Full check only before PR
make check-compliance
```

---

## Summary

**Use this compliance check:**
- ✅ At start of work session
- ✅ Before each commit (`make review-commit`)
- ✅ After significant changes
- ✅ Before creating PR (`make check-compliance`)
- ✅ When token count is high (refocus)

**Key Principles:**
- Follow ALL documented rules
- Use tools to verify compliance
- Fix violations immediately
- Keep token usage efficient
- Maintain quality standards

---

## Quick Reference Commands

```bash
# Code quality
make fmt                    # Format code
make vet                    # Static analysis
make test                   # Run tests

# Compliance
make check-compliance       # Full check
make review-commit          # Commit check

# Coverage
make coverage               # Generate coverage
open coverage/index.html    # View coverage

# Pre-commit
make pre-commit            # Quick checks
```

---

**Related Documentation:**
- [Atomic Commits](./atomic-commits.md)
- [Senior Engineer Guidelines](./senior-engineer-guidelines.md)
- [Go Guidelines](./go-guidelines.md)
- [Process Task List](./process-task-list.md)
- [Token Efficiency Guide](./token-efficiency.md)

---

*Last Updated: 2025-12-23*
*Version: 1.0*


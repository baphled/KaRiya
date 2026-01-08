# Task 35: Test Coverage Improvement (62.5% → 80%)

**Status**: Ready for Implementation  
**Priority**: HIGH  
**Estimated Time**: 2-3 days  
**Current Coverage**: 62.5%  
**Target Coverage**: 80%  
**Gap**: 17.5%

---

## Overview

Increase test coverage from 62.5% to 80% to meet compliance requirements. This task focuses on identifying and testing uncovered code paths, particularly in CLI/TUI components, edge cases, and error handling paths.

---

## Current Situation

### Compliance Check Results
```
Test Coverage: 62.4926% ❌ (Target: 80%)
Status: ❌ Fail - Coverage significantly below 80%
```

### Known Coverage Stats (from AGENTS.md)
- **Overall**: 62.5% (needs improvement)
- **Intent Framework**: 88.1% ✅
- **GlobalContext**: 100% ✅
- **Domain Models**: >95% ✅
- **Repository**: >90% ✅
- **Service**: >85% ✅
- **CV Service**: 100% ✅

### Gap Analysis
The coverage gap is likely in:
1. **CLI/TUI components** (intents, models, app)
2. **Error handling paths** (untested error cases)
3. **Edge cases** (boundary conditions, nil checks)
4. **Integration points** (message passing, state transitions)

---

## Phase 1: Identify Coverage Gaps

### Acceptance Criteria
- [ ] Generate detailed coverage report with file-level breakdown
- [ ] Identify files with <60% coverage
- [ ] Identify specific uncovered functions/lines
- [ ] Create prioritized list of files to improve

### Steps

1. **Generate detailed coverage report**
   ```bash
   go test -coverprofile=coverage.out ./...
   go tool cover -html=coverage.out -o coverage.html
   go tool cover -func=coverage.out | sort -k3 -n > coverage-by-file.txt
   ```

2. **Analyze by package**
   ```bash
   go test -cover ./internal/cli/... -coverprofile=cli-coverage.out
   go test -cover ./internal/domain/... -coverprofile=domain-coverage.out
   go test -cover ./internal/repository/... -coverprofile=repo-coverage.out
   go test -cover ./internal/service/... -coverprofile=service-coverage.out
   ```

3. **Identify lowest coverage files**
   ```bash
   # Files with <60% coverage
   go tool cover -func=coverage.out | awk '$3 < 60 {print}'
   ```

4. **Document findings** in this task file

---

## Phase 2: Prioritize Test Improvements

### Categorize Gaps

#### High Priority (Business Logic)
Files in:
- `internal/cli/intents/` - Intent state machines
- `internal/cli/models/` - TUI models
- `internal/service/career/` - Business logic

#### Medium Priority (Infrastructure)
Files in:
- `internal/cli/app/` - Application root
- `internal/cli/components/` - Reusable components
- `internal/cli/context/` - Global context

#### Low Priority (Auxiliary)
Files in:
- `cmd/` - CLI entry points (hard to test)
- Helper utilities (may be covered indirectly)

### Create Test Plan
For each low-coverage file:
- [ ] Identify untested functions
- [ ] Determine test strategy (unit vs integration)
- [ ] Estimate time to test
- [ ] Assign to appropriate test suite

---

## Phase 3: Write Missing Tests

### Testing Strategy

#### Unit Tests (Ginkgo/Gomega)
For functions with:
- Pure logic (no side effects)
- Deterministic behavior
- Clear inputs/outputs

Example structure:
```go
var _ = Describe("FunctionName", func() {
    Context("when condition A", func() {
        It("should produce result X", func() {
            // Arrange
            // Act
            // Assert
        })
    })
    
    Context("when condition B", func() {
        It("should handle error Y", func() {
            // Test error path
        })
    })
})
```

#### Integration Tests
For components with:
- State machines (intents)
- Message passing
- Multi-component interactions

Use test harnesses from `internal/cli/intents/testing.go`

#### Edge Case Tests
Focus on:
- Nil checks
- Empty inputs
- Boundary conditions
- Error paths
- Concurrent access

### Implementation Checklist

#### CLI/App Package
- [ ] `internal/cli/app/app.go` - Add tests for init, update, view
- [ ] `internal/cli/app/messages.go` - Test all message types

#### Intent Models (if <80%)
- [ ] `capture_event_intent.go` - Test all states
- [ ] `browse_timeline_intent.go` - Test all states
- [ ] `generate_cv_intent.go` - Test all states
- [ ] `export_artifact_intent.go` - Test all states
- [ ] `configure_system_intent.go` - Test all states
- [ ] `burst_management_intent.go` - Test all states

#### CLI Models
- [ ] `form.go` - Test form validation, submission
- [ ] `burst_editor.go` - Test editor operations
- [ ] `fact_editor.go` - Test editor operations
- [ ] `metadata_editor.go` - Test editor operations

#### Components
- [ ] `standard_view.go` - Additional edge cases
- [ ] `modal.go` - Additional modal types
- [ ] `loading_messages.go` - Edge cases

#### Context
- [ ] `global.go` - Additional concurrent access tests

### Acceptance Criteria
- [ ] All new tests pass
- [ ] No regressions in existing tests
- [ ] Coverage improves by at least 5% per package

---

## Phase 4: Verify Coverage Improvement

### Steps

1. **Run full test suite with coverage**
   ```bash
   go test -v -cover ./... -coverprofile=coverage-new.out
   ```

2. **Compare before/after**
   ```bash
   # Before
   go tool cover -func=coverage.out | grep "total"
   # After
   go tool cover -func=coverage-new.out | grep "total"
   ```

3. **Generate updated report**
   ```bash
   go tool cover -html=coverage-new.out -o coverage-after.html
   ```

4. **Verify compliance**
   ```bash
   make check-compliance
   ```

### Acceptance Criteria
- [ ] Overall coverage ≥80%
- [ ] No package <70% coverage
- [ ] All tests passing (100% pass rate)
- [ ] Zero race conditions
- [ ] Compliance check passes

---

## Phase 5: Document and Update

### Tasks
- [ ] Update AGENTS.md with new coverage stats
- [ ] Update test documentation if patterns changed
- [ ] Commit all new tests with proper messages
- [ ] Update this task file with final results

### Acceptance Criteria
- [ ] AGENTS.md reflects new coverage (80%+)
- [ ] All commits follow atomic commit guidelines
- [ ] Task marked as complete

---

## Notes & Considerations

### Quick Wins
Focus on files with many untested lines but simple logic:
- Error handling paths (add error injection tests)
- View methods (test rendering with various states)
- Helper functions (straightforward unit tests)

### Difficult Areas
Some areas may be hard to test:
- TUI rendering (can use test harnesses)
- Terminal interactions (can mock)
- Concurrent operations (use race detector)

Don't force 100% - some code is legitimately hard to test (e.g., main() functions, panic handlers). Aim for reasonable coverage of business logic.

### Tools
- `go test -cover` - Basic coverage
- `go tool cover -html` - Visual coverage report
- `go test -race` - Race condition detection
- `make check-compliance` - Full compliance check

### Time Estimates
- Phase 1 (Analysis): 2-3 hours
- Phase 2 (Planning): 1-2 hours
- Phase 3 (Implementation): 1-2 days (depends on gap size)
- Phase 4 (Verification): 1 hour
- Phase 5 (Documentation): 1 hour

**Total**: 2-3 days

---

## Success Criteria

- [x] Task file created
- [ ] Coverage gaps identified and documented
- [ ] Test plan created
- [ ] Missing tests written
- [ ] Coverage ≥80% achieved
- [ ] All tests passing
- [ ] Zero race conditions
- [ ] Compliance check passes
- [ ] Documentation updated

---

## Resources

### Existing Test Patterns
- See `internal/cli/intents/*_test.go` for Ginkgo test examples
- See `internal/cli/intents/testing.go` for test harnesses
- See `internal/service/career/cv/*_test.go` for comprehensive test suites

### Coverage Documentation
- Go coverage tutorial: https://go.dev/blog/cover
- Ginkgo documentation: https://onsi.github.io/ginkgo/
- Gomega matchers: https://onsi.github.io/gomega/

### Related Files
- `Makefile` - `make test`, `make coverage` targets
- `docs/rules/go-guidelines.md` - Testing standards
- `AGENTS.md` - Current test coverage stats

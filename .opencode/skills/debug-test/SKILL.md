---
name: debug-test
description: Debug failing tests and common test issues in KaRiya
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

## What I do

Help diagnose and fix failing tests in the KaRiya codebase.

## When to use me

Use this skill when tests fail unexpectedly or you need to debug test issues.

## Running Tests

### All Tests
```bash
make test
```

### Specific Suite
```bash
make test-suite SUITE=./internal/cli/intents/myfeature/...
```

### Single Test
```bash
make individual-test TEST="should display items"
```

### Verbose Output
```bash
go test -v ./... -run "TestName"
```

### With Race Detection
```bash
make test-race
```

## Common Issues and Solutions

### 1. Race Conditions

**Symptom:** Tests pass individually but fail when run together.

**Diagnosis:**
```bash
go test -race ./path/to/package/...
```

**Common causes:**
- Shared state without mutex
- GlobalContext access
- Concurrent map access

**Fix:** Add proper synchronization or use thread-safe patterns.

### 2. Multiple Ginkgo Entry Points

**Symptom:** `Found more than one test suite file`

**Fix:** Each package should have exactly ONE `*_suite_test.go`:

```go
// package_suite_test.go
package mypackage_test

import (
    "testing"
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)

func TestMyPackage(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "MyPackage Suite")
}
```

### 3. Fixture Issues

**Symptom:** Tests fail with nil pointer or missing data.

**Check fixture usage:**
```bash
make check-fixtures
```

**Use fixtures, not inline structs:**
```go
// WRONG
event := &career.Event{Title: "Test"}

// CORRECT
event := fixtures.NewEvent().WithTitle("Test").Build()
```

### 4. Flaky Tests

**Symptom:** Tests pass sometimes, fail other times.

**Diagnosis:**
```bash
# Run multiple times
for i in {1..10}; do go test ./path/... || break; done
```

**Common causes:**
- Time-dependent logic
- Unordered map iteration
- External dependencies

### 5. Coverage Too Low

**Diagnosis:**
```bash
go test -coverprofile=/tmp/cover.out ./path/... && \
  go tool cover -func=/tmp/cover.out | grep -v "100.0%"
```

**View HTML report:**
```bash
go tool cover -html=/tmp/cover.out
```

### 6. Test Timeout

**Symptom:** Test hangs or times out.

**Run with timeout:**
```bash
go test -timeout 30s ./path/...
```

**Common causes:**
- Infinite loop
- Blocking channel operation
- Deadlock

### 7. Form Alignment Issues

**Symptom:** Huh form not centering/aligning properly.

**Fix:** Use form wrapper pattern:
```go
// Use models.*Form wrapper, not raw *huh.Form
form := models.NewCaptureForm(theme, config)
```

## Debugging Techniques

### Print Debugging
```go
fmt.Printf("DEBUG: value = %+v\n", value)
```

### Ginkgo Focused Tests
```go
FIt("focused test", func() { ... })
FDescribe("focused suite", func() { ... })
```

**Remember to remove `F` prefix before committing!**

### Breakpoints with Delve
```bash
dlv test ./path/to/package -- -test.run "TestName"
```

## Test Patterns

### Table-Driven Tests
```go
DescribeTable("validation",
    func(input string, expected bool) {
        result := Validate(input)
        Expect(result).To(Equal(expected))
    },
    Entry("valid input", "good", true),
    Entry("empty input", "", false),
)
```

### BDD Style
```go
Describe("MyFeature", func() {
    Context("when condition", func() {
        BeforeEach(func() { ... })
        
        It("should behave correctly", func() { ... })
    })
})
```

## Related skills

- `tdd-workflow` - Write tests properly
- `check-compliance` - Validate test coverage
- `create-screen` - Proper screen testing patterns

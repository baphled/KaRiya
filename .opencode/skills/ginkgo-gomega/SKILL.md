---
name: ginkgo-gomega
description: Ginkgo v2 BDD testing framework and Gomega assertion library for Go
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

# Ginkgo & Gomega Testing Skill

You are an expert in Ginkgo v2 BDD testing framework and Gomega assertion library for Go.

## Overview

KaRiya uses Ginkgo for all tests. Understanding these patterns is essential for TDD workflow.

## Test Suite Structure

### One Suite Per Package (REQUIRED)

```go
// mypackage_suite_test.go
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

**CRITICAL**: Only ONE `TestXxx` function per package. Multiple entry points cause failures.

---

## Ginkgo Structure Blocks

### Describe - Group Related Tests

```go
var _ = Describe("EventService", func() {
    // Tests for EventService
})
```

### Context - Describe Conditions

```go
Describe("Save", func() {
    Context("when event is valid", func() {
        It("persists the event", func() { ... })
    })
    
    Context("when event is invalid", func() {
        It("returns validation error", func() { ... })
    })
})
```

### It - Individual Test Case

```go
It("creates an event with the given text", func() {
    event, err := service.Create(ctx, "Started new project")
    
    Expect(err).ToNot(HaveOccurred())
    Expect(event.Text).To(Equal("Started new project"))
})
```

### When - Alias for Context

```go
When("the repository returns an error", func() {
    It("propagates the error", func() { ... })
})
```

---

## Setup and Teardown

### BeforeEach / AfterEach

```go
var _ = Describe("EventService", func() {
    var (
        service *EventService
        repo    *MockEventRepository
        ctrl    *gomock.Controller
    )
    
    BeforeEach(func() {
        ctrl = gomock.NewController(GinkgoT())
        repo = NewMockEventRepository(ctrl)
        service = NewEventService(repo)
    })
    
    AfterEach(func() {
        ctrl.Finish()
    })
    
    It("...", func() { ... })
})
```

### BeforeSuite / AfterSuite

```go
var sharedDB *gorm.DB

var _ = BeforeSuite(func() {
    var err error
    sharedDB, err = setupTestDatabase()
    Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {
    if sharedDB != nil {
        sqlDB, _ := sharedDB.DB()
        sqlDB.Close()
    }
})
```

---

## Gomega Assertions

### Basic Matchers

```go
// Equality
Expect(value).To(Equal(expected))
Expect(value).ToNot(Equal(other))

// Nil checks
Expect(err).To(BeNil())
Expect(err).ToNot(HaveOccurred())  // Preferred for errors
Expect(result).ToNot(BeNil())

// Boolean
Expect(isValid).To(BeTrue())
Expect(isEmpty).To(BeFalse())

// Numeric
Expect(count).To(BeNumerically(">", 0))
Expect(count).To(BeNumerically(">=", 5))
Expect(value).To(BeZero())
```

### Collection Matchers

```go
// Length
Expect(items).To(HaveLen(5))
Expect(items).To(BeEmpty())
Expect(items).ToNot(BeEmpty())

// Contains
Expect(items).To(ContainElement(expected))
Expect(items).To(ContainElements(a, b, c))

// Key checks (maps)
Expect(myMap).To(HaveKey("name"))
Expect(myMap).To(HaveKeyWithValue("name", "John"))
```

### String Matchers

```go
Expect(text).To(ContainSubstring("error"))
Expect(text).To(HavePrefix("Error:"))
Expect(text).To(HaveSuffix(".go"))
Expect(text).To(MatchRegexp(`^\d{4}-\d{2}-\d{2}$`))
```

### Error Matchers

```go
// Error occurred
Expect(err).To(HaveOccurred())
Expect(err).ToNot(HaveOccurred())

// Error type
Expect(err).To(MatchError("specific error message"))
Expect(err).To(MatchError(ContainSubstring("failed")))

// Wrapped errors
Expect(errors.Is(err, ErrNotFound)).To(BeTrue())
```

### Pointer Matchers

```go
Expect(ptr).To(BeNil())
Expect(ptr).ToNot(BeNil())
Expect(ptr).To(PointTo(Equal(expectedValue)))
```

### Struct Matchers

```go
Expect(event).To(Equal(&Event{
    ID:   "123",
    Text: "Test",
}))

// Field matching
Expect(event.Text).To(Equal("Test"))
Expect(event.CreatedAt).ToNot(BeZero())
```

---

## Table-Driven Tests

### DescribeTable Pattern

```go
DescribeTable("Validate",
    func(input string, expectedValid bool) {
        err := Validate(input)
        if expectedValid {
            Expect(err).ToNot(HaveOccurred())
        } else {
            Expect(err).To(HaveOccurred())
        }
    },
    Entry("valid input", "hello", true),
    Entry("empty input", "", false),
    Entry("too long", strings.Repeat("x", 1001), false),
)
```

### With Complex Structs

```go
DescribeTable("Event creation",
    func(text string, company string, expectedErr error) {
        event, err := NewEvent(text, company)
        
        if expectedErr != nil {
            Expect(err).To(MatchError(expectedErr))
            Expect(event).To(BeNil())
        } else {
            Expect(err).ToNot(HaveOccurred())
            Expect(event.Text).To(Equal(text))
        }
    },
    Entry("valid event", "Did something", "Acme", nil),
    Entry("empty text", "", "Acme", ErrEmptyText),
    Entry("missing company", "Did something", "", ErrMissingCompany),
)
```

---

## Focus and Skip

### Focus Tests (Development Only)

```go
FDescribe("focused describe", func() { ... })
FContext("focused context", func() { ... })
FIt("focused test", func() { ... })
FWhen("focused when", func() { ... })
```

### Skip Tests

```go
XDescribe("skipped describe", func() { ... })
XContext("skipped context", func() { ... })
XIt("skipped test", func() { ... })

// Programmatic skip
It("skips conditionally", func() {
    if condition {
        Skip("Skipping because...")
    }
})
```

**WARNING**: Never commit focused tests (F-prefix). CI will fail.

---

## Async Testing

### Eventually / Consistently

```go
// Wait for condition
Eventually(func() int {
    return len(service.GetItems())
}).Should(BeNumerically(">", 0))

// With timeout
Eventually(func() bool {
    return service.IsReady()
}, 5*time.Second, 100*time.Millisecond).Should(BeTrue())

// Consistently (stays true)
Consistently(func() int {
    return service.Count()
}, 1*time.Second).Should(Equal(5))
```

### Channel Matchers

```go
Expect(channel).To(Receive())
Expect(channel).To(Receive(Equal(expectedValue)))
Expect(channel).To(BeClosed())
```

---

## Test Organization Patterns

### KaRiya Naming Conventions

```go
// Unit tests: *_test.go
event_service_test.go

// E2E tests: *_e2e_test.go
capture_e2e_test.go

// Navigation tests: *_navigation_test.go
browse_navigation_test.go
```

### Documentation Comments in Tests

```go
var _ = Describe("FactService", func() {
    Describe("ExtractFacts", func() {
        Context("given an event with clear competencies", func() {
            It("extracts skill-based facts", func() {
                // Test implementation
            })
        })
    })
})
```

---

## Common Patterns

### Testing Services with Mocks

```go
var _ = Describe("EventService", func() {
    var (
        service *EventService
        repo    *MockEventRepository
        ctrl    *gomock.Controller
        ctx     context.Context
    )
    
    BeforeEach(func() {
        ctx = context.Background()
        ctrl = gomock.NewController(GinkgoT())
        repo = NewMockEventRepository(ctrl)
        service = NewEventService(repo)
    })
    
    AfterEach(func() {
        ctrl.Finish()
    })
    
    Describe("Create", func() {
        It("saves the event to repository", func() {
            repo.EXPECT().
                Save(ctx, gomock.Any()).
                Return(nil)
            
            event, err := service.Create(ctx, "Test event")
            
            Expect(err).ToNot(HaveOccurred())
            Expect(event).ToNot(BeNil())
        })
    })
})
```

### Testing Validation

```go
Describe("Validate", func() {
    Context("when text is empty", func() {
        It("returns validation error", func() {
            event := &Event{Text: ""}
            
            err := event.Validate()
            
            Expect(err).To(HaveOccurred())
            Expect(err.Error()).To(ContainSubstring("text"))
        })
    })
})
```

---

## Anti-Patterns

### DON'T: Multiple Test Entry Points

```go
// WRONG - Multiple TestXxx functions
func TestEventService(t *testing.T) { ... }
func TestFactService(t *testing.T) { ... }  // Will cause issues!

// CORRECT - One suite entry point
func TestMyPackage(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "MyPackage Suite")
}
```

### DON'T: Commit Focused Tests

```go
// WRONG - Will cause CI failure
FIt("my test", func() { ... })

// CORRECT
It("my test", func() { ... })
```

### DON'T: Use testify/assert

```go
// WRONG - Don't mix assertion libraries
assert.Equal(t, expected, actual)

// CORRECT - Use Gomega
Expect(actual).To(Equal(expected))
```

### DON'T: Skip BeforeEach/AfterEach

```go
// WRONG - Setup in test
It("does something", func() {
    ctrl := gomock.NewController(GinkgoT())
    repo := NewMockRepo(ctrl)
    // ... test ...
    ctrl.Finish()  // Easy to forget!
})

// CORRECT - Setup in BeforeEach
BeforeEach(func() {
    ctrl = gomock.NewController(GinkgoT())
    repo = NewMockRepo(ctrl)
})
AfterEach(func() {
    ctrl.Finish()
})
```

---

## Running Tests

```bash
# Run all tests
make test

# Run specific suite
make test-suite SUITE=./internal/service/career/...

# Run single test by name
make individual-test TEST="creates an event"

# Verbose output
go test -v ./... -ginkgo.v

# Run focused tests only (development)
go test ./... -ginkgo.focus="my test"

# Skip tests matching pattern
go test ./... -ginkgo.skip="slow"
```

---

## Related Skills

- `gomock` - Mock generation and usage
- `test-fixtures` - Factory patterns for test data
- `e2e-testing` - End-to-end test patterns
- `tdd-workflow` - Red-Green-Refactor cycle

---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# Task 22: Add Comprehensive Strategy System Tests

**Created**: 2026-01-07
**Status**: Ready for Implementation
**Priority**: LOW (Nice to Have)
**Estimated Time**: 2-3 hours
**Related**: Task 17 (Form Refactoring), Task 19 (Test Fixes), Task 21 (Documentation)

---

## Overview

Add comprehensive tests for the strategy system (quick/manual) to ensure robust coverage of field visibility, strategy selection, and toggle behavior.

**Current Status**:
- Strategy system implemented and working (Task 17)
- All 2,078 tests passing
- Basic functionality verified through persistence tests
- Need: Dedicated strategy-focused tests

---

## Test Coverage Goals

### Current Coverage
- ✅ Form submission works (persistence tests)
- ✅ Field visibility basic behavior (form tests)
- ✅ Manual strategy with showOptionalFields=true (default)

### Missing Coverage
- ⚠️ Quick strategy behavior (minimal fields)
- ⚠️ Field toggle ('t' key) in both strategies
- ⚠️ Strategy selection workflow
- ⚠️ Field navigation skipping hidden fields
- ⚠️ Edge cases (toggle on focused hidden field, etc.)

---

## Test Categories

### 1. Strategy Selection Tests (30 min)

**File**: `internal/cli/models/form_strategy_test.go` (new)

**Tests to Add**:
```go
Describe("FormModel Strategy System", func() {
    Describe("Default Strategy", func() {
        It("should default to manual strategy", func() {})
        It("should show all fields by default in manual strategy", func() {})
    })

    Describe("SetStrategy", func() {
        It("should set strategy to quick", func() {})
        It("should set strategy to manual", func() {})
        It("should hide optional fields when strategy is quick", func() {})
        It("should show optional fields when strategy is manual", func() {})
    })

    Describe("Strategy Persistence", func() {
        It("should maintain strategy across form lifecycle", func() {})
        It("should not change strategy on field navigation", func() {})
    })
})
```

**Estimated**: 8-10 test specs

### 2. Field Visibility Tests (45 min)

**File**: `internal/cli/models/form_visibility_test.go` (new)

**Tests to Add**:
```go
Describe("Field Visibility", func() {
    Describe("Quick Strategy", func() {
        It("should only show TextField and DateField", func() {})
        It("should hide Company, Project, Tags, Categories", func() {})
        It("should always show SubmitButton", func() {})
    })

    Describe("Manual Strategy", func() {
        Context("with showOptionalFields=true", func() {
            It("should show all fields", func() {})
        })

        Context("with showOptionalFields=false", func() {
            It("should hide optional fields", func() {})
            It("should still show TextField and DateField", func() {})
        })
    })

    Describe("Field Toggle (t key)", func() {
        Context("in Quick strategy", func() {
            It("should show optional fields when t pressed", func() {})
            It("should hide optional fields when t pressed again", func() {})
        })

        Context("in Manual strategy", func() {
            It("should hide optional fields when t pressed", func() {})
            It("should show optional fields when t pressed again", func() {})
        })
    })

    Describe("isFieldVisible helper", func() {
        It("should return true for TextField in all cases", func() {})
        It("should return true for SubmitButton in all cases", func() {})
        It("should return correct visibility for optional fields", func() {})
    })
})
```

**Estimated**: 12-15 test specs

### 3. Field Navigation Tests (45 min)

**File**: `internal/cli/models/form_navigation_test.go` (new)

**Tests to Add**:
```go
Describe("Field Navigation with Hidden Fields", func() {
    Describe("Tab Navigation", func() {
        Context("in Quick strategy (fields hidden)", func() {
            It("should skip from TextField to DateField", func() {})
            It("should skip from DateField to SubmitButton", func() {})
            It("should not focus on hidden Company field", func() {})
        })

        Context("in Manual strategy (fields visible)", func() {
            It("should navigate through all fields in order", func() {})
            It("should focus on Company, Project, Tags, Categories", func() {})
        })

        Context("when toggling visibility", func() {
            It("should update navigation to include newly visible fields", func() {})
            It("should skip newly hidden fields", func() {})
        })
    })

    Describe("Shift+Tab (Backward Navigation)", func() {
        It("should skip hidden fields when navigating backward", func() {})
        It("should wrap around correctly with hidden fields", func() {})
    })

    Describe("Edge Cases", func() {
        It("should handle toggle when focused on TextField", func() {})
        It("should not get stuck in hidden field loop", func() {})
        It("should handle wrap-around with all fields hidden except required", func() {})
    })
})
```

**Estimated**: 10-12 test specs

### 4. Integration Tests (30 min)

**File**: `internal/cli/models/form_strategy_integration_test.go` (new)

**Tests to Add**:
```go
Describe("Strategy System Integration", func() {
    Describe("Quick Capture Workflow", func() {
        It("should allow submission with only TextField and DateField", func() {})
        It("should not require Company or Project", func() {})
        It("should save event with minimal data", func() {})
    })

    Describe("Manual Capture Workflow", func() {
        It("should accept event with all fields filled", func() {})
        It("should save event with complete metadata", func() {})
    })

    Describe("Toggle During Entry", func() {
        It("should allow toggling from Quick to Manual mid-entry", func() {})
        It("should preserve entered data when toggling", func() {})
        It("should allow toggling from Manual to Quick mid-entry", func() {})
    })

    Describe("Edit Mode", func() {
        It("should use manual strategy for editing", func() {})
        It("should show all event data when editing", func() {})
    })
})
```

**Estimated**: 8-10 test specs

### 5. Edge Case Tests (30 min)

**File**: Add to existing `internal/cli/models/form_test.go`

**Tests to Add**:
```go
Describe("Strategy Edge Cases", func() {
    It("should handle rapid toggle presses", func() {})
    It("should handle toggle during form submission", func() {})
    It("should handle window resize with hidden fields", func() {})
    It("should render correctly with minimal terminal size", func() {})
    It("should handle validation errors on hidden-then-shown fields", func() {})
})
```

**Estimated**: 5-7 test specs

---

## Implementation Plan

### Phase 1: Setup and Basic Tests (45 min)

**Tasks**:
- [ ] Create `form_strategy_test.go`
- [ ] Add strategy selection tests (8-10 specs)
- [ ] Verify tests pass
- [ ] Run with race detector

### Phase 2: Field Visibility Tests (45 min)

**Tasks**:
- [ ] Create `form_visibility_test.go`
- [ ] Add visibility tests (12-15 specs)
- [ ] Add toggle tests
- [ ] Verify all tests pass

### Phase 3: Navigation Tests (45 min)

**Tasks**:
- [ ] Create `form_navigation_test.go`
- [ ] Add tab navigation tests
- [ ] Add shift+tab tests
- [ ] Add edge case navigation tests
- [ ] Verify navigation logic works correctly

### Phase 4: Integration Tests (45 min)

**Tasks**:
- [ ] Create `form_strategy_integration_test.go`
- [ ] Add end-to-end workflow tests
- [ ] Add toggle during entry tests
- [ ] Add edit mode tests
- [ ] Verify complete workflows

### Phase 5: Verification (30 min)

**Tasks**:
- [ ] Run full test suite
- [ ] Check coverage for form_test.go
- [ ] Run with race detector
- [ ] Run benchmarks
- [ ] Document test coverage improvements

---

## Acceptance Criteria

### Must Have
- [ ] 40-50 new test specs for strategy system
- [ ] All tests passing (maintain 100% pass rate)
- [ ] Zero race conditions
- [ ] Coverage for all strategy scenarios

### Should Have
- [ ] Tests cover edge cases
- [ ] Tests cover both quick and manual strategies
- [ ] Tests verify field toggle behavior
- [ ] Tests verify navigation with hidden fields

### Nice to Have
- [ ] Performance benchmarks for toggle operation
- [ ] Tests for rapid toggle scenarios
- [ ] Tests for all terminal sizes

---

## Test Structure Example

```go
package models_test

import (
    "testing"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"

    "github.com/baphled/kariya/internal/cli/models"
    "github.com/baphled/kariya/internal/cli/service"
    careerservice "github.com/baphled/kariya/internal/service/career"
    careerrepo "github.com/baphled/kariya/internal/repository/career"
    tea "github.com/charmbracelet/bubbletea"
)

var _ = Describe("FormModel Strategy System", func() {
    var (
        form   *models.FormModel
        cliSvc *service.CLIEventService
        repo   *careerrepo.MemoryRepository
    )

    BeforeEach(func() {
        repo = careerrepo.NewMemoryRepository()
        svc := careerservice.NewService(repo)
        cliSvc = service.NewCLIEventService(svc)
        form = models.NewFormModel(cliSvc)
    })

    Describe("Default Behavior", func() {
        It("should default to manual strategy", func() {
            Expect(form.GetStrategy()).To(Equal("manual"))
        })

        It("should show all fields by default", func() {
            Expect(form.IsFieldVisible(models.CompanyField)).To(BeTrue())
            Expect(form.IsFieldVisible(models.ProjectField)).To(BeTrue())
        })
    })

    Describe("Quick Strategy", func() {
        BeforeEach(func() {
            form.SetStrategy("quick")
        })

        It("should hide optional fields", func() {
            Expect(form.IsFieldVisible(models.CompanyField)).To(BeFalse())
            Expect(form.IsFieldVisible(models.ProjectField)).To(BeFalse())
        })

        It("should show required fields", func() {
            Expect(form.IsFieldVisible(models.TextField)).To(BeTrue())
            Expect(form.IsFieldVisible(models.DateField)).To(BeTrue())
        })
    })
})
```

---

## Performance Benchmarks

Add benchmarks to measure toggle performance:

```go
func BenchmarkFieldToggle(b *testing.B) {
    repo := careerrepo.NewMemoryRepository()
    svc := careerservice.NewService(repo)
    cliSvc := service.NewCLIEventService(svc)
    form := models.NewFormModel(cliSvc)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        form.ToggleOptionalFields()
    }
}

func BenchmarkFieldVisibilityCheck(b *testing.B) {
    repo := careerrepo.NewMemoryRepository()
    svc := careerservice.NewService(repo)
    cliSvc := service.NewCLIEventService(svc)
    form := models.NewFormModel(cliSvc)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = form.IsFieldVisible(models.CompanyField)
    }
}
```

---

## Success Metrics

- ✅ 40-50 new tests for strategy system
- ✅ All 2,100+ tests passing (100% pass rate)
- ✅ Strategy system coverage >95%
- ✅ Zero race conditions
- ✅ All edge cases covered
- ✅ Performance benchmarks documented

---

## Notes

### Implementation Tips
1. Use Ginkgo BeforeEach for test setup
2. Test both strategies exhaustively
3. Verify field visibility after every toggle
4. Test navigation paths through hidden fields
5. Use race detector frequently

### Things to Watch
- Field focus when toggling visibility
- Navigation wrap-around with hidden fields
- Form submission with minimal vs full data
- Edit mode always uses manual strategy

---

**Last Updated**: 2026-01-07
**Author**: AI Assistant (via OpenCode)
**Status**: Optional - Project is production ready without these tests

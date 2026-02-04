# Capture Workflow BDD Test Plan

## Purpose
This document outlines the comprehensive BDD test coverage needed for all models in the capture workflow to enable confident migration to the new screens-based architecture.

## Test Strategy
- **Behavior-focused**: Tests describe WHAT the model does, not HOW
- **Ginkgo/Gomega**: Use BDD style with descriptive test names
- **Complete coverage**: All public methods, all edge cases, all state transitions
- **Migration support**: Tests serve as behavioral specification for new architecture

## Models to Test

### 1. CaptureForm (`capture_form.go`)
**Status**: ✅ Tests created (`capture_form_test.go`)

**Behaviors Covered**:
- Initialization with default strategy
- Strategy management (Quick vs Manual)
- Event loading for editing
- Form submission (via Ctrl+S and form completion)
- Cancel behavior (Escape key handling)
- Window resizing
- View rendering
- Date parsing and validation
- Empty form handling
- Special characters preservation
- State transitions

**Test Count**: ~40 test cases

---

### 2. MetadataEditorModelNew (`metadata_editor_new.go`)
**Status**: ✅ Partial tests exist (`metadata_editor_new_test.go`)

**Behaviors to Add**:
- ✅ Basic initialization
- ✅ Form completion handling
- ✅ Cancel handling
- ⚠️ **MISSING**: Tag selector integration
- ⚠️ **MISSING**: Category selector integration  
- ⚠️ **MISSING**: Skill selector integration
- ⚠️ **MISSING**: Multi-select behavior
- ⚠️ **MISSING**: Validation errors
- ⚠️ **MISSING**: Persistence to CLIEventService
- ⚠️ **MISSING**: Revert functionality
- ⚠️ **MISSING**: Modal overlay integration (GetTitle, GetContent, GetFooter)
- ⚠️ **MISSING**: Error handling and display
- ⚠️ **MISSING**: Dimensions handling

**Additional Tests Needed**: ~30 test cases

---

### 3. BurstSuggestionModelNew (`burst_suggestion_new.go`)
**Status**: ✅ Tests implemented (`burst_suggestion_new_test.go`)

**Behaviors to Test**:
- Initialization with suggestions
- Navigation through suggestions (up/down, j/k)
- Confirm suggestion (y key)
- Reject suggestion (n key)
- Edit mode activation (e key)
- Editing name/description via huh form
- Saving edits
- Cancelling edits (Escape in edit mode)
- Progress tracking (confirmed vs rejected)
- Related events loading and caching
- Completion detection (all suggestions processed)
- Confidence score visualization
- Modal overlay integration
- Message generation (ConfirmBurstMsg, RejectBurstSuggestionMsg)
- BurstProcessingCompleteMsg
- Empty suggestions handling
- Window resizing
- Theme integration

**Test File**: `burst_suggestion_new_test.go`  
**Test Count**: ~45 test cases

---

### 4. FactEditorModelNew (`fact_editor_new.go`)
**Status**: ❌ No tests exist

**Behaviors to Test**:
- Initialization with fact
- Form completion and validation
- Fact field editing (Text, CompetencyCategories, RoleFit, etc.)
- Form submission with repository update
- Cancel handling
- Revert functionality
- Error handling
- Modal overlay integration (GetTitle, GetContent, GetFooter)
- Timestamp updates
- Window resizing
- Quit handling
- Theme integration
- Form state tracking (submitted, cancelled)
- Validation errors

**Test File**: `fact_editor_new_test.go`  
**Test Count**: ~35 test cases

---

## Test Organization

### File Structure
```
internal/cli/models/
├── models_suite_test.go          # Ginkgo suite setup
├── capture_form_test.go          # ✅ Created
├── metadata_editor_new_test.go   # ✅ Exists (needs expansion)
├── burst_suggestion_new_test.go  # ❌ Needs creation
└── fact_editor_new_test.go       # ❌ Needs creation
```

### Test Pattern
Each test file follows this structure:

```go
var _ = Describe("ModelName", func() {
    var (
        model *models.ModelName
        // Dependencies
    )

    BeforeEach(func() {
        // Setup fresh model
    })

    Describe("Behavior Category", func() {
        Context("when specific condition", func() {
            It("should exhibit expected behavior", func() {
                // Arrange, Act, Assert
            })
        })
    })
})
```

## Test Coverage Requirements

### Minimum Coverage Targets
- **Models package**: 90%+ (currently at 9.1%)
- **Each model**: 85%+ line coverage
- **Critical paths**: 100% (submission, validation, state transitions)

### Critical Behaviors (Must Have 100% Coverage)
1. **Form submission flows**
   - CaptureForm submission
   - MetadataEditor save
   - FactEditor save
   - BurstSuggestion confirm/reject

2. **Cancel/Revert flows**
   - Escape key handling
   - Revert functionality
   - State cleanup

3. **Validation logic**
   - Field validation
   - Error handling
   - User feedback

4. **State transitions**
   - Form state changes
   - Edit mode toggles
   - Completion detection

5. **Data integrity**
   - Persistence to repository
   - Data copying (not referencing)
   - Timestamp management

## BDD Scenarios

### Example: BurstSuggestionModelNew

```gherkin
Feature: Burst Suggestion Review
  As a user reviewing burst suggestions
  I want to navigate, edit, accept, and reject suggestions
  So that I can curate my career bursts

  Scenario: Navigate through multiple suggestions
    Given I have 3 burst suggestions
    When I press the down arrow
    Then I should see the second suggestion
    And the progress bar should show "2/3"

  Scenario: Edit a suggestion name
    Given I am viewing a burst suggestion
    When I press 'e' to edit
    Then I should see the edit form
    When I enter a new name "Leadership Period"
    And I press Enter to save
    Then the suggestion should have the updated name
    And I should return to the review view

  Scenario: Confirm a suggestion
    Given I am viewing a burst suggestion
    When I press 'y' to confirm
    Then a ConfirmBurstMsg should be sent
    And the suggestion should be added to confirmed list
    And I should move to the next suggestion

  Scenario: Complete review with mixed acceptances
    Given I have 3 burst suggestions
    When I confirm the first suggestion
    And I reject the second suggestion
    And I confirm the third suggestion
    Then a BurstProcessingCompleteMsg should be sent
    And confirmed count should be 2
    And rejected count should be 1
```

## Integration Points

### Models interact with:
1. **Forms package** (`internal/cli/forms/`)
   - Form builders and validators
   - FormData structures
   - Field configuration

2. **Service layer** (`internal/cli/service/`)
   - CLIEventService for persistence
   - CareerService for enrichment

3. **Repositories** (`internal/repository/career/`)
   - EventRepository
   - BurstRepository
   - FactRepository
   - SkillRepository

4. **UIKit** (`internal/cli/uikit/`)
   - Theme integration
   - Layout components
   - Container components

### Test Isolation
- Use memory repositories for fast, isolated tests
- Mock services where appropriate
- Test UI kit integration at boundaries only

## Migration Strategy

### Phase 1: Complete BDD Coverage (This PR)
- ✅ CaptureForm tests
- ⏳ Expand MetadataEditor tests
- ⏳ Create BurstSuggestion tests
- ⏳ Create FactEditor tests

### Phase 2: Architecture Migration (Future PR)
- Extract screens from models
- Use BDD tests to validate behavior preservation
- Models become lightweight wrappers
- Screens handle UI rendering

### Phase 3: Deprecate Legacy (Future PR)
- Mark models as deprecated
- Migrate all intents to screens
- Remove legacy models

## Success Criteria

✅ **Definition of Done**:
1. All 4 model files have comprehensive BDD test files
2. Test coverage for models package > 85%
3. All critical paths have 100% coverage
4. Tests run fast (< 5 seconds for full suite)
5. Tests are maintainable and well-documented
6. No flaky tests
7. Tests serve as executable specification

## References

- [BDD Workflow Guide](../development/BDD_WORKFLOW.md)
- [Navigation Testing Guide](../development/NAVIGATION_TESTING_GUIDE.md)
- [Forms Guide](../FORMS_GUIDE.md)
- [Intent Architecture Guide](../INTENT_ARCHITECTURE_GUIDE.md)

## Notes

- Tests use Ginkgo v2 BDD framework
- Gomega matchers provide expressive assertions
- Memory repositories enable fast, isolated tests
- Tests describe behavior, not implementation
- Tests will guide the migration to screens architecture

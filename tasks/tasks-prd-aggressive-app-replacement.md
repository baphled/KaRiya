# Task List: Aggressive app.go Replacement with Intent-Based Architecture

**PRD Reference**: `prd-aggressive-app-replacement.md`
**Timeline**: 3 weeks (121 hours)
**Status**: Phase 2 - Detailed Sub-Tasks Generated

---

## Relevant Files

### Core Application Files
- `internal/cli/app/app.go` - Root application model (to be completely rebuilt, ~1,766 lines → ~200 lines)
- `internal/cli/app/messages.go` - Application-level message types

### Intent Framework (Already Exists)
- `internal/cli/intents/contract.go` - Intent and IntentRouter interfaces
- `internal/cli/intents/result.go` - IntentResult[T] and IntentError types
- `internal/cli/intents/router.go` - IntentRouter implementation
- `internal/cli/intents/testing.go` - Test utilities and harnesses

### Existing Intent Implementations (5 - Already Complete)
- `internal/cli/intents/capture_event.go` - Capture event context and state
- `internal/cli/intents/capture_event_intent.go` - Capture event intent implementation
- `internal/cli/intents/capture_event_views_test.go` - Capture event tests
- `internal/cli/intents/browse_timeline.go` - Browse timeline context and state
- `internal/cli/intents/browse_timeline_intent.go` - Browse timeline intent implementation
- `internal/cli/intents/browse_timeline_test.go` - Browse timeline tests
- `internal/cli/intents/generate_cv.go` - Generate CV context and state
- `internal/cli/intents/generate_cv_intent.go` - Generate CV intent implementation
- `internal/cli/intents/generate_cv_test.go` - Generate CV tests
- `internal/cli/intents/export_artifact.go` - Export artifact context and state
- `internal/cli/intents/export_artifact_intent.go` - Export artifact intent implementation
- `internal/cli/intents/export_artifact_test.go` - Export artifact tests
- `internal/cli/intents/configure_system.go` - Configure system context and state
- `internal/cli/intents/configure_system_intent.go` - Configure system intent implementation
- `internal/cli/intents/configure_system_test.go` - Configure system tests

### New Intent Implementations (5 - To Be Created)
- `internal/cli/intents/burst_management.go` - Burst management context and state
- `internal/cli/intents/burst_management_intent.go` - Burst management intent implementation
- `internal/cli/intents/burst_management_test.go` - Burst management tests
- `internal/cli/intents/fact_management.go` - Fact management context and state
- `internal/cli/intents/fact_management_intent.go` - Fact management intent implementation
- `internal/cli/intents/fact_management_test.go` - Fact management tests
- `internal/cli/intents/import_wizard.go` - Import wizard context and state
- `internal/cli/intents/import_wizard_intent.go` - Import wizard intent implementation
- `internal/cli/intents/import_wizard_test.go` - Import wizard tests
- `internal/cli/intents/metadata_editor.go` - Metadata editor context and state
- `internal/cli/intents/metadata_editor_intent.go` - Metadata editor intent implementation
- `internal/cli/intents/metadata_editor_test.go` - Metadata editor tests
- `internal/cli/intents/bulk_operations.go` - Bulk operations context and state
- `internal/cli/intents/bulk_operations_intent.go` - Bulk operations intent implementation
- `internal/cli/intents/bulk_operations_test.go` - Bulk operations tests

### Legacy Model Files (To Be Removed - 30+ Files)
- `internal/cli/models/cv_config_manager.go` - (Handled by ConfigureSystem intent)
- `internal/cli/models/cv_generator.go` - (Handled by GenerateCV intent)
- `internal/cli/models/cv_preview.go` - (Handled by GenerateCV intent)
- `internal/cli/models/cv_export_dialog.go` - (Handled by ExportArtifact intent)
- `internal/cli/models/cv_list.go` - (Handled by BrowseTimeline intent)
- `internal/cli/models/list.go` - (Handled by various intents)
- `internal/cli/models/details.go` - (Handled by BrowseTimeline intent)
- `internal/cli/models/burst_list.go` - (Handled by BurstManagement intent)
- `internal/cli/models/burst_editor.go` - (Handled by BurstManagement intent)
- `internal/cli/models/burst_details.go` - (Handled by BurstManagement intent)
- `internal/cli/models/burst_suggestion.go` - (Handled by BurstManagement intent)
- `internal/cli/models/fact_list.go` - (Handled by FactManagement intent)
- `internal/cli/models/fact_editor.go` - (Handled by FactManagement intent)
- `internal/cli/models/fact_details.go` - (Handled by FactManagement intent)
- `internal/cli/models/facts_results.go` - (Handled by FactManagement intent)
- `internal/cli/models/import_review.go` - (Handled by ImportWizard intent)
- `internal/cli/models/metadata_review.go` - (Handled by MetadataEditor intent)
- `internal/cli/models/metadata_editor.go` - (Handled by MetadataEditor intent)
- `internal/cli/models/bulk_operations.go` - (Handled by BulkOperations intent)
- And 11+ other legacy screen models and utilities

### Testing Files
- `internal/cli/intents/contract_test.go` - Intent framework tests
- `internal/cli/intents/router_test.go` - Router implementation tests
- `internal/cli/intents/result_test.go` - Result type tests
- `internal/cli/intents/testing_test.go` - Testing utilities tests
- `internal/cli/intents/benchmarks_test.go` - Performance benchmarks

### Configuration & Documentation
- `internal/cli/app/messages.go` - Root application message types
- `internal/cli/context/global.go` - GlobalContext for shared state
- `internal/cli/styles/styles.go` - Centralized styling
- `internal/cli/components/` - Reusable UI components

### Notes

- The task list follows the 4-phase approach outlined in the PRD: Preparation (9h), Create Missing Intents (64h), Rebuild app.go (14h), Testing & Validation (34h)
- All test files should use Ginkgo v2 + Gomega framework, consistent with existing tests
- Performance benchmarks must be included in `internal/cli/intents/benchmarks_test.go`
- Legacy model files will be gradually removed as corresponding intents are completed
- The root `app.go` rebuild should be done last to ensure all intents are ready for integration

---

## Tasks

### 1.0 Preparation & Codebase Audit (9 hours)

- [ ] 1.1 Document legacy screen-to-intent mapping
  - [ ] 1.1.1 Create comprehensive audit document listing all 29+ legacy screens
  - [ ] 1.1.2 Map each legacy screen to corresponding intent (5 existing + 5 new)
  - [ ] 1.1.3 Note special logic, edge cases, and data transformations for each screen
  - [ ] 1.1.4 Identify cross-model dependencies and shared utilities
  - [ ] 1.1.5 Document any custom styling or UI patterns unique to each screen

- [ ] 1.2 Analyze feature gaps and requirements
  - [ ] 1.2.1 Review BurstManagement features from legacy burst_list, burst_editor, burst_details, burst_suggestion models
  - [ ] 1.2.2 Review FactManagement features from legacy fact_list, fact_editor, fact_details, facts_results models
  - [ ] 1.2.3 Review ImportWizard features from legacy import_review model
  - [ ] 1.2.4 Review MetadataEditor features from legacy metadata_review, metadata_editor models
  - [ ] 1.2.5 Review BulkOperations features from legacy bulk_operations model
  - [ ] 1.2.6 Identify any missing features not yet covered by existing intents

- [ ] 1.3 Plan new intent implementations
  - [ ] 1.3.1 Define state machines for each of 5 new intents
  - [ ] 1.3.2 Document context and result structures for each intent
  - [ ] 1.3.3 Identify reusable components needed for new intents
  - [ ] 1.3.4 Plan test coverage targets (30-40 tests per intent)
  - [ ] 1.3.5 Document any async operations or external dependencies

- [ ] 1.4 Verify intent framework completeness
  - [ ] 1.4.1 Review IntentRouter implementation for all required methods
  - [ ] 1.4.2 Verify IntentResult[T] type safety implementation
  - [ ] 1.4.3 Check testing utilities in testing.go for adequacy
  - [ ] 1.4.4 Confirm all 5 existing intents follow consistent patterns
  - [ ] 1.4.5 Document any framework enhancements needed before new intents

---

### 2.0 Implement 5 New Missing Intents (64 hours)

#### 2.1 BurstManagement Intent (16 hours)

- [ ] 2.1.1 Create burst_management.go with context and result structures
  - [ ] 2.1.1.1 Define BurstManagementState enum (List, View, Edit, Suggest, Confirm)
  - [ ] 2.1.1.2 Define BurstManagementContext struct with required fields
  - [ ] 2.1.1.3 Define BurstManagementResult struct with output fields
  - [ ] 2.1.1.4 Add comprehensive documentation and examples

- [ ] 2.1.2 Create burst_management_intent.go with full implementation
  - [ ] 2.1.2.1 Implement Init() to load bursts from service
  - [ ] 2.1.2.2 Implement Update() with state machine for all 5 states
  - [ ] 2.1.2.3 Implement View() with rendering for each state
  - [ ] 2.1.2.4 Implement Result() returning typed IntentResult
  - [ ] 2.1.2.5 Implement state transition handlers (handleListState, handleViewState, etc.)
  - [ ] 2.1.2.6 Implement helper methods for burst operations (create, edit, delete, suggest)

- [ ] 2.1.3 Create burst_management_test.go with comprehensive tests
  - [ ] 2.1.3.1 Test state machine transitions (List → View → Edit → Confirm)
  - [ ] 2.1.3.2 Test view rendering for each state (10+ tests)
  - [ ] 2.1.3.3 Test burst CRUD operations (create, read, update, delete)
  - [ ] 2.1.3.4 Test burst suggestion logic
  - [ ] 2.1.3.5 Test result handling and metadata preservation
  - [ ] 2.1.3.6 Test error handling and edge cases
  - [ ] 2.1.3.7 Achieve 90%+ code coverage

#### 2.2 FactManagement Intent (16 hours)

- [ ] 2.2.1 Create fact_management.go with context and result structures
  - [ ] 2.2.1.1 Define FactManagementState enum (List, View, Edit, Review, Confirm)
  - [ ] 2.2.1.2 Define FactManagementContext struct with required fields
  - [ ] 2.2.1.3 Define FactManagementResult struct with output fields
  - [ ] 2.2.1.4 Add comprehensive documentation and examples

- [ ] 2.2.2 Create fact_management_intent.go with full implementation
  - [ ] 2.2.2.1 Implement Init() to load facts from service
  - [ ] 2.2.2.2 Implement Update() with state machine for all 5 states
  - [ ] 2.2.2.3 Implement View() with rendering for each state
  - [ ] 2.2.2.4 Implement Result() returning typed IntentResult
  - [ ] 2.2.2.5 Implement state transition handlers (handleListState, handleViewState, etc.)
  - [ ] 2.2.2.6 Implement helper methods for fact operations (create, edit, delete, search, filter)

- [ ] 2.2.3 Create fact_management_test.go with comprehensive tests
  - [ ] 2.2.3.1 Test state machine transitions (List → View → Edit → Review → Confirm)
  - [ ] 2.2.3.2 Test view rendering for each state (12+ tests)
  - [ ] 2.2.3.3 Test fact CRUD operations
  - [ ] 2.2.3.4 Test fact search and filtering
  - [ ] 2.2.3.5 Test fact quality assessment
  - [ ] 2.2.3.6 Test result handling and metadata preservation
  - [ ] 2.2.3.7 Test error handling and edge cases
  - [ ] 2.2.3.8 Achieve 90%+ code coverage

#### 2.3 ImportWizard Intent (12 hours)

- [ ] 2.3.1 Create import_wizard.go with context and result structures
  - [ ] 2.3.1.1 Define ImportWizardState enum (SelectFile, Review, Progress, Complete, Confirm)
  - [ ] 2.3.1.2 Define ImportWizardContext struct with required fields
  - [ ] 2.3.1.3 Define ImportWizardResult struct with output fields
  - [ ] 2.3.1.4 Add comprehensive documentation and examples

- [ ] 2.3.2 Create import_wizard_intent.go with full implementation
  - [ ] 2.3.2.1 Implement Init() for file selection UI
  - [ ] 2.3.2.2 Implement Update() with state machine for all 5 states
  - [ ] 2.3.2.3 Implement View() with rendering for each state
  - [ ] 2.3.2.4 Implement Result() returning typed IntentResult
  - [ ] 2.3.2.5 Implement file selection and validation
  - [ ] 2.3.2.6 Implement import progress tracking and reporting

- [ ] 2.3.3 Create import_wizard_test.go with comprehensive tests
  - [ ] 2.3.3.1 Test state machine transitions
  - [ ] 2.3.3.2 Test file selection and validation (8+ tests)
  - [ ] 2.3.3.3 Test import progress tracking
  - [ ] 2.3.3.4 Test error handling for invalid files
  - [ ] 2.3.3.5 Test result handling and metadata preservation
  - [ ] 2.3.3.6 Achieve 90%+ code coverage

#### 2.4 MetadataEditor Intent (10 hours)

- [ ] 2.4.1 Create metadata_editor.go with context and result structures
  - [ ] 2.4.1.1 Define MetadataEditorState enum (Review, Edit, Confirm)
  - [ ] 2.4.1.2 Define MetadataEditorContext struct with required fields
  - [ ] 2.4.1.3 Define MetadataEditorResult struct with output fields
  - [ ] 2.4.1.4 Add comprehensive documentation and examples

- [ ] 2.4.2 Create metadata_editor_intent.go with full implementation
  - [ ] 2.4.2.1 Implement Init() to load metadata for editing
  - [ ] 2.4.2.2 Implement Update() with state machine for all 3 states
  - [ ] 2.4.2.3 Implement View() with rendering for each state
  - [ ] 2.4.2.4 Implement Result() returning typed IntentResult
  - [ ] 2.4.2.5 Implement metadata editing with validation
  - [ ] 2.4.2.6 Implement change tracking for audit purposes

- [ ] 2.4.3 Create metadata_editor_test.go with comprehensive tests
  - [ ] 2.4.3.1 Test state machine transitions (Review → Edit → Confirm)
  - [ ] 2.4.3.2 Test metadata editing and validation (8+ tests)
  - [ ] 2.4.3.3 Test change tracking
  - [ ] 2.4.3.4 Test result handling and metadata preservation
  - [ ] 2.4.3.5 Achieve 90%+ code coverage

#### 2.5 BulkOperations Intent (10 hours)

- [ ] 2.5.1 Create bulk_operations.go with context and result structures
  - [ ] 2.5.1.1 Define BulkOperationsState enum (SelectOp, Configure, Execute, Confirm)
  - [ ] 2.5.1.2 Define BulkOperationsContext struct with required fields
  - [ ] 2.5.1.3 Define BulkOperationsResult struct with output fields
  - [ ] 2.5.1.4 Add comprehensive documentation and examples

- [ ] 2.5.2 Create bulk_operations_intent.go with full implementation
  - [ ] 2.5.2.1 Implement Init() to show available bulk operations
  - [ ] 2.5.2.2 Implement Update() with state machine for all 4 states
  - [ ] 2.5.2.3 Implement View() with rendering for each state
  - [ ] 2.5.2.4 Implement Result() returning typed IntentResult
  - [ ] 2.5.2.5 Implement operation selection and configuration
  - [ ] 2.5.2.6 Implement bulk operation execution with progress

- [ ] 2.5.3 Create bulk_operations_test.go with comprehensive tests
  - [ ] 2.5.3.1 Test state machine transitions
  - [ ] 2.5.3.2 Test operation selection and configuration (8+ tests)
  - [ ] 2.5.3.3 Test bulk operation execution
  - [ ] 2.5.3.4 Test progress tracking
  - [ ] 2.5.3.5 Test result handling and metadata preservation
  - [ ] 2.5.3.6 Achieve 90%+ code coverage

---

### 3.0 Rebuild Root Application (app.go) (14 hours)

- [ ] 3.1 Audit current app.go implementation
  - [ ] 3.1.1 Document all 31 screen constants
  - [ ] 3.1.2 Document all 37+ model fields
  - [ ] 3.1.3 Document all message types in Update()
  - [ ] 3.1.4 Document all screen-specific logic
  - [ ] 3.1.5 Identify critical business logic to preserve

- [ ] 3.2 Design new minimal app.go structure
  - [ ] 3.2.1 Define root model with only 5-7 core fields (services, router, state, width, height)
  - [ ] 3.2.2 Design message routing strategy (intent mode vs. home screen)
  - [ ] 3.2.3 Plan menu selection and intent activation flow
  - [ ] 3.2.4 Define back navigation and context preservation
  - [ ] 3.2.5 Document global shortcuts (Ctrl+C, ?, Esc, Home)

- [ ] 3.3 Create new app.go implementation (~200 lines)
  - [ ] 3.3.1 Create NewModel() factory function with minimal initialization
  - [ ] 3.3.2 Implement Init() method for root model setup
  - [ ] 3.3.3 Implement Update() method with intent routing (~30 lines)
  - [ ] 3.3.4 Implement View() method with simple screen selection (~20 lines)
  - [ ] 3.3.5 Implement menu selection handling
  - [ ] 3.3.6 Implement intent activation and result handling
  - [ ] 3.3.7 Implement back navigation logic
  - [ ] 3.3.8 Implement global shortcuts (quit, help, home)

- [ ] 3.4 Register all 10 intents with router
  - [ ] 3.4.1 Register CaptureEvent intent with factory
  - [ ] 3.4.2 Register BrowseTimeline intent with factory
  - [ ] 3.4.3 Register GenerateCV intent with factory
  - [ ] 3.4.4 Register ExportArtifact intent with factory
  - [ ] 3.4.5 Register ConfigureSystem intent with factory
  - [ ] 3.4.6 Register BurstManagement intent with factory
  - [ ] 3.4.7 Register FactManagement intent with factory
  - [ ] 3.4.8 Register ImportWizard intent with factory
  - [ ] 3.4.9 Register MetadataEditor intent with factory
  - [ ] 3.4.10 Register BulkOperations intent with factory

- [ ] 3.5 Implement menu and navigation
  - [ ] 3.5.1 Create menu with all 10 intent options
  - [ ] 3.5.2 Implement menu selection keyboard handling
  - [ ] 3.5.3 Implement intent activation from menu
  - [ ] 3.5.4 Implement back navigation to menu
  - [ ] 3.5.5 Implement context restoration after navigation

- [ ] 3.6 Remove all legacy code
  - [ ] 3.6.1 Remove all 31 screen constants
  - [ ] 3.6.2 Remove all 37+ screen-specific model fields
  - [ ] 3.6.3 Remove all screen-specific Update() logic
  - [ ] 3.6.4 Remove all screen-specific View() logic
  - [ ] 3.6.5 Remove all screen-specific helper methods
  - [ ] 3.6.6 Remove all breadcrumb and workflow management code

- [ ] 3.7 Verify app.go compilation and structure
  - [ ] 3.7.1 Ensure code compiles without errors
  - [ ] 3.7.2 Verify app.go is under 250 lines
  - [ ] 3.7.3 Verify Update() is under 50 lines
  - [ ] 3.7.4 Verify View() is under 30 lines
  - [ ] 3.7.5 Verify all model fields are documented

---

### 4.0 Comprehensive Testing & Quality Assurance (34 hours)

- [ ] 4.1 Unit tests for new intents (20 hours)
  - [ ] 4.1.1 Run BurstManagement tests and verify 30+ tests pass
  - [ ] 4.1.2 Run FactManagement tests and verify 40+ tests pass
  - [ ] 4.1.3 Run ImportWizard tests and verify 25+ tests pass
  - [ ] 4.1.4 Run MetadataEditor tests and verify 20+ tests pass
  - [ ] 4.1.5 Run BulkOperations tests and verify 20+ tests pass
  - [ ] 4.1.6 Verify all new intent tests achieve 90%+ coverage
  - [ ] 4.1.7 Run race detector on all new intent tests
  - [ ] 4.1.8 Fix any race conditions or coverage gaps

- [ ] 4.2 Integration tests (8 hours)
  - [ ] 4.2.1 Test all intents activate from menu correctly
  - [ ] 4.2.2 Test navigation between intents works properly
  - [ ] 4.2.3 Test back navigation preserves context
  - [ ] 4.2.4 Test result handling from each intent
  - [ ] 4.2.5 Test global shortcuts (quit, help, home)
  - [ ] 4.2.6 Test menu selection keyboard handling
  - [ ] 4.2.7 Verify 100% of workflows are completable

- [ ] 4.3 End-to-end workflow testing (6 hours)
  - [ ] 4.3.1 Test complete capture event workflow
  - [ ] 4.3.2 Test complete burst management workflow
  - [ ] 4.3.3 Test complete fact management workflow
  - [ ] 4.3.4 Test complete import wizard workflow
  - [ ] 4.3.5 Test complete CV generation and export workflow
  - [ ] 4.3.6 Test complete system configuration workflow
  - [ ] 4.3.7 Verify all data operations work correctly

- [ ] 4.4 Performance benchmarking (4 hours)
  - [ ] 4.4.1 Run intent Init benchmarks (target: < 50ms)
  - [ ] 4.4.2 Run view render benchmarks (target: < 100ms)
  - [ ] 4.4.3 Run state transition benchmarks (target: < 10ms)
  - [ ] 4.4.4 Run router operation benchmarks (target: < 1ms)
  - [ ] 4.4.5 Run full test suite with race detector (target: < 5s)
  - [ ] 4.4.6 Document any performance issues and optimizations

- [ ] 4.5 Code quality verification (4 hours)
  - [ ] 4.5.1 Run gofmt on all modified files
  - [ ] 4.5.2 Run golangci-lint on all modified files
  - [ ] 4.5.3 Verify zero lint issues
  - [ ] 4.5.4 Run full test suite: `go test -v ./...`
  - [ ] 4.5.5 Verify 164+ tests passing
  - [ ] 4.5.6 Verify 87%+ overall code coverage
  - [ ] 4.5.7 Verify zero race conditions detected

- [ ] 4.6 Documentation and cleanup (4 hours)
  - [ ] 4.6.1 Update README with new intent architecture overview
  - [ ] 4.6.2 Update TUI_DEVELOPER_GUIDE.md with new patterns
  - [ ] 4.6.3 Create phase completion report documenting changes
  - [ ] 4.6.4 Document any breaking changes or migration notes
  - [ ] 4.6.5 Update AGENTS.md with implementation summary
  - [ ] 4.6.6 Clean up any temporary files or branches

---

## Implementation Notes

### State Machine Pattern Template

Each intent should follow this pattern:

```go
type IntentState string

const (
    StateInitial IntentState = "initial"
    StateWorking IntentState = "working"
    StateFinal   IntentState = "final"
)

func (i *IntentModel) Update(msg tea.Msg) tea.Cmd {
    switch i.state {
    case StateInitial:
        return i.handleInitial(msg)
    case StateWorking:
        return i.handleWorking(msg)
    case StateFinal:
        return i.handleFinal(msg)
    }
    return nil
}

func (i *IntentModel) View() string {
    switch i.state {
    case StateInitial:
        return i.viewInitial()
    case StateWorking:
        return i.viewWorking()
    case StateFinal:
        return i.viewFinal()
    }
    return ""
}
```

### Testing Pattern Template

Each intent test should include:

```go
var _ = Describe("YourIntent", func() {
    var (
        intent   *YourIntentModel
        ctx      context.Context
        service  *MockService
    )

    BeforeEach(func() {
        service = NewMockService()
        intent = NewYourIntent(service)
        ctx = context.Background()
    })

    Describe("State Transitions", func() {
        It("should transition from initial to working state", func() {
            // Test implementation
        })
    })

    Describe("View Rendering", func() {
        It("should render initial state view", func() {
            // Test implementation
        })
    })

    Describe("Result Handling", func() {
        It("should return typed IntentResult", func() {
            // Test implementation
        })
    })
})
```

### Component Reuse

Leverage existing components in `internal/cli/components/`:
- Card - Display data in card format
- List - Render scrollable lists
- Form - Handle form input
- Modal - Display modal dialogs
- Progress - Show progress indicators
- Menu - Display menu options

### Styling

Use centralized styling from `internal/cli/styles/styles.go`:
- Consistent color scheme across all intents
- Uniform border styles
- Standardized spacing and padding
- Responsive to terminal width changes

---

## Success Criteria Checklist

### Code Reduction
- [ ] `app.go` reduced to < 250 lines (from 1,766)
- [ ] `Update()` method < 50 lines (from 906)
- [ ] `View()` method < 30 lines (from 180)
- [ ] Model fields < 10 (from 37+)
- [ ] All 31 screen constants removed
- [ ] All 30+ legacy model files deleted

### Intent Implementation
- [ ] 5 new intents fully implemented
- [ ] 10 total intents registered with router
- [ ] All intents follow state machine pattern
- [ ] All intents use type-safe `IntentResult[T]`
- [ ] All intents support back navigation

### Testing & Quality
- [ ] 164+ tests passing (100% pass rate)
- [ ] 87%+ code coverage
- [ ] 0 race conditions detected
- [ ] All lint checks passing
- [ ] All performance benchmarks met
- [ ] Code formatted with gofmt

### Feature Preservation
- [ ] All 29+ legacy screen features preserved
- [ ] All CRUD operations working
- [ ] All menu items functional
- [ ] All workflows completable
- [ ] All data operations unchanged
- [ ] All error handling preserved

### Integration
- [ ] All intents activate from menu
- [ ] Navigation between intents works
- [ ] Results properly handled
- [ ] Context preserved on back navigation
- [ ] Services properly initialized
- [ ] Global state properly managed

### Documentation
- [ ] Implementation guide updated
- [ ] Intent patterns documented
- [ ] Code comments clear and complete
- [ ] README reflects new architecture
- [ ] Phase completion report written
- [ ] Known issues documented

---

**Document Version**: 1.0
- **Process Guide**: `/docs/rules/master-task-prompt.md`

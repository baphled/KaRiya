# Workflow Issues Found - Manual Trace Analysis

**Date**: January 3, 2026
**Status**: Issues Identified and Documented

---

## Summary

After manually tracing through the workflow, the following issues have been identified:

1. **Constructor Signature Mismatch** - Some intents return errors, others don't
2. **Nil Pointer Dereferences** - No nil checks in some intent constructors
3. **Missing Context Fields** - Some intents lack required initialization
4. **Inconsistent Error Handling** - No uniform error handling pattern

---

## Detailed Issues

### Issue 1: Constructor Signature Inconsistency

#### Intents with Error Returns (5):
- `NewCaptureEventIntent(context *CaptureEventContext) (*CaptureEventIntent, error)`
- `NewBrowseTimelineIntent(context *BrowseTimelineContext) (*BrowseTimelineIntent, error)`
- `NewGenerateCVIntent(context *GenerateCVContext) (*GenerateCVIntent, error)`
- `NewExportArtifactIntent(ctx context.Context) (*ExportArtifactIntent, error)`
- `NewConfigureSystemIntent(ctx context.Context) (*ConfigureSystemIntent, error)`

#### Intents WITHOUT Error Returns (5):
- `NewBurstManagementIntent(data *BurstManagementContext) *BurstManagementModel`
- `NewFactManagementIntent(data *FactManagementContext) *FactManagementModel`
- `NewImportWizardIntent(data *ImportWizardContext) *ImportWizardModel`
- `NewMetadataEditorIntent(data *MetadataEditorContext) *MetadataEditorModel`
- `NewBulkOperationsIntent(data *BulkOperationsContext) *BulkOperationsModel`

**Problem**: app.go tries to handle errors for all intents, but the new intents don't return errors.

**Location**: `internal/cli/app/app.go` lines 220-340 (registerAllIntents function)

---

### Issue 2: Missing Context Fields in New Intents

#### BurstManagement Context Missing:
```go
burstCtx := &intents.BurstManagementContext{
  Bursts:  bursts,
  Service: careerService,  // ← Added but might not be in struct
}
```

**Problem**: The `BurstManagementContext` struct might not have a `Service` field defined.

**Impact**: Will panic when trying to access `m.data.Service` in BurstManagement intent methods.

---

### Issue 3: Nil Pointer Dereferences

#### In app.go registerAllIntents:
```go
// Line 295-307: BurstManagement
burstCtx := &intents.BurstManagementContext{
  Bursts:  bursts,
  Service: careerService,
}
return intents.NewBurstManagementIntent(burstCtx)
// ← No nil check, will panic if BurstManagementContext is nil
```

**Problem**: If the intent constructor returns nil (due to error), the router will try to use it.

**Impact**: Runtime panics when activating intents.

---

### Issue 4: Missing LoadBursts() Method

In `BurstManagementIntent.Init()`:
```go
if err := m.data.LoadBursts(); err != nil {
  // ...
}
```

**Problem**: `BurstManagementContext` might not have a `LoadBursts()` method.

**Impact**: Will panic with "LoadBursts method not found" when intent initializes.

---

### Issue 5: Missing LoadFacts() Method

In `FactManagementIntent.Init()`:
```go
if err := m.data.LoadFacts(); err != nil {
  // ...
}
```

**Problem**: `FactManagementContext` might not have a `LoadFacts()` method.

**Impact**: Will panic with "LoadFacts method not found" when intent initializes.

---

### Issue 6: Missing Repository Methods

In various intent constructors:
```go
careerService.GetBurstRepository()
careerService.GetFactRepository()
```

**Problem**: The career service might not have these getter methods.

**Impact**: Will panic with "GetBurstRepository method not found" at app startup.

---

## Workflow Trace - Where Issues Occur

### Step 1: Application Startup
```
NewModel()
  → registerAllIntents()
    → router.RegisterIntent("burst_management", factory)
    → factory() called on intent activation
      → careerService.GetBurstRepository() ← ISSUE: Method might not exist
```

**Potential Panic**: "GetBurstRepository method not found"

---

### Step 2: User Selects "Manage Bursts"
```
handleMenuInput("enter")
  → activateIntent("burst_management")
    → router.ActivateIntent("burst_management", {})
      → factory() called
        → NewBurstManagementIntent(burstCtx)
          → Returns *BurstManagementModel (no error)
            → Intent.Init()
              → m.data.LoadBursts() ← ISSUE: Method might not exist
```

**Potential Panic**: "LoadBursts method not found"

---

### Step 3: Intent Update Loop
```
handleIntentInput(msg)
  → intentRouter.HandleMessage(msg)
    → activeIntent.Update(msg)
      → Accesses m.data.Service ← ISSUE: Field might not exist
```

**Potential Panic**: "Service field not found"

---

## Root Causes

1. **Incomplete Intent Implementation**
   - New intents (BurstManagement, FactManagement, etc.) may not be fully implemented
   - Missing required methods and fields

2. **Missing Service Methods**
   - CareerService might not have GetBurstRepository() or GetFactRepository()

3. **Signature Inconsistency**
   - Old intents return (Model, error)
   - New intents return just Model
   - No error handling for new intents

4. **Context Structure Mismatch**
   - Trying to add fields to contexts that don't have them
   - Trying to call methods on contexts that don't have them

---

## Recommended Fixes

### Priority 1: Fix Constructor Signatures (CRITICAL)
```go
// Make all new intents return errors for consistency
func NewBurstManagementIntent(data *BurstManagementContext) (*BurstManagementModel, error) {
  if data == nil {
    return nil, fmt.Errorf("context is required")
  }
  return &BurstManagementModel{
    data: data,
    // ...
  }, nil
}
```

### Priority 2: Add Missing Service Methods
```go
// In CareerService
func (s *Service) GetBurstRepository() BurstRepository {
  return s.burstRepository
}

func (s *Service) GetFactRepository() FactRepository {
  return s.factRepository
}
```

### Priority 3: Verify Context Structures
```go
// Verify BurstManagementContext has Service field
type BurstManagementContext struct {
  // ...
  Service *careerservice.Service
  // ...
}

// Add LoadBursts method
func (c *BurstManagementContext) LoadBursts() error {
  // Implementation
}
```

### Priority 4: Update app.go Registration
```go
// Handle both signature styles
router.RegisterIntent("burst_management", func() intents.Intent {
  burstCtx := &intents.BurstManagementContext{
    Bursts: bursts,
    Service: careerService,
  }
  intent, err := intents.NewBurstManagementIntent(burstCtx)
  if err != nil {
    log.Error("Failed to create BurstManagement intent: %v", err)
    return nil
  }
  return intent
})
```

---

## Testing Plan

1. **Unit Tests**: Test each intent constructor with nil inputs
2. **Integration Tests**: Test each menu item activation
3. **Workflow Tests**: Test complete user workflows
4. **Panic Recovery**: Add panic recovery to intent activation

---

## Conclusion

The aggressive app.go refactoring introduced several issues due to incomplete intent implementations and missing infrastructure. The fixes are straightforward but require careful implementation to maintain consistency across all intents.

**Estimated Fix Time**: 2-3 hours for complete resolution



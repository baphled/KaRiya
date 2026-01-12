# Escape Key Test Coverage - Quick Summary

## Current Status: 5/10 Complete (50%)

### ✅ Intents with Complete Escape Tests

1. **capture_event** - Including context-aware edit tests ✨
2. **browse_timeline**
3. **configure_system**
4. **export_artifact**
5. **generate_cv**

### ❌ Missing Escape Tests

#### High Priority (Has Edit Context)
- **burst_management** - Needs context-aware tests for EditingBurst
- **fact_management** - Needs context-aware tests for EditingFact

#### Medium Priority
- **import_wizard** - Async operations, needs careful testing
- **bulk_operations** - Async operations, needs careful testing

#### Low Priority
- **metadata_editor** - Simple states, less critical

---

## Key Concepts

### Context-Aware Navigation

When an intent supports both **creating new items** and **editing existing items**, escape behavior must be context-aware:

**New Item**:
```
Form → Esc → Previous Screen (within intent)
```

**Edit Existing Item**:
```
Form → Esc → Cancel Intent (return to calling intent)
```

**Implementation Pattern**:
```go
case KeyBack:
    if i.context.EditingItem != nil {
        // Editing - cancel and return to caller
        i.setCancelled()
        return nil
    }
    // New item - go back to previous state
    i.state.currentState = PreviousState
    return nil
```

### Test Pattern

**For each intent, test**:
1. Root state → Esc → Cancel intent
2. Intermediate states → Esc → Previous state
3. **If has edit context**: New vs Edit behavior
4. Final states → Esc → Cancel/Complete

**Example Test Structure**:
```go
Describe("Form State", func() {
    Context("when creating new item", func() {
        It("should go back to list", func() {
            // EditingItem = nil
            // Esc → Previous state
        })
    })
    
    Context("when editing existing item", func() {
        It("should cancel intent", func() {
            // EditingItem != nil
            // Esc → Cancel entire intent
        })
    })
})
```

---

## Quick Commands

```bash
# Check current test coverage
ls internal/cli/intents/*_escape_test.go

# Run all escape tests
ginkgo --focus="Escape Key" ./internal/cli/intents/

# Run specific intent escape tests
ginkgo --focus="burst_management.*Escape" ./internal/cli/intents/

# Create new escape test file
cp internal/cli/intents/capture_event_escape_test.go \
   internal/cli/intents/burst_management_escape_test.go
```

---

## Next Steps

1. **Review** full test plan: `bugs/escape-test-plan.md`
2. **Start with** burst_management (highest priority)
3. **Use template** from test plan document
4. **Verify** implementation handles context correctly
5. **Test** manually before committing

---

**See Also**: 
- Full test plan: `bugs/escape-test-plan.md`
- Implementation example: `internal/cli/intents/capture_event_intent.go:310-319`
- Test example: `internal/cli/intents/capture_event_escape_test.go:54-92`

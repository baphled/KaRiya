# Migration Task: Move Modals from intents/modals.go

## Overview

The file `internal/cli/intents/modals.go` contains modal structs that violate the architecture rules. This task tracks the migration of these modals to their correct locations.

## Current Violations

| Check | Violation | File |
|-------|-----------|------|
| #24 | Modal structs in intents/ package | `intents/modals.go` |
| #25 | Direct huh import in intents | `intents/modals.go:9` |
| #28 | Modal structs in wrong location | `intents/modals.go` |

## Modals to Migrate

### 1. EditBurstModal

| Attribute | Current | Target |
|-----------|---------|--------|
| File | `intents/modals.go` | `screens/burst_management/modals/edit_modal.go` |
| Package | `intents` | `modals` |
| Struct | `EditBurstModal` | `EditModal` |
| Constructor | `NewEditBurstModal()` | `NewEditModal()` |

### 2. EditFactModal

| Attribute | Current | Target |
|-----------|---------|--------|
| File | `intents/modals.go` | `screens/fact_management/modals/edit_modal.go` |
| Package | `intents` | `modals` |
| Struct | `EditFactModal` | `EditModal` |
| Constructor | `NewEditFactModal()` | `NewEditModal()` |

## Migration Steps

### Phase 1: Create Target Directories

```bash
mkdir -p internal/cli/screens/burst_management/modals
mkdir -p internal/cli/screens/fact_management/modals
```

### Phase 2: Migrate EditBurstModal

1. Create new file:
   ```go
   // File: internal/cli/screens/burst_management/modals/edit_modal.go
   package modals
   
   import (
       "github.com/baphled/kariya/internal/cli/forms"
       "github.com/baphled/kariya/internal/domain/career"
   )
   
   // EditModal handles inline editing of burst details.
   type EditModal struct {
       // ... fields from EditBurstModal
   }
   
   func NewEditModal(burst *career.Burst) *EditModal {
       // ... implementation
   }
   ```

2. Refactor to use `forms/` package instead of direct `huh` import

3. Update references in `intents/burst_management/`:
   ```go
   // OLD
   import "github.com/baphled/kariya/internal/cli/intents"
   editModal *intents.EditBurstModal
   
   // NEW
   import "github.com/baphled/kariya/internal/cli/screens/burst_management/modals"
   editModal *modals.EditModal
   ```

4. Update constructor calls:
   ```go
   // OLD
   i.editModal = intents.NewEditBurstModal(burst)
   
   // NEW
   i.editModal = modals.NewEditModal(burst)
   ```

### Phase 3: Migrate EditFactModal

1. Create new file:
   ```go
   // File: internal/cli/screens/fact_management/modals/edit_modal.go
   package modals
   
   import (
       "github.com/baphled/kariya/internal/cli/forms"
       "github.com/baphled/kariya/internal/domain/career"
   )
   
   // EditModal handles inline editing of fact details.
   type EditModal struct {
       // ... fields from EditFactModal
   }
   
   func NewEditModal(fact *career.Fact) *EditModal {
       // ... implementation
   }
   ```

2. Refactor to use `forms/` package instead of direct `huh` import

3. Update references in `intents/fact_management/`:
   ```go
   // OLD
   import "github.com/baphled/kariya/internal/cli/intents"
   editModal *intents.EditFactModal
   
   // NEW
   import "github.com/baphled/kariya/internal/cli/screens/fact_management/modals"
   editModal *modals.EditModal
   ```

### Phase 4: Cleanup

1. Remove `intents/modals.go` after all migrations complete
2. Verify no remaining references to old modal types
3. Run `make check-intent-architecture` to confirm violations resolved

### Phase 5: Additional Fact Management Work

The fact_management intent also needs screens extracted (Check #26, #27):

1. Create screens:
   ```
   screens/fact_management/
   ├── list.go       # FactListScreen (from helpers.go getStateContent)
   ├── detail.go     # FactDetailScreen (from helpers.go getViewFactContent)
   ├── form.go       # FactFormScreen (from helpers.go getEditorContent)
   ├── delete.go     # FactDeleteScreen (from helpers.go getDeleteConfirmContent)
   └── modals/
       └── edit_modal.go  # EditModal (migrated above)
   ```

2. Update `fact_management/helpers.go` to delegate to screens
3. Update `fact_management/intent.go` to orchestrate screens

## Acceptance Criteria

- [ ] `internal/cli/intents/modals.go` deleted
- [ ] `EditBurstModal` migrated to `screens/burst_management/modals/edit_modal.go`
- [ ] `EditFactModal` migrated to `screens/fact_management/modals/edit_modal.go`
- [ ] No direct `huh` imports in migrated modals (use `forms/` package)
- [ ] All references updated in intents
- [ ] `make check-intent-architecture` passes with no violations
- [ ] All tests pass

## Dependencies

This migration should be completed BEFORE:
- Any new features added to burst_management intent
- Any new features added to fact_management intent
- Enabling checks #24, #25, #27, #28 as blocking in CI

## Related Documentation

- [MODAL_NAMING.md](../conventions/MODAL_NAMING.md) - Modal naming conventions
- [SCREEN_NAMING.md](../conventions/SCREEN_NAMING.md) - Screen naming conventions
- [INTENT_ARCHITECTURE_GUIDE.md](../INTENT_ARCHITECTURE_GUIDE.md) - Architecture rules

## Priority

**HIGH** - Blocks enforcement of new architecture checks.

## Estimated Effort

- EditBurstModal migration: 1-2 hours
- EditFactModal migration: 1-2 hours
- Fact management screens extraction: 4-6 hours
- Testing and verification: 1-2 hours

**Total: 7-12 hours**

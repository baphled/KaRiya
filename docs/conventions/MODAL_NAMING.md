# Modal Naming Convention

This document defines the standard naming conventions for modals in the KaRiya codebase.

## Overview

Modals are overlay UI components that appear above the main screen content. They are used for confirmations, forms, filters, and other focused interactions. Modals must be placed in approved locations and follow strict naming conventions.

## Allowed Locations

Modals may ONLY be defined in the following locations:

| Location | Purpose | When to Use |
|----------|---------|-------------|
| `internal/cli/uikit/feedback/` | Reusable/generic modals | Confirm, Error, Loading, Success, Warning |
| `internal/cli/screens/{feature}/modals/` | Feature-specific modals | Filter, Sort, Search, Edit for specific feature |
| `internal/cli/components/` | Legacy (deprecated) | Do NOT add new modals here |

### Forbidden Locations

Modals MUST NOT be defined in:

| Location | Reason |
|----------|--------|
| `internal/cli/intents/` | Violates layer separation |
| `internal/cli/models/` | Package has been deleted |
| Any other location | Not in approved list |

## Directory Structure

### Feature-Specific Modals

```
internal/cli/screens/{feature}/modals/
├── edit_modal.go      # EditModal
├── filter_modal.go    # FilterModal
├── sort_modal.go      # SortModal
├── search_modal.go    # SearchModal
├── delete_modal.go    # DeleteModal
├── detail_modal.go    # DetailModal
└── helpers.go         # Shared modal utilities
```

### UIKit Feedback Modals

```
internal/cli/uikit/feedback/
├── confirm_modal.go   # ConfirmModal
├── error_modal.go     # (via NewErrorModal)
├── help_modal.go      # HelpModal
├── info_modal.go      # InfoModal
├── detail_modal.go    # DetailModal
└── modal.go           # Base modal types
```

## Naming Rules

### File Naming

| Rule | Convention | Example |
|------|------------|---------|
| Pattern | `{action}_modal.go` | `edit_modal.go`, `filter_modal.go` |
| Lowercase | All lowercase with underscores | `quick_add_modal.go` |
| Suffix | Always ends with `_modal.go` | `delete_modal.go` |

### Struct Naming

| Rule | Convention | Example |
|------|------------|---------|
| Pattern | `{Action}Modal` | `EditModal`, `FilterModal` |
| No feature prefix | Package provides context | `EditModal` NOT `FactEditModal` |
| Suffix | Always ends with `Modal` | `FilterModal` NOT `Filter` |
| PascalCase | Use PascalCase | `QuickAddModal` |

### Constructor Naming

| Rule | Convention | Example |
|------|------------|---------|
| Pattern | `New{StructName}()` | `NewEditModal()` |
| Match struct | Constructor matches struct name | `NewFilterModal()` for `FilterModal` |

## Modal Types

| Type | Purpose | File Name | Struct Example |
|------|---------|-----------|----------------|
| `Edit` | Edit existing item | `edit_modal.go` | `EditModal` |
| `Filter` | Filter list items | `filter_modal.go` | `FilterModal` |
| `Sort` | Sort options | `sort_modal.go` | `SortModal` |
| `Search` | Text search | `search_modal.go` | `SearchModal` |
| `Confirm` | Confirmation dialog | `confirm_modal.go` | `ConfirmModal` |
| `Delete` | Delete confirmation | `delete_modal.go` | `DeleteModal` |
| `Detail` | View item details | `detail_modal.go` | `DetailModal` |
| `QuickAdd` | Quick add form | `quick_add_modal.go` | `QuickAddModal` |
| `Help` | Help overlay | `help_modal.go` | `HelpModal` |

## One Struct Per File Rule

Each modal file must contain exactly ONE modal struct (except `helpers.go`).

```go
// CORRECT: One struct per file
// File: screens/timeline/modals/filter_modal.go
type FilterModal struct {
    // ...
}

// WRONG: Multiple structs in one file
// File: screens/timeline/modals/detail_modal.go
type EventDetailModal struct { ... }
type SkillsDetailModal struct { ... }  // VIOLATION - split into separate files
```

## Examples

### Feature-Specific Modal

```go
// File: screens/fact_management/modals/edit_modal.go
package modals

import (
    "github.com/baphled/kariya/internal/cli/forms"
    "github.com/baphled/kariya/internal/domain/career"
)

// EditModal handles editing of fact details.
type EditModal struct {
    visible  bool
    form     forms.Form
    original *career.Fact
    modified *career.Fact
}

// NewEditModal creates a new fact editing modal.
func NewEditModal(fact *career.Fact) *EditModal {
    return &EditModal{
        original: fact,
        modified: copyFact(fact),
        form:     forms.NewFactEditForm(fact),
    }
}

func (m *EditModal) IsVisible() bool {
    return m.visible
}

func (m *EditModal) Show() {
    m.visible = true
}

func (m *EditModal) Hide() {
    m.visible = false
}
```

### Using UIKit Feedback Modals

```go
// In intent code
import "github.com/baphled/kariya/internal/cli/uikit/feedback"

// For confirmations, use feedback package
deleteModal := feedback.NewConfirmModal(
    "Delete Fact",
    "Are you sure you want to delete this fact?",
).WithVariant(feedback.ConfirmDestructive)

// For errors
errorModal := feedback.NewErrorModal("Operation failed", err.Error())

// For success
successModal := feedback.NewSuccessModal("Fact saved successfully")
```

## huh Import Restriction

The `huh` library (form library) should NOT be directly imported in modal files. Use the `forms/` package instead.

```go
// WRONG: Direct huh import
import "github.com/charmbracelet/huh"

type EditModal struct {
    form *huh.Form  // VIOLATION
}

// CORRECT: Use forms package
import "github.com/baphled/kariya/internal/cli/forms"

type EditModal struct {
    form forms.Form  // Uses forms package wrapper
}
```

## Enforcement

These conventions are enforced by multiple checks in `scripts/check-intent-architecture.sh`:

| Check | What It Catches |
|-------|-----------------|
| #24 | Modal structs in intents/ package |
| #25 | Direct huh import in wrong locations |
| #28 | Modal structs in wrong locations |
| #29 | Modal naming convention violations |

Run to verify:
```bash
make check-intent-architecture
```

## Migration Guide

If you have modals in the wrong location:

1. Create the target directory:
   ```bash
   mkdir -p internal/cli/screens/{feature}/modals/
   ```

2. Move the modal file:
   ```bash
   mv internal/cli/intents/modals.go internal/cli/screens/{feature}/modals/edit_modal.go
   ```

3. Update the package declaration:
   ```go
   package modals  // NOT package intents
   ```

4. Update imports in intents:
   ```go
   import "github.com/baphled/kariya/internal/cli/screens/{feature}/modals"
   ```

5. Update struct references:
   ```go
   // OLD
   editModal *intents.EditFactModal
   
   // NEW
   editModal *modals.EditModal
   ```

## Related Documentation

- [SCREEN_NAMING.md](SCREEN_NAMING.md) - Screen naming conventions
- [INTENT_ARCHITECTURE_GUIDE.md](../INTENT_ARCHITECTURE_GUIDE.md) - Overall architecture
- [FORMS_GUIDE.md](../FORMS_GUIDE.md) - Forms usage guide

# Screen Naming Convention

This document defines the standard naming conventions for screens in the KaRiya codebase.

## Overview

Screens are UI components that represent distinct views within an intent workflow. They are located in the `internal/cli/screens/` package and follow strict naming conventions to ensure consistency and maintainability.

## Directory Structure

```
internal/cli/screens/{feature}/
├── list.go          # {Entity}ListScreen
├── detail.go        # {Entity}DetailScreen
├── form.go          # {Entity}FormScreen
├── delete.go        # {Entity}DeleteScreen
├── select.go        # {Entity}SelectScreen
├── confirm.go       # {Entity}ConfirmScreen
├── preview.go       # {Entity}PreviewScreen
├── review.go        # {Entity}ReviewScreen
└── modals/          # Feature-specific modals (see MODAL_NAMING.md)
    ├── edit_modal.go
    ├── filter_modal.go
    └── helpers.go
```

## Naming Rules

### Package Naming

| Rule | Convention | Example |
|------|------------|---------|
| Directory name | `{feature_name}/` (underscore) | `fact_management/`, `burst_management/` |
| Package name | Same as directory | `package fact_management` |
| No hyphens | Use underscore, not hyphen | `fact_management/` NOT `fact-management/` |

### File Naming

| Rule | Convention | Example |
|------|------------|---------|
| Pattern | `{type}.go` | `list.go`, `detail.go`, `form.go` |
| Alternative | `{entity}_{type}.go` | `event_list.go`, `skill_form.go` |
| Lowercase | All lowercase with underscores | `delete_confirm.go` |
| No suffix | Do not add `_screen.go` suffix | `list.go` NOT `list_screen.go` |

### Struct Naming

| Rule | Convention | Example |
|------|------------|---------|
| Pattern | `{Entity}{Type}Screen` | `FactListScreen`, `SkillDetailScreen` |
| Entity | Short, singular entity name | `Fact`, `Skill`, `Burst`, `Event` |
| Type | Screen purpose | `List`, `Detail`, `Form`, `Delete`, `Select` |
| Suffix | Always ends with `Screen` | `FactListScreen` NOT `FactList` |
| PascalCase | Use PascalCase | `FactListScreen` NOT `Fact_List_Screen` |

### Constructor Naming

| Rule | Convention | Example |
|------|------------|---------|
| Pattern | `New{StructName}()` | `NewFactListScreen()` |
| Match struct | Constructor matches struct name | `NewSkillDetailScreen()` for `SkillDetailScreen` |

## Screen Types

| Type | Purpose | File Name | Struct Example |
|------|---------|-----------|----------------|
| `List` | Display list/table of items | `list.go` | `FactListScreen` |
| `Detail` | Display single item details | `detail.go` | `FactDetailScreen` |
| `Form` | Create or edit item | `form.go` | `FactFormScreen` |
| `Delete` | Delete confirmation | `delete.go` | `FactDeleteScreen` |
| `Select` | Selection picker | `select.go` | `FactSelectScreen` |
| `Confirm` | Generic confirmation | `confirm.go` | `FactConfirmScreen` |
| `Preview` | Preview before action | `preview.go` | `FactPreviewScreen` |
| `Review` | Review changes | `review.go` | `FactReviewScreen` |

## Examples

### Fact Management Screens

```
screens/fact_management/
├── list.go       # FactListScreen
├── detail.go     # FactDetailScreen
├── form.go       # FactFormScreen
├── delete.go     # FactDeleteScreen
└── modals/
    └── edit_modal.go  # EditModal
```

```go
// File: screens/fact_management/list.go
package fact_management

type FactListScreen struct {
    *base.BaseScreen
    table *behaviors.TableBehavior[*domain.Fact]
}

func NewFactListScreen(facts []*domain.Fact) *FactListScreen {
    // ...
}
```

### Timeline Screens (Reference Implementation)

```
screens/timeline/
├── event_list.go         # TimelineEventListScreen
├── event_detail.go       # TimelineEventDetailScreen
├── event_delete_confirm.go # EventDeleteConfirmScreen
└── modals/
    ├── filter_modal.go   # FilterModal
    ├── sort_modal.go     # SortModal
    ├── search_modal.go   # SearchModal
    ├── edit_modal.go     # EditModal
    └── helpers.go
```

## Enforcement

These conventions are enforced by Check #29 in `scripts/check-intent-architecture.sh`:

- Screen structs must end with `Screen`
- Package names must use underscore (not hyphen)
- Violations block commits

Run to verify:
```bash
make check-intent-architecture
```

## Related Documentation

- [MODAL_NAMING.md](MODAL_NAMING.md) - Modal naming conventions
- [INTENT_ARCHITECTURE_GUIDE.md](../INTENT_ARCHITECTURE_GUIDE.md) - Overall architecture
- [INTENT_DEVELOPMENT_CHECKLIST.md](../checklists/INTENT_DEVELOPMENT_CHECKLIST.md) - Development checklist

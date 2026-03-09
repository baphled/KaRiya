# Skill: kariya-generic-modal-intent-patterns

## What I do

I provide standard patterns for implementing generic modal interactions (Search, Sort, Filter) within KaRiya Intents. I focus on reusable logic where the primary difference is the mapping between domain entities and presenters.

## When to use me

- Adding Search/Sort/Filter to an Intent (e.g., Timeline, Skills, Facts)
- Implementing "Quick Add" or "Edit" modals using the `FormModalAdapter`
- Managing modal lifecycle using the `ModalRegistry`

## Core Patterns

1. **Adapter-based Integration**: Use `intents.NewFormModalAdapter` to wrap form-based modals. This normalizes the `Update` and `View` signatures for the intent.
2. **Registry Management**: Use `intents.NewModalRegistry()` in your intent to track which modal is currently visible and handle overlays.
3. **Domain Mapping**: Keep modal data structures focused on presentation; perform entity-to-form mapping in the intent handler.
4. **Action Handlers**: Decouple modal submission from visibility; the intent Update loop should handle the `Applied` result from the adapter.

## Implementation Checklist

- [ ] Initialize `i.modalRegistry = intents.NewModalRegistry()` in `NewIntent`
- [ ] Create modal instances (e.g., `i.searchModal = modals.NewSearchModal(...)`)
- [ ] Wrap modals in adapters: `i.modalRegistry.Register(modalKeySearch, intents.NewFormModalAdapter(...))`
- [ ] In `Update()`: Check `i.modalRegistry.HasVisibleModal()` before other key handling
- [ ] Handle `ModalUpdateResult.Applied` to trigger domain-specific updates (e.g., `i.RefreshData()`)
- [ ] Use `i.modalRegistry.RenderOverlay(baseView)` in `View()` to composite the final UI

## Boundary Rules

| Component | Responsibility | Mapping Location |
|-----------|----------------|------------------|
| **Entity** | Domain business rules | `internal/domain/` |
| **Form** | Presentation-only state | `internal/cli/screens/*/modals/` |
| **Intent** | Syncing Entity ↔ Form | `internal/cli/intents/*/handlers.go` |

## Forbidden Patterns

- ❌ **Intent logic in Modal**: Modals should only return their data; the intent decides how to use it.
- ❌ **Direct Modal Rendering**: Use `ModalRegistry.RenderOverlay` for consistent dimming and positioning.
- ❌ **Raw string keys for search**: Use typed constants (e.g., `FilterLayerSearch`) with `behaviors.FilterStack`.

## KB Reference

See ADRs in vault:
- `[[ADR 003: View Abstraction]]`
- `[[ADR 005: Callback Ownership Boundary]]`
- `[[ADR 006: Intent Architecture]]`

## Related Skills

- `architecture` - Layer boundaries
- `clean-code` - Generic patterns and naming
- `documentation-writing` - Clear implementation guides

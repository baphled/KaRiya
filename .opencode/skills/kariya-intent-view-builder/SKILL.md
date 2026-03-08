# Skill: kariya-intent-view-builder

## What I do

I provide a structured approach to building KaRiya Intents using the View-based abstraction pattern. I enforce clear boundaries between orchestration (Intent), presentation (Screen/View), and user feedback (Modals).

## When to use me

- Creating a new Intent workflow (e.g., `internal/cli/intents/myintent/`)
- Designing complex views with multi-level navigation
- Implementing state-driven modal overlays

## Core Principles

1. **Separation of Concerns**: Intents manage state and orchestration; Views/Screens manage rendering.
2. **Standardised Layout**: Use `intents.CreateStandardView` or `intents.CreateStandardViewWithBreadcrumbs` for consistent TUI structure.
3. **Implicit Feedback**: BaseIntent state (Error, Loading, Progress, Success) automatically triggers modal overlays in the standard view.
4. **Typed Action Vocabulary**: Define a typed `Action` enum and `Nav` struct near the view/screen to formalise intent-screen communication.

## Implementation Checklist

- [ ] Define `State` enums in `constants.go`
- [ ] Define `Action` enums and `Nav` result struct in `types.go` or near the screen
- [ ] Embed `*intents.BaseIntent` in your Intent struct
- [ ] Use `intents.NewModalRegistry()` for complex modal management
- [ ] Implement `View()` using `intents.CreateStandardView` or `intents.CreateStandardViewWithBreadcrumbs`
- [ ] Use `intents.MessageInterceptor` in `Update()` to protect global keys (Esc, ?, q)
- [ ] Delegate `Update()` to `activeScreen` and handle `screens.ScreenResult`

## Boundary Rules

| Component | Responsibility | Ownership |
|-----------|----------------|-----------|
| **Intent** | State Machine, Service Calls | Orchestration |
| **Screen** | UI Layout, Internal Selection | Presentation |
| **View/Nav** | Action Vocabulary (Typed) | Schema |
| **Modal** | Transient Interaction | Feedback |

## Forbidden Patterns

- ❌ **Screen importing Intent**: Screens must remain generic or communicate via `ScreenResult`.
- ❌ **Service logic in View**: Business rules and data fetching belong in the Intent or Service layer.
- ❌ **Raw map payloads**: Use typed `Nav` structs for `ResultData` in `NavigateResult`.
- ❌ **Direct Modal Rendering**: Use `BaseIntent` state or `ModalRegistry` to overlay modals onto the base view.

## KB Reference

See ADRs in vault:
- `[[ADR 003: View Abstraction]]`
- `[[ADR 004: Action Layer Pattern]]`
- `[[ADR 005: Callback Ownership Boundary]]`
- `[[ADR 006: Intent Architecture]]`

## Related Skills

- `architecture` - Layer separation rules
- `clean-code` - Naming and structure
- `obsidian-mermaid-expert` - Mapping state transitions

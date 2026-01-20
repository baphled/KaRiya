# KaRiya AI Agent Instructions

## Session Start (Required)

```bash
make session-start   # MUST run first - validates environment
```

## Critical Rules

1. **PRs target `next`** - Never `main`. Only `next→main` for releases.
2. **TDD** - Write test FIRST, then implementation (Red→Green→Refactor)
3. **Commits** - Use `make ai-commit FILE=<path>` only
4. **Compliance** - Run `make check-compliance` before AND after tasks
5. **One task** - One logical change per commit

## Essential Commands

| Command | Purpose |
|---------|---------|
| `make session-start` | Start every session (required) |
| `make session-end` | End session (cleanup) |
| `make pre-task` | Checklist before any task |
| `make check-compliance` | Validate before/after tasks |
| `make what-to-use NEED="x"` | Component lookup (table, form, modal...) |
| `make check-patterns` | Detect pattern violations |
| `make check-patterns-strict` | Strict pattern check (blocking) |
| `make ai-commit FILE=...` | Commit with AI attribution |
| `make pre-pr` | Validate before creating PR |

## TDD Workflow Commands

| Command | Purpose |
|---------|---------|
| `make tdd-red` | Start TDD: write failing test |
| `make tdd-green` | Make test pass with minimal code |
| `make tdd-refactor` | Improve code quality |
| `make tdd-document` | Finalize and commit |

## Task Management

| Command | Purpose |
|---------|---------|
| `make new-feature TASK="x"` | Create new feature task |
| `make new-bug BUG="x"` | Create new bug report |

## Component Patterns (Enforced)

| Need | Use | Not |
|------|-----|-----|
| Table | `behaviors.TableBehavior[T]` | `table.New()` |
| Form in intent | `models.*Form` wrapper | `*huh.Form` directly |
| Text/titles | `primitives.Title()`, `primitives.Body()` | Raw lipgloss |
| Badges | `primitives.HelpKeyBadge()` | `components.KeyBadge` |
| Colors | `theme.Primary()` etc | `lipgloss.Color("#xxx")` |
| Layout | `layout.ScreenLayout` | Manual composition |
| Modals | `feedback.Modal`, `behaviors.RenderModalOverlay()` | Custom modal code |
| Intent | Embed `*BaseIntent` | Custom base |

Run `make what-to-use NEED="keyword"` for detailed usage and examples.

## Workflow

```
1. make session-start
2. Pick ONE task from tasks/
3. make pre-task
4. Write test FIRST
5. Implement minimal code to pass
6. make check-compliance
7. make ai-commit FILE=/tmp/commit.txt
```

## When to Refuse

- Skipping `make session-start`
- Implementation before test
- Multiple changes per commit
- Using `git commit` directly
- Skipping compliance checks
- PRs targeting `main`
- Hardcoded colors/styles
- Raw `*huh.Form` in intents

## Finding Documentation

```bash
make what-to-use NEED="keyword"   # Component help with examples
```

### Development Guides (docs/development/)

| Topic | Document |
|-------|----------|
| Session protocol | [docs/development/SESSION_PROTOCOL.md](docs/development/SESSION_PROTOCOL.md) |
| BDD workflow | [docs/development/BDD_WORKFLOW.md](docs/development/BDD_WORKFLOW.md) |
| Architecture | [docs/development/ARCHITECTURE_OVERVIEW.md](docs/development/ARCHITECTURE_OVERVIEW.md) |
| Development workflow | [docs/development/DEVELOPMENT_WORKFLOW.md](docs/development/DEVELOPMENT_WORKFLOW.md) |
| Common tasks | [docs/development/COMMON_TASKS.md](docs/development/COMMON_TASKS.md) |
| Intent patterns | [docs/development/INTENT_PATTERNS_LIBRARY.md](docs/development/INTENT_PATTERNS_LIBRARY.md) |
| Keyboard system | [docs/development/KEYBOARD_SYSTEM_GUIDE.md](docs/development/KEYBOARD_SYSTEM_GUIDE.md) |

### Reference Guides

| Topic | Document |
|-------|----------|
| Intent architecture | [docs/INTENT_ARCHITECTURE_GUIDE.md](docs/INTENT_ARCHITECTURE_GUIDE.md) |
| UIKit components | [docs/UIKIT_GUIDE.md](docs/UIKIT_GUIDE.md) |
| Code standards | [docs/rules/senior-engineer-guidelines.md](docs/rules/senior-engineer-guidelines.md) |
| Commit rules | [docs/rules/AI_COMMIT_ATTRIBUTION.md](docs/rules/AI_COMMIT_ATTRIBUTION.md) |
| Branching | [docs/BRANCHING_STRATEGY.md](docs/BRANCHING_STRATEGY.md) |

## Code Examples

See `examples/` directory:
- `intent_template.go.example` - Intent pattern
- `form_wrapper_template.go.example` - Form wrapper pattern  
- `behavior_usage.go.example` - Behavior patterns

## Project Structure

```
internal/cli/
├── intents/     # Workflows (state machines)
├── behaviors/   # Reusable behaviors (TableBehavior, CRUD)
├── components/  # LEGACY - migrate to uikit/ (see UIKIT_GUIDE.md)
├── uikit/       # UIKit component library
│   ├── primitives/  # Text, Button, Badge, Input
│   ├── containers/  # Box, Overlay
│   ├── feedback/    # Modal, ModalContainer, HelpModal
│   ├── layout/      # ScreenLayout, Header, Footer
│   └── theme/       # Theme integration
├── models/      # Form wrappers
└── forms/       # Form configs
```

---

**When in doubt**: `make what-to-use NEED="keyword"` or `make check-patterns`

# KaRiya AI Agent Instructions

Go TUI application using Bubble Tea. Follow these rules strictly.

## Foundation Rules (MANDATORY)

### ALWAYS
- Follow the architecture and code style guidelines
- Write tests FIRST (TDD)
- Capture knowledge in memory (MCP) and obsidian (MCP) when DISCOVERING
  or CHANGING something
- Use `make ai-commit` for commits (not `git commit`)
- Ask if confidence is ASSUMED or UNKNOWN
- NEVER skip checks or tests
  - Unless explicitly approved by user, NEVER skip `make session-start` or any
    tests
- Refuse if about to do something that violates these rules
- **NEVER declare a task "done" or "complete" - only the USER decides when work is finished**
- **NEVER decide to skip, postpone, or deprioritize work - only the USER makes these decisions**

### Token Efficiency
- Be concise and precise
- No unnecessary words
- Specific, not vague
- Structured output

### Pre-Action Framework

Before ANY action, complete this checklist:

```
1. STOP  - What am I being asked to do?
2. THINK - What do I KNOW vs ASSUME vs UNKNOWN?
3. INVESTIGATE - Use tools, skills, git history to find out
4. CHOOSE - Which skills/tools are appropriate?
5. CONFIDENCE - VERIFIED / SUPPORTED / ASSUMED / UNKNOWN?
6. ACT or ASK - If ASSUMED/UNKNOWN: ask user first
```

## Memory Capture (MANDATORY)

Capture knowledge when you **DISCOVER** or **CHANGE** something:

| Trigger | What to Capture | Example |
|---------|-----------------|---------|
| DISCOVERY | Validation rules, patterns, gotchas | "Event text requires 10+ chars" |
| CHANGE | New rules, modified behaviour | "Changed min from 0 to 10 chars" |

Use `mcp_memory_create_entities` or `mcp_memory_add_observations`.

## Commands

| Task | Command |
|------|---------|
| All tests | `make test` |
| Single test | `make individual-test TEST="description"` |
| BDD tests | `make bdd` |
| Build | `make build` |
| Lint | `make vet && make staticcheck` |
| Compliance | `make check-compliance` |
| Commit | `AI_MODEL="claude-opus-4-6" make ai-commit FILE=/tmp/commit.txt` |

**Note:** The `AI_AGENT` is auto-detected from the `OPENCODE` environment variable and should NOT be overridden.

## Code Style

### Imports (grouped, alphabetical)
```go
import (
    "context"                                          // stdlib
    tea "github.com/charmbracelet/bubbletea"          // external
    "github.com/baphled/kariya/internal/domain/career" // internal
)
```

### Naming
- Files: `snake_case.go` | Packages: `lowercase`
- Types/Public: `PascalCase` | Private: `camelCase`

### Forbidden
- Comments inside function bodies
- `TODO`, `FIXME`, `HACK` markers
- Hardcoded colors (use `theme.Primary()`)
- Direct `huh.*` (use `forms/` wrappers)

## Architecture

```
App -> Intents -> Views/Modals -> UIKit -> Behaviors
```

- `views/` NEVER imports `intents/`
- `uikit/` NEVER imports `views/` or `intents/`
- All domain types MUST be converted to display types at intent boundaries
- `internal/ui/display/` defines common types (Event, Skill, Burst, CVView, etc.)
- Screens and views MUST NOT import `domain/career` or `service/career`
- Intents use a 5-file pattern: `constants.go`, `types.go`, `intent.go`, `handlers.go`, `helpers.go`
- Async service calls MUST return a message with an `Err error` field

### Modal Orchestration

- `ModalRegistry` is the central hub for modal coordination, rebuilt every render cycle
- Priority-based rendering ensures critical alerts override standard views
- Four standard adapters: `ErrorModalAdapter`, `FormModalAdapter[T]`, `ConfirmModalAdapter`, `ViewModalAdapter`
- `FormViewAdapter` bridges `widgets.View` to `ManagedModal` for search/sort/filter overlays

### Components

| Need | Use | NOT |
|------|-----|-----|
| Table | `behaviors.TableBehavior[T]` | `table.New()` |
| Forms | `forms.NewInput()` | Direct `huh.*` |
| Form Fields | `forms.FormValue[TD, TBind]` | Monolithic `Form[T]` |
| Modals | `feedback.Modal` | Custom code |
| Colors | `theme.Primary()` | `lipgloss.Color()` |
| Badges | `primitives.BadgeBuilder` | Custom styling |

- Views define package-level badge function slices for canonical patterns
- Action keys are defined in `views/*/actions.go` with `ActionKey` constants
- Navigation payloads use typed `Nav` structs handled by intents

### Interface Contracts

- `widgets.View` requires `RenderContent()`, `HelpText()`, `Init()`, and `Update()`
- `Update()` returns `(tea.Cmd, ViewResult)` where result is `Navigate`, `Cancel`, `Submit`, or `Error`
- Service operations are wrapped in `func() tea.Msg {}` closures
- Handlers MUST check `msg.Err` before proceeding with state transitions

## Testing

### TDD Cycle
1. Write failing test FIRST
2. Implement minimal code
3. Refactor
4. `make check-compliance`
5. `make ai-commit`

### Test Location
- E2E: `internal/testutil/e2e/{workflow}_e2e_test.go`
- BDD: `features/*.feature`

### BDD Step Rules
- Steps MUST use service-layer data creation
- NEVER use `env.TypeText()`, `env.Tab()`, or `env.ClearTextField()` in `When` steps

## Git

- Feature branches only (never commit to next/main)
- PRs target `next` branch
- Use `make ai-commit` (not `git commit`)

## When to REFUSE

- Skip `make session-start`
- Commit to `next`/`main` directly
- Write code before tests
- Use `git commit` instead of `make ai-commit`
- Add TODO/FIXME comments
- Import `intents/` from `views/`
- **Commit code when tests are failing or hanging without explicit user approval**
- **Declare task complete when the USER has not confirmed completion**
- **Make scope decisions (what's "separate", what's "additional", what can wait)**

## When to ASK

- Confidence is ASSUMED or UNKNOWN
- Multiple valid interpretations exist
- Tests fail but unsure if test or code is wrong
- About to "fix" something that might be intentional
- **When you encounter ANY blocker, failure, or issue - NEVER decide independently to skip it**
- **When tests don't pass - NEVER commit and declare it "mostly done" without user approval**
- **When you want to mark something as "needs investigation" - ASK FIRST**

## Key Docs

#MV|- Architecture: See AGENTS.md architecture section above and Obsidian KB
#HN|- Testing: See AGENTS.md testing section above and Obsidian KB
#QZ|- Intents: See AGENTS.md architecture section above and Obsidian KB

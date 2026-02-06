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
| Commit | `make ai-commit FILE=/tmp/commit.txt` |

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
App -> Intents -> Screens/Modals -> UIKit -> Behaviors
```

- `screens/` NEVER imports `intents/`
- `uikit/` NEVER imports `screens/` or `intents/`

### Components

| Need | Use | NOT |
|------|-----|-----|
| Table | `behaviors.TableBehavior[T]` | `table.New()` |
| Forms | `forms.NewInput()` | Direct `huh.*` |
| Modals | `feedback.Modal` | Custom code |
| Colors | `theme.Primary()` | `lipgloss.Color()` |

## Testing

### TDD Cycle
1. Write failing test FIRST
2. Implement minimal code
3. Refactor
4. `make check-compliance`
5. `make ai-commit`

## 🎬 VHS Demo Generation

When creating or modifying TUI workflows, generate visual documentation.

**IMPORTANT**: When a feature task is marked as "done", the [VHS Demo Generation Prompt](docs/prompts/VHS_DEMO_GENERATION_PROMPT.md) should be triggered to ensure visual documentation exists.

**Trigger Conditions**:
| Change Type | Action | Location |
|-------------|--------|----------|
| **New Feature** | Create 3 new tapes | `demos/vhs/features/<feature>/` |
| **Bug Fix** | Find and update existing tape | Existing tape location |
| **Enhancement** | Find and update existing tape | Existing tape location |

**Required for**:
- New intents/workflows
- UI changes that affect user experience
- Bug fixes with visual impact

**Demo Scenarios** (all three required for new features):

| Scenario | File | Purpose |
|----------|------|---------|
| **Happy path** | `happy-path.tape` | Successful workflow completion |
| **Sad path** | `sad-path.tape` | Error handling, validation |
| **Edge cases** | `edge-cases.tape` | Cancel, back nav, empty states |

**Workflow**:
```bash
# 1. Create feature demo directory
mkdir -p demos/vhs/features/your-feature

# 2. Copy templates
cp demos/vhs/features/template/*.tape demos/vhs/features/your-feature/

# 3. Customize tapes for your feature

# 4. Generate demos
make vhs-feature FEATURE=your-feature

# 5. Include GIFs in PR description
```

**Quick Commands**:
```bash
make vhs-demos                    # Generate all workflow demos
make vhs-{workflow}               # Generate specific workflow (vhs-capture, vhs-browse, etc.)
make vhs-feature FEATURE=name     # Generate feature demo for PR
make vhs-golden-compare           # Visual regression test
```

**📖 VHS Demo Prompt:** [VHS_DEMO_GENERATION_PROMPT.md](docs/prompts/VHS_DEMO_GENERATION_PROMPT.md)  
**📖 Workflow Guide:** [WORKFLOW_DOCUMENTATION_GUIDE.md](docs/workflows/WORKFLOW_DOCUMENTATION_GUIDE.md)  
**📖 VHS README:** [demos/vhs/README.md](demos/vhs/README.md)

## Testing

### Test Location
- E2E: `internal/testutil/e2e/{workflow}_e2e_test.go`
- BDD: `features/*.feature`

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
- Import `intents/` from `screens/`
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

- Architecture: `docs/development/ARCHITECTURE_OVERVIEW.md`
- Testing: `docs/development/BDD_WORKFLOW.md`
- Intents: `docs/INTENT_ARCHITECTURE_GUIDE.md`

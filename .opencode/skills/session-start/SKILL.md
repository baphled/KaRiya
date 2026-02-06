---
name: session-start
description: Initialize a KaRiya development session with mandatory environment validation and rules acknowledgment
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

## What I do

Initialize a development session for KaRiya by running mandatory environment validation and acknowledging project rules.

## When to use me

Use this skill at the **start of every work session** before making any changes.

## Procedure

1. Run the session start command:
   ```bash
   make session-start
   ```

2. If it fails, fix the issues before proceeding. Common issues:
   - Git hooks not installed: `make install-git-hooks`
   - Missing tools: `make ci-install-tools`

3. Acknowledge these **Critical Rules (Zero Tolerance)**:
   - **Feature branches** - NEVER commit directly to `next` or `main`
   - **PRs target `next`** - Never `main`. Only `next->main` for releases
   - **TDD** - Write test FIRST, then implementation (Red->Green->Refactor)
   - **Commits** - Use `make ai-commit FILE=<path>` only (not `git commit`)
   - **Compliance** - Run `make check-compliance` before AND after tasks
   - **One task** - One logical change per commit
   - **Architecture** - Follow layer hierarchy, no shortcuts
   - **Documentation** - Required for all exported symbols

## AI Attribution Setup

**Before any commits**, know your correct attribution:

| If running in | Set `AI_AGENT` to | Set `AI_MODEL` to |
|---------------|-------------------|-------------------|
| Opencode with Opus | `Opencode` | `Claude Opus 4.5` |
| Opencode with Sonnet | `Opencode` | `Claude Sonnet 4` |
| Claude Code | `Claude Code` | `Claude Opus 4.5` or `Claude Sonnet 4` |

**Use when committing:**
```bash
AI_AGENT="Opencode" AI_MODEL="Claude Opus 4.5" make ai-commit FILE=/tmp/commit.txt
```

The script auto-detects the agent but **defaults to wrong models**. Always set `AI_MODEL` explicitly.

## Session workflow

```
make session-start
  -> work on task
  -> make check-compliance
  -> AI_AGENT="Opencode" AI_MODEL="Claude Opus 4.5" make ai-commit FILE=/tmp/commit.txt
  -> repeat
make session-end
```

## If session-start fails

**REFUSE to proceed** until issues are fixed. The session validation ensures:
- Git hooks are installed
- All CI tools are available
- Pattern compliance is checked
- Architecture rules are enforced

## Related skills

- `tdd-workflow` - Follow TDD cycle for implementation
- `check-compliance` - Validate code before commits
- `ai-commit` - Create properly attributed commits

# AI Agent Prompts

This directory contains specialized prompts for AI agents working on the KaRiya project. These prompts are triggered automatically during specific workflow events.

## Available Prompts

| Prompt | Trigger | Purpose |
|--------|---------|---------|
| [VHS_DEMO_GENERATION_PROMPT.md](VHS_DEMO_GENERATION_PROMPT.md) | Task marked "done" | Generate visual documentation for UI features |

## How Prompts Work

AI agents (Claude Code, Cursor, etc.) should automatically trigger these prompts when specific conditions are met. The trigger conditions are defined in `AGENTS.md`.

### VHS Demo Generation

**Trigger Conditions** (from AGENTS.md):

| Change Type | Action | Location |
|-------------|--------|----------|
| **New Feature** | Create 3 new tapes | `demos/vhs/features/<feature>/` |
| **Bug Fix** | Find and update existing tape | Existing tape location |
| **Enhancement** | Find and update existing tape | Existing tape location |

When a feature task is marked as "done", the AI agent should:
1. Analyze what was changed
2. Follow the VHS Demo Generation Prompt
3. Create or update VHS tapes as needed
4. Include demo GIFs in PR description

## Adding New Prompts

When adding a new prompt:

1. Create the prompt file in this directory
2. Add an entry to this README
3. Update `AGENTS.md` with trigger conditions
4. Document the expected behavior

## Directory Structure

```
docs/prompts/
├── README.md                      # This file
└── VHS_DEMO_GENERATION_PROMPT.md  # VHS demo automation
```

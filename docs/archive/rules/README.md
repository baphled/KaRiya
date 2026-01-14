---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# Archived Rules

This directory contains rules and guidelines that are no longer actively used but are preserved for historical reference.

## Why Archive?

Files are archived when:
- The workflow they describe has been replaced
- They are redundant with other active documentation
- They are no longer referenced by active code or processes

## Archived Files

### generate-prd.md (Archived 2026-01-06)
**Reason**: PRD generation workflow was replaced with direct task file creation
**Replacement**: Tasks are now created directly as `tasks-XX-name.md` without separate PRDs

### generate-tasks.md (Archived 2026-01-06)
**Reason**: Task generation from PRDs workflow was replaced
**Replacement**: Single comprehensive task files with all phases together

### task-instructions.md (Archived 2026-01-06)
**Reason**: Redundant with master-task-prompt.md
**Replacement**: docs/rules/master-task-prompt.md contains all task instructions

### process-task-list.md (Archived 2026-01-08)
**Reason**: Consolidated into master-task-prompt.md for single source of truth
**Replacement**: docs/rules/master-task-prompt.md "Task Processing" section contains all task processing rules
**Key concepts preserved**: Authority order, completion criteria, refactoring constraints, one-task-at-a-time execution

## Restoration

If you need to reference these files, they are preserved here unchanged. To restore a file to active use:
1. Move it back to docs/rules/
2. Update docs/rules/README.md to list it
3. Update AGENTS.md to reference it
4. Ensure it's integrated into current workflows

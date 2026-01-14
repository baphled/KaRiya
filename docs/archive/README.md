---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# Archived Documentation

This directory contains historical documentation from completed phases, deprecated features, and superseded implementation plans. These documents are preserved for reference but are no longer actively maintained.

## Directory Structure

### `phases/` (20+ files)
**Purpose**: Completed phase reports and summaries

Historical completion reports from project phases:
- Phase 1-12 completion reports (Dec 2025 - Jan 2026)
- Phase implementation guides
- Phase structure analyses
- Phase preparation documents
- ATOMIC_TASKS_PHASE3.md (Phase 3 atomic tasks)

**Why archived**: All project phases are complete. These documents serve as historical record of implementation progress.

### `plans/` (6 files)
**Purpose**: Superseded implementation plans

Completed implementation plans:
- CV_GENERATION_IMPLEMENTATION_PLAN.md
- NEW_INTENT_IMPLEMENTATION_PLAN.md
- COMPONENT_EXTRACTION_PLAN.md
- INTEGRATION_ACTION_PLAN.md
- PHASE_12_CV_EXPORT_IMPLEMENTATION_PLAN.md
- PHASE_2_ENHANCED_BULLET_GENERATOR_PLAN.md

**Why archived**: Implementation work complete. Plans superseded by actual code and updated documentation.

### `app-migration/` (6 files)
**Purpose**: App.go migration documentation

Complete migration documentation from legacy app.go to intent-driven architecture:
- APP_GO_MIGRATION_PLAN.md
- APP_GO_AGGRESSIVE_REPLACEMENT_PLAN.md
- APP_GO_MIGRATION_CODE_EXAMPLES.md
- APP_GO_MIGRATION_QUICK_START.md
- APP_GO_INTENT_INTEGRATION_REVIEW.md
- AGGRESSIVE_REPLACEMENT_START_NOW.md

**Why archived**: Migration complete (Phase 6). Intent architecture is now production standard.

### `investigations/` (10 files)
**Purpose**: One-time analysis and investigation reports

Historical analysis documents:
- BURST_FACTS_INTEGRATION_REPORT.md
- BURST_FACTS_VERIFICATION.md
- FEATURE_REVIEW_BURST_FACT_EXTRACTION.md
- FEATURE_REVIEW_NEW_INTENTS.md
- CSV_IMPORT_ANALYSIS.md
- CSV_IMPORT_SUCCESS_SUMMARY.md
- BUG_FIX_CV_GENERATION_FILTER.md
- LEGACY_SCREEN_SPECIAL_LOGIC.md
- NAVIGATION_TEST_COVERAGE_REPORT.md
- DOCUMENTATION_REVIEW_SUMMARY.md

**Why archived**: One-time investigations. Issues resolved or features implemented. Kept for reference.

### `sessions/` (6 files)
**Purpose**: Session summaries and completion reports

Historical session work summaries:
- FINAL_SESSION_SUMMARY.md
- CLEANUP_SUMMARY.md
- TASK17-COMPLETE-SUMMARY.md
- PRIORITY2_COMPLETION_REPORT.md
- PRIORITY_5_VERIFICATION_REPORT.md
- ESCAPE_KEY_STANDARDIZATION_COMPLETE.md

**Why archived**: Session-specific work complete. Outcomes integrated into main documentation.

### `summaries/` (4 files)
**Purpose**: Implementation summaries

- IMPLEMENTATION_SUMMARY.md (burst/fact persistence fix)
- BURST_SUGGESTION_MIGRATION_COMPLETE.md (huh forms migration Phase 5)

**Why archived**: Implementation work documented and complete.

### `proposals/` (2 files)
**Purpose**: Architecture proposals (accepted/rejected)

Historical architecture proposals:
- WORKFLOW_REFACTORING_STRATEGY.md (navigation state machine proposal)
- REFINED_WORKFLOW_STRATEGY.md (intent-driven navigation proposal)

**Why archived**: Proposals accepted and implemented in current intent architecture.

### `rules/` (5 files)
**Purpose**: Superseded or deprecated rule documents

Archived development rules:
- process-task-list.md (superseded by master-task-prompt.md)

**Why archived**: Rules superseded by updated versions in `docs/rules/`.

### `reports/` (3 files)
**Purpose**: Historical performance and analysis reports

- CLI_PERFORMANCE_BENCHMARKS_2025-12.md (older CLI benchmarks)

**Why archived**: Superseded by PERFORMANCE_BENCHMARKS.md (intent-focused, 2026-01).

## Archive Policy

### What Gets Archived

Documents are archived when they meet one or more of these criteria:

1. **Work Complete**: Implementation finished, outcomes documented elsewhere
2. **Superseded**: Replaced by newer, more accurate documentation
3. **One-Time**: Analysis or investigation completed, issues resolved
4. **Historical**: Phase/session-specific, no longer current
5. **Deprecated**: Feature removed or approach changed

### What Stays Active

Documents remain in main `docs/` when:

1. **Currently Used**: Referenced in active development workflows
2. **Living Documentation**: Continuously updated (guides, standards, references)
3. **Authoritative**: Primary source of truth for a topic
4. **Discovery**: Essential for onboarding or daily operations

## How to Use Archived Docs

### Finding Historical Information

If you need to understand:
- **"Why was X designed this way?"** → Check `proposals/` and `plans/`
- **"How was feature Y implemented?"** → Check `phases/` and `summaries/`
- **"What issues did we encounter?"** → Check `investigations/` and `sessions/`
- **"What changed in migration Z?"** → Check specific migration directories

### Referencing Archived Docs

When referencing archived documentation:
- Clearly mark as "ARCHIVED" or "Historical"
- Include archive date if available
- Link to current replacement documentation if it exists
- Explain why the archived doc is still relevant to the discussion

### Maintaining Archives

**Do not delete** archived documentation unless:
- Content is truly obsolete with no historical value
- Duplicates active documentation exactly
- Contains sensitive or incorrect information

**Do update** archived documentation to:
- Add "ARCHIVED [DATE]" prefix to title
- Include note explaining why it was archived
- Link to replacement documentation

## Related Documentation

**Active Documentation**:
- [`docs/README.md`](../README.md) - Current documentation index
- [`AGENTS.md`](../../AGENTS.md) - Comprehensive project handover
- [`docs/IMPLEMENTATION_ROADMAP.md`](../IMPLEMENTATION_ROADMAP.md) - Current implementation status

**Other Archives**:
- [`tasks/`](../../tasks/) - All task files (completed tasks remain for reference)
- [`CHANGELOG.md`](../../CHANGELOG.md) - Release history

---

**Last Updated**: 2026-01-08  
**Total Archived Files**: 52+  
**Archive Status**: Active (regularly maintained)

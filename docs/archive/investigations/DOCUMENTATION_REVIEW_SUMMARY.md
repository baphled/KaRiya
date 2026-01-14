---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# Documentation Review & Reorganization Plan
**Date:** 2025-12-23
**Status:** Ready for implementation

## Executive Summary

Reviewed all 42 markdown files across the KaRiya project. Found well-structured documentation with minor organizational improvements needed. Main issues:
- 6 setup guides cluttering root directory
- PRD files scattered across 3 locations
- Editor-specific documentation mixed with core docs
- Test reports without version history

## Current State

### Root Directory (10 files)
- ✅ **Keep**: README.md, AGENTS.md
- ⚠️ **Move**: 6 setup guides to `docs/setup/`
- ⚠️ **Move**: 2 editor-specific guides to `docs/tools/editor-setup/`
- ⚠️ **Move/Consolidate**: TESTING.md (redundant)
- ⚠️ **Archive**: TEST_REPORT.md to `docs/reports/` with date

### docs/rules/ (16 files)
- ✅ **Status**: Excellent organization, no changes needed
- Comprehensive coverage: commit guidelines, AI attribution, task workflows, compliance checks
- Clear hierarchy: full guides + quick references

### docs/features/ (8 files)
- ✅ **Status**: Well-structured feature specifications
- ⚠️ **Note**: Missing file 02 (sequence jumps from 01 to 03)

### docs/ (4 files)
- ✅ CI/CD documentation well-organized
- ⚠️ **Rename**: KaRiya.md → PRD_MASTER.md (clarify it's the authoritative PRD)

### features/ (1 file)
- ⚠️ **Move**: prd-kariya-career-journal.md → docs/PRD_USER_STORIES.md

### tasks/ (3 files)
- ✅ **Keep**: tasks-01-career-event-capture.md, tasks-prd-career-entry-cli.md
- ⚠️ **Move**: prd-career-entry-cli.md → docs/PRD_CLI.md

## Proposed Structure

```
KaRiya/
├── README.md                           # Project overview (keep)
├── AGENTS.md                           # Handover document (keep)
├── TESTING.md                          # DELETE (redundant with neotest docs)
│
├── docs/
│   ├── README.md                       # NEW: Documentation index
│   ├── PRD_MASTER.md                   # RENAMED: from KaRiya.md
│   ├── PRD_USER_STORIES.md             # MOVED: from features/
│   ├── PRD_CLI.md                      # MOVED: from tasks/
│   ├── integration-test-strategy.md    # Keep
│   ├── CI_CD_PIPELINE.md               # Keep
│   ├── CI_CD_QUICK_REF.md              # Keep
│   │
│   ├── setup/                          # NEW DIRECTORY
│   │   ├── README.md                   # NEW: Setup guide index
│   │   ├── AI_COMMIT_ATTRIBUTION_SETUP.md
│   │   ├── AI_COMMIT_SETUP.md
│   │   ├── ATOMIC_COMMITS_SETUP.md
│   │   ├── CI_CD_SETUP.md
│   │   ├── MASTER_TASK_PROMPT_SETUP.md
│   │   └── RULES_COMPLIANCE_SETUP.md
│   │
│   ├── tools/                          # NEW DIRECTORY
│   │   └── editor-setup/
│   │       ├── README.md               # NEW: Editor setup guide index
│   │       ├── NEOTEST_SETUP.md        # CONSOLIDATED: from 2 files
│   │       └── NEOVIM_CONFIG.md        # OPTIONAL: broader config guide
│   │
│   ├── reports/                        # NEW DIRECTORY
│   │   ├── README.md                   # NEW: Reports index
│   │   └── TEST_REPORT_2025-12-23.md   # MOVED + DATED
│   │
│   ├── rules/                          # EXISTING: No changes
│   │   ├── README.md                   # NEW: Rules index
│   │   └── [16 existing rule files]
│   │
│   └── features/                       # EXISTING: Minor cleanup
│       ├── README.md                   # NEW: Features index
│       ├── 01-career-event-capture.md
│       ├── 02-PLACEHOLDER.md           # FIX: Add missing or renumber
│       ├── 03-burst-fact-extraction.md
│       ├── 04-cv-generation.md
│       ├── 05-data-model-schema.md
│       ├── 06-user-experience.md
│       ├── 07-role-audience-filtering.md
│       ├── 08-metadata-validation.md
│       └── 09-export-integration.md
│
├── features/                           # EXISTING: Remove PRD
│   └── (empty or delete directory)
│
└── tasks/                              # EXISTING: Remove PRD
    ├── tasks-01-career-event-capture.md
    └── tasks-prd-career-entry-cli.md
```

## Action Items

### Phase 1: Create Directory Structure (Priority: High)
- [ ] Create `docs/setup/`
- [ ] Create `docs/tools/editor-setup/`
- [ ] Create `docs/reports/`

### Phase 2: Move Setup Guides (Priority: High)
- [ ] Move `AI_COMMIT_ATTRIBUTION_SETUP.md` → `docs/setup/`
- [ ] Move `AI_COMMIT_SETUP.md` → `docs/setup/`
- [ ] Move `ATOMIC_COMMITS_SETUP.md` → `docs/setup/`
- [ ] Move `CI_CD_SETUP.md` → `docs/setup/`
- [ ] Move `MASTER_TASK_PROMPT_SETUP.md` → `docs/setup/`
- [ ] Move `RULES_COMPLIANCE_SETUP.md` → `docs/setup/`

### Phase 3: Move & Consolidate Editor Docs (Priority: Medium)
- [ ] Consolidate `NEOTEST-GO-SETUP.md` + `NEOTEST_TROUBLESHOOTING.md` → `docs/tools/editor-setup/NEOTEST_SETUP.md`
- [ ] Delete `TESTING.md` (redundant content)

### Phase 4: Move & Rename PRD Files (Priority: High)
- [ ] Rename `docs/KaRiya.md` → `docs/PRD_MASTER.md`
- [ ] Move `features/prd-kariya-career-journal.md` → `docs/PRD_USER_STORIES.md`
- [ ] Move `tasks/prd-career-entry-cli.md` → `docs/PRD_CLI.md`

### Phase 5: Archive Reports (Priority: Low)
- [ ] Move `TEST_REPORT.md` → `docs/reports/TEST_REPORT_2025-12-23.md`

### Phase 6: Create Index Files (Priority: Medium)
- [ ] Create `docs/setup/README.md` (setup guide index)
- [ ] Create `docs/rules/README.md` (rules index)
- [ ] Create `docs/features/README.md` (features index)
- [ ] Create `docs/tools/editor-setup/README.md` (editor setup index)
- [ ] Create `docs/reports/README.md` (reports index)

### Phase 7: Update Cross-References (Priority: Medium)
- [ ] Update internal links in moved files
- [ ] Update AGENTS.md references to new paths
- [ ] Update README.md to point to new documentation structure
- [ ] Update any Makefile references to moved files

### Phase 8: Cleanup (Priority: Low)
- [ ] Consider removing empty `features/` directory
- [ ] Fix missing feature file 02 or renumber sequence
- [ ] Add cross-references between related documents

## Files Summary

### Files to Move (10)
1. Root → docs/setup/ (6 files)
2. Root → docs/tools/editor-setup/ (2 files consolidated)
3. Root → docs/reports/ (1 file)
4. features/ → docs/ (1 file)
5. tasks/ → docs/ (1 file)

### Files to Rename (1)
1. docs/KaRiya.md → docs/PRD_MASTER.md

### Files to Delete (1)
1. TESTING.md (redundant)

### Files to Create (6)
1. docs/README.md
2. docs/setup/README.md
3. docs/rules/README.md
4. docs/features/README.md
5. docs/tools/editor-setup/README.md
6. docs/reports/README.md

### Files to Keep in Root (2 + config files)
1. README.md
2. AGENTS.md
3. All configuration files (go.mod, Makefile, package.json, etc.)

## Benefits

### Immediate Benefits
1. **Cleaner Root**: Only essential files at project root
2. **Better Discovery**: Related docs grouped by purpose
3. **Professional Structure**: Matches industry standards
4. **Easier Onboarding**: Clear navigation for new developers

### Long-term Benefits
1. **Maintainability**: Easier to update organized docs
2. **Scalability**: Clear structure for adding new docs
3. **Version Control**: Historical reports with dates
4. **Cross-referencing**: Logical document relationships

## Implementation Order

1. **Phase 1**: Create directories (5 minutes)
2. **Phase 2 + 4**: Move critical docs (setup guides + PRDs) (15 minutes)
3. **Phase 6**: Create index files (30 minutes)
4. **Phase 3 + 5 + 7**: Consolidate, archive, update references (20 minutes)
5. **Phase 8**: Cleanup (10 minutes)

**Total Time**: ~80 minutes

## Validation Checklist

After implementation, verify:
- [ ] All markdown files have valid paths
- [ ] No broken internal links
- [ ] AGENTS.md references updated
- [ ] README.md updated
- [ ] Makefile targets work
- [ ] Git history preserved for moved files
- [ ] Cross-references functional

## Notes

- Use `git mv` for file moves to preserve history
- Update AGENTS.md section 18.1 with new documentation paths
- Consider adding table of contents to large documents
- Keep this review document for future reference

---

**Next Step**: Proceed with Phase 1 (create directories)


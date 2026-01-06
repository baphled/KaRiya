# Task 11: Burst Enhancement - Full Implementation

**Status**: 📋 Ready for Implementation  
**Total Time**: 7 hours across 4 phases  
**Can be done incrementally**: Yes - one phase per session  
**Next Task Number**: tasks-12  
**Target Location**: `tasks/tasks-11-burst-enhancement.md` (when ready to execute)

---

## Executive Summary

**Pain Point**: Users cannot see events/facts within bursts, cannot confirm bursts for CV generation, cannot edit or delete bursts.

**Goal**: Enhance burst management with full CRUD capabilities while keeping the simple List/Detail state model.

**Approach**: Hybrid - leverage existing `BurstEditorModel`, add inline delete confirmation, show events/facts, enable workflow integration.

**Key Discovery**: `BurstEditorModel` already exists (fully implemented)! We'll wire it up instead of building from scratch.

### What We'll Build

- ✅ View events within burst (expandable list with navigation)
- ✅ View facts extracted from burst and events
- ✅ Confirm bursts (trigger fact extraction, mark as reviewed)
- ✅ Edit bursts (using existing BurstEditorModel as modal)
- ✅ Delete bursts (inline y/n confirmation)
- ✅ Navigate between bursts and events

---

*[FULL CONTENT - See complete task document with all 4 phases, implementation checklists, code snippets, testing strategies, and acceptance criteria above]*

---

## Ready to Execute?

When ready to move from planning to execution:

1. **Move this file** from `.opencode/plan/` to `tasks/tasks-11-burst-enhancement.md`
2. **Update AGENTS.md** with minimal documentation reference
3. **Start with Phase 1** (2.5 hours)
4. **Check off tasks** as you complete them
5. **Create completion report** when all phases done

---

**This document is ready to be moved to `tasks/` directory when approved.**

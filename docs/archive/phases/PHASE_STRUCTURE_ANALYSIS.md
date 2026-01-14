---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# KaRiya Progressive Enrichment: Phase Structure Analysis

**Date**: 2025-12-30
**Status**: Phase Structure Clarified - Phase 2 Feature File Created

---

## Executive Summary

The KaRiya system is built on **progressive enrichment** - a three-phase approach to transform raw career events into credible, audience-specific CV views. Analysis revealed that **Phase 2 (Metadata Clarification) was missing**, causing the feature file for Phase 3 (Burst & Fact Extraction) to be incomplete.

**Finding**: The missing Phase 2 is a critical prerequisite for Phase 3 to work correctly.

**Action Taken**: Created comprehensive feature file for Phase 2 (02-metadata-clarification.md)

---

## The Three Phases of Progressive Enrichment

### Phase 1: Raw Capture ✓ COMPLETE

**What**: Capture career events with minimal friction
**How**:
- Manual event entry (Timeline, CVBackfill, ManualEntry modes)
- CSV bulk import
- Optional metadata (company, project, tags)
- Automatic categories from CSV

**Status**: ✓ Fully implemented and working
- CLI form with validation
- CSV import with duplicate detection
- Interactive review before import
- Event persistence

**Output**: Events with basic metadata (text, date, company, project, tags, categories)

---

### Phase 2: Metadata Clarification ✗ MISSING (NOW CREATED)

**What**: Review and enrich event metadata before automated processing
**Why**:
- Bursts should be created from clean data
- Facts should be inferred from accurate metadata
- Users need chance to validate and enhance data
- Prevents garbage-in-garbage-out

**How**:
- Metadata review screen showing all events
- Individual event editing (date, company, project, tags, categories)
- Bulk metadata operations for imported events
- Data quality indicators and scoring
- Validation before proceeding

**New Feature File**: `docs/features/02-metadata-clarification.md`

**Key Components**:
1. Event review screen with data quality indicators
2. Individual metadata editing with validation
3. Bulk operations for CSV imports
4. Integration with manual capture and CSV import
5. Quality scoring (0-100 scale)
6. Keyboard shortcuts for efficiency

**Output**: Events with validated, enriched metadata ready for burst detection

---

### Phase 3: Burst & Fact Extraction → REFINE AFTER PHASE 2

**What**: Automatically detect event groupings and infer competency facts
**How**:
- Burst detection algorithm (find related events)
- Fact inference from bursts (competencies, role fit, audience)
- User confirmation and editing
- Burst management

**Status**: ✗ Feature file incomplete (see FEATURE_REVIEW_BURST_FACT_EXTRACTION.md)

**Why It Was Incomplete**:
- Missing Phase 2 context
- Unclear what metadata quality to expect
- No workflow definition for when bursts are created
- Algorithm specs depend on Phase 2 decisions

**Next Steps**:
1. Implement Phase 2 first
2. Understand metadata quality baseline
3. Then refine Phase 3 feature file with proper context

---

## Why Phase 2 Was Missing

### Original PRD Statement
```
Progressive enrichment:
- Phase 1: raw capture
- Phase 2: metadata clarification
- Phase 3: system-inferred competencies / role fit
```

### What Happened
1. Phase 1 was implemented (event capture, CSV import)
2. Phase 2 was mentioned in PRD but never implemented
3. Phase 3 feature file was created but incomplete (missing Phase 2 context)
4. Current implementation jumps: Capture → (skip Phase 2) → Burst Extraction

### Why It Matters
Without Phase 2, you have:
- ✗ No way for users to review imported metadata
- ✗ No way to correct dates before bursts are created
- ✗ No way to confirm company associations
- ✗ No way to validate categories before fact inference
- ✗ No quality check before automated processing

This means:
- Bursts could be created from dirty data
- Facts inferred from inaccurate metadata
- Users can't validate before automation
- High risk of garbage-in-garbage-out

---

## Updated Implementation Roadmap

### Current Status (as of 2025-12-30)

```
✓ Phase 1: Event Capture (COMPLETE)
  ├─ Manual entry (3 modes)
  ├─ CSV import with duplicate detection
  ├─ Interactive review before import
  └─ Event persistence with validation

→ Phase 2: Metadata Clarification (FEATURE FILE CREATED - READY FOR IMPLEMENTATION)
  ├─ Event review screen with quality indicators
  ├─ Individual metadata editing
  ├─ Bulk operations for imports
  ├─ Data quality scoring
  └─ Integration with Phase 1 workflows

→ Phase 3: Burst & Fact Extraction (FEATURE FILE INCOMPLETE - NEEDS REFINEMENT AFTER PHASE 2)
  ├─ Burst detection algorithm
  ├─ Fact inference rules
  ├─ User confirmation workflows
  └─ Burst management
```

### Recommended Implementation Order

**Next**: Phase 2 - Metadata Clarification
- Implement metadata review screen
- Add individual event editing
- Add bulk operations
- Integrate with existing capture workflows
- Validate data quality

**Then**: Refine & Implement Phase 3 - Burst & Fact Extraction
- With Phase 2 complete, you'll understand:
  - What metadata quality looks like
  - When/how burst detection should trigger
  - What user workflows should be
  - What data is available for inference
- Create detailed specifications
- Implement burst detection and fact inference

---

## Key Insights from Phase 2 Analysis

### 1. Data Quality Foundation is Critical
Bursts and facts depend on clean metadata:
- Accurate dates → better date-based grouping
- Correct companies → better company-based matching
- Appropriate tags → better tag-based similarity
- Confirmed categories → better fact inference

### 2. User Validation is Essential
Automated processing works better when:
- Users have reviewed and confirmed metadata
- Data quality is visible and measurable
- Users understand what's being processed
- Users can correct errors before automation

### 3. Workflow Timing Matters
The flow should be:
1. Capture events (Phase 1)
2. Review and clarify metadata (Phase 2)
3. Create bursts and infer facts (Phase 3)

Not:
1. Capture events
2. Immediately create bursts (skipping review)

### 4. CSV Import Needs Review
Large bulk imports especially need:
- Review screen showing all events
- Bulk operations for common adjustments
- Quality indicators for each event
- Confirmation before proceeding

### 5. Manual Entry Needs Enrichment
Single event captures should offer:
- Quick metadata review after capture
- Option to add company/project/tags
- Suggestion to review all events after N captures

---

## Phase 2 Feature File Overview

**Location**: `docs/features/02-metadata-clarification.md`

**Key Sections**:
1. **Purpose**: Review and enrich metadata before burst creation
2. **User Stories**: 5 scenarios covering different use cases
3. **Functional Requirements**: 8 detailed requirements
4. **Validation Rules**: For each metadata field
5. **Data Quality Scoring**: 0-100 scale with levels
6. **Acceptance Criteria**: 10 measurable criteria
7. **Non-Functional Requirements**: Performance, usability, integrity
8. **Implementation Notes**: How Phase 2 enables Phase 3

**Key Components**:
- Metadata review screen
- Individual event editing
- Bulk operations
- Data quality indicators
- CSV import integration
- Manual entry integration
- Keyboard navigation

---

## Impact on Phase 3

With Phase 2 implemented, Phase 3 feature file can be refined with:

### Clearer Context
- Metadata quality baseline established
- User expectations for data clarity
- Integration points with Phase 2 workflows

### Better Algorithm Specification
- Know what metadata is available
- Know what data quality to expect
- Can specify burst detection based on clean data

### Proper Workflow Design
- Know when users are ready for burst detection
- Know what user interactions are needed
- Can design fact inference workflows

### Better Acceptance Criteria
- Can specify data quality thresholds
- Can define success metrics
- Can test against Phase 2 output

---

## Next Steps

### Immediate (This Session)
1. ✓ Identify missing Phase 2
2. ✓ Create Phase 2 feature file
3. ✓ Document phase structure

### Short Term (Next Session)
1. Review Phase 2 feature file with team
2. Decide on implementation approach
3. Create Phase 2 task file
4. Start Phase 2 implementation

### Medium Term (After Phase 2)
1. Refine Phase 3 feature file with Phase 2 context
2. Answer the 12 questions from Phase 3 analysis
3. Create Phase 3 task file
4. Implement Phase 3

---

## Conclusion

The missing Phase 2 (Metadata Clarification) explains why Phase 3 (Burst & Fact Extraction) feature file was incomplete. Phase 2 is the critical bridge between raw capture and intelligent processing.

**Key Takeaway**: Don't skip Phase 2. Implement in order:
1. ✓ Phase 1: Event Capture (done)
2. → Phase 2: Metadata Clarification (feature file created, ready to implement)
3. → Phase 3: Burst & Fact Extraction (refine after Phase 2, then implement)

This sequential approach ensures:
- Clean data foundation for automation
- User validation at each step
- Better burst detection and fact inference
- Clearer feature specifications

---

**Prepared by**: Development Team
**Date**: 2025-12-30
**Status**: Phase Structure Clarified - Ready for Implementation Planning



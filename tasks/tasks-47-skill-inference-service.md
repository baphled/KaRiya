# Task 47: Skill Inference Service

## Overview
- **Goal**: Create skill inference service that detects technologies in event text and suggests skills to users
- **Time Estimate**: 12-16 hours (3 trigger points + generic modal refactoring)
- **Prerequisites**: Task 45 (burst detection UI for pattern familiarity), understanding of skill domain model, text analysis basics

## Session Contract Acknowledgment
- [ ] Ran `make session-start` and it passed
- [ ] Acknowledge and commit to following all workflow rules
- [ ] Token count: _____ (must be < 50k to start)

## Pre-Task Checklist (MUST COMPLETE BEFORE STARTING)
- [ ] `make check-compliance` passes
- [ ] Reviewed existing patterns in:
  - `internal/domain/career/skill.go` (skill model)
  - `internal/repository/career/skill_repository.go` (persistence)
  - `internal/service/career/burstfact/detector.go` (similar detection pattern)
  - `internal/cli/intents/burst_management/` (subdirectory structure, suggestion UI pattern)
  - `internal/cli/intents/skillsmanagement/` (existing skills UI, subdirectory structure)
- [ ] Confirmed this is ONE atomic task (skill inference system)
- [ ] Identified which test files will be created/modified

## Current Status

**TASK 47 NOT STARTED** - Ready to Begin

---

## Architecture Overview

### Complete System Architecture

```
┌─────────────────────────────────────────────────────┐
│  USER TRIGGER POINTS (3 locations)                  │
├─────────────────────────────────────────────────────┤
│ 1. Burst Confirmation (automatic after facts)      │
│ 2. Burst Detail Modal ("i" key - manual)           │
│ 3. ManageSkills List ("i" key - manual)            │
└─────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────┐
│  SERVICE LAYER                                      │
├─────────────────────────────────────────────────────┤
│ SkillInferenceService                               │
│  ├─ InferSkillsFromBurst(burst, events)           │
│  ├─ InferSkillsFromEvents(events)                 │
│  └─ CreateSkillsFromSuggestions(suggestions)      │
│                                                     │
│ TechnologyKeywords (~100 entries)                   │
│  ├─ Backend (Go, Python, Ruby, Java...)           │
│  ├─ Frontend (React, Vue, TypeScript...)          │
│  ├─ Database (PostgreSQL, MongoDB, Redis...)      │
│  ├─ DevOps (Docker, Kubernetes, Terraform...)     │
│  ├─ Cloud (AWS, GCP, Azure, Lambda...)            │
│  ├─ Mobile (iOS, Swift, React Native...)          │
│  └─ Tooling (Git, GraphQL, REST, gRPC...)         │
└─────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────┐
│  DETECTION ALGORITHM                                │
├─────────────────────────────────────────────────────┤
│ 1. Regex word boundary matching (\bkeyword\b)     │
│ 2. Context extraction (~80 chars)                  │
│ 3. Confidence scoring (0.5 / 0.75 / 0.95)         │
│ 4. Deduplication (merge same skill)                │
└─────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────┐
│  UI LAYER (Generic Modal)                          │
├─────────────────────────────────────────────────────┤
│ SuggestionReviewModal (reused)                      │
│  ├─ Format burst suggestions                       │
│  └─ Format skill suggestions                       │
│     • Name: "Go" (Category: backend)               │
│     • Confidence: ████████████░░ 95%               │
│     • Contexts: "...built API using Go..."         │
│     • [a] Accept  [r] Reject  [A] Accept all       │
└─────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────┐
│  PERSISTENCE                                        │
├─────────────────────────────────────────────────────┤
│ 1. Create/reuse skill (SkillRepository)            │
│ 2. Link to events (junction table via event.Skills)│
│ 3. Update skill.LastUsed (most recent event date)  │
└─────────────────────────────────────────────────────┘
```

---

## Architecture Decisions

### Decision 1: Modal Reusability

**Decision**: Reuse `SuggestionReviewModal` for both burst and skill suggestions.

**Rationale**:
- Consistent UX across suggestion types
- Reduced code duplication
- Same keyboard shortcuts and patterns
- Single component to maintain and test

**Implementation**:
```go
type SuggestionReviewModal struct {
    suggestions     interface{}  // []burstfact.BurstSuggestion OR []SkillSuggestion
    suggestionType  string       // "burst" or "skill"
    currentIndex    int
    accepted        []interface{}
}

func (m *SuggestionReviewModal) View() string {
    switch m.suggestionType {
    case "burst":
        return m.renderBurstSuggestion()
    case "skill":
        return m.renderSkillSuggestion()
    }
}
```

### Decision 2: Three Trigger Points

**Decision**: THREE trigger points for skill inference (not two).

**Locations**:

#### 1. Automatic (Burst Confirmation)
- **When**: After fact extraction completes
- **What**: Analyzes burst events automatically
- **Why**: Seamless skill capture as users work

```
User confirms burst
    ↓
Facts extracted (3 facts)
    ↓
Skill inference runs automatically
    ↓
Modal shows suggestions (if any found)
```

#### 2. Manual (Burst Detail Modal)
- **When**: "i" key in burst detail modal
- **What**: Analyzes specific burst on-demand
- **Why**: Allows re-inference after event edits, targeted project analysis

```
User views burst detail
    ↓
Presses "i" key
    ↓
Skill inference runs for this burst
    ↓
Modal shows suggestions
```

#### 3. Manual (ManageSkills List)
- **When**: "i" key in skills list
- **What**: Analyzes ALL events globally
- **Why**: One-time bulk skill import for historical data

```
User in Skills List
    ↓
Presses "i" key
    ↓
Skill inference runs on ALL events
    ↓
Modal shows suggestions
```

**Benefits**:
- **Automatic**: No user action needed, captures skills naturally
- **Manual (burst)**: Per-project/company analysis (targeted)
- **Manual (global)**: Historical skill import (comprehensive)

### Decision 3: Timeline Adjustment

**Original**: 8-10 hours  
**Updated**: 12-16 hours

**Additions**:
- +30m: Phase 0 (pre-task setup - remove skipped test)
- +1h: Make SuggestionReviewModal generic
- +1h: Third trigger point (burst detail modal)
- +1h: Testing all three trigger points
- +1h: Documentation updates

---

## Context

Currently, skills are **manually managed** through the ManageSkills intent. Users must:
1. Manually create each skill
2. Manually link skills to events
3. Remember which technologies they've used

This task creates **automatic skill inference** that:
1. Analyzes event text for technology mentions (e.g., "built API with Go and PostgreSQL")
2. Suggests skills with confidence scores (0.0-1.0)
3. Auto-links inferred skills to source events via junction table
4. **Integrates into THREE locations**:
   - Automatic: burst confirmation flow (after fact extraction)
   - Manual: burst detail modal ("i" key)
   - Manual: ManageSkills list ("i" key)

### Why This Matters for CV Generation

From our analysis:
- CV "Core Competencies" currently uses generic categories from Facts
- The Skills table has rich data (Level, YearsUsed) but isn't used in CV generation
- Task 48 will use Skills table for professional CV skills sections
- **This task populates the Skills table automatically**

---

## Technology Dictionary Requirements

Based on user preference:
- **~100 comprehensive technology keywords** from the start
- Organized by category for easy maintenance
- File-based dictionary that can be internally updated
- Supports aliases (e.g., "golang" → "Go", "k8s" → "Kubernetes")

### Dictionary Coverage

```
~100 keywords across 7 categories:

├─ Backend (25 keywords)
│  Go, Python, Ruby, Java, Node.js, PHP, Rust, Scala, etc.
│
├─ Frontend (20 keywords)
│  React, Vue, Angular, TypeScript, Tailwind CSS, etc.
│
├─ Database (15 keywords)
│  PostgreSQL, MySQL, MongoDB, Redis, Elasticsearch, etc.
│
├─ DevOps (15 keywords)
│  Docker, Kubernetes, Terraform, Jenkins, CircleCI, etc.
│
├─ Cloud (15 keywords)
│  AWS, GCP, Azure, Lambda, S3, CloudFormation, etc.
│
├─ Mobile (5 keywords)
│  iOS, Android, Swift, Kotlin, React Native, Flutter, etc.
│
└─ Tooling (15 keywords)
   Git, GitHub, GraphQL, REST, gRPC, Kafka, etc.
```

---

## Files to Create

### Service Layer
- [ ] `internal/service/career/technology_keywords.go` - Technology dictionary (~100 entries)
- [ ] `internal/service/career/technology_keywords_test.go` - Dictionary tests
- [ ] `internal/service/career/skill_inference.go` - Service interface
- [ ] `internal/service/career/skill_inference_impl.go` - Implementation
- [ ] `internal/service/career/skill_inference_test.go` - Service tests

### Documentation
- [ ] `docs/features/SKILL_INFERENCE.md` - Comprehensive feature guide
- [ ] Update `docs/SKILLS_GUIDE.md` - Add skill inference section
- [ ] Update `docs/workflows/MANAGE_SKILLS_WORKFLOW.md` - Add inference workflow
- [ ] Update `docs/development/MODAL_OVERLAY_PATTERN.md` - Add generic modal pattern
- [ ] `docs/design/SKILL_INFERENCE_UI_DESIGN.md` - UI design specification

---

## Files to Modify

### Intent Layer (Subdirectory Structure)
- [ ] `internal/cli/intents/burst_management/constants.go` - Add skill suggestion states
- [ ] `internal/cli/intents/burst_management/messages.go` - Add SkillSuggestionsLoadedMsg
- [ ] `internal/cli/intents/burst_management/helpers.go` - Add inferSkillsFromBurst command
- [ ] `internal/cli/intents/burst_management/handlers.go` - Add handlers for skill suggestions
- [ ] `internal/cli/intents/skillsmanagement/constants.go` - Add StateInferring, StateSuggestionReview
- [ ] `internal/cli/intents/skillsmanagement/messages.go` - Add SkillSuggestionsLoadedMsg
- [ ] `internal/cli/intents/skillsmanagement/handlers.go` - Add "i" key handler

### UI Layer (Generic Modal)
- [ ] `internal/cli/screens/burst_management/modals/suggestion_review_modal.go` - Make generic for bursts AND skills

---

## Implementation Timeline

| Phase | Description | Time | Total |
|-------|-------------|------|-------|
| 0 | Pre-Task Setup (remove skipped test) | 15-30m | 15-30m |
| 1 | Technology Dictionary (~100 keywords) | 2-3h | 2.25-3.5h |
| 2 | Service Interface | 1h | 3.25-4.5h |
| 3 | Keyword Detection (word boundaries) | 2-3h | 5.25-7.5h |
| 4 | Confidence Scoring | 1-2h | 6.25-9.5h |
| 5 | Skill Creation & Event Linking | 2h | 8.25-11.5h |
| 6 | Burst Integration (2 triggers + generic modal) | 3h | 11.25-14.5h |
| 7 | ManageSkills Integration (1 trigger) | 1-2h | 12.25-16.5h |

**Total**: 12-16 hours (realistic for 3 trigger points + generic modal)

---

## Phase 0: Pre-Task Setup

**Goal**: Unblock `session-start` and establish baseline

**Why**: Skipped test blocks `session-start` (PROHIBITED in pre-commit hooks)

**Actions**:
1. Identify and remove (or implement) skipped test in `internal/cli/intents/consistency_test.go`
2. Run `make session-start` to verify it passes
3. Run `make check-compliance` to establish baseline
4. Update this task file (task number verified)

**TDD Checklist - Phase 0**:
- [ ] Locate skipped test: `grep -r "Skip(" internal/cli/intents/consistency_test.go`
- [ ] Decision: Remove test or implement it
- [ ] Run `make session-start` - should pass
- [ ] Run `make check-compliance` - establish baseline

**Commits** (1):
- [ ] `test(intents): remove skipped consistency test` OR `test(intents): implement consistency test`

**Time Estimate**: 15-30 minutes

---

## Detailed Phases

*Note: Due to length constraints, I'll include a summary. Full implementation details for Phases 1-7 remain as in the original task file with these key updates:*

### Key Updates Across All Phases:

**Phase 1-5**: Service layer implementation (unchanged from original spec)
**Phase 6**: Now includes:
- Subphase 6.1: Make SuggestionReviewModal generic
- Subphase 6.2: Add states and messages
- Subphase 6.3: Automatic trigger (after fact extraction)
- Subphase 6.4: Manual trigger (burst detail modal "i" key)

**Phase 7**: ManageSkills integration with "i" key (third trigger point)

---

## Expected UX Examples

### After burst confirmation (automatic):
```
┌─────────────────────────────────────────────────┐
│ Skill Suggestions (from 5 burst events)         │
│                                                 │
│ Suggestion 1 of 3                               │
│                                                 │
│ Go (backend) - 95% confidence                   │
│ ████████████████████░░░░                        │
│                                                 │
│ Found in 5 events:                              │
│   • "...built API using Go and gRPC for..."    │
│   • "...migrated service to Go for better..."  │
│   • "...wrote Go microservices that handle..." │
│                                                 │
│ [a] Accept  [r] Reject  [A] Accept all  [Esc]  │
└─────────────────────────────────────────────────┘
```

### In burst detail modal after "i" key (manual):
```
┌─────────────────────────────────────────────────┐
│ Analyzing Burst Events                          │
│                                                 │
│ ⏳ Detecting skills from 8 events...            │
│                                                 │
│ This may take a moment...                       │
└─────────────────────────────────────────────────┘
```

### In ManageSkills after pressing "i" (manual):
```
┌─────────────────────────────────────────────────┐
│ Inferring Skills from All Events                │
│                                                 │
│ ⏳ Analyzing 42 events for technology mentions...│
│                                                 │
│ This may take a few moments...                  │
└─────────────────────────────────────────────────┘
```

---

## Acceptance Criteria
- [ ] Technology dictionary has ~100 comprehensive entries organized by category
- [ ] Skill inference uses word boundary regex (prevents partial matches like "goal" matching "go")
- [ ] Confidence scoring differentiates high/medium/low usage patterns (0.95/0.75/0.5)
- [ ] Deduplication merges same skill from multiple events
- [ ] Created skills are linked to source events via junction table
- [ ] LastUsed is set to most recent event date
- [ ] **Automatic trigger**: Integration into burst confirmation works (after facts)
- [ ] **Manual trigger 1**: "i" key in burst detail modal works
- [ ] **Manual trigger 2**: "i" key in ManageSkills list works
- [ ] Generic modal handles both burst and skill suggestions
- [ ] Context snippets show actual technology usage (~80 chars)
- [ ] Empty suggestions handled gracefully (no error)

---

## Documentation

See additional documentation:
- **UI Design**: `docs/design/SKILL_INFERENCE_UI_DESIGN.md` - Complete UI analysis and reusability assessment
- **Feature Guide**: `docs/features/SKILL_INFERENCE.md` - Technical reference (to be created)
- **User Guide**: `docs/SKILLS_GUIDE.md` - End-user documentation (to be updated)
- **Workflow**: `docs/workflows/MANAGE_SKILLS_WORKFLOW.md` - Operational procedures (to be updated)

---

## Rollback Plan
1. Remove service files:
   - `internal/service/career/technology_keywords.go`
   - `internal/service/career/technology_keywords_test.go`
   - `internal/service/career/skill_inference.go`
   - `internal/service/career/skill_inference_impl.go`
   - `internal/service/career/skill_inference_test.go`
2. Revert intent changes:
   - `internal/cli/intents/burst_management/`
   - `internal/cli/intents/skillsmanagement/`
3. Revert modal changes:
   - `internal/cli/screens/burst_management/modals/suggestion_review_modal.go`
4. Run `make check-compliance`
5. All tests should pass

---

## Dependencies
- Task 45 (burst detection UI) - for similar pattern understanding ✅ Complete
- Skill domain model and repository ✅ Already exists
- SuggestionReviewModal ✅ Already exists (will be made generic)

---

## Next Steps After Completion
- Task 48: Enhanced CV Skills Section (will use skills created by this task)
- Task 44: Wire Enhanced Bullet Generator (independent, can be done in parallel)

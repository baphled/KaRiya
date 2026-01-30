# Task 47: Skill Inference Service

## Overview
- **Goal**: Create skill inference service that detects technologies AND soft skills in event text and suggests skills to users
- **Time Estimate**: 12-16 hours (service + UI) + 8-11 hours (enhancements: expanded keywords + soft skills)
- **Prerequisites**: Task 45 (burst detection UI for pattern familiarity), understanding of skill domain model, text analysis basics, familiarity with CompetencyCategory system

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
  - `internal/constants/constants.go` (CompetencyCategory enum for soft skills)
- [ ] Confirmed this is ONE atomic task (skill inference system)
- [ ] Identified which test files will be created/modified

## Current Status

**PHASE 6 COMPLETE** - Service Layer Production-Ready, Continuing with UI Integration and Enhancements

### Progress Summary

| Phase | Status | Tests | Description |
|-------|--------|-------|-------------|
| 0 | ✅ COMPLETE | - | Pre-task setup (removed skipped test) |
| 1 | ✅ COMPLETE | 33 | Technology Dictionary (150 keywords, 7 categories) |
| 2 | ✅ COMPLETE | 9 | Service Interface & Validation |
| 3 | ✅ COMPLETE | 21 | Keyword Detection (word boundaries, context extraction) |
| 4 | ✅ COMPLETE | 28 | Confidence Scoring (3-tier, proximity-based) |
| 5 | ✅ COMPLETE | 17 | Skill Persistence (create/update, event linking) |
| 6 | ✅ COMPLETE | 4 | Integration Tests & Documentation |
| 7 | 🔄 IN PROGRESS | - | UI Integration (burst_management) |
| 8 | ⏳ PENDING | - | UI Integration (skillsmanagement) |
| 9 | ⏳ PENDING | - | E2E Testing |
| **10** | **⏳ PENDING** | **~10** | **Expand Keyword Dictionary (150 → 230 keywords, 8 new categories)** |
| **11** | **⏳ PENDING** | **50-75** | **Add Soft Skills Detection (5 new competency categories)** |

**Total Tests**: 112 passing → 172-197 passing (after enhancements)  
**Total Lines**: 3,287 → 4,500-5,000 (after enhancements)  
**Commits**: 15 (all TDD, all AI-attributed) → ~25-30 after enhancements  
**Time Invested**: ~9.5 hours (service + docs) + TBD (UI + enhancements)  

### Service Layer: Production-Ready ✅

The complete service layer is implemented, tested, and documented:
- ✅ Keyword detection with word boundaries
- ✅ Confidence scoring (0.5/0.75/0.95)
- ✅ Skill persistence with event linking
- ✅ Case-insensitive deduplication
- ✅ Context extraction for UI display
- ✅ Integration tests demonstrating full workflow
- ✅ Comprehensive UI integration guide (`docs/guides/SKILL_INFERENCE_INTEGRATION.md`)

### Next: Parallel Tracks

**Track 1 (Phase 7-9)**: UI Integration  
**Track 2 (Phase 10)**: Expand keyword dictionary (can run in parallel with UI work)  
**Track 3 (Phase 11)**: Add soft skills detection (can run in parallel with UI work)  

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
│ TechnologyKeywords (~230 entries after Phase 10)    │
│  TECHNICAL SKILLS:                                 │
│  ├─ Backend (32): Go, Python, Ruby, Deno, Bun...  │
│  ├─ Frontend (38): React, Vue, Astro, Remix...    │
│  ├─ Database (19): PostgreSQL, MongoDB, Redis...  │
│  ├─ DevOps (24): Docker, Kubernetes, Terraform... │
│  ├─ Cloud (23): AWS, GCP, Azure, Lambda...        │
│  ├─ Mobile (10): iOS, Swift, React Native...      │
│  ├─ Tooling (19): Git, GraphQL, REST, gRPC...     │
│  ├─ Testing (15): Jest, Cypress, Pytest...        │
│  ├─ Build (12): Maven, Gradle, npm, yarn...       │
│  ├─ ML/Data (19): TensorFlow, Spark, Airflow...   │
│  ├─ Monitoring (7): Splunk, ELK, Sentry...        │
│  ├─ Documentation (6): Swagger, OpenAPI...        │
│  └─ OS (7): Linux, Ubuntu, macOS...               │
│                                                     │
│  SOFT SKILLS (Phase 11):                           │
│  ├─ Leadership: lead, manage, strategic...         │
│  ├─ Communication: present, document, explain...   │
│  ├─ Collaboration: team, cross-functional...       │
│  ├─ Problem Solving: debug, analyze, optimize...   │
│  ├─ Project Management: plan, deliver, roadmap...  │
│  └─ Architecture: design, scalable, distributed... │
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
**Phase 0-6 Complete**: ~9.5 hours  
**Phase 7-9 Remaining**: 3-6 hours  
**Phase 10-11 Enhancements**: 8-11 hours  
**Total**: 20.5-26.5 hours

**Additions**:
- +2-3h: Phase 10 (expanded keyword dictionary)
- +6-8h: Phase 11 (soft skills detection with CompetencyCategory integration)

### Decision 4: Soft Skills via CompetencyCategory

**Decision**: Extend `CompetencyCategory` enum (not `Skill.Category`) for soft skills.

**Rationale**:
- Soft skills are **competency areas** (leadership, communication), not technical tools
- CompetencyCategory already used by Facts system for CV generation
- Two existing soft skills: `CompetencyLeadership`, `CompetencyMentoring`
- Adding 5 more creates comprehensive soft skill detection
- Integrates with existing fact extraction and competency inference

**Implementation**: Phase 11 extends `internal/constants/constants.go` with 5 new categories.

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
5. **Phase 10**: Expands to 230 technologies across 14 categories
6. **Phase 11**: Detects soft skills (leadership, communication, etc.) via CompetencyCategory system

### Why This Matters for CV Generation

From our analysis:
- CV "Core Competencies" currently uses generic categories from Facts
- The Skills table has rich data (Level, YearsUsed) but isn't used in CV generation
- Task 48 will use Skills table for professional CV skills sections
- **This task populates the Skills table automatically**
- **Phase 11 enhances Facts with soft skill competencies for CV strengths sections**

---

## Technology Dictionary Requirements

Based on user preference:
- **Phase 1**: ~150 comprehensive technology keywords across 7 categories ✅ Complete
- **Phase 10**: Expand to ~230 keywords across 14 categories
- **Phase 11**: Add ~50 soft skill keyword patterns across 5 competency categories
- Organized by category for easy maintenance
- File-based dictionary that can be internally updated
- Supports aliases (e.g., "golang" → "Go", "k8s" → "Kubernetes")

### Dictionary Coverage (After Phase 10-11)

```
Total: ~280 detection patterns across 19 categories

TECHNICAL SKILLS (~230 keywords after Phase 10):

├─ Backend (32 keywords)
│  Go, Python, Ruby, Java, Node.js, PHP, Rust, Scala, Deno, Bun, Hono, etc.
│
├─ Frontend (38 keywords)
│  React, Vue, Angular, TypeScript, Next.js, Astro, Remix, SolidJS, Qwik, etc.
│
├─ Database (19 keywords)
│  PostgreSQL, MySQL, MongoDB, Redis, Elasticsearch, Cassandra, DynamoDB, etc.
│
├─ DevOps (24 keywords)
│  Docker, Kubernetes, Terraform, Ansible, Jenkins, CircleCI, GitHub Actions, etc.
│
├─ Cloud (23 keywords)
│  AWS, GCP, Azure, Lambda, S3, EC2, CloudFormation, Heroku, Vercel, etc.
│
├─ Mobile (10 keywords)
│  iOS, Android, Swift, Kotlin, React Native, Flutter, Xamarin, Ionic, etc.
│
├─ Tooling (19 keywords)
│  Git, GitHub, GraphQL, REST, gRPC, Kafka, RabbitMQ, OAuth, JWT, etc.
│
├─ Testing (15 keywords - NEW Phase 10)
│  Jest, Mocha, Chai, Jasmine, Pytest, JUnit, RSpec, Cypress, Selenium, Playwright, etc.
│
├─ Build Tools (12 keywords - NEW Phase 10)
│  Maven, Gradle, Make, Bazel, npm, yarn, pnpm, pip, poetry, Cargo, CMake, etc.
│
├─ Machine Learning (10 keywords - NEW Phase 10)
│  TensorFlow, PyTorch, scikit-learn, Pandas, NumPy, Keras, Jupyter, OpenAI, LangChain, etc.
│
├─ Data Engineering (9 keywords - NEW Phase 10)
│  Apache Spark, Airflow, Databricks, Snowflake, dbt, Hadoop, Hive, Presto, Redshift, etc.
│
├─ Monitoring/Observability (7 keywords - NEW Phase 10)
│  Splunk, ELK Stack, Kibana, Logstash, Sentry, PagerDuty, Honeycomb, etc.
│
├─ Documentation (6 keywords - NEW Phase 10)
│  Swagger, OpenAPI, Redoc, Docusaurus, MkDocs, Sphinx, etc.
│
└─ Operating Systems (7 keywords - NEW Phase 10)
   Linux, Unix, Ubuntu, Debian, CentOS, macOS, Windows Server, etc.

SOFT SKILLS (~50 keyword patterns after Phase 11):

├─ Leadership (10 patterns)
│  lead, manage, direct, strategic, vision, initiative, guide, empower, delegate, etc.
│
├─ Communication (12 patterns)
│  present, document, explain, write, articulate, stakeholder, meeting, report, clarify, brief, etc.
│
├─ Collaboration (9 patterns)
│  collaborate, team, cross-functional, partner, coordinate, facilitate, align, joint, etc.
│
├─ Problem Solving (10 patterns)
│  debug, analyze, troubleshoot, investigate, diagnose, optimize, fix, resolve, identify, etc.
│
├─ Project Management (7 patterns)
│  plan, estimate, schedule, deliver, milestone, sprint, roadmap, prioritize, etc.
│
└─ Architecture (6 patterns)
   architect, design, scalable, distributed, microservices, pattern, infrastructure, platform, etc.
```

---

## Files Created (Service Layer - Phase 0-6) ✅

### Service Layer (Complete)
- [x] `internal/service/career/technology/keywords.go` - Technology dictionary (150 entries)
- [x] `internal/service/career/technology/keywords_test.go` - Dictionary tests (33 specs)
- [x] `internal/service/career/skillinference/inference.go` - Service interface
- [x] `internal/service/career/skillinference/detector.go` - Implementation (464 lines)
- [x] `internal/service/career/skillinference/inference_test.go` - Interface tests (9 specs)
- [x] `internal/service/career/skillinference/detector_test.go` - Detection tests (21 specs)
- [x] `internal/service/career/skillinference/confidence_test.go` - Confidence tests (28 specs)
- [x] `internal/service/career/skillinference/persistence_test.go` - Persistence tests (17 specs)
- [x] `internal/service/career/skillinference/integration_test.go` - Integration tests (4 specs)

### Documentation (Complete)
- [x] `docs/guides/SKILL_INFERENCE_INTEGRATION.md` - Comprehensive UI integration guide

### UI Layer (To Create - Phase 7-9)
- [ ] `internal/cli/screens/burst_management/modals/skill_suggestion_modal.go` - Skill suggestion modal
- [ ] `internal/cli/screens/burst_management/modals/skill_suggestion_modal_test.go` - Modal tests
- [ ] Update `docs/features/SKILL_INFERENCE.md` - Comprehensive feature guide (if needed)
- [ ] Update `docs/SKILLS_GUIDE.md` - Add skill inference section (if needed)
- [ ] Update `docs/workflows/MANAGE_SKILLS_WORKFLOW.md` - Add inference workflow (if needed)

### Phase 10 Files (To Modify)
- [ ] `internal/service/career/technology/keywords.go` - Add 74 new keywords (8 categories)
- [ ] `docs/guides/SKILL_INFERENCE_INTEGRATION.md` - Update with expanded dictionary

### Phase 11 Files (To Create/Modify)
- [ ] `internal/constants/constants.go` - Add 5 soft skill competency constants
- [ ] `internal/service/career/classification/classifier.go` - Add soft skill keyword lists
- [ ] `internal/service/career/burstfact/classifier.go` - Extend competency inference
- [ ] `internal/cli/uikit/selectors/category_selector.go` - Add soft skill categories
- [ ] `internal/service/career/cv/profile_inference.go` - Add soft skill strength mappings
- [ ] `internal/service/career/classification/soft_skills_test.go` - NEW: Keyword tests
- [ ] `internal/service/career/burstfact/soft_skills_classifier_test.go` - NEW: Inference tests

---

## Files to Modify

### Intent Layer (Subdirectory Structure - Phase 7-9)
- [ ] `internal/cli/intents/burst_management/constants.go` - Add skill suggestion states
- [ ] `internal/cli/intents/burst_management/messages.go` - Add SkillSuggestionsLoadedMsg
- [ ] `internal/cli/intents/burst_management/helpers.go` - Add inferSkillsFromBurst command
- [ ] `internal/cli/intents/burst_management/handlers.go` - Add handlers for skill suggestions
- [ ] `internal/cli/intents/skillsmanagement/constants.go` - Add StateInferring, StateSuggestionReview
- [ ] `internal/cli/intents/skillsmanagement/messages.go` - Add SkillSuggestionsLoadedMsg
- [ ] `internal/cli/intents/skillsmanagement/handlers.go` - Add "i" key handler

### UI Layer (Generic Modal - Phase 7-9)
- [ ] `internal/cli/screens/burst_management/modals/suggestion_review_modal.go` - Make generic for bursts AND skills

---

## Implementation Timeline

| Phase | Description | Time | Cumulative |
|-------|-------------|------|------------|
| 0-6 | Service Layer ✅ COMPLETE | ~9.5h | 9.5h |
| 7-9 | UI Integration (can run parallel with 10-11) | 3-6h | 12.5-15.5h |
| **10** | **Expand Keywords (150→230) - 8 new categories** | **2-3h** | **14.5-18.5h** |
| **11** | **Add Soft Skills Detection - 5 competency categories** | **6-8h** | **20.5-26.5h** |

**Total**: 20.5-26.5 hours (full system with enhancements)

**Parallel Execution**: Phases 7-9, 10, and 11 can be worked on simultaneously:
- **UI Track**: Phase 7-9 (3-6 hours)
- **Keywords Track**: Phase 10 (2-3 hours)
- **Soft Skills Track**: Phase 11 (6-8 hours)
- **Wall Clock Time**: max(3-6h, 2-3h, 6-8h) = 6-8 hours if fully parallelized

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

## Phase 10: Expand Keyword Dictionary (150 → 230 Keywords)

**Goal**: Add 74 new technology keywords across 8 categories to improve detection coverage from ~70% to ~85-90%

**Status**: ⏳ PENDING

**Time Estimate**: 2-3 hours

**Prerequisites**: Phase 1-6 complete (keyword detection infrastructure exists)

**Can Run in Parallel With**: Phase 7-9 (UI work), Phase 11 (soft skills)

---

### Subphase 10.1: Add Testing Frameworks (15 keywords)

**TDD Workflow**:
1. **RED**: Verify current tests pass with 150 keywords baseline
2. **GREEN**: Add 15 testing framework keywords to `keywords.go`:
   ```go
   // Testing Frameworks (15)
   {"jest", "Jest", "testing"},
   {"mocha", "Mocha", "testing"},
   {"chai", "Chai", "testing"},
   {"jasmine", "Jasmine", "testing"},
   {"pytest", "Pytest", "testing"},
   {"junit", "JUnit", "testing"},
   {"testng", "TestNG", "testing"},
   {"rspec", "RSpec", "testing"},
   {"cypress", "Cypress", "testing"},
   {"selenium", "Selenium", "testing"},
   {"playwright", "Playwright", "testing"},
   {"webdriverio", "WebDriverIO", "testing"},
   {"cucumber", "Cucumber", "testing"},
   {"postman", "Postman", "testing"},
   {"insomnia", "Insomnia", "testing"},
   ```

**Files to Modify**:
- `internal/service/career/technology/keywords.go` (add after line 191)

**Test Verification**:
```bash
make test-suite SUITE=./internal/service/career/technology/...
# Should pass with 33 tests (no new tests needed - generic detection)
```

**Commits** (1):
- [ ] `feat(service): add 15 testing framework keywords to dictionary`

**Time**: 20 minutes

---

### Subphase 10.2: Add Build Tools (12 keywords)

**TDD Workflow**:
1. **GREEN**: Add 12 build tool keywords:
   ```go
   // Build Tools & Package Managers (12)
   {"maven", "Maven", "build"},
   {"gradle", "Gradle", "build"},
   {"make", "Make", "build"},
   {"bazel", "Bazel", "build"},
   {"npm", "npm", "build"},
   {"yarn", "Yarn", "build"},
   {"pnpm", "pnpm", "build"},
   {"pip", "pip", "build"},
   {"poetry", "Poetry", "build"},
   {"bundler", "Bundler", "build"},
   {"cargo", "Cargo", "build"},
   {"cmake", "CMake", "build"},
   ```

**Commits** (1):
- [ ] `feat(service): add 12 build tool keywords to dictionary`

**Time**: 15 minutes

---

### Subphase 10.3: Add ML/Data Engineering (19 keywords)

**TDD Workflow**:
1. **GREEN**: Add ML/Data keywords:
   ```go
   // Machine Learning & Data Science (10)
   {"tensorflow", "TensorFlow", "ml"},
   {"pytorch", "PyTorch", "ml"},
   {"scikit-learn", "scikit-learn", "ml"},
   {"sklearn", "scikit-learn", "ml"},
   {"pandas", "Pandas", "ml"},
   {"numpy", "NumPy", "ml"},
   {"keras", "Keras", "ml"},
   {"jupyter", "Jupyter", "ml"},
   {"openai", "OpenAI", "ml"},
   {"langchain", "LangChain", "ml"},
   
   // Data Engineering (9)
   {"apache spark", "Apache Spark", "data"},
   {"spark", "Apache Spark", "data"},
   {"airflow", "Apache Airflow", "data"},
   {"databricks", "Databricks", "data"},
   {"snowflake", "Snowflake", "data"},
   {"dbt", "dbt", "data"},
   {"hadoop", "Hadoop", "data"},
   {"hive", "Hive", "data"},
   {"presto", "Presto", "data"},
   ```

**Commits** (1):
- [ ] `feat(service): add 19 ML and data engineering keywords to dictionary`

**Time**: 25 minutes

---

### Subphase 10.4: Add Monitoring/Documentation/OS (20 keywords)

**TDD Workflow**:
1. **GREEN**: Add infrastructure keywords:
   ```go
   // Monitoring & Observability (7)
   {"splunk", "Splunk", "monitoring"},
   {"elk stack", "ELK Stack", "monitoring"},
   {"kibana", "Kibana", "monitoring"},
   {"logstash", "Logstash", "monitoring"},
   {"sentry", "Sentry", "monitoring"},
   {"pagerduty", "PagerDuty", "monitoring"},
   {"honeycomb", "Honeycomb", "monitoring"},
   
   // Documentation Tools (6)
   {"swagger", "Swagger", "documentation"},
   {"openapi", "OpenAPI", "documentation"},
   {"redoc", "Redoc", "documentation"},
   {"docusaurus", "Docusaurus", "documentation"},
   {"mkdocs", "MkDocs", "documentation"},
   {"sphinx", "Sphinx", "documentation"},
   
   // Operating Systems (7)
   {"linux", "Linux", "os"},
   {"unix", "Unix", "os"},
   {"ubuntu", "Ubuntu", "os"},
   {"debian", "Debian", "os"},
   {"centos", "CentOS", "os"},
   {"macos", "macOS", "os"},
   {"windows server", "Windows Server", "os"},
   ```

**Commits** (1):
- [ ] `feat(service): add 20 monitoring, documentation, and OS keywords`

**Time**: 30 minutes

---

### Subphase 10.5: Add Modern Frameworks (8 keywords)

**TDD Workflow**:
1. **GREEN**: Add newer runtime/framework keywords:
   ```go
   // Modern Runtimes & Frameworks (8)
   {"deno", "Deno", "backend"},
   {"bun", "Bun", "backend"},
   {"astro", "Astro", "frontend"},
   {"remix", "Remix", "frontend"},
   {"solidjs", "SolidJS", "frontend"},
   {"qwik", "Qwik", "frontend"},
   {"fresh", "Fresh", "frontend"},
   {"hono", "Hono", "backend"},
   ```

**Commits** (1):
- [ ] `feat(service): add 8 modern runtime and framework keywords`

**Time**: 15 minutes

---

### Subphase 10.6: Update Documentation

**TDD Workflow**:
1. **REFACTOR**: Update docs to reflect new coverage
   - Update `keywords.go` header comment: "~100 keywords" → "~230 keywords"
   - Update `docs/guides/SKILL_INFERENCE_INTEGRATION.md` with new categories
   - Update section counts in comments

**Commits** (1):
- [ ] `docs(service): update skill inference docs for 230 keyword dictionary`

**Time**: 20 minutes

---

### Phase 10 Acceptance Criteria
- [ ] 74 new keywords added (150 → 224 total technical keywords)
- [ ] 8 new categories: testing, build, ml, data, monitoring, documentation, os, modern
- [ ] All existing tests pass (33 technology tests)
- [ ] Manual verification: New keywords detected in test events
- [ ] Documentation updated with new keyword counts
- [ ] No breaking changes to existing API

**Phase 10 Total**: 2-3 hours

---

## Phase 11: Add Soft Skills Detection

**Goal**: Extend system to detect soft skills from event text using CompetencyCategory system

**Status**: ⏳ PENDING

**Time Estimate**: 6-8 hours

**Prerequisites**: Phase 10 complete (or can run in parallel), understanding of `CompetencyCategory` system

**Can Run in Parallel With**: Phase 7-9 (UI work), Phase 10 (keywords)

---

### Context: Soft Skills Architecture

**Approach**: Extend `CompetencyCategory` enum (used by Facts/CV system) to include soft skills

**Current Competencies** (2 soft skills exist):
- `CompetencyTechnical` - Technical skills ✅
- `CompetencyLeadership` ✅ (soft skill - already exists)
- `CompetencyProduct` - Product thinking ✅
- `CompetencyConsulting` - Client advisory ✅
- `CompetencyResearch` - Investigation ✅
- `CompetencyMentoring` ✅ (soft skill - already exists)

**New Soft Skills to Add** (5 categories):
- `CompetencyCommunication` - Presenting, documenting, explaining
- `CompetencyCollaboration` - Teamwork, cross-functional work
- `CompetencyProblemSolving` - Debugging, troubleshooting, analysis
- `CompetencyProjectManagement` - Planning, delivery, roadmaps
- `CompetencyArchitecture` - System design, scalability, technical decisions

**Why This Matters**:
- Soft skills detected from event text automatically
- Integrated with fact extraction for CV generation
- Facts tagged with soft skill competencies
- CV "Core Competencies" includes soft skills
- Profile inference generates soft skill strength descriptions

---

### Subphase 11.1: Extend Domain Constants (30 min)

**TDD Workflow**:
1. **RED**: Test `IsValidCompetencyCategory("communication")` → Fails
2. **GREEN**: Add 5 new `CompetencyCategory` constants
3. **GREEN**: Update `AllCompetencyCategories()`
4. **GREEN**: Update `IsValidCompetencyCategory()`
5. **GREEN**: Update `GetCompetencyDescription()`
6. **REFACTOR**: Organize constants by type (technical vs soft)

**Files to Modify**:
```go
// internal/constants/constants.go (lines 65-119)

// Add after line 70 (after CompetencyMentoring):
const (
    // ... existing constants ...
    CompetencyMentoring      CompetencyCategory = "mentoring"
    
    // NEW: Soft Skills (Phase 11)
    CompetencyCommunication     CompetencyCategory = "communication"
    CompetencyCollaboration     CompetencyCategory = "collaboration"
    CompetencyProblemSolving    CompetencyCategory = "problem-solving"
    CompetencyProjectManagement CompetencyCategory = "project-management"
    CompetencyArchitecture      CompetencyCategory = "architecture"
)

// Update AllCompetencyCategories() function:
func AllCompetencyCategories() []CompetencyCategory {
    return []CompetencyCategory{
        CompetencyTechnical,
        CompetencyLeadership,
        CompetencyProduct,
        CompetencyConsulting,
        CompetencyResearch,
        CompetencyMentoring,
        CompetencyCommunication,     // NEW
        CompetencyCollaboration,     // NEW
        CompetencyProblemSolving,    // NEW
        CompetencyProjectManagement, // NEW
        CompetencyArchitecture,      // NEW
    }
}

// Update IsValidCompetencyCategory() - add new cases

// Update GetCompetencyDescription() - add new descriptions:
case CompetencyCommunication:
    return "Communication and documentation skills"
case CompetencyCollaboration:
    return "Cross-functional collaboration and teamwork"
case CompetencyProblemSolving:
    return "Analytical and problem-solving abilities"
case CompetencyProjectManagement:
    return "Project planning and delivery management"
case CompetencyArchitecture:
    return "System architecture and technical design"
```

**Commits** (1):
- [ ] `feat(constants): add 5 soft skill competency categories`

**Time**: 30 minutes

---

### Subphase 11.2: Add Soft Skills Keywords (45 min)

**TDD Workflow**:
1. **RED**: Test keyword detection for communication keywords → Not detected
2. **GREEN**: Add keyword lists to classifier
3. **REFACTOR**: Organize keywords by competency

**Files to Modify**:
```go
// internal/service/career/classification/classifier.go (add after line 61)

var (
    // ... existing keyword maps ...
    
    // NEW: Soft Skills Keywords (Phase 11)
    communicationKeywords = []string{
        "communicate", "present", "document", "explain", "write",
        "articulate", "stakeholder", "meeting", "update", "report",
        "clarify", "brief",
    }
    
    collaborationKeywords = []string{
        "collaborate", "team", "cross-functional", "partner",
        "coordinate", "facilitate", "align", "together", "joint",
    }
    
    problemSolvingKeywords = []string{
        "solve", "debug", "analyze", "troubleshoot", "investigate",
        "diagnose", "optimize", "fix", "resolve", "identify",
    }
    
    projectManagementKeywords = []string{
        "plan", "estimate", "schedule", "deliver", "milestone",
        "sprint", "roadmap", "prioritize",
    }
    
    architectureKeywords = []string{
        "architect", "design", "scalable", "distributed",
        "microservices", "pattern", "infrastructure", "platform",
    }
)
```

**Commits** (1):
- [ ] `feat(service): add soft skill keyword lists to classifier`

**Time**: 45 minutes

---

### Subphase 11.3: Extend Competency Inference (2 hours)

**TDD Workflow**:
1. **RED**: Test fact extraction with communication keywords → No competency assigned
2. **GREEN**: Update `InferCompetencies()` to detect 5 new soft skills
3. **GREEN**: Add confidence scoring for soft skills
4. **REFACTOR**: Clean up competency detection logic

**Files to Modify**:
```go
// internal/service/career/burstfact/classifier.go (extend lines 177-251)

func (c *Classifier) InferCompetencies(text string, category string) []string {
    competencies := make(map[string]bool)
    lowerText := strings.ToLower(text)
    
    // ... existing technical, leadership, mentoring detection ...
    
    // NEW: Communication detection
    if c.hasKeywords(lowerText, classification.CommunicationKeywords()) {
        competencies[constants.CompetencyCommunication] = true
    }
    
    // NEW: Collaboration detection
    if c.hasKeywords(lowerText, classification.CollaborationKeywords()) {
        competencies[constants.CompetencyCollaboration] = true
    }
    
    // NEW: Problem Solving detection
    if c.hasKeywords(lowerText, classification.ProblemSolvingKeywords()) {
        competencies[constants.CompetencyProblemSolving] = true
    }
    
    // NEW: Project Management detection
    if c.hasKeywords(lowerText, classification.ProjectManagementKeywords()) {
        competencies[constants.CompetencyProjectManagement] = true
    }
    
    // NEW: Architecture detection
    if c.hasKeywords(lowerText, classification.ArchitectureKeywords()) {
        competencies[constants.CompetencyArchitecture] = true
    }
    
    return mapToSlice(competencies)
}
```

**Commits** (1):
- [ ] `feat(service): extend competency inference for 5 soft skills`

**Time**: 2 hours

---

### Subphase 11.4: Update UI Category Selector (30 min)

**TDD Workflow**:
1. **RED**: Test category selector with "communication" → Not in options
2. **GREEN**: Add soft skills to `AllowedCategories`
3. **GREEN**: Add descriptions to `GetCategoryDescription()`
4. **REFACTOR**: Group categories by type

**Files to Modify**:
```go
// internal/cli/uikit/selectors/category_selector.go (lines 11-19, 135-148)

var AllowedCategories = map[string]bool{
    // ... existing categories ...
    
    // NEW: Soft Skills (Phase 11)
    "communication":       true,
    "collaboration":       true,
    "problem-solving":     true,
    "project-management":  true,
    "architecture":        true,
}

// Update GetCategoryDescription():
case "communication":
    return "Communication and documentation skills"
case "collaboration":
    return "Cross-functional collaboration and teamwork"
case "problem-solving":
    return "Analytical and problem-solving abilities"
case "project-management":
    return "Project planning and delivery management"
case "architecture":
    return "System architecture and technical design"
```

**Commits** (1):
- [ ] `feat(ui): add soft skill categories to fact category selector`

**Time**: 30 minutes

---

### Subphase 11.5: Extend Profile Inference (30 min)

**TDD Workflow**:
1. **RED**: Test profile inference with soft skills → No strength text
2. **GREEN**: Add soft skill mappings to `categoryStrengthMapping`
3. **REFACTOR**: Organize strength descriptions

**Files to Modify**:
```go
// internal/service/career/cv/profile_inference.go (lines 25-35)

var categoryStrengthMapping = map[string]string{
    // ... existing mappings ...
    
    // NEW: Soft Skills (Phase 11)
    "communication":       "Strong communication and documentation skills with stakeholder engagement",
    "collaboration":       "Effective cross-functional collaboration and team leadership",
    "problem-solving":     "Strong analytical and problem-solving abilities with systematic approach",
    "project-management":  "Project planning and delivery management with milestone tracking",
    "architecture":        "System architecture and technical design expertise with scalability focus",
}
```

**Commits** (1):
- [ ] `feat(service): add soft skill strength mappings for CV generation`

**Time**: 30 minutes

---

### Subphase 11.6: Comprehensive Testing (2-3 hours)

**TDD Workflow**:
1. Unit tests for each soft skill competency
2. Integration tests for fact extraction with soft skills
3. E2E tests for full workflow

**New Test Files**:
```go
// internal/service/career/classification/soft_skills_test.go (~200 lines)
var _ = Describe("Soft Skills Keyword Detection", func() {
    Describe("Communication Keywords", func() {
        It("should detect presentation mentions", func() {
            text := "Presented technical design to stakeholders"
            result := classifier.DetectCommunication(text)
            Expect(result).To(BeTrue())
        })
        // ... 10-15 more tests
    })
    
    // Repeat for each soft skill category
})

// internal/service/career/burstfact/soft_skills_classifier_test.go (~300 lines)
var _ = Describe("Soft Skills Competency Inference", func() {
    It("should infer communication competency from presentation text", func() {
        text := "Presented quarterly roadmap to executive team"
        competencies := classifier.InferCompetencies(text, "")
        Expect(competencies).To(ContainElement("communication"))
    })
    
    It("should infer multiple soft skills from complex text", func() {
        text := "Led cross-functional team to debug production issues and delivered fixes ahead of schedule"
        competencies := classifier.InferCompetencies(text, "")
        Expect(competencies).To(ContainElement("leadership"))
        Expect(competencies).To(ContainElement("collaboration"))
        Expect(competencies).To(ContainElement("problem-solving"))
    })
    
    // ... 50-75 more tests
})
```

**Test Coverage**:
- 10-15 tests per soft skill category = 50-75 tests
- Integration tests: 5-10 tests
- E2E tests: 3-5 tests

**Commits** (3):
- [ ] `test(service): add soft skill keyword detection tests`
- [ ] `test(service): add soft skill competency inference tests`
- [ ] `test(e2e): add soft skill fact extraction workflow tests`

**Time**: 2-3 hours

---

### Subphase 11.7: Update Documentation (30 min)

**Files to Update**:
- `docs/guides/SKILL_INFERENCE_INTEGRATION.md` - Add soft skills section
- `docs/SKILLS_GUIDE.md` - Mention soft skill detection
- `docs/workflows/MANAGE_SKILLS_WORKFLOW.md` - Update with soft skill examples

**Documentation Sections to Add**:
```markdown
## Soft Skills Detection (Phase 11)

The skill inference system detects soft skills through the CompetencyCategory system:

### Soft Skill Categories

- **Leadership**: "led team of 5 engineers", "managed project deliverables"
- **Communication**: "presented to executive team", "documented architecture decisions"
- **Collaboration**: "worked with cross-functional teams", "partnered with design"
- **Problem Solving**: "debugged production issue", "optimized query performance"
- **Project Management**: "planned sprint roadmap", "delivered features on schedule"
- **Architecture**: "designed scalable microservices", "architected distributed system"

### How It Works

1. Event text analyzed for soft skill keywords
2. Competency categories assigned to facts
3. Facts used in CV generation for "Core Competencies" section
4. Profile inference generates soft skill strength descriptions
```

**Commits** (1):
- [ ] `docs(guides): document soft skill detection capabilities`

**Time**: 30 minutes

---

### Phase 11 Acceptance Criteria
- [ ] 5 new soft skill competency categories added to constants
- [ ] 40-50 soft skill keyword patterns added to classifier
- [ ] Competency inference detects soft skills from event text
- [ ] UI category selector includes all 5 soft skills
- [ ] Profile inference generates soft skill strength descriptions
- [ ] 50-75 new tests passing (keyword detection + inference)
- [ ] Documentation updated with soft skills examples
- [ ] Fact extraction E2E with soft skills works
- [ ] CV generation includes soft skill competencies

**Phase 11 Total**: 6-8 hours

---

## Expected UX Examples

### After burst confirmation (automatic):
```
┌─────────────────────────────────────────────────┐
│ Skill Suggestions (from 5 burst events)         │
│                                                 │
│ Suggestion 1 of 5                               │
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
│                                                 │
│ Also detected (Phase 10-11):                    │
│ • Cypress (testing) - 85%                       │
│ • Jest (testing) - 90%                          │
│ • Leadership (soft skill) - 80%                 │
│ • Problem Solving (soft skill) - 75%            │
└─────────────────────────────────────────────────┘
```

### In burst detail modal after "i" key (manual):
```
┌─────────────────────────────────────────────────┐
│ Analyzing Burst Events                          │
│                                                 │
│ ⏳ Detecting technical and soft skills...       │
│                                                 │
│ This may take a moment...                       │
└─────────────────────────────────────────────────┘
```

### In ManageSkills after pressing "i" (manual):
```
┌─────────────────────────────────────────────────┐
│ Inferring Skills from All Events                │
│                                                 │
│ ⏳ Analyzing 42 events...                        │
│   - Technical skills: 230 keyword patterns      │
│   - Soft skills: 5 competency categories        │
│                                                 │
│ This may take a few moments...                  │
└─────────────────────────────────────────────────┘
```

---

## Acceptance Criteria

### Phase 0-9 (Original)
- [ ] Technology dictionary has ~100 comprehensive entries organized by category ✅
- [ ] Skill inference uses word boundary regex (prevents partial matches) ✅
- [ ] Confidence scoring differentiates high/medium/low usage patterns ✅
- [ ] Deduplication merges same skill from multiple events ✅
- [ ] Created skills are linked to source events via junction table ✅
- [ ] LastUsed is set to most recent event date ✅
- [ ] **Automatic trigger**: Integration into burst confirmation works (after facts)
- [ ] **Manual trigger 1**: "i" key in burst detail modal works
- [ ] **Manual trigger 2**: "i" key in ManageSkills list works
- [ ] Generic modal handles both burst and skill suggestions
- [ ] Context snippets show actual technology usage (~80 chars) ✅
- [ ] Empty suggestions handled gracefully (no error) ✅

### Phase 10 Criteria (Expanded Keywords)
- [ ] Technology dictionary expanded to 230 keywords (74 new entries)
- [ ] 8 new categories: testing, build, ml, data, monitoring, documentation, os, modern
- [ ] All existing tests pass with expanded dictionary
- [ ] Manual verification: New keywords detected in test events
- [ ] Documentation updated with new keyword counts
- [ ] No breaking changes to existing API

### Phase 11 Criteria (Soft Skills)
- [ ] 5 new soft skill competency categories added (communication, collaboration, problem-solving, project-management, architecture)
- [ ] 40-50 soft skill keyword patterns integrated into classifier
- [ ] Competency inference detects soft skills from event text
- [ ] UI category selector supports all 5 soft skill categories
- [ ] Profile inference generates soft skill strength descriptions
- [ ] 50-75 new tests covering soft skill detection and inference
- [ ] Documentation includes soft skills examples and workflow
- [ ] Fact extraction E2E with soft skills works end-to-end
- [ ] CV generation includes soft skill competencies in Core Competencies section

---

## Documentation

See additional documentation:
- **UI Design**: `docs/design/SKILL_INFERENCE_UI_DESIGN.md` - Complete UI analysis and reusability assessment
- **Feature Guide**: `docs/features/SKILL_INFERENCE.md` - Technical reference (to be created)
- **User Guide**: `docs/SKILLS_GUIDE.md` - End-user documentation (to be updated with Phase 10-11)
- **Workflow**: `docs/workflows/MANAGE_SKILLS_WORKFLOW.md` - Operational procedures (to be updated with Phase 10-11)
- **Integration Guide**: `docs/guides/SKILL_INFERENCE_INTEGRATION.md` - Comprehensive UI integration guide (to be updated with Phase 10-11)

---

## Rollback Plan

1. Remove service files:
   - `internal/service/career/technology_keywords.go`
   - `internal/service/career/technology_keywords_test.go`
   - `internal/service/career/skill_inference.go`
   - `internal/service/career/skill_inference_impl.go`
   - `internal/service/career/skill_inference_test.go`
2. Revert Phase 11 changes:
   - `internal/constants/constants.go` (remove soft skill constants)
   - `internal/service/career/classification/classifier.go` (remove soft skill keywords)
   - `internal/service/career/burstfact/classifier.go` (revert competency inference)
   - Remove `internal/service/career/classification/soft_skills_test.go`
   - Remove `internal/service/career/burstfact/soft_skills_classifier_test.go`
3. Revert intent changes:
   - `internal/cli/intents/burst_management/`
   - `internal/cli/intents/skillsmanagement/`
4. Revert modal changes:
   - `internal/cli/screens/burst_management/modals/suggestion_review_modal.go`
5. Run `make check-compliance`
6. All tests should pass

---

## Dependencies

- Task 45 (burst detection UI) - for similar pattern understanding ✅ Complete
- Skill domain model and repository ✅ Already exists
- SuggestionReviewModal ✅ Already exists (will be made generic)
- CompetencyCategory system ✅ Already exists (will be extended in Phase 11)

---

## Next Steps After Completion

- Task 48: Enhanced CV Skills Section (will use skills created by this task)
- Task 54: Ollama LLM Integration (hybrid keyword + semantic inference)
- Task 44: Wire Enhanced Bullet Generator (independent, can be done in parallel)

---

## Service Layer Achievements (Phases 0-6)

### Implementation Highlights

#### 1. Word Boundary Detection
**Problem**: Substring matching caused false positives ("goal" matched "go")  
**Solution**: `\bkeyword\b` regex with case-insensitive matching  
**Result**: Zero false positives in 111 tests

#### 2. Proximity-Based Confidence Scoring  
**Problem**: Distant word matches inflated confidence  
**Solution**: Pattern words must be within 3 words of each other  
**Example**: 
```
Text: "Built API with PostgreSQL after working with MongoDB"
PostgreSQL: 0.95 (matches "built...with...postgresql" - close)
MongoDB:    0.75 (matches "working with mongodb" - correct pattern)
```

#### 3. Three-Tier Confidence System
- **High (0.95)**: Active usage ("built with X", "using X", "X developer")
- **Medium (0.75)**: Passive mention ("worked with X", "X project")
- **Low (0.5)**: Simple keyword presence

#### 4. Repository Pattern with Full Test Coverage
- Clean service/repository separation
- Mock repositories for isolated testing
- 112 tests covering all scenarios
- Integration tests demonstrating complete workflow

### Files Statistics

```
Production Code: 1,020 lines → 1,500-1,800 lines (after Phase 10-11)
├── internal/service/career/technology/keywords.go (224 → 350 lines)
├── internal/service/career/skillinference/inference.go (74 lines)
├── internal/service/career/skillinference/detector.go (464 lines)
├── internal/constants/constants.go (+30 lines Phase 11)
├── internal/service/career/classification/classifier.go (+80 lines Phase 11)
└── internal/service/career/burstfact/classifier.go (+100 lines Phase 11)

Test Code: 2,267 lines → 3,500-4,000 lines (after Phase 10-11)
├── internal/service/career/technology/keywords_test.go (208 lines)
├── internal/service/career/skillinference/inference_test.go (199 lines)
├── internal/service/career/skillinference/detector_test.go (450 lines)
├── internal/service/career/skillinference/confidence_test.go (507 lines)
├── internal/service/career/skillinference/persistence_test.go (457 lines)
├── internal/service/career/skillinference/integration_test.go (223 lines)
├── internal/service/career/classification/soft_skills_test.go (NEW: ~200 lines)
└── internal/service/career/burstfact/soft_skills_classifier_test.go (NEW: ~300 lines)

Documentation: 471 lines → 600-700 lines (after Phase 10-11)
└── docs/guides/SKILL_INFERENCE_INTEGRATION.md
```

### Service API

```go
// Primary interface
type SkillInferenceService interface {
    // Infer skills from events
    InferSkillsFromEvents(ctx context.Context, events []*career.Event) ([]SkillSuggestion, error)
    
    // Infer skills from specific burst
    InferSkillsFromBurst(ctx context.Context, burst *career.Burst, events []*career.Event) ([]SkillSuggestion, error)
    
    // Create skills from accepted suggestions
    CreateSkillsFromSuggestions(ctx context.Context, suggestions []SkillSuggestion) ([]*career.Skill, error)
}

// Usage example
service := skillinference.NewSkillInferenceService(skillRepo, eventRepo)
suggestions, err := service.InferSkillsFromEvents(ctx, events)
skills, err := service.CreateSkillsFromSuggestions(ctx, acceptedSuggestions)
```

### Commits (15 total → ~25-30 after Phase 10-11)

All commits follow TDD Red-Green-Refactor and include AI attribution:

**Phase 0-6 (Complete):**
```
23cc88b7 docs(guides): add comprehensive UI integration guide
b84f2033 test(service): add end-to-end integration tests
85056a1e test(service): update tests for constructor signature
5e9019d6 feat(service): implement skill persistence
2badf3ce test(service): add skill persistence tests
6a6c14f0 refactor(service): remove redundant nil check
79980c83 feat(service): implement confidence scoring
ad804b0d test(service): add confidence scoring tests
ddbd2f20 feat(service): implement keyword detection
464acef9 test(service): add keyword detection tests
03f434f6 feat(service): add service interface
bcf1c57d test(service): add service interface tests
54a942cb feat(service): add technology keyword dictionary
7dcfa753 test(service): add technology dictionary tests
81319fe5 test(intents): remove redundant skipped test
```

**Phase 10 (Pending - 6 commits):**
- feat(service): add 15 testing framework keywords
- feat(service): add 12 build tool keywords
- feat(service): add 19 ML and data engineering keywords
- feat(service): add 20 monitoring, documentation, and OS keywords
- feat(service): add 8 modern runtime and framework keywords
- docs(service): update skill inference docs for 230 keyword dictionary

**Phase 11 (Pending - 7 commits):**
- feat(constants): add 5 soft skill competency categories
- feat(service): add soft skill keyword lists to classifier
- feat(service): extend competency inference for 5 soft skills
- feat(ui): add soft skill categories to fact category selector
- feat(service): add soft skill strength mappings for CV generation
- test(service): add soft skill keyword detection and inference tests
- docs(guides): document soft skill detection capabilities

### Next: UI Integration + Enhancements (Phases 7-11)

With the service layer complete, the next phases will:

**UI Track (Phase 7-9)**: 3-6 hours
1. Create SkillSuggestionModal component
2. Integrate with burst_management intent (automatic + manual triggers)
3. Integrate with skillsmanagement intent (global manual trigger)
4. Add E2E tests with real TUI navigation

**Keywords Track (Phase 10)**: 2-3 hours
1. Add 74 new technology keywords (8 categories)
2. Update documentation with expanded coverage
3. Verify detection of new keywords

**Soft Skills Track (Phase 11)**: 6-8 hours
1. Extend CompetencyCategory with 5 soft skill categories
2. Add keyword detection for soft skills
3. Integrate with fact extraction and CV generation
4. Comprehensive testing (50-75 tests)
5. Update documentation

**Updated**: January 30, 2026 - Service layer complete, planning UI integration and enhancements (Phase 7-11)

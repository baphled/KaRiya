# Task List: Burst and Fact Extraction Feature

**PRD Reference**: `docs/features/03-burst-fact-extraction.md`

**Purpose**: Automatically group and enrich career events by detecting bursts (related event groupings) and extracting facts (inferred competencies, role fit, audience relevance, and strength signals).

**Status**: 🔄 **IN PLANNING** (Ready for Phase 1 execution)

---

## Tasks

### Phase 1: Foundation & Core Components

#### 1.0 Create Burst Domain Model and Persistence
- [x] 1.1 Define Burst struct with fields: ID, Name, Description, EventIDs (≥2 required), CreatedAt, UpdatedAt, CompetencyFocus
- [x] 1.2 Implement Burst validation rules (≥2 related CareerEvents required, no duplicate event IDs)
- [x] 1.3 Create repository interface methods for Burst (Create, GetByID, Update, Delete, List, Count)
- [x] 1.4 Implement MemoryRepository for Burst operations (thread-safe with sync.RWMutex)
- [x] 1.5 Implement SQLiteRepository for Burst persistence with schema:
  ```sql
  CREATE TABLE IF NOT EXISTS bursts (
      id TEXT PRIMARY KEY,
      name TEXT NOT NULL,
      description TEXT,
      event_ids TEXT NOT NULL,
      competency_focus TEXT,
      created_at DATETIME NOT NULL,
      updated_at DATETIME NOT NULL
  )
  ```
- [x] 1.6 Write comprehensive unit tests for Burst model validation (edge cases, boundary conditions)
- [x] 1.7 Write repository tests for all CRUD operations and filtering

#### 2.0 Create Fact Domain Model and Persistence
- [x] 2.1 Define Fact struct with fields: ID, Text, CompetencyCategories ([]string), RoleFit (Principal/EM/Staff/Senior IC), AudienceRelevance ([]string), StrengthSignal (string), SourceEventID (optional), SourceBurstID (optional), CreatedAt, UpdatedAt
- [x] 2.2 Implement Fact validation rules (must have valid references, no aspirational language, no ungrounded metrics)
- [x] 2.3 Create repository interface methods for Fact (Create, GetByID, Update, Delete, List, Count, GetBySourceEventID, GetBySourceBurstID)
- [x] 2.4 Implement MemoryRepository for Fact operations (thread-safe)
- [x] 2.5 Implement SQLiteRepository for Fact persistence with schema:
  ```sql
  CREATE TABLE IF NOT EXISTS facts (
      id TEXT PRIMARY KEY,
      text TEXT NOT NULL,
      competencies TEXT,
      role_fit TEXT,
      audience_relevance TEXT,
      strength_signal TEXT,
      source_event_id TEXT,
      source_burst_id TEXT,
      created_at DATETIME NOT NULL,
      updated_at DATETIME NOT NULL
  )
  ```
- [x] 2.6 Write comprehensive unit tests for Fact model validation
- [x] 2.7 Write repository tests for all CRUD operations and filtering

#### 3.0 Create Classification and Inference System
- [x] 3.1 Implement role fit classifier (Principal, EM, Staff, Senior IC) based on event keywords and context
- [x] 3.2 Implement audience relevance analyzer (Hiring Manager, Recruiter, Peer) based on fact type
- [x] 3.3 Implement strength signal extractor (identifies key achievements and impact indicators)
- [x] 3.4 Create inference rules engine for fact generation from events/bursts
- [x] 3.5 Implement validation for aspirational language detection (reject "will", "should", "could")
- [x] 3.6 Implement metrics validation (ensure metrics are grounded/measured, not speculative)
- [x] 3.7 Write comprehensive unit tests for all classifiers and validators (edge cases, false positives)
- [x] 3.8 Write integration tests for inference rules engine

### Phase 2: Burst Detection and Management

#### 4.0 Create Burst Detection Engine
- [x] 4.1 Implement burst detection algorithm using event text similarity and temporal proximity
- [x] 4.2 Create similarity scorer (text matching, keyword overlap, company/project matching)
- [x] 4.3 Implement temporal grouping (events within 6 months considered related)
- [x] 4.4 Create burst suggestion generation (returns list of suggested bursts with confidence scores)
- [x] 4.5 Implement user confirmation workflow (suggest burst, user confirms or rejects)
- [x] 4.6 Create service method: `SuggestBursts(ctx, eventIDs) -> []BurstSuggestion`
- [x] 4.7 Create service method: `ConfirmBurst(ctx, burst) -> error` (validates and persists)
- [x] 4.8 Create service method: `RejectBurstSuggestion(ctx, eventIDs) -> error` (records rejection to prevent re-suggesting)
- [x] 4.9 Write unit tests for similarity scoring and temporal grouping
- [x] 4.10 Write integration tests for burst suggestion workflow

#### 5.0 Create Burst Display Component
- [x] 5.1 Implement BurstListModel as BubbleTea Model for displaying burst list
- [x] 5.2 Display burst name, event count, competency focus, creation date
- [x] 5.3 Implement scrolling (up/down arrows) through burst list
- [x] 5.4 Implement expand/collapse to show related events in burst
- [x] 5.5 Implement filtering by competency focus
- [x] 5.6 Implement sorting by creation date, event count, or name
- [x] 5.7 Implement keyboard navigation (↑/↓ for bursts, Enter to view details, Space for select)
- [x] 5.8 Add visual indicators for burst health/completeness
- [x] 5.9 Write comprehensive unit tests for all interactions
- [x] 5.10 Test edge cases (empty list, single burst, large burst list)

#### 6.0 Create Burst Suggestion and Confirmation Screen
- [x] 6.1 Implement BurstSuggestionModel as BubbleTea Model for reviewing burst suggestions
- [x] 6.2 Display suggested burst with related events and confidence score
- [x] 6.3 Allow user to confirm or reject suggestion
- [x] 6.4 Show related events with text preview
- [x] 6.5 Implement optional burst name/description editing before confirmation
- [x] 6.6 Implement keyboard navigation (↑/↓ for suggestions, 'y'/'n' for confirm/reject, 'e' to edit)
- [x] 6.7 Add visual feedback for confidence score (color-coded, percentage)
- [x] 6.8 Write comprehensive unit tests for suggestion workflow
- [x] 6.9 Test edge cases (low confidence, single event suggestions, multiple suggestions)

### Phase 3: Fact Extraction and Inference

#### 7.0 Create Fact Extraction Engine
- [x] 7.1 Implement fact extraction from single CareerEvent
- [x] 7.2 Implement fact extraction from Burst (multiple related events)
- [x] 7.3 Create competency inference from event text and tags
- [x] 7.4 Create role fit inference (Principal/EM/Staff/Senior IC) from facts and context
- [x] 7.5 Create audience relevance inference (Hiring Manager, Recruiter, Peer)
- [x] 7.6 Create strength signal extraction (key achievements, impact indicators)
- [x] 7.7 Implement fact validation and filtering (no aspirational language, grounded metrics)
- [x] 7.8 Create service method: `ExtractFacts(ctx, event/burst) -> []Fact`
- [x] 7.9 Create service method: `ValidateFact(ctx, fact) -> error` (validation rules)
- [x] 7.10 Write unit tests for all extraction and inference components
- [x] 7.11 Write integration tests for complete fact extraction workflow

#### 8.0 Create Fact Display Components
- [x] 8.1 Implement FactCardComponent for displaying individual fact
- [x] 8.2 Display fact text, competencies (as badges), role fit (icon), audience relevance
- [x] 8.3 Show strength signal and source reference (event ID or burst ID)
- [x] 8.4 Implement FactListModel for displaying facts for an event/burst
- [x] 8.5 Implement scrolling through fact list
- [x] 8.6 Implement filtering by competency, role fit, or audience
- [x] 8.7 Implement sorting by creation date or relevance
- [x] 8.8 Add visual indicators for fact confidence/strength
- [x] 8.9 Write comprehensive unit tests for display components
- [x] 8.10 Test edge cases (no facts, single fact, large fact list)

#### 9.0 Create Fact Management UI
- [x] 9.1 Implement FactEditorModel for reviewing and editing extracted facts
- [x] 9.2 Allow editing fact text with validation
- [x] 9.3 Allow editing competency categories (from AllowedCompetencies)
- [x] 9.4 Allow editing role fit (Principal/EM/Staff/Senior IC)
- [x] 9.5 Allow editing audience relevance (Hiring Manager, Recruiter, Peer)
- [x] 9.6 Implement accept/reject workflow for extracted facts
- [x] 9.7 Implement keyboard navigation (Tab for fields, Enter to confirm, Escape to cancel)
- [x] 9.8 Add helpful validation error messages
- [x] 9.9 Implement undo/revert to original extracted fact
- [x] 9.10 Write comprehensive unit tests for editor interactions

### Phase 4: Integration with Existing Features

#### 10.0 Integrate Burst Suggestions with Metadata Review
- [x] 10.1 Add burst suggestion trigger after metadata clarification
- [x] 10.2 Display burst suggestions in dedicated screen after user confirms metadata
- [x] 10.3 Allow user to accept/reject each burst suggestion
- [x] 10.4 Show related events for each suggested burst
- [x] 10.5 Allow editing burst name/description before confirmation
- [x] 10.6 Persist accepted bursts to database
- [x] 10.7 Record rejected suggestions to prevent re-suggesting
- [x] 10.8 Navigate back to metadata review or home after burst workflow
- [x] 10.9 Write integration tests for metadata review → burst suggestion workflow

#### 11.0 Integrate Fact Display with Event Details
- [x] 11.1 Add facts section to event detail view
- [x] 11.2 Display extracted facts for selected event
- [x] 11.3 Show facts grouped by source (inferred from event, inferred from burst)
- [x] 11.4 Allow user to review and confirm facts
- [x] 11.5 Allow user to edit individual facts (through FactEditorModel)
- [x] 11.6 Allow user to reject facts (mark as not applicable)
- [x] 11.7 Implement keyboard navigation for fact review
- [x] 11.8 Persist fact confirmations to database
- [x] 11.9 Write integration tests for event detail → facts view

#### 12.0 Create User Confirmation and Workflow Integration
- [x] 12.1 Design workflow: Capture/Import → Metadata Review → Burst Suggestions → Fact Extraction → Confirmation
- [x] 12.2 Implement screen navigation for complete workflow
- [x] 12.3 Add progress indicator showing current step in workflow
- [x] 12.4 Allow skipping burst suggestions (user can enable/disable burst detection)
- [x] 12.5 Allow skipping fact extraction (user can enable/disable fact extraction)
- [x] 12.6 Implement "review later" option for bursts and facts
- [x] 12.7 Create home screen menu option to review pending bursts and facts
- [x] 12.8 Write integration tests for complete end-to-end workflow

### Phase 5: Testing and Documentation

#### 13.0 Comprehensive Testing Suite
- [x] 13.1 Write end-to-end tests for complete burst detection workflow
- [x] 13.2 Write end-to-end tests for fact extraction and confirmation
- [x] 13.3 Test burst detection with various event similarity scenarios (high/medium/low similarity)
- [x] 13.4 Test fact extraction with various event types and contexts
- [x] 13.5 Test inference rules with edge cases (aspirational language, ungrounded metrics)
- [x] 13.6 Test role fit classification with different event contexts
- [x] 13.7 Test audience relevance inference for different fact types
- [x] 13.8 Test integration with metadata review workflow
- [x] 13.9 Test keyboard navigation for all burst and fact screens
- [x] 13.10 Test edge cases (no similar events, low confidence suggestions, conflicting facts)
- [x] 13.11 Run race detector: `go test -race ./...` (verify 0 race conditions)
- [x] 13.12 Verify code coverage meets 80%+ threshold across all new code
- [x] 13.13 Performance test: burst detection ≤2s for ≤500 events
- [x] 13.14 Performance test: fact extraction ≤1s per event/burst

#### 14.0 Documentation and User Guidance
- [x] 14.1 Create BURST_FACT_EXTRACTION_GUIDE.md with comprehensive feature overview
- [x] 14.2 Document burst detection algorithm and how it works
- [x] 14.3 Document fact extraction rules and inference process
- [x] 14.4 Provide examples of burst suggestions (good/bad examples)
- [x] 14.5 Provide examples of extracted facts (good/bad examples)
- [x] 14.6 Document keyboard shortcuts for burst and fact screens
- [x] 14.7 Document role fit classification and audience relevance
- [x] 14.8 Update README.md with burst and fact features
- [x] 14.9 Update CLI_GUIDE.md with burst/fact workflow shortcuts
- [x] 14.10 Update CHANGELOG.md with feature description and test results
- [x] 14.11 Create troubleshooting guide for common burst/fact issues
- [x] 14.12 Document allowed competencies and their definitions

---

## Implementation Guidelines

### Architecture Patterns

1. **Separation of Concerns**
   - Domain models (Burst, Fact) in `internal/domain/career/`
   - Inference logic in `internal/service/career/burst_fact/`
   - UI components in `internal/cli/models/`
   - Repository implementations in `internal/repository/career/`

2. **Reuse Existing Patterns**
   - Follow BubbleTea Model pattern from existing screens
   - Use existing styling system from `internal/cli/styles/`
   - Follow validation patterns from metadata clarification
   - Build on existing service layer architecture

3. **Domain-Driven Design**
   - Burst and Fact as first-class domain concepts
   - Validation at domain level
   - Service layer orchestrates inference
   - Repository handles persistence

4. **Inference System Design**
   - Modular inference rules (easily extensible)
   - Confidence scoring for suggestions
   - User confirmation workflow for all inferences
   - Traceability to source events/bursts

### Key Files to Create

**Domain Models**:
- `internal/domain/career/burst.go` (Burst model)
- `internal/domain/career/burst_test.go` (Burst tests)
- `internal/domain/career/fact.go` (Fact model)
- `internal/domain/career/fact_test.go` (Fact tests)

**Inference Engine**:
- `internal/service/career/burst_fact/detector.go` (Burst detection)
- `internal/service/career/burst_fact/detector_test.go` (Detection tests)
- `internal/service/career/burst_fact/extractor.go` (Fact extraction)
- `internal/service/career/burst_fact/extractor_test.go` (Extraction tests)
- `internal/service/career/burst_fact/classifier.go` (Role fit, audience, etc.)
- `internal/service/career/burst_fact/classifier_test.go` (Classification tests)

**Repository**:
- `internal/repository/career/burst_repository.go` (Burst interface)
- `internal/repository/career/memory_burst_repository.go` (In-memory impl)
- `internal/repository/career/sqlite_burst_repository.go` (SQLite impl)
- `internal/repository/career/fact_repository.go` (Fact interface)
- `internal/repository/career/memory_fact_repository.go` (In-memory impl)
- `internal/repository/career/sqlite_fact_repository.go` (SQLite impl)

**UI Components**:
- `internal/cli/models/burst_list.go` (Burst list screen)
- `internal/cli/models/burst_list_test.go` (List tests)
- `internal/cli/models/burst_suggestion.go` (Burst suggestion screen)
- `internal/cli/models/burst_suggestion_test.go` (Suggestion tests)
- `internal/cli/models/fact_editor.go` (Fact editor screen)
- `internal/cli/models/fact_editor_test.go` (Editor tests)
- `internal/cli/models/fact_list.go` (Fact list screen)
- `internal/cli/models/fact_list_test.go` (List tests)

**Documentation**:
- `docs/BURST_FACT_EXTRACTION_GUIDE.md` (User guide)
- Updated `README.md`, `CLI_GUIDE.md`, `CHANGELOG.md`

### Success Criteria (All Must Be Met)

- [x] Users can see suggested bursts after metadata clarification
- [x] Users can accept/reject burst suggestions
- [x] Burst detection algorithm works with various event similarities
- [x] Facts are extracted from events and bursts
- [x] Users can review and confirm extracted facts
- [x] Role fit classification works correctly
- [x] Audience relevance inference is accurate
- [x] All inferences are traceable to source events/bursts
- [x] Aspirational language is rejected
- [x] Metrics are validated for being grounded
- [x] All changes persisted to database
- [x] Code coverage ≥ 80%
- [x] All tests passing (100% pass rate)
- [x] Race detector passes (0 conditions)
- [x] Performance targets met (≤2s for 500 events)

---

## Phase Dependencies

This feature builds on:
- **Phase 1-2**: Event capture and metadata clarification ✅
- **Phase 2**: Event listing and filtering ✅
- **Phase 3**: Metadata review workflow ✅

This feature enables:
- **Phase 6**: Portfolio/case study generation
- **Phase 7**: CV generation with burst grouping
- **Phase 8**: Advanced analytics

---

## Estimated Effort

- Phase 1: 8-10 hours (domain models, repositories, inference engine foundation)
- Phase 2: 6-8 hours (burst detection, UI components)
- Phase 3: 6-8 hours (fact extraction, UI components)
- Phase 4: 4-6 hours (integration with existing features)
- Phase 5: 4-6 hours (testing, documentation)

**Total**: 28-38 hours

---

## Completion Tracking

- **Phase 1**: ⏳ Ready for execution
- **Phase 2**: ⏳ Awaiting Phase 1 completion
- **Phase 3**: ⏳ Awaiting Phase 2 completion
- **Phase 4**: ⏳ Awaiting Phase 3 completion
- **Phase 5**: ⏳ Awaiting Phase 4 completion

---

**Document Version**: 1.0
**Created**: 2025-12-30
**Status**: Ready for Phase 1 Execution
**Template Source**: tasks-03-metadata-clarification.md
**Process Guide**: docs/rules/process-task-list.md


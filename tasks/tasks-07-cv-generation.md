# Task List: CV Generation Feature

**PRD Reference**: `docs/features/05-cv-generation.md`

**Purpose**: Transform raw career events into credible, audience- and role-specific CV views with full traceability, conservative defaults, and no text rewriting.

**Status**: ⏳ **NOT STARTED** (0%)

---

## Tasks

### Phase 1: Foundation & Core Components

#### 1.0 Create CV Domain Model and Persistence
- [ ] 1.1 Define CVView struct with fields: ID, Name, TargetRole (Principal/Staff/EM/Senior IC), TargetAudience ([]string: hiring_manager, recruiter, peer), GeneratedAt, Metadata (map[string]interface{})
- [ ] 1.2 Define CVSection struct with fields: ID, CVViewID, SectionType (string: experience, skills, summary), Title, Order, Content ([]CVBullet)
- [ ] 1.3 Define CVBullet struct with fields: ID, Text, SourceEventIDs ([]string), SourceFactIDs ([]string), Rank (float64), InclusionReason (string), CreatedAt
- [ ] 1.4 Implement CVView validation rules (target role required, at least one audience, valid role/audience values)
- [ ] 1.5 Implement CVBullet validation rules (≥1 source event/fact, single-claim bullets, no aspirational language, no inferred metrics)
- [ ] 1.6 Create repository interface methods for CVView (Create, GetByID, Update, Delete, List, Count)
- [ ] 1.7 Create repository interface methods for CVSection (Create, GetByID, GetByCVViewID, Update, Delete)
- [ ] 1.8 Create repository interface methods for CVBullet (Create, GetByID, GetBySectionID, Update, Delete)
- [ ] 1.9 Implement MemoryRepository for CVView, CVSection, CVBullet operations (thread-safe with sync.RWMutex)
- [ ] 1.10 Implement SQLiteRepository for CVView persistence with schema:
  ```sql
  CREATE TABLE IF NOT EXISTS cv_views (
      id TEXT PRIMARY KEY,
      name TEXT NOT NULL,
      target_role TEXT NOT NULL,
      target_audience TEXT NOT NULL,
      generated_at DATETIME NOT NULL,
      metadata TEXT
  )
  ```
- [ ] 1.11 Implement SQLiteRepository for CVSection persistence with schema:
  ```sql
  CREATE TABLE IF NOT EXISTS cv_sections (
      id TEXT PRIMARY KEY,
      cv_view_id TEXT NOT NULL,
      section_type TEXT NOT NULL,
      title TEXT NOT NULL,
      section_order INTEGER NOT NULL,
      FOREIGN KEY (cv_view_id) REFERENCES cv_views(id) ON DELETE CASCADE
  )
  ```
- [ ] 1.12 Implement SQLiteRepository for CVBullet persistence with schema:
  ```sql
  CREATE TABLE IF NOT EXISTS cv_bullets (
      id TEXT PRIMARY KEY,
      section_id TEXT NOT NULL,
      text TEXT NOT NULL,
      source_event_ids TEXT NOT NULL,
      source_fact_ids TEXT,
      rank REAL NOT NULL,
      inclusion_reason TEXT,
      created_at DATETIME NOT NULL,
      FOREIGN KEY (section_id) REFERENCES cv_sections(id) ON DELETE CASCADE
  )
  ```
- [ ] 1.13 Write comprehensive unit tests for CVView, CVSection, CVBullet validation (edge cases, boundary conditions)
- [ ] 1.14 Write repository tests for all CRUD operations and filtering

#### 2.0 Create Bullet Generation Engine
- [ ] 2.1 Implement BulletGenerator service in `internal/service/career/cv/bullet_generator.go`
- [ ] 2.2 Create method: `GenerateBullets(ctx, events, facts, targetRole, targetAudience) -> []CVBullet`
- [ ] 2.3 Implement inclusion criteria filter (bullet must trace to ≥1 event, single-claim only, prefer repeated signals)
- [ ] 2.4 Implement exclusion criteria filter (no inferred metrics, no aspirational language, no role inflation)
- [ ] 2.5 Implement ranking algorithm with priority: Ownership > Contribution > Strategy > Execution > Outcome > Activity
- [ ] 2.6 Create role-specific bullet generation rules:
  - Principal: Strategic ownership, cross-team leadership (3-4 bullets max)
  - Staff: Technical leadership, high-complexity implementation (4-5 bullets max)
  - EM: Team leadership, mentorship, delivery accountability (3-4 bullets max)
  - Senior IC: Deep technical contribution, system design (4-5 bullets max)
- [ ] 2.7 Create audience-specific filtering:
  - Hiring Manager: Outcomes, ownership, business impact
  - Recruiter: Skills, competencies, high-level achievements
  - Peer: Technical depth, collaboration, problem-solving
- [ ] 2.8 Implement compression logic (older roles compress first, respect bullet caps)
- [ ] 2.9 Implement traceability tracking (store source event/fact IDs in CVBullet)
- [ ] 2.10 Write comprehensive unit tests for bullet generation with various event/fact combinations
- [ ] 2.11 Write tests for inclusion/exclusion criteria edge cases
- [ ] 2.12 Write tests for ranking algorithm with different event types

#### 3.0 Create CV Section Builder
- [ ] 3.1 Implement SectionBuilder service in `internal/service/career/cv/section_builder.go`
- [ ] 3.2 Create method: `BuildSections(ctx, bullets, targetRole) -> []CVSection`
- [ ] 3.3 Implement experience section generation (group bullets by company/project/timeframe)
- [ ] 3.4 Implement skills section generation (extract competencies from bullets/facts)
- [ ] 3.5 Implement summary section generation (optional, based on role fit and top achievements)
- [ ] 3.6 Implement section ordering logic (experience first, then skills, then summary)
- [ ] 3.7 Implement content organization within sections (chronological for experience, grouped for skills)
- [ ] 3.8 Write comprehensive unit tests for section building
- [ ] 3.9 Write tests for section ordering and content organization

### Phase 2: CV Generation Service

#### 4.0 Create CV Generation Orchestrator
- [ ] 4.1 Implement CVGenerationService in `internal/service/career/cv/cv_generation_service.go`
- [ ] 4.2 Create method: `GenerateCV(ctx, cvViewConfig) -> CVView` (orchestrates bullet generation, section building, validation)
- [ ] 4.3 Implement event retrieval from repository based on filters (date range, tags, companies)
- [ ] 4.4 Implement fact retrieval from repository for selected events
- [ ] 4.5 Call BulletGenerator to create bullets from events/facts
- [ ] 4.6 Call SectionBuilder to organize bullets into sections
- [ ] 4.7 Create CVView with metadata (generation timestamp, source count, filters applied)
- [ ] 4.8 Validate generated CV against acceptance criteria (traceability, bullet caps, factual accuracy)
- [ ] 4.9 Persist CVView, CVSections, CVBullets to repository
- [ ] 4.10 Write comprehensive unit tests for CV generation orchestration
- [ ] 4.11 Write integration tests for complete CV generation workflow
- [ ] 4.12 Write performance tests (CV generation ≤2s for ≤500 events)

#### 5.0 Create CV View Management Service
- [ ] 5.1 Implement CVViewService in `internal/service/career/cv/cv_view_service.go`
- [ ] 5.2 Create method: `ListCVViews(ctx, filters) -> []CVView` (list all generated CV views)
- [ ] 5.3 Create method: `GetCVView(ctx, cvViewID) -> CVView` (retrieve specific CV view with sections/bullets)
- [ ] 5.4 Create method: `DeleteCVView(ctx, cvViewID) -> error` (delete CV view and related sections/bullets)
- [ ] 5.5 Create method: `UpdateCVViewMetadata(ctx, cvViewID, metadata) -> error` (update view name, description)
- [ ] 5.6 Create method: `RegenerateCVView(ctx, cvViewID) -> CVView` (regenerate CV with same config but latest events/facts)
- [ ] 5.7 Create method: `GetSourceEvents(ctx, cvViewID) -> []CareerEvent` (retrieve all events used in CV)
- [ ] 5.8 Create method: `GetSourceFacts(ctx, cvViewID) -> []Fact` (retrieve all facts used in CV)
- [ ] 5.9 Write comprehensive unit tests for CV view management
- [ ] 5.10 Write integration tests for view retrieval and updates

#### 6.0 Create Traceability System
- [ ] 6.1 Implement TraceabilityService in `internal/service/career/cv/traceability_service.go`
- [ ] 6.2 Create method: `GetBulletSources(ctx, bulletID) -> ([]CareerEvent, []Fact)` (retrieve source events/facts for bullet)
- [ ] 6.3 Create method: `GetEventUsage(ctx, eventID) -> []CVBullet` (find all bullets using this event)
- [ ] 6.4 Create method: `GetFactUsage(ctx, factID) -> []CVBullet` (find all bullets using this fact)
- [ ] 6.5 Create method: `ValidateTraceability(ctx, cvViewID) -> ValidationReport` (ensure all bullets trace to sources)
- [ ] 6.6 Implement visualization data for source event → bullet mapping
- [ ] 6.7 Implement audit trail for CV generation (track when generated, from which events/facts)
- [ ] 6.8 Write comprehensive unit tests for traceability operations
- [ ] 6.9 Write integration tests for source tracking and validation

### Phase 3: UI Components for CV Generation

#### 7.0 Create CV Configuration Screen
- [ ] 7.1 Implement CVConfigModel as BubbleTea Model in `internal/cli/models/cv_config.go`
- [ ] 7.2 Display form fields: CV name, target role (dropdown: Principal/Staff/EM/Senior IC), target audience (multi-select: hiring_manager/recruiter/peer)
- [ ] 7.3 Add optional filters: date range, companies, tags, competencies
- [ ] 7.4 Implement field validation (required fields, valid selections)
- [ ] 7.5 Implement keyboard navigation (Tab for fields, Space for multi-select, Enter to generate)
- [ ] 7.6 Add preview of event count matching filters
- [ ] 7.7 Display estimated bullet count based on role and events
- [ ] 7.8 Implement "Generate CV" button with confirmation
- [ ] 7.9 Write comprehensive unit tests for CV config interactions
- [ ] 7.10 Test edge cases (no events matching filters, invalid role/audience combinations)

#### 8.0 Create CV Preview Screen
- [ ] 8.1 Implement CVPreviewModel as BubbleTea Model in `internal/cli/models/cv_preview.go`
- [ ] 8.2 Display generated CV with sections (experience, skills, summary)
- [ ] 8.3 Render bullets with formatting (bullet points, proper indentation)
- [ ] 8.4 Implement scrolling through long CV content (up/down arrows)
- [ ] 8.5 Add visual indicators for bullet source (event count, fact count)
- [ ] 8.6 Implement "View Sources" action for selected bullet (press 's' to see source events/facts)
- [ ] 8.7 Implement "Export CV" action (press 'e' to export to file)
- [ ] 8.8 Implement "Regenerate CV" action (press 'r' to regenerate with latest data)
- [ ] 8.9 Implement "Edit Config" action (press 'c' to modify CV configuration)
- [ ] 8.10 Write comprehensive unit tests for CV preview interactions
- [ ] 8.11 Test edge cases (empty CV, single bullet, very long CV)

#### 9.0 Create CV List Screen
- [ ] 9.1 Implement CVListModel as BubbleTea Model in `internal/cli/models/cv_list.go`
- [ ] 9.2 Display list of generated CV views with name, role, audience, generation date
- [ ] 9.3 Implement scrolling through CV list (up/down arrows)
- [ ] 9.4 Implement selection and preview (Enter to view selected CV)
- [ ] 9.5 Implement filtering by target role or audience
- [ ] 9.6 Implement sorting by generation date, name, or role
- [ ] 9.7 Implement "Delete CV" action (press 'd' with confirmation)
- [ ] 9.8 Implement "Regenerate CV" action (press 'r' to regenerate selected CV)
- [ ] 9.9 Add visual indicators for CV staleness (outdated if events/facts changed since generation)
- [ ] 9.10 Write comprehensive unit tests for CV list interactions
- [ ] 9.11 Test edge cases (empty list, single CV, large list with 100+ CVs)

#### 10.0 Create Bullet Source Viewer
- [ ] 10.1 Implement BulletSourceModel as BubbleTea Model in `internal/cli/models/bullet_source.go`
- [ ] 10.2 Display selected bullet text at top
- [ ] 10.3 Show list of source events with preview (text, date, company)
- [ ] 10.4 Show list of source facts with competency badges
- [ ] 10.5 Implement scrolling through sources (up/down arrows)
- [ ] 10.6 Implement "View Full Event" action (press 'e' to see complete event details)
- [ ] 10.7 Implement "View Full Fact" action (press 'f' to see complete fact details)
- [ ] 10.8 Add visual indicators for source contribution (primary vs. supporting evidence)
- [ ] 10.9 Write comprehensive unit tests for source viewer interactions
- [ ] 10.10 Test edge cases (single source, multiple sources, mixed event/fact sources)

### Phase 4: Integration with Existing Features

#### 11.0 Integrate CV Generation with Event Timeline
- [ ] 11.1 Add "Generate CV" option to main menu (press 'g' from home screen)
- [ ] 11.2 Navigate to CV configuration screen when option selected
- [ ] 11.3 Pass event filters from current view to CV config (if on filtered event list)
- [ ] 11.4 Navigate to CV preview screen after generation completes
- [ ] 11.5 Allow navigation back to event timeline from CV preview
- [ ] 11.6 Write integration tests for event timeline → CV generation workflow

#### 12.0 Integrate CV Generation with Fact Extraction
- [ ] 12.1 Include extracted facts in bullet generation automatically
- [ ] 12.2 Prioritize bullets with fact support over event-only bullets
- [ ] 12.3 Display fact source indicator in CV preview
- [ ] 12.4 Allow viewing facts from bullet source viewer
- [ ] 12.5 Write integration tests for fact → CV bullet workflow

#### 13.0 Integrate CV Generation with Burst Detection
- [ ] 13.1 Group related bullets from same burst in CV sections
- [ ] 13.2 Use burst metadata (competency focus) to enhance bullet ranking
- [ ] 13.3 Display burst context in bullet source viewer
- [ ] 13.4 Allow generating CV focused on specific burst (filter by burst ID)
- [ ] 13.5 Write integration tests for burst → CV generation workflow

#### 14.0 Create CV Export Functionality
- [ ] 14.1 Implement CVExporter in `internal/service/career/cv/cv_exporter.go`
- [ ] 14.2 Create method: `ExportToMarkdown(ctx, cvView) -> string` (export CV to Markdown format)
- [ ] 14.3 Create method: `ExportToJSON(ctx, cvView) -> string` (export CV to JSON with traceability)
- [ ] 14.4 Create method: `ExportToPlainText(ctx, cvView) -> string` (export CV to plain text for ATS)
- [ ] 14.5 Implement file writing with user-specified path
- [ ] 14.6 Add export format selection to CV preview screen
- [ ] 14.7 Add confirmation message after successful export
- [ ] 14.8 Write comprehensive unit tests for export formats
- [ ] 14.9 Write integration tests for complete export workflow

### Phase 5: Testing and Documentation

#### 15.0 Comprehensive Testing Suite
- [ ] 15.1 Write end-to-end tests for complete CV generation workflow (config → generate → preview → export)
- [ ] 15.2 Test bullet generation with various event/fact combinations (high/medium/low quality)
- [ ] 15.3 Test inclusion/exclusion criteria with edge cases (aspirational language, inferred metrics, role inflation)
- [ ] 15.4 Test ranking algorithm with different priority scenarios
- [ ] 15.5 Test compression logic (bullet caps, older roles compress first)
- [ ] 15.6 Test role-specific generation for all 4 roles (Principal, Staff, EM, Senior IC)
- [ ] 15.7 Test audience-specific filtering for all 3 audiences (hiring_manager, recruiter, peer)
- [ ] 15.8 Test traceability system (all bullets trace to sources, no orphaned bullets)
- [ ] 15.9 Test CV view management (create, list, delete, regenerate)
- [ ] 15.10 Test integration with event timeline, fact extraction, burst detection
- [ ] 15.11 Run race detector: `go test -race ./...` (verify 0 race conditions)
- [ ] 15.12 Verify code coverage meets 80%+ threshold across all new code
- [ ] 15.13 Performance test: CV generation ≤2s for ≤500 events
- [ ] 15.14 Performance test: bullet ranking ≤100ms for 1000 bullets
- [ ] 15.15 Performance test: traceability lookup ≤50ms per bullet

#### 16.0 Documentation and User Guidance
- [ ] 16.1 Create CV_GENERATION_GUIDE.md with comprehensive feature overview
- [ ] 16.2 Document bullet generation rules (inclusion/exclusion criteria, ranking priority)
- [ ] 16.3 Document role-specific generation rules and bullet caps
- [ ] 16.4 Document audience-specific filtering rules
- [ ] 16.5 Provide examples of generated CVs for each role and audience
- [ ] 16.6 Document compression logic and how older roles are compressed
- [ ] 16.7 Document traceability system and how to view sources
- [ ] 16.8 Document export formats and use cases
- [ ] 16.9 Document keyboard shortcuts for CV screens
- [ ] 16.10 Update README.md with CV generation features
- [ ] 16.11 Update CLI_GUIDE.md with CV workflow shortcuts
- [ ] 16.12 Update CHANGELOG.md with feature description and test results
- [ ] 16.13 Create troubleshooting guide for common CV generation issues
- [ ] 16.14 Document conservative defaults and generation principles

---

## Implementation Guidelines

### Architecture Patterns

1. **Separation of Concerns**
   - Domain models (CVView, CVSection, CVBullet) in `internal/domain/career/`
   - Bullet generation logic in `internal/service/career/cv/`
   - UI components in `internal/cli/models/`
   - Repository implementations in `internal/repository/career/`

2. **Reuse Existing Patterns**
   - Follow BubbleTea Model pattern from existing screens
   - Use existing styling system from `internal/cli/styles/`
   - Follow validation patterns from metadata clarification and burst-fact extraction
   - Build on existing service layer architecture

3. **Domain-Driven Design**
   - CVView, CVSection, CVBullet as first-class domain concepts
   - Validation at domain level (no aspirational language, no inferred metrics)
   - Service layer orchestrates generation
   - Repository handles persistence

4. **Read-Only Generation**
   - CVs are generated on-demand, not stored permanently (optional persistence for caching)
   - No modification of source events or facts during generation
   - Regeneration always reflects latest data
   - Traceability preserved for all bullets

5. **Conservative Defaults**
   - Prefer explicit user choices over inferred preferences
   - Err on side of fewer bullets rather than more
   - Strict validation (reject questionable content)
   - Clear source attribution

### Key Files to Create

**Domain Models**:
- `internal/domain/career/cv_view.go` - CVView struct and validation
- `internal/domain/career/cv_view_test.go` - CVView tests
- `internal/domain/career/cv_section.go` - CVSection struct and validation
- `internal/domain/career/cv_section_test.go` - CVSection tests
- `internal/domain/career/cv_bullet.go` - CVBullet struct and validation
- `internal/domain/career/cv_bullet_test.go` - CVBullet tests

**Generation Engine**:
- `internal/service/career/cv/bullet_generator.go` - Bullet generation logic
- `internal/service/career/cv/bullet_generator_test.go` - Bullet generation tests
- `internal/service/career/cv/section_builder.go` - Section building logic
- `internal/service/career/cv/section_builder_test.go` - Section building tests
- `internal/service/career/cv/cv_generation_service.go` - CV generation orchestrator
- `internal/service/career/cv/cv_generation_service_test.go` - CV generation tests
- `internal/service/career/cv/cv_view_service.go` - CV view management
- `internal/service/career/cv/cv_view_service_test.go` - CV view management tests
- `internal/service/career/cv/traceability_service.go` - Traceability tracking
- `internal/service/career/cv/traceability_service_test.go` - Traceability tests
- `internal/service/career/cv/cv_exporter.go` - Export functionality
- `internal/service/career/cv/cv_exporter_test.go` - Export tests
- `internal/service/career/cv/integration_test.go` - End-to-end integration tests

**Repository**:
- `internal/repository/career/cv_view_repository.go` - CVView repository interface
- `internal/repository/career/memory_cv_view_repository.go` - In-memory implementation
- `internal/repository/career/sqlite_cv_view_repository.go` - SQLite implementation
- `internal/repository/career/cv_section_repository.go` - CVSection repository interface
- `internal/repository/career/cv_bullet_repository.go` - CVBullet repository interface

**UI Components**:
- `internal/cli/models/cv_config.go` - CV configuration screen
- `internal/cli/models/cv_config_test.go` - CV config tests
- `internal/cli/models/cv_preview.go` - CV preview screen
- `internal/cli/models/cv_preview_test.go` - CV preview tests
- `internal/cli/models/cv_list.go` - CV list screen
- `internal/cli/models/cv_list_test.go` - CV list tests
- `internal/cli/models/bullet_source.go` - Bullet source viewer
- `internal/cli/models/bullet_source_test.go` - Bullet source tests

**Documentation**:
- `docs/CV_GENERATION_GUIDE.md` - Comprehensive user guide
- Updated `README.md`, `CLI_GUIDE.md`, `CHANGELOG.md`

### Success Criteria (All Must Be Met)

- [ ] Users can generate role-specific CVs (Principal, Staff, EM, Senior IC)
- [ ] Users can generate audience-specific CVs (hiring_manager, recruiter, peer)
- [ ] Bullet generation respects inclusion/exclusion criteria
- [ ] Bullet ranking follows priority: Ownership > Contribution > Strategy > Execution > Outcome > Activity
- [ ] Compression logic enforces bullet caps per role (Principal: 3-4, Staff: 4-5, EM: 3-4, Senior IC: 4-5)
- [ ] All bullets trace to ≥1 source event or fact
- [ ] No aspirational language in generated bullets
- [ ] No inferred metrics in generated bullets
- [ ] No role inflation in generated bullets
- [ ] Traceability system allows viewing sources for any bullet
- [ ] CV generation completes in ≤2s for ≤500 events
- [ ] All CVs can be exported to Markdown, JSON, and plain text
- [ ] Code coverage ≥ 80%
- [ ] All tests passing (100% pass rate)
- [ ] Race detector passes (0 conditions detected)

---

## Phase Dependencies

This feature builds on:
- **Phase 1-2**: Event capture and metadata clarification ✅
- **Phase 3-4**: Burst detection and fact extraction ✅

This feature enables:
- **Phase 7**: Portfolio/case study generation
- **Phase 8**: Advanced analytics and career insights

---

## Estimated Effort

- Phase 1: 10-12 hours (domain models, repositories, bullet generation engine)
- Phase 2: 8-10 hours (CV generation service, view management, traceability)
- Phase 3: 10-12 hours (UI components: config, preview, list, source viewer)
- Phase 4: 6-8 hours (integration with existing features, export functionality)
- Phase 5: 6-8 hours (testing, documentation)

**Total**: 40-50 hours
**Completed**: 0 hours
**Remaining**: 40-50 hours

---

## Completion Tracking

- **Phase 1**: ⏳ **NOT STARTED (0%)**
  - [ ] CVView, CVSection, CVBullet domain models (1.1-1.5)
  - [ ] Repository interfaces and implementations (1.6-1.12)
  - [ ] Domain model tests (1.13-1.14)
  - [ ] Bullet generation engine (2.1-2.12)
  - [ ] Section builder (3.1-3.9)

- **Phase 2**: ⏳ **NOT STARTED (0%)**
  - [ ] CV generation orchestrator (4.1-4.12)
  - [ ] CV view management service (5.1-5.10)
  - [ ] Traceability system (6.1-6.9)

- **Phase 3**: ⏳ **NOT STARTED (0%)**
  - [ ] CV configuration screen (7.1-7.10)
  - [ ] CV preview screen (8.1-8.11)
  - [ ] CV list screen (9.1-9.11)
  - [ ] Bullet source viewer (10.1-10.10)

- **Phase 4**: ⏳ **NOT STARTED (0%)**
  - [ ] Event timeline integration (11.1-11.6)
  - [ ] Fact extraction integration (12.1-12.5)
  - [ ] Burst detection integration (13.1-13.5)
  - [ ] CV export functionality (14.1-14.9)

- **Phase 5**: ⏳ **NOT STARTED (0%)**
  - [ ] Comprehensive testing (15.1-15.15)
  - [ ] Documentation (16.1-16.14)

---

## Key Design Decisions

### 1. Read-Only Generation
CVs are generated on-demand from source events and facts. No permanent storage of generated content (except optional caching for performance). This ensures CVs always reflect latest data.

### 2. Strict Validation
- No aspirational language (will, should, could, might, may, want, wish, hope, plan, intend, attempt, try, would)
- No inferred metrics (must be grounded in source events)
- No role inflation (must match actual role in source events)
- Single-claim bullets only (one achievement per bullet)

### 3. Ranking Algorithm
Priority: Ownership > Contribution > Strategy > Execution > Outcome > Activity

This ensures bullets showing ownership and strategic impact are ranked higher than purely execution-focused bullets.

### 4. Compression Logic
Older roles compress first when approaching bullet caps. This preserves recency bias while respecting hard limits.

### 5. Traceability
Every bullet stores source event/fact IDs. Users can view sources at any time to understand bullet origin and validate accuracy.

### 6. Conservative Defaults
- Minimum bullets rather than maximum (quality over quantity)
- Explicit user choices over inferred preferences
- Clear attribution to sources
- Reversible and explainable generation

---

## Performance Targets

- CV generation: ≤2s for ≤500 events
- Bullet ranking: ≤100ms for 1000 bullets
- Traceability lookup: ≤50ms per bullet
- Section building: ≤50ms per CV
- Export to Markdown/JSON/plain text: ≤100ms per CV

---

## Testing Strategy

### Unit Tests
- Domain model validation (CVView, CVSection, CVBullet)
- Bullet generation (inclusion/exclusion criteria, ranking)
- Section building (organization, ordering)
- Traceability (source tracking, validation)
- Export formats (Markdown, JSON, plain text)

### Integration Tests
- Complete CV generation workflow (config → generate → preview → export)
- Integration with event timeline
- Integration with fact extraction
- Integration with burst detection
- Traceability across system boundaries

### Performance Tests
- CV generation with 100, 500, 1000 events
- Bullet ranking with 100, 500, 1000 bullets
- Traceability lookup with 100, 500, 1000 bullets

### Edge Case Tests
- No events matching filters
- Single event
- Events with no facts
- Events with multiple facts
- All bullets fail validation
- Maximum bullet cap reached
- Empty sections
- Very long bullet text

---

## Next Steps

1. **Phase 1.0**: Create domain models (CVView, CVSection, CVBullet) with validation
2. **Phase 1.0**: Implement repositories (memory and SQLite)
3. **Phase 1.0**: Create bullet generation engine
4. **Phase 1.0**: Create section builder
5. **Phase 2.0**: Implement CV generation orchestrator
6. **Phase 2.0**: Implement traceability system
7. **Phase 3.0**: Create UI components (config, preview, list, source viewer)
8. **Phase 4.0**: Integrate with existing features
9. **Phase 5.0**: Comprehensive testing and documentation

---

**Document Version**: 1.0
**Created**: 2025-12-31
**Status**: NOT STARTED (0%)
**Template Source**: tasks-05-burst-fact-extraction.md
**Process Guide**: docs/rules/master-task-prompt.md


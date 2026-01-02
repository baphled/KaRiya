# Task List: CV Generation Feature

**PRD Reference**: `docs/features/05-cv-generation.md`

**Purpose**: Transform raw career events into credible, audience- and role-specific CV views with full traceability, conservative defaults, and no text rewriting.

**Status**: ⏳ **NOT STARTED** (0%)

**Version**: 2.0 - Updated with YAML Configuration & Ephemeral Generation

---

## Key Architecture Decisions (UPDATED)

- ✅ **CVs are ephemeral** (in-memory only, NOT stored in database)
- ✅ **Configurations stored as YAML files** (NOT in database)
- ✅ **Location**: `$HOME/.kariya/cv_configs/`
- ✅ **Generation**: On-the-fly from events, always fresh
- ✅ **Export**: Users can export to text/markdown
- ✅ **Traceability**: Full source event/fact tracking maintained

---

## Tasks

### Phase 1: Foundation & Core Components

#### 1.0 Create CV Domain Models

- [ ] 1.1 Define `CVView` struct in `internal/domain/career/cv.go`
  - Fields: ID, Name, TargetRole, TargetAudience, EventFilters, GeneratedAt, SourceEventCount, SourceFactCount
  - Implement Validate() method
  - No CreatedAt/UpdatedAt (ephemeral only)

- [ ] 1.2 Define `CVSection` struct in `internal/domain/career/cv.go`
  - Fields: ID, CVViewID, SectionType, Title, Order, Content
  - Implement Validate() method

- [ ] 1.3 Define `CVBullet` struct in `internal/domain/career/cv.go`
  - Fields: ID, SectionID, Text, SourceEventIDs, SourceFactIDs, Rank, InclusionReason, Confidence
  - Implement Validate() method

- [ ] 1.4 Define `CVConfig` struct in `internal/domain/career/cv.go`
  - Fields: Name, TargetRole, TargetAudience, EventFilters, CreatedAt, UpdatedAt
  - Implement YAML serialization/deserialization
  - Implement Validate() method

- [ ] 1.5 Add validation helpers in `internal/domain/career/cv.go`
  - IsAspirationLanguage(text string) bool
  - IsSingleClaimBullet(text string) bool
  - HasInferredMetrics(text string) bool
  - IsRoleInflation(text string, roleFit RoleFit) bool

- [ ] 1.6 Create error types in `internal/domain/career/cv.go`
  - ErrInvalidCVRole
  - ErrInvalidAudience
  - ErrBulletMultipleClaims
  - ErrBulletAspirationLanguage
  - ErrBulletInferredMetrics
  - ErrNoSourceEvents

- [ ] 1.7 Write comprehensive unit tests in `internal/domain/career/cv_test.go`
  - CVView validation tests
  - CVSection validation tests
  - CVBullet validation tests
  - CVConfig YAML serialization tests
  - Helper function tests
  - Edge cases and boundary conditions
  - Use Ginkgo/Gomega pattern

#### 2.0 Create YAML Configuration System

- [ ] 2.1 Define `CVConfigManager` interface in `internal/service/career/cv/config_manager.go`
  - Methods: LoadConfig, SaveConfig, DeleteConfig, ListConfigs, GetConfigPath
  - Support context.Context for all operations

- [ ] 2.2 Implement `YAMLConfigManager` in `internal/service/career/cv/yaml_config_manager.go`
  - Store configs in `$HOME/.kariya/cv_configs/` directory
  - Use `gopkg.in/yaml.v3` for serialization
  - Create directory if it doesn't exist
  - Handle file I/O errors gracefully
  - Support atomic writes (write to temp, then rename)
  - Implement LoadConfig(ctx, name) (*CVConfig, error)
  - Implement SaveConfig(ctx, config) error
  - Implement DeleteConfig(ctx, name) error
  - Implement ListConfigs(ctx) ([]*CVConfig, error)
  - Implement GetConfigPath(name) string

- [ ] 2.3 Write comprehensive tests in `internal/service/career/cv/config_manager_test.go`
  - Load/save/delete config tests
  - List configs tests
  - YAML serialization/deserialization tests
  - Directory creation tests
  - Atomic write tests
  - Error handling tests
  - Use Ginkgo/Gomega pattern

#### 3.0 Remove Database Persistence (CORRECTED)

- [ ] 3.1 **DO NOT create** database tables for CVs
  - ❌ cv_views table NOT needed
  - ❌ cv_sections table NOT needed
  - ❌ cv_bullets table NOT needed

- [ ] 3.2 **DO NOT create** repository interfaces for CVs
  - ❌ CVViewRepository NOT needed
  - ❌ CVSectionRepository NOT needed
  - ❌ CVBulletRepository NOT needed

- [ ] 3.3 **DO NOT create** SQLite repository implementations
  - ❌ SQLiteCVViewRepository NOT needed
  - ❌ SQLiteCVSectionRepository NOT needed
  - ❌ SQLiteCVBulletRepository NOT needed

- [ ] 3.4 **DO NOT create** Memory repository implementations
  - ❌ MemoryCVViewRepository NOT needed
  - ❌ MemoryCVSectionRepository NOT needed
  - ❌ MemoryCVBulletRepository NOT needed

### Phase 2: CV Generation Service

#### 4.0 Create Bullet Generation Engine

- [ ] 4.1 Implement `BulletGenerator` service in `internal/service/career/cv/bullet_generator.go`
  - Constructor: NewBulletGenerator(eventRepo, factRepo, logger)
  - Method: GenerateBullets(ctx, events, facts, targetRole, targetAudience) -> []CVBullet

- [ ] 4.2 Implement inclusion criteria filter
  - Bullet must trace to ≥1 event
  - Single-claim bullets only
  - Prefer repeated signals
  - No aspirational language (use AspirationKeywords)
  - No inferred metrics
  - No role inflation

- [ ] 4.3 Implement ranking algorithm
  - Priority: Ownership > Contribution > Strategy > Execution > Outcome > Activity
  - Scoring: 0.0-1.0 based on priority and signals
  - Older events score lower (temporal decay)
  - Repeated signals boost score

- [ ] 4.4 Implement role-specific bullet caps
  - Principal: 3-4 bullets max
  - Staff: 4-5 bullets max
  - EM: 3-4 bullets max
  - Senior IC: 4-5 bullets max
  - Compression: Remove lower-ranked bullets first

- [ ] 4.5 Implement audience-specific filtering
  - Hiring Manager: Outcomes, ownership, business impact
  - Recruiter: Skills, competencies, high-level achievements
  - Peer: Technical depth, collaboration, problem-solving
  - Multi-audience: Include bullets relevant to ANY audience

- [ ] 4.6 Implement traceability tracking
  - Store source event IDs in CVBullet.SourceEventIDs
  - Store source fact IDs in CVBullet.SourceFactIDs
  - Store inclusion reason in CVBullet.InclusionReason
  - Calculate confidence score

- [ ] 4.7 Write comprehensive unit tests in `internal/service/career/cv/bullet_generator_test.go`
  - Inclusion criteria tests
  - Exclusion criteria tests
  - Ranking algorithm tests
  - Role-specific cap tests
  - Audience filtering tests
  - Edge cases (no events, no facts, mixed signals)
  - Performance tests (1000+ events)
  - Use Ginkgo/Gomega pattern

#### 5.0 Create CV Section Builder

- [ ] 5.1 Implement `SectionBuilder` service in `internal/service/career/cv/section_builder.go`
  - Constructor: NewSectionBuilder(logger)
  - Method: BuildSections(ctx, bullets, events, targetRole) -> []CVSection

- [ ] 5.2 Implement experience section generation
  - Group bullets by company/project/timeframe
  - Sort chronologically (newest first)
  - Create section with title "Experience"

- [ ] 5.3 Implement skills section generation
  - Extract competency categories from facts
  - Group bullets by category
  - Create section with title "Core Competencies"
  - Optional: only if bullets with fact sources

- [ ] 5.4 Implement summary section generation
  - Optional, based on role fit
  - Create 1-2 sentence summary from top bullets
  - Create section with title "Professional Summary"

- [ ] 5.5 Implement section ordering logic
  - Order: Experience → Skills → Summary
  - Sections without bullets are skipped

- [ ] 5.6 Implement content organization within sections
  - Experience: Chronological (newest first)
  - Skills: Grouped by category, alphabetical
  - Summary: High-priority bullets first

- [ ] 5.7 Write comprehensive unit tests in `internal/service/career/cv/section_builder_test.go`
  - Section generation tests
  - Section ordering tests
  - Content organization tests
  - Edge cases
  - Use Ginkgo/Gomega pattern

#### 6.0 Create CV Generation Orchestrator

- [ ] 6.1 Implement `CVGenerationService` in `internal/service/career/cv/cv_generation_service.go`
  - Constructor: NewCVGenerationService(eventRepo, factRepo, bulletGen, sectionBuilder, configManager, logger)
  - Method: GenerateCV(ctx, configName) -> *CVView, error
  - Method: GenerateCVFromConfig(ctx, config) -> *CVView, error

- [ ] 6.2 Implement event retrieval with filtering
  - Support filters: date range, tags, companies, categories
  - Apply filters to repository.List()

- [ ] 6.3 Implement fact retrieval for selected events
  - Get facts for each event
  - Filter by role fit if applicable

- [ ] 6.4 Orchestrate bullet generation and section building
  - Call BulletGenerator.GenerateBullets()
  - Call SectionBuilder.BuildSections()
  - Validate results

- [ ] 6.5 Create CVView with metadata
  - Set TargetRole, TargetAudience
  - Set GeneratedAt timestamp
  - Set SourceEventCount, SourceFactCount
  - Store EventFilters

- [ ] 6.6 **No persistence** - CVView is ephemeral only
  - Return in-memory CVView
  - Do NOT save to database
  - Do NOT save to files
  - CV exists only for current session

- [ ] 6.7 Write comprehensive unit tests in `internal/service/career/cv/cv_generation_service_test.go`
  - Complete CV generation workflow
  - Filtering and event retrieval
  - Ephemeral nature verification
  - Error handling
  - Use Ginkgo/Gomega pattern

#### 7.0 Create Traceability Service

- [ ] 7.1 Implement `TraceabilityService` in `internal/service/career/cv/traceability_service.go`
  - Constructor: NewTraceabilityService(eventRepo, factRepo, logger)

- [ ] 7.2 Implement GetBulletSources(ctx, bulletID) -> ([]CareerEvent, []Fact, error)
  - Retrieve source events for bullet
  - Retrieve source facts for bullet
  - Return both with full details

- [ ] 7.3 Implement GetEventUsage(ctx, eventID) -> []CVBullet, error
  - Find all bullets using specific event
  - Return bullets with context

- [ ] 7.4 Implement GetFactUsage(ctx, factID) -> []CVBullet, error
  - Find all bullets using specific fact
  - Return bullets with context

- [ ] 7.5 Implement ValidateTraceability(ctx, cvViewID) -> ValidationReport, error
  - Verify all bullets have valid source events/facts
  - Verify all source IDs reference existing entities
  - Return validation report with any issues

- [ ] 7.6 Implement visualization data methods
  - GetEventBulletMapping(ctx) - map events to bullets
  - GetFactBulletMapping(ctx) - map facts to bullets
  - Return data suitable for visualization

- [ ] 7.7 Write comprehensive unit tests in `internal/service/career/cv/traceability_service_test.go`
  - Source retrieval tests
  - Usage tracking tests
  - Validation tests
  - Edge cases
  - Use Ginkgo/Gomega pattern

### Phase 3: UI Components for CV Generation

#### 8.0 Create CV Configuration Management UI

- [ ] 8.1 Implement `CVConfigManagerModel` in `internal/cli/models/cv_config_manager.go`
  - Embed BaseStandardModel
  - Display list of YAML configs
  - Fields: configManager, configs, selectedIdx, header, helpFooter, footer

- [ ] 8.2 Implement config list display
  - Show table with columns: Name, Role, Audiences, Last Updated
  - Support pagination
  - Highlight selected config

- [ ] 8.3 Implement navigation
  - j/k or arrow keys: navigate list
  - Enter: generate CV from selected config
  - n: new config (go to editor)
  - e: edit config
  - d: delete config (with confirmation)
  - Esc: back to main menu

- [ ] 8.4 Write comprehensive unit tests in `internal/cli/models/cv_config_manager_test.go`
  - List display tests
  - Navigation tests
  - Action tests
  - Use Ginkgo/Gomega pattern

#### 9.0 Create CV Configuration Editor

- [ ] 9.1 Implement `CVConfigEditorModel` in `internal/cli/models/cv_config_editor.go`
  - Embed BaseStandardModel
  - Form fields: name, role, audience, date range, companies, tags, categories
  - Fields: inputs, roleSelector, audienceSelector, focusIndex, formErrors

- [ ] 9.2 Implement form fields
  - CV Name (required, text input)
  - Target Role (required, dropdown)
  - Target Audience (required, multi-select)
  - Date Range (optional, start/end date)
  - Companies (optional, multi-select)
  - Tags (optional, multi-select)
  - Competencies (optional, multi-select)

- [ ] 9.3 Implement field validation
  - CV Name: not empty, valid filename characters
  - Target Role: must select one
  - Target Audience: at least one
  - Date Range: end >= start if both provided

- [ ] 9.4 Implement keyboard navigation
  - Tab: move to next field
  - Shift+Tab: move to previous field
  - Arrow keys: navigate within multi-select
  - Space: toggle multi-select items
  - Enter: save config
  - Esc: cancel without saving

- [ ] 9.5 Implement save functionality
  - Call configManager.SaveConfig()
  - Show confirmation message
  - Return to config list

- [ ] 9.6 Write comprehensive unit tests in `internal/cli/models/cv_config_editor_test.go`
  - Field validation tests
  - Keyboard navigation tests
  - Save/load tests
  - Error handling tests
  - Use Ginkgo/Gomega pattern

#### 10.0 Create CV Generation & Preview Screen

- [ ] 10.1 Implement `CVGeneratorModel` in `internal/cli/models/cv_generator.go`
  - Embed BaseStandardModel
  - Show config summary before generating
  - Generate CV from selected config
  - Show loading indicator during generation
  - Handle errors gracefully
  - Display generated CV on completion

- [ ] 10.2 Implement `CVPreviewModel` in `internal/cli/models/cv_preview.go`
  - Embed BaseStandardModel
  - Display CV name, metadata (role, audience, date)
  - Display sections in order
  - Display bullets within sections
  - Use consistent styling

- [ ] 10.3 Implement section navigation
  - Arrow keys: navigate between sections
  - j/k: navigate bullets within section
  - Enter: show bullet sources
  - Esc: go back

- [ ] 10.4 Implement source event display
  - Show which events contributed to bullet
  - Show fact contributions if applicable
  - Display in expandable panel

- [ ] 10.5 Implement traceability indicators
  - Show source count next to each bullet
  - Highlight high-confidence bullets
  - Show inclusion reason on hover/expand

- [ ] 10.6 Implement export options
  - "Export as Text" (saves to file)
  - "Export as Markdown" (saves to file)
  - "Copy to Clipboard" (copies to clipboard)
  - Show file save location

- [ ] 10.7 Write comprehensive unit tests in `internal/cli/models/cv_generator_test.go` and `cv_preview_test.go`
  - Generation tests
  - Display tests
  - Navigation tests
  - Source tracing tests
  - Export tests
  - Use Ginkgo/Gomega pattern

#### 11.0 Create Supporting Models

- [ ] 11.1 Implement `RoleSelectorModel` in `internal/cli/models/role_selector.go`
  - Dropdown for role selection
  - Support keyboard navigation
  - Return selected role

- [ ] 11.2 Implement `AudienceConfiguratorModel` in `internal/cli/models/audience_configurator.go`
  - Multi-select for audience
  - Support keyboard navigation
  - Return selected audiences

- [ ] 11.3 Implement `SourceEventTracerModel` in `internal/cli/models/source_event_tracer.go`
  - Display source events/facts for bullet
  - Show detailed event information
  - Support navigation and scrolling

- [ ] 11.4 Write comprehensive unit tests for all supporting models
  - Use Ginkgo/Gomega pattern

### Phase 4: Integration with Existing Features

#### 12.0 Integrate CV Generation with Event Timeline

- [ ] 12.1 Add "Generate CV" option to event action menu
- [ ] 12.2 Support multi-event selection for CV generation
- [ ] 12.3 Navigate to CV configuration screen when option selected
- [ ] 12.4 Navigate to CV preview screen after generation completes
- [ ] 12.5 Allow navigation back to event timeline from CV preview
- [ ] 12.6 Write integration tests for event timeline → CV generation workflow

#### 13.0 Add CV Generation to Main Menu

- [ ] 13.1 Add "Manage CV Configs" option in main menu
- [ ] 13.2 Add "Generate CV" option in main menu
- [ ] 13.3 Add keyboard shortcut (e.g., 'c' for CV)
- [ ] 13.4 Add breadcrumb navigation support
- [ ] 13.5 Write integration tests for main menu → CV workflow

#### 14.0 Create CV Export Functionality

- [ ] 14.1 Implement `ExportService` in `internal/service/career/cv/export_service.go`

- [ ] 14.2 Implement plain text export
  - Export CV to plain text format
  - Include section headers and bullets
  - Include metadata (role, audience, date)
  - Optional: include source event references

- [ ] 14.3 Implement markdown export
  - Export CV to markdown format
  - Use markdown headers for sections
  - Use markdown lists for bullets
  - Include metadata as markdown comments

- [ ] 14.4 Implement file save functionality
  - Determine export directory (`$HOME/.kariya/cv_exports/` or similar)
  - Generate filename with timestamp
  - Handle file I/O errors
  - Show success message with file path

- [ ] 14.5 Implement clipboard copy functionality
  - Copy generated CV to clipboard
  - Show success message

- [ ] 14.6 Write comprehensive unit tests in `internal/service/career/cv/export_service_test.go`
  - Plain text export tests
  - Markdown export tests
  - Format validation tests
  - File save tests
  - Use Ginkgo/Gomega pattern

### Phase 5: Testing and Documentation

#### 15.0 Comprehensive Testing Suite

- [ ] 15.1 Write end-to-end tests for complete CV generation workflow
  - Config creation → generation → preview → export

- [ ] 15.2 Test bullet generation with various event/fact combinations
  - High/medium/low quality events

- [ ] 15.3 Test inclusion/exclusion criteria with edge cases
  - Aspirational language, inferred metrics, role inflation

- [ ] 15.4 Test ranking algorithm with different priority scenarios

- [ ] 15.5 Test compression logic
  - Bullet caps, older roles compress first

- [ ] 15.6 Test role-specific generation for all 4 roles
  - Principal, Staff, EM, Senior IC

- [ ] 15.7 Test audience-specific filtering for all 3 audiences
  - Hiring_manager, recruiter, peer

- [ ] 15.8 Test traceability system
  - All bullets trace to sources, no orphaned bullets

- [ ] 15.9 Test CV view management
  - Create, list, delete, regenerate

- [ ] 15.10 Test integration with event timeline, fact extraction, burst detection

- [ ] 15.11 Run race detector: `go test -race ./...`
  - Verify 0 race conditions

- [ ] 15.12 Verify code coverage meets 80%+ threshold across all new code

- [ ] 15.13 Performance test: CV generation ≤2s for ≤500 events

- [ ] 15.14 Performance test: bullet ranking ≤100ms for 1000 bullets

- [ ] 15.15 Performance test: traceability lookup ≤50ms per bullet

#### 16.0 Documentation and User Guidance

- [ ] 16.1 Create CV_GENERATION_GUIDE.md with comprehensive feature overview

- [ ] 16.2 Document YAML configuration format with examples

- [ ] 16.3 Document bullet generation rules
  - Inclusion/exclusion criteria, ranking priority

- [ ] 16.4 Document role-specific generation rules and bullet caps

- [ ] 16.5 Document audience-specific filtering rules

- [ ] 16.6 Provide examples of generated CVs for each role and audience

- [ ] 16.7 Document compression logic and how older roles are compressed

- [ ] 16.8 Document traceability system and how to view sources

- [ ] 16.9 Document export formats and use cases

- [ ] 16.10 Document keyboard shortcuts for CV screens

- [ ] 16.11 Update README.md with CV generation features

- [ ] 16.12 Update CLI_GUIDE.md with CV workflow shortcuts

- [ ] 16.13 Update CHANGELOG.md with feature description and test results

- [ ] 16.14 Create troubleshooting guide for common CV generation issues

---

## Implementation Guidelines

### Architecture Patterns

1. **Separation of Concerns**
   - Domain models in `internal/domain/career/`
   - Configuration management in `internal/service/career/cv/`
   - Bullet generation logic in `internal/service/career/cv/`
   - UI components in `internal/cli/models/`
   - **NO repository layer** for CVs (ephemeral only)

2. **Reuse Existing Patterns**
   - Follow BubbleTea Model pattern from existing screens
   - Use existing styling system
   - Follow validation patterns
   - Build on existing service layer architecture

3. **Domain-Driven Design**
   - CVView, CVSection, CVBullet as first-class domain concepts
   - Validation at domain level
   - Service layer orchestrates generation
   - **No persistence layer** (ephemeral only)

4. **Read-Only, Ephemeral Generation**
   - CVs generated on-demand, NOT stored permanently
   - No modification of source events or facts
   - Regeneration always reflects latest data
   - Traceability preserved for all bullets

5. **YAML Configuration**
   - Stored in `$HOME/.kariya/cv_configs/`
   - Human-readable format
   - Can be version-controlled
   - Can be edited directly

6. **Conservative Defaults**
   - Prefer explicit user choices
   - Err on side of fewer bullets
   - Strict validation
   - Clear source attribution

### Success Criteria (All Must Be Met)

- [ ] Users can generate role-specific CVs
- [ ] Users can generate audience-specific CVs
- [ ] Bullet generation respects inclusion/exclusion criteria
- [ ] Bullet ranking follows priority order
- [ ] Compression logic enforces bullet caps
- [ ] All bullets trace to ≥1 source
- [ ] No aspirational language in bullets
- [ ] No inferred metrics in bullets
- [ ] No role inflation in bullets
- [ ] Traceability system allows viewing sources
- [ ] CV generation ≤2s for ≤500 events
- [ ] All CVs can be exported to text and markdown
- [ ] **CVs are NOT stored in database** (ephemeral only)
- [ ] **Configurations are stored as YAML files**
- [ ] Code coverage ≥ 80%
- [ ] All tests passing (100% pass rate)
- [ ] Race detector passes (0 conditions)

---

## Estimated Effort

- Phase 1: 1 week (domain models, YAML configuration)
- Phase 2: 1.5 weeks (services)
- Phase 3: 1.5 weeks (UI components)
- Phase 4: 1 week (integration, export)
- Phase 5: 1 week (testing, documentation)

**Total**: 5-6 weeks

---

**Document Version**: 2.0 (Updated with YAML Configuration & Ephemeral Generation)
**Updated**: 2026-01-02 - Corrected to use YAML config and ephemeral generation (no database storage)

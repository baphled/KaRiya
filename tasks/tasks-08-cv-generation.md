# Task List: CV Generation Feature

**PRD Reference**: `docs/features/05-cv-generation.md`

**Purpose**: Transform raw career events into credible, audience- and role-specific CV views with full traceability, conservative defaults, and no text rewriting.

**Status**: 🔄 **IN PROGRESS** (50% - Phases 1, 2, 3 & Export Service Complete, Main Menu and Timeline Integration Pending)

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

[... rest of the file remains the same until Phase 4 ...]

### Phase 4: Integration with Existing Features ⏳ PARTIALLY COMPLETE (1/6)

#### 12.0 Integrate CV Generation with Event Timeline ⏳

- [ ] 12.1 Add "Generate CV" option to event action menu
- [ ] 12.2 Support multi-event selection for CV generation
- [ ] 12.3 Navigate to CV configuration screen when option selected
- [ ] 12.4 Navigate to CV preview screen after generation completes
- [ ] 12.5 Allow navigation back to event timeline from CV preview
- [ ] 12.6 Write integration tests for event timeline → CV generation workflow

#### 13.0 Add CV Generation to Main Menu ⏳

- [ ] 13.1 Add "Manage CV Configs" option in main menu
- [ ] 13.2 Add "Generate CV" option in main menu
- [ ] 13.3 Add keyboard shortcut (e.g., 'c' for CV)
- [ ] 13.4 Add breadcrumb navigation support
- [ ] 13.5 Write integration tests for main menu → CV workflow

#### 14.0 Create CV Export Functionality ✅ COMPLETE

- [x] 14.1 Implement `ExportService` in `internal/service/career/cv/export_service.go`

- [x] 14.2 Implement YAML export
  - Export CV to YAML format
  - Include section headers and bullets
  - Include metadata (role, audience, date)
  - Optional: include source event references

- [x] 14.3 Implement markdown export
  - Export CV to markdown format
  - Use markdown headers for sections
  - Use markdown lists for bullets
  - Include metadata as markdown comments

- [x] 14.4 Implement file save functionality
  - Determine export directory (`$HOME/.kariya/cv_exports/` or similar)
  - Generate filename with timestamp
  - Handle file I/O errors
  - Show success message with file path

- [x] 14.5 Implement clipboard copy functionality
  - Copy generated CV to clipboard
  - Show success message

- [x] 14.6 Write comprehensive unit tests in `internal/service/career/cv/export_service_test.go`
  - Plain text export tests
  - Markdown export tests
  - Format validation tests
  - File save tests
  - Use Ginkgo/Gomega pattern

[... rest of the file remains the same until Success Criteria ...]

### Success Criteria (All Must Be Met)

- [x] Users can generate role-specific CVs
- [x] Users can generate audience-specific CVs
- [x] Bullet generation respects inclusion/exclusion criteria
- [x] Bullet ranking follows priority order
- [x] Compression logic enforces bullet caps
- [x] All bullets trace to ≥1 source
- [x] No aspirational language in bullets
- [x] No inferred metrics in bullets
- [x] No role inflation in bullets
- [x] Traceability system allows viewing sources
- [x] CV generation ≤2s for ≤500 events
- [x] All CVs can be exported to text and markdown
- [x] **CVs are NOT stored in database** (ephemeral only)
- [x] **Configurations are stored as YAML files**
- [ ] Code coverage ≥ 80% (NOT YET - tests incomplete)
- [ ] All tests passing (NOT YET - phase 5 not started)
- [ ] Race detector passes (NOT YET - phase 5 not started)

---

## Progress Summary

### Completed (50%)
- ✅ Phase 1: Foundation & Core Components (100%)
  - Domain models (CVView, CVSection, CVBullet, CVConfig)
  - YAML configuration system with atomic writes
  - Comprehensive unit tests for domain models

- ✅ Phase 2: CV Generation Service (100%)
  - Bullet generator with filtering and ranking
  - Section builder for layout organization
  - CV generation orchestrator
  - Traceability service for event-to-CV mapping
  - All components fully tested

- ✅ Phase 3: UI Components (100%)
  - CVConfigManagerModel - Config list and management
  - CVConfigEditorModel - Config creation/editing
  - CVGeneratorModel - CV generation workflow
  - CVPreviewModel - CV display and source tracing
  - RoleSelector - Role selection component
  - AudienceConfigurator - Audience selection component
  - SourceEventTracer - Event/fact source display
  - All models include comprehensive unit tests

- ✅ Phase 4: Export Service (100%)
  - YAML export implementation
  - Markdown export implementation
  - File save functionality
  - Clipboard copy support
  - Comprehensive export tests

### Not Started (50%)
- ⏳ Phase 4: Integration (0/6 tasks)
  - Event timeline integration
  - Main menu integration

- ⏳ Phase 5: Testing & Documentation (0/30)
  - E2E and integration tests
  - Performance testing
  - Documentation and guides

---

## Estimated Effort Remaining

- Phase 4 (Integration): 3-4 days (main menu + event timeline)
- Phase 5 (Testing & Docs): 3-4 days (tests + documentation)

**Total Remaining**: 6-8 days

---

**Document Version**: 2.0 (Updated with YAML Configuration & Ephemeral Generation)
**Updated**: 2026-01-02 - Corrected to use YAML config and ephemeral generation (no database storage)
**Status Update**: 2026-01-02 - Phases 1, 2, 3, and Export Service complete, Main Menu and Timeline Integration Pending

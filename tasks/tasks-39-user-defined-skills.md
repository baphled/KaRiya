# Task 39: User-Defined Skills Management

## Overview
- **Goal**: Enable users to define, manage, and associate skills with career events as a prerequisite for technology-focused CV generation
- **Time Estimate**: 2-3 days
- **Prerequisites**: Understanding of domain models, repository pattern, intent architecture, huh forms

## Session Contract Acknowledgment
- [ ] Ran `make session-start` and it passed
- [ ] Acknowledge and commit to following all workflow rules
- [ ] Token count: _____ (must be < 50k to start)

## Pre-Task Checklist (MUST COMPLETE BEFORE STARTING)
- [ ] `make check-compliance` passes
- [ ] Reviewed existing patterns in:
  - `internal/domain/career/event.go` (domain model example)
  - `internal/repository/career/sqlite_repository.go` (repository pattern)
  - `internal/cli/intents/burst_management_intent.go` (intent example)
  - `internal/cli/forms/` (huh forms examples)
- [ ] Confirmed this is ONE atomic task (skills management system)
- [ ] Identified which test files will be created/modified

## Context

Currently, KaRiya has no user-defined skills. Skills are derived automatically from event tags and fact competencies during CV generation. This task creates a dedicated skills management system that will later enable technology-focused CV generation (Task 40).

## Files to Create

### Domain Layer
- [ ] `internal/domain/career/skill.go` - Skill domain model with validation
- [ ] `internal/domain/career/skill_test.go` - Domain model tests

### Repository Layer
- [ ] `internal/repository/career/skill_repository.go` - Repository interface
- [ ] `internal/repository/career/sqlite_skill_repository.go` - SQLite implementation
- [ ] `internal/repository/career/sqlite_skill_repository_test.go` - Repository tests
- [ ] `internal/repository/career/migrations/005_create_skills.sql` - Skills table migration
- [ ] `internal/repository/career/migrations/006_create_event_skills.sql` - Junction table migration

### Intent Layer
- [ ] `internal/cli/intents/manage_skills.go` - Intent types and states
- [ ] `internal/cli/intents/manage_skills_intent.go` - Intent implementation
- [ ] `internal/cli/intents/manage_skills_test.go` - Intent tests

### Forms
- [ ] `internal/cli/forms/skill_form.go` - Skill add/edit form configuration

## Files to Modify

- [ ] `internal/domain/career/event.go` - Add Skills field
- [ ] `internal/repository/career/sqlite_repository.go` - Handle skill associations in GetByID, List, Create, Update
- [ ] `internal/cli/models/form.go` - Add optional skills multi-select field
- [ ] `internal/cli/forms/metadata_form.go` - Add skills field to metadata editor
- [ ] `internal/cli/app/app.go` - Register ManageSkills intent, add to menu
- [ ] `internal/service/career/cv/data_processing_service.go` - Use user-defined skills

## Implementation Plan

### Phase 1: Domain Model & Migrations

**Goal**: Create Skill domain model and database schema

#### Domain Model
```go
// internal/domain/career/skill.go

type Skill struct {
    ID          string     `json:"id"`
    Name        string     `json:"name"`         // "Ruby", "Kubernetes", "React"
    Category    string     `json:"category"`     // User-definable, suggestions: "backend", "frontend", "devops", "database", "other"
    Level       string     `json:"level"`        // "beginner", "intermediate", "advanced", "expert" (optional, can be derived)
    YearsUsed   *int       `json:"years_used"`   // Optional
    LastUsed    *time.Time `json:"last_used"`    // Optional, can be derived from events
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
}

// Common categories (suggestions, not enforced - user can define custom)
var CommonSkillCategories = []string{
    "backend",
    "frontend",
    "devops",
    "database",
    "cloud",
    "tooling",
    "other",
}
```

#### Validation Rules
- Name: required, 1-100 characters, unique (case-insensitive)
- Category: required, 1-50 characters
- Level: optional, one of: beginner, intermediate, advanced, expert (empty string allowed)
- YearsUsed: optional, 0-50 range
- LastUsed: optional, cannot be in future

#### Migration 005: Skills Table
```sql
-- +goose Up
CREATE TABLE IF NOT EXISTS skills (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    category TEXT NOT NULL,
    level TEXT,
    years_used INTEGER,
    last_used DATETIME,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_skills_category ON skills(category);
CREATE INDEX IF NOT EXISTS idx_skills_name ON skills(name);

-- +goose Down
DROP INDEX IF EXISTS idx_skills_name;
DROP INDEX IF EXISTS idx_skills_category;
DROP TABLE IF EXISTS skills;
```

#### Migration 006: Event-Skills Junction Table
```sql
-- +goose Up
CREATE TABLE IF NOT EXISTS event_skills (
    event_id TEXT NOT NULL,
    skill_id TEXT NOT NULL,
    PRIMARY KEY (event_id, skill_id),
    FOREIGN KEY (event_id) REFERENCES career_events(id) ON DELETE CASCADE,
    FOREIGN KEY (skill_id) REFERENCES skills(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_event_skills_event ON event_skills(event_id);
CREATE INDEX IF NOT EXISTS idx_event_skills_skill ON event_skills(skill_id);

-- +goose Down
DROP INDEX IF EXISTS idx_event_skills_skill;
DROP INDEX IF EXISTS idx_event_skills_event;
DROP TABLE IF EXISTS event_skills;
```

#### Update CareerEvent Domain Model
```go
// internal/domain/career/event.go

type CareerEvent struct {
    // ... existing fields ...
    Skills     []string  `json:"skills,omitempty"`     // Skill IDs (NEW)
}
```

**TDD Checklist - Phase 1:**
- [ ] Write failing test: Skill validation (name required)
- [ ] Test passes
- [ ] Write failing test: Skill validation (name unique, case-insensitive)
- [ ] Test passes
- [ ] Write failing test: Skill validation (category required)
- [ ] Test passes
- [ ] Write failing test: Skill validation (level optional, enum)
- [ ] Test passes
- [ ] Write failing test: Skill validation (years range)
- [ ] Test passes
- [ ] Write failing test: Skill validation (lastUsed not future)
- [ ] Test passes
- [ ] Write failing test: Migration 005 creates skills table
- [ ] Test passes
- [ ] Write failing test: Migration 006 creates junction table
- [ ] Test passes
- [ ] Commit: `test(skill): add domain validation tests`
- [ ] Commit: `feat(skill): add Skill domain model`
- [ ] Commit: `feat(migrations): add skills and event_skills tables`

### Phase 2: Repository Layer

**Goal**: Implement persistence for skills and event-skill associations

#### SkillRepository Interface
```go
// internal/repository/career/skill_repository.go

type SkillRepository interface {
    Create(ctx context.Context, skill *career.Skill) error
    GetByID(ctx context.Context, id string) (*career.Skill, error)
    List(ctx context.Context, filters *SkillFilters) ([]*career.Skill, error)
    Update(ctx context.Context, skill *career.Skill) error
    Delete(ctx context.Context, id string) error
    GetByName(ctx context.Context, name string) (*career.Skill, error)
    GetByCategory(ctx context.Context, category string) ([]*career.Skill, error)
    GetSkillsForEvent(ctx context.Context, eventID string) ([]*career.Skill, error)
    GetEventCountsForSkills(ctx context.Context) (map[string]int, error)
}

type SkillFilters struct {
    Category string
    Level    string
}
```

#### Update EventRepository
- Modify `GetByID()` to load associated skill IDs
- Modify `List()` to load associated skill IDs
- Modify `Create()` to save skill associations
- Modify `Update()` to update skill associations (delete old, insert new)
- Modify `Delete()` to cascade delete associations (handled by FK)

**TDD Checklist - Phase 2:**
- [ ] Write failing test: SkillRepository.Create
- [ ] Test passes
- [ ] Write failing test: SkillRepository.Create duplicate name fails
- [ ] Test passes
- [ ] Write failing test: SkillRepository.GetByID
- [ ] Test passes
- [ ] Write failing test: SkillRepository.List (no filters)
- [ ] Test passes
- [ ] Write failing test: SkillRepository.List (category filter)
- [ ] Test passes
- [ ] Write failing test: SkillRepository.List (level filter)
- [ ] Test passes
- [ ] Write failing test: SkillRepository.Update
- [ ] Test passes
- [ ] Write failing test: SkillRepository.Delete
- [ ] Test passes
- [ ] Write failing test: SkillRepository.GetByName
- [ ] Test passes
- [ ] Write failing test: SkillRepository.GetByCategory
- [ ] Test passes
- [ ] Write failing test: SkillRepository.GetSkillsForEvent
- [ ] Test passes
- [ ] Write failing test: SkillRepository.GetEventCountsForSkills
- [ ] Test passes
- [ ] Write failing test: EventRepository.GetByID loads skill IDs
- [ ] Test passes
- [ ] Write failing test: EventRepository.Create saves skill associations
- [ ] Test passes
- [ ] Write failing test: EventRepository.Update updates skill associations
- [ ] Test passes
- [ ] Commit: `test(skill): add repository tests`
- [ ] Commit: `feat(skill): implement SkillRepository`
- [ ] Commit: `feat(event): add skill associations to EventRepository`

### Phase 3: Skill Form Configuration

**Goal**: Create huh form for adding/editing skills with suggestions

#### Skill Form
```go
// internal/cli/forms/skill_form.go

type SkillFormData struct {
    Name      string
    Category  string
    Level     string
    YearsUsed string
}

// NewSkillForm creates form with category suggestions and skill name suggestions
func NewSkillForm(existingSkill *career.Skill, skillSuggestions []string) *huh.Form

// ApplySkillFormData validates and applies form data to skill
func ApplySkillFormData(skill *career.Skill, data *SkillFormData) error
```

#### Validators
- SkillName: 1-100 chars, required, trimmed
- SkillCategory: 1-50 chars, required, trimmed
- SkillLevel: optional (empty allowed), one of: beginner, intermediate, advanced, expert
- YearsUsed: optional (empty allowed), integer 0-50

#### Skill Suggestions
- Extract skill-like terms from event text (similar to technology extraction planned for Task 40)
- Show as autocomplete suggestions when adding new skill
- User can type custom skill name

**TDD Checklist - Phase 3:**
- [ ] Write failing test: SkillForm creation (new skill)
- [ ] Test passes
- [ ] Write failing test: SkillForm creation (edit existing)
- [ ] Test passes
- [ ] Write failing test: SkillName validator (required)
- [ ] Test passes
- [ ] Write failing test: SkillName validator (length)
- [ ] Test passes
- [ ] Write failing test: SkillName validator (trimmed)
- [ ] Test passes
- [ ] Write failing test: SkillCategory validator
- [ ] Test passes
- [ ] Write failing test: SkillLevel validator (optional)
- [ ] Test passes
- [ ] Write failing test: SkillLevel validator (enum)
- [ ] Test passes
- [ ] Write failing test: YearsUsed validator (optional)
- [ ] Test passes
- [ ] Write failing test: YearsUsed validator (range)
- [ ] Test passes
- [ ] Write failing test: ApplySkillFormData (new skill)
- [ ] Test passes
- [ ] Write failing test: ApplySkillFormData (update existing)
- [ ] Test passes
- [ ] Commit: `test(forms): add skill form tests`
- [ ] Commit: `feat(forms): add skill form configuration`

### Phase 4: Manage Skills Intent

**Goal**: Create TUI for managing skills

#### States
```go
const (
    SkillsStateList   SkillsState = "list"     // View all skills grouped by category
    SkillsStateAdd    SkillsState = "add"      // Add new skill (huh form)
    SkillsStateEdit   SkillsState = "edit"     // Edit existing skill (huh form)
    SkillsStateDelete SkillsState = "delete"   // Confirm deletion
)
```

#### View Structure
- List: Skills grouped by category, sorted alphabetically, shows level and years if present
- Add: Huh form for new skill with suggestions
- Edit: Huh form pre-populated with existing skill
- Delete: Confirmation modal

#### Keyboard Shortcuts
- List: `n` add, `e` edit, `d` delete, `j/k` navigate, `Esc` back
- Add/Edit: Form navigation (Tab, Enter, Esc)
- Delete: `y` confirm, `n` cancel, `Esc` cancel

**TDD Checklist - Phase 4:**
- [ ] Write failing test: ManageSkillsIntent.Init loads skills
- [ ] Test passes
- [ ] Write failing test: Navigate list with j/k
- [ ] Test passes
- [ ] Write failing test: Press n transitions to Add state
- [ ] Test passes
- [ ] Write failing test: Add skill form submission creates skill
- [ ] Test passes
- [ ] Write failing test: Add skill form cancel returns to list
- [ ] Test passes
- [ ] Write failing test: Press e transitions to Edit state
- [ ] Test passes
- [ ] Write failing test: Edit skill form submission updates skill
- [ ] Test passes
- [ ] Write failing test: Edit skill form cancel returns to list
- [ ] Test passes
- [ ] Write failing test: Press d transitions to Delete state
- [ ] Test passes
- [ ] Write failing test: Delete confirmation removes skill
- [ ] Test passes
- [ ] Write failing test: Delete cancel returns to list
- [ ] Test passes
- [ ] Write failing test: List view groups skills by category
- [ ] Test passes
- [ ] Write failing test: List view shows empty state
- [ ] Test passes
- [ ] Write failing test: Escape from list returns result
- [ ] Test passes
- [ ] Commit: `test(skills): add ManageSkills intent tests`
- [ ] Commit: `feat(skills): implement ManageSkills intent`

### Phase 5: Event Capture Integration

**Goal**: Allow associating skills with events during capture (optional field)

#### Form Updates
- Add skills multi-select field after Categories
- Visible in both quick and manual modes, but optional
- Load user's skills for selection
- Save skill associations on event creation/update

#### Metadata Editor
- Add skills field to metadata form
- Allow editing skill associations on existing events

**TDD Checklist - Phase 5:**
- [ ] Write failing test: FormModel includes skills field
- [ ] Test passes
- [ ] Write failing test: Skills field visible in quick mode
- [ ] Test passes
- [ ] Write failing test: Skills field visible in manual mode
- [ ] Test passes
- [ ] Write failing test: Skills field is optional
- [ ] Test passes
- [ ] Write failing test: Skills saved on event creation
- [ ] Test passes
- [ ] Write failing test: Skills loaded on event edit
- [ ] Test passes
- [ ] Write failing test: MetadataForm includes skills field
- [ ] Test passes
- [ ] Write failing test: Skills updated via metadata editor
- [ ] Test passes
- [ ] Commit: `test(forms): add skill selection to event forms`
- [ ] Commit: `feat(forms): add skill selection to event capture`
- [ ] Commit: `feat(forms): add skill selection to metadata editor`

### Phase 6: App Integration

**Goal**: Register intent and add to main menu

#### Menu Integration
- Add "Manage Skills" option to main menu (after "Browse Timeline")
- Register ManageSkills intent with router
- Add keyboard shortcut (s for skills)

**TDD Checklist - Phase 6:**
- [ ] Write failing test: ManageSkills appears in menu
- [ ] Test passes
- [ ] Write failing test: ManageSkills intent registered with router
- [ ] Test passes
- [ ] Write failing test: Navigate to ManageSkills from menu
- [ ] Test passes
- [ ] Write failing test: Result handler for ManageSkills
- [ ] Test passes
- [ ] Commit: `test(app): add ManageSkills intent integration tests`
- [ ] Commit: `feat(app): integrate ManageSkills intent and menu`

### Phase 7: CV Generation Integration

**Goal**: Use user-defined skills in CV generation

#### DataProcessingService Updates
- ExtractSkills() prioritizes user's defined skills over derived skills
- Skills section shows user-defined skills with event counts
- Derive LastUsed from event dates (most recent event with that skill)
- Show skills grouped by category

**TDD Checklist - Phase 7:**
- [ ] Write failing test: ExtractSkills includes user-defined skills
- [ ] Test passes
- [ ] Write failing test: User-defined skills prioritized over derived
- [ ] Test passes
- [ ] Write failing test: Skills section shows defined skills grouped by category
- [ ] Test passes
- [ ] Write failing test: Event count per skill
- [ ] Test passes
- [ ] Write failing test: LastUsed derived from events
- [ ] Test passes
- [ ] Commit: `test(cv): add user-defined skills to CV generation`
- [ ] Commit: `feat(cv): integrate user-defined skills in CV generation`

## Pre-Commit Checklist (BEFORE EACH COMMIT)
- [ ] `make review-commit` passes
- [ ] AI attribution included (if AI-generated)
- [ ] Commit message explains **WHY**, not just WHAT
- [ ] Commit is atomic (ONE logical change)

## Post-Task Checklist (MUST COMPLETE BEFORE NEXT TASK)
- [ ] `make check-compliance` passes
- [ ] All checkboxes above completed
- [ ] Task marked complete `[x]` in task file
- [ ] Token count: _____ (< 100k to continue)

## Acceptance Criteria
- [ ] Users can add, edit, and delete skills via Manage Skills intent
- [ ] Skills are persisted in database
- [ ] Skills can be associated with events during capture (optional field, visible in both quick and manual modes)
- [ ] Skills can be edited via metadata editor
- [ ] Skills appear in CV generation (skills section, grouped by category)
- [ ] Skill suggestions shown when adding new skill (extracted from event text)
- [ ] All tests pass (100% pass rate)
- [ ] Coverage maintained ≥ 80%
- [ ] Zero staticcheck warnings
- [ ] Zero race conditions
- [ ] Documentation updated (add docs/SKILLS_GUIDE.md)

## Rollback Plan
- Migrations can be rolled back via `goose down`
- New tables (skills, event_skills) can be dropped without affecting existing data
- CareerEvent.Skills field is optional (existing events unaffected)
- FormModel skills field is optional (existing capture flow unaffected)

## Documentation to Create

### docs/SKILLS_GUIDE.md
- How to manage skills
- How to associate skills with events
- How skills appear in CVs
- Category suggestions
- Best practices for skill management

## Notes

### UI Design Principles
- Follow TUI_STANDARDS.md for keyboard shortcuts
- Use StandardView for all screens
- Use huh forms for add/edit operations
- Group skills by category in list view
- Show skill count and event count per category
- Empty state shows helpful message with "n to add skill"

### Category Suggestions
Common categories to suggest (user can define custom):
- backend - Server-side languages and frameworks
- frontend - UI technologies and frameworks
- devops - Infrastructure, deployment, CI/CD
- database - Database systems and query languages
- cloud - Cloud platforms and services
- mobile - Mobile development technologies
- tooling - Development tools and utilities
- other - Uncategorized skills

### Skill Derivation (Implemented in Phase 7)
- LastUsed: Most recent event date where skill is associated
- Event count: Number of events with that skill
- Level: Optional, user can set manually or leave for future auto-derivation

### Skill Suggestions Algorithm
When adding a new skill, suggest technologies found in event text:
1. Scan all event text for capitalized technical terms
2. Scan for acronyms (API, CI/CD, AWS, etc.)
3. Scan for framework patterns (.js, SQL suffixes)
4. Filter out company names and common words
5. Show top 10 suggestions sorted by frequency
6. User can select from suggestions or type custom

### Integration with Task 40
This task creates the foundation for Task 40 (Role Emphasis Redesign), which will:
- Extract technologies from user's defined skills
- Allow selecting technologies for CV generation (Generalist: 2-5, Specialist: 1)
- Filter/prioritize CV bullets based on technology selection
- Determine Focus Area from skill categories

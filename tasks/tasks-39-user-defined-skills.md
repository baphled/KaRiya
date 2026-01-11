# Task 39: User-Defined Skills Management

## Overview
- **Goal**: Enable users to define, manage, and associate skills with career events as a prerequisite for technology-focused CV generation
- **Time Estimate**: 2-3 days
- **Prerequisites**: Understanding of domain models, repository pattern, intent architecture, huh forms

## Session Contract Acknowledgment
- [x] Ran `make session-start` and it passed
- [x] Acknowledge and commit to following all workflow rules
- [x] Token count: 37866 (must be < 50k to start)

## Pre-Task Checklist (MUST COMPLETE BEFORE STARTING)
- [x] `make check-compliance` passes
- [x] Reviewed existing patterns in:
  - `internal/domain/career/event.go` (domain model example)
  - `internal/repository/career/sqlite_repository.go` (repository pattern)
  - `internal/cli/intents/burst_management_intent.go` (intent example)
  - `internal/cli/forms/` (huh forms examples)
- [x] Confirmed this is ONE atomic task (skills management system)
- [x] Identified which test files will be created/modified

## Context

Currently, KaRiya has no user-defined skills. Skills are derived automatically from event tags and fact competencies during CV generation. This task creates a dedicated skills management system that will later enable technology-focused CV generation (Task 40).

## Files to Create

### Domain Layer
- [x] `internal/domain/career/skill.go` - Skill domain model with validation
- [x] `internal/domain/career/skill_test.go` - Domain model tests

### Repository Layer
- [x] `internal/repository/career/skill_repository.go` - Repository interface
- [x] `internal/repository/career/sqlite_skill_repository.go` - SQLite implementation
- [x] `internal/repository/career/sqlite_skill_repository_test.go` - Repository tests
- [x] `internal/repository/career/migrations/005_create_skills.sql` - Skills table migration
- [x] `internal/repository/career/migrations/006_create_event_skills.sql` - Junction table migration

### Intent Layer
- [x] `internal/cli/intents/manage_skills.go` - Intent types and states
- [x] `internal/cli/intents/manage_skills_intent.go` - Intent implementation
- [x] `internal/cli/intents/manage_skills_test.go` - Intent tests

### Forms
- [x] `internal/cli/forms/skill_form.go` - Skill add/edit form configuration

### Components (NEW - for detail and events views)
- [ ] `internal/cli/intents/manage_skills_detail_view.go` - Detail view rendering helper
- [ ] `internal/cli/intents/manage_skills_events_view.go` - Events view rendering helper

## Files to Modify

- [x] `internal/domain/career/event.go` - Add Skills field
- [x] `internal/repository/career/sqlite_repository.go` - Handle skill associations in GetByID, List, Create, Update
- [ ] `internal/repository/career/skill_repository.go` - Add GetEventCountsForSkills, GetLastUsedForSkills (for detail view)
- [ ] `internal/cli/models/form.go` - Add optional skills multi-select field
- [ ] `internal/cli/forms/metadata_form.go` - Add skills field to metadata editor
- [x] `internal/cli/app/app.go` - Register ManageSkills intent, add to menu
- [x] `internal/service/career/service.go` - Add SkillRepository to Service
- [x] `cmd/cli/main.go` - Initialize SkillRepository
- [ ] `internal/service/career/cv/data_processing_service.go` - Use user-defined skills
- [ ] `internal/cli/importer/parser.go` - Add Skills column parsing
- [ ] `internal/cli/importer/parser_test.go` - Add Skills parsing tests
- [ ] `docs/CSV_FORMAT_GUIDE.md` - Document Skills column
- [ ] `docs/CSV_IMPORT_GUIDE.md` - Add Skills import examples

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
- [x] Write failing test: Skill validation (name required)
- [x] Test passes
- [x] Write failing test: Skill validation (name unique, case-insensitive) - deferred to repository layer
- [x] Test passes
- [x] Write failing test: Skill validation (category required)
- [x] Test passes
- [x] Write failing test: Skill validation (level optional, enum)
- [x] Test passes
- [x] Write failing test: Skill validation (years range)
- [x] Test passes
- [x] Write failing test: Skill validation (lastUsed not future)
- [x] Test passes
- [x] Write failing test: Migration 005 creates skills table - will test via migration test updates
- [x] Test passes
- [x] Write failing test: Migration 006 creates junction table - will test via migration test updates
- [x] Test passes
- [x] Commit: `test(domain): add Skill validation tests` (519d149)
- [x] Commit: `feat(domain): add Skill domain model` - merged with tests commit (519d149)
- [x] Commit: `feat(repo): add skills and event_skills tables` (510b890)

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
- [x] Write failing test: SkillRepository.Create
- [x] Test passes
- [x] Write failing test: SkillRepository.Create duplicate name fails
- [x] Test passes
- [x] Write failing test: SkillRepository.GetByID
- [x] Test passes
- [x] Write failing test: SkillRepository.List (no filters)
- [x] Test passes
- [x] Write failing test: SkillRepository.List (category filter)
- [x] Test passes
- [x] Write failing test: SkillRepository.List (level filter)
- [x] Test passes
- [x] Write failing test: SkillRepository.Update
- [x] Test passes
- [x] Write failing test: SkillRepository.Delete
- [x] Test passes
- [x] Write failing test: SkillRepository.GetByName
- [x] Test passes
- [x] Write failing test: SkillRepository.GetByCategory
- [x] Test passes
- [x] Write failing test: SkillRepository.GetSkillsForEvent
- [x] Test passes
- [x] Write failing test: SkillRepository.GetEventCountsForSkills
- [x] Test passes
- [x] Write failing test: EventRepository.GetByID loads skill IDs
- [x] Test passes
- [x] Write failing test: EventRepository.Create saves skill associations
- [x] Test passes
- [x] Write failing test: EventRepository.Update updates skill associations
- [x] Test passes
- [x] Commit: `test(repo): add SkillRepository tests (TDD RED)` (f07b7a1)
- [x] Commit: `feat(repo): implement SkillRepository` - merged with above commit
- [x] Commit: `feat(repo): add event-skill associations` (1f5c8f3)

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
- [x] Write failing test: SkillForm creation (new skill)
- [x] Test passes
- [x] Write failing test: SkillForm creation (edit existing)
- [x] Test passes
- [x] Write failing test: SkillName validator (required)
- [x] Test passes
- [x] Write failing test: SkillName validator (length)
- [x] Test passes
- [x] Write failing test: SkillName validator (trimmed)
- [x] Test passes
- [x] Write failing test: SkillCategory validator
- [x] Test passes
- [x] Write failing test: SkillLevel validator (optional)
- [x] Test passes
- [x] Write failing test: SkillLevel validator (enum)
- [x] Test passes
- [x] Write failing test: YearsUsed validator (optional)
- [x] Test passes
- [x] Write failing test: YearsUsed validator (range)
- [x] Test passes
- [x] Write failing test: ApplySkillFormData (new skill)
- [x] Test passes
- [x] Write failing test: ApplySkillFormData (update existing)
- [x] Test passes
- [x] Commit: `test(cli): add skill form tests` (ea92cc7)
- [x] Commit: `feat(cli): implement skill form configuration` (8eb7999)

### Phase 4: Manage Skills Intent

**Goal**: Create TUI for managing skills

#### States
```go
const (
    SkillsStateList         SkillsState = "list"          // View all skills grouped by category
    SkillsStateDetail       SkillsState = "detail"        // View single skill details (NEW - matches bursts/facts)
    SkillsStateDetailEvents SkillsState = "detail_events" // View events using this skill (NEW - matches bursts)
    SkillsStateAdd          SkillsState = "add"           // Add new skill (huh form)
    SkillsStateEdit         SkillsState = "edit"          // Edit existing skill (huh form)
    SkillsStateDelete       SkillsState = "delete"        // Confirm deletion
)
```

#### View Structure

**List View** (SkillsStateList):
- Skills grouped by category, sorted alphabetically
- Table shows: Name, Level, Years Used, Event Count
- Pagination: 15 skills per page
- Focus indicator: ▶ marker on selected row
- Empty state: "No skills yet. Press 'n' to add your first skill."

**Detail View** (SkillsStateDetail) - **NEW**:
- Shows single skill with full details:
  - Name
  - Category
  - Level (if set)
  - Years Used (if set)
  - Last Used (derived from events)
  - Event Count (number of events with this skill)
  - Created/Updated timestamps
- Actions available: View events (`Enter`), Edit (`e`), Delete (`d`)

**Detail Events View** (SkillsStateDetailEvents) - **NEW**:
- Shows all events that use this skill
- Table format: Date, Event Text (truncated), Company
- Pagination: 10 events per page
- Allows navigation back to detail view (Esc)

**Add/Edit** (SkillsStateAdd/Edit):
- Huh form for skill input with suggestions
- Fields: Name, Category, Level, Years Used

**Delete** (SkillsStateDelete):
- Confirmation modal showing skill name and event count
- Warning if skill is used by events
- Options: Confirm (`y`), Cancel (`n/Esc`)

#### Keyboard Shortcuts

**List View**:
- `j/k` or `↑/↓` - Navigate skills
- `Enter` - View skill detail (NEW - matches bursts/facts)
- `n` - Add new skill
- `e` - Edit selected skill (quick edit from list)
- `d` - Delete selected skill (quick delete from list)
- `Esc` - Back to main menu

**Detail View** (NEW):
- `Enter` - View events using this skill (NEW - natural progression, like list→detail)
- `e` - Edit skill
- `d` - Delete skill
- `Esc` - Back to list

**Detail Events View** (NEW):
- `j/k` or `↑/↓` - Navigate events
- `Esc` - Back to detail view

**Add/Edit Form**:
- `Tab` - Next field
- `Shift+Tab` - Previous field
- `Enter` - Submit (when on submit button)
- `Ctrl+S` - Save (from any field)
- `Esc` - Cancel

**Delete Confirmation**:
- `y` - Confirm deletion
- `n/Esc` - Cancel

**TDD Checklist - Phase 4:**

#### Basic List Operations (Completed):
- [x] Write failing test: ManageSkillsIntent.Init loads skills
- [x] Test passes
- [x] Write failing test: Navigate list with j/k
- [x] Test passes
- [x] Write failing test: Press n transitions to Add state
- [x] Test passes
- [x] Write failing test: Add skill form submission creates skill
- [x] Test passes
- [x] Write failing test: Add skill form cancel returns to list
- [x] Test passes
- [x] Write failing test: Press e transitions to Edit state (quick edit from list)
- [x] Test passes
- [x] Write failing test: Edit skill form submission updates skill
- [x] Test passes
- [x] Write failing test: Edit skill form cancel returns to list
- [x] Test passes
- [x] Write failing test: Press d transitions to Delete state (quick delete from list)
- [x] Test passes
- [x] Write failing test: Delete confirmation removes skill
- [x] Test passes
- [x] Write failing test: Delete cancel returns to list
- [x] Test passes
- [x] Write failing test: List view groups skills by category
- [x] Test passes
- [x] Write failing test: List view shows empty state
- [x] Test passes
- [x] Write failing test: Escape from list returns result
- [x] Test passes

#### Detail View Operations (NEW - To Match Bursts/Facts):
- [ ] Write failing test: Press Enter from list transitions to Detail state
- [ ] Test passes
- [ ] Write failing test: Detail view shows skill name, category, level, years used
- [ ] Test passes
- [ ] Write failing test: Detail view shows event count for skill
- [ ] Test passes
- [ ] Write failing test: Detail view shows last used date (derived from events)
- [ ] Test passes
- [ ] Write failing test: Press Enter from detail transitions to DetailEvents state
- [ ] Test passes
- [ ] Write failing test: DetailEvents view shows all events using skill
- [ ] Test passes
- [ ] Write failing test: DetailEvents view paginates events (10 per page)
- [ ] Test passes
- [ ] Write failing test: Press Esc from DetailEvents returns to Detail
- [ ] Test passes
- [ ] Write failing test: Press e from detail transitions to Edit state
- [ ] Test passes
- [ ] Write failing test: Press d from detail transitions to Delete state
- [ ] Test passes
- [ ] Write failing test: Press Esc from detail returns to List
- [ ] Test passes
- [ ] Write failing test: List view shows event count per skill
- [ ] Test passes

#### Commits:
- [x] Commit: `test(skills): add ManageSkills intent tests (basic list)` (already done)
- [x] Commit: `feat(skills): implement ManageSkills intent (basic list)` (already done)
- [ ] Commit: `test(skills): add detail view and events view tests (TDD RED)`
- [ ] Commit: `feat(skills): add detail view and events view states`

### Phase 4B: Filter and Sort (Optional Enhancement)

**Goal**: Add filtering and sorting capabilities to skill list (matches planned burst features)

**Priority**: ⚠️ **OPTIONAL** - Can be deferred to future task if time-constrained

#### Filter Options
- **By Category**: Show only skills in specific category (backend, frontend, devops, etc.)
- **By Level**: Show only skills with specific level (beginner, intermediate, advanced, expert)
- **By Usage**: Show only skills with >0 events (hide unused skills)

#### Sort Options
- **By Name** (A-Z, Z-A) - Default
- **By Event Count** (Most used first, Least used first)
- **By Last Used** (Most recent first, Oldest first)
- **By Category** (Grouped view - already default)

#### UI Implementation
- Add filter/sort bar below header (above skill table)
- Show active filters with clear button
- Keyboard shortcuts:
  - `f` - Open filter menu
  - `s` - Open sort menu
  - `x` - Clear all filters

#### Repository Enhancement
```go
type SkillFilters struct {
    Category    string   // Filter by category
    Level       string   // Filter by level
    MinEvents   int      // Minimum event count (e.g., >0 for used skills only)
    SortBy      string   // "name", "events", "last_used", "category"
    SortOrder   string   // "asc", "desc"
}

// Update List signature
List(ctx context.Context, filters *SkillFilters) ([]*career.Skill, error)
```

**TDD Checklist - Phase 4B (Optional):**
- [ ] Write failing test: Filter by category
- [ ] Test passes
- [ ] Write failing test: Filter by level
- [ ] Test passes
- [ ] Write failing test: Filter by min events (used skills only)
- [ ] Test passes
- [ ] Write failing test: Sort by name (A-Z, Z-A)
- [ ] Test passes
- [ ] Write failing test: Sort by event count
- [ ] Test passes
- [ ] Write failing test: Sort by last used
- [ ] Test passes
- [ ] Write failing test: Clear filters returns all skills
- [ ] Test passes
- [ ] Write failing test: Press f opens filter menu
- [ ] Test passes
- [ ] Write failing test: Press s opens sort menu
- [ ] Test passes
- [ ] Commit: `test(skills): add filter and sort tests (TDD RED)`
- [ ] Commit: `feat(skills): add filter and sort capabilities`

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

### Phase 5B: CSV Import Integration

**Goal**: Allow importing skills with events via CSV files

#### CSV Format Extension

Add optional **Skills** column to CSV import format:

```csv
Text,Date,Categories,Tags,Project,Company,Skills
"Architected platform migration",2024-01,Technical,technical;architecture,Platform,TechCorp,"Go;Kubernetes;PostgreSQL"
"Built React dashboard",2024-02,Technical,technical,Dashboard,TechCorp,"React;TypeScript;Redux"
```

#### Skills Column Specification

- **Format**: Semicolon-separated list of skill names
- **Optional**: Empty column is valid
- **Matching**: Case-insensitive skill name lookup
- **Auto-Creation**: Skills not found are created with default category "other"
- **Examples**:
  - `Go;Kubernetes;PostgreSQL` (3 skills)
  - `React;TypeScript` (2 skills)
  - `Ruby` (1 skill)
  - (empty - optional)

#### Implementation Strategy

**Approach: Auto-Create Skills with Default Category** ✅ **RECOMMENDED**

When parsing Skills column:
1. Split by semicolon, trim whitespace
2. For each skill name:
   - Look up existing skill by name (case-insensitive)
   - If found: Use existing skill ID
   - If not found: Create new skill with:
     - Name: Trimmed skill name
     - Category: "other" (default for auto-created skills)
     - Level: Empty (optional)
   - Add skill ID to event.Skills array
3. User can later edit skill categories via Manage Skills intent

**Benefits**:
- ✅ Seamless bulk imports (no need to pre-create skills)
- ✅ Skill names from CSV are preserved exactly
- ✅ User can refine categories later via Manage Skills
- ✅ Simple, predictable behavior

#### Parser Changes

**Update CSVParser struct**:
```go
type CSVParser struct {
    existingEvents  *[]*career.CareerEvent
    skillRepository career.SkillRepository  // NEW: need skill repository
    dateFormats     []string
    categoryMapper  *CategoryMapper
    tagMapper       *TagMapper
    mapData         bool
    ctx             context.Context         // NEW: need context for repo calls
}
```

**Update Constructors**:
```go
func NewCSVParser(ctx context.Context, existingEvents []*career.CareerEvent, skillRepo career.SkillRepository) *CSVParser
func NewCSVParserWithMapping(ctx context.Context, existingEvents []*career.CareerEvent, skillRepo career.SkillRepository) *CSVParser
```

**Add Skills Parsing** (in `parseRow` function, after Categories parsing):
```go
// Parse Skills (optional, semicolon-separated skill names)
if skillsStr, ok := rawData["Skills"]; ok && strings.TrimSpace(skillsStr) != "" {
    rawSkills := strings.Split(skillsStr, ";")
    
    for _, skillName := range rawSkills {
        skillName = strings.TrimSpace(skillName)
        if skillName == "" {
            continue
        }
        
        // Look up skill by name (case-insensitive)
        skill, err := p.skillRepository.GetByName(p.ctx, skillName)
        if err != nil || skill == nil {
            // Skill doesn't exist - create it with default category
            newSkill := &career.Skill{
                Name:      skillName,
                Category:  "other", // Default category for auto-created skills
                Level:     "",      // Optional, user can set later
                CreatedAt: time.Now(),
                UpdatedAt: time.Now(),
            }
            
            // Generate ID
            newSkill.ID = generateID() // Use same ID generation as events
            
            if err := p.skillRepository.Create(p.ctx, newSkill); err != nil {
                parsedRow.ValidationErrors = append(parsedRow.ValidationErrors,
                    fmt.Sprintf("Failed to create skill '%s': %v", skillName, err))
                parsedRow.IsValid = false
                continue
            }
            skill = newSkill
        }
        
        // Add skill ID to event
        event.Skills = append(event.Skills, skill.ID)
    }
}
```

#### Integration Points

**Update Import Wizard Intent**:
```go
// internal/cli/intents/import_wizard_intent.go

// Pass skillRepository to parser
parser := importer.NewCSVParserWithMapping(
    context.Background(),
    existingEvents,
    i.context.SkillRepository, // NEW: pass skill repository
)
```

**Update Service**:
```go
// internal/service/career/service.go

// Ensure SkillRepository is available for import (already added in Phase 2)
type Service struct {
    // ... existing fields
    skillRepository career.SkillRepository // Already added
}
```

#### Documentation Updates

**CSV_FORMAT_GUIDE.md** - Add Skills section (after Company field):

```markdown
### Skills Field

- **Format**: Semicolon-separated list of skill names
- **Maximum**: No limit on number of skills per event
- **Examples**:
  - `Go;Kubernetes;PostgreSQL` (3 skills)
  - `React;TypeScript;Redux` (3 skills)
  - `Ruby` (1 skill)
  - (empty - optional)

- **Skill Matching**:
  Your CSV uses skill names like:
  - `Go;Kubernetes;Docker`
  - `React;TypeScript`
  - `Ruby;Rails;PostgreSQL`

- **KaRiya Behavior**:
  - Skills are matched by name (case-insensitive)
  - If skill doesn't exist, it's automatically created with category "other"
  - You can later edit skill categories via "Manage Skills" (press 's' from main menu)
  - Skills are associated with the imported event

- **Storage**: Skills stored as IDs, names matched automatically
```

**CSV_IMPORT_GUIDE.md** - Add Skills section (after Tags section):

```markdown
### Skills (Optional)

Skills can be included in CSV imports to associate technical skills with events.

**Format**: Semicolon-separated list of skill names

**Example**:
```csv
Text,Date,Categories,Tags,Project,Company,Skills
"Architected microservices platform",2024-01,Technical,technical;architecture,Platform,TechCorp,"Go;Kubernetes;PostgreSQL;Docker"
"Built React dashboard with TypeScript",2024-02,Technical,technical,Dashboard,TechCorp,"React;TypeScript;Redux;CSS"
```

**Behavior**:
- Skills are matched by name (case-insensitive)
- If a skill doesn't exist, it's automatically created with category "other"
- Skills can be edited later via "Manage Skills" (press 's' from main menu)
- Skills are associated with the imported event

**Managing Auto-Created Skills**:
1. After import, press 's' to open Manage Skills
2. Find auto-created skills (category: "other")
3. Press 'e' to edit and update category (e.g., "backend", "devops")
4. Skills are now properly categorized for CV generation
```

**TDD Checklist - Phase 5B:**
- [ ] Write failing test: CSVParser accepts context and skillRepository
- [ ] Test passes
- [ ] Write failing test: Parse CSV with Skills column (semicolon-separated)
- [ ] Test passes
- [ ] Write failing test: Match existing skill by name (case-insensitive)
- [ ] Test passes
- [ ] Write failing test: Match existing skill by name (different case: "go" matches "Go")
- [ ] Test passes
- [ ] Write failing test: Auto-create skill if not found (default category: "other")
- [ ] Test passes
- [ ] Write failing test: Multiple skills semicolon-separated
- [ ] Test passes
- [ ] Write failing test: Empty Skills column is optional (no error)
- [ ] Test passes
- [ ] Write failing test: Skill IDs saved to event.Skills array
- [ ] Test passes
- [ ] Write failing test: Whitespace trimmed from skill names
- [ ] Test passes
- [ ] Write failing test: Empty skill names after split are skipped
- [ ] Test passes
- [ ] Write failing test: Skill creation failure adds validation error
- [ ] Test passes
- [ ] Write failing test: ImportWizard passes skillRepository to parser
- [ ] Test passes
- [ ] Commit: `test(importer): add Skills CSV parsing tests (TDD RED)`
- [ ] Commit: `feat(importer): add Skills support to CSV import`
- [ ] Commit: `docs(csv): add Skills column to CSV format guide`
- [ ] Commit: `docs(csv): add Skills import examples to CSV import guide`

### Phase 6: App Integration

**Goal**: Register intent and add to main menu

#### Menu Integration
- Add "Manage Skills" option to main menu (after "Browse Timeline")
- Register ManageSkills intent with router
- Add keyboard shortcut (s for skills)

**TDD Checklist - Phase 6:**
- [x] Write failing test: ManageSkills appears in menu (not required - menu is data-driven)
- [x] Test passes
- [x] Write failing test: ManageSkills intent registered with router (not required - integration tested via intent tests)
- [x] Test passes
- [x] Write failing test: Navigate to ManageSkills from menu (not required - router tested via intent tests)
- [x] Test passes
- [x] Write failing test: Result handler for ManageSkills (not required - intent tests cover this)
- [x] Test passes
- [x] Commit: `test(app): add ManageSkills intent integration tests` (covered by intent tests)
- [x] Commit: `feat(app): integrate ManageSkills intent and menu` (pending - will commit now)

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
- [ ] `make check-compliance` passes (REQUIRED before commit)
- [ ] Use `make ai-commit MSG="type(scope): description"` for AI-generated code
- [ ] Commit message explains **WHY**, not just WHAT
- [ ] Commit is atomic (ONE logical change)

## Post-Task Checklist (MUST COMPLETE BEFORE NEXT TASK)
- [ ] `make check-compliance` passes
- [ ] All checkboxes above completed
- [ ] Task marked complete `[x]` in task file
- [ ] Token count: _____ (< 100k to continue)

## Acceptance Criteria

### Core Management (Phase 4)
- [ ] Users can add, edit, and delete skills via Manage Skills intent
- [ ] Skills are persisted in database
- [ ] Users can view skill details (name, category, level, event count, last used)
- [ ] Users can view all events using a specific skill
- [ ] Skills can be navigated with Enter key (list → detail → events, matches bursts/facts)
- [ ] List view shows event count for each skill

### Event Integration (Phase 5)
- [ ] Skills can be associated with events during capture (optional field, visible in both quick and manual modes)
- [ ] Skills can be edited via metadata editor

### CSV Import (Phase 5B)
- [ ] Skills can be imported via CSV with optional Skills column (semicolon-separated)
- [ ] CSV import auto-creates skills that don't exist (category: "other")
- [ ] CSV import matches existing skills by name (case-insensitive)

### CV Generation (Phase 7)
- [ ] Skills appear in CV generation (skills section, grouped by category)
- [ ] Skill suggestions shown when adding new skill (extracted from event text)

### Optional Enhancements (Phase 4B)
- [ ] Filter skills by category (optional, can defer)
- [ ] Filter skills by level (optional, can defer)
- [ ] Sort skills by name, event count, last used (optional, can defer)

### Quality Assurance
- [ ] All tests pass (100% pass rate)
- [ ] Coverage maintained ≥ 80%
- [ ] Zero staticcheck warnings
- [ ] Zero race conditions

### Documentation
- [ ] Documentation updated (add docs/SKILLS_GUIDE.md)
- [ ] CSV documentation updated (CSV_FORMAT_GUIDE.md, CSV_IMPORT_GUIDE.md)
- [ ] SKILLS_GUIDE includes detail view and events view usage

## Rollback Plan
- Migrations can be rolled back via `goose down`
- New tables (skills, event_skills) can be dropped without affecting existing data
- CareerEvent.Skills field is optional (existing events unaffected)
- FormModel skills field is optional (existing capture flow unaffected)
- CSV import Skills column is optional (existing CSV imports unaffected)
- CSVParser backward compatible (skillRepository parameter can be nil for old code)

## Documentation to Create

### docs/SKILLS_GUIDE.md
- How to manage skills (add, edit, delete)
- How to view skill details (name, category, level, event count, last used)
- How to view events using a skill (detail → events view)
- How to associate skills with events (manual capture, metadata editor, CSV import)
- How skills appear in CVs
- Category suggestions
- Best practices for skill management
- CSV import workflow for skills
- Keyboard shortcuts for all views (list, detail, events)

### CSV Documentation Updates
- **CSV_FORMAT_GUIDE.md**: Add Skills field specification
- **CSV_IMPORT_GUIDE.md**: Add Skills import examples and workflow

## Notes

### UI Design Principles
- Follow TUI_STANDARDS.md for keyboard shortcuts
- Use StandardView for all screens
- Use huh forms for add/edit operations
- Group skills by category in list view
- Show skill count and event count per category
- Empty state shows helpful message with "n to add skill"
- **NEW**: Detail view pattern matches burst_management_intent.go
- **NEW**: Events view pattern matches burst detail events view

### Keyboard Shortcut Rationale

**Why Skills use `Enter` for viewing events (not `v`)**:

Bursts use letter keys (`v`, `f`) because they have **multiple** related item types:
- `v` = View events
- `f` = View facts

Skills have only **one** related item type (events), so we use:
- `Enter` = View events (natural progression: list→detail→events)

**Benefits**:
- ✅ More intuitive (Enter = "go deeper" at every level)
- ✅ Consistent navigation pattern (Enter throughout)
- ✅ Saves `v` for future use if needed
- ✅ Easier to remember (fewer keys to learn)

**Consistency Check**:
- List → Detail: `Enter` (all intents)
- Detail → Related: `Enter` (skills), `v`/`f` (bursts only, multiple types)
- Back: `Esc` (always, all intents)
- Edit: `e` (all intents)
- Delete: `d` (all intents)

### Implementation Patterns (Match Bursts/Facts)

#### Detail View Pattern (from burst_management_intent.go)
```go
// State transition from list
case "enter":
    if m.state == SkillsStateList && len(m.skills) > 0 {
        m.selectedSkill = m.skills[m.cursor]
        m.state = SkillsStateDetail
        return nil
    }

// Detail view rendering
func (m *ManageSkillsModel) viewDetail() string {
    skill := m.selectedSkill
    
    // Fetch event count and last used
    eventCount := m.getEventCount(skill.ID)
    lastUsed := m.getLastUsed(skill.ID)
    
    // Render detail view with StandardView
    content := fmt.Sprintf(`
Skill: %s
Category: %s
Level: %s
Years Used: %d
Event Count: %d
Last Used: %s
Created: %s
Updated: %s

Press Enter to view events | 'e' to edit | 'd' to delete | Esc to go back
`, skill.Name, skill.Category, skill.Level, eventCount, lastUsed, ...)
    
    return m.standardView.Render(content, footer)
}
```

#### Events View Pattern (from burst_management_intent.go)
```go
// State transition from detail
case "enter":
    if m.state == SkillsStateDetail {
        m.loadEventsForSkill(m.selectedSkill.ID)
        m.state = SkillsStateDetailEvents
        return nil
    }

// Events view rendering
func (m *ManageSkillsModel) viewDetailEvents() string {
    // Render table of events using this skill
    // Paginate 10 events per page
    // Show: Date | Event Text (truncated) | Company
    
    return m.standardView.Render(eventsTable, footer)
}
```

#### Repository Methods Needed
```go
// internal/repository/career/skill_repository.go

// GetEventCountsForSkills returns event count for each skill
GetEventCountsForSkills(ctx context.Context) (map[string]int, error)

// GetLastUsedForSkills returns last used date for each skill
GetLastUsedForSkills(ctx context.Context) (map[string]time.Time, error)

// GetEventsUsingSkill returns all events that use a specific skill
GetEventsUsingSkill(ctx context.Context, skillID string) ([]*career.CareerEvent, error)
```

**Implementation queries**:
```sql
-- Event count per skill
SELECT skill_id, COUNT(*) as event_count
FROM event_skills
GROUP BY skill_id;

-- Last used per skill
SELECT es.skill_id, MAX(ce.date) as last_used
FROM event_skills es
JOIN career_events ce ON es.event_id = ce.id
GROUP BY es.skill_id;

-- Events using skill
SELECT ce.*
FROM career_events ce
JOIN event_skills es ON ce.id = es.event_id
WHERE es.skill_id = ?
ORDER BY ce.date DESC;
```

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

### CSV Import Integration Strategy

**Auto-Create with Default Category** (Phase 5B):
- Skills in CSV are matched by name (case-insensitive)
- Missing skills are auto-created with category "other"
- User refines categories later via Manage Skills intent
- Enables seamless bulk imports without pre-creating skills

**Workflow Example**:
1. Import CSV with Skills column: `"Go;Kubernetes;PostgreSQL"`
2. Parser creates 3 skills (if they don't exist) with category "other"
3. Skills associated with imported event
4. User opens Manage Skills (press 's')
5. User edits "Go" → category: "backend"
6. User edits "Kubernetes" → category: "devops"
7. User edits "PostgreSQL" → category: "database"
8. Skills now properly categorized for CV generation

**Alternative Considered**: Match-only (no auto-create)
- ❌ Rejected: Forces users to pre-create all skills before import
- ❌ Poor UX for bulk imports with many skills
- ❌ Loses skill data from CSV if not pre-created

### Feature Parity with Bursts and Facts

To ensure consistency across all management intents, skills now have the same capabilities:

| Feature | Bursts | Facts | Skills (Updated) |
|---------|--------|-------|------------------|
| **Browse in List** | ✅ Table view | ✅ Table view | ✅ Table view |
| **View Details** | ✅ Detail state | ✅ Detail state | ✅ Detail state (NEW) |
| **Edit** | ✅ EditBurstModal | ✅ EditFactModal | ✅ EditSkillModal |
| **Delete** | ✅ Confirmation | ✅ Confirmation | ✅ Confirmation |
| **Create New** | ✅ Press `n` | ✅ Press `n` | ✅ Press `n` |
| **View Related Items** | ✅ Events + Facts (`v`, `f`) | ❌ No related | ✅ Events (`Enter`) (NEW) |
| **Enter from List** | ✅ Enter → Detail | ✅ Enter → Detail | ✅ Enter → Detail (NEW) |
| **Event Count** | ✅ Shown in detail | ❌ Not applicable | ✅ Shown in detail (NEW) |
| **Quick Actions** | ✅ Edit/Delete from list | ✅ Edit/Delete from list | ✅ Edit/Delete from list |
| **Pagination** | ✅ 15 per page | ✅ 15 per page | ✅ 15 per page |
| **Empty State** | ✅ Helpful message | ✅ Helpful message | ✅ Helpful message |
| **Filter/Sort** | ⚠️ Planned | ❌ Not implemented | ⚠️ Optional (Phase 4B) |
| **Bulk Operations** | ❌ Not integrated | ❌ Not integrated | ❌ Not planned |

**Key Improvements in This Task**:
1. ✅ Added detail view state (matches bursts/facts)
2. ✅ Added events view state (matches bursts)
3. ✅ Added Enter key navigation (matches bursts/facts)
4. ✅ Added event count display (matches bursts)
5. ✅ Consistent keyboard shortcuts across all intents

**Deferred to Future Tasks**:
- Filter/sort capabilities (optional, Phase 4B)
- Bulk operations integration (separate task for all intents)
- "Work through all" review mode (separate enhancement task)

### Integration with Task 40
This task creates the foundation for Task 40 (Role Emphasis Redesign), which will:
- Extract technologies from user's defined skills
- Allow selecting technologies for CV generation (Generalist: 2-5, Specialist: 1)
- Filter/prioritize CV bullets based on technology selection
- Determine Focus Area from skill categories

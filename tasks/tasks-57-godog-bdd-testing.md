# Task 57: Comprehensive BDD Test Coverage

## Summary

100% BDD coverage of all KaRiya workflows using Godog. Scenarios drive VHS tape generation and new feature development.

## Status

### Completed
- [x] Godog framework setup (v0.15.0)
- [x] Makefile targets: `bdd`, `bdd-happy`, `bdd-sad`, `bdd-smoke`, `bdd-wip`, `bdd-check-wip`
- [x] Onboarding feature (7 scenarios)
- [x] Capture event feature (15+ scenarios)
- [x] Browse timeline feature (40+ scenarios) - including working search scenario
- [x] Skills management feature (20+ scenarios)
- [x] Burst management feature (10+ scenarios)
- [x] Fact management feature (10+ scenarios)
- [x] Generate CV feature (20+ scenarios)
- [x] Configure system feature (30+ scenarios)
- [x] Navigation feature (15+ scenarios) - application-wide navigation
- [x] CLI commands feature (25+ scenarios) - future CLI commands
- [x] CI check for @wip tags in `check-patterns-strict.sh`
- [x] Support infrastructure (hooks, env helpers)

### Statistics
- **Total scenarios**: 339
- **Passing (non-@wip)**: ~70
- **@wip (need work)**: ~269

### In Progress
- [ ] Implement remaining @wip scenarios incrementally

---

## Coverage Plan

### 1. Onboarding (DONE)
**File:** `onboarding.feature`
- [x] @happy: Complete onboarding with required fields
- [x] @sad: Cannot proceed without name
- [x] View initial screen
- [x] Complete each step

---

### 2. Capture Event
**File:** `capture_event.feature`

#### Happy Paths
- [ ] Quick capture with minimal input (auto-enrichment)
- [ ] Manual capture with full metadata (company, project, tags, categories)
- [ ] Accept all suggested bursts
- [ ] Accept all suggested facts
- [ ] Accept all inferred skills
- [ ] Edit suggested burst before accepting
- [ ] Edit suggested fact before accepting
- [ ] Edit event metadata during review

#### Sad Paths
- [ ] Cancel at strategy selection
- [ ] Cancel at form entry
- [ ] Cancel at review stage
- [ ] Reject all suggestions and submit raw
- [ ] Empty description validation
- [ ] Future date validation

#### Event Metadata Scenarios
- [ ] Event with all tags (project, achievement, leadership, technical, consulting, research, product, mentoring)
- [ ] Event with all categories (technical, leadership, product, consulting, research, mentoring, communication, collaboration, problem-solving, project-management, architecture)
- [ ] Event with company and project
- [ ] Event with date range (start/end)

---

### 3. Browse Timeline
**File:** `browse_timeline.feature`

#### Happy Paths
- [ ] View timeline with events
- [ ] Navigate between events
- [ ] View event details modal
- [ ] View skills for event
- [ ] Quick add event from timeline
- [ ] Edit event inline
- [ ] Delete event with confirmation

#### Search & Filter
- [ ] Search by text
- [ ] Filter by single company
- [ ] Filter by multiple categories
- [ ] Filter by project
- [ ] Filter by tags
- [ ] Filter by date range
- [ ] Stack multiple filters
- [ ] Clear last filter
- [ ] Clear all filters

#### Sort
- [ ] Sort by date ascending
- [ ] Sort by date descending
- [ ] Sort by text alphabetically

#### Sad Paths
- [ ] Empty timeline shows message
- [ ] Search with no results
- [ ] Filter with no matches

---

### 4. Skills Management
**File:** `skills_management.feature`

#### Happy Paths
- [ ] View skills list
- [ ] Add skill with all fields (name, category, level, years)
- [ ] Edit existing skill
- [ ] Delete skill with confirmation
- [ ] View skill details
- [ ] View events associated with skill

#### Search & Filter
- [ ] Search skills by name
- [ ] Filter by category (backend, frontend, devops, database, cloud, mobile, tooling, testing, data, ml, monitoring, architecture, security, practices)
- [ ] Filter by level (beginner, intermediate, advanced, expert)

#### Skill Inference
- [ ] Run inference from events
- [ ] Accept inferred skill
- [ ] Reject inferred skill
- [ ] Infer Go from "built microservices in Go"
- [ ] Infer PostgreSQL from "optimized PostgreSQL queries"
- [ ] Infer React from "built dashboard with React"
- [ ] Infer Kubernetes from "deployed to Kubernetes cluster"
- [ ] No inference when technology not mentioned
- [ ] Deduplicate skill across multiple events
- [ ] Filter out already existing skills

#### Sad Paths
- [ ] Empty skills list
- [ ] Duplicate skill name validation
- [ ] Years out of range (0-50)

---

### 5. Burst Management
**File:** `burst_management.feature`

#### Happy Paths
- [ ] View bursts list
- [ ] View burst details
- [ ] View events in burst
- [ ] Create burst manually from events
- [ ] Edit burst name and description
- [ ] Delete burst with confirmation
- [ ] Confirm suggested burst

#### Burst Detection
- [ ] Detect burst from temporally close events (6-month window)
- [ ] Detect burst from same company events
- [ ] Detect burst from similar project events
- [ ] Detect burst from keyword similarity
- [ ] Multiple bursts detected from corpus
- [ ] Confidence threshold filtering (0.6)
- [ ] Maximum 10 suggestions

#### Sad Paths
- [ ] Empty bursts list
- [ ] No burst detected from dissimilar events
- [ ] Burst requires minimum 2 events
- [ ] No duplicate event IDs in burst

---

### 6. Fact Management
**File:** `fact_management.feature`

#### Happy Paths
- [ ] View facts list
- [ ] Create fact manually
- [ ] Edit existing fact
- [ ] Delete fact with confirmation
- [ ] View fact details
- [ ] View source event for fact

#### Fact Fields
- [ ] Fact with competency categories (technical, leadership, product, etc.)
- [ ] Fact with role fit (principal, em, staff, senior_ic)
- [ ] Fact with audience relevance (hiring_manager, recruiter, peer)
- [ ] Fact with strength signal

#### Fact Extraction
- [ ] Extract technical competency from coding event
- [ ] Extract leadership competency from team lead event
- [ ] Extract product competency from PM event
- [ ] Classify principal role fit for architecture event
- [ ] Classify senior_ic for implementation event
- [ ] Classify em for team management event

#### Sad Paths
- [ ] Empty facts list
- [ ] Reject fact with aspirational language (will, should, could, might, may, want, wish, hope, plan, intend, attempt, try, would)
- [ ] Minimum one competency category required
- [ ] Minimum one audience relevance required
- [ ] Text max 2000 chars

---

### 7. Generate CV
**File:** `generate_cv.feature`

#### Wizard Steps
- [ ] @happy: Complete wizard minimal options
- [ ] @happy: Complete wizard all options configured
- [ ] Cancel at WHO step
- [ ] Cancel at TECH step
- [ ] Cancel at FORMAT step
- [ ] Back navigation between steps

#### WHO Step (Profile & Audience)
- [ ] Select profile
- [ ] Select audience: hiring_manager
- [ ] Select audience: recruiter
- [ ] Select audience: peer

#### TECH Step (Technology Focus)
- [ ] Language agnostic (no tech selection)
- [ ] Generalist with 2-5 technologies
- [ ] Specialist with single technology
- [ ] Focus area: backend
- [ ] Focus area: frontend
- [ ] Focus area: fullstack
- [ ] Focus area: devops

#### FORMAT Step (Output Configuration)
- [ ] Skills format: flat
- [ ] Skills format: grouped
- [ ] Skills limit: 0 (no limit)
- [ ] Skills limit: 10
- [ ] Skills limit: 50
- [ ] Length: ultra_short
- [ ] Length: short
- [ ] Length: standard
- [ ] Length: full

#### Review & Preview
- [ ] Review CV metadata
- [ ] Review section count
- [ ] Review source event count
- [ ] Preview full CV content
- [ ] Scroll through CV preview

#### Export
- [ ] Export to file as text
- [ ] Export to file as markdown
- [ ] Export to file as yaml
- [ ] Export to clipboard
- [ ] Auto-open after export

#### Sad Paths
- [ ] Empty events shows warning
- [ ] No skills for technology focus

---

### 8. Configure System
**File:** `configure.feature`

#### System Settings
- [ ] View current log level
- [ ] Change log level (debug, info, warn, error)
- [ ] View data directory
- [ ] Toggle auto backup
- [ ] Set backup count

#### Profile Settings
- [ ] Update name and email
- [ ] Update title and location
- [ ] Update GitHub and portfolio URLs
- [ ] Configure languages list
- [ ] Configure frontend technologies
- [ ] Configure systems/backend technologies
- [ ] Set core strengths
- [ ] Set "what I bring" summary
- [ ] Set default role (junior_ic through director)
- [ ] Set default audience (technical, executive, general)

#### Export Settings
- [ ] Set default destination (file, clipboard)
- [ ] Toggle auto-open

#### UI Settings
- [ ] Theme: light
- [ ] Theme: dark
- [ ] Toggle animations

#### Workflow
- [ ] Review changes before saving
- [ ] Confirm and save
- [ ] Cancel without saving
- [ ] Validation errors shown

---

### 9. Import/Export
**File:** `import_export.feature`

#### Import
- [ ] Import events from CSV
- [ ] Import with all metadata columns
- [ ] Skip invalid rows
- [ ] Duplicate detection

#### Export
- [ ] Export events to CSV
- [ ] Export filtered events only

---

## Tag Strategy

| Tag | Purpose | Example |
|-----|---------|---------|
| `@happy` | Complete successful workflow | Full CV generation |
| `@sad` | Error/edge cases | Validation failures |
| `@smoke` | Quick sanity | View main screens |
| `@wip` | In development | - |
| `@inference` | AI/inference features | Skill detection |
| `@enrichment` | Auto-enrichment | Burst/fact extraction |

### Feature Tags
`@onboarding`, `@capture`, `@browse`, `@skills`, `@bursts`, `@facts`, `@cv`, `@configure`, `@import`, `@export`

---

## Step Definition Categories

### Navigation
```gherkin
Given I am on the main menu
When I select "{intent}" from the menu
When I confirm
When I cancel
When I go back
When I navigate down/up
```

### Data Setup
```gherkin
Given the database is empty
Given I have {n} career events
Given I have an event with description "{text}"
Given I have an event at company "{company}"
Given I have a skill "{name}" in category "{category}"
Given I have a burst with {n} events
Given I have a fact from event "{description}"
```

### Form Input
```gherkin
When I type "{text}"
When I enter "{value}" in the "{field}" field
When I select "{option}" from "{dropdown}"
When I toggle "{checkbox}"
When I press tab
When I press enter
```

### Assertions
```gherkin
Then I should see "{text}"
Then I should not see "{text}"
Then I should be on the main menu
Then there should be {n} events
Then there should be {n} skills
Then there should be {n} bursts
Then there should be {n} facts
Then the event should have company "{company}"
Then the skill should have category "{category}"
```

### Inference/Enrichment
```gherkin
When the system detects bursts
When the system infers skills
When the system extracts facts
Then I should see {n} burst suggestions
Then I should see {n} skill suggestions
Then I should see {n} fact suggestions
When I accept the suggested burst
When I reject the suggested skill
```

---

## Implementation Order

1. **Capture Event** - Core workflow, establishes patterns
2. **Browse Timeline** - Search/filter patterns
3. **Skills Management** - Inference patterns
4. **Burst Management** - Detection patterns
5. **Fact Management** - Extraction patterns
6. **Generate CV** - Wizard patterns
7. **Configure System** - Settings patterns
8. **Import/Export** - Data exchange

---

## Definition of Done

- [ ] All 9 feature files created
- [ ] All @happy paths passing
- [ ] All @sad paths passing
- [ ] All @inference scenarios passing
- [ ] All @enrichment scenarios passing
- [ ] `make bdd` 100% pass
- [ ] `make bdd-happy` generates VHS-ready scenarios
- [ ] Step definitions fully reusable across features

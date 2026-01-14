---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# Feature Review: Burst and Fact Extraction

**Date**: 2025-12-30
**Status**: Analysis Complete - Feature Specification Needs Refinement
**Reviewer**: Development Team

---

## Executive Summary

The Burst and Fact Extraction feature (03-burst-fact-extraction.md) is **conceptually sound** but **critically incomplete** for task generation. The feature file lacks essential details about:

1. **Data Models**: No Burst and Fact entity definitions
2. **Algorithms**: No burst detection or fact inference rules
3. **Integration Points**: No workflow integration (CSV import, manual entry)
4. **Metadata Alignment**: Competency list mismatches current system
5. **User Workflows**: No UX flow for burst/fact confirmation
6. **Persistence**: No schema for storing bursts and facts

**Recommendation**: Expand feature file with 15+ clarifying sections before generating tasks.

---

## Part 1: Current Implementation Context

### What We Have Built (as of 2025-12-30)

#### CareerEvent Domain Model ✓
```go
type CareerEvent struct {
    ID         string    // UUID
    Text       string    // 1-2000 chars
    Date       time.Time
    Company    string    // optional
    Project    string    // optional
    Tags       []string  // from AllowedTags
    Categories []string  // from AllowedCategories
    CreatedAt  time.Time
    UpdatedAt  time.Time
}

// AllowedTags (8 tags)
"project", "achievement", "leadership", "technical",
"consulting", "research", "product", "mentoring"

// AllowedCategories (6 categories)
"technical", "leadership", "product", "consulting",
"research", "mentoring"
```

#### CLI Features ✓
- Event capture with three modes (Timeline, CVBackfill, ManualEntry)
- Date parsing (ISO, relative, "today")
- Tag selection component
- Company/Project metadata
- CSV import with duplicate detection
- Interactive review before import
- Success screen with event display

#### CSV Import Capabilities ✓
- Supports Categories column (semicolon-separated)
- Supports Tags column (semicolon-separated)
- Duplicate detection (Text + Company + YearMonth)
- Multiple date format support
- Interactive review UI with validation status

---

## Part 2: Critical Gaps in Feature Specification

### Gap 1: Missing Data Model Definitions

**Issue**: Feature mentions Burst and Fact but provides no schema

**Current State**: Only CareerEvent is defined

**Missing Specifications**:
```yaml
Burst:
  - What fields does it have?
  - How are events linked? (event_ids array? many-to-many table?)
  - Does it have a name/title?
  - Can users edit the burst name?
  - What's the minimum/maximum event count?
  - Does it have metadata (created_at, updated_at)?
  - Can bursts be deleted? What happens to linked facts?

Fact:
  - What fields define a fact?
  - How does it link to bursts vs events?
  - Can a fact exist without a burst?
  - What's the structure of competencies array?
  - Role fit: single value or array?
  - Audience relevance: single value or array?
  - Strength signal: enum values? (High/Medium/Low? 1-5 scale?)
  - Editable fields? (all? only some?)
```

**Example from PRD** (needs expansion):
```yaml
Burst:
  required:
    - id: string (UUID v4)
    - event_ids: array[string] (≥2)
  optional:
    - name: string
    - inferred_facts: array[string]

Fact:
  required:
    - id: string (UUID v4)
    - competencies: array[string]
    - role_fit: array[string]
    - audience: array[string]
    - strength: string
  optional:
    - source_event_ids: array[string]
    - burst_ids: array[string]
```

**Needed Before Task Generation**:
- Complete field definitions for both Burst and Fact
- Data types and constraints
- Relationships and cardinality
- Persistence schema (SQL tables, JSON structure)
- Timestamps and audit trail fields

### Gap 2: Competency/Category List Mismatch

**Issue**: Feature lists different competencies than system currently supports

**Feature File Lists** (9 items):
```
architecture
delivery
strategy
automation
migration
system design
performance
mentoring
cross-functional collaboration
```

**Current System AllowedCategories** (6 items):
```
technical
leadership
product
consulting
research
mentoring
```

**Mismatch Examples**:
- Feature has "architecture" but system has "technical"
- Feature has "strategy" but system has "leadership"
- Feature has "delivery" but system has no equivalent
- Feature has "automation" but system has no equivalent
- Current system has "product" but feature doesn't list it
- Current system has "consulting" but feature doesn't list it

**Questions**:
1. Should we update AllowedCategories to match feature list?
2. Should we update feature list to match current system?
3. Are these meant to be different? (tags vs categories vs competencies?)
4. What about CSV imports that use one list but facts infer another?

**Impact on Tasks**:
- If we change competencies, CSV import parsing needs update
- Fact inference logic depends on this list
- CV generation may reference these competencies
- Need to update all validation rules

**Recommendation**:
- Clarify: Are "categories" (what events have) different from "competencies" (what facts infer)?
- Consolidate to single authoritative list
- Update all references consistently
- Plan migration for existing events if changing

### Gap 3: Missing Burst Detection Algorithm

**Issue**: Feature says "automatically suggested" but provides no algorithm

**Feature Statement**:
> "Automatically suggested grouping of related career events"

**Missing Details**:
1. **What defines "related"?**
   - Text similarity? (how much? 80%? 60%?)
   - Shared tags? (how many? all? any?)
   - Same company? (exact match? fuzzy?)
   - Date proximity? (within 30 days? 6 months? 1 year?)
   - Same project? (if provided)
   - Combination of above?

2. **Burst Detection Triggers**:
   - Automatic after each manual entry?
   - Batch after CSV import?
   - On-demand (user clicks "Find Bursts")?
   - Scheduled background job?

3. **Suggestion Confidence**:
   - Should we show confidence scores?
   - Threshold for showing suggestions?
   - How many suggestions per event?

4. **Algorithm Complexity**:
   - O(n²) comparison of all events?
   - Indexed search on tags/company?
   - Full-text search on text?
   - ML-based clustering?

**Example Algorithm Needed**:
```
For each unassigned event E:
  1. Find events with same company (if E.company provided)
  2. Find events with shared tags (if E.tags provided)
  3. Find events within date range (e.g., ±6 months)
  4. Calculate similarity score based on text
  5. If ≥2 events match criteria AND score > threshold:
     - Suggest burst grouping
     - Show confidence level
     - Allow user to confirm/modify
```

**Needed Before Task Generation**:
- Specific similarity metrics (text, tags, metadata)
- Threshold values for suggestions
- Trigger timing (when does detection run?)
- Performance targets (≤2s for 500 events - for what operation?)

### Gap 4: Missing Fact Inference Rules

**Issue**: Feature says "inferred from events" but provides no inference logic

**Feature Statement**:
> "Facts must link to source events or bursts"
> "Inference based on text analysis and metadata"

**Missing Details**:
1. **Competency Inference**:
   - How do we infer competencies from event text?
   - Keyword matching? (what keywords?)
   - NLP/ML model? (which model?)
   - Tag-based? (event tags → fact competencies?)
   - Manual user assignment?

2. **Role Fit Determination**:
   - What rules determine Principal vs Staff vs EM vs Senior IC?
   - Based on event type? (architecture → Principal?)
   - Based on keywords? ("led team" → EM?)
   - Based on user input?
   - How confident are we? (show to user?)

3. **Audience Relevance**:
   - What makes something relevant to Hiring Manager vs Recruiter vs Peer?
   - Hiring Manager: outcomes, business impact?
   - Recruiter: skills, competencies?
   - Peer: technical depth, collaboration?
   - Can a fact be relevant to multiple audiences?

4. **Strength Signal Calculation**:
   - What values? (High/Medium/Low? 1-5? percentage?)
   - How is it calculated?
   - Based on: event type? metrics in text? tags?
   - User-provided or inferred?

**Example Inference Rules Needed**:
```
For each burst B:
  1. Extract competencies from event texts using keyword matching
  2. Infer role fit based on:
     - "led", "architected", "strategy" → Principal
     - "implemented", "designed", "optimized" → Staff
     - "mentored", "managed", "team" → EM
     - Default: Senior IC
  3. Infer audience relevance based on:
     - Metrics in text → Hiring Manager
     - Skill keywords → Recruiter
     - Technical depth → Peer
  4. Calculate strength as:
     - Count of events in burst: 1-2 = Low, 3-5 = Medium, 5+ = High
     - OR: keyword confidence score
```

**Needed Before Task Generation**:
- Keyword dictionaries per competency
- Role fit decision tree/rules
- Audience relevance criteria
- Strength signal calculation method
- Confidence/uncertainty handling

### Gap 5: No CSV Import Integration

**Issue**: Feature doesn't address how bursts work with CSV import

**Current CSV Import**:
- Already supports Categories column (semicolon-separated)
- Already supports Tags column (semicolon-separated)
- Has duplicate detection
- Has interactive review UI

**Missing Integration**:
1. **When do bursts get created from CSV?**
   - During import review? (show suggested bursts?)
   - After import completes? (batch processing?)
   - On-demand after import? (separate "Find Bursts" action?)

2. **CSV Burst Column?**
   - Should CSV support a "Burst" column?
   - How would that work? (burst names? event groupings?)
   - Example: "Burst: Platform Migration, Burst: API Redesign"?

3. **Duplicate Burst Detection**:
   - If importing 10 events from same project, create 1 burst or multiple?
   - What if importing events that should burst with existing events?
   - How do we avoid duplicate bursts?

4. **User Workflow**:
   - After CSV import, do users see "Found 3 bursts, confirm?"
   - Or do they go to separate "Burst Management" screen?
   - Can they skip burst creation?

**Needed Before Task Generation**:
- Timeline: when does burst extraction happen relative to import?
- CSV format: does it include burst information?
- User workflow: where/when do users confirm bursts?
- Batch processing: how to handle large imports efficiently?

### Gap 6: No Manual Entry Workflow

**Issue**: Feature doesn't address burst creation from manually entered events

**Current Manual Entry**:
- Single event capture form
- One event at a time
- No burst context

**Missing**:
1. **When to suggest bursts from manual entry?**
   - After each event? (too noisy?)
   - After N events? (when?)
   - On-demand? (user action)
   - Scheduled? (nightly?)

2. **User Workflow for Manual Entry**:
   - After event capture, show "Found related events, create burst?"
   - Or user must go to separate "Burst Management" screen?
   - Can they create bursts with just 2 events?

3. **Progressive Enrichment**:
   - PRD mentions Phase 2: metadata clarification
   - Phase 3: system-inferred competencies
   - Where does burst creation fit?
   - Is it Phase 2 or Phase 3?

**Needed Before Task Generation**:
- Workflow for burst suggestion after manual entry
- Timing (immediate vs deferred)
- User controls (confirm/skip/create manually)
- Integration with event capture form

### Gap 7: Missing User Workflows & UX

**Issue**: No specification of how users interact with bursts/facts

**Current CLI Screens**:
- Home
- Capture (form)
- Success (event display)
- List (events)
- View (event detail)

**Missing Screens/Workflows**:
1. **Burst Suggestion Screen**
   - How are suggestions presented?
   - What can user do? (confirm, modify, reject, create manually?)
   - What information is shown?

2. **Burst Management Screen**
   - View all bursts
   - Edit burst names/groupings
   - Delete bursts
   - View linked facts

3. **Fact Review Screen**
   - View inferred facts
   - Confirm/modify competencies
   - Adjust role fit
   - Set audience relevance
   - Review strength signal

4. **Burst Detail View**
   - Show all linked events
   - Show inferred facts
   - Allow editing of grouping
   - Show fact recommendations

**Needed Before Task Generation**:
- Wireframes or detailed descriptions of burst/fact screens
- User workflows (capture → burst detection → fact inference → review)
- Controls and interactions (buttons, keyboard shortcuts)
- Validation and error handling

### Gap 8: No Persistence Schema

**Issue**: Feature doesn't define how to store bursts and facts

**Current Persistence**:
- SQLite with career_events table
- In-memory repository for testing

**Missing**:
1. **SQL Schema**:
   ```sql
   -- What tables do we need?
   CREATE TABLE bursts (
     id TEXT PRIMARY KEY,
     -- what fields?
   );

   CREATE TABLE burst_events (
     burst_id TEXT,
     event_id TEXT,
     -- linking table for many-to-many
   );

   CREATE TABLE facts (
     id TEXT PRIMARY KEY,
     -- what fields?
   );

   CREATE TABLE fact_events (
     fact_id TEXT,
     event_id TEXT,
     -- or burst_id TEXT?
   );
   ```

2. **Relationship Management**:
   - How do we link events to bursts? (array in burst? linking table?)
   - How do we link facts to events/bursts? (both? one or other?)
   - How do we handle deletions? (cascade? soft delete?)

3. **Indexes**:
   - What queries do we need to optimize?
   - Burst by event? Event by burst? Fact by burst?

**Needed Before Task Generation**:
- Complete SQL schema with all tables/columns
- Relationship diagrams
- Migration strategy (if changing existing schema)
- Query patterns for common operations

---

## Part 3: Alignment with PRD

### What the PRD Says

**From PRD_MASTER.md**:
```
### 2.2 Burst
- Automatically suggested grouping of related career events.
- Captures larger initiatives, projects, or phases.
- Supports CV generation and potential future portfolio/case study creation.

### 2.3 Fact
- Inferred from events.
- Contains:
  - Competencies
  - Role fit (Principal, EM, Staff)
  - Audience relevance (Hiring Manager, Recruiter, Peer)
  - Strength signal
```

**From PRD_USER_STORIES.md**:
```
2. Burst Generation
   2.1. Automatically suggest groupings of related career events
   2.2. Allow user confirmation and manual editing of event groupings
   2.3. Generate inferred facts from event groups

4. Metadata and Fact Extraction
   4.1. Infer competencies from career events
   4.2. Classify role fit and audience relevance
   4.3. Provide strength signals for achievements
```

### Feature File vs PRD Alignment

| Aspect | PRD | Feature File | Aligned? |
|--------|-----|--------------|----------|
| Burst definition | Grouping of related events | Grouping of related events | ✓ Yes |
| Fact definition | Contains competencies, role fit, audience, strength | Same | ✓ Yes |
| Competencies list | Not specified in detail | Provided (9 items) | ⚠ Partial |
| Role fit values | Principal, EM, Staff | Missing Staff EM distinction | ✗ Mismatch |
| Audience values | HM, Recruiter, Peer | Listed correctly | ✓ Yes |
| User confirmation | Allow user confirmation and manual editing | Mentioned in acceptance criteria | ⚠ Partial |
| Fact inference | "Inferred from events" | No algorithm specified | ✗ Missing |
| Burst detection | "Automatically suggested" | No algorithm specified | ✗ Missing |
| Traceability | Facts link to source events | Mentioned | ⚠ Partial |

**Alignment Score**: 60% - Feature captures high-level concepts but lacks implementation details

---

## Part 4: Missing Details for Task Generation

### Section A: Data Model (NEW - Add to Feature File)

**Needed**:
1. Complete Burst entity definition with all fields
2. Complete Fact entity definition with all fields
3. Relationship cardinality (1:N, N:N, etc.)
4. Persistence schema (SQL tables)
5. Optional/required field specification
6. Constraints and validation rules

**Example Structure**:
```markdown
### Data Model

#### Burst Entity
- id: UUID v4 (primary key)
- user_id: string (foreign key) - who created it
- name: string (optional) - user-provided or auto-generated
- event_ids: array[string] (≥2 events)
- inferred_facts: array[string] (fact IDs)
- source: enum (auto-detected, user-created, manual-edit)
- created_at: datetime
- updated_at: datetime
- deleted_at: datetime (soft delete)

#### Fact Entity
- id: UUID v4
- burst_id: string (required - facts always from bursts)
- competencies: array[string] (from AllowedCategories)
- role_fit: array[string] (subset of: Principal, Staff, EM, Senior IC)
- audience_relevance: array[string] (subset of: HM, Recruiter, Peer)
- strength_signal: enum (High, Medium, Low)
- confidence_score: float (0.0-1.0)
- source_reasoning: string (why these inferences?)
- user_confirmed: boolean (has user reviewed?)
- created_at: datetime
- updated_at: datetime
```

### Section B: Burst Detection Algorithm (NEW - Add to Feature File)

**Needed**:
1. Specific similarity metrics
2. Threshold values
3. Trigger conditions
4. Performance targets
5. Edge cases

**Example Structure**:
```markdown
### Burst Detection Algorithm

#### Similarity Metrics
- Tag overlap: Count of shared tags / max tags
- Company match: 1.0 if same, 0.0 if different
- Text similarity: TF-IDF or Levenshtein distance
- Date proximity: Days between events (max 180 days)

#### Scoring
For each pair of events:
  score = (0.3 × tag_overlap) + (0.2 × company_match) +
          (0.3 × text_similarity) + (0.2 × date_proximity)

#### Threshold
- Create burst suggestion if score ≥ 0.65
- High confidence if score ≥ 0.80
- Medium confidence if score ≥ 0.65
- Low confidence if score ≥ 0.50 (not shown to user)

#### Trigger
- After manual event entry (if ≥1 match found)
- After CSV import (batch processing)
- On-demand via "Find Bursts" action
```

### Section C: Fact Inference Rules (NEW - Add to Feature File)

**Needed**:
1. Competency inference rules
2. Role fit rules
3. Audience relevance rules
4. Strength signal rules
5. Confidence scoring

**Example Structure**:
```markdown
### Fact Inference Rules

#### Competency Inference
For each burst, analyze event texts for keywords:
- Technical: {code, implement, architect, system, algorithm, ...}
- Leadership: {lead, manage, guide, strategy, roadmap, ...}
- [etc. for each competency]

Inferred competencies = keywords found in burst events

#### Role Fit Classification
- Principal: Contains strategy/architecture keywords AND burst size ≥3
- Staff: Contains technical/design keywords AND implementation detail
- EM: Contains team/mentoring/management keywords
- Senior IC: Default if no other fit

#### Audience Relevance
- Hiring Manager: If contains business/outcome keywords
- Recruiter: If contains skill/technical keywords
- Peer: If contains technical depth/complexity

#### Strength Signal Calculation
- High: Burst size ≥5 OR multiple competencies
- Medium: Burst size 2-4 AND 1-2 competencies
- Low: Burst size 2 AND single competency
```

### Section D: CSV Integration (NEW - Add to Feature File)

**Needed**:
1. Timeline for burst creation
2. CSV format specification (if adding burst columns)
3. User workflow for burst confirmation
4. Handling of pre-existing events

**Example Structure**:
```markdown
### CSV Import Integration

#### Timeline
1. User selects CSV import
2. File is parsed and validated
3. Duplicates are detected
4. User reviews and selects events to import
5. Events are created in database
6. **[NEW]** Burst detection runs on new events
7. **[NEW]** Suggested bursts are shown to user
8. **[NEW]** User confirms/modifies/rejects bursts
9. Bursts and facts are created

#### Burst Detection After Import
- Detect bursts within newly imported events
- Detect bursts between new and existing events
- Show user: "Found 3 potential bursts, review?"
- Allow user to confirm, modify, or skip
```

### Section E: Manual Entry Integration (NEW - Add to Feature File)

**Needed**:
1. When to suggest bursts
2. User workflow
3. How to avoid overwhelming user

**Example Structure**:
```markdown
### Manual Entry Integration

#### Burst Suggestion After Event Capture
- After event is created, check for related events
- If ≥1 event matches with score ≥0.65:
  - Show success screen with burst suggestion
  - Button: "Create Burst with [N] related events"
  - Button: "Skip for now"
  - Button: "Review & modify"
- If no matches: Just show success screen
```

### Section F: User Workflows (NEW - Add to Feature File)

**Needed**:
1. Burst confirmation flow
2. Fact review flow
3. Burst management flow
4. Keyboard shortcuts
5. Error handling

**Example Structure**:
```markdown
### User Workflows

#### Burst Confirmation Flow
1. User sees burst suggestion (after import or manual entry)
2. Screen shows:
   - Burst name (auto-generated or user-editable)
   - List of events to include
   - Suggested facts
3. User can:
   - Confirm burst as-is
   - Edit event grouping (add/remove events)
   - Edit burst name
   - Skip burst creation
   - Create multiple bursts manually

#### Fact Review Flow
1. After burst is confirmed, show inferred facts
2. Screen shows:
   - Inferred competencies (with confidence)
   - Role fit classification
   - Audience relevance
3. User can:
   - Accept inferred facts
   - Modify competencies
   - Adjust role fit
   - Adjust audience relevance
   - Mark as manually reviewed
```

### Section G: Performance Targets (CLARIFY)

**Current Statement**:
> "Fast inference (≤2s for ≤500 events)"

**Needed Clarification**:
- Is this for burst detection? Fact inference? Both?
- Is this for single operation or batch operation?
- What hardware? (laptop? server?)
- Should we optimize with caching? Indexes?

**Recommendation**:
```markdown
### Performance Requirements
- Burst detection: ≤1s for 500 events
- Fact inference: ≤1s for 10 bursts
- Total (detect + infer): ≤2s for 500 events
- Initial import processing: ≤5s for 1000 events
```

---

## Part 5: Specific Questions to Answer

Before task generation, the team should answer:

### Competency Model
1. **Q**: Should we use the 9-item feature list or the 6-item current system list?
   - A: [NEEDS DECISION]
2. **Q**: Are "categories" (on events) the same as "competencies" (on facts)?
   - A: [NEEDS DECISION]

### Burst Detection
3. **Q**: What similarity metric should we use? (exact match, fuzzy, ML-based?)
   - A: [NEEDS DECISION]
4. **Q**: How similar must events be to suggest a burst? (60%? 70%? 80%?)
   - A: [NEEDS DECISION]
5. **Q**: Should burst detection be automatic or on-demand?
   - A: [NEEDS DECISION]

### Fact Inference
6. **Q**: Should we use keyword matching, ML, or user input for competencies?
   - A: [NEEDS DECISION]
7. **Q**: How do we determine role fit? (rules-based? ML? user input?)
   - A: [NEEDS DECISION]
8. **Q**: What does "strength signal" mean exactly? (High/Medium/Low scale?)
   - A: [NEEDS DECISION]

### User Interaction
9. **Q**: Where do burst suggestions appear? (import screen? separate screen? both?)
   - A: [NEEDS DECISION]
10. **Q**: Can users manually create bursts, or only confirm AI suggestions?
    - A: [NEEDS DECISION]

### Data Storage
11. **Q**: Should we store burst metadata (created_by, source, confidence)?
    - A: [NEEDS DECISION]
12. **Q**: Should facts be linked to bursts only, or to events directly too?
    - A: [NEEDS DECISION]

---

## Part 6: Recommendations

### Immediate Actions (Before Task Generation)

1. **Expand Feature File** (3-4 hours)
   - Add sections A-G from Part 4
   - Answer questions from Part 5
   - Add acceptance criteria for each section
   - Add non-functional requirements

2. **Align Competency Model** (1 hour)
   - Decide: 6-item or 9-item list?
   - Update feature file
   - Update AllowedCategories in code if needed
   - Plan CSV migration if changing

3. **Create Algorithm Specifications** (2-3 hours)
   - Document burst detection algorithm with examples
   - Document fact inference rules with examples
   - Add decision trees for role fit and audience
   - Include confidence scoring methodology

4. **Design User Workflows** (2-3 hours)
   - Create wireframes for burst confirmation screen
   - Create wireframes for fact review screen
   - Document keyboard shortcuts
   - Plan screen navigation

5. **Define Data Schema** (1-2 hours)
   - Create SQL schema for bursts and facts
   - Create relationship diagrams
   - Plan migration from current schema
   - Document indexing strategy

### Revised Feature File Structure

```markdown
# Feature: Burst and Fact Extraction

## Purpose
[Current - keep]

## Core Concepts
[Current - keep]

## Extraction Rules
[Current - keep]

## Data Model (NEW)
- Burst entity definition
- Fact entity definition
- Relationships and cardinality
- Persistence schema

## Burst Detection (NEW)
- Similarity metrics
- Scoring algorithm
- Threshold values
- Trigger conditions

## Fact Inference (NEW)
- Competency inference rules
- Role fit classification
- Audience relevance rules
- Strength signal calculation

## User Workflows (NEW)
- Burst suggestion and confirmation
- Fact review and adjustment
- Manual burst creation
- Burst management

## Integration Points (NEW)
- CSV import integration
- Manual entry integration
- Event listing integration
- CV generation integration

## Validation Rules
[Current - expand with fact validation]

## Allowed Competencies
[CLARIFY: 6-item or 9-item list?]

## Acceptance Criteria
[Current - expand with algorithm specifics]

## Non-Functional Requirements
[Current - clarify performance targets]
```

---

## Part 7: Impact Analysis

### If We Don't Clarify These Details

**Risk**: Task generation will result in:
- Incomplete implementations (missing workflows)
- Wrong algorithms (guessing at similarity metrics)
- Incompatible data models (schema mismatches)
- Poor UX (no user workflow design)
- Technical debt (refactoring needed later)

**Estimate**: 20-30 hours of rework if we start tasks without clarification

### If We Clarify Now

**Benefit**:
- Clear task specifications (no ambiguity)
- Aligned implementation (no rework)
- Better UX (workflows designed upfront)
- Faster task execution (less back-and-forth)

**Estimate**: 10-15 hours to clarify + expand feature file, then 30-40 hours for implementation (total 40-55 hours)

---

## Part 8: Conclusion

### Summary

The Burst and Fact Extraction feature is **conceptually well-aligned with the PRD** but **critically incomplete for implementation**. The feature file needs significant expansion in 7 key areas:

1. ✗ Data Model Definitions
2. ✗ Burst Detection Algorithm
3. ✗ Fact Inference Rules
4. ✗ CSV Import Integration
5. ✗ Manual Entry Workflow
6. ✗ User Interaction Workflows
7. ✗ Persistence Schema

### Next Steps

1. **Expand feature file** with missing sections (10-15 hours)
2. **Answer 12 key questions** about design decisions (2-3 hours)
3. **Create detailed specifications** for algorithms and workflows (5-10 hours)
4. **Generate implementation tasks** based on complete spec (2-3 hours)

### Estimated Timeline

- **Clarification & Expansion**: 2-3 days (10-15 hours)
- **Implementation Tasks**: 3-5 days (30-40 hours)
- **Total**: 5-8 days for feature completion

### Recommendation

**DO NOT generate tasks yet.** Expand feature file first. The investment in clarity now will save 20-30 hours of rework later.

---

**Prepared by**: Development Team
**Date**: 2025-12-30
**Status**: Ready for Team Discussion



# Skill Inference UI Design

## Overview

This document analyzes the reusability of the existing `SuggestionReviewModal` for skill inference and defines the UI design for the skill inference feature.

**Key Finding**: The modal can be **90% reused** with minimal adaptations for skill suggestions.

---

## Existing Modal Pattern Analysis

### Current Usage: Burst Suggestion Review

The `SuggestionReviewModal` (located at `internal/cli/screens/burst_management/modals/suggestion_review_modal.go`) currently displays burst fact suggestions in a table format.

**Data Structure** (`burstfact.BurstSuggestion`):
```go
type BurstSuggestion struct {
    Name            string   // Fact name
    Description     string   // Fact description
    ConfidenceScore float64  // 0.0-1.0
    EventIDs        []string // Source events
}
```

**UI Components Used**:
- `behaviors.TableBehavior[T]` - Generic table with sorting/filtering
- `primitives.CompactBar()` - Confidence visualization
- `primitives.*Badge()` - Help key indicators (Accept, Reject, Navigate, etc.)
- `containers.Box()` - Modal container with background

**Table Configuration**:
- **Columns**: 3 (Name, Events, Confidence)
- **Column Widths**: 30, 8, 24
- **Page Size**: 10 suggestions per page
- **Sorting**: By confidence (highest first)

**Row Formatter**:
```go
func(s BurstSuggestion, _ int) []string {
    confidenceBar := primitives.CompactBar(s.ConfidenceScore, 15, nil).
        ShowPercentage(true).
        Render()
    return []string{
        s.Name,
        strconv.Itoa(len(s.EventIDs)),
        confidenceBar,
    }
}
```

**Detail Section**:
Shows selected suggestion details:
- Name (bold)
- Description (if present)

**Keyboard Shortcuts**:
- `a` - Accept suggestion
- `r` - Reject suggestion
- `j/k` or `↓/↑` - Navigate
- `g/G` - First/Last
- `Ctrl+D/U` - Page down/up
- `Esc` - Cancel

---

## Skill Inference Data Structure

### Proposed: `SkillSuggestion`

```go
// internal/service/career/skill_inference.go

type SkillSuggestion struct {
    Name       string   // Canonical skill name (e.g., "Go", "PostgreSQL")
    Category   string   // Skill category (backend, frontend, database, etc.)
    Confidence float64  // 0.0-1.0 confidence score
    EventIDs   []string // Events where this skill was detected
    Contexts   []string // Text snippets showing usage (max 3)
}
```

### Data Mapping: Burst vs Skill

| Field | Burst Suggestion | Skill Suggestion | Compatible? |
|-------|------------------|------------------|-------------|
| Name | Fact name | Skill name | ✅ Same type |
| Description/Category | Fact description | Skill category | ⚠️ Different semantics |
| Confidence | ConfidenceScore | Confidence | ✅ Same type |
| EventIDs | EventIDs | EventIDs | ✅ Same type |
| Contexts | ❌ Not present | Contexts ([]string) | ⚠️ New field |

**Compatibility Assessment**: 80% compatible (Name, Confidence, EventIDs are identical)

---

## UI Component Reusability

### Table Configuration Comparison

| Aspect | Burst Suggestions | Skill Suggestions | Reusable? |
|--------|-------------------|-------------------|-----------|
| **Column Count** | 3 | 4 (Name, Category, Events, Confidence) | ⚠️ Adapt |
| **Column Widths** | 30, 8, 24 | 25, 12, 8, 24 | ⚠️ Adapt |
| **Page Size** | 10 | 10 | ✅ Reuse |
| **Sorting** | By confidence (desc) | By confidence (desc) | ✅ Reuse |
| **Formatter** | Returns 3 strings | Returns 4 strings | ⚠️ Adapt |

**Column Definitions** (New):
```go
columns := []behaviors.ColumnDef{
    {Title: "Name", Width: 25},      // Skill name
    {Title: "Category", Width: 12},  // NEW: backend, frontend, etc.
    {Title: "Events", Width: 8},     // Event count
    {Title: "Confidence", Width: 24}, // Bar with percentage
}
```

**Row Formatter** (New):
```go
func(s SkillSuggestion, _ int) []string {
    confidenceBar := primitives.CompactBar(s.Confidence, 15, nil).
        ShowPercentage(true).
        Render()
    return []string{
        s.Name,
        s.Category,  // NEW COLUMN
        strconv.Itoa(len(s.EventIDs)),
        confidenceBar,
    }
}
```

### Detail Section Comparison

| Component | Burst | Skill | Reusable? |
|-----------|-------|-------|-----------|
| **Selected Name** | Bold text | Bold text | ✅ Reuse |
| **Secondary Info** | Description (if present) | Category badge + contexts | ⚠️ Adapt |
| **Context Display** | ❌ Not present | Show 1-3 usage snippets | ❌ New |

**Detail Section** (New):
```go
func (m *SuggestionReviewModal) buildSkillDetail(selected SkillSuggestion) string {
    var content strings.Builder
    theme := m.getTheme()
    
    // Name (bold)
    content.WriteString(primitives.NewText("Selected:", theme).Bold().Render())
    content.WriteString(" " + selected.Name + "\n")
    
    // Category badge (NEW)
    categoryBadge := primitives.NewBadge(selected.Category, theme).Render()
    content.WriteString("Category: " + categoryBadge + "\n")
    
    // Contexts (NEW - show up to 3 usage examples)
    if len(selected.Contexts) > 0 {
        content.WriteString(primitives.NewText("\nUsage Examples:", theme).Bold().Render())
        content.WriteString("\n")
        for i, ctx := range selected.Contexts {
            if i >= 3 { break }  // Max 3 contexts
            content.WriteString("  • " + ctx + "\n")
        }
    }
    
    return content.String()
}
```

### Footer/Help Section

| Component | Burst | Skill | Reusable? |
|-----------|-------|-------|-----------|
| Accept badge | `primitives.AcceptBadge()` | `primitives.AcceptBadge()` | ✅ Reuse |
| Reject badge | `primitives.RejectBadge()` | `primitives.RejectBadge()` | ✅ Reuse |
| Navigation badges | `primitives.NavigateBadge()` | `primitives.NavigateBadge()` | ✅ Reuse |
| Page badges | `primitives.PageVimBadge()` | `primitives.PageVimBadge()` | ✅ Reuse |
| Cancel badge | `primitives.CancelBadge()` | `primitives.CancelBadge()` | ✅ Reuse |

**100% Reusable** - No changes needed!

---

## Generic Modal Design

### Strategy: Type Detection

Instead of creating a separate modal, we'll make `SuggestionReviewModal` **generic** to support both types.

**Design Pattern**: Type union with `interface{}`

```go
type SuggestionReviewModal struct {
    // Polymorphic suggestion storage
    suggestions     interface{}  // []burstfact.BurstSuggestion OR []SkillSuggestion
    suggestionType  string       // "burst" or "skill"
    
    // Separate TableBehavior instances (type-safe)
    burstTable      *behaviors.TableBehavior[burstfact.BurstSuggestion]
    skillTable      *behaviors.TableBehavior[SkillSuggestion]
    
    // Shared state
    theme           themes.Theme
    action          SuggestionAction
    accepted        []interface{}  // Polymorphic accepted list
    visible         bool
    width           int
    height          int
}
```

### Constructor Overloading

**For Burst Suggestions**:
```go
func NewBurstSuggestionModal(
    suggestions []burstfact.BurstSuggestion,
    theme themes.Theme,
) *SuggestionReviewModal {
    // ... existing implementation
    return &SuggestionReviewModal{
        suggestions:    suggestions,
        suggestionType: "burst",
        burstTable:     table,
        theme:          theme,
        // ...
    }
}
```

**For Skill Suggestions**:
```go
func NewSkillSuggestionModal(
    suggestions []SkillSuggestion,
    theme themes.Theme,
) *SuggestionReviewModal {
    // Sort by confidence (same as burst)
    sortedSuggestions := make([]SkillSuggestion, len(suggestions))
    copy(sortedSuggestions, suggestions)
    sort.Slice(sortedSuggestions, func(i, j int) bool {
        return sortedSuggestions[i].Confidence > sortedSuggestions[j].Confidence
    })
    
    // Define 4 columns (Name, Category, Events, Confidence)
    columns := []behaviors.ColumnDef{
        {Title: "Name", Width: 25},
        {Title: "Category", Width: 12},
        {Title: "Events", Width: 8},
        {Title: "Confidence", Width: 24},
    }
    
    // Formatter with category column
    formatter := func(s SkillSuggestion, _ int) []string {
        confidenceBar := primitives.CompactBar(s.Confidence, 15, nil).
            ShowPercentage(true).
            Render()
        return []string{
            s.Name,
            s.Category,
            strconv.Itoa(len(s.EventIDs)),
            confidenceBar,
        }
    }
    
    table := behaviors.NewTableBehavior(theme, columns, formatter).
        EmptyMessage("No skill suggestions available").
        PaginationPrefix("Skill Suggestions").
        PageSize(10)
    
    table.SetItems(sortedSuggestions)
    
    return &SuggestionReviewModal{
        suggestions:    suggestions,
        suggestionType: "skill",
        skillTable:     table,
        theme:          theme,
        accepted:       []interface{}{},
        visible:        true,
        width:          80,
        height:         24,
    }
}
```

### Rendering with Type Switch

```go
func (m *SuggestionReviewModal) buildContent() string {
    switch m.suggestionType {
    case "burst":
        return m.buildBurstContent()
    case "skill":
        return m.buildSkillContent()
    default:
        return "Unknown suggestion type"
    }
}

func (m *SuggestionReviewModal) buildSkillContent() string {
    if m.skillTable == nil {
        return "No skill suggestions available"
    }
    
    theme := m.getTheme()
    var content strings.Builder
    
    // Table
    content.WriteString(m.skillTable.Render())
    content.WriteString("\n\n")
    
    // Detail section (show selected skill)
    selected := m.skillTable.GetSelectedItem()
    if selected != nil {
        content.WriteString(primitives.NewText("Selected:", theme).Bold().Render())
        content.WriteString(" " + selected.Name + "\n")
        
        // Category badge
        categoryBadge := primitives.NewBadge(selected.Category, theme).Render()
        content.WriteString("Category: " + categoryBadge + "\n")
        
        // Usage contexts (NEW)
        if len(selected.Contexts) > 0 {
            content.WriteString("\n")
            content.WriteString(primitives.NewText("Usage Examples:", theme).Bold().Render())
            content.WriteString("\n")
            for i, ctx := range selected.Contexts {
                if i >= 3 { break }
                content.WriteString("  • " + ctx + "\n")
            }
        }
    }
    
    return content.String()
}
```

### Update Method with Type Switch

```go
func (m *SuggestionReviewModal) Update(msg tea.Msg) tea.Cmd {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "a", "A":
            return m.handleAccept()
        case "r", "R":
            return m.handleReject()
        // ... navigation keys
        }
    }
    
    // Delegate to correct table
    switch m.suggestionType {
    case "burst":
        if m.burstTable != nil {
            m.burstTable.Update(msg)
        }
    case "skill":
        if m.skillTable != nil {
            m.skillTable.Update(msg)
        }
    }
    
    return nil
}
```

---

## Implementation Phases

### Phase 1: Refactor Existing Modal

**Goal**: Make `SuggestionReviewModal` generic without breaking burst suggestions.

**Changes**:
1. Add `suggestionType` field
2. Add `skillTable` field (alongside `burstTable`)
3. Rename constructor to `NewBurstSuggestionModal()`
4. Add type switching to `buildContent()`, `Update()`, `View()`

**Tests**:
- Existing burst suggestion tests should still pass
- No behavioral changes for burst flow

### Phase 2: Add Skill Support

**Goal**: Add `NewSkillSuggestionModal()` constructor and skill-specific rendering.

**Changes**:
1. Create `NewSkillSuggestionModal()` constructor
2. Implement `buildSkillContent()` with 4 columns + contexts
3. Add type handling in `Update()` and `View()`

**Tests**:
- Create skill suggestion modal
- Display 4-column table
- Show context snippets
- Accept/reject skill suggestions

### Phase 3: Integration

**Goal**: Wire skill modal into intents.

**Changes**:
1. Add skill suggestion states to `burst_management/constants.go`
2. Add skill suggestion messages to `burst_management/messages.go`
3. Add handlers to `burst_management/handlers.go`
4. Add "i" key handler to `skillsmanagement/handlers.go`

---

## Loading States

### During Inference

**Burst Detail Modal** (after "i" key):
```
┌─────────────────────────────────────────────────┐
│ Analyzing Burst Events                          │
│                                                 │
│ ⏳ Detecting skills from 8 events...            │
│                                                 │
│ This may take a moment...                       │
└─────────────────────────────────────────────────┘
```

**ManageSkills List** (after "i" key):
```
┌─────────────────────────────────────────────────┐
│ Inferring Skills from All Events                │
│                                                 │
│ ⏳ Analyzing 42 events for technology mentions...│
│                                                 │
│ This may take a few moments...                  │
└─────────────────────────────────────────────────┘
```

**Component**: `primitives.Spinner()` + `primitives.Title()` + `primitives.Body()`

---

## Example UI Flows

### Flow 1: Automatic (After Burst Confirmation)

```
1. User confirms burst
   ├─ Facts extracted: "Backend refactor to microservices" (3 facts)
   └─ Skill inference triggers automatically

2. Skill suggestions found (3):
   ┌─────────────────────────────────────────────────────────────┐
   │ Skill Suggestions (from 5 burst events)                     │
   ├─────────────────────────────────────────────────────────────┤
   │ Name           Category    Events  Confidence               │
   │ Go             backend     5       ████████████████████░░ 95%│
   │ PostgreSQL     database    3       ███████████████░░░░░░ 75%│
   │ Docker         devops      2       ██████████░░░░░░░░░░ 50%│
   ├─────────────────────────────────────────────────────────────┤
   │ Selected: Go                                                │
   │ Category: backend                                           │
   │                                                             │
   │ Usage Examples:                                             │
   │   • "...built API using Go and gRPC for microservices..."  │
   │   • "...migrated monolith to Go for better performance..." │
   │   • "...wrote Go services that handle 10k req/s..."        │
   ├─────────────────────────────────────────────────────────────┤
   │ [a] Accept [r] Reject [↑↓] Navigate [Esc] Cancel           │
   └─────────────────────────────────────────────────────────────┘

3. User presses 'a' (accept)
   ├─ Skill "Go" created/reused
   ├─ Linked to 5 burst events
   ├─ LastUsed set to most recent event date
   └─ Next suggestion shown (PostgreSQL)

4. User presses 'A' (accept all)
   ├─ All 3 skills created/linked
   └─ Modal closes, returns to burst list
```

### Flow 2: Manual (Burst Detail Modal)

```
1. User views burst detail modal
   ├─ Shows burst facts, events, dates
   └─ Help badge shows: [i] Infer skills

2. User presses 'i'
   ├─ Loading modal appears
   ├─ Service analyzes 8 burst events
   └─ Returns 2 skill suggestions

3. Skill suggestion modal shown
   ┌─────────────────────────────────────────────────────────────┐
   │ Skill Suggestions (from 8 events)                           │
   ├─────────────────────────────────────────────────────────────┤
   │ Name           Category    Events  Confidence               │
   │ Kubernetes     devops      8       ████████████████████░░ 95%│
   │ Helm           devops      4       ███████████████░░░░░░ 75%│
   ├─────────────────────────────────────────────────────────────┤
   │ Selected: Kubernetes                                        │
   │ Category: devops                                            │
   │                                                             │
   │ Usage Examples:                                             │
   │   • "...deployed to Kubernetes clusters using Helm..."     │
   │   • "...managed k8s infrastructure for 20+ services..."    │
   │   • "...wrote Kubernetes operators for custom resources..."│
   ├─────────────────────────────────────────────────────────────┤
   │ [a] Accept [r] Reject [↑↓] Navigate [Esc] Cancel           │
   └─────────────────────────────────────────────────────────────┘

4. User accepts both
   └─ Returns to burst detail modal
```

### Flow 3: Manual (ManageSkills List)

```
1. User in ManageSkills list
   ├─ Shows existing skills (if any)
   └─ Help badge shows: [i] Infer from all events

2. User presses 'i'
   ├─ Loading modal appears
   ├─ Service analyzes ALL 42 events
   └─ Returns 15 skill suggestions

3. Skill suggestion modal shown (paginated)
   ┌─────────────────────────────────────────────────────────────┐
   │ Skill Suggestions (Page 1/2)                                │
   ├─────────────────────────────────────────────────────────────┤
   │ Name           Category    Events  Confidence               │
   │ Go             backend     22      ████████████████████░░ 95%│
   │ PostgreSQL     database    18      ███████████████░░░░░░ 75%│
   │ React          frontend    15      ███████████████░░░░░░ 75%│
   │ Docker         devops      12      ██████████░░░░░░░░░░ 50%│
   │ AWS            cloud       10      ██████████░░░░░░░░░░ 50%│
   │ Kubernetes     devops      10      ████████████████████░░ 95%│
   │ TypeScript     frontend    8       ███████████████░░░░░░ 75%│
   │ GraphQL        tooling     5       ██████████░░░░░░░░░░ 50%│
   │ Redis          database    4       ███████████████░░░░░░ 75%│
   │ Terraform      devops      3       ██████████░░░░░░░░░░ 50%│
   ├─────────────────────────────────────────────────────────────┤
   │ [Ctrl+D] Next page                                          │
   │ [a] Accept [r] Reject [A] Accept all [Esc] Cancel          │
   └─────────────────────────────────────────────────────────────┘

4. User presses 'A' (accept all)
   ├─ All 15 skills created/linked
   └─ Returns to skills list showing new skills
```

---

## Component Reuse Summary

### 100% Reusable Components (No Changes)

| Component | Location | Usage |
|-----------|----------|-------|
| `TableBehavior[T]` | `behaviors/table_behavior.go` | Generic table with type parameter |
| `CompactBar()` | `uikit/primitives/bar.go` | Confidence visualization |
| `AcceptBadge()` | `uikit/primitives/badge.go` | "Accept" help key |
| `RejectBadge()` | `uikit/primitives/badge.go` | "Reject" help key |
| `NavigateBadge()` | `uikit/primitives/badge.go` | "Navigate" help key |
| `PageVimBadge()` | `uikit/primitives/badge.go` | Page navigation help |
| `CancelBadge()` | `uikit/primitives/badge.go` | "Cancel" help key |
| `Box()` | `uikit/containers/box.go` | Modal container |
| `Overlay()` | `uikit/containers/overlay.go` | Modal overlay background |
| `NewText()` | `uikit/primitives/text.go` | Styled text |
| `NewBadge()` | `uikit/primitives/badge.go` | Category badge |

### Adapted Components (Minor Changes)

| Component | Original | Adaptation | Effort |
|-----------|----------|------------|--------|
| Column definitions | 3 columns | 4 columns (add Category) | 5 min |
| Row formatter | 3 strings | 4 strings (add category field) | 5 min |
| Detail section | Name + Description | Name + Category + Contexts | 15 min |
| Constructor | `NewSuggestionReviewModal()` | Split into `NewBurst...()` + `NewSkill...()` | 20 min |

### New Components (To Build)

| Component | Purpose | Complexity | Effort |
|-----------|---------|------------|--------|
| Type switching logic | Route to burst/skill rendering | Low | 30 min |
| `buildSkillContent()` | Render skill table + contexts | Medium | 45 min |
| Context rendering | Display 1-3 usage snippets | Low | 15 min |
| Loading modals | Show "Analyzing..." state | Low | 20 min |

---

## Code Reuse Statistics

| Category | Lines of Code | Percentage |
|----------|---------------|------------|
| **Unchanged** (table, badges, box) | ~200 LOC | 70% |
| **Adapted** (columns, formatter, detail) | ~60 LOC | 20% |
| **New** (type switch, contexts, loading) | ~30 LOC | 10% |
| **TOTAL** | ~290 LOC | 100% |

**Reuse Rate**: **90%** (70% unchanged + 20% adapted)

---

## Testing Strategy

### Unit Tests

**Modal Tests**:
- `TestNewBurstSuggestionModal` - Create burst modal
- `TestNewSkillSuggestionModal` - Create skill modal
- `TestModalTypeSwitch` - Correct rendering based on type
- `TestSkillTableColumns` - 4 columns (Name, Category, Events, Confidence)
- `TestSkillDetailContext` - Context snippets displayed (max 3)

**Service Tests** (separate from modal):
- `TestInferSkillsFromEvents` - Detection algorithm
- `TestConfidenceScoring` - High/medium/low patterns
- `TestCreateSkillsFromSuggestions` - Persistence and linking

### Integration Tests

**Burst Integration**:
- `TestBurstConfirmationSkillInference` - Automatic trigger
- `TestBurstDetailModalInference` - Manual "i" key trigger
- `TestSkillSuggestionAcceptance` - Accept flow
- `TestSkillSuggestionRejection` - Reject flow

**ManageSkills Integration**:
- `TestManageSkillsInference` - "i" key in skills list
- `TestGlobalEventAnalysis` - Analyze all events
- `TestAcceptAllSkills` - Bulk acceptance

---

## Accessibility & UX Considerations

### Keyboard Navigation

**Same as Burst Suggestions** (100% consistent):
- `j/k` or `↓/↑` - Navigate suggestions
- `g/G` - First/Last suggestion
- `Ctrl+D/U` - Page down/up
- `a` - Accept current suggestion
- `r` - Reject current suggestion
- `A` - Accept all remaining suggestions
- `Esc` - Cancel and close modal

### Visual Hierarchy

1. **Title**: "Skill Suggestions (from X events)" - Bold, top
2. **Table**: 4 columns with clear headers
3. **Detail Section**: Selected skill with category badge + contexts
4. **Footer**: Help badges (Accept, Reject, Navigate, Cancel)

### Empty States

**No Suggestions Found**:
```
┌─────────────────────────────────────────────────┐
│ Skill Suggestions                               │
│                                                 │
│ No technology keywords detected in events.      │
│                                                 │
│ Try adding descriptions with technologies like: │
│   • Programming languages (Go, Python, Java)   │
│   • Frameworks (React, Django, Spring)         │
│   • Databases (PostgreSQL, MongoDB, Redis)     │
│                                                 │
│ [Esc] Close                                     │
└─────────────────────────────────────────────────┘
```

---

## Performance Considerations

### Lazy Loading

**Current**: All suggestions loaded at once  
**Future Optimization**: Paginated backend if >100 suggestions

### Caching

**Keyword Map**: `GetKeywordMap()` creates map once, reused across service calls  
**Event Loading**: Batch load events instead of one-by-one

---

## Future Enhancements

### 1. Custom Keyword Management

Allow users to add/edit technology keywords without code changes:
- Store keywords in DB instead of hardcoded
- UI to add/remove/edit keywords
- Import/export keyword dictionaries

### 2. Machine Learning Integration

Replace regex with ML model for better detection:
- Train on user's accepted/rejected suggestions
- Improve confidence scoring over time
- Detect skills not in keyword dictionary

### 3. Skill Trend Visualization

Show skill usage over time:
- Timeline graph of skill mentions
- Detect emerging vs declining skills
- Highlight skills not used in 6+ months

---

## Summary

### Key Takeaways

1. **90% Code Reuse**: Existing `SuggestionReviewModal` is highly reusable with minor adaptations
2. **No New UI Components Needed**: All primitives (`TableBehavior`, `CompactBar`, badges) already exist
3. **Consistent UX**: Same keyboard shortcuts and patterns as burst suggestions
4. **Minimal Implementation Effort**: ~30 LOC new code, ~60 LOC adaptations, ~200 LOC reused

### Implementation Effort Estimate

| Task | Effort |
|------|--------|
| Refactor modal to generic pattern | 1-2h |
| Add skill-specific rendering | 1h |
| Add loading states | 30m |
| Integration (3 trigger points) | 2-3h |
| Testing (unit + integration) | 2-3h |
| **TOTAL** | **6.5-9.5h** (out of 12-16h total task time) |

**UI work is ~60% of total task** (remaining 40% is service layer implementation).

---

## References

- **Existing Pattern**: `internal/cli/screens/burst_management/modals/suggestion_review_modal.go`
- **UIKit Components**: `internal/cli/uikit/primitives/`, `internal/cli/uikit/containers/`
- **TableBehavior**: `internal/cli/behaviors/table_behavior.go`
- **Service Interface**: `internal/service/career/skill_inference.go` (to be created)
- **Data Structure**: `internal/service/career/skill_inference.go` (SkillSuggestion)

---

**Document Version**: 1.0  
**Last Updated**: 2026-01-30  
**Status**: Ready for Implementation

# Skill Inference Service - UI Integration Guide

## Overview

This guide explains how to integrate the Skill Inference Service with the UI layer, specifically the `burst_management` intent.

**Service Status**: ✅ **Production-Ready** (107 tests, 100% passing)

---

## Service Capabilities

The Skill Inference Service (`internal/service/career/skillinference`) provides:

✅ **Keyword Detection**: Word boundary regex matching for 224 technologies across 14 categories  
✅ **Confidence Scoring**: 3-tier pattern-based confidence (0.5/0.75/0.95)  
✅ **Skill Persistence**: Create/update skills with event linking  
✅ **Case-Insensitive**: Skill name matching and deduplication  
✅ **Context Extraction**: Usage examples for UI display  

---

## Integration Workflow

### Recommended Flow (After Burst Confirmation)

```
User confirms burst
     ↓
Extract facts (existing)
     ↓
**Infer skills** ← NEW STEP
     ↓
Show skill suggestions modal
     ↓
User accepts/rejects suggestions
     ↓
Create accepted skills
     ↓
Return to burst list
```

---

## Implementation Steps

### 1. Add Service to IntentContext

**File**: `internal/cli/intents/burst_management/context.go`

```go
type IntentContext struct {
    // ... existing fields ...
    
    // SkillInferenceService for detecting skills from events
    SkillInferenceService skillinference.SkillInferenceService
    
    // SkillRepository for skill persistence
    SkillRepository careerrepo.SkillRepository
    
    // EventRepository for event-skill linking
    EventRepository careerrepo.EventRepository
}
```

### 2. Add New States

**File**: `internal/cli/intents/burst_management/constants.go`

```go
const (
    // ... existing states ...
    
    // StateInferringSkills shows loading state during skill inference.
    StateInferringSkills State = "inferring_skills"
    
    // StateSkillSuggestionReview shows skill suggestions for review.
    StateSkillSuggestionReview State = "skill_suggestion_review"
)
```

### 3. Add Fields to Intent

**File**: `internal/cli/intents/burst_management/types.go`

```go
type Intent struct {
    // ... existing fields ...
    
    // --- Skill Inference ---
    
    // skillSuggestions holds detected skill suggestions
    skillSuggestions []skillinference.SkillSuggestion
    
    // skillSuggestionsLoading indicates if skill inference is in progress
    skillSuggestionsLoading bool
    
    // skillSuggestionsError stores any error from skill inference
    skillSuggestionsError error
    
    // skillSuggestionModal holds the skill suggestion review modal
    skillSuggestionModal *SkillSuggestionModal // Create this modal
}
```

### 4. Add Message Types

**File**: `internal/cli/intents/burst_management/messages.go`

```go
// SkillSuggestionsLoadedMsg is sent when skill suggestions are ready.
type SkillSuggestionsLoadedMsg struct {
    Suggestions []skillinference.SkillSuggestion
}

// SkillSuggestionsErrorMsg is sent when skill inference fails.
type SkillSuggestionsErrorMsg struct {
    Err error
}

// SkillSuggestionsAcceptedMsg is sent when user accepts skill suggestions.
type SkillSuggestionsAcceptedMsg struct {
    AcceptedSuggestions []skillinference.SkillSuggestion
}
```

### 5. Add Helper Function

**File**: `internal/cli/intents/burst_management/helpers.go`

```go
// inferSkillsFromBurst runs skill inference on burst events.
// Returns a Bubble Tea command that sends SkillSuggestionsLoadedMsg or SkillSuggestionsErrorMsg.
func (i *Intent) inferSkillsFromBurst(burst *career.Burst) tea.Cmd {
    return func() tea.Msg {
        // Load events for the burst
        filters := careerrepo.EventListFilters{BurstID: burst.ID}
        events, err := i.context.Service.ListEvents(i.getContext(), filters)
        if err != nil {
            return SkillSuggestionsErrorMsg{Err: err}
        }
        
        // Infer skills
        suggestions, err := i.context.SkillInferenceService.InferSkillsFromBurst(
            i.getContext(),
            burst,
            events,
        )
        if err != nil {
            return SkillSuggestionsErrorMsg{Err: err}
        }
        
        return SkillSuggestionsLoadedMsg{Suggestions: suggestions}
    }
}

// createSkillsFromSuggestions persists accepted skill suggestions.
func (i *Intent) createSkillsFromSuggestions(suggestions []skillinference.SkillSuggestion) tea.Cmd {
    return func() tea.Msg {
        skills, err := i.context.SkillInferenceService.CreateSkillsFromSuggestions(
            i.getContext(),
            suggestions,
        )
        if err != nil {
            return SkillSuggestionsErrorMsg{Err: err}
        }
        
        // Return success message
        return SkillsCreatedMsg{Skills: skills}
    }
}
```

### 6. Add Handler

**File**: `internal/cli/intents/burst_management/handlers.go`

```go
// handleSkillSuggestionsLoaded handles skill suggestions being loaded.
func (i *Intent) handleSkillSuggestionsLoaded(msg SkillSuggestionsLoadedMsg) tea.Cmd {
    i.skillSuggestions = msg.Suggestions
    i.skillSuggestionsLoading = false
    
    if len(msg.Suggestions) == 0 {
        // No skills detected - show info modal and return to list
        i.showInfoModal("No Skills Detected", "No technology skills were detected in this burst.")
        i.state = StateList
        return nil
    }
    
    // Show skill suggestion review modal
    i.skillSuggestionModal = NewSkillSuggestionModal(msg.Suggestions, i.theme)
    i.skillSuggestionModal.Show()
    i.state = StateSkillSuggestionReview
    
    return nil
}

// handleSkillSuggestionsAccepted handles user accepting skill suggestions.
func (i *Intent) handleSkillSuggestionsAccepted(msg SkillSuggestionsAcceptedMsg) tea.Cmd {
    // Create skills from accepted suggestions
    return i.createSkillsFromSuggestions(msg.AcceptedSuggestions)
}
```

### 7. Update Intent Flow

**File**: `internal/cli/intents/burst_management/intent.go`

```go
// In handleFactExtractionComplete(), add skill inference:

func (i *Intent) handleFactExtractionComplete(msg FactExtractionCompleteMsg) tea.Cmd {
    // ... existing fact extraction handling ...
    
    // After facts are extracted, infer skills
    i.skillSuggestionsLoading = true
    i.state = StateInferringSkills
    return i.inferSkillsFromBurst(i.selectedBurst)
}
```

---

## Create Skill Suggestion Modal

**File**: `internal/cli/screens/burst_management/modals/skill_suggestion_modal.go`

```go
package modals

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/baphled/kariya/internal/service/career/skillinference"
    "github.com/baphled/kariya/internal/themes"
)

// SkillSuggestionModal displays detected skill suggestions with checkboxes.
type SkillSuggestionModal struct {
    suggestions      []skillinference.SkillSuggestion
    selectedIndices  map[int]bool // Track checkbox state
    cursorPosition   int
    visible          bool
    theme            *themes.Theme
}

// NewSkillSuggestionModal creates a new skill suggestion modal.
func NewSkillSuggestionModal(suggestions []skillinference.SkillSuggestion, theme *themes.Theme) *SkillSuggestionModal {
    // Initialize with all suggestions selected by default
    selected := make(map[int]bool)
    for i := range suggestions {
        selected[i] = true
    }
    
    return &SkillSuggestionModal{
        suggestions:     suggestions,
        selectedIndices: selected,
        cursorPosition:  0,
        theme:           theme,
    }
}

// Update handles input events.
func (m *SkillSuggestionModal) Update(msg tea.Msg) (tea.Cmd, interface{}) {
    if !m.visible {
        return nil, nil
    }
    
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.Type {
        case tea.KeyUp, tea.KeyDown:
            // Navigate suggestions
        case tea.KeySpace:
            // Toggle checkbox
            m.selectedIndices[m.cursorPosition] = !m.selectedIndices[m.cursorPosition]
        case tea.KeyEnter:
            // Accept selected suggestions
            accepted := m.getAcceptedSuggestions()
            return nil, SkillSuggestionsAcceptedMsg{AcceptedSuggestions: accepted}
        case tea.KeyEsc:
            // Cancel - don't create any skills
            m.Hide()
            return nil, nil
        }
    }
    
    return nil, nil
}

// View renders the modal.
func (m *SkillSuggestionModal) View() string {
    if !m.visible {
        return ""
    }
    
    // Render title, suggestions with checkboxes, confidence scores, help text
    // Use UIKit primitives for rendering
    // Show:
    // - Skill name
    // - Category
    // - Confidence bar/score
    // - Usage context (truncated)
    // - Checkbox state
    // - Help: Space=toggle, Enter=accept, Esc=cancel
}
```

---

## Testing Strategy

### Unit Tests

**File**: `internal/cli/intents/burst_management/helpers_test.go`

```go
Describe("Skill Inference", func() {
    It("should infer skills from burst events", func() {
        // Setup: Create burst with events
        // Call: inferSkillsFromBurst()
        // Assert: SkillSuggestionsLoadedMsg received with suggestions
    })
    
    It("should handle skill inference errors", func() {
        // Setup: Mock service to return error
        // Call: inferSkillsFromBurst()
        // Assert: SkillSuggestionsErrorMsg received
    })
    
    It("should create skills from accepted suggestions", func() {
        // Setup: Create suggestions
        // Call: createSkillsFromSuggestions()
        // Assert: Skills created in repository
    })
})
```

### Integration Tests

See `internal/service/career/skillinference/integration_test.go` for complete workflow tests:
- Burst confirmation → skill inference → acceptance → skill creation
- Partial acceptance of suggestions
- Reuse of existing skills

---

## UI/UX Considerations

### Modal Design

**Skill Suggestion Modal** should display:

1. **Header**: "Detected Skills from [Burst Name]"
2. **Subtitle**: "[N] skills detected with confidence scores"
3. **Skill List** (scrollable):
   - ☑ Checkbox (checked by default)
   - **Skill Name** (bold)
   - Category badge (e.g., [Backend], [Database])
   - Confidence bar: `████████░░ 80%` (visual + percentage)
   - Context snippet (1 line, truncated)
4. **Footer**: 
   - Help: `Space=Toggle • Enter=Accept • Esc=Cancel`
   - Summary: `[N] of [Total] selected`

### Keyboard Navigation

| Key | Action |
|-----|--------|
| `↑`/`↓` | Navigate suggestions |
| `Space` | Toggle checkbox for current suggestion |
| `a` | Select/deselect all |
| `Enter` | Accept selected suggestions and create skills |
| `Esc` | Cancel without creating skills |
| `?` | Show help |

### Loading States

- **StateInferringSkills**: Show spinner with text "Detecting skills from burst..."
- **Creating Skills**: Show progress "Creating N skills..." (if needed)

### Error Handling

- **No skills detected**: Show info modal, return to burst list
- **Service error**: Show error modal with retry option
- **Repository error**: Show error modal, log details

---

## Performance Considerations

### Async Operations

Skill inference runs asynchronously (via `tea.Cmd`) to avoid blocking the UI:

```go
// Good: Non-blocking
return i.inferSkillsFromBurst(burst)

// Bad: Blocking UI
suggestions, _ := service.InferSkillsFromBurst(...)
```

### Caching

Consider caching skill suggestions per burst:

```go
type Intent struct {
    // Cache suggestions to avoid re-inference on modal re-open
    cachedSuggestions map[string][]skillinference.SkillSuggestion
}
```

### Cancellation

Support cancelling long-running inference:

```go
// Store cancel function
ctx, cancel := context.WithCancel(i.getContext())
i.cancelFunc = cancel

// Cancel on Esc
if i.state == StateInferringSkills && keyMsg.Type == tea.KeyEsc {
    if i.cancelFunc != nil {
        i.cancelFunc()
    }
}
```

---

## Example: Complete Integration

See `docs/examples/skill_inference_integration_example.go` for a complete working example (to be created).

---

## Soft Skills Detection (Phase 11)

The skill inference system detects soft skills through the `CompetencyCategory` system.

### Soft Skill Categories

| Category | String Value | Example Text |
|----------|-------------|--------------|
| Communication | `communication` | "Presented technical design to stakeholders" |
| Collaboration | `collaboration` | "Collaborated with cross-functional partners" |
| Problem Solving | `problem-solving` | "Debugged and diagnosed root cause" |
| Project Management | `project-management` | "Planned sprint milestones and delivery schedule" |
| Architecture | `architecture` | "Architected distributed microservices platform" |

### How Soft Skills Are Detected

1. **Tag-based**: Events tagged with soft skill categories are detected immediately.
2. **Keyword-based**: Event text is analyzed for soft skill keywords (8-12 per category).
3. **Competency inference**: `InferCompetencies()` assigns soft skill categories to facts.
4. **Profile inference**: `categoryStrengthMapping` generates CV strength descriptions.
5. **UI selector**: `AllowedCategories` includes all 11 competency types.

### Integration Points

- **Constants**: `internal/constants/constants.go` - 5 new `CompetencyCategory` values.
- **Classification**: `internal/service/career/classification/classifier.go` - Keyword lists and tag mapping.
- **Fact Inference**: `internal/service/career/burstfact/classifier.go` - `InferCompetencies()` with soft skills.
- **UI Selector**: `internal/cli/uikit/selectors/category_selector.go` - 11 categories in `AllowedCategories`.
- **Profile**: `internal/service/career/cv/profile_inference.go` - Strength and value proposition mappings.

---

## Current Status

✅ **Service Layer**: Complete (112+ tests, 100% passing)
✅ **Technology Dictionary**: 224 keywords across 14 categories
✅ **UI Integration**: Complete (3 triggers, suggestion modal, events drill-down)
✅ **Modal Components**: Complete (skill suggestion review modal)
✅ **E2E Tests**: 364 passing
✅ **Soft Skills**: 5 categories, 37 new tests across 5 packages

---

## Completed Phases

1. ✅ **Phases 0-6**: Service layer with persistence (112 tests)
2. ✅ **Phases 7-9**: UI integration (3 triggers, suggestion modal, drill-down)
3. ✅ **Phase 10**: Expanded keyword dictionary (224 keywords, 14 categories)
4. ✅ **Phase 11**: Soft skills detection (5 categories, full stack integration)

---

## References

- Service Implementation: `internal/service/career/skillinference/`
- Technology Dictionary: `internal/service/career/technology/`
- Classification: `internal/service/career/classification/`
- Burstfact Classifier: `internal/service/career/burstfact/classifier.go`
- Profile Inference: `internal/service/career/cv/profile_inference.go`
- UI Category Selector: `internal/cli/uikit/selectors/category_selector.go`
- Constants: `internal/constants/constants.go`
- Integration Tests: `internal/service/career/skillinference/integration_test.go`
- E2E Tests: `internal/testutil/e2e/skill_inference_e2e_test.go`
- Modal Usage Guide: `docs/guides/MODAL_USAGE_GUIDE.md`

---

**Document Version**: 2.0
**Last Updated**: January 30, 2026
**Author**: AI Assistant (Claude Sonnet 4)

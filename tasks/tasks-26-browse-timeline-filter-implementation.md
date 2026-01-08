# Task 26: BrowseTimeline Filter Implementation

**Created**: 2026-01-08
**Status**: Ready for Implementation
**Priority**: HIGH
**Estimated Time**: 3-4 hours
**Related**: Codebase Audit (2026-01-08)

---

## Overview

BrowseTimeline advertises filter and search functionality in footer but pressing 'f' or '/' does nothing. Four filter types are defined but never implemented. Pagination displays all events instead of current page.

**Issues**:
1. Footer shows "f Filter" but no 'f' key handler exists
2. Footer shows "/ Search" but no '/' key handler exists
3. Companies, Categories, DateFrom, DateTo filters defined but not implemented
4. Pagination test failures - table shows ALL events instead of current page

---

## Files to Modify

- [ ] `internal/cli/intents/browse_timeline_intent.go`
- [ ] `internal/cli/intents/browse_timeline.go`
- [ ] `internal/cli/intents/browse_timeline_test.go`

---

## Implementation Plan

### Phase 1: Implement Filter UI (2 hours)

**Add 'f' key handler**:
```go
func (i *BrowseTimelineIntent) updateTimelineView(msg tea.Msg) tea.Cmd {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "f":
            i.state.showingFilterPanel = true
            i.initializeFilterInputs()
            return nil
        // ... existing handlers
```

**Add filter panel state**:
```go
type TimelineViewState struct {
    // ... existing fields
    showingFilterPanel bool
    filterInputs       map[string]textinput.Model
    focusedFilter      int
}

func (i *BrowseTimelineIntent) initializeFilterInputs() {
    i.state.filterInputs = make(map[string]textinput.Model)
    
    // Company filter
    companyInput := textinput.New()
    companyInput.Placeholder = "Filter by company"
    i.state.filterInputs["company"] = companyInput
    
    // Category filter
    categoryInput := textinput.New()
    categoryInput.Placeholder = "Filter by category"
    i.state.filterInputs["category"] = categoryInput
    
    // Date range
    dateFromInput := textinput.New()
    dateFromInput.Placeholder = "From (YYYY-MM-DD)"
    i.state.filterInputs["date_from"] = dateFromInput
    
    dateToInput := textinput.New()
    dateToInput.Placeholder = "To (YYYY-MM-DD)"
    i.state.filterInputs["date_to"] = dateToInput
}
```

**Implement missing filters in applyFilters()** (line 315-380):
```go
func (i *BrowseTimelineIntent) applyFilters() []*career.CareerEvent {
    // ... existing SearchText and Tags filters
    
    // Company filter
    if i.state.filters.Companies != nil && len(i.state.filters.Companies) > 0 {
        var companyFiltered []*career.CareerEvent
        for _, event := range filtered {
            for _, company := range i.state.filters.Companies {
                if strings.Contains(strings.ToLower(event.Company), 
                    strings.ToLower(company)) {
                    companyFiltered = append(companyFiltered, event)
                    break
                }
            }
        }
        filtered = companyFiltered
    }
    
    // Category filter
    if i.state.filters.Categories != nil && len(i.state.filters.Categories) > 0 {
        var categoryFiltered []*career.CareerEvent
        for _, event := range filtered {
            for _, filterCat := range i.state.filters.Categories {
                for _, eventCat := range event.Categories {
                    if strings.EqualFold(eventCat, filterCat) {
                        categoryFiltered = append(categoryFiltered, event)
                        break
                    }
                }
            }
        }
        filtered = categoryFiltered
    }
    
    // Date range filter
    if i.state.filters.DateFrom != nil {
        var dateFiltered []*career.CareerEvent
        for _, event := range filtered {
            if event.Date.After(*i.state.filters.DateFrom) || 
               event.Date.Equal(*i.state.filters.DateFrom) {
                dateFiltered = append(dateFiltered, event)
            }
        }
        filtered = dateFiltered
    }
    
    if i.state.filters.DateTo != nil {
        var dateFiltered []*career.CareerEvent
        for _, event := range filtered {
            if event.Date.Before(*i.state.filters.DateTo) || 
               event.Date.Equal(*i.state.filters.DateTo) {
                dateFiltered = append(dateFiltered, event)
            }
        }
        filtered = dateFiltered
    }
    
    return filtered
}
```

**Tasks**:
- [ ] Add 'f' key handler to show filter panel
- [ ] Add filter panel state to TimelineViewState
- [ ] Initialize filter inputs
- [ ] Implement Companies filter in applyFilters()
- [ ] Implement Categories filter in applyFilters()
- [ ] Implement DateFrom filter in applyFilters()
- [ ] Implement DateTo filter in applyFilters()
- [ ] Render filter panel in view
- [ ] Apply filters and refresh list
- [ ] Test each filter type

---

### Phase 2: Implement Search UI (1 hour)

**Add '/' key handler**:
```go
case "/":
    i.state.searchMode = true
    i.state.searchInput = textinput.New()
    i.state.searchInput.Placeholder = "Search events..."
    i.state.searchInput.Focus()
    return nil
```

**Handle search input**:
```go
if i.state.searchMode {
    var cmd tea.Cmd
    i.state.searchInput, cmd = i.state.searchInput.Update(msg)
    
    // Update search filter as user types
    i.state.filters.SearchText = i.state.searchInput.Value()
    i.refreshFilteredEvents()
    
    // Exit search on Enter or Esc
    if msg.String() == "enter" || msg.String() == "esc" {
        i.state.searchMode = false
    }
    
    return cmd
}
```

**Tasks**:
- [ ] Add '/' key handler
- [ ] Add searchMode and searchInput to state
- [ ] Handle search input updates
- [ ] Update SearchText filter as user types
- [ ] Refresh filtered events on each keystroke
- [ ] Exit search mode on Enter/Esc
- [ ] Test real-time search filtering

---

### Phase 3: Fix Pagination (1 hour)

**Location**: Test failures at lines 478-497

**Issue**: Table shows ALL events instead of slicing by page

**Fix pagination in viewTimelineView()**:
```go
func (i *BrowseTimelineIntent) viewTimelineView() string {
    filtered := i.applyFilters()
    
    // Calculate pagination
    totalEvents := len(filtered)
    eventsPerPage := 20
    totalPages := (totalEvents + eventsPerPage - 1) / eventsPerPage
    
    // Ensure currentPage is in bounds
    if i.state.currentPage < 1 {
        i.state.currentPage = 1
    }
    if i.state.currentPage > totalPages {
        i.state.currentPage = totalPages
    }
    
    // Calculate slice bounds for current page
    startIdx := (i.state.currentPage - 1) * eventsPerPage
    endIdx := startIdx + eventsPerPage
    if endIdx > totalEvents {
        endIdx = totalEvents
    }
    
    // Slice events for current page ONLY
    pageEvents := filtered[startIdx:endIdx]
    
    // Convert to table rows (pageEvents, not filtered!)
    rows := eventsToRows(pageEvents)
    
    // ... render table with rows
}
```

**Tasks**:
- [ ] Calculate correct page bounds
- [ ] Slice filtered events by current page
- [ ] Pass only current page events to table
- [ ] Update pagination display to show current/total pages
- [ ] Test pagination shows correct events per page
- [ ] Fix failing pagination tests

---

### Phase 4: Vim Navigation (30 min)

**Add g/G for jump to top/bottom**:
```go
case "g":
    i.state.currentPage = 1
    i.state.selectedIndex = 0
    return nil
    
case "G":
    filtered := i.applyFilters()
    totalPages := (len(filtered) + 19) / 20
    i.state.currentPage = totalPages
    i.state.selectedIndex = 0
    return nil
```

**Tasks**:
- [ ] Add 'g' handler (jump to first page)
- [ ] Add 'G' handler (jump to last page)
- [ ] Test vim jump navigation

---

## Acceptance Criteria

### Must Have
- [ ] Pressing 'f' opens filter panel
- [ ] Can filter by company (contains match)
- [ ] Can filter by category (exact match)
- [ ] Can filter by date range (from/to)
- [ ] Pressing '/' opens search input
- [ ] Search filters events in real-time
- [ ] Pagination displays correct page of events (not all)
- [ ] All tests passing (2,078/2,078)
- [ ] Failing pagination tests now pass

### Should Have
- [ ] Vim g/G navigation for jump to top/bottom
- [ ] Filter panel shows current filter values
- [ ] Can clear filters
- [ ] Pagination shows "Page X of Y"

---

## Testing Strategy

```go
It("should filter by company", func() {
    intent.state.filters.Companies = []string{"Acme Corp"}
    filtered := intent.applyFilters()
    for _, event := range filtered {
        Expect(event.Company).To(ContainSubstring("Acme"))
    }
})

It("should show only current page of events", func() {
    // Create 50 events
    intent.events = createTestEvents(50)
    intent.state.currentPage = 2
    
    view := intent.View()
    rows := extractTableRows(view)
    
    // Should show 20 events (page 2: events 21-40)
    Expect(rows).To(HaveLen(20))
})
```

---

## References

- `internal/cli/intents/browse_timeline_intent.go` - Implementation
- `internal/cli/intents/browse_timeline_test.go` - Tests (lines 478-497 failing)
- `docs/TUI_STANDARDS.md` - Filter UI patterns

---

**Last Updated**: 2026-01-08
**Status**: Ready for implementation

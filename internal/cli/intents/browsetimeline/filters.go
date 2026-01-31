package browsetimeline

import (
	"sort"
	"time"

	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/domain/career"
)

// applyFilters filters the events based on current filter state.
func (i *Intent) applyFilters() {
	filtered := make([]*career.Event, 0)

	for _, evt := range i.context.Events {
		if !i.eventMatchesFilters(evt) {
			continue
		}
		filtered = append(filtered, evt)
	}

	i.sortEvents(filtered)
	i.filteredEvents = filtered
}

// eventMatchesFilters checks if an event passes all current filters.
func (i *Intent) eventMatchesFilters(evt *career.Event) bool {
	// Apply text search (if specified).
	if i.filters.SearchText != "" {
		categoriesStr := ""
		if len(evt.Categories) > 0 {
			categoriesStr = evt.Categories[0]
		}
		if !behaviors.SearchableFields(i.filters.SearchText, evt.Text, evt.Company, categoriesStr, evt.Project) {
			return false
		}
	}

	// Apply date range filter (if specified).
	if i.filters.DateFrom != "" {
		dateFrom, err := time.Parse("2006-01-02", i.filters.DateFrom)
		if err == nil && evt.Date.Before(dateFrom) {
			return false
		}
	}
	if i.filters.DateTo != "" {
		dateTo, err := time.Parse("2006-01-02", i.filters.DateTo)
		if err == nil && evt.Date.After(dateTo) {
			return false
		}
	}

	// Apply tag filters.
	if len(i.filters.Tags) > 0 && !i.eventHasAnyTag(evt) {
		return false
	}

	// Apply company filter.
	if len(i.filters.Companies) > 0 && !i.eventHasAnyCompany(evt) {
		return false
	}

	// Apply category filter.
	if len(i.filters.Categories) > 0 && !i.eventHasAnyCategory(evt) {
		return false
	}

	// Apply project filter.
	if len(i.filters.Projects) > 0 && !i.eventHasAnyProject(evt) {
		return false
	}

	return true
}

// eventHasAnyTag checks if event has any of the filter tags.
func (i *Intent) eventHasAnyTag(evt *career.Event) bool {
	for _, filterTag := range i.filters.Tags {
		for _, evtTag := range evt.Tags {
			if evtTag == filterTag {
				return true
			}
		}
	}
	return false
}

// eventHasAnyCompany checks if event matches any filter company.
func (i *Intent) eventHasAnyCompany(evt *career.Event) bool {
	for _, filterCompany := range i.filters.Companies {
		if evt.Company == filterCompany {
			return true
		}
	}
	return false
}

// eventHasAnyCategory checks if event has any of the filter categories.
func (i *Intent) eventHasAnyCategory(evt *career.Event) bool {
	for _, filterCat := range i.filters.Categories {
		for _, evtCat := range evt.Categories {
			if evtCat == filterCat {
				return true
			}
		}
	}
	return false
}

// eventHasAnyProject checks if event matches any filter project.
func (i *Intent) eventHasAnyProject(evt *career.Event) bool {
	for _, filterProject := range i.filters.Projects {
		if evt.Project == filterProject {
			return true
		}
	}
	return false
}

// sortEvents sorts the filtered events based on current sort settings.
func (i *Intent) sortEvents(filtered []*career.Event) {
	switch i.filters.SortBy {
	case "date":
		sort.Slice(filtered, func(a, b int) bool {
			if i.filters.SortOrder == "asc" {
				return filtered[a].Date.Before(filtered[b].Date)
			}
			return filtered[a].Date.After(filtered[b].Date)
		})
	case "text":
		sort.Slice(filtered, func(a, b int) bool {
			if i.filters.SortOrder == "asc" {
				return filtered[a].Text < filtered[b].Text
			}
			return filtered[a].Text > filtered[b].Text
		})
	}
}

// HasActiveFilters inspects the current filter state to determine whether
// any user-applied filters differ from defaults.
//
// Returns:
//   - True if any filter field has a non-default value, false otherwise.
//
// Side effects:
//   - None.
func (i *Intent) HasActiveFilters() bool {
	f := i.filters
	if f == nil {
		return false
	}

	return f.SearchText != "" ||
		len(f.Tags) > 0 ||
		len(f.Companies) > 0 ||
		len(f.Categories) > 0 ||
		len(f.Projects) > 0 ||
		f.DateFrom != "" ||
		f.DateTo != "" ||
		(f.SortBy != "" && f.SortBy != "date") ||
		(f.SortOrder != "" && f.SortOrder != "desc")
}

// ClearFilters removes the most recently applied filter layer, or resets
// all filters to defaults if the filter stack is empty.
//
// Side effects:
//   - Pops the top filter layer from the stack and clears its corresponding
//     filter field, or resets all filters when the stack is empty.
func (i *Intent) ClearFilters() {
	if i.filterStack == nil || i.filterStack.IsEmpty() {
		i.clearAllFilters()
		return
	}

	layer := i.filterStack.Pop()

	switch layer {
	case behaviors.FilterLayerSearch:
		i.filters.SearchText = ""
	case behaviors.FilterLayerCompany:
		i.filters.Companies = []string{}
	case behaviors.FilterLayerCategory:
		i.filters.Categories = []string{}
	case behaviors.FilterLayerProject:
		i.filters.Projects = []string{}
	case behaviors.FilterLayerTags:
		i.filters.Tags = []string{}
	case behaviors.FilterLayerSort:
		i.filters.SortBy = "date"
		i.filters.SortOrder = "desc"
	}

	if !i.HasActiveFilters() {
		i.filterStack.Clear()
	}
}

// clearAllFilters resets all filters to default state.
func (i *Intent) clearAllFilters() {
	i.filters = &Filters{
		SearchText: "",
		Tags:       []string{},
		Companies:  []string{},
		Categories: []string{},
		Projects:   []string{},
		DateFrom:   "",
		DateTo:     "",
		SortBy:     "date",
		SortOrder:  "desc",
	}
	if i.filterStack != nil {
		i.filterStack.Clear()
	}
}

// ApplyFilters is the public entry point for re-evaluating all events
// against the current filter and sort criteria.
//
// Side effects:
//   - Rebuilds the filtered events list and re-sorts it in place.
func (i *Intent) ApplyFilters() {
	i.applyFilters()
}

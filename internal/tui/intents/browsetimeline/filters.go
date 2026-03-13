package browsetimeline

import (
	"sort"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/ui/behaviors"
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
	if !i.matchesSearchText(evt) {
		return false
	}
	if !i.matchesDateRange(evt) {
		return false
	}
	if !i.matchesTagFilters(evt) {
		return false
	}
	if !i.matchesCompanyFilters(evt) {
		return false
	}
	if !i.matchesCategoryFilters(evt) {
		return false
	}
	return i.matchesProjectFilters(evt)
}

// matchesSearchText checks whether the event matches the search text filter.
//
// Expected:
//   - evt is a non-nil event.
//
// Returns:
//   - True when the event matches the search text or no search is set.
//
// Side effects:
//   - None.
func (i *Intent) matchesSearchText(evt *career.Event) bool {
	if i.filters.SearchText == "" {
		return true
	}

	categoriesStr := ""
	if len(evt.Categories) > 0 {
		categoriesStr = evt.Categories[0]
	}
	return behaviors.SearchableFields(i.filters.SearchText, evt.Text, evt.Company, categoriesStr, evt.Project)
}

// matchesDateRange checks whether the event is inside the date range filters.
//
// Expected:
//   - evt is a non-nil event.
//
// Returns:
//   - True when the event matches date filters or no date filters are set.
//
// Side effects:
//   - None.
func (i *Intent) matchesDateRange(evt *career.Event) bool {
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
	return true
}

// matchesTagFilters checks whether the event matches tag filters.
//
// Expected:
//   - evt is a non-nil event.
//
// Returns:
//   - True when the event matches tag filters or no tags are set.
//
// Side effects:
//   - None.
func (i *Intent) matchesTagFilters(evt *career.Event) bool {
	if len(i.filters.Tags) == 0 {
		return true
	}
	return i.eventHasAnyTag(evt)
}

// matchesCompanyFilters checks whether the event matches company filters.
//
// Expected:
//   - evt is a non-nil event.
//
// Returns:
//   - True when the event matches company filters or no companies are set.
//
// Side effects:
//   - None.
func (i *Intent) matchesCompanyFilters(evt *career.Event) bool {
	if len(i.filters.Companies) == 0 {
		return true
	}
	return i.eventHasAnyCompany(evt)
}

// matchesCategoryFilters checks whether the event matches category filters.
//
// Expected:
//   - evt is a non-nil event.
//
// Returns:
//   - True when the event matches category filters or no categories are set.
//
// Side effects:
//   - None.
func (i *Intent) matchesCategoryFilters(evt *career.Event) bool {
	if len(i.filters.Categories) == 0 {
		return true
	}
	return i.eventHasAnyCategory(evt)
}

// matchesProjectFilters checks whether the event matches project filters.
//
// Expected:
//   - evt is a non-nil event.
//
// Returns:
//   - True when the event matches project filters or no projects are set.
//
// Side effects:
//   - None.
func (i *Intent) matchesProjectFilters(evt *career.Event) bool {
	if len(i.filters.Projects) == 0 {
		return true
	}
	return i.eventHasAnyProject(evt)
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
//
// Returns:
//   - A bool value.
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
//
// Side effects:
//   - None.
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
//
// Side effects:
//   - None.
func (i *Intent) ApplyFilters() {
	i.applyFilters()
}

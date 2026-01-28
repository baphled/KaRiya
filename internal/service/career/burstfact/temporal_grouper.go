package burstfact

import (
	"sort"
	"time"
)

// TemporalGrouper provides methods to group events based on temporal proximity
// Events within 6 months are considered temporally related
type TemporalGrouper struct{}

// NewTemporalGrouper creates a new TemporalGrouper instance
func NewTemporalGrouper() *TemporalGrouper {
	return &TemporalGrouper{}
}

// IsTemporallyRelated checks if two events are within 6 months of each other
// Returns true if the absolute difference between dates is ≤ 6 months
func (tg *TemporalGrouper) IsTemporallyRelated(date1, date2 time.Time) bool {
	if date1.IsZero() || date2.IsZero() {
		return false
	}

	monthsDiff := tg.MonthsDifference(date1, date2)
	return monthsDiff <= 6
}

// MonthsDifference calculates the absolute number of months between two dates
// Returns the difference as an integer number of months (always positive)
// If date2 has remaining days after reaching the month anniversary from date1,
// those days are counted as requiring an additional month (ceiling semantics)
func (tg *TemporalGrouper) MonthsDifference(date1, date2 time.Time) int {
	// Ensure date1 is earlier than date2 for consistent calculation
	if date1.After(date2) {
		date1, date2 = date2, date1
	}

	months := 0
	current := date1

	// Count complete months
	for {
		nextMonth := tg.addMonths(current, 1)
		if nextMonth.After(date2) {
			break
		}
		months++
		current = nextMonth
	}

	// Check if there are remaining days that exceed the anniversary
	// If date2 is after the current month anniversary, and has days beyond the anniversary day,
	// count it as needing one more month
	if current.Before(date2) && date2.Day() > current.Day() {
		months++
	}

	return months
}

// addMonths adds a specified number of months to a date, handling day-of-month overflow
// For example, Jan 31 + 1 month = Feb 28 (or Feb 29 in leap years)
func (tg *TemporalGrouper) addMonths(date time.Time, months int) time.Time {
	newMonth := int(date.Month()) + months
	newYear := date.Year()

	// Handle year overflow
	for newMonth > 12 {
		newMonth -= 12
		newYear++
	}

	// Handle day overflow (e.g., Jan 31 + 1 month should be Feb 28, not Feb 31)
	day := date.Day()
	maxDay := tg.daysInMonth(newYear, time.Month(newMonth))
	if day > maxDay {
		day = maxDay
	}

	return time.Date(newYear, time.Month(newMonth), day, date.Hour(), date.Minute(), date.Second(), date.Nanosecond(), date.Location())
}

// daysInMonth returns the number of days in a given month
func (tg *TemporalGrouper) daysInMonth(year int, month time.Month) int {
	// Use the fact that time.Date with day=0 gives the last day of the previous month
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// GroupEventsByTemporal groups events into clusters where each cluster contains
// events that are within 6 months of each other in chronological order.
// Events are sorted chronologically before grouping.
func (tg *TemporalGrouper) GroupEventsByTemporal(events []time.Time) [][]time.Time {
	if len(events) == 0 {
		return [][]time.Time{}
	}

	if len(events) == 1 {
		return [][]time.Time{{events[0]}}
	}

	// Sort events chronologically
	sortedEvents := make([]time.Time, len(events))
	copy(sortedEvents, events)
	sort.Slice(sortedEvents, func(i, j int) bool {
		return sortedEvents[i].Before(sortedEvents[j])
	})

	clusters := [][]time.Time{}
	currentCluster := []time.Time{sortedEvents[0]}

	for i := 1; i < len(sortedEvents); i++ {
		// Check if current event is within 6 months of the last event in cluster
		if tg.IsTemporallyRelated(currentCluster[len(currentCluster)-1], sortedEvents[i]) {
			currentCluster = append(currentCluster, sortedEvents[i])
		} else {
			// Start a new cluster
			clusters = append(clusters, currentCluster)
			currentCluster = []time.Time{sortedEvents[i]}
		}
	}

	// Add the last cluster
	clusters = append(clusters, currentCluster)

	return clusters
}

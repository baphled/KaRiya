package fixtures

import (
	"time"

	repo "github.com/baphled/kariya/internal/repository/career"
)

// EventListFilters creates a minimal EventListFilters with sensible defaults.
//
// Returns:
//   - A fully initialized repo.EventListFilters ready for use.
//
// Side effects:
//   - None.
func EventListFilters() *repo.EventListFilters {
	return &repo.EventListFilters{}
}

// EventListFiltersWithLimit creates an EventListFilters with pagination.
//
// Expected:
//   - int must be valid.
//
// Returns:
//   - A fully initialized repo.EventListFilters ready for use.
//
// Side effects:
//   - None.
func EventListFiltersWithLimit(offset, limit int) *repo.EventListFilters {
	return &repo.EventListFilters{
		Offset: offset,
		Limit:  limit,
	}
}

// EventListFiltersWithTags creates an EventListFilters filtered by tags.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A fully initialized repo.EventListFilters ready for use.
//
// Side effects:
//   - None.
func EventListFiltersWithTags(tags []string) *repo.EventListFilters {
	return &repo.EventListFilters{
		Tags: tags,
	}
}

// EventListFiltersWithSort creates an EventListFilters with sorting.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A fully initialized repo.EventListFilters ready for use.
//
// Side effects:
//   - None.
func EventListFiltersWithSort(sortBy, sortOrder string) *repo.EventListFilters {
	return &repo.EventListFilters{
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}
}

// EventListFiltersWithDateRange creates an EventListFilters with date range.
//
// Expected:
//   - time must be valid.
//
// Returns:
//   - A fully initialized repo.EventListFilters ready for use.
//
// Side effects:
//   - None.
func EventListFiltersWithDateRange(startDate, endDate *time.Time) *repo.EventListFilters {
	return &repo.EventListFilters{
		StartDate: startDate,
		EndDate:   endDate,
	}
}

// FactListFilters creates a minimal FactListFilters with sensible defaults.
//
// Returns:
//   - A fully initialized repo.FactListFilters ready for use.
//
// Side effects:
//   - None.
func FactListFilters() *repo.FactListFilters {
	return &repo.FactListFilters{}
}

// FactListFiltersWithCategory creates a FactListFilters filtered by competency category.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A fully initialized repo.FactListFilters ready for use.
//
// Side effects:
//   - None.
func FactListFiltersWithCategory(category string) *repo.FactListFilters {
	return &repo.FactListFilters{
		CompetencyCategory: category,
	}
}

// FactListFiltersWithRole creates a FactListFilters filtered by role fit.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A fully initialized repo.FactListFilters ready for use.
//
// Side effects:
//   - None.
func FactListFiltersWithRole(roleFit string) *repo.FactListFilters {
	return &repo.FactListFilters{
		RoleFit: roleFit,
	}
}

// FactListFiltersWithLimit creates a FactListFilters with pagination.
//
// Expected:
//   - int must be valid.
//
// Returns:
//   - A fully initialized repo.FactListFilters ready for use.
//
// Side effects:
//   - None.
func FactListFiltersWithLimit(offset, limit int) *repo.FactListFilters {
	return &repo.FactListFilters{
		Offset: offset,
		Limit:  limit,
	}
}

// FactListFiltersWithSort creates a FactListFilters with sorting.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A fully initialized repo.FactListFilters ready for use.
//
// Side effects:
//   - None.
func FactListFiltersWithSort(sortBy, sortOrder string) *repo.FactListFilters {
	return &repo.FactListFilters{
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}
}

// FactListFiltersWithAudience creates a FactListFilters filtered by audience relevance.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A fully initialized repo.FactListFilters ready for use.
//
// Side effects:
//   - None.
func FactListFiltersWithAudience(audience string) *repo.FactListFilters {
	return &repo.FactListFilters{
		AudienceRelevance: audience,
	}
}

// BurstListFilters creates a minimal BurstListFilters with sensible defaults.
//
// Returns:
//   - A fully initialized repo.BurstListFilters ready for use.
//
// Side effects:
//   - None.
func BurstListFilters() *repo.BurstListFilters {
	return &repo.BurstListFilters{}
}

// BurstListFiltersWithLimit creates a BurstListFilters with pagination.
//
// Expected:
//   - int must be valid.
//
// Returns:
//   - A fully initialized repo.BurstListFilters ready for use.
//
// Side effects:
//   - None.
func BurstListFiltersWithLimit(offset, limit int) *repo.BurstListFilters {
	return &repo.BurstListFilters{
		Offset: offset,
		Limit:  limit,
	}
}

// BurstListFiltersWithSort creates a BurstListFilters with sorting.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A fully initialized repo.BurstListFilters ready for use.
//
// Side effects:
//   - None.
func BurstListFiltersWithSort(sortBy, sortOrder string) *repo.BurstListFilters {
	return &repo.BurstListFilters{
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}
}

// SkillListFilters creates a minimal SkillListFilters with sensible defaults.
//
// Returns:
//   - A fully initialized repo.SkillListFilters ready for use.
//
// Side effects:
//   - None.
func SkillListFilters() *repo.SkillListFilters {
	return &repo.SkillListFilters{}
}

// SkillListFiltersWithCategory creates a SkillListFilters filtered by category.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A fully initialized repo.SkillListFilters ready for use.
//
// Side effects:
//   - None.
func SkillListFiltersWithCategory(category string) *repo.SkillListFilters {
	return &repo.SkillListFilters{
		Category: category,
	}
}

// SkillListFiltersWithLimit creates a SkillListFilters with pagination.
//
// Expected:
//   - int must be valid.
//
// Returns:
//   - A fully initialized repo.SkillListFilters ready for use.
//
// Side effects:
//   - None.
func SkillListFiltersWithLimit(offset, limit int) *repo.SkillListFilters {
	return &repo.SkillListFilters{
		Offset: offset,
		Limit:  limit,
	}
}

// SkillListFiltersWithLevel creates a SkillListFilters filtered by level.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A fully initialized repo.SkillListFilters ready for use.
//
// Side effects:
//   - None.
func SkillListFiltersWithLevel(level string) *repo.SkillListFilters {
	return &repo.SkillListFilters{
		Level: level,
	}
}

// SkillListFiltersWithSort creates a SkillListFilters with sorting.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A fully initialized repo.SkillListFilters ready for use.
//
// Side effects:
//   - None.
func SkillListFiltersWithSort(sortBy, sortOrder string) *repo.SkillListFilters {
	return &repo.SkillListFilters{
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}
}

// SkillListFiltersWithMinEvents creates a SkillListFilters filtered by minimum event count.
//
// Expected:
//   - int must be valid.
//
// Returns:
//   - A fully initialized repo.SkillListFilters ready for use.
//
// Side effects:
//   - None.
func SkillListFiltersWithMinEvents(minEvents int) *repo.SkillListFilters {
	return &repo.SkillListFilters{
		MinEvents: minEvents,
	}
}

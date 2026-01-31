package fixtures

import (
	repo "github.com/baphled/kariya/internal/repository/career"
)

// EventListFilters creates a minimal EventListFilters with sensible defaults.
func EventListFilters() *repo.EventListFilters {
	return &repo.EventListFilters{}
}

// EventListFiltersWithLimit creates an EventListFilters with pagination.
func EventListFiltersWithLimit(offset, limit int) *repo.EventListFilters {
	return &repo.EventListFilters{
		Offset: offset,
		Limit:  limit,
	}
}

// EventListFiltersWithTags creates an EventListFilters filtered by tags.
func EventListFiltersWithTags(tags []string) *repo.EventListFilters {
	return &repo.EventListFilters{
		Tags: tags,
	}
}

// EventListFiltersWithSort creates an EventListFilters with sorting.
func EventListFiltersWithSort(sortBy, sortOrder string) *repo.EventListFilters {
	return &repo.EventListFilters{
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}
}

// FactListFilters creates a minimal FactListFilters with sensible defaults.
func FactListFilters() *repo.FactListFilters {
	return &repo.FactListFilters{}
}

// FactListFiltersWithCategory creates a FactListFilters filtered by competency category.
func FactListFiltersWithCategory(category string) *repo.FactListFilters {
	return &repo.FactListFilters{
		CompetencyCategory: category,
	}
}

// FactListFiltersWithRole creates a FactListFilters filtered by role fit.
func FactListFiltersWithRole(roleFit string) *repo.FactListFilters {
	return &repo.FactListFilters{
		RoleFit: roleFit,
	}
}

// FactListFiltersWithLimit creates a FactListFilters with pagination.
func FactListFiltersWithLimit(offset, limit int) *repo.FactListFilters {
	return &repo.FactListFilters{
		Offset: offset,
		Limit:  limit,
	}
}

// BurstListFilters creates a minimal BurstListFilters with sensible defaults.
func BurstListFilters() *repo.BurstListFilters {
	return &repo.BurstListFilters{}
}

// BurstListFiltersWithLimit creates a BurstListFilters with pagination.
func BurstListFiltersWithLimit(offset, limit int) *repo.BurstListFilters {
	return &repo.BurstListFilters{
		Offset: offset,
		Limit:  limit,
	}
}

// BurstListFiltersWithSort creates a BurstListFilters with sorting.
func BurstListFiltersWithSort(sortBy, sortOrder string) *repo.BurstListFilters {
	return &repo.BurstListFilters{
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}
}

// SkillListFilters creates a minimal SkillListFilters with sensible defaults.
func SkillListFilters() *repo.SkillListFilters {
	return &repo.SkillListFilters{}
}

// SkillListFiltersWithCategory creates a SkillListFilters filtered by category.
func SkillListFiltersWithCategory(category string) *repo.SkillListFilters {
	return &repo.SkillListFilters{
		Category: category,
	}
}

// SkillListFiltersWithLimit creates a SkillListFilters with pagination.
func SkillListFiltersWithLimit(offset, limit int) *repo.SkillListFilters {
	return &repo.SkillListFilters{
		Offset: offset,
		Limit:  limit,
	}
}

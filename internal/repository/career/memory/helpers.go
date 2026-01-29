package memory

import (
	"sync"
	"time"
)

// storeNew inserts entity into the store if its key is not already present.
// The caller must validate the entity and set timestamps before calling this.
func storeNew[T any](mu *sync.RWMutex, store map[string]T, id string, entity T, dupErr error) error {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := store[id]; exists {
		return dupErr
	}

	store[id] = entity

	return nil
}

// filterByDateRange returns items whose timestamp falls within the optional
// [start, end] range. The getTime function extracts the relevant timestamp
// from each item. Both bounds are inclusive; nil means unbounded.
func filterByDateRange[T any](items []T, getTime func(T) time.Time, start, end *time.Time) []T {
	if start == nil && end == nil {
		return items
	}

	var filtered []T
	for _, item := range items {
		t := getTime(item)
		if start != nil && t.Before(*start) {
			continue
		}
		if end != nil && t.After(*end) {
			continue
		}
		filtered = append(filtered, item)
	}

	return filtered
}

// paginate returns a subset of items based on offset and limit.
// A limit of 0 or less means no limit (return everything from offset).
func paginate[T any](items []T, offset, limit int) []T {
	if offset > len(items) {
		return nil
	}

	if limit <= 0 {
		return items[offset:]
	}

	end := offset + limit
	if end > len(items) {
		end = len(items)
	}

	return items[offset:end]
}

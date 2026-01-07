package components

import (
	"testing"
)

func TestNewPaginationHelper(t *testing.T) {
	ph := NewPaginationHelper(10)
	if ph.GetCurrentPage() != 1 {
		t.Errorf("Expected currentPage to be 1, got %d", ph.GetCurrentPage())
	}
	if ph.GetPageSize() != 10 {
		t.Errorf("Expected pageSize to be 10, got %d", ph.GetPageSize())
	}
	if ph.GetTotalCount() != 0 {
		t.Errorf("Expected totalCount to be 0, got %d", ph.GetTotalCount())
	}
}

func TestPaginationHelper_SetTotalCount(t *testing.T) {
	ph := NewPaginationHelper(10)
	ph.SetTotalCount(25)
	if ph.GetTotalCount() != 25 {
		t.Errorf("Expected totalCount to be 25, got %d", ph.GetTotalCount())
	}
}

func TestPaginationHelper_GetTotalPages(t *testing.T) {
	tests := []struct {
		name       string
		pageSize   int
		totalCount int
		want       int
	}{
		{"Empty list", 10, 0, 1},
		{"Exact fit", 10, 20, 2},
		{"Partial page", 10, 25, 3},
		{"Single item", 10, 1, 1},
		{"Page size boundary", 10, 10, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ph := NewPaginationHelper(tt.pageSize)
			ph.SetTotalCount(tt.totalCount)
			got := ph.GetTotalPages()
			if got != tt.want {
				t.Errorf("GetTotalPages() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestPaginationHelper_GetPageStartIndex(t *testing.T) {
	ph := NewPaginationHelper(10)
	ph.SetTotalCount(35)

	tests := []struct {
		page      int
		wantStart int
	}{
		{1, 0},
		{2, 10},
		{3, 20},
		{4, 30},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			// Navigate to page
			ph.GoToFirstPage()
			for i := 1; i < tt.page; i++ {
				ph.NextPage()
			}
			got := ph.GetPageStartIndex()
			if got != tt.wantStart {
				t.Errorf("GetPageStartIndex() on page %d = %d, want %d", tt.page, got, tt.wantStart)
			}
		})
	}
}

func TestPaginationHelper_GetPageEndIndex(t *testing.T) {
	ph := NewPaginationHelper(10)
	ph.SetTotalCount(35)

	tests := []struct {
		page    int
		wantEnd int
	}{
		{1, 10},
		{2, 20},
		{3, 30},
		{4, 40}, // End index is exclusive, so 40 even though there are only 35 items
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			// Navigate to page
			ph.GoToFirstPage()
			for i := 1; i < tt.page; i++ {
				ph.NextPage()
			}
			got := ph.GetPageEndIndex()
			if got != tt.wantEnd {
				t.Errorf("GetPageEndIndex() on page %d = %d, want %d", tt.page, got, tt.wantEnd)
			}
		})
	}
}

func TestPaginationHelper_GetPaginationInfo(t *testing.T) {
	tests := []struct {
		name       string
		pageSize   int
		totalCount int
		page       int
		itemName   string
		want       string
	}{
		{"Empty list", 10, 0, 1, "items", "Showing 0 of 0 items"},
		{"First page", 10, 25, 1, "events", "Showing 1-10 of 25 events"},
		{"Middle page", 10, 25, 2, "events", "Showing 11-20 of 25 events"},
		{"Last page partial", 10, 25, 3, "events", "Showing 21-25 of 25 events"},
		{"Exact fit", 10, 20, 2, "items", "Showing 11-20 of 20 items"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ph := NewPaginationHelper(tt.pageSize)
			ph.SetTotalCount(tt.totalCount)
			// Navigate to page
			for i := 1; i < tt.page; i++ {
				ph.NextPage()
			}
			got := ph.GetPaginationInfo(tt.itemName)
			if got != tt.want {
				t.Errorf("GetPaginationInfo() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPaginationHelper_NextPage(t *testing.T) {
	ph := NewPaginationHelper(10)
	ph.SetTotalCount(25)

	// Start on page 1
	if ph.GetCurrentPage() != 1 {
		t.Fatalf("Expected to start on page 1")
	}

	// Move to page 2
	ph.NextPage()
	if ph.GetCurrentPage() != 2 {
		t.Errorf("Expected page 2, got %d", ph.GetCurrentPage())
	}

	// Move to page 3 (last page)
	ph.NextPage()
	if ph.GetCurrentPage() != 3 {
		t.Errorf("Expected page 3, got %d", ph.GetCurrentPage())
	}

	// Try to move beyond last page - should stay on page 3
	ph.NextPage()
	if ph.GetCurrentPage() != 3 {
		t.Errorf("Expected to stay on page 3, got %d", ph.GetCurrentPage())
	}
}

func TestPaginationHelper_PrevPage(t *testing.T) {
	ph := NewPaginationHelper(10)
	ph.SetTotalCount(25)
	ph.GoToLastPage()

	// Start on last page (3)
	if ph.GetCurrentPage() != 3 {
		t.Fatalf("Expected to start on page 3")
	}

	// Move to page 2
	ph.PrevPage()
	if ph.GetCurrentPage() != 2 {
		t.Errorf("Expected page 2, got %d", ph.GetCurrentPage())
	}

	// Move to page 1
	ph.PrevPage()
	if ph.GetCurrentPage() != 1 {
		t.Errorf("Expected page 1, got %d", ph.GetCurrentPage())
	}

	// Try to move before first page - should stay on page 1
	ph.PrevPage()
	if ph.GetCurrentPage() != 1 {
		t.Errorf("Expected to stay on page 1, got %d", ph.GetCurrentPage())
	}
}

func TestPaginationHelper_GoToFirstPage(t *testing.T) {
	ph := NewPaginationHelper(10)
	ph.SetTotalCount(35)
	ph.GoToLastPage()

	if ph.GetCurrentPage() != 4 {
		t.Fatalf("Expected to be on last page (4)")
	}

	ph.GoToFirstPage()
	if ph.GetCurrentPage() != 1 {
		t.Errorf("Expected page 1, got %d", ph.GetCurrentPage())
	}
}

func TestPaginationHelper_GoToLastPage(t *testing.T) {
	ph := NewPaginationHelper(10)
	ph.SetTotalCount(35)

	if ph.GetCurrentPage() != 1 {
		t.Fatalf("Expected to start on page 1")
	}

	ph.GoToLastPage()
	if ph.GetCurrentPage() != 4 {
		t.Errorf("Expected page 4, got %d", ph.GetCurrentPage())
	}
}

func TestPaginationHelper_IsFirstPage(t *testing.T) {
	ph := NewPaginationHelper(10)
	ph.SetTotalCount(25)

	if !ph.IsFirstPage() {
		t.Error("Expected IsFirstPage() to be true on page 1")
	}

	ph.NextPage()
	if ph.IsFirstPage() {
		t.Error("Expected IsFirstPage() to be false on page 2")
	}
}

func TestPaginationHelper_IsLastPage(t *testing.T) {
	ph := NewPaginationHelper(10)
	ph.SetTotalCount(25)

	if ph.IsLastPage() {
		t.Error("Expected IsLastPage() to be false on page 1")
	}

	ph.GoToLastPage()
	if !ph.IsLastPage() {
		t.Error("Expected IsLastPage() to be true on last page")
	}
}

func TestPaginationHelper_CanGoNext(t *testing.T) {
	ph := NewPaginationHelper(10)
	ph.SetTotalCount(25)

	if !ph.CanGoNext() {
		t.Error("Expected CanGoNext() to be true on page 1")
	}

	ph.GoToLastPage()
	if ph.CanGoNext() {
		t.Error("Expected CanGoNext() to be false on last page")
	}
}

func TestPaginationHelper_CanGoPrev(t *testing.T) {
	ph := NewPaginationHelper(10)
	ph.SetTotalCount(25)

	if ph.CanGoPrev() {
		t.Error("Expected CanGoPrev() to be false on page 1")
	}

	ph.NextPage()
	if !ph.CanGoPrev() {
		t.Error("Expected CanGoPrev() to be true on page 2")
	}
}

func TestPaginationHelper_SetTotalCount_AdjustsCurrentPage(t *testing.T) {
	ph := NewPaginationHelper(10)
	ph.SetTotalCount(50)
	ph.GoToLastPage() // Page 5

	if ph.GetCurrentPage() != 5 {
		t.Fatalf("Expected to be on page 5")
	}

	// Reduce total count so current page is invalid
	ph.SetTotalCount(25) // Now only 3 pages

	// Should auto-adjust to last valid page (3)
	if ph.GetCurrentPage() != 3 {
		t.Errorf("Expected current page to adjust to 3, got %d", ph.GetCurrentPage())
	}
}

func TestPaginationHelper_SetTotalCount_ZeroResetsPage(t *testing.T) {
	ph := NewPaginationHelper(10)
	ph.SetTotalCount(50)
	ph.NextPage()
	ph.NextPage() // Page 3

	ph.SetTotalCount(0)

	if ph.GetCurrentPage() != 1 {
		t.Errorf("Expected page to reset to 1 when totalCount is 0, got %d", ph.GetCurrentPage())
	}
}

// Edge case tests
func TestPaginationHelper_EdgeCases(t *testing.T) {
	t.Run("Page size of 1", func(t *testing.T) {
		ph := NewPaginationHelper(1)
		ph.SetTotalCount(3)
		if ph.GetTotalPages() != 3 {
			t.Errorf("Expected 3 pages, got %d", ph.GetTotalPages())
		}
	})

	t.Run("Large page size", func(t *testing.T) {
		ph := NewPaginationHelper(1000)
		ph.SetTotalCount(50)
		if ph.GetTotalPages() != 1 {
			t.Errorf("Expected 1 page, got %d", ph.GetTotalPages())
		}
	})

	t.Run("Method chaining", func(t *testing.T) {
		ph := NewPaginationHelper(10)
		ph.SetTotalCount(50).NextPage().NextPage()
		if ph.GetCurrentPage() != 3 {
			t.Errorf("Expected page 3 after chaining, got %d", ph.GetCurrentPage())
		}
	})
}

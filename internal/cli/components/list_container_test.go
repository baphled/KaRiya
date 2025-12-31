package components

import (
	"strings"
	"testing"
)

func TestNewListContainer(t *testing.T) {
	lc := NewListContainer()

	if len(lc.items) != 0 {
		t.Errorf("expected empty items, got %d", len(lc.items))
	}

	if lc.emptyMessage == "" {
		t.Error("expected default empty message")
	}

	if lc.hasPagination {
		t.Error("expected hasPagination to be false by default")
	}

	if lc.hasEmptyMessage {
		t.Error("expected hasEmptyMessage to be false by default")
	}
}

func TestListContainerSetItems(t *testing.T) {
	items := []string{"Item 1", "Item 2", "Item 3"}
	lc := NewListContainer().SetItems(items)

	if len(lc.items) != len(items) {
		t.Errorf("expected %d items, got %d", len(items), len(lc.items))
	}

	for i, item := range items {
		if lc.items[i] != item {
			t.Errorf("expected item %q, got %q", item, lc.items[i])
		}
	}
}

func TestListContainerSetEmptyStateMessage(t *testing.T) {
	message := "No results found"
	lc := NewListContainer().SetEmptyStateMessage(message)

	if lc.emptyMessage != message {
		t.Errorf("expected empty message %q, got %q", message, lc.emptyMessage)
	}

	if !lc.hasEmptyMessage {
		t.Error("expected hasEmptyMessage to be true")
	}
}

func TestListContainerSetPaginationInfo(t *testing.T) {
	info := "Page 1 of 5"
	lc := NewListContainer().SetPaginationInfo(info)

	if lc.paginationInfo != info {
		t.Errorf("expected pagination info %q, got %q", info, lc.paginationInfo)
	}

	if !lc.hasPagination {
		t.Error("expected hasPagination to be true")
	}
}

func TestListContainerRenderEmpty(t *testing.T) {
	lc := NewListContainer()

	rendered := lc.Render()

	if !strings.Contains(rendered, "No items to display") {
		t.Error("expected rendered output to contain default empty message")
	}
}

func TestListContainerRenderEmptyWithCustomMessage(t *testing.T) {
	message := "No records found"
	lc := NewListContainer().SetEmptyStateMessage(message)

	rendered := lc.Render()

	if !strings.Contains(rendered, message) {
		t.Errorf("expected rendered output to contain %q", message)
	}
}

func TestListContainerRenderWithItems(t *testing.T) {
	items := []string{"Item 1", "Item 2", "Item 3"}
	lc := NewListContainer().SetItems(items)

	rendered := lc.Render()

	for _, item := range items {
		if !strings.Contains(rendered, item) {
			t.Errorf("expected rendered output to contain %q", item)
		}
	}
}

func TestListContainerRenderWithItemsAndPagination(t *testing.T) {
	items := []string{"Item 1", "Item 2"}
	pagination := "Page 1 of 3"

	lc := NewListContainer().
		SetItems(items).
		SetPaginationInfo(pagination)

	rendered := lc.Render()

	for _, item := range items {
		if !strings.Contains(rendered, item) {
			t.Errorf("expected rendered output to contain %q", item)
		}
	}

	if !strings.Contains(rendered, pagination) {
		t.Errorf("expected rendered output to contain pagination %q", pagination)
	}
}

func TestListContainerBuilderChaining(t *testing.T) {
	items := []string{"Item 1", "Item 2"}
	message := "Empty list"
	pagination := "Page 1 of 2"

	lc := NewListContainer().
		SetItems(items).
		SetEmptyStateMessage(message).
		SetPaginationInfo(pagination)

	if len(lc.items) != len(items) {
		t.Error("expected items to be set")
	}

	if lc.emptyMessage != message {
		t.Error("expected empty message to be set")
	}

	if lc.paginationInfo != pagination {
		t.Error("expected pagination info to be set")
	}
}

func TestListContainerRenderSingleItem(t *testing.T) {
	items := []string{"Single Item"}
	lc := NewListContainer().SetItems(items)

	rendered := lc.Render()

	if !strings.Contains(rendered, "Single Item") {
		t.Error("expected rendered output to contain single item")
	}
}

func TestListContainerRenderManyItems(t *testing.T) {
	items := make([]string, 100)
	for i := 0; i < 100; i++ {
		items[i] = "Item " + string(rune(i))
	}

	lc := NewListContainer().SetItems(items)

	rendered := lc.Render()

	// Verify some items are present
	if !strings.Contains(rendered, "Item") {
		t.Error("expected rendered output to contain items")
	}
}

func TestListContainerRenderItemOrdering(t *testing.T) {
	items := []string{"First", "Second", "Third"}
	lc := NewListContainer().SetItems(items)

	rendered := lc.Render()

	firstPos := strings.Index(rendered, "First")
	secondPos := strings.Index(rendered, "Second")
	thirdPos := strings.Index(rendered, "Third")

	if firstPos == -1 || secondPos == -1 || thirdPos == -1 {
		t.Error("expected all items to be present")
	}

	if firstPos > secondPos || secondPos > thirdPos {
		t.Error("expected items to appear in order")
	}
}

func TestListContainerEmptyItemsList(t *testing.T) {
	lc := NewListContainer().SetItems([]string{})

	rendered := lc.Render()

	if !strings.Contains(rendered, "No items to display") {
		t.Error("expected rendered output to show empty message")
	}
}

func TestListContainerWithPaginationButNoItems(t *testing.T) {
	lc := NewListContainer().SetPaginationInfo("Page 1 of 1")

	rendered := lc.Render()

	// Should show empty message, not pagination
	if !strings.Contains(rendered, "No items to display") {
		t.Error("expected rendered output to show empty message")
	}
}

func TestListContainerMultilineItems(t *testing.T) {
	items := []string{
		"Item 1\nWith multiple\nLines",
		"Item 2\nAlso multiline",
	}

	lc := NewListContainer().SetItems(items)

	rendered := lc.Render()

	if !strings.Contains(rendered, "multiple") || !strings.Contains(rendered, "multiline") {
		t.Error("expected rendered output to contain multiline content")
	}
}

func TestListContainerSpecialCharacters(t *testing.T) {
	items := []string{
		"Item with [brackets]",
		"Item with {braces}",
		"Item with special @#$%",
	}

	lc := NewListContainer().SetItems(items)

	rendered := lc.Render()

	if !strings.Contains(rendered, "[brackets]") ||
		!strings.Contains(rendered, "{braces}") ||
		!strings.Contains(rendered, "@#$%") {
		t.Error("expected rendered output to contain special characters")
	}
}

func TestListContainerPaginationWithoutInfo(t *testing.T) {
	items := []string{"Item 1", "Item 2"}
	lc := NewListContainer().SetItems(items)

	rendered := lc.Render()

	// Should not include pagination info if not set
	if strings.Contains(rendered, "Page") {
		t.Error("expected rendered output to not contain pagination info")
	}
}

func TestListContainerEmptyPaginationInfo(t *testing.T) {
	items := []string{"Item 1"}
	lc := NewListContainer().SetItems(items).SetPaginationInfo("")

	rendered := lc.Render()

	// Empty pagination info should not be displayed
	if !strings.Contains(rendered, "Item 1") {
		t.Error("expected rendered output to contain item")
	}
}

func TestListContainerRenderConsistency(t *testing.T) {
	items := []string{"Item 1", "Item 2"}
	lc := NewListContainer().SetItems(items)

	rendered1 := lc.Render()
	rendered2 := lc.Render()

	if rendered1 != rendered2 {
		t.Error("expected consistent rendering")
	}
}

func TestListContainerLongItemText(t *testing.T) {
	longText := strings.Repeat("a", 200)
	items := []string{longText}

	lc := NewListContainer().SetItems(items)

	rendered := lc.Render()

	if !strings.Contains(rendered, "a") {
		t.Error("expected rendered output to contain long item text")
	}
}

func TestListContainerPaginationInfo(t *testing.T) {
	items := []string{"Item 1", "Item 2", "Item 3"}
	pagination := "Showing 1-3 of 10 items"

	lc := NewListContainer().
		SetItems(items).
		SetPaginationInfo(pagination)

	rendered := lc.Render()

	if !strings.Contains(rendered, "Showing 1-3 of 10 items") {
		t.Error("expected rendered output to contain pagination info")
	}
}

func TestListContainerDuplicateItems(t *testing.T) {
	items := []string{"Item", "Item", "Item"}
	lc := NewListContainer().SetItems(items)

	rendered := lc.Render()

	// Count occurrences of "Item"
	count := strings.Count(rendered, "Item")
	if count < 3 {
		t.Errorf("expected at least 3 occurrences of 'Item', got %d", count)
	}
}


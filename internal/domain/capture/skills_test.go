package capture_test

import (
	"testing"

	"github.com/baphled/kariya/internal/domain/capture"
)

func TestFilterNewSkillSuggestions(t *testing.T) {
	t.Run("returns all when no existing names", func(t *testing.T) {
		suggestions := []capture.SkillSuggestion{
			{Name: "Go", EventIDs: []string{"e1"}},
			{Name: "Python", EventIDs: []string{"e2"}},
		}

		result := capture.FilterNewSkillSuggestions(suggestions, nil)

		if len(result) != 2 {
			t.Errorf("length = %d, want 2", len(result))
		}
	})

	t.Run("filters existing skills case-insensitively", func(t *testing.T) {
		suggestions := []capture.SkillSuggestion{
			{Name: "Go", EventIDs: []string{"e1"}},
			{Name: "Python", EventIDs: []string{"e2"}},
			{Name: "Rust", EventIDs: []string{"e3"}},
		}
		existing := []string{"go", "PYTHON"}

		result := capture.FilterNewSkillSuggestions(suggestions, existing)

		if len(result) != 1 {
			t.Fatalf("length = %d, want 1", len(result))
		}
		if result[0].Name != "Rust" {
			t.Errorf("Name = %q, want %q", result[0].Name, "Rust")
		}
	})

	t.Run("returns empty when all exist", func(t *testing.T) {
		suggestions := []capture.SkillSuggestion{
			{Name: "Go", EventIDs: []string{"e1"}},
		}
		existing := []string{"Go"}

		result := capture.FilterNewSkillSuggestions(suggestions, existing)

		if len(result) != 0 {
			t.Errorf("length = %d, want 0", len(result))
		}
	})

	t.Run("handles empty suggestions", func(t *testing.T) {
		result := capture.FilterNewSkillSuggestions(nil, []string{"Go"})

		if result == nil {
			t.Error("expected non-nil empty slice")
		}
		if len(result) != 0 {
			t.Errorf("length = %d, want 0", len(result))
		}
	})

	t.Run("preserves event IDs in filtered results", func(t *testing.T) {
		suggestions := []capture.SkillSuggestion{
			{Name: "Rust", EventIDs: []string{"e1", "e2", "e3"}},
		}

		result := capture.FilterNewSkillSuggestions(suggestions, []string{"Go"})

		if len(result) != 1 {
			t.Fatalf("length = %d, want 1", len(result))
		}
		if len(result[0].EventIDs) != 3 {
			t.Errorf("EventIDs length = %d, want 3", len(result[0].EventIDs))
		}
	})
}

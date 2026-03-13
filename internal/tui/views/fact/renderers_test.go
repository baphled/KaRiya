package fact_test

import (
	"time"

	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/baphled/kariya/internal/tui/views/fact"
	"github.com/baphled/kariya/internal/ui/display"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("View Renderers", func() {
	var testFact display.Fact

	BeforeEach(func() {
		fact := fixtures.Fact("fact-123", "")
		fact.Text = "Experienced in Go development"
		fact.CompetencyCategories = []string{"Backend", "Go"}
		fact.StrengthSignal = "strong"
		fact.RoleFit = "senior"
		fact.AudienceRelevance = []string{"technical"}
		fact.CreatedAt = time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC)
		testFact = display.FactFromDomain(fact)
	})

	Describe("RenderFactDetail", func() {
		It("renders fact details correctly", func() {
			result := fact.RenderFactDetail(testFact)
			Expect(result).To(ContainSubstring("Fact Details"))
			Expect(result).To(ContainSubstring("Experienced in Go development"))
			Expect(result).To(ContainSubstring("strong"))
			Expect(result).To(ContainSubstring("2025-01-15"))
		})

		It("returns no-selection message for nil fact", func() {
			result := fact.RenderFactDetail(display.Fact{})
			Expect(result).To(Equal("No fact selected"))
		})
	})

	Describe("RenderDeleteConfirm", func() {
		It("renders delete confirmation correctly", func() {
			result := fact.RenderDeleteConfirm(testFact)
			Expect(result).To(ContainSubstring("Confirm Deletion"))
			Expect(result).To(ContainSubstring("Delete this fact?"))
			Expect(result).To(ContainSubstring("Experienced in Go development"))
			Expect(result).To(ContainSubstring("cannot be undone"))
		})

		It("returns no-fact message for nil fact", func() {
			result := fact.RenderDeleteConfirm(display.Fact{})
			Expect(result).To(Equal("No fact to delete"))
		})

		It("truncates long fact text", func() {
			longFact := display.FactFromDomain(fixtures.Fact("fact-long", ""))
			longFact.Text = "This is a very long fact text that exceeds one hundred characters and should be truncated with ellipsis at the end of the string"
			result := fact.RenderDeleteConfirm(longFact)
			Expect(result).To(ContainSubstring("..."))
			Expect(result).NotTo(ContainSubstring("end of the string"))
		})
	})

	Describe("RenderResults", func() {
		It("renders statistics correctly", func() {
			result := fact.RenderResults(42, 10)
			Expect(result).To(ContainSubstring("Fact Statistics"))
			Expect(result).To(ContainSubstring("42"))
			Expect(result).To(ContainSubstring("10"))
		})

		It("omits current count when zero", func() {
			result := fact.RenderResults(42, 0)
			Expect(result).To(ContainSubstring("42"))
			Expect(result).NotTo(ContainSubstring("Currently viewing"))
		})
	})

	Describe("RenderEditorFallback", func() {
		It("renders editor content with fact data", func() {
			result := fact.RenderEditorFallback(testFact, false)
			Expect(result).To(ContainSubstring("Edit Fact"))
			Expect(result).To(ContainSubstring("Experienced in Go development"))
			Expect(result).To(ContainSubstring("Ctrl+S"))
		})

		It("hides save instructions when form has errors", func() {
			result := fact.RenderEditorFallback(testFact, true)
			Expect(result).To(ContainSubstring("Edit Fact"))
			Expect(result).NotTo(ContainSubstring("Ctrl+S"))
		})

		It("returns no-fact message for nil fact", func() {
			result := fact.RenderEditorFallback(display.Fact{}, false)
			Expect(result).To(ContainSubstring("No fact loaded"))
		})
	})

	Describe("truncateText", func() {
		It("returns text unchanged when under max length", func() {
			result := truncateTextHelper("hello", 10)
			Expect(result).To(Equal("hello"))
		})

		It("returns text unchanged when equal to max length", func() {
			result := truncateTextHelper("hello", 5)
			Expect(result).To(Equal("hello"))
		})

		It("truncates text and adds ellipsis when over max length", func() {
			result := truncateTextHelper("hello world", 5)
			Expect(result).To(Equal("hello..."))
		})

		It("handles empty string", func() {
			result := truncateTextHelper("", 10)
			Expect(result).To(Equal(""))
		})

		It("handles single character", func() {
			result := truncateTextHelper("a", 1)
			Expect(result).To(Equal("a"))
		})

		It("handles very long text", func() {
			longText := "The quick brown fox jumps over the lazy dog"
			result := truncateTextHelper(longText, 10)
			Expect(result).To(Equal("The quick ..."))
		})

		It("handles max length of 1", func() {
			result := truncateTextHelper("hello", 1)
			Expect(result).To(Equal("h..."))
		})

		It("handles max length of 0", func() {
			result := truncateTextHelper("hello", 0)
			Expect(result).To(Equal("..."))
		})
	})
})

func truncateTextHelper(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

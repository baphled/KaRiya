package facts_test

import (
	"time"

	"github.com/baphled/kariya/internal/cli/screens/facts"
	"github.com/baphled/kariya/internal/domain/career"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("View Renderers", func() {
	var testFact *career.Fact

	BeforeEach(func() {
		testFact = &career.Fact{
			ID:                   "fact-123",
			Text:                 "Experienced in Go development",
			CompetencyCategories: []string{"Backend", "Go"},
			StrengthSignal:       "strong",
			RoleFit:              "senior",
			AudienceRelevance:    []string{"technical"},
			CreatedAt:            time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC),
		}
	})

	Describe("RenderFactDetail", func() {
		It("renders fact details correctly", func() {
			result := facts.RenderFactDetail(testFact)
			Expect(result).To(ContainSubstring("Fact Details"))
			Expect(result).To(ContainSubstring("Experienced in Go development"))
			Expect(result).To(ContainSubstring("strong"))
			Expect(result).To(ContainSubstring("2025-01-15"))
		})

		It("returns no-selection message for nil fact", func() {
			result := facts.RenderFactDetail(nil)
			Expect(result).To(Equal("No fact selected"))
		})
	})

	Describe("RenderDeleteConfirm", func() {
		It("renders delete confirmation correctly", func() {
			result := facts.RenderDeleteConfirm(testFact)
			Expect(result).To(ContainSubstring("Confirm Deletion"))
			Expect(result).To(ContainSubstring("Delete this fact?"))
			Expect(result).To(ContainSubstring("Experienced in Go development"))
			Expect(result).To(ContainSubstring("cannot be undone"))
		})

		It("returns no-fact message for nil fact", func() {
			result := facts.RenderDeleteConfirm(nil)
			Expect(result).To(Equal("No fact to delete"))
		})
	})

	Describe("RenderResults", func() {
		It("renders statistics correctly", func() {
			result := facts.RenderResults(42, 10)
			Expect(result).To(ContainSubstring("Fact Statistics"))
			Expect(result).To(ContainSubstring("42"))
			Expect(result).To(ContainSubstring("10"))
		})

		It("omits current count when zero", func() {
			result := facts.RenderResults(42, 0)
			Expect(result).To(ContainSubstring("42"))
			Expect(result).NotTo(ContainSubstring("Currently viewing"))
		})
	})

	Describe("RenderEditorFallback", func() {
		It("renders editor content with fact data", func() {
			result := facts.RenderEditorFallback(testFact, false)
			Expect(result).To(ContainSubstring("Edit Fact"))
			Expect(result).To(ContainSubstring("Experienced in Go development"))
			Expect(result).To(ContainSubstring("Ctrl+S"))
		})

		It("hides save instructions when form has errors", func() {
			result := facts.RenderEditorFallback(testFact, true)
			Expect(result).To(ContainSubstring("Edit Fact"))
			Expect(result).NotTo(ContainSubstring("Ctrl+S"))
		})

		It("returns no-fact message for nil fact", func() {
			result := facts.RenderEditorFallback(nil, false)
			Expect(result).To(ContainSubstring("No fact loaded"))
		})
	})
})

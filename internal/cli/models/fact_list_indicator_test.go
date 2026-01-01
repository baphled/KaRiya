package models

import (
	"context"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("FactListModel Indicator Navigation", func() {
	var (
		model *FactListModel
		repo  *careerrepo.MemoryRepository
		svc   *careerservice.Service
		ctx   context.Context
		facts []*career.Fact
	)

	ginkgo.BeforeEach(func() {
		repo = careerrepo.NewMemoryRepository()
		svc = careerservice.NewService(repo)
		ctx = context.Background()
		model = NewFactListModel(svc, ctx)
		model.width = 100
		model.height = 20

		// Create facts with timestamps that don't change during the test
		now := time.Now()
		facts = []*career.Fact{
			{ID: "fact-1", Text: "First fact", CompetencyCategories: []string{"technical"}, RoleFit: career.RoleFitPrincipal, AudienceRelevance: []string{"hiring_manager"}, SourceEventID: "event-1", CreatedAt: now.Add(-10 * time.Hour), UpdatedAt: now.Add(-10 * time.Hour)},
			{ID: "fact-2", Text: "Second fact", CompetencyCategories: []string{"mentoring"}, RoleFit: career.RoleFitStaff, AudienceRelevance: []string{"peer"}, SourceEventID: "event-2", CreatedAt: now.Add(-5 * time.Hour), UpdatedAt: now.Add(-5 * time.Hour)},
			{ID: "fact-3", Text: "Third fact", CompetencyCategories: []string{"product"}, RoleFit: career.RoleFitEM, AudienceRelevance: []string{"hiring_manager"}, SourceBurstID: "burst-1", CreatedAt: now, UpdatedAt: now},
			{ID: "fact-4", Text: "Fourth fact", CompetencyCategories: []string{"leadership"}, RoleFit: career.RoleFitPrincipal, AudienceRelevance: []string{"recruiter"}, SourceEventID: "event-4", CreatedAt: now.Add(1 * time.Hour), UpdatedAt: now.Add(1 * time.Hour)},
			{ID: "fact-5", Text: "Fifth fact", CompetencyCategories: []string{"communication"}, RoleFit: career.RoleFitStaff, AudienceRelevance: []string{"peer"}, SourceEventID: "event-5", CreatedAt: now.Add(2 * time.Hour), UpdatedAt: now.Add(2 * time.Hour)},
		}
	})

	ginkgo.Describe("Indicator Position Tracking", func() {
		ginkgo.It("should show indicator on the first table row at cursor position 0", func() {
			model.SetFacts(facts)

			// Check table rows directly
			rows := model.table.Rows()
			gomega.Expect(len(rows)).To(gomega.BeNumerically(">", 0), "Should have rows")

			// The first row (index 0) should have the indicator
			firstRowContent := rows[0][0]
			gomega.Expect(firstRowContent).To(gomega.HavePrefix("▶ "), "First row should have indicator at cursor position")
		})

		ginkgo.It("should move indicator to specific position", func() {
			model.SetFacts(facts)

			// Move to third item using container method
			model.listContainer.SetSelectedIdx(2)
			model.updateTableRows()

			rows := model.table.Rows()
			gomega.Expect(len(rows)).To(gomega.BeNumerically(">", 2), "Should have at least 3 rows")

			// Third row should have indicator
			thirdRowContent := rows[2][0]
			gomega.Expect(thirdRowContent).To(gomega.HavePrefix("▶ "), "Third row should have indicator")
		})

		ginkgo.It("should not skip items when navigating", func() {
			model.SetFacts(facts)

			// Navigate through all items using container methods
			for i := 0; i < len(facts); i++ {
				model.listContainer.SetSelectedIdx(i)
				model.updateTableRows()

				rows := model.table.Rows()
				gomega.Expect(len(rows)).To(gomega.Equal(len(facts)), "Should have all facts as rows")

				// Verify the indicator is at the correct position
				rowContent := rows[i][0]
				gomega.Expect(rowContent).To(gomega.HavePrefix("▶ "), "Row %d should have indicator at cursor position", i)
			}
		})

		ginkgo.It("should reset cursor when filtering", func() {
			model.SetFacts(facts)

			// Move to second item using container method
			model.listContainer.SetSelectedIdx(1)

			// Apply filter
			model.SetCompetencyFilter("leadership")

			// Should reset to first item of filtered list
			gomega.Expect(model.listContainer.GetSelectedIdx()).To(gomega.Equal(0), "Cursor should reset to 0 after filtering")
		})

		ginkgo.It("should display exactly one indicator per update", func() {
			model.SetFacts(facts)

			// Navigate through items and verify only one indicator
			for i := 0; i < len(facts); i++ {
				model.listContainer.SetSelectedIdx(i)
				model.updateTableRows()

				rows := model.table.Rows()

				// Count indicators
				indicatorCount := 0
				for _, row := range rows {
					if strings.HasPrefix(row[0], "▶ ") {
						indicatorCount++
					}
				}

				gomega.Expect(indicatorCount).To(gomega.Equal(1), "Should have exactly one indicator at position %d", i)
			}
		})
	})

	ginkgo.Describe("Table Row Rendering", func() {
		ginkgo.It("should render all facts in the table", func() {
			model.SetFacts(facts)

			// Check that all facts are in the table rows
			gomega.Expect(len(model.table.Rows())).To(gomega.Equal(len(facts)))
		})

		ginkgo.It("should have correct number of rows after filtering", func() {
			model.SetFacts(facts)
			model.SetCompetencyFilter("leadership")

			// Should only have facts with "leadership" competency
			gomega.Expect(len(model.table.Rows())).To(gomega.Equal(1))
		})

		ginkgo.It("should not have gaps in row rendering", func() {
			model.SetFacts(facts)

			// Verify each row contains text
			rows := model.table.Rows()
			gomega.Expect(len(rows)).To(gomega.Equal(len(facts)))

			// Check that rows are not empty and contain fact content
			for i, row := range rows {
				rowContent := row[0] // First column contains the text with indicator
				gomega.Expect(rowContent).NotTo(gomega.BeEmpty(), "Row %d should not be empty", i)
				// Verify it contains either the indicator or spaces
				gomega.Expect(len(rowContent) > 2).To(gomega.BeTrue(), "Row %d should contain fact text", i)
			}
		})

		ginkgo.It("should maintain indicator prefix format", func() {
			model.SetFacts(facts)

			rows := model.table.Rows()
			for i, row := range rows {
				rowContent := row[0]
				if i == model.listContainer.GetSelectedIdx() {
					gomega.Expect(rowContent).To(gomega.HavePrefix("▶ "), "Row %d should have indicator prefix", i)
				} else {
					gomega.Expect(rowContent).To(gomega.HavePrefix("  "), "Row %d should have space prefix", i)
				}
			}
		})
	})
})

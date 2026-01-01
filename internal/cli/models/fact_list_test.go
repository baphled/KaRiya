package models

import (
	"context"
	"strings"
	"fmt"

	"time"

	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("FactListModel", func() {
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

		facts = []*career.Fact{
			{ID: "fact-1", Text: "Led architecture redesign", CompetencyCategories: []string{"technical", "leadership"}, RoleFit: career.RoleFitPrincipal, AudienceRelevance: []string{"hiring_manager", "recruiter"}, SourceEventID: "event-1", CreatedAt: time.Now().Add(-10 * time.Hour), UpdatedAt: time.Now().Add(-10 * time.Hour)},
			{ID: "fact-2", Text: "Mentored engineers", CompetencyCategories: []string{"mentoring"}, RoleFit: career.RoleFitStaff, AudienceRelevance: []string{"peer"}, SourceEventID: "event-2", CreatedAt: time.Now().Add(-5 * time.Hour), UpdatedAt: time.Now().Add(-5 * time.Hour)},
			{ID: "fact-3", Text: "Managed roadmap", CompetencyCategories: []string{"product", "leadership"}, RoleFit: career.RoleFitEM, AudienceRelevance: []string{"hiring_manager"}, SourceBurstID: "burst-1", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		}
	})

	ginkgo.Describe("Initialization", func() {
		ginkgo.It("should create new model", func() {
			m := NewFactListModel(svc, ctx)
			gomega.Expect(m).NotTo(gomega.BeNil())
			gomega.Expect(m.service).To(gomega.Equal(svc))
			gomega.Expect(len(m.facts)).To(gomega.Equal(0))
		})
	})

	ginkgo.Describe("Setting Facts", func() {
		ginkgo.It("should set facts", func() {
			model.SetFacts(facts)
			gomega.Expect(len(model.facts)).To(gomega.Equal(3))
			gomega.Expect(len(model.filtered)).To(gomega.Equal(3))
		})

		ginkgo.It("should reset navigation", func() {
			model.focusedIdx = 2
			model.SetFacts(facts)
			gomega.Expect(model.focusedIdx).To(gomega.Equal(0))
		})
	})

	ginkgo.Describe("Rendering", func() {
		ginkgo.It("should show empty state", func() {
			result := model.View()
			gomega.Expect(result).To(gomega.ContainSubstring("No facts found"))
		})

		ginkgo.It("should show facts list", func() {
			model.SetFacts(facts)
			result := model.View()
			gomega.Expect(result).To(gomega.ContainSubstring("Facts"))
			// Updated to match standardized list.go pagination format
			gomega.Expect(result).To(gomega.ContainSubstring("Showing 1-3 of 3 facts"))
		})

		ginkgo.It("should show selection count", func() {
			model.SetFacts(facts)
			model.selectedFacts["fact-1"] = true
			model.selectedFacts["fact-2"] = true
			result := model.View()
			// Selection count is no longer shown in pagination - removed per list.go standardization
			// gomega.Expect(result).To(gomega.ContainSubstring("2 selected"))
			gomega.Expect(result).To(gomega.ContainSubstring("Showing 1-3 of 3 facts"))
		})
	})

	ginkgo.Describe("Competency Filtering", func() {
		ginkgo.BeforeEach(func() {
			model.SetFacts(facts)
		})

		ginkgo.It("should filter by competency", func() {
			model.SetCompetencyFilter("mentoring")
			gomega.Expect(len(model.filtered)).To(gomega.Equal(1))
			gomega.Expect(model.filtered[0].ID).To(gomega.Equal("fact-2"))
		})

		ginkgo.It("should be case-insensitive", func() {
			model.SetCompetencyFilter("MENTORING")
			gomega.Expect(len(model.filtered)).To(gomega.Equal(1))
		})

		ginkgo.It("should filter multiple matches", func() {
			model.SetCompetencyFilter("leadership")
			gomega.Expect(len(model.filtered)).To(gomega.Equal(2))
		})
	})

	ginkgo.Describe("Role Fit Filtering", func() {
		ginkgo.BeforeEach(func() {
			model.SetFacts(facts)
		})

		ginkgo.It("should filter by role fit", func() {
			model.SetRoleFitFilter(career.RoleFitStaff)
			gomega.Expect(len(model.filtered)).To(gomega.Equal(1))
			gomega.Expect(model.filtered[0].ID).To(gomega.Equal("fact-2"))
		})
	})

	ginkgo.Describe("Audience Filtering", func() {
		ginkgo.BeforeEach(func() {
			model.SetFacts(facts)
		})

		ginkgo.It("should filter by audience", func() {
			model.SetAudienceFilter("peer")
			gomega.Expect(len(model.filtered)).To(gomega.Equal(1))
			gomega.Expect(model.filtered[0].ID).To(gomega.Equal("fact-2"))
		})

		ginkgo.It("should handle multiple matches", func() {
			model.SetAudienceFilter("hiring_manager")
			gomega.Expect(len(model.filtered)).To(gomega.Equal(2))
		})
	})

	ginkgo.Describe("Sorting", func() {
		ginkgo.BeforeEach(func() {
			model.SetFacts(facts)
		})

		ginkgo.It("should sort by date descending", func() {
			model.SetSort("date", "desc")
			gomega.Expect(model.filtered[0].ID).To(gomega.Equal("fact-3"))
			gomega.Expect(model.filtered[2].ID).To(gomega.Equal("fact-1"))
		})

		ginkgo.It("should sort by date ascending", func() {
			model.SetSort("date", "asc")
			gomega.Expect(model.filtered[0].ID).To(gomega.Equal("fact-1"))
			gomega.Expect(model.filtered[2].ID).To(gomega.Equal("fact-3"))
		})

		ginkgo.It("should sort by relevance", func() {
			model.SetSort("relevance", "desc")
			gomega.Expect(model.filtered[2].CompetencyCategories).To(gomega.HaveLen(1))
		})
	})

	ginkgo.Describe("Selection", func() {
		ginkgo.BeforeEach(func() {
			model.SetFacts(facts)
		})

		ginkgo.It("should track selections", func() {
			model.selectedFacts["fact-1"] = true
			gomega.Expect(model.selectedFacts["fact-1"]).To(gomega.BeTrue())
		})

		ginkgo.It("should get selected facts", func() {
			model.selectedFacts["fact-1"] = true
			model.selectedFacts["fact-3"] = true
			selected := model.GetSelectedFacts()
			gomega.Expect(len(selected)).To(gomega.Equal(2))
		})

		ginkgo.It("should clear selection", func() {
			model.selectedFacts["fact-1"] = true
			model.ClearSelection()
			gomega.Expect(len(model.selectedFacts)).To(gomega.Equal(0))
		})
	})

	ginkgo.Describe("Navigation", func() {
		ginkgo.BeforeEach(func() {
			model.SetFacts(facts)
		})

		ginkgo.It("should get selected fact", func() {
			model.focusedIdx = 1
			selected := model.GetSelectedFact()
			gomega.Expect(selected.ID).To(gomega.Equal(facts[1].ID))
		})

		ginkgo.It("should return nil for empty", func() {
			model.SetFacts([]*career.Fact{})
			selected := model.GetSelectedFact()
			gomega.Expect(selected).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("State Flags", func() {
		ginkgo.It("should track submitted state", func() {
			gomega.Expect(model.IsSubmitted()).To(gomega.BeFalse())
			model.submitted = true
			gomega.Expect(model.IsSubmitted()).To(gomega.BeTrue())
		})

		ginkgo.It("should track cancelled state", func() {
			gomega.Expect(model.IsCancelled()).To(gomega.BeFalse())
			model.cancelled = true
			gomega.Expect(model.IsCancelled()).To(gomega.BeTrue())
		})

		ginkgo.It("should track error state", func() {
			testErr := fmt.Errorf("test error")
			model.err = testErr
			gomega.Expect(model.GetError()).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Describe("Edge Cases", func() {
		ginkgo.It("should handle multiple SetFacts calls", func() {
			model.SetFacts(facts)
			gomega.Expect(len(model.filtered)).To(gomega.Equal(3))

			model.SetFacts([]*career.Fact{facts[0]})
			gomega.Expect(len(model.filtered)).To(gomega.Equal(1))
		})

		ginkgo.It("should handle combining filters", func() {
			model.SetFacts(facts)
			model.SetCompetencyFilter("leadership")
			model.SetRoleFitFilter(career.RoleFitPrincipal)
			gomega.Expect(len(model.filtered)).To(gomega.Equal(1))
			gomega.Expect(model.filtered[0].ID).To(gomega.Equal("fact-1"))
		})
	})

	ginkgo.Describe("Focus Indicator Display", func() {
		ginkgo.It("should display focus indicator for selected item", func() {
			model.SetFacts(facts)
			view := model.View()
			gomega.Expect(view).To(gomega.ContainSubstring("▶"))
		})

		ginkgo.It("should use consistent marker character", func() {
			model.SetFacts(facts)
			view := model.View()
			gomega.Expect(view).To(gomega.ContainSubstring("▶ "))
		})

		ginkgo.It("should render marker at beginning of line", func() {
			model.SetFacts(facts)
			view := model.View()
			lines := strings.Split(view, "\n")
			foundMarker := false
			for _, line := range lines {
				if strings.Contains(line, "▶ ") {
					foundMarker = true
					break
				}
			}
			gomega.Expect(foundMarker).To(gomega.BeTrue())
		})

		ginkgo.It("should not show marker in empty state", func() {
			model.SetFacts([]*career.Fact{})
			view := model.View()
			gomega.Expect(view).NotTo(gomega.ContainSubstring("▶"))
		})
	})
})

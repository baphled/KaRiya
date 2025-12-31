package models

import (
	"context"
	"testing"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/google/uuid"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	tea "github.com/charmbracelet/bubbletea"
)

func TestBurstListModel(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Burst List Model Suite")
}

var _ = ginkgo.Describe("BurstListModel", func() {
	var (
		model      *BurstListModel
		service    *careerservice.Service
		repo       *careerrepo.MemoryRepository
		ctx        context.Context
		testBursts []*career.Burst
	)

	ginkgo.BeforeEach(func() {
		ctx = context.Background()
		repo = careerrepo.NewMemoryRepository()
		service = careerservice.NewService(repo)

		testBursts = []*career.Burst{
			{
				ID:              uuid.New().String(),
				Name:            "Platform Migration",
				EventIDs:        []string{uuid.New().String(), uuid.New().String()},
				CompetencyFocus: "Technical",
				CreatedAt:       time.Now().Add(-30 * 24 * time.Hour),
				UpdatedAt:       time.Now().Add(-30 * 24 * time.Hour),
			},
			{
				ID:              uuid.New().String(),
				Name:            "Team Leadership",
				EventIDs:        []string{uuid.New().String(), uuid.New().String(), uuid.New().String()},
				CompetencyFocus: "Leadership",
				CreatedAt:       time.Now().Add(-20 * 24 * time.Hour),
				UpdatedAt:       time.Now().Add(-20 * 24 * time.Hour),
			},
			{
				ID:              uuid.New().String(),
				Name:            "Architecture Design",
				EventIDs:        []string{uuid.New().String()},
				CompetencyFocus: "Technical",
				CreatedAt:       time.Now().Add(-10 * 24 * time.Hour),
				UpdatedAt:       time.Now().Add(-10 * 24 * time.Hour),
			},
		}

		model = NewBurstListModel(service, ctx)
	})

	ginkgo.It("should initialize with selectedIdx=0", func() {
		gomega.Expect(model.selectedIdx).To(gomega.Equal(0))
	})

	ginkgo.It("should move down on KeyDown", func() {
		model.bursts = testBursts
		initial := model.selectedIdx
		model.Update(tea.KeyMsg{Type: tea.KeyDown})
		gomega.Expect(model.selectedIdx).To(gomega.Equal(initial + 1))
	})

	ginkgo.It("should move up on KeyUp", func() {
		model.bursts = testBursts
		model.selectedIdx = 2
		model.Update(tea.KeyMsg{Type: tea.KeyUp})
		gomega.Expect(model.selectedIdx).To(gomega.Equal(1))
	})

	ginkgo.It("should wrap to end on KeyUp at start", func() {
		model.bursts = testBursts
		model.selectedIdx = 0
		model.Update(tea.KeyMsg{Type: tea.KeyUp})
		gomega.Expect(model.selectedIdx).To(gomega.Equal(len(testBursts) - 1))
	})

	ginkgo.It("should wrap to start on KeyDown at end", func() {
		model.bursts = testBursts
		model.selectedIdx = len(testBursts) - 1
		model.Update(tea.KeyMsg{Type: tea.KeyDown})
		gomega.Expect(model.selectedIdx).To(gomega.Equal(0))
	})

	ginkgo.It("should toggle expanded on Space", func() {
		model.bursts = testBursts
		initial := model.expandedIndices[0]
		model.Update(tea.KeyMsg{Type: tea.KeySpace})
		gomega.Expect(model.expandedIndices[0]).To(gomega.Equal(!initial))
	})

	ginkgo.It("should render burst name", func() {
		model.bursts = testBursts
		output := model.View()
		gomega.Expect(output).To(gomega.ContainSubstring("Platform Migration"))
	})

	ginkgo.It("should render event count", func() {
		model.bursts = testBursts
		output := model.View()
		gomega.Expect(output).To(gomega.ContainSubstring("2"))
	})

	ginkgo.It("should render competency", func() {
		model.bursts = testBursts
		output := model.View()
		gomega.Expect(output).To(gomega.ContainSubstring("Technical"))
	})

	ginkgo.It("should filter by competency", func() {
		model.bursts = testBursts
		model.filterBy = "Leadership"
		output := model.View()
		gomega.Expect(output).To(gomega.ContainSubstring("Team Leadership"))
	})

	ginkgo.It("should sort by event count", func() {
		model.bursts = testBursts
		model.sortBy = "event_count"
		displayed := model.getDisplayedBursts()
		first := testBursts[displayed[0]]
		gomega.Expect(len(first.EventIDs)).To(gomega.Equal(1))
	})

	ginkgo.It("should sort by name", func() {
		model.bursts = testBursts
		model.sortBy = "name"
		displayed := model.getDisplayedBursts()
		for i := 0; i < len(displayed)-1; i++ {
			current := testBursts[displayed[i]].Name
			next := testBursts[displayed[i+1]].Name
			gomega.Expect(current <= next).To(gomega.BeTrue())
		}
	})

	ginkgo.It("should sort by date", func() {
		model.bursts = testBursts
		model.sortBy = "date"
		displayed := model.getDisplayedBursts()
		for i := 0; i < len(displayed)-1; i++ {
			current := testBursts[displayed[i]].CreatedAt
			next := testBursts[displayed[i+1]].CreatedAt
			gomega.Expect(current.Before(next) || current.Equal(next)).To(gomega.BeTrue())
		}
	})

	ginkgo.It("should show events when expanded", func() {
		model.bursts = testBursts
		model.expandedIndices[0] = true
		output := model.View()
		gomega.Expect(output).To(gomega.ContainSubstring("Events:"))
	})

	ginkgo.It("should show empty message", func() {
		model.bursts = []*career.Burst{}
		output := model.View()
		gomega.Expect(output).To(gomega.ContainSubstring("No bursts found"))
	})

	ginkgo.It("should handle single burst", func() {
		model.bursts = []*career.Burst{testBursts[0]}
		output := model.View()
		gomega.Expect(output).To(gomega.ContainSubstring("Platform Migration"))
	})

	ginkgo.It("should handle large list", func() {
		largeBursts := make([]*career.Burst, 100)
		for i := 0; i < 100; i++ {
			largeBursts[i] = &career.Burst{
				ID:              uuid.New().String(),
				Name:            "Burst",
				EventIDs:        []string{uuid.New().String()},
				CompetencyFocus: "Technical",
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			}
		}
		model.bursts = largeBursts
		output := model.View()
		gomega.Expect(len(output) > 0).To(gomega.BeTrue())
	})

	ginkgo.It("should show no matching message", func() {
		model.filterBy = "NonExistent"
		output := model.View()
		gomega.Expect(output).To(gomega.ContainSubstring("No matching bursts"))
	})

	ginkgo.It("should update size on resize", func() {
		msg := tea.WindowSizeMsg{Width: 200, Height: 40}
		model.Update(msg)
		gomega.Expect(model.width).To(gomega.Equal(200))
		gomega.Expect(model.height).To(gomega.Equal(40))
	})

	ginkgo.It("should handle long name", func() {
		model.bursts = []*career.Burst{
			{
				ID:              uuid.New().String(),
				Name:            "This is a very long burst name that should be truncated gracefully",
				EventIDs:        []string{uuid.New().String()},
				CompetencyFocus: "Technical",
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			},
		}
		model.width = 80
		output := model.View()
		gomega.Expect(len(output) > 0).To(gomega.BeTrue())
	})

	ginkgo.It("should handle 2 events", func() {
		model.bursts = []*career.Burst{testBursts[0]}
		output := model.View()
		gomega.Expect(output).To(gomega.ContainSubstring("2"))
	})

	ginkgo.It("should handle 50+ events", func() {
		burst := &career.Burst{
			ID:              uuid.New().String(),
			Name:            "Large Burst",
			EventIDs:        make([]string, 50),
			CompetencyFocus: "Technical",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		for i := 0; i < 50; i++ {
			burst.EventIDs[i] = uuid.New().String()
		}
		model.bursts = []*career.Burst{burst}
		output := model.View()
		gomega.Expect(output).To(gomega.ContainSubstring("50"))
	})
})

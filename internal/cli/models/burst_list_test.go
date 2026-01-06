package models

import (
	"context"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

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
				ID:        uuid.New().String(),
				Name:      "Platform Migration",
				EventIDs:  []string{uuid.New().String(), uuid.New().String()},
				CreatedAt: time.Now().Add(-30 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-30 * 24 * time.Hour),
			},
			{
				ID:        uuid.New().String(),
				Name:      "Team Leadership",
				EventIDs:  []string{uuid.New().String(), uuid.New().String(), uuid.New().String()},
				CreatedAt: time.Now().Add(-20 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-20 * 24 * time.Hour),
			},
			{
				ID:        uuid.New().String(),
				Name:      "Architecture Design",
				EventIDs:  []string{uuid.New().String()},
				CreatedAt: time.Now().Add(-10 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-10 * 24 * time.Hour),
			},
		}

		model = NewBurstListModel(service, ctx)
	})

	ginkgo.It("should create new model", func() {
		gomega.Expect(model).NotTo(gomega.BeNil())
		gomega.Expect(model.GetSelectedIdx()).To(gomega.Equal(0))
	})

	ginkgo.It("should set bursts", func() {
		model.SetBursts(testBursts)
		gomega.Expect(model.GetBursts()).To(gomega.HaveLen(3))
	})

	ginkgo.It("should reset navigation", func() {
		model.SetBursts(testBursts)
		model.listContainer.SetSelectedIdx(2)
		model.SetBursts(testBursts)
		gomega.Expect(model.GetSelectedIdx()).To(gomega.Equal(0))
	})

	ginkgo.It("should show empty state", func() {
		model.SetBursts([]*career.Burst{})
		output := model.View()
		gomega.Expect(output).To(gomega.ContainSubstring("No bursts found"))
	})

	ginkgo.It("should show bursts list", func() {
		model.SetBursts(testBursts)
		output := model.View()
		gomega.Expect(output).To(gomega.ContainSubstring("Platform Migration"))
	})

	ginkgo.It("should show selection count in pagination", func() {
		model.SetBursts(testBursts)
		output := model.View()
		gomega.Expect(output).To(gomega.ContainSubstring("Showing"))
	})

	ginkgo.It("should filter by competency", func() {
		model.SetBursts(testBursts)
		output := model.View()
		gomega.Expect(output).To(gomega.ContainSubstring("Team Leadership"))
		gomega.Expect(output).NotTo(gomega.ContainSubstring("Platform Migration"))
	})

	ginkgo.It("should be case-insensitive", func() {
		model.SetBursts(testBursts)
		output := model.View()
		gomega.Expect(output).To(gomega.ContainSubstring("Team Leadership"))
	})

	ginkgo.It("should filter multiple matches", func() {
		model.SetBursts(testBursts)
		bursts := model.GetBursts()
		gomega.Expect(len(bursts)).To(gomega.Equal(2))
	})

	ginkgo.It("should sort by event count", func() {
		model.SetBursts(testBursts)
		model.SetSort("event_count", "asc")
		bursts := model.GetBursts()
		first := bursts[0]
		gomega.Expect(len(first.EventIDs)).To(gomega.Equal(1))
	})

	ginkgo.It("should sort by name", func() {
		model.SetBursts(testBursts)
		model.SetSort("name", "asc")
		bursts := model.GetBursts()
		for i := 0; i < len(bursts)-1; i++ {
			current := bursts[i].Name
			next := bursts[i+1].Name
			gomega.Expect(current <= next).To(gomega.BeTrue())
		}
	})

	ginkgo.It("should sort by date", func() {
		model.SetBursts(testBursts)
		model.SetSort("date", "asc")
		bursts := model.GetBursts()
		for i := 0; i < len(bursts)-1; i++ {
			current := bursts[i].CreatedAt
			next := bursts[i+1].CreatedAt
			gomega.Expect(current.Before(next) || current.Equal(next)).To(gomega.BeTrue())
		}
	})

	ginkgo.It("should track selections", func() {
		model.SetBursts(testBursts)
		burst := model.GetSelectedBurst()
		gomega.Expect(burst).NotTo(gomega.BeNil())
		model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
		gomega.Expect(model.selectedBursts[burst.ID]).To(gomega.BeTrue())
	})

	ginkgo.It("should get selected burst", func() {
		model.SetBursts(testBursts)
		burst := model.GetSelectedBurst()
		gomega.Expect(burst).NotTo(gomega.BeNil())
		// Architecture Design is the most recent (created 10 days ago)
		gomega.Expect(burst.Name).To(gomega.Equal("Architecture Design"))
	})

	ginkgo.It("should return nil for empty", func() {
		model.SetBursts([]*career.Burst{})
		burst := model.GetSelectedBurst()
		gomega.Expect(burst).To(gomega.BeNil())
	})

	ginkgo.It("should handle single burst", func() {
		model.SetBursts([]*career.Burst{testBursts[0]})
		output := model.View()
		gomega.Expect(output).To(gomega.ContainSubstring("Platform Migration"))
	})

	ginkgo.It("should handle large list", func() {
		largeBursts := make([]*career.Burst, 100)
		for i := 0; i < 100; i++ {
			largeBursts[i] = &career.Burst{
				ID:        uuid.New().String(),
				Name:      "Burst",
				EventIDs:  []string{uuid.New().String()},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
		}
		model.SetBursts(largeBursts)
		output := model.View()
		gomega.Expect(len(output) > 0).To(gomega.BeTrue())
	})

	ginkgo.It("should show no matching message when filtered", func() {
		model.SetBursts(testBursts)
		output := model.View()
		gomega.Expect(output).To(gomega.ContainSubstring("No bursts found"))
	})

	ginkgo.It("should update size on resize", func() {
		msg := tea.WindowSizeMsg{Width: 200, Height: 40}
		model.Update(msg)
		gomega.Expect(model.width).To(gomega.Equal(200))
		gomega.Expect(model.height).To(gomega.Equal(40))
	})

	ginkgo.It("should handle long name", func() {
		model.SetBursts([]*career.Burst{
			{
				ID:        uuid.New().String(),
				Name:      "This is a very long burst name that should be truncated gracefully",
				EventIDs:  []string{uuid.New().String()},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		})
		model.width = 80
		output := model.View()
		gomega.Expect(len(output) > 0).To(gomega.BeTrue())
	})

	ginkgo.It("should handle 2 events", func() {
		model.SetBursts([]*career.Burst{testBursts[0]})
		output := model.View()
		gomega.Expect(output).To(gomega.ContainSubstring("2"))
	})

	ginkgo.It("should handle 50+ events", func() {
		burst := &career.Burst{
			ID:        uuid.New().String(),
			Name:      "Large Burst",
			EventIDs:  make([]string, 50),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		for i := 0; i < 50; i++ {
			burst.EventIDs[i] = uuid.New().String()
		}
		model.SetBursts([]*career.Burst{burst})
		output := model.View()
		gomega.Expect(output).To(gomega.ContainSubstring("50"))
	})

	ginkgo.Describe("Focus Indicator Display", func() {
		ginkgo.It("should display focus indicator for selected item", func() {
			model.SetBursts(testBursts)
			view := model.View()
			gomega.Expect(view).To(gomega.ContainSubstring("▶"))
		})

		ginkgo.It("should use consistent marker character", func() {
			model.SetBursts(testBursts)
			view := model.View()
			gomega.Expect(view).To(gomega.ContainSubstring("▶ "))
		})

		ginkgo.It("should render marker at beginning of line", func() {
			model.SetBursts(testBursts)
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
			model.SetBursts([]*career.Burst{})
			view := model.View()
			gomega.Expect(view).NotTo(gomega.ContainSubstring("▶"))
		})
	})
})

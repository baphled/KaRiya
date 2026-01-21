package app_test

import (
	"context"
	"time"

	"github.com/baphled/kariya/internal/cli/app"
	"github.com/baphled/kariya/internal/cli/service"
	career "github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("List Container Navigation Integration - From Main Menu", func() {
	var (
		model      *app.Model
		repo       *careerrepo.MemoryRepository
		svc        *careerservice.Service
		cliService *service.CLIEventService
	)

	BeforeEach(func() {
		repo = careerrepo.NewMemoryRepository()
		burstRepo := careerrepo.NewMemoryBurstRepository()
		factRepo := careerrepo.NewMemoryFactRepository()
		svc = careerservice.NewService(repo)
		svc.SetBurstRepository(burstRepo)
		svc.SetFactRepository(factRepo)
		cliService = service.NewCLIEventService(svc)

		// Pre-populate repository with test events using fixtures
		event1 := fixtures.EventWith("event1", "First event - Learned Go", "Company A", "")
		event1.Date = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		event1.Tags = []string{"technical", "project"}
		event1.Categories = []string{"technical"}

		event2 := fixtures.EventWith("event2", "Second event - Led team meeting", "Company B", "")
		event2.Date = time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC)
		event2.Tags = []string{"leadership"}
		event2.Categories = []string{"leadership"}

		event3 := fixtures.EventWith("event3", "Third event - Deployed to production", "Company C", "")
		event3.Date = time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC)
		event3.Tags = []string{"product", "technical"}
		event3.Categories = []string{"technical"}

		for _, event := range []*career.CareerEvent{event1, event2, event3} {
			_ = repo.Create(context.Background(), event)
		}

		// Pre-populate burst/fact repositories using fixtures
		_ = burstRepo.Create(context.Background(), fixtures.Burst("b1", "e1", "e2"))
		_ = factRepo.Create(context.Background(), fixtures.Fact("f1", "e1"))

		model = app.NewModel(cliService, svc)
		model.SkipOnboarding() // Skip onboarding for tests
		Expect(model).NotTo(BeNil())
	})

	Describe("BrowseTimeline navigation from menu", func() {
		It("should activate BrowseTimeline intent when selected from menu", func() {
			// Start at menu
			Expect(model.GetState()).To(Equal(app.StateMenu))

			// Navigate to BrowseTimeline (second menu item)
			modelInterface, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
			model = modelInterface.(*app.Model)

			// Select the BrowseTimeline intent
			modelInterface, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			model = modelInterface.(*app.Model)

			// Should now be in intent state
			Expect(model.GetState()).To(Equal(app.StateIntent))

			// The active intent should be browse_timeline
			activeIntent := model.GetActiveIntent()
			Expect(activeIntent).NotTo(BeNil())
		})

		It("should allow navigation down in BrowseTimeline after activation", func() {
			// Activate BrowseTimeline intent
			modelInterface, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
			model = modelInterface.(*app.Model)
			modelInterface, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			model = modelInterface.(*app.Model)

			// Verify we're in the intent
			Expect(model.GetState()).To(Equal(app.StateIntent))

			// Try to navigate down in the timeline
			modelInterface, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
			model = modelInterface.(*app.Model)

			// Should still be in intent state
			Expect(model.GetState()).To(Equal(app.StateIntent))

			// The view should render without errors
			view := model.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should allow navigation up in BrowseTimeline after activation", func() {
			// Activate BrowseTimeline intent
			modelInterface, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
			model = modelInterface.(*app.Model)
			modelInterface, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			model = modelInterface.(*app.Model)

			// Navigate down first
			modelInterface, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
			model = modelInterface.(*app.Model)

			// Then navigate up
			modelInterface, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})
			model = modelInterface.(*app.Model)

			// Should still be in intent state
			Expect(model.GetState()).To(Equal(app.StateIntent))

			// The view should render without errors
			view := model.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle rapid navigation in BrowseTimeline", func() {
			// Activate BrowseTimeline intent
			modelInterface, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
			model = modelInterface.(*app.Model)
			modelInterface, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			model = modelInterface.(*app.Model)

			// Rapidly navigate down multiple times
			for i := 0; i < 3; i++ {
				modelInterface, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
				model = modelInterface.(*app.Model)
			}

			// Should still be in intent state
			Expect(model.GetState()).To(Equal(app.StateIntent))

			// Rapidly navigate up
			for i := 0; i < 3; i++ {
				modelInterface, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})
				model = modelInterface.(*app.Model)
			}

			// Should still be in intent state
			Expect(model.GetState()).To(Equal(app.StateIntent))

			// The view should render without errors
			view := model.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should use vim-style navigation keys (j/k) in BrowseTimeline", func() {
			// Activate BrowseTimeline intent
			modelInterface, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
			model = modelInterface.(*app.Model)
			modelInterface, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			model = modelInterface.(*app.Model)

			// Navigate down with 'j'
			modelInterface, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
			model = modelInterface.(*app.Model)

			// Should still be in intent state
			Expect(model.GetState()).To(Equal(app.StateIntent))

			// Navigate up with 'k'
			modelInterface, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
			model = modelInterface.(*app.Model)

			// Should still be in intent state
			Expect(model.GetState()).To(Equal(app.StateIntent))

			// The view should render without errors
			view := model.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should maintain state across multiple navigation commands", func() {
			// Activate BrowseTimeline intent
			modelInterface, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
			model = modelInterface.(*app.Model)
			modelInterface, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			model = modelInterface.(*app.Model)

			initialView := model.View()
			Expect(initialView).NotTo(BeEmpty())

			// Navigate down
			modelInterface, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
			model = modelInterface.(*app.Model)
			viewAfterDown := model.View()
			Expect(viewAfterDown).NotTo(BeEmpty())

			// Navigate down again
			modelInterface, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
			model = modelInterface.(*app.Model)
			viewAfterSecondDown := model.View()
			Expect(viewAfterSecondDown).NotTo(BeEmpty())

			// Navigate up
			modelInterface, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})
			model = modelInterface.(*app.Model)
			viewAfterUp := model.View()
			Expect(viewAfterUp).NotTo(BeEmpty())

			// All views should be different (showing different selected items)
			// This is a basic check that navigation actually changes what's displayed
			Expect(viewAfterDown).NotTo(Equal(initialView))
			Expect(viewAfterSecondDown).NotTo(Equal(viewAfterDown))
			Expect(viewAfterUp).NotTo(Equal(viewAfterSecondDown))
		})
	})

	Describe("BrowseTimeline list container synchronization bug", func() {
		It("should properly synchronize list container cursor with intent state", func() {
			// Activate BrowseTimeline intent
			modelInterface, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
			model = modelInterface.(*app.Model)
			modelInterface, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			model = modelInterface.(*app.Model)

			// Get the active intent
			activeIntent := model.GetActiveIntent()
			Expect(activeIntent).NotTo(BeNil())

			// The intent should have proper list container state
			view := model.View()
			Expect(view).NotTo(BeEmpty())

			// Navigate down
			modelInterface, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
			model = modelInterface.(*app.Model)

			// Get the view after navigation
			viewAfterNav := model.View()
			Expect(viewAfterNav).NotTo(BeEmpty())

			// The view should reflect the navigation
			Expect(viewAfterNav).NotTo(Equal(view))
		})

		It("should NOT lose navigation state when rendering multiple times", func() {
			// Activate BrowseTimeline intent
			modelInterface, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
			model = modelInterface.(*app.Model)
			modelInterface, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			model = modelInterface.(*app.Model)

			// Navigate down
			modelInterface, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
			model = modelInterface.(*app.Model)

			// Get view multiple times (simulating multiple renders)
			view1 := model.View()
			view2 := model.View()
			view3 := model.View()

			// All three views should be identical (same state)
			Expect(view1).To(Equal(view2))
			Expect(view2).To(Equal(view3))

			// Navigate down again
			modelInterface, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
			model = modelInterface.(*app.Model)

			// Get view again
			view4 := model.View()

			// This should be different from the previous view
			Expect(view4).NotTo(Equal(view1))
		})
	})

	Describe("Navigation in other list-based intents", func() {
		It("should support navigation in ConfigureSystem intent", func() {
			// Navigate to ConfigureSystem (5th menu item, index 4)
			for i := 0; i < 4; i++ {
				modelInterface, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
				model = modelInterface.(*app.Model)
			}

			// Activate ConfigureSystem intent
			modelInterface, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			model = modelInterface.(*app.Model)

			// Should be in intent state
			Expect(model.GetState()).To(Equal(app.StateIntent))

			// Try to navigate
			modelInterface, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
			model = modelInterface.(*app.Model)

			// Should still be in intent state
			Expect(model.GetState()).To(Equal(app.StateIntent))

			// View should render
			view := model.View()
			Expect(view).NotTo(BeEmpty())
		})
	})
})

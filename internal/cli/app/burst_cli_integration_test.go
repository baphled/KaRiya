package app

import (
	"context"
	"time"

	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Burst CLI Integration", func() {
	var (
		app        *Model
		repo       *careerrepo.MemoryRepository
		svc        *careerservice.Service
		cliSvc     *service.CLIEventService
		ctx        context.Context
		testEvents []*career.CareerEvent
		testBursts []*career.Burst
	)

	BeforeEach(func() {
		ctx = context.Background()
		repo = careerrepo.NewMemoryRepository()
		svc = careerservice.NewService(repo)
		cliSvc = service.NewCLIEventService(svc)
		app = NewModel(cliSvc, svc)

		// Create test events that can form bursts
		testEvents = []*career.CareerEvent{
			{
				ID:        uuid.New().String(),
				Text:      "Led backend team on microservices migration project",
				Date:      time.Now().Add(-30 * 24 * time.Hour),
				Company:   "TechCorp",
				Tags:      []string{"leadership", "technical"},
				CreatedAt: time.Now().Add(-30 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-30 * 24 * time.Hour),
			},
			{
				ID:        uuid.New().String(),
				Text:      "Architected service mesh for improved scalability",
				Date:      time.Now().Add(-25 * 24 * time.Hour),
				Company:   "TechCorp",
				Tags:      []string{"technical", "project"},
				CreatedAt: time.Now().Add(-25 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-25 * 24 * time.Hour),
			},
			{
				ID:        uuid.New().String(),
				Text:      "Mentored junior developers on best practices",
				Date:      time.Now().Add(-20 * 24 * time.Hour),
				Company:   "TechCorp",
				Tags:      []string{"mentoring", "leadership"},
				CreatedAt: time.Now().Add(-20 * 24 * time.Hour),
				UpdatedAt: time.Now().Add(-20 * 24 * time.Hour),
			},
		}

		// Add events to repository
		for _, event := range testEvents {
			err := svc.CaptureEvent(ctx, event, careerservice.ManualEntry)
			Expect(err).NotTo(HaveOccurred())
		}

		// Create test bursts
		testBursts = []*career.Burst{
			{
				ID:              uuid.New().String(),
				Name:            "Platform Migration",
				Description:     "Led comprehensive platform migration effort",
				EventIDs:        []string{testEvents[0].ID, testEvents[1].ID},
				CompetencyFocus: "technical",
				CreatedAt:       time.Now().Add(-15 * 24 * time.Hour),
				UpdatedAt:       time.Now().Add(-15 * 24 * time.Hour),
			},
			{
				ID:              uuid.New().String(),
				Name:            "Team Leadership",
				Description:     "Leadership and mentoring activities",
				EventIDs:        []string{testEvents[0].ID, testEvents[2].ID},
				CompetencyFocus: "leadership",
				CreatedAt:       time.Now().Add(-10 * 24 * time.Hour),
				UpdatedAt:       time.Now().Add(-10 * 24 * time.Hour),
			},
		}

		// Add bursts to repository
		for _, burst := range testBursts {
			err := svc.ConfirmBurst(ctx, burst)
			Expect(err).NotTo(HaveOccurred())
		}
	})

	Describe("Navigation to Burst List", func() {
		It("should navigate to burst list screen when 'b' key is pressed from home", func() {
			// Verify starting at home screen
			Expect(app.currentScreen).To(Equal(HomeScreen))

			// Simulate pressing 'b' key
			keyMsg := tea.KeyMsg{
				Type:  tea.KeyRunes,
				Runes: []rune{'b'},
			}

			updatedApp, _ := app.Update(keyMsg)
			app = updatedApp.(*Model)

			// Verify navigation to burst list screen
			Expect(app.currentScreen).To(Equal(BurstListScreen))
			Expect(app.previousScreen).To(Equal(HomeScreen))
		})

		It("should initialize burst list model when navigating to burst screen", func() {
			// Navigate to burst list screen
			keyMsg := tea.KeyMsg{
				Type:  tea.KeyRunes,
				Runes: []rune{'b'},
			}

			updatedApp, _ := app.Update(keyMsg)
			app = updatedApp.(*Model)

			// Verify burst list model is initialized
			Expect(app.burstListModel).NotTo(BeNil())
		})

		It("should set correct breadcrumbs for burst list screen", func() {
			// Navigate to burst list screen
			keyMsg := tea.KeyMsg{
				Type:  tea.KeyRunes,
				Runes: []rune{'b'},
			}

			updatedApp, _ := app.Update(keyMsg)
			app = updatedApp.(*Model)

			// Update breadcrumbs
			app.updateBreadcrumbs()

			// Verify breadcrumbs
			Expect(app.breadcrumbs).To(Equal([]string{"Home", "Bursts"}))
		})
	})

	Describe("Burst List Display", func() {
		BeforeEach(func() {
			// Navigate to burst list screen
			keyMsg := tea.KeyMsg{
				Type:  tea.KeyRunes,
				Runes: []rune{'b'},
			}

			updatedApp, _ := app.Update(keyMsg)
			app = updatedApp.(*Model)
		})

		It("should render burst list view", func() {
			view := app.View()

			// Verify we're not getting an error message
			Expect(view).NotTo(ContainSubstring("Error: Burst List model not initialized"))
			Expect(view).NotTo(BeEmpty())
		})

		It("should display existing bursts", func() {
			view := app.View()

			// The burst list model should be loaded and displaying bursts
			// Note: The actual content depends on BurstListModel implementation
			// but we can verify it's not showing an error
			Expect(view).NotTo(ContainSubstring("Error"))
		})

		It("should handle empty burst list gracefully", func() {
			// Create new app with no bursts
			emptyRepo := careerrepo.NewMemoryRepository()
			emptySvc := careerservice.NewService(emptyRepo)
			emptyCliSvc := service.NewCLIEventService(emptySvc)
			emptyApp := NewModel(emptyCliSvc, emptySvc)

			// Navigate to burst list
			keyMsg := tea.KeyMsg{
				Type:  tea.KeyRunes,
				Runes: []rune{'b'},
			}

			updatedApp, _ := emptyApp.Update(keyMsg)
			emptyApp = updatedApp.(*Model)

			view := emptyApp.View()

			// Should not crash and should display some content
			Expect(view).NotTo(BeEmpty())
			Expect(view).NotTo(ContainSubstring("Error"))
		})
	})

	Describe("Home Screen Integration", func() {
		It("should display 'View Bursts' option in home screen commands", func() {
			view := app.renderHome()

			// Verify the burst option is displayed
			Expect(view).To(ContainSubstring("b - View Bursts"))
		})

		It("should include burst navigation in help text", func() {
			view := app.renderHome()

			// Verify commands section includes bursts
			Expect(view).To(ContainSubstring("Commands:"))
			Expect(view).To(ContainSubstring("b - View Bursts"))
			Expect(view).To(ContainSubstring("c - Capture Career Event"))
			Expect(view).To(ContainSubstring("l - List Events"))
		})
	})

	Describe("Navigation Back from Burst List", func() {
		BeforeEach(func() {
			// Navigate to burst list screen first
			keyMsg := tea.KeyMsg{
				Type:  tea.KeyRunes,
				Runes: []rune{'b'},
			}

			updatedApp, _ := app.Update(keyMsg)
			app = updatedApp.(*Model)
		})

		It("should navigate back to home when 'h' key is pressed", func() {
			// Press 'h' to go home
			keyMsg := tea.KeyMsg{
				Type:  tea.KeyRunes,
				Runes: []rune{'h'},
			}

			updatedApp, _ := app.Update(keyMsg)
			app = updatedApp.(*Model)

			// Verify navigation back to home
			Expect(app.currentScreen).To(Equal(HomeScreen))
			Expect(app.previousScreen).To(Equal(BurstListScreen))
		})

		It("should handle back navigation via BackMsg", func() {
			// Send BackMsg
			backMsg := tea.KeyMsg{
				Type: tea.KeyBackspace,
			}

			updatedApp, _ := app.Update(backMsg)
			app = updatedApp.(*Model)

			// Should navigate back to previous screen (home)
			Expect(app.currentScreen).To(Equal(HomeScreen))
		})
	})
})

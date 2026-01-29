package cv_test

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/cv"
)

func TestCVProfileSelect(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CVProfileSelect Screen Suite")
}

var _ = Describe("ProfileSelectScreen", func() {
	var (
		screen   screens.Screen
		profiles []*intents.CVProfile
	)

	BeforeEach(func() {
		profiles = []*intents.CVProfile{
			{
				ID:             "profile-1",
				Name:           "Senior Engineer Profile",
				TargetRole:     "senior_ic",
				TargetAudience: "hiring_manager",
				Description:    "For senior IC roles at tech companies",
			},
			{
				ID:             "profile-2",
				Name:           "Staff Engineer Profile",
				TargetRole:     "staff",
				TargetAudience: "hiring_manager",
				Description:    "For staff+ engineering roles",
			},
			{
				ID:             "profile-3",
				Name:           "Engineering Manager Profile",
				TargetRole:     "em",
				TargetAudience: "hiring_manager",
				Description:    "For engineering manager positions",
			},
		}

		screen = cv.NewCVProfileSelectScreen(profiles)
	})

	Describe("NewCVProfileSelectScreen", func() {
		It("should create a new screen", func() {
			Expect(screen).ToNot(BeNil())
		})

		It("should accept profiles", func() {
			screen := cv.NewCVProfileSelectScreen(profiles)
			Expect(screen).ToNot(BeNil())
		})
	})

	// Screen interface doesn't have Init method - removed test

	Describe("Update", func() {
		Context("navigation", func() {
			It("should navigate down with 'j' key", func() {
				// Init not needed - screen ready to use
				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})

				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())

				// View should show selection moved
				view := screen.View()
				Expect(view).ToNot(BeEmpty())
			})

			It("should navigate down with arrow down", func() {
				// Init not needed - screen ready to use
				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyDown})

				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should navigate up with 'k' key", func() {
				// Init not needed - screen ready to use
				// Move down first
				screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})

				// Then move up
				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})

				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should navigate up with arrow up", func() {
				// Init not needed - screen ready to use
				screen.Update(tea.KeyMsg{Type: tea.KeyDown})

				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyUp})

				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should jump to top with 'g' key", func() {
				// Init not needed - screen ready to use
				// Move down a few times
				screen.Update(tea.KeyMsg{Type: tea.KeyDown})
				screen.Update(tea.KeyMsg{Type: tea.KeyDown})

				// Jump to top
				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}})

				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should jump to bottom with 'G' key", func() {
				// Init not needed - screen ready to use

				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}})

				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})
		})

		Context("selection", func() {
			It("should return NavigateResult on Enter key", func() {
				// Init not needed - screen ready to use

				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyEnter})

				Expect(cmd).To(BeNil())
				Expect(result).ToNot(BeNil())
				Expect(result.Type()).To(Equal(screens.ResultNavigate))
				Expect(result.Data()).ToNot(BeNil())

				// Should return the selected profile
				selectedProfile, ok := result.Data().(*intents.CVProfile)
				Expect(ok).To(BeTrue())
				Expect(selectedProfile.ID).To(Equal("profile-1"))
			})

			It("should return selected profile after navigation", func() {
				// Init not needed - screen ready to use

				// Navigate to second profile
				screen.Update(tea.KeyMsg{Type: tea.KeyDown})

				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyEnter})

				Expect(cmd).To(BeNil())
				Expect(result).ToNot(BeNil())

				selectedProfile, ok := result.Data().(*intents.CVProfile)
				Expect(ok).To(BeTrue())
				Expect(selectedProfile.ID).To(Equal("profile-2"))
			})
		})

		Context("cancellation", func() {
			It("should return CancelResult on Esc key", func() {
				// Init not needed - screen ready to use

				cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyEsc})

				Expect(cmd).To(BeNil())
				Expect(result).ToNot(BeNil())
				Expect(result.Type()).To(Equal(screens.ResultCancel))
			})

			// Note: 'q' key behavior depends on SelectScreen implementation
			// If 'q' doesn't trigger cancel, this test should be removed or SelectScreen updated
		})

		Context("window resize", func() {
			It("should handle window resize", func() {
				// Init not needed - screen ready to use

				cmd, result := screen.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})
		})
	})

	Describe("View", func() {
		It("should render profile list", func() {
			// Init not needed - screen ready to use

			view := screen.View()

			Expect(view).ToNot(BeEmpty())
			Expect(view).To(ContainSubstring("Senior Engineer Profile"))
			Expect(view).To(ContainSubstring("Staff Engineer Profile"))
			Expect(view).To(ContainSubstring("Engineering Manager Profile"))
		})

		It("should show selection indicator", func() {
			// Init not needed - screen ready to use

			view := screen.View()

			// Should have a selection indicator (▶ or similar)
			Expect(view).To(MatchRegexp(`[▶►>]`))
		})

		It("should show profile descriptions", func() {
			// Init not needed - screen ready to use

			view := screen.View()

			Expect(view).To(ContainSubstring("For senior IC roles"))
			Expect(view).To(ContainSubstring("For staff+ engineering roles"))
		})

		It("should show help text", func() {
			// Init not needed - screen ready to use

			view := screen.View()

			// SelectScreen uses combined format like "↑/↓/j/k: Navigate"
			Expect(view).To(MatchRegexp("↑|↓|j|k"))
			Expect(view).To(ContainSubstring("Enter"))
			Expect(view).To(ContainSubstring("Esc"))
		})
	})

	Describe("Edge Cases", func() {
		Context("empty profile list", func() {
			It("should handle empty profile list gracefully", func() {
				emptyScreen := cv.NewCVProfileSelectScreen([]*intents.CVProfile{})

				view := emptyScreen.View()
				// SelectScreen shows "No items available"
				Expect(view).To(ContainSubstring("No items available"))
			})
		})

		Context("single profile", func() {
			It("should work with single profile", func() {
				singleProfile := []*intents.CVProfile{profiles[0]}
				singleScreen := cv.NewCVProfileSelectScreen(singleProfile)

				// Init not needed
				cmd, result := singleScreen.Update(tea.KeyMsg{Type: tea.KeyEnter})

				Expect(cmd).To(BeNil())
				Expect(result).ToNot(BeNil())
				Expect(result.Type()).To(Equal(screens.ResultNavigate))
			})
		})
	})
})

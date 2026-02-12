package burst_management_test

import (
	"time"

	"github.com/baphled/kariya/internal/cli/screens"
	burstscreens "github.com/baphled/kariya/internal/cli/screens/burst_management"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("BurstDetailScreen", func() {
	var (
		screen *burstscreens.BurstDetailScreen
		burst  *career.Burst
	)

	BeforeEach(func() {
		burst = fixtures.Burst("burst-1", "event-1", "event-2", "event-3")
		burst.Name = "Senior Backend Engineer Career Growth"
		burst.Description = "Progressed from mid-level to senior backend engineer, focusing on system architecture and team leadership."
		burst.CreatedAt = time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
		burst.UpdatedAt = time.Date(2024, 1, 20, 14, 45, 0, 0, time.UTC)
	})

	Describe("Construction", func() {
		It("should create a burst detail screen", func() {
			screen = burstscreens.NewBurstDetailScreen(burst)
			Expect(screen).NotTo(BeNil())
		})

		It("should store burst", func() {
			screen = burstscreens.NewBurstDetailScreen(burst)
			Expect(screen.GetBurst()).To(Equal(burst))
		})

		It("should handle nil burst gracefully", func() {
			screen = burstscreens.NewBurstDetailScreen(nil)
			Expect(screen).NotTo(BeNil())
			Expect(screen.GetBurst()).To(BeNil())
		})
	})

	Describe("Actions", func() {
		BeforeEach(func() {
			screen = burstscreens.NewBurstDetailScreen(burst)
			screen.SetTerminalInfo(120, 40)
		})

		It("should return view_events action on 'v' key", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultNavigate))
			navResult := result.(*screens.NavigateResult)
			data := navResult.ResultData.(map[string]interface{})
			Expect(data["action"]).To(Equal("view_events"))
			Expect(data["burst"]).To(Equal(burst))
		})

		It("should return view_facts action on 'f' key", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultNavigate))
			navResult := result.(*screens.NavigateResult)
			data := navResult.ResultData.(map[string]interface{})
			Expect(data["action"]).To(Equal("view_facts"))
			Expect(data["burst"]).To(Equal(burst))
		})

		It("should return edit action on 'e' key", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultNavigate))
			navResult := result.(*screens.NavigateResult)
			data := navResult.ResultData.(map[string]interface{})
			Expect(data["action"]).To(Equal("edit"))
			Expect(data["burst"]).To(Equal(burst))
		})

		It("should return delete action on 'd' key", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultNavigate))
			navResult := result.(*screens.NavigateResult)
			data := navResult.ResultData.(map[string]interface{})
			Expect(data["action"]).To(Equal("delete"))
			Expect(data["burst"]).To(Equal(burst))
		})

		It("should return confirm action on 'c' key", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultNavigate))
			navResult := result.(*screens.NavigateResult)
			data := navResult.ResultData.(map[string]interface{})
			Expect(data["action"]).To(Equal("confirm"))
			Expect(data["burst"]).To(Equal(burst))
		})

		It("should return skills action on 's' key", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultNavigate))
			navResult := result.(*screens.NavigateResult)
			data := navResult.ResultData.(map[string]interface{})
			Expect(data["action"]).To(Equal("skills"))
			Expect(data["burst"]).To(Equal(burst))
		})

		It("should return infer action on 'i' key", func() {
			// Create a confirmed burst for this test
			confirmedBurst := fixtures.Burst("burst-confirmed", "event-1", "event-2", "event-3")
			confirmedBurst.Name = "Senior Backend Engineer"
			confirmedBurst.Confirmed = true
			confirmedScreen := burstscreens.NewBurstDetailScreen(confirmedBurst)
			confirmedScreen.SetTerminalInfo(120, 40)

			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}}
			_, result := confirmedScreen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultNavigate))
			navResult := result.(*screens.NavigateResult)
			data := navResult.ResultData.(map[string]interface{})
			Expect(data["action"]).To(Equal("infer"))
			Expect(data["burst"]).To(Equal(confirmedBurst))
		})

		It("should return help action on '?' key", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultNavigate))
			navResult := result.(*screens.NavigateResult)
			data := navResult.ResultData.(map[string]interface{})
			Expect(data["action"]).To(Equal("help"))
		})

		It("should return nil for 'i' on unconfirmed burst", func() {
			unconfirmedBurst := fixtures.Burst("burst-unconfirmed", "event-1", "event-2")
			unconfirmedBurst.Confirmed = false
			unconfirmedScreen := burstscreens.NewBurstDetailScreen(unconfirmedBurst)
			unconfirmedScreen.SetTerminalInfo(120, 40)

			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}}
			_, result := unconfirmedScreen.Update(msg)

			Expect(result).To(BeNil())
		})
	})

	Describe("Navigation", func() {
		BeforeEach(func() {
			screen = burstscreens.NewBurstDetailScreen(burst)
		})

		It("should return CancelResult on escape", func() {
			msg := tea.KeyMsg{Type: tea.KeyEsc}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultCancel))
		})

		It("should return CancelResult on backspace", func() {
			msg := tea.KeyMsg{Type: tea.KeyBackspace}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultCancel))
		})

		It("should return nil result for unhandled keys", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}}
			_, result := screen.Update(msg)

			Expect(result).To(BeNil())
		})
	})

	Describe("Rendering", func() {
		BeforeEach(func() {
			screen = burstscreens.NewBurstDetailScreen(burst)
			screen.SetTerminalInfo(120, 40)
		})

		It("should render content", func() {
			content := screen.RenderContent()
			Expect(content).NotTo(BeEmpty())
		})

		It("should display burst name", func() {
			content := screen.RenderContent()
			Expect(content).To(ContainSubstring("Senior Backend Engineer Career Growth"))
		})

		It("should display burst description", func() {
			content := screen.RenderContent()
			Expect(content).To(ContainSubstring("Progressed from mid-level to senior"))
		})

		It("should display event count", func() {
			content := screen.RenderContent()
			Expect(content).To(ContainSubstring("3")) // 3 events
		})

		It("should display confirmed status when false", func() {
			content := screen.RenderContent()
			Expect(content).To(ContainSubstring("No"))
		})

		It("should display confirmed status when true", func() {
			burst.Confirmed = true
			screen = burstscreens.NewBurstDetailScreen(burst)
			screen.SetTerminalInfo(120, 40)

			content := screen.RenderContent()
			Expect(content).To(ContainSubstring("Yes"))
		})

		It("should display created date", func() {
			content := screen.RenderContent()
			Expect(content).To(ContainSubstring("2024-01-15"))
		})

		It("should display updated date", func() {
			content := screen.RenderContent()
			Expect(content).To(ContainSubstring("2024-01-20"))
		})

		It("should handle nil burst in rendering", func() {
			screen = burstscreens.NewBurstDetailScreen(nil)
			screen.SetTerminalInfo(120, 40)

			content := screen.RenderContent()
			Expect(content).To(ContainSubstring("No burst selected"))
		})
	})

	Describe("Window Size Updates", func() {
		BeforeEach(func() {
			screen = burstscreens.NewBurstDetailScreen(burst)
		})

		It("should handle window size changes", func() {
			msg := tea.WindowSizeMsg{Width: 100, Height: 30}
			cmd, result := screen.Update(msg)

			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
		})

		It("should update terminal info on window size change", func() {
			initialWidth := 80
			initialHeight := 24
			screen.SetTerminalInfo(initialWidth, initialHeight)

			newWidth := 120
			newHeight := 40
			msg := tea.WindowSizeMsg{Width: newWidth, Height: newHeight}
			screen.Update(msg)

			// Verify the screen received the new dimensions.
			// We can't directly check terminal info, but we can verify it doesn't panic.
			content := screen.RenderContent()
			Expect(content).NotTo(BeEmpty())
		})
	})

	Describe("Help Toggle", func() {
		BeforeEach(func() {
			screen = burstscreens.NewBurstDetailScreen(burst)
			screen.SetTerminalInfo(120, 40)
		})

		It("should return help action on '?' key", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
			_, result := screen.Update(msg)

			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultNavigate))
			navResult := result.(*screens.NavigateResult)
			data := navResult.ResultData.(map[string]interface{})
			Expect(data["action"]).To(Equal("help"))
		})
	})
})

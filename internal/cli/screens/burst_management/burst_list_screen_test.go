package burst_management_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/screens"
	burst_management "github.com/baphled/kariya/internal/cli/screens/burst_management"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	tea "github.com/charmbracelet/bubbletea"
)

var _ = Describe("BurstListScreen", func() {
	var (
		screen *burst_management.BurstListScreen
		bursts []*career.Burst
	)

	BeforeEach(func() {
		now := time.Now()
		b1 := fixtures.BurstConfirmed("burst-1", "e1", "e2", "e3")
		b1.Name = "First Achievement Period"
		b1.Description = "Led major platform migration"
		b1.CreatedAt = now

		b2 := fixtures.Burst("burst-2", "e4", "e5")
		b2.Name = "Team Leadership Growth"
		b2.Description = "Grew team from 3 to 12 engineers"
		b2.CreatedAt = now.Add(-24 * time.Hour)

		bursts = []*career.Burst{b1, b2}
	})

	Describe("NewBurstListScreen", func() {
		It("should create a new screen with bursts", func() {
			screen = burst_management.NewBurstListScreen(bursts)
			Expect(screen).NotTo(BeNil())
			Expect(screen.GetBursts()).To(HaveLen(2))
		})

		It("should handle empty burst list", func() {
			screen = burst_management.NewBurstListScreen([]*career.Burst{})
			Expect(screen).NotTo(BeNil())
			Expect(screen.GetBursts()).To(BeEmpty())
		})

		It("should handle nil burst list", func() {
			screen = burst_management.NewBurstListScreen(nil)
			Expect(screen).NotTo(BeNil())
		})
	})

	Describe("Update", func() {
		BeforeEach(func() {
			screen = burst_management.NewBurstListScreen(bursts)
		})

		Context("window resize", func() {
			It("should handle window size message", func() {
				msg := tea.WindowSizeMsg{Width: 120, Height: 40}
				cmd, result := screen.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})
		})

		Context("escape key", func() {
			It("should return cancel result on escape", func() {
				msg := tea.KeyMsg{Type: tea.KeyEsc}
				cmd, result := screen.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).NotTo(BeNil())

				cancelResult, ok := result.(*screens.CancelResult)
				Expect(ok).To(BeTrue())
				Expect(cancelResult).NotTo(BeNil())
			})
		})

		Context("navigation keys", func() {
			It("should handle up arrow", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("up")}
				cmd, result := screen.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should handle down arrow", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("down")}
				cmd, result := screen.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should handle j (vim down)", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
				cmd, result := screen.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should handle k (vim up)", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
				cmd, result := screen.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should handle ctrl+d (page down)", func() {
				msg := tea.KeyMsg{Type: tea.KeyCtrlD}
				cmd, result := screen.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should handle ctrl+u (page up)", func() {
				msg := tea.KeyMsg{Type: tea.KeyCtrlU}
				cmd, result := screen.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should handle home key (go to first)", func() {
				msg := tea.KeyMsg{Type: tea.KeyHome}
				cmd, result := screen.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should handle end key (go to last)", func() {
				msg := tea.KeyMsg{Type: tea.KeyEnd}
				cmd, result := screen.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should handle g key (vim go to first)", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}}
				cmd, result := screen.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should handle G key (vim go to last)", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}}
				cmd, result := screen.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})
		})

		Context("enter key - view details", func() {
			It("should return navigate result with selected burst", func() {
				msg := tea.KeyMsg{Type: tea.KeyEnter}
				cmd, result := screen.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).NotTo(BeNil())

				navResult, ok := result.(*screens.NavigateResult)
				Expect(ok).To(BeTrue())

				burst, ok := navResult.ResultData.(*career.Burst)
				Expect(ok).To(BeTrue())
				Expect(burst.ID).To(Equal("burst-1"))
			})

			It("should handle enter with no selection", func() {
				emptyScreen := burst_management.NewBurstListScreen([]*career.Burst{})
				msg := tea.KeyMsg{Type: tea.KeyEnter}
				cmd, result := emptyScreen.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})
		})

		Context("action keys", func() {
			It("should handle 'a' for add", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
				cmd, result := screen.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).NotTo(BeNil())

				navResult, ok := result.(*screens.NavigateResult)
				Expect(ok).To(BeTrue())

				actionData, ok := navResult.ResultData.(map[string]interface{})
				Expect(ok).To(BeTrue())
				Expect(actionData["action"]).To(Equal("add"))
			})

			It("should handle 'e' for edit with selection", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
				cmd, result := screen.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).NotTo(BeNil())

				navResult, ok := result.(*screens.NavigateResult)
				Expect(ok).To(BeTrue())

				actionData, ok := navResult.ResultData.(map[string]interface{})
				Expect(ok).To(BeTrue())
				Expect(actionData["action"]).To(Equal("edit"))
				Expect(actionData["burst"]).NotTo(BeNil())
			})

			It("should handle 'e' with no selection", func() {
				emptyScreen := burst_management.NewBurstListScreen([]*career.Burst{})
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
				cmd, result := emptyScreen.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should handle 'd' for delete with selection", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}
				cmd, result := screen.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).NotTo(BeNil())

				navResult, ok := result.(*screens.NavigateResult)
				Expect(ok).To(BeTrue())

				actionData, ok := navResult.ResultData.(map[string]interface{})
				Expect(ok).To(BeTrue())
				Expect(actionData["action"]).To(Equal("delete"))
				Expect(actionData["burst"]).NotTo(BeNil())
			})

			It("should handle 'd' with no selection", func() {
				emptyScreen := burst_management.NewBurstListScreen([]*career.Burst{})
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}
				cmd, result := emptyScreen.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should handle 's' for suggest", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
				cmd, result := screen.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).NotTo(BeNil())

				navResult, ok := result.(*screens.NavigateResult)
				Expect(ok).To(BeTrue())

				actionData, ok := navResult.ResultData.(map[string]interface{})
				Expect(ok).To(BeTrue())
				Expect(actionData["action"]).To(Equal("suggest"))
			})
		})

		Context("unhandled keys", func() {
			It("should return nil for unhandled keys", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}}
				cmd, result := screen.Update(msg)
				Expect(cmd).To(BeNil())
				Expect(result).To(BeNil())
			})
		})
	})

	Describe("RenderContent", func() {
		BeforeEach(func() {
			screen = burst_management.NewBurstListScreen(bursts)
		})

		It("should render table content", func() {
			content := screen.RenderContent()
			Expect(content).NotTo(BeEmpty())
			Expect(content).To(ContainSubstring("Name"))
			Expect(content).To(ContainSubstring("Description"))
		})

		It("should show burst names in content", func() {
			content := screen.RenderContent()
			Expect(content).To(ContainSubstring("First Achievement Period"))
			Expect(content).To(ContainSubstring("Team Leadership Growth"))
		})

		It("should show confirmed status", func() {
			content := screen.RenderContent()
			Expect(content).To(ContainSubstring("✓ Yes"))
			Expect(content).To(ContainSubstring("✗ No"))
		})

		It("should show event counts", func() {
			content := screen.RenderContent()
			Expect(content).To(ContainSubstring("3"))
			Expect(content).To(ContainSubstring("2"))
		})

		It("should show empty message when no bursts", func() {
			emptyScreen := burst_management.NewBurstListScreen([]*career.Burst{})
			content := emptyScreen.RenderContent()
			Expect(content).To(ContainSubstring("No bursts found"))
		})
	})

	Describe("View", func() {
		BeforeEach(func() {
			screen = burst_management.NewBurstListScreen(bursts)
		})

		It("should render complete view", func() {
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should include breadcrumbs", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Main Menu"))
			Expect(view).To(ContainSubstring("Burst Management"))
		})

		It("should include help footer", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("Navigate"))
			Expect(view).To(ContainSubstring("Add"))
			Expect(view).To(ContainSubstring("Edit"))
			Expect(view).To(ContainSubstring("Delete"))
			Expect(view).To(ContainSubstring("Suggest"))
		})

		It("should include content", func() {
			view := screen.View()
			Expect(view).To(ContainSubstring("First Achievement Period"))
		})
	})

	Describe("SetTheme", func() {
		BeforeEach(func() {
			screen = burst_management.NewBurstListScreen(bursts)
		})

		It("should accept theme", func() {
			theme := themes.NewDefaultTheme()
			Expect(func() {
				screen.SetTheme(theme)
			}).NotTo(Panic())
		})

		It("should handle nil theme", func() {
			Expect(func() {
				screen.SetTheme(nil)
			}).NotTo(Panic())
		})
	})

	Describe("GetSelectedIndex", func() {
		BeforeEach(func() {
			screen = burst_management.NewBurstListScreen(bursts)
		})

		It("should return initial index of 0", func() {
			Expect(screen.GetSelectedIndex()).To(Equal(0))
		})

		It("should update index after navigation", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
			screen.Update(msg)
			Expect(screen.GetSelectedIndex()).To(Equal(1))
		})
	})

	Describe("Empty List Safeguards - Burst Management (Phase 2-Tier 3)", func() {
		It("should return nil for view on empty list", func() {
			emptyScreen := burst_management.NewBurstListScreen([]*career.Burst{})

			msg := tea.KeyMsg{Type: tea.KeyEnter}
			_, result := emptyScreen.Update(msg)

			Expect(result).To(BeNil())
		})
	})
	Describe("Update with unhandled message", func() {
		BeforeEach(func() {
			screen = burst_management.NewBurstListScreen(bursts)
		})

		It("should return nil for non-key non-window messages", func() {
			type customMsg struct{}
			cmd, result := screen.Update(customMsg{})
			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
		})
	})

	Describe("Arrow key navigation", func() {
		BeforeEach(func() {
			screen = burst_management.NewBurstListScreen(bursts)
		})

		It("should handle up arrow key type", func() {
			msg := tea.KeyMsg{Type: tea.KeyUp}
			cmd, result := screen.Update(msg)
			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
		})

		It("should handle down arrow key type", func() {
			msg := tea.KeyMsg{Type: tea.KeyDown}
			cmd, result := screen.Update(msg)
			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
		})

		It("should handle page down key type", func() {
			msg := tea.KeyMsg{Type: tea.KeyPgDown}
			cmd, result := screen.Update(msg)
			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
		})

		It("should handle page up key type", func() {
			msg := tea.KeyMsg{Type: tea.KeyPgUp}
			cmd, result := screen.Update(msg)
			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
		})
	})

	Describe("View with theme", func() {
		It("should render with explicit theme set", func() {
			screen = burst_management.NewBurstListScreen(bursts)
			screen.SetTheme(themes.NewDefaultTheme())
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
			Expect(view).To(ContainSubstring("First Achievement Period"))
		})

		It("should render with non-Theme interface value", func() {
			screen = burst_management.NewBurstListScreen(bursts)
			screen.SetTheme("not-a-theme")
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("burstRowFormatter edge cases", func() {
		It("should truncate long burst names", func() {
			longNameBurst := fixtures.Burst("burst-long", "e1")
			longNameBurst.Name = "This Is A Very Long Burst Name That Exceeds Twenty Seven Characters"
			longNameBurst.Description = "Short desc"
			screen = burst_management.NewBurstListScreen([]*career.Burst{longNameBurst})
			content := screen.RenderContent()
			Expect(content).NotTo(ContainSubstring("This Is A Very Long Burst Name That Exceeds Twenty Seven Characters"))
		})

		It("should show dash for empty description", func() {
			emptyDescBurst := fixtures.Burst("burst-empty-desc", "e1")
			emptyDescBurst.Name = "Test Burst"
			emptyDescBurst.Description = ""
			screen = burst_management.NewBurstListScreen([]*career.Burst{emptyDescBurst})
			content := screen.RenderContent()
			Expect(content).To(ContainSubstring("-"))
		})

		It("should truncate long descriptions", func() {
			longDescBurst := fixtures.Burst("burst-long-desc", "e1")
			longDescBurst.Name = "Test Burst"
			longDescBurst.Description = "This description is quite long and exceeds thirty two characters easily"
			screen = burst_management.NewBurstListScreen([]*career.Burst{longDescBurst})
			content := screen.RenderContent()
			Expect(content).To(ContainSubstring("..."))
		})

		It("should replace newlines in descriptions", func() {
			newlineBurst := fixtures.Burst("burst-newline", "e1")
			newlineBurst.Name = "Test Burst"
			newlineBurst.Description = "Line one\nLine two\rLine three"
			screen = burst_management.NewBurstListScreen([]*career.Burst{newlineBurst})
			content := screen.RenderContent()
			Expect(content).To(ContainSubstring("Line one Line two Line three"))
		})

		It("should handle whitespace-only description", func() {
			whitespaceBurst := fixtures.Burst("burst-ws", "e1")
			whitespaceBurst.Name = "Test Burst"
			whitespaceBurst.Description = "   \t  "
			screen = burst_management.NewBurstListScreen([]*career.Burst{whitespaceBurst})
			content := screen.RenderContent()
			Expect(content).To(ContainSubstring("-"))
		})
	})
})

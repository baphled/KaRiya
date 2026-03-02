package facts_test

import (
	"testing"
	"time"

	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/facts"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFactListScreen(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "FactListScreen Suite")
}

var _ = Describe("FactListScreen", func() {
	var (
		screen    *facts.FactListScreen
		factList  []*career.Fact
		createdAt time.Time
	)

	BeforeEach(func() {
		createdAt = time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

		f1 := fixtures.FactWith("fact-1", "Led development of microservices architecture")
		f1.CreatedAt = createdAt
		f1.UpdatedAt = createdAt

		f2 := fixtures.FactWith("fact-2", "Reduced API response time by 40%")
		f2.CreatedAt = createdAt.Add(24 * time.Hour)
		f2.UpdatedAt = createdAt.Add(24 * time.Hour)

		f3 := fixtures.FactWith("fact-3", "Mentored 5 junior engineers")
		f3.CreatedAt = createdAt.Add(48 * time.Hour)
		f3.UpdatedAt = createdAt.Add(48 * time.Hour)

		factList = []*career.Fact{f1, f2, f3}
	})

	Describe("Construction", func() {
		It("should create a fact list screen", func() {
			screen = facts.NewFactListScreen(factList)
			Expect(screen).NotTo(BeNil())
		})

		It("should handle empty fact list", func() {
			screen = facts.NewFactListScreen([]*career.Fact{})
			Expect(screen).NotTo(BeNil())
		})

		It("should handle nil fact list", func() {
			screen = facts.NewFactListScreen(nil)
			Expect(screen).NotTo(BeNil())
		})
	})

	Describe("Navigation", func() {
		BeforeEach(func() {
			screen = facts.NewFactListScreen(factList)
			screen.SetTerminalInfo(120, 40)
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

		It("should navigate down with arrow key", func() {
			msg := tea.KeyMsg{Type: tea.KeyDown}
			_, result := screen.Update(msg)

			Expect(result).To(BeNil())
		})

		It("should navigate up with arrow key", func() {
			msg := tea.KeyMsg{Type: tea.KeyUp}
			_, result := screen.Update(msg)

			Expect(result).To(BeNil())
		})

		It("should handle home key (go to first)", func() {
			msg := tea.KeyMsg{Type: tea.KeyHome}
			_, result := screen.Update(msg)
			Expect(result).To(BeNil())
		})

		It("should handle end key (go to last)", func() {
			msg := tea.KeyMsg{Type: tea.KeyEnd}
			_, result := screen.Update(msg)
			Expect(result).To(BeNil())
		})

		It("should handle g key (vim go to first)", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}}
			_, result := screen.Update(msg)
			Expect(result).To(BeNil())
		})

		It("should handle G key (vim go to last)", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}}
			_, result := screen.Update(msg)
			Expect(result).To(BeNil())
		})

		It("should handle Ctrl+D for page down", func() {
			msg := tea.KeyMsg{Type: tea.KeyCtrlD}
			_, result := screen.Update(msg)
			Expect(result).To(BeNil())
		})

		It("should handle Ctrl+U for page up", func() {
			msg := tea.KeyMsg{Type: tea.KeyCtrlU}
			_, result := screen.Update(msg)
			Expect(result).To(BeNil())
		})

		It("should navigate down with j key (vim)", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
			_, result := screen.Update(msg)
			Expect(result).To(BeNil())
		})

		It("should navigate up with k key (vim)", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
			_, result := screen.Update(msg)
			Expect(result).To(BeNil())
		})

		It("should ignore unhandled rune keys", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}}
			_, result := screen.Update(msg)
			Expect(result).To(BeNil())
		})
	})

	Describe("Rendering", func() {
		BeforeEach(func() {
			screen = facts.NewFactListScreen(factList)
			screen.SetTerminalInfo(120, 40)
		})

		It("should render content", func() {
			content := screen.RenderContent()
			Expect(content).NotTo(BeEmpty())
		})

		It("should display fact text", func() {
			content := screen.RenderContent()
			Expect(content).To(ContainSubstring("Led development of microservices"))
		})

		It("should display multiple facts", func() {
			content := screen.RenderContent()
			Expect(content).To(ContainSubstring("Led development"))
			Expect(content).To(ContainSubstring("Reduced API response time"))
			Expect(content).To(ContainSubstring("Mentored 5 junior"))
		})

		It("should show empty message when no facts", func() {
			screen = facts.NewFactListScreen([]*career.Fact{})
			screen.SetTerminalInfo(120, 40)

			content := screen.RenderContent()
			Expect(content).To(ContainSubstring("No facts found"))
		})

		It("should handle nil theme gracefully", func() {
			screen = facts.NewFactListScreen(factList)
			// Don't set theme

			content := screen.RenderContent()
			Expect(content).NotTo(BeEmpty())
		})

		It("should return output from View", func() {
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

	})

	Describe("Window Size Updates", func() {
		BeforeEach(func() {
			screen = facts.NewFactListScreen(factList)
		})

		It("should handle window size changes", func() {
			msg := tea.WindowSizeMsg{Width: 100, Height: 30}
			cmd, result := screen.Update(msg)

			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
		})
	})

	Describe("SetTheme", func() {
		BeforeEach(func() {
			screen = facts.NewFactListScreen(factList)
		})

		It("should apply a valid theme", func() {
			theme := themes.NewDefaultTheme()
			screen.SetTheme(theme)
			content := screen.RenderContent()
			Expect(content).NotTo(BeEmpty())
		})

		It("should handle non-theme value gracefully", func() {
			screen.SetTheme("not-a-theme")
			content := screen.RenderContent()
			Expect(content).NotTo(BeEmpty())
		})
	})

})

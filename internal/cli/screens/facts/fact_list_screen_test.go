package facts_test

import (
	"testing"
	"time"

	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/facts"
	"github.com/baphled/kariya/internal/domain/career"
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
		factList = []*career.Fact{
			{
				ID:        "fact-1",
				Text:      "Led development of microservices architecture",
				CreatedAt: createdAt,
				UpdatedAt: createdAt,
			},
			{
				ID:        "fact-2",
				Text:      "Reduced API response time by 40%",
				CreatedAt: createdAt.Add(24 * time.Hour),
				UpdatedAt: createdAt.Add(24 * time.Hour),
			},
			{
				ID:        "fact-3",
				Text:      "Mentored 5 junior engineers",
				CreatedAt: createdAt.Add(48 * time.Hour),
				UpdatedAt: createdAt.Add(48 * time.Hour),
			},
		}
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
})

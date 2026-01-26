package modals_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/screens/timeline/modals"
	tea "github.com/charmbracelet/bubbletea"
)

var _ = Describe("SearchModal", func() {
	var modal *modals.SearchModal

	Describe("NewSearchModal", func() {
		It("creates a modal with empty search text", func() {
			modal = modals.NewSearchModal("", 80, 24)

			Expect(modal).NotTo(BeNil())
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("creates a modal with existing search text", func() {
			modal = modals.NewSearchModal("existing search", 80, 24)

			Expect(modal).NotTo(BeNil())
		})

		It("respects minimum width", func() {
			modal = modals.NewSearchModal("", 20, 24)

			Expect(modal).NotTo(BeNil())
		})

		It("respects maximum width", func() {
			modal = modals.NewSearchModal("", 200, 24)

			Expect(modal).NotTo(BeNil())
		})
	})

	Describe("Init", func() {
		BeforeEach(func() {
			modal = modals.NewSearchModal("", 80, 24)
		})

		It("makes modal visible", func() {
			modal.Init()

			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("returns a command", func() {
			cmd := modal.Init()

			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("Update", func() {
		BeforeEach(func() {
			modal = modals.NewSearchModal("", 80, 24)
			modal.Init()
		})

		It("closes on escape key without applying", func() {
			escMsg := tea.KeyMsg{Type: tea.KeyEsc}

			cmd, applied, data := modal.Update(escMsg)

			Expect(modal.IsVisible()).To(BeFalse())
			Expect(applied).To(BeFalse())
			Expect(data).To(BeNil())
			Expect(cmd).To(BeNil())
		})

		It("handles window resize", func() {
			resizeMsg := tea.WindowSizeMsg{Width: 100, Height: 30}

			cmd, applied, data := modal.Update(resizeMsg)

			Expect(modal.IsVisible()).To(BeTrue())
			Expect(applied).To(BeFalse())
			Expect(data).To(BeNil())
			Expect(cmd).NotTo(BeNil())
		})

		It("forwards other key messages to form", func() {
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}

			cmd, applied, data := modal.Update(keyMsg)

			Expect(modal.IsVisible()).To(BeTrue())
			Expect(applied).To(BeFalse())
			Expect(data).To(BeNil())
			_ = cmd
		})
	})

	Describe("View", func() {
		BeforeEach(func() {
			modal = modals.NewSearchModal("", 80, 24)
		})

		It("returns empty string when not visible", func() {
			Expect(modal.IsVisible()).To(BeFalse())

			view := modal.View()

			Expect(view).To(BeEmpty())
		})

		It("returns content when visible", func() {
			modal.Init()

			view := modal.View()

			Expect(view).NotTo(BeEmpty())
			Expect(view).To(ContainSubstring("Search"))
		})
	})

	Describe("Visibility methods", func() {
		BeforeEach(func() {
			modal = modals.NewSearchModal("", 80, 24)
		})

		It("Show makes modal visible", func() {
			Expect(modal.IsVisible()).To(BeFalse())

			modal.Show()

			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("Hide makes modal invisible", func() {
			modal.Init()
			Expect(modal.IsVisible()).To(BeTrue())

			modal.Hide()

			Expect(modal.IsVisible()).To(BeFalse())
		})
	})

	Describe("RenderOverlay", func() {
		BeforeEach(func() {
			modal = modals.NewSearchModal("", 80, 24)
		})

		It("returns base view when not visible", func() {
			baseView := "Base View Content"

			result := modal.RenderOverlay(baseView)

			Expect(result).To(Equal(baseView))
		})

		It("renders overlay when visible", func() {
			modal.Init()
			baseView := "Base View Content"

			result := modal.RenderOverlay(baseView)

			Expect(result).NotTo(Equal(baseView))
			Expect(result).To(ContainSubstring("Search"))
		})
	})
})

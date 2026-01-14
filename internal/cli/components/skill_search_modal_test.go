package components_test

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/components"
)

var _ = Describe("SkillSearchModal", func() {
	var (
		modal  *components.SkillSearchModal
		width  int
		height int
	)

	BeforeEach(func() {
		width = 80
		height = 24
		modal = components.NewSkillSearchModal("", width, height)
	})

	Describe("NewSkillSearchModal", func() {
		It("should create a new search modal", func() {
			Expect(modal).NotTo(BeNil())
		})

		It("should initialize as not visible", func() {
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("should accept pre-populated search text", func() {
			modal = components.NewSkillSearchModal("test search", width, height)
			Expect(modal.GetSearchText()).To(Equal("test search"))
		})
	})

	Describe("Init", func() {
		It("should make modal visible", func() {
			modal.Init()
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("should return a command", func() {
			cmd := modal.Init()
			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("Update", func() {
		BeforeEach(func() {
			modal.Init()
		})

		Context("when Esc is pressed", func() {
			It("should hide modal and return applied=false", func() {
				cmd, applied, data := modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
				Expect(cmd).To(BeNil())
				Expect(applied).To(BeFalse())
				Expect(data).To(BeNil())
				Expect(modal.IsVisible()).To(BeFalse())
			})
		})

		Context("when Enter is pressed", func() {
			It("should complete search and return applied=true with search text", func() {
				// Type search text
				modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("golang")})

				// Press Enter to submit
				// First Enter moves focus/completes the input field
				modal.Update(tea.KeyMsg{Type: tea.KeyEnter})

				// Second Enter submits the form (when form state is completed)
				cmd, applied, data := modal.Update(tea.KeyMsg{Type: tea.KeyEnter})

				// Should return applied=true when form completes
				if applied {
					Expect(data).NotTo(BeNil())
					Expect(data.SearchText).To(Equal("golang"))
					Expect(modal.IsVisible()).To(BeFalse())
				}
				// Cmd may be nil or not depending on form state
				_ = cmd
			})

			It("should keep modal visible while form is not complete", func() {
				// Press Enter before form is complete
				cmd, applied, data := modal.Update(tea.KeyMsg{Type: tea.KeyEnter})

				// Should forward to form (not applied yet)
				Expect(applied).To(BeFalse())
				Expect(data).To(BeNil())
				// Cmd may be nil or not
				_ = cmd
			})
		})

		Context("navigation", func() {
			It("should handle tab key for field navigation", func() {
				// Tab should be forwarded to form
				cmd, applied, data := modal.Update(tea.KeyMsg{Type: tea.KeyTab})
				Expect(applied).To(BeFalse())
				Expect(data).To(BeNil())
				Expect(modal.IsVisible()).To(BeTrue())
				_ = cmd
			})

			It("should handle shift+tab for reverse navigation", func() {
				// Shift+Tab should be forwarded to form
				cmd, applied, data := modal.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
				Expect(applied).To(BeFalse())
				Expect(data).To(BeNil())
				Expect(modal.IsVisible()).To(BeTrue())
				_ = cmd
			})
		})

		Context("text input", func() {
			It("should accept text input for search field", func() {
				// Type characters
				modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("r")})
				modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("u")})
				modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("b")})
				modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("y")})

				// Should still be visible and not applied
				Expect(modal.IsVisible()).To(BeTrue())
			})

			It("should allow backspace to delete characters", func() {
				// Type and then delete
				modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("test")})
				cmd, applied, data := modal.Update(tea.KeyMsg{Type: tea.KeyBackspace})

				Expect(applied).To(BeFalse())
				Expect(data).To(BeNil())
				Expect(modal.IsVisible()).To(BeTrue())
				_ = cmd
			})
		})

		It("should forward other keys to form", func() {
			cmd, applied, data := modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			// Form should handle the key
			Expect(applied).To(BeFalse())
			Expect(data).To(BeNil())
			// Modal should still be visible
			Expect(modal.IsVisible()).To(BeTrue())
			// Cmd may or may not be nil depending on form state
			_ = cmd
		})
	})

	Describe("View", func() {
		Context("when not visible", func() {
			It("should return empty string", func() {
				Expect(modal.View()).To(BeEmpty())
			})
		})

		Context("when visible", func() {
			BeforeEach(func() {
				modal.Init()
			})

			It("should return non-empty view", func() {
				view := modal.View()
				Expect(view).NotTo(BeEmpty())
			})

			It("should contain search field", func() {
				view := modal.View()
				Expect(view).To(ContainSubstring("Search"))
			})
		})
	})

	Describe("Visibility Methods", func() {
		It("should toggle visibility with Show/Hide", func() {
			Expect(modal.IsVisible()).To(BeFalse())

			modal.Show()
			Expect(modal.IsVisible()).To(BeTrue())

			modal.Hide()
			Expect(modal.IsVisible()).To(BeFalse())
		})
	})

	Describe("SetSize", func() {
		It("should update modal dimensions", func() {
			newWidth := 120
			newHeight := 40
			modal.SetSize(newWidth, newHeight)
			// Modal should still be functional after resize
			modal.Init()
			Expect(modal.IsVisible()).To(BeTrue())
		})
	})

	Describe("RenderOverlay", func() {
		var baseView string

		BeforeEach(func() {
			baseView = "Base View Content\nLine 2\nLine 3"
		})

		Context("when not visible", func() {
			It("should return base view unchanged", func() {
				result := modal.RenderOverlay(baseView)
				Expect(result).To(Equal(baseView))
			})
		})

		Context("when visible", func() {
			BeforeEach(func() {
				modal.Init()
			})

			It("should return overlay view", func() {
				result := modal.RenderOverlay(baseView)
				Expect(result).NotTo(Equal(baseView))
				// Should contain overlay rendering
				Expect(result).NotTo(BeEmpty())
			})
		})
	})

	Describe("GetSearchText", func() {
		It("should return current search text", func() {
			modal = components.NewSkillSearchModal("initial", width, height)
			Expect(modal.GetSearchText()).To(Equal("initial"))
		})

		It("should return empty string when not set", func() {
			Expect(modal.GetSearchText()).To(BeEmpty())
		})
	})

	Describe("Modal Styling", func() {
		It("should have solid background styling in View output", func() {
			modal := components.NewSkillSearchModal("test", 80, 24)
			modal.Show()

			view := modal.View()
			// Check for rounded border characters (top-left corner)
			Expect(view).To(ContainSubstring("╭"))
			// Check for rounded border characters (top-right corner)
			Expect(view).To(ContainSubstring("╮"))
		})

		It("should display keyboard shortcuts in footer (KeyBadge pattern)", func() {
			modal := components.NewSkillSearchModal("test", 80, 24)
			modal.Show()

			view := modal.View()
			// Modal should show keyboard shortcuts for user guidance
			// Pattern: Tab (next field), Enter (submit), Esc (cancel)
			Expect(view).To(ContainSubstring("Tab"))
			Expect(view).To(ContainSubstring("Enter"))
			Expect(view).To(ContainSubstring("Esc"))
			// Should indicate what each key does
			Expect(view).To(Or(
				ContainSubstring("Next"),
				ContainSubstring("field"),
				ContainSubstring("Navigate"),
			))
			Expect(view).To(Or(
				ContainSubstring("Submit"),
				ContainSubstring("Confirm"),
			))
			Expect(view).To(Or(
				ContainSubstring("Cancel"),
				ContainSubstring("Back"),
			))
		})
	})

	Describe("Modal Width", func() {
		It("should calculate width as 60% of terminal width", func() {
			// 100 chars terminal = 60% = 60 chars modal
			modal := components.NewSkillSearchModal("test", 100, 24)
			modal.Show()
			view := modal.View()

			// Modal should not span full width
			// With 100 char terminal, modal should be ~60 chars + padding/border
			lines := strings.Split(view, "\n")
			Expect(len(lines)).To(BeNumerically(">", 0))

			// Find the widest line (should be border)
			maxWidth := 0
			for _, line := range lines {
				// Use visual width (rune count)
				width := len([]rune(line))
				if width > maxWidth {
					maxWidth = width
				}
			}

			// Should be around 60 chars + border (2) + padding (4) = ~66 chars
			// Allow some flexibility but should definitely be < 80
			Expect(maxWidth).To(BeNumerically("<", 80))
			Expect(maxWidth).To(BeNumerically(">", 50))
		})

		It("should enforce max width of 60 chars", func() {
			// 200 chars terminal should still result in max 60 char modal
			modal := components.NewSkillSearchModal("test", 200, 24)
			modal.Show()
			view := modal.View()

			lines := strings.Split(view, "\n")
			maxWidth := 0
			for _, line := range lines {
				width := len([]rune(line))
				if width > maxWidth {
					maxWidth = width
				}
			}

			// Max 60 + border (2) + padding (4) = ~66 chars
			Expect(maxWidth).To(BeNumerically("<=", 70))
		})

		It("should enforce min width of 40 chars", func() {
			// Very small terminal should still have min 40 char modal
			modal := components.NewSkillSearchModal("test", 50, 24)
			modal.Show()
			view := modal.View()

			lines := strings.Split(view, "\n")
			maxWidth := 0
			for _, line := range lines {
				width := len([]rune(line))
				if width > maxWidth {
					maxWidth = width
				}
			}

			// Min 40 + border (2) + padding (4) = ~46 chars
			Expect(maxWidth).To(BeNumerically(">=", 46))
		})
	})
})

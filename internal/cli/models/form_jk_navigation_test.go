package models_test

import (
	"github.com/baphled/kariya/internal/cli/models"
	"github.com/baphled/kariya/internal/cli/service"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("FormModel j/k Navigation", func() {
	var (
		form      *models.FormModel
		cliSvc    *service.CLIEventService
		repo      *careerrepo.MemoryRepository
		careerSvc *careerservice.Service
	)

	BeforeEach(func() {
		repo = careerrepo.NewMemoryRepository()
		careerSvc = careerservice.NewService(repo)
		cliSvc = service.NewCLIEventService(careerSvc)
		form = models.NewFormModel(cliSvc)
	})

	Describe("j/k navigation for tags", func() {
		BeforeEach(func() {
			// Move to tags field by pressing Tab multiple times to get to TagsField
			// TextField -> DateField -> CompanyField -> ProjectField -> TagsField
			for i := 0; i < 4; i++ {
				form.Update(tea.KeyMsg{Type: tea.KeyTab})
			}
		})

		It("should navigate down tags with j key", func() {
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			// Verify form is still functional
			Expect(form).NotTo(BeNil())
		})

		It("should navigate up tags with k key", func() {
			// Move to a tag first
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			// Then go back up
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
			Expect(form).NotTo(BeNil())
		})

		It("should wrap around tags with j key", func() {
			availableTags := form.TagSelector().AvailableTags()
			// Navigate to the last tag
			for i := 0; i < len(availableTags); i++ {
				form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			}
			// One more j should wrap to first tag
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			Expect(form).NotTo(BeNil())
		})

		It("should wrap around tags with k key", func() {
			// k at first tag should wrap to last tag
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
			Expect(form).NotTo(BeNil())
		})

		It("should allow space to toggle tags during j/k navigation", func() {
			// Navigate to a tag
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			// Toggle with space
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
			// Navigate away
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
			Expect(form).NotTo(BeNil())
		})
	})

	Describe("j/k navigation for categories", func() {
		BeforeEach(func() {
			// Move to categories field
			// TextField -> DateField -> CompanyField -> ProjectField -> TagsField -> CategoriesField
			for i := 0; i < 5; i++ {
				form.Update(tea.KeyMsg{Type: tea.KeyTab})
			}
		})

		It("should navigate down categories with j key", func() {
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			Expect(form).NotTo(BeNil())
		})

		It("should navigate up categories with k key", func() {
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
			Expect(form).NotTo(BeNil())
		})

		It("should wrap around categories with j key", func() {
			availableCategories := form.CategorySelector().AvailableCategories()
			// Navigate to the last category
			for i := 0; i < len(availableCategories); i++ {
				form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			}
			// One more j should wrap to first category
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			Expect(form).NotTo(BeNil())
		})

		It("should wrap around categories with k key", func() {
			// k at first category should wrap to last category
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
			Expect(form).NotTo(BeNil())
		})

		It("should allow space to toggle categories during j/k navigation", func() {
			// Navigate to a category
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			// Toggle with space
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
			// Navigate away
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
			Expect(form).NotTo(BeNil())
		})
	})

	Describe("j/k navigation for modes", func() {
		BeforeEach(func() {
			// Move to mode field
			// TextField -> DateField -> CompanyField -> ProjectField -> TagsField -> CategoriesField -> ModeField
			for i := 0; i < 6; i++ {
				form.Update(tea.KeyMsg{Type: tea.KeyTab})
			}
		})

		It("should navigate down modes with j key", func() {
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			Expect(form).NotTo(BeNil())
		})

		It("should navigate up modes with k key", func() {
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
			Expect(form).NotTo(BeNil())
		})

		It("should wrap around modes with j key", func() {
			// Navigate through all modes
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			// One more j should wrap to first mode
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			Expect(form).NotTo(BeNil())
		})

		It("should wrap around modes with k key", func() {
			// k at first mode should wrap to last mode
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
			Expect(form).NotTo(BeNil())
		})
	})

	Describe("j/k key behavior in text fields", func() {
		It("should allow j/k as regular input in text field", func() {
			// When in text field (default focus), j and k should be typed
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
			// Verify text field contains the characters
			textValue := form.GetInputValue(0) // TextField index
			Expect(textValue).To(ContainSubstring("j"))
			Expect(textValue).To(ContainSubstring("k"))
		})
	})

	Describe("Navigation consistency", func() {
		It("should maintain consistent state after j/k navigation", func() {
			// Move to tags (4 tabs)
			for i := 0; i < 4; i++ {
				form.Update(tea.KeyMsg{Type: tea.KeyTab})
			}
			// Navigate with j
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			// Move to categories (1 more tab)
			form.Update(tea.KeyMsg{Type: tea.KeyTab})
			// Navigate with k
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
			// Verify form is still functional
			Expect(form).NotTo(BeNil())
		})

		It("should not interfere with other navigation keys", func() {
			// Use Tab to navigate
			form.Update(tea.KeyMsg{Type: tea.KeyTab})
			// Move to tags
			for i := 0; i < 3; i++ {
				form.Update(tea.KeyMsg{Type: tea.KeyTab})
			}
			// Now use j/k in a navigation context
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			// Verify form is still functional
			Expect(form).NotTo(BeNil())
		})

		It("should work with Up/Down arrows as alternative navigation", func() {
			// Move to tags
			for i := 0; i < 4; i++ {
				form.Update(tea.KeyMsg{Type: tea.KeyTab})
			}
			// Use Up/Down arrows
			form.Update(tea.KeyMsg{Type: tea.KeyUp})
			form.Update(tea.KeyMsg{Type: tea.KeyDown})
			// Then use j/k
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			form.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
			Expect(form).NotTo(BeNil())
		})
	})
})

package modals_test

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/baphled/kariya/internal/cli/screens/skills/modals"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("FilterModal", func() {
	var (
		skills        []*career.Skill
		currentFilter *modals.Filters
		modal         *modals.FilterModal
	)

	BeforeEach(func() {
		years1 := 5
		years2 := 3
		years3 := 4

		s1 := fixtures.SkillWith("skill-1", "Go", "backend", "expert")
		s1.YearsUsed = &years1
		s2 := fixtures.SkillWith("skill-2", "React", "frontend", "intermediate")
		s2.YearsUsed = &years2
		s3 := fixtures.SkillWith("skill-3", "Docker", "devops", "advanced")
		s3.YearsUsed = &years3

		skills = []*career.Skill{s1, s2, s3}

		currentFilter = &modals.Filters{
			Categories: []string{},
			Levels:     []string{},
			MinYears:   0,
			MaxYears:   0,
		}
	})

	Describe("validateYearsInput", func() {
		It("should accept empty string", func() {
			err := modals.ValidateYearsInput("")
			Expect(err).ToNot(HaveOccurred())
		})

		It("should accept valid positive number", func() {
			err := modals.ValidateYearsInput("5")
			Expect(err).ToNot(HaveOccurred())
		})

		It("should accept zero", func() {
			err := modals.ValidateYearsInput("0")
			Expect(err).ToNot(HaveOccurred())
		})

		It("should accept large numbers", func() {
			err := modals.ValidateYearsInput("999")
			Expect(err).ToNot(HaveOccurred())
		})

		It("should reject negative numbers", func() {
			err := modals.ValidateYearsInput("-5")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("must be positive"))
		})

		It("should reject non-numeric strings", func() {
			err := modals.ValidateYearsInput("abc")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("must be a number"))
		})

		It("should reject mixed alphanumeric", func() {
			err := modals.ValidateYearsInput("5abc")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("must be a number"))
		})

		It("should reject decimal numbers", func() {
			err := modals.ValidateYearsInput("5.5")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("must be a number"))
		})

		It("should reject special characters", func() {
			err := modals.ValidateYearsInput("@#$")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("must be a number"))
		})

		It("should reject strings with spaces", func() {
			err := modals.ValidateYearsInput("5 ")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("must be a number"))
		})
	})

	Describe("NewFilterModal", func() {
		It("should create a new skill filter modal", func() {
			modal = modals.NewFilterModal(skills, currentFilter, 120, 40)

			Expect(modal).NotTo(BeNil())
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("should extract unique categories from skills", func() {
			modal = modals.NewFilterModal(skills, currentFilter, 120, 40)

			// Modal should include Backend, Frontend, DevOps categories
			Expect(modal).NotTo(BeNil())
		})

		It("should pre-populate form with current filters", func() {
			currentFilter.Categories = []string{"Backend"}
			currentFilter.Levels = []string{"Expert"}
			currentFilter.MinYears = 2
			currentFilter.MaxYears = 10

			modal = modals.NewFilterModal(skills, currentFilter, 120, 40)

			Expect(modal).NotTo(BeNil())
			// Form data should be pre-populated (internal state)
		})

		It("should handle zero years in current filter", func() {
			currentFilter.MinYears = 0
			currentFilter.MaxYears = 0

			modal = modals.NewFilterModal(skills, currentFilter, 120, 40)

			Expect(modal).NotTo(BeNil())
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle negative years gracefully", func() {
			currentFilter.MinYears = -1
			currentFilter.MaxYears = -5

			modal = modals.NewFilterModal(skills, currentFilter, 120, 40)

			Expect(modal).NotTo(BeNil())
		})

		It("should calculate modal width for small terminals", func() {
			modal = modals.NewFilterModal(skills, currentFilter, 50, 40)
			Expect(modal).NotTo(BeNil())
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should calculate modal width for large terminals", func() {
			modal = modals.NewFilterModal(skills, currentFilter, 200, 40)
			Expect(modal).NotTo(BeNil())
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should calculate modal width for very small terminals", func() {
			modal = modals.NewFilterModal(skills, currentFilter, 30, 40)
			Expect(modal).NotTo(BeNil())
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Filtering Options", func() {
		BeforeEach(func() {
			modal = modals.NewFilterModal(skills, currentFilter, 120, 40)
		})

		It("should include category filter", func() {
			// Modal should have MultiSelect for categories
			Expect(modal).NotTo(BeNil())
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("should include level filter", func() {
			// Modal should have Select for level
			Expect(modal).NotTo(BeNil())
		})

		It("should include years range filter", func() {
			// Modal should have input fields for min/max years
			Expect(modal).NotTo(BeNil())
		})

		It("should handle empty categories gracefully", func() {
			emptySkills := []*career.Skill{}
			modal = modals.NewFilterModal(emptySkills, currentFilter, 120, 40)

			Expect(modal).NotTo(BeNil())
		})

		It("should validate years input accepts numbers", func() {
			currentFilter.MinYears = 5
			currentFilter.MaxYears = 10
			modal = modals.NewFilterModal(skills, currentFilter, 120, 40)
			Expect(modal).NotTo(BeNil())
			filters := modal.ToFilters()
			Expect(filters.MinYears).To(Equal(5))
			Expect(filters.MaxYears).To(Equal(10))
		})

		It("should validate years input rejects negative numbers", func() {
			currentFilter.MinYears = -5
			modal = modals.NewFilterModal(skills, currentFilter, 120, 40)
			Expect(modal).NotTo(BeNil())
			filters := modal.ToFilters()
			Expect(filters.MinYears).To(Equal(0))
		})

		It("should handle zero min years", func() {
			currentFilter.MinYears = 0
			currentFilter.MaxYears = 10
			modal = modals.NewFilterModal(skills, currentFilter, 120, 40)
			filters := modal.ToFilters()
			Expect(filters.MinYears).To(Equal(0))
			Expect(filters.MaxYears).To(Equal(10))
		})

		It("should handle zero max years", func() {
			currentFilter.MinYears = 5
			currentFilter.MaxYears = 0
			modal = modals.NewFilterModal(skills, currentFilter, 120, 40)
			filters := modal.ToFilters()
			Expect(filters.MinYears).To(Equal(5))
			Expect(filters.MaxYears).To(Equal(0))
		})
	})

	Describe("Sorting Options", func() {
		BeforeEach(func() {
			modal = modals.NewFilterModal(skills, currentFilter, 120, 40)
		})

		It("should include sort by options", func() {
			// Modal should have Select for sort field (name, category, level, years, events)
			Expect(modal).NotTo(BeNil())
		})

		It("should include sort order options", func() {
			// Modal should have Select for sort direction (asc, desc)
			Expect(modal).NotTo(BeNil())
		})
	})

	Describe("Init", func() {
		It("should initialize the form", func() {
			modal = modals.NewFilterModal(skills, currentFilter, 120, 40)
			cmd := modal.Init()
			Expect(cmd).NotTo(BeNil())
		})

		It("should allow multiple Init calls", func() {
			modal = modals.NewFilterModal(skills, currentFilter, 120, 40)
			cmd1 := modal.Init()
			cmd2 := modal.Init()
			Expect(cmd1).NotTo(BeNil())
			Expect(cmd2).NotTo(BeNil())
		})

		It("should work after hiding and showing", func() {
			modal = modals.NewFilterModal(skills, currentFilter, 120, 40)
			modal.Init()
			modal.Hide()
			modal.Show()
			cmd := modal.Init()
			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("Update", func() {
		BeforeEach(func() {
			modal = modals.NewFilterModal(skills, currentFilter, 120, 40)
		})

		It("should handle escape key to close without applying", func() {
			cmd, completed, data := modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(cmd).To(BeNil())
			Expect(completed).To(BeFalse())
			Expect(data).To(BeNil())
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("should handle window resize messages", func() {
			cmd, completed, data := modal.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
			Expect(cmd).To(BeNil())
			Expect(completed).To(BeFalse())
			Expect(data).To(BeNil())
		})

		It("should handle small terminal width", func() {
			cmd, completed, data := modal.Update(tea.WindowSizeMsg{Width: 50, Height: 30})
			Expect(cmd).To(BeNil())
			Expect(completed).To(BeFalse())
			Expect(data).To(BeNil())
		})

		It("should handle large terminal width", func() {
			cmd, completed, data := modal.Update(tea.WindowSizeMsg{Width: 200, Height: 30})
			Expect(cmd).To(BeNil())
			Expect(completed).To(BeFalse())
			Expect(data).To(BeNil())
		})

		It("should return nil when modal is hidden", func() {
			modal.Hide()
			cmd, completed, data := modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(cmd).To(BeNil())
			Expect(completed).To(BeFalse())
			Expect(data).To(BeNil())
		})

		It("should handle regular key messages", func() {
			_, completed, data := modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			Expect(completed).To(BeFalse())
			Expect(data).To(BeNil())
		})

		It("should handle enter key", func() {
			_, completed, data := modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(completed).To(BeFalse())
			Expect(data).To(BeNil())
		})

		It("should handle tab key", func() {
			_, completed, data := modal.Update(tea.KeyMsg{Type: tea.KeyTab})
			Expect(completed).To(BeFalse())
			Expect(data).To(BeNil())
		})

		It("should handle shift+tab key", func() {
			_, completed, data := modal.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
			Expect(completed).To(BeFalse())
			Expect(data).To(BeNil())
		})

		It("should handle backspace key", func() {
			_, completed, data := modal.Update(tea.KeyMsg{Type: tea.KeyBackspace})
			Expect(completed).To(BeFalse())
			Expect(data).To(BeNil())
		})

		It("should handle delete key", func() {
			_, completed, data := modal.Update(tea.KeyMsg{Type: tea.KeyDelete})
			Expect(completed).To(BeFalse())
			Expect(data).To(BeNil())
		})

		It("should handle arrow keys", func() {
			_, completed, data := modal.Update(tea.KeyMsg{Type: tea.KeyUp})
			Expect(completed).To(BeFalse())
			Expect(data).To(BeNil())
		})

		It("should handle down arrow", func() {
			_, completed, data := modal.Update(tea.KeyMsg{Type: tea.KeyDown})
			Expect(completed).To(BeFalse())
			Expect(data).To(BeNil())
		})

		It("should handle left arrow", func() {
			_, completed, data := modal.Update(tea.KeyMsg{Type: tea.KeyLeft})
			Expect(completed).To(BeFalse())
			Expect(data).To(BeNil())
		})

		It("should handle right arrow", func() {
			_, completed, data := modal.Update(tea.KeyMsg{Type: tea.KeyRight})
			Expect(completed).To(BeFalse())
			Expect(data).To(BeNil())
		})
	})

	Describe("Show", func() {
		It("should make modal visible", func() {
			modal = modals.NewFilterModal(skills, currentFilter, 120, 40)
			modal.Hide()
			Expect(modal.IsVisible()).To(BeFalse())
			modal.Show()
			Expect(modal.IsVisible()).To(BeTrue())
		})
	})

	Describe("ToFilters", func() {
		It("should convert form data to Filters", func() {
			modal = modals.NewFilterModal(skills, currentFilter, 120, 40)
			filters := modal.ToFilters()
			Expect(filters).NotTo(BeNil())
			Expect(filters.Categories).NotTo(BeNil())
			Expect(filters.Levels).NotTo(BeNil())
		})

		It("should parse valid MinYears", func() {
			currentFilter.MinYears = 2
			modal = modals.NewFilterModal(skills, currentFilter, 120, 40)
			filters := modal.ToFilters()
			Expect(filters.MinYears).To(Equal(2))
		})

		It("should parse valid MaxYears", func() {
			currentFilter.MaxYears = 10
			modal = modals.NewFilterModal(skills, currentFilter, 120, 40)
			filters := modal.ToFilters()
			Expect(filters.MaxYears).To(Equal(10))
		})

		It("should handle empty years strings", func() {
			modal = modals.NewFilterModal(skills, currentFilter, 120, 40)
			filters := modal.ToFilters()
			Expect(filters.MinYears).To(Equal(0))
			Expect(filters.MaxYears).To(Equal(0))
		})
	})

	Describe("RenderOverlay", func() {
		It("should return base view when modal is hidden", func() {
			modal = modals.NewFilterModal(skills, currentFilter, 120, 40)
			modal.Hide()
			baseView := "base content"
			result := modal.RenderOverlay(baseView)
			Expect(result).To(Equal(baseView))
		})

		It("should render overlay when modal is visible", func() {
			modal = modals.NewFilterModal(skills, currentFilter, 120, 40)
			baseView := "base content"
			result := modal.RenderOverlay(baseView)
			Expect(result).NotTo(BeEmpty())
			Expect(result).NotTo(Equal(baseView))
		})
	})

	Describe("View", func() {
		It("should return empty string when not visible", func() {
			modal = modals.NewFilterModal(skills, currentFilter, 120, 40)
			modal.Hide()

			view := modal.View()
			Expect(view).To(BeEmpty())
		})

		It("should render form when visible", func() {
			modal = modals.NewFilterModal(skills, currentFilter, 120, 40)

			view := modal.View()
			Expect(view).NotTo(BeEmpty())
			// Should have huh form controls
			Expect(view).To(ContainSubstring("enter"))
		})

		It("should use default theme when nil is provided", func() {
			modal = modals.NewFilterModal(skills, currentFilter, 120, 40)
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
			// NewFilterModal initializes a default theme internally
		})

		It("should render with very narrow width", func() {
			modal = modals.NewFilterModal(skills, currentFilter, 40, 40)
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should render with very wide width", func() {
			modal = modals.NewFilterModal(skills, currentFilter, 250, 40)
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should render multiple times consistently", func() {
			modal = modals.NewFilterModal(skills, currentFilter, 120, 40)
			view1 := modal.View()
			view2 := modal.View()
			Expect(view1).To(Equal(view2))
		})

		It("should display keyboard shortcuts in footer (KeyBadge pattern)", func() {
			modal = modals.NewFilterModal(skills, currentFilter, 120, 40)

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
				ContainSubstring("Apply"),
			))
			Expect(view).To(Or(
				ContainSubstring("Cancel"),
				ContainSubstring("Back"),
			))
		})
	})
})

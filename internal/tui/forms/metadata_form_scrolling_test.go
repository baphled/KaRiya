package forms_test

import (
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/baphled/kariya/internal/tui/forms"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Metadata Form Scrolling", func() {
	var (
		event           *career.Event
		availableTags   []string
		availableCats   []string
		availableSkills []*career.Skill
	)

	BeforeEach(func() {
		event = fixtures.EventWith("test-event", "Built REST API with Go and PostgreSQL", "Test Company", "Test Project")
		event.Tags = []string{"backend", "go"}
		event.Categories = []string{"development"}
		event.Skills = []string{}

		availableTags = []string{"backend", "frontend", "database", "api"}
		availableCats = []string{"development", "testing", "deployment"}
		availableSkills = []*career.Skill{
			fixtures.SkillWith("skill-1", "Go", "backend", "advanced"),
			fixtures.SkillWith("skill-2", "PostgreSQL", "backend", "intermediate"),
			fixtures.SkillWith("skill-3", "Redis", "backend", "intermediate"),
			fixtures.SkillWith("skill-4", "Docker", "devops", "intermediate"),
		}
	})

	Describe("Form with constrained height", func() {
		It("should show all field titles including Skills", func() {
			// Create form with constrained height (simulating modal overlay)
			data := forms.GetMetadataFormData(event)
			form := forms.NewMetadataForm(data, forms.MetadataFormConfig{
				AvailableTags:       availableTags,
				AvailableCategories: availableCats,
				AvailableSkills:     availableSkills,
			})

			// Initialize the form
			cmd := form.Init()
			Expect(cmd).NotTo(BeNil())

			// Render the form view
			view := form.View()

			// Debug output
			GinkgoWriter.Printf("Form view:\n%s\n", view)
			GinkgoWriter.Printf("View length: %d\n", len(view))

			// All field titles should be visible
			Expect(view).To(ContainSubstring("Date"))
			Expect(view).To(ContainSubstring("Company"))
			Expect(view).To(ContainSubstring("Project"))
			Expect(view).To(ContainSubstring("Tags"))
			Expect(view).To(ContainSubstring("Categories"))

			// This is the critical assertion - Skills should be visible
			Expect(view).To(ContainSubstring("Skills"),
				"Skills field should be visible in the metadata form. "+
					"If this fails, the form is constrained and not showing all fields.")
		})

		It("should render a scrollable viewport with constrained dimensions", func() {
			data := forms.GetMetadataFormData(event)
			form := forms.NewMetadataForm(
				data,
				forms.MetadataFormConfig{
					AvailableTags:       availableTags,
					AvailableCategories: availableCats,
					AvailableSkills:     availableSkills,
					Width:               74, // ModalFormWidth(80)
					Height:              18, // ModalFormHeight(40) = 40-20-2
				},
			)

			cmd := form.Init()
			Expect(cmd).NotTo(BeNil())

			view := form.View()

			GinkgoWriter.Printf("Constrained form view:\n%s\n", view)
			GinkgoWriter.Printf("View length: %d\n", len(view))

			Expect(view).To(ContainSubstring("Date"))
			Expect(view).To(ContainSubstring("Company"))
			Expect(view).To(ContainSubstring("Project"))
			Expect(view).To(ContainSubstring("Tags"))

			Expect(data.Tags).NotTo(BeNil())
			Expect(data.Categories).NotTo(BeNil())
			Expect(data.Skills).NotTo(BeNil())
		})
	})

	Describe("Form field count", func() {
		It("should have 6 fields total", func() {
			// The metadata form should have:
			// 1. Date
			// 2. Company
			// 3. Project
			// 4. Tags
			// 5. Categories
			// 6. Skills

			data := forms.GetMetadataFormData(event)
			Expect(data.Date).NotTo(BeNil())
			Expect(data.Company).NotTo(BeNil())
			Expect(data.Project).NotTo(BeNil())
			Expect(data.Tags).NotTo(BeNil())
			Expect(data.Categories).NotTo(BeNil())
			Expect(data.Skills).NotTo(BeNil())
		})
	})
})

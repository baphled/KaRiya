package models_test

import (
	"strings"

	"github.com/baphled/kariya/internal/cli/models"
	"github.com/baphled/kariya/internal/cli/service"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Form Model - Label Styling Consistency", func() {
	var (
		form *models.FormModel
	)

	BeforeEach(func() {
		repo := careerrepo.NewMemoryRepository()
		svc := careerservice.NewService(repo)
		cliSvc := service.NewCLIEventService(svc)
		form = models.NewFormModel(cliSvc)
	})

	Context("Label Rendering", func() {
		It("should render all field labels with consistent styling", func() {
			view := form.View()
			Expect(view).To(ContainSubstring("Event Text (required):"))
			Expect(view).To(ContainSubstring("Date (optional):"))
			Expect(view).To(ContainSubstring("Company (optional):"))
			Expect(view).To(ContainSubstring("Project (optional):"))
			Expect(view).To(ContainSubstring("Tags"))
			Expect(view).To(ContainSubstring("Categories"))
			Expect(view).To(ContainSubstring("Capture Mode:"))
		})

		It("should have consistent label presence in view output", func() {
			view := form.View()
			lines := strings.Split(view, "\n")
			labelCount := 0
			for _, line := range lines {
				if strings.Contains(line, ":") &&
					(strings.Contains(line, "Event") ||
						strings.Contains(line, "Date") ||
						strings.Contains(line, "Company") ||
						strings.Contains(line, "Project") ||
						strings.Contains(line, "Mode")) {
					labelCount++
				}
			}
			Expect(labelCount).To(BeNumerically(">", 0))
		})

		It("should render labels without errors", func() {
			view := form.View()
			Expect(view).NotTo(BeEmpty())
			Expect(len(view)).To(BeNumerically(">", 0))
		})

		It("should have consistent label formatting for all fields", func() {
			view := form.View()
			// All labels end with colon
			Expect(view).To(ContainSubstring("Event Text (required):"))
			Expect(view).To(ContainSubstring("Date (optional):"))
			Expect(view).To(ContainSubstring("Company (optional):"))
			Expect(view).To(ContainSubstring("Project (optional):"))
			Expect(view).To(ContainSubstring("Capture Mode:"))
		})
	})

	Context("Label Style Properties", func() {
		It("should display form without rendering errors", func() {
			// Verify the form can be rendered
			Expect(func() {
				_ = form.View()
			}).NotTo(Panic())
		})

		It("should maintain label visibility across multiple renders", func() {
			view1 := form.View()
			view2 := form.View()
			// Both renders should contain labels
			Expect(view1).To(ContainSubstring("Event Text"))
			Expect(view2).To(ContainSubstring("Event Text"))
		})
	})
})


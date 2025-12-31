package components_test

import (
	"github.com/baphled/kariya/internal/cli/components"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Spinner Component", func() {
	Context("when creating a spinner", func() {
		It("should create spinner with default style", func() {
			spinner := components.NewSpinner()
			Expect(spinner).NotTo(BeNil())
		})

		It("should have a View method that returns a string", func() {
			spinner := components.NewSpinner()
			view := spinner.View()
			Expect(view).To(BeAssignableToTypeOf(""))
		})
	})

	Context("when rendering the spinner", func() {
		It("should render spinner with styled output", func() {
			spinner := components.NewSpinner()
			view := spinner.View()
			Expect(view).NotTo(BeEmpty())
		})
	})
})

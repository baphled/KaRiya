package feedback_test

import (
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ModalContainer", func() {
	var container *feedback.ModalContainer

	BeforeEach(func() {
		container = feedback.NewModalContainer()
	})

	Describe("NewModalContainer", func() {
		It("creates a non-nil container", func() {
			Expect(container).NotTo(BeNil())
		})

		It("renders without error in default state", func() {
			output := container.Render()
			Expect(output).NotTo(BeEmpty())
		})
	})

	Describe("SetTitle", func() {
		It("returns the container for chaining", func() {
			result := container.SetTitle("Test Title")
			Expect(result).To(BeIdenticalTo(container))
		})

		It("includes the title in rendered output", func() {
			container.SetTitle("My Title")
			output := container.Render()
			Expect(output).To(ContainSubstring("My Title"))
		})
	})

	Describe("SetMessage", func() {
		It("returns the container for chaining", func() {
			result := container.SetMessage("Test message")
			Expect(result).To(BeIdenticalTo(container))
		})

		It("includes the message in rendered output", func() {
			container.SetMessage("Important message")
			output := container.Render()
			Expect(output).To(ContainSubstring("Important message"))
		})
	})

	Describe("SetButtons", func() {
		It("returns the container for chaining", func() {
			result := container.SetButtons([]string{"OK", "Cancel"})
			Expect(result).To(BeIdenticalTo(container))
		})

		It("includes buttons in rendered output", func() {
			container.SetButtons([]string{"Confirm", "Cancel"})
			output := container.Render()
			Expect(output).To(ContainSubstring("Confirm"))
			Expect(output).To(ContainSubstring("Cancel"))
		})
	})

	Describe("SetInstructions", func() {
		It("returns the container for chaining", func() {
			result := container.SetInstructions("Press Enter")
			Expect(result).To(BeIdenticalTo(container))
		})

		It("includes instructions in rendered output", func() {
			container.SetInstructions("Use arrow keys to navigate")
			output := container.Render()
			Expect(output).To(ContainSubstring("Use arrow keys to navigate"))
		})
	})

	Describe("WithDestructiveStyle", func() {
		It("returns the container for chaining", func() {
			result := container.WithDestructiveStyle()
			Expect(result).To(BeIdenticalTo(container))
		})

		It("renders with destructive styling applied", func() {
			output := container.
				SetTitle("Delete Item").
				SetMessage("This cannot be undone").
				WithDestructiveStyle().
				Render()
			Expect(output).To(ContainSubstring("Delete Item"))
			Expect(output).To(ContainSubstring("This cannot be undone"))
		})
	})

	Describe("WithTheme", func() {
		It("returns the container for chaining", func() {
			theme := themes.NewDefaultTheme()
			result := container.WithTheme(theme)
			Expect(result).To(BeIdenticalTo(container))
		})

		It("renders correctly with a custom theme", func() {
			theme := themes.NewDefaultTheme()
			output := container.WithTheme(theme).SetTitle("Themed").Render()
			Expect(output).To(ContainSubstring("Themed"))
		})

		It("falls back to default theme when nil", func() {
			output := container.WithTheme(nil).SetTitle("Nil Theme").Render()
			Expect(output).To(ContainSubstring("Nil Theme"))
		})
	})

	Describe("WithWidth", func() {
		It("returns the container for chaining", func() {
			result := container.WithWidth(80)
			Expect(result).To(BeIdenticalTo(container))
		})

		It("renders with specified width", func() {
			output := container.SetTitle("Wide").WithWidth(100).Render()
			Expect(output).NotTo(BeEmpty())
		})
	})

	Describe("WithScrollHint", func() {
		It("returns the container for chaining", func() {
			result := container.WithScrollHint(true)
			Expect(result).To(BeIdenticalTo(container))
		})

		It("includes scroll hint when enabled", func() {
			output := container.WithScrollHint(true).Render()
			Expect(output).To(ContainSubstring("Scroll"))
		})

		It("excludes scroll hint when disabled", func() {
			output := container.WithScrollHint(false).Render()
			Expect(output).NotTo(ContainSubstring("Scroll"))
		})
	})

	Describe("Render", func() {
		Context("with all components", func() {
			It("renders title, message, buttons, and instructions", func() {
				output := container.
					SetTitle("Confirm Action").
					SetMessage("This will delete the item permanently").
					SetButtons([]string{"Delete", "Cancel"}).
					SetInstructions("Press Enter to confirm").
					Render()

				Expect(output).To(ContainSubstring("Confirm Action"))
				Expect(output).To(ContainSubstring("This will delete the item permanently"))
				Expect(output).To(ContainSubstring("Delete"))
				Expect(output).To(ContainSubstring("Cancel"))
				Expect(output).To(ContainSubstring("Press Enter to confirm"))
			})
		})

		Context("with destructive style and all components", func() {
			It("renders destructive modal with full content", func() {
				output := container.
					SetTitle("Remove Entry").
					SetMessage("Are you sure?").
					SetButtons([]string{"Yes", "No"}).
					SetInstructions("This action is permanent").
					WithDestructiveStyle().
					Render()

				Expect(output).To(ContainSubstring("Remove Entry"))
				Expect(output).To(ContainSubstring("Are you sure?"))
				Expect(output).To(ContainSubstring("Yes"))
			})
		})

		Context("with scroll hint and width", func() {
			It("renders with scroll hint and specified width", func() {
				output := container.
					SetTitle("Scrollable Content").
					SetMessage("Long content here").
					WithWidth(60).
					WithScrollHint(true).
					Render()

				Expect(output).To(ContainSubstring("Scrollable Content"))
				Expect(output).To(ContainSubstring("Scroll"))
			})
		})
	})
})

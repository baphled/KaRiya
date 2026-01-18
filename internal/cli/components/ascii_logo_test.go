package components_test

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/components"
)

var _ = Describe("ASCIILogo", func() {
	var logo *components.ASCIILogo

	Describe("NewASCIILogo", func() {
		It("should create a new logo with animation enabled", func() {
			logo = components.NewASCIILogo(true, 120)
			Expect(logo).ToNot(BeNil())
		})

		It("should create a new logo without animation", func() {
			logo = components.NewASCIILogo(false, 120)
			Expect(logo).ToNot(BeNil())
		})

		It("should create logo with default tagline and version", func() {
			logo = components.NewASCIILogo(false, 120)
			view := logo.ViewStatic()
			Expect(view).To(ContainSubstring("Career Event Management System"))
			Expect(view).To(ContainSubstring("v1.0.0"))
		})
	})

	Describe("Setters", func() {
		BeforeEach(func() {
			logo = components.NewASCIILogo(false, 120)
		})

		It("should set width", func() {
			logo.SetWidth(80)
			// Width should be set (no direct getter for private field, but method shouldn't panic)
			Expect(func() { logo.SetWidth(80) }).NotTo(Panic())
		})

		It("should set tagline", func() {
			logo.SetTagline("Custom Tagline")
			view := logo.ViewStatic()
			Expect(view).To(ContainSubstring("Custom Tagline"))
			Expect(view).NotTo(ContainSubstring("Career Event Management System"))
		})

		It("should set version", func() {
			logo.SetVersion("v2.0.0")
			view := logo.ViewStatic()
			Expect(view).To(ContainSubstring("v2.0.0"))
			Expect(view).NotTo(ContainSubstring("v1.0.0"))
		})

		It("should show/hide tagline", func() {
			logo.ShowTagline(false)
			view := logo.ViewStatic()
			Expect(view).NotTo(ContainSubstring("Career Event Management System"))

			logo.ShowTagline(true)
			view = logo.ViewStatic()
			Expect(view).To(ContainSubstring("Career Event Management System"))
		})

		It("should show/hide version", func() {
			logo.ShowVersion(false)
			view := logo.ViewStatic()
			Expect(view).NotTo(ContainSubstring("v1.0.0"))

			logo.ShowVersion(true)
			view = logo.ViewStatic()
			Expect(view).To(ContainSubstring("v1.0.0"))
		})
	})

	Describe("Init", func() {
		Context("with animation disabled", func() {
			It("should initialize without animation command", func() {
				logo = components.NewASCIILogo(false, 120)
				cmd := logo.Init()
				Expect(cmd).To(BeNil())
			})

			It("should render at full opacity immediately", func() {
				logo = components.NewASCIILogo(false, 120)
				logo.Init()
				view := logo.View()
				// ASCII art contains box-drawing characters, not plain text
				Expect(view).To(ContainSubstring("██"))
			})
		})

		Context("with animation enabled", func() {
			It("should initialize with tick command", func() {
				logo = components.NewASCIILogo(true, 120)
				cmd := logo.Init()
				Expect(cmd).ToNot(BeNil())
			})
		})
	})

	Describe("Update", func() {
		Context("with animation", func() {
			BeforeEach(func() {
				logo = components.NewASCIILogo(true, 120)
				logo.Init()
			})

			It("should update fade progress on TickMsg", func() {
				// Send multiple tick messages to progress animation
				// Animation increments by 0.1 each tick: 0.0, 0.1, 0.2... 0.9, 1.0
				// So it needs 10 ticks to reach 1.0, and on the 11th update it returns nil
				for i := 0; i < 11; i++ {
					model, cmd := logo.Update(components.TickMsg{})
					logo = model.(*components.ASCIILogo)
					if i < 10 {
						Expect(cmd).ToNot(BeNil()) // Should continue ticking
					} else {
						Expect(cmd).To(BeNil()) // Should stop at full opacity
					}
				}

				// Final view should be fully rendered with actual ASCII art logo
				finalView := logo.View()

				// Check for the actual ASCII art characters (verbatim from logoArt constant)
				Expect(finalView).To(ContainSubstring("██╗  ██╗ █████╗ ██████╗ ██╗██╗   ██╗ █████╗"))
				Expect(finalView).To(ContainSubstring("██║ ██╔╝██╔══██╗██╔══██╗██║╚██╗ ██╔╝██╔══██╗"))
				Expect(finalView).To(ContainSubstring("█████╔╝ ███████║██████╔╝██║ ╚████╔╝ ███████║"))
				Expect(finalView).To(ContainSubstring("██╔═██╗ ██╔══██║██╔══██╗██║  ╚██╔╝  ██╔══██║"))
				Expect(finalView).To(ContainSubstring("██║  ██╗██║  ██║██║  ██║██║   ██║   ██║  ██║"))
				Expect(finalView).To(ContainSubstring("╚═╝  ╚═╝╚═╝  ╚═╝╚═╝  ╚═╝╚═╝   ╚═╝   ╚═╝  ╚═╝"))

				// Check for tagline and version
				Expect(finalView).To(ContainSubstring("Career Event Management System"))
				Expect(finalView).To(ContainSubstring("v1.0.0"))
			})

			It("should stop animation when fade progress reaches 1.0", func() {
				// Progress through animation
				for i := 0; i < 15; i++ { // More than needed
					model, cmd := logo.Update(components.TickMsg{})
					logo = model.(*components.ASCIILogo)
					if cmd == nil {
						break
					}
				}

				// Further updates should not produce commands
				_, cmd := logo.Update(components.TickMsg{})
				Expect(cmd).To(BeNil())
			})

			It("should handle non-TickMsg gracefully", func() {
				model, cmd := logo.Update(tea.KeyMsg{Type: tea.KeyEnter})
				Expect(model).ToNot(BeNil())
				Expect(cmd).To(BeNil())
			})
		})

		Context("without animation", func() {
			BeforeEach(func() {
				logo = components.NewASCIILogo(false, 120)
				logo.Init()
			})

			It("should not respond to TickMsg", func() {
				_, cmd := logo.Update(components.TickMsg{})
				Expect(cmd).To(BeNil())
			})
		})
	})

	Describe("View", func() {
		BeforeEach(func() {
			logo = components.NewASCIILogo(false, 120)
			logo.Init()
		})

		It("should render the logo", func() {
			view := logo.View()
			Expect(view).ToNot(BeEmpty())
			Expect(view).To(ContainSubstring("██"))
		})

		It("should include tagline when enabled", func() {
			logo.ShowTagline(true)
			view := logo.View()
			Expect(view).To(ContainSubstring("Career Event Management System"))
		})

		It("should exclude tagline when disabled", func() {
			logo.ShowTagline(false)
			view := logo.View()
			Expect(view).NotTo(ContainSubstring("Career Event Management System"))
		})

		It("should include version when enabled", func() {
			logo.ShowVersion(true)
			view := logo.View()
			Expect(view).To(ContainSubstring("v1.0.0"))
		})

		It("should exclude version when disabled", func() {
			logo.ShowVersion(false)
			view := logo.View()
			Expect(view).NotTo(ContainSubstring("v1.0.0"))
		})

		It("should render logo with custom tagline", func() {
			logo.SetTagline("Test Tagline")
			view := logo.View()
			Expect(view).To(ContainSubstring("Test Tagline"))
		})

		It("should render logo with custom version", func() {
			logo.SetVersion("v99.99.99")
			view := logo.View()
			Expect(view).To(ContainSubstring("v99.99.99"))
		})
	})

	Describe("ViewStatic", func() {
		It("should render at full opacity regardless of animation state", func() {
			logo = components.NewASCIILogo(true, 120)
			// Don't init - animation not started
			staticView := logo.ViewStatic()
			Expect(staticView).To(ContainSubstring("██"))
			Expect(staticView).To(ContainSubstring("Career Event Management System"))
		})

		It("should not affect current fade progress", func() {
			logo = components.NewASCIILogo(true, 120)
			logo.Init()

			// Render static view
			staticView := logo.ViewStatic()

			// Regular view should still show animation state
			animatedView := logo.View()

			// Both should contain logo
			Expect(staticView).To(ContainSubstring("██"))
			Expect(animatedView).To(ContainSubstring("██"))
		})
	})

	Describe("GetHeight", func() {
		It("should return correct height with tagline and version", func() {
			logo = components.NewASCIILogo(false, 120)
			logo.ShowTagline(true)
			logo.ShowVersion(true)
			height := logo.GetHeight()
			Expect(height).To(Equal(9)) // 6 (logo) + 2 (tagline) + 1 (version)
		})

		It("should return correct height without tagline", func() {
			logo = components.NewASCIILogo(false, 120)
			logo.ShowTagline(false)
			logo.ShowVersion(true)
			height := logo.GetHeight()
			Expect(height).To(Equal(7)) // 6 (logo) + 1 (version)
		})

		It("should return correct height without version", func() {
			logo = components.NewASCIILogo(false, 120)
			logo.ShowTagline(true)
			logo.ShowVersion(false)
			height := logo.GetHeight()
			Expect(height).To(Equal(8)) // 6 (logo) + 2 (tagline)
		})

		It("should return logo-only height when both disabled", func() {
			logo = components.NewASCIILogo(false, 120)
			logo.ShowTagline(false)
			logo.ShowVersion(false)
			height := logo.GetHeight()
			Expect(height).To(Equal(6)) // 6 (logo only)
		})
	})

	Describe("GetWidth", func() {
		It("should return the logo width", func() {
			logo = components.NewASCIILogo(false, 120)
			width := logo.GetWidth()
			Expect(width).To(Equal(51)) // Logo art is 51 characters wide
		})
	})

	Describe("Edge Cases", func() {
		It("should handle empty tagline", func() {
			logo = components.NewASCIILogo(false, 120)
			logo.SetTagline("")
			logo.ShowTagline(true)
			view := logo.ViewStatic()
			Expect(view).NotTo(ContainSubstring("Career Event Management System"))
		})

		It("should handle empty version", func() {
			logo = components.NewASCIILogo(false, 120)
			logo.SetVersion("")
			logo.ShowVersion(true)
			view := logo.ViewStatic()
			Expect(view).NotTo(ContainSubstring("v1.0.0"))
		})

		It("should handle very small width", func() {
			logo = components.NewASCIILogo(false, 1)
			view := logo.ViewStatic()
			Expect(view).ToNot(BeEmpty())
		})

		It("should handle very large width", func() {
			logo = components.NewASCIILogo(false, 10000)
			view := logo.ViewStatic()
			Expect(view).ToNot(BeEmpty())
		})

		It("should render consistently across multiple calls", func() {
			logo = components.NewASCIILogo(false, 120)
			logo.Init()

			view1 := logo.View()
			view2 := logo.View()
			view3 := logo.View()

			Expect(view1).To(Equal(view2))
			Expect(view2).To(Equal(view3))
		})
	})

	Describe("Animation Lifecycle", func() {
		It("should complete full animation cycle", func() {
			logo = components.NewASCIILogo(true, 120)
			cmd := logo.Init()
			Expect(cmd).ToNot(BeNil())

			tickCount := 0
			for cmd != nil && tickCount < 20 {
				model, nextCmd := logo.Update(components.TickMsg{})
				logo = model.(*components.ASCIILogo)
				cmd = nextCmd
				tickCount++
			}

			Expect(tickCount).To(BeNumerically("<=", 11)) // Should complete in ~10 ticks (allow 11 for edge)
			finalView := logo.View()
			Expect(finalView).To(ContainSubstring("██"))
		})
	})

	Describe("Logo Content", func() {
		It("should contain ASCII art with box-drawing characters", func() {
			logo = components.NewASCIILogo(false, 120)
			view := logo.ViewStatic()

			// Logo should contain box-drawing characters
			Expect(view).To(ContainSubstring("██"))
			Expect(view).To(ContainSubstring("╗"))
			Expect(view).To(ContainSubstring("╚"))
		})

		It("should be multi-line", func() {
			logo = components.NewASCIILogo(false, 120)
			view := logo.ViewStatic()
			lines := strings.Split(view, "\n")
			Expect(len(lines)).To(BeNumerically(">=", 6)) // At least 6 lines for logo art
		})
	})
})

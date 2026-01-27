package intents_test

import (
	"errors"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/terminal"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/display"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("View Helpers", func() {
	var base *intents.BaseIntent
	var info *terminal.Info

	BeforeEach(func() {
		base = intents.NewBaseIntent()
		info = terminal.NewInfo()
		info.Width = 100
		info.Height = 30
		info.IsValid = true
		base.UpdateTerminalInfo(info)
	})

	Describe("CreateStandardView", func() {
		It("should return non-nil view", func() {
			view := intents.CreateStandardView(base)
			Expect(view).NotTo(BeNil())
		})

		It("should have correct terminal info", func() {
			view := intents.CreateStandardView(base)
			Expect(view.TerminalInfo).To(Equal(info))
		})

		It("should use full width", func() {
			view := intents.CreateStandardView(base)
			Expect(view.UseFullWidth).To(BeTrue())
		})

		Context("with logo", func() {
			It("should include logo from BaseIntent", func() {
				logo := display.NewLogo(false, 100)
				base.SetLogo(logo)
				base.SetLogoSpacing(3)

				view := intents.CreateStandardView(base)
				Expect(view.Logo).To(Equal(logo))
				Expect(view.LogoSpacing).To(Equal(3))
			})
		})
	})

	Describe("CreateStandardViewWithBreadcrumbs", func() {
		It("should return non-nil view", func() {
			view := intents.CreateStandardViewWithBreadcrumbs(base, "Home", "Settings", "Display")
			Expect(view).NotTo(BeNil())
		})

		It("should include all breadcrumbs", func() {
			view := intents.CreateStandardViewWithBreadcrumbs(base, "Home", "Settings", "Display")
			Expect(view.Breadcrumbs).To(HaveLen(3))
			Expect(view.Breadcrumbs).To(Equal([]string{"Home", "Settings", "Display"}))
		})
	})

	Describe("Modal Priority", func() {
		Context("with error state", func() {
			It("should show error modal", func() {
				base.SetError(errors.New("test error"))
				view := intents.CreateStandardView(base)

				Expect(view.ShowModal).To(BeTrue())
				Expect(view.Modal).NotTo(BeNil())
				modal, ok := view.Modal.(*feedback.Modal)
				Expect(ok).To(BeTrue(), "Modal should be *feedback.Modal")
				Expect(modal.Type).To(Equal(feedback.ModalError))
			})
		})

		Context("with loading state", func() {
			It("should show loading modal", func() {
				base.SetLoading("Loading...")
				view := intents.CreateStandardView(base)

				Expect(view.ShowModal).To(BeTrue())
				Expect(view.Modal).NotTo(BeNil())
				modal, ok := view.Modal.(*feedback.Modal)
				Expect(ok).To(BeTrue(), "Modal should be *feedback.Modal")
				Expect(modal.Type).To(Equal(feedback.ModalLoading))
			})
		})

		Context("with progress state", func() {
			It("should show progress modal", func() {
				base.SetProgress("Processing", "Step 1", 0.5)
				view := intents.CreateStandardView(base)

				Expect(view.ShowModal).To(BeTrue())
				Expect(view.Modal).NotTo(BeNil())
				modal, ok := view.Modal.(*feedback.Modal)
				Expect(ok).To(BeTrue(), "Modal should be *feedback.Modal")
				Expect(modal.Type).To(Equal(feedback.ModalProgress))
			})
		})

		Context("with success state", func() {
			It("should show success modal", func() {
				base.SetSuccess("Success!")
				view := intents.CreateStandardView(base)

				Expect(view.ShowModal).To(BeTrue())
				Expect(view.Modal).NotTo(BeNil())
				modal, ok := view.Modal.(*feedback.Modal)
				Expect(ok).To(BeTrue(), "Modal should be *feedback.Modal")
				Expect(modal.Type).To(Equal(feedback.ModalSuccess))
			})
		})

		Context("error over loading priority", func() {
			It("should show error modal when both set", func() {
				base.SetLoading("Loading...")
				base.SetError(errors.New("error occurred"))
				view := intents.CreateStandardView(base)

				modal, ok := view.Modal.(*feedback.Modal)
				Expect(ok).To(BeTrue(), "Modal should be *feedback.Modal")
				Expect(modal.Type).To(Equal(feedback.ModalError))
			})
		})

		Context("loading over success priority", func() {
			It("should show loading modal when both set", func() {
				base.SetLoading("Loading...")
				base.SetSuccess("Success!")
				view := intents.CreateStandardView(base)

				modal, ok := view.Modal.(*feedback.Modal)
				Expect(ok).To(BeTrue(), "Modal should be *feedback.Modal")
				Expect(modal.Type).To(Equal(feedback.ModalLoading))
			})
		})
	})

	Describe("Error Title Extraction", func() {
		It("should extract validation error title", func() {
			base.SetError(errors.New("validation failed: field is required"))
			view := intents.CreateStandardView(base)
			modal, ok := view.Modal.(*feedback.Modal)
			Expect(ok).To(BeTrue(), "Modal should be *feedback.Modal")
			Expect(modal.Title).To(ContainSubstring("Validation"))
		})

		It("should extract database error title", func() {
			base.SetError(errors.New("database connection failed"))
			view := intents.CreateStandardView(base)
			modal, ok := view.Modal.(*feedback.Modal)
			Expect(ok).To(BeTrue(), "Modal should be *feedback.Modal")
			Expect(modal.Title).To(ContainSubstring("Database"))
		})

		It("should use default error title", func() {
			base.SetError(errors.New("something went wrong"))
			view := intents.CreateStandardView(base)
			modal, ok := view.Modal.(*feedback.Modal)
			Expect(ok).To(BeTrue(), "Modal should be *feedback.Modal")
			Expect(modal.Title).To(Equal("Error"))
		})
	})

	Describe("Themed Footer Helpers", func() {
		var theme themes.Theme

		BeforeEach(func() {
			theme = themes.NewDefaultTheme()
		})

		Describe("ThemedNavigationFooter", func() {
			It("should contain navigation elements", func() {
				footer := intents.ThemedNavigationFooter(theme)

				Expect(footer).NotTo(BeEmpty())
				expectedParts := []string{"Navigate", "Select", "Back"}
				for _, part := range expectedParts {
					Expect(footer).To(ContainSubstring(part))
				}
			})

			It("should work with nil theme", func() {
				footer := intents.ThemedNavigationFooter(nil)
				Expect(footer).NotTo(BeEmpty())
				Expect(footer).To(ContainSubstring("Navigate"))
			})
		})

		Describe("ThemedFormFooter", func() {
			It("should contain form elements", func() {
				footer := intents.ThemedFormFooter(theme)

				Expect(footer).NotTo(BeEmpty())
				expectedParts := []string{"Next", "Previous", "Submit", "Cancel"}
				for _, part := range expectedParts {
					Expect(footer).To(ContainSubstring(part))
				}
			})
		})

		Describe("ThemedListFooter", func() {
			It("should contain list elements including search", func() {
				footer := intents.ThemedListFooter(theme)

				Expect(footer).NotTo(BeEmpty())
				expectedParts := []string{"Navigate", "Select", "Search", "Back"}
				for _, part := range expectedParts {
					Expect(footer).To(ContainSubstring(part))
				}
			})
		})

		Describe("ThemedDetailViewFooter", func() {
			It("should contain scroll and back elements", func() {
				footer := intents.ThemedDetailViewFooter(theme)

				Expect(footer).NotTo(BeEmpty())
				expectedParts := []string{"Scroll", "Back"}
				for _, part := range expectedParts {
					Expect(footer).To(ContainSubstring(part))
				}
			})
		})

		Describe("ThemedCustomFooter", func() {
			It("should contain custom badges", func() {
				footer := intents.ThemedCustomFooter(theme,
					primitives.HelpKeyBadge("x", "Custom1", theme),
					primitives.HelpKeyBadge("y", "Custom2", theme),
				)

				Expect(footer).NotTo(BeEmpty())
				expectedParts := []string{"Custom1", "Custom2", "x", "y"}
				for _, part := range expectedParts {
					Expect(footer).To(ContainSubstring(part))
				}
			})

			It("should return empty for no badges", func() {
				footer := intents.ThemedCustomFooter(theme)
				Expect(footer).To(BeEmpty())
			})
		})

		Describe("ThemedGlobalBadges", func() {
			It("should contain quit badge", func() {
				footer := intents.ThemedGlobalBadges(theme)

				Expect(footer).NotTo(BeEmpty())
				Expect(footer).To(ContainSubstring("Quit"))
			})
		})

		Describe("CombineThemedFooters", func() {
			It("should combine themed footers", func() {
				footer1 := intents.ThemedNavigationFooter(theme)
				footer2 := intents.ThemedGlobalBadges(theme)

				combined := intents.CombineThemedFooters(footer1, footer2)

				Expect(combined).NotTo(BeEmpty())
				Expect(combined).To(ContainSubstring("Navigate"))
				Expect(combined).To(ContainSubstring("Quit"))
			})

			It("should return empty for no footers", func() {
				combined := intents.CombineThemedFooters()
				Expect(combined).To(BeEmpty())
			})

			It("should skip empty strings", func() {
				footer1 := intents.ThemedNavigationFooter(theme)
				combined := intents.CombineThemedFooters(footer1, "", "   ")

				Expect(combined).NotTo(BeEmpty())
				Expect(combined).To(ContainSubstring("Navigate"))
			})
		})
	})
})

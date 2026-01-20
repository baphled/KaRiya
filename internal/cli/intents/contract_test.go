package intents_test

import (
	"errors"
	"time"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/terminal"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/display"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("BaseIntent", func() {
	var base *intents.BaseIntent

	BeforeEach(func() {
		base = intents.NewBaseIntent()
	})

	Describe("NewBaseIntent", func() {
		It("should return non-nil", func() {
			Expect(base).NotTo(BeNil())
		})

		It("should initialize terminal info", func() {
			Expect(base.GetTerminalInfo()).NotTo(BeNil())
		})

		It("should set default logo spacing to 2", func() {
			Expect(base.GetLogoSpacing()).To(Equal(2))
		})
	})

	Describe("TerminalInfo", func() {
		It("should update and retrieve terminal info", func() {
			info := terminal.NewInfo()
			info.Width = 120
			info.Height = 40
			info.IsValid = true

			base.UpdateTerminalInfo(info)

			retrieved := base.GetTerminalInfo()
			Expect(retrieved.Width).To(Equal(120))
			Expect(retrieved.Height).To(Equal(40))
			Expect(retrieved.IsValid).To(BeTrue())
		})
	})

	Describe("GetMinimumSize", func() {
		It("should return positive minimum size", func() {
			width, height := base.GetMinimumSize()
			Expect(width).To(BeNumerically(">", 0))
			Expect(height).To(BeNumerically(">", 0))
		})
	})

	Describe("GetModalDimensions", func() {
		It("should return defaults when terminal info has zero dimensions", func() {
			width, height := base.GetModalDimensions()
			Expect(width).To(Equal(120))
			Expect(height).To(Equal(40))
		})

		It("should return terminal dimensions when available", func() {
			info := terminal.NewInfo()
			info.Width = 200
			info.Height = 60
			info.IsValid = true
			base.UpdateTerminalInfo(info)

			width, height := base.GetModalDimensions()
			Expect(width).To(Equal(200))
			Expect(height).To(Equal(60))
		})

		It("should return defaults when terminal info is nil", func() {
			base.UpdateTerminalInfo(nil)
			width, height := base.GetModalDimensions()
			Expect(width).To(Equal(120))
			Expect(height).To(Equal(40))
		})

		It("should return defaults when only width is zero", func() {
			info := terminal.NewInfo()
			info.Width = 0
			info.Height = 60
			base.UpdateTerminalInfo(info)

			width, height := base.GetModalDimensions()
			Expect(width).To(Equal(120))
			Expect(height).To(Equal(40))
		})

		It("should return defaults when only height is zero", func() {
			info := terminal.NewInfo()
			info.Width = 200
			info.Height = 0
			base.UpdateTerminalInfo(info)

			width, height := base.GetModalDimensions()
			Expect(width).To(Equal(120))
			Expect(height).To(Equal(40))
		})
	})

	Describe("Logo Management", func() {
		It("should return nil logo initially", func() {
			Expect(base.GetLogo()).To(BeNil())
		})

		It("should set and get logo", func() {
			logo := display.NewLogo(false, 80)
			base.SetLogo(logo)
			Expect(base.GetLogo()).To(Equal(logo))
		})
	})

	Describe("Logo Spacing", func() {
		It("should return default spacing of 2", func() {
			Expect(base.GetLogoSpacing()).To(Equal(2))
		})

		It("should set and get custom spacing", func() {
			base.SetLogoSpacing(5)
			Expect(base.GetLogoSpacing()).To(Equal(5))
		})
	})

	Describe("Loading State", func() {
		It("should not be loading initially", func() {
			Expect(base.IsLoading()).To(BeFalse())
		})

		It("should set loading state and message", func() {
			base.SetLoading("Processing...")
			Expect(base.IsLoading()).To(BeTrue())
			Expect(base.GetLoadingMessage()).To(Equal("Processing..."))
		})

		It("should clear loading state", func() {
			base.SetLoading("Processing...")
			base.ClearLoading()
			Expect(base.IsLoading()).To(BeFalse())
			Expect(base.GetLoadingMessage()).To(BeEmpty())
		})
	})

	Describe("Error State", func() {
		It("should not have error initially", func() {
			Expect(base.HasError()).To(BeFalse())
			Expect(base.GetError()).To(BeNil())
		})

		It("should set and get error", func() {
			testErr := errors.New("test error")
			base.SetError(testErr)
			Expect(base.HasError()).To(BeTrue())
			Expect(base.GetError()).To(Equal(testErr))
		})

		It("should clear error", func() {
			base.SetError(errors.New("test error"))
			base.ClearError()
			Expect(base.HasError()).To(BeFalse())
			Expect(base.GetError()).To(BeNil())
		})
	})

	Describe("Success State", func() {
		It("should not show success initially", func() {
			Expect(base.ShouldShowSuccess()).To(BeFalse())
		})

		It("should set and get success message", func() {
			base.SetSuccess("Operation completed!")
			Expect(base.ShouldShowSuccess()).To(BeTrue())
			Expect(base.GetSuccessMessage()).To(Equal("Operation completed!"))
		})

		It("should clear success state", func() {
			base.SetSuccess("Success!")
			base.ClearSuccess()
			Expect(base.ShouldShowSuccess()).To(BeFalse())
			Expect(base.GetSuccessMessage()).To(BeEmpty())
		})
	})

	Describe("Success Expiry", func() {
		It("should expire success after 3 seconds", func() {
			base.SetSuccess("Test message")
			Expect(base.ShouldShowSuccess()).To(BeTrue())

			time.Sleep(3500 * time.Millisecond)

			Expect(base.ShouldShowSuccess()).To(BeFalse())
		})
	})

	Describe("Progress State", func() {
		It("should not be enabled initially", func() {
			Expect(base.IsProgressEnabled()).To(BeFalse())
		})

		It("should set and get progress", func() {
			base.SetProgress("Processing", "Step 1 of 3", 0.33)
			Expect(base.IsProgressEnabled()).To(BeTrue())

			title, message, value := base.GetProgress()
			Expect(title).To(Equal("Processing"))
			Expect(message).To(Equal("Step 1 of 3"))
			Expect(value).To(Equal(0.33))
		})

		It("should clear progress state", func() {
			base.SetProgress("Processing", "Step 1", 0.5)
			base.ClearProgress()
			Expect(base.IsProgressEnabled()).To(BeFalse())

			title, message, value := base.GetProgress()
			Expect(title).To(BeEmpty())
			Expect(message).To(BeEmpty())
			Expect(value).To(Equal(0.0))
		})
	})

	Describe("CreateView", func() {
		BeforeEach(func() {
			info := terminal.NewInfo()
			info.Width = 100
			info.Height = 30
			info.IsValid = true
			base.UpdateTerminalInfo(info)
		})

		It("should return non-nil view", func() {
			view := base.CreateView()
			Expect(view).NotTo(BeNil())
		})

		It("should include terminal info", func() {
			view := base.CreateView()
			Expect(view.TerminalInfo).To(Equal(base.GetTerminalInfo()))
		})

		Context("with logo", func() {
			It("should include logo in view", func() {
				logo := display.NewLogo(false, 100)
				base.SetLogo(logo)

				view := base.CreateView()
				Expect(view.Logo).To(Equal(logo))
				Expect(view.LogoSpacing).To(Equal(2))
			})
		})

		Context("with error state", func() {
			It("should show modal when error is set", func() {
				base.SetError(errors.New("test error"))
				view := base.CreateView()
				Expect(view.ShowModal).To(BeTrue())
				Expect(view.Modal).NotTo(BeNil())
			})
		})
	})

	Describe("CreateViewWithBreadcrumbs", func() {
		BeforeEach(func() {
			info := terminal.NewInfo()
			info.Width = 100
			info.Height = 30
			info.IsValid = true
			base.UpdateTerminalInfo(info)
		})

		It("should include breadcrumbs", func() {
			view := base.CreateViewWithBreadcrumbs("Main", "Settings", "Display")
			Expect(view).NotTo(BeNil())
			Expect(view.Breadcrumbs).To(HaveLen(3))
			Expect(view.Breadcrumbs).To(Equal([]string{"Main", "Settings", "Display"}))
		})
	})

	Describe("State Independence", func() {
		It("should maintain independent states", func() {
			base.SetLoading("Loading...")
			base.SetError(errors.New("error occurred"))
			base.SetSuccess("Success!")
			base.SetProgress("Processing", "Step 1", 0.5)

			Expect(base.IsLoading()).To(BeTrue())
			Expect(base.HasError()).To(BeTrue())
			Expect(base.ShouldShowSuccess()).To(BeTrue())
			Expect(base.IsProgressEnabled()).To(BeTrue())
		})

		It("should not affect other states when clearing one", func() {
			base.SetLoading("Loading...")
			base.SetError(errors.New("error"))
			base.SetSuccess("Success!")
			base.SetProgress("Processing", "Step 1", 0.5)

			base.ClearLoading()

			Expect(base.HasError()).To(BeTrue())
			Expect(base.ShouldShowSuccess()).To(BeTrue())
			Expect(base.IsProgressEnabled()).To(BeTrue())
		})
	})

	Describe("Theme Management", func() {
		It("should return nil theme manager initially", func() {
			Expect(base.GetThemeManager()).To(BeNil())
		})

		It("should return nil theme when no manager is set", func() {
			Expect(base.Theme()).To(BeNil())
		})

		It("should set and get theme manager", func() {
			tm := themes.NewThemeManager()
			base.SetThemeManager(tm)
			Expect(base.GetThemeManager()).To(Equal(tm))
		})

		It("should return active theme when manager is set", func() {
			tm := themes.NewThemeManager()
			base.SetThemeManager(tm)

			theme := base.Theme()
			Expect(theme).NotTo(BeNil())
			Expect(theme.Name()).To(Equal("default"))
		})
	})

	Describe("Theme Styles", func() {
		BeforeEach(func() {
			tm := themes.NewThemeManager()
			base.SetThemeManager(tm)
		})

		It("should access palette", func() {
			theme := base.Theme()
			Expect(theme).NotTo(BeNil())
			Expect(theme.Palette()).NotTo(BeNil())
		})

		It("should access styles", func() {
			theme := base.Theme()
			Expect(theme).NotTo(BeNil())
			Expect(theme.Styles()).NotTo(BeNil())
		})
	})

	Describe("Help Modal", func() {
		It("should be hidden initially", func() {
			Expect(base.IsHelpVisible()).To(BeFalse())
		})

		It("should toggle visibility", func() {
			base.ToggleHelp()
			Expect(base.IsHelpVisible()).To(BeTrue())

			base.ToggleHelp()
			Expect(base.IsHelpVisible()).To(BeFalse())
		})

		It("should show and hide", func() {
			base.ShowHelp()
			Expect(base.IsHelpVisible()).To(BeTrue())

			base.HideHelp()
			Expect(base.IsHelpVisible()).To(BeFalse())
		})

		It("should work with terminal dimensions", func() {
			info := terminal.NewInfo()
			info.Width = 120
			info.Height = 40
			info.IsValid = true
			base.UpdateTerminalInfo(info)

			base.ShowHelp()
			Expect(base.IsHelpVisible()).To(BeTrue())
		})
	})
})

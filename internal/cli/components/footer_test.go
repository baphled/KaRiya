package components

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Footer", func() {
	var footer FooterModel

	BeforeEach(func() {
		footer = NewFooter(80)
	})

	Describe("Creation", func() {
		It("creates a footer with width", func() {
			Expect(footer.width).To(Equal(80))
		})

		It("initializes with empty status message", func() {
			Expect(footer.GetStatusMessage()).To(Equal(""))
		})

		It("initializes with empty mode context", func() {
			Expect(footer.GetModeContext()).To(Equal(""))
		})

		It("initializes with status and mode shown by default", func() {
			Expect(footer.showStatus).To(BeTrue())
			Expect(footer.showMode).To(BeTrue())
		})
	})

	Describe("Status Message", func() {
		It("sets status message", func() {
			footer.SetStatusMessage("5/10 events")
			Expect(footer.GetStatusMessage()).To(Equal("5/10 events"))
		})

		It("renders status message", func() {
			footer.SetStatusMessage("Processing...")
			view := footer.View()
			Expect(view).To(ContainSubstring("Processing..."))
		})

		It("hides status when showStatus is false", func() {
			footer.SetStatusMessage("5/10 events")
			footer.SetShowStatus(false)
			view := footer.View()
			Expect(view).NotTo(ContainSubstring("5/10 events"))
		})

		It("truncates long status message", func() {
			longStatus := "This is a very long status message that exceeds the width limit"
			footer.SetStatusMessage(longStatus)
			footer.SetWidth(40)
			view := footer.View()
			Expect(view).To(ContainSubstring("..."))
		})
	})

	Describe("Mode Context", func() {
		It("sets mode context", func() {
			footer.SetModeContext("Capture Mode: Timeline")
			Expect(footer.GetModeContext()).To(Equal("Capture Mode: Timeline"))
		})

		It("renders mode context", func() {
			footer.SetModeContext("List View")
			view := footer.View()
			Expect(view).To(ContainSubstring("List View"))
		})

		It("hides mode when showMode is false", func() {
			footer.SetModeContext("Capture Mode: Timeline")
			footer.SetShowMode(false)
			view := footer.View()
			Expect(view).NotTo(ContainSubstring("Capture Mode: Timeline"))
		})

		It("truncates long mode context", func() {
			longMode := "This is a very long mode context that exceeds the width limit"
			footer.SetModeContext(longMode)
			footer.SetWidth(40)
			view := footer.View()
			Expect(view).To(ContainSubstring("..."))
		})
	})

	Describe("Status and Mode Together", func() {
		It("renders both status and mode", func() {
			footer.SetStatusMessage("3/10 events")
			footer.SetModeContext("Capture Mode: Timeline")
			view := footer.View()
			Expect(view).To(ContainSubstring("3/10 events"))
			Expect(view).To(ContainSubstring("Capture Mode: Timeline"))
		})

		It("includes separator between status and mode", func() {
			footer.SetStatusMessage("3/10 events")
			footer.SetModeContext("Capture Mode: Timeline")
			view := footer.View()
			Expect(view).To(ContainSubstring("|"))
		})

		It("shows only status when mode is empty", func() {
			footer.SetStatusMessage("3/10 events")
			footer.SetModeContext("")
			view := footer.View()
			Expect(view).To(ContainSubstring("3/10 events"))
		})

		It("shows only mode when status is empty", func() {
			footer.SetStatusMessage("")
			footer.SetModeContext("Capture Mode: Timeline")
			view := footer.View()
			Expect(view).To(ContainSubstring("Capture Mode: Timeline"))
		})

		It("shows nothing when both status and mode are empty", func() {
			footer.SetStatusMessage("")
			footer.SetModeContext("")
			view := footer.View()
			Expect(view).To(Equal(""))
		})
	})

	Describe("Help Footer Integration", func() {
		It("accepts help footer component", func() {
			helpFooter := NewHelpFooter("form", 80)
			footer.SetHelpFooter(&helpFooter)
			Expect(footer.helpFooter).NotTo(BeNil())
			Expect(footer.showHelp).To(BeTrue())
		})

		It("renders with help footer", func() {
			helpFooter := NewHelpFooter("form", 80)
			footer.SetStatusMessage("Status")
			footer.SetHelpFooter(&helpFooter)
			view := footer.View()
			Expect(len(view)).To(BeNumerically(">", 0))
		})

		It("toggles help footer display", func() {
			helpFooter := NewHelpFooter("form", 80)
			footer.SetHelpFooter(&helpFooter)
			footer.SetShowHelp(false)
			Expect(footer.showHelp).To(BeFalse())
		})

		It("handles nil help footer gracefully", func() {
			footer.SetShowHelp(true) // Should not crash even if helpFooter is nil
			view := footer.View()
			Expect(len(view)).To(BeNumerically(">=", 0))
		})
	})

	Describe("Configuration", func() {
		It("sets width", func() {
			footer.SetWidth(100)
			Expect(footer.width).To(Equal(100))
		})

		It("sets height", func() {
			footer.SetHeight(2)
			Expect(footer.height).To(Equal(2))
		})

		It("propagates width to help footer", func() {
			helpFooter := NewHelpFooter("form", 80)
			footer.SetHelpFooter(&helpFooter)
			footer.SetWidth(120)
			Expect(footer.helpFooter.width).To(Equal(120))
		})

		It("toggles status display", func() {
			footer.SetStatusMessage("Status")
			footer.SetShowStatus(false)
			Expect(footer.showStatus).To(BeFalse())
		})

		It("toggles mode display", func() {
			footer.SetModeContext("Mode")
			footer.SetShowMode(false)
			Expect(footer.showMode).To(BeFalse())
		})

		It("toggles help display", func() {
			footer.SetShowHelp(true)
			Expect(footer.showHelp).To(BeFalse()) // Still false because helpFooter is nil
		})
	})

	Describe("Rendering", func() {
		It("renders empty footer when width is 0", func() {
			footer.SetWidth(0)
			view := footer.View()
			Expect(view).To(Equal(""))
		})

		It("renders empty footer when width is negative", func() {
			footer.SetWidth(-1)
			view := footer.View()
			Expect(view).To(Equal(""))
		})

		It("renders status only when mode is disabled", func() {
			footer.SetStatusMessage("Status")
			footer.SetModeContext("Mode")
			footer.SetShowMode(false)
			view := footer.View()
			Expect(view).To(ContainSubstring("Status"))
			Expect(view).NotTo(ContainSubstring("Mode"))
		})

		It("renders mode only when status is disabled", func() {
			footer.SetStatusMessage("Status")
			footer.SetModeContext("Mode")
			footer.SetShowStatus(false)
			view := footer.View()
			Expect(view).NotTo(ContainSubstring("Status"))
			Expect(view).To(ContainSubstring("Mode"))
		})

		It("handles responsive sizing for narrow terminals", func() {
			narrowFooter := NewFooter(20)
			narrowFooter.SetStatusMessage("Status")
			narrowFooter.SetModeContext("Mode")
			view := narrowFooter.View()
			Expect(len(view)).To(BeNumerically(">", 0))
		})

		It("handles responsive sizing for wide terminals", func() {
			wideFooter := NewFooter(200)
			wideFooter.SetStatusMessage("Status")
			wideFooter.SetModeContext("Mode")
			view := wideFooter.View()
			Expect(view).To(ContainSubstring("Status"))
			Expect(view).To(ContainSubstring("Mode"))
		})
	})

	Describe("Edge Cases", func() {
		It("handles empty status and mode strings", func() {
			footer.SetStatusMessage("")
			footer.SetModeContext("")
			view := footer.View()
			Expect(view).To(Equal(""))
		})

		It("handles unicode in status message", func() {
			footer.SetStatusMessage("événements: 5")
			view := footer.View()
			Expect(view).To(ContainSubstring("événements"))
		})

		It("handles unicode in mode context", func() {
			footer.SetModeContext("Aperçu: Chronologie")
			view := footer.View()
			Expect(view).To(ContainSubstring("Aperçu"))
		})

		It("handles special characters in status", func() {
			footer.SetStatusMessage("Status [3/10] - Ready!")
			view := footer.View()
			Expect(view).To(ContainSubstring("Ready"))
		})

		It("handles very long status message", func() {
			longStatus := "This is an extremely long status message that definitely exceeds the width of the footer and should be truncated with an ellipsis to indicate there is more content"
			footer.SetStatusMessage(longStatus)
			footer.SetWidth(40)
			view := footer.View()
			Expect(view).To(ContainSubstring("..."))
		})

		It("handles very long mode context", func() {
			longMode := "This is an extremely long mode context message that definitely exceeds the width and should be truncated with ellipsis"
			footer.SetModeContext(longMode)
			footer.SetWidth(40)
			view := footer.View()
			Expect(view).To(ContainSubstring("..."))
		})
	})

	Describe("Multiple Footers", func() {
		It("creates independent footer instances", func() {
			footer1 := NewFooter(80)
			footer2 := NewFooter(100)

			footer1.SetStatusMessage("Status 1")
			footer2.SetStatusMessage("Status 2")

			Expect(footer1.GetStatusMessage()).To(Equal("Status 1"))
			Expect(footer2.GetStatusMessage()).To(Equal("Status 2"))
		})
	})
})

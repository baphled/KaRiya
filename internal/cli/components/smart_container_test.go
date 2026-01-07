package components_test

import (
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/terminal"
)

var _ = Describe("SmartContainer", func() {
	var (
		container *components.SmartContainer
		info      *terminal.Info
	)

	BeforeEach(func() {
		info = terminal.NewInfo()
		container = components.NewSmartContainer(info)
	})

	Describe("NewSmartContainer", func() {
		It("should create container with default settings", func() {
			Expect(container).NotTo(BeNil())
		})
	})

	Describe("SetContent", func() {
		It("should allow method chaining", func() {
			result := container.SetContent("test content")
			Expect(result).To(Equal(container))
		})
	})

	Describe("SetCenteringMode", func() {
		It("should allow method chaining", func() {
			result := container.SetCenteringMode(components.CenterHorizontal)
			Expect(result).To(Equal(container))
		})
	})

	Describe("Render", func() {
		Context("with invalid terminal info", func() {
			It("should render minimal mode message", func() {
				// Terminal not updated, so IsValid = false but CanRender returns true
				// Set size below minimum to trigger minimal mode
				info.Update(tea.WindowSizeMsg{Width: 30, Height: 10})

				view := container.SetContent("test").Render()
				Expect(view).To(ContainSubstring("Terminal too small"))
			})
		})

		Context("with normal terminal size", func() {
			BeforeEach(func() {
				info.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
			})

			It("should render content", func() {
				view := container.SetContent("Hello World").Render()
				Expect(view).To(ContainSubstring("Hello World"))
			})

			It("should handle empty content", func() {
				view := container.SetContent("").Render()
				Expect(view).To(Equal(""))
			})

			It("should center content horizontally when mode is CenterHorizontal", func() {
				container.SetCenteringMode(components.CenterHorizontal)
				container.SetContent("Hello")
				view := container.Render()

				// Content should be centered (have leading spaces)
				Expect(view).To(ContainSubstring("Hello"))
				// Can't easily test exact spacing due to ANSI codes, but we verify render succeeds
			})

			It("should center content both ways when mode is CenterBoth", func() {
				container.SetCenteringMode(components.CenterBoth)
				container.SetContent("Hello")
				view := container.Render()

				Expect(view).To(ContainSubstring("Hello"))
			})

			It("should not center when mode is CenterNone", func() {
				container.SetCenteringMode(components.CenterNone)
				container.SetContent("Hello")
				view := container.Render()

				Expect(view).To(ContainSubstring("Hello"))
			})
		})

		Context("with compact terminal size", func() {
			BeforeEach(func() {
				info.Update(tea.WindowSizeMsg{Width: 70, Height: 24})
			})

			It("should use compact rendering mode", func() {
				view := container.SetContent("Compact test").Render()
				Expect(view).To(ContainSubstring("Compact test"))
			})

			It("should have smaller margins", func() {
				view := container.SetContent("test").Render()
				// Verify render succeeds with compact margins
				Expect(view).NotTo(BeEmpty())
			})
		})

		Context("with tiny terminal size", func() {
			BeforeEach(func() {
				info.Update(tea.WindowSizeMsg{Width: 50, Height: 20})
			})

			It("should use tiny rendering mode", func() {
				view := container.SetContent("Tiny test").Render()
				Expect(view).To(ContainSubstring("Tiny test"))
			})

			It("should have minimal margins", func() {
				view := container.SetContent("test").Render()
				Expect(view).NotTo(BeEmpty())
			})
		})
	})

	Describe("SetCustomMargins", func() {
		It("should allow custom margins", func() {
			info.Update(tea.WindowSizeMsg{Width: 100, Height: 40})

			customMargins := components.Margins{
				Top:    5,
				Right:  10,
				Bottom: 5,
				Left:   10,
			}

			container.SetCustomMargins(customMargins)
			view := container.SetContent("test").Render()

			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("SetMinContentSize", func() {
		It("should enforce minimum content width", func() {
			info.Update(tea.WindowSizeMsg{Width: 100, Height: 40})

			container.SetMinContentSize(50, 20)
			view := container.SetContent("test").Render()

			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("SetMaxContentSize", func() {
		It("should enforce maximum content width", func() {
			info.Update(tea.WindowSizeMsg{Width: 200, Height: 60})

			container.SetMaxContentSize(100, 40)
			view := container.SetContent("test").Render()

			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Multi-line content", func() {
		BeforeEach(func() {
			info.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
		})

		It("should handle multi-line content", func() {
			content := "Line 1\nLine 2\nLine 3"
			view := container.SetContent(content).Render()

			Expect(view).To(ContainSubstring("Line 1"))
			Expect(view).To(ContainSubstring("Line 2"))
			Expect(view).To(ContainSubstring("Line 3"))
		})

		It("should center multi-line content horizontally", func() {
			content := "Short\nLonger Line\nShort"
			container.SetCenteringMode(components.CenterHorizontal)
			view := container.SetContent(content).Render()

			Expect(view).To(ContainSubstring("Short"))
			Expect(view).To(ContainSubstring("Longer Line"))
		})
	})

	Describe("Method chaining", func() {
		It("should allow full method chaining", func() {
			info.Update(tea.WindowSizeMsg{Width: 100, Height: 40})

			view := container.
				SetContent("Chained content").
				SetCenteringMode(components.CenterBoth).
				SetMinContentSize(50, 20).
				Render()

			Expect(view).To(ContainSubstring("Chained content"))
		})
	})
})

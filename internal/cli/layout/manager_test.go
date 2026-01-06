package layout_test

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/layout"
	"github.com/baphled/kariya/internal/cli/terminal"
)

func TestLayout(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Layout Suite")
}

var _ = Describe("Manager", func() {
	var (
		manager *layout.Manager
		info    *terminal.Info
	)

	BeforeEach(func() {
		info = terminal.NewInfo()
		manager = layout.NewManager(info)
	})

	Describe("NewManager", func() {
		It("should create manager with default config", func() {
			Expect(manager).NotTo(BeNil())
		})

		It("should use provided terminal info", func() {
			info.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
			manager = layout.NewManager(info)

			content := manager.GetContentArea()
			Expect(content.Width).To(BeNumerically(">", 0))
		})
	})

	Describe("GetContentArea", func() {
		Context("with normal terminal size", func() {
			BeforeEach(func() {
				info.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
			})

			It("should calculate content area with margins", func() {
				content := manager.GetContentArea()

				Expect(content.Width).To(BeNumerically(">", 0))
				Expect(content.Height).To(BeNumerically(">", 0))
				Expect(content.X).To(BeNumerically(">=", 0))
				Expect(content.Y).To(BeNumerically(">=", 0))
			})

			It("should subtract margins from total area", func() {
				content := manager.GetContentArea()

				// Content should be less than terminal size
				Expect(content.Width).To(BeNumerically("<", 100))
				Expect(content.Height).To(BeNumerically("<", 40))
			})
		})

		Context("with compact terminal size", func() {
			BeforeEach(func() {
				info.Update(tea.WindowSizeMsg{Width: 70, Height: 24})
			})

			It("should use smaller margins", func() {
				content := manager.GetContentArea()

				// Compact should have more usable space percentage-wise
				Expect(content.Width).To(BeNumerically(">", 60))
			})
		})

		Context("with tiny terminal size", func() {
			BeforeEach(func() {
				info.Update(tea.WindowSizeMsg{Width: 50, Height: 20})
			})

			It("should use minimal margins", func() {
				content := manager.GetContentArea()

				// Tiny should maximize usable space
				Expect(content.Width).To(BeNumerically(">", 45))
			})
		})
	})

	Describe("GetMargins", func() {
		It("should return tiny margins for tiny terminals", func() {
			info.Update(tea.WindowSizeMsg{Width: 50, Height: 20})
			margins := manager.GetMargins()

			// Tiny margins should be minimal
			Expect(margins.Left + margins.Right).To(BeNumerically("<=", 2))
		})

		It("should return compact margins for compact terminals", func() {
			info.Update(tea.WindowSizeMsg{Width: 70, Height: 24})
			margins := manager.GetMargins()

			// Compact margins should be modest
			Expect(margins.Left + margins.Right).To(BeNumerically("<=", 4))
		})

		It("should return normal margins for normal terminals", func() {
			info.Update(tea.WindowSizeMsg{Width: 90, Height: 30})
			margins := manager.GetMargins()

			// Normal margins should be comfortable
			Expect(margins.Left + margins.Right).To(Equal(8))
		})

		It("should return large margins for large terminals", func() {
			info.Update(tea.WindowSizeMsg{Width: 140, Height: 50})
			margins := manager.GetMargins()

			// Large margins should be spacious
			Expect(margins.Left + margins.Right).To(BeNumerically(">=", 10))
		})
	})

	Describe("CalculateColumns", func() {
		BeforeEach(func() {
			info.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
		})

		It("should divide width into equal columns", func() {
			columns := manager.CalculateColumns(3, 2)

			Expect(columns).To(HaveLen(3))
			// Total width should approximately equal content width
			total := columns[0] + columns[1] + columns[2] + 4 // +4 for gutters
			content := manager.GetContentArea()
			Expect(total).To(BeNumerically("~", content.Width, 2))
		})

		It("should handle 2 columns", func() {
			columns := manager.CalculateColumns(2, 2)

			Expect(columns).To(HaveLen(2))
			Expect(columns[0]).To(BeNumerically(">", 0))
			Expect(columns[1]).To(BeNumerically(">", 0))
		})

		It("should handle 4 columns", func() {
			columns := manager.CalculateColumns(4, 2)

			Expect(columns).To(HaveLen(4))
			for _, col := range columns {
				Expect(col).To(BeNumerically(">", 0))
			}
		})

		It("should distribute remainder to first columns", func() {
			columns := manager.CalculateColumns(3, 2)

			// Some columns may be 1 char wider due to remainder
			Expect(columns[0]).To(BeNumerically(">=", columns[1]))
			Expect(columns[1]).To(BeNumerically(">=", columns[2]))
		})
	})

	Describe("ShouldUseCompactLayout", func() {
		It("should return true for tiny terminals", func() {
			info.Update(tea.WindowSizeMsg{Width: 50, Height: 20})
			Expect(manager.ShouldUseCompactLayout()).To(BeTrue())
		})

		It("should return true for compact terminals", func() {
			info.Update(tea.WindowSizeMsg{Width: 70, Height: 24})
			Expect(manager.ShouldUseCompactLayout()).To(BeTrue())
		})

		It("should return false for normal terminals", func() {
			info.Update(tea.WindowSizeMsg{Width: 90, Height: 30})
			Expect(manager.ShouldUseCompactLayout()).To(BeFalse())
		})

		It("should return false for large terminals", func() {
			info.Update(tea.WindowSizeMsg{Width: 140, Height: 50})
			Expect(manager.ShouldUseCompactLayout()).To(BeFalse())
		})
	})

	Describe("ShouldUseListLayout", func() {
		It("should return true for very narrow terminals", func() {
			info.Update(tea.WindowSizeMsg{Width: 50, Height: 20})
			Expect(manager.ShouldUseListLayout()).To(BeTrue())
		})

		It("should return false for normal width terminals", func() {
			info.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
			Expect(manager.ShouldUseListLayout()).To(BeFalse())
		})
	})

	Describe("UpdateTerminalInfo", func() {
		It("should update internal terminal info", func() {
			newInfo := terminal.NewInfo()
			newInfo.Update(tea.WindowSizeMsg{Width: 150, Height: 50})

			manager.UpdateTerminalInfo(newInfo)

			content := manager.GetContentArea()
			// Should reflect new larger size
			Expect(content.Width).To(BeNumerically(">", 100))
		})
	})
})

var _ = Describe("Rectangle", func() {
	It("should define rectangle structure", func() {
		rect := layout.Rectangle{
			X:      10,
			Y:      5,
			Width:  80,
			Height: 30,
		}

		Expect(rect.X).To(Equal(10))
		Expect(rect.Y).To(Equal(5))
		Expect(rect.Width).To(Equal(80))
		Expect(rect.Height).To(Equal(30))
	})
})

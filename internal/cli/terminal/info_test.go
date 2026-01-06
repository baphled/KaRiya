package terminal_test

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/terminal"
)

func TestTerminal(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Terminal Suite")
}

var _ = Describe("Info", func() {
	var info *terminal.Info

	BeforeEach(func() {
		info = terminal.NewInfo()
	})

	Describe("Initial state", func() {
		It("should not be valid initially", func() {
			Expect(info.IsValid).To(BeFalse())
		})

		It("should have zero dimensions initially", func() {
			Expect(info.Width).To(Equal(0))
			Expect(info.Height).To(Equal(0))
		})
	})

	Describe("Update", func() {
		It("should update dimensions from WindowSizeMsg", func() {
			msg := tea.WindowSizeMsg{Width: 100, Height: 40}
			info.Update(msg)

			Expect(info.Width).To(Equal(100))
			Expect(info.Height).To(Equal(40))
			Expect(info.IsValid).To(BeTrue())
		})

		It("should update LastUpdated timestamp", func() {
			before := time.Now()
			msg := tea.WindowSizeMsg{Width: 80, Height: 24}
			info.Update(msg)

			Expect(info.LastUpdated).To(BeTemporally(">=", before))
		})

		It("should allow multiple updates", func() {
			info.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
			Expect(info.Width).To(Equal(80))

			info.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
			Expect(info.Width).To(Equal(120))
			Expect(info.Height).To(Equal(40))
		})
	})

	Describe("GetCategory", func() {
		It("should return SizeNormal for uninitialized info", func() {
			Expect(info.GetCategory()).To(Equal(terminal.SizeNormal))
		})

		It("should return SizeTiny for width < 60", func() {
			info.Update(tea.WindowSizeMsg{Width: 50, Height: 20})
			Expect(info.GetCategory()).To(Equal(terminal.SizeTiny))
		})

		It("should return SizeCompact for width 60-79", func() {
			info.Update(tea.WindowSizeMsg{Width: 70, Height: 24})
			Expect(info.GetCategory()).To(Equal(terminal.SizeCompact))
		})

		It("should return SizeNormal for width 80-119", func() {
			info.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
			Expect(info.GetCategory()).To(Equal(terminal.SizeNormal))
		})

		It("should return SizeLarge for width 120-159", func() {
			info.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
			Expect(info.GetCategory()).To(Equal(terminal.SizeLarge))
		})

		It("should return SizeXLarge for width >= 160", func() {
			info.Update(tea.WindowSizeMsg{Width: 200, Height: 60})
			Expect(info.GetCategory()).To(Equal(terminal.SizeXLarge))
		})
	})

	Describe("GetSafeDimensions", func() {
		It("should return defaults when not valid", func() {
			config := terminal.DefaultConfig
			width, height := info.GetSafeDimensions(config)

			Expect(width).To(Equal(config.DefaultWidth))
			Expect(height).To(Equal(config.DefaultHeight))
		})

		It("should return actual dimensions when valid and above minimum", func() {
			info.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
			config := terminal.DefaultConfig

			width, height := info.GetSafeDimensions(config)
			Expect(width).To(Equal(100))
			Expect(height).To(Equal(40))
		})

		It("should enforce minimum width", func() {
			info.Update(tea.WindowSizeMsg{Width: 30, Height: 40})
			config := terminal.DefaultConfig

			width, _ := info.GetSafeDimensions(config)
			Expect(width).To(Equal(config.MinWidth))
		})

		It("should enforce minimum height", func() {
			info.Update(tea.WindowSizeMsg{Width: 100, Height: 10})
			config := terminal.DefaultConfig

			_, height := info.GetSafeDimensions(config)
			Expect(height).To(Equal(config.MinHeight))
		})
	})

	Describe("CanRender", func() {
		It("should return true when not valid (assume defaults)", func() {
			Expect(info.CanRender(terminal.DefaultConfig)).To(BeTrue())
		})

		It("should return true when dimensions meet minimum", func() {
			info.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
			Expect(info.CanRender(terminal.DefaultConfig)).To(BeTrue())
		})

		It("should return false when width below minimum", func() {
			info.Update(tea.WindowSizeMsg{Width: 30, Height: 24})
			Expect(info.CanRender(terminal.DefaultConfig)).To(BeFalse())
		})

		It("should return false when height below minimum", func() {
			info.Update(tea.WindowSizeMsg{Width: 80, Height: 10})
			Expect(info.CanRender(terminal.DefaultConfig)).To(BeFalse())
		})
	})

	Describe("ContentArea", func() {
		It("should calculate available space after margins", func() {
			info.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
			margins := terminal.Margins{Top: 2, Right: 4, Bottom: 2, Left: 4}

			width, height := info.ContentArea(margins)
			Expect(width).To(Equal(92))  // 100 - 4 - 4
			Expect(height).To(Equal(36)) // 40 - 2 - 2
		})

		It("should enforce minimum content area width", func() {
			info.Update(tea.WindowSizeMsg{Width: 30, Height: 40})
			margins := terminal.Margins{Top: 2, Right: 10, Bottom: 2, Left: 10}

			width, _ := info.ContentArea(margins)
			Expect(width).To(BeNumerically(">=", 20))
		})

		It("should enforce minimum content area height", func() {
			info.Update(tea.WindowSizeMsg{Width: 100, Height: 20})
			margins := terminal.Margins{Top: 10, Right: 4, Bottom: 10, Left: 4}

			_, height := info.ContentArea(margins)
			Expect(height).To(BeNumerically(">=", 5))
		})

		It("should use defaults when not valid", func() {
			margins := terminal.Margins{Top: 2, Right: 4, Bottom: 2, Left: 4}
			width, height := info.ContentArea(margins)

			Expect(width).To(Equal(72))  // 80 - 4 - 4
			Expect(height).To(Equal(20)) // 24 - 2 - 2
		})
	})
})

var _ = Describe("Config", func() {
	Describe("DefaultConfig", func() {
		It("should have sensible defaults", func() {
			config := terminal.DefaultConfig

			Expect(config.MinWidth).To(Equal(40))
			Expect(config.MinHeight).To(Equal(15))
			Expect(config.DefaultWidth).To(Equal(80))
			Expect(config.DefaultHeight).To(Equal(24))
		})
	})
})

var _ = Describe("Margins", func() {
	It("should define margin structure", func() {
		margins := terminal.Margins{
			Top:    2,
			Right:  4,
			Bottom: 2,
			Left:   4,
		}

		Expect(margins.Top).To(Equal(2))
		Expect(margins.Right).To(Equal(4))
		Expect(margins.Bottom).To(Equal(2))
		Expect(margins.Left).To(Equal(4))
	})
})

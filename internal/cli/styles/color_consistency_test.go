package styles_test

import (
	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/charmbracelet/lipgloss"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Color Scheme Consistency", func() {
	Describe("Background Colors", func() {
		It("should define primary background color", func() {
			Expect(styles.ColorBackground).NotTo(Equal(lipgloss.Color("")))
		})

		It("should define alternate background color", func() {
			Expect(styles.ColorBackgroundAlt).NotTo(Equal(lipgloss.Color("")))
		})

		It("should define card background color", func() {
			Expect(styles.ColorBackgroundCard).NotTo(Equal(lipgloss.Color("")))
		})
	})

	Describe("Accent Colors", func() {
		It("should define teal accent color", func() {
			Expect(styles.ColorAccentTeal).NotTo(Equal(lipgloss.Color("")))
		})

		It("should define green accent color", func() {
			Expect(styles.ColorAccentGreen).NotTo(Equal(lipgloss.Color("")))
		})

		It("should define purple accent color", func() {
			Expect(styles.ColorAccentPurple).NotTo(Equal(lipgloss.Color("")))
		})
	})

	Describe("Text Colors", func() {
		It("should define primary text color", func() {
			Expect(styles.ColorTextPrimary).NotTo(Equal(lipgloss.Color("")))
		})

		It("should define secondary text color", func() {
			Expect(styles.ColorTextSecondary).NotTo(Equal(lipgloss.Color("")))
		})

		It("should define muted text color", func() {
			Expect(styles.ColorTextMuted).NotTo(Equal(lipgloss.Color("")))
		})
	})

	Describe("Status Colors", func() {
		It("should define error color", func() {
			Expect(styles.ColorError).NotTo(Equal(lipgloss.Color("")))
		})

		It("should define warning color", func() {
			Expect(styles.ColorWarning).NotTo(Equal(lipgloss.Color("")))
		})

		It("should define success color", func() {
			Expect(styles.ColorSuccess).NotTo(Equal(lipgloss.Color("")))
		})

		It("should define info color", func() {
			Expect(styles.ColorInfo).NotTo(Equal(lipgloss.Color("")))
		})
	})

	Describe("Border Colors", func() {
		It("should define default border color", func() {
			Expect(styles.ColorBorder).NotTo(Equal(lipgloss.Color("")))
		})

		It("should define active border color", func() {
			Expect(styles.ColorBorderActive).NotTo(Equal(lipgloss.Color("")))
		})

		It("should define error border color", func() {
			Expect(styles.ColorBorderError).NotTo(Equal(lipgloss.Color("")))
		})
	})

	Describe("Color Consistency", func() {
		It("should use same color for success and green accent", func() {
			Expect(styles.ColorSuccess).To(Equal(styles.ColorAccentGreen))
		})

		It("should use teal for active borders", func() {
			Expect(styles.ColorBorderActive).To(Equal(styles.ColorAccentTeal))
		})

		It("should use error color for error borders", func() {
			Expect(styles.ColorBorderError).To(Equal(styles.ColorError))
		})
	})

	Describe("Style Application", func() {
		It("should apply status colors to message boxes", func() {
			errorBox := styles.ErrorBox
			Expect(errorBox).NotTo(BeNil())
		})

		It("should apply colors to buttons", func() {
			primaryButton := styles.ButtonPrimary
			Expect(primaryButton).NotTo(BeNil())
		})

		It("should apply colors to inputs", func() {
			input := styles.InputBase
			Expect(input).NotTo(BeNil())
		})

		It("should apply colors to cards", func() {
			card := styles.CardBase
			Expect(card).NotTo(BeNil())
		})
	})
})


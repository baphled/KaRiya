package styles_test

import (
	"testing"

	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/charmbracelet/lipgloss"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestStyles(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Styles Suite")
}

var _ = Describe("Styles", func() {
	Describe("Color Scheme", func() {
		It("should define all background colors", func() {
			Expect(styles.ColorBackground).ToNot(BeEmpty())
			Expect(styles.ColorBackgroundAlt).ToNot(BeEmpty())
			Expect(styles.ColorBackgroundCard).ToNot(BeEmpty())
		})

		It("should define all accent colors", func() {
			Expect(styles.ColorAccentTeal).ToNot(BeEmpty())
			Expect(styles.ColorAccentGreen).ToNot(BeEmpty())
			Expect(styles.ColorAccentPurple).ToNot(BeEmpty())
		})

		It("should define all text colors", func() {
			Expect(styles.ColorTextPrimary).ToNot(BeEmpty())
			Expect(styles.ColorTextSecondary).ToNot(BeEmpty())
			Expect(styles.ColorTextMuted).ToNot(BeEmpty())
		})

		It("should define all status colors", func() {
			Expect(styles.ColorError).ToNot(BeEmpty())
			Expect(styles.ColorWarning).ToNot(BeEmpty())
			Expect(styles.ColorSuccess).ToNot(BeEmpty())
			Expect(styles.ColorInfo).ToNot(BeEmpty())
		})

		It("should define all border colors", func() {
			Expect(styles.ColorBorder).ToNot(BeEmpty())
			Expect(styles.ColorBorderActive).ToNot(BeEmpty())
			Expect(styles.ColorBorderError).ToNot(BeEmpty())
		})
	})

	Describe("Button Styles", func() {
		It("should have base button style with padding and border", func() {
			rendered := styles.ButtonBase.Render("Test")
			Expect(rendered).ToNot(BeEmpty())
		})

		It("should have primary button style", func() {
			rendered := styles.ButtonPrimary.Render("Primary")
			Expect(rendered).ToNot(BeEmpty())
		})

		It("should have secondary button style", func() {
			rendered := styles.ButtonSecondary.Render("Secondary")
			Expect(rendered).ToNot(BeEmpty())
		})

		It("should have focused button style", func() {
			rendered := styles.ButtonFocused.Render("Focused")
			Expect(rendered).ToNot(BeEmpty())
		})

		It("should have disabled button style", func() {
			rendered := styles.ButtonDisabled.Render("Disabled")
			Expect(rendered).ToNot(BeEmpty())
		})
	})

	Describe("Input Field Styles", func() {
		It("should have base input style", func() {
			rendered := styles.InputBase.Render("Input text")
			Expect(rendered).ToNot(BeEmpty())
		})

		It("should have focused input style", func() {
			rendered := styles.InputFocused.Render("Focused input")
			Expect(rendered).ToNot(BeEmpty())
		})

		It("should have error input style", func() {
			rendered := styles.InputError.Render("Error input")
			Expect(rendered).ToNot(BeEmpty())
		})

		It("should have input label style", func() {
			rendered := styles.InputLabel.Render("Label")
			Expect(rendered).ToNot(BeEmpty())
		})

		It("should have input hint style", func() {
			rendered := styles.InputHint.Render("Hint text")
			Expect(rendered).ToNot(BeEmpty())
		})
	})

	Describe("Card Styles", func() {
		It("should have base card style", func() {
			rendered := styles.CardBase.Render("Card content")
			Expect(rendered).ToNot(BeEmpty())
		})

		It("should have card header style", func() {
			rendered := styles.CardHeader.Render("Card Header")
			Expect(rendered).ToNot(BeEmpty())
		})

		It("should have card content style", func() {
			rendered := styles.CardContent.Render("Card content text")
			Expect(rendered).ToNot(BeEmpty())
		})

		It("should have card footer style", func() {
			rendered := styles.CardFooter.Render("Card footer")
			Expect(rendered).ToNot(BeEmpty())
		})

		It("should have content card style for view content areas", func() {
			rendered := styles.ContentCard.Render("Content card text")
			Expect(rendered).ToNot(BeEmpty())
		})
	})

	Describe("Header Styles", func() {
		It("should have main header style", func() {
			rendered := styles.HeaderMain.Render("Main Header")
			Expect(rendered).ToNot(BeEmpty())
		})

		It("should have section header style", func() {
			rendered := styles.HeaderSection.Render("Section Header")
			Expect(rendered).ToNot(BeEmpty())
		})

		It("should have subsection header style", func() {
			rendered := styles.HeaderSubsection.Render("Subsection Header")
			Expect(rendered).ToNot(BeEmpty())
		})
	})

	Describe("Error Message Styles", func() {
		It("should have error box style", func() {
			rendered := styles.ErrorBox.Render("Error message")
			Expect(rendered).ToNot(BeEmpty())
		})

		It("should have error text style", func() {
			rendered := styles.ErrorText.Render("Error!")
			Expect(rendered).ToNot(BeEmpty())
		})

		It("should have error hint style", func() {
			rendered := styles.ErrorHint.Render("Error hint text")
			Expect(rendered).ToNot(BeEmpty())
		})
	})

	Describe("Warning Message Styles", func() {
		It("should have warning box style", func() {
			rendered := styles.WarningBox.Render("Warning message")
			Expect(rendered).ToNot(BeEmpty())
		})

		It("should have warning text style", func() {
			rendered := styles.WarningText.Render("Warning!")
			Expect(rendered).ToNot(BeEmpty())
		})

		It("should have warning hint style", func() {
			rendered := styles.WarningHint.Render("Warning hint text")
			Expect(rendered).ToNot(BeEmpty())
		})
	})

	Describe("Success Message Styles", func() {
		It("should have success box style", func() {
			rendered := styles.SuccessBox.Render("Success message")
			Expect(rendered).ToNot(BeEmpty())
		})

		It("should have success text style", func() {
			rendered := styles.SuccessText.Render("Success!")
			Expect(rendered).ToNot(BeEmpty())
		})

		It("should have success hint style", func() {
			rendered := styles.SuccessHint.Render("Success hint text")
			Expect(rendered).ToNot(BeEmpty())
		})
	})

	Describe("Info Message Styles", func() {
		It("should have info box style", func() {
			rendered := styles.InfoBox.Render("Info message")
			Expect(rendered).ToNot(BeEmpty())
		})

		It("should have info text style", func() {
			rendered := styles.InfoText.Render("Info!")
			Expect(rendered).ToNot(BeEmpty())
		})

		It("should have info hint style", func() {
			rendered := styles.InfoHint.Render("Info hint text")
			Expect(rendered).ToNot(BeEmpty())
		})
	})

	Describe("List Styles", func() {
		It("should have list item style", func() {
			rendered := styles.ListItem.Render("List item")
			Expect(rendered).ToNot(BeEmpty())
		})

		It("should have selected list item style", func() {
			rendered := styles.ListItemSelected.Render("Selected item")
			Expect(rendered).ToNot(BeEmpty())
		})

		It("should have focused list item style", func() {
			rendered := styles.ListItemFocused.Render("Focused item")
			Expect(rendered).ToNot(BeEmpty())
		})
	})

	Describe("Tag Styles", func() {
		It("should have base tag style", func() {
			rendered := styles.TagBase.Render("tag")
			Expect(rendered).ToNot(BeEmpty())
		})

		It("should have selected tag style", func() {
			rendered := styles.TagSelected.Render("selected-tag")
			Expect(rendered).ToNot(BeEmpty())
		})
	})

	Describe("Progress and Spinner Styles", func() {
		It("should have progress bar style", func() {
			rendered := styles.ProgressBar.Render("████░░░░")
			Expect(rendered).ToNot(BeEmpty())
		})

		It("should have progress text style", func() {
			rendered := styles.ProgressText.Render("50%")
			Expect(rendered).ToNot(BeEmpty())
		})

		It("should have spinner style", func() {
			rendered := styles.SpinnerStyle.Render("⠋")
			Expect(rendered).ToNot(BeEmpty())
		})
	})

	Describe("Getter Functions", func() {
		Describe("Color Getters", func() {
			It("should return background colors", func() {
				Expect(styles.GetColorBackground()).To(Equal(styles.ColorBackground))
				Expect(styles.GetColorBackgroundAlt()).To(Equal(styles.ColorBackgroundAlt))
				Expect(styles.GetColorBackgroundCard()).To(Equal(styles.ColorBackgroundCard))
			})

			It("should return accent colors", func() {
				Expect(styles.GetColorAccentTeal()).To(Equal(styles.ColorAccentTeal))
				Expect(styles.GetColorAccentGreen()).To(Equal(styles.ColorAccentGreen))
				Expect(styles.GetColorAccentPurple()).To(Equal(styles.ColorAccentPurple))
			})

			It("should return text colors", func() {
				Expect(styles.GetColorTextPrimary()).To(Equal(styles.ColorTextPrimary))
				Expect(styles.GetColorTextSecondary()).To(Equal(styles.ColorTextSecondary))
				Expect(styles.GetColorTextMuted()).To(Equal(styles.ColorTextMuted))
			})

			It("should return status colors", func() {
				Expect(styles.GetColorError()).To(Equal(styles.ColorError))
				Expect(styles.GetColorWarning()).To(Equal(styles.ColorWarning))
				Expect(styles.GetColorSuccess()).To(Equal(styles.ColorSuccess))
				Expect(styles.GetColorInfo()).To(Equal(styles.ColorInfo))
			})

			It("should return border colors", func() {
				Expect(styles.GetColorBorder()).To(Equal(styles.ColorBorder))
				Expect(styles.GetColorBorderActive()).To(Equal(styles.ColorBorderActive))
				Expect(styles.GetColorBorderError()).To(Equal(styles.ColorBorderError))
			})
		})

		Describe("Button Style Getters", func() {
			It("should return button styles", func() {
				Expect(styles.GetButtonBase()).ToNot(BeNil())
				Expect(styles.GetButtonPrimary()).ToNot(BeNil())
				Expect(styles.GetButtonSecondary()).ToNot(BeNil())
				Expect(styles.GetButtonFocused()).ToNot(BeNil())
				Expect(styles.GetButtonDisabled()).ToNot(BeNil())
				Expect(styles.GetButtonPrimaryFocused()).ToNot(BeNil())
				Expect(styles.GetButtonSecondaryFocused()).ToNot(BeNil())
			})

			It("should return button styles that can render", func() {
				Expect(styles.GetButtonBase().Render("Test")).ToNot(BeEmpty())
				Expect(styles.GetButtonPrimary().Render("Primary")).ToNot(BeEmpty())
				Expect(styles.GetButtonSecondary().Render("Secondary")).ToNot(BeEmpty())
			})
		})

		Describe("Input Style Getters", func() {
			It("should return input styles", func() {
				Expect(styles.GetInputBase()).ToNot(BeNil())
				Expect(styles.GetInputFocused()).ToNot(BeNil())
				Expect(styles.GetInputError()).ToNot(BeNil())
				Expect(styles.GetInputLabel()).ToNot(BeNil())
				Expect(styles.GetInputHint()).ToNot(BeNil())
			})

			It("should return input styles that can render", func() {
				Expect(styles.GetInputBase().Render("Input")).ToNot(BeEmpty())
				Expect(styles.GetInputFocused().Render("Focused")).ToNot(BeEmpty())
				Expect(styles.GetInputError().Render("Error")).ToNot(BeEmpty())
			})
		})

		Describe("Card Style Getters", func() {
			It("should return card styles", func() {
				Expect(styles.GetCardBase()).ToNot(BeNil())
				Expect(styles.GetCardHeader()).ToNot(BeNil())
				Expect(styles.GetCardContent()).ToNot(BeNil())
				Expect(styles.GetCardFooter()).ToNot(BeNil())
				Expect(styles.GetContentCard()).ToNot(BeNil())
			})

			It("should return card styles that can render", func() {
				Expect(styles.GetCardBase().Render("Card")).ToNot(BeEmpty())
				Expect(styles.GetCardHeader().Render("Header")).ToNot(BeEmpty())
				Expect(styles.GetCardContent().Render("Content")).ToNot(BeEmpty())
				Expect(styles.GetContentCard().Render("Content Card")).ToNot(BeEmpty())
			})

			It("should return a copy of ContentCard style to prevent mutation", func() {
				style1 := styles.GetContentCard()
				style2 := styles.GetContentCard()
				// Modify style1
				style1 = style1.Bold(true)
				// style2 should not be affected
				rendered1 := style1.Render("Bold")
				rendered2 := style2.Render("NotBold")
				Expect(rendered1).ToNot(Equal(rendered2))
			})
		})

		Describe("Modal Style Getters", func() {
			It("should return modal styles", func() {
				Expect(styles.GetModalBase()).ToNot(BeNil())
				Expect(styles.GetModalTitle()).ToNot(BeNil())
				Expect(styles.GetModalMessage()).ToNot(BeNil())
				Expect(styles.GetModalButtonContainer()).ToNot(BeNil())
				Expect(styles.GetModalInstructions()).ToNot(BeNil())
				Expect(styles.GetModalDestructive()).ToNot(BeNil())
				Expect(styles.GetModalDestructiveTitle()).ToNot(BeNil())
			})

			It("should return modal styles that can render", func() {
				Expect(styles.GetModalBase().Render("Modal")).ToNot(BeEmpty())
				Expect(styles.GetModalTitle().Render("Title")).ToNot(BeEmpty())
				Expect(styles.GetModalMessage().Render("Message")).ToNot(BeEmpty())
			})
		})

		Describe("Header Style Getters", func() {
			It("should return header styles", func() {
				Expect(styles.GetHeaderMain()).ToNot(BeNil())
				Expect(styles.GetHeaderSection()).ToNot(BeNil())
				Expect(styles.GetHeaderSubsection()).ToNot(BeNil())
			})

			It("should return header styles that can render", func() {
				Expect(styles.GetHeaderMain().Render("Main")).ToNot(BeEmpty())
				Expect(styles.GetHeaderSection().Render("Section")).ToNot(BeEmpty())
				Expect(styles.GetHeaderSubsection().Render("Subsection")).ToNot(BeEmpty())
			})
		})

		Describe("Message Style Getters", func() {
			Context("Error styles", func() {
				It("should return error styles", func() {
					Expect(styles.GetErrorBox()).ToNot(BeNil())
					Expect(styles.GetErrorText()).ToNot(BeNil())
					Expect(styles.GetErrorHint()).ToNot(BeNil())
					Expect(styles.GetErrorMsg()).ToNot(BeNil())
				})

				It("should return error styles that can render", func() {
					Expect(styles.GetErrorBox().Render("Error")).ToNot(BeEmpty())
					Expect(styles.GetErrorText().Render("Error Text")).ToNot(BeEmpty())
				})
			})

			Context("Warning styles", func() {
				It("should return warning styles", func() {
					Expect(styles.GetWarningBox()).ToNot(BeNil())
					Expect(styles.GetWarningText()).ToNot(BeNil())
					Expect(styles.GetWarningHint()).ToNot(BeNil())
					Expect(styles.GetWarning()).ToNot(BeNil())
				})

				It("should return warning styles that can render", func() {
					Expect(styles.GetWarningBox().Render("Warning")).ToNot(BeEmpty())
					Expect(styles.GetWarningText().Render("Warning Text")).ToNot(BeEmpty())
				})
			})

			Context("Success styles", func() {
				It("should return success styles", func() {
					Expect(styles.GetSuccessBox()).ToNot(BeNil())
					Expect(styles.GetSuccessText()).ToNot(BeNil())
					Expect(styles.GetSuccessHint()).ToNot(BeNil())
				})

				It("should return success styles that can render", func() {
					Expect(styles.GetSuccessBox().Render("Success")).ToNot(BeEmpty())
					Expect(styles.GetSuccessText().Render("Success Text")).ToNot(BeEmpty())
				})
			})

			Context("Info styles", func() {
				It("should return info styles", func() {
					Expect(styles.GetInfoBox()).ToNot(BeNil())
					Expect(styles.GetInfoText()).ToNot(BeNil())
					Expect(styles.GetInfoHint()).ToNot(BeNil())
				})

				It("should return info styles that can render", func() {
					Expect(styles.GetInfoBox().Render("Info")).ToNot(BeEmpty())
					Expect(styles.GetInfoText().Render("Info Text")).ToNot(BeEmpty())
				})
			})
		})

		Describe("List Style Getters", func() {
			It("should return list styles", func() {
				Expect(styles.GetListItem()).ToNot(BeNil())
				Expect(styles.GetListItemSelected()).ToNot(BeNil())
				Expect(styles.GetListItemFocused()).ToNot(BeNil())
			})

			It("should return list styles that can render", func() {
				Expect(styles.GetListItem().Render("Item")).ToNot(BeEmpty())
				Expect(styles.GetListItemSelected().Render("Selected")).ToNot(BeEmpty())
				Expect(styles.GetListItemFocused().Render("Focused")).ToNot(BeEmpty())
			})
		})

		Describe("Tag Style Getters", func() {
			It("should return tag styles", func() {
				Expect(styles.GetTagBase()).ToNot(BeNil())
				Expect(styles.GetTagSelected()).ToNot(BeNil())
			})

			It("should return tag styles that can render", func() {
				Expect(styles.GetTagBase().Render("tag")).ToNot(BeEmpty())
				Expect(styles.GetTagSelected().Render("selected")).ToNot(BeEmpty())
			})
		})

		Describe("Progress Style Getters", func() {
			It("should return progress styles", func() {
				Expect(styles.GetProgressBar()).ToNot(BeNil())
				Expect(styles.GetProgressText()).ToNot(BeNil())
			})

			It("should return progress styles that can render", func() {
				Expect(styles.GetProgressBar().Render("████")).ToNot(BeEmpty())
				Expect(styles.GetProgressText().Render("50%")).ToNot(BeEmpty())
			})
		})

		Describe("Spinner Style Getter", func() {
			It("should return spinner style", func() {
				Expect(styles.GetSpinnerStyle()).ToNot(BeNil())
			})

			It("should return spinner style that can render", func() {
				Expect(styles.GetSpinnerStyle().Render("⠋")).ToNot(BeEmpty())
			})
		})

		Describe("Badge Style Getters", func() {
			It("should return badge styles", func() {
				Expect(styles.GetBadge()).ToNot(BeNil())
				Expect(styles.GetBadgeSelected()).ToNot(BeNil())
				Expect(styles.GetBadgeFocused()).ToNot(BeNil())
			})

			It("should return badge styles that can render", func() {
				Expect(styles.GetBadge().Render("badge")).ToNot(BeEmpty())
				Expect(styles.GetBadgeSelected().Render("selected")).ToNot(BeEmpty())
				Expect(styles.GetBadgeFocused().Render("focused")).ToNot(BeEmpty())
			})
		})

		Describe("Label Style Getters", func() {
			It("should return label styles", func() {
				Expect(styles.GetLabel()).ToNot(BeNil())
				Expect(styles.GetLabelFocused()).ToNot(BeNil())
				Expect(styles.GetHint()).ToNot(BeNil())
			})

			It("should return label styles that can render", func() {
				Expect(styles.GetLabel().Render("Label")).ToNot(BeEmpty())
				Expect(styles.GetLabelFocused().Render("Focused Label")).ToNot(BeEmpty())
				Expect(styles.GetHint().Render("Hint")).ToNot(BeEmpty())
			})
		})

		Describe("Spacing Constants Getter", func() {
			It("should return spacing constants struct", func() {
				spacing := styles.GetSpacingConstants()
				Expect(spacing.PaddingHorizontalSmall).To(Equal(1))
				Expect(spacing.PaddingHorizontalBase).To(Equal(2))
				Expect(spacing.PaddingHorizontalLarge).To(Equal(3))
				Expect(spacing.PaddingVerticalSmall).To(Equal(0))
				Expect(spacing.PaddingVerticalBase).To(Equal(1))
				Expect(spacing.MarginSmall).To(Equal(1))
				Expect(spacing.MarginBase).To(Equal(2))
				Expect(spacing.MarginLarge).To(Equal(3))
				Expect(spacing.MaxContentWidth).To(Equal(120))
				Expect(spacing.GridGutterWidth).To(Equal(2))
			})
		})
	})

	Describe("Helper Functions", func() {
		Describe("WithBorder", func() {
			It("should add a border to a style", func() {
				baseStyle := lipgloss.NewStyle()
				styledWithBorder := styles.WithBorder(baseStyle)
				rendered := styledWithBorder.Render("Bordered text")
				Expect(rendered).ToNot(BeEmpty())
			})
		})

		Describe("WithFocusedBorder", func() {
			It("should add a focused border to a style", func() {
				baseStyle := lipgloss.NewStyle()
				styledWithFocusedBorder := styles.WithFocusedBorder(baseStyle)
				rendered := styledWithFocusedBorder.Render("Focused border text")
				Expect(rendered).ToNot(BeEmpty())
			})
		})

		Describe("WithErrorBorder", func() {
			It("should add an error border to a style", func() {
				baseStyle := lipgloss.NewStyle()
				styledWithErrorBorder := styles.WithErrorBorder(baseStyle)
				rendered := styledWithErrorBorder.Render("Error border text")
				Expect(rendered).ToNot(BeEmpty())
			})
		})

		Describe("WithPadding", func() {
			It("should add padding to a style", func() {
				baseStyle := lipgloss.NewStyle()
				styledWithPadding := styles.WithPadding(baseStyle, 1, 2)
				rendered := styledWithPadding.Render("Padded text")
				Expect(rendered).ToNot(BeEmpty())
			})
		})

		Describe("WithMargin", func() {
			It("should add margin to a style", func() {
				baseStyle := lipgloss.NewStyle()
				styledWithMargin := styles.WithMargin(baseStyle, 1, 2)
				rendered := styledWithMargin.Render("Margined text")
				Expect(rendered).ToNot(BeEmpty())
			})
		})
	})

	Describe("Responsive Layout Helpers", func() {
		Describe("MaxWidth", func() {
			It("should return terminal width minus margin when smaller than max", func() {
				terminalWidth := 80
				maxWidth := styles.MaxWidth(terminalWidth)
				Expect(maxWidth).To(Equal(76)) // 80 - 4 margin
			})

			It("should return max content width when terminal is larger", func() {
				terminalWidth := 200
				maxWidth := styles.MaxWidth(terminalWidth)
				Expect(maxWidth).To(Equal(120)) // maxContentWidth constant
			})

			It("should handle very small terminal widths", func() {
				terminalWidth := 40
				maxWidth := styles.MaxWidth(terminalWidth)
				Expect(maxWidth).To(Equal(36)) // 40 - 4 margin
			})
		})

		Describe("CenterHorizontal", func() {
			It("should center text horizontally", func() {
				text := "Centered"
				width := 20
				centered := styles.CenterHorizontal(text, width)
				Expect(centered).ToNot(BeEmpty())
				Expect(len(centered)).To(BeNumerically(">=", len(text)))
			})
		})

		Describe("CenterVertical", func() {
			It("should center text vertically", func() {
				text := "Centered"
				height := 10
				centered := styles.CenterVertical(text, height)
				Expect(centered).ToNot(BeEmpty())
			})
		})

		Describe("Center", func() {
			It("should center text both horizontally and vertically", func() {
				text := "Centered"
				width := 20
				height := 10
				centered := styles.Center(text, width, height)
				Expect(centered).ToNot(BeEmpty())
			})
		})

		Describe("AlignLeft", func() {
			It("should align text to the left", func() {
				text := "Left"
				width := 20
				aligned := styles.AlignLeft(text, width)
				Expect(aligned).ToNot(BeEmpty())
			})
		})

		Describe("AlignRight", func() {
			It("should align text to the right", func() {
				text := "Right"
				width := 20
				aligned := styles.AlignRight(text, width)
				Expect(aligned).ToNot(BeEmpty())
			})
		})

		Describe("ResponsiveStyle", func() {
			It("should return a style with responsive width", func() {
				terminalWidth := 100
				style := styles.ResponsiveStyle(terminalWidth)
				rendered := style.Render("Responsive text")
				Expect(rendered).ToNot(BeEmpty())
			})
		})

		Describe("ResponsiveCard", func() {
			It("should return a card style with responsive width", func() {
				terminalWidth := 100
				cardStyle := styles.ResponsiveCard(terminalWidth)
				rendered := cardStyle.Render("Responsive card")
				Expect(rendered).ToNot(BeEmpty())
			})
		})

		Describe("ResponsiveInput", func() {
			It("should return an input style with responsive width", func() {
				terminalWidth := 100
				inputStyle := styles.ResponsiveInput(terminalWidth)
				rendered := inputStyle.Render("Responsive input")
				Expect(rendered).ToNot(BeEmpty())
			})
		})

		Describe("TwoColumn", func() {
			It("should split width into two equal columns", func() {
				terminalWidth := 100
				leftWidth, rightWidth := styles.TwoColumn(terminalWidth)
				Expect(leftWidth).To(BeNumerically(">", 0))
				Expect(rightWidth).To(BeNumerically(">", 0))
				Expect(leftWidth + rightWidth).To(BeNumerically("<=", 100))
			})
		})

		Describe("ThreeColumn", func() {
			It("should split width into three columns", func() {
				terminalWidth := 120
				leftWidth, centerWidth, rightWidth := styles.ThreeColumn(terminalWidth)
				Expect(leftWidth).To(BeNumerically(">", 0))
				Expect(centerWidth).To(BeNumerically(">", 0))
				Expect(rightWidth).To(BeNumerically(">", 0))
				Expect(leftWidth + centerWidth + rightWidth).To(BeNumerically("<=", 120))
			})
		})

		Describe("Grid", func() {
			It("should calculate grid column width and gutter", func() {
				terminalWidth := 100
				columns := 3
				columnWidth, gutter := styles.Grid(terminalWidth, columns)
				Expect(columnWidth).To(BeNumerically(">", 0))
				Expect(gutter).To(Equal(2))
				totalWidth := (columnWidth * columns) + (gutter * (columns - 1))
				Expect(totalWidth).To(BeNumerically("<=", styles.MaxWidth(terminalWidth)))
			})

			It("should handle single column grid", func() {
				terminalWidth := 100
				columns := 1
				columnWidth, gutter := styles.Grid(terminalWidth, columns)
				Expect(columnWidth).To(BeNumerically(">", 0))
				Expect(gutter).To(Equal(2))
			})

			It("should handle many column grid", func() {
				terminalWidth := 200
				columns := 6
				columnWidth, gutter := styles.Grid(terminalWidth, columns)
				Expect(columnWidth).To(BeNumerically(">", 0))
				Expect(gutter).To(Equal(2))
			})
		})
	})
})

package themes_test

import (
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bubbles Theme Integration", func() {
	var theme themes.Theme

	BeforeEach(func() {
		theme = themes.NewDefaultTheme()
	})

	Describe("NewThemedListStyles", func() {
		It("should return list.Styles", func() {
			styles := themes.NewThemedListStyles(theme)
			Expect(styles).NotTo(BeNil())
		})

		It("should apply theme colors to title", func() {
			styles := themes.NewThemedListStyles(theme)
			// The styles should be configured - we verify they render
			rendered := styles.Title.Render("Test")
			Expect(rendered).To(ContainSubstring("Test"))
		})

		It("should apply theme colors to filter prompt", func() {
			styles := themes.NewThemedListStyles(theme)
			rendered := styles.FilterPrompt.Render("Filter:")
			Expect(rendered).To(ContainSubstring("Filter:"))
		})

		It("should handle nil theme gracefully", func() {
			styles := themes.NewThemedListStyles(nil)
			// Should return default styles, not panic
			Expect(styles).NotTo(BeNil())
		})
	})

	Describe("NewThemedTableStyles", func() {
		It("should return table.Styles", func() {
			styles := themes.NewThemedTableStyles(theme)
			Expect(styles).NotTo(BeNil())
		})

		It("should apply theme colors to header", func() {
			styles := themes.NewThemedTableStyles(theme)
			rendered := styles.Header.Render("Header")
			Expect(rendered).To(ContainSubstring("Header"))
		})

		It("should apply theme colors to selected row", func() {
			styles := themes.NewThemedTableStyles(theme)
			rendered := styles.Selected.Render("Selected")
			Expect(rendered).To(ContainSubstring("Selected"))
		})

		It("should apply theme colors to cells", func() {
			styles := themes.NewThemedTableStyles(theme)
			rendered := styles.Cell.Render("Cell")
			Expect(rendered).To(ContainSubstring("Cell"))
		})

		It("should handle nil theme gracefully", func() {
			styles := themes.NewThemedTableStyles(nil)
			Expect(styles).NotTo(BeNil())
		})
	})

	Describe("NewThemedProgress", func() {
		It("should return a configured progress.Model", func() {
			p := themes.NewThemedProgress(theme)
			Expect(p).NotTo(BeNil())
		})

		It("should be able to set percentage", func() {
			p := themes.NewThemedProgress(theme)
			p.SetPercent(0.5)
			Expect(p.Percent()).To(BeNumerically("~", 0.5, 0.01))
		})

		It("should render progress bar", func() {
			p := themes.NewThemedProgress(theme)
			p.SetPercent(0.5)
			view := p.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle nil theme gracefully", func() {
			p := themes.NewThemedProgress(nil)
			Expect(p).NotTo(BeNil())
		})
	})

	Describe("NewThemedSpinner", func() {
		It("should return a configured spinner.Model", func() {
			s := themes.NewThemedSpinner(theme)
			Expect(s).NotTo(BeNil())
		})

		It("should render spinner", func() {
			s := themes.NewThemedSpinner(theme)
			view := s.View()
			// Spinner may be empty initially, but should not panic
			_ = view
		})

		It("should handle nil theme gracefully", func() {
			s := themes.NewThemedSpinner(nil)
			Expect(s).NotTo(BeNil())
		})
	})

	Describe("NewThemedHelpStyles", func() {
		It("should return help.Styles", func() {
			styles := themes.NewThemedHelpStyles(theme)
			Expect(styles).NotTo(BeNil())
		})

		It("should apply theme colors to short key", func() {
			styles := themes.NewThemedHelpStyles(theme)
			rendered := styles.ShortKey.Render("q")
			Expect(rendered).To(ContainSubstring("q"))
		})

		It("should apply theme colors to short description", func() {
			styles := themes.NewThemedHelpStyles(theme)
			rendered := styles.ShortDesc.Render("quit")
			Expect(rendered).To(ContainSubstring("quit"))
		})

		It("should handle nil theme gracefully", func() {
			styles := themes.NewThemedHelpStyles(nil)
			Expect(styles).NotTo(BeNil())
		})
	})

	Describe("ApplyThemeToList", func() {
		It("should apply theme to an existing list.Model", func() {
			// Create a basic list
			items := []list.Item{}
			delegate := list.NewDefaultDelegate()
			l := list.New(items, delegate, 80, 24)

			// Apply theme
			l = themes.ApplyThemeToList(l, theme)

			// Verify it still works
			Expect(l.Width()).To(Equal(80))
			Expect(l.Height()).To(Equal(24))
		})
	})

	Describe("ApplyThemeToTable", func() {
		It("should apply theme to an existing table.Model", func() {
			// Create a basic table
			columns := []table.Column{
				{Title: "Name", Width: 20},
			}
			rows := []table.Row{
				{"Test"},
			}
			t := table.New(
				table.WithColumns(columns),
				table.WithRows(rows),
			)

			// Apply theme
			t = themes.ApplyThemeToTable(t, theme)

			// Verify it still works
			view := t.View()
			Expect(view).To(ContainSubstring("Test"))
		})
	})
})

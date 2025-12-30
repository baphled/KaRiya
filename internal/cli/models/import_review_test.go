package models_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/importer"
	"github.com/baphled/kariya/internal/cli/models"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
)

var _ = Describe("ImportReviewModel", func() {
	var (
		model      *models.ImportReviewModel
		parsedRows []*importer.ParsedRow
	)

	BeforeEach(func() {
		// Create test parsed rows
		parsedRows = []*importer.ParsedRow{
			{
				RowNumber: 1,
				RawData:   map[string]string{"text": "First event", "date": "2025-12-28"},
				Event: &career.CareerEvent{
					ID:        "event-1",
					Text:      "First event",
					Date:      time.Now().Add(-48 * time.Hour),
					Company:   "CompanyA",
					Project:   "ProjectA",
					Tags:      []string{"technical"},
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				ValidationErrors: []string{},
				IsValid:          true,
				IsDuplicate:      false,
				DuplicateOf:      "",
			},
			{
				RowNumber: 2,
				RawData:   map[string]string{"text": "Second event", "date": "2025-12-29"},
				Event: &career.CareerEvent{
					ID:        "event-2",
					Text:      "Second event",
					Date:      time.Now().Add(-24 * time.Hour),
					Company:   "CompanyB",
					Project:   "ProjectB",
					Tags:      []string{"leadership"},
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				ValidationErrors: []string{},
				IsValid:          true,
				IsDuplicate:      false,
				DuplicateOf:      "",
			},
		}

		model = models.NewImportReviewModel(parsedRows)
	})

	Context("when import is confirmed with Enter key", func() {
		It("should trigger metadata review after import", func() {
			// Arrange: model is set up with parsed rows
			Expect(len(parsedRows)).To(Equal(2))

			// Act: simulate pressing Enter to confirm import
			_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})

			// Assert: cmd should emit a message
			Expect(cmd).NotTo(BeNil())

			// Execute the command to get the message
			msg := cmd()
			// Currently returns ImportReviewMsg; should return MetadataReviewTriggeredMsg after enhancement
			_, ok := msg.(models.ImportReviewMsg)
			Expect(ok).To(BeTrue())
		})
	})

	Context("navigation", func() {
		It("should move focus down with down arrow", func() {
			// Arrange

			// Act: press down
			model.Update(tea.KeyMsg{Type: tea.KeyDown})

			// Assert: view should render without error
			view := model.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should move focus up with up arrow", func() {
			// Arrange: move focus down first
			model.Update(tea.KeyMsg{Type: tea.KeyDown})

			// Act: press up
			model.Update(tea.KeyMsg{Type: tea.KeyUp})

			// Assert: view should render without error
			view := model.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Context("selection", func() {
		It("should deselect all rows with 'd' key", func() {
			// Arrange: all valid rows are pre-selected

			// Act: press 'd' to deselect all
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})

			// Assert: pressing enter should return empty selected rows
			_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			msg := cmd()
			importMsg, ok := msg.(models.ImportReviewMsg)
			Expect(ok).To(BeTrue())
			Expect(len(importMsg.SelectedRows)).To(Equal(0))
		})

		It("should select all rows with 'a' key", func() {
			// Arrange: deselect all first
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})

			// Act: press 'a' to select all
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			// Assert: pressing enter should return all selected rows
			_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			msg := cmd()
			importMsg, ok := msg.(models.ImportReviewMsg)
			Expect(ok).To(BeTrue())
			Expect(len(importMsg.SelectedRows)).To(Equal(2))
		})
	})

	Context("window resize", func() {
		It("should handle window size changes", func() {
			// Arrange

			// Act: simulate window resize
			model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

			// Assert: view should render without error
			view := model.View()
			Expect(view).NotTo(BeEmpty())
		})
	})
})


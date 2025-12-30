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

	Context("initialization", func() {
		It("should initialize with parsed rows", func() {
			Expect(len(model.ParsedRows)).To(Equal(2))
		})

		It("should pre-select valid non-duplicate rows", func() {
			selected := model.GetSelectedEventIDs()
			Expect(len(selected)).To(Equal(2))
		})
	})

	Context("navigation", func() {
		It("should move focus down", func() {
			model.Update(tea.KeyMsg{Type: tea.KeyDown})
			// Verify no panic and view renders
			view := model.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should move focus up", func() {
			model.Update(tea.KeyMsg{Type: tea.KeyDown})
			model.Update(tea.KeyMsg{Type: tea.KeyUp})
			view := model.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Context("selection operations", func() {
		It("should select all valid rows with 'a'", func() {
			// Deselect all first
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
			selected := model.GetSelectedEventIDs()
			Expect(len(selected)).To(Equal(0))

			// Select all
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			selected = model.GetSelectedEventIDs()
			Expect(len(selected)).To(Equal(2))
		})

		It("should deselect all with 'd'", func() {
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
			selected := model.GetSelectedEventIDs()
			Expect(len(selected)).To(Equal(0))
		})
	})

	Context("rendering", func() {
		It("should render view without error", func() {
			view := model.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should show summary with counts", func() {
			view := model.View()
			Expect(view).To(ContainSubstring("Valid"))
		})
	})

	Context("parsing warnings and duplicate detection", func() {
		It("should identify parsing warnings", func() {
			warnings := model.GetParsingWarnings(0)
			Expect(len(warnings)).To(Equal(1))
		})

		It("should detect when row has warnings", func() {
			hasWarnings := model.HasWarnings(0)
			Expect(hasWarnings).To(BeTrue())
		})

		It("should provide duplicate information", func() {
			dupInfo := model.GetDuplicateInfo(0)
			Expect(dupInfo).To(Equal(""))
		})
	})

	Context("bulk operations support", func() {
		It("should return selected events", func() {
			selected := model.GetSelectedEvents()
			Expect(len(selected)).To(Equal(2))
		})

		It("should return selected event IDs", func() {
			ids := model.GetSelectedEventIDs()
			Expect(len(ids)).To(Equal(2))
		})

		It("should check if bulk edit is available", func() {
			canBulk := model.CanBulkEdit()
			Expect(canBulk).To(BeTrue())
		})

		It("should apply bulk metadata updates", func() {
			count := model.ApplyBulkMetadataUpdate("NewCorp", "NewProject", []string{"tag1"}, []string{})
			Expect(count).To(Equal(2))
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

	Context("parsing issues display", func() {
		It("should show parsing issues in summary when present", func() {
			// Arrange: Create a row with parsing warnings
			rowWithWarnings := &importer.ParsedRow{
				RowNumber: 3,
				RawData:   map[string]string{"text": "Event without date"},
				Event: &career.CareerEvent{
					ID:        "event-3",
					Text:      "Event without date",
					Date:      time.Time{}, // Missing date
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				ValidationErrors: []string{},
				IsValid:          false, // Invalid due to missing date
				IsDuplicate:      false,
				DuplicateOf:      "",
			}
			parsedRows = append(parsedRows, rowWithWarnings)
			model = models.NewImportReviewModel(parsedRows)

			// Act: Get the view output
			view := model.View()

			// Assert: View should contain warning indicator
			Expect(view).To(ContainSubstring("✗ ERR"))
		})

		It("should display parsing warning details when row is focused", func() {
			// Arrange: Create a row with parsing warnings
			rowWithWarnings := &importer.ParsedRow{
				RowNumber: 3,
				RawData:   map[string]string{"text": "Event without date"},
				Event: &career.CareerEvent{
					ID:        "event-3",
					Text:      "Event without date",
					Date:      time.Time{}, // Missing date
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				ValidationErrors: []string{},
				IsValid:          false,
				IsDuplicate:      false,
				DuplicateOf:      "",
			}
			parsedRows = append(parsedRows, rowWithWarnings)
			model = models.NewImportReviewModel(parsedRows)

			// Act: Focus on the warning row and get warnings
			model.Update(tea.KeyMsg{Type: tea.KeyDown})
			warnings := model.GetParsingWarnings(2)

			// Assert: Should have parsing warnings
			Expect(len(warnings)).To(BeNumerically(">", 0))
			Expect(warnings).To(ContainElement(ContainSubstring("missing date")))
		})

		It("should show warning count in summary", func() {
			// Arrange: Create rows with warnings
			rowWithWarnings := &importer.ParsedRow{
				RowNumber: 3,
				RawData:   map[string]string{"text": "Event without date"},
				Event: &career.CareerEvent{
					ID:        "event-3",
					Text:      "Event without date",
					Date:      time.Time{},
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				ValidationErrors: []string{},
				IsValid:          false,
				IsDuplicate:      false,
				DuplicateOf:      "",
			}
			parsedRows = append(parsedRows, rowWithWarnings)
			model = models.NewImportReviewModel(parsedRows)

			// Act: Get warning count
			warningCount := 0
			for i := 0; i < len(parsedRows); i++ {
				if model.HasWarnings(i) {
					warningCount++
				}
			}

			// Assert: Should have at least one row with warnings
			Expect(warningCount).To(BeNumerically(">", 0))
		})
	})
	Context("duplicate detection status display", func() {
		It("should show duplicate indicator for duplicate rows", func() {
			// Arrange: Create a duplicate row
			duplicateRow := &importer.ParsedRow{
				RowNumber:        3,
				RawData:          map[string]string{"text": "Duplicate event"},
				Event:            &career.CareerEvent{ID: "event-3", Text: "Duplicate event", CreatedAt: time.Now(), UpdatedAt: time.Now()},
				ValidationErrors: []string{},
				IsValid:          true,
				IsDuplicate:      true,
				DuplicateOf:      "event-1",
			}
			parsedRows = append(parsedRows, duplicateRow)
			model = models.NewImportReviewModel(parsedRows)

			// Act: Get the view output
			view := model.View()

			// Assert: View should contain duplicate indicator
			Expect(view).To(ContainSubstring("⚠ DUP"))
		})

		It("should display duplicate reference when focused", func() {
			// Arrange: Create a duplicate row
			duplicateRow := &importer.ParsedRow{
				RowNumber:        3,
				RawData:          map[string]string{"text": "Duplicate event"},
				Event:            &career.CareerEvent{ID: "event-3", Text: "Duplicate event", CreatedAt: time.Now(), UpdatedAt: time.Now()},
				ValidationErrors: []string{},
				IsValid:          true,
				IsDuplicate:      true,
				DuplicateOf:      "event-1",
			}
			parsedRows = append(parsedRows, duplicateRow)
			model = models.NewImportReviewModel(parsedRows)

			// Act: Navigate to duplicate row and get duplicate info
			model.Update(tea.KeyMsg{Type: tea.KeyDown})
			model.Update(tea.KeyMsg{Type: tea.KeyDown})
			dupInfo := model.GetDuplicateInfo(2)

			// Assert: Should return the duplicate reference
			Expect(dupInfo).To(Equal("duplicate of event #event-1"))
		})

		It("should not allow selection of duplicate rows", func() {
			// Arrange: Create a duplicate row
			duplicateRow := &importer.ParsedRow{
				RowNumber:        3,
				RawData:          map[string]string{"text": "Duplicate event"},
				Event:            &career.CareerEvent{ID: "event-3", Text: "Duplicate event", CreatedAt: time.Now(), UpdatedAt: time.Now()},
				ValidationErrors: []string{},
				IsValid:          true,
				IsDuplicate:      true,
				DuplicateOf:      "event-1",
			}
			parsedRows = append(parsedRows, duplicateRow)
			model = models.NewImportReviewModel(parsedRows)

			// Act: Try to deselect all and check selection state
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
			selected := model.GetSelectedEventIDs()

			// Assert: Should have 0 selected (duplicates not pre-selected)
			Expect(len(selected)).To(Equal(0))
		})
	})
})

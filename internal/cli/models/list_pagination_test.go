package models

import (
	"context"
	"fmt"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ListModel Pagination Bug", func() {
	var (
		repo  *careerrepo.MemoryRepository
		svc   *careerservice.Service
		ctx   context.Context
		model *ListModel
	)

	BeforeEach(func() {
		repo = careerrepo.NewMemoryRepository()
		svc = careerservice.NewService(repo)
		ctx = context.Background()
	})

	Context("when table pagination updates", func() {
		BeforeEach(func() {
			// Create 25 events
			for i := 1; i <= 25; i++ {
				event := &career.CareerEvent{
					Text: fmt.Sprintf("Item%02d", i),
					Date: time.Now().Add(-time.Duration(i*24) * time.Hour),
				}
				err := svc.CaptureEvent(ctx, event, careerservice.TimelineJournaling)
				Expect(err).NotTo(HaveOccurred())
			}

			model = NewListModel(svc, ctx)
		})

		It("should have cursor visible on page 1", func() {
			view := model.View()
			// The cursor indicator should be visible
			Expect(view).To(ContainSubstring("▶"))
		})

		It("should have cursor visible on page 2 after navigation", func() {
			// Navigate to page 2
			model.NextPage()
			
			view := model.View()
			
			// The cursor indicator should still be visible on page 2
			// BUG: The cursor might not be visible after pagination
			Expect(view).To(ContainSubstring("▶"))
		})

		It("table cursor should be in sync with pagination after navigation", func() {
			// Navigate to next page
			model.NextPage()
			
			// After navigation, the table should have updated
			// The table's cursor should still point to a valid row
			tableRowCount := len(model.table.Rows())
			tableCursor := model.table.Cursor()
			
			// The table cursor should be less than the number of rows
			// BUG: The cursor might be out of sync with the table rows
			Expect(tableCursor).To(BeNumerically("<", tableRowCount), "Table cursor should be within valid row range")
			Expect(model.pagination.GetCurrentPage()).To(Equal(2), "Should be on page 2")
		})

		It("table rows should change when navigating pages", func() {
			// Get rows from page 1
			page1Rows := model.table.Rows()
			page1RowCount := len(page1Rows)
			
			// Get the text from first row of page 1
			var page1FirstRowText string
			if page1RowCount > 0 {
				page1FirstRowText = page1Rows[0][0]
			}
			
			// Navigate to page 2
			model.NextPage()
			
			// Get rows from page 2
			page2Rows := model.table.Rows()
			page2RowCount := len(page2Rows)
			
			// Get the text from first row of page 2
			var page2FirstRowText string
			if page2RowCount > 0 {
				page2FirstRowText = page2Rows[0][0]
			}
			
			// The first row should be different
			// BUG: The table rows might not update, so they could be the same
			Expect(page2FirstRowText).NotTo(Equal(page1FirstRowText), "First row of page 2 should be different from page 1")
		})

		It("should display different items in view after pagination", func() {
			// Get page 1 view
			view1 := model.View()
			Expect(view1).To(ContainSubstring("Item01"))
			
			// Navigate to page 2
			model.NextPage()
			
			// Get page 2 view
			view2 := model.View()
			
			// Page 2 should show items from page 2, not page 1
			// BUG: The view might not update and still show page 1 items
			Expect(view2).To(ContainSubstring("Item11"))
			Expect(view2).NotTo(ContainSubstring("Item01"))
		})
	})
})

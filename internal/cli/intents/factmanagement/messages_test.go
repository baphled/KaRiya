package factmanagement_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents/factmanagement"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
)

var _ = Describe("Messages", func() {
	Describe("FactSelectedMsg", func() {
		It("should store fact and index", func() {
			fact := fixtures.FactWith("fact-1", "Test fact")
			msg := factmanagement.FactSelectedMsg{
				Fact:  fact,
				Index: 5,
			}
			Expect(msg.Fact).To(Equal(fact))
			Expect(msg.Index).To(Equal(5))
		})

		It("should allow nil fact", func() {
			msg := factmanagement.FactSelectedMsg{
				Fact:  nil,
				Index: -1,
			}
			Expect(msg.Fact).To(BeNil())
			Expect(msg.Index).To(Equal(-1))
		})
	})

	Describe("FactSavedMsg", func() {
		It("should store fact and metadata", func() {
			fact := fixtures.FactWith("fact-1", "Test fact")
			msg := factmanagement.FactSavedMsg{
				Fact:    fact,
				IsNew:   true,
				Message: "Fact created successfully",
			}
			Expect(msg.Fact).To(Equal(fact))
			Expect(msg.IsNew).To(BeTrue())
			Expect(msg.Message).To(Equal("Fact created successfully"))
		})

		It("should indicate update when not new", func() {
			fact := fixtures.FactWith("fact-1", "Updated fact")
			msg := factmanagement.FactSavedMsg{
				Fact:    fact,
				IsNew:   false,
				Message: "Fact updated successfully",
			}
			Expect(msg.IsNew).To(BeFalse())
		})
	})

	Describe("FactDeletedMsg", func() {
		It("should store deleted fact ID", func() {
			msg := factmanagement.FactDeletedMsg{
				FactID: "fact-123",
			}
			Expect(msg.FactID).To(Equal("fact-123"))
		})
	})

	Describe("FactsLoadedMsg", func() {
		It("should store facts and total count", func() {
			facts := []*career.Fact{
				fixtures.FactWith("fact-1", "Fact 1"),
				fixtures.FactWith("fact-2", "Fact 2"),
			}
			msg := factmanagement.FactsLoadedMsg{
				Facts: facts,
				Total: 100,
			}
			Expect(msg.Facts).To(HaveLen(2))
			Expect(msg.Total).To(Equal(100))
		})

		It("should allow empty facts", func() {
			msg := factmanagement.FactsLoadedMsg{
				Facts: []*career.Fact{},
				Total: 0,
			}
			Expect(msg.Facts).To(BeEmpty())
			Expect(msg.Total).To(Equal(0))
		})
	})

	Describe("ErrorMsg", func() {
		It("should store error details", func() {
			err := errors.New("database connection failed")
			msg := factmanagement.ErrorMsg{
				Code:    "DB_ERROR",
				Message: "Failed to connect to database",
				Err:     err,
			}
			Expect(msg.Code).To(Equal("DB_ERROR"))
			Expect(msg.Message).To(Equal("Failed to connect to database"))
			Expect(msg.Err).To(Equal(err))
		})

		It("should allow nil error", func() {
			msg := factmanagement.ErrorMsg{
				Code:    "VALIDATION_ERROR",
				Message: "Invalid input",
				Err:     nil,
			}
			Expect(msg.Err).To(BeNil())
		})
	})
})

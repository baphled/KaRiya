package burst_management_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents/burst_management"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	"github.com/baphled/kariya/internal/testutil/mocks"
)

var _ = Describe("Interfaces", func() {
	Describe("BurstService", func() {
		Describe("Interface Compliance", func() {
			It("should be implemented by BurstServiceMock", func() {
				// Verify the mock implements the interface at compile time.
				var _ burst_management.BurstService = (*mocks.BurstServiceMock)(nil)

				// Also verify at runtime that we can use it.
				mock := mocks.NewBurstServiceMock()
				Expect(mock).NotTo(BeNil())

				// Cast to interface to ensure it works.
				var service burst_management.BurstService = mock
				Expect(service).NotTo(BeNil())
			})

			It("should define GetEventByID method", func() {
				mock := mocks.NewBurstServiceMock()
				var service burst_management.BurstService = mock

				// Method should exist and be callable.
				_, err := service.GetEventByID(context.TODO(), "test-id")
				// Error expected since event doesn't exist, but method is callable.
				Expect(err).To(HaveOccurred())
			})

			It("should define ListEvents method", func() {
				mock := mocks.NewBurstServiceMock()
				var service burst_management.BurstService = mock

				events, err := service.ListEvents(context.TODO(), careerrepo.ListFilters{})
				Expect(err).NotTo(HaveOccurred())
				Expect(events).NotTo(BeNil())
			})

			It("should define GetFactsBySourceBurstID method", func() {
				mock := mocks.NewBurstServiceMock()
				var service burst_management.BurstService = mock

				facts, err := service.GetFactsBySourceBurstID(context.TODO(), "test-burst-id")
				Expect(err).NotTo(HaveOccurred())
				Expect(facts).NotTo(BeNil())
			})

			It("should define ExtractFactsFromBurst method", func() {
				mock := mocks.NewBurstServiceMock()
				var service burst_management.BurstService = mock

				// Method should exist and be callable.
				// The mock returns nil slice by default, which is valid.
				_, err := service.ExtractFactsFromBurst(context.TODO(), nil)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should define SaveFact method", func() {
				mock := mocks.NewBurstServiceMock()
				var service burst_management.BurstService = mock

				err := service.SaveFact(context.TODO(), nil)
				// No error expected for nil fact in mock.
				Expect(err).NotTo(HaveOccurred())
			})

			It("should define SuggestBursts method", func() {
				mock := mocks.NewBurstServiceMock()
				var service burst_management.BurstService = mock

				// Method should exist and be callable.
				// The mock returns nil slice by default, which is valid.
				_, err := service.SuggestBursts(context.TODO(), []string{"event-1"})
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Describe("Interface Method Signatures", func() {
			It("should have 6 methods for cohesive burst operations", func() {
				// This test documents that the interface has 6 methods:
				// Event operations (2): GetEventByID, ListEvents
				// Fact operations (3): GetFactsBySourceBurstID, ExtractFactsFromBurst, SaveFact
				// Suggestion operations (1): SuggestBursts
				//
				// The nolint:interfacebloat directive is justified because these
				// methods are all necessary for burst management workflows.
				mock := mocks.NewBurstServiceMock()
				var service burst_management.BurstService = mock
				Expect(service).NotTo(BeNil())
			})
		})
	})
})

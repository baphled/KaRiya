package career_test

import (
	"context"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Project and Company Fields", func() {
	var (
		repo *careerrepo.MemoryRepository
		ctx  context.Context
	)

	BeforeEach(func() {
		repo = careerrepo.NewMemoryRepository()
		ctx = context.Background()
	})

	Context("when saving events with company and project", func() {
		It("should persist both fields correctly", func() {
			event := &career.CareerEvent{
				Text:    "Developed new microservice architecture",
				Date:    time.Now().Add(-24 * time.Hour),
				Company: "TechCorp Inc.",
				Project: "Platform Modernization",
				Tags:    []string{"technical"},
			}

			// Create the event
			err := repo.Create(ctx, event)
			Expect(err).ToNot(HaveOccurred())

			// Retrieve the event
			retrieved, err := repo.GetByID(ctx, event.ID)
			Expect(err).ToNot(HaveOccurred())

			// Verify company and project
			Expect(retrieved.Company).To(Equal("TechCorp Inc."))
			Expect(retrieved.Project).To(Equal("Platform Modernization"))
		})

		It("should handle empty company and project", func() {
			event := &career.CareerEvent{
				Text:    "Simple event without company or project",
				Date:    time.Now().Add(-24 * time.Hour),
				Company: "",
				Project: "",
				Tags:    []string{"technical"},
			}

			// Create the event
			err := repo.Create(ctx, event)
			Expect(err).ToNot(HaveOccurred())

			// Retrieve the event
			retrieved, err := repo.GetByID(ctx, event.ID)
			Expect(err).ToNot(HaveOccurred())

			// Verify empty strings are preserved
			Expect(retrieved.Company).To(Equal(""))
			Expect(retrieved.Project).To(Equal(""))
		})

		It("should update company and project correctly", func() {
			event := &career.CareerEvent{
				Text:    "Initial event",
				Date:    time.Now().Add(-24 * time.Hour),
				Company: "OldCorp",
				Project: "Old Project",
				Tags:    []string{"technical"},
			}

			// Create the event
			err := repo.Create(ctx, event)
			Expect(err).ToNot(HaveOccurred())

			// Update with new company and project
			event.Company = "NewCorp Ltd."
			event.Project = "Cloud Migration"
			err = repo.Update(ctx, event)
			Expect(err).ToNot(HaveOccurred())

			// Retrieve and verify
			updated, err := repo.GetByID(ctx, event.ID)
			Expect(err).ToNot(HaveOccurred())

			Expect(updated.Company).To(Equal("NewCorp Ltd."))
			Expect(updated.Project).To(Equal("Cloud Migration"))
		})

		It("should list events with company and project", func() {
			events := []*career.CareerEvent{
				{
					Text:    "Event 1",
					Date:    time.Now().Add(-48 * time.Hour),
					Company: "Company A",
					Project: "Project Alpha",
					Tags:    []string{"technical"},
				},
				{
					Text:    "Event 2",
					Date:    time.Now().Add(-24 * time.Hour),
					Company: "Company B",
					Project: "Project Beta",
					Tags:    []string{"leadership"},
				},
			}

			// Create events
			for _, e := range events {
				err := repo.Create(ctx, e)
				Expect(err).ToNot(HaveOccurred())
			}

			// List all events
			filters := careerrepo.ListFilters{
				Limit:     10,
				SortBy:    "date",
				SortOrder: "desc",
			}
			retrieved, err := repo.List(ctx, filters)
			Expect(err).ToNot(HaveOccurred())
			Expect(retrieved).To(HaveLen(2))

			// Verify first event (most recent)
			Expect(retrieved[0].Company).To(Equal("Company B"))
			Expect(retrieved[0].Project).To(Equal("Project Beta"))

			// Verify second event
			Expect(retrieved[1].Company).To(Equal("Company A"))
			Expect(retrieved[1].Project).To(Equal("Project Alpha"))
		})
	})
})

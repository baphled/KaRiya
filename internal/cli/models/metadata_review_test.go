package models_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/models"
	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
)

var _ = Describe("MetadataReviewModel", func() {
	var (
		model   *models.MetadataReviewModel
		service *careerservice.Service
		repo    *careerrepo.MemoryRepository
		ctx     context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		repo = careerrepo.NewMemoryRepository()
		service = careerservice.NewService(repo)

		// Create some test events
		event1 := &career.CareerEvent{
			ID:         "event-1",
			Text:       "First event",
			Date:       time.Now().Add(-48 * time.Hour),
			Company:    "CompanyA",
			Project:    "ProjectA",
			Tags:       []string{"technical"},
			Categories: []string{"Technical"},
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		event2 := &career.CareerEvent{
			ID:         "event-2",
			Text:       "Second event",
			Date:       time.Now().Add(-24 * time.Hour),
			Company:    "CompanyB",
			Project:    "ProjectB",
			Tags:       []string{"leadership"},
			Categories: []string{"Leadership"},
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		repo.Create(ctx, event1)
		repo.Create(ctx, event2)

		model = models.NewMetadataReviewModel(service, ctx)
	})

	Describe("Creation", func() {
		It("should create a metadata review model", func() {
			Expect(model).NotTo(BeNil())
		})
	})

	Describe("View", func() {
		It("should render without panicking", func() {
			Expect(func() {
				_ = model.View()
			}).NotTo(Panic())
		})

		It("should show header in view", func() {
			view := model.View()
			Expect(view).To(ContainSubstring("Metadata Review"))
		})
	})

	Describe("Init", func() {
		It("should initialize without error", func() {
			cmd := model.Init()
			Expect(cmd).To(BeNil())
		})
	})

	Describe("Update", func() {
		It("should handle key messages", func() {
			_, cmd := model.Update(nil)
			Expect(cmd).To(BeNil())
		})
	})
})

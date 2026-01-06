package cv

import (
	"context"
	"io"
	"time"

	career "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CVGenerationService - Sections Should Be Populated", func() {
	var (
		log *logger.Logger
		ctx context.Context
	)

	BeforeEach(func() {
		log = logger.New(io.Discard, logger.InfoLevel)
		ctx = context.Background()
	})

	It("should populate sections when generating CV with real bullet generator and section builder", func() {
		// Create test events
		events := []*career.CareerEvent{
			{
				ID:         uuid.New().String(),
				Text:       "Architected microservices platform",
				Date:       time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				Company:    "TechCorp",
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
				Tags:       []string{"architecture"},
				Categories: []string{"technical"},
			},
		}

		facts := []*career.Fact{
			{
				ID:                   uuid.New().String(),
				Text:                 "Expert in distributed systems",
				CompetencyCategories: []string{"technical"},
				RoleFit:              "staff",
				AudienceRelevance:    []string{"hiring_manager"},
				SourceEventID:        events[0].ID,
				CreatedAt:            time.Now(),
				UpdatedAt:            time.Now(),
			},
		}

		config := &career.CVConfig{
			Name:           "test-cv",
			TargetRole:     "staff",
			TargetAudience: "hiring_manager",
			EventFilters:   make(map[string]interface{}),
		}

		// Create repositories that return test data
		eventRepo := &TestEventRepository{events: events}
		factRepo := &TestFactRepository{facts: facts}

		// Create services with REAL implementations (not empty mocks)
		configManager := NewMemoryConfigManager()
		bulletGenerator := NewBulletGenerator(eventRepo, factRepo, log)
		sectionBuilder := NewSectionBuilder(log) // Use REAL section builder, not empty

		service := NewCVGenerationService(
			eventRepo,
			factRepo,
			configManager,
			bulletGenerator,
			sectionBuilder,
			log,
		)

		// Generate CV
		cv, err := service.GenerateCVFromConfig(ctx, config)
		Expect(err).NotTo(HaveOccurred())
		Expect(cv).NotTo(BeNil())

		// THIS IS THE KEY TEST - sections should be populated
		Expect(cv.Sections).NotTo(BeEmpty(), "CV SECTIONS MUST NOT BE EMPTY when generating from events")
		Expect(len(cv.Sections)).To(BeNumerically(">", 0), "CV should have at least one section")

		// Verify sections have content
		for _, section := range cv.Sections {
			Expect(section.ID).NotTo(BeEmpty())
			Expect(section.Title).NotTo(BeEmpty())
			// Summary sections use Summary field, others use Content array
			if section.SectionType == "summary" {
				Expect(section.Summary).NotTo(BeEmpty(), "Summary section should have Summary text")
			} else {
				Expect(len(section.Content)).To(BeNumerically(">", 0), "Section should have content groups")
			}
		}
	})
})

// TestEventRepository provides test events
type TestEventRepository struct {
	events []*career.CareerEvent
}

func (r *TestEventRepository) List(ctx context.Context, filters careerrepo.ListFilters) ([]*career.CareerEvent, error) {
	return r.events, nil
}

func (r *TestEventRepository) GetByID(ctx context.Context, id string) (*career.CareerEvent, error) {
	return nil, nil
}

func (r *TestEventRepository) Create(ctx context.Context, event *career.CareerEvent) error {
	return nil
}

func (r *TestEventRepository) Update(ctx context.Context, event *career.CareerEvent) error {
	return nil
}

func (r *TestEventRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (r *TestEventRepository) Count(ctx context.Context, filters careerrepo.ListFilters) (int, error) {
	return len(r.events), nil
}

// TestFactRepository provides test facts
type TestFactRepository struct {
	facts []*career.Fact
}

func (r *TestFactRepository) List(ctx context.Context, filters careerrepo.FactListFilters) ([]*career.Fact, error) {
	return r.facts, nil
}

func (r *TestFactRepository) GetByID(ctx context.Context, id string) (*career.Fact, error) {
	return nil, nil
}

func (r *TestFactRepository) Create(ctx context.Context, fact *career.Fact) error {
	return nil
}

func (r *TestFactRepository) Update(ctx context.Context, fact *career.Fact) error {
	return nil
}

func (r *TestFactRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (r *TestFactRepository) Count(ctx context.Context, filters careerrepo.FactListFilters) (int, error) {
	return len(r.facts), nil
}

func (r *TestFactRepository) GetBySourceEventID(ctx context.Context, eventID string) ([]*career.Fact, error) {
	return []*career.Fact{}, nil
}

func (r *TestFactRepository) GetBySourceBurstID(ctx context.Context, burstID string) ([]*career.Fact, error) {
	return []*career.Fact{}, nil
}

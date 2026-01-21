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

var _ = Describe("DefaultCVGenerationService", func() {
	var (
		log *logger.Logger
		ctx context.Context
	)

	BeforeEach(func() {
		log = logger.New(io.Discard, logger.InfoLevel)
		ctx = context.Background()
	})

	Describe("GenerateCV", func() {
		It("should generate CV from saved configuration by name", func() {
			config := &career.CVConfig{
				Name:           "test-cv",
				TargetRole:     "principal",
				TargetAudience: "hiring_manager",
				EventFilters:   make(map[string]interface{}),
			}

			configManager := NewMockConfigManager()
			configManager.configs["test-cv"] = config

			service := NewCVGenerationService(
				NewEmptyRepository(),
				NewEmptyFactRepository(),
				configManager,
				NewEmptyBulletGenerator(),
				NewEmptySectionBuilder(),
				log,
			)

			cv, err := service.GenerateCV(ctx, "test-cv")
			Expect(err).NotTo(HaveOccurred())
			Expect(cv).NotTo(BeNil())
			Expect(cv.Name).To(Equal("test-cv"))
			Expect(cv.TargetRole).To(Equal("principal"))
		})

		It("should return error when configuration not found", func() {
			configManager := NewMockConfigManager()

			service := NewCVGenerationService(
				NewEmptyRepository(),
				NewEmptyFactRepository(),
				configManager,
				NewEmptyBulletGenerator(),
				NewEmptySectionBuilder(),
				log,
			)

			_, err := service.GenerateCV(ctx, "nonexistent-cv")
			Expect(err).To(HaveOccurred())
		})

		It("should handle context cancellation", func() {
			config := &career.CVConfig{
				Name:           "test-cv",
				TargetRole:     "principal",
				TargetAudience: "hiring_manager",
				EventFilters:   make(map[string]interface{}),
			}

			configManager := NewMockConfigManager()
			configManager.configs["test-cv"] = config

			cancelCtx, cancel := context.WithCancel(context.Background())
			cancel()

			service := NewCVGenerationService(
				NewEmptyRepository(),
				NewEmptyFactRepository(),
				configManager,
				NewEmptyBulletGenerator(),
				NewEmptySectionBuilder(),
				log,
			)

			_, err := service.GenerateCV(cancelCtx, "test-cv")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("GenerateCVFromConfig", func() {
		It("should return error when configuration is nil", func() {
			service := NewCVGenerationService(
				nil,
				nil,
				nil,
				nil,
				nil,
				log,
			)
			_, err := service.GenerateCVFromConfig(ctx, nil)
			Expect(err).To(HaveOccurred())
		})

		It("should require valid target role", func() {
			config := &career.CVConfig{
				Name: "test-cv",
				// Missing TargetRole
			}

			service := NewCVGenerationService(
				nil,
				nil,
				nil,
				nil,
				nil,
				log,
			)
			_, err := service.GenerateCVFromConfig(ctx, config)
			Expect(err).To(HaveOccurred())
		})

		It("should require at least one target audience", func() {
			config := &career.CVConfig{
				Name:       "test-cv",
				TargetRole: "principal",
				// Missing TargetAudience
			}

			service := NewCVGenerationService(
				nil,
				nil,
				nil,
				nil,
				nil,
				log,
			)
			_, err := service.GenerateCVFromConfig(ctx, config)
			Expect(err).To(HaveOccurred())
		})

		It("should generate CV with valid configuration", func() {
			config := &career.CVConfig{
				Name:           "test-cv",
				TargetRole:     "principal",
				TargetAudience: "hiring_manager",
				EventFilters:   make(map[string]interface{}),
			}

			service := NewCVGenerationService(
				NewEmptyRepository(),
				NewEmptyFactRepository(),
				NewMockConfigManager(),
				NewEmptyBulletGenerator(),
				NewEmptySectionBuilder(),
				log,
			)

			cv, err := service.GenerateCVFromConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())
			Expect(cv).NotTo(BeNil())
			Expect(cv.Name).To(Equal("test-cv"))
			Expect(cv.TargetRole).To(Equal("principal"))
		})

		It("should set generated timestamp", func() {
			config := &career.CVConfig{
				Name:           "test-cv",
				TargetRole:     "principal",
				TargetAudience: "hiring_manager",
				EventFilters:   make(map[string]interface{}),
			}

			service := NewCVGenerationService(
				NewEmptyRepository(),
				NewEmptyFactRepository(),
				NewMockConfigManager(),
				NewEmptyBulletGenerator(),
				NewEmptySectionBuilder(),
				log,
			)

			beforeGeneration := time.Now()
			cv, err := service.GenerateCVFromConfig(ctx, config)
			afterGeneration := time.Now()

			Expect(err).NotTo(HaveOccurred())
			Expect(cv.GeneratedAt).To(BeTemporally(">=", beforeGeneration))
			Expect(cv.GeneratedAt).To(BeTemporally("<=", afterGeneration))
		})

		It("should handle context cancellation", func() {
			config := &career.CVConfig{
				Name:           "test-cv",
				TargetRole:     "principal",
				TargetAudience: "hiring_manager",
				EventFilters:   make(map[string]interface{}),
			}

			cancelCtx, cancel := context.WithCancel(context.Background())
			cancel()

			service := NewCVGenerationService(
				NewEmptyRepository(),
				NewEmptyFactRepository(),
				NewMockConfigManager(),
				NewEmptyBulletGenerator(),
				NewEmptySectionBuilder(),
				log,
			)

			_, err := service.GenerateCVFromConfig(cancelCtx, config)
			Expect(err).To(HaveOccurred())
		})

		It("should generate unique CV IDs", func() {
			config := &career.CVConfig{
				Name:           "test-cv",
				TargetRole:     "principal",
				TargetAudience: "hiring_manager",
				EventFilters:   make(map[string]interface{}),
			}

			service := NewCVGenerationService(
				NewEmptyRepository(),
				NewEmptyFactRepository(),
				NewMockConfigManager(),
				NewEmptyBulletGenerator(),
				NewEmptySectionBuilder(),
				log,
			)

			cv1, err := service.GenerateCVFromConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())

			cv2, err := service.GenerateCVFromConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())

			Expect(cv1.ID).NotTo(Equal(cv2.ID))
		})

		It("should support single target audience", func() {
			config := &career.CVConfig{
				Name:           "test-cv",
				TargetRole:     "principal",
				TargetAudience: "hiring_manager",
				EventFilters:   make(map[string]interface{}),
			}

			service := NewCVGenerationService(
				NewEmptyRepository(),
				NewEmptyFactRepository(),
				NewMockConfigManager(),
				NewEmptyBulletGenerator(),
				NewEmptySectionBuilder(),
				log,
			)

			cv, err := service.GenerateCVFromConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())
			Expect(cv.TargetAudience).To(Equal("hiring_manager"))
		})

		It("should track source event and fact counts", func() {
			config := &career.CVConfig{
				Name:           "test-cv",
				TargetRole:     "principal",
				TargetAudience: "hiring_manager",
				EventFilters:   make(map[string]interface{}),
			}

			service := NewCVGenerationService(
				NewCountingRepository(5),
				NewCountingFactRepository(3),
				NewMockConfigManager(),
				NewEmptyBulletGenerator(),
				NewEmptySectionBuilder(),
				log,
			)

			cv, err := service.GenerateCVFromConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())
			Expect(cv.SourceEventCount).To(Equal(5))
			Expect(cv.SourceFactCount).To(Equal(3))
		})

		It("should preserve event filters in generated CV", func() {
			filters := map[string]interface{}{
				"companies": []string{"Google", "Meta"},
				"tags":      []string{"leadership", "technical"},
			}

			config := &career.CVConfig{
				Name:           "test-cv",
				TargetRole:     "principal",
				TargetAudience: "hiring_manager",
				EventFilters:   filters,
			}

			service := NewCVGenerationService(
				NewEmptyRepository(),
				NewEmptyFactRepository(),
				NewMockConfigManager(),
				NewEmptyBulletGenerator(),
				NewEmptySectionBuilder(),
				log,
			)

			cv, err := service.GenerateCVFromConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())
			Expect(cv.EventFilters).To(Equal(filters))
		})

		It("should support all valid target roles", func() {
			validRoles := []string{"principal", "staff", "em", "senior_ic"}

			for _, role := range validRoles {
				config := &career.CVConfig{
					Name:           "test-cv",
					TargetRole:     role,
					TargetAudience: "hiring_manager",
					EventFilters:   make(map[string]interface{}),
				}

				service := NewCVGenerationService(
					NewEmptyRepository(),
					NewEmptyFactRepository(),
					NewMockConfigManager(),
					NewEmptyBulletGenerator(),
					NewEmptySectionBuilder(),
					log,
				)

				cv, err := service.GenerateCVFromConfig(ctx, config)
				Expect(err).NotTo(HaveOccurred())
				Expect(cv.TargetRole).To(Equal(role))
			}
		})

		It("should return ephemeral CVs (not persisted)", func() {
			config := &career.CVConfig{
				Name:           "test-cv",
				TargetRole:     "principal",
				TargetAudience: "hiring_manager",
				EventFilters:   make(map[string]interface{}),
			}

			service := NewCVGenerationService(
				NewEmptyRepository(),
				NewEmptyFactRepository(),
				NewMockConfigManager(),
				NewEmptyBulletGenerator(),
				NewEmptySectionBuilder(),
				log,
			)

			cv, err := service.GenerateCVFromConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())

			// CV should exist in memory but not be retrievable from config manager
			// (since it's ephemeral and not persisted)
			Expect(cv).NotTo(BeNil())
			Expect(cv.ID).NotTo(BeEmpty())
		})
	})

	Describe("Edge Cases", func() {
		It("should handle empty event repository gracefully", func() {
			config := &career.CVConfig{
				Name:           "test-cv",
				TargetRole:     "principal",
				TargetAudience: "hiring_manager",
				EventFilters:   make(map[string]interface{}),
			}

			service := NewCVGenerationService(
				NewEmptyRepository(),
				NewEmptyFactRepository(),
				NewMockConfigManager(),
				NewEmptyBulletGenerator(),
				NewEmptySectionBuilder(),
				log,
			)

			cv, err := service.GenerateCVFromConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())
			Expect(cv).NotTo(BeNil())
		})

		It("should handle empty fact repository gracefully", func() {
			config := &career.CVConfig{
				Name:           "test-cv",
				TargetRole:     "principal",
				TargetAudience: "hiring_manager",
				EventFilters:   make(map[string]interface{}),
			}

			service := NewCVGenerationService(
				NewEmptyRepository(),
				NewEmptyFactRepository(),
				NewMockConfigManager(),
				NewEmptyBulletGenerator(),
				NewEmptySectionBuilder(),
				log,
			)

			cv, err := service.GenerateCVFromConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())
			Expect(cv.SourceFactCount).To(Equal(0))
		})

		It("should handle nil event filters", func() {
			config := &career.CVConfig{
				Name:           "test-cv",
				TargetRole:     "principal",
				TargetAudience: "hiring_manager",
				EventFilters:   nil,
			}

			service := NewCVGenerationService(
				NewEmptyRepository(),
				NewEmptyFactRepository(),
				NewMockConfigManager(),
				NewEmptyBulletGenerator(),
				NewEmptySectionBuilder(),
				log,
			)

			cv, err := service.GenerateCVFromConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())
			Expect(cv).NotTo(BeNil())
		})
	})
})

// Mock implementations for testing

type MockConfigManager struct {
	configs map[string]*career.CVConfig
}

func NewMockConfigManager() *MockConfigManager {
	return &MockConfigManager{
		configs: make(map[string]*career.CVConfig),
	}
}

func (m *MockConfigManager) LoadConfig(ctx context.Context, name string) (*career.CVConfig, error) {
	if config, ok := m.configs[name]; ok {
		return config, nil
	}
	return nil, ErrConfigNotFound
}

func (m *MockConfigManager) SaveConfig(ctx context.Context, config *career.CVConfig) error {
	m.configs[config.Name] = config
	return nil
}

func (m *MockConfigManager) DeleteConfig(ctx context.Context, name string) error {
	delete(m.configs, name)
	return nil
}

func (m *MockConfigManager) ListConfigs(ctx context.Context) ([]*career.CVConfig, error) {
	configs := make([]*career.CVConfig, 0, len(m.configs))
	for _, config := range m.configs {
		configs = append(configs, config)
	}
	return configs, nil
}

func (m *MockConfigManager) GetConfigPath(name string) string {
	return "/tmp/" + name
}

func (m *MockConfigManager) ConfigExists(ctx context.Context, name string) (bool, error) {
	_, ok := m.configs[name]
	return ok, nil
}

type EmptyRepository struct{}

func NewEmptyRepository() *EmptyRepository {
	return &EmptyRepository{}
}

func (r *EmptyRepository) List(ctx context.Context, filters careerrepo.ListFilters) ([]*career.CareerEvent, error) {
	return []*career.CareerEvent{}, nil
}

func (r *EmptyRepository) GetByID(ctx context.Context, id string) (*career.CareerEvent, error) {
	return nil, nil //nolint:nilnil // test stub - method not used in these tests
}

func (r *EmptyRepository) Create(ctx context.Context, event *career.CareerEvent) error {
	return nil
}

func (r *EmptyRepository) Update(ctx context.Context, event *career.CareerEvent) error {
	return nil
}

func (r *EmptyRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (r *EmptyRepository) Count(ctx context.Context, filters careerrepo.ListFilters) (int, error) {
	return 0, nil
}

type CountingRepository struct {
	count int
}

func NewCountingRepository(count int) *CountingRepository {
	return &CountingRepository{count: count}
}

func (r *CountingRepository) List(ctx context.Context, filters careerrepo.ListFilters) ([]*career.CareerEvent, error) {
	events := make([]*career.CareerEvent, r.count)
	for i := 0; i < r.count; i++ {
		events[i] = &career.CareerEvent{
			ID:   uuid.New().String(),
			Text: "Sample event",
			Date: time.Now(),
		}
	}
	return events, nil
}

func (r *CountingRepository) GetByID(ctx context.Context, id string) (*career.CareerEvent, error) {
	return nil, nil //nolint:nilnil // test stub - method not used in these tests
}

func (r *CountingRepository) Create(ctx context.Context, event *career.CareerEvent) error {
	return nil
}

func (r *CountingRepository) Update(ctx context.Context, event *career.CareerEvent) error {
	return nil
}

func (r *CountingRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (r *CountingRepository) Count(ctx context.Context, filters careerrepo.ListFilters) (int, error) {
	return r.count, nil
}

type EmptyFactRepository struct{}

func NewEmptyFactRepository() *EmptyFactRepository {
	return &EmptyFactRepository{}
}

func (r *EmptyFactRepository) List(ctx context.Context, filters careerrepo.FactListFilters) ([]*career.Fact, error) {
	return []*career.Fact{}, nil
}

func (r *EmptyFactRepository) GetByID(ctx context.Context, id string) (*career.Fact, error) {
	return nil, nil //nolint:nilnil // test stub - method not used in these tests
}

func (r *EmptyFactRepository) Create(ctx context.Context, fact *career.Fact) error {
	return nil
}

func (r *EmptyFactRepository) Update(ctx context.Context, fact *career.Fact) error {
	return nil
}

func (r *EmptyFactRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (r *EmptyFactRepository) Count(ctx context.Context, filters careerrepo.FactListFilters) (int, error) {
	return 0, nil
}

func (r *EmptyFactRepository) GetBySourceEventID(ctx context.Context, eventID string) ([]*career.Fact, error) {
	return []*career.Fact{}, nil
}

func (r *EmptyFactRepository) GetBySourceBurstID(ctx context.Context, burstID string) ([]*career.Fact, error) {
	return []*career.Fact{}, nil
}

type CountingFactRepository struct {
	count int
}

func NewCountingFactRepository(count int) *CountingFactRepository {
	return &CountingFactRepository{count: count}
}

func (r *CountingFactRepository) List(ctx context.Context, filters careerrepo.FactListFilters) ([]*career.Fact, error) {
	facts := make([]*career.Fact, r.count)
	for i := 0; i < r.count; i++ {
		facts[i] = &career.Fact{
			ID:   uuid.New().String(),
			Text: "Sample fact",
		}
	}
	return facts, nil
}

func (r *CountingFactRepository) GetByID(ctx context.Context, id string) (*career.Fact, error) {
	return nil, nil //nolint:nilnil // test stub - method not used in these tests
}

func (r *CountingFactRepository) Create(ctx context.Context, fact *career.Fact) error {
	return nil
}

func (r *CountingFactRepository) Update(ctx context.Context, fact *career.Fact) error {
	return nil
}

func (r *CountingFactRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (r *CountingFactRepository) Count(ctx context.Context, filters careerrepo.FactListFilters) (int, error) {
	return r.count, nil
}

func (r *CountingFactRepository) GetBySourceEventID(ctx context.Context, eventID string) ([]*career.Fact, error) {
	return []*career.Fact{}, nil
}

func (r *CountingFactRepository) GetBySourceBurstID(ctx context.Context, burstID string) ([]*career.Fact, error) {
	return []*career.Fact{}, nil
}

type EmptyBulletGenerator struct{}

func NewEmptyBulletGenerator() *EmptyBulletGenerator {
	return &EmptyBulletGenerator{}
}

func (g *EmptyBulletGenerator) GenerateBullets(ctx context.Context, events []*career.CareerEvent, facts []*career.Fact, targetRole string, targetAudience string) ([]*career.CVBullet, error) {
	return []*career.CVBullet{}, nil
}

func (g *EmptyBulletGenerator) FilterByTechnologies(bullets []*career.CVBullet, events []*career.CareerEvent, techFocus TechnologyFocus, technologies []string) []*career.CVBullet {
	return bullets // No-op for tests
}

type EmptySectionBuilder struct{}

func NewEmptySectionBuilder() *EmptySectionBuilder {
	return &EmptySectionBuilder{}
}

func (b *EmptySectionBuilder) BuildSections(ctx context.Context, bullets []*career.CVBullet, events []*career.CareerEvent, facts []*career.Fact, targetRole string, skillsConfig *SkillsFormatConfig) ([]*career.CVSection, error) {
	return []*career.CVSection{}, nil
}

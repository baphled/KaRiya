package cv

import (
	"context"

	career "github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	"github.com/baphled/kariya/internal/logger"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DefaultCVGenerationService", func() {
	var (
		log *logger.Logger
		ctx context.Context
	)

	BeforeEach(func() {
		log = logger.DefaultLogger()
		ctx = context.Background()
	})

	It("should handle nil configuration", func() {
		service := NewCVGenerationService(nil, nil, nil, nil, nil, log)
		_, err := service.GenerateCVFromConfig(ctx, nil)
		Expect(err).To(HaveOccurred())
	})

	It("should require valid target role", func() {
		config := &career.CVConfig{
			Name: "test-cv",
		}

		service := NewCVGenerationService(nil, nil, nil, nil, nil, log)
		_, err := service.GenerateCVFromConfig(ctx, config)
		Expect(err).To(HaveOccurred())
	})

	It("should handle context cancellation", func() {
		config := &career.CVConfig{
			Name:           "test-cv",
			TargetRole:     "principal",
			TargetAudience: []string{"hiring_manager"},
		}

		cancelCtx, cancel := context.WithCancel(context.Background())
		cancel()

		service := NewCVGenerationService(nil, nil, nil, nil, nil, log)
		_, err := service.GenerateCVFromConfig(cancelCtx, config)
		Expect(err).To(HaveOccurred())
	})

	It("should load configuration by name", func() {
		config := &career.CVConfig{
			Name:           "test-cv",
			TargetRole:     "principal",
			TargetAudience: []string{"hiring_manager"},
			EventFilters:   make(map[string]interface{}),
		}

		configManager := NewMockConfigManager()
		configManager.configs["test-cv"] = config

		service := NewCVGenerationService(
			NewEmptyRepository(),
			nil,
			configManager,
			NewEmptyBulletGenerator(),
			NewEmptySectionBuilder(),
			log,
		)

		_, err := service.GenerateCV(ctx, "test-cv")
		Expect(err).NotTo(HaveOccurred())
	})
})

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
	return nil, nil
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
	return nil, nil
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

type EmptyBulletGenerator struct{}

func NewEmptyBulletGenerator() *EmptyBulletGenerator {
	return &EmptyBulletGenerator{}
}

func (g *EmptyBulletGenerator) GenerateBullets(ctx context.Context, events []*career.CareerEvent, facts []*career.Fact, targetRole string, targetAudiences []string) ([]*career.CVBullet, error) {
	return []*career.CVBullet{}, nil
}

type EmptySectionBuilder struct{}

func NewEmptySectionBuilder() *EmptySectionBuilder {
	return &EmptySectionBuilder{}
}

func (b *EmptySectionBuilder) BuildSections(ctx context.Context, bullets []*career.CVBullet, events []*career.CareerEvent, targetRole string) ([]*career.CVSection, error) {
	return []*career.CVSection{}, nil
}

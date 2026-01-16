package cv

import (
	"context"
	"database/sql"
	"io"
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
)

var _ = Describe("CV Generation Integration Tests", func() {
	var (
		log          *logger.Logger
		ctx          context.Context
		tempDir      string
		dbPath       string
		eventRepo    *careerrepo.SQLiteRepository
		factRepo     *careerrepo.SQLiteFactRepository
		db           *sql.DB
		cvGenService CVGenerationService
	)

	BeforeEach(func() {
		log = logger.New(io.Discard, logger.InfoLevel)
		ctx = context.Background()

		// Create a temporary directory for the test database
		var err error
		tempDir, err = os.MkdirTemp("", "kariya-cv-integration-test-")
		Expect(err).NotTo(HaveOccurred())

		dbPath = filepath.Join(tempDir, "test_events.db")

		// Create the repositories
		eventRepo, err = careerrepo.NewSQLiteRepository(dbPath)
		Expect(err).NotTo(HaveOccurred())

		db, err = sql.Open("sqlite", dbPath)
		Expect(err).NotTo(HaveOccurred())

		factRepo, err = careerrepo.NewSQLiteFactRepository(db)
		Expect(err).NotTo(HaveOccurred())

		// Create CV services
		bulletGenerator := NewBulletGenerator(eventRepo, factRepo, log)
		sectionBuilder := NewSectionBuilder(log)
		configManager := NewMemoryConfigManager()

		cvGenService = NewCVGenerationService(
			eventRepo,
			factRepo,
			configManager,
			bulletGenerator,
			sectionBuilder,
			log,
		)
	})

	AfterEach(func() {
		if eventRepo != nil {
			eventRepo.Close()
		}
		if db != nil {
			db.Close()
		}
		if tempDir != "" {
			os.RemoveAll(tempDir)
		}
	})

	// Helper function to seed test events
	seedTestEvents := func(count int, category string) []*career.CareerEvent {
		events := make([]*career.CareerEvent, count)
		for i := 0; i < count; i++ {
			event := &career.CareerEvent{
				Text:       "Test event for " + category,
				Date:       time.Now().AddDate(0, 0, -i),
				Tags:       []string{category, "project"},
				Company:    "Test Company",
				Categories: []string{category},
			}
			err := eventRepo.Create(ctx, event)
			Expect(err).NotTo(HaveOccurred())
			events[i] = event
		}
		return events
	}

	Describe("CV Generation from Database", func() {
		It("should generate CV from events in the database", func() {
			// Seed test data
			seedTestEvents(5, "technical")
			seedTestEvents(3, "product")

			// Create a test CV configuration
			config := &career.CVConfig{
				Name:           "test-senior-ic",
				TargetRole:     "senior_ic",
				TargetAudience: "recruiter",
				EventFilters: map[string]interface{}{
					"categories": []string{"technical", "product"},
				},
			}

			// Generate CV
			cv, err := cvGenService.GenerateCVFromConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())
			Expect(cv).NotTo(BeNil())

			// Verify CV structure
			Expect(cv.Name).To(Equal("test-senior-ic"))
			Expect(cv.TargetRole).To(Equal("senior_ic"))
			Expect(cv.TargetAudience).To(Equal("recruiter"))

			// Verify CV has content
			Expect(cv.ID).NotTo(BeEmpty())
			Expect(cv.GeneratedAt).NotTo(BeZero())
		})

		It("should handle CV generation with single target audience", func() {
			// Seed test data
			seedTestEvents(5, "technical")
			seedTestEvents(3, "leadership")

			// Create a test CV configuration with single audience
			config := &career.CVConfig{
				Name:           "test-principal",
				TargetRole:     "principal",
				TargetAudience: "hiring_manager",
				EventFilters: map[string]interface{}{
					"categories": []string{"technical", "achievement", "leadership"},
				},
			}

			// Generate CV
			cv, err := cvGenService.GenerateCVFromConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())
			Expect(cv).NotTo(BeNil())

			// Verify single audience is preserved
			Expect(cv.TargetAudience).To(Equal("hiring_manager"))
		})

		It("should generate CV with all available events when no filters applied", func() {
			// Seed test data with multiple categories
			seedTestEvents(3, "technical")
			seedTestEvents(2, "product")
			seedTestEvents(2, "leadership")

			// Retrieve all events first
			allEvents, err := eventRepo.List(ctx, careerrepo.ListFilters{})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(allEvents)).To(Equal(7))

			// Create a test CV configuration with no filters
			config := &career.CVConfig{
				Name:           "test-no-filters",
				TargetRole:     "senior_ic",
				TargetAudience: "recruiter",
				EventFilters:   make(map[string]interface{}),
			}

			// Generate CV
			cv, err := cvGenService.GenerateCVFromConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())
			Expect(cv).NotTo(BeNil())
		})

		It("should preserve CV metadata during generation", func() {
			// Seed test data
			seedTestEvents(3, "leadership")

			// Create a test CV configuration
			config := &career.CVConfig{
				Name:           "test-metadata",
				TargetRole:     "principal",
				TargetAudience: "hiring_manager",
				EventFilters: map[string]interface{}{
					"categories": []string{"leadership"},
				},
			}

			beforeTime := time.Now()

			// Generate CV
			cv, err := cvGenService.GenerateCVFromConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())
			Expect(cv).NotTo(BeNil())

			afterTime := time.Now()

			// Verify metadata
			Expect(cv.ID).NotTo(BeEmpty())
			Expect(cv.Name).To(Equal("test-metadata"))
			Expect(cv.TargetRole).To(Equal("principal"))
			Expect(cv.GeneratedAt).To(BeTemporally(">=", beforeTime))
			Expect(cv.GeneratedAt).To(BeTemporally("<=", afterTime))
		})

		It("should handle CV generation with category filters", func() {
			// Seed test data
			seedTestEvents(5, "technical")
			seedTestEvents(3, "product")
			seedTestEvents(2, "leadership")

			// Create a test CV configuration with category filters
			config := &career.CVConfig{
				Name:           "test-technical-only",
				TargetRole:     "senior_ic",
				TargetAudience: "recruiter",
				EventFilters: map[string]interface{}{
					"categories": []string{"technical"},
				},
			}

			// Generate CV
			cv, err := cvGenService.GenerateCVFromConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())
			Expect(cv).NotTo(BeNil())
		})

		It("should validate CV configuration before generation", func() {
			// Seed some test data (needed for valid case comparison)
			seedTestEvents(2, "technical")

			// Create an invalid configuration (missing target role)
			invalidConfig := &career.CVConfig{
				Name:           "test-invalid",
				TargetRole:     "", // Invalid: empty target role
				TargetAudience: "recruiter",
				EventFilters:   make(map[string]interface{}),
			}

			// Generation should fail
			_, err := cvGenService.GenerateCVFromConfig(ctx, invalidConfig)
			Expect(err).To(HaveOccurred())
		})

		It("should handle empty database gracefully", func() {
			// Don't seed any data - database is empty

			config := &career.CVConfig{
				Name:           "test-empty-db",
				TargetRole:     "senior_ic",
				TargetAudience: "recruiter",
				EventFilters:   make(map[string]interface{}),
			}

			// Generate CV - should still work but with empty content
			cv, err := cvGenService.GenerateCVFromConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())
			Expect(cv).NotTo(BeNil())
			Expect(cv.SourceEventCount).To(Equal(0))
		})

		It("should handle CV generation with multiple audience types", func() {
			// Seed test data
			seedTestEvents(4, "technical")

			// Test with different audience types
			audiences := []string{"recruiter", "hiring_manager", "peer"}

			for _, audience := range audiences {
				config := &career.CVConfig{
					Name:           "test-" + audience,
					TargetRole:     "senior_ic",
					TargetAudience: audience,
					EventFilters:   make(map[string]interface{}),
				}

				cv, err := cvGenService.GenerateCVFromConfig(ctx, config)
				Expect(err).NotTo(HaveOccurred())
				Expect(cv).NotTo(BeNil())
				Expect(cv.TargetAudience).To(Equal(audience))
			}
		})
	})
})

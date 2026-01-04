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

var _ = Describe("CV Generation Integration Tests with Real Database", func() {
	var (
		log       *logger.Logger
		ctx       context.Context
		homeDir   string
		dbPath    string
	)

	BeforeEach(func() {
		log = logger.New(io.Discard, logger.InfoLevel)
		ctx = context.Background()

		// Get the home directory
		var err error
		homeDir, err = os.UserHomeDir()
		Expect(err).NotTo(HaveOccurred())

		dbPath = filepath.Join(homeDir, ".kariya", "events.db")
	})

	Describe("CV Generation from Real Database", func() {
		It("should generate CV from production database at ~/.kariya/events.db", func() {
			// Check if the database exists
			if _, err := os.Stat(dbPath); os.IsNotExist(err) {
				Skip("Database not found at " + dbPath)
			}

			// Open the real database
			eventRepo, err := careerrepo.NewSQLiteRepository(dbPath)
			Expect(err).NotTo(HaveOccurred())
			defer eventRepo.Close()

			// Open the database directly for fact repository
			db, err := sql.Open("sqlite", dbPath)
			Expect(err).NotTo(HaveOccurred())
			defer db.Close()

			factRepo, err := careerrepo.NewSQLiteFactRepository(db)
			Expect(err).NotTo(HaveOccurred())

			// Create CV services
			bulletGenerator := NewBulletGenerator(eventRepo, factRepo, log)
			sectionBuilder := NewSectionBuilder(log)
			configManager := NewMemoryConfigManager()

			cvGenService := NewCVGenerationService(
				eventRepo,
				factRepo,
				configManager,
				bulletGenerator,
				sectionBuilder,
				log,
			)

			// Create a test CV configuration
			config := &career.CVConfig{
				Name:           "test-senior-ic",
				TargetRole:     "senior_ic",
				TargetAudience: []string{"recruiter", "hiring_manager"},
				EventFilters: map[string]interface{}{
					"categories": []string{"technical", "achievement"},
				},
			}

			// Generate CV
			cv, err := cvGenService.GenerateCVFromConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())
			Expect(cv).NotTo(BeNil())

			// Verify CV structure
			Expect(cv.Name).To(Equal("test-senior-ic"))
			Expect(cv.TargetRole).To(Equal("senior_ic"))
			Expect(cv.TargetAudience).To(Equal([]string{"recruiter", "hiring_manager"}))

			// Verify CV has content
			Expect(cv.ID).NotTo(BeEmpty())
			Expect(cv.GeneratedAt).NotTo(BeZero())

			// Log CV statistics
			GinkgoWriter.Printf("CV Generated Successfully:\n")
			GinkgoWriter.Printf("  Name: %s\n", cv.Name)
			GinkgoWriter.Printf("  Target Role: %s\n", cv.TargetRole)
			GinkgoWriter.Printf("  Target Audience: %v\n", cv.TargetAudience)
			GinkgoWriter.Printf("  Source Event Count: %d\n", cv.SourceEventCount)
			GinkgoWriter.Printf("  Source Fact Count: %d\n", cv.SourceFactCount)
		})

		It("should handle CV generation with multiple target audiences", func() {
			// Check if the database exists
			if _, err := os.Stat(dbPath); os.IsNotExist(err) {
				Skip("Database not found at " + dbPath)
			}

			// Open the real database
			eventRepo, err := careerrepo.NewSQLiteRepository(dbPath)
			Expect(err).NotTo(HaveOccurred())
			defer eventRepo.Close()

			// Open the database directly for fact repository
			db, err := sql.Open("sqlite", dbPath)
			Expect(err).NotTo(HaveOccurred())
			defer db.Close()

			factRepo, err := careerrepo.NewSQLiteFactRepository(db)
			Expect(err).NotTo(HaveOccurred())

			// Create CV services
			bulletGenerator := NewBulletGenerator(eventRepo, factRepo, log)
			sectionBuilder := NewSectionBuilder(log)
			configManager := NewMemoryConfigManager()

			cvGenService := NewCVGenerationService(
				eventRepo,
				factRepo,
				configManager,
				bulletGenerator,
				sectionBuilder,
				log,
			)

			// Create a test CV configuration with multiple audiences
			config := &career.CVConfig{
				Name:           "test-principal",
				TargetRole:     "principal",
				TargetAudience: []string{"hiring_manager", "peer", "recruiter"},
				EventFilters: map[string]interface{}{
					"categories": []string{"technical", "achievement", "leadership"},
				},
			}

			// Generate CV
			cv, err := cvGenService.GenerateCVFromConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())
			Expect(cv).NotTo(BeNil())

			// Verify multiple audiences are preserved
			Expect(cv.TargetAudience).To(HaveLen(3))
			Expect(cv.TargetAudience).To(ContainElement("hiring_manager"))
			Expect(cv.TargetAudience).To(ContainElement("peer"))
			Expect(cv.TargetAudience).To(ContainElement("recruiter"))

			GinkgoWriter.Printf("CV with Multiple Audiences Generated:\n")
			GinkgoWriter.Printf("  Audiences: %v\n", cv.TargetAudience)
			GinkgoWriter.Printf("  Source Event Count: %d\n", cv.SourceEventCount)
			GinkgoWriter.Printf("  Source Fact Count: %d\n", cv.SourceFactCount)
		})

		It("should generate CV with all available events when no filters applied", func() {
			// Check if the database exists
			if _, err := os.Stat(dbPath); os.IsNotExist(err) {
				Skip("Database not found at " + dbPath)
			}

			// Open the real database
			eventRepo, err := careerrepo.NewSQLiteRepository(dbPath)
			Expect(err).NotTo(HaveOccurred())
			defer eventRepo.Close()

			// Open the database directly for fact repository
			db, err := sql.Open("sqlite", dbPath)
			Expect(err).NotTo(HaveOccurred())
			defer db.Close()

			factRepo, err := careerrepo.NewSQLiteFactRepository(db)
			Expect(err).NotTo(HaveOccurred())

			// Create CV services
			bulletGenerator := NewBulletGenerator(eventRepo, factRepo, log)
			sectionBuilder := NewSectionBuilder(log)
			configManager := NewMemoryConfigManager()

			cvGenService := NewCVGenerationService(
				eventRepo,
				factRepo,
				configManager,
				bulletGenerator,
				sectionBuilder,
				log,
			)

			// Retrieve all events first
			allEvents, err := eventRepo.List(ctx, careerrepo.ListFilters{})
			Expect(err).NotTo(HaveOccurred())

			GinkgoWriter.Printf("Total Events in Database: %d\n", len(allEvents))

			// Create a test CV configuration with no filters
			config := &career.CVConfig{
				Name:           "test-no-filters",
				TargetRole:     "senior_ic",
				TargetAudience: []string{"recruiter"},
				EventFilters:   make(map[string]interface{}),
			}

			// Generate CV
			cv, err := cvGenService.GenerateCVFromConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())
			Expect(cv).NotTo(BeNil())

			GinkgoWriter.Printf("CV with No Filters Generated:\n")
			GinkgoWriter.Printf("  Total Events Available: %d\n", len(allEvents))
			GinkgoWriter.Printf("  Events Used in CV: %d\n", cv.SourceEventCount)
			GinkgoWriter.Printf("  Facts Used in CV: %d\n", cv.SourceFactCount)
		})

		It("should preserve CV metadata during generation", func() {
			// Check if the database exists
			if _, err := os.Stat(dbPath); os.IsNotExist(err) {
				Skip("Database not found at " + dbPath)
			}

			// Open the real database
			eventRepo, err := careerrepo.NewSQLiteRepository(dbPath)
			Expect(err).NotTo(HaveOccurred())
			defer eventRepo.Close()

			// Open the database directly for fact repository
			db, err := sql.Open("sqlite", dbPath)
			Expect(err).NotTo(HaveOccurred())
			defer db.Close()

			factRepo, err := careerrepo.NewSQLiteFactRepository(db)
			Expect(err).NotTo(HaveOccurred())

			// Create CV services
			bulletGenerator := NewBulletGenerator(eventRepo, factRepo, log)
			sectionBuilder := NewSectionBuilder(log)
			configManager := NewMemoryConfigManager()

			cvGenService := NewCVGenerationService(
				eventRepo,
				factRepo,
				configManager,
				bulletGenerator,
				sectionBuilder,
				log,
			)

			// Create a test CV configuration
			config := &career.CVConfig{
				Name:           "test-metadata",
				TargetRole:     "principal",
				TargetAudience: []string{"hiring_manager", "peer"},
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

			GinkgoWriter.Printf("CV Metadata Preserved:\n")
			GinkgoWriter.Printf("  ID: %s\n", cv.ID)
			GinkgoWriter.Printf("  Name: %s\n", cv.Name)
			GinkgoWriter.Printf("  Generated At: %v\n", cv.GeneratedAt)
		})

		It("should handle CV generation with category filters", func() {
			// Check if the database exists
			if _, err := os.Stat(dbPath); os.IsNotExist(err) {
				Skip("Database not found at " + dbPath)
			}

			// Open the real database
			eventRepo, err := careerrepo.NewSQLiteRepository(dbPath)
			Expect(err).NotTo(HaveOccurred())
			defer eventRepo.Close()

			// Open the database directly for fact repository
			db, err := sql.Open("sqlite", dbPath)
			Expect(err).NotTo(HaveOccurred())
			defer db.Close()

			factRepo, err := careerrepo.NewSQLiteFactRepository(db)
			Expect(err).NotTo(HaveOccurred())

			// Create CV services
			bulletGenerator := NewBulletGenerator(eventRepo, factRepo, log)
			sectionBuilder := NewSectionBuilder(log)
			configManager := NewMemoryConfigManager()

			cvGenService := NewCVGenerationService(
				eventRepo,
				factRepo,
				configManager,
				bulletGenerator,
				sectionBuilder,
				log,
			)

			// Create a test CV configuration with category filters
			config := &career.CVConfig{
				Name:           "test-technical-only",
				TargetRole:     "senior_ic",
				TargetAudience: []string{"recruiter"},
				EventFilters: map[string]interface{}{
					"categories": []string{"technical"},
				},
			}

			// Generate CV
			cv, err := cvGenService.GenerateCVFromConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())
			Expect(cv).NotTo(BeNil())

			GinkgoWriter.Printf("CV with Technical Category Filter:\n")
			GinkgoWriter.Printf("  Events Used: %d\n", cv.SourceEventCount)
			GinkgoWriter.Printf("  Facts Used: %d\n", cv.SourceFactCount)
			GinkgoWriter.Printf("  Filters Applied: %v\n", cv.EventFilters)
		})

		It("should validate CV configuration before generation", func() {
			// Check if the database exists
			if _, err := os.Stat(dbPath); os.IsNotExist(err) {
				Skip("Database not found at " + dbPath)
			}

			// Open the real database
			eventRepo, err := careerrepo.NewSQLiteRepository(dbPath)
			Expect(err).NotTo(HaveOccurred())
			defer eventRepo.Close()

			// Open the database directly for fact repository
			db, err := sql.Open("sqlite", dbPath)
			Expect(err).NotTo(HaveOccurred())
			defer db.Close()

			factRepo, err := careerrepo.NewSQLiteFactRepository(db)
			Expect(err).NotTo(HaveOccurred())

			// Create CV services
			bulletGenerator := NewBulletGenerator(eventRepo, factRepo, log)
			sectionBuilder := NewSectionBuilder(log)
			configManager := NewMemoryConfigManager()

			cvGenService := NewCVGenerationService(
				eventRepo,
				factRepo,
				configManager,
				bulletGenerator,
				sectionBuilder,
				log,
			)

			// Create an invalid configuration (missing target role)
			invalidConfig := &career.CVConfig{
				Name:           "test-invalid",
				TargetRole:     "", // Invalid: empty target role
				TargetAudience: []string{"recruiter"},
				EventFilters:   make(map[string]interface{}),
			}

			// Generation should fail
			_, err = cvGenService.GenerateCVFromConfig(ctx, invalidConfig)
			Expect(err).To(HaveOccurred())

			GinkgoWriter.Printf("Invalid configuration correctly rejected\n")
		})
	})
})


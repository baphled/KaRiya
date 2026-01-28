package cv

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
)

var _ = Describe("CV Generation Integration Tests", func() {
	var (
		log          *logger.Logger
		ctx          context.Context
		tempDir      string
		repos        *careerrepo.Repositories
		cvGenService CVGenerationService
	)

	BeforeEach(func() {
		log = logger.New(io.Discard, logger.InfoLevel)
		ctx = context.Background()

		// Create a temporary directory for the test database
		var err error
		tempDir, err = os.MkdirTemp("", "kariya-cv-integration-test-")
		Expect(err).NotTo(HaveOccurred())

		dbPath := filepath.Join(tempDir, "test_events.db")

		// Create ORM repositories
		repos, err = careerrepo.NewRepositoriesFromPath(dbPath)
		Expect(err).NotTo(HaveOccurred())

		// Create CV services with role-based scoring.
		bulletGenerator := NewBulletGenerator(log, nil) // nil uses default scoring config
		sectionBuilder := NewSectionBuilder(nil, log)
		configManager := NewMemoryConfigManager()
		dataProcessor := NewDataProcessingService(log)

		cvGenService = NewCVGenerationService(
			repos.Event,
			repos.Fact,
			configManager,
			bulletGenerator,
			dataProcessor,
			sectionBuilder,
			log,
		)
	})

	AfterEach(func() {
		if repos != nil {
			repos.Close()
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
			err := repos.Event.Create(ctx, event)
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
			allEvents, err := repos.Event.List(ctx, careerrepo.ListFilters{})
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

	// BUG-003: Regression test to ensure all companies are included in generated CVs
	// Previously, BulletGenerator applied a total cap (e.g., 40 bullets) before grouping
	// by company, causing later companies to be completely excluded from the CV.
	Describe("BUG-003: Multi-Company CV Generation", func() {
		// Helper to seed events for multiple companies using fixtures
		seedMultiCompanyEvents := func(companies []string, eventsPerCompany int) map[string][]*career.CareerEvent {
			result := make(map[string][]*career.CareerEvent)
			for _, company := range companies {
				events := make([]*career.CareerEvent, eventsPerCompany)
				for i := 0; i < eventsPerCompany; i++ {
					eventID := fmt.Sprintf("%s-event-%d", company, i)
					event := fixtures.EventWith(
						eventID,
						fmt.Sprintf("Implemented feature %c at %s", rune('A'+i), company),
						company,
						"",
					)
					event.Tags = []string{"technical", "project"}
					event.Categories = []string{"technical"}
					event.Date = time.Now().AddDate(0, 0, -i)
					err := repos.Event.Create(ctx, event)
					Expect(err).NotTo(HaveOccurred())
					events[i] = event
				}
				result[company] = events
			}
			return result
		}

		// Helper to extract company names from generated CV
		extractCompaniesFromCV := func(cv *career.CVView) map[string]bool {
			companies := make(map[string]bool)
			for _, section := range cv.Sections {
				if section.SectionType == "experience" {
					for _, group := range section.Content {
						if group.Header != "" {
							companies[group.Header] = true
						}
					}
				}
			}
			return companies
		}

		It("should include ALL companies in generated CV when many events exist (BUG-003 fix)", func() {
			// This test reproduces the BUG-003 scenario:
			// - Multiple companies with events
			// - Previously, a total bullet cap would exclude later companies

			companies := []string{
				"Company Alpha",
				"Company Beta",
				"Company Gamma",
				"Company Delta",
				"Company Epsilon",
			}

			// Seed 10 events per company (50 total events)
			// Before the fix, a cap of 40 bullets would exclude some companies
			seedMultiCompanyEvents(companies, 10)

			// Verify all events were created
			allEvents, err := repos.Event.List(ctx, careerrepo.ListFilters{})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(allEvents)).To(Equal(50))

			// Generate CV with staff role (was previously capped at 40 bullets)
			config := &career.CVConfig{
				Name:           "bug-003-test",
				TargetRole:     "staff",
				TargetAudience: "hiring_manager",
				EventFilters:   make(map[string]interface{}),
			}

			cv, err := cvGenService.GenerateCVFromConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())
			Expect(cv).NotTo(BeNil())

			// CRITICAL ASSERTION: All 5 companies must be represented
			companiesInCV := extractCompaniesFromCV(cv)
			Expect(len(companiesInCV)).To(Equal(5),
				"BUG-003 regression: All 5 companies should appear in CV, got %d: %v",
				len(companiesInCV), companiesInCV)

			// Verify each specific company is present
			for _, company := range companies {
				Expect(companiesInCV).To(HaveKey(company),
					"BUG-003 regression: Company '%s' should appear in CV", company)
			}
		})

		It("should include all companies even with large event counts (stress test)", func() {
			// Stress test with more companies and events
			// Use unique company names to avoid interference with other tests
			companies := []string{
				"Stress Corp A", "Stress Corp B", "Stress Corp C",
				"Stress Corp D", "Stress Corp E", "Stress Corp F",
				"Stress Corp G", "Stress Corp H", "Stress Corp I",
				"Stress Corp J",
			}

			// Seed 15 events per company
			seedMultiCompanyEvents(companies, 15)

			// Generate CV
			config := &career.CVConfig{
				Name:           "bug-003-stress-test",
				TargetRole:     "principal",
				TargetAudience: "recruiter",
				EventFilters:   make(map[string]interface{}),
			}

			cv, err := cvGenService.GenerateCVFromConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())
			Expect(cv).NotTo(BeNil())

			// CRITICAL: All 10 stress companies must be represented
			companiesInCV := extractCompaniesFromCV(cv)
			for _, company := range companies {
				Expect(companiesInCV).To(HaveKey(company),
					"BUG-003 regression: Company '%s' should appear in CV", company)
			}
		})

		It("should maintain per-company bullet caps while including all companies", func() {
			// Verify that while all companies are included, per-company caps still apply
			companies := []string{"Alpha Inc", "Beta Inc", "Gamma Inc"}

			// Seed 20 events per company
			seedMultiCompanyEvents(companies, 20)

			config := &career.CVConfig{
				Name:           "bug-003-caps-test",
				TargetRole:     "senior_ic",
				TargetAudience: "peer",
				EventFilters:   make(map[string]interface{}),
			}

			cv, err := cvGenService.GenerateCVFromConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())
			Expect(cv).NotTo(BeNil())

			// All 3 companies must be represented
			companiesInCV := extractCompaniesFromCV(cv)
			Expect(len(companiesInCV)).To(Equal(3))

			// Verify per-company bullet caps are reasonable (not all 20 events per company)
			for _, section := range cv.Sections {
				if section.SectionType == "experience" {
					for _, group := range section.Content {
						// Per-company cap should limit bullets (typically 4-5 per company)
						Expect(len(group.Bullets)).To(BeNumerically("<=", 10),
							"Per-company bullet cap should limit bullets for %s", group.Header)
					}
				}
			}
		})
	})
})

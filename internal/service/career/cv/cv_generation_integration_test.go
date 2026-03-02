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
	careersql "github.com/baphled/kariya/internal/repository/career/sql"
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
		repos, err = careersql.NewRepositoriesFromPath(dbPath)
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

	seedTestEvents := func(count int, category string) []*career.Event {
		events := make([]*career.Event, count)
		for i := range count {
			eventID := fmt.Sprintf("%s-seed-%d", category, i)
			event := fixtures.EventWith(eventID, "Test event for "+category, "Test Company", "")
			event.Date = time.Now().AddDate(0, 0, -i)
			event.Tags = []string{category, "project"}
			event.Categories = []string{category}
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

			config := fixtures.CVConfigWithFilters("test-senior-ic", map[string]interface{}{
				"categories": []string{"technical", "product"},
			})
			config.TargetRole = "senior_ic"
			config.TargetAudience = "recruiter"

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

			config := fixtures.CVConfigWithFilters("test-principal", map[string]interface{}{
				"categories": []string{"technical", "achievement", "leadership"},
			})
			config.TargetRole = "principal"
			config.TargetAudience = "hiring_manager"

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
			allEvents, err := repos.Event.List(ctx, *fixtures.EventListFilters())
			Expect(err).NotTo(HaveOccurred())
			Expect(allEvents).To(HaveLen(7))

			config := fixtures.CVConfigWithFilters("test-no-filters", make(map[string]interface{}))
			config.TargetRole = "senior_ic"
			config.TargetAudience = "recruiter"

			// Generate CV
			cv, err := cvGenService.GenerateCVFromConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())
			Expect(cv).NotTo(BeNil())
		})

		It("should preserve CV metadata during generation", func() {
			// Seed test data
			seedTestEvents(3, "leadership")

			config := fixtures.CVConfigWithFilters("test-metadata", map[string]interface{}{
				"categories": []string{"leadership"},
			})
			config.TargetRole = "principal"
			config.TargetAudience = "hiring_manager"

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

			config := fixtures.CVConfigWithFilters("test-technical-only", map[string]interface{}{
				"categories": []string{"technical"},
			})
			config.TargetRole = "senior_ic"
			config.TargetAudience = "recruiter"

			// Generate CV
			cv, err := cvGenService.GenerateCVFromConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())
			Expect(cv).NotTo(BeNil())
		})

		It("should validate CV configuration before generation", func() {
			// Seed some test data (needed for valid case comparison)
			seedTestEvents(2, "technical")

			invalidConfig := fixtures.CVConfigWithFilters("test-invalid", make(map[string]interface{}))
			invalidConfig.TargetRole = ""
			invalidConfig.TargetAudience = "recruiter"

			// Generation should fail
			_, err := cvGenService.GenerateCVFromConfig(ctx, invalidConfig)
			Expect(err).To(HaveOccurred())
		})

		It("should handle empty database gracefully", func() {
			// Don't seed any data - database is empty

			config := fixtures.CVConfigWithFilters("test-empty-db", make(map[string]interface{}))
			config.TargetRole = "senior_ic"
			config.TargetAudience = "recruiter"

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
				config := fixtures.CVConfigWithFilters("test-"+audience, make(map[string]interface{}))
				config.TargetRole = "senior_ic"
				config.TargetAudience = audience

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
		seedMultiCompanyEvents := func(companies []string, eventsPerCompany int) map[string][]*career.Event {
			result := make(map[string][]*career.Event)
			for _, company := range companies {
				events := make([]*career.Event, eventsPerCompany)
				for i := range eventsPerCompany {
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
			allEvents, err := repos.Event.List(ctx, *fixtures.EventListFilters())
			Expect(err).NotTo(HaveOccurred())
			Expect(allEvents).To(HaveLen(50))

			config := fixtures.CVConfigWithFilters("bug-003-test", make(map[string]interface{}))

			cv, err := cvGenService.GenerateCVFromConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())
			Expect(cv).NotTo(BeNil())

			// CRITICAL ASSERTION: All 5 companies must be represented
			companiesInCV := extractCompaniesFromCV(cv)
			Expect(companiesInCV).To(HaveLen(5),
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

			config := fixtures.CVConfigWithFilters("bug-003-stress-test", make(map[string]interface{}))
			config.TargetRole = "principal"
			config.TargetAudience = "recruiter"

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

			config := fixtures.CVConfigWithFilters("bug-003-caps-test", make(map[string]interface{}))
			config.TargetRole = "senior_ic"
			config.TargetAudience = "peer"

			cv, err := cvGenService.GenerateCVFromConfig(ctx, config)
			Expect(err).NotTo(HaveOccurred())
			Expect(cv).NotTo(BeNil())

			// All 3 companies must be represented
			companiesInCV := extractCompaniesFromCV(cv)
			Expect(companiesInCV).To(HaveLen(3))

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

var _ = Describe("GenerateCVFromConfig additional branches", func() {
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

		var err error
		tempDir, err = os.MkdirTemp("", "kariya-cv-branch-test-")
		Expect(err).NotTo(HaveOccurred())

		dbPath := filepath.Join(tempDir, "test_events.db")
		repos, err = careersql.NewRepositoriesFromPath(dbPath)
		Expect(err).NotTo(HaveOccurred())

		bulletGenerator := NewBulletGenerator(log, nil)
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

	It("should return error for nil config", func() {
		cv, err := cvGenService.GenerateCVFromConfig(ctx, nil)
		Expect(err).To(HaveOccurred())
		Expect(cv).To(BeNil())
	})

	It("should return error for cancelled context", func() {
		cancelCtx, cancel := context.WithCancel(ctx)
		cancel()
		cv, err := cvGenService.GenerateCVFromConfig(cancelCtx, fixtures.CVConfig("cancelled-test"))
		Expect(err).To(HaveOccurred())
		Expect(cv).To(BeNil())
	})

	It("should return empty CV when no events match", func() {
		config := fixtures.CVConfigWithFilters("empty-cv", map[string]interface{}{
			"categories": []string{"nonexistent"},
		})
		config.TargetRole = "senior_ic"
		config.TargetAudience = "recruiter"

		cv, err := cvGenService.GenerateCVFromConfig(ctx, config)
		Expect(err).NotTo(HaveOccurred())
		Expect(cv).NotTo(BeNil())
		Expect(cv.SourceEventCount).To(Equal(0))
	})

	It("should apply length format constraints", func() {
		for i := range 5 {
			eventID := fmt.Sprintf("length-test-%d", i)
			event := fixtures.EventWith(eventID, "Built microservices architecture", "TechCorp", "")
			event.Date = time.Now().AddDate(0, 0, -i)
			event.Tags = []string{"technical"}
			event.Categories = []string{"technical"}
			err := repos.Event.Create(ctx, event)
			Expect(err).NotTo(HaveOccurred())
		}

		config := fixtures.CVConfigWithFilters("length-test", map[string]interface{}{})
		config.TargetRole = "senior_ic"
		config.TargetAudience = "recruiter"
		config.LengthFormat = "one_page"

		cv, err := cvGenService.GenerateCVFromConfig(ctx, config)
		Expect(err).NotTo(HaveOccurred())
		Expect(cv).NotTo(BeNil())
	})

	It("should apply technology focus filtering", func() {
		for i := range 5 {
			eventID := fmt.Sprintf("tech-focus-%d", i)
			event := fixtures.EventWith(eventID, "Developed Go microservices", "TechCorp", "")
			event.Date = time.Now().AddDate(0, 0, -i)
			event.Tags = []string{"go", "microservices"}
			event.Categories = []string{"technical"}
			err := repos.Event.Create(ctx, event)
			Expect(err).NotTo(HaveOccurred())
		}

		config := fixtures.CVConfigWithFilters("tech-focus-test", map[string]interface{}{})
		config.TargetRole = "senior_ic"
		config.TargetAudience = "recruiter"
		config.TechnologyFocus = string(TechnologyFocusSpecialist)
		config.SelectedTechnologies = []string{"go"}

		cv, err := cvGenService.GenerateCVFromConfig(ctx, config)
		Expect(err).NotTo(HaveOccurred())
		Expect(cv).NotTo(BeNil())
	})

	It("should apply skills format config", func() {
		for i := range 3 {
			eventID := fmt.Sprintf("skills-fmt-%d", i)
			event := fixtures.EventWith(eventID, "Led team delivery", "TechCorp", "")
			event.Date = time.Now().AddDate(0, 0, -i)
			event.Tags = []string{"leadership"}
			event.Categories = []string{"leadership"}
			err := repos.Event.Create(ctx, event)
			Expect(err).NotTo(HaveOccurred())
		}

		config := fixtures.CVConfigWithFilters("skills-fmt-test", map[string]interface{}{})
		config.TargetRole = "principal"
		config.TargetAudience = "hiring_manager"
		config.SkillsFormat = "grouped"
		config.SkillsLimit = 5

		cv, err := cvGenService.GenerateCVFromConfig(ctx, config)
		Expect(err).NotTo(HaveOccurred())
		Expect(cv).NotTo(BeNil())
	})
})

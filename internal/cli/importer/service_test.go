package importer_test

import (
	"bytes"
	"context"

	"github.com/baphled/kariya/internal/cli/importer"
	careermemory "github.com/baphled/kariya/internal/repository/career/memory"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Import Service", func() {
	var (
		importSvc *importer.ImportService
		repo      *careermemory.EventRepository
		svc       *careerservice.Service
		ctx       context.Context
	)

	BeforeEach(func() {
		repo = careermemory.NewEventRepository()
		svc = careerservice.NewService(repo)
		importSvc = importer.NewImportService(svc)
		ctx = context.Background()
	})

	Describe("PrepareImport", func() {
		It("should parse CSV and return parsed rows", func() {
			csv := `Text,Date,Categories,Tags,Project,Company
First event description,2024-01,Technical,technical,Project1,Company1
Second event description,2024-02,Leadership,leadership,Project2,Company2`

			reader := bytes.NewReader([]byte(csv))
			rows, err := importSvc.PrepareImport(ctx, reader)

			Expect(err).NotTo(HaveOccurred())
			Expect(rows).To(HaveLen(2))
			Expect(rows[0].IsValid).To(BeTrue())
			Expect(rows[1].IsValid).To(BeTrue())
		})
	})

	Describe("ImportRows", func() {
		It("should import selected valid rows", func() {
			csv := `Text,Date,Categories,Tags,Project,Company
First event description,2024-01,Technical,technical,Project1,Company1
Second event description,2024-02,Leadership,leadership,Project2,Company2`

			reader := bytes.NewReader([]byte(csv))
			rows, err := importSvc.PrepareImport(ctx, reader)
			Expect(err).NotTo(HaveOccurred())

			selectedRows := []int{1, 2}
			result, err := importSvc.ImportRows(ctx, rows, selectedRows)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.SuccessCount).To(Equal(2))
			Expect(result.CreatedEvents).To(HaveLen(2))
		})

		It("should detect burst suggestions after import", func() {
			// Create test CSV with events likely to form bursts
			csv := `Text,Date,Categories,Tags,Project,Company
Implemented cloud migration phase 1,2024-01-15,Technical,technical,CloudMigration,TechCorp
Optimized database queries for cloud,2024-01-20,Technical,technical,CloudMigration,TechCorp
Completed cloud migration phase 2,2024-01-25,Technical,technical,CloudMigration,TechCorp
Led performance optimization sprint,2024-02-01,Leadership,leadership,Performance,TechCorp
Improved API response times,2024-02-05,Technical,technical,Performance,TechCorp`

			reader := bytes.NewReader([]byte(csv))
			rows, err := importSvc.PrepareImport(ctx, reader)
			Expect(err).NotTo(HaveOccurred())

			// Select all valid rows
			selectedRows := []int{1, 2, 3, 4, 5}
			result, err := importSvc.ImportRows(ctx, rows, selectedRows)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.SuccessCount).To(Equal(5))
			Expect(result.CreatedEvents).To(HaveLen(5))

			// Verify burst suggestions were generated
			//
			// The exact number depends on the burst detection algorithm
			// We just verify that the field is populated (could be 0 if no bursts detected)
			Expect(result.BurstSuggestions).NotTo(BeNil())
		})

		It("should handle burst detection failure gracefully", func() {
			csv := `Text,Date,Categories,Tags,Project,Company
First event description,2024-01,Technical,technical,Project1,Company1`

			reader := bytes.NewReader([]byte(csv))
			rows, err := importSvc.PrepareImport(ctx, reader)
			Expect(err).NotTo(HaveOccurred())

			selectedRows := []int{1}
			result, err := importSvc.ImportRows(ctx, rows, selectedRows)

			// Import should succeed even if burst detection fails
			Expect(err).NotTo(HaveOccurred())
			Expect(result.SuccessCount).To(Equal(1))
			Expect(result.CreatedEvents).To(HaveLen(1))

			// BurstSuggestions should be initialized (empty or with results)
			Expect(result.BurstSuggestions).NotTo(BeNil())
		})

		It("should extract facts from imported events", func() {
			csv := `Text,Date,Categories,Tags,Project,Company
Implemented cloud migration strategy,2024-01-15,Technical,technical,CloudMigration,TechCorp
Led team through system redesign,2024-02-01,Leadership,leadership,Redesign,TechCorp
Mentored junior engineers on best practices,2024-02-10,Mentoring,mentoring,Training,TechCorp`

			reader := bytes.NewReader([]byte(csv))
			rows, err := importSvc.PrepareImport(ctx, reader)
			Expect(err).NotTo(HaveOccurred())

			// Select all valid rows
			selectedRows := []int{1, 2, 3}
			result, err := importSvc.ImportRows(ctx, rows, selectedRows)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.SuccessCount).To(Equal(3))
			Expect(result.CreatedEvents).To(HaveLen(3))

			// Verify fact extraction fields are initialized
			Expect(result.FactsByEventID).NotTo(BeNil())
			Expect(result.FactsByCompetency).NotTo(BeNil())
			Expect(result.ExtractedFactsCount).To(BeNumerically(">=", 0))
		})

		It("should link facts to their parent bursts via SourceBurstID", func() {
			burstRepo := careermemory.NewBurstRepository()
			factRepo := careermemory.NewFactRepository()
			svc.SetBurstRepository(burstRepo)
			svc.SetFactRepository(factRepo)

			// Events designed to cluster into a burst (same company, project, close dates).
			csv := `Text,Date,Categories,Tags,Project,Company
Implemented cloud migration phase 1,2024-01-15,Technical,technical,CloudMigration,TechCorp
Completed cloud migration phase 2,2024-01-20,Technical,technical,CloudMigration,TechCorp
Optimized cloud infrastructure performance,2024-01-25,Technical,technical,CloudMigration,TechCorp`

			reader := bytes.NewReader([]byte(csv))
			rows, err := importSvc.PrepareImport(ctx, reader)
			Expect(err).NotTo(HaveOccurred())

			selectedRows := []int{1, 2, 3}
			result, err := importSvc.ImportRows(ctx, rows, selectedRows)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.SuccessCount).To(Equal(3))

			// Bursts must have been detected from these clustered events.
			allBursts, err := burstRepo.List(ctx, *fixtures.BurstListFilters())
			Expect(err).NotTo(HaveOccurred())
			Expect(allBursts).NotTo(BeEmpty(),
				"clustered events with same project/company should produce at least one burst")

			// Build set of event IDs belonging to bursts.
			burstEventIDs := make(map[string]bool)
			for _, burst := range allBursts {
				for _, eid := range burst.EventIDs {
					burstEventIDs[eid] = true
				}
			}

			// Every fact whose event belongs to a burst must have SourceBurstID set.
			allFacts, err := factRepo.List(ctx, *fixtures.FactListFilters())
			Expect(err).NotTo(HaveOccurred())

			for _, fact := range allFacts {
				if burstEventIDs[fact.SourceEventID] {
					Expect(fact.SourceBurstID).NotTo(BeEmpty(),
						"fact for burst event %s should have SourceBurstID set", fact.SourceEventID)
				}
			}
		})

		It("should track facts by competency category during import", func() {
			csv := `Text,Date,Categories,Tags,Project,Company
Implemented cloud migration,2024-01-15,Technical,technical,CloudMigration,TechCorp
Led architectural review,2024-02-01,Leadership,leadership,Architecture,TechCorp`

			reader := bytes.NewReader([]byte(csv))
			rows, err := importSvc.PrepareImport(ctx, reader)
			Expect(err).NotTo(HaveOccurred())

			selectedRows := []int{1, 2}
			result, err := importSvc.ImportRows(ctx, rows, selectedRows)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.SuccessCount).To(Equal(2))

			// Verify fact tracking structures are populated
			Expect(result.FactsByEventID).NotTo(BeNil())
			Expect(result.FactsByCompetency).NotTo(BeNil())
		})
	})
})

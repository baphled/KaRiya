package importer_test

import (
	"bytes"
	"context"
	"time"

	"github.com/baphled/kariya/internal/cli/importer"
	careerepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Import Service", func() {
	var (
		importSvc *importer.ImportService
		repo      *careerepo.MemoryRepository
		svc       *careerservice.Service
		ctx       context.Context
	)

	BeforeEach(func() {
		repo = careerepo.NewMemoryRepository()
		svc = careerservice.NewService(repo)
		importSvc = importer.NewImportService(svc)
		ctx = context.Background()
	})

	Describe("PrepareImport", func() {
		It("should parse CSV and return parsed rows", func() {
			csv := `Text,Date,Categories,Tags,Project,Company
Event 1,2024-01,Technical,technical,Project1,Company1
Event 2,2024-02,Leadership,leadership,Project2,Company2`

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
Event 1,2024-01,Technical,technical,Project1,Company1
Event 2,2024-02,Leadership,leadership,Project2,Company2`

			reader := bytes.NewReader([]byte(csv))
			rows, err := importSvc.PrepareImport(ctx, reader)
			Expect(err).NotTo(HaveOccurred())

			selectedRows := []int{1, 2}
			result, err := importSvc.ImportRows(ctx, rows, selectedRows)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.SuccessCount).To(Equal(2))
			Expect(result.CreatedEvents).To(HaveLen(2))
		})
	})
})

func parseDate(dateStr string) time.Time {
	t, _ := time.Parse("2006-01", dateStr)
	return t
}

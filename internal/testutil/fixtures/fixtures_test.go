package fixtures_test

import (
	"testing"

	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/brianvoe/gofakeit/v7"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFixtures(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Fixtures Suite")
}

var _ = Describe("Reproducibility", func() {
	It("should produce same data with same seed", func() {
		fixtures.SetSeed(12345)
		company1 := gofakeit.Company()
		project1 := gofakeit.BuzzWord()

		fixtures.SetSeed(12345)
		company2 := gofakeit.Company()
		project2 := gofakeit.BuzzWord()

		Expect(company1).To(Equal(company2))
		Expect(project1).To(Equal(project2))
	})
})

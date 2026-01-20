package fixtures_test

import (
	"testing"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
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
		event1 := fixtures.EventFactory.MustCreate().(*career.CareerEvent)

		fixtures.SetSeed(12345)
		event2 := fixtures.EventFactory.MustCreate().(*career.CareerEvent)

		// With same seed, random parts should match
		// Note: IDs are sequential so they reset with new factory
		Expect(event1.Company).To(Equal(event2.Company))
		Expect(event1.Project).To(Equal(event2.Project))
	})
})

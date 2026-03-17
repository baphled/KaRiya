package captureevent_test

import (
	"github.com/baphled/kariya/internal/domain/career"

	"github.com/baphled/kariya/internal/testutil/fixtures"
	ce "github.com/baphled/kariya/internal/tui/intents/captureevent"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Result", func() {
	Describe("Construction", func() {
		It("should create a result with all fields", func() {
			event := fixtures.EventWith("", "test event", "", "")
			bursts := []*career.Burst{fixtures.Burst("burst-1")}
			facts := []*career.Fact{fixtures.FactWith("fact-1", "fact-1")}
			accepted := map[string]bool{"burst-1": true}
			rejected := map[string]string{"fact-2": "not relevant"}

			result := &ce.Result{
				Event:          event,
				Bursts:         bursts,
				Facts:          facts,
				AcceptedFields: accepted,
				RejectedFields: rejected,
			}

			Expect(result.Event).To(Equal(event))
			Expect(result.Bursts).To(HaveLen(1))
			Expect(result.Facts).To(HaveLen(1))
			Expect(result.AcceptedFields).To(HaveKeyWithValue("burst-1", true))
			Expect(result.RejectedFields).To(HaveKeyWithValue("fact-2", "not relevant"))
		})

		It("should allow nil fields for empty results", func() {
			result := &ce.Result{}
			Expect(result.Event).To(BeNil())
			Expect(result.Bursts).To(BeNil())
			Expect(result.Facts).To(BeNil())
			Expect(result.AcceptedFields).To(BeNil())
			Expect(result.RejectedFields).To(BeNil())
		})
	})
})

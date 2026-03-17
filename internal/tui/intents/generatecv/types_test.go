package generatecv_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/baphled/kariya/internal/tui/intents/generatecv"
)

var _ = Describe("Types", func() {
	Describe("Intent construction via NewIntent", func() {
		It("should produce a non-nil intent with valid context", func() {
			ctx := &generatecv.IntentValidator{
				AvailableProfiles: []*generatecv.CVProfile{{ID: "p1"}},
				Events:            []*career.Event{fixtures.Event("e1")},
			}

			intent, err := generatecv.NewIntent(ctx)

			Expect(err).ToNot(HaveOccurred())
			Expect(intent).NotTo(BeNil())
		})
	})
})

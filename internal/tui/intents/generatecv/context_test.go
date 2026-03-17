package generatecv_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/baphled/kariya/internal/tui/intents/generatecv"
)

var _ = Describe("IntentValidator", func() {
	Describe("Construction", func() {
		It("should create a context with profiles and events", func() {
			ctx := &generatecv.IntentValidator{
				AvailableProfiles: []*generatecv.CVProfile{
					{ID: "p1", Name: "Senior Dev"},
				},
				Events: []*career.Event{
					fixtures.Event("event-1"),
				},
			}
			Expect(ctx.AvailableProfiles).To(HaveLen(1))
			Expect(ctx.Events).To(HaveLen(1))
		})

		It("should store a default profile", func() {
			profile := &generatecv.CVProfile{ID: "p1", Name: "Default"}
			ctx := &generatecv.IntentValidator{
				AvailableProfiles: []*generatecv.CVProfile{profile},
				Events:            []*career.Event{fixtures.Event("e1")},
				DefaultProfile:    profile,
			}
			Expect(ctx.DefaultProfile).To(Equal(profile))
		})

		It("should store facts", func() {
			ctx := &generatecv.IntentValidator{
				AvailableProfiles: []*generatecv.CVProfile{{ID: "p1"}},
				Events:            []*career.Event{fixtures.Event("e1")},
				Facts:             []*career.Fact{fixtures.Fact("f1", "e1")},
			}
			Expect(ctx.Facts).To(HaveLen(1))
		})
	})

	Describe("Validate", func() {
		Context("with valid context", func() {
			It("should return nil", func() {
				ctx := &generatecv.IntentValidator{
					AvailableProfiles: []*generatecv.CVProfile{
						{ID: "p1", Name: "Senior Dev"},
					},
					Events: []*career.Event{
						fixtures.Event("event-1"),
					},
				}
				Expect(ctx.Validate()).To(Succeed())
			})
		})

		Context("with nil context", func() {
			It("should return error", func() {
				var ctx *generatecv.IntentValidator
				Expect(ctx.Validate()).To(MatchError("IntentValidator cannot be nil"))
			})
		})

		Context("with empty profiles", func() {
			It("should return error", func() {
				ctx := &generatecv.IntentValidator{
					AvailableProfiles: []*generatecv.CVProfile{},
					Events:            []*career.Event{fixtures.Event("e1")},
				}
				Expect(ctx.Validate()).To(MatchError("IntentValidator must have at least one available profile"))
			})
		})

		Context("with nil profiles", func() {
			It("should return error", func() {
				ctx := &generatecv.IntentValidator{
					Events: []*career.Event{fixtures.Event("e1")},
				}
				Expect(ctx.Validate()).To(MatchError("IntentValidator must have at least one available profile"))
			})
		})

		Context("with empty events", func() {
			It("should return error", func() {
				ctx := &generatecv.IntentValidator{
					AvailableProfiles: []*generatecv.CVProfile{{ID: "p1"}},
					Events:            []*career.Event{},
				}
				Expect(ctx.Validate()).To(MatchError("IntentValidator must have at least one event"))
			})
		})

		Context("with nil events", func() {
			It("should return error", func() {
				ctx := &generatecv.IntentValidator{
					AvailableProfiles: []*generatecv.CVProfile{{ID: "p1"}},
				}
				Expect(ctx.Validate()).To(MatchError("IntentValidator must have at least one event"))
			})
		})
	})
})

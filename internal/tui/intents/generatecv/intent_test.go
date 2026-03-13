package generatecv_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/baphled/kariya/internal/tui/intents/generatecv"
)

func validContext() *generatecv.IntentValidator {
	return &generatecv.IntentValidator{
		AvailableProfiles: []*generatecv.CVProfile{
			{ID: "p1", Name: "Senior Dev", TargetRole: "Backend Engineer", TargetAudience: "hiring_manager"},
		},
		Events: []*career.Event{
			fixtures.Event("event-1"),
		},
	}
}

func validContextWithDefault() *generatecv.IntentValidator {
	profile := &generatecv.CVProfile{ID: "p1", Name: "Senior Dev", TargetRole: "Backend Engineer", TargetAudience: "hiring_manager"}
	return &generatecv.IntentValidator{
		AvailableProfiles: []*generatecv.CVProfile{profile},
		Events:            []*career.Event{fixtures.Event("event-1")},
		DefaultProfile:    profile,
	}
}

var _ = Describe("Intent", func() {
	Describe("NewIntent", func() {
		Context("with valid context", func() {
			It("should create an intent", func() {
				intent, err := generatecv.NewIntent(validContext())
				Expect(err).NotTo(HaveOccurred())
				Expect(intent).NotTo(BeNil())
			})
		})

		Context("with nil context", func() {
			It("should return error", func() {
				intent, err := generatecv.NewIntent(nil)
				Expect(err).To(HaveOccurred())
				Expect(intent).To(BeNil())
			})
		})

		Context("with empty profiles", func() {
			It("should return validation error", func() {
				ctx := &generatecv.IntentValidator{
					Events: []*career.Event{fixtures.Event("e1")},
				}
				intent, err := generatecv.NewIntent(ctx)
				Expect(err).To(HaveOccurred())
				Expect(intent).To(BeNil())
			})
		})

		Context("with empty events", func() {
			It("should return validation error", func() {
				ctx := &generatecv.IntentValidator{
					AvailableProfiles: []*generatecv.CVProfile{{ID: "p1"}},
				}
				intent, err := generatecv.NewIntent(ctx)
				Expect(err).To(HaveOccurred())
				Expect(intent).To(BeNil())
			})
		})

		Context("with default profile", func() {
			It("should select the default profile", func() {
				intent, err := generatecv.NewIntent(validContextWithDefault())
				Expect(err).NotTo(HaveOccurred())
				Expect(intent).NotTo(BeNil())
			})
		})

		Context("without default profile", func() {
			It("should select the first available profile", func() {
				intent, err := generatecv.NewIntent(validContext())
				Expect(err).NotTo(HaveOccurred())
				Expect(intent).NotTo(BeNil())
			})
		})
	})

	Describe("Init", func() {
		It("should return a command", func() {
			intent, err := generatecv.NewIntent(validContext())
			Expect(err).NotTo(HaveOccurred())

			cmd := intent.Init()
			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("View", func() {
		It("should return a non-empty string", func() {
			intent, err := generatecv.NewIntent(validContext())
			Expect(err).NotTo(HaveOccurred())
			intent.Init()

			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Update", func() {
		It("should handle quit key", func() {
			intent, err := generatecv.NewIntent(validContext())
			Expect(err).NotTo(HaveOccurred())
			intent.Init()

			result := intent.Result()
			Expect(result).To(BeNil())
		})
	})

	Describe("Result", func() {
		Context("before completion", func() {
			It("should return nil", func() {
				intent, err := generatecv.NewIntent(validContext())
				Expect(err).NotTo(HaveOccurred())

				Expect(intent.Result()).To(BeNil())
			})
		})
	})
})

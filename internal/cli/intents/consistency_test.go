package intents

import (
	"context"
	"strings"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// testStandardViewConsistency checks that an intent's view follows StandardView patterns
func testStandardViewConsistency(intentName, view string) {
	Expect(view).NotTo(BeEmpty(), "[%s] View is empty", intentName)

	// Check view is substantial (not just whitespace)
	trimmed := strings.TrimSpace(view)
	Expect(len(trimmed)).To(BeNumerically(">=", 50),
		"[%s] View seems too short (%d chars), may be incomplete", intentName, len(trimmed))
}

var _ = Describe("StandardView Consistency", func() {
	Describe("CaptureEvent", func() {
		It("should use StandardView patterns", func() {
			ctx := &CaptureEventContext{
				CaptureStrategy: "manual",
				PreviousEvent:   nil,
				Metadata:        make(map[string]string),
			}

			intent, err := NewCaptureEventIntent(ctx)
			Expect(err).NotTo(HaveOccurred())

			intent.Init()
			view := intent.View()

			testStandardViewConsistency("CaptureEvent", view)
		})
	})

	// BrowseTimeline consistency tests are in internal/cli/intents/browse_timeline/

	Describe("GenerateCV", func() {
		It("should use StandardView patterns", func() {
			ctx := &GenerateCVContext{
				AvailableProfiles: []*CVProfile{
					{
						ID:             "default",
						Name:           "Default Profile",
						TargetRole:     "staff",
						TargetAudience: "hiring_manager",
					},
				},
				Events: []*career.CareerEvent{
					fixtures.EventWith(uuid.New().String(), "Implemented test feature for CV generation", "TechCorp", "Platform"),
				},
			}

			intent, err := NewGenerateCVIntent(ctx)
			Expect(err).NotTo(HaveOccurred())

			intent.Init()
			view := intent.View()

			testStandardViewConsistency("GenerateCV", view)
		})
	})

	Describe("ConfigureSystem", func() {
		It("should use StandardView patterns", func() {
			ctx := context.Background()

			intent, err := NewConfigureSystemIntent(ctx)
			Expect(err).NotTo(HaveOccurred())

			intent.Init()
			view := intent.View()

			testStandardViewConsistency("ConfigureSystem", view)
		})
	})

	Describe("FactManagement", func() {
		It("should use StandardView patterns", func() {
			ctx := NewFactManagementContext(nil, context.Background())

			intent := NewFactManagementIntent(ctx)
			Expect(intent).NotTo(BeNil())

			intent.Init()
			view := intent.View()

			testStandardViewConsistency("FactManagement", view)
		})
	})

	Describe("BurstManagement", func() {
		It("should use StandardView patterns", func() {
			ctx := NewBurstManagementContext(nil, nil, context.Background())

			intent, err := NewBurstManagementIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(intent).NotTo(BeNil())

			intent.Init()
			view := intent.View()

			testStandardViewConsistency("BurstManagement", view)
		})
	})
})

var _ = Describe("All Intents Initialization", func() {
	DescribeTable("should initialize successfully",
		func(name string, createFunc func() (interface{}, error), initFunc func(interface{}), viewFunc func(interface{}) string) {
			intent, err := createFunc()
			Expect(err).NotTo(HaveOccurred(), "Failed to create %s intent", name)

			initFunc(intent)
			view := viewFunc(intent)

			Expect(view).NotTo(BeEmpty(), "%s produced empty view", name)
			Expect(len(strings.TrimSpace(view))).To(BeNumerically(">=", 10),
				"%s view is too short", name)
		},
		Entry("CaptureEvent",
			"CaptureEvent",
			func() (interface{}, error) {
				return NewCaptureEventIntent(&CaptureEventContext{
					CaptureStrategy: "manual",
					Metadata:        make(map[string]string),
				})
			},
			func(i interface{}) { _ = i.(*CaptureEventIntent).Init() }, //nolint:errcheck // Init returns tea.Cmd which is intentionally discarded in tests
			func(i interface{}) string { return i.(*CaptureEventIntent).View() },
		),
		// BrowseTimeline Entry is in internal/cli/intents/browse_timeline/
		Entry("GenerateCV",
			"GenerateCV",
			func() (interface{}, error) {
				return NewGenerateCVIntent(&GenerateCVContext{
					AvailableProfiles: []*CVProfile{
						{
							ID:             "default",
							Name:           "Default Profile",
							TargetRole:     "staff",
							TargetAudience: "hiring_manager",
						},
					},
					Events: []*career.CareerEvent{
						fixtures.EventWith(uuid.New().String(), "Implemented test feature for CV generation", "TechCorp", "Platform"),
					},
				})
			},
			func(i interface{}) { _ = i.(*GenerateCVIntent).Init() }, //nolint:errcheck // Init returns tea.Cmd which is intentionally discarded in tests
			func(i interface{}) string { return i.(*GenerateCVIntent).View() },
		),
		Entry("ConfigureSystem",
			"ConfigureSystem",
			func() (interface{}, error) { return NewConfigureSystemIntent(context.Background()) },
			func(i interface{}) { _ = i.(*ConfigureSystemIntent).Init() }, //nolint:errcheck // Init returns tea.Cmd which is intentionally discarded in tests
			func(i interface{}) string { return i.(*ConfigureSystemIntent).View() },
		),
	)
})

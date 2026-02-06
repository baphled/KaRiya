package intents

import (
	"context"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// testStandardViewConsistency checks that an intent's view follows StandardView patterns.
func testStandardViewConsistency(intentName, view string) {
	Expect(view).NotTo(BeEmpty(), "[%s] View is empty", intentName)

	// Check view is substantial (not just whitespace)
	trimmed := strings.TrimSpace(view)
	Expect(len(trimmed)).To(BeNumerically(">=", 50),
		"[%s] View seems too short (%d chars), may be incomplete", intentName, len(trimmed))
}

var _ = Describe("StandardView Consistency", func() {
	// CaptureEvent consistency tests are in internal/cli/intents/captureevent/

	// BrowseTimeline consistency tests are in internal/cli/intents/browsetimeline/

	// GenerateCV consistency tests are in internal/cli/intents/generatecv/

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
		// CaptureEvent Entry is in internal/cli/intents/captureevent/
		// BrowseTimeline Entry is in internal/cli/intents/browsetimeline/
		// GenerateCV Entry is in internal/cli/intents/generatecv/
		Entry("ConfigureSystem",
			"ConfigureSystem",
			func() (interface{}, error) { return NewConfigureSystemIntent(context.Background()) },
			func(i interface{}) { _ = i.(*ConfigureSystemIntent).Init() }, //nolint:errcheck // Init returns tea.Cmd which is intentionally discarded in tests
			func(i interface{}) string { return i.(*ConfigureSystemIntent).View() },
		),
	)
})

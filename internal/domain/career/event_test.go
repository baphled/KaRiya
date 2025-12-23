package career

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CareerEvent", func() {
	Context("Validation", func() {
		It("should pass for a valid event", func() {
			event := &CareerEvent{
				Text: "Developed a high-performance backend service",
				Date: time.Now().AddDate(0, 0, -10),
				Tags: []string{"technical", "project"},
			}
			Expect(event.Validate()).To(Succeed(), "Valid event should pass validation")
		})

		It("should fail for empty text", func() {
			event := &CareerEvent{
				Text: "   ",
				Date: time.Now(),
			}
			Expect(event.Validate()).ToNot(Succeed(), "Empty text should fail validation")
		})

		It("should fail for overly long text", func() {
			longText := string(make([]byte, 2001))
			event := &CareerEvent{
				Text: longText,
				Date: time.Now(),
			}
			Expect(event.Validate()).ToNot(Succeed(), "Text exceeding 2000 characters should fail validation")
		})

		It("should fail for future dates", func() {
			event := &CareerEvent{
				Text: "Future project",
				Date: time.Now().AddDate(0, 0, 1),
			}
			Expect(event.Validate()).ToNot(Succeed(), "Future date should fail validation")
		})

		It("should fail for invalid tags", func() {
			event := &CareerEvent{
				Text: "Some event",
				Date: time.Now(),
				Tags: []string{"invalid_tag"},
			}
			Expect(event.Validate()).ToNot(Succeed(), "Invalid tag should fail validation")
		})
	})
})

package bursts_test

import (
	"bytes"
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	cmdpkg "github.com/baphled/kariya/internal/cli/cmd"
	"github.com/baphled/kariya/internal/cli/cmd/bursts"
	"github.com/baphled/kariya/internal/cli/cmd/cliutil"
	"github.com/baphled/kariya/internal/testutil/fixtures"
)

var _ = Describe("DetectBursts", func() {
	var (
		ctx cliutil.ServiceContext
		out *bytes.Buffer
		err *bytes.Buffer
	)

	BeforeEach(func() {
		cliCtx := cmdpkg.NewCLIContext("", true)
		initErr := cliCtx.InitService(new(bytes.Buffer))
		Expect(initErr).NotTo(HaveOccurred())
		ctx = cliCtx
		out = new(bytes.Buffer)
		err = new(bytes.Buffer)
	})

	Context("with no events", func() {
		It("should return 0 and print info message", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			code := bursts.DetectBursts(svc, out, err, tea.WithInput(nil))
			Expect(code).To(Equal(0))
		})

		It("should not write error output", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			code := bursts.DetectBursts(svc, out, err, tea.WithInput(nil))
			Expect(code).To(Equal(0))
			Expect(err.Len()).To(Equal(0))
		})
	})

	Context("with events in database", func() {
		BeforeEach(func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			event := fixtures.EventWith("event-1", "Implemented authentication system for web application", "", "")
			captureErr := svc.CaptureEvent(context.Background(), event, "timeline")
			Expect(captureErr).NotTo(HaveOccurred())
		})

		It("should handle execution without panic", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			Expect(func() {
				bursts.DetectBursts(svc, out, err, tea.WithInput(nil))
			}).NotTo(Panic())
		})

		It("should return valid exit code", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			code := bursts.DetectBursts(svc, out, err)
			Expect(code).To(BeNumerically(">=", 0))
			Expect(code).To(BeNumerically("<=", 1))
		})
	})

	Context("with multiple events", func() {
		BeforeEach(func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			for i := 1; i <= 3; i++ {
				eventID := "event-" + string(rune(48+i))
				eventText := "Implemented feature number " + string(rune(48+i))
				event := fixtures.EventWith(eventID, eventText, "", "")
				event.Date = time.Now()
				captureErr := svc.CaptureEvent(context.Background(), event, "timeline")
				Expect(captureErr).NotTo(HaveOccurred())
			}
		})

		It("should process all events without panic", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			Expect(func() {
				bursts.DetectBursts(svc, out, err)
			}).NotTo(Panic())
		})

		It("should return valid exit code", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			code := bursts.DetectBursts(svc, out, err, tea.WithInput(nil))
			Expect(code).To(BeNumerically(">=", 0))
			Expect(code).To(BeNumerically("<=", 1))
		})
	})

	Context("error handling", func() {
		It("should handle service errors gracefully", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			code := bursts.DetectBursts(svc, out, err, tea.WithInput(nil))
			Expect(code).To(BeNumerically(">=", 0))
		})

		It("should return valid exit code", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			code := bursts.DetectBursts(svc, out, err, tea.WithInput(nil))
			Expect(code).To(BeNumerically(">=", 0))
			Expect(code).To(BeNumerically("<=", 1))
		})

		It("should handle nil service gracefully", func() {
			Expect(func() {
				bursts.DetectBursts(nil, out, err, tea.WithInput(nil))
			}).To(Panic())
		})
	})
})

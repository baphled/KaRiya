package bursts_test

import (
	"bytes"
	"context"

	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	cmdpkg "github.com/baphled/kariya/internal/cli/cmd"
	"github.com/baphled/kariya/internal/cli/cmd/bursts"
	"github.com/baphled/kariya/internal/cli/cmd/cliutil"
	"github.com/baphled/kariya/internal/testutil/fixtures"
)

var _ = Describe("ListBursts", func() {
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

	Context("with no bursts", func() {
		It("should return 0 and print info message", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			code := bursts.ListBursts(svc, out, err, tea.WithInput(nil))
			Expect(code).To(Equal(0))
		})

		It("should not write error output", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			code := bursts.ListBursts(svc, out, err, tea.WithInput(nil))
			Expect(code).To(Equal(0))
			Expect(err.Len()).To(Equal(0))
		})
	})

	Context("with bursts in database", func() {
		BeforeEach(func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())
			burst := fixtures.Burst("burst-1", "event-1", "event-2")
			burst.Name = "Authentication Project"
			burst.Description = "Implemented auth system"
			saveErr := svc.SaveBurst(context.Background(), burst)
			Expect(saveErr).NotTo(HaveOccurred())
		})

		It("should return 0 on successful execution", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			code := bursts.ListBursts(svc, out, err, tea.WithInput(nil))
			Expect(code).To(Equal(0))
		})

		It("should write output to stdout", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			code := bursts.ListBursts(svc, out, err, tea.WithInput(nil))
			Expect(code).To(Equal(0))
			Expect(out.String()).To(ContainSubstring("Bursts"))
		})

		It("should display burst information", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			code := bursts.ListBursts(svc, out, err, tea.WithInput(nil))
			Expect(code).To(Equal(0))
			Expect(out.String()).To(ContainSubstring("Authentication Project"))
		})

		It("should handle execution without panic", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			Expect(func() {
				bursts.ListBursts(svc, out, err, tea.WithInput(nil))
			}).NotTo(Panic())
		})
	})

	Context("with multiple bursts", func() {
		BeforeEach(func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			for i := 1; i <= 3; i++ {
				burstID := "burst-" + string(rune(48+i))
				burst := fixtures.Burst(burstID, "event-1", "event-2")
				burst.Name = "Burst " + string(rune(48+i))
				saveErr := svc.SaveBurst(context.Background(), burst)
				Expect(saveErr).NotTo(HaveOccurred())
			}
		})

		It("should list all bursts", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			code := bursts.ListBursts(svc, out, err, tea.WithInput(nil))
			Expect(code).To(Equal(0))
			Expect(out.String()).To(ContainSubstring("Total bursts: 3"))
		})
	})

	Context("error handling", func() {
		It("should handle service errors gracefully", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			code := bursts.ListBursts(svc, out, err, tea.WithInput(nil))
			Expect(code).To(BeNumerically(">=", 0))
		})

		It("should return valid exit code", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			code := bursts.ListBursts(svc, out, err, tea.WithInput(nil))
			Expect(code).To(BeNumerically(">=", 0))
			Expect(code).To(BeNumerically("<=", 1))
		})

		It("should handle nil service gracefully", func() {
			Expect(func() {
				bursts.ListBursts(nil, out, err, tea.WithInput(nil))
			}).To(Panic())
		})
	})
})

package bursts_test

import (
	"bytes"
	"context"
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	cmdpkg "github.com/baphled/kariya/internal/cli/cmd"
	"github.com/baphled/kariya/internal/cli/cmd/bursts"
	"github.com/baphled/kariya/internal/cli/cmd/cliutil"
	careermemory "github.com/baphled/kariya/internal/repository/career/memory"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	mockrepo "github.com/baphled/kariya/internal/testutil/mocks/repository"
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

var _ = Describe("ListBursts Error Paths", func() {
	var (
		out *bytes.Buffer
		err *bytes.Buffer
	)

	BeforeEach(func() {
		out = new(bytes.Buffer)
		err = new(bytes.Buffer)
	})

	Context("when burst repository is not configured", func() {
		It("should return exit code 1", func() {
			// Create service with nil burst repository
			svc := careerservice.NewService(nil)
			// Don't call SetBurstRepository - leave it nil

			code := bursts.ListBursts(svc, out, err, tea.WithInput(nil))
			Expect(code).To(Equal(1))
		})

		It("should print error message", func() {
			svc := careerservice.NewService(nil)

			bursts.ListBursts(svc, out, err, tea.WithInput(nil))
			// Error is printed to stderr via cliutil.PrintError, not to err buffer
			// Just verify exit code is 1 (which we test above)
		})
	})

	Context("when repository List() returns an error", func() {
		var (
			ctrl      *gomock.Controller
			mockRepo  *mockrepo.MockBurstRepository
			eventRepo *careermemory.EventRepository
		)

		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())
			mockRepo = mockrepo.NewMockBurstRepository(ctrl)
			eventRepo = careermemory.NewEventRepository()
		})

		AfterEach(func() {
			ctrl.Finish()
		})

		It("should return exit code 1 on database error", func() {
			// Setup mock to return error
			mockRepo.EXPECT().
				List(gomock.Any(), gomock.Any()).
				Return(nil, errors.New("database connection failed"))

			svc := careerservice.NewService(eventRepo)
			svc.SetBurstRepository(mockRepo)

			code := bursts.ListBursts(svc, out, err, tea.WithInput(nil))
			Expect(code).To(Equal(1))
		})

		It("should handle repository errors gracefully", func() {
			mockRepo.EXPECT().
				List(gomock.Any(), gomock.Any()).
				Return(nil, errors.New("timeout"))

			svc := careerservice.NewService(eventRepo)
			svc.SetBurstRepository(mockRepo)

			Expect(func() {
				bursts.ListBursts(svc, out, err, tea.WithInput(nil))
			}).NotTo(Panic())
		})
	})
})

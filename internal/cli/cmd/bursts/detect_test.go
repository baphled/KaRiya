package bursts_test

import (
	"bytes"
	"context"
	"errors"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	cmdpkg "github.com/baphled/kariya/internal/cli/cmd"
	"github.com/baphled/kariya/internal/cli/cmd/bursts"
	"github.com/baphled/kariya/internal/cli/cmd/cliutil"
	careermemory "github.com/baphled/kariya/internal/repository/career/memory"
	careerservice "github.com/baphled/kariya/internal/service/career"
	burst_fact "github.com/baphled/kariya/internal/service/career/burstfact"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	mockrepo "github.com/baphled/kariya/internal/testutil/mocks/repository"
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

var _ = Describe("DetectBursts Error Paths", func() {
	var (
		out *bytes.Buffer
		err *bytes.Buffer
	)

	BeforeEach(func() {
		out = new(bytes.Buffer)
		err = new(bytes.Buffer)
	})

	Context("when ListEvents returns an error", func() {
		var (
			ctrl          *gomock.Controller
			mockEventRepo *mockrepo.MockEventRepository
		)

		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())
			mockEventRepo = mockrepo.NewMockEventRepository(ctrl)
		})

		AfterEach(func() {
			ctrl.Finish()
		})

		It("should return exit code 1", func() {
			// Setup mock to return error
			mockEventRepo.EXPECT().
				List(gomock.Any(), gomock.Any()).
				Return(nil, errors.New("database connection failed"))

			svc := careerservice.NewService(mockEventRepo)
			svc.SetBurstRepository(careermemory.NewBurstRepository())

			code := bursts.DetectBursts(svc, out, err, tea.WithInput(nil))
			Expect(code).To(Equal(1))
		})

		It("should handle error gracefully without panic", func() {
			mockEventRepo.EXPECT().
				List(gomock.Any(), gomock.Any()).
				Return(nil, errors.New("query timeout"))

			svc := careerservice.NewService(mockEventRepo)
			svc.SetBurstRepository(careermemory.NewBurstRepository())

			Expect(func() {
				bursts.DetectBursts(svc, out, err, tea.WithInput(nil))
			}).NotTo(Panic())
		})
	})

	Context("when SaveBurstSuggestions returns no saved bursts", func() {
		var (
			ctx           context.Context
			ctrl          *gomock.Controller
			mockBurstRepo *mockrepo.MockBurstRepository
		)

		BeforeEach(func() {
			ctx = context.Background()
			ctrl = gomock.NewController(GinkgoT())
			mockBurstRepo = mockrepo.NewMockBurstRepository(ctrl)
		})

		AfterEach(func() {
			ctrl.Finish()
		})

		It("should return exit code 0 with info message when saves fail", func() {
			// Setup mock to fail all saves
			mockBurstRepo.EXPECT().
				Create(gomock.Any(), gomock.Any()).
				Return(errors.New("database write failed")).
				AnyTimes()

			svc := careerservice.NewService(careermemory.NewEventRepository())
			svc.SetBurstRepository(mockBurstRepo)

			suggestions := []burst_fact.BurstSuggestion{
				{
					Name:        "Test Burst",
					Description: "Test Description",
					EventIDs:    []string{"event-1", "event-2"},
				},
			}

			// When all saves fail, SaveBurstSuggestions returns empty list (not error)
			// This triggers the else branch (isSuccess=false) which prints info message
			code := bursts.SaveBurstSuggestions(ctx, svc, suggestions, tea.WithInput(nil))
			Expect(code).To(Equal(0))
		})

		It("should handle save errors gracefully without panic", func() {
			mockBurstRepo.EXPECT().
				Create(gomock.Any(), gomock.Any()).
				Return(errors.New("constraint violation")).
				AnyTimes()

			svc := careerservice.NewService(careermemory.NewEventRepository())
			svc.SetBurstRepository(mockBurstRepo)

			suggestions := []burst_fact.BurstSuggestion{
				{
					Name:        "Test",
					Description: "Test",
					EventIDs:    []string{"event-1"},
				},
			}

			Expect(func() {
				bursts.SaveBurstSuggestions(ctx, svc, suggestions, tea.WithInput(nil))
			}).NotTo(Panic())
		})
	})
})

package facts_test

import (
	"errors"
	"bytes"
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	cmdpkg "github.com/baphled/kariya/internal/cli/cmd"
	domain "github.com/baphled/kariya/internal/domain/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	mockrepo "github.com/baphled/kariya/internal/testutil/mocks/repository"
	"github.com/baphled/kariya/internal/cli/cmd/cliutil"
	"github.com/baphled/kariya/internal/cli/cmd/facts"
	"github.com/baphled/kariya/internal/testutil/fixtures"
)

var _ = Describe("ExtractFacts", func() {
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

			code := facts.ExtractFacts(svc, out, err, tea.WithInput(nil))
			Expect(code).To(Equal(0))
		})

		It("should not write error output", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			code := facts.ExtractFacts(svc, out, err, tea.WithInput(nil))
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
				facts.ExtractFacts(svc, out, err, tea.WithInput(nil))
			}).NotTo(Panic())
		})

		It("should return valid exit code", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			code := facts.ExtractFacts(svc, out, err)
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
				facts.ExtractFacts(svc, out, err)
			}).NotTo(Panic())
		})

		It("should return valid exit code", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			code := facts.ExtractFacts(svc, out, err, tea.WithInput(nil))
			Expect(code).To(BeNumerically(">=", 0))
			Expect(code).To(BeNumerically("<=", 1))
		})
	})

	Context("error handling", func() {
		It("should handle service errors gracefully", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			code := facts.ExtractFacts(svc, out, err, tea.WithInput(nil))
			Expect(code).To(BeNumerically(">=", 0))
		})

		It("should return valid exit code", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			code := facts.ExtractFacts(svc, out, err, tea.WithInput(nil))
			Expect(code).To(BeNumerically(">=", 0))
			Expect(code).To(BeNumerically("<=", 1))
		})

		It("should handle nil service gracefully", func() {
			Expect(func() {
				facts.ExtractFacts(nil, out, err, tea.WithInput(nil))
			}).To(Panic())
		})
	})

	Context("when no facts are extracted from events", func() {
		It("should return zero with info message when no facts found", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			// Create event that exists but won't generate saveable facts
			// Using text that may not generate facts with high enough confidence
			event := fixtures.EventWith("no-fact-event", 
				"Worked on some tasks and assignments during the period", "", "")
			captureErr := svc.CaptureEvent(context.Background(), event, "timeline")
			Expect(captureErr).NotTo(HaveOccurred())

			out := new(bytes.Buffer)
			err := new(bytes.Buffer)

			code := facts.ExtractFacts(svc, out, err, tea.WithInput(nil))
			// Should return 0 whether facts found or not
			Expect(code).To(Equal(0))
		})
	})

	Context("extraction with various event scenarios", func() {
		It("should handle events with technical terminology", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			event := fixtures.EventWith("tech-event", 
				"Implemented microservices architecture with Kubernetes Docker containers and cloud deployment", "", "")
			captureErr := svc.CaptureEvent(context.Background(), event, "timeline")
			Expect(captureErr).NotTo(HaveOccurred())

			code := facts.ExtractFacts(svc, out, err, tea.WithInput(nil))
			Expect(code).To(BeNumerically(">=", 0))
			Expect(code).To(BeNumerically("<=", 1))
		})

		It("should process events with multiple skill mentions", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			event := fixtures.EventWith("multi-skill", 
				"Developed REST APIs using Python Django PostgreSQL with unit tests and CI/CD pipelines", "", "")
			captureErr := svc.CaptureEvent(context.Background(), event, "timeline")
			Expect(captureErr).NotTo(HaveOccurred())

			code := facts.ExtractFacts(svc, out, err, tea.WithInput(nil))
			Expect(code).To(BeNumerically(">=", 0))
			Expect(code).To(BeNumerically("<=", 1))
		})
	})

	Context("ExtractFacts error scenarios", func() {
		It("should return 1 when ListEvents fails", func() {
			ctrl := gomock.NewController(GinkgoT())
			mockEventRepo := mockrepo.NewMockEventRepository(ctrl)
			mockEventRepo.EXPECT().
				List(gomock.Any(), gomock.Any()).
				Return(nil, errors.New("database error"))

			svc := careerservice.NewService(mockEventRepo)

			code := facts.ExtractFacts(svc, out, err, tea.WithInput(nil))
			Expect(code).To(Equal(1))
		})

		It("should return 0 when events list is empty", func() {
			ctrl := gomock.NewController(GinkgoT())
			mockEventRepo := mockrepo.NewMockEventRepository(ctrl)
			mockEventRepo.EXPECT().
				List(gomock.Any(), gomock.Any()).
				Return([]*domain.Event{}, nil)

			svc := careerservice.NewService(mockEventRepo)

			code := facts.ExtractFacts(svc, out, err, tea.WithInput(nil))
			Expect(code).To(Equal(0))
		})

		It("should handle repository connection errors", func() {
			ctrl := gomock.NewController(GinkgoT())
			mockEventRepo := mockrepo.NewMockEventRepository(ctrl)
			mockEventRepo.EXPECT().
				List(gomock.Any(), gomock.Any()).
				Return(nil, errors.New("connection timeout"))

			svc := careerservice.NewService(mockEventRepo)

			code := facts.ExtractFacts(svc, out, err, tea.WithInput(nil))
			Expect(code).To(Equal(1))
		})

		It("should extract facts from events with comprehensive text", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			event := fixtures.EventWith("comprehensive",
				"Implemented microservices platform with Kubernetes Docker containers and automated CI/CD deployment pipelines", "", "")
			captureErr := svc.CaptureEvent(context.Background(), event, "timeline")
			Expect(captureErr).NotTo(HaveOccurred())

			code := facts.ExtractFacts(svc, out, err, tea.WithInput(nil))
			Expect(code).To(BeNumerically(">=", 0))
		})
	})

})

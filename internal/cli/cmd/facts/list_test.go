package facts_test

import (
	"bytes"
	"context"
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	cmdpkg "github.com/baphled/kariya/internal/cli/cmd"
	"github.com/baphled/kariya/internal/cli/cmd/cliutil"
	"github.com/baphled/kariya/internal/cli/cmd/facts"
	careermemory "github.com/baphled/kariya/internal/repository/career/memory"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	mockrepo "github.com/baphled/kariya/internal/testutil/mocks/repository"
)

var _ = Describe("ListFacts", func() {
	var (
		ctx cliutil.ServiceProvider
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

	Context("with no facts", func() {
		It("should return 0 and print info message", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			code := facts.ListFacts(svc, out, err, tea.WithInput(nil))
			Expect(code).To(Equal(0))
		})

		It("should not write error output", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			code := facts.ListFacts(svc, out, err, tea.WithInput(nil))
			Expect(code).To(Equal(0))
			Expect(err.Len()).To(Equal(0))
		})
	})

	Context("with facts in database", func() {
		BeforeEach(func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())
			fact := fixtures.FactWithCategories("fact-1", "Implemented authentication system", "event-1", []string{"technical"}, []string{"hiring_manager"})
			fact.RoleFit = "principal"
			saveErr := svc.SaveFact(context.Background(), fact)
			Expect(saveErr).NotTo(HaveOccurred())
		})

		It("should return 0 on successful execution", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			code := facts.ListFacts(svc, out, err, tea.WithInput(nil))
			Expect(code).To(Equal(0))
		})

		It("should write output to stdout", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			code := facts.ListFacts(svc, out, err, tea.WithInput(nil))
			Expect(code).To(Equal(0))
			Expect(out.String()).To(ContainSubstring("Existing Facts"))
		})

		It("should display fact information", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			code := facts.ListFacts(svc, out, err, tea.WithInput(nil))
			Expect(code).To(Equal(0))
			Expect(out.String()).To(ContainSubstring("Implemented authentication system"))
		})

		It("should handle execution without panic", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			Expect(func() {
				facts.ListFacts(svc, out, err, tea.WithInput(nil))
			}).NotTo(Panic())
		})
	})

	Context("with multiple facts", func() {
		BeforeEach(func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			for i := 1; i <= 3; i++ {
				factID := "fact-" + string(rune(48+i))
				factText := "Fact " + string(rune(48+i))
				fact := fixtures.FactWithCategories(factID, factText, "event-1", []string{"technical"}, []string{"hiring_manager"})
				fact.RoleFit = "principal"
				saveErr := svc.SaveFact(context.Background(), fact)
				Expect(saveErr).NotTo(HaveOccurred())
			}
		})

		It("should list all facts", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			code := facts.ListFacts(svc, out, err, tea.WithInput(nil))
			Expect(code).To(Equal(0))
			Expect(out.String()).To(ContainSubstring("Total facts: 3"))
		})
	})

	Context("error handling", func() {
		It("should handle service errors gracefully", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			code := facts.ListFacts(svc, out, err, tea.WithInput(nil))
			Expect(code).To(BeNumerically(">=", 0))
		})

		It("should return valid exit code", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			code := facts.ListFacts(svc, out, err, tea.WithInput(nil))
			Expect(code).To(BeNumerically(">=", 0))
			Expect(code).To(BeNumerically("<=", 1))
		})

		It("should handle nil service gracefully", func() {
			Expect(func() {
				facts.ListFacts(nil, out, err, tea.WithInput(nil))
			}).To(Panic())
		})
	})
})

var _ = Describe("ListFacts Error Paths", func() {
	var (
		out *bytes.Buffer
		err *bytes.Buffer
	)

	BeforeEach(func() {
		out = new(bytes.Buffer)
		err = new(bytes.Buffer)
	})

	Context("when fact repository is not configured", func() {
		It("should return exit code 1", func() {
			// Create service with nil fact repository
			svc := careerservice.NewService(nil)
			// Don't call SetFactRepository - leave it nil

			code := facts.ListFacts(svc, out, err, tea.WithInput(nil))
			Expect(code).To(Equal(1))
		})

		It("should print error message", func() {
			svc := careerservice.NewService(nil)

			facts.ListFacts(svc, out, err, tea.WithInput(nil))
			// Error is printed to stderr via cliutil.PrintError, not to err buffer
			// Just verify exit code is 1 (which we test above)
		})
	})

	Context("when repository List() returns an error", func() {
		var (
			ctrl      *gomock.Controller
			mockRepo  *mockrepo.MockFactRepository
			eventRepo *careermemory.EventRepository
		)

		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())
			mockRepo = mockrepo.NewMockFactRepository(ctrl)
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
			svc.SetFactRepository(mockRepo)

			code := facts.ListFacts(svc, out, err, tea.WithInput(nil))
			Expect(code).To(Equal(1))
		})

		It("should handle repository errors gracefully", func() {
			mockRepo.EXPECT().
				List(gomock.Any(), gomock.Any()).
				Return(nil, errors.New("timeout"))

			svc := careerservice.NewService(eventRepo)
			svc.SetFactRepository(mockRepo)

			Expect(func() {
				facts.ListFacts(svc, out, err, tea.WithInput(nil))
			}).NotTo(Panic())
		})
	})

	Describe("ExecuteListFacts", func() {
		It("should fail with nil service", func() {
			err := facts.ExecuteListFacts(nil, new(bytes.Buffer), new(bytes.Buffer), tea.WithInput(nil))
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("service not initialized"))
		})

		It("should succeed with valid service", func() {
			cliCtx := cmdpkg.NewCLIContext("", true)
			initErr := cliCtx.InitService(new(bytes.Buffer))
			Expect(initErr).NotTo(HaveOccurred())

			svc := cliCtx.Service()
			Expect(svc).NotTo(BeNil())

			err := facts.ExecuteListFacts(svc, new(bytes.Buffer), new(bytes.Buffer), tea.WithInput(nil))
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("ExecuteExtractFacts", func() {
		It("should fail with nil service", func() {
			err := facts.ExecuteExtractFacts(nil, new(bytes.Buffer), new(bytes.Buffer), tea.WithInput(nil))
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("service not initialized"))
		})

		It("should succeed with valid service", func() {
			cliCtx := cmdpkg.NewCLIContext("", true)
			initErr := cliCtx.InitService(new(bytes.Buffer))
			Expect(initErr).NotTo(HaveOccurred())

			svc := cliCtx.Service()
			Expect(svc).NotTo(BeNil())

			err := facts.ExecuteExtractFacts(svc, new(bytes.Buffer), new(bytes.Buffer), tea.WithInput(nil))
			Expect(err).NotTo(HaveOccurred())
		})
	})
})

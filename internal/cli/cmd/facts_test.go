package cmd_test

import (
	"bytes"
	"context"
	"errors"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	cmdpkg "github.com/baphled/kariya/internal/cli/cmd"
	"github.com/baphled/kariya/internal/cli/cmd/cliutil"
	"github.com/baphled/kariya/internal/domain/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	mockrepo "github.com/baphled/kariya/internal/testutil/mocks/repository"
)

var _ = Describe("Facts Commands", func() {
	var (
		ctx    *cmdpkg.CLIContext
		out    *bytes.Buffer
		errOut *bytes.Buffer
		runner cliutil.MockProgressRunner
	)

	BeforeEach(func() {
		ctx = cmdpkg.NewCLIContext("", true)
		initErr := ctx.InitService(errOut)
		Expect(initErr).NotTo(HaveOccurred())

		out = new(bytes.Buffer)
		errOut = new(bytes.Buffer)
		runner = cliutil.MockProgressRunner{}
	})

	Describe("ExtractFacts", func() {
		It("returns 0 when no events exist", func() {
			code := cmdpkg.ExtractFacts(ctx.Service(), out, errOut, runner)
			Expect(code).To(Equal(0))
			Expect(out.String()).To(ContainSubstring("No facts extracted"))
		})

		It("returns 1 when ExtractFactsFromAllEvents fails", func() {
			ctrl := gomock.NewController(GinkgoT())
			defer ctrl.Finish()

			mockEventRepo := mockrepo.NewMockEventRepository(ctrl)
			mockEventRepo.EXPECT().
				List(gomock.Any(), gomock.Any()).
				Return(nil, errors.New("database error")).
				Times(1)

			mockService := careerservice.NewService(mockEventRepo)
			code := cmdpkg.ExtractFacts(mockService, out, errOut, runner)
			Expect(code).To(Equal(1))
			Expect(errOut.String()).To(ContainSubstring("Error during fact extraction"))
		})

		It("returns 0 when events exist and facts are extracted successfully", func() {
			event := fixtures.EventFactory.MustCreateWithOption(map[string]interface{}{
				"ID":   "event-1",
				"Date": time.Now(),
			}).(*career.Event)
			captureErr := ctx.Service().CaptureEvent(context.Background(), event, "timeline")
			Expect(captureErr).NotTo(HaveOccurred())

			code := cmdpkg.ExtractFacts(ctx.Service(), out, errOut, runner)
			Expect(code).To(Equal(0))
			output := out.String()
			Expect(output).To(ContainSubstring("Fact Extraction Results"))
			Expect(output).To(ContainSubstring("Extracted"))
			Expect(output).To(ContainSubstring("Fact extraction complete"))
		})
	})

	Describe("ListFacts", func() {
		It("returns 0 when no facts exist", func() {
			code := cmdpkg.ListFacts(ctx.Service(), out, errOut)
			Expect(code).To(Equal(0))
			Expect(out.String()).To(ContainSubstring("No facts found"))
		})

		It("returns 1 when fact repository is not configured", func() {
			ctx.Service().SetFactRepository(nil)

			code := cmdpkg.ListFacts(ctx.Service(), out, errOut)
			Expect(code).To(Equal(1))
			Expect(errOut.String()).To(ContainSubstring("Fact repository not configured"))
		})

		It("displays all fact fields when facts exist", func() {
			ctrl := gomock.NewController(GinkgoT())
			defer ctrl.Finish()

			now := time.Now()
			mockFactRepo := mockrepo.NewMockFactRepository(ctrl)
			mockFactRepo.EXPECT().
				List(gomock.Any(), gomock.Any()).
				Return([]*career.Fact{
					{
						ID:                   "fact-1",
						Text:                 "Led API migration project",
						SourceEventID:        "event-1",
						CompetencyCategories: []string{"technical", "leadership"},
						RoleFit:              "Senior Engineer",
						AudienceRelevance:    []string{"hiring-manager", "recruiter"},
						CreatedAt:            now,
					},
				}, nil).
				Times(1)

			ctx.Service().SetFactRepository(mockFactRepo)

			code := cmdpkg.ListFacts(ctx.Service(), out, errOut)
			Expect(code).To(Equal(0))
			output := out.String()
			Expect(output).To(ContainSubstring("=== Existing Facts ==="))
			Expect(output).To(ContainSubstring("Total facts: 1"))
			Expect(output).To(ContainSubstring("Led API migration project"))
			Expect(output).To(ContainSubstring("fact-1"))
			Expect(output).To(ContainSubstring("event-1"))
			Expect(output).To(ContainSubstring("technical, leadership"))
			Expect(output).To(ContainSubstring("Senior Engineer"))
			Expect(output).To(ContainSubstring("hiring-manager, recruiter"))
		})
	})
})

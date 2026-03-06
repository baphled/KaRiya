package cmd_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
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

var _ = Describe("Bursts Commands", func() {
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

	Describe("DetectBursts", func() {
		Context("when DetectAndSaveBursts returns an error", func() {
			It("returns 1 and displays error message", func() {
				ctrl := gomock.NewController(GinkgoT())
				defer ctrl.Finish()

				mockEventRepo := mockrepo.NewMockEventRepository(ctrl)
				mockEventRepo.EXPECT().
					List(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("database connection failed")).
					Times(1)

				mockService := careerservice.NewService(mockEventRepo)
				mockService.SetFactRepository(ctx.Service().GetFactRepository())
				mockService.SetBurstRepository(ctx.Service().GetBurstRepository())

				code := cmdpkg.DetectBursts(mockService, out, errOut, runner)
				Expect(code).To(Equal(1))
				Expect(errOut.String()).To(ContainSubstring("Error detecting bursts"))
			})
		})

		Context("when no events exist", func() {
			It("returns 0 and displays no events message", func() {
				code := cmdpkg.DetectBursts(ctx.Service(), out, errOut, runner)
				Expect(code).To(Equal(0))
				Expect(out.String()).To(ContainSubstring("No bursts detected"))
			})
		})

		Context("when bursts are detected and saved successfully", func() {
			It("returns 0 and displays burst detection results", func() {
				// Create events with similar content to trigger burst detection
				baseDate := time.Now()
				similarText := "Led API redesign for performance improvement"
				for i := range 5 {
					event := fixtures.EventFactory.MustCreateWithOption(map[string]interface{}{
						"ID":      fmt.Sprintf("event-%d", i),
						"Date":    baseDate.AddDate(0, 0, -i),
						"Text":    similarText,
						"Company": "TechCorp",
						"Project": "API Modernisation",
						"Tags":    []string{"technical", "achievement", "leadership"},
					}).(*career.Event)
					captureErr := ctx.Service().CaptureEvent(context.Background(), event, "timeline")
					Expect(captureErr).NotTo(HaveOccurred())
				}

				code := cmdpkg.DetectBursts(ctx.Service(), out, errOut, runner)
				Expect(code).To(Equal(0))
				output := out.String()
				Expect(output).To(ContainSubstring("Detected"))
				Expect(output).To(ContainSubstring("Burst detection complete"))
			})
		})
	})

	Describe("ListBursts", func() {
		Context("when no bursts exist", func() {
			It("returns 0 and displays no bursts message", func() {
				code := cmdpkg.ListBursts(ctx.Service(), out, errOut)
				Expect(code).To(Equal(0))
				Expect(out.String()).To(ContainSubstring("No bursts found"))
			})
		})

		Context("when burst repository is nil", func() {
			It("returns 1 and displays error message", func() {
				ctx.Service().SetBurstRepository(nil)

				code := cmdpkg.ListBursts(ctx.Service(), out, errOut)
				Expect(code).To(Equal(1))
				Expect(errOut.String()).To(ContainSubstring("Burst repository not configured"))
			})
		})

		Context("when bursts exist", func() {
			It("returns 0 and displays burst header", func() {
				ctrl := gomock.NewController(GinkgoT())
				defer ctrl.Finish()

				mockBurstRepo := mockrepo.NewMockBurstRepository(ctrl)
				bursts := []*career.Burst{fixtures.Burst("burst-1", "event-1", "event-2")}
				mockBurstRepo.EXPECT().
					List(gomock.Any(), gomock.Any()).
					Return(bursts, nil).
					Times(1)

				ctx.Service().SetBurstRepository(mockBurstRepo)

				code := cmdpkg.ListBursts(ctx.Service(), out, errOut)
				Expect(code).To(Equal(0))
				output := out.String()
				Expect(output).To(ContainSubstring("=== Bursts ==="))
				Expect(output).To(ContainSubstring("Total bursts"))
			})
		})
	})
})

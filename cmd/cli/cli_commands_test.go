package main

import (
	"bytes"
	"context"
	"time"

	careermemory "github.com/baphled/kariya/internal/repository/career/memory"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CLI Command Handlers", func() {
	var (
		svc         *careerservice.Service
		repo        *careermemory.EventRepository
		burstRepo   *careermemory.BurstRepository
		factRepo    *careermemory.FactRepository
		ctx         context.Context
		out, errOut bytes.Buffer
	)

	BeforeEach(func() {
		ctx = context.Background()
		repo = careermemory.NewEventRepository()
		burstRepo = careermemory.NewBurstRepository()
		factRepo = careermemory.NewFactRepository()
		svc = careerservice.NewService(repo)
		svc.SetBurstRepository(burstRepo)
		svc.SetFactRepository(factRepo)
		out.Reset()
		errOut.Reset()
	})

	Describe("handleDetectBursts", func() {
		Context("when no events exist", func() {
			It("should return success with informative message", func() {
				exitCode := handleDetectBursts(svc, &out, &errOut)
				Expect(exitCode).To(Equal(0))
				Expect(out.String()).To(ContainSubstring("No events found"))
			})
		})

		Context("when events exist", func() {
			BeforeEach(func() {
				// Create clustered events that will trigger burst detection
				baseDate := time.Now().AddDate(0, 0, -30)
				for i := 0; i < 5; i++ {
					evt := fixtures.EventWith("evt"+string(rune('1'+i)), "Cloud migration task", "TechCorp", "CloudMigration")
					evt.Date = baseDate.AddDate(0, 0, i)
					evt.Tags = []string{"technical", "project"}
					err := repo.Create(ctx, evt)
					Expect(err).ToNot(HaveOccurred())
				}
			})

			It("should detect and display burst information", func() {
				exitCode := handleDetectBursts(svc, &out, &errOut)
				Expect(exitCode).To(Equal(0))
				output := out.String()
				Expect(output).To(ContainSubstring("Detecting bursts"))
				Expect(output).To(ContainSubstring("Burst Detection Results"))
			})

			It("should show confidence scores", func() {
				handleDetectBursts(svc, &out, &errOut)
				Expect(out.String()).To(ContainSubstring("Confidence:"))
			})

			It("should save burst suggestions", func() {
				handleDetectBursts(svc, &out, &errOut)
				output := out.String()
				Expect(output).To(Or(
					ContainSubstring("Saved"),
					ContainSubstring("complete"),
				))
			})
		})
	})

	Describe("handleExtractFacts", func() {
		Context("when no events exist", func() {
			It("should return success with informative message", func() {
				exitCode := handleExtractFacts(svc, &out, &errOut)
				Expect(exitCode).To(Equal(0))
				Expect(out.String()).To(ContainSubstring("No events found"))
			})
		})

		Context("when events exist", func() {
			BeforeEach(func() {
				// Create events with factual content
				evt := fixtures.EventWith("e1", "Led a team of 5 engineers to deliver microservices architecture", "TechCorp", "")
				evt.Date = time.Now().AddDate(0, 0, -10)
				evt.Tags = []string{"leadership", "technical"}
				err := repo.Create(ctx, evt)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should extract facts from events", func() {
				exitCode := handleExtractFacts(svc, &out, &errOut)
				Expect(exitCode).To(Equal(0))
				Expect(out.String()).To(ContainSubstring("Extracting facts"))
			})

			It("should show extraction results", func() {
				handleExtractFacts(svc, &out, &errOut)
				output := out.String()
				Expect(output).To(Or(
					ContainSubstring("facts"),
					ContainSubstring("Extracted"),
				))
			})

			It("should show completion message", func() {
				handleExtractFacts(svc, &out, &errOut)
				Expect(out.String()).To(ContainSubstring("complete"))
			})
		})
	})

	Describe("handleShowBursts", func() {
		Context("when burst repository is not configured", func() {
			It("should return error", func() {
				svcWithoutBursts := careerservice.NewService(repo)
				exitCode := handleShowBursts(svcWithoutBursts, &out, &errOut)
				Expect(exitCode).To(Equal(1))
				Expect(errOut.String()).To(ContainSubstring("Burst repository not configured"))
			})
		})

		Context("when no bursts exist", func() {
			It("should show no bursts message", func() {
				exitCode := handleShowBursts(svc, &out, &errOut)
				Expect(exitCode).To(Equal(0))
				Expect(out.String()).To(ContainSubstring("No bursts found"))
			})
		})

		Context("when bursts exist", func() {
			BeforeEach(func() {
				// Create sample burst
				b := fixtures.Burst("b1", "e1", "e2", "e3")
				b.Name = "Cloud Migration Sprint"
				b.Description = "Complete cloud infrastructure migration"
				err := burstRepo.Create(ctx, b)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should display all bursts", func() {
				exitCode := handleShowBursts(svc, &out, &errOut)
				Expect(exitCode).To(Equal(0))
				Expect(out.String()).To(ContainSubstring("Bursts"))
			})

			It("should show burst count", func() {
				handleShowBursts(svc, &out, &errOut)
				output := out.String()
				Expect(output).To(Or(
					ContainSubstring("Total bursts"),
					ContainSubstring("1"),
				))
			})

			It("should show burst names", func() {
				handleShowBursts(svc, &out, &errOut)
				Expect(out.String()).To(ContainSubstring("Cloud Migration Sprint"))
			})

			It("should show burst IDs", func() {
				handleShowBursts(svc, &out, &errOut)
				Expect(out.String()).To(ContainSubstring("b1"))
			})

			It("should show event counts", func() {
				handleShowBursts(svc, &out, &errOut)
				Expect(out.String()).To(ContainSubstring("Events:"))
			})

			It("should show descriptions", func() {
				handleShowBursts(svc, &out, &errOut)
				Expect(out.String()).To(ContainSubstring("Complete cloud infrastructure migration"))
			})
		})
	})

	Describe("handleShowFacts", func() {
		Context("when fact repository is not configured", func() {
			It("should return error", func() {
				svcWithoutFacts := careerservice.NewService(repo)
				exitCode := handleShowFacts(svcWithoutFacts, &out, &errOut)
				Expect(exitCode).To(Equal(1))
				Expect(errOut.String()).To(ContainSubstring("Fact repository not configured"))
			})
		})

		Context("when no facts exist", func() {
			It("should show no facts message", func() {
				exitCode := handleShowFacts(svc, &out, &errOut)
				Expect(exitCode).To(Equal(0))
				Expect(out.String()).To(ContainSubstring("No facts found"))
			})
		})

		Context("when facts exist", func() {
			BeforeEach(func() {
				// Create sample fact
				f := fixtures.FactWithCategories("f1", "Led team of 5 engineers in microservices migration", "e1", []string{"leadership", "technical"}, []string{"peer", "hiring_manager"})
				err := factRepo.Create(ctx, f)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should display all facts", func() {
				exitCode := handleShowFacts(svc, &out, &errOut)
				Expect(exitCode).To(Equal(0))
				Expect(out.String()).To(ContainSubstring("Facts"))
			})

			It("should show fact count", func() {
				handleShowFacts(svc, &out, &errOut)
				output := out.String()
				Expect(output).To(Or(
					ContainSubstring("Total facts"),
					ContainSubstring("1"),
				))
			})

			It("should show fact text", func() {
				handleShowFacts(svc, &out, &errOut)
				Expect(out.String()).To(ContainSubstring("Led team of 5 engineers"))
			})

			It("should show fact IDs", func() {
				handleShowFacts(svc, &out, &errOut)
				Expect(out.String()).To(ContainSubstring("f1"))
			})

			It("should show source event IDs", func() {
				handleShowFacts(svc, &out, &errOut)
				Expect(out.String()).To(ContainSubstring("Source Event"))
			})

			It("should show competency categories", func() {
				handleShowFacts(svc, &out, &errOut)
				output := out.String()
				Expect(output).To(ContainSubstring("Competencies"))
			})

			It("should show role fit", func() {
				handleShowFacts(svc, &out, &errOut)
				Expect(out.String()).To(ContainSubstring("Role Fit"))
			})

			It("should show target audiences", func() {
				handleShowFacts(svc, &out, &errOut)
				Expect(out.String()).To(ContainSubstring("Target Audiences"))
			})
		})
	})
})

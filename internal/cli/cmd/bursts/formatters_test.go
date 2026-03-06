package bursts_test

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/cmd/bursts"
	domain "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/burstfact"
	"github.com/baphled/kariya/internal/testutil/fixtures"
)

var _ = Describe("Bursts Formatters", func() {
	Describe("FormatBurstDetectionResults", func() {
		Context("when no bursts are detected", func() {
			It("should return header with zero count", func() {
				output := bursts.FormatBurstDetectionResults([]burstfact.BurstSuggestion{}, 10)
				Expect(output).To(ContainSubstring("Burst Detection Results"))
				Expect(output).To(ContainSubstring("Detected 0 bursts from 10 events"))
			})
		})

		Context("when bursts are detected", func() {
			It("should format all bursts with details", func() {
				suggestions := []burstfact.BurstSuggestion{
					{
						Name:            "Q1 Initiative",
						EventIDs:        []string{"event-1", "event-2"},
						ConfidenceScore: 0.85,
						Description:     "First quarter project",
					},
					{
						Name:            "Q2 Migration",
						EventIDs:        []string{"event-3", "event-4", "event-5"},
						ConfidenceScore: 0.92,
					},
				}

				output := bursts.FormatBurstDetectionResults(suggestions, 20)
				Expect(output).To(ContainSubstring("Detected 2 bursts from 20 events"))
				Expect(output).To(ContainSubstring("1. Q1 Initiative"))
				Expect(output).To(ContainSubstring("2. Q2 Migration"))
				Expect(output).To(ContainSubstring("Events: 2"))
				Expect(output).To(ContainSubstring("Events: 3"))
				Expect(output).To(ContainSubstring("Confidence: 85.0%"))
				Expect(output).To(ContainSubstring("Confidence: 92.0%"))
			})
		})

		Context("when bursts have names", func() {
			It("should display names", func() {
				suggestions := []burstfact.BurstSuggestion{
					{
						Name:            "Named Burst",
						EventIDs:        []string{"event-1", "event-2"},
						ConfidenceScore: 0.75,
					},
				}

				output := bursts.FormatBurstDetectionResults(suggestions, 5)
				Expect(output).To(ContainSubstring("1. Named Burst"))
			})
		})

		Context("when bursts have no names", func() {
			It("should use default burst numbering", func() {
				suggestions := []burstfact.BurstSuggestion{
					{
						Name:            "",
						EventIDs:        []string{"event-1", "event-2"},
						ConfidenceScore: 0.80,
					},
				}

				output := bursts.FormatBurstDetectionResults(suggestions, 5)
				Expect(output).To(ContainSubstring("1. Burst 1"))
			})
		})

		Context("when bursts have descriptions", func() {
			It("should display descriptions", func() {
				suggestions := []burstfact.BurstSuggestion{
					{
						Name:            "Project",
						EventIDs:        []string{"event-1", "event-2"},
						ConfidenceScore: 0.90,
						Description:     "Migration effort",
					},
				}

				output := bursts.FormatBurstDetectionResults(suggestions, 5)
				Expect(output).To(ContainSubstring("Description: Migration effort"))
			})
		})
	})

	Describe("FormatSaveResults", func() {
		Context("when bursts were saved", func() {
			It("should return success message", func() {
				savedBursts := []*domain.Burst{
					fixtures.Burst("1", "event-1", "event-2"),
					fixtures.Burst("2", "event-3", "event-4"),
					fixtures.Burst("3", "event-5", "event-6"),
				}

				message, isSuccess := bursts.FormatSaveResults(savedBursts, 5)
				Expect(isSuccess).To(BeTrue())
				Expect(message).To(Equal("Burst detection complete! Saved 3 of 5 bursts to database."))
			})
		})

		Context("when no bursts were saved", func() {
			It("should return info message", func() {
				savedBursts := []*domain.Burst{}

				message, isSuccess := bursts.FormatSaveResults(savedBursts, 3)
				Expect(isSuccess).To(BeFalse())
				Expect(message).To(ContainSubstring("no bursts were saved"))
				Expect(message).To(ContainSubstring("repository may not be configured"))
			})
		})

		Context("when all bursts were saved", func() {
			It("should show correct counts", func() {
				savedBursts := []*domain.Burst{
					fixtures.Burst("1", "event-1", "event-2"),
					fixtures.Burst("2", "event-3", "event-4"),
				}

				message, isSuccess := bursts.FormatSaveResults(savedBursts, 2)
				Expect(isSuccess).To(BeTrue())
				Expect(message).To(Equal("Burst detection complete! Saved 2 of 2 bursts to database."))
			})
		})
	})

	Describe("FormatBurstList", func() {
		Context("with empty list", func() {
			It("should show zero count", func() {
				output := bursts.FormatBurstList([]*domain.Burst{})
				Expect(output).To(ContainSubstring("Total bursts: 0"))
			})
		})

		Context("with bursts having names", func() {
			It("should format with numbered list", func() {
				burstList := []*domain.Burst{
					fixtures.Burst("burst-1", "event-1", "event-2"),
					fixtures.Burst("burst-2", "event-3", "event-4"),
				}
				burstList[0].Name = "Q1 Migration"
				burstList[1].Name = "Q2 Refactoring"

				output := bursts.FormatBurstList(burstList)
				Expect(output).To(ContainSubstring("Total bursts: 2"))
				Expect(output).To(ContainSubstring("1. Q1 Migration"))
				Expect(output).To(ContainSubstring("2. Q2 Refactoring"))
				Expect(output).To(ContainSubstring("Events: 2"))
			})
		})

		Context("with unnamed bursts", func() {
			It("should use default numbering", func() {
				burst := fixtures.Burst("burst-1", "event-1", "event-2")
				burst.Name = ""
				burstList := []*domain.Burst{burst}

				output := bursts.FormatBurstList(burstList)
				Expect(output).To(ContainSubstring("1. Burst 1"))
			})
		})

		Context("with bursts having descriptions", func() {
			It("should show descriptions", func() {
				burst := fixtures.Burst("burst-1", "event-1", "event-2")
				burst.Name = "Project"
				burst.Description = "Major migration effort"
				burstList := []*domain.Burst{burst}

				output := bursts.FormatBurstList(burstList)
				Expect(output).To(ContainSubstring("Description: Major migration effort"))
			})
		})

		Context("with bursts without descriptions", func() {
			It("should not show description line", func() {
				burst := fixtures.Burst("burst-1", "event-1", "event-2")
				burst.Name = "Project"
				burst.Description = ""
				burstList := []*domain.Burst{burst}

				output := bursts.FormatBurstList(burstList)
				lines := strings.Split(output, "\n")
				descriptionFound := false
				for _, line := range lines {
					if strings.Contains(line, "Description:") {
						descriptionFound = true
						break
					}
				}
				Expect(descriptionFound).To(BeFalse())
			})
		})
	})
})

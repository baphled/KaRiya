package models

import (
	"strings"

	"time"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("FactCard", func() {
	var (
		fact *career.Fact
		card *FactCard
	)

	ginkgo.BeforeEach(func() {
		fact = &career.Fact{
			ID:                   "fact-1",
			Text:                 "Led cross-functional team to deliver critical platform migration",
			CompetencyCategories: []string{"leadership", "technical"},
			RoleFit:              career.RoleFitPrincipal,
			AudienceRelevance:    []string{"hiring_manager", "recruiter"},
			StrengthSignal:       "leadership",
			SourceEventID:        "event-123",
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		}
		card = NewFactCard(fact)
		card.SetWidth(80)
	})

	ginkgo.Describe("NewFactCard", func() {
		ginkgo.It("should create a new fact card", func() {
			c := NewFactCard(fact)
			gomega.Expect(c).NotTo(gomega.BeNil())
			gomega.Expect(c.fact).To(gomega.Equal(fact))
			gomega.Expect(c.width).To(gomega.Equal(80))
		})

		ginkgo.It("should handle nil fact", func() {
			c := NewFactCard(nil)
			gomega.Expect(c).NotTo(gomega.BeNil())
			gomega.Expect(c.fact).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("SetWidth", func() {
		ginkgo.It("should set the width", func() {
			card.SetWidth(100)
			gomega.Expect(card.width).To(gomega.Equal(100))
		})

		ginkgo.It("should enforce minimum width of 20", func() {
			card.SetWidth(5)
			gomega.Expect(card.width).To(gomega.Equal(20))
		})

		ginkgo.It("should handle large widths", func() {
			card.SetWidth(200)
			gomega.Expect(card.width).To(gomega.Equal(200))
		})
	})

	ginkgo.Describe("Render", func() {
		ginkgo.It("should render nil fact as empty string", func() {
			card.fact = nil
			result := card.Render()
			gomega.Expect(result).To(gomega.Equal(""))
		})

		ginkgo.It("should render fact text", func() {
			result := card.Render()
			gomega.Expect(result).To(gomega.ContainSubstring(fact.Text))
		})

		ginkgo.It("should render all competencies", func() {
			result := card.Render()
			gomega.Expect(result).To(gomega.ContainSubstring("leadership"))
			gomega.Expect(result).To(gomega.ContainSubstring("technical"))
		})

		ginkgo.It("should render role fit with icon", func() {
			result := card.Render()
			gomega.Expect(result).To(gomega.ContainSubstring("principal"))
		})

		ginkgo.It("should render audience relevance", func() {
			result := card.Render()
			gomega.Expect(result).To(gomega.ContainSubstring("hiring_manager"))
			gomega.Expect(result).To(gomega.ContainSubstring("recruiter"))
		})

		ginkgo.It("should render strength signal", func() {
			result := card.Render()
			gomega.Expect(result).To(gomega.ContainSubstring("leadership"))
			gomega.Expect(result).To(gomega.ContainSubstring("Strength"))
		})

		ginkgo.It("should render source event reference", func() {
			result := card.Render()
			gomega.Expect(result).To(gomega.ContainSubstring("event-123"))
		})

		ginkgo.It("should render creation date", func() {
			result := card.Render()
			gomega.Expect(result).To(gomega.ContainSubstring("Created"))
		})

		ginkgo.It("should omit empty strength signal", func() {
			fact.StrengthSignal = ""
			result := card.Render()
			gomega.Expect(strings.Count(result, "Strength:")).To(gomega.Equal(0))
		})

		ginkgo.It("should render burst source instead of event source", func() {
			fact.SourceEventID = ""
			fact.SourceBurstID = "burst-456"
			result := card.Render()
			gomega.Expect(result).To(gomega.ContainSubstring("burst-456"))
		})
	})

	ginkgo.Describe("renderText", func() {
		ginkgo.It("should wrap long text", func() {
			fact.Text = "This is a very long fact text that should be wrapped to fit the specified width constraint for proper display in the terminal interface and readability"
			card.SetWidth(30)
			result := card.renderText()
			lines := strings.Split(result, "\n")
			gomega.Expect(len(lines) > 1).To(gomega.BeTrue())
		})

		ginkgo.It("should preserve short text without wrapping", func() {
			fact.Text = "Short text"
			result := card.renderText()
			gomega.Expect(result).To(gomega.ContainSubstring("Short text"))
		})
	})

	ginkgo.Describe("renderCompetencies", func() {
		ginkgo.It("should render competency badges", func() {
			result := card.renderCompetencies()
			gomega.Expect(result).NotTo(gomega.Equal(""))
		})

		ginkgo.It("should handle empty competencies", func() {
			fact.CompetencyCategories = []string{}
			result := card.renderCompetencies()
			gomega.Expect(result).To(gomega.Equal(""))
		})

		ginkgo.It("should render multiple competencies", func() {
			fact.CompetencyCategories = []string{"technical", "leadership", "product"}
			result := card.renderCompetencies()
			gomega.Expect(result).To(gomega.ContainSubstring("technical"))
			gomega.Expect(result).To(gomega.ContainSubstring("leadership"))
			gomega.Expect(result).To(gomega.ContainSubstring("product"))
		})
	})

	ginkgo.Describe("renderRoleFit", func() {
		ginkgo.It("should render principal role fit", func() {
			fact.RoleFit = career.RoleFitPrincipal
			result := card.renderRoleFit()
			gomega.Expect(result).To(gomega.ContainSubstring("principal"))
		})

		ginkgo.It("should render em role fit", func() {
			fact.RoleFit = career.RoleFitEM
			result := card.renderRoleFit()
			gomega.Expect(result).To(gomega.ContainSubstring("em"))
		})

		ginkgo.It("should render staff role fit", func() {
			fact.RoleFit = career.RoleFitStaff
			result := card.renderRoleFit()
			gomega.Expect(result).To(gomega.ContainSubstring("staff"))
		})

		ginkgo.It("should render senior ic role fit", func() {
			fact.RoleFit = career.RoleFitSeniorIC
			result := card.renderRoleFit()
			gomega.Expect(result).To(gomega.ContainSubstring("senior_ic"))
		})
	})

	ginkgo.Describe("getRoleFitIcon", func() {
		ginkgo.It("should return crown icon for principal", func() {
			icon := getRoleFitIcon(career.RoleFitPrincipal)
			gomega.Expect(icon).To(gomega.Equal("👑"))
		})

		ginkgo.It("should return briefcase icon for em", func() {
			icon := getRoleFitIcon(career.RoleFitEM)
			gomega.Expect(icon).To(gomega.Equal("👔"))
		})

		ginkgo.It("should return star icon for staff", func() {
			icon := getRoleFitIcon(career.RoleFitStaff)
			gomega.Expect(icon).To(gomega.Equal("⭐"))
		})

		ginkgo.It("should return medal icon for senior ic", func() {
			icon := getRoleFitIcon(career.RoleFitSeniorIC)
			gomega.Expect(icon).To(gomega.Equal("🎖️"))
		})

		ginkgo.It("should return sparkle icon for unknown", func() {
			icon := getRoleFitIcon("")
			gomega.Expect(icon).To(gomega.Equal("✨"))
		})
	})

	ginkgo.Describe("wrapText", func() {
		ginkgo.It("should wrap text to specified width", func() {
			text := "This is a test text that should be wrapped properly to fit the width constraint"
			result := wrapText(text, 20)
			lines := strings.Split(result, "\n")
			for _, line := range lines {
				gomega.Expect(len(line) <= 20).To(gomega.BeTrue())
			}
		})

		ginkgo.It("should handle very short width", func() {
			text := "This is a test"
			result := wrapText(text, 5)
			lines := strings.Split(result, "\n")
			gomega.Expect(len(lines) == 2).To(gomega.BeTrue())
		})

		ginkgo.It("should preserve single words", func() {
			text := "SingleLongWord AnotherWord"
			result := wrapText(text, 15)
			gomega.Expect(result).To(gomega.ContainSubstring("SingleLongWord"))
			gomega.Expect(result).To(gomega.ContainSubstring("AnotherWord"))
		})

		ginkgo.It("should enforce minimum width of 10", func() {
			text := "Test text"
			result := wrapText(text, 2)
			lines := strings.Split(result, "\n")
			for _, line := range lines {
				gomega.Expect(len(line) >= 4).To(gomega.BeTrue())
			}
		})

		ginkgo.It("should handle empty text", func() {
			result := wrapText("", 20)
			gomega.Expect(result).To(gomega.Equal(""))
		})

		ginkgo.It("should handle single word", func() {
			result := wrapText("Hello", 20)
			gomega.Expect(result).To(gomega.Equal("Hello"))
		})
	})

	ginkgo.Describe("renderSource", func() {
		ginkgo.It("should render event source", func() {
			fact.SourceEventID = "event-123"
			fact.SourceBurstID = ""
			result := card.renderSource()
			gomega.Expect(result).To(gomega.ContainSubstring("event-123"))
		})

		ginkgo.It("should render burst source", func() {
			fact.SourceEventID = ""
			fact.SourceBurstID = "burst-456"
			result := card.renderSource()
			gomega.Expect(result).To(gomega.ContainSubstring("burst-456"))
		})

		ginkgo.It("should render unknown source", func() {
			fact.SourceEventID = ""
			fact.SourceBurstID = ""
			result := card.renderSource()
			gomega.Expect(result).To(gomega.ContainSubstring("Unknown"))
		})

		ginkgo.It("should prefer event source over burst", func() {
			fact.SourceEventID = "event-123"
			fact.SourceBurstID = "burst-456"
			result := card.renderSource()
			gomega.Expect(result).To(gomega.ContainSubstring("event-123"))
			gomega.Expect(result).NotTo(gomega.ContainSubstring("burst-456"))
		})
	})

	ginkgo.Describe("renderAudience", func() {
		ginkgo.It("should render audience relevance", func() {
			fact.AudienceRelevance = []string{"hiring_manager", "recruiter"}
			result := card.renderAudience()
			gomega.Expect(result).To(gomega.ContainSubstring("hiring_manager"))
			gomega.Expect(result).To(gomega.ContainSubstring("recruiter"))
		})

		ginkgo.It("should handle empty audience", func() {
			fact.AudienceRelevance = []string{}
			result := card.renderAudience()
			gomega.Expect(result).To(gomega.Equal(""))
		})

		ginkgo.It("should handle single audience", func() {
			fact.AudienceRelevance = []string{"peer"}
			result := card.renderAudience()
			gomega.Expect(result).To(gomega.ContainSubstring("peer"))
		})
	})

	ginkgo.Describe("renderStrengthSignal", func() {
		ginkgo.It("should render strength signal", func() {
			fact.StrengthSignal = "leadership"
			result := card.renderStrengthSignal()
			gomega.Expect(result).To(gomega.ContainSubstring("leadership"))
		})

		ginkgo.It("should handle empty strength signal", func() {
			fact.StrengthSignal = ""
			result := card.renderStrengthSignal()
			gomega.Expect(result).To(gomega.Equal(""))
		})
	})

	ginkgo.Describe("renderMetadata", func() {
		ginkgo.It("should render creation date", func() {
			result := card.renderMetadata()
			gomega.Expect(result).To(gomega.ContainSubstring("Created"))
		})

		ginkgo.It("should include current date", func() {
			result := card.renderMetadata()
			today := time.Now().Format("2006-01-02")
			gomega.Expect(result).To(gomega.ContainSubstring(today))
		})
	})
})

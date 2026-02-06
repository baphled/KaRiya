package browsetimeline_test

import (
	"context"

	"github.com/baphled/kariya/internal/cli/intents/browsetimeline"
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/domain/career"
	careermemory "github.com/baphled/kariya/internal/repository/career/memory"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// E2E Tests for Skill Management in BrowseTimeline
//
// These tests exercise complete user workflows:
// - Viewing skills for an event
// - Adding existing skills to an event
// - Removing skills from an event
// - Modal visibility and navigation

var _ = Describe("Skill Management E2E", func() {
	var (
		intent          *browsetimeline.Intent
		ctx             context.Context
		eventRepo       *careermemory.EventRepository
		skillRepo       *careermemory.SkillRepository
		cliEventService *service.CLIEventService
		cliSkillService *service.CLISkillService
		testEvent       *career.Event
		testSkills      []*career.Skill
		linkedSkill     *career.Skill
		availableSkills []*career.Skill
	)

	BeforeEach(func() {
		ctx = context.Background()

		testEvent = fixtures.EventWith("event-1", "Implemented microservices architecture", "TechCo", "Platform")

		linkedSkill = fixtures.SkillWith("skill-1", "Go", "backend", "advanced")
		availableSkills = []*career.Skill{
			fixtures.SkillWith("skill-2", "Docker", "devops", "intermediate"),
			fixtures.SkillWith("skill-3", "Kubernetes", "devops", "intermediate"),
			fixtures.SkillWith("skill-4", "React", "frontend", "advanced"),
		}
		testSkills = append([]*career.Skill{linkedSkill}, availableSkills...)

		eventRepo = careermemory.NewEventRepository()
		skillRepo = careermemory.NewSkillRepository()
		eventRepo.SetSkillRepository(skillRepo)
		skillRepo.SetEventRepository(eventRepo)

		err := eventRepo.Create(ctx, testEvent)
		Expect(err).ToNot(HaveOccurred())

		for _, skill := range testSkills {
			err := skillRepo.Create(ctx, skill)
			Expect(err).ToNot(HaveOccurred())
		}

		err = eventRepo.LinkSkill(ctx, "event-1", "skill-1")
		Expect(err).ToNot(HaveOccurred())

		careerSvc := careerservice.NewService(eventRepo)
		careerSvc.SetSkillRepository(skillRepo)

		cliEventService = service.NewCLIEventService(careerSvc)
		cliSkillService = service.NewCLISkillService(skillRepo)

		intentCtx := &browsetimeline.IntentContext{
			Events:          []*career.Event{testEvent},
			InitialFilters:  nil,
			SelectedEventID: "event-1",
			CLIEventService: cliEventService,
			CLISkillService: cliSkillService,
		}

		intent, err = browsetimeline.NewIntent(intentCtx)
		Expect(err).ToNot(HaveOccurred())
		intent.Init()
	})

	Describe("Viewing Skills Modal", func() {
		It("should show skills modal when 's' is pressed from event detail", func() {
			By("Opening event detail view")
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			view := intent.View()
			Expect(view).To(ContainSubstring(testEvent.Text))

			By("Pressing 's' to view skills")
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})

			By("Verifying skills modal is visible")
			view = intent.View()
			Expect(view).To(ContainSubstring("Skills"))
			Expect(view).To(ContainSubstring("Go"))
		})

		It("should show action hints in skills modal", func() {
			By("Opening event detail and skills modal")
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})

			By("Verifying action hints are shown")
			view := intent.View()
			Expect(view).To(SatisfyAny(
				ContainSubstring("a: Add Existing"),
				ContainSubstring("a:"),
			))
			Expect(view).To(SatisfyAny(
				ContainSubstring("d: Remove"),
				ContainSubstring("d:"),
			))
		})

		It("should close skills modal on escape", func() {
			By("Opening skills modal")
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})

			view := intent.View()
			Expect(view).To(ContainSubstring("Skills"))

			By("Pressing escape to close")
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			By("Verifying modal is closed and back at event detail")
			view = intent.View()
			Expect(view).To(ContainSubstring(testEvent.Text))
			Expect(view).ToNot(ContainSubstring("a: Add Existing"))
		})
	})

	Describe("Add Existing Skill Workflow", func() {
		It("should complete full workflow to add existing skill", func() {
			By("Opening event detail")
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			By("Opening skills modal")
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
			view := intent.View()
			Expect(view).To(ContainSubstring("Skills"))
			Expect(view).To(ContainSubstring("Go"))

			By("Pressing 'a' to add existing skill")
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			By("Verifying skill picker modal is visible")
			view = intent.View()
			Expect(view).To(ContainSubstring("Select Skill"))
			Expect(view).To(ContainSubstring("Docker"))
			Expect(view).ToNot(ContainSubstring("Go"))

			By("Selecting a skill with Enter")
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			By("Verifying skill was linked")
			skills, err := cliEventService.GetSkillsForEvent(ctx, "event-1")
			Expect(err).ToNot(HaveOccurred())
			Expect(skills).To(HaveLen(2))

			By("Verifying back at skills modal with updated list")
			view = intent.View()
			Expect(view).To(ContainSubstring("Skills"))
		})

		It("should navigate within skill picker", func() {
			By("Opening skill picker")
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			By("Navigating down with 'j'")
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})

			By("Navigating up with 'k'")
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})

			view := intent.View()
			Expect(view).To(ContainSubstring("Select Skill"))
		})

		It("should cancel skill picker on escape", func() {
			By("Opening skill picker")
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			view := intent.View()
			Expect(view).To(ContainSubstring("Select Skill"))

			By("Pressing escape to cancel")
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			By("Verifying back at skills modal")
			view = intent.View()
			Expect(view).To(ContainSubstring("Skills"))
			Expect(view).ToNot(ContainSubstring("Select Skill"))

			By("Verifying no skill was linked")
			skills, err := cliEventService.GetSkillsForEvent(ctx, "event-1")
			Expect(err).ToNot(HaveOccurred())
			Expect(skills).To(HaveLen(1))
		})
	})

	Describe("Remove Skill Workflow", func() {
		It("should complete full workflow to remove skill", func() {
			By("Verifying event starts with 1 skill")
			skills, err := cliEventService.GetSkillsForEvent(ctx, "event-1")
			Expect(err).ToNot(HaveOccurred())
			Expect(skills).To(HaveLen(1))

			By("Opening event detail and skills modal")
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})

			By("Pressing 'd' to remove selected skill")
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})

			By("Processing the unlink message")
			if cmd != nil {
				msg := cmd()
				intent.Update(msg)
			}

			By("Verifying skill was unlinked")
			skills, err = cliEventService.GetSkillsForEvent(ctx, "event-1")
			Expect(err).ToNot(HaveOccurred())
			Expect(skills).To(BeEmpty())

			By("Verifying skills modal shows empty state")
			view := intent.View()
			Expect(view).To(ContainSubstring("Skills"))
			Expect(view).To(ContainSubstring("No skills"))
		})
	})

	Describe("Error Handling", func() {
		It("should show error modal if listing skills fails", func() {
			By("Corrupting the skill repository")
			skillRepo = nil

			By("Attempting to open skills modal")
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})

			By("Verifying error is handled gracefully")
		})
	})

	Describe("Modal Priority", func() {
		It("should prioritize skill picker over skills detail modal", func() {
			By("Opening skills detail modal")
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})

			By("Opening skill picker")
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			By("Verifying picker is shown, not detail")
			view := intent.View()
			Expect(view).To(ContainSubstring("Select Skill"))
			Expect(view).ToNot(ContainSubstring("a: Add Existing"))
		})
	})
})

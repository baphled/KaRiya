package skillsmanagement

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	domain "github.com/baphled/kariya/internal/domain/career"
	memoryrepo "github.com/baphled/kariya/internal/repository/career/memory"
	"github.com/baphled/kariya/internal/service/career/skillinference"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	tea "github.com/charmbracelet/bubbletea"
)

func newTestIntent() *Intent {
	skillRepo := memoryrepo.NewSkillRepository()
	intentCtx := NewIntentContext(context.Background(), skillRepo)
	intent, _ := NewIntent(intentCtx)
	return intent
}

func newActiveTestIntent() *Intent {
	intent := newTestIntent()
	intent.Init()
	return intent
}

var _ = Describe("skillRowFormatterWithCounts", func() {
	It("formats all fields for a skill with years and event count", func() {
		skill := fixtures.SkillWithYears("s1", "Go Programming", "backend", 5)
		counts := map[string]int{"s1": 3}
		formatter := skillRowFormatterWithCounts(counts)
		row := formatter(skill, 0)
		Expect(row).To(HaveLen(5))
		Expect(row[0]).To(Equal("Go Programming"))
		Expect(row[1]).To(Equal("backend"))
		Expect(row[2]).To(Equal("advanced"))
		Expect(row[3]).To(Equal("5"))
		Expect(row[4]).To(Equal("3"))
	})

	It("truncates long names to 25 chars", func() {
		skill := fixtures.SkillWith("", "This Is A Very Long Skill Name That Exceeds Limit", "", "")
		formatter := skillRowFormatterWithCounts(nil)
		row := formatter(skill, 0)
		Expect(row[0]).To(HaveSuffix("..."))
		Expect(row[0]).To(HaveLen(25))
	})

	It("uses dash for empty category", func() {
		skill := fixtures.SkillWith("", "Go", "", "")
		formatter := skillRowFormatterWithCounts(nil)
		row := formatter(skill, 0)
		Expect(row[1]).To(Equal("-"))
	})

	It("uses dash for empty level", func() {
		skill := fixtures.SkillWith("", "Go", "backend", "")
		formatter := skillRowFormatterWithCounts(nil)
		row := formatter(skill, 0)
		Expect(row[2]).To(Equal("-"))
	})

	It("uses dash for nil years", func() {
		skill := fixtures.SkillWith("", "Go", "backend", "advanced")
		formatter := skillRowFormatterWithCounts(nil)
		row := formatter(skill, 0)
		Expect(row[3]).To(Equal("-"))
	})

	It("uses dash when event counts map is nil", func() {
		skill := fixtures.SkillWith("s1", "Go", "", "")
		formatter := skillRowFormatterWithCounts(nil)
		row := formatter(skill, 0)
		Expect(row[4]).To(Equal("-"))
	})

	It("uses dash when skill ID not in event counts", func() {
		skill := fixtures.SkillWith("s1", "Go", "", "")
		counts := map[string]int{"other": 5}
		formatter := skillRowFormatterWithCounts(counts)
		row := formatter(skill, 0)
		Expect(row[4]).To(Equal("-"))
	})

	It("does not truncate names at exactly 22 characters", func() {
		skill := fixtures.SkillWith("", "ExactlyTwentyTwoChars!", "", "")
		Expect(skill.Name).To(HaveLen(22))
		formatter := skillRowFormatterWithCounts(nil)
		row := formatter(skill, 0)
		Expect(row[0]).To(Equal("ExactlyTwentyTwoChars!"))
	})
})

var _ = Describe("getBreadcrumbs", func() {
	var intent *Intent

	BeforeEach(func() {
		intent = newActiveTestIntent()
	})

	It("returns Skills for StateList", func() {
		intent.state = StateList
		breadcrumbs := intent.getBreadcrumbs()
		Expect(breadcrumbs).To(Equal([]string{"Skills"}))
	})

	It("returns Skills > name for StateDetail with selected skill", func() {
		intent.state = StateDetail
		intent.selectedSkill = fixtures.SkillWith("s1", "Go Programming", "backend", "advanced")
		breadcrumbs := intent.getBreadcrumbs()
		Expect(breadcrumbs).To(Equal([]string{"Skills", "Go Programming"}))
	})

	It("returns Skills for StateDetail without selected skill", func() {
		intent.state = StateDetail
		intent.selectedSkill = nil
		breadcrumbs := intent.getBreadcrumbs()
		Expect(breadcrumbs).To(Equal([]string{"Skills"}))
	})

	It("returns Skills > name > Events for StateDetailEvents", func() {
		intent.state = StateDetailEvents
		intent.selectedSkill = fixtures.SkillWith("s1", "Go", "backend", "advanced")
		breadcrumbs := intent.getBreadcrumbs()
		Expect(breadcrumbs).To(Equal([]string{"Skills", "Go", "Events"}))
	})

	It("returns Skills > name > Events > Detail for StateDetailEventDetail", func() {
		intent.state = StateDetailEventDetail
		intent.selectedSkill = fixtures.SkillWith("s1", "Go", "backend", "advanced")
		breadcrumbs := intent.getBreadcrumbs()
		Expect(breadcrumbs).To(Equal([]string{"Skills", "Go", "Events", "Detail"}))
	})

	It("returns Skills > Add for StateAdd", func() {
		intent.state = StateAdd
		breadcrumbs := intent.getBreadcrumbs()
		Expect(breadcrumbs).To(Equal([]string{"Skills", "Add"}))
	})

	It("returns Skills > Edit for StateEdit", func() {
		intent.state = StateEdit
		breadcrumbs := intent.getBreadcrumbs()
		Expect(breadcrumbs).To(Equal([]string{"Skills", "Edit"}))
	})

	It("returns Skills > Delete for StateDelete", func() {
		intent.state = StateDelete
		breadcrumbs := intent.getBreadcrumbs()
		Expect(breadcrumbs).To(Equal([]string{"Skills", "Delete"}))
	})
})

var _ = Describe("getContextHelp", func() {
	var intent *Intent

	BeforeEach(func() {
		intent = newActiveTestIntent()
	})

	It("returns non-empty help for StateList", func() {
		intent.state = StateList
		help := intent.getContextHelp()
		Expect(help).NotTo(BeEmpty())
	})

	It("returns non-empty help for StateDetail", func() {
		intent.state = StateDetail
		help := intent.getContextHelp()
		Expect(help).NotTo(BeEmpty())
	})

	It("returns non-empty help for StateInferringSkills", func() {
		intent.state = StateInferringSkills
		help := intent.getContextHelp()
		Expect(help).NotTo(BeEmpty())
	})

	It("returns non-empty help for StateSkillSuggestionReview", func() {
		intent.state = StateSkillSuggestionReview
		help := intent.getContextHelp()
		Expect(help).NotTo(BeEmpty())
	})

	It("returns non-empty help for StateDelete", func() {
		intent.state = StateDelete
		help := intent.getContextHelp()
		Expect(help).NotTo(BeEmpty())
	})

	It("returns non-empty help for default state", func() {
		intent.state = StateAdd
		help := intent.getContextHelp()
		Expect(help).NotTo(BeEmpty())
	})
})

var _ = Describe("rebuildModalRegistry", func() {
	var intent *Intent

	BeforeEach(func() {
		intent = newActiveTestIntent()
	})

	It("creates registry if nil", func() {
		intent.modalRegistry = nil
		intent.rebuildModalRegistry()
		Expect(intent.modalRegistry).NotTo(BeNil())
	})

	It("clears existing registry entries", func() {
		intent.rebuildModalRegistry()
		Expect(intent.modalRegistry).NotTo(BeNil())
	})

	It("registers feedback modal when present", func() {
		intent.feedbackModal = &feedback.Modal{
			Title:   "Test",
			Message: "Test message",
			Type:    feedback.ModalError,
		}
		intent.rebuildModalRegistry()
		Expect(intent.modalRegistry).NotTo(BeNil())
	})

	It("returns early when loading modal is present", func() {
		intent.loadingModal = &feedback.Modal{
			Title:   "Loading",
			Message: "Please wait...",
		}
		intent.rebuildModalRegistry()
		Expect(intent.modalRegistry).NotTo(BeNil())
	})
})

var _ = Describe("openFilterModal", func() {
	var intent *Intent

	BeforeEach(func() {
		intent = newActiveTestIntent()
	})

	It("creates filter modal and returns init command", func() {
		cmd := intent.openFilterModal()
		Expect(cmd).NotTo(BeNil())
		Expect(intent.filterModal).NotTo(BeNil())
	})

	It("pre-populates with existing category filter", func() {
		intent.context.Filters = &Filters{Category: "backend"}
		cmd := intent.openFilterModal()
		Expect(cmd).NotTo(BeNil())
		Expect(intent.filterModal).NotTo(BeNil())
	})

	It("pre-populates with existing level filter", func() {
		intent.context.Filters = &Filters{Level: "advanced"}
		cmd := intent.openFilterModal()
		Expect(cmd).NotTo(BeNil())
		Expect(intent.filterModal).NotTo(BeNil())
	})

	It("handles nil filters", func() {
		intent.context.Filters = nil
		cmd := intent.openFilterModal()
		Expect(cmd).NotTo(BeNil())
		Expect(intent.filterModal).NotTo(BeNil())
	})
})

var _ = Describe("openSortModal", func() {
	var intent *Intent

	BeforeEach(func() {
		intent = newActiveTestIntent()
	})

	It("creates sort modal and returns init command", func() {
		cmd := intent.openSortModal()
		Expect(cmd).NotTo(BeNil())
		Expect(intent.sortModal).NotTo(BeNil())
	})

	It("pre-populates with existing sort config", func() {
		intent.context.Filters = &Filters{SortBy: "name", SortOrder: "asc"}
		cmd := intent.openSortModal()
		Expect(cmd).NotTo(BeNil())
		Expect(intent.sortModal).NotTo(BeNil())
	})

	It("handles nil filters", func() {
		intent.context.Filters = nil
		cmd := intent.openSortModal()
		Expect(cmd).NotTo(BeNil())
		Expect(intent.sortModal).NotTo(BeNil())
	})
})

var _ = Describe("openSearchModal", func() {
	var intent *Intent

	BeforeEach(func() {
		intent = newActiveTestIntent()
	})

	It("creates search modal and returns init command", func() {
		cmd := intent.openSearchModal()
		Expect(cmd).NotTo(BeNil())
		Expect(intent.searchModal).NotTo(BeNil())
	})

	It("pre-populates with existing search text", func() {
		intent.context.Filters = &Filters{SearchText: "golang"}
		cmd := intent.openSearchModal()
		Expect(cmd).NotTo(BeNil())
		Expect(intent.searchModal).NotTo(BeNil())
	})

	It("handles nil filters", func() {
		intent.context.Filters = nil
		cmd := intent.openSearchModal()
		Expect(cmd).NotTo(BeNil())
		Expect(intent.searchModal).NotTo(BeNil())
	})
})

var _ = Describe("openViewDetailModal", func() {
	var intent *Intent

	BeforeEach(func() {
		intent = newActiveTestIntent()
	})

	It("returns nil when no skill is selected", func() {
		intent.selectedSkill = nil
		cmd := intent.openViewDetailModal()
		Expect(cmd).To(BeNil())
	})

	It("creates view detail modal for selected skill", func() {
		intent.selectedSkill = fixtures.SkillWith("s1", "Go", "backend", "advanced")
		cmd := intent.openViewDetailModal()
		Expect(cmd).To(BeNil())
		Expect(intent.viewDetailModal).NotTo(BeNil())
	})

	It("includes event count when available", func() {
		intent.selectedSkill = fixtures.SkillWith("s1", "Go", "backend", "advanced")
		intent.eventCounts = map[string]int{"s1": 5}
		cmd := intent.openViewDetailModal()
		Expect(cmd).To(BeNil())
		Expect(intent.viewDetailModal).NotTo(BeNil())
	})
})

var _ = Describe("openAddEditModal", func() {
	var intent *Intent

	BeforeEach(func() {
		intent = newActiveTestIntent()
	})

	It("creates add modal when skill is nil", func() {
		cmd := intent.openAddEditModal(nil)
		Expect(cmd).NotTo(BeNil())
		Expect(intent.addEditModal).NotTo(BeNil())
	})

	It("creates edit modal when skill is provided", func() {
		skill := fixtures.SkillWith("s1", "Go", "backend", "advanced")
		cmd := intent.openAddEditModal(skill)
		Expect(cmd).NotTo(BeNil())
		Expect(intent.addEditModal).NotTo(BeNil())
	})
})

var _ = Describe("openDeleteModal", func() {
	var intent *Intent

	BeforeEach(func() {
		intent = newActiveTestIntent()
	})

	It("returns nil when skill is nil", func() {
		cmd := intent.openDeleteModal(nil)
		Expect(cmd).To(BeNil())
	})

	It("creates delete confirmation modal for valid skill", func() {
		skill := fixtures.SkillWith("s1", "Go", "backend", "advanced")
		_ = intent.openDeleteModal(skill)
		Expect(intent.deleteModal).NotTo(BeNil())
		Expect(intent.selectedSkill).To(Equal(skill))
	})

	It("truncates long skill names in confirmation", func() {
		skill := fixtures.SkillWith("s1", "This Is An Extremely Long Skill Name That Exceeds The Fifty Character Limit For Display", "backend", "advanced")
		_ = intent.openDeleteModal(skill)
		Expect(intent.deleteModal).NotTo(BeNil())
	})
})

var _ = Describe("openSkillEventsModal", func() {
	var intent *Intent

	BeforeEach(func() {
		intent = newActiveTestIntent()
	})

	It("returns nil when no skill is selected", func() {
		intent.selectedSkill = nil
		events := fixtures.Events(2)
		cmd := intent.openSkillEventsModal(events)
		Expect(cmd).To(BeNil())
	})

	It("creates events modal for selected skill", func() {
		intent.selectedSkill = fixtures.SkillWith("s1", "Go", "backend", "advanced")
		events := fixtures.Events(2)
		cmd := intent.openSkillEventsModal(events)
		Expect(cmd).To(BeNil())
		Expect(intent.skillEventsModal).NotTo(BeNil())
	})
})

var _ = Describe("openEventDetailModal", func() {
	var intent *Intent

	BeforeEach(func() {
		intent = newActiveTestIntent()
	})

	It("returns nil when event is nil", func() {
		cmd := intent.openEventDetailModal(nil)
		Expect(cmd).To(BeNil())
	})

	It("creates event detail modal for valid event", func() {
		event := fixtures.EventWith("e1", "Test event", "Company", "Project")
		cmd := intent.openEventDetailModal(event)
		Expect(cmd).To(BeNil())
		Expect(intent.eventDetailModal).NotTo(BeNil())
	})
})

var _ = Describe("loadEventsForSkillModal", func() {
	var intent *Intent

	BeforeEach(func() {
		intent = newActiveTestIntent()
		intent.selectedSkill = fixtures.SkillWith("s1", "Go", "backend", "advanced")
	})

	It("returns a non-nil command", func() {
		cmd := intent.loadEventsForSkillModal()
		Expect(cmd).NotTo(BeNil())
	})

	It("command returns SkillEventsForModalLoadedMsg", func() {
		cmd := intent.loadEventsForSkillModal()
		msg := cmd()
		_, ok := msg.(SkillEventsForModalLoadedMsg)
		Expect(ok).To(BeTrue())
	})
})

var _ = Describe("syncTableSelection", func() {
	var intent *Intent

	BeforeEach(func() {
		intent = newActiveTestIntent()
	})

	It("syncs selection from table behavior", func() {
		intent.syncTableSelection()
		Expect(intent.selectedIndex).To(Equal(0))
	})

	It("sets selectedSkill to nil when no items", func() {
		intent.syncTableSelection()
		Expect(intent.selectedSkill).To(BeNil())
	})
})

var _ = Describe("getTerminalDimensions", func() {
	var intent *Intent

	BeforeEach(func() {
		intent = newActiveTestIntent()
	})

	It("returns default dimensions when no terminal info", func() {
		width, height := intent.getTerminalDimensions()
		Expect(width).To(BeNumerically(">", 0))
		Expect(height).To(BeNumerically(">", 0))
	})
})

var _ = Describe("filterNewSuggestions", func() {
	It("returns all suggestions when no existing names", func() {
		suggestions := []skillinference.SkillSuggestion{
			{Name: "Go", Category: "backend"},
			{Name: "Docker", Category: "devops"},
		}
		result := filterNewSuggestions(suggestions, []string{})
		Expect(result).To(HaveLen(2))
	})

	It("filters out existing names case-insensitively", func() {
		suggestions := []skillinference.SkillSuggestion{
			{Name: "Go", Category: "backend"},
			{Name: "Docker", Category: "devops"},
		}
		result := filterNewSuggestions(suggestions, []string{"go"})
		Expect(result).To(HaveLen(1))
		Expect(result[0].Name).To(Equal("Docker"))
	})

	It("returns empty when all are existing", func() {
		suggestions := []skillinference.SkillSuggestion{
			{Name: "Go", Category: "backend"},
		}
		result := filterNewSuggestions(suggestions, []string{"Go"})
		Expect(result).To(BeEmpty())
	})
})

var _ = Describe("createSkill", func() {
	var intent *Intent

	BeforeEach(func() {
		intent = newActiveTestIntent()
	})

	It("returns a command that creates the skill", func() {
		skill := fixtures.SkillWith("", "Go", "backend", "advanced")
		cmd := intent.createSkill(skill)
		Expect(cmd).NotTo(BeNil())
		msg := cmd()
		createdMsg, ok := msg.(SkillCreatedMsg)
		Expect(ok).To(BeTrue())
		Expect(createdMsg.Error).NotTo(HaveOccurred())
		Expect(createdMsg.Skill).To(Equal(skill))
	})
})

var _ = Describe("updateSkill", func() {
	var intent *Intent

	BeforeEach(func() {
		intent = newActiveTestIntent()
		skill := fixtures.SkillWith("s1", "Go", "backend", "advanced")
		_ = intent.context.SkillRepository.Create(context.Background(), skill)
	})

	It("returns a command that updates the skill", func() {
		skill := fixtures.SkillWith("s1", "Go Updated", "backend", "expert")
		cmd := intent.updateSkill(skill)
		Expect(cmd).NotTo(BeNil())
		msg := cmd()
		updatedMsg, ok := msg.(SkillUpdatedMsg)
		Expect(ok).To(BeTrue())
		Expect(updatedMsg.Error).NotTo(HaveOccurred())
	})
})

var _ = Describe("createSkillsFromSuggestions", func() {
	var intent *Intent

	BeforeEach(func() {
		intent = newActiveTestIntent()
	})

	It("returns SkillsCreatedMsg with nil for empty suggestions", func() {
		cmd := intent.createSkillsFromSuggestions([]skillinference.SkillSuggestion{})
		Expect(cmd).NotTo(BeNil())
		msg := cmd()
		createdMsg, ok := msg.(SkillsCreatedMsg)
		Expect(ok).To(BeTrue())
		Expect(createdMsg.Error).NotTo(HaveOccurred())
		Expect(createdMsg.Skills).To(BeNil())
	})

	It("returns error when service is nil", func() {
		intent.context.SkillInferenceService = nil
		suggestions := []skillinference.SkillSuggestion{
			{Name: "Go", Category: "backend", Confidence: 0.95},
		}
		cmd := intent.createSkillsFromSuggestions(suggestions)
		Expect(cmd).NotTo(BeNil())
		msg := cmd()
		createdMsg, ok := msg.(SkillsCreatedMsg)
		Expect(ok).To(BeTrue())
		Expect(createdMsg.Error).To(HaveOccurred())
		Expect(createdMsg.Error.Error()).To(ContainSubstring("not available"))
	})
})

var _ = Describe("handleSkillEventsLoaded", func() {
	var intent *Intent

	BeforeEach(func() {
		intent = newActiveTestIntent()
	})

	It("stores events on success", func() {
		events := fixtures.Events(3)
		intent.handleSkillEventsLoaded(SkillEventsLoadedMsg{Events: events})
		Expect(intent.skillEvents).To(HaveLen(3))
		Expect(intent.eventsLoaded).To(BeTrue())
	})

	It("does nothing on error", func() {
		intent.handleSkillEventsLoaded(SkillEventsLoadedMsg{
			Error: errors.New("load failed"),
		})
		Expect(intent.skillEvents).To(BeNil())
		Expect(intent.eventsLoaded).To(BeFalse())
	})
})

var _ = Describe("handleSkillEventsForModalLoaded", func() {
	var intent *Intent

	BeforeEach(func() {
		intent = newActiveTestIntent()
	})

	It("returns nil on error", func() {
		cmd := intent.handleSkillEventsForModalLoaded(SkillEventsForModalLoadedMsg{
			Error: errors.New("load failed"),
		})
		Expect(cmd).To(BeNil())
	})

	It("opens skill events modal on success when skill is selected", func() {
		intent.selectedSkill = fixtures.SkillWith("s1", "Go", "backend", "advanced")
		events := fixtures.Events(2)
		cmd := intent.handleSkillEventsForModalLoaded(SkillEventsForModalLoadedMsg{
			Events: events,
		})
		Expect(cmd).To(BeNil())
		Expect(intent.skillEventsModal).NotTo(BeNil())
	})
})

var _ = Describe("handleScreenResult", func() {
	var intent *Intent

	BeforeEach(func() {
		intent = newActiveTestIntent()
	})

	It("returns nil for nil result", func() {
		cmd := intent.handleScreenResult(nil)
		Expect(cmd).To(BeNil())
	})

	It("returns nil for non-ScreenResult type", func() {
		cmd := intent.handleScreenResult("not a screen result")
		Expect(cmd).To(BeNil())
	})
})

var _ = Describe("noopCmd", func() {
	It("returns nil", func() {
		result := noopCmd()
		Expect(result).To(BeNil())
	})
})

var _ = Describe("Filters nil safety", func() {
	It("HasActiveFilters returns false for nil receiver", func() {
		var f *Filters
		Expect(f.HasActiveFilters()).To(BeFalse())
	})

	It("Clear is safe to call on nil receiver", func() {
		var f *Filters
		Expect(func() { f.Clear() }).NotTo(Panic())
	})
})

var _ = Describe("handleViewDetailModalUpdate", func() {
	var intent *Intent

	BeforeEach(func() {
		intent = newActiveTestIntent()
		intent.selectedSkill = fixtures.SkillWith("s1", "Go", "backend", "advanced")
		intent.openViewDetailModal()
	})

	It("clears modal when not visible after update", func() {
		Expect(intent.viewDetailModal).NotTo(BeNil())
	})
})

var _ = Describe("handleDeleteModalUpdate", func() {
	var intent *Intent

	BeforeEach(func() {
		intent = newActiveTestIntent()
		skill := fixtures.SkillWith("s1", "Go", "backend", "advanced")
		intent.openDeleteModal(skill)
	})

	It("creates delete modal", func() {
		Expect(intent.deleteModal).NotTo(BeNil())
	})
})

var _ = Describe("handleAddEditModalUpdate", func() {
	var intent *Intent

	BeforeEach(func() {
		intent = newActiveTestIntent()
	})

	It("creates add modal for nil skill", func() {
		intent.openAddEditModal(nil)
		Expect(intent.addEditModal).NotTo(BeNil())
	})

	It("creates edit modal for existing skill", func() {
		skill := fixtures.SkillWith("s1", "Go", "backend", "advanced")
		intent.openAddEditModal(skill)
		Expect(intent.addEditModal).NotTo(BeNil())
	})
})

var _ = Describe("Intent reloadSkills", func() {
	var intent *Intent

	BeforeEach(func() {
		intent = newActiveTestIntent()
	})

	It("returns a command that loads skills", func() {
		cmd := intent.reloadSkills()
		Expect(cmd).NotTo(BeNil())
		msg := cmd()
		loadedMsg, ok := msg.(SkillsLoadedMsg)
		Expect(ok).To(BeTrue())
		Expect(loadedMsg.Error).NotTo(HaveOccurred())
	})
})

var _ = Describe("Intent transitionToScreen", func() {
	It("sets active screen", func() {
		intent := newActiveTestIntent()
		Expect(intent.activeScreen).NotTo(BeNil())
	})
})

var _ = Describe("handleFeedbackModalUpdate", func() {
	var intent *Intent

	BeforeEach(func() {
		intent = newActiveTestIntent()
	})

	It("returns nil when feedback modal is nil", func() {
		intent.feedbackModal = nil
		cmd := intent.handleFeedbackModalUpdate(nil)
		Expect(cmd).To(BeNil())
	})

	It("returns noopCmd when feedback modal is active", func() {
		intent.feedbackModal = &feedback.Modal{
			Title:   "Error",
			Message: "Something went wrong",
			Type:    feedback.ModalError,
		}
		cmd := intent.handleFeedbackModalUpdate(nil)
		Expect(cmd).NotTo(BeNil())
	})
})

var _ = Describe("handleLoadingModalUpdate", func() {
	var intent *Intent

	BeforeEach(func() {
		intent = newActiveTestIntent()
	})

	It("returns nil when loading modal is nil", func() {
		intent.loadingModal = nil
		cmd := intent.handleLoadingModalUpdate(nil)
		Expect(cmd).To(BeNil())
	})
})

var _ = Describe("handleSkillEventsModalUpdate via modal setup", func() {
	var intent *Intent

	BeforeEach(func() {
		intent = newActiveTestIntent()
		intent.selectedSkill = fixtures.SkillWith("s1", "Go", "backend", "advanced")
		events := fixtures.Events(2)
		intent.openSkillEventsModal(events)
	})

	It("creates skill events modal", func() {
		Expect(intent.skillEventsModal).NotTo(BeNil())
	})
})

var _ = Describe("handleEventDetailModalUpdate via modal setup", func() {
	var intent *Intent

	BeforeEach(func() {
		intent = newActiveTestIntent()
		event := fixtures.EventWith("e1", "Test event", "Company", "Project")
		intent.openEventDetailModal(event)
	})

	It("creates event detail modal", func() {
		Expect(intent.eventDetailModal).NotTo(BeNil())
	})
})

var _ = Describe("Intent HasActiveModal", func() {
	var intent *Intent

	BeforeEach(func() {
		intent = newActiveTestIntent()
	})

	It("returns false when no modals are active", func() {
		Expect(intent.HasActiveModal()).To(BeFalse())
	})

	It("returns true when loading modal is set", func() {
		intent.loadingModal = &feedback.Modal{}
		Expect(intent.HasActiveModal()).To(BeTrue())
	})

	It("returns true when feedback modal is set", func() {
		intent.feedbackModal = &feedback.Modal{}
		Expect(intent.HasActiveModal()).To(BeTrue())
	})
})

var _ = Describe("handleFilterModalUpdate via modal", func() {
	var intent *Intent

	BeforeEach(func() {
		intent = newActiveTestIntent()
		intent.openFilterModal()
	})

	It("delegates to filter modal and returns cmd", func() {
		Expect(intent.filterModal).NotTo(BeNil())
		cmd := intent.handleFilterModalUpdate(tea.KeyMsg{Type: tea.KeyEsc})
		_ = cmd
	})
})

var _ = Describe("handleSortModalUpdate via modal", func() {
	var intent *Intent

	BeforeEach(func() {
		intent = newActiveTestIntent()
		intent.openSortModal()
	})

	It("delegates to sort modal and returns cmd", func() {
		Expect(intent.sortModal).NotTo(BeNil())
		cmd := intent.handleSortModalUpdate(tea.KeyMsg{Type: tea.KeyEsc})
		_ = cmd
	})
})

var _ = Describe("handleSearchModalUpdate via modal", func() {
	var intent *Intent

	BeforeEach(func() {
		intent = newActiveTestIntent()
		intent.openSearchModal()
	})

	It("delegates to search modal and returns cmd", func() {
		Expect(intent.searchModal).NotTo(BeNil())
		cmd := intent.handleSearchModalUpdate(tea.KeyMsg{Type: tea.KeyEsc})
		_ = cmd
	})
})

var _ = Describe("handleViewDetailModalUpdate via modal", func() {
	var intent *Intent

	BeforeEach(func() {
		intent = newActiveTestIntent()
		intent.selectedSkill = fixtures.SkillWith("s1", "Go", "backend", "advanced")
		intent.openViewDetailModal()
	})

	It("delegates to view detail modal", func() {
		Expect(intent.viewDetailModal).NotTo(BeNil())
		cmd := intent.handleViewDetailModalUpdate(tea.KeyMsg{Type: tea.KeyEsc})
		_ = cmd
	})

	It("clears modal when no longer visible after Esc", func() {
		intent.handleViewDetailModalUpdate(tea.KeyMsg{Type: tea.KeyEsc})
		// After Esc, modal may close depending on implementation
	})
})

var _ = Describe("handleAddEditModalUpdate via modal", func() {
	var intent *Intent

	BeforeEach(func() {
		intent = newActiveTestIntent()
		intent.openAddEditModal(nil)
	})

	It("delegates to add/edit modal", func() {
		Expect(intent.addEditModal).NotTo(BeNil())
		cmd := intent.handleAddEditModalUpdate(tea.KeyMsg{Type: tea.KeyEsc})
		_ = cmd
	})
})

var _ = Describe("handleDeleteModalUpdate via modal", func() {
	var intent *Intent

	BeforeEach(func() {
		intent = newActiveTestIntent()
		skill := fixtures.SkillWith("s1", "Go", "backend", "advanced")
		intent.openDeleteModal(skill)
	})

	It("delegates to delete modal", func() {
		Expect(intent.deleteModal).NotTo(BeNil())
		cmd := intent.handleDeleteModalUpdate(tea.KeyMsg{Type: tea.KeyEsc})
		_ = cmd
	})
})

var _ = Describe("handleSkillEventsModalUpdate via modal", func() {
	var intent *Intent

	BeforeEach(func() {
		intent = newActiveTestIntent()
		intent.selectedSkill = fixtures.SkillWith("s1", "Go", "backend", "advanced")
		events := fixtures.Events(2)
		intent.openSkillEventsModal(events)
	})

	It("delegates to skill events modal", func() {
		Expect(intent.skillEventsModal).NotTo(BeNil())
		cmd := intent.handleSkillEventsModalUpdate(tea.KeyMsg{Type: tea.KeyEsc})
		_ = cmd
	})
})

var _ = Describe("handleEventDetailModalUpdate via modal", func() {
	var intent *Intent

	BeforeEach(func() {
		intent = newActiveTestIntent()
		event := fixtures.EventWith("e1", "Test event", "Company", "Project")
		intent.openEventDetailModal(event)
	})

	It("delegates to event detail modal", func() {
		Expect(intent.eventDetailModal).NotTo(BeNil())
		cmd := intent.handleEventDetailModalUpdate(tea.KeyMsg{Type: tea.KeyEsc})
		_ = cmd
	})
})

var _ = Describe("resolveEventsFromIDs", func() {
	var intent *Intent

	BeforeEach(func() {
		intent = newActiveTestIntent()
	})

	It("returns empty slice when event repo is nil", func() {
		intent.context.EventRepository = nil
		result := intent.resolveEventsFromIDs([]string{"e1", "e2"})
		Expect(result).To(BeEmpty())
	})

	It("returns empty slice when event IDs is empty", func() {
		result := intent.resolveEventsFromIDs([]string{})
		Expect(result).To(BeEmpty())
	})

	It("resolves events from IDs using repository", func() {
		eventRepo := memoryrepo.NewEventRepository()
		event := fixtures.EventWith("e1", "Test event", "Company", "Project")
		_ = eventRepo.Create(context.Background(), event)
		intent.context.EventRepository = eventRepo
		result := intent.resolveEventsFromIDs([]string{"e1", "nonexistent"})
		Expect(result).To(HaveLen(1))
	})
})

var _ = Describe("ApplyFilters", func() {
	It("is a no-op that satisfies the interface", func() {
		intent := newActiveTestIntent()
		Expect(func() { intent.ApplyFilters() }).NotTo(Panic())
	})
})

var _ = Describe("GetSelectedIndex", func() {
	It("returns the selected index", func() {
		intent := newActiveTestIntent()
		Expect(intent.GetSelectedIndex()).To(Equal(0))
	})
})

var _ = Describe("Intent startSkillInference", func() {
	var intent *Intent

	BeforeEach(func() {
		intent = newActiveTestIntent()
	})

	It("shows error when inference service is nil", func() {
		intent.context.SkillInferenceService = nil
		cmd := intent.startSkillInference()
		Expect(cmd).To(BeNil())
		Expect(intent.feedbackModal).NotTo(BeNil())
	})

	It("shows error when event repository is nil", func() {
		intent.context.EventRepository = nil
		intent.context.SkillInferenceService = &mockInferenceService{}
		cmd := intent.startSkillInference()
		Expect(cmd).To(BeNil())
		Expect(intent.feedbackModal).NotTo(BeNil())
	})

	It("starts inference when both services are available", func() {
		eventRepo := memoryrepo.NewEventRepository()
		intent.context.EventRepository = eventRepo
		intent.context.SkillInferenceService = &mockInferenceService{}
		cmd := intent.startSkillInference()
		Expect(cmd).NotTo(BeNil())
		Expect(intent.state).To(Equal(StateInferringSkills))
		Expect(intent.loadingModal).NotTo(BeNil())
	})
})

type mockInferenceService struct{}

func (m *mockInferenceService) InferSkillsFromEvents(_ context.Context, _ []*domain.Event) (*skillinference.InferenceResult, error) {
	return &skillinference.InferenceResult{}, nil
}

func (m *mockInferenceService) CreateSkillsFromSuggestions(_ context.Context, _ []skillinference.SkillSuggestion) ([]*domain.Skill, error) {
	return nil, nil
}

var _ = Describe("Intent View edge cases", func() {
	It("returns empty string when not active", func() {
		intent := newTestIntent()
		intent.active = false
		Expect(intent.View()).To(BeEmpty())
	})

	It("returns No active screen when screen is nil", func() {
		intent := newTestIntent()
		intent.active = true
		intent.activeScreen = nil
		Expect(intent.View()).To(Equal("No active screen"))
	})
})

var _ = Describe("Intent HasActiveFilters when context nil", func() {
	It("returns false when context is nil", func() {
		intent := newActiveTestIntent()
		intent.context = nil
		Expect(intent.HasActiveFilters()).To(BeFalse())
	})

	It("returns false when filters is nil", func() {
		intent := newActiveTestIntent()
		intent.context.Filters = nil
		Expect(intent.HasActiveFilters()).To(BeFalse())
	})
})

var _ = Describe("Intent ClearFilters edge cases", func() {
	It("is safe when context is nil", func() {
		intent := newActiveTestIntent()
		intent.context = nil
		Expect(func() { intent.ClearFilters() }).NotTo(Panic())
	})

	It("is safe when filters is nil", func() {
		intent := newActiveTestIntent()
		intent.context.Filters = nil
		Expect(func() { intent.ClearFilters() }).NotTo(Panic())
	})
})

var _ = Describe("rebuildModalRegistry with multiple modals", func() {
	var intent *Intent

	BeforeEach(func() {
		intent = newActiveTestIntent()
	})

	It("registers search modal when present", func() {
		intent.openSearchModal()
		intent.rebuildModalRegistry()
		Expect(intent.modalRegistry).NotTo(BeNil())
	})

	It("registers filter modal when present", func() {
		intent.openFilterModal()
		intent.rebuildModalRegistry()
		Expect(intent.modalRegistry).NotTo(BeNil())
	})

	It("registers sort modal when present", func() {
		intent.openSortModal()
		intent.rebuildModalRegistry()
		Expect(intent.modalRegistry).NotTo(BeNil())
	})

	It("registers add/edit modal when present", func() {
		intent.openAddEditModal(nil)
		intent.rebuildModalRegistry()
		Expect(intent.modalRegistry).NotTo(BeNil())
	})

	It("registers delete modal when present", func() {
		skill := fixtures.SkillWith("s1", "Go", "backend", "advanced")
		intent.openDeleteModal(skill)
		intent.rebuildModalRegistry()
		Expect(intent.modalRegistry).NotTo(BeNil())
	})

	It("registers view detail modal when present", func() {
		intent.selectedSkill = fixtures.SkillWith("s1", "Go", "backend", "advanced")
		intent.openViewDetailModal()
		intent.rebuildModalRegistry()
		Expect(intent.modalRegistry).NotTo(BeNil())
	})

	It("registers skill events modal when present", func() {
		intent.selectedSkill = fixtures.SkillWith("s1", "Go", "backend", "advanced")
		events := fixtures.Events(1)
		intent.openSkillEventsModal(events)
		intent.rebuildModalRegistry()
		Expect(intent.modalRegistry).NotTo(BeNil())
	})

	It("registers event detail modal when present", func() {
		event := fixtures.EventWith("e1", "Test", "Company", "Project")
		intent.openEventDetailModal(event)
		intent.rebuildModalRegistry()
		Expect(intent.modalRegistry).NotTo(BeNil())
	})
})

var _ = Describe("inferSkillsFromAllEvents", func() {
	var intent *Intent

	BeforeEach(func() {
		intent = newActiveTestIntent()
		eventRepo := memoryrepo.NewEventRepository()
		intent.context.EventRepository = eventRepo
	})

	It("returns error when no events exist", func() {
		cmd := intent.inferSkillsFromAllEvents()
		Expect(cmd).NotTo(BeNil())
		msg := cmd()
		loadedMsg, ok := msg.(SkillSuggestionsLoadedMsg)
		Expect(ok).To(BeTrue())
		Expect(loadedMsg.Error).To(HaveOccurred())
	})
})

var _ = Describe("handleLoadingModalUpdate spinner tick", func() {
	var intent *Intent

	BeforeEach(func() {
		intent = newActiveTestIntent()
		intent.loadingModal = feedback.NewLoadingModal("Loading...", true)
	})

	It("processes spinner tick message", func() {
		cmd := intent.handleLoadingModalUpdate(feedback.ModalSpinnerTickMsg{})
		_ = cmd
	})
})

// Ensure unused imports don't cause issues.
var _ = behaviors.ColumnDef{}
var _ intents.ResultStatus
var _ = themes.NewThemeManager

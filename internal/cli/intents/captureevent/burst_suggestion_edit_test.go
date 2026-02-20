package captureevent

import (
	"context"

	"github.com/baphled/kariya/internal/domain/career"
	burstfact "github.com/baphled/kariya/internal/service/career/burstfact"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("BurstSuggestionModelNew edit flow", func() {
	var (
		model       *BurstSuggestionModelNew
		suggestions []burstfact.BurstSuggestion
	)

	BeforeEach(func() {
		suggestions = []burstfact.BurstSuggestion{
			{Name: "API Work", Description: "REST API development", EventIDs: []string{"e1", "e2"}, ConfidenceScore: 0.9},
			{Name: "CI/CD Work", Description: "Pipeline setup", EventIDs: []string{"e3"}, ConfidenceScore: 0.7},
		}
		model = NewBurstSuggestionModelNew(context.Background(), nil, suggestions)
	})

	Describe("startEdit", func() {
		It("sets editing to true", func() {
			model.startEdit()
			Expect(model.IsEditing()).To(BeTrue())
		})

		It("creates an edit form", func() {
			model.startEdit()
			Expect(model.editForm).NotTo(BeNil())
		})

		It("populates editFormData from current suggestion", func() {
			model.startEdit()
			Expect(model.editFormData).NotTo(BeNil())
			Expect(model.editFormData.Name).To(Equal("API Work"))
			Expect(model.editFormData.Description).To(Equal("REST API development"))
		})

		It("uses previously edited name if it exists", func() {
			model.SetEditedName(0, "Edited Name")
			model.startEdit()
			Expect(model.editFormData.Name).To(Equal("Edited Name"))
		})

		It("uses previously edited description if it exists", func() {
			model.SetEditedDescription(0, "Edited Desc")
			model.startEdit()
			Expect(model.editFormData.Description).To(Equal("Edited Desc"))
		})

		It("does nothing when suggestions is empty", func() {
			empty := NewBurstSuggestionModelNew(context.Background(), nil, []burstfact.BurstSuggestion{})
			result, cmd := empty.startEdit()
			Expect(result).NotTo(BeNil())
			Expect(cmd).To(BeNil())
			Expect(empty.IsEditing()).To(BeFalse())
		})
	})

	Describe("handleEditKeyMsg", func() {
		Context("when editForm is nil", func() {
			It("resets editing to false and returns no command", func() {
				model.editing = true
				model.editForm = nil
				result, cmd := model.handleEditKeyMsg(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
				Expect(result).NotTo(BeNil())
				Expect(cmd).To(BeNil())
				Expect(model.editing).To(BeFalse())
			})
		})

		Context("when editForm is active", func() {
			BeforeEach(func() {
				model.startEdit()
			})

			It("forwards key messages to the form", func() {
				Expect(func() {
					model.handleEditKeyMsg(tea.KeyMsg{Type: tea.KeyTab})
				}).NotTo(Panic())
			})

			It("returns non-nil model", func() {
				result, _ := model.handleEditKeyMsg(tea.KeyMsg{Type: tea.KeyTab})
				Expect(result).NotTo(BeNil())
			})
		})

		Context("when e key is pressed in browse mode", func() {
			It("activates edit mode via Update", func() {
				Expect(model.IsEditing()).To(BeFalse())
				model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
				Expect(model.IsEditing()).To(BeTrue())
			})
		})
	})

	Describe("saveEdits", func() {
		BeforeEach(func() {
			model.startEdit()
		})

		It("saves edited name to editedNames map", func() {
			model.editFormData.Name = "New Name"
			model.saveEdits()
			Expect(model.editedNames[0]).To(Equal("New Name"))
		})

		It("saves edited description to editedDescs map", func() {
			model.editFormData.Description = "New Desc"
			model.saveEdits()
			Expect(model.editedDescs[0]).To(Equal("New Desc"))
		})

		It("resets editing to false", func() {
			model.saveEdits()
			Expect(model.IsEditing()).To(BeFalse())
		})

		It("clears editForm", func() {
			model.saveEdits()
			Expect(model.editForm).To(BeNil())
		})

		It("clears editFormData", func() {
			model.saveEdits()
			Expect(model.editFormData).To(BeNil())
		})

		It("does not panic when editFormData is nil", func() {
			model.editFormData = nil
			Expect(func() { model.saveEdits() }).NotTo(Panic())
		})
	})

	Describe("renderEditView", func() {
		It("renders the edit form header", func() {
			model.startEdit()
			view := model.renderEditView()
			Expect(view).To(ContainSubstring("Edit Burst Name"))
		})

		It("renders a help text line", func() {
			model.startEdit()
			view := model.renderEditView()
			Expect(view).To(ContainSubstring("Tab"))
		})

		It("renders error message when editForm is nil", func() {
			model.editing = true
			model.editForm = nil
			view := model.renderEditView()
			Expect(view).To(ContainSubstring("not initialized"))
		})
	})

	Describe("View in edit mode", func() {
		It("calls renderEditView when editing is true", func() {
			model.startEdit()
			view := model.View()
			Expect(view).To(ContainSubstring("Edit Burst Name"))
		})
	})

	Describe("confirmCurrent — last suggestion path", func() {
		It("returns a non-nil command when confirming the final suggestion", func() {
			single := []burstfact.BurstSuggestion{
				{Name: "Solo Burst", Description: "Only one", EventIDs: []string{"e1"}, ConfidenceScore: 0.8},
			}
			m := NewBurstSuggestionModelNew(context.Background(), nil, single)
			_, cmd := m.confirmCurrent()
			Expect(cmd).NotTo(BeNil())
		})

		It("adds the burst to confirmed when confirming the final suggestion", func() {
			single := []burstfact.BurstSuggestion{
				{Name: "Solo Burst", Description: "Only one", EventIDs: []string{"e1"}},
			}
			m := NewBurstSuggestionModelNew(context.Background(), nil, single)
			m.confirmCurrent()
			Expect(m.GetConfirmed()).To(HaveLen(1))
			Expect(m.GetConfirmed()[0].Name).To(Equal("Solo Burst"))
		})

		It("marks the model as done after confirming the final suggestion", func() {
			single := []burstfact.BurstSuggestion{
				{Name: "Solo Burst", Description: "Only one", EventIDs: []string{"e1"}},
			}
			m := NewBurstSuggestionModelNew(context.Background(), nil, single)
			m.confirmCurrent()
			Expect(m.IsDone()).To(BeTrue())
		})

		It("uses edited name when confirming", func() {
			model.SetEditedName(0, "Custom Name")
			model.confirmCurrent()
			Expect(model.GetConfirmed()[0].Name).To(Equal("Custom Name"))
		})

		It("uses edited description when confirming", func() {
			model.SetEditedDescription(0, "Custom Desc")
			model.confirmCurrent()
			Expect(model.GetConfirmed()[0].Description).To(Equal("Custom Desc"))
		})

		It("increments cursor and adds to confirmed when not on last suggestion", func() {
			model.confirmCurrent()
			Expect(model.currentIdx).To(Equal(1))
			Expect(model.GetConfirmed()).To(HaveLen(1))
		})
	})

	Describe("rejectCurrent — last suggestion path", func() {
		It("returns a non-nil command when rejecting the final suggestion", func() {
			single := []burstfact.BurstSuggestion{
				{Name: "Solo Burst", Description: "Only one", EventIDs: []string{"e1"}},
			}
			m := NewBurstSuggestionModelNew(context.Background(), nil, single)
			_, cmd := m.rejectCurrent()
			Expect(cmd).NotTo(BeNil())
		})

		It("increments cursor when not on last suggestion", func() {
			model.rejectCurrent()
			Expect(model.currentIdx).To(Equal(1))
		})

		It("marks the model as done when final suggestion is rejected", func() {
			single := []burstfact.BurstSuggestion{
				{Name: "Last", Description: "desc", EventIDs: []string{"e1"}},
			}
			m := NewBurstSuggestionModelNew(context.Background(), nil, single)
			Expect(m.IsDone()).To(BeFalse())
			m.rejectCurrent()
			Expect(m.IsDone()).To(BeTrue())
		})

		It("adds the suggestion to rejected", func() {
			single := []burstfact.BurstSuggestion{
				{Name: "Last", Description: "desc", EventIDs: []string{"e1"}},
			}
			m := NewBurstSuggestionModelNew(context.Background(), nil, single)
			m.rejectCurrent()
			Expect(m.GetRejected()).To(HaveLen(1))
		})
	})

	Describe("renderRelatedEvents", func() {
		Context("with no event IDs", func() {
			It("renders no-related-events message", func() {
				s := burstfact.BurstSuggestion{Name: "Test", EventIDs: []string{}}
				view := model.renderRelatedEvents(s)
				Expect(view).To(ContainSubstring("No related events"))
			})
		})

		Context("with cached events", func() {
			It("renders cached event text", func() {
				cached := []*career.Event{
					fixtures.EventWith("e1", "Led payment API migration", "", ""),
				}
				model.relatedEvents[0] = cached
				s := burstfact.BurstSuggestion{Name: "Test", EventIDs: []string{"e1"}}
				view := model.renderRelatedEvents(s)
				Expect(view).To(ContainSubstring("Led payment API migration"))
			})

			It("truncates long event text to 57 characters plus ellipsis", func() {
				long := "A very long event description that exceeds sixty characters definitely"
				cached := []*career.Event{
					fixtures.EventWith("e1", long, "", ""),
				}
				model.relatedEvents[0] = cached
				s := burstfact.BurstSuggestion{Name: "Test", EventIDs: []string{"e1"}}
				view := model.renderRelatedEvents(s)
				Expect(view).To(ContainSubstring("..."))
			})

			It("shows '... and N more' for more than 3 events", func() {
				cached := []*career.Event{
					fixtures.EventWith("e1", "Event one content here displayed", "", ""),
					fixtures.EventWith("e2", "Event two content here displayed", "", ""),
					fixtures.EventWith("e3", "Event three content here displayed", "", ""),
					fixtures.EventWith("e4", "Event four content here displayed", "", ""),
				}
				model.relatedEvents[0] = cached
				s := burstfact.BurstSuggestion{Name: "Test", EventIDs: []string{"e1", "e2", "e3", "e4"}}
				view := model.renderRelatedEvents(s)
				Expect(view).To(ContainSubstring("and 1 more"))
			})
		})

		Context("with events pre-cached as empty slice", func() {
			It("renders the related events header with no event lines", func() {
				m := NewBurstSuggestionModelNew(context.Background(), nil, suggestions)
				m.relatedEvents[0] = []*career.Event{}
				s := burstfact.BurstSuggestion{Name: "Test", EventIDs: []string{"unknown-id"}}
				view := m.renderRelatedEvents(s)
				Expect(view).To(ContainSubstring("Related Events:"))
			})
		})
	})

	Describe("renderConfidenceScore", func() {
		It("renders BoxSuccess style for high confidence", func() {
			view := model.renderConfidenceScore(0.9)
			Expect(view).To(ContainSubstring("90%"))
		})

		It("renders BoxInfo style for medium confidence", func() {
			view := model.renderConfidenceScore(0.65)
			Expect(view).NotTo(BeEmpty())
		})

		It("renders BoxWarning style for low-medium confidence", func() {
			view := model.renderConfidenceScore(0.45)
			Expect(view).NotTo(BeEmpty())
		})

		It("renders BoxDestructive style for very low confidence", func() {
			view := model.renderConfidenceScore(0.2)
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("renderBurstDetails", func() {
		It("renders empty name warning when name is blank", func() {
			s := burstfact.BurstSuggestion{Name: "", Description: "desc", EventIDs: []string{"e1"}}
			view := model.renderBurstDetails(s)
			Expect(view).To(ContainSubstring("not set"))
		})

		It("renders the name when set", func() {
			s := burstfact.BurstSuggestion{Name: "My Burst", Description: "desc", EventIDs: []string{"e1"}}
			view := model.renderBurstDetails(s)
			Expect(view).To(ContainSubstring("My Burst"))
		})

		It("omits description section when description is empty", func() {
			s := burstfact.BurstSuggestion{Name: "Burst", Description: "", EventIDs: []string{"e1"}}
			view := model.renderBurstDetails(s)
			Expect(view).NotTo(ContainSubstring("Description:"))
		})

		It("renders description when present", func() {
			s := burstfact.BurstSuggestion{Name: "Burst", Description: "Some desc", EventIDs: []string{"e1"}}
			view := model.renderBurstDetails(s)
			Expect(view).To(ContainSubstring("Some desc"))
		})

		It("uses edited name when one exists for current index", func() {
			model.SetEditedName(0, "Edited Burst Name")
			s := model.suggestions[0]
			view := model.renderBurstDetails(s)
			Expect(view).To(ContainSubstring("Edited Burst Name"))
		})
	})

	Describe("createBurstFromSuggestion — empty name fallback", func() {
		It("generates a name from event count when name is empty", func() {
			s := burstfact.BurstSuggestion{Name: "", EventIDs: []string{"e1", "e2", "e3"}}
			burst := model.createBurstFromSuggestion(s)
			Expect(burst.Name).To(ContainSubstring("3 events"))
		})
	})
})

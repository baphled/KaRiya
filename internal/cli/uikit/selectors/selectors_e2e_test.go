package selectors_test

import (
	"github.com/baphled/kariya/internal/cli/uikit/selectors"
	domain "github.com/baphled/kariya/internal/domain/career"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Selectors E2E", func() {
	Describe("TagSelector - Event Tagging Workflow", func() {
		var selector *selectors.TagSelector

		BeforeEach(func() {
			selector = selectors.NewTagSelector()
		})

		It("should support complete event tagging workflow", func() {
			// User creates a new event and wants to add tags.

			// Step 1: View available tags.
			available := selector.AvailableTags()
			Expect(available).To(HaveLen(8))
			Expect(available).To(ContainElements("technical", "leadership", "achievement"))

			// Step 2: Select relevant tags for a technical achievement.
			Expect(selector.SelectTag("technical")).To(Succeed())
			Expect(selector.SelectTag("achievement")).To(Succeed())

			// Step 3: Verify selections.
			Expect(selector.SelectedTags()).To(HaveLen(2))
			Expect(selector.IsSelected("technical")).To(BeTrue())
			Expect(selector.IsSelected("achievement")).To(BeTrue())

			// Step 4: User realizes they also want to add leadership.
			Expect(selector.SelectTag("leadership")).To(Succeed())
			Expect(selector.SelectedTags()).To(HaveLen(3))

			// Step 5: User removes technical (wrong choice).
			Expect(selector.DeselectTag("technical")).To(Succeed())
			Expect(selector.SelectedTags()).To(HaveLen(2))
			Expect(selector.IsSelected("technical")).To(BeFalse())

			// Step 6: Final selection is consistent.
			final := selector.SelectedTags()
			Expect(final).To(Equal([]string{"achievement", "leadership"}))
		})

		It("should support editing existing event tags", func() {
			// Simulate loading tags from an existing event.
			existingTags := []string{"technical", "project"}
			selector.SetSelectedTags(existingTags)

			// Verify loaded state.
			Expect(selector.SelectedTags()).To(HaveLen(2))
			Expect(selector.IsSelected("technical")).To(BeTrue())
			Expect(selector.IsSelected("project")).To(BeTrue())

			// User modifies: remove project, add achievement.
			Expect(selector.DeselectTag("project")).To(Succeed())
			Expect(selector.SelectTag("achievement")).To(Succeed())

			// Final state.
			Expect(selector.SelectedTags()).To(Equal([]string{"achievement", "technical"}))
		})

		It("should support tag search/filter workflow", func() {
			// User types "le" to filter tags.
			filtered := selector.FilterTags("le")
			Expect(filtered).To(HaveLen(1))
			Expect(filtered).To(ContainElement("leadership"))

			// User clears filter and types "p" - matches product and project.
			filtered = selector.FilterTags("p")
			Expect(filtered).To(HaveLen(2))
			Expect(filtered).To(ContainElements("product", "project"))

			// User types "pro" to narrow down.
			filtered = selector.FilterTags("pro")
			Expect(filtered).To(HaveLen(2))
			Expect(filtered).To(ContainElements("product", "project"))

			// User types "proj" to get just project.
			filtered = selector.FilterTags("proj")
			Expect(filtered).To(HaveLen(1))
			Expect(filtered).To(ContainElement("project"))

			// User clears filter - sees all.
			filtered = selector.FilterTags("")
			Expect(filtered).To(HaveLen(8))
		})

		It("should enforce business rules", func() {
			// Cannot select invalid tag.
			err := selector.SelectTag("invalid-tag")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("not a valid tag"))

			// Cannot select same tag twice.
			Expect(selector.SelectTag("technical")).To(Succeed())
			err = selector.SelectTag("technical")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("already selected"))

			// Cannot deselect unselected tag.
			err = selector.DeselectTag("leadership")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("not selected"))
		})
	})

	Describe("CategorySelector - Event Categorization Workflow", func() {
		var selector *selectors.CategorySelector

		BeforeEach(func() {
			selector = selectors.NewCategorySelector()
		})

		It("should support complete categorization workflow", func() {
			// User categorizes a career event.

			// Step 1: View available categories.
			available := selector.AvailableCategories()
			Expect(available).To(HaveLen(6))
			Expect(available).To(ContainElements("technical", "leadership", "product"))

			// Step 2: This is a technical leadership event.
			Expect(selector.SelectCategory("technical")).To(Succeed())
			Expect(selector.SelectCategory("leadership")).To(Succeed())

			// Step 3: Verify.
			Expect(selector.SelectedCategories()).To(HaveLen(2))

			// Step 4: User toggles mentoring on/off.
			Expect(selector.ToggleCategory("mentoring")).To(Succeed())
			Expect(selector.IsSelected("mentoring")).To(BeTrue())

			Expect(selector.ToggleCategory("mentoring")).To(Succeed())
			Expect(selector.IsSelected("mentoring")).To(BeFalse())

			// Step 5: Clear and start over.
			selector.Clear()
			Expect(selector.SelectedCategories()).To(BeEmpty())
		})

		It("should support bulk category setting", func() {
			// Load categories from existing event.
			err := selector.SetSelected([]string{"technical", "consulting"})
			Expect(err).NotTo(HaveOccurred())

			Expect(selector.SelectedCategories()).To(HaveLen(2))
			Expect(selector.IsSelected("technical")).To(BeTrue())
			Expect(selector.IsSelected("consulting")).To(BeTrue())

			// Replace with different categories.
			err = selector.SetSelected([]string{"leadership", "mentoring"})
			Expect(err).NotTo(HaveOccurred())

			Expect(selector.SelectedCategories()).To(HaveLen(2))
			Expect(selector.IsSelected("technical")).To(BeFalse())
			Expect(selector.IsSelected("leadership")).To(BeTrue())
		})

		It("should handle case-insensitive input", func() {
			// User might type in different cases.
			Expect(selector.SelectCategory("Technical")).To(Succeed())
			Expect(selector.SelectCategory("LEADERSHIP")).To(Succeed())

			Expect(selector.IsSelected("technical")).To(BeTrue())
			Expect(selector.IsSelected("leadership")).To(BeTrue())
		})

		It("should provide category descriptions for UI", func() {
			desc := selectors.GetCategoryDescription("technical")
			Expect(desc).To(ContainSubstring("Technical"))

			desc = selectors.GetCategoryDescription("leadership")
			Expect(desc).To(ContainSubstring("Leadership"))

			// Unknown category returns the name.
			desc = selectors.GetCategoryDescription("unknown")
			Expect(desc).To(Equal("unknown"))
		})
	})

	Describe("SkillSelector - Event Skill Association Workflow", func() {
		var (
			selector        *selectors.SkillSelector
			availableSkills []*domain.Skill
		)

		BeforeEach(func() {
			availableSkills = []*domain.Skill{
				{ID: "sk-go", Name: "Go", Category: "backend"},
				{ID: "sk-react", Name: "React", Category: "frontend"},
				{ID: "sk-postgres", Name: "PostgreSQL", Category: "database"},
				{ID: "sk-k8s", Name: "Kubernetes", Category: "devops"},
				{ID: "sk-docker", Name: "Docker", Category: "devops"},
				{ID: "sk-aws", Name: "AWS", Category: "cloud"},
			}
			selector = selectors.NewSkillSelector(availableSkills)
		})

		It("should support complete skill association workflow", func() {
			// User is associating skills with a career event.

			// Step 1: View available skills (sorted by name).
			skills := selector.AvailableSkills()
			Expect(skills).To(HaveLen(6))
			Expect(skills[0].Name).To(Equal("AWS"))
			Expect(skills[1].Name).To(Equal("Docker"))

			// Step 2: Select skills used in this event.
			Expect(selector.SelectSkill("sk-go")).To(Succeed())
			Expect(selector.SelectSkill("sk-k8s")).To(Succeed())
			Expect(selector.SelectSkill("sk-docker")).To(Succeed())

			// Step 3: Verify selections.
			selected := selector.SelectedSkillIDs()
			Expect(selected).To(HaveLen(3))
			Expect(selected).To(ContainElements("sk-go", "sk-k8s", "sk-docker"))

			// Step 4: Get skill details for display.
			skill, err := selector.GetSkillByID("sk-go")
			Expect(err).NotTo(HaveOccurred())
			Expect(skill.Name).To(Equal("Go"))
			Expect(skill.Category).To(Equal("backend"))

			// Step 5: Toggle a skill.
			Expect(selector.ToggleSkill("sk-docker")).To(Succeed())
			Expect(selector.IsSelected("sk-docker")).To(BeFalse())

			Expect(selector.ToggleSkill("sk-aws")).To(Succeed())
			Expect(selector.IsSelected("sk-aws")).To(BeTrue())
		})

		It("should support editing existing event skills", func() {
			// Load skills from existing event.
			existingSkillIDs := []string{"sk-go", "sk-postgres"}
			selector.SetSelectedSkills(existingSkillIDs)

			Expect(selector.SelectedSkillIDs()).To(HaveLen(2))
			Expect(selector.IsSelected("sk-go")).To(BeTrue())
			Expect(selector.IsSelected("sk-postgres")).To(BeTrue())

			// Modify: remove postgres, add react.
			Expect(selector.DeselectSkill("sk-postgres")).To(Succeed())
			Expect(selector.SelectSkill("sk-react")).To(Succeed())

			selected := selector.SelectedSkillIDs()
			Expect(selected).To(ContainElements("sk-go", "sk-react"))
			Expect(selected).NotTo(ContainElement("sk-postgres"))
		})

		It("should support skill search workflow", func() {
			// User types to filter skills.

			// Filter by "Do" - should find Docker.
			filtered := selector.FilterSkills("Do")
			Expect(filtered).To(HaveLen(1))
			Expect(filtered[0].Name).To(Equal("Docker"))

			// Filter by "k" - should find Kubernetes.
			filtered = selector.FilterSkills("k")
			Expect(filtered).To(HaveLen(1))
			Expect(filtered[0].Name).To(Equal("Kubernetes"))

			// Empty filter shows all.
			filtered = selector.FilterSkills("")
			Expect(filtered).To(HaveLen(6))
		})

		It("should handle invalid skill IDs gracefully", func() {
			// Cannot select non-existent skill.
			err := selector.SelectSkill("sk-nonexistent")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("not a valid skill"))

			// Cannot get non-existent skill.
			_, err = selector.GetSkillByID("sk-nonexistent")
			Expect(err).To(HaveOccurred())

			// SetSelectedSkills ignores invalid IDs.
			selector.SetSelectedSkills([]string{"sk-go", "sk-invalid", "sk-react"})
			Expect(selector.SelectedSkillIDs()).To(HaveLen(2))
			Expect(selector.SelectedSkillIDs()).To(ContainElements("sk-go", "sk-react"))
		})

		It("should support reset workflow", func() {
			selector.SetSelectedSkills([]string{"sk-go", "sk-react", "sk-k8s"})
			Expect(selector.SelectedSkillIDs()).To(HaveLen(3))

			selector.Reset()
			Expect(selector.SelectedSkillIDs()).To(BeEmpty())
		})
	})

	Describe("AudienceRelevanceSelector - Fact Audience Workflow", func() {
		var selector *selectors.AudienceRelevanceSelector

		BeforeEach(func() {
			selector = selectors.NewAudienceRelevanceSelector()
		})

		It("should support complete audience selection workflow", func() {
			// User is marking which audiences a fact is relevant for.

			// Step 1: Initially no audiences selected.
			Expect(selector.GetSelected()).To(BeEmpty())

			// Step 2: This fact is relevant for hiring managers and peers.
			selector.SetSelected([]string{"hiring_manager", "peer"})

			// Step 3: Verify.
			selected := selector.GetSelected()
			Expect(selected).To(HaveLen(2))
			Expect(selected).To(ContainElements("hiring_manager", "peer"))

			// Step 4: Check individual selections.
			Expect(selector.IsSelected("hiring_manager")).To(BeTrue())
			Expect(selector.IsSelected("recruiter")).To(BeFalse())
			Expect(selector.IsSelected("peer")).To(BeTrue())

			// Step 5: Render for display.
			rendered := selector.Render()
			Expect(rendered).To(ContainSubstring("[✓ hiring_manager]"))
			Expect(rendered).To(ContainSubstring("[ recruiter ]"))
			Expect(rendered).To(ContainSubstring("[✓ peer]"))
		})

		It("should support editing existing fact audiences", func() {
			// Load from existing fact.
			selector.SetSelected([]string{"recruiter"})
			Expect(selector.GetSelected()).To(Equal([]string{"recruiter"}))

			// Change to different audiences.
			selector.SetSelected([]string{"hiring_manager", "recruiter", "peer"})
			Expect(selector.GetSelected()).To(HaveLen(3))
		})

		It("should maintain consistent ordering", func() {
			// Set in different order.
			selector.SetSelected([]string{"peer", "hiring_manager"})

			// GetSelected returns in options order.
			selected := selector.GetSelected()
			Expect(selected[0]).To(Equal("hiring_manager"))
			Expect(selected[1]).To(Equal("peer"))
		})

		It("should handle edge cases", func() {
			// Empty selection.
			selector.SetSelected([]string{})
			Expect(selector.GetSelected()).To(BeEmpty())

			// Nil selection.
			selector.SetSelected(nil)
			Expect(selector.GetSelected()).To(BeEmpty())

			// Invalid audiences are stored but not returned (filtered by options).
			selector.SetSelected([]string{"hiring_manager", "invalid_audience"})
			selected := selector.GetSelected()
			Expect(selected).To(ContainElement("hiring_manager"))
			Expect(selected).NotTo(ContainElement("invalid_audience"))
		})

		It("should render correctly with no selections", func() {
			rendered := selector.Render()
			Expect(rendered).To(ContainSubstring("[ hiring_manager ]"))
			Expect(rendered).To(ContainSubstring("[ recruiter ]"))
			Expect(rendered).To(ContainSubstring("[ peer ]"))
			Expect(rendered).NotTo(ContainSubstring("✓"))
		})

		It("should render correctly with all selections", func() {
			selector.SetSelected([]string{"hiring_manager", "recruiter", "peer"})
			rendered := selector.Render()
			Expect(rendered).To(ContainSubstring("[✓ hiring_manager]"))
			Expect(rendered).To(ContainSubstring("[✓ recruiter]"))
			Expect(rendered).To(ContainSubstring("[✓ peer]"))
		})
	})

	Describe("Cross-Selector Integration", func() {
		It("should support event creation with multiple selectors", func() {
			// Simulate creating an event that uses all selectors.
			tagSelector := selectors.NewTagSelector()
			categorySelector := selectors.NewCategorySelector()
			skillSelector := selectors.NewSkillSelector([]*domain.Skill{
				{ID: "sk-go", Name: "Go", Category: "backend"},
				{ID: "sk-k8s", Name: "Kubernetes", Category: "devops"},
			})

			// User fills out event form.
			Expect(tagSelector.SelectTag("technical")).To(Succeed())
			Expect(tagSelector.SelectTag("achievement")).To(Succeed())

			Expect(categorySelector.SelectCategory("technical")).To(Succeed())

			Expect(skillSelector.SelectSkill("sk-go")).To(Succeed())
			Expect(skillSelector.SelectSkill("sk-k8s")).To(Succeed())

			// Collect all selections for event creation.
			eventTags := tagSelector.SelectedTags()
			eventCategories := categorySelector.SelectedCategories()
			eventSkillIDs := skillSelector.SelectedSkillIDs()

			Expect(eventTags).To(Equal([]string{"achievement", "technical"}))
			Expect(eventCategories).To(Equal([]string{"technical"}))
			Expect(eventSkillIDs).To(ContainElements("sk-go", "sk-k8s"))
		})
	})
})

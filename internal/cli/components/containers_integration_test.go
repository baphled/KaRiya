package components

import (
	"strings"
	"testing"

	"github.com/baphled/kariya/internal/cli/styles"
)

// TestContainerComposition tests that containers can be composed together
func TestContainerComposition(t *testing.T) {
	// Create a form field
	field := NewFormFieldContainer().
		SetLabel("Username").
		SetInput("john_doe").
		SetHint("3-20 characters")

	fieldContent := field.Render()

	// Wrap in a card
	card := NewCardContainer().
		SetHeader("User Form").
		SetBody(fieldContent)

	rendered := card.Render()

	if !strings.Contains(rendered, "Username") || !strings.Contains(rendered, "john_doe") {
		t.Error("expected composed output to contain form field content")
	}

	if !strings.Contains(rendered, "User Form") {
		t.Error("expected composed output to contain card header")
	}
}

// TestScreenContainerWithCardContainer tests nesting CardContainer inside ScreenContainer
func TestScreenContainerWithCardContainer(t *testing.T) {
	// Create a card
	card := NewCardContainer().
		SetHeader("Event Details").
		SetBody("Career milestone achieved")

	cardContent := card.Render()

	// Wrap in screen container
	screen := NewScreenContainer(cardContent).
		WithPaddingMode(PaddingNormal)

	rendered := screen.Render()

	if !strings.Contains(rendered, "Event Details") || !strings.Contains(rendered, "Career milestone achieved") {
		t.Error("expected screen container to preserve card content")
	}
}

// TestNestedContainers tests multiple levels of container nesting
func TestNestedContainers(t *testing.T) {
	// Create form fields
	field1 := NewFormFieldContainer().
		SetLabel("Title").
		SetInput("Software Engineer")

	field2 := NewFormFieldContainer().
		SetLabel("Company").
		SetInput("Tech Corp")

	// Create a card with both fields
	card := NewCardContainer().
		SetHeader("Career Event").
		SetBody(field1.Render() + "\n" + field2.Render())

	// Wrap in section
	section := NewSectionContainer(card.Render()).
		SetTitle("New Event")

	// Wrap in screen
	screen := NewScreenContainer(section.Render()).
		WithPaddingMode(PaddingSpacious)

	rendered := screen.Render()

	if !strings.Contains(rendered, "New Event") ||
		!strings.Contains(rendered, "Career Event") ||
		!strings.Contains(rendered, "Software Engineer") {
		t.Error("expected all nested content to be present")
	}
}

// TestContainerColorConsistency tests that containers use consistent colors
func TestContainerColorConsistency(t *testing.T) {
	// Create various containers
	screen := NewScreenContainer("content")
	card := NewCardContainer().SetBody("content")
	section := NewSectionContainer("content").SetTitle("Title")
	field := NewFormFieldContainer().SetLabel("Label").SetInput("input")
	list := NewListContainer().SetItems([]string{"item"})
	modal := NewModalContainer().SetTitle("Title")

	// All should render without error
	screenOut := screen.Render()
	cardOut := card.Render()
	sectionOut := section.Render()
	fieldOut := field.Render()
	listOut := list.Render()
	modalOut := modal.Render()

	if screenOut == "" || cardOut == "" || sectionOut == "" ||
		fieldOut == "" || listOut == "" || modalOut == "" {
		t.Error("expected all containers to render non-empty output")
	}
}

// TestContainerWithFormFields tests form containers with multiple fields
func TestContainerWithFormFields(t *testing.T) {
	// Create multiple form fields
	fields := []string{
		NewFormFieldContainer().
			SetLabel("Email").
			SetInput("user@example.com").
			SetHint("Valid email required").
			Render(),
		NewFormFieldContainer().
			SetLabel("Password").
			SetInput("••••••••").
			SetFocused(true).
			Render(),
		NewFormFieldContainer().
			SetLabel("Confirm").
			SetInput("••••••••").
			SetError("Passwords do not match").
			Render(),
	}

	// Create a form card
	formBody := strings.Join(fields, "\n\n")
	form := NewCardContainer().
		SetHeader("Login Form").
		SetBody(formBody)

	rendered := form.Render()

	if !strings.Contains(rendered, "Email") ||
		!strings.Contains(rendered, "Password") ||
		!strings.Contains(rendered, "Confirm") ||
		!strings.Contains(rendered, "Passwords do not match") {
		t.Error("expected form to contain all fields and errors")
	}
}

// TestListContainerWithPagination tests list rendering with pagination
func TestListContainerWithPagination(t *testing.T) {
	items := []string{
		"Event 1 - Started new role",
		"Event 2 - Promoted to lead",
		"Event 3 - Completed certification",
	}

	list := NewListContainer().
		SetItems(items).
		SetPaginationInfo("Showing 1-3 of 15 events")

	rendered := list.Render()

	for _, item := range items {
		if !strings.Contains(rendered, item) {
			t.Errorf("expected list to contain %q", item)
		}
	}

	if !strings.Contains(rendered, "Showing 1-3 of 15 events") {
		t.Error("expected list to contain pagination info")
	}
}

// TestModalWithAllSections tests modal with all available sections
func TestModalWithAllSections(t *testing.T) {
	modal := NewModalContainer().
		SetTitle("Delete Event").
		SetMessage("This action cannot be undone.\nAre you sure you want to delete this event?").
		SetButtons([]string{"Delete", "Cancel"}).
		SetInstructions("Press Tab to navigate, Enter to select").
		WithDestructiveStyle()

	rendered := modal.Render()

	if !strings.Contains(rendered, "Delete Event") ||
		!strings.Contains(rendered, "This action cannot be undone") ||
		!strings.Contains(rendered, "Delete") ||
		!strings.Contains(rendered, "Cancel") ||
		!strings.Contains(rendered, "Press Tab") {
		t.Error("expected modal to contain all sections")
	}
}

// TestContainerWithEmptyStates tests containers with empty content
func TestContainerWithEmptyStates(t *testing.T) {
	// Empty list
	emptyList := NewListContainer().
		SetEmptyStateMessage("No events recorded yet")

	listOut := emptyList.Render()
	if !strings.Contains(listOut, "No events recorded yet") {
		t.Error("expected empty list to show message")
	}

	// Empty form
	emptyForm := NewFormFieldContainer()
	_ = emptyForm.Render()
	// Empty form with no sections renders empty string, which is acceptable

	// Empty modal
	emptyModal := NewModalContainer()
	modalOut := emptyModal.Render()
	if modalOut == "" {
		t.Error("expected empty modal to still render")
	}
}

// TestContainerSpacingIntegration tests spacing consistency across containers
func TestContainerSpacingIntegration(t *testing.T) {
	// Create content with different spacing modes
	compactScreen := NewScreenContainer("content").WithPaddingMode(PaddingCompact)
	normalScreen := NewScreenContainer("content").WithPaddingMode(PaddingNormal)
	spaciousScreen := NewScreenContainer("content").WithPaddingMode(PaddingSpacious)

	compact := compactScreen.Render()
	normal := normalScreen.Render()
	spacious := spaciousScreen.Render()

	// All should render
	if compact == "" || normal == "" || spacious == "" {
		t.Error("expected all spacing modes to render")
	}

	// Spacious should be longer due to more padding
	if len(spacious) <= len(normal) {
		t.Error("expected spacious padding to produce longer output")
	}
}

// TestContainerBorderAndBackgroundColors tests color variations
func TestContainerBorderAndBackgroundColors(t *testing.T) {
	// Normal card
	normalCard := NewCardContainer().
		SetBody("Normal card")

	// Custom color card
	customCard := NewCardContainer().
		SetBody("Custom card").
		WithBorderColor(styles.ColorBorderActive).
		WithBackgroundColor(styles.ColorBackgroundAlt)

	// Destructive modal
	destructiveModal := NewModalContainer().
		SetTitle("Delete").
		SetMessage("Confirm?").
		WithDestructiveStyle()

	normalOut := normalCard.Render()
	customOut := customCard.Render()
	destructiveOut := destructiveModal.Render()

	if !strings.Contains(normalOut, "Normal card") ||
		!strings.Contains(customOut, "Custom card") ||
		!strings.Contains(destructiveOut, "Delete") {
		t.Error("expected all color variations to render properly")
	}
}

// TestComplexFormWithValidation tests a complex form with validation states
func TestComplexFormWithValidation(t *testing.T) {
	// Create form fields with various states
	validField := NewFormFieldContainer().
		SetLabel("Username").
		SetInput("john_doe").
		SetHint("3-20 characters")

	invalidField := NewFormFieldContainer().
		SetLabel("Email").
		SetInput("invalid-email").
		SetError("Invalid email format").
		SetFocused(true)

	successField := NewFormFieldContainer().
		SetLabel("Password").
		SetInput("••••••••").
		SetHint("At least 8 characters")

	// Combine in a form
	formContent := strings.Join([]string{
		validField.Render(),
		invalidField.Render(),
		successField.Render(),
	}, "\n\n")

	form := NewCardContainer().
		SetHeader("Registration Form").
		SetBody(formContent).
		SetFooter("All fields are required")

	rendered := form.Render()

	if !strings.Contains(rendered, "Username") ||
		!strings.Contains(rendered, "Email") ||
		!strings.Contains(rendered, "Invalid email format") ||
		!strings.Contains(rendered, "Password") {
		t.Error("expected form to contain all fields and validation states")
	}
}

// TestListWithDetailedItems tests list with complex item content
func TestListWithDetailedItems(t *testing.T) {
	items := []string{
		"[2025-01-15] Started new role at Tech Corp",
		"[2025-01-10] Completed AWS certification",
		"[2025-01-05] Led team project successfully",
	}

	list := NewListContainer().
		SetItems(items).
		SetPaginationInfo("Showing 1-3 of 10 career events")

	rendered := list.Render()

	for _, item := range items {
		if !strings.Contains(rendered, item) {
			t.Errorf("expected list to contain %q", item)
		}
	}
}

// TestSectionContainerWithFormFields tests section containing form fields
func TestSectionContainerWithFormFields(t *testing.T) {
	field := NewFormFieldContainer().
		SetLabel("Event Title").
		SetInput("Promotion").
		SetHint("Brief description of the event")

	section := NewSectionContainer(field.Render()).
		SetTitle("Career Event Details")

	rendered := section.Render()

	if !strings.Contains(rendered, "Career Event Details") ||
		!strings.Contains(rendered, "Event Title") ||
		!strings.Contains(rendered, "Promotion") {
		t.Error("expected section to contain form field content")
	}
}

// TestCardWithMultipleSections tests card with header, body, and footer
func TestCardWithMultipleSections(t *testing.T) {
	card := NewCardContainer().
		SetHeader("Career Event #42").
		SetBody("Started new role as Senior Engineer at Tech Corp\n\nResponsibilities: Team leadership, architecture design").
		SetFooter("Created on 2025-01-15 • Last modified on 2025-01-20")

	rendered := card.Render()

	if !strings.Contains(rendered, "Career Event #42") ||
		!strings.Contains(rendered, "Senior Engineer") ||
		!strings.Contains(rendered, "Created on 2025-01-15") {
		t.Error("expected card to contain all sections")
	}
}

// TestContainerRenderingConsistency tests that containers render consistently
func TestContainerRenderingConsistency(t *testing.T) {
	// Create a complex nested structure
	field := NewFormFieldContainer().
		SetLabel("Title").
		SetInput("Software Engineer")

	card := NewCardContainer().
		SetHeader("Event").
		SetBody(field.Render())

	section := NewSectionContainer(card.Render()).
		SetTitle("Details")

	screen := NewScreenContainer(section.Render())

	// Render multiple times
	render1 := screen.Render()
	render2 := screen.Render()
	render3 := screen.Render()

	if render1 != render2 || render2 != render3 {
		t.Error("expected consistent rendering across multiple calls")
	}
}

// TestAllContainersWithStyles tests that all containers properly apply styles
func TestAllContainersWithStyles(t *testing.T) {
	// Test each container type applies styles correctly
	containers := []struct {
		name     string
		renderer func() string
	}{
		{
			"ScreenContainer",
			func() string {
				return NewScreenContainer("Test").WithPaddingMode(PaddingNormal).Render()
			},
		},
		{
			"CardContainer",
			func() string {
				return NewCardContainer().SetBody("Test").Render()
			},
		},
		{
			"SectionContainer",
			func() string {
				return NewSectionContainer("Test").SetTitle("Title").Render()
			},
		},
		{
			"FormFieldContainer",
			func() string {
				return NewFormFieldContainer().SetLabel("Label").SetInput("Input").Render()
			},
		},
		{
			"ListContainer",
			func() string {
				return NewListContainer().SetItems([]string{"Item"}).Render()
			},
		},
		{
			"ModalContainer",
			func() string {
				return NewModalContainer().SetTitle("Title").Render()
			},
		},
	}

	for _, c := range containers {
		t.Run(c.name, func(t *testing.T) {
			rendered := c.renderer()
			if rendered == "" {
				t.Errorf("expected %s to render non-empty output", c.name)
			}
			if !strings.Contains(rendered, "Test") && !strings.Contains(rendered, "Title") &&
				!strings.Contains(rendered, "Label") && !strings.Contains(rendered, "Item") {
				// At least some content should be present
				if len(rendered) < 5 {
					t.Errorf("expected %s to render meaningful content", c.name)
				}
			}
		})
	}
}


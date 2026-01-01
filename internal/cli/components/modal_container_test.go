package components

import (
	"strings"
	"testing"
)

func TestNewModalContainer(t *testing.T) {
	mc := NewModalContainer()

	if mc.isDestructive {
		t.Error("expected isDestructive to be false by default")
	}

	if mc.hasTitle || mc.hasMessage || mc.hasButtons || mc.hasInstructions {
		t.Error("expected new modal to have no sections by default")
	}
}

func TestModalContainerSetTitle(t *testing.T) {
	title := "Confirm Action"
	mc := NewModalContainer().SetTitle(title)

	if mc.title != title {
		t.Errorf("expected title %q, got %q", title, mc.title)
	}

	if !mc.hasTitle {
		t.Error("expected hasTitle to be true")
	}
}

func TestModalContainerSetMessage(t *testing.T) {
	message := "Are you sure you want to continue?"
	mc := NewModalContainer().SetMessage(message)

	if mc.message != message {
		t.Errorf("expected message %q, got %q", message, mc.message)
	}

	if !mc.hasMessage {
		t.Error("expected hasMessage to be true")
	}
}

func TestModalContainerSetButtons(t *testing.T) {
	buttons := []string{"Yes", "No"}
	mc := NewModalContainer().SetButtons(buttons)

	if len(mc.buttons) != len(buttons) {
		t.Errorf("expected %d buttons, got %d", len(buttons), len(mc.buttons))
	}

	for i, btn := range buttons {
		if mc.buttons[i] != btn {
			t.Errorf("expected button %q, got %q", btn, mc.buttons[i])
		}
	}

	if !mc.hasButtons {
		t.Error("expected hasButtons to be true")
	}
}

func TestModalContainerSetInstructions(t *testing.T) {
	instructions := "Press Tab to navigate, Enter to confirm"
	mc := NewModalContainer().SetInstructions(instructions)

	if mc.instructions != instructions {
		t.Errorf("expected instructions %q, got %q", instructions, mc.instructions)
	}

	if !mc.hasInstructions {
		t.Error("expected hasInstructions to be true")
	}
}

func TestModalContainerWithDestructiveStyle(t *testing.T) {
	mc := NewModalContainer().WithDestructiveStyle()

	if !mc.isDestructive {
		t.Error("expected isDestructive to be true")
	}
}

func TestModalContainerRenderWithTitle(t *testing.T) {
	title := "Delete Item"
	mc := NewModalContainer().SetTitle(title)

	rendered := mc.Render()

	if !strings.Contains(rendered, title) {
		t.Errorf("expected rendered output to contain title %q", title)
	}
}

func TestModalContainerRenderWithMessage(t *testing.T) {
	message := "This action cannot be undone"
	mc := NewModalContainer().SetMessage(message)

	rendered := mc.Render()

	if !strings.Contains(rendered, message) {
		t.Errorf("expected rendered output to contain message %q", message)
	}
}

func TestModalContainerRenderWithButtons(t *testing.T) {
	buttons := []string{"Delete", "Cancel"}
	mc := NewModalContainer().SetButtons(buttons)

	rendered := mc.Render()

	for _, btn := range buttons {
		if !strings.Contains(rendered, btn) {
			t.Errorf("expected rendered output to contain button %q", btn)
		}
	}
}

func TestModalContainerRenderWithInstructions(t *testing.T) {
	instructions := "Press y to confirm"
	mc := NewModalContainer().SetInstructions(instructions)

	rendered := mc.Render()

	if !strings.Contains(rendered, instructions) {
		t.Errorf("expected rendered output to contain instructions %q", instructions)
	}
}

func TestModalContainerRenderAllSections(t *testing.T) {
	title := "Confirm"
	message := "Are you sure?"
	buttons := []string{"Yes", "No"}
	instructions := "Use arrow keys"

	mc := NewModalContainer().
		SetTitle(title).
		SetMessage(message).
		SetButtons(buttons).
		SetInstructions(instructions)

	rendered := mc.Render()

	if !strings.Contains(rendered, title) {
		t.Errorf("expected rendered output to contain title %q", title)
	}

	if !strings.Contains(rendered, message) {
		t.Errorf("expected rendered output to contain message %q", message)
	}

	for _, btn := range buttons {
		if !strings.Contains(rendered, btn) {
			t.Errorf("expected rendered output to contain button %q", btn)
		}
	}

	if !strings.Contains(rendered, instructions) {
		t.Errorf("expected rendered output to contain instructions %q", instructions)
	}
}

func TestModalContainerBuilderChaining(t *testing.T) {
	title := "Confirm"
	message := "Are you sure?"
	buttons := []string{"Yes", "No"}
	instructions := "Press Enter"

	mc := NewModalContainer().
		SetTitle(title).
		SetMessage(message).
		SetButtons(buttons).
		SetInstructions(instructions).
		WithDestructiveStyle()

	if mc.title != title || mc.message != message || mc.instructions != instructions {
		t.Error("expected all sections to be set after chaining")
	}

	if !mc.isDestructive {
		t.Error("expected destructive style to be set")
	}
}

func TestModalContainerRenderDestructiveStyle(t *testing.T) {
	title := "Delete Permanently"
	message := "This cannot be undone"

	mc := NewModalContainer().
		SetTitle(title).
		SetMessage(message).
		WithDestructiveStyle()

	rendered := mc.Render()

	if !strings.Contains(rendered, title) || !strings.Contains(rendered, message) {
		t.Error("expected rendered output to contain title and message")
	}
}

func TestModalContainerRenderNormalStyle(t *testing.T) {
	title := "Confirm Action"
	message := "Continue?"

	mc := NewModalContainer().
		SetTitle(title).
		SetMessage(message)

	rendered := mc.Render()

	if !strings.Contains(rendered, title) || !strings.Contains(rendered, message) {
		t.Error("expected rendered output to contain title and message")
	}
}

func TestModalContainerMultipleButtons(t *testing.T) {
	buttons := []string{"Save", "Don't Save", "Cancel"}
	mc := NewModalContainer().SetButtons(buttons)

	rendered := mc.Render()

	for _, btn := range buttons {
		if !strings.Contains(rendered, btn) {
			t.Errorf("expected rendered output to contain button %q", btn)
		}
	}
}

func TestModalContainerSingleButton(t *testing.T) {
	buttons := []string{"OK"}
	mc := NewModalContainer().SetButtons(buttons)

	rendered := mc.Render()

	if !strings.Contains(rendered, "OK") {
		t.Error("expected rendered output to contain OK button")
	}
}

func TestModalContainerEmptyButtons(t *testing.T) {
	mc := NewModalContainer().SetButtons([]string{})

	rendered := mc.Render()

	// Should render without buttons
	if rendered == "" {
		t.Error("expected rendered output to be non-empty")
	}
}

func TestModalContainerMultilineMessage(t *testing.T) {
	message := "This is a long message\nthat spans multiple\nlines"
	mc := NewModalContainer().SetMessage(message)

	rendered := mc.Render()

	if !strings.Contains(rendered, "long message") || !strings.Contains(rendered, "spans multiple") {
		t.Error("expected rendered output to contain multiline message")
	}
}

func TestModalContainerRenderEmpty(t *testing.T) {
	mc := NewModalContainer()

	rendered := mc.Render()

	// Should still render a modal, even with no content
	if rendered == "" {
		t.Error("expected rendered output to be non-empty")
	}
}

func TestModalContainerSectionOrdering(t *testing.T) {
	// Sections should appear in order: title, message, buttons, instructions
	title := "TITLE"
	message := "MESSAGE"
	buttons := []string{"BUTTON"}
	instructions := "INSTRUCTIONS"

	mc := NewModalContainer().
		SetTitle(title).
		SetMessage(message).
		SetButtons(buttons).
		SetInstructions(instructions)

	rendered := mc.Render()

	titlePos := strings.Index(rendered, title)
	messagePos := strings.Index(rendered, message)
	buttonPos := strings.Index(rendered, "BUTTON")
	instructionsPos := strings.Index(rendered, instructions)

	if titlePos == -1 || messagePos == -1 || buttonPos == -1 || instructionsPos == -1 {
		t.Error("expected all sections to be present")
	}

	if titlePos > messagePos || messagePos > buttonPos || buttonPos > instructionsPos {
		t.Error("expected sections to appear in order: title, message, buttons, instructions")
	}
}

func TestModalContainerWithoutTitle(t *testing.T) {
	message := "Confirm?"
	buttons := []string{"Yes", "No"}

	mc := NewModalContainer().
		SetMessage(message).
		SetButtons(buttons)

	rendered := mc.Render()

	if !strings.Contains(rendered, message) || !strings.Contains(rendered, "Yes") {
		t.Error("expected rendered output to contain message and buttons")
	}
}

func TestModalContainerWithoutMessage(t *testing.T) {
	title := "Confirm"
	buttons := []string{"OK"}

	mc := NewModalContainer().
		SetTitle(title).
		SetButtons(buttons)

	rendered := mc.Render()

	if !strings.Contains(rendered, title) || !strings.Contains(rendered, "OK") {
		t.Error("expected rendered output to contain title and buttons")
	}
}

func TestModalContainerSpecialCharacters(t *testing.T) {
	title := "Delete [Item]"
	message := "Permanently remove {item}?"
	buttons := []string{"Delete @#$", "Cancel"}

	mc := NewModalContainer().
		SetTitle(title).
		SetMessage(message).
		SetButtons(buttons)

	rendered := mc.Render()

	if !strings.Contains(rendered, "[Item]") || !strings.Contains(rendered, "{item}") || !strings.Contains(rendered, "@#$") {
		t.Error("expected rendered output to contain special characters")
	}
}

func TestModalContainerRenderConsistency(t *testing.T) {
	title := "Confirm"
	message := "Are you sure?"

	mc := NewModalContainer().
		SetTitle(title).
		SetMessage(message)

	rendered1 := mc.Render()
	rendered2 := mc.Render()

	if rendered1 != rendered2 {
		t.Error("expected consistent rendering")
	}
}

func TestModalContainerLongContent(t *testing.T) {
	title := "Long Title " + strings.Repeat("Extended ", 10)
	message := strings.Repeat("This is a message. ", 10)

	mc := NewModalContainer().
		SetTitle(title).
		SetMessage(message)

	rendered := mc.Render()

	if !strings.Contains(rendered, "Extended") || !strings.Contains(rendered, "This is a message") {
		t.Error("expected rendered output to contain long content")
	}
}

func TestModalContainerDestructiveWithAllSections(t *testing.T) {
	title := "Delete Permanently"
	message := "This cannot be undone"
	buttons := []string{"Delete", "Cancel"}
	instructions := "Press Tab to navigate"

	mc := NewModalContainer().
		SetTitle(title).
		SetMessage(message).
		SetButtons(buttons).
		SetInstructions(instructions).
		WithDestructiveStyle()

	rendered := mc.Render()

	if !strings.Contains(rendered, title) ||
		!strings.Contains(rendered, message) ||
		!strings.Contains(rendered, "Delete") ||
		!strings.Contains(rendered, instructions) {
		t.Error("expected rendered output to contain all sections")
	}
}

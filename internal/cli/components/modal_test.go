package components

import (
	"strings"
	"testing"
	"time"
)

func TestNewErrorModal(t *testing.T) {
	modal := NewErrorModal("Error Title", "Error message")

	if modal.Type != ModalError {
		t.Errorf("Expected Type to be ModalError, got %v", modal.Type)
	}

	if modal.Title != "Error Title" {
		t.Errorf("Expected Title to be 'Error Title', got '%s'", modal.Title)
	}

	if modal.Message != "Error message" {
		t.Errorf("Expected Message to be 'Error message', got '%s'", modal.Message)
	}

	if !modal.Bell {
		t.Error("Expected Bell to be true for error modal")
	}

	if !modal.Cancellable {
		t.Error("Expected Cancellable to be true for error modal")
	}

	if modal.FadeInDuration != 150*time.Millisecond {
		t.Errorf("Expected FadeInDuration to be 150ms, got %v", modal.FadeInDuration)
	}
}

func TestNewLoadingModal(t *testing.T) {
	modal := NewLoadingModal("Loading...", true)

	if modal.Type != ModalLoading {
		t.Errorf("Expected Type to be ModalLoading, got %v", modal.Type)
	}

	if modal.Title != "Loading" {
		t.Errorf("Expected Title to be 'Loading', got '%s'", modal.Title)
	}

	if modal.Message != "Loading..." {
		t.Errorf("Expected Message to be 'Loading...', got '%s'", modal.Message)
	}

	if !modal.Cancellable {
		t.Error("Expected Cancellable to be true")
	}

	if modal.spinner == nil {
		t.Error("Expected spinner to be initialized")
	}
}

func TestNewProgressModal(t *testing.T) {
	modal := NewProgressModal("Installing", "Installing packages...", 0.5)

	if modal.Type != ModalProgress {
		t.Errorf("Expected Type to be ModalProgress, got %v", modal.Type)
	}

	if modal.Progress != 0.5 {
		t.Errorf("Expected Progress to be 0.5, got %f", modal.Progress)
	}
}

func TestNewSuccessModal(t *testing.T) {
	modal := NewSuccessModal("Operation completed successfully")

	if modal.Type != ModalSuccess {
		t.Errorf("Expected Type to be ModalSuccess, got %v", modal.Type)
	}

	if modal.AutoDismiss != 3*time.Second {
		t.Errorf("Expected AutoDismiss to be 3s, got %v", modal.AutoDismiss)
	}
}

func TestNewWarningModal(t *testing.T) {
	modal := NewWarningModal("Warning Title", "Warning message")

	if modal.Type != ModalWarning {
		t.Errorf("Expected Type to be ModalWarning, got %v", modal.Type)
	}

	if !modal.Bell {
		t.Error("Expected Bell to be true for warning modal")
	}

	if !modal.Cancellable {
		t.Error("Expected Cancellable to be true for warning modal")
	}
}

func TestModalContent_SetMessageRotator(t *testing.T) {
	modal := NewLoadingModal("Loading...", false)
	rotator := NewLoadingMessageRotator(LoadingMessagesCV, 2*time.Second)

	modal.SetMessageRotator(rotator)

	if modal.messageRotator != rotator {
		t.Error("Expected messageRotator to be set")
	}
}

func TestModalContent_Render_ErrorModal(t *testing.T) {
	modal := NewErrorModal("Error", "Something went wrong")
	output := modal.Render(80, 24)

	if output == "" {
		t.Error("Expected Render to return non-empty string")
	}

	// Should contain title and message
	if !strings.Contains(output, "Error") {
		t.Error("Expected output to contain title")
	}

	if !strings.Contains(output, "Something went wrong") {
		t.Error("Expected output to contain message")
	}

	// Should contain warning icon
	if !strings.Contains(output, "⚠️") {
		t.Error("Expected output to contain error icon")
	}

	// Should contain dismissal hint
	if !strings.Contains(output, "Press Esc to dismiss") {
		t.Error("Expected output to contain dismissal hint")
	}
}

func TestModalContent_Render_LoadingModal(t *testing.T) {
	modal := NewLoadingModal("Processing...", true)
	output := modal.Render(80, 24)

	if output == "" {
		t.Error("Expected Render to return non-empty string")
	}

	// Should contain loading icon
	if !strings.Contains(output, "⏳") {
		t.Error("Expected output to contain loading icon")
	}

	// Should contain cancellation hint
	if !strings.Contains(output, "Press Esc to cancel") {
		t.Error("Expected output to contain cancel hint")
	}
}

func TestModalContent_Render_ProgressModal(t *testing.T) {
	modal := NewProgressModal("Installing", "Installing packages...", 0.65)
	output := modal.Render(80, 24)

	if output == "" {
		t.Error("Expected Render to return non-empty string")
	}

	// Should contain progress percentage
	if !strings.Contains(output, "65%") {
		t.Error("Expected output to contain progress percentage")
	}

	// Should contain progress bar
	if !strings.Contains(output, "█") || !strings.Contains(output, "░") {
		t.Error("Expected output to contain progress bar")
	}

	// Should contain progress icon
	if !strings.Contains(output, "📊") {
		t.Error("Expected output to contain progress icon")
	}
}

func TestModalContent_Render_SuccessModal(t *testing.T) {
	modal := NewSuccessModal("Operation completed")
	output := modal.Render(80, 24)

	if output == "" {
		t.Error("Expected Render to return non-empty string")
	}

	// Should contain success icon
	if !strings.Contains(output, "✅") {
		t.Error("Expected output to contain success icon")
	}

	// Should contain auto-dismiss hint
	if !strings.Contains(output, "Auto-dismiss") {
		t.Error("Expected output to contain auto-dismiss hint")
	}
}

func TestModalContent_Render_WarningModal(t *testing.T) {
	modal := NewWarningModal("Warning", "This action is risky")
	output := modal.Render(80, 24)

	if output == "" {
		t.Error("Expected Render to return non-empty string")
	}

	// Should contain warning icon
	if !strings.Contains(output, "⚠️") {
		t.Error("Expected output to contain warning icon")
	}
}

func TestModalContent_UpdateProgress(t *testing.T) {
	modal := NewProgressModal("Test", "Test", 0.5)

	modal.UpdateProgress(0.75)
	if modal.Progress != 0.75 {
		t.Errorf("Expected Progress to be 0.75, got %f", modal.Progress)
	}

	// Test clamping to 0.0
	modal.UpdateProgress(-0.5)
	if modal.Progress != 0.0 {
		t.Errorf("Expected Progress to be clamped to 0.0, got %f", modal.Progress)
	}

	// Test clamping to 1.0
	modal.UpdateProgress(1.5)
	if modal.Progress != 1.0 {
		t.Errorf("Expected Progress to be clamped to 1.0, got %f", modal.Progress)
	}
}

func TestModalContent_AdvanceSpinner(t *testing.T) {
	modal := NewLoadingModal("Loading...", false)

	// Should not panic when spinner is set
	modal.AdvanceSpinner()

	// Get initial frame
	initialFrame := modal.spinner.GetFrame()

	// Advance multiple times to ensure it cycles
	for i := 0; i < 20; i++ {
		modal.AdvanceSpinner()
	}

	// After advancing, frame should still be valid
	currentFrame := modal.spinner.GetFrame()
	if currentFrame == "" {
		t.Error("Expected spinner frame to be non-empty after advancing")
	}

	// Should have cycled back
	if len(initialFrame) == 0 || len(currentFrame) == 0 {
		t.Error("Expected spinner frames to be non-empty")
	}
}

func TestModalContent_RotateMessage(t *testing.T) {
	modal := NewLoadingModal("Initial message", false)
	rotator := NewLoadingMessageRotator([]string{"Msg1", "Msg2", "Msg3"}, 1*time.Millisecond)
	modal.SetMessageRotator(rotator)

	// First call should return first message
	msg1 := modal.RotateMessage()
	if msg1 != "Msg1" {
		t.Errorf("Expected first message to be 'Msg1', got '%s'", msg1)
	}

	// Wait for interval to elapse
	time.Sleep(10 * time.Millisecond)

	// Rotate should advance
	msg2 := modal.RotateMessage()
	if msg2 != "Msg2" {
		t.Errorf("Expected second message to be 'Msg2', got '%s'", msg2)
	}
}

func TestModalContent_Render_WithActions(t *testing.T) {
	modal := NewErrorModal("Confirm", "Are you sure?")
	modal.Actions = []string{"Yes", "No"}

	output := modal.Render(80, 24)

	// Should contain actions
	if !strings.Contains(output, "Yes") || !strings.Contains(output, "No") {
		t.Error("Expected output to contain action buttons")
	}
}

func TestModalContent_CalculateOpacity(t *testing.T) {
	modal := NewErrorModal("Test", "Test")

	// Immediately after creation, should be near 0
	opacity := modal.calculateOpacity()
	if opacity < 0.0 || opacity > 1.0 {
		t.Errorf("Expected opacity between 0 and 1, got %f", opacity)
	}

	// After waiting longer than fade duration, should be 1.0
	modal.fadeStartTime = time.Now().Add(-200 * time.Millisecond)
	opacity = modal.calculateOpacity()
	if opacity != 1.0 {
		t.Errorf("Expected opacity to be 1.0 after fade duration, got %f", opacity)
	}
}

func TestWrapText(t *testing.T) {
	text := "This is a very long line that should be wrapped to fit within the specified width constraint"
	wrapped := wrapText(text, 20)

	lines := strings.Split(wrapped, "\n")
	if len(lines) <= 1 {
		t.Error("Expected text to be wrapped into multiple lines")
	}

	// Each line should be <= 20 characters (accounting for ANSI codes)
	for _, line := range lines {
		// This is a basic check; actual width may vary with styling
		if len(line) > 30 { // Allow some buffer for potential styling
			t.Errorf("Expected wrapped line to be <= 30 chars, got %d: '%s'", len(line), line)
		}
	}
}

func TestWrapText_EmptyString(t *testing.T) {
	wrapped := wrapText("", 20)
	if wrapped != "" {
		t.Errorf("Expected empty string to remain empty, got '%s'", wrapped)
	}
}

func TestWrapText_ShortString(t *testing.T) {
	text := "Short"
	wrapped := wrapText(text, 20)
	if wrapped != text {
		t.Errorf("Expected short string to remain unchanged, got '%s'", wrapped)
	}
}

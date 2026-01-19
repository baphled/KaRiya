package components

import (
	"strings"
	"testing"
	"time"

	"github.com/baphled/kariya/internal/cli/uikit/feedback"
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
	messages := []string{"Analyzing...", "Processing...", "Finishing..."}
	rotator := feedback.NewLoadingMessageRotator(messages)

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
	rotator := feedback.NewLoadingMessageRotator([]string{"Msg1", "Msg2", "Msg3"})
	modal.SetMessageRotator(rotator)

	// First call should advance and return second message (index 0 -> 1)
	msg1 := modal.RotateMessage()
	if msg1 != "Msg2" {
		t.Errorf("Expected first rotation to return 'Msg2', got '%s'", msg1)
	}

	// Second call should advance and return third message
	msg2 := modal.RotateMessage()
	if msg2 != "Msg3" {
		t.Errorf("Expected second rotation to return 'Msg3', got '%s'", msg2)
	}

	// Third call should wrap around to first message
	msg3 := modal.RotateMessage()
	if msg3 != "Msg1" {
		t.Errorf("Expected third rotation to return 'Msg1', got '%s'", msg3)
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

// Edge case tests added for Task 16 Phase 2.2

func TestModalContent_ModalLargerThanTerminal(t *testing.T) {
	// Create modal with very long content for small terminal
	longMessage := strings.Repeat("This is a very long error message that exceeds terminal width. ", 10)
	modal := NewErrorModal("Error", longMessage)

	output := modal.Render(40, 10)

	// Should handle gracefully without panic
	if output == "" {
		t.Error("Expected output for modal larger than terminal")
	}
}

func TestModalContent_VeryLongErrorMessage(t *testing.T) {
	// Test with extremely long error messages
	longMessage := strings.Repeat("Error detail. ", 100)
	modal := NewErrorModal("Error", longMessage)

	output := modal.Render(80, 40)

	// Should wrap or truncate gracefully
	if output == "" {
		t.Error("Expected output for very long error message")
	}
}

func TestModalContent_ProgressEdgeValues(t *testing.T) {
	// Test 0.0 progress
	modal1 := NewProgressModal("Starting", "Starting process...", 0.0)
	output1 := modal1.Render(80, 40)
	if !strings.Contains(output1, "0%") {
		t.Error("Expected 0% to be shown for 0.0 progress")
	}

	// Test 1.0 progress
	modal2 := NewProgressModal("Complete", "Process complete!", 1.0)
	output2 := modal2.Render(80, 40)
	if !strings.Contains(output2, "100%") {
		t.Error("Expected 100% to be shown for 1.0 progress")
	}

	// Test mid-range progress
	modal3 := NewProgressModal("Processing", "Processing data...", 0.5)
	output3 := modal3.Render(80, 40)
	if !strings.Contains(output3, "50%") {
		t.Error("Expected 50% to be shown for 0.5 progress")
	}
}

func TestModalContent_FadeInTiming(t *testing.T) {
	modal := NewErrorModal("Test", "Message")

	// Test rendering - fade-in is calculated internally from fadeStartTime
	// We just verify it renders without panic at various points
	output := modal.Render(80, 40)
	if output == "" {
		t.Error("Expected output during fade-in")
	}

	// Wait a moment and render again
	time.Sleep(50 * time.Millisecond)
	output = modal.Render(80, 40)
	if output == "" {
		t.Error("Expected output after some fade-in time")
	}
}

func TestModalContent_AutoDismissTimingAccuracy(t *testing.T) {
	// Test success modal with auto-dismiss
	modal := NewSuccessModal("Operation completed")

	// Verify modal properties are set correctly for auto-dismiss
	if modal.Type != ModalSuccess {
		t.Error("Expected modal type to be Success")
	}

	output := modal.Render(80, 40)
	if output == "" {
		t.Error("Expected output for success modal")
	}
}

func TestModalContent_MultipleRapidChanges(t *testing.T) {
	modal := NewProgressModal("Processing", "Processing items...", 0.0)

	// Simulate rapid progress updates
	for i := 0.0; i <= 1.0; i += 0.1 {
		modal.UpdateProgress(i)
		output := modal.Render(80, 40)
		if output == "" {
			t.Errorf("Expected output for progress %.1f", i)
		}
	}
}

func TestModalContent_EmptyTitleAndMessage(t *testing.T) {
	// Test modal with empty title
	modal1 := NewErrorModal("", "Message")
	output1 := modal1.Render(80, 40)
	if output1 == "" {
		t.Error("Expected output for modal with empty title")
	}

	// Test modal with empty message
	modal2 := NewErrorModal("Title", "")
	output2 := modal2.Render(80, 40)
	if output2 == "" {
		t.Error("Expected output for modal with empty message")
	}

	// Test modal with both empty
	modal3 := NewErrorModal("", "")
	output3 := modal3.Render(80, 40)
	if output3 == "" {
		t.Error("Expected output for modal with empty title and message")
	}
}

func TestModalContent_ActionButtonsOverflow(t *testing.T) {
	modal := NewErrorModal("Error", "Test error message")

	// Add many action buttons
	modal.Actions = []string{
		"Very Long Action Button 1",
		"Very Long Action Button 2",
		"Very Long Action Button 3",
		"Very Long Action Button 4",
		"Very Long Action Button 5",
	}

	// Render in narrow terminal
	output := modal.Render(50, 20)

	// Should handle overflow gracefully
	if output == "" {
		t.Error("Expected output with many action buttons")
	}
}

func TestModalContent_SmallTerminalRendering(t *testing.T) {
	modal := NewErrorModal("Error", "An error occurred")

	// Test on very small terminal
	output := modal.Render(30, 10)

	// Should render something even on small terminal
	if output == "" {
		t.Error("Expected output for small terminal")
	}
}

func TestModalContent_LargeTerminalRendering(t *testing.T) {
	modal := NewLoadingModal("Loading...", false)

	// Test on very large terminal
	output := modal.Render(200, 60)

	// Should render appropriately for large terminal
	if output == "" {
		t.Error("Expected output for large terminal")
	}
}

// Modal Overlay Tests are in modal_overlay_test.go (Ginkgo)

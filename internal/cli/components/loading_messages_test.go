package components

import (
	"testing"
	"time"
)

func TestNewLoadingMessageRotator(t *testing.T) {
	messages := []string{"Msg1", "Msg2", "Msg3"}
	rotator := NewLoadingMessageRotator(messages, 2*time.Second)

	if rotator == nil {
		t.Fatal("Expected NewLoadingMessageRotator to return non-nil rotator")
	}

	if len(rotator.messages) != 3 {
		t.Errorf("Expected 3 messages, got %d", len(rotator.messages))
	}

	if rotator.currentIndex != 0 {
		t.Errorf("Expected currentIndex to be 0, got %d", rotator.currentIndex)
	}

	if rotator.rotateInterval != 2*time.Second {
		t.Errorf("Expected rotateInterval to be 2s, got %v", rotator.rotateInterval)
	}
}

func TestNewLoadingMessageRotator_EmptyMessages(t *testing.T) {
	rotator := NewLoadingMessageRotator([]string{}, 0)

	// Should fall back to LoadingMessagesGeneric
	if len(rotator.messages) == 0 {
		t.Error("Expected rotator to have fallback messages")
	}
}

func TestNewLoadingMessageRotator_ZeroInterval(t *testing.T) {
	rotator := NewLoadingMessageRotator([]string{"Msg1"}, 0)

	// Should default to 2 seconds
	if rotator.rotateInterval != 2*time.Second {
		t.Errorf("Expected default rotateInterval to be 2s, got %v", rotator.rotateInterval)
	}
}

func TestLoadingMessageRotator_GetCurrent(t *testing.T) {
	messages := []string{"Msg1", "Msg2", "Msg3"}
	rotator := NewLoadingMessageRotator(messages, 2*time.Second)

	msg := rotator.GetCurrent()
	if msg != "Msg1" {
		t.Errorf("Expected first message to be 'Msg1', got '%s'", msg)
	}

	// Should not advance
	msg2 := rotator.GetCurrent()
	if msg2 != "Msg1" {
		t.Errorf("Expected GetCurrent to not advance, got '%s'", msg2)
	}
}

func TestLoadingMessageRotator_Rotate(t *testing.T) {
	messages := []string{"Msg1", "Msg2", "Msg3"}
	// Use very small interval to allow quick rotation in tests
	rotator := NewLoadingMessageRotator(messages, 1*time.Millisecond)

	// First call should return first message
	msg1 := rotator.Rotate()
	if msg1 != "Msg1" {
		t.Errorf("Expected first message to be 'Msg1', got '%s'", msg1)
	}

	// Wait for interval to elapse
	time.Sleep(10 * time.Millisecond)

	// Second call should advance
	msg2 := rotator.Rotate()
	if msg2 != "Msg2" {
		t.Errorf("Expected second message to be 'Msg2', got '%s'", msg2)
	}

	time.Sleep(10 * time.Millisecond)

	// Third call should advance
	msg3 := rotator.Rotate()
	if msg3 != "Msg3" {
		t.Errorf("Expected third message to be 'Msg3', got '%s'", msg3)
	}

	time.Sleep(10 * time.Millisecond)

	// Fourth call should wrap around
	msg4 := rotator.Rotate()
	if msg4 != "Msg1" {
		t.Errorf("Expected fourth message to wrap to 'Msg1', got '%s'", msg4)
	}
}

func TestLoadingMessageRotator_Rotate_WithInterval(t *testing.T) {
	messages := []string{"Msg1", "Msg2"}
	rotator := NewLoadingMessageRotator(messages, 100*time.Millisecond)

	// First call
	msg1 := rotator.Rotate()
	if msg1 != "Msg1" {
		t.Errorf("Expected first message to be 'Msg1', got '%s'", msg1)
	}

	// Immediate second call should not advance (interval not elapsed)
	msg2 := rotator.Rotate()
	if msg2 != "Msg1" {
		t.Errorf("Expected message to not advance before interval, got '%s'", msg2)
	}

	// Wait for interval
	time.Sleep(150 * time.Millisecond)

	// Now should advance
	msg3 := rotator.Rotate()
	if msg3 != "Msg2" {
		t.Errorf("Expected message to advance after interval, got '%s'", msg3)
	}
}

func TestLoadingMessageRotator_Reset(t *testing.T) {
	messages := []string{"Msg1", "Msg2", "Msg3"}
	rotator := NewLoadingMessageRotator(messages, 0)

	// Advance a few times
	time.Sleep(10 * time.Millisecond)
	rotator.Rotate()
	time.Sleep(10 * time.Millisecond)
	rotator.Rotate()

	// Reset should go back to first message
	rotator.Reset()

	msg := rotator.GetCurrent()
	if msg != "Msg1" {
		t.Errorf("Expected Reset to return to first message, got '%s'", msg)
	}

	if rotator.currentIndex != 0 {
		t.Errorf("Expected currentIndex to be 0 after Reset, got %d", rotator.currentIndex)
	}
}

func TestLoadingMessageRotator_SetMessages(t *testing.T) {
	rotator := NewLoadingMessageRotator([]string{"Msg1"}, 2*time.Second)

	newMessages := []string{"NewMsg1", "NewMsg2"}
	rotator.SetMessages(newMessages)

	if len(rotator.messages) != 2 {
		t.Errorf("Expected 2 messages after SetMessages, got %d", len(rotator.messages))
	}

	msg := rotator.GetCurrent()
	if msg != "NewMsg1" {
		t.Errorf("Expected first new message to be 'NewMsg1', got '%s'", msg)
	}

	// Should reset index
	if rotator.currentIndex != 0 {
		t.Errorf("Expected currentIndex to be reset to 0, got %d", rotator.currentIndex)
	}
}

func TestLoadingMessageRotator_SetMessages_Empty(t *testing.T) {
	rotator := NewLoadingMessageRotator([]string{"Msg1"}, 2*time.Second)

	rotator.SetMessages([]string{})

	// Should not change messages
	if len(rotator.messages) != 1 {
		t.Errorf("Expected messages to remain unchanged, got %d", len(rotator.messages))
	}
}

func TestNewSimpleSpinner(t *testing.T) {
	spinner := NewSimpleSpinner()

	if spinner == nil {
		t.Fatal("Expected NewSimpleSpinner to return non-nil spinner")
	}

	if len(spinner.frames) == 0 {
		t.Error("Expected spinner to have frames")
	}

	if spinner.currentFrame != 0 {
		t.Errorf("Expected currentFrame to be 0, got %d", spinner.currentFrame)
	}

	if spinner.frameInterval != 80*time.Millisecond {
		t.Errorf("Expected frameInterval to be 80ms, got %v", spinner.frameInterval)
	}
}

func TestSimpleSpinner_GetFrame(t *testing.T) {
	spinner := NewSimpleSpinner()

	frame := spinner.GetFrame()
	if frame == "" {
		t.Error("Expected GetFrame to return non-empty frame")
	}

	// Should be one of the default frames
	found := false
	for _, f := range defaultSpinnerFrames {
		if frame == f {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected frame to be from default frames, got '%s'", frame)
	}
}

func TestSimpleSpinner_Advance(t *testing.T) {
	spinner := NewSimpleSpinner()

	initialFrame := spinner.GetFrame()

	// Advance immediately - should not change (interval not elapsed)
	spinner.Advance()
	currentFrame := spinner.GetFrame()
	if currentFrame != initialFrame {
		t.Error("Expected frame to not change before interval elapsed")
	}

	// Wait for interval
	spinner.lastUpdate = time.Now().Add(-100 * time.Millisecond)
	spinner.Advance()
	newFrame := spinner.GetFrame()

	// Should have advanced to next frame
	if newFrame == initialFrame {
		t.Error("Expected frame to change after interval elapsed")
	}
}

func TestSimpleSpinner_Advance_Wrapping(t *testing.T) {
	spinner := NewSimpleSpinner()

	// Advance through all frames
	for i := 0; i < len(defaultSpinnerFrames)*2; i++ {
		spinner.lastUpdate = time.Now().Add(-100 * time.Millisecond)
		spinner.Advance()
	}

	// Should still return a valid frame (wrapped around)
	frame := spinner.GetFrame()
	if frame == "" {
		t.Error("Expected GetFrame to return valid frame after wrapping")
	}

	// currentFrame should be within bounds
	if spinner.currentFrame >= len(spinner.frames) {
		t.Errorf("Expected currentFrame to wrap around, got %d", spinner.currentFrame)
	}
}

func TestSimpleSpinner_SetFrames(t *testing.T) {
	spinner := NewSimpleSpinner()

	customFrames := []string{"A", "B", "C"}
	spinner.SetFrames(customFrames)

	if len(spinner.frames) != 3 {
		t.Errorf("Expected 3 frames, got %d", len(spinner.frames))
	}

	frame := spinner.GetFrame()
	if frame != "A" {
		t.Errorf("Expected first frame to be 'A', got '%s'", frame)
	}

	// Should reset index
	if spinner.currentFrame != 0 {
		t.Errorf("Expected currentFrame to be reset to 0, got %d", spinner.currentFrame)
	}
}

func TestSimpleSpinner_SetFrames_Empty(t *testing.T) {
	spinner := NewSimpleSpinner()
	originalFrames := spinner.frames

	spinner.SetFrames([]string{})

	// Should not change frames
	if len(spinner.frames) != len(originalFrames) {
		t.Error("Expected frames to remain unchanged when setting empty frames")
	}
}

func TestSimpleSpinner_SetFrameInterval(t *testing.T) {
	spinner := NewSimpleSpinner()

	spinner.SetFrameInterval(50 * time.Millisecond)

	if spinner.frameInterval != 50*time.Millisecond {
		t.Errorf("Expected frameInterval to be 50ms, got %v", spinner.frameInterval)
	}
}

func TestPredefinedMessageSets(t *testing.T) {
	// Test that all predefined message sets are non-empty
	if len(LoadingMessagesCV) == 0 {
		t.Error("Expected LoadingMessagesCV to be non-empty")
	}

	if len(LoadingMessagesExport) == 0 {
		t.Error("Expected LoadingMessagesExport to be non-empty")
	}

	if len(LoadingMessagesGeneric) == 0 {
		t.Error("Expected LoadingMessagesGeneric to be non-empty")
	}

	if len(LoadingMessagesFetch) == 0 {
		t.Error("Expected LoadingMessagesFetch to be non-empty")
	}
}

func TestDefaultSpinnerFrames(t *testing.T) {
	if len(defaultSpinnerFrames) != 10 {
		t.Errorf("Expected 10 default spinner frames, got %d", len(defaultSpinnerFrames))
	}

	// Each frame should be non-empty
	for i, frame := range defaultSpinnerFrames {
		if frame == "" {
			t.Errorf("Expected frame %d to be non-empty", i)
		}
	}
}

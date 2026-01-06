package components

import (
	"time"
)

// LoadingMessageRotator rotates through a set of loading messages
type LoadingMessageRotator struct {
	messages       []string
	currentIndex   int
	rotateInterval time.Duration
	lastRotation   time.Time
}

// Predefined message sets for common operations
var (
	// LoadingMessagesCV for CV generation operations
	LoadingMessagesCV = []string{
		"Analyzing career events...",
		"Extracting relevant facts...",
		"Organizing content by audience...",
		"Crafting professional bullet points...",
		"Finalizing your CV...",
	}

	// LoadingMessagesExport for export operations
	LoadingMessagesExport = []string{
		"Preparing export...",
		"Formatting content...",
		"Writing file...",
		"Export complete!",
	}

	// LoadingMessagesGeneric for general operations
	LoadingMessagesGeneric = []string{
		"Processing...",
		"Almost there...",
		"Just a moment...",
		"Finishing up...",
		"Done!",
	}

	// LoadingMessagesFetch for data fetching operations
	LoadingMessagesFetch = []string{
		"Fetching data...",
		"Loading records...",
		"Preparing results...",
		"Ready!",
	}
)

// NewLoadingMessageRotator creates a new loading message rotator
func NewLoadingMessageRotator(messages []string, interval time.Duration) *LoadingMessageRotator {
	if len(messages) == 0 {
		messages = LoadingMessagesGeneric
	}
	if interval == 0 {
		interval = 2 * time.Second
	}

	return &LoadingMessageRotator{
		messages:       messages,
		currentIndex:   0,
		rotateInterval: interval,
		lastRotation:   time.Now(),
	}
}

// GetCurrent returns the current message without rotating
func (r *LoadingMessageRotator) GetCurrent() string {
	if r.currentIndex >= len(r.messages) {
		r.currentIndex = 0
	}
	return r.messages[r.currentIndex]
}

// Rotate advances to the next message if the interval has elapsed
func (r *LoadingMessageRotator) Rotate() string {
	now := time.Now()
	if now.Sub(r.lastRotation) >= r.rotateInterval {
		r.currentIndex++
		if r.currentIndex >= len(r.messages) {
			r.currentIndex = 0
		}
		r.lastRotation = now
	}
	return r.GetCurrent()
}

// Reset resets to the first message
func (r *LoadingMessageRotator) Reset() {
	r.currentIndex = 0
	r.lastRotation = time.Now()
}

// SetMessages updates the message set
func (r *LoadingMessageRotator) SetMessages(messages []string) {
	if len(messages) > 0 {
		r.messages = messages
		r.Reset()
	}
}

// SimpleSpinner provides a simple spinner animation for modals
type SimpleSpinner struct {
	frames        []string
	currentFrame  int
	lastUpdate    time.Time
	frameInterval time.Duration
}

// Default spinner frames (Braille dots)
var defaultSpinnerFrames = []string{
	"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏",
}

// NewSimpleSpinner creates a new simple spinner with default frames
func NewSimpleSpinner() *SimpleSpinner {
	return &SimpleSpinner{
		frames:        defaultSpinnerFrames,
		currentFrame:  0,
		lastUpdate:    time.Now(),
		frameInterval: 80 * time.Millisecond,
	}
}

// GetFrame returns the current spinner frame
func (s *SimpleSpinner) GetFrame() string {
	if s.currentFrame >= len(s.frames) {
		s.currentFrame = 0
	}
	return s.frames[s.currentFrame]
}

// Advance advances to the next frame
func (s *SimpleSpinner) Advance() {
	now := time.Now()
	if now.Sub(s.lastUpdate) >= s.frameInterval {
		s.currentFrame++
		if s.currentFrame >= len(s.frames) {
			s.currentFrame = 0
		}
		s.lastUpdate = now
	}
}

// SetFrames sets custom spinner frames
func (s *SimpleSpinner) SetFrames(frames []string) {
	if len(frames) > 0 {
		s.frames = frames
		s.currentFrame = 0
	}
}

// SetFrameInterval sets the interval between frame changes
func (s *SimpleSpinner) SetFrameInterval(interval time.Duration) {
	s.frameInterval = interval
}

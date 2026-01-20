package components

import (
	"sync"
	"time"
)

// LoadingMessageRotator rotates through a set of loading messages
type LoadingMessageRotator struct {
	mu             sync.RWMutex
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
	r.mu.RLock()
	defer r.mu.RUnlock()

	currentIndex := r.currentIndex
	if currentIndex >= len(r.messages) {
		currentIndex = 0
	}
	return r.messages[currentIndex]
}

// Rotate advances to the next message if the interval has elapsed, then returns current
func (r *LoadingMessageRotator) Rotate() string {
	r.mu.Lock()
	now := time.Now()
	if now.Sub(r.lastRotation) >= r.rotateInterval {
		r.currentIndex++
		if r.currentIndex >= len(r.messages) {
			r.currentIndex = 0
		}
		r.lastRotation = now
	}
	r.mu.Unlock()

	return r.GetCurrent()
}

// Reset resets to the first message
func (r *LoadingMessageRotator) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.currentIndex = 0
	r.lastRotation = time.Now()
}

// SetMessages updates the message set
func (r *LoadingMessageRotator) SetMessages(messages []string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(messages) > 0 {
		r.messages = messages
		r.currentIndex = 0
		r.lastRotation = time.Now()
	}
}

// SimpleSpinner provides a simple spinner animation for modals
type SimpleSpinner struct {
	mu            sync.RWMutex
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
	s.mu.RLock()
	defer s.mu.RUnlock()

	currentFrame := s.currentFrame
	if currentFrame >= len(s.frames) {
		currentFrame = 0
	}
	return s.frames[currentFrame]
}

// Advance advances to the next frame
func (s *SimpleSpinner) Advance() {
	s.mu.Lock()
	defer s.mu.Unlock()

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
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(frames) > 0 {
		s.frames = frames
		s.currentFrame = 0
	}
}

// SetFrameInterval sets the interval between frame changes
func (s *SimpleSpinner) SetFrameInterval(interval time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.frameInterval = interval
}

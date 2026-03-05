package harness

import (
	"sync"
)

// StubClipboardWriter provides a deterministic clipboard implementation for BDD tests.
// It stores content in memory instead of using the system clipboard, avoiding
// failures in headless/CI environments where no display is available.
type StubClipboardWriter struct {
	mu      sync.Mutex
	content string
}

// WriteAll writes text to the stub clipboard, storing it in memory.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (s *StubClipboardWriter) WriteAll(text string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.content = text
	return nil
}

// IsUnsupported returns false since the stub clipboard is always available.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (s *StubClipboardWriter) IsUnsupported() bool {
	return false
}

// GetContent returns the current clipboard content for test assertions.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (s *StubClipboardWriter) GetContent() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.content
}

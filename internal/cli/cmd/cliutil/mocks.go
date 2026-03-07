package cliutil

import (
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// MockFileOpener is a mock implementation of FileOpener.
type MockFileOpener struct {
	OpenFn func(path string) (io.ReadCloser, error)
	StatFn func(path string) (os.FileInfo, error)
}

// Open calls OpenFn if set, otherwise returns os.ErrNotExist.
//
// Expected:
//   - path must be a valid file path.
//
// Returns:
//   - An io.ReadCloser and nil error if OpenFn is set, otherwise nil and os.ErrNotExist.
//
// Side effects:
//   - None (mock implementation).
func (m MockFileOpener) Open(path string) (io.ReadCloser, error) {
	if m.OpenFn != nil {
		return m.OpenFn(path)
	}
	return nil, os.ErrNotExist
}

// Stat calls StatFn if set, otherwise returns os.ErrNotExist.
//
// Expected:
//   - path must be a valid file path.
//
// Returns:
//   - os.FileInfo and nil error if StatFn is set, otherwise nil and os.ErrNotExist.
//
// Side effects:
//   - None (mock implementation).
func (m MockFileOpener) Stat(path string) (os.FileInfo, error) {
	if m.StatFn != nil {
		return m.StatFn(path)
	}
	return nil, os.ErrNotExist
}

// MockProgressRunner is a mock implementation of ProgressRunner.
type MockProgressRunner struct {
	RunWithSpinnerFn  func(message string, fn func() error, opts ...tea.ProgramOption) error
	RunWithProgressFn func(message string, total int, fn func(update func(current int)) error, opts ...tea.ProgramOption) error
}

// RunWithSpinner calls RunWithSpinnerFn if set, otherwise executes fn directly.
//
// Expected:
//   - message must be a non-empty string.
//   - fn must be a non-nil function.
//
// Returns:
//   - An error if the function fails, nil otherwise.
//
// Side effects:
//   - None (mock implementation).
func (m MockProgressRunner) RunWithSpinner(message string, fn func() error, opts ...tea.ProgramOption) error {
	if m.RunWithSpinnerFn != nil {
		return m.RunWithSpinnerFn(message, fn, opts...)
	}
	return fn()
}

// RunWithProgress calls RunWithProgressFn if set, otherwise executes fn directly with no-op update.
//
// Expected:
//   - message must be a non-empty string.
//   - fn must be a non-nil function that calls update to report progress.
//
// Returns:
//   - An error if the function fails, nil otherwise.
//
// Side effects:
//   - None (mock implementation).
func (m MockProgressRunner) RunWithProgress(
	message string,
	total int,
	fn func(update func(current int)) error,
	opts ...tea.ProgramOption,
) error {
	if m.RunWithProgressFn != nil {
		return m.RunWithProgressFn(message, total, fn, opts...)
	}
	return fn(func(_ int) {})
}

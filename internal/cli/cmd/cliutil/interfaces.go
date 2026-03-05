// Package cliutil provides common interfaces and utilities for CLI commands.
package cliutil

import (
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// FileOpener decouples filesystem access for better testability.
type FileOpener interface {
	Open(path string) (io.ReadCloser, error)
	Stat(path string) (os.FileInfo, error)
}

// OSFileOpener is the default implementation of FileOpener using the os package.
type OSFileOpener struct{}

// Open opens a file for reading using os.Open.
func (o OSFileOpener) Open(path string) (io.ReadCloser, error) {
	return os.Open(path)
}

// Stat returns FileInfo for a given path using os.Stat.
func (o OSFileOpener) Stat(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

// ProgressRunner defines the interface for running tasks with progress feedback.
type ProgressRunner interface {
	RunWithSpinner(message string, fn func() error, opts ...tea.ProgramOption) error
	RunWithProgress(message string, total int, fn func(update func(current int)) error, opts ...tea.ProgramOption) error
}

// DefaultProgressRunner is the standard implementation of ProgressRunner.
type DefaultProgressRunner struct{}

// RunWithSpinner runs a task with a spinner using the default implementation.
func (d DefaultProgressRunner) RunWithSpinner(message string, fn func() error, opts ...tea.ProgramOption) error {
	return RunWithSpinner(message, fn, opts...)
}

// RunWithProgress runs a task with a progress bar using the default implementation.
func (d DefaultProgressRunner) RunWithProgress(
	message string,
	total int,
	fn func(update func(current int)) error,
	opts ...tea.ProgramOption,
) error {
	return RunWithProgress(message, total, fn, opts...)
}

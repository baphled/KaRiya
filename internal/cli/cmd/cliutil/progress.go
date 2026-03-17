package cliutil

import (
	"fmt"
	"time"

	"github.com/baphled/kariya/internal/ui/uikit/feedback"
	"github.com/baphled/kariya/internal/ui/uikit/primitives"
	"github.com/baphled/kariya/internal/ui/uikit/theme"
	tea "github.com/charmbracelet/bubbletea"
)

// ============================================================================
// Spinner Utility
// ============================================================================

// RunWithSpinner runs a function with a spinner animation.
// It clears the spinner line upon completion.
func RunWithSpinner(message string, fn func() error, opts ...tea.ProgramOption) error {
	resultChan := make(chan error, 1)

	go func() {
		resultChan <- fn()
		close(resultChan)
	}()

	p := tea.NewProgram(newSpinnerModel(message, resultChan), opts...)
	m, err := p.Run()
	if err != nil {
		return err
	}

	if sm, ok := m.(spinnerModel); ok {
		return sm.err
	}
	return nil
}

type spinnerModel struct {
	message    string
	spinner    *feedback.SimpleSpinner
	resultChan chan error
	err        error
	done       bool
	quitting   bool
}

type spinnerTickMsg time.Time
type spinnerDoneMsg struct{ err error }

func newSpinnerModel(message string, resultChan chan error) spinnerModel {
	return spinnerModel{
		message:    message,
		spinner:    feedback.NewSimpleSpinner(),
		resultChan: resultChan,
	}
}

func (m spinnerModel) Init() tea.Cmd {
	return tea.Batch(
		m.tick(),
		m.waitForResult(),
	)
}

func (m spinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}
	case spinnerDoneMsg:
		m.done = true
		m.err = msg.err
		return m, tea.Quit
	case spinnerTickMsg:
		if m.done {
			return m, nil
		}
		m.spinner.Advance()
		return m, m.tick()
	}
	return m, nil
}

func (m spinnerModel) tick() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return spinnerTickMsg(t)
	})
}

func (m spinnerModel) waitForResult() tea.Cmd {
	return func() tea.Msg {
		err := <-m.resultChan
		return spinnerDoneMsg{err: err}
	}
}

func (m spinnerModel) View() string {
	if m.done || m.quitting {
		return ""
	}
	return fmt.Sprintf("%s %s", m.spinner.GetFrame(), m.message)
}

// ============================================================================
// Progress Bar Utility
// ============================================================================

// RunWithProgress runs a function with progress updates.
// It clears the progress line upon completion.
func RunWithProgress(message string, total int, fn func(update func(current int)) error, opts ...tea.ProgramOption) error {
	progressChan := make(chan int)
	resultChan := make(chan error, 1)

	go func() {
		// Provide an update function that sends to the channel
		err := fn(func(current int) {
			progressChan <- current
		})
		resultChan <- err
		close(progressChan)
		close(resultChan)
	}()

	p := tea.NewProgram(newProgressModel(message, total, progressChan, resultChan), opts...)
	m, err := p.Run()
	if err != nil {
		return err
	}

	if pm, ok := m.(progressModel); ok {
		return pm.err
	}
	return nil
}

type progressModel struct {
	message      string
	total        int
	current      int
	progressChan chan int
	resultChan   chan error
	err          error
	done         bool
	quitting     bool
	th           theme.Theme
}

type progressMsg int
type progressDoneMsg struct{ err error }

func newProgressModel(message string, total int, progressChan chan int, resultChan chan error) progressModel {
	return progressModel{
		message:      message,
		total:        total,
		progressChan: progressChan,
		resultChan:   resultChan,
		th:           theme.Default(),
	}
}

func (m progressModel) Init() tea.Cmd {
	return tea.Batch(
		m.waitForProgress(),
		m.waitForResult(),
	)
}

func (m progressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}
	case progressMsg:
		m.current = int(msg)
		return m, m.waitForProgress() // Re-subscribe for next update
	case progressDoneMsg:
		m.done = true
		m.err = msg.err
		return m, tea.Quit
	}
	return m, nil
}

func (m progressModel) waitForProgress() tea.Cmd {
	return func() tea.Msg {
		p, ok := <-m.progressChan
		if !ok {
			return nil // Channel closed, no more updates
		}
		return progressMsg(p)
	}
}

func (m progressModel) waitForResult() tea.Cmd {
	return func() tea.Msg {
		err := <-m.resultChan
		return progressDoneMsg{err: err}
	}
}

func (m progressModel) View() string {
	if m.done || m.quitting {
		return ""
	}

	percent := 0.0
	if m.total > 0 {
		percent = float64(m.current) / float64(m.total)
	}

	// Create progress bar using primitives
	pb := primitives.NewProgressBar(percent, m.th).
		Width(30).
		ShowPercentage(true)

	return fmt.Sprintf("%s %s", pb.Render(), m.message)
}

// ============================================================================
// Feedback Utilities
// ============================================================================

// PrintSuccess prints a success message.
func PrintSuccess(message string) {
	th := theme.Default()
	fmt.Println(primitives.SuccessText(message, th).Render())
}

// PrintError prints an error message.
func PrintError(message string) {
	th := theme.Default()
	fmt.Println(primitives.ErrorText(message, th).Render())
}

// PrintInfo prints an info message.
func PrintInfo(message string) {
	th := theme.Default()
	fmt.Println(primitives.InfoText(message, th).Render())
}

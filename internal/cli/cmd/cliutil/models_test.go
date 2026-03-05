package cliutil

import (
	"errors"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestModels(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CLI Util Models Suite")
}

var _ = Describe("spinnerModel", func() {
	var (
		resultChan chan error
		model      spinnerModel
	)

	BeforeEach(func() {
		resultChan = make(chan error, 1)
		model = newSpinnerModel("Loading...", resultChan)
	})

	Describe("Construction", func() {
		It("initializes with correct message", func() {
			Expect(model.message).To(Equal("Loading..."))
		})

		It("initializes with spinner", func() {
			Expect(model.spinner).NotTo(BeNil())
		})

		It("initializes with result channel", func() {
			Expect(model.resultChan).To(Equal(resultChan))
		})

		It("initializes with done=false", func() {
			Expect(model.done).To(BeFalse())
		})

		It("initializes with err=nil", func() {
			Expect(model.err).NotTo(HaveOccurred())
		})

		It("initializes with quitting=false", func() {
			Expect(model.quitting).To(BeFalse())
		})
	})

	Describe("Init", func() {
		It("returns a batch command", func() {
			cmd := model.Init()
			Expect(cmd).NotTo(BeNil())
		})

		It("returns a command that can be executed", func() {
			cmd := model.Init()
			msg := cmd()
			Expect(msg).NotTo(BeNil())
		})
	})

	Describe("Update", func() {
		Context("with spinnerTickMsg", func() {
			It("advances spinner when not done", func() {
				initialFrame := model.spinner.GetFrame()
				updated, cmd := model.Update(spinnerTickMsg(time.Now()))
				result := updated.(spinnerModel)

				Expect(result.spinner.GetFrame()).NotTo(Equal(initialFrame))
				Expect(cmd).NotTo(BeNil())
			})

			It("returns nil command when done", func() {
				model.done = true
				_, cmd := model.Update(spinnerTickMsg(time.Now()))
				Expect(cmd).To(BeNil())
			})

			It("does not advance spinner when done", func() {
				model.done = true
				initialFrame := model.spinner.GetFrame()
				updated, _ := model.Update(spinnerTickMsg(time.Now()))
				result := updated.(spinnerModel)

				Expect(result.spinner.GetFrame()).To(Equal(initialFrame))
			})
		})

		Context("with spinnerDoneMsg", func() {
			It("marks done=true", func() {
				updated, _ := model.Update(spinnerDoneMsg{err: nil})
				result := updated.(spinnerModel)
				Expect(result.done).To(BeTrue())
			})

			It("stores error from message", func() {
				testErr := errors.New("test error")
				updated, _ := model.Update(spinnerDoneMsg{err: testErr})
				result := updated.(spinnerModel)
				Expect(result.err).To(Equal(testErr))
			})

			It("returns tea.Quit command", func() {
				_, cmd := model.Update(spinnerDoneMsg{err: nil})
				Expect(cmd).NotTo(BeNil())
				msg := cmd()
				_, isQuit := msg.(tea.QuitMsg)
				Expect(isQuit).To(BeTrue())
			})

			It("handles nil error", func() {
				updated, _ := model.Update(spinnerDoneMsg{err: nil})
				result := updated.(spinnerModel)
				Expect(result.err).NotTo(HaveOccurred())
			})
		})

		Context("with ctrl+c key", func() {
			It("marks quitting=true", func() {
				updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
				result := updated.(spinnerModel)
				Expect(result.quitting).To(BeTrue())
			})

			It("returns tea.Quit command", func() {
				_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
				Expect(cmd).NotTo(BeNil())
				msg := cmd()
				_, isQuit := msg.(tea.QuitMsg)
				Expect(isQuit).To(BeTrue())
			})
		})

		Context("with other key messages", func() {
			It("ignores other keys", func() {
				updated, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
				result := updated.(spinnerModel)
				Expect(result.quitting).To(BeFalse())
				Expect(cmd).To(BeNil())
			})
		})

		Context("with unknown message type", func() {
			It("returns model unchanged", func() {
				updated, cmd := model.Update("unknown")
				result := updated.(spinnerModel)
				Expect(result).To(Equal(model))
				Expect(cmd).To(BeNil())
			})
		})
	})

	Describe("View", func() {
		It("returns empty string when done", func() {
			model.done = true
			view := model.View()
			Expect(view).To(Equal(""))
		})

		It("returns empty string when quitting", func() {
			model.quitting = true
			view := model.View()
			Expect(view).To(Equal(""))
		})

		It("returns spinner and message when active", func() {
			view := model.View()
			Expect(view).To(ContainSubstring("Loading..."))
		})

		It("includes spinner frame in view", func() {
			view := model.View()
			Expect(len(view)).To(BeNumerically(">", len("Loading...")))
		})
	})
})

var _ = Describe("progressModel", func() {
	var (
		progressChan chan int
		resultChan   chan error
		model        progressModel
	)

	BeforeEach(func() {
		progressChan = make(chan int)
		resultChan = make(chan error, 1)
		model = newProgressModel("Processing...", 100, progressChan, resultChan)
	})

	Describe("Construction", func() {
		It("initializes with correct message", func() {
			Expect(model.message).To(Equal("Processing..."))
		})

		It("initializes with correct total", func() {
			Expect(model.total).To(Equal(100))
		})

		It("initializes with current=0", func() {
			Expect(model.current).To(Equal(0))
		})

		It("initializes with progress channel", func() {
			Expect(model.progressChan).To(Equal(progressChan))
		})

		It("initializes with result channel", func() {
			Expect(model.resultChan).To(Equal(resultChan))
		})

		It("initializes with done=false", func() {
			Expect(model.done).To(BeFalse())
		})

		It("initializes with err=nil", func() {
			Expect(model.err).NotTo(HaveOccurred())
		})

		It("initializes with quitting=false", func() {
			Expect(model.quitting).To(BeFalse())
		})

		It("initializes with theme", func() {
			Expect(model.th).NotTo(BeNil())
		})
	})

	Describe("Init", func() {
		It("returns a batch command", func() {
			cmd := model.Init()
			Expect(cmd).NotTo(BeNil())
		})

		It("returns a command that can be executed", func() {
			cmd := model.Init()
			msg := cmd()
			Expect(msg).NotTo(BeNil())
		})
	})

	Describe("Update", func() {
		Context("with progressMsg", func() {
			It("updates current value", func() {
				updated, _ := model.Update(progressMsg(50))
				result := updated.(progressModel)
				Expect(result.current).To(Equal(50))
			})

			It("returns waitForProgress command", func() {
				_, cmd := model.Update(progressMsg(50))
				Expect(cmd).NotTo(BeNil())
			})

			It("handles zero progress", func() {
				updated, _ := model.Update(progressMsg(0))
				result := updated.(progressModel)
				Expect(result.current).To(Equal(0))
			})

			It("handles progress equal to total", func() {
				updated, _ := model.Update(progressMsg(100))
				result := updated.(progressModel)
				Expect(result.current).To(Equal(100))
			})

			It("handles progress exceeding total", func() {
				updated, _ := model.Update(progressMsg(150))
				result := updated.(progressModel)
				Expect(result.current).To(Equal(150))
			})
		})

		Context("with progressDoneMsg", func() {
			It("marks done=true", func() {
				updated, _ := model.Update(progressDoneMsg{err: nil})
				result := updated.(progressModel)
				Expect(result.done).To(BeTrue())
			})

			It("stores error from message", func() {
				testErr := errors.New("progress error")
				updated, _ := model.Update(progressDoneMsg{err: testErr})
				result := updated.(progressModel)
				Expect(result.err).To(Equal(testErr))
			})

			It("returns tea.Quit command", func() {
				_, cmd := model.Update(progressDoneMsg{err: nil})
				Expect(cmd).NotTo(BeNil())
				msg := cmd()
				_, isQuit := msg.(tea.QuitMsg)
				Expect(isQuit).To(BeTrue())
			})

			It("handles nil error", func() {
				updated, _ := model.Update(progressDoneMsg{err: nil})
				result := updated.(progressModel)
				Expect(result.err).NotTo(HaveOccurred())
			})
		})

		Context("with ctrl+c key", func() {
			It("marks quitting=true", func() {
				updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
				result := updated.(progressModel)
				Expect(result.quitting).To(BeTrue())
			})

			It("returns tea.Quit command", func() {
				_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
				Expect(cmd).NotTo(BeNil())
				msg := cmd()
				_, isQuit := msg.(tea.QuitMsg)
				Expect(isQuit).To(BeTrue())
			})
		})

		Context("with other key messages", func() {
			It("ignores other keys", func() {
				updated, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
				result := updated.(progressModel)
				Expect(result.quitting).To(BeFalse())
				Expect(cmd).To(BeNil())
			})
		})

		Context("with unknown message type", func() {
			It("returns model unchanged", func() {
				updated, cmd := model.Update("unknown")
				result := updated.(progressModel)
				Expect(result).To(Equal(model))
				Expect(cmd).To(BeNil())
			})
		})
	})

	Describe("View", func() {
		It("returns empty string when done", func() {
			model.done = true
			view := model.View()
			Expect(view).To(Equal(""))
		})

		It("returns empty string when quitting", func() {
			model.quitting = true
			view := model.View()
			Expect(view).To(Equal(""))
		})

		It("returns progress bar and message when active", func() {
			view := model.View()
			Expect(view).To(ContainSubstring("Processing..."))
		})

		It("shows 0% when current is 0", func() {
			model.current = 0
			view := model.View()
			Expect(view).To(ContainSubstring("0%"))
		})

		It("shows 50% when current is 50 and total is 100", func() {
			model.current = 50
			view := model.View()
			Expect(view).To(ContainSubstring("50%"))
		})

		It("shows 100% when current equals total", func() {
			model.current = 100
			view := model.View()
			Expect(view).To(ContainSubstring("100%"))
		})

		It("handles zero total gracefully", func() {
			model.total = 0
			model.current = 0
			view := model.View()
			Expect(view).To(ContainSubstring("0%"))
		})

		It("includes message in view", func() {
			view := model.View()
			Expect(view).To(ContainSubstring("Processing..."))
		})
	})
})

var _ = Describe("Message Types", func() {
	Describe("spinnerTickMsg", func() {
		It("can be created from time.Time", func() {
			now := time.Now()
			msg := spinnerTickMsg(now)
			Expect(time.Time(msg)).To(Equal(now))
		})
	})

	Describe("spinnerDoneMsg", func() {
		It("can store nil error", func() {
			msg := spinnerDoneMsg{err: nil}
			Expect(msg.err).NotTo(HaveOccurred())
		})

		It("can store error", func() {
			err := errors.New("test")
			msg := spinnerDoneMsg{err: err}
			Expect(msg.err).To(Equal(err))
		})
	})

	Describe("progressMsg", func() {
		It("can be created from int", func() {
			msg := progressMsg(42)
			Expect(int(msg)).To(Equal(42))
		})
	})

	Describe("progressDoneMsg", func() {
		It("can store nil error", func() {
			msg := progressDoneMsg{err: nil}
			Expect(msg.err).NotTo(HaveOccurred())
		})

		It("can store error", func() {
			err := errors.New("test")
			msg := progressDoneMsg{err: err}
			Expect(msg.err).To(Equal(err))
		})
	})
})

var _ = Describe("Command Functions", func() {
	Describe("spinnerModel.tick", func() {
		It("returns a command", func() {
			model := newSpinnerModel("test", make(chan error, 1))
			cmd := model.tick()
			Expect(cmd).NotTo(BeNil())
		})

		It("command returns spinnerTickMsg", func() {
			model := newSpinnerModel("test", make(chan error, 1))
			cmd := model.tick()
			msg := cmd()
			_, ok := msg.(spinnerTickMsg)
			Expect(ok).To(BeTrue())
		})
	})

	Describe("spinnerModel.waitForResult", func() {
		It("returns a command", func() {
			resultChan := make(chan error, 1)
			model := newSpinnerModel("test", resultChan)
			cmd := model.waitForResult()
			Expect(cmd).NotTo(BeNil())
		})

		It("command blocks until result available", func() {
			resultChan := make(chan error, 1)
			model := newSpinnerModel("test", resultChan)
			cmd := model.waitForResult()

			go func() {
				time.Sleep(10 * time.Millisecond)
				resultChan <- nil
			}()

			msg := cmd()
			_, ok := msg.(spinnerDoneMsg)
			Expect(ok).To(BeTrue())
		})

		It("command returns spinnerDoneMsg with error", func() {
			resultChan := make(chan error, 1)
			testErr := errors.New("test error")
			resultChan <- testErr

			model := newSpinnerModel("test", resultChan)
			cmd := model.waitForResult()
			msg := cmd()

			doneMsg, ok := msg.(spinnerDoneMsg)
			Expect(ok).To(BeTrue())
			Expect(doneMsg.err).To(Equal(testErr))
		})
	})

	Describe("progressModel.waitForProgress", func() {
		It("returns a command", func() {
			progressChan := make(chan int)
			resultChan := make(chan error, 1)
			model := newProgressModel("test", 100, progressChan, resultChan)
			cmd := model.waitForProgress()
			Expect(cmd).NotTo(BeNil())
		})

		It("command blocks until progress available", func() {
			progressChan := make(chan int)
			resultChan := make(chan error, 1)
			model := newProgressModel("test", 100, progressChan, resultChan)
			cmd := model.waitForProgress()

			go func() {
				time.Sleep(10 * time.Millisecond)
				progressChan <- 50
			}()

			msg := cmd()
			progressMsg, ok := msg.(progressMsg)
			Expect(ok).To(BeTrue())
			Expect(int(progressMsg)).To(Equal(50))
		})

		It("command returns nil when channel closed", func() {
			progressChan := make(chan int)
			resultChan := make(chan error, 1)
			model := newProgressModel("test", 100, progressChan, resultChan)
			cmd := model.waitForProgress()

			close(progressChan)
			msg := cmd()
			Expect(msg).To(BeNil())
		})
	})

	Describe("progressModel.waitForResult", func() {
		It("returns a command", func() {
			progressChan := make(chan int)
			resultChan := make(chan error, 1)
			model := newProgressModel("test", 100, progressChan, resultChan)
			cmd := model.waitForResult()
			Expect(cmd).NotTo(BeNil())
		})

		It("command blocks until result available", func() {
			progressChan := make(chan int)
			resultChan := make(chan error, 1)
			model := newProgressModel("test", 100, progressChan, resultChan)
			cmd := model.waitForResult()

			go func() {
				time.Sleep(10 * time.Millisecond)
				resultChan <- nil
			}()

			msg := cmd()
			_, ok := msg.(progressDoneMsg)
			Expect(ok).To(BeTrue())
		})

		It("command returns progressDoneMsg with error", func() {
			progressChan := make(chan int)
			resultChan := make(chan error, 1)
			testErr := errors.New("test error")
			resultChan <- testErr

			model := newProgressModel("test", 100, progressChan, resultChan)
			cmd := model.waitForResult()
			msg := cmd()

			doneMsg, ok := msg.(progressDoneMsg)
			Expect(ok).To(BeTrue())
			Expect(doneMsg.err).To(Equal(testErr))
		})
	})
})

var _ = Describe("Model Lifecycle Tests", func() {
	Describe("spinnerModel full lifecycle", func() {
		It("transitions from active to done", func() {
			resultChan := make(chan error, 1)
			model := newSpinnerModel("Loading", resultChan)

			Expect(model.done).To(BeFalse())
			Expect(model.quitting).To(BeFalse())

			updated, cmd := model.Update(spinnerTickMsg(time.Now()))
			result := updated.(spinnerModel)
			Expect(result.done).To(BeFalse())
			Expect(cmd).NotTo(BeNil())

			updated, cmd = result.Update(spinnerDoneMsg{err: nil})
			result = updated.(spinnerModel)
			Expect(result.done).To(BeTrue())
			Expect(cmd).NotTo(BeNil())

			Expect(result.View()).To(Equal(""))
		})

		It("transitions from active to quitting", func() {
			resultChan := make(chan error, 1)
			model := newSpinnerModel("Loading", resultChan)

			updated, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
			result := updated.(spinnerModel)
			Expect(result.quitting).To(BeTrue())
			Expect(cmd).NotTo(BeNil())

			Expect(result.View()).To(Equal(""))
		})
	})

	Describe("progressModel full lifecycle", func() {
		It("transitions through progress updates to done", func() {
			progressChan := make(chan int)
			resultChan := make(chan error, 1)
			model := newProgressModel("Processing", 100, progressChan, resultChan)

			Expect(model.current).To(Equal(0))
			Expect(model.done).To(BeFalse())

			updated, cmd := model.Update(progressMsg(25))
			result := updated.(progressModel)
			Expect(result.current).To(Equal(25))
			Expect(cmd).NotTo(BeNil())

			updated, _ = result.Update(progressMsg(50))
			result = updated.(progressModel)
			Expect(result.current).To(Equal(50))

			updated, _ = result.Update(progressDoneMsg{err: nil})
			result = updated.(progressModel)
			Expect(result.done).To(BeTrue())
			Expect(result.View()).To(Equal(""))
		})

		It("transitions from active to quitting", func() {
			progressChan := make(chan int)
			resultChan := make(chan error, 1)
			model := newProgressModel("Processing", 100, progressChan, resultChan)

			updated, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
			result := updated.(progressModel)
			Expect(result.quitting).To(BeTrue())
			Expect(cmd).NotTo(BeNil())

			Expect(result.View()).To(Equal(""))
		})
	})

	Describe("Error propagation", func() {
		It("spinner propagates errors through model", func() {
			resultChan := make(chan error, 1)
			testErr := errors.New("spinner error")
			model := newSpinnerModel("Loading", resultChan)

			updated, _ := model.Update(spinnerDoneMsg{err: testErr})
			result := updated.(spinnerModel)
			Expect(result.err).To(Equal(testErr))
		})

		It("progress propagates errors through model", func() {
			progressChan := make(chan int)
			resultChan := make(chan error, 1)
			testErr := errors.New("progress error")
			model := newProgressModel("Processing", 100, progressChan, resultChan)

			updated, _ := model.Update(progressDoneMsg{err: testErr})
			result := updated.(progressModel)
			Expect(result.err).To(Equal(testErr))
		})
	})
})

var _ = Describe("Print Utilities", func() {
	Describe("PrintSuccess", func() {
		It("can be called with a message", func() {
			Expect(func() {
				PrintSuccess("Operation completed")
			}).NotTo(Panic())
		})

		It("handles empty message", func() {
			Expect(func() {
				PrintSuccess("")
			}).NotTo(Panic())
		})

		It("handles long messages", func() {
			longMsg := "This is a very long message that contains a lot of text and should be handled correctly by the print functions without any issues or truncation"
			Expect(func() {
				PrintSuccess(longMsg)
			}).NotTo(Panic())
		})

		It("handles special characters", func() {
			Expect(func() {
				PrintSuccess("✓ Success with special chars: @#$%")
			}).NotTo(Panic())
		})
	})

	Describe("PrintError", func() {
		It("can be called with a message", func() {
			Expect(func() {
				PrintError("An error occurred")
			}).NotTo(Panic())
		})

		It("handles empty message", func() {
			Expect(func() {
				PrintError("")
			}).NotTo(Panic())
		})

		It("handles long messages", func() {
			longMsg := "This is a very long error message that contains a lot of text and should be handled correctly by the print functions without any issues or truncation"
			Expect(func() {
				PrintError(longMsg)
			}).NotTo(Panic())
		})

		It("handles special characters", func() {
			Expect(func() {
				PrintError("✗ Error with special chars: @#$%")
			}).NotTo(Panic())
		})
	})

	Describe("PrintInfo", func() {
		It("can be called with a message", func() {
			Expect(func() {
				PrintInfo("Information message")
			}).NotTo(Panic())
		})

		It("handles empty message", func() {
			Expect(func() {
				PrintInfo("")
			}).NotTo(Panic())
		})

		It("handles long messages", func() {
			longMsg := "This is a very long info message that contains a lot of text and should be handled correctly by the print functions without any issues or truncation"
			Expect(func() {
				PrintInfo(longMsg)
			}).NotTo(Panic())
		})

		It("handles special characters", func() {
			Expect(func() {
				PrintInfo("ℹ Info with special chars: @#$%")
			}).NotTo(Panic())
		})
	})
})

var _ = Describe("Integration Tests - Wrapper Functions", func() {
	Describe("RunWithSpinner", func() {
		It("executes function successfully with nil error", func() {
			executed := false
			err := RunWithSpinner("Loading", func() error {
				executed = true
				return nil
			}, tea.WithInput(nil))

			Expect(err).NotTo(HaveOccurred())
			Expect(executed).To(BeTrue())
		})

		It("returns error from function", func() {
			testErr := errors.New("function failed")
			executed := false
			err := RunWithSpinner("Loading", func() error {
				executed = true
				return testErr
			}, tea.WithInput(nil))

			Expect(err).To(Equal(testErr))
			Expect(executed).To(BeTrue())
		})

		It("handles immediate completion", func() {
			executed := false
			err := RunWithSpinner("Quick task", func() error {
				executed = true
				return nil
			}, tea.WithInput(nil))

			Expect(err).NotTo(HaveOccurred())
			Expect(executed).To(BeTrue())
		})

		It("handles immediate error", func() {
			testErr := errors.New("immediate error")
			executed := false
			err := RunWithSpinner("Quick task", func() error {
				executed = true
				return testErr
			}, tea.WithInput(nil))

			Expect(err).To(Equal(testErr))
			Expect(executed).To(BeTrue())
		})

		It("executes function in goroutine", func() {
			executed := false
			err := RunWithSpinner("Async task", func() error {
				executed = true
				time.Sleep(10 * time.Millisecond)
				return nil
			}, tea.WithInput(nil))

			Expect(err).NotTo(HaveOccurred())
			Expect(executed).To(BeTrue())
		})

		It("handles error in function", func() {
			// Test that errors from the function are properly propagated
			testErr := errors.New("function error")
			executed := false
			err := RunWithSpinner("Error task", func() error {
				executed = true
				return testErr
			}, tea.WithInput(nil))

			Expect(err).To(Equal(testErr))
			Expect(executed).To(BeTrue())
		})

		It("handles empty message", func() {
			err := RunWithSpinner("", func() error {
				return nil
			}, tea.WithInput(nil))

			Expect(err).NotTo(HaveOccurred())
		})

		It("handles long message", func() {
			longMsg := "This is a very long spinner message that should be displayed correctly"
			err := RunWithSpinner(longMsg, func() error {
				return nil
			}, tea.WithInput(nil))

			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("RunWithProgress", func() {
		It("executes function successfully with progress updates", func() {
			executed := false
			progressUpdates := 0
			err := RunWithProgress("Processing", 100, func(update func(int)) error {
				executed = true
				update(25)
				progressUpdates++
				update(50)
				progressUpdates++
				update(75)
				progressUpdates++
				update(100)
				progressUpdates++
				return nil
			}, tea.WithInput(nil))

			Expect(err).NotTo(HaveOccurred())
			Expect(executed).To(BeTrue())
			Expect(progressUpdates).To(Equal(4))
		})

		It("returns error from function", func() {
			testErr := errors.New("progress function failed")
			executed := false
			err := RunWithProgress("Processing", 100, func(update func(int)) error {
				executed = true
				update(50)
				return testErr
			}, tea.WithInput(nil))

			Expect(err).To(Equal(testErr))
			Expect(executed).To(BeTrue())
		})

		It("handles immediate completion", func() {
			executed := false
			err := RunWithProgress("Quick progress", 100, func(update func(int)) error {
				executed = true
				return nil
			}, tea.WithInput(nil))

			Expect(err).NotTo(HaveOccurred())
			Expect(executed).To(BeTrue())
		})

		It("handles partial progress on error", func() {
			testErr := errors.New("error after partial progress")
			executed := false
			progressUpdates := 0
			err := RunWithProgress("Partial", 100, func(update func(int)) error {
				executed = true
				update(25)
				progressUpdates++
				update(50)
				progressUpdates++
				return testErr
			}, tea.WithInput(nil))

			Expect(err).To(Equal(testErr))
			Expect(executed).To(BeTrue())
			Expect(progressUpdates).To(Equal(2))
		})

		It("handles zero total edge case", func() {
			executed := false
			err := RunWithProgress("Zero total", 0, func(update func(int)) error {
				executed = true
				update(0)
				return nil
			}, tea.WithInput(nil))

			Expect(err).NotTo(HaveOccurred())
			Expect(executed).To(BeTrue())
		})

		It("handles negative total edge case", func() {
			executed := false
			err := RunWithProgress("Negative total", -1, func(update func(int)) error {
				executed = true
				update(0)
				return nil
			}, tea.WithInput(nil))

			Expect(err).NotTo(HaveOccurred())
			Expect(executed).To(BeTrue())
		})

		It("handles progress exceeding total", func() {
			executed := false
			err := RunWithProgress("Overflow", 100, func(update func(int)) error {
				executed = true
				update(50)
				update(150)
				return nil
			}, tea.WithInput(nil))

			Expect(err).NotTo(HaveOccurred())
			Expect(executed).To(BeTrue())
		})

		It("handles rapid progress updates", func() {
			executed := false
			progressUpdates := 0
			err := RunWithProgress("Rapid", 100, func(update func(int)) error {
				executed = true
				for i := 0; i <= 100; i += 10 {
					update(i)
					progressUpdates++
				}
				return nil
			}, tea.WithInput(nil))

			Expect(err).NotTo(HaveOccurred())
			Expect(executed).To(BeTrue())
			Expect(progressUpdates).To(Equal(11))
		})

		It("handles empty message", func() {
			err := RunWithProgress("", 100, func(update func(int)) error {
				update(50)
				return nil
			}, tea.WithInput(nil))

			Expect(err).NotTo(HaveOccurred())
		})

		It("handles long message", func() {
			longMsg := "This is a very long progress message that should be displayed correctly during progress updates"
			err := RunWithProgress(longMsg, 100, func(update func(int)) error {
				update(50)
				return nil
			}, tea.WithInput(nil))

			Expect(err).NotTo(HaveOccurred())
		})

		It("handles large total value", func() {
			executed := false
			err := RunWithProgress("Large total", 1000000, func(update func(int)) error {
				executed = true
				update(500000)
				return nil
			}, tea.WithInput(nil))

			Expect(err).NotTo(HaveOccurred())
			Expect(executed).To(BeTrue())
		})

		It("handles no progress updates", func() {
			executed := false
			err := RunWithProgress("No updates", 100, func(update func(int)) error {
				executed = true
				// Don't call update at all
				return nil
			}, tea.WithInput(nil))

			Expect(err).NotTo(HaveOccurred())
			Expect(executed).To(BeTrue())
		})
	})

	Describe("RunWithSpinner and RunWithProgress error handling", func() {
		It("spinner handles wrapped error", func() {
			wrappedErr := errors.New("wrapped error")
			err := RunWithSpinner("Wrapped", func() error {
				return wrappedErr
			}, tea.WithInput(nil))

			Expect(err).To(Equal(wrappedErr))
		})

		It("progress handles wrapped error", func() {
			wrappedErr := errors.New("wrapped error")
			err := RunWithProgress("Wrapped", 100, func(update func(int)) error {
				return wrappedErr
			}, tea.WithInput(nil))

			Expect(err).To(Equal(wrappedErr))
		})

		It("spinner handles nil error explicitly", func() {
			err := RunWithSpinner("Nil error", func() error {
				return nil
			}, tea.WithInput(nil))

			Expect(err).NotTo(HaveOccurred())
		})

		It("progress handles nil error explicitly", func() {
			err := RunWithProgress("Nil error", 100, func(update func(int)) error {
				return nil
			}, tea.WithInput(nil))

			Expect(err).NotTo(HaveOccurred())
		})
	})
})

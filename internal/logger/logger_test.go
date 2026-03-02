package logger

import (
	"bytes"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestLogger(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Logger Suite")
}

var _ = Describe("Logger", func() {
	Describe("Logging Levels", func() {
		var (
			buf    bytes.Buffer
			logger *Logger
		)

		testCases := []struct {
			name           string
			logLevel       LogLevel
			logFunc        func(l *Logger, format string, args ...interface{})
			logMessage     string
			expectedOutput bool
		}{
			{
				name:           "Debug logging at Debug level",
				logLevel:       DebugLevel,
				logFunc:        func(l *Logger, format string, args ...interface{}) { l.Debug(format, args...) },
				logMessage:     "Debug message",
				expectedOutput: true,
			},
			{
				name:           "Info logging at Info level",
				logLevel:       InfoLevel,
				logFunc:        func(l *Logger, format string, args ...interface{}) { l.Info(format, args...) },
				logMessage:     "Info message",
				expectedOutput: true,
			},
			{
				name:           "Warn logging at Warn level",
				logLevel:       WarnLevel,
				logFunc:        func(l *Logger, format string, args ...interface{}) { l.Warn(format, args...) },
				logMessage:     "Warning message",
				expectedOutput: true,
			},
			{
				name:           "Error logging at Error level",
				logLevel:       ErrorLevel,
				logFunc:        func(l *Logger, format string, args ...interface{}) { l.Error(format, args...) },
				logMessage:     "Error message",
				expectedOutput: true,
			},
			{
				name:           "Debug logging at Info level (should not log)",
				logLevel:       InfoLevel,
				logFunc:        func(l *Logger, format string, args ...interface{}) { l.Debug(format, args...) },
				logMessage:     "Debug message",
				expectedOutput: false,
			},
		}

		BeforeEach(func() {
			buf.Reset()
		})

		for _, tc := range testCases {
			It(tc.name, func() {
				logger = New(&buf, tc.logLevel)
				tc.logFunc(logger, tc.logMessage)
				logOutput := buf.String()

				if tc.expectedOutput {
					Expect(logOutput).To(ContainSubstring(tc.logMessage), "Log message should be present")
				} else {
					Expect(logOutput).To(BeEmpty(), "Log message should not be present")
				}
			})
		}
	})

	Describe("Logger Context", func() {
		var (
			buf    bytes.Buffer
			logger *Logger
		)

		BeforeEach(func() {
			buf.Reset()
			logger = New(&buf, DebugLevel)
		})

		It("should add and log context", func() {
			logger.SetContext("user_id", "123")
			logger.SetContext("operation", "test")
			logger.Info("Test log message")

			logOutput := buf.String()
			Expect(logOutput).To(ContainSubstring("user_id=123"))
			Expect(logOutput).To(ContainSubstring("operation=test"))
		})

		It("should create logger with fields", func() {
			loggerWithFields := logger.WithFields(map[string]string{
				"user_id": "456",
				"action":  "create",
			})

			loggerWithFields.Info("Test log message")

			logOutput := buf.String()
			Expect(logOutput).To(ContainSubstring("user_id=456"))
			Expect(logOutput).To(ContainSubstring("action=create"))
		})

		It("should clear context", func() {
			logger.SetContext("user_id", "123")
			logger.SetContext("operation", "test")
			logger.ClearContext()
			logger.Info("Test log message")

			logOutput := buf.String()
			Expect(logOutput).NotTo(ContainSubstring("user_id=123"))
			Expect(logOutput).NotTo(ContainSubstring("operation=test"))
		})
	})

	Describe("Log Level String Representation", func() {
		testCases := []struct {
			level    LogLevel
			expected string
		}{
			{DebugLevel, "DEBUG"},
			{InfoLevel, "INFO"},
			{WarnLevel, "WARN"},
			{ErrorLevel, "ERROR"},
			{FatalLevel, "FATAL"},
		}

		for _, tc := range testCases {
			It(tc.expected+" should have correct string representation", func() {
				Expect(tc.level.String()).To(Equal(tc.expected))
			})
		}

		It("returns UNKNOWN for undefined log levels", func() {
			unknownLevel := LogLevel(99)
			Expect(unknownLevel.String()).To(Equal("UNKNOWN"))
		})
	})

	Describe("Constructor Functions", func() {
		It("creates a console logger at Info level", func() {
			l := ConsoleLogger()
			Expect(l).NotTo(BeNil())
			Expect(l.level).To(Equal(InfoLevel))
		})

		It("creates a default logger", func() {
			l := DefaultLogger()
			Expect(l).NotTo(BeNil())
			Expect(l.level).To(Equal(InfoLevel))
		})

		It("creates a file logger", func() {
			l := FileLogger()
			Expect(l).NotTo(BeNil())
			Expect(l.level).To(Equal(InfoLevel))
		})
	})

	Describe("WithFields with existing context", func() {
		It("preserves existing context when adding new fields", func() {
			var buf bytes.Buffer
			l := New(&buf, DebugLevel)
			l.SetContext("existing_key", "existing_value")

			newLogger := l.WithFields(map[string]string{
				"new_key": "new_value",
			})

			newLogger.Info("test message")
			output := buf.String()
			Expect(output).To(ContainSubstring("existing_key=existing_value"))
			Expect(output).To(ContainSubstring("new_key=new_value"))
		})
	})
})

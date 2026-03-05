package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestRootCommand(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		expectError    bool
	}{
		{
			name:           "version flag long",
			args:           []string{"--version"},
			expectedOutput: "kariya version",
			expectError:    false,
		},
		{
			name:           "version flag short",
			args:           []string{"-v"},
			expectedOutput: "kariya version",
			expectError:    false,
		},
		{
			name:           "help flag long",
			args:           []string{"--help"},
			expectedOutput: "KaRiya is a terminal user interface",
			expectError:    false,
		},
		{
			name:           "help flag short",
			args:           []string{"-h"},
			expectedOutput: "KaRiya is a terminal user interface",
			expectError:    false,
		},
		{
			name:           "no args shows help",
			args:           []string{},
			expectedOutput: "KaRiya is a terminal user interface",
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewRootCmd()
			out := new(bytes.Buffer)
			cmd.SetOut(out)
			cmd.SetErr(out)
			cmd.SetArgs(tt.args)

			err := cmd.Execute()
			if tt.expectError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			output := out.String()
			if !strings.Contains(output, tt.expectedOutput) {
				t.Errorf("output %q does not contain %q", output, tt.expectedOutput)
			}
		})
	}
}

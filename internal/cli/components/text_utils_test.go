package components

import (
	"testing"
)

func TestTruncateText(t *testing.T) {
	tests := []struct {
		name   string
		text   string
		maxLen int
		want   string
	}{
		{
			name:   "Text shorter than max",
			text:   "Hello",
			maxLen: 10,
			want:   "Hello",
		},
		{
			name:   "Text exactly at max",
			text:   "Hello World",
			maxLen: 11,
			want:   "Hello World",
		},
		{
			name:   "Text longer than max",
			text:   "Hello World, this is a long text",
			maxLen: 20,
			want:   "Hello World, this...", // Truncated at byte 17 + "..."
		},
		{
			name:   "Empty text",
			text:   "",
			maxLen: 10,
			want:   "",
		},
		{
			name:   "MaxLen of 3",
			text:   "Hello",
			maxLen: 3,
			want:   "...", // maxLen 3 means ellipsis only (3-3=0 chars)
		},
		{
			name:   "MaxLen of 4",
			text:   "Hello World",
			maxLen: 4,
			want:   "H...",
		},
		{
			name:   "MaxLen very small (1)",
			text:   "Hello",
			maxLen: 1,
			want:   "H",
		},
		{
			name:   "MaxLen very small (2)",
			text:   "Hello",
			maxLen: 2,
			want:   "He",
		},
		{
			name:   "Unicode text",
			text:   "Hello 世界", // "世界" is 6 bytes, total = 5+1+6 = 12 bytes
			maxLen: 10,
			want:   "Hello \xe4...", // Truncates mid-unicode char (byte-based)
		},
		{
			name:   "Unicode text truncated",
			text:   "Hello World 世界 extra text",
			maxLen: 15,
			want:   "Hello World ...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TruncateText(tt.text, tt.maxLen)
			if got != tt.want {
				t.Errorf("TruncateText(%q, %d) = %q, want %q", tt.text, tt.maxLen, got, tt.want)
			}
		})
	}
}

func TestTruncateText_Length(t *testing.T) {
	tests := []struct {
		name   string
		text   string
		maxLen int
	}{
		{"Normal truncation", "This is a very long text that needs truncation", 20},
		{"Small maxLen", "Short", 3},
		{"Large text", string(make([]byte, 1000)), 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TruncateText(tt.text, tt.maxLen)
			if len(result) > tt.maxLen {
				t.Errorf("TruncateText result length %d exceeds maxLen %d", len(result), tt.maxLen)
			}
		})
	}
}

func TestTruncateText_EdgeCases(t *testing.T) {
	t.Run("Zero maxLen", func(t *testing.T) {
		result := TruncateText("Hello", 0)
		if len(result) != 0 {
			t.Errorf("Expected empty string for maxLen 0, got %q", result)
		}
	})

	t.Run("Exactly 3 chars with long text", func(t *testing.T) {
		result := TruncateText("Hello World", 3)
		if result != "..." {
			t.Errorf("Expected '...' for maxLen 3, got %q", result)
		}
	})

	t.Run("Exactly 4 chars with long text", func(t *testing.T) {
		result := TruncateText("Hello World", 4)
		if result != "H..." {
			t.Errorf("Expected 'H...' for maxLen 4, got %q", result)
		}
	})
}

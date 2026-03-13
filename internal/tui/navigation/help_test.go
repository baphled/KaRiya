package navigation

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Navigation Help Text", func() {
	Describe("GetHelpText", func() {
		It("should return empty string for empty keys", func() {
			result := GetHelpText([]NavigationKey{})
			Expect(result).To(Equal(""))
		})

		It("should format single key with full format", func() {
			result := GetHelpText([]NavigationKey{KeyBack})
			Expect(result).To(ContainSubstring("Esc"))
			Expect(result).To(ContainSubstring("Go back to previous screen"))
			Expect(result).To(ContainSubstring("→"))
		})

		It("should format multiple keys with newline separation", func() {
			result := GetHelpText([]NavigationKey{KeyUp, KeyDown, KeySelect})
			Expect(result).To(ContainSubstring("↑/k"))
			Expect(result).To(ContainSubstring("↓/j"))
			Expect(result).To(ContainSubstring("Enter"))
			lines := strings.Count(result, "\n")
			Expect(lines).To(Equal(2))
		})

		It("should include descriptions for all keys", func() {
			result := GetHelpText([]NavigationKey{KeyBack, KeySelect})
			Expect(result).To(ContainSubstring("Go back to previous screen"))
			Expect(result).To(ContainSubstring("Confirm selection or submit"))
		})
	})

	Describe("GetHelpTextCompact", func() {
		It("should return empty string for empty keys", func() {
			result := GetHelpTextCompact([]NavigationKey{}, true)
			Expect(result).To(Equal(""))
		})

		It("should format single key in compact format", func() {
			result := GetHelpTextCompact([]NavigationKey{KeyBack}, true)
			Expect(result).To(ContainSubstring("Esc:"))
			Expect(result).To(ContainSubstring("Go back to previous screen"))
			Expect(result).NotTo(ContainSubstring("→"))
		})

		It("should use pipe separator in compact format", func() {
			result := GetHelpTextCompact([]NavigationKey{KeyUp, KeyDown}, true)
			Expect(result).To(ContainSubstring("|"))
		})

		It("should use arrow separator in full format", func() {
			result := GetHelpTextCompact([]NavigationKey{KeyUp, KeyDown}, false)
			Expect(result).To(ContainSubstring("→"))
			Expect(result).NotTo(ContainSubstring("|"))
		})

		It("should use newlines in full format with multiple keys", func() {
			result := GetHelpTextCompact([]NavigationKey{KeyUp, KeyDown, KeyLeft}, false)
			lines := strings.Count(result, "\n")
			Expect(lines).To(Equal(2))
		})
	})

	Describe("GetContextualHelp", func() {
		It("should return form-specific help for form context", func() {
			result := GetContextualHelp("form")
			Expect(result).To(ContainSubstring("↑/k"))
			Expect(result).To(ContainSubstring("↓/j"))
			Expect(result).To(ContainSubstring("Enter"))
			Expect(result).To(ContainSubstring("Esc"))
		})

		It("should return list-specific help for list context", func() {
			result := GetContextualHelp("list")
			Expect(result).To(ContainSubstring("↑/k"))
			Expect(result).To(ContainSubstring("↓/j"))
			Expect(result).To(ContainSubstring("e:"))
			Expect(result).To(ContainSubstring("d:"))
		})

		It("should return metadata_review-specific help", func() {
			result := GetContextualHelp("metadata_review")
			Expect(result).To(ContainSubstring("↑/k"))
			Expect(result).To(ContainSubstring("↓/j"))
			Expect(result).To(ContainSubstring("e:"))
			Expect(result).To(ContainSubstring("b:"))
		})

		It("should return home-specific help", func() {
			result := GetContextualHelp("home")
			Expect(result).To(ContainSubstring("c:"))
			Expect(result).To(ContainSubstring("l:"))
			Expect(result).To(ContainSubstring("m:"))
			Expect(result).To(ContainSubstring("?:"))
		})

		It("should return default help for unknown context", func() {
			result := GetContextualHelp("unknown_context")
			Expect(result).To(ContainSubstring("↑/k"))
			Expect(result).To(ContainSubstring("↓/j"))
			Expect(result).To(ContainSubstring("Enter"))
		})

		It("should use compact format for contextual help", func() {
			result := GetContextualHelp("form")
			Expect(result).To(ContainSubstring("|"))
		})
	})

	Describe("GetFullHelp", func() {
		It("should return help text for all 19 navigation keys", func() {
			result := GetFullHelp()
			keys := AllNavigationKeys()
			for _, key := range keys {
				Expect(result).To(ContainSubstring(string(key)))
			}
		})

		It("should use full format", func() {
			result := GetFullHelp()
			Expect(result).To(ContainSubstring("→"))
		})

		It("should include all descriptions", func() {
			result := GetFullHelp()
			for _, description := range KeyDescription {
				Expect(result).To(ContainSubstring(description))
			}
		})
	})

	Describe("GetCompactHelp", func() {
		It("should return help text for all navigation keys", func() {
			result := GetCompactHelp()
			keys := AllNavigationKeys()
			for _, key := range keys {
				Expect(result).To(ContainSubstring(string(key)))
			}
		})

		It("should use compact format", func() {
			result := GetCompactHelp()
			Expect(result).To(ContainSubstring("|"))
		})

		It("should fit on single line", func() {
			result := GetCompactHelp()
			lines := strings.Count(result, "\n")
			Expect(lines).To(Equal(0))
		})
	})

	Describe("GetGroupedHelp", func() {
		It("should organize keys by group", func() {
			result := GetGroupedHelp()
			Expect(result).To(ContainSubstring("Navigation:"))
			Expect(result).To(ContainSubstring("Actions:"))
			Expect(result).To(ContainSubstring("Modes:"))
			Expect(result).To(ContainSubstring("Tools:"))
			Expect(result).To(ContainSubstring("Global:"))
		})

		It("should include all navigation keys", func() {
			result := GetGroupedHelp()
			keys := AllNavigationKeys()
			for _, key := range keys {
				Expect(result).To(ContainSubstring(string(key)))
			}
		})

		It("should use full format with arrows", func() {
			result := GetGroupedHelp()
			Expect(result).To(ContainSubstring("→"))
		})

		It("should include descriptions for each key", func() {
			result := GetGroupedHelp()
			for _, description := range KeyDescription {
				Expect(result).To(ContainSubstring(description))
			}
		})

		It("should have multiple lines (grouped format)", func() {
			result := GetGroupedHelp()
			lines := strings.Count(result, "\n")
			Expect(lines).To(BeNumerically(">", 10))
		})
	})

	Describe("Help text with custom keys", func() {
		It("should handle subsets of keys", func() {
			customKeys := []NavigationKey{KeyUp, KeyDown, KeySelect}
			result := GetHelpText(customKeys)
			Expect(result).To(ContainSubstring("↑/k"))
			Expect(result).To(ContainSubstring("↓/j"))
			Expect(result).To(ContainSubstring("Enter"))
			Expect(result).NotTo(ContainSubstring("Esc"))
		})

		It("should preserve key order", func() {
			customKeys := []NavigationKey{KeyBack, KeyUp, KeyDown}
			result := GetHelpText(customKeys)
			backIdx := strings.Index(result, "Esc")
			upIdx := strings.Index(result, "↑/k")
			downIdx := strings.Index(result, "↓/j")
			Expect(backIdx).To(BeNumerically("<", upIdx))
			Expect(upIdx).To(BeNumerically("<", downIdx))
		})
	})
})

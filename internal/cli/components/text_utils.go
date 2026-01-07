package components

// TruncateText truncates text to maxLen with ellipsis if needed.
// If the text is longer than maxLen, it will be truncated to maxLen-3
// characters and "..." will be appended.
func TruncateText(text string, maxLen int) string {
	if len(text) > maxLen {
		if maxLen < 3 {
			// If maxLen is too small for ellipsis, just return truncated text
			return text[:maxLen]
		}
		return text[:maxLen-3] + "..."
	}
	return text
}

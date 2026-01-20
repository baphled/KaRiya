// Package widgets provides higher-level composite components that combine
// primitives and behaviors into reusable UI patterns.
//
// # DetailView Widget
//
// DetailView renders structured key-value detail information with consistent
// styling, text wrapping, sections, and list support. It solves the problem
// of inconsistent detail screen rendering across the application.
//
// ## Features
//
//   - Field-value pairs with styled labels
//   - Automatic text wrapping when width is set
//   - Section headers for organizing related fields
//   - List rendering (comma-separated, bulleted, custom separator)
//   - Conditional fields/lists (only show if non-empty)
//   - Theme-aware styling via theme.Aware embedding
//
// ## Usage
//
// Basic usage with fields:
//
//	dv := widgets.NewDetailView(theme).
//	    Title("Event Details").
//	    Width(60).
//	    Field("Name", event.Name).
//	    Field("Date", event.Date.Format("2006-01-02")).
//	    FieldIf("Company", event.Company).  // Only if not empty
//	    List("Tags", event.Tags)
//	rendered := dv.Render()
//
// With sections:
//
//	dv := widgets.NewDetailView(theme).
//	    Title("User Profile").
//	    Section("Personal Info").
//	    Field("Name", user.Name).
//	    Field("Email", user.Email).
//	    Section("Work Info").
//	    Field("Company", user.Company).
//	    Field("Role", user.Role)
//
// With lists:
//
//	dv := widgets.NewDetailView(theme).
//	    List("Tags", []string{"go", "tui"}).           // "go, tui"
//	    ListWithSeparator("Skills", skills, " | ").     // "go | python | rust"
//	    BulletList("Steps", steps)                      // Bulleted list
//
// ## API Reference
//
// Constructors:
//
//	NewDetailView(theme Theme) *DetailView
//
// Configuration methods (chainable):
//
//	Title(title string) *DetailView       // Set title at top
//	Width(width int) *DetailView          // Set max width for wrapping
//
// Content methods (chainable):
//
//	Section(title string) *DetailView                        // Start new section
//	Field(label, value string) *DetailView                   // Add field
//	FieldIf(label, value string) *DetailView                 // Add field if value non-empty
//	List(label string, values []string) *DetailView          // Add comma-separated list
//	ListIf(label string, values []string) *DetailView        // Add list if non-empty
//	ListWithSeparator(label string, values []string, sep string) *DetailView
//	BulletList(label string, values []string) *DetailView    // Add bulleted list
//
// Rendering:
//
//	Render() string   // Generate styled output
//
// Theme integration (inherited from theme.Aware):
//
//	SetTheme(theme Theme)   // Update theme
//	Theme() Theme           // Get current theme (default if nil)
//
// ## Text Wrapping
//
// When Width() is set, long text values are automatically wrapped to fit
// within the specified width. The wrapping algorithm:
//   - Splits text on word boundaries
//   - Handles very long words gracefully (no breaking)
//   - Preserves existing line breaks in multi-line values
//   - Minimum width of 10 characters enforced
//
// ## Best Practices
//
//  1. Always set Width() for responsive layouts
//  2. Use FieldIf/ListIf for optional fields to avoid empty lines
//  3. Group related fields using Section()
//  4. Use BulletList for step-by-step instructions or long lists
//  5. Use List with custom separator for inline display of items
package widgets

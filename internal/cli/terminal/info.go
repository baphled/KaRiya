package terminal

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Info holds terminal dimension information
type Info struct {
	Width       int
	Height      int
	IsValid     bool      // Has received at least one size update
	LastUpdated time.Time // Track when last updated
}

// Config defines terminal size constraints and defaults
type Config struct {
	MinWidth      int // Minimum supported width
	MinHeight     int // Minimum supported height
	DefaultWidth  int // Default width when unknown
	DefaultHeight int // Default height when unknown
}

// DefaultConfig provides sensible defaults for terminal configuration
var DefaultConfig = Config{
	MinWidth:      40,
	MinHeight:     15,
	DefaultWidth:  80,
	DefaultHeight: 24,
}

// Margins represents spacing around content
type Margins struct {
	Top    int
	Right  int
	Bottom int
	Left   int
}

// SizeCategory represents the terminal size category for responsive design
type SizeCategory int

const (
	// SizeTiny represents terminals < 60 columns
	SizeTiny SizeCategory = iota
	// SizeCompact represents terminals 60-79 columns
	SizeCompact
	// SizeNormal represents terminals 80-119 columns
	SizeNormal
	// SizeLarge represents terminals 120-159 columns
	SizeLarge
	// SizeXLarge represents terminals >= 160 columns
	SizeXLarge
)

// Layout overhead constants for consistent UI sizing across all components.
// These values represent the vertical space consumed by UI chrome.
const (
	// StandardViewOverhead is the space used by StandardView (logo + breadcrumbs + footer).
	// Logo spacing (2) + logo (~9) + blank (1) + breadcrumbs (~2) + blank (1) + footer (~3) = ~18 lines
	StandardViewOverhead = 18

	// ModalOverhead is the additional space for modal chrome on top of StandardViewOverhead.
	// Modal border (2) + title (1) + padding (1) = ~4 lines
	ModalOverhead = 4

	// FormOverhead is the space for form chrome (header, confirm button area).
	// Confirm button group takes ~5 lines.
	FormOverhead = 5

	// TableHeaderOverhead is the space for table headers and borders.
	TableHeaderOverhead = 3

	// MinContentHeight is the minimum height for any scrollable content area.
	MinContentHeight = 10

	// MinContentWidth is the minimum width for any content area.
	MinContentWidth = 40

	// DefaultPageSize is the fallback page size when terminal size is unknown.
	DefaultPageSize = 15
)

// NewInfo creates a new Info instance with default values
func NewInfo() *Info {
	return &Info{
		Width:   0,
		Height:  0,
		IsValid: false,
	}
}

// Update updates the terminal info from a WindowSizeMsg
func (i *Info) Update(msg tea.WindowSizeMsg) {
	i.Width = msg.Width
	i.Height = msg.Height
	i.IsValid = true
	i.LastUpdated = time.Now()
}

// GetCategory returns the size category based on terminal width
func (i *Info) GetCategory() SizeCategory {
	if !i.IsValid {
		return SizeNormal // Safe default
	}

	switch {
	case i.Width < 60:
		return SizeTiny
	case i.Width < 80:
		return SizeCompact
	case i.Width < 120:
		return SizeNormal
	case i.Width < 160:
		return SizeLarge
	default:
		return SizeXLarge
	}
}

// GetSafeDimensions returns dimensions with enforced minimums and fallback to defaults
func (i *Info) GetSafeDimensions(config Config) (width, height int) {
	if !i.IsValid {
		return config.DefaultWidth, config.DefaultHeight
	}

	width = max(i.Width, config.MinWidth)
	height = max(i.Height, config.MinHeight)
	return
}

// CanRender returns true if the terminal can render content at minimum size
func (i *Info) CanRender(config Config) bool {
	if !i.IsValid {
		return true // Assume we can render with defaults
	}
	return i.Width >= config.MinWidth && i.Height >= config.MinHeight
}

// ContentArea calculates available space after accounting for margins
// Returns width and height with enforced minimums
func (i *Info) ContentArea(margins Margins) (width, height int) {
	safeWidth, safeHeight := i.GetSafeDimensions(DefaultConfig)

	width = safeWidth - margins.Left - margins.Right
	height = safeHeight - margins.Top - margins.Bottom

	// Ensure minimum content area
	width = max(width, 20)
	height = max(height, 5)

	return
}

// ContentHeight returns the available height for content after subtracting overhead.
// The overhead parameter should be the sum of all UI chrome above/below the content.
// Returns at least MinContentHeight.
func (i *Info) ContentHeight(overhead int) int {
	_, safeHeight := i.GetSafeDimensions(DefaultConfig)
	h := safeHeight - overhead
	if h < MinContentHeight {
		h = MinContentHeight
	}
	return h
}

// ContentWidth returns the available width for content after subtracting margins.
// Returns at least MinContentWidth.
func (i *Info) ContentWidth(horizontalMargin int) int {
	safeWidth, _ := i.GetSafeDimensions(DefaultConfig)
	w := safeWidth - horizontalMargin
	if w < MinContentWidth {
		w = MinContentWidth
	}
	return w
}

// FormContentHeight returns the height available for form content.
// This accounts for StandardView overhead only - forms handle their own internal chrome.
func (i *Info) FormContentHeight() int {
	return i.ContentHeight(StandardViewOverhead)
}

// ModalContentHeight returns the height available for modal content (not forms).
// This accounts for StandardView overhead and modal chrome.
func (i *Info) ModalContentHeight() int {
	return i.ContentHeight(StandardViewOverhead + ModalOverhead)
}

// ModalFormContentHeight returns the height available for a form inside a modal.
// This accounts for StandardView and modal chrome - forms handle their own internal chrome.
func (i *Info) ModalFormContentHeight() int {
	return i.ContentHeight(StandardViewOverhead + ModalOverhead)
}

// PageSize returns the dynamic page size based on available content height.
// This is useful for paginated lists and tables.
// The itemHeight parameter specifies how many lines each item takes.
func (i *Info) PageSize(itemHeight int) int {
	if itemHeight <= 0 {
		itemHeight = 1
	}

	// Available height for list items (after standard view and table header)
	availableHeight := i.ContentHeight(StandardViewOverhead + TableHeaderOverhead)

	// Calculate how many items fit
	pageSize := availableHeight / itemHeight
	if pageSize < 5 {
		pageSize = 5 // Minimum 5 items per page
	}
	if pageSize > 50 {
		pageSize = 50 // Maximum 50 items per page
	}
	return pageSize
}

// TablePageSize returns the page size optimized for table views.
// Assumes each row is 1 line tall.
func (i *Info) TablePageSize() int {
	return i.PageSize(1)
}

// max returns the larger of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

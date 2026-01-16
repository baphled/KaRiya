package containers

// Overlay provides centered modal overlay with optional background dimming.
type Overlay struct {
}

// NewOverlay creates a new overlay with the given dimensions.
func NewOverlay(width, height int) *Overlay {
	return &Overlay{}
}

// Content sets the content to display in the center.
func (o *Overlay) Content(content string) *Overlay {
	return o
}

// Dimmed enables background dimming.
func (o *Overlay) Dimmed() *Overlay {
	return o
}

// DimmedWith sets a custom dim character.
func (o *Overlay) DimmedWith(char rune) *Overlay {
	return o
}

// Render returns the rendered overlay as a string.
func (o *Overlay) Render() string {
	return ""
}

package capture

// MetadataInput holds the raw metadata fields for event metadata editing.
// This decouples metadata application from any form library.
type MetadataInput struct {
	Date       string
	Company    string
	Project    string
	Tags       []string
	Categories []string
	Skills     []string
}

// BurstEditInputFromFields creates a BurstEditInput from raw field values.
//
// Expected: Name and description are plain strings with no TUI dependencies.
// Returns: A BurstEditInput struct ready for ApplyBurstEdit.
// Side effects: None.
func BurstEditInputFromFields(name, description string) BurstEditInput {
	return BurstEditInput{
		Name:        name,
		Description: description,
	}
}

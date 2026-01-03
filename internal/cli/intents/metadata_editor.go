package intents

import (
	"context"
)

type MetadataEditorState string

const (
	MetadataReviewState MetadataEditorState = "review"
	MetadataEditState   MetadataEditorState = "edit"
	MetadataConfirmState MetadataEditorState = "confirm"
)

type MetadataEditorContext struct {
	CurrentState MetadataEditorState
	EntityType string
	EntityID string
	OriginalMetadata map[string]interface{}
	EditedMetadata map[string]interface{}
	ChangedFields map[string]bool
	FieldValues map[string]interface{}
	FormErrors map[string]string
	SelectedFieldIndex int
	ScrollPosition int
	Context context.Context
	PreviousState MetadataEditorState
}

type MetadataEditorResult struct {
	Action string
	EntityType string
	EntityID string
	Changes map[string]interface{}
	OriginalValues map[string]interface{}
	Error error
	Message string
}

func NewMetadataEditorContext(ctx context.Context) *MetadataEditorContext {
	return &MetadataEditorContext{
		CurrentState: MetadataReviewState,
		EntityType: "",
		EntityID: "",
		OriginalMetadata: make(map[string]interface{}),
		EditedMetadata: make(map[string]interface{}),
		ChangedFields: make(map[string]bool),
		FieldValues: make(map[string]interface{}),
		FormErrors: make(map[string]string),
		SelectedFieldIndex: 0,
		ScrollPosition: 0,
		Context: ctx,
		PreviousState: MetadataReviewState,
	}
}

func (c *MetadataEditorContext) ClearFormErrors() {
	c.FormErrors = make(map[string]string)
}

func (c *MetadataEditorContext) SetFormError(field, message string) {
	c.FormErrors[field] = message
}

func (c *MetadataEditorContext) HasFormErrors() bool {
	return len(c.FormErrors) > 0
}

func (c *MetadataEditorContext) SetFieldValue(field string, value interface{}) {
	c.FieldValues[field] = value
	c.ChangedFields[field] = true
	c.EditedMetadata[field] = value
}

func (c *MetadataEditorContext) GetChanges() map[string]interface{} {
	changes := make(map[string]interface{})
	for field := range c.ChangedFields {
		changes[field] = c.EditedMetadata[field]
	}
	return changes
}

func (c *MetadataEditorContext) HasChanges() bool {
	return len(c.ChangedFields) > 0
}

func (c *MetadataEditorContext) ResetChanges() {
	c.ChangedFields = make(map[string]bool)
	c.EditedMetadata = make(map[string]interface{})
	for k, v := range c.OriginalMetadata {
		c.EditedMetadata[k] = v
	}
	c.ClearFormErrors()
}

func (c *MetadataEditorContext) LoadMetadata(metadata map[string]interface{}) {
	c.OriginalMetadata = make(map[string]interface{})
	c.EditedMetadata = make(map[string]interface{})
	for k, v := range metadata {
		c.OriginalMetadata[k] = v
		c.EditedMetadata[k] = v
	}
}


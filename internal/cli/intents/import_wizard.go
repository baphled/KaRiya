package intents

import (
	"context"
	"time"
)

type ImportWizardState string

const (
	ImportFileSelectState ImportWizardState = "file_select"
	ImportPreviewState    ImportWizardState = "preview"
	ImportProgressState   ImportWizardState = "progress"
	ImportCompleteState   ImportWizardState = "complete"
)

type ImportWizardContext struct {
	CurrentState    ImportWizardState
	FilePath        string
	FileSize        int64
	TotalRows       int
	ProcessedRows   int
	SuccessfulRows  int
	ErrorRows       int
	Errors          []string
	IsImporting     bool
	IsPaused        bool
	ImportStartTime time.Time
	FormErrors      map[string]string
	Context         context.Context
	PreviousState   ImportWizardState
}

type ImportWizardResult struct {
	Action         string
	FilePath       string
	ProcessedRows  int
	SuccessfulRows int
	ErrorRows      int
	Errors         []string
	Error          error
	Message        string
}

func NewImportWizardContext(ctx context.Context) *ImportWizardContext {
	return &ImportWizardContext{
		CurrentState:   ImportFileSelectState,
		FilePath:       "",
		FileSize:       0,
		TotalRows:      0,
		ProcessedRows:  0,
		SuccessfulRows: 0,
		ErrorRows:      0,
		Errors:         make([]string, 0),
		IsImporting:    false,
		IsPaused:       false,
		FormErrors:     make(map[string]string),
		Context:        ctx,
		PreviousState:  ImportFileSelectState,
	}
}

func (c *ImportWizardContext) ClearFormErrors() {
	c.FormErrors = make(map[string]string)
}

func (c *ImportWizardContext) SetFormError(field, message string) {
	c.FormErrors[field] = message
}

func (c *ImportWizardContext) HasFormErrors() bool {
	return len(c.FormErrors) > 0
}

func (c *ImportWizardContext) StartImport() {
	c.IsImporting = true
	c.IsPaused = false
	c.ImportStartTime = time.Now()
	c.ProcessedRows = 0
	c.SuccessfulRows = 0
	c.ErrorRows = 0
	c.Errors = make([]string, 0)
}

func (c *ImportWizardContext) PauseImport() {
	c.IsPaused = true
}

func (c *ImportWizardContext) ResumeImport() {
	c.IsPaused = false
}

func (c *ImportWizardContext) CancelImport() {
	c.IsImporting = false
	c.IsPaused = false
}

func (c *ImportWizardContext) CompleteImport() {
	c.IsImporting = false
	c.IsPaused = false
}

func (c *ImportWizardContext) AddError(message string) {
	c.Errors = append(c.Errors, message)
	c.ErrorRows++
}

func (c *ImportWizardContext) IncrementProcessed() {
	c.ProcessedRows++
}

func (c *ImportWizardContext) IncrementSuccessful() {
	c.SuccessfulRows++
}

func (c *ImportWizardContext) GetProgress() float64 {
	if c.TotalRows == 0 {
		return 0
	}
	return float64(c.ProcessedRows) / float64(c.TotalRows)
}

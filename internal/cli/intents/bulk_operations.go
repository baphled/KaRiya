package intents

import (
	"context"
	"time"
)

type BulkOperationsState string

const (
	BulkSelectOpState  BulkOperationsState = "select_op"
	BulkConfigureState BulkOperationsState = "configure"
	BulkExecuteState   BulkOperationsState = "execute"
	BulkCompleteState  BulkOperationsState = "complete"
)

type BulkOperationsContext struct {
	CurrentState       BulkOperationsState
	AvailableOps       []string
	SelectedOp         string
	ScopeType          string
	AffectedItemCount  int
	ProcessedCount     int
	SuccessCount       int
	FailureCount       int
	SkippedCount       int
	CurrentItemIndex   int
	IsExecuting        bool
	IsPaused           bool
	ExecutionStartTime time.Time
	Results            map[string]string
	Errors             []string
	FormErrors         map[string]string
	Context            context.Context
	PreviousState      BulkOperationsState
}

type BulkOperationsResult struct {
	Operation      string
	ProcessedCount int
	SuccessCount   int
	FailureCount   int
	SkippedCount   int
	Results        map[string]string
	Errors         []string
	Error          error
	Message        string
}

func NewBulkOperationsContext(ctx context.Context) *BulkOperationsContext {
	return &BulkOperationsContext{
		CurrentState:      BulkSelectOpState,
		AvailableOps:      []string{"delete", "tag", "archive", "export"},
		SelectedOp:        "",
		ScopeType:         "selected",
		AffectedItemCount: 0,
		ProcessedCount:    0,
		SuccessCount:      0,
		FailureCount:      0,
		SkippedCount:      0,
		CurrentItemIndex:  0,
		IsExecuting:       false,
		IsPaused:          false,
		Results:           make(map[string]string),
		Errors:            make([]string, 0),
		FormErrors:        make(map[string]string),
		Context:           ctx,
		PreviousState:     BulkSelectOpState,
	}
}

func (c *BulkOperationsContext) ClearFormErrors() {
	c.FormErrors = make(map[string]string)
}

func (c *BulkOperationsContext) SetFormError(field, message string) {
	c.FormErrors[field] = message
}

func (c *BulkOperationsContext) HasFormErrors() bool {
	return len(c.FormErrors) > 0
}

func (c *BulkOperationsContext) SelectOperation(op string) {
	for _, availOp := range c.AvailableOps {
		if availOp == op {
			c.SelectedOp = op
			break
		}
	}
}

func (c *BulkOperationsContext) StartExecution() {
	c.IsExecuting = true
	c.IsPaused = false
	c.ExecutionStartTime = time.Now()
	c.ProcessedCount = 0
	c.SuccessCount = 0
	c.FailureCount = 0
	c.SkippedCount = 0
	c.Results = make(map[string]string)
	c.Errors = make([]string, 0)
}

func (c *BulkOperationsContext) PauseExecution() {
	c.IsPaused = true
}

func (c *BulkOperationsContext) ResumeExecution() {
	c.IsPaused = false
}

func (c *BulkOperationsContext) CompleteExecution() {
	c.IsExecuting = false
	c.IsPaused = false
}

func (c *BulkOperationsContext) AddResult(itemID, status string) {
	c.Results[itemID] = status
	c.ProcessedCount++
	if status == "success" {
		c.SuccessCount++
	} else if status == "error" {
		c.FailureCount++
	} else if status == "skipped" {
		c.SkippedCount++
	}
}

func (c *BulkOperationsContext) AddError(message string) {
	c.Errors = append(c.Errors, message)
}

func (c *BulkOperationsContext) GetProgress() float64 {
	if c.AffectedItemCount == 0 {
		return 0
	}
	return float64(c.ProcessedCount) / float64(c.AffectedItemCount)
}

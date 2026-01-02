package intents

import (
	"testing"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/stretchr/testify/assert"
)

func TestCaptureEventIntent_NewCaptureEventIntent(t *testing.T) {
	tests := []struct {
		name    string
		ctx     *CaptureEventContext
		wantErr bool
	}{
		{
			name: "creates new intent with valid context",
			ctx: &CaptureEventContext{
				CaptureStrategy: "manual",
				PreviousEvent:   nil,
				Metadata:        make(map[string]string),
			},
			wantErr: false,
		},
		{
			name: "fails with empty strategy",
			ctx: &CaptureEventContext{
				CaptureStrategy: "",
				PreviousEvent:   nil,
				Metadata:        make(map[string]string),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			intent, err := NewCaptureEventIntent(tt.ctx)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, intent)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, intent)
				assert.True(t, intent.active)
				assert.Equal(t, CaptureStateChooseStrategy, intent.state.currentState)
				assert.Nil(t, intent.result)
			}
		})
	}
}

func TestCaptureEventIntent_Init(t *testing.T) {
	ctx := &CaptureEventContext{
		CaptureStrategy: "manual",
		PreviousEvent:   nil,
		Metadata:        make(map[string]string),
	}
	intent, err := NewCaptureEventIntent(ctx)
	assert.NoError(t, err)

	cmd := intent.Init()
	// Init returns nil for now, which is valid
	assert.Nil(t, cmd)
}

func TestCaptureEventIntent_View(t *testing.T) {
	ctx := &CaptureEventContext{
		CaptureStrategy: "manual",
		PreviousEvent:   nil,
		Metadata:        make(map[string]string),
	}
	intent, err := NewCaptureEventIntent(ctx)
	assert.NoError(t, err)

	tests := []struct {
		name     string
		state    string
		wantView string
	}{
		{
			name:     "renders view for ChooseStrategy state",
			state:    CaptureStateChooseStrategy,
			wantView: "Choose Capture Strategy",
		},
		{
			name:     "renders view for Form state",
			state:    CaptureStateForm,
			wantView: "Capture Event Form",
		},
		{
			name:     "renders view for Review state",
			state:    CaptureStateReview,
			wantView: "Review Inferred Event",
		},
		{
			name:     "renders view for Submit state",
			state:    CaptureStateSubmit,
			wantView: "Submitting event",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			intent.state.currentState = tt.state
			view := intent.View()
			assert.NotEmpty(t, view)
			assert.Contains(t, view, tt.wantView)
		})
	}
}

func TestCaptureEventIntent_View_Inactive(t *testing.T) {
	ctx := &CaptureEventContext{
		CaptureStrategy: "manual",
		PreviousEvent:   nil,
		Metadata:        make(map[string]string),
	}
	intent, err := NewCaptureEventIntent(ctx)
	assert.NoError(t, err)

	intent.active = false
	view := intent.View()
	assert.Contains(t, view, "not active")
}

func TestCaptureEventIntent_Result(t *testing.T) {
	ctx := &CaptureEventContext{
		CaptureStrategy: "manual",
		PreviousEvent:   nil,
		Metadata:        make(map[string]string),
	}
	intent, err := NewCaptureEventIntent(ctx)
	assert.NoError(t, err)

	// Should return nil while active
	result := intent.Result()
	assert.Nil(t, result)

	// Should return completed result after completion
	eventResult := &CaptureEventResult{
		Event:          nil,
		Bursts:         make([]*career.Burst, 0),
		Facts:          make([]*career.Fact, 0),
		AcceptedFields: make(map[string]bool),
		RejectedFields: make(map[string]string),
	}
	intent.setCompleted(eventResult)
	result = intent.Result()
	assert.NotNil(t, result)
	assert.Equal(t, Completed, result.Status)
}

func TestCaptureEventIntent_SetCompleted(t *testing.T) {
	ctx := &CaptureEventContext{
		CaptureStrategy: "manual",
		PreviousEvent:   nil,
		Metadata:        make(map[string]string),
	}
	intent, err := NewCaptureEventIntent(ctx)
	assert.NoError(t, err)

	eventResult := &CaptureEventResult{
		AcceptedFields: make(map[string]bool),
		RejectedFields: make(map[string]string),
	}
	intent.setCompleted(eventResult)

	assert.Equal(t, Completed, intent.result.Status)
	assert.False(t, intent.active)
}

func TestCaptureEventIntent_SetCancelled(t *testing.T) {
	ctx := &CaptureEventContext{
		CaptureStrategy: "manual",
		PreviousEvent:   nil,
		Metadata:        make(map[string]string),
	}
	intent, err := NewCaptureEventIntent(ctx)
	assert.NoError(t, err)

	intent.setCancelled()

	assert.Equal(t, Cancelled, intent.result.Status)
	assert.False(t, intent.active)
}

func TestCaptureEventIntent_SetFailed(t *testing.T) {
	ctx := &CaptureEventContext{
		CaptureStrategy: "manual",
		PreviousEvent:   nil,
		Metadata:        make(map[string]string),
	}
	intent, err := NewCaptureEventIntent(ctx)
	assert.NoError(t, err)

	intent.setFailed("error_code", "Error message", nil)

	assert.Equal(t, Failed, intent.result.Status)
	assert.NotNil(t, intent.result.Error)
	assert.Equal(t, "error_code", intent.result.Error.Code)
	assert.Equal(t, "Error message", intent.result.Error.Message)
	assert.False(t, intent.active)
}

func TestCaptureEventIntent_SetPartial(t *testing.T) {
	ctx := &CaptureEventContext{
		CaptureStrategy: "manual",
		PreviousEvent:   nil,
		Metadata:        make(map[string]string),
	}
	intent, err := NewCaptureEventIntent(ctx)
	assert.NoError(t, err)

	eventResult := &CaptureEventResult{
		AcceptedFields: make(map[string]bool),
		RejectedFields: make(map[string]string),
	}
	intent.setPartial(eventResult, "partial_code", "Partial message")

	assert.Equal(t, Partial, intent.result.Status)
	assert.NotNil(t, intent.result.Error)
	assert.Equal(t, "partial_code", intent.result.Error.Code)
	assert.Equal(t, "Partial message", intent.result.Error.Message)
	assert.False(t, intent.active)
}

func TestCaptureEventIntent_Update_Inactive(t *testing.T) {
	ctx := &CaptureEventContext{
		CaptureStrategy: "manual",
		PreviousEvent:   nil,
		Metadata:        make(map[string]string),
	}
	intent, err := NewCaptureEventIntent(ctx)
	assert.NoError(t, err)

	intent.active = false
	cmd := intent.Update(nil)
	assert.Nil(t, cmd)
}

func TestCaptureEventIntent_Update_ChooseStrategy(t *testing.T) {
	ctx := &CaptureEventContext{
		CaptureStrategy: "manual",
		PreviousEvent:   nil,
		Metadata:        make(map[string]string),
	}
	intent, err := NewCaptureEventIntent(ctx)
	assert.NoError(t, err)

	intent.state.currentState = CaptureStateChooseStrategy
	cmd := intent.Update(nil)
	// Command will be nil for now since updateChooseStrategy returns nil
	assert.Nil(t, cmd)
}


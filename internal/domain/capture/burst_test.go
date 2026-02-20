package capture_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/baphled/kariya/internal/domain/capture"
	"github.com/baphled/kariya/internal/testutil/fixtures"
)

func TestApplyBurstEdit_AppliesChanges(t *testing.T) {
	burst := fixtures.Burst("burst-1")
	burst.Name = "Old Name"
	burst.Description = "Old description"

	input := capture.BurstEditInput{Name: "New Name", Description: "New description"}

	result, err := capture.ApplyBurstEdit(burst, input)
	require.NoError(t, err)

	assert.Equal(t, "New Name", result.Name)
	assert.Equal(t, "New description", result.Description)
	assert.False(t, result.UpdatedAt.IsZero(), "UpdatedAt should not be zero")
}

func TestApplyBurstEdit_RejectsEmptyName(t *testing.T) {
	burst := fixtures.Burst("burst-1")

	_, err := capture.ApplyBurstEdit(burst, capture.BurstEditInput{Name: "", Description: "desc"})
	require.Error(t, err)
	assert.Equal(t, "burst name cannot be empty", err.Error())
}

func TestApplyBurstEdit_RejectsWhitespaceName(t *testing.T) {
	burst := fixtures.Burst("burst-1")

	_, err := capture.ApplyBurstEdit(burst, capture.BurstEditInput{Name: "   ", Description: "desc"})
	require.Error(t, err)
}

func TestApplyBurstEdit_PreservesOtherFields(t *testing.T) {
	burst := fixtures.BurstConfirmed("burst-1", "e1", "e2")

	result, err := capture.ApplyBurstEdit(burst, capture.BurstEditInput{Name: "New", Description: "Updated"})
	require.NoError(t, err)

	assert.Equal(t, "burst-1", result.ID)
	assert.Len(t, result.EventIDs, 2)
	assert.True(t, result.Confirmed)
}

func TestCreateBurstFromSuggestion_AllFields(t *testing.T) {
	input := capture.BurstSuggestionInput{
		Name:        "API Development Sprint",
		Description: "Related events about API work",
		EventIDs:    []string{"e1", "e2", "e3"},
	}

	burst := capture.CreateBurstFromSuggestion(input)

	assert.Equal(t, input.Name, burst.Name)
	assert.Equal(t, input.Description, burst.Description)
	assert.Len(t, burst.EventIDs, 3)
	assert.True(t, burst.Confirmed)
	assert.Empty(t, burst.ID, "ID should be empty (assigned by repo)")
}

func TestCreateBurstFromSuggestion_DefaultsName(t *testing.T) {
	input := capture.BurstSuggestionInput{Name: "", EventIDs: []string{"e1"}}

	burst := capture.CreateBurstFromSuggestion(input)

	assert.Equal(t, "Untitled burst", burst.Name)
}

func TestConfirmBurstTimestamps_SetsFlags(t *testing.T) {
	burst := fixtures.Burst("burst-1")
	burst.Confirmed = false

	before := time.Now()
	result := capture.ConfirmBurstTimestamps(burst)
	after := time.Now()

	assert.True(t, result.Confirmed)
	require.NotNil(t, result.ConfirmedAt)
	assert.False(t, result.ConfirmedAt.Before(before) || result.ConfirmedAt.After(after))
	assert.False(t, result.UpdatedAt.Before(before) || result.UpdatedAt.After(after))
}

func TestConfirmBurstTimestamps_ModifiesInPlace(t *testing.T) {
	burst := fixtures.Burst("burst-1")
	burst.Confirmed = false

	result := capture.ConfirmBurstTimestamps(burst)

	assert.Same(t, burst, result)
}

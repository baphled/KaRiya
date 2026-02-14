package capture_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/baphled/kariya/internal/domain/capture"
)

func TestNewEventFromInput_AllFields(t *testing.T) {
	input := capture.EventInput{
		Text:       "Led API redesign for performance improvement",
		Date:       "2024-06-15",
		Company:    "TechCorp",
		Project:    "API Modernization",
		Tags:       []string{"technical", "achievement"},
		Categories: []string{"leadership"},
		Skills:     []string{"go", "api-design"},
	}

	event, err := capture.NewEventFromInput(input)
	require.NoError(t, err)

	assert.Equal(t, input.Text, event.Text)
	assert.Equal(t, input.Company, event.Company)
	assert.Equal(t, input.Project, event.Project)
	assert.Len(t, event.Tags, 2)
	assert.Len(t, event.Categories, 1)
	assert.Len(t, event.Skills, 2)

	expectedDate := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
	assert.True(t, event.Date.Equal(expectedDate), "Date = %v, want %v", event.Date, expectedDate)
	assert.False(t, event.CreatedAt.IsZero(), "CreatedAt should not be zero")
	assert.False(t, event.UpdatedAt.IsZero(), "UpdatedAt should not be zero")
}

func TestNewEventFromInput_DefaultsDateToNow(t *testing.T) {
	input := capture.EventInput{Text: "Some event text for testing purposes"}

	before := time.Now()
	event, err := capture.NewEventFromInput(input)
	after := time.Now()

	require.NoError(t, err)
	assert.False(t, event.Date.Before(before) || event.Date.After(after),
		"Date = %v, want between %v and %v", event.Date, before, after)
}

func TestNewEventFromInput_ParsesRelativeDates(t *testing.T) {
	t.Run("today", func(t *testing.T) {
		input := capture.EventInput{Text: "Some event text for testing purposes", Date: "today"}

		before := time.Now()
		event, err := capture.NewEventFromInput(input)
		after := time.Now()

		require.NoError(t, err)
		assert.False(t, event.Date.Before(before) || event.Date.After(after))
	})

	t.Run("yesterday", func(t *testing.T) {
		input := capture.EventInput{Text: "Some event text for testing purposes", Date: "yesterday"}

		event, err := capture.NewEventFromInput(input)
		require.NoError(t, err)

		yesterday := time.Now().AddDate(0, 0, -1)
		assert.Equal(t, yesterday.Day(), event.Date.Day())
	})
}

func TestNewEventFromInput_InvalidDate(t *testing.T) {
	input := capture.EventInput{Text: "Some event text for testing purposes", Date: "not-a-date"}

	_, err := capture.NewEventFromInput(input)
	require.Error(t, err)

	var dateErr *capture.DateParseError
	assert.ErrorAs(t, err, &dateErr)
}

func TestNewEventFromInput_NilTagsAndCategories(t *testing.T) {
	input := capture.EventInput{Text: "Some event text for testing purposes", Date: "2024-01-01"}

	event, err := capture.NewEventFromInput(input)
	require.NoError(t, err)

	assert.Nil(t, event.Tags)
	assert.Nil(t, event.Categories)
}

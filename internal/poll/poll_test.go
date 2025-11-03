package poll

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewDatePoll(t *testing.T) {
	parameters := []struct {
		name            string
		year            int
		month           time.Month
		weekdays        []time.Weekday
		additionalDays  []int
		excludedDays    []int
		expectedAnswers []time.Time
		expectedExpiry  time.Time
	}{
		{
			name:           "poll in the past",
			year:           2025,
			month:          time.May,
			weekdays:       []time.Weekday{time.Friday, time.Saturday},
			additionalDays: []int{},
			excludedDays:   []int{},
			expectedAnswers: []time.Time{
				time.Date(2025, 5, 2, 20, 0, 0, 0, time.UTC),
				time.Date(2025, 5, 3, 20, 0, 0, 0, time.UTC),
				time.Date(2025, 5, 9, 20, 0, 0, 0, time.UTC),
				time.Date(2025, 5, 10, 20, 0, 0, 0, time.UTC),
				time.Date(2025, 5, 16, 20, 0, 0, 0, time.UTC),
				time.Date(2025, 5, 17, 20, 0, 0, 0, time.UTC),
				time.Date(2025, 5, 23, 20, 0, 0, 0, time.UTC),
				time.Date(2025, 5, 24, 20, 0, 0, 0, time.UTC),
				time.Date(2025, 5, 30, 20, 0, 0, 0, time.UTC),
				time.Date(2025, 5, 31, 20, 0, 0, 0, time.UTC),
			},
			expectedExpiry: time.Date(2025, 4, 30, 12, 0, 0, 0, time.UTC),
		},
		{
			name:           "poll before new year",
			year:           2025,
			month:          time.December,
			weekdays:       []time.Weekday{time.Friday},
			additionalDays: []int{},
			excludedDays:   []int{},
			expectedAnswers: []time.Time{
				time.Date(2025, 12, 5, 20, 0, 0, 0, time.UTC),
				time.Date(2025, 12, 12, 20, 0, 0, 0, time.UTC),
				time.Date(2025, 12, 19, 20, 0, 0, 0, time.UTC),
				time.Date(2025, 12, 26, 20, 0, 0, 0, time.UTC),
			},
			expectedExpiry: time.Date(2025, 11, 30, 12, 0, 0, 0, time.UTC),
		},
		{
			name:           "poll after new year",
			year:           2026,
			month:          time.January,
			weekdays:       []time.Weekday{time.Saturday},
			additionalDays: []int{},
			excludedDays:   []int{},
			expectedAnswers: []time.Time{
				time.Date(2026, 1, 3, 20, 0, 0, 0, time.UTC),
				time.Date(2026, 1, 10, 20, 0, 0, 0, time.UTC),
				time.Date(2026, 1, 17, 20, 0, 0, 0, time.UTC),
				time.Date(2026, 1, 24, 20, 0, 0, 0, time.UTC),
				time.Date(2026, 1, 31, 20, 0, 0, 0, time.UTC),
			},
			expectedExpiry: time.Date(2025, 12, 31, 12, 0, 0, 0, time.UTC),
		},
		{
			name:           "poll with expiration on leap day",
			year:           2028,
			month:          time.March,
			weekdays:       []time.Weekday{time.Wednesday},
			additionalDays: []int{},
			excludedDays:   []int{},
			expectedAnswers: []time.Time{
				time.Date(2028, 3, 1, 20, 0, 0, 0, time.UTC),
				time.Date(2028, 3, 8, 20, 0, 0, 0, time.UTC),
				time.Date(2028, 3, 15, 20, 0, 0, 0, time.UTC),
				time.Date(2028, 3, 22, 20, 0, 0, 0, time.UTC),
				time.Date(2028, 3, 29, 20, 0, 0, 0, time.UTC),
			},
			expectedExpiry: time.Date(2028, 2, 29, 12, 0, 0, 0, time.UTC),
		},
		{
			name:           "poll with additional and excluded days",
			year:           2025,
			month:          time.December,
			weekdays:       []time.Weekday{time.Friday, time.Saturday},
			additionalDays: []int{26, 27, 28, 29, 30},
			excludedDays:   []int{23, 24, 25, 31},
			expectedAnswers: []time.Time{
				time.Date(2025, 12, 5, 20, 0, 0, 0, time.UTC),
				time.Date(2025, 12, 6, 20, 0, 0, 0, time.UTC),
				time.Date(2025, 12, 12, 20, 0, 0, 0, time.UTC),
				time.Date(2025, 12, 13, 20, 0, 0, 0, time.UTC),
				time.Date(2025, 12, 19, 20, 0, 0, 0, time.UTC),
				time.Date(2025, 12, 20, 20, 0, 0, 0, time.UTC),
				time.Date(2025, 12, 26, 20, 0, 0, 0, time.UTC),
				time.Date(2025, 12, 27, 20, 0, 0, 0, time.UTC),
				time.Date(2025, 12, 28, 20, 0, 0, 0, time.UTC),
				time.Date(2025, 12, 29, 20, 0, 0, 0, time.UTC),
			},
			expectedExpiry: time.Date(2025, 11, 30, 12, 0, 0, 0, time.UTC),
		},
	}

	for _, parameter := range parameters {
		t.Run(parameter.name, func(t *testing.T) {
			poll := NewDatePoll("TestQuestion", parameter.year, parameter.month, parameter.weekdays, time.UTC, parameter.additionalDays, parameter.excludedDays)

			assert.NotNil(t, poll)
			assert.Equal(t, "TestQuestion", poll.Question)
			assert.Equal(t, parameter.expectedAnswers, poll.Answers)
			assert.Equal(t, parameter.expectedExpiry, poll.Expiry)
		})
	}
}

func TestNewDatePollResult(t *testing.T) {
	parameters := []struct {
		name           string
		pollID         string
		winningAnswers []time.Time
		finalized      bool
	}{
		{
			name:           "no winning answers, not finalized",
			pollID:         "poll-123",
			winningAnswers: []time.Time{},
			finalized:      false,
		},
		{
			name:   "multiple winning answers, finalized",
			pollID: "poll-xyz",
			winningAnswers: []time.Time{
				time.Date(2025, 12, 5, 20, 0, 0, 0, time.UTC),
				time.Date(2025, 12, 12, 20, 0, 0, 0, time.UTC),
			},
			finalized: true,
		},
	}

	for _, parameter := range parameters {
		t.Run(parameter.name, func(t *testing.T) {
			result := NewDatePollResult(parameter.pollID, parameter.winningAnswers, parameter.finalized)

			assert.NotNil(t, result)
			assert.Equal(t, parameter.pollID, result.PollID)
			assert.Equal(t, parameter.winningAnswers, result.WinningAnswers)
			assert.Equal(t, parameter.finalized, result.Finalized)
		})
	}
}

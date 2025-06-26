package poll

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewDatePoll(t *testing.T) {
	parameters := []struct {
		name            string
		year            int
		month           time.Month
		weekdays        []time.Weekday
		expectedAnswers []time.Time
		expectedExpiry  time.Time
	}{
		{
			name:     "poll in the past",
			year:     2025,
			month:    time.May,
			weekdays: []time.Weekday{time.Friday, time.Saturday},
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
			name:     "poll before new year",
			year:     2025,
			month:    time.December,
			weekdays: []time.Weekday{time.Friday},
			expectedAnswers: []time.Time{
				time.Date(2025, 12, 5, 20, 0, 0, 0, time.UTC),
				time.Date(2025, 12, 12, 20, 0, 0, 0, time.UTC),
				time.Date(2025, 12, 19, 20, 0, 0, 0, time.UTC),
				time.Date(2025, 12, 26, 20, 0, 0, 0, time.UTC),
			},
			expectedExpiry: time.Date(2025, 11, 30, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "poll after new year",
			year:     2026,
			month:    time.January,
			weekdays: []time.Weekday{time.Saturday},
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
			name:     "poll with expiration on leap day",
			year:     2028,
			month:    time.March,
			weekdays: []time.Weekday{time.Wednesday},
			expectedAnswers: []time.Time{
				time.Date(2028, 3, 1, 20, 0, 0, 0, time.UTC),
				time.Date(2028, 3, 8, 20, 0, 0, 0, time.UTC),
				time.Date(2028, 3, 15, 20, 0, 0, 0, time.UTC),
				time.Date(2028, 3, 22, 20, 0, 0, 0, time.UTC),
				time.Date(2028, 3, 29, 20, 0, 0, 0, time.UTC),
			},
			expectedExpiry: time.Date(2028, 2, 29, 12, 0, 0, 0, time.UTC),
		},
	}

	for _, parameter := range parameters {
		t.Run(parameter.name, func(t *testing.T) {
			poll := NewDatePoll("TestQuestion", parameter.year, parameter.month, parameter.weekdays, time.UTC)

			assert.NotNil(t, poll)
			assert.Equal(t, "TestQuestion", poll.Question)
			assert.Equal(t, parameter.expectedAnswers, poll.Answers)
			assert.Equal(t, parameter.expectedExpiry, poll.Expiry)
		})
	}
}

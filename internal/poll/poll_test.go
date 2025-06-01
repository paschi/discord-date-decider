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
				time.Date(2025, 5, 2, 12, 0, 0, 0, time.UTC),
				time.Date(2025, 5, 3, 12, 0, 0, 0, time.UTC),
				time.Date(2025, 5, 9, 12, 0, 0, 0, time.UTC),
				time.Date(2025, 5, 10, 12, 0, 0, 0, time.UTC),
				time.Date(2025, 5, 16, 12, 0, 0, 0, time.UTC),
				time.Date(2025, 5, 17, 12, 0, 0, 0, time.UTC),
				time.Date(2025, 5, 23, 12, 0, 0, 0, time.UTC),
				time.Date(2025, 5, 24, 12, 0, 0, 0, time.UTC),
				time.Date(2025, 5, 30, 12, 0, 0, 0, time.UTC),
				time.Date(2025, 5, 31, 12, 0, 0, 0, time.UTC),
			},
			expectedExpiry: time.Date(2025, 4, 30, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "poll before new year",
			year:     2025,
			month:    time.December,
			weekdays: []time.Weekday{time.Friday},
			expectedAnswers: []time.Time{
				time.Date(2025, 12, 5, 12, 0, 0, 0, time.UTC),
				time.Date(2025, 12, 12, 12, 0, 0, 0, time.UTC),
				time.Date(2025, 12, 19, 12, 0, 0, 0, time.UTC),
				time.Date(2025, 12, 26, 12, 0, 0, 0, time.UTC),
			},
			expectedExpiry: time.Date(2025, 11, 30, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "poll after new year",
			year:     2026,
			month:    time.January,
			weekdays: []time.Weekday{time.Saturday},
			expectedAnswers: []time.Time{
				time.Date(2026, 1, 3, 12, 0, 0, 0, time.UTC),
				time.Date(2026, 1, 10, 12, 0, 0, 0, time.UTC),
				time.Date(2026, 1, 17, 12, 0, 0, 0, time.UTC),
				time.Date(2026, 1, 24, 12, 0, 0, 0, time.UTC),
				time.Date(2026, 1, 31, 12, 0, 0, 0, time.UTC),
			},
			expectedExpiry: time.Date(2025, 12, 31, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "poll with expiration on leap day",
			year:     2028,
			month:    time.March,
			weekdays: []time.Weekday{time.Wednesday},
			expectedAnswers: []time.Time{
				time.Date(2028, 3, 1, 12, 0, 0, 0, time.UTC),
				time.Date(2028, 3, 8, 12, 0, 0, 0, time.UTC),
				time.Date(2028, 3, 15, 12, 0, 0, 0, time.UTC),
				time.Date(2028, 3, 22, 12, 0, 0, 0, time.UTC),
				time.Date(2028, 3, 29, 12, 0, 0, 0, time.UTC),
			},
			expectedExpiry: time.Date(2028, 2, 29, 12, 0, 0, 0, time.UTC),
		},
	}

	for _, parameter := range parameters {
		t.Run(parameter.name, func(t *testing.T) {
			poll := NewDatePoll("TestQuestion", parameter.year, parameter.month, parameter.weekdays)

			assert.NotNil(t, poll)
			assert.Equal(t, "TestQuestion", poll.Question)
			assert.Equal(t, parameter.expectedAnswers, poll.Answers)
			assert.Equal(t, parameter.expectedExpiry, poll.Expiry)
		})
	}
}

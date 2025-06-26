package poll

import (
	"time"
)

type DatePoll struct {
	Question string
	Answers  []time.Time
	Expiry   time.Time
}

type DatePollResult struct {
	PollID         string
	WinningAnswers []time.Time
	Finalized      bool
}

func NewDatePoll(question string, year int, month time.Month, weekdays []time.Weekday, location *time.Location) *DatePoll {
	return &DatePoll{
		Question: question,
		Expiry:   time.Date(year, month, 0, 12, 0, 0, 0, location),
		Answers:  getDates(year, month, weekdays, location),
	}
}

func NewDatePollResult(pollID string, winningAnswers []time.Time, finalized bool) *DatePollResult {
	return &DatePollResult{
		PollID:         pollID,
		WinningAnswers: winningAnswers,
		Finalized:      finalized,
	}
}

func getDates(year int, month time.Month, weekdays []time.Weekday, location *time.Location) []time.Time {
	var dates []time.Time
	firstDay := time.Date(year, month, 1, 20, 0, 0, 0, location)
	lastDay := time.Date(year, month+1, 0, 20, 0, 0, 0, location)
	for day := firstDay; day.Before(lastDay.AddDate(0, 0, 1)); day = day.AddDate(0, 0, 1) {
		for _, weekday := range weekdays {
			if day.Weekday() == weekday {
				dates = append(dates, day)
				break
			}
		}
	}
	return dates
}

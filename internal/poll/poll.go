package poll

import (
	"time"
)

const maxAnswers = 10

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

func NewDatePoll(question string, year int, month time.Month, weekdays []time.Weekday, location *time.Location, additionalDays []int, excludedDays []int) *DatePoll {
	return &DatePoll{
		Question: question,
		Expiry:   time.Date(year, month, 0, 12, 0, 0, 0, location),
		Answers:  getDates(year, month, weekdays, location, additionalDays, excludedDays),
	}
}

func NewDatePollResult(pollID string, winningAnswers []time.Time, finalized bool) *DatePollResult {
	return &DatePollResult{
		PollID:         pollID,
		WinningAnswers: winningAnswers,
		Finalized:      finalized,
	}
}

func getDates(year int, month time.Month, weekdays []time.Weekday, location *time.Location, additionalDays []int, excludedDays []int) []time.Time {
	var dates []time.Time
	count := 0
	firstDay := time.Date(year, month, 1, 20, 0, 0, 0, location)
	lastDay := time.Date(year, month+1, 0, 20, 0, 0, 0, location)
	for day := firstDay; day.Before(lastDay.AddDate(0, 0, 1)); day = day.AddDate(0, 0, 1) {
		if contains(excludedDays, day.Day()) {
			continue
		}
		if contains(weekdays, day.Weekday()) || contains(additionalDays, day.Day()) {
			if count < maxAnswers {
				dates = append(dates, day)
				count++
			}
		}
	}
	return dates
}

func contains[T comparable](s []T, e T) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

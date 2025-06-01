package poll

import (
	"time"
)

type DatePoll struct {
	Question string
	Answers  []time.Time
	Expiry   time.Time
}

func NewDatePoll(question string, year int, month time.Month, weekdays []time.Weekday) *DatePoll {
	return &DatePoll{
		Question: question,
		Expiry:   time.Date(year, month, 0, 12, 0, 0, 0, time.UTC),
		Answers:  getDates(year, month, weekdays),
	}
}

func getDates(year int, month time.Month, weekdays []time.Weekday) []time.Time {
	var dates []time.Time
	firstDay := time.Date(year, month, 1, 12, 0, 0, 0, time.UTC)
	lastDay := time.Date(year, month+1, 0, 12, 0, 0, 0, time.UTC)
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

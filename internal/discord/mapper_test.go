package discord

import (
	"fmt"
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/paschi/discord-date-decider/internal/message"
	"github.com/paschi/discord-date-decider/internal/poll"
	"github.com/stretchr/testify/assert"
)

func TestToDiscordMessage(t *testing.T) {
	parameters := []struct {
		name    string
		message *message.Message
	}{
		{
			name:    "message mentioning everyone",
			message: message.NewMessage("testMessage", true),
		},
		{
			name:    "message not mentioning everyone",
			message: message.NewMessage("anotherTestMessage", false),
		},
	}

	for _, parameter := range parameters {
		t.Run(parameter.name, func(t *testing.T) {
			discordMessage := toDiscordMessage(parameter.message)

			assert.NotNil(t, discordMessage)
			assert.Equal(t, parameter.message.Content, discordMessage.Content)
			if parameter.message.MentionsEveryone {
				assert.NotNil(t, discordMessage.AllowedMentions)
				assert.Contains(t, discordMessage.AllowedMentions.Parse, discordgo.AllowedMentionTypeEveryone)
			} else {
				assert.Nil(t, discordMessage.AllowedMentions)
			}
		})
	}
}

func TestToDiscordPollMessage(t *testing.T) {
	futureDate := time.Now().AddDate(0, 1, 0)
	pastDate := time.Now().AddDate(0, -1, 0)

	parameters := []struct {
		name        string
		poll        *poll.DatePoll
		expectError bool
	}{
		{
			name:        "valid poll",
			poll:        poll.NewDatePoll("Test Poll Question", futureDate.Year(), futureDate.Month(), []time.Weekday{time.Friday, time.Saturday}, time.UTC),
			expectError: false,
		},
		{
			name:        "expired poll",
			poll:        poll.NewDatePoll("Expired Poll Question", pastDate.Year(), pastDate.Month(), []time.Weekday{time.Friday, time.Saturday}, time.UTC),
			expectError: true,
		},
	}

	for _, parameter := range parameters {
		t.Run(parameter.name, func(t *testing.T) {
			discordPoll, err := toDiscordPollMessage(parameter.poll)

			if parameter.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "poll is already expired")
				assert.Nil(t, discordPoll)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, discordPoll)
				assert.Equal(t, parameter.poll.Question, discordPoll.Poll.Question.Text)
				assert.Equal(t, len(parameter.poll.Answers), len(discordPoll.Poll.Answers))
				assert.True(t, discordPoll.Poll.AllowMultiselect)
				for i, answer := range parameter.poll.Answers {
					expectedText := fmt.Sprintf("%s, %02d.%02d.%d", answer.Weekday().String(), answer.Day(), answer.Month(), answer.Year())
					assert.Equal(t, expectedText, discordPoll.Poll.Answers[i].Media.Text)
				}
			}
		})
	}
}

func TestToDatePollResult(t *testing.T) {
	loc := time.UTC

	makeAnswer := func(id int, d time.Time) discordgo.PollAnswer {
		text := fmt.Sprintf("%s, %02d.%02d.%d", d.Weekday().String(), d.Day(), d.Month(), d.Year())
		return discordgo.PollAnswer{AnswerID: id, Media: &discordgo.PollMedia{Text: text}}
	}

	t.Run("single winner", func(t *testing.T) {
		date1 := time.Date(2025, 8, 15, 0, 0, 0, 0, loc) // Friday
		date2 := time.Date(2025, 8, 16, 0, 0, 0, 0, loc) // Saturday
		msg := &discordgo.Message{
			ID: "123",
			Poll: &discordgo.Poll{
				Answers: []discordgo.PollAnswer{
					makeAnswer(0, date1),
					makeAnswer(1, date2),
				},
				Results: &discordgo.PollResults{
					Finalized:    true,
					AnswerCounts: []*discordgo.PollAnswerCount{{ID: 0, Count: 5}, {ID: 1, Count: 3}},
				},
			},
		}

		result, err := toDatePollResult(msg, loc)
		assert.NoError(t, err)
		if assert.NotNil(t, result) {
			assert.Equal(t, "123", result.PollID)
			assert.True(t, result.Finalized)
			if assert.Len(t, result.WinningAnswers, 1) {
				expected := time.Date(2025, 8, 15, 20, 0, 0, 0, loc)
				assert.True(t, expected.Equal(result.WinningAnswers[0]))
			}
		}
	})

	t.Run("tie between two answers", func(t *testing.T) {
		date1 := time.Date(2025, 12, 5, 0, 0, 0, 0, loc)  // Friday
		date2 := time.Date(2025, 12, 12, 0, 0, 0, 0, loc) // Friday
		msg := &discordgo.Message{
			ID: "456",
			Poll: &discordgo.Poll{
				Answers: []discordgo.PollAnswer{
					makeAnswer(10, date1),
					makeAnswer(11, date2),
				},
				Results: &discordgo.PollResults{
					Finalized:    false,
					AnswerCounts: []*discordgo.PollAnswerCount{{ID: 10, Count: 7}, {ID: 11, Count: 7}},
				},
			},
		}

		result, err := toDatePollResult(msg, loc)
		assert.NoError(t, err)
		if assert.NotNil(t, result) {
			assert.Equal(t, "456", result.PollID)
			assert.False(t, result.Finalized)
			if assert.Len(t, result.WinningAnswers, 2) {
				expected1 := time.Date(2025, 12, 5, 20, 0, 0, 0, loc)
				expected2 := time.Date(2025, 12, 12, 20, 0, 0, 0, loc)
				// The function preserves the order of Answers slice
				assert.True(t, expected1.Equal(result.WinningAnswers[0]))
				assert.True(t, expected2.Equal(result.WinningAnswers[1]))
			}
		}
	})

	t.Run("no poll in message returns error", func(t *testing.T) {
		msg := &discordgo.Message{ID: "789"}
		result, err := toDatePollResult(msg, loc)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("no answer counts returns error", func(t *testing.T) {
		date1 := time.Date(2025, 1, 3, 0, 0, 0, 0, loc)
		msg := &discordgo.Message{
			ID: "999",
			Poll: &discordgo.Poll{
				Answers: []discordgo.PollAnswer{makeAnswer(0, date1)},
				Results: &discordgo.PollResults{Finalized: true, AnswerCounts: []*discordgo.PollAnswerCount{}},
			},
		}
		result, err := toDatePollResult(msg, loc)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

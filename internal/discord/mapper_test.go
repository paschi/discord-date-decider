package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/paschi/discord-date-decider/internal/message"
	"github.com/paschi/discord-date-decider/internal/poll"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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
			poll:        poll.NewDatePoll("Test Poll Question", futureDate.Year(), futureDate.Month(), []time.Weekday{time.Friday, time.Saturday}),
			expectError: false,
		},
		{
			name:        "expired poll",
			poll:        poll.NewDatePoll("Expired Poll Question", pastDate.Year(), pastDate.Month(), []time.Weekday{time.Friday, time.Saturday}),
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
					expectedText := fmt.Sprintf("%s, %02d.%02d.", answer.Weekday().String(), answer.Day(), answer.Month())
					assert.Equal(t, expectedText, discordPoll.Poll.Answers[i].Media.Text)
				}
			}
		})
	}
}

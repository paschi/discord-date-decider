package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/paschi/discord-date-decider/internal/message"
	"github.com/paschi/discord-date-decider/internal/poll"
	"math"
	"time"
)

func toDiscordMessage(message *message.Message) *discordgo.MessageSend {
	var allowedMentions *discordgo.MessageAllowedMentions
	if message.MentionsEveryone {
		allowedMentions = &discordgo.MessageAllowedMentions{
			Parse: []discordgo.AllowedMentionType{discordgo.AllowedMentionTypeEveryone},
		}
	}
	return &discordgo.MessageSend{
		Content:         message.Content,
		AllowedMentions: allowedMentions,
	}
}

func toDiscordPollMessage(poll *poll.DatePoll) (*discordgo.MessageSend, error) {
	if time.Now().After(poll.Expiry) {
		return nil, fmt.Errorf("poll is already expired")
	}
	hoursUntilExpiry := int(math.Floor(poll.Expiry.Sub(time.Now()).Hours()))
	return &discordgo.MessageSend{
		Poll: &discordgo.Poll{
			Question:         discordgo.PollMedia{Text: poll.Question},
			Answers:          toDiscordAnswers(poll.Answers),
			AllowMultiselect: true,
			Duration:         hoursUntilExpiry,
		},
	}, nil
}

func toDiscordAnswers(answers []time.Time) []discordgo.PollAnswer {
	var discordAnswers []discordgo.PollAnswer
	for _, answer := range answers {
		discordAnswers = append(discordAnswers, discordgo.PollAnswer{Media: &discordgo.PollMedia{
			Text: fmt.Sprintf("%s, %02d.%02d.", answer.Weekday().String(), answer.Day(), answer.Month()),
		}})
	}
	return discordAnswers
}

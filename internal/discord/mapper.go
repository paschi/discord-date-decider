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
			Text: fmt.Sprintf("%s, %02d.%02d.%d", answer.Weekday().String(), answer.Day(), answer.Month(), answer.Year()),
		}})
	}
	return discordAnswers
}

func toDatePollResult(discordMessage *discordgo.Message) (*poll.DatePollResult, error) {
	discordPoll := discordMessage.Poll
	if discordPoll == nil {
		return nil, fmt.Errorf("no poll message")
	}
	var highestCount int
	var winningAnswerID int
	for i, answerCount := range discordPoll.Results.AnswerCounts {
		if i == 0 || answerCount.Count > highestCount {
			highestCount = answerCount.Count
			winningAnswerID = answerCount.ID
		}
	}
	for _, answer := range discordPoll.Answers {
		if answer.AnswerID == winningAnswerID {
			dateStr := answer.Media.Text
			var weekday string
			var day, month, year int
			_, err := fmt.Sscanf(dateStr, "%s %02d.%02d.%d", &weekday, &day, &month, &year)
			if err != nil {
				return nil, err
			}
			winningDate := time.Date(year, time.Month(month), day, 12, 0, 0, 0, time.UTC)
			return poll.NewDatePollResult(discordMessage.ID, winningDate, discordPoll.Results.Finalized), nil
		}
	}
	return nil, fmt.Errorf("could not find a winning answer for this poll: %+v, %+v, %+v, %+v", discordPoll, discordPoll.Answers, discordPoll.Results, discordPoll.Results.AnswerCounts)
}

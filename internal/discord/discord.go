package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/paschi/discord-date-decider/internal/poll"
	"log"
	"math"
	"time"
)

type DiscordService struct {
	session *discordgo.Session
}

func NewDiscordService(token string) (*DiscordService, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	err = session.Open()
	if err != nil {
		return nil, err
	}
	return &DiscordService{session: session}, nil
}

func (d *DiscordService) Close() {
	err := d.session.Close()
	if err != nil {
		log.Printf("error closing session: %v", err)
	}
}

func (d *DiscordService) SendPoll(channel string, poll *poll.DatePoll) (*discordgo.Message, error) {
	discordPoll, err := toDiscordPoll(poll)
	if err != nil {
		log.Printf("error converting poll to discord poll: %v", err)
		return nil, err
	}
	log.Printf("sending converted poll to channel %s", channel)
	message, err := d.session.ChannelMessageSendComplex(channel, &discordgo.MessageSend{
		Poll: discordPoll,
	})
	if err != nil {
		log.Printf("error sending poll: %v", err)
		return nil, err
	}
	return message, nil
}

func (d *DiscordService) SendMessage(channel string, content string) (*discordgo.Message, error) {
	message, err := d.session.ChannelMessageSendComplex(channel, &discordgo.MessageSend{
		Content: content,
		AllowedMentions: &discordgo.MessageAllowedMentions{
			Parse: []discordgo.AllowedMentionType{discordgo.AllowedMentionTypeEveryone},
		},
	})
	if err != nil {
		log.Printf("error sending message: %v", err)
		return nil, err
	}
	return message, nil
}

func toDiscordPoll(p *poll.DatePoll) (*discordgo.Poll, error) {
	if time.Now().After(p.Expiry) {
		return nil, fmt.Errorf("poll is already expired")
	}
	hoursUntilExpiry := int(math.Round(p.Expiry.Sub(time.Now()).Hours()))
	return &discordgo.Poll{
		Question:         discordgo.PollMedia{Text: p.Question},
		Answers:          toDiscordAnswers(p.Answers),
		AllowMultiselect: true,
		Duration:         hoursUntilExpiry,
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

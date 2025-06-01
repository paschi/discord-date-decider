package discord

import (
	"fmt"
	"github.com/paschi/discord-date-decider/internal/message"
	"github.com/paschi/discord-date-decider/internal/poll"
)

type Service interface {
	Open() error
	Close() error
	SendMessage(channelID string, message *message.Message) (string, error)
	SendPoll(channelID string, poll *poll.DatePoll) (string, error)
	PinPoll(channelID string, pollID string) error
}

type DefaultService struct {
	client Client
}

func NewDefaultService(client Client) *DefaultService {
	return &DefaultService{client}
}

func (d *DefaultService) Open() error {
	err := d.client.Open()
	if err != nil {
		return fmt.Errorf("could not open discord client: %w", err)
	}
	return nil
}

func (d *DefaultService) Close() error {
	err := d.client.Close()
	if err != nil {
		return fmt.Errorf("could not close discord client: %w", err)
	}
	return nil
}

func (d *DefaultService) SendMessage(channelID string, message *message.Message) (string, error) {
	discordMessage, err := d.client.ChannelMessageSend(channelID, toDiscordMessage(message))
	if err != nil {
		return "", fmt.Errorf("client could not send message to channel: %w", err)
	}
	return discordMessage.ID, nil
}

func (d *DefaultService) SendPoll(channelID string, poll *poll.DatePoll) (string, error) {
	discordPoll, err := toDiscordPollMessage(poll)
	if err != nil {
		return "", fmt.Errorf("client could not convert poll to discord poll: %w", err)
	}
	discordMessage, err := d.client.ChannelMessageSend(channelID, discordPoll)
	if err != nil {
		return "", fmt.Errorf("client could not send poll to channel: %w", err)
	}
	return discordMessage.ID, nil
}

func (d *DefaultService) PinPoll(channel string, message string) error {
	err := d.client.ChannelMessagePin(channel, message)
	if err != nil {
		return fmt.Errorf("client could not pin message to channel: %w", err)
	}
	return nil
}

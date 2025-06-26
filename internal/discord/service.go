package discord

import (
	"fmt"
	"time"

	"github.com/paschi/discord-date-decider/internal/message"
	"github.com/paschi/discord-date-decider/internal/poll"
)

type Service interface {
	Open() error
	Close() error
	SendMessage(channelID string, message *message.Message) (string, error)
	SendPoll(channelID string, poll *poll.DatePoll) (string, error)
	PinPoll(channelID string, pollID string) error
	UnpinPoll(channelID string, pollID string) error
	GetLastPinnedPollResult(channelID string, location *time.Location) (*poll.DatePollResult, error)
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
		return "", fmt.Errorf("could not send message to channel: %w", err)
	}
	return discordMessage.ID, nil
}

func (d *DefaultService) SendPoll(channelID string, poll *poll.DatePoll) (string, error) {
	discordPoll, err := toDiscordPollMessage(poll)
	if err != nil {
		return "", fmt.Errorf("could not convert poll to discord poll: %w", err)
	}
	discordMessage, err := d.client.ChannelMessageSend(channelID, discordPoll)
	if err != nil {
		return "", fmt.Errorf("could not send poll to channel: %w", err)
	}
	return discordMessage.ID, nil
}

func (d *DefaultService) PinPoll(channelID string, pollID string) error {
	err := d.client.ChannelMessagePin(channelID, pollID)
	if err != nil {
		return fmt.Errorf("could not pin message to channel: %w", err)
	}
	return nil
}

func (d *DefaultService) UnpinPoll(channelID string, pollID string) error {
	err := d.client.ChannelMessageUnpin(channelID, pollID)
	if err != nil {
		return fmt.Errorf("could not unpin message from channel: %w", err)
	}
	return nil
}

func (d *DefaultService) GetLastPinnedPollResult(channelID string, location *time.Location) (*poll.DatePollResult, error) {
	pinnedMessages, err := d.client.ChannelMessagesPinned(channelID)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve pinned messages: %w", err)
	}
	for _, pinnedMessage := range pinnedMessages {
		if pinnedMessage.Poll == nil {
			continue
		}
		result, err := toDatePollResult(pinnedMessage, location)
		if err != nil {
			return nil, fmt.Errorf("could not convert message to date poll result: %w", err)
		}
		return result, nil
	}
	return nil, fmt.Errorf("could not find last pinned poll")
}

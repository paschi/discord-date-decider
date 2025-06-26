package discord

import (
	"github.com/bwmarrin/discordgo"
)

type Client interface {
	Open() error
	Close() error
	ChannelMessageSend(channelID string, data *discordgo.MessageSend) (*discordgo.Message, error)
	ChannelMessagePin(channelID string, messageID string) error
	ChannelMessageUnpin(channelID string, messageID string) error
	ChannelMessagesPinned(channelID string) ([]*discordgo.Message, error)
}

type DefaultClient struct {
	session *discordgo.Session
}

func NewDefaultClient(token string) (*DefaultClient, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	return &DefaultClient{session: session}, nil
}

func (c *DefaultClient) Open() error {
	return c.session.Open()
}

func (c *DefaultClient) Close() error {
	return c.session.Close()
}

func (c *DefaultClient) ChannelMessageSend(channelID string, data *discordgo.MessageSend) (*discordgo.Message, error) {
	return c.session.ChannelMessageSendComplex(channelID, data)
}

func (c *DefaultClient) ChannelMessagePin(channelID string, messageID string) error {
	return c.session.ChannelMessagePin(channelID, messageID)
}

func (c *DefaultClient) ChannelMessageUnpin(channelID string, messageID string) error {
	return c.session.ChannelMessageUnpin(channelID, messageID)
}

func (c *DefaultClient) ChannelMessagesPinned(channelID string) ([]*discordgo.Message, error) {
	return c.session.ChannelMessagesPinned(channelID)
}

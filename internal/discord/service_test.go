package discord

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"github.com/paschi/discord-date-decider/internal/message"
	"github.com/paschi/discord-date-decider/internal/poll"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type MockClient struct {
	mock.Mock
}

func (m *MockClient) Open() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockClient) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockClient) ChannelMessageSend(channelID string, data *discordgo.MessageSend) (*discordgo.Message, error) {
	args := m.Called(channelID, data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*discordgo.Message), args.Error(1)
}

func (m *MockClient) ChannelMessagePin(channelID string, messageID string) error {
	args := m.Called(channelID, messageID)
	return args.Error(0)
}

func TestNewDefaultService(t *testing.T) {
	t.Run("successful initialization", func(t *testing.T) {
		mockClient := new(MockClient)

		service := NewDefaultService(mockClient)

		assert.NotNil(t, service)
		mockClient.AssertExpectations(t)
	})
}

func TestNewDefaultService_Open(t *testing.T) {
	t.Run("successful open", func(t *testing.T) {
		mockClient := new(MockClient)
		mockClient.On("Open").Return(nil)

		service := NewDefaultService(mockClient)
		err := service.Open()

		assert.NoError(t, err)
		mockClient.AssertExpectations(t)
	})

	t.Run("error during open", func(t *testing.T) {
		mockClient := new(MockClient)
		expectedErr := errors.New("open error")
		mockClient.On("Open").Return(expectedErr)

		service := NewDefaultService(mockClient)
		err := service.Open()

		assert.Error(t, err)
		assert.Equal(t, expectedErr, errors.Unwrap(err))
		mockClient.AssertExpectations(t)
	})
}

func TestDefaultService_Close(t *testing.T) {
	t.Run("successful close", func(t *testing.T) {
		mockClient := new(MockClient)
		mockClient.On("Close").Return(nil)

		service := NewDefaultService(mockClient)
		err := service.Close()

		assert.NoError(t, err)
		mockClient.AssertExpectations(t)
	})

	t.Run("error during close", func(t *testing.T) {
		mockClient := new(MockClient)
		expectedErr := errors.New("close error")
		mockClient.On("Close").Return(expectedErr)

		service := NewDefaultService(mockClient)
		err := service.Close()

		assert.Error(t, err)
		assert.Equal(t, expectedErr, errors.Unwrap(err))
		mockClient.AssertExpectations(t)
	})
}

func TestDefaultService_SendMessage(t *testing.T) {
	t.Run("successful message send", func(t *testing.T) {
		mockClient := new(MockClient)
		channelID := "test-channel"
		testMessage := message.NewMessage("test content", false)
		discordMessage := &discordgo.Message{ID: "message-id"}
		mockClient.On("ChannelMessageSend", channelID, mock.AnythingOfType("*discordgo.MessageSend")).
			Return(discordMessage, nil)

		service := NewDefaultService(mockClient)
		messageID, err := service.SendMessage(channelID, testMessage)

		assert.NoError(t, err)
		assert.Equal(t, discordMessage.ID, messageID)
		mockClient.AssertExpectations(t)
	})

	t.Run("error during message send", func(t *testing.T) {
		mockClient := new(MockClient)
		channelID := "test-channel"
		testMessage := message.NewMessage("test content", true)
		expectedErr := errors.New("send error")
		mockClient.On("ChannelMessageSend", channelID, mock.AnythingOfType("*discordgo.MessageSend")).
			Return(nil, expectedErr)

		service := NewDefaultService(mockClient)
		messageID, err := service.SendMessage(channelID, testMessage)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, errors.Unwrap(err))
		assert.Equal(t, "", messageID)
		mockClient.AssertExpectations(t)
	})
}

func TestDefaultService_SendPoll(t *testing.T) {
	t.Run("successful poll send", func(t *testing.T) {
		mockClient := new(MockClient)
		channelID := "test-channel"
		futureDate := time.Now().AddDate(0, 1, 0)
		testPoll := poll.NewDatePoll("Test Poll", futureDate.Year(), futureDate.Month(), []time.Weekday{time.Friday, time.Saturday})
		discordMessage := &discordgo.Message{ID: "poll-id"}
		mockClient.On("ChannelMessageSend", channelID, mock.AnythingOfType("*discordgo.MessageSend")).
			Return(discordMessage, nil)

		service := NewDefaultService(mockClient)
		pollID, err := service.SendPoll(channelID, testPoll)

		assert.NoError(t, err)
		assert.Equal(t, discordMessage.ID, pollID)
		mockClient.AssertExpectations(t)
	})

	t.Run("error during poll send", func(t *testing.T) {
		mockClient := new(MockClient)
		channelID := "test-channel"
		futureDate := time.Now().AddDate(0, 1, 0)
		testPoll := poll.NewDatePoll("Test Poll", futureDate.Year(), futureDate.Month(), []time.Weekday{time.Friday, time.Saturday})
		expectedErr := errors.New("send error")

		mockClient.On("ChannelMessageSend", channelID, mock.AnythingOfType("*discordgo.MessageSend")).
			Return(nil, expectedErr)

		service := NewDefaultService(mockClient)
		pollID, err := service.SendPoll(channelID, testPoll)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, errors.Unwrap(err))
		assert.Equal(t, "", pollID)
		mockClient.AssertExpectations(t)
	})

	t.Run("error during convert poll", func(t *testing.T) {
		mockClient := new(MockClient)
		channelID := "test-channel"
		pastDate := time.Now().AddDate(0, -1, 0)
		testPoll := poll.NewDatePoll("Test Poll", pastDate.Year(), pastDate.Month(), []time.Weekday{time.Friday, time.Saturday})

		service := NewDefaultService(mockClient)
		pollID, err := service.SendPoll(channelID, testPoll)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "poll is already expired")
		assert.Equal(t, "", pollID)
		mockClient.AssertNotCalled(t, "ChannelMessageSend")
	})
}

func TestDefaultService_PinPoll(t *testing.T) {
	t.Run("successful pin poll", func(t *testing.T) {
		mockClient := new(MockClient)
		channelID := "test-channel"
		messageID := "message-id"
		mockClient.On("ChannelMessagePin", channelID, messageID).Return(nil)

		service := NewDefaultService(mockClient)
		err := service.PinPoll(channelID, messageID)

		assert.NoError(t, err)
		mockClient.AssertExpectations(t)
	})

	t.Run("error during pin poll", func(t *testing.T) {
		mockClient := new(MockClient)
		channelID := "test-channel"
		messageID := "message-id"
		expectedErr := errors.New("pin error")
		mockClient.On("ChannelMessagePin", channelID, messageID).Return(expectedErr)

		service := NewDefaultService(mockClient)
		err := service.PinPoll(channelID, messageID)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, errors.Unwrap(err))
		mockClient.AssertExpectations(t)
	})
}

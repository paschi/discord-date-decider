package main

import (
	"errors"
	"github.com/paschi/discord-date-decider/internal/message"
	"github.com/paschi/discord-date-decider/internal/poll"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) Open() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockService) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockService) SendMessage(channelID string, message *message.Message) (string, error) {
	args := m.Called(channelID, message)
	return args.String(0), args.Error(1)
}

func (m *MockService) SendPoll(channelID string, poll *poll.DatePoll) (string, error) {
	args := m.Called(channelID, poll)
	return args.String(0), args.Error(1)
}

func (m *MockService) PinPoll(channelID string, pollID string) error {
	args := m.Called(channelID, pollID)
	return args.Error(0)
}

func TestStartPoll(t *testing.T) {
	t.Run("successful poll start", func(t *testing.T) {
		mockService := new(MockService)
		bot := NewBot(mockService)
		pollChannelID := "poll-channel-id"
		announcementChannelID := "announcement-channel-id"
		pollID := "poll-id"
		messageID := "message-id"
		request := PollRequest{
			Action:                "startPoll",
			PollChannelID:         pollChannelID,
			AnnouncementChannelID: announcementChannelID,
		}
		mockService.On("Open").Return(nil)
		mockService.On("SendPoll", pollChannelID, mock.AnythingOfType("*poll.DatePoll")).Return(pollID, nil)
		mockService.On("PinPoll", pollChannelID, pollID).Return(nil)
		mockService.On("SendMessage", announcementChannelID, mock.AnythingOfType("*message.Message")).Return(messageID, nil)
		mockService.On("Close").Return(nil)

		err := bot.StartPoll(request)

		assert.NoError(t, err)
		mockService.AssertExpectations(t)
	})

	t.Run("error during open", func(t *testing.T) {
		mockService := new(MockService)
		bot := NewBot(mockService)
		request := PollRequest{
			Action:                "startPoll",
			PollChannelID:         "poll-channel-id",
			AnnouncementChannelID: "announcement-channel-id",
		}
		mockService.On("Open").Return(assert.AnError)

		err := bot.StartPoll(request)

		assert.Error(t, err)
		mockService.AssertExpectations(t)
		mockService.AssertNotCalled(t, "SendPoll")
		mockService.AssertNotCalled(t, "PinPoll")
		mockService.AssertNotCalled(t, "SendMessage")
		mockService.AssertNotCalled(t, "Close")
	})

	t.Run("error during send poll", func(t *testing.T) {
		mockService := new(MockService)
		bot := NewBot(mockService)
		pollChannelID := "poll-channel-id"
		announcementChannelID := "announcement-channel-id"
		request := PollRequest{
			Action:                "startPoll",
			PollChannelID:         pollChannelID,
			AnnouncementChannelID: announcementChannelID,
		}
		mockService.On("Open").Return(nil)
		mockService.On("SendPoll", pollChannelID, mock.AnythingOfType("*poll.DatePoll")).Return("", assert.AnError)
		mockService.On("Close").Return(nil)

		err := bot.StartPoll(request)

		assert.Error(t, err)
		mockService.AssertExpectations(t)
		mockService.AssertNotCalled(t, "PinPoll")
		mockService.AssertNotCalled(t, "SendMessage")
	})

	t.Run("error during pin poll", func(t *testing.T) {
		mockService := new(MockService)
		bot := NewBot(mockService)
		pollChannelID := "poll-channel-id"
		announcementChannelID := "announcement-channel-id"
		pollID := "poll-id"
		request := PollRequest{
			Action:                "startPoll",
			PollChannelID:         pollChannelID,
			AnnouncementChannelID: announcementChannelID,
		}
		mockService.On("Open").Return(nil)
		mockService.On("SendPoll", pollChannelID, mock.AnythingOfType("*poll.DatePoll")).Return(pollID, nil)
		mockService.On("PinPoll", pollChannelID, pollID).Return(assert.AnError)
		mockService.On("Close").Return(nil)

		err := bot.StartPoll(request)

		assert.Error(t, err)
		mockService.AssertExpectations(t)
		mockService.AssertNotCalled(t, "SendMessage")
	})

	t.Run("error during send message", func(t *testing.T) {
		mockService := new(MockService)
		bot := NewBot(mockService)
		pollChannelID := "poll-channel-id"
		announcementChannelID := "announcement-channel-id"
		pollID := "poll-id"
		request := PollRequest{
			Action:                "startPoll",
			PollChannelID:         pollChannelID,
			AnnouncementChannelID: announcementChannelID,
		}
		mockService.On("Open").Return(nil)
		mockService.On("SendPoll", pollChannelID, mock.AnythingOfType("*poll.DatePoll")).Return(pollID, nil)
		mockService.On("PinPoll", pollChannelID, pollID).Return(nil)
		mockService.On("SendMessage", announcementChannelID, mock.AnythingOfType("*message.Message")).Return("", assert.AnError)
		mockService.On("Close").Return(nil)

		err := bot.StartPoll(request)

		assert.Error(t, err)
		mockService.AssertExpectations(t)
	})

	t.Run("error during close", func(t *testing.T) {
		mockService := new(MockService)
		bot := NewBot(mockService)
		pollChannelID := "poll-channel-id"
		announcementChannelID := "announcement-channel-id"
		pollID := "poll-id"
		messageID := "message-id"
		request := PollRequest{
			Action:                "startPoll",
			PollChannelID:         pollChannelID,
			AnnouncementChannelID: announcementChannelID,
		}
		mockService.On("Open").Return(nil)
		mockService.On("SendPoll", pollChannelID, mock.AnythingOfType("*poll.DatePoll")).Return(pollID, nil)
		mockService.On("PinPoll", pollChannelID, pollID).Return(nil)
		mockService.On("SendMessage", announcementChannelID, mock.AnythingOfType("*message.Message")).Return(messageID, nil)
		mockService.On("Close").Return(assert.AnError)

		err := bot.StartPoll(request)

		assert.Error(t, err)
		mockService.AssertExpectations(t)
	})

	t.Run("error during close does not override other errors", func(t *testing.T) {
		mockService := new(MockService)
		bot := NewBot(mockService)
		pollChannelID := "poll-channel-id"
		announcementChannelID := "announcement-channel-id"
		pollID := "poll-id"
		request := PollRequest{
			Action:                "startPoll",
			PollChannelID:         pollChannelID,
			AnnouncementChannelID: announcementChannelID,
		}
		expectedErr := errors.New("some error")
		closeErr := errors.New("error during close")
		mockService.On("Open").Return(nil)
		mockService.On("SendPoll", pollChannelID, mock.AnythingOfType("*poll.DatePoll")).Return(pollID, nil)
		mockService.On("PinPoll", pollChannelID, pollID).Return(nil)
		mockService.On("SendMessage", announcementChannelID, mock.AnythingOfType("*message.Message")).Return("", expectedErr)
		mockService.On("Close").Return(closeErr)

		err := bot.StartPoll(request)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		mockService.AssertExpectations(t)
	})
}

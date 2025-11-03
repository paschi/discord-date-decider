package discord

import (
	"errors"
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/paschi/discord-date-decider/internal/message"
	"github.com/paschi/discord-date-decider/internal/poll"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func (m *MockClient) ChannelMessageUnpin(channelID string, messageID string) error {
	args := m.Called(channelID, messageID)
	return args.Error(0)
}

func (m *MockClient) ChannelMessagesPinned(channelID string) ([]*discordgo.Message, error) {
	args := m.Called(channelID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*discordgo.Message), args.Error(1)
}

func (m *MockClient) ExpirePoll(channelID string, messageID string) error {
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
		testPoll := poll.NewDatePoll("Test Poll", futureDate.Year(), futureDate.Month(), []time.Weekday{time.Friday, time.Saturday}, time.UTC, []int{}, []int{})
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
		testPoll := poll.NewDatePoll("Test Poll", futureDate.Year(), futureDate.Month(), []time.Weekday{time.Friday, time.Saturday}, time.UTC, []int{}, []int{})
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
		testPoll := poll.NewDatePoll("Test Poll", pastDate.Year(), pastDate.Month(), []time.Weekday{time.Friday, time.Saturday}, time.UTC, []int{}, []int{})

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

func TestDefaultService_UnpinPoll(t *testing.T) {
	t.Run("successful unpin poll", func(t *testing.T) {
		mockClient := new(MockClient)
		channelID := "test-channel"
		messageID := "message-id"
		mockClient.On("ChannelMessageUnpin", channelID, messageID).Return(nil)

		service := NewDefaultService(mockClient)
		err := service.UnpinPoll(channelID, messageID)

		assert.NoError(t, err)
		mockClient.AssertExpectations(t)
	})

	t.Run("error during unpin poll", func(t *testing.T) {
		mockClient := new(MockClient)
		channelID := "test-channel"
		messageID := "message-id"
		expectedErr := errors.New("unpin error")
		mockClient.On("ChannelMessageUnpin", channelID, messageID).Return(expectedErr)

		service := NewDefaultService(mockClient)
		err := service.UnpinPoll(channelID, messageID)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, errors.Unwrap(err))
		mockClient.AssertExpectations(t)
	})
}

func TestDefaultService_GetLastPinnedPollResult(t *testing.T) {
	location := time.UTC

	t.Run("successful get last pinned poll result", func(t *testing.T) {
		mockClient := new(MockClient)
		channelID := "test-channel"

		nonPollMsg := &discordgo.Message{ID: "no-poll", Poll: nil}

		answer1 := discordgo.PollAnswer{AnswerID: 0, Media: &discordgo.PollMedia{Text: "Friday, 01.11.2025"}}
		answer2 := discordgo.PollAnswer{AnswerID: 1, Media: &discordgo.PollMedia{Text: "Saturday, 02.11.2025"}}
		pollMsg := &discordgo.Message{
			ID: "poll-123",
			Poll: &discordgo.Poll{
				Question: discordgo.PollMedia{Text: "Test Poll"},
				Answers:  []discordgo.PollAnswer{answer1, answer2},
				Results: &discordgo.PollResults{
					AnswerCounts: []*discordgo.PollAnswerCount{{ID: 0, Count: 5}, {ID: 1, Count: 3}},
					Finalized:    true,
				},
			},
		}

		mockClient.On("ChannelMessagesPinned", channelID).Return([]*discordgo.Message{nonPollMsg, pollMsg}, nil)

		service := NewDefaultService(mockClient)
		result, err := service.GetLastPinnedPollResult(channelID, location)

		assert.NoError(t, err)
		if assert.NotNil(t, result) {
			assert.Equal(t, "poll-123", result.PollID)
			assert.True(t, result.Finalized)
			expectedDate := time.Date(2025, time.November, 1, 20, 0, 0, 0, location)
			assert.Len(t, result.WinningAnswers, 1)
			assert.Equal(t, expectedDate, result.WinningAnswers[0])
		}
		mockClient.AssertExpectations(t)
	})

	t.Run("error when retrieving pinned messages", func(t *testing.T) {
		mockClient := new(MockClient)
		channelID := "test-channel"
		expectedErr := errors.New("pinned messages error")
		mockClient.On("ChannelMessagesPinned", channelID).Return(([]*discordgo.Message)(nil), expectedErr)

		service := NewDefaultService(mockClient)
		result, err := service.GetLastPinnedPollResult(channelID, location)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, errors.Unwrap(err))
		mockClient.AssertExpectations(t)
	})

	t.Run("no pinned poll found", func(t *testing.T) {
		mockClient := new(MockClient)
		channelID := "test-channel"
		mockClient.On("ChannelMessagesPinned", channelID).Return([]*discordgo.Message{{ID: "m1"}, {ID: "m2"}}, nil)

		service := NewDefaultService(mockClient)
		result, err := service.GetLastPinnedPollResult(channelID, location)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "could not find last pinned poll")
		mockClient.AssertExpectations(t)
	})
}

package main

import (
	"context"
	"fmt"
	"github.com/paschi/discord-date-decider/internal/discord"
	"github.com/paschi/discord-date-decider/internal/message"
	"github.com/paschi/discord-date-decider/internal/poll"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

type Bot struct {
	service discord.Service
}

type PollRequest struct {
	Action                string `json:"action"`
	PollChannelID         string `json:"pollChannelId"`
	AnnouncementChannelID string `json:"announcementChannelId"`
}

func main() {
	lambda.Start(handleRequest)
}

func initBot() (*Bot, error) {
	discordToken := os.Getenv("DISCORD_TOKEN")
	if discordToken == "" {
		return nil, fmt.Errorf("could not find discord token in environment")
	}
	client, err := discord.NewDefaultClient(discordToken)
	if err != nil {
		return nil, fmt.Errorf("could not create discord client: %w", err)
	}
	service := discord.NewDefaultService(client)
	return NewBot(service), nil
}

func handleRequest(_ context.Context, request PollRequest) error {
	bot, err := initBot()
	if err != nil {
		return fmt.Errorf("could not initialize bot: %w", err)
	}
	switch request.Action {
	case "startPoll":
		return bot.StartPoll(request)
	case "endPoll":
		return bot.EndPoll(request)
	default:
		log.Printf("unknown action: %s", request.Action)
		return fmt.Errorf("unknown action: %s", request.Action)
	}
}

func NewBot(service discord.Service) *Bot {
	return &Bot{
		service: service,
	}
}

func (b *Bot) StartPoll(request PollRequest) (err error) {
	log.Printf("executing 'startPoll' request: %+v", request)
	err = b.service.Open()
	if err != nil {
		log.Printf("could not open service: %v", err)
		return
	}
	defer func() {
		if closeErr := b.service.Close(); closeErr != nil {
			log.Printf("could not close service: %v", closeErr)
			if err == nil {
				err = closeErr
			}
		}
	}()
	year, month := getNextMonth(time.Now())
	datePoll := poll.NewDatePoll(fmt.Sprintf("Poll for %s %d", month, year), year, month, []time.Weekday{time.Friday, time.Saturday})
	pollID, err := b.service.SendPoll(request.PollChannelID, datePoll)
	if err != nil {
		log.Printf("service could not send poll to poll channel: %v", err)
		return
	}
	log.Printf("service successfully sent poll to poll channel: %s", pollID)
	err = b.service.PinPoll(request.PollChannelID, pollID)
	if err != nil {
		log.Printf("service could not pin poll to poll channel: %v", err)
		return
	}
	log.Printf("service successfully pinned poll to poll channel")
	announcement := message.NewMessage(fmt.Sprintf(`@here :wave: Hey! I justed posted a new poll for %s :calendar:. Check it out! :eyes:
-# Beep boop. I'm a bot. :robot:`, month), true)
	messageID, err := b.service.SendMessage(request.AnnouncementChannelID, announcement)
	if err != nil {
		log.Printf("service could not send message to announcement channel: %v", err)
		return
	}
	log.Printf("service successfully sent message to announcement channel: %s", messageID)
	return
}

func (b *Bot) EndPoll(request PollRequest) (err error) {
	log.Printf("executing 'endPoll' request: %+v", request)
	err = b.service.Open()
	if err != nil {
		log.Printf("could not open service: %v", err)
		return
	}
	defer func() {
		if closeErr := b.service.Close(); closeErr != nil {
			log.Printf("could not close service: %v", closeErr)
			if err == nil {
				err = closeErr
			}
		}
	}()
	result, err := b.service.GetLastPinnedPollResult(request.PollChannelID)
	if err != nil {
		log.Printf("could not retrieve last poll result: %v", err)
		return
	}
	log.Printf("successfully retrieved last poll result: %+v", result)
	if !result.Finalized {
		err = b.service.ExpirePoll(request.PollChannelID, result.PollID)
		if err != nil {
			log.Printf("could not expire poll result: %v", err)
			return
		}
		log.Printf("successfully expired poll")
	}
	err = b.service.UnpinPoll(request.PollChannelID, result.PollID)
	if err != nil {
		log.Printf("could not unpin poll from poll channel: %v", err)
		return
	}
	log.Printf("service successfully unpinned poll from poll channel")
	announcement := message.NewMessage(fmt.Sprintf(`@here We have a winner :trophy:! The next event happens at %s :calendar:. See you then!
-# Beep boop. I'm a bot. :robot:`, result.WinningAnswer), true)
	messageID, err := b.service.SendMessage(request.AnnouncementChannelID, announcement)
	if err != nil {
		log.Printf("service could not send message to announcement channel: %v", err)
		return
	}
	log.Printf("service successfully sent message to announcement channel: %s", messageID)
	return
}

func getNextMonth(now time.Time) (int, time.Month) {
	month := now.AddDate(0, 1, 0)
	return month.Year(), month.Month()
}

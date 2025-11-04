package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/paschi/discord-date-decider/internal/discord"
	"github.com/paschi/discord-date-decider/internal/message"
	"github.com/paschi/discord-date-decider/internal/poll"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/klauspost/lctime"
)

const (
	defaultLocale           = "en_US"
	defaultPollTitle        = "Poll for %s %d"
	defaultStartPollMessage = "@here :wave: Hey! I just posted a new poll for %s :calendar:. Check it out! :eyes:"
	defaultEndPollMessage   = "@here We have a winner :trophy:! The next event happens on <t:%d:F> :calendar:. See you then!"
)

type Bot struct {
	service discord.Service
}

type PollRequest struct {
	Action                string `json:"action"`
	PollChannelID         string `json:"pollChannelId"`
	AnnouncementChannelID string `json:"announcementChannelId"`
	TimeZone              string `json:"timeZone"`
	Locale                string `json:"locale"`
	Title                 string `json:"title"`
	Message               string `json:"message"`
	AdditionalDays        []int  `json:"additionalDays"`
	ExcludedDays          []int  `json:"excludedDays"`
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
	nextMonth := time.Now().AddDate(0, 1, 0)
	year, month := nextMonth.Year(), nextMonth.Month()
	location, err := time.LoadLocation(request.TimeZone)
	if err != nil {
		log.Printf("could not load location '%s': %v", request.TimeZone, err)
		return
	}
	locale := getOrDefault(request.Locale, defaultLocale)
	err = lctime.SetLocale(locale)
	if err != nil {
		log.Printf("could not load locale: %s", locale)
		return
	}
	pollTitle := fmt.Sprintf(getOrDefault(request.Title, defaultPollTitle), lctime.Strftime("%B", nextMonth), year)
	datePoll := poll.NewDatePoll(pollTitle, year, month, []time.Weekday{time.Friday, time.Saturday}, location, request.AdditionalDays, request.ExcludedDays)
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
	messageText := fmt.Sprintf(getOrDefault(request.Message, defaultStartPollMessage), lctime.Strftime("%B", nextMonth))
	announcement := message.NewMessage(messageText, true)
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
	location, err := time.LoadLocation(request.TimeZone)
	if err != nil {
		log.Printf("could not load location '%s': %v", request.TimeZone, err)
		return
	}
	result, err := b.service.GetLastPinnedPollResult(request.PollChannelID, location)
	if err != nil {
		log.Printf("could not retrieve last poll result: %v", err)
		return
	}
	log.Printf("successfully retrieved last poll result: %+v", result)
	if !result.Finalized {
		log.Printf("poll is not yet finalized")
		return fmt.Errorf("poll is not yet finalized")
	}
	err = b.service.UnpinPoll(request.PollChannelID, result.PollID)
	if err != nil {
		log.Printf("could not unpin poll from poll channel: %v", err)
		return
	}
	log.Printf("service successfully unpinned poll from poll channel")
	messageText := fmt.Sprintf(getOrDefault(request.Message, defaultEndPollMessage), getEarliestTime(result.WinningAnswers).Unix())
	announcement := message.NewMessage(messageText, true)
	messageID, err := b.service.SendMessage(request.AnnouncementChannelID, announcement)
	if err != nil {
		log.Printf("service could not send message to announcement channel: %v", err)
		return
	}
	log.Printf("service successfully sent message to announcement channel: %s", messageID)
	return
}

func getEarliestTime(times []time.Time) time.Time {
	var earliest time.Time
	for i, t := range times {
		if i == 0 {
			earliest = t
			continue
		}
		if t.Before(earliest) {
			earliest = t
		}
	}
	return earliest
}

func getOrDefault(value string, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

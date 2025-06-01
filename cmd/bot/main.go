package main

import (
	"context"
	"fmt"
	"github.com/paschi/discord-date-decider/internal/discord"
	"github.com/paschi/discord-date-decider/internal/poll"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

type PollRequest struct {
	Action              string `json:"action"`
	PollChannel         string `json:"pollChannel"`
	AnnouncementChannel string `json:"announcementChannel"`
}

var (
	discordToken string
)

func init() {
	discordToken = os.Getenv("DISCORD_TOKEN")
	if discordToken == "" {
		log.Fatalf("discord token not found in environment")
	}
}

func main() {
	lambda.Start(handleRequest)
}

func handleRequest(context context.Context, request PollRequest) error {
	switch request.Action {
	case "startPoll":
		return startPoll(request)
	default:
		log.Printf("unknown action: %s", request.Action)
		return fmt.Errorf("unknown action: %s", request.Action)
	}
}

func startPoll(request PollRequest) error {
	service, err := discord.NewDiscordService(discordToken)
	if err != nil {
		return fmt.Errorf("could not create discord service: %v", err)
	}
	defer service.Close()
	year, month := getNextMonth(time.Now())
	datePoll := poll.NewDatePoll(fmt.Sprintf("Poll for %s %d", month, year), year, month, []time.Weekday{time.Friday, time.Saturday})
	log.Printf("creating poll in channel %s: %+v", request.PollChannel, datePoll)
	_, err = service.SendPoll(request.PollChannel, datePoll)
	if err != nil {
		log.Printf("error sending poll to channel %s: %v", request.PollChannel, err)
		return err
	}
	log.Printf("poll sent to channel %s", request.PollChannel)
	_, err = service.SendMessage(request.AnnouncementChannel, fmt.Sprintf(
		`@here :wave: Hey! I justed posted a new poll for %s :calendar:. Check it out! :eyes:
-# *Beep boop.* I'm a bot. :robot:`, month))
	if err != nil {
		log.Printf("error sending message to channel %s: %v", request.AnnouncementChannel, err)
		return err
	}
	return nil
}

func getNextMonth(now time.Time) (int, time.Month) {
	month := now.AddDate(0, 1, 0)
	return month.Year(), month.Month()
}

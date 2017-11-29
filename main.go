package main

import (
	"log"
	"os"

	"github.com/nlopes/slack"
	"github.com/unimicro/unibot/auth"
	"github.com/unimicro/unibot/gitter"
	"github.com/unimicro/unibot/webhooks"
)

const (
	tokenFileLocation = "./tokens.json"
)

var botID = "N/A"

func main() {
	authTokens := auth.ReadTokenFile(tokenFileLocation)
	slackToken := authTokens.Slack
	gitterToken := authTokens.Gitter

	api := slack.New(string(slackToken))
	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)
	api.SetDebug(false)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	go webhooks.StartWebhooksServer(rtm)
	go gitter.Listen(rtm, gitterToken)
	readRtmStream(rtm, logger)
}

func readRtmStream(rtm *slack.RTM, logger *log.Logger) {
	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.ConnectedEvent:
			botID = ev.Info.User.ID
			log.Printf("Listening for slack messages (%d)...\n", ev.ConnectionCount)
		case *slack.MessageEvent:
			if ev.User != botID {
				handleMessage(ev, rtm)
			}
		case *slack.RTMError:
			logger.Printf("Error: %s\n", ev.Error())
		case *slack.InvalidAuthEvent:
			logger.Println("Invalid credentials")
			return
		}
	}
}

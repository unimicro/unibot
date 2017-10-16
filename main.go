package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/nlopes/slack"
	"github.com/unimicro/unibot/webhooks"
)

const (
	slackTokenFilePath = "./slack_api.token"
)

var botID = "N/A"

func main() {
	slackToken, err := ioutil.ReadFile(slackTokenFilePath)
	if err != nil {
		panic("Couldn't read the slack token file:" + string(slackTokenFilePath))
	}

	api := slack.New(strings.TrimSpace(string(slackToken)))
	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)
	api.SetDebug(false)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	go webhooks.StartWebhooksServer(rtm)

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore hello

		case *slack.ConnectedEvent:
			botID = ev.Info.User.ID
			//logger.Println("Infos:", ev.Info)
			logger.Println("Connection counter:", ev.ConnectionCount)

		case *slack.MessageEvent:
			if ev.User != botID {
				handleMessage(ev, rtm)
			}

		case *slack.PresenceChangeEvent:
			// logger.Printf("Presence Change: %v\n", ev)

		case *slack.LatencyReport:
			// logger.Printf("Current latency: %v\n", ev.Value)

		case *slack.RTMError:
			logger.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			logger.Printf("Invalid credentials")
			return

		default:

			// Ignore other events..
			// logger.Printf("Unexpected: %v\n", msg.Data)
		}
	}
}

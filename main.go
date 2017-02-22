package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/nlopes/slack"
)

const (
	SLACK_TOKEN_FILE_PATH = "./slack_api.token"
)

var botID = "N/A"

func main() {
	slackToken, err := ioutil.ReadFile(SLACK_TOKEN_FILE_PATH)
	if err != nil {
		panic("Couldn't read the slack token file:" + string(SLACK_TOKEN_FILE_PATH))
	}

	api := slack.New(string(slackToken))
	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)
	api.SetDebug(false)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore hello

		case *slack.ConnectedEvent:
			botID = ev.Info.User.ID
			//fmt.Println("Infos:", ev.Info)
			fmt.Println("Connection counter:", ev.ConnectionCount)

		case *slack.MessageEvent:
			if ev.User != botID {
				HandleMessage(ev, rtm)
			}

		case *slack.PresenceChangeEvent:
			// fmt.Printf("Presence Change: %v\n", ev)

		case *slack.LatencyReport:
			// fmt.Printf("Current latency: %v\n", ev.Value)

		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return

		default:

			// Ignore other events..
			// fmt.Printf("Unexpected: %v\n", msg.Data)
		}
	}
}

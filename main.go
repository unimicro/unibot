package main

import (
    "log"
    "os"

    "github.com/nlopes/slack"
    "github.com/unimicro/unibot/auth"
    "github.com/unimicro/unibot/gitter"
    "github.com/unimicro/unibot/logger"
)

var botID = "N/A"

func main() {
    tokens := auth.GetTokens()

    api := slack.New(tokens.Slack.AsString())
    slackLogger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
    slack.SetLogger(slackLogger)
    api.SetDebug(false)

    rtm := api.NewRTM()
    go rtm.ManageConnection()

    go gitter.Listen(rtm, tokens.Gitter)

    logger.SetRTM(rtm)

    readRtmStream(rtm, slackLogger)
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

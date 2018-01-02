package gitter

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/nlopes/slack"
	"github.com/unimicro/unibot/auth"
	"github.com/unimicro/unibot/constants"
)

const (
	economyRoomID   = "5a1d5911d73408ce4f80a64d"
	economyRoomName = "Economy"
	unibotRoomID    = "5a1d8778d73408ce4f80af80"
	unibotRoomName  = "unibot-test-room"
	roomUrl         = "https://gitter.im/unimicro/%s"
	roomApiUrl      = "https://stream.gitter.im/v1/rooms/%s/chatMessages"
	heartbeat       = " \n"
	quietPeriod     = time.Duration(time.Minute * 10)
)

func Listen(rtm *slack.RTM, gitterToken auth.Token) {
	client := &http.Client{}
	var (
		gitterUrl           string
		roomName            string
		slackReceiverRoomID string
	)
	if isDevelopEnvironment() {
		gitterUrl = fmt.Sprintf(roomApiUrl, unibotRoomID)
		slackReceiverRoomID = constants.UnibotLogChannelID
		roomName = unibotRoomName
	} else {
		gitterUrl = fmt.Sprintf(roomApiUrl, economyRoomID)
		slackReceiverRoomID = constants.DevelopersChannelID
		roomName = economyRoomName
	}
	request, err := http.NewRequest("GET", gitterUrl, nil)
	if err != nil {
		panic("ERROR starting initial gitter request: " + err.Error())
	}
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", gitterToken))
	var (
		lastMessageWasSent time.Time
		waitMultiplier     = 1
	)
	for {
		response, err := client.Do(request)
		if err != nil {
			logToSlackTestChannel(
				rtm,
				fmt.Sprintf(
					"ERROR making http request to gitter (retrying in %d seconds): %s\n",
					waitMultiplier,
					err.Error(),
				),
			)
			time.Sleep(time.Second * time.Duration(waitMultiplier))
			waitMultiplier *= 2
		}
		defer response.Body.Close()

		log.Println("Listening for gitter messages...")

		reader := bufio.NewReader(response.Body)
		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				logToSlackTestChannel(
					rtm,
					"ERROR reading message from gitter: '"+err.Error()+"'",
				)
				break
			}
			if string(line) == heartbeat {
				// Dropping heartbeat
				continue
			}
			if time.Now().Before(lastMessageWasSent.Add(quietPeriod)) {
				// Dropping because still inside quited period
				continue
			}

			var messageFromGitter GitterMessage
			err = json.Unmarshal(line, &messageFromGitter)
			if err != nil {
				logToSlackTestChannel(
					rtm,
					fmt.Sprintf("ERROR parsing JSON message from gitter: %s\nThe JSON:\n%s", err.Error(), line),
				)
				continue
			}

			lastMessageWasSent = time.Now()
			messageToSlack := fmt.Sprintf(
				"%s: %s?at=%s\n%s",
				messageFromGitter.FromUser.Username,
				fmt.Sprintf(roomUrl, roomName),
				messageFromGitter.ID,
				messageFromGitter.Text,
			)
			rtm.SendMessage(rtm.NewOutgoingMessage(messageToSlack, slackReceiverRoomID))
		}
	}
}

func isDevelopEnvironment() bool {
	return os.Getenv("DEVELOP") != ""
}

func logToSlackTestChannel(rtm *slack.RTM, message string) {
	log.Printf("Sendt this message to slack log: %s\n", message)
	rtm.SendMessage(rtm.NewOutgoingMessage(
		message,
		constants.UnibotLogChannelID,
	))
}

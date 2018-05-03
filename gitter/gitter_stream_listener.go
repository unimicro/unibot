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
	quietPeriod     = time.Duration(time.Hour * 1)
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
		lastFullMessageWasSent time.Time
		lastFullMessageID      string
		waitMultiplier         = 1
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
			continue
		}
		waitMultiplier *= 1
		defer response.Body.Close()

		log.Println("Listening for gitter messages...")

		reader := bufio.NewReader(response.Body)
		for {
			isInQuietPeriod := time.Now().Before(lastFullMessageWasSent.Add(quietPeriod))
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

			var messageFromGitter GitterMessage
			err = json.Unmarshal(line, &messageFromGitter)
			if err != nil {
				logToSlackTestChannel(
					rtm,
					fmt.Sprintf("ERROR parsing JSON message from gitter: %s\nThe JSON:\n%s", err.Error(), line),
				)
				continue
			}

			heading := fmt.Sprintf(
				"<%s?at=%s|open in gitter>:",
				fmt.Sprintf(roomUrl, roomName),
				messageFromGitter.ID,
			)
			params := slack.PostMessageParameters{
				IconURL:  messageFromGitter.FromUser.AvatarURLSmall,
				Username: messageFromGitter.FromUser.DisplayName + " (gitter)",
				Attachments: []slack.Attachment{
					{
						Pretext:    heading,
						Text:       messageFromGitter.Text,
						MarkdownIn: []string{"pretext"},
					},
				},
			}

			// Send messages as replies to last "full" message if in quiet period
			if isInQuietPeriod {
				params.ThreadTimestamp = lastFullMessageID
			}

			_, postedMessageID, err := rtm.Client.PostMessage(slackReceiverRoomID, "", params)
			if err != nil {
				logToSlackTestChannel(rtm, "Got an error back when posting a message to slack: "+err.Error())
				continue
			}

			if !isInQuietPeriod {
				lastFullMessageID = postedMessageID
				lastFullMessageWasSent = time.Now()
			}
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

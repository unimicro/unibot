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
	economyGitterRoomID  = "5a1d5911d73408ce4f80a64d"
	economyGitterRoomUrl = "https://gitter.im/unimicro/Economy"
	unibotGitterRoomID   = "5a1d8778d73408ce4f80af80"
	gitterEconomyRoomUrl = "https://stream.gitter.im/v1/rooms/" + economyGitterRoomID + "/chatMessages"
	gitterUnibotRoomUrl  = "https://stream.gitter.im/v1/rooms/" + unibotGitterRoomID + "/chatMessages"
	gitterHeartbeat      = " \n"
	quietPeriod          = time.Duration(time.Minute * 10)
)

func Listen(rtm *slack.RTM, gitterToken auth.Token) {
	client := &http.Client{}
	var (
		request *http.Request
		err     error
		roomUrl string
		roomID  string
	)
	if os.Getenv("DEVELOP") != "" {
		request, err = http.NewRequest("GET", gitterUnibotRoomUrl, nil)
		roomID = constants.UnibotLogChannelID
		roomUrl = gitterUnibotRoomUrl
	} else {
		request, err = http.NewRequest("GET", gitterEconomyRoomUrl, nil)
		roomID = constants.DevelopersChannelID
		roomUrl = economyGitterRoomUrl
	}
	if err != nil {
		panic("ERROR starting initial gitter request: " + err.Error())
	}
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", gitterToken))
	var (
		lastMessageWasSent time.Time
		waitMultiplier     = 16
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
			if string(line) == gitterHeartbeat {
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
				"%s on gitter: %s?at=%s",
				messageFromGitter.FromUser.DisplayName,
				roomUrl,
				messageFromGitter.ID,
			)
			rtm.SendMessage(rtm.NewOutgoingMessage(messageToSlack, roomID))
		}
	}
}

func logToSlackTestChannel(rtm *slack.RTM, message string) {
	log.Printf("Sendt this message to slack log: %s\n", message)
	rtm.SendMessage(rtm.NewOutgoingMessage(
		message,
		constants.UnibotLogChannelID,
	))
}

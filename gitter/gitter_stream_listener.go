package gitter

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math"
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

var (
	client               = &http.Client{}
	isDevelopEnvironment = os.Getenv("DEVELOP") != ""
)

func Listen(rtm *slack.RTM, gitterToken auth.Token) {
	var (
		gitterUrl           string
		roomName            string
		slackReceiverRoomID string
	)
	if isDevelopEnvironment {
		gitterUrl = fmt.Sprintf(roomApiUrl, unibotRoomID)
		slackReceiverRoomID = constants.UnibotLogChannelID
		roomName = unibotRoomName
	} else {
		gitterUrl = fmt.Sprintf(roomApiUrl, economyRoomID)
		slackReceiverRoomID = constants.DevelopersChannelID
		roomName = economyRoomName
	}
	var (
		lastFullMessageWasSent time.Time
		lastFullMessageID      string
		retries                = 0
	)
	for {
		time.Sleep(time.Second * time.Duration(pow(2, retries)))
		response, err := makeHttpRequest(gitterUrl, gitterToken)
		if err != nil {
			logToSlackTestChannel(
				rtm,
				fmt.Sprintf(
					"ERROR making http request to gitter (retrying in %d seconds): %s\n",
					pow(2, retries),
					err,
				),
			)
			retries++
			continue
		}
		defer response.Body.Close()

		log.Println("Listening for gitter messages...")

		reader := bufio.NewReader(response.Body)
		for {
			isInQuietPeriod := time.Now().Before(lastFullMessageWasSent.Add(quietPeriod))
			line, err := reader.ReadBytes('\n')
			if err != nil {
				logToSlackTestChannel(
					rtm,
					fmt.Sprintf(
						"ERROR reading message from gitter (retrying in %d seconds): '%s'",
						pow(2, retries),
						err,
					),
				)
				retries++
				break
			}
			retries = 0
			if string(line) == heartbeat {
				// Dropping heartbeat
				continue
			}

			var messageFromGitter GitterMessage
			err = json.Unmarshal(line, &messageFromGitter)
			if err != nil {
				logToSlackTestChannel(
					rtm,
					fmt.Sprintf("ERROR parsing JSON message from gitter: %s\nThe JSON:\n%s", err, line),
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

func logToSlackTestChannel(rtm *slack.RTM, message string) {
	log.Printf("Sendt this message to slack log: %s\n", message)
	rtm.SendMessage(rtm.NewOutgoingMessage(
		message,
		constants.UnibotLogChannelID,
	))
}

func makeHttpRequest(url string, bearerToken auth.Token) (*http.Response, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Couldn't start initial gitter request: %s", err)
	}
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", bearerToken))
	response, err := client.Do(request)
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("%d statuscode returned from gitter\n", response.StatusCode)
	}
	return response, nil
}

func pow(x, y int) int {
	return int(math.Pow(float64(x), float64(y)))
}

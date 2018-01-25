package logger

import (
	"log"
	"strings"

	"github.com/nlopes/slack"
	"github.com/unimicro/unibot/constants"
)

var _rtm *slack.RTM

func SetRTM(rtm *slack.RTM) {
	_rtm = rtm
}

func Log(message ...string) {
	log.Printf("Sendt this message to slack log: %s\n", strings.Join(message, " "))
	_rtm.SendMessage(_rtm.NewOutgoingMessage(
		strings.Join(message, " "),
		constants.UnibotLogChannelID,
	))

}

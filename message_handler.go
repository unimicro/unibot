package main

import (
	"strings"

	"github.com/nlopes/slack"
)

const (
	DIRECT_CHANNEL_PREFIX = "D"
)

func HandleMessage(ev *slack.MessageEvent, rtm *slack.RTM) {
	if strings.HasPrefix(strings.ToLower(ev.Text), TRAVEL_TEXT_COMMAND_PREFIX) && strings.HasPrefix(ev.Channel, DIRECT_CHANNEL_PREFIX) {
		HandleTravelTextCommand(ev, rtm)
	}
}

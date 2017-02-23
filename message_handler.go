package main

import (
	"strings"

	"github.com/nlopes/slack"
	"github.com/unimicro/unibot/traveltext"
)

const (
	DIRECT_CHANNEL_PREFIX = "D"
)

func handleMessage(ev *slack.MessageEvent, rtm *slack.RTM) {
	if strings.HasPrefix(strings.ToLower(ev.Text), traveltext.TravelTextCommandPrefix) && strings.HasPrefix(ev.Channel, DIRECT_CHANNEL_PREFIX) {
		traveltext.HandleTravelTextCommand(ev, rtm)
	}
}

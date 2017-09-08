package main

import (
	"regexp"
	"strings"

	"github.com/nlopes/slack"
	"github.com/unimicro/unibot/traveltext"
)

const (
	directChannelPrefix = "D"
	publicChannelPrefix = "C"
)

func handleMessage(ev *slack.MessageEvent, rtm *slack.RTM) {
	if strings.HasPrefix(strings.ToLower(ev.Text), traveltext.TravelTextCommandPrefix) && strings.HasPrefix(ev.Channel, directChannelPrefix) {
		traveltext.HandleTravelTextCommand(ev, rtm)
	}
	if strings.HasPrefix(ev.Text, "<@U3C7XCU5S> tt ") && strings.HasPrefix(ev.Channel, publicChannelPrefix) {
		rtm.SendMessage(rtm.NewOutgoingMessage("You're trying to send a TravelText command in a public channel, please talk to me privatly for that: https://unimicro.slack.com/messages/@unibot", ev.Channel))
	}

	r, _ := regexp.Compile("^[[:alpha:]]{2,3}-\\d+|\\s[[:alpha:]]{1,3}-\\d+")
	matches := r.FindAllString(ev.Text, -1)
	if len(matches) > 0 && strings.HasPrefix(ev.Channel, publicChannelPrefix) {
		for _, m := range matches {
			rtm.SendMessage(rtm.NewOutgoingMessage(
				"https://unimicro.atlassian.net/browse/"+strings.ToUpper(strings.TrimSpace(m)),
				ev.Channel,
			))
		}
	}
}

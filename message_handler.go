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
	uniBotUserId        = "<@U3C7XCU5S>"
)

var jiraKeyRegex = regexp.MustCompile("^[[:alpha:]]{2,3}-\\d+|\\s[[:alpha:]]{1,3}-\\d+")

func handleMessage(ev *slack.MessageEvent, rtm *slack.RTM) {

	switch {
	case strings.HasPrefix(strings.ToLower(ev.Text), traveltext.TravelTextCommandPrefix) && strings.HasPrefix(ev.Channel, directChannelPrefix):
		traveltext.HandleTravelTextCommand(ev, rtm)
	case strings.HasPrefix(ev.Text, uniBotUserId+" tt ") && strings.HasPrefix(ev.Channel, publicChannelPrefix):
		handleTravelTextInPulicChannel(ev, rtm)
	case strings.HasPrefix(ev.Channel, publicChannelPrefix):
		handleJiraKeyInMessage(ev, rtm)
	}
}

func handleTravelTextInPulicChannel(ev *slack.MessageEvent, rtm *slack.RTM) {
	rtm.SendMessage(rtm.NewOutgoingMessage("You're trying to send a TravelText command in a public channel, please talk to me privatly for that: https://unimicro.slack.com/messages/@unibot", ev.Channel))
}

func handleJiraKeyInMessage(ev *slack.MessageEvent, rtm *slack.RTM) {
	matches := jiraKeyRegex.FindAllString(ev.Text, -1)
	for _, m := range matches {
		rtm.SendMessage(rtm.NewOutgoingMessage(
			"https://unimicro.atlassian.net/browse/"+strings.ToUpper(strings.TrimSpace(m)),
			ev.Channel,
		))
	}
}

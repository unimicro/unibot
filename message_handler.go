package main

import (
	"strings"

	"github.com/nlopes/slack"
	"github.com/unimicro/unibot/jenkins"
	"github.com/unimicro/unibot/jira"
	"github.com/unimicro/unibot/traveltext"
)

const (
	directChannelPrefix   = "D"
	publicChannelPrefix   = "C"
	uniBotUserId          = "<@U3C7XCU5S>"
	pulicTTCommandWarning = "You're trying to send a TravelText command in a public channel" +
		", please talk to me privately for that: https://unimicro.slack.com/messages/@unibot"
)

func handleMessage(ev *slack.MessageEvent, rtm *slack.RTM) {
	switch {
	case isTravelTextInPublicChannel(ev):
		rtm.SendMessage(rtm.NewOutgoingMessage(pulicTTCommandWarning, ev.Channel))
	case isDMTravelTextCommand(ev):
		traveltext.HandleTravelTextCommand(ev, rtm)
	case strings.HasPrefix(strings.ToLower(ev.Text), jenkins.JenkinsCommandPrefix):
		jenkins.HandleJenkinsCommand(ev, rtm)
	case strings.HasPrefix(ev.Channel, publicChannelPrefix):
		jira.HandleJiraKeyInMessage(ev, rtm)
	}
}

func isDMTravelTextCommand(ev *slack.MessageEvent) bool {
	return strings.HasPrefix(strings.ToLower(ev.Text), traveltext.TravelTextCommandPrefix) &&
		strings.HasPrefix(ev.Channel, directChannelPrefix)
}

func isTravelTextInPublicChannel(ev *slack.MessageEvent) bool {
	return strings.HasPrefix(ev.Text, uniBotUserId+" tt ") &&
		strings.HasPrefix(ev.Channel, publicChannelPrefix)
}

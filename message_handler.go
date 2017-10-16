package main

import (
	"log"
	"regexp"
	"strings"

	"github.com/nlopes/slack"
	"github.com/unimicro/unibot/jenkins"
	"github.com/unimicro/unibot/traveltext"
)

const (
	directChannelPrefix   = "D"
	publicChannelPrefix   = "C"
	uniBotUserId          = "<@U3C7XCU5S>"
	pulicTTCommandWarning = "You're trying to send a TravelText command in a public channel" +
		", please talk to me privately for that: https://unimicro.slack.com/messages/@unibot"
)

var jiraKeyRegex = regexp.MustCompile("^[[:alpha:]]{2,3}-\\d+|\\s[[:alpha:]]{1,3}-\\d+")

func handleMessage(ev *slack.MessageEvent, rtm *slack.RTM) {
	log.Println(ev.Text)
	log.Println(jenkins.JenkinsCommandPrefix)
	switch {
	case isTravelTextInPublicChannel(ev):
		rtm.SendMessage(rtm.NewOutgoingMessage(pulicTTCommandWarning, ev.Channel))
	case isDMTravelTextCommand(ev):
		traveltext.HandleTravelTextCommand(ev, rtm)
	case strings.HasPrefix(strings.ToLower(ev.Text), jenkins.JenkinsCommandPrefix):
		jenkins.HandleJenkinsCommand(ev, rtm)
	case strings.HasPrefix(ev.Channel, publicChannelPrefix):
		handleJiraKeyInMessage(ev, rtm)
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

func handleJiraKeyInMessage(ev *slack.MessageEvent, rtm *slack.RTM) {
	matches := jiraKeyRegex.FindAllString(ev.Text, -1)
	for _, m := range matches {
		rtm.SendMessage(rtm.NewOutgoingMessage(
			"https://unimicro.atlassian.net/browse/"+strings.ToUpper(strings.TrimSpace(m)),
			ev.Channel,
		))
	}
}

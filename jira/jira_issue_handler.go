package jira

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/nlopes/slack"
	"github.com/unimicro/unibot/auth"
	"github.com/unimicro/unibot/logger"
)

var jiraKeyRegex = regexp.MustCompile("^[[:alpha:]]{2,3}-\\d+|\\s[[:alpha:]]{1,3}-\\d+")

func HandleJiraKeyInMessage(ev *slack.MessageEvent, rtm *slack.RTM) {
	issueKeys := findIssueKeys(ev.Text)
	for _, jiraIssue := range getJiraIssueInfo(issueKeys) {
		heading := fmt.Sprintf(
			"<https://unimicro.atlassian.net/browse/%[1]s|%[1]s> %[2]s",
			jiraIssue.Key,
			jiraIssue.Fields.Summary,
		)

		params := slack.PostMessageParameters{
			IconURL:  jiraIssue.Fields.Issuetype.IconURL,
			Username: "Jira",
			Attachments: []slack.Attachment{
				{
					Pretext:    heading,
					MarkdownIn: []string{"pretext"},
				},
			},
		}

		var footerInfo []string
		if jiraIssue.Fields.Assignee.DisplayName != "" {
			footerInfo = append(footerInfo, "Assignee: "+jiraIssue.Fields.Assignee.DisplayName)
			params.Attachments[0].FooterIcon = jiraIssue.Fields.Assignee.AvatarUrls.One6X16
		}
		footerInfo = append(footerInfo, "Status: "+jiraIssue.Fields.Status.Name)
		if jiraIssue.Fields.Resolution.Name != "" {
			footerInfo = append(footerInfo, "Resolution: "+jiraIssue.Fields.Resolution.Name)
			params.Attachments[0].Footer = strings.Join(footerInfo, ", ")
		}

		params.Attachments[0].Footer = strings.Join(footerInfo, ", ")

		_, _, err := rtm.Client.PostMessage(ev.Channel, "", params)
		if err != nil {
			logger.Log("[JIRA] Error when posting message", err.Error())
		}
	}
}

func findIssueKeys(text string) []string {
	matches := jiraKeyRegex.FindAllString(text, -1)
	issueKeys := Map(Map(matches, strings.ToUpper), strings.TrimSpace)
	return removeDuplicatesUnordered(issueKeys)
}

func Map(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func removeDuplicatesUnordered(elements []string) []string {
	encountered := map[string]bool{}

	for v := range elements {
		encountered[elements[v]] = true
	}

	result := []string{}
	for key, _ := range encountered {
		result = append(result, key)
	}
	return result
}

func getJiraIssueInfo(keys []string) (issues []JiraIssue) {
	tokens := auth.GetTokens()

	client := &http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	for _, key := range keys {
		url := "https://unimicro.atlassian.net/rest/api/2/issue/" + key
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			continue
		}
		req.Header.Add(
			"Authorization",
			"Basic "+base64.StdEncoding.EncodeToString([]byte("karlgustav@unimicro.no:"+tokens.Jira)),
		)
		resp, err := client.Do(req)
		if err != nil {
			logger.Log("[JIRA] Couldn't make request:", err.Error())
			continue
		}
		if resp.StatusCode == 404 {
			// Issue not found, ignoring
			continue
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Log("[JIRA] Couldn't read body from request:", err.Error())
			continue
		}
		var issue JiraIssue
		err = json.Unmarshal(body, &issue)
		if err != nil {
			logger.Log("[JIRA] Couldn't unmarshal body into jira issue:", err.Error(), "\nBody:", string(body))
			continue
		}

		issues = append(issues, issue)
	}

	return
}

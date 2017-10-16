package jenkins

import (
	"encoding/json"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/nlopes/slack"
	"github.com/unimicro/unibot/storage"
)

const (
	JenkinsCommandPrefix       = "jenkins "
	jenkinsAddCommandPrefix    = JenkinsCommandPrefix + "add "
	jenkinsRemoveCommandPrefix = JenkinsCommandPrefix + "remove "
	jenkinsListCommandPrefix   = JenkinsCommandPrefix + "list"
	helpMessage                = "Available commands:\n" +
		"help => this message\n" +
		"listen => write the jenkins job name after this command to notify this channel when the job builds\n" +
		"stop => write the jenkins job name after this command to stop notifying this channel about the job \n" +
		"list => list all jobs that is sendt to this channel"
)

func HandleJenkinsCommand(ev *slack.MessageEvent, rtm *slack.RTM) {
	var message string
	switch {
	case strings.HasPrefix(ev.Text, JenkinsCommandPrefix+"help"):
		message = helpMessage
	case strings.HasPrefix(ev.Text, jenkinsAddCommandPrefix):
		jobName := ev.Text[len(jenkinsAddCommandPrefix):len(ev.Text)]

		err := storage.DB.Update(func(tx *bolt.Tx) error {
			bucket, err := tx.CreateBucketIfNotExists([]byte(storage.JenkinsBucket))
			if err != nil {
				return err
			}
			jsonString := bucket.Get([]byte(ev.Channel))
			var jobs map[string]bool
			err = json.Unmarshal(jsonString, &jobs)
			if err != nil {
				jobs = make(map[string]bool)
			}
			if jobs[jobName] {
				message = "You already have this jenkins job in this channel: " + jobName
				return nil
			}
			jobs[jobName] = true
			jsonString, err = json.Marshal(jobs)
			if err != nil {
				return err
			}
			err = bucket.Put([]byte(ev.Channel), jsonString)
			if err != nil {
				return err
			}

			message = "Added \"" + jobName + "\" to the list of jobs sendt to this channel"
			return nil
		})

		if err != nil {
			message = err.Error()
		}
	case strings.HasPrefix(ev.Text, jenkinsListCommandPrefix):
		err := storage.DB.View(func(tx *bolt.Tx) error {
			var jobs map[string]bool
			bucket := tx.Bucket([]byte(storage.JenkinsBucket))
			if bucket == nil {
				jobs = make(map[string]bool)
			} else {
				jsonString := bucket.Get([]byte(ev.Channel))
				err := json.Unmarshal(jsonString, &jobs)
				if err != nil {
					jobs = make(map[string]bool)
				}
			}
			if len(jobs) > 0 {
				var keys []string
				for k := range jobs {
					keys = append(keys, k)
				}
				message = "All jobs sendt to this channel:\n" + strings.Join(keys, "\n")
			} else {
				message = "there are no jobs sendt to this channel"
			}
			return nil
		})

		if err != nil {
			message = "ERROR: " + err.Error()
		}
	case strings.HasPrefix(ev.Text, jenkinsRemoveCommandPrefix):
		jobName := ev.Text[len(jenkinsRemoveCommandPrefix):len(ev.Text)]

		err := storage.DB.Update(func(tx *bolt.Tx) error {
			bucket, err := tx.CreateBucketIfNotExists([]byte(storage.JenkinsBucket))
			if err != nil {
				return err
			}
			jsonString := bucket.Get([]byte(ev.Channel))
			var jobs map[string]bool
			err = json.Unmarshal(jsonString, &jobs)
			if err != nil {
				jobs = make(map[string]bool)
			}
			if !jobs[jobName] {
				message = jobName + " is not in this channel"
				return nil
			}
			delete(jobs, jobName)
			jsonString, err = json.Marshal(jobs)
			if err != nil {
				return err
			}
			err = bucket.Put([]byte(ev.Channel), jsonString)
			if err != nil {
				return err
			}

			message = "Removed " + jobName + " from the list of jobs sendt to this channel"
			return nil
		})

		if err != nil {
			message = err.Error()
		}
	}
	rtm.SendMessage(rtm.NewOutgoingMessage(message, ev.Channel))
}

func stringInSlice(list []string, a string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

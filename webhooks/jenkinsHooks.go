package webhooks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/unimicro/unibot/storage"
)

const (
	unibotTestChannel = "C47FWA3S4"
	backendChannel    = "C1BAXE95X"
	frontendChannel   = "C0FQMKYA1"
)

func jenkinsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var message string
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			message = "ERROR: Could not get request.body!"

		}
		obj := &JenkinsNotification{}
		err = json.Unmarshal(data, obj)

		if err != nil {
			message = fmt.Sprintf(
				"ERROR: Could not parse json body as Jenkins notification. Error:\n%s\nData:\n%s",
				err.Error(),
				data,
			)
		} else {
			err = storage.DB.View(func(tx *bolt.Tx) error {
				bucket := tx.Bucket([]byte(storage.JenkinsBucket))
				if bucket == nil {
					return nil
				}
				bucket.ForEach(func(channelName, jsonString []byte) error {
					var jobs map[string]bool
					err = json.Unmarshal(jsonString, &jobs)
					if err != nil {
						return nil
					}
					for job, _ := range jobs {
						if string(job) == obj.Name {
							message = fmt.Sprintf(
								"Name: \"%v\", Phase: %v, Status: %v",
								obj.DisplayName,
								obj.Build.Phase,
								obj.Build.Status,
							)
							messageBus.SendMessage(messageBus.NewOutgoingMessage(message, string(channelName)))
							break
						}
					}
					return nil
				})
				return nil
			})

			if err != nil {
				messageBus.SendMessage(messageBus.NewOutgoingMessage(err.Error(), unibotTestChannel))
			}

		}

	default:
		w.WriteHeader(404)
		w.Write([]byte("This endpoint only supports POST"))
	}
}

package traveltext

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/nlopes/slack"
	"github.com/unimicro/unibot/storage"
)

const (
	secretRandLenght        = 32
	baseUrl                 = "https://traveltext-jarvis.azurewebsites.net"
	ttLoginUrl              = baseUrl + "/Signup"
	ttCommandUrl            = baseUrl + "/api/UniHour/Command"
	TravelTextCommandPrefix = "tt "
)

var (
	errUnAuthenticated = errors.New("User is un-authenticated")
)

func HandleTravelTextCommand(ev *slack.MessageEvent, rtm *slack.RTM) {
	var message string
	var userString []byte
	err := storage.DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(storage.TraveltextBucket))
		if err != nil {
			return err
		}
		userString = bucket.Get([]byte(ev.User))
		return nil
	})
	rtm.SendMessage(rtm.NewOutgoingMessage("...", ev.Channel))
	if userString == nil {
		message = authenticateUser(ev, rtm)
	} else {
		var user User
		err = json.Unmarshal(userString, &user)
		if err != nil {
			message = fmt.Sprintf("Could not unmarshal bytestring to user object: %v", userString)
		} else {
			cmd := strings.TrimPrefix(strings.ToLower(ev.Text), TravelTextCommandPrefix)
			response, err := runTTCommand(cmd, user)
			switch {
			case err == errUnAuthenticated:
				message = authenticateUser(ev, rtm)
			case err != nil:
				message = "Got error when trying to send command to TravelText: " + err.Error()
			default:
				message = response.FriendlyResponse
			}
		}
	}

	rtm.SendMessage(rtm.NewOutgoingMessage(message, ev.Channel))

	if err != nil {
		rtm.SendMessage(rtm.NewOutgoingMessage(err.Error(), ev.Channel))
	}

}

func authenticateUser(ev *slack.MessageEvent, rtm *slack.RTM) string {
	user, err := addUser(ev.User)
	if err != nil {
		return "Got error when trying to write DB: " + err.Error()
	}
	url := ttLoginUrl + "?source=unislack&thirdpartyid=" + user.Secret
	return "Login med unieconomy brukeren din: " + url
}

func addUser(userID string) (User, error) {
	user := User{
		ID:     userID,
		Secret: userID + "-" + randomString(secretRandLenght),
	}
	err := storage.DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(storage.TraveltextBucket))
		if err != nil {
			return err
		}
		userString, err := json.Marshal(user)
		if err != nil {
			return err
		}
		err = bucket.Put([]byte(userID), userString)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return user, err
	}
	return user, nil
}

func runTTCommand(cmd string, user User) (*TTResponse, error) {
	payload := TTRequestPayload{
		UniToken:                "",
		ThirdPartyIdentificator: user.Secret,
		Command:                 cmd,
	}

	return postRequest(ttCommandUrl, payload)
}

func postRequest(url string, payload TTRequestPayload) (ttResponse *TTResponse, err error) {
	ttResponse = &TTResponse{}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("Couldn't marshal the payload to json: " + err.Error())
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request error: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		return nil, errUnAuthenticated
	} else if resp.StatusCode < 100 || resp.StatusCode > 299 {
		return nil, fmt.Errorf(resp.Status)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, ttResponse)
	if err != nil {
		return nil, fmt.Errorf("Couldn't unmarshal the response payload from json: " + string(body))
	}
	return ttResponse, nil
}

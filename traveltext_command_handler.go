package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	scribble "github.com/nanobox-io/golang-scribble"
	"github.com/nlopes/slack"
)

const (
	DB_FOLDER                  = "./unibot.db"
	DB_NAME                    = "unibot"
	SECRET_RAND_LENGHT         = 32
	BASE_URL                   = "https://traveltext-jarvis.azurewebsites.net"
	TT_LOGIN_URL               = BASE_URL + "/Signup"
	TT_COMMAND_URL             = BASE_URL + "/api/UniHour/Command"
	TRAVEL_TEXT_COMMAND_PREFIX = "tt "
)

var (
	DB                 *scribble.Driver
	errUnAuthenticated = errors.New("User is un-authenticated")
)

func init() {
	var err error
	DB, err = scribble.New(DB_FOLDER, nil)
	if err != nil {
		panic("Couldn't read DB file: " + DB_FOLDER)
	}
}

func HandleTravelTextCommand(ev *slack.MessageEvent, rtm *slack.RTM) {
	user := User{}
	var message string
	getExistingUserError := DB.Read(DB_NAME, ev.User, &user)

	if getExistingUserError != nil {
		message = authenticateUser(ev, rtm)
	} else {
		cmd := strings.TrimPrefix(strings.ToLower(ev.Text), TRAVEL_TEXT_COMMAND_PREFIX)
		response, err := runTTCommand(cmd, user)
		if err != nil {
			if err == errUnAuthenticated {
				message = authenticateUser(ev, rtm)
			} else {
				message = "Got error when trying to send command to TravelText: " + err.Error()
			}
		} else {
			message = response.FriendlyResponse
		}
	}
	rtm.SendMessage(rtm.NewOutgoingMessage(message, ev.Channel))
}

func authenticateUser(ev *slack.MessageEvent, rtm *slack.RTM) string {
	user, err := addUser(ev.User)
	if err != nil {
		return "Got error when trying to write DB: " + err.Error()
	} else {
		url := TT_LOGIN_URL + "?source=unislack&thirdpartyid=" + user.Secret
		return "Login med unieconomy brukeren din: " + url
	}
}

func addUser(userID string) (User, error) {
	user := User{
		ID:     userID,
		Secret: userID + "-" + RandomString(SECRET_RAND_LENGHT),
	}
	err := DB.Write(DB_NAME, userID, user)

	if err != nil {
		return user, err
	} else {
		return user, nil
	}
}

func runTTCommand(cmd string, user User) (TTResponse, error) {
	payload := TTRequestPayload{
		UniToken:                "",
		ThirdPartyIdentificator: user.Secret,
		Command:                 cmd,
	}

	return postRequest(TT_COMMAND_URL, payload)
}

func postRequest(url string, payload TTRequestPayload) (TTResponse, error) {
	ttResponse := TTResponse{}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return ttResponse, fmt.Errorf("Couldn't marshal the payload to json: " + err.Error())
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ttResponse, fmt.Errorf("HTTP request error: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		return ttResponse, errUnAuthenticated
	}

	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &ttResponse)
	if err != nil {
		return ttResponse, fmt.Errorf("Couldn't unmarshal the response payload from json: " + string(body))
	}
	return ttResponse, nil
}

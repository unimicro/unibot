package traveltext

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
	dbFolder                = "./unibot.db"
	dbName                  = "unibot"
	secretRandLenght        = 32
	baseUrl                 = "https://traveltext-jarvis.azurewebsites.net"
	ttLoginUrl              = baseUrl + "/Signup"
	ttCommandUrl            = baseUrl + "/api/UniHour/Command"
	TravelTextCommandPrefix = "tt "
)

var (
	db                 *scribble.Driver
	errUnAuthenticated = errors.New("User is un-authenticated")
)

func init() {
	var err error
	db, err = scribble.New(dbFolder, nil)
	if err != nil {
		panic("Couldn't read the DB folder: " + dbFolder)
	}
}

func HandleTravelTextCommand(ev *slack.MessageEvent, rtm *slack.RTM) {
	user := User{}
	var message string
	getExistingUserError := db.Read(dbName, ev.User, &user)

	rtm.SendMessage(rtm.NewOutgoingMessage("loading traveltext...", ev.Channel))
	if getExistingUserError != nil {
		message = authenticateUser(ev, rtm)
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

	rtm.SendMessage(rtm.NewOutgoingMessage(message, ev.Channel))
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
	err := db.Write(dbName, userID, user)

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

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request error: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		return nil, errUnAuthenticated
	}

	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, ttResponse)
	if err != nil {
		return nil, fmt.Errorf("Couldn't unmarshal the response payload from json: " + string(body))
	}
	return ttResponse, nil
}

package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/unimicro/unibot/constants"
)

var tokens *ServiceAuthentication

func GetTokens() ServiceAuthentication {
	if tokens != nil {
		return *tokens
	} else {
		t := readTokenFile(constants.TokenFileLocation)
		tokens = &t
		return *tokens
	}
}

func readTokenFile(path string) (tokens ServiceAuthentication) {

	authFile, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("Couldn't read the auth/token file: %s\n", path))
	}

	err = json.Unmarshal(authFile, &tokens)
	if err != nil {
		panic(fmt.Sprintf(
			"Couldn't parse the auth/token file (%s) as JSON: %s\n",
			path,
			authFile,
		))
	}

	if tokens.Slack == "" {
		panic("\"slack\" key missing from tokens file " + path)
	}

	if tokens.Gitter == "" {
		panic("\"gitter\" key missing from tokens file " + path)
	}

	if tokens.Jira == "" {
		panic("\"jira\" key missing from tokens file " + path)
	}
	return
}

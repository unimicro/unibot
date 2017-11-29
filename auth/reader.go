package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const (
	slackKey  = "slack"
	gitterKey = "gitter"
)

func ReadTokenFile(path string) (token ServiceAuthentication) {
	authFile, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("Couldn't read the auth/token file: %s\n", path))
	}

	err = json.Unmarshal(authFile, &token)
	if err != nil {
		panic(fmt.Sprintf(
			"Couldn't parse the auth/token file (%s) as JSON: %s\n",
			path,
			authFile,
		))
	}

	if token.Slack == "" {
		panic("\"slack\" key missing from tokens file " + path)
	}

	if token.Gitter == "" {
		panic("\"gitter\" key missing from tokens file " + path)
	}
	return
}

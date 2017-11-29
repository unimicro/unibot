package auth

type ServiceAuthentication struct {
	Slack  Token `json:"slack"`
	Gitter Token `json:"gitter"`
}

type Token string

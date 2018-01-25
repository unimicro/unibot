package auth

type ServiceAuthentication struct {
	Slack  Token `json:"slack"`
	Gitter Token `json:"gitter"`
	Jira   Token `json:"jira"`
}

type Token string

func (t *Token) AsString() string {
	return string(*t)
}

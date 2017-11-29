package gitter

import "time"

type GitterMessage struct {
	ID       string    `json:"id"`
	Text     string    `json:"text"`
	HTML     string    `json:"html"`
	Sent     time.Time `json:"sent"`
	FromUser struct {
		ID              string `json:"id"`
		Username        string `json:"username"`
		DisplayName     string `json:"displayName"`
		URL             string `json:"url"`
		AvatarURL       string `json:"avatarUrl"`
		AvatarURLSmall  string `json:"avatarUrlSmall"`
		AvatarURLMedium string `json:"avatarUrlMedium"`
		V               int    `json:"v"`
		Gv              string `json:"gv"`
	} `json:"fromUser"`
	Unread   bool          `json:"unread"`
	ReadBy   int           `json:"readBy"`
	Urls     []interface{} `json:"urls"`
	Mentions []interface{} `json:"mentions"`
	Issues   []interface{} `json:"issues"`
	Meta     []interface{} `json:"meta"`
	V        int           `json:"v"`
}

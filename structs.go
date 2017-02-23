package main

type User struct {
	ID     string
	Secret string
}

type TTRequestPayload struct {
	UniToken                string `json:"UniToken"`
	ThirdPartyIdentificator string `json:"ThirdPartyIdentificator"`
	Command                 string `json:"Command"`
}

type TTResponse struct {
	FriendlyResponse string     `json:"FriendlyResponse"`
	Intent           Intent     `json:"Intent"`
	Entities         []Entities `json:"Entities"`
}

type Intent struct {
	Name  string   `json:"Name"`
	Score *float64 `json:"Score"`
}

type Entities struct {
	Name  string   `json:"Name"`
	Value string   `json:"Value"`
	Score *float64 `json:"Score"`
}

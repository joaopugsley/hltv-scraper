package models

type TeamData struct {
	Logo *string `json:"logo,omitempty"`
	Name string  `json:"name,omitempty"`
}

type EventData struct {
	Logo *string `json:"logo,omitempty"`
	Name string  `json:"name,omitempty"`
}

type MatchData struct {
	HltvURL string      `json:"hltv_url"`
	Format  string      `json:"format"`
	Time    string      `json:"time"`
	Teams   [2]TeamData `json:"teams"`
	Event   EventData   `json:"event"`
}

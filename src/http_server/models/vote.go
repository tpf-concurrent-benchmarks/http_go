package models

type Vote struct {
	PollID   string    `json:"poll_id"`
	Option   int    `json:"option"`
}

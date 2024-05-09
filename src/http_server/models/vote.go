package models

type Vote struct {
	Username string `json:"username"`
	PollID   string    `json:"poll_id"`
	Option   int    `json:"option"`
}

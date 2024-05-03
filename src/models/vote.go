type Vote struct {
	Username string `json:"username"`
	PollID   int    `json:"poll_id"`
	Option   int    `json:"option"`
}

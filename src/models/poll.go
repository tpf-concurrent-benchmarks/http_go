type Option struct {
	Name  string `json:"name"`
	Votes int    `json:"votes,omitempty"`
}

type PollWithVotes struct {
	Title   string   `json:"title"`
	Options []Option `json:"options"`
}

type PollInDB struct {
	PollWithVotes
	ID int `json:"id"`
}

type Poll struct {
	Title   string   `json:"title"`
	Options []string `json:"options"`
}

type PollGet struct {
	Poll
	ID int `json:"id"`
}
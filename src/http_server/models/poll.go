package models

type Option struct {
	Name  string `json:"name"`
	Votes int    `json:"votes"`
}

type PollWithVotes struct {
	Title   string   `json:"title"`
	Options []Option `json:"options"`
}

type Poll struct {
	Title   string   `json:"title"`
	Options []string `json:"options"`
}

type PollMeta struct {
	ID string `json:"id"`
	Title string `json:"title"`
}
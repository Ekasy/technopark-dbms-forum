package models

type Vote struct {
	ThreadId   int64
	ThreadSlug string
	Nickname   string `json:"nickname"`
	Voice      int64  `json:"voice"`
}

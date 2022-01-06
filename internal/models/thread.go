package models

import (
	"time"
)

type Thread struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	Forum   string `json:"forum"`
	Message string `json:"message"`
	Votes   int64  `json:"votes"`
	Slug    string `json:"slug"`
	Created string `json:"created"`
}

type ThreadInput struct {
	Title   string `json:"title"`
	Author  string `json:"author"`
	Forum   string `json:"forum"`
	Message string `json:"message"`
	Slug    string `json:"slug"`
	Created string `json:"created"`
}

const (
	Layout string = "2006-01-02T15:04:05.000-07:00"
)

func (ti *ThreadInput) ToThread(forumSlug string) *Thread {
	dt := time.Now().Format(Layout)
	if ti.Created == "" {
		ti.Created = dt
	}
	return &Thread{
		Title:   ti.Title,
		Author:  ti.Author,
		Forum:   forumSlug,
		Message: ti.Message,
		Slug:    ti.Slug,
		Created: ti.Created,
	}
}

type ThreadUpdate struct {
	Id      int64
	Slug    string
	Message string `json:"message"`
	Title   string `json:"title"`
}

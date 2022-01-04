package models

import "time"

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

func (ti *ThreadInput) ToThread(forumSlug string) *Thread {
	dt := time.Now().Format("2006-01-02T15:03:05.527+00:00")
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

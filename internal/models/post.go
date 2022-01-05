package models

type Post struct {
	Id       int64  `json:"id"`
	Parent   int64  `json:"parent,omitempty"`
	Author   string `json:"author"`
	Message  string `json:"message"`
	IsEdited bool   `json:"isEdited,omitempty"`
	Forum    string `json:"forum"`
	Thread   int64  `json:"thread"`
	Created  string `json:"created"`
}

type PostInput struct {
	Parent  int64  `json:"parent,omitempty"`
	Author  string `json:"author"`
	Message string `json:"message"`
}

type PostsInput struct {
	PI []*PostInput
}

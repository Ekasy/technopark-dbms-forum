package models

type Forum struct {
	Title   string `json:"title"`
	User    string `json:"user"`
	Slug    string `json:"slug"`
	Posts   int64  `json:"posts"`
	Threads int64  `json:"threads"`
}

type ForumInput struct {
	Title string `json:"title"`
	User  string `json:"user"`
	Slug  string `json:"slug"`
}

func (fi *ForumInput) ToForum(posts int64, threads int64) *Forum {
	return &Forum{
		Title:   fi.Title,
		User:    fi.User,
		Slug:    fi.Slug,
		Posts:   posts,
		Threads: threads,
	}
}

func (fi *ForumInput) ToDefaultForum() *Forum {
	return fi.ToForum(0, 0)
}

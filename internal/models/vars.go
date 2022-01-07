package models

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type ThreadsVars struct {
	ForumSlug string
	Limit     int64
	Since     string
	Sorting   string
	Sign      string
}

func NewThreadsVars(vars map[string]string, query url.Values) *ThreadsVars {
	tv := &ThreadsVars{
		ForumSlug: vars["slug"],
		Limit:     100,
		Since:     "",
		Sorting:   "ASC",
		Sign:      ">=",
	}

	limit, err := strconv.ParseInt(query.Get("limit"), 10, 64)
	if err == nil {
		tv.Limit = limit
	}

	since := query.Get("since")
	if since != "" {
		tv.Since = since
	}

	// desc sorting
	sorting, err := strconv.ParseBool(query.Get("desc"))
	if err == nil {
		if sorting {
			tv.Sorting = "DESC"
			tv.Sign = "<="
		} else {
			tv.Sign = ">="
		}
	}

	return tv
}

type ThreadsQuery struct {
	ThreadId   int64
	ThreadSlug string
	Limit      int64
	Since      int64
	Sort       string
	Sign       string
	Sorting    string
}

func NewThreadQuery(vars map[string]string, query url.Values) *ThreadsQuery {
	slug := vars["slug_or_id"]
	id, err := strconv.ParseInt(slug, 10, 64)
	if err == nil {
		slug = ""
	} else {
		id = 0
	}

	tq := &ThreadsQuery{
		ThreadId:   id,
		ThreadSlug: slug,
		Limit:      100,
		Since:      0,
		Sort:       "flat",
		Sign:       ">",
		Sorting:    "ASC",
	}

	limit, err := strconv.ParseInt(query.Get("limit"), 10, 64)
	if err == nil {
		tq.Limit = limit
	}

	since, err := strconv.ParseInt(query.Get("since"), 10, 64)
	if err == nil {
		tq.Since = since
	}

	sort := query.Get("sort")
	if sort != "" {
		tq.Sort = sort
	}

	// desc sorting
	sorting, err := strconv.ParseBool(query.Get("desc"))
	if err == nil {
		if sorting {
			tq.Sorting = "DESC"
			tq.Sign = "<"
		} else {
			tq.Sign = ">"
		}
	}
	return tq
}

type ForumUsersQuery struct {
	ForumSlug string
	Limit     int64
	Since     string
	Sorting   string
	Sign      string
}

func NewForumUsersQuery(vars map[string]string, query url.Values) *ForumUsersQuery {
	tv := &ForumUsersQuery{
		ForumSlug: vars["slug"],
		Limit:     100,
		Since:     "",
		Sorting:   "ASC",
		Sign:      ">",
	}

	limit, err := strconv.ParseInt(query.Get("limit"), 10, 64)
	if err == nil {
		tv.Limit = limit
	}

	since := query.Get("since")
	if since != "" {
		tv.Since = since
	}

	// desc sorting
	sorting, err := strconv.ParseBool(query.Get("desc"))
	if err == nil {
		if sorting {
			tv.Sorting = "DESC"
			tv.Sign = "<"
		} else {
			tv.Sign = ">"
		}
	}

	return tv
}

type PostQuery struct {
	PostId  int64
	Related []string
}

func NewPostQuery(vars map[string]string, r *http.Request) *PostQuery {
	pq := &PostQuery{
		PostId:  0,
		Related: make([]string, 0),
	}

	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err == nil {
		pq.PostId = id
	}

	related := r.URL.RawQuery
	if strings.Contains(related, "user") {
		pq.Related = append(pq.Related, "user")
	}

	if strings.Contains(related, "forum") {
		pq.Related = append(pq.Related, "forum")
	}

	if strings.Contains(related, "thread") {
		pq.Related = append(pq.Related, "thread")
	}

	return pq
}

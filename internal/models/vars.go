package models

import (
	"net/url"
	"strconv"
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

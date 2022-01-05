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

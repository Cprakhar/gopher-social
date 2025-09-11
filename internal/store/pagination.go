package store

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type PaginatedFeedQuery struct {
	Limit  int    `validate:"min=1,max=20"`
	Offset int    `validate:"min=0"`
	Sort   string `validate:"oneof=asc desc"`
	Search string
	Tags   []string
	Since  *time.Time
	Until  *time.Time
}

func (p PaginatedFeedQuery) Parse(ctx *gin.Context) (PaginatedFeedQuery, error) {
	limit := ctx.Query("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return p, err
		}
		p.Limit = l
	}

	offset := ctx.Query("offset")
	if offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			return p, err
		}
		p.Offset = o
	}

	sort := ctx.Query("sort")
	if sort != "" {
		p.Sort = sort
	}

	tags := ctx.QueryArray("tags")
	if tags != nil {
		p.Tags = tags
	}

	search := ctx.Query("search")
	if search != "" {
		p.Search = search
	}

	since := ctx.Query("since")
	if since != "" {
		s, err := time.Parse(time.RFC3339, since)
		if err != nil {
			return p, err
		}
		p.Since = &s
	}

	until := ctx.Query("until")
	if until != "" {
		u, err := time.Parse(time.RFC3339, until)
		if err != nil {
			return p, err
		}
		p.Until = &u
	}

	return p, nil
}

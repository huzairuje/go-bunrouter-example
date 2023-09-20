package primitive

import (
	"time"
)

type Article struct {
	ID        int64     `bun:"id,pk,autoincrement"`
	Author    string    `bun:"author"`
	Title     string    `bun:"title"`
	Body      string    `bun:"body"`
	CreatedAt time.Time `bun:"created_at"`
	UpdatedAt time.Time `bun:"updated_at"`
	DeletedAt time.Time `bun:"deleted_at"`
}

type ParameterFindArticle struct {
	Query     string
	Author    string
	PageSize  int
	Offset    int
	SortBy    string
	SortOrder string
}

type ParameterArticleHandler struct {
	Query  string
	Author string
}

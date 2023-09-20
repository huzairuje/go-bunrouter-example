package health

import (
	"github.com/uptrace/bun"
)

type RepositoryInterface interface {
	CheckUpTimeDB() (err error)
}

type Repository struct {
	db *bun.DB
}

func NewRepository(db *bun.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r Repository) CheckUpTimeDB() (err error) {
	err = r.db.Ping()
	if err != nil {
		return err
	}

	return
}

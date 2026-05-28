package repository

import (
	"errors"

	"github.com/wb-go/wbf/dbpg"
)

var (
	ErrAliasNotFound    = errors.New("need to create short_url first")
	ErrUniqueConstraint = errors.New("short_url already exists in db")
)

type Repository struct {
	db *dbpg.DB
}

func New(db *dbpg.DB) *Repository {
	return &Repository{
		db: db,
	}
}

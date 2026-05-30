package repository

import (
	"errors"

	"github.com/wb-go/wbf/dbpg"
)

var (
	ErrNoSuchItem = errors.New("no item with such id")
	ErrNoSuchUser = errors.New("no such user")
)

type Repository struct {
	db *dbpg.DB
}

func New(db *dbpg.DB) *Repository {
	return &Repository{
		db: db,
	}
}

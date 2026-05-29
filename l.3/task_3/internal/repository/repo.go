package repository

import (
	"errors"

	"github.com/wb-go/wbf/dbpg"
)

var (
	ErrNotSuchComment = errors.New("there is not comment with such id")
	ErrInvalidParenID = errors.New("there is not parent comment with provided id")
)

type Repository struct {
	db *dbpg.DB
}

func New(db *dbpg.DB) *Repository {
	return &Repository{
		db: db,
	}
}

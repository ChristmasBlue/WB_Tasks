package service

import (
	"task-3/internal/dto"
	"task-3/internal/model"
)

type Storage interface {
	GetCommentsById(id int) ([]*model.Comment, error)
	GetCommentsPaginated(config dto.CommentsPagination) ([]*model.Comment, error)
	GetAllComments() ([]*model.Comment, error)
	GetCommentsByTextSearch(text string) ([]*model.Comment, error)
	CreateComment(comment dto.CreateComment) (*dto.CreateComment, error)
	DeleteCommentById(id int) error
}

type Service struct {
	storage Storage
}

func New(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

package handler

import (
	"task-3/internal/dto"
	"task-3/internal/model"
)

type CommentService interface {
	GetAllComments() ([]*model.Comment, error)
	GetCommentsById(int) ([]*model.Comment, error)
	GetCommentsPaginated(dto.CommentsPagination) ([]*model.Comment, error)
	GetCommentsByTextSearch(string) ([]*model.Comment, error)
	CreateComment(dto.CreateComment) (*dto.CreateComment, error)
	DeleteCommentById(int) error
}

type Handler struct {
	service CommentService
}

func New(service CommentService) *Handler {
	return &Handler{
		service: service,
	}
}

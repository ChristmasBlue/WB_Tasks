package service

import (
	"task-3/internal/dto"
	"task-3/internal/model"
)

func (s *Service) GetAllComments() ([]*model.Comment, error) {
	return s.storage.GetAllComments()
}

func (s *Service) GetCommentsById(id int) ([]*model.Comment, error) {
	return s.storage.GetCommentsById(id)
}

func (s *Service) GetCommentsPaginated(config dto.CommentsPagination) ([]*model.Comment, error) {
	return s.storage.GetCommentsPaginated(config)
}

func (s *Service) GetCommentsByTextSearch(text string) ([]*model.Comment, error) {
	return s.storage.GetCommentsByTextSearch(text)
}

package service

import (
	"task-3/internal/dto"
)

func (s *Service) CreateComment(comment dto.CreateComment) (*dto.CreateComment, error) {
	return s.storage.CreateComment(comment)
}

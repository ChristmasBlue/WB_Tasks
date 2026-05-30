package service

import (
	"context"
	"task-7/internal/dto"
	"task-7/internal/model"
)

func (s *Service) CreateItem(ctx context.Context, createItem dto.CreateItem) (*model.Item, error) {
	return s.storage.CreateItem(ctx, createItem)
}

func (s *Service) CreateUser(ctx context.Context, createUser dto.CreateUser) (*model.User, error) {
	return s.storage.CreateUser(ctx, createUser)
}

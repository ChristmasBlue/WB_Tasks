package service

import (
	"context"
	"task-7/internal/model"
)

func (s *Service) GetUsersWithChanges(ctx context.Context) ([]model.UserHistory, error) {
	return s.storage.GetUsersWithChanges(ctx)
}

func (s *Service) GetAllItems(ctx context.Context) ([]model.Item, error) {
	return s.storage.GetAllItems(ctx)
}

func (s *Service) GetUserRole(ctx context.Context, id int) (string, error) {
	return s.storage.GetUserRole(ctx, id)
}

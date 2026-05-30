package service

import (
	"context"

	"task-6/internal/dto"
	"task-6/internal/model"
)

func (r *Service) CreateItem(ctx context.Context, item dto.CreateItem) (*model.Item, error) {
	return r.storage.CreateItem(ctx, item)
}

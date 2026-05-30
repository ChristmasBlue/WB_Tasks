package service

import (
	"context"

	"task-6/internal/dto"
)

func (s *Service) UpdateItem(ctx context.Context, id int, item dto.UpdateItem) error {
	return s.storage.UpdateItem(ctx, id, item)
}

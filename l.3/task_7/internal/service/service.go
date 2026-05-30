package service

import (
	"context"
	"task-7/internal/dto"
	"task-7/internal/model"
)

type Storage interface {
	CreateItem(ctx context.Context, createItem dto.CreateItem) (*model.Item, error)
	CreateUser(ctx context.Context, createUser dto.CreateUser) (*model.User, error)
	GetUsersWithChanges(ctx context.Context) ([]model.UserHistory, error)
	GetAllItems(ctx context.Context) ([]model.Item, error)
	GetUserRole(ctx context.Context, id int) (string, error)
	UpdateItem(ctx context.Context, updateItem *dto.UpdateItem) error
	DeleteItem(ctx context.Context, id int) error
}

type Service struct {
	storage Storage
}

func New(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

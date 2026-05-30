package handler

import (
	"context"
	"task-7/internal/dto"
	"task-7/internal/model"
)

type Service interface {
	CreateItem(ctx context.Context, createItem dto.CreateItem) (*model.Item, error)
	CreateUser(ctx context.Context, createUser dto.CreateUser) (*model.User, error)
	GetUsersWithChanges(ctx context.Context) ([]model.UserHistory, error)
	GetAllItems(ctx context.Context) ([]model.Item, error)
	GetUserRole(ctx context.Context, id int) (string, error)
	UpdateItem(ctx context.Context, updateItem *dto.UpdateItem) error
	DeleteItem(ctx context.Context, id int) error
}

type Handler struct {
	ctx     context.Context
	service Service
}

func New(ctx context.Context, service Service) *Handler {
	return &Handler{
		ctx:     ctx,
		service: service,
	}
}

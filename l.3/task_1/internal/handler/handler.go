package handler

import (
	"context"

	"task-1/internal/dto"
	"task-1/internal/model"
)

type NotifierService interface {
	GetNotificationStatus(int) (*dto.NotificationStatus, error)
	GetAllNotifications() ([]model.Notification, error)
	CreateNotification(model.Notification) (*model.Notification, error)
	UpdateNotificationStatus(int, string) error
	PublishReadyNotifications(context.Context) error
	ConsumeMessages(ctx context.Context) error
}

type Handler struct {
	service NotifierService
}

func New(service NotifierService) *Handler {
	return &Handler{
		service: service,
	}
}

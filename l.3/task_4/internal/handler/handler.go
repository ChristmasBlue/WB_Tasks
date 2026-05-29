package handler

import (
	"github.com/google/uuid"
	"task-4/internal/dto"
	"task-4/internal/model"
)

type ImageProcessorService interface {
	ProcessImage(dto.Message) error
	GetImageStatus(uuid.UUID) (*model.Image, error)
	GetImageById(uuid.UUID) (string, error)
	CreateImage([]byte, dto.Message) (*uuid.UUID, error)
	DeleteImage(uuid.UUID) error
}

type Handler struct {
	service ImageProcessorService
}

func New(service ImageProcessorService) *Handler {
	return &Handler{
		service: service,
	}
}

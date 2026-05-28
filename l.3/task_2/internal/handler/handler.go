package handler

import (
	"task-2/internal/dto"
	"task-2/internal/model"
)

type ShortnerServcie interface {
	GetAnalytics(string) ([]dto.RedirectInfo, error)
	GetUrlByShort(string, model.RedirectInfo) (*model.Url, error)
	CreateShortUrl(model.Url) (*model.Url, error)
	AggregateByUserAgent() ([]dto.UserAgentDTO, error)
	AggregateByDate() ([]dto.DateDTO, error)
	AggregateByMonth() ([]dto.MonthDTO, error)
}

type Handler struct {
	service ShortnerServcie
}

func New(service ShortnerServcie) *Handler {
	return &Handler{
		service: service,
	}
}

package service

import (
	"task-2/internal/dto"
	"task-2/internal/model"
)

type Storage interface {
	CreateShortUrl(model.Url) (*model.Url, error)
	CreateRedirectInfo(model.RedirectInfo) error
	GetUrlByShort(string, model.RedirectInfo) (*model.Url, error)
	GetAnalytics(string) ([]dto.RedirectInfo, error)
	AggregateByUserAgent() ([]dto.UserAgentDTO, error)
	AggregateByDate() ([]dto.DateDTO, error)
	AggregateByMonth() ([]dto.MonthDTO, error)
}

type Cache interface {
	Get(string) (string, error)
	Set(string, interface{}) error
}

type Service struct {
	storage Storage
	cache   Cache
}

func New(storage Storage, cache Cache) *Service {
	return &Service{
		storage: storage,
		cache:   cache,
	}
}

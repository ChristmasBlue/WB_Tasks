package service

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/wb-go/wbf/zlog"
	"task-2/internal/dto"
	"task-2/internal/model"
)

func (s *Service) GetAnalytics(short_url string) ([]dto.RedirectInfo, error) {
	analytics, err := s.storage.GetAnalytics(short_url)
	if err != nil {
		return nil, err
	}

	for _, a := range analytics {
		if a.RedirectCount >= 5 {
			if err := s.cache.Set(a.Url, a.ShortUrl); err != nil {
				zlog.Logger.Error().Msg("could not save url to cache: " + err.Error())
			}
		}
	}

	return analytics, nil
}

func (s *Service) GetUrlByShort(short_url string, redirectInfo model.RedirectInfo) (*model.Url, error) {
	url, err := s.cache.Get(short_url)
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("could not get short_url from redis: %w", err)

	}

	var urlInfo *model.Url
	if err == redis.Nil {
		urlInfo, err = s.storage.GetUrlByShort(short_url, redirectInfo)
		url = urlInfo.Url
		if err != nil {
			return nil, err
		}
	}

	if err := s.storage.CreateRedirectInfo(redirectInfo); err != nil {
		return nil, err
	}

	urlInfo.ShortUrl = short_url
	urlInfo.Url = url

	return urlInfo, nil
}

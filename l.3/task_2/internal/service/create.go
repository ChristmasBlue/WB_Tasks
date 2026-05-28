package service

import (
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
	"task-2/internal/model"
	"task-2/internal/repository"
)

func (s *Service) CreateShortUrl(url model.Url) (*model.Url, error) {
	url.Url = validateUrlScheme(url.Url)

	short_url, err := s.cache.Get(url.Url)
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("could not get value from redis: %w", err)
	}

	var urlInfo *model.Url
	if err == redis.Nil {
		for {
			url.ShortUrl = generateShortLink()
			urlInfo, err = s.storage.CreateShortUrl(url)
			short_url = urlInfo.ShortUrl
			if err != nil && !errors.Is(err, repository.ErrUniqueConstraint) {
				return nil, err
			}

			if errors.Is(err, repository.ErrUniqueConstraint) {
				continue
			}
			break
		}
	}

	urlInfo.Url = url.Url
	urlInfo.ShortUrl = short_url

	return urlInfo, nil
}

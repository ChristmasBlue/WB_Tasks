package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"task-2/internal/model"
)

func (r *Repository) CreateShortUrl(urlInfo model.Url) (*model.Url, error) {
	query := "SELECT * FROM urls WHERE url=$1"
	rows, err := r.db.QueryContext(
		context.Background(),
		query,
		urlInfo.Url,
	)
	if err != nil {
		return nil, fmt.Errorf("could not get url info from db: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&urlInfo.Id, &urlInfo.ShortUrl, &urlInfo.Url)
		if err == nil {
			return &urlInfo, nil
		}
	}

	query = "INSERT INTO urls(url, short_url) VALUES($1, $2)"
	_, err = r.db.ExecContext(
		context.Background(),
		query,
		urlInfo.Url,
		urlInfo.ShortUrl,
	)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrUniqueConstraint
		}
		return nil, fmt.Errorf("could not save url info in db: %w", err)
	}

	return &urlInfo, nil
}

func (r *Repository) CreateRedirectInfo(redirectInfo model.RedirectInfo) error {
	query := `INSERT INTO redirect_analytics (short_url, user_agent) VALUES ($1, $2);`
	_, err := r.db.ExecContext(
		context.Background(),
		query,
		redirectInfo.ShortUrl,
		redirectInfo.UserAgent,
	)
	if err != nil {
		return fmt.Errorf("could not insert redirect info to db: %w", err)
	}

	return nil
}

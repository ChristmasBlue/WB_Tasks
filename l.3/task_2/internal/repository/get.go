package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"task-2/internal/dto"
	"task-2/internal/model"
)

func (r *Repository) GetUrlByShort(short_url string, redirectInfo model.RedirectInfo) (*model.Url, error) {
	tx, err := r.db.Master.Begin()
	if err != nil {
		return nil, fmt.Errorf("could not start transcation: %w", err)
	}
	defer tx.Rollback()

	query := "SELECT * FROM urls WHERE short_url=$1"
	rows, err := tx.Query(
		query,
		short_url,
	)
	if err != nil {
		return nil, fmt.Errorf("could not get alias from db: %w", err)
	}
	defer rows.Close()

	var urlInfo model.Url
	hasNext := false
	for rows.Next() {
		err = rows.Scan(&urlInfo.Id, &urlInfo.ShortUrl, &urlInfo.Url)
		if err != nil {
			return nil, fmt.Errorf("could not scan rows result: %w", err)
		}
		hasNext = true
	}

	if !hasNext {
		return nil, ErrAliasNotFound
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("could not commit transaction: %w", err)
	}

	return &urlInfo, nil
}

func (r *Repository) GetAnalytics(short_url string) ([]dto.RedirectInfo, error) {
	query := `SELECT r.short_url, u.url, COUNT(r.short_url) as redirect_count,
	ARRAY_AGG(DISTINCT r.user_agent ORDER BY r.user_agent) AS all_user_agents,
	ARRAY_AGG(DISTINCT r.request_time ORDER BY r.request_time) AS all_request_times
	FROM redirect_analytics r
	JOIN urls u ON u.short_url = r.short_url
	WHERE r.short_url = $1
    GROUP BY r.short_url, u.url;`

	rows, err := r.db.QueryContext(
		context.Background(),
		query,
		short_url,
	)
	if err != nil {
		return nil, fmt.Errorf("could not get analytics results from db: %w", err)
	}

	var redirectInfo []dto.RedirectInfo
	for rows.Next() {
		var redirect dto.RedirectInfo
		err := rows.Scan(
			&redirect.ShortUrl,
			&redirect.Url,
			&redirect.RedirectCount,
			pq.Array(&redirect.UserAgent),
			pq.Array(&redirect.RequestTime),
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, ErrAliasNotFound
			}
			return nil, fmt.Errorf("could not scan redirectInfo result to model: %w", err)
		}

		redirectInfo = append(redirectInfo, redirect)
	}

	return redirectInfo, nil
}

package repository

import (
	"context"
	"fmt"

	"github.com/lib/pq"
	"task-2/internal/dto"
)

func (r *Repository) AggregateByUserAgent() ([]dto.UserAgentDTO, error) {
	query := `SELECT short_url, COUNT(short_url) AS count,
	ARRAY_AGG(DISTINCT user_agent) AS user_agent
	FROM redirect_analytics 
	GROUP BY short_url;`

	rows, err := r.db.QueryContext(
		context.Background(),
		query,
	)
	if err != nil {
		return nil, fmt.Errorf("could not send request to get aggregated data from db: %w", err)
	}

	var analytics []dto.UserAgentDTO
	for rows.Next() {
		var next dto.UserAgentDTO
		if err := rows.Scan(&next.ShortUrl, &next.RedirectCount, pq.Array(&next.UserAgent)); err != nil {
			return nil, fmt.Errorf("could not scan aggregated data from db: %w", err)
		}
		analytics = append(analytics, next)
	}

	return analytics, nil
}

func (r *Repository) AggregateByDate() ([]dto.DateDTO, error) {
	query := `SELECT COUNT(short_url),
    EXTRACT(DAY FROM request_time) AS day,
    EXTRACT(MONTH FROM request_time) AS month,
    EXTRACT(YEAR FROM request_time) AS year,
    ARRAY_AGG(short_url) AS short_urls,
    ARRAY_AGG(request_time ORDER BY request_time) AS request_times
	FROM redirect_analytics
	GROUP BY day, month, year;`
	rows, err := r.db.QueryContext(
		context.Background(),
		query,
	)
	if err != nil {
		return nil, fmt.Errorf("could not send request to get aggregated data from db: %w", err)
	}

	var analytics []dto.DateDTO
	for rows.Next() {
		time := []string{}
		short_url := []string{}
		var next dto.DateDTO
		err := rows.Scan(
			&next.RedirectCount,
			&next.Day,
			&next.Month,
			&next.Year,
			pq.Array(&short_url),
			pq.Array(&time),
		)
		if err != nil {
			return nil, fmt.Errorf("could not scan aggregated data from db: %w", err)
		}

		for i := range short_url {
			var info dto.UrlInfo
			info.ShortUrl = short_url[i]

			if i < len(time) {
				info.Time = time[i]
			}
			next.UrlInfo = append(next.UrlInfo, info)

		}
		analytics = append(analytics, next)
	}

	return analytics, nil
}

func (r *Repository) AggregateByMonth() ([]dto.MonthDTO, error) {
	query := `SELECT COUNT(short_url),
    EXTRACT(MONTH FROM request_time) AS month,
    EXTRACT(YEAR FROM request_time) AS year,
    ARRAY_AGG(short_url) AS short_urls,
    ARRAY_AGG(request_time ORDER BY request_time) AS request_times
	FROM redirect_analytics
	GROUP BY month, year;`
	rows, err := r.db.QueryContext(
		context.Background(),
		query,
	)
	if err != nil {
		return nil, fmt.Errorf("could not send request to get aggregated data from db: %w", err)
	}

	var analytics []dto.MonthDTO
	for rows.Next() {
		short_url := []string{}
		time := []string{}
		var next dto.MonthDTO
		err := rows.Scan(
			&next.RedirectCount,
			&next.Month,
			&next.Year,
			pq.Array(&short_url),
			pq.Array(&time),
		)
		if err != nil {
			return nil, fmt.Errorf("could not scan aggregated data from db: %w", err)
		}

		for i := range short_url {
			var info dto.UrlInfo
			info.ShortUrl = short_url[i]

			if i < len(time) {
				info.Time = time[i]
			}
			next.UrlInfo = append(next.UrlInfo, info)

		}
		analytics = append(analytics, next)
	}

	return analytics, nil
}

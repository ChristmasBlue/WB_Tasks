package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"task-7/internal/model"
)

func (r *Repository) GetUsersWithChanges(ctx context.Context) ([]model.UserHistory, error) {
	query := `SELECT u.id, u.name, u.role, u.created_at,
  	COALESCE(
    JSON_AGG(
      JSON_BUILD_OBJECT(
        'item_id', ih.item_id,
        'changed_column', ih.changed_column,
        'changed_from', ih.changed_from,
        'change_time', ih.change_time
      )
      ORDER BY ih.change_time
    ) FILTER (WHERE ih.id IS NOT NULL),
    '[]') AS history
	FROM users u
	LEFT JOIN items_history ih ON ih.changed_by_id = u.id
	GROUP BY u.id, u.name, u.role, u.created_at;`

	rows, err := r.db.Master.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not get change history from db: %w", err)
	}
	defer rows.Close()

	var history []model.UserHistory
	for rows.Next() {
		var historyJSON []byte
		var userHistory model.UserHistory
		var changes []model.Change
		err := rows.Scan(
			&userHistory.ID,
			&userHistory.Name,
			&userHistory.Role,
			&userHistory.CreatedAt,
			&historyJSON,
		)
		if err != nil {
			return nil, fmt.Errorf("could not get changes history from db: %w", err)
		}

		if err := json.Unmarshal(historyJSON, &changes); err != nil {
			return nil, fmt.Errorf("could not get changes history from db: %w", err)
		}

		userHistory.History = changes

		history = append(history, userHistory)
	}

	return history, nil
}

func (r *Repository) GetAllItems(ctx context.Context) ([]model.Item, error) {
	query := "SELECT * FROM items"

	rows, err := r.db.Master.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not get items from db: %w", err)
	}

	var items []model.Item
	for rows.Next() {
		var item model.Item
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Count,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("could not scan row result to model: %w", err)
		}

		items = append(items, item)
	}

	return items, nil
}

func (r *Repository) GetUserRole(ctx context.Context, id int) (string, error) {
	query := "SELECT role FROM users WHERE id = $1"

	var role string
	err := r.db.Master.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(&role)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", ErrNoSuchUser
		}
		return "", fmt.Errorf("could not get user role from db: %w", err)
	}

	return role, nil
}

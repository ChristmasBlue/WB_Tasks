package repository

import (
	"context"
	"fmt"
	"task-7/internal/dto"
	"task-7/internal/model"
)

func (r *Repository) CreateItem(ctx context.Context, createItem dto.CreateItem) (*model.Item, error) {
	query := `INSERT INTO items(name, count) VALUES ($1, $2)
	RETURNING id, created_at`

	var item model.Item
	err := r.db.Master.QueryRowContext(
		ctx,
		query,
		createItem.Name,
		createItem.Count,
	).Scan(&item.ID, &item.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("could not create user: %w", err)
	}

	item.Name = createItem.Name
	item.Count = createItem.Count

	return &item, nil
}

func (r *Repository) CreateUser(ctx context.Context, createUser dto.CreateUser) (*model.User, error) {
	query := `INSERT INTO users(name, role) VALUES ($1, $2)
	RETURNING id, created_at`

	var user model.User
	err := r.db.Master.QueryRowContext(
		ctx,
		query,
		createUser.Name,
		createUser.Role,
	).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("could not create user: %w", err)
	}

	user.Name = createUser.Name
	user.Role = createUser.Role

	return &user, nil
}

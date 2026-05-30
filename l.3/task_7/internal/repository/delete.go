package repository

import (
	"context"
	"fmt"
)

func (r *Repository) DeleteItem(ctx context.Context, id int) error {
	query := "DELETE FROM items WHERE id = $1"

	_, err := r.db.Master.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("could not delete item by id: %w", err)
	}

	return nil
}

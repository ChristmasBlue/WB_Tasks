package repository

import (
	"context"
	"fmt"
	"task-7/internal/dto"
)

func (r *Repository) UpdateItem(ctx context.Context, updateItem *dto.UpdateItem) error {
	// I know about risk of sql injection here but it is not possible
	// to use placehlder because of dirver and userID is int also
	// it guarantees than sql injection string will not be placed
	tx, err := r.db.Master.Begin()
	if err != nil {
		return fmt.Errorf("could not start transaction: %w", err)
	}
	defer tx.Rollback()

	query := fmt.Sprintf("SET LOCAL app.current_user_id = %d;", updateItem.UserID)
	_, err = tx.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("could not update item in db: %w", err)
	}

	query = `UPDATE items SET 
		name = COALESCE($1, name),
		count = COALESCE($2, count)
		WHERE id = $3;`
	result, err := tx.ExecContext(
		ctx,
		query,
		updateItem.Name,
		updateItem.Count,
		updateItem.ID,
	)
	if err != nil {
		return fmt.Errorf("could not update item in db: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not get number of affected rows in db: %w", err)
	}

	if rowsAffected == 0 {
		return ErrNoSuchItem
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transcation: %w", err)
	}

	return nil
}

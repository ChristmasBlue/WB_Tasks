package repository

import "fmt"

func (r *Repository) DeleteCommentById(id int) error {
	query := "DELETE FROM comments WHERE id = $1"

	_, err := r.db.Master.Exec(query, id)
	if err != nil {
		return fmt.Errorf("could not delete notification from db: %w", err)
	}

	return nil
}

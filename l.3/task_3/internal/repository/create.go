package repository

import (
	"fmt"

	"task-3/internal/dto"
)

func (r *Repository) CreateComment(comment dto.CreateComment) (*dto.CreateComment, error) {
	query := `INSERT INTO comments(parent_id, text) 
	VALUES ($1, $2) 
	RETURNING id, created_at`

	err := r.db.Master.QueryRow(query, comment.ParentID, comment.Text).Scan(
		&comment.ID,
		&comment.CreatedAt,
	)
	if err != nil {
		if isForeignKeyViolation(err) {
			return nil, ErrInvalidParenID
		}
		return nil, fmt.Errorf("could not create comment in db: %w", err)
	}

	return &comment, nil
}

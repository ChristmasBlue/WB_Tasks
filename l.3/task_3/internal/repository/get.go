package repository

import (
	"fmt"

	"task-3/internal/dto"
	"task-3/internal/model"
)

var (
	defaultLimit = 10
)

func (r *Repository) GetCommentsById(id int) ([]*model.Comment, error) {
	query := `WITH RECURSIVE comment_tree AS (
  	SELECT id, parent_id, text, created_at
 	FROM comments
  	WHERE id = $1 
  
 	UNION

  	SELECT c.id, c.parent_id, c.text, c.created_at
  	FROM comments c
  	INNER JOIN comment_tree ct ON c.parent_id = ct.id
	)
	SELECT * FROM comment_tree;`

	rows, err := r.db.Master.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("could not get comments from db: %w", err)
	}
	defer rows.Close()

	var comments []model.Comment
	for rows.Next() {
		var comment model.Comment
		err := rows.Scan(
			&comment.ID,
			&comment.ParentID,
			&comment.Text,
			&comment.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("could not scan row to model: %w", err)
		}

		comments = append(comments, comment)
	}

	if len(comments) == 0 {
		return nil, ErrNotSuchComment
	}

	commentTree := buildTree(comments)

	return commentTree, nil
}

func (r *Repository) GetCommentsPaginated(config dto.CommentsPagination) ([]*model.Comment, error) {
	if config.Limit != 0 {
		defaultLimit = config.Limit
	}

	offset := (config.Page - 1) * defaultLimit

	query := `WITH RECURSIVE comment_tree AS (
  	SELECT id, parent_id, text, created_at
 	FROM comments
	WHERE id = $1

 	UNION

  	SELECT c.id, c.parent_id, c.text, c.created_at
  	FROM comments c
  	INNER JOIN comment_tree ct ON c.parent_id = ct.id
	)
	SELECT * FROM comment_tree 
	ORDER BY created_at ASC
	LIMIT $2 OFFSET $3;`

	rows, err := r.db.Master.Query(query, config.ParentID, defaultLimit, offset)
	if err != nil {
		return nil, fmt.Errorf("could not get comments from db: %w", err)
	}
	defer rows.Close()

	var comments []model.Comment
	for rows.Next() {
		var comment model.Comment
		err := rows.Scan(
			&comment.ID,
			&comment.ParentID,
			&comment.Text,
			&comment.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("could not scan row to model: %w", err)
		}

		comments = append(comments, comment)
	}

	commentTree := buildTree(comments)

	return commentTree, nil
}

func (r *Repository) GetAllComments() ([]*model.Comment, error) {
	query := "SELECT id, parent_id, text, created_at FROM comments"

	rows, err := r.db.Master.Query(query)
	if err != nil {
		return nil, fmt.Errorf("could not get comments from db: %w", err)
	}
	defer rows.Close()

	var comments []model.Comment
	for rows.Next() {
		var comment model.Comment
		err := rows.Scan(
			&comment.ID,
			&comment.ParentID,
			&comment.Text,
			&comment.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("could not scan row to model: %w", err)
		}

		comments = append(comments, comment)
	}

	commentTree := buildTree(comments)

	return commentTree, nil
}

func (r *Repository) GetCommentsByTextSearch(text string) ([]*model.Comment, error) {
	query := `SELECT id, parent_id, text, created_at FROM comments
	WHERE search_vector @@ plainto_tsquery('russian', $1)
	ORDER BY ts_rank(search_vector, plainto_tsquery('russian', $1));`

	rows, err := r.db.Master.Query(query, text)
	if err != nil {
		return nil, fmt.Errorf("could not get comments from db: %w", err)
	}
	defer rows.Close()

	var comments []model.Comment
	for rows.Next() {
		var comment model.Comment
		err := rows.Scan(
			&comment.ID,
			&comment.ParentID,
			&comment.Text,
			&comment.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("could not scan row to model: %w", err)
		}

		comments = append(comments, comment)
	}

	commentTree := buildTree(comments)

	return commentTree, nil
}

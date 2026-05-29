package repository

import (
	"errors"

	"github.com/lib/pq"
	"task-3/internal/model"
)

func buildTree(comments []model.Comment) []*model.Comment {
	commentMap := make(map[int]*model.Comment)
	var roots []*model.Comment

	for i := range comments {
		commentMap[comments[i].ID] = &comments[i]
	}

	for i := range comments {
		c := &comments[i]
		if c.ParentID == nil {
			roots = append(roots, c)
		} else {
			parent, ok := commentMap[*c.ParentID]
			if ok {
				parent.Children = append(parent.Children, c)
			}
		}
	}

	return roots
}

func isForeignKeyViolation(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "23503"
	}
	return false
}

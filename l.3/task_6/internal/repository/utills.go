package repository

import (
	"strings"

	"task-6/internal/dto"
)

func prepareParams(params dto.GetItemsParams) string {
	var orderByBuilder strings.Builder

	if len(params.SortBy) > 0 {
		orderByBuilder.WriteString(" ORDER BY ")
		for i, field := range params.SortBy {
			orderByBuilder.WriteString(field)
			if i != len(params.SortBy)-1 {
				orderByBuilder.WriteString(", ")
			}
		}
	}
	return orderByBuilder.String()
}

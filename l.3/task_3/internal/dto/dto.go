package dto

import "time"

type CreateComment struct {
	ID        int       `json:"id"`
	ParentID  *int      `json:"parent_id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

type CommentsPagination struct {
	ParentID int
	Page     int
	Limit    int
}

type SearchText struct {
	Text string `json:"text"`
}

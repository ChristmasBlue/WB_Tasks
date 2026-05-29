package model

import "time"

type Comment struct {
	ID        int        `json:"id"`
	ParentID  *int       `json:"parent_id"`
	Text      string     `json:"text"`
	CreatedAt time.Time  `json:"created_at"`
	Children  []*Comment `json:"children"`
}

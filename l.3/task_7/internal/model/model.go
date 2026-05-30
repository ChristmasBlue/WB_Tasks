package model

import "time"

type Item struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Count     int       `json:"count"`
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	Token     string    `json:"token,omitempty"`
	CreatedAt time.Time `json:"create_at"`
}

type Change struct {
	ItemID        int    `json:"item_id"`
	ChangedColumn string `json:"changed_column"`
	ChangedFrom   string `json:"changed_from"`
	ChangeTime    string `json:"change_time"`
}

type UserHistory struct {
	User
	History []Change `json:"history"`
}

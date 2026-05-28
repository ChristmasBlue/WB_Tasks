package model

import "time"

type Url struct {
	Id       int    `json:"-"`
	Url      string `json:"url,omitempty"`
	ShortUrl string `json:"short_url"`
}

type RedirectInfo struct {
	Id          int       `json:"-"`
	ShortUrl    string    `json:"short_url"`
	RequestTime time.Time `json:"request_time"`
	UserAgent   string    `json:"user_agent"`
}

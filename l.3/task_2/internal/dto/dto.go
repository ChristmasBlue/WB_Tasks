package dto

type UserAgentDTO struct {
	ShortUrl      string   `json:"short_url"`
	UserAgent     []string `json:"user_agent"`
	RedirectCount int      `json:"redirect_count"`
}

type UrlInfo struct {
	ShortUrl string `json:"short_url"`
	Time     string `json:"time"`
}

type DateDTO struct {
	Day           int       `json:"day"`
	Month         int       `json:"month"`
	Year          int       `json:"year"`
	UrlInfo       []UrlInfo `json:"url_info"`
	RedirectCount int       `json:"redirect_count"`
}

type MonthDTO struct {
	Month   int `json:"month"`
	Year    int `json:"year"`
	UrlInfo []struct {
		ShortUrl string `json:"short_url"`
		Time     string `json:"time"`
	} `json:"url_info"`
	RedirectCount int `json:"redirect_count"`
}

type RedirectInfo struct {
	Id            int      `json:"-"`
	Url           string   `json:"url"`
	ShortUrl      string   `json:"short_url"`
	RedirectCount int      `json:"redirect_count"`
	RequestTime   []string `json:"request_time"`
	UserAgent     []string `json:"user_agent"`
}

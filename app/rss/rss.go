package rss

import "time"

type RSS struct {
	Source    string    `json:"source" bson:"source"`
	Title     string    `json:"title" bson:"title"`
	Subtitle  string    `json:"subtitle" bson:"subtitle"`
	URL       string    `json:"url" bson:"url"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
}

type SearchResult struct {
	Feeds  []RSS              `json:"rsss"`
	Result SearchResultResult `json:"_result"`
}

type SearchResultResult struct {
	Offset int64 `json:"offset"`
	Limit  int64 `json:"limit"`
	Total  int64 `json:"total"`
}

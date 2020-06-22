package rss

type RSS struct {
	Source    string `json:"source" bson:"source"`
	Title     string `json:"title" bson:"title"`
	Subtitle  string `json:"subtitle" bson:"subtitle"`
	URL       string `json:"url" bson:"url"`
	Timestamp int64  `json:"timestamp" bson:"timestamp"`
}

type SearchResult struct {
	Feeds  []RSS `json:"RSSs"`
	Result struct {
		Offset int64 `json:"offset"`
		Limit  int64 `json:"limit"`
		Total  int64 `json:"total"`
	} `json:"_result"`
}
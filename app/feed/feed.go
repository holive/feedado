package feed

type Feed struct {
	Source      string    `json:"source"`
	Description string    `json:"description,omitempty"`
	Sections    []Section `json:"sections,omitempty"`
}

type Section struct {
	ParentBlockClass string `json:"parent_block_class,omitempty"`
	EachBlockClass   string `json:"each_block_class,omitempty"`
	Title            string `json:"title,omitempty"`
	Subtitle         string `json:"subtitle,omitempty"`
	Url              string `json:"url,omitempty"`
}

type SearchResult struct {
	Feeds  []Feed `json:"feeds"`
	Result struct {
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
		Total  int `json:"total"`
	} `json:"_result"`
}

package feed

type Feed struct {
	Source      string    `bson:"source",json:"source"`
	Description string    `bson:"description",json:"description,omitempty"`
	Sections    []Section `bson:"sections",json:"sections,omitempty"`
}

type Section struct {
	ParentBlockClass string `bson:"parent_block_class",json:"parent_block_class,omitempty"`
	EachBlockClass   string `bson:"each_block_class",json:"each_block_class,omitempty"`
	Title            string `bson:"title",json:"title,omitempty"`
	Subtitle         string `bson:"subtitle",json:"subtitle,omitempty"`
	Url              string `bson:"url",json:"url,omitempty"`
}

type SearchResult struct {
	Feeds  []Feed `json:"feeds"`
	Result struct {
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
		Total  int `json:"total"`
	} `json:"_result"`
}

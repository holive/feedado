package feed

type Feed struct {
	Source      string    `json:"source" bson:"source"`
	Description string    `json:"description,omitempty" bson:"description"`
	Sections    []Section `json:"sections,omitempty,"bson:"sections"`
}

type Section struct {
	ParentBlockClass string `json:"parent_block_class,omitempty" bson:"parent_block_class"`
	EachBlockClass   string `json:"each_block_class,omitempty" bson:"each_block_class"`
	Title            string `json:"title,omitempty" bson:"title"`
	Subtitle         string `json:"subtitle,omitempty" bson:"subtitle"`
	Url              string `json:"url,omitempty" bson:"url"`
}

type SearchResult struct {
	Feeds  []Feed `json:"feeds"`
	Result struct {
		Offset int64 `json:"offset"`
		Limit  int64 `json:"limit"`
		Total  int64 `json:"total"`
	} `json:"_result"`
}

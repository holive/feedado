package feed

type Feed struct {
	Source      string    `json:"source" bson:"source"`
	Description string    `json:"description,omitempty" bson:"description"`
	Sections    []Section `json:"sections,omitempty,"bson:"sections"`
}

type Section struct {
	SectionSelector  string `json:"section_selector" bson:"section_selector"`
	TitleSelector    string `json:"title_selector,omitempty" bson:"title_selector"`
	SubtitleSelector string `json:"subtitle_selector,omitempty" bson:"subtitle_selector"`
	UrlSelector      string `json:"url_selector,omitempty" bson:"url_selector"`
}

type SearchResult struct {
	Feeds  []Feed `json:"feeds"`
	Result struct {
		Offset int64 `json:"offset"`
		Limit  int64 `json:"limit"`
		Total  int64 `json:"total"`
	} `json:"_result"`
}

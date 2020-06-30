package feed

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Feed struct {
	Id          primitive.ObjectID `bson:"_id,omitempty"`
	Source      string             `json:"source" bson:"source"`
	Description string             `json:"description,omitempty" bson:"description"`
	Category    string             `json:"category,omitempty" bson:"category"`
	Sections    []Section          `json:"sections,omitempty,"bson:"sections"`
}

type Section struct {
	SectionSelector     string `json:"section_selector" bson:"section_selector"`
	TitleSelector       string `json:"title_selector,omitempty" bson:"title_selector"`
	TitleMustContain    string `json:"title_must_contain,omitempty" bson:"title_must_contain"`
	SubtitleSelector    string `json:"subtitle_selector,omitempty" bson:"subtitle_selector"`
	SubtitleMustContain string `json:"subtitle_must_contain,omitempty" bson:"subtitle_must_contain"`
	UrlSelector         string `json:"url_selector,omitempty" bson:"url_selector"`
}

type SearchResult struct {
	Feeds  []Feed             `json:"feeds"`
	Result SearchResultResult `json:"_result"`
}

type SearchResultResult struct {
	Offset int64 `json:"offset"`
	Limit  int64 `json:"limit"`
	Total  int64 `json:"total"`
}

type SQS struct {
	ID string `json:"_id"`
}

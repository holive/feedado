package user

type User struct {
	Email string `json:"email" bson:"email"`
	Feeds []Feed `json:"feeds" bson:"feeds"`
}

type Feed struct {
	Name    string   `json:"name" bson:"name"`
	Sources []string `json:"sources" bson:"sources"`
}

type SearchResult struct {
	Users  []User `json:"users"`
	Result struct {
		Offset int64 `json:"offset"`
		Limit  int64 `json:"limit"`
		Total  int64 `json:"total"`
	} `json:"_result"`
}

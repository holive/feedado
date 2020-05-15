package mongo

type FeedRepository struct {
	conn *Client
}

func NewFeedRepository(conn *Client) *FeedRepository {
	return &FeedRepository{
		conn: conn,
	}
}

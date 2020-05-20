package feed

type Repository interface {
	Create(feed *Feed) (*Feed, error)
	Update(feed *Feed) (*Feed, error)
	DeleteBySource(source string) error
	FindBySource(source string) (*Feed, error)
	FindAll(limit string, offset string) (SearchResult, error)
}

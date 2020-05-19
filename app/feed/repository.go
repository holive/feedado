package feed

type Repository interface {
	FindAll(limit string, offset string) (SearchResult, error)
	Create(seller *Feed) (*Feed, error)
	Delete(ID string) error
}

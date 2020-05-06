package metric

import (
	"net/http"

	newrelic "github.com/newrelic/go-agent"
	"github.com/pkg/errors"
)

type contextKey string

// ContextTransaction is used to store and get a transaction from context.
var ContextTransaction = contextKey("transaction")

// Metricer interface is the contract a metric must have to be used.
type Metricer interface {
	StartTransaction(string, http.ResponseWriter, *http.Request) Transaction
	StartExternalSegment(Transaction, *http.Request) *ExternalSegment
	StartSegment(Transaction, string) *Segment
	StartDatastoreSegment(Transaction) *DatastoreSegment
}

// Client is responsible to generate the metric.
type Client struct {
	appname     string
	license     string
	labels      map[string]string
	config      newrelic.Config
	application newrelic.Application
}

// Transaction implements newrelic transaction interface.
type Transaction newrelic.Transaction

// ExternalSegment implements newrelic external segment struct.
type ExternalSegment struct {
	Es *newrelic.ExternalSegment
}

// Segment implements newrelic segment struct.
type Segment struct {
	s *newrelic.Segment
}

// DatastoreSegment implements newrelic datastore segment struct.
type DatastoreSegment struct {
	ds *newrelic.DatastoreSegment
}

// SetProduct defines the datastore type.
// See the constants in https://github.com/newrelic/go-agent/blob/master/datastore.go.
func (s *DatastoreSegment) SetProduct(v string) *DatastoreSegment {
	s.ds.Product = newrelic.DatastoreProduct(v)
	return s
}

// SetCollection defines the table or group.
func (s *DatastoreSegment) SetCollection(v string) *DatastoreSegment {
	s.ds.Collection = v
	return s
}

// SetOperation defines the relevant action, e.g. "SELECT" or "GET".
func (s *DatastoreSegment) SetOperation(v string) *DatastoreSegment {
	s.ds.Operation = v
	return s
}

// SetParameterizedQuery may be set to the query being performed. It must
// not contain any raw parameters, only placeholders.
func (s *DatastoreSegment) SetParameterizedQuery(v string) *DatastoreSegment {
	s.ds.ParameterizedQuery = v
	return s
}

// SetHost defines the name of the server hosting the datastore.
func (s *DatastoreSegment) SetHost(v string) *DatastoreSegment {
	s.ds.Host = v
	return s
}

// SetPortPathOrID to represent either the port, path, or id of the
// datastore being connected to.
func (s *DatastoreSegment) SetPortPathOrID(v string) *DatastoreSegment {
	s.ds.PortPathOrID = v
	return s
}

// SetDatabaseName defines the name of database where the current query is being
// executed.
func (s *DatastoreSegment) SetDatabaseName(v string) *DatastoreSegment {
	s.ds.DatabaseName = v
	return s
}

// End implements newrelic datastore segment end function.
func (s *DatastoreSegment) End() error {
	return s.ds.End()
}

func (c *Client) newConfig() {
	c.config = newrelic.NewConfig(c.appname, c.license)
	c.config.Labels = c.labels
}

func (c *Client) newApplication() error {
	app, err := newrelic.NewApplication(c.config)
	if err != nil {
		return err
	}

	c.application = app

	return nil
}

// StartTransaction implements newrelic StartTransaction function.
func (c *Client) StartTransaction(name string, w http.ResponseWriter, r *http.Request) Transaction {
	return c.application.StartTransaction(name, w, r)
}

// StartExternalSegment starts a new external segment.
func (c *Client) StartExternalSegment(txn Transaction, r *http.Request) *ExternalSegment {
	s := newrelic.StartExternalSegment(txn, r)
	return &ExternalSegment{s}
}

// StartSegment starts a new segment.
func StartSegment(txn Transaction, name string) *Segment {
	s := newrelic.StartSegment(txn, name)
	return &Segment{s}
}

// StartSegment starts a new segment.
func (c *Client) StartSegment(txn Transaction, name string) *Segment {
	s := newrelic.StartSegment(txn, name)
	return &Segment{s}
}

// StartDatastoreSegment starts a new datastore segment.
func StartDatastoreSegment(txn Transaction) *DatastoreSegment {
	ds := &newrelic.DatastoreSegment{
		StartTime: txn.StartSegmentNow(),
	}
	return &DatastoreSegment{ds}
}

// StartDatastoreSegment starts a new datastore segment.
func (c *Client) StartDatastoreSegment(txn Transaction) *DatastoreSegment {
	ds := &newrelic.DatastoreSegment{
		StartTime: txn.StartSegmentNow(),
	}
	return &DatastoreSegment{ds}
}

// SetURL defines the url of this external segment.
func (s *ExternalSegment) SetURL(v string) {
	s.Es.URL = v
}

// SetResponse defines the response of this external segment.
func (s *ExternalSegment) SetResponse(r *http.Response) {
	s.Es.Response = r
}

// End implements newrelic external segment end function.
func (s *ExternalSegment) End() error {
	return s.Es.End()
}

// End implements newrelic segment end function.
func (s *Segment) End() error {
	return s.s.End()
}

// NewClient returns a Client ready to generate metrics.
func NewClient(options ...func(*Client) error) (*Client, error) {
	client := new(Client)

	for _, option := range options {
		if err := option(client); err != nil {
			return nil, errors.Wrap(err, "error during initialization")
		}
	}

	if client.appname == "" || client.license == "" {
		return nil, errors.New("appname and license are required")
	}

	client.newConfig()

	if err := client.newApplication(); err != nil {
		return nil, errors.Wrap(err, "error during initialization")
	}

	return client, nil
}

// SetAppname sets the appname on Client.
func SetAppname(appname string) func(*Client) error {
	return func(c *Client) error {
		c.appname = appname
		return nil
	}
}

// SetLicense sets the license on Client.
func SetLicense(license string) func(*Client) error {
	return func(c *Client) error {
		c.license = license
		return nil
	}
}

// SetLabels sets the labels on Client.
func SetLabels(labels map[string]string) func(*Client) error {
	return func(c *Client) error {
		c.labels = labels
		return nil
	}
}

// DatastoreProduct return a given newrelic datastore product from a text.
func DatastoreProduct(key string) newrelic.DatastoreProduct {
	return newrelic.DatastoreProduct(key)
}

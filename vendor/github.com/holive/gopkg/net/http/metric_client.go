package http

import (
	"net/http"

	"github.com/holive/gopkg/metric"
	"github.com/pkg/errors"
)

// MetricClient is responsible to generate the external segment for metric.
type MetricClient struct {
	name   string
	metric metric.Metricer
	runner Runner
}

// Do executes an external segment metric request.
func (mc *MetricClient) Do(r *http.Request) (*http.Response, error) {
	txn, ok := r.Context().Value(metric.ContextTransaction).(metric.Transaction)
	if !ok {
		txn = nil
	}

	s := mc.metric.StartExternalSegment(txn, r)
	defer s.End()

	s.SetURL(r.URL.String())
	response, err := mc.runner.Do(r)
	s.SetResponse(response)

	return response, err
}

// NewMetricClient returns a client ready to generate metrics with external segments.
func NewMetricClient(options ...func(*MetricClient) error) (*MetricClient, error) {
	client := new(MetricClient)

	for _, option := range options {
		if err := option(client); err != nil {
			return nil, errors.Wrap(err, "error during initialization")
		}
	}

	if client.name == "" {
		return nil, errors.New("segment name is required")
	}
	if client.metric == nil {
		return nil, errors.New("metric client is required")
	}
	if client.runner == nil {
		return nil, errors.New("http runner is required")
	}

	return client, nil
}

// SetName sets the segment name on MetricClient.
func SetName(name string) func(*MetricClient) error {
	return func(mc *MetricClient) error {
		mc.name = name
		return nil
	}
}

// SetMetric sets the metric client on MetricClient.
func SetMetric(metric metric.Metricer) func(*MetricClient) error {
	return func(mc *MetricClient) error {
		mc.metric = metric
		return nil
	}
}

// SetHTTPClient sets the runner on MetricClient.
func SetHTTPClient(runner Runner) func(*MetricClient) error {
	return func(mc *MetricClient) error {
		mc.runner = runner
		return nil
	}
}

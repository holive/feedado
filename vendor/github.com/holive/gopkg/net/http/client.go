package http

import (
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type transporter struct {
	rt      http.RoundTripper
	headers http.Header
}

func (transporter *transporter) RoundTrip(req *http.Request) (*http.Response, error) {
	for key, values := range transporter.headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	return transporter.rt.RoundTrip(req)
}

// NewClient returns a initialized http client.
func NewClient(options ...func(*http.Client) error) (*http.Client, error) {
	dialer, err := NewCacheDial()
	if err != nil {
		return nil, errors.Wrap(err, "cache dialer error")
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			DialContext:           dialer.DialContext,
			MaxIdleConnsPerHost:   500,
			MaxIdleConns:          1000,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			ResponseHeaderTimeout: 30 * time.Second,
		},
	}

	for _, option := range options {
		if err := option(client); err != nil {
			return nil, errors.Wrap(err, "could not initialize http client")
		}
	}

	return client, nil
}

// ClientTimeout defines the http.Runner timeout.
func ClientTimeout(duration time.Duration) func(*http.Client) error {
	return func(client *http.Client) error {
		client.Timeout = duration
		return nil
	}
}

// ClientHeaders defines default headers used in every request.
func ClientHeaders(headers http.Header) func(*http.Client) error {
	return func(client *http.Client) error {
		client.Transport = &transporter{client.Transport, headers}
		return nil
	}
}

// ClientTransport defines default headers used in every request.
func ClientTransport(transport *http.Transport) func(*http.Client) error {
	return func(client *http.Client) error {
		client.Transport = transport
		return nil
	}
}

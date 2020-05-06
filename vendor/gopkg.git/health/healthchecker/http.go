package healthchecker

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"gopkg.git/health"
	infraHTTP "gopkg.git/net/http"
)

// HTTP checker that check a given URL
type HTTP struct {
	name   string
	u      *url.URL
	runner infraHTTP.Runner
	logger log.Logger
}

// Check the dependency health status.
func (h *HTTP) Check(ctx context.Context) health.DependencyResult {
	req, err := http.NewRequest(http.MethodGet, h.u.String(), nil)
	if err != nil {
		return h.error(err)
	}

	resp, err := h.runner.Do(req.WithContext(ctx))
	if err != nil {
		return h.error(err)
	}
	defer h.close(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return h.error(errors.Errorf("unexpected statuscode %d", resp.StatusCode))
	}

	return h.result()
}

func (h *HTTP) close(closer io.Closer) {
	if err := closer.Close(); err != nil {
		level.Error(h.logger).Log(
			"message", errors.Wrap(err, "error during body close").Error(),
		)
	}
}

func (h *HTTP) result() health.DependencyResult {
	return health.DependencyResult{
		Name:      h.name,
		Status:    health.StatusOK,
		Reference: h.u.String(),
	}
}

func (h *HTTP) error(err error) health.DependencyResult {
	r := h.result()
	r.Status = health.DependencyFail
	r.Description = err.Error()
	return r
}

// NewHTTP returns a initialized HTTP checker.
func NewHTTP(options ...func(*HTTP) error) (*HTTP, error) {
	h := &HTTP{}

	for _, option := range options {
		if err := option(h); err != nil {
			return nil, errors.Wrap(err, "invalid option")
		}
	}

	if h.name == "" {
		return nil, errors.New("invalid name")
	}

	if h.u == nil {
		return nil, errors.New("invalid url")
	}

	if h.runner == nil {
		return nil, errors.New("invalid runner")
	}

	return h, nil
}

// HTTPName defines the dependency name.
func HTTPName(name string) func(*HTTP) error {
	return func(h *HTTP) error {
		h.name = name
		return nil
	}
}

// HTTPUrl defines the dependency url to check.
func HTTPUrl(rawurl string) func(*HTTP) error {
	return func(h *HTTP) error {
		u, err := url.Parse(rawurl)
		if err != nil {
			return err
		}
		h.u = u
		return nil
	}
}

// HTTPRunner defines the http runner which will be used to make http requests.
func HTTPRunner(runner infraHTTP.Runner) func(*HTTP) error {
	return func(h *HTTP) error {
		h.runner = runner
		return nil
	}
}

// HTTPLogger defines the http runner which will be used to make http requests.
func HTTPLogger(logger log.Logger) func(*HTTP) error {
	return func(h *HTTP) error {
		h.logger = logger
		return nil
	}
}

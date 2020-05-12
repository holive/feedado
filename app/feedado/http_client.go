package feedado

import (
	"net/http"
	"time"

	"github.com/holive/feedado/app/config"
	infraHTTP "github.com/holive/gopkg/net/http"
	"github.com/pkg/errors"
)

func initHTTPClient(cfg *config.Config) (infraHTTP.Runner, error) {
	defaultHeader := make(http.Header)
	if userAgent := cfg.HTTPClient.UserAgent; userAgent != "" {
		defaultHeader.Add("User-Agent", userAgent)
	}

	dialer, err := infraHTTP.NewCacheDial()
	if err != nil {
		return nil, errors.Wrap(err, "cache dialer error")
	}

	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		MaxIdleConnsPerHost:   50,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
	}

	runner, err := infraHTTP.NewClient(
		infraHTTP.ClientTimeout(cfg.HTTPClient.Timeout),
		infraHTTP.ClientHeaders(defaultHeader),
		infraHTTP.ClientTransport(transport),
	)
	if err != nil {
		return nil, errors.Wrap(err, "Could not create HTTP Client")
	}

	return infraHTTP.NewRetryHTTP(runner), nil
}

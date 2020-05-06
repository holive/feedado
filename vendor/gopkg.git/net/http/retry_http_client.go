package http

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"time"
)

var (
	// Default retry configuration
	defaultRetryWaitMin = 1 * time.Second
	defaultRetryWaitMax = 30 * time.Second
	defaultRetryMax     = 4
)

// CheckRetry specifies a policy for handling retries. It is called
// following each request with the response and error values returned by
// the http.Client. If CheckRetry returns false, the Client stops retrying
// and returns the response to the caller. If CheckRetry returns an error,
// that error value is returned in lieu of the error from the request. The
// Client will close any response body when retrying, but if the retry is
// aborted it is up to the CheckResponse callback to properly close any
// response body before returning.
type CheckRetry func(ctx context.Context, resp *http.Response, err error) (bool, error)

// Backoff specifies a policy for how long to wait between retries.
// It is called after a failing request to determine the amount of time
// that should pass before trying again.
type Backoff func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration

// RetryableClientHTTP does some magic.,
type RetryableClientHTTP struct {
	Runner Runner

	RetryWaitMin time.Duration // Minimum time to wait
	RetryWaitMax time.Duration // Maximum time to wait
	RetryMax     int           // Maximum number of retries

	// CheckRetry specifies the policy for handling retries, and is called
	// after each request. The default policy is DefaultRetryPolicy.
	CheckRetry CheckRetry

	// Backoff specifies the policy for how long to wait between retries
	Backoff Backoff
}

// Do executes an external segment metric request.
func (c *RetryableClientHTTP) Do(req *http.Request) (*http.Response, error) {

	var resp *http.Response
	var err error

	for i := 0; ; i++ {
		var code int // HTTP response code

		// Attempt the request
		resp, err = c.Runner.Do(req)
		if resp != nil {
			code = resp.StatusCode
		}

		// Check if we should continue with retries.
		checkOK, checkErr := c.CheckRetry(req.Context(), resp, err)

		// Now decide if we should continue.
		if !checkOK {
			if checkErr != nil {
				err = checkErr
			}
			return resp, err
		}

		// We do this before drainBody beause there's no need for the I/O if
		// we're breaking out
		remain := c.RetryMax - i
		if remain <= 0 {
			break
		}

		wait := c.Backoff(c.RetryWaitMin, c.RetryWaitMax, i, resp)
		// TODO: Log description of the error.
		// desc := fmt.Sprintf("%s %s", req.Method, req.URL)
		// if code > 0 {
		// 	desc = fmt.Sprintf("%s (status: %d)", desc, code)
		// }
		_ = code

		time.Sleep(wait)
	}

	// By default, we close the response body and return an error without
	// returning the response
	if resp != nil {
		resp.Body.Close()
	}
	return nil, fmt.Errorf("%s %s giving up after %d attempts",
		req.Method, req.URL, c.RetryMax+1)
}

// DefaultRetryPolicy provides a default callback for Client.CheckRetry, which
// will retry on connection errors and server errors.
func DefaultRetryPolicy(ctx context.Context, resp *http.Response, err error) (bool, error) {
	// do not retry on context.Canceled or context.DeadlineExceeded
	if ctx.Err() != nil {
		return false, ctx.Err()
	}

	if err != nil {
		return true, err
	}

	// Check the response code. We retry on 500-range responses to allow
	// the server time to recover, as 500's are typically not permanent
	// errors and may relate to outages on the server side. This will catch
	// invalid response codes as well, like 0 and 999.
	if resp.StatusCode == 0 || (resp.StatusCode >= 500 && resp.StatusCode != 501) {
		return true, nil
	}

	return false, nil
}

// DefaultBackoff provides a default callback for Client.Backoff which
// will perform exponential backoff based on the attempt number and limited
// by the provided minimum and maximum durations.
func DefaultBackoff(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
	mult := math.Pow(2, float64(attemptNum)) * float64(min)
	sleep := time.Duration(mult)
	if float64(sleep) != mult || sleep > max {
		sleep = max
	}
	return sleep
}

// NewRetryHTTP .
func NewRetryHTTP(runner Runner) *RetryableClientHTTP {
	return &RetryableClientHTTP{
		Runner:       runner,
		RetryWaitMin: defaultRetryWaitMin,
		RetryWaitMax: defaultRetryWaitMax,
		RetryMax:     defaultRetryMax,
		CheckRetry:   DefaultRetryPolicy,
		Backoff:      DefaultBackoff,
	}
}

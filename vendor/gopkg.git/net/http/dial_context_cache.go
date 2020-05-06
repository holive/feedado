package http

import (
	"context"
	"net"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/viki-org/dnscache"
)

// A CacheDial contains options for connecting to an address caching the DNS resolution.
type CacheDial struct {
	Resolver *dnscache.Resolver
	Dialer   *net.Dialer
}

// NewCacheDial returns a initialized CacheDial.
func NewCacheDial(options ...func(cd *CacheDial) error) (*CacheDial, error) {
	resolver := dnscache.New(5 * time.Minute)
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
		Resolver: &net.Resolver{
			PreferGo: false,
		},
	}
	cd := &CacheDial{resolver, dialer}
	for _, option := range options {
		if err := option(cd); err != nil {
			return nil, errors.Wrap(err, "could not initialize cache dialer")
		}
	}

	return cd, nil
}

// DialContext connects to the address on the named network using the provided context.
//
// The provided Context must be non-nil. If the context expires before the connection is complete,
// an error is returned. Once successfully connected, any expiration of the context will not affect
// the connection.
//
// When using TCP, and the host in the address parameter resolves to multiple network addresses,
// any dial timeout (from d.Timeout or ctx) is spread over each consecutive dial, such that each is
// given an appropriate fraction of the time to connect. For example, if a host has 4 IP addresses
// and the timeout is 1 minute, the connect to each single address will be given 15 seconds to
// complete before trying the next one.
//
// See func Dial for a description of the network and address parameters.
func (cd *CacheDial) DialContext(
	ctx context.Context, network string, address string,
) (conn net.Conn, err error) {
	separator := strings.LastIndex(address, ":")
	ips, err := cd.Resolver.Fetch(address[:separator])
	if err != nil {
		return nil, err
	}
	for _, ip := range ips {
		conn, err = cd.Dialer.Dial(network, ip.String()+address[separator:])
		if err == nil {
			break
		}
	}
	return
}

// SetCacheTime .
func SetCacheTime(interval time.Duration) func(*CacheDial) error {
	return func(cd *CacheDial) error {
		cd.Resolver = dnscache.New(interval)
		return nil
	}
}

// SetTimeout .
func SetTimeout(interval time.Duration) func(*CacheDial) error {
	return func(cd *CacheDial) error {
		cd.Dialer.Timeout = interval
		return nil
	}
}

// SetDualStack .
func SetDualStack(dualStack bool) func(*CacheDial) error {
	return func(cd *CacheDial) error {
		cd.Dialer.DualStack = dualStack
		return nil
	}
}

// SetKeepAlive .
func SetKeepAlive(interval time.Duration) func(*CacheDial) error {
	return func(cd *CacheDial) error {
		cd.Dialer.KeepAlive = interval
		return nil
	}
}

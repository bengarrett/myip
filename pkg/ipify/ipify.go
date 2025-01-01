// Package myipcom returns your Internet-facing IPv4 or IPv6
// address, sourced from the ipify API.
// https://www.ipify.org
// © Ben Garrett https://github.com/bengarrett/myip
package ipify

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

// https://api.ipify.org
//
// Output:
// 1.1.1.1

var (
	ErrNoIP    = errors.New("ip address is empty")
	ErrInvalid = errors.New("ip address is invalid")
	ErrRequest = errors.New("ipify.org error")
	ErrStatus  = errors.New("unusual ipify.org server response")
)

const (
	domain = "ipify.org"
	Linkv4 = "https://api.ipify.org"
	Linkv6 = "https://api6.ipify.org"
)

// IPv4 returns the clients online IP address.
func IPv4(ctx context.Context, cancel context.CancelFunc) (string, error) {
	return Request(ctx, cancel, Linkv4)
}

// IPv6 returns the clients online IP address. Using this on a network
// that does not support IPv6 will result in an error.
func IPv6(ctx context.Context, cancel context.CancelFunc) (string, error) {
	return Request(ctx, cancel, Linkv6)
}

// Request the ipify API and return a valid IPv4 or IPv6 address.
func Request(ctx context.Context, cancel context.CancelFunc, url string) (string, error) {
	b, err := RequestB(ctx, cancel, url)
	if b == nil && err == nil && errors.Is(ctx.Err(), context.Canceled) {
		return "", nil
	}
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Printf("\n%s: timeout", domain)
			return "", nil
		}
		if err == nil && errors.Is(ctx.Err(), context.Canceled) {
			return "", nil
		}
		if errors.Is(errors.Unwrap(err), context.Canceled) {
			return "", nil
		}
		e := fmt.Errorf("%w: %w", ErrRequest, err)
		return "", e
	}

	ip := string(b)
	if err := Valid(ip); err != nil {
		return ip, err
	}

	return ip, nil
}

// RequestB requests the ipify API and return the response body.
func RequestB(ctx context.Context, cancel context.CancelFunc, url string) ([]byte, error) {
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s, %w", strings.ToLower(resp.Status), ErrStatus)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return b, err
	}

	return b, nil
}

// Valid returns nil if ip is a valid textual representation of an IP address.
func Valid(ip string) error {
	if ip == "" {
		return ErrNoIP
	}
	if net.ParseIP(ip) == nil {
		return ErrInvalid
	}

	return nil
}

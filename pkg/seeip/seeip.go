// Package seeip returns your Internet-facing IPv4 or IPv6
// address, sourced from the SeeIP API.
// https://seeip.org
// © Ben Garrett https://github.com/bengarrett/myip
package seeip

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

// https://seeip.org
//
// Output:
// {"ip":"1.1.1.1"}

var (
	ErrNoIP    = errors.New("ip address is empty")
	ErrInvalid = errors.New("ip address is invalid")
	ErrRequest = errors.New("myip.com error")
	ErrStatus  = errors.New("unusual seeip.org server response")
)

const (
	domain = "ip4.seeip.org"
	linkv4 = "https://ip4.seeip.org"
	linkv6 = "https://ip6.seeip.org"
)

// IPv4 returns the clients online IP address.
func IPv4(ctx context.Context, cancel context.CancelFunc) (string, error) {
	return Request(ctx, cancel, linkv4)
}

// IPv6 returns the clients online IP address. Using this on a network
// that does not support IPv6 will result in an error.
func IPv6(ctx context.Context, cancel context.CancelFunc) (string, error) {
	return Request(ctx, cancel, linkv6)
}

// Request the seeip API URL and return a valid IPv4 or IPv6 address.
func Request(ctx context.Context, cancel context.CancelFunc, url string) (string, error) {
	b, err := request(ctx, cancel, url)
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
		e := fmt.Errorf("%w: %s", ErrRequest, err)
		return "", e
	}

	ip := string(b)
	if ok, err := valid(ip); !ok {
		return ip, err
	}

	return ip, nil
}

func request(ctx context.Context, cancel context.CancelFunc, url string) ([]byte, error) {
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
		return nil, err
	}

	return b, nil
}

func valid(ip string) (bool, error) {
	if ip == "" {
		return false, ErrNoIP
	}

	if net.ParseIP(ip) == nil {
		return false, ErrInvalid
	}

	return true, nil
}

// Package myipcom returns your Internet-facing IPv4 or IPv6
// address, sourced from the MYIP.com API.
// https://www.myip.com
// © Ben Garrett https://github.com/bengarrett/myip
package myipcom

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

// https://api.myip.com
//
// Output:
// {"ip":"118.209.50.85","country":"Australia","cc":"AU"}
// Note: The returned IP address can be either v4 or v6.

// Result of query.
type Result struct {
	IP      string `json:"ip"`
	Country string `json:"country"`
	ISOCode string `json:"cc"`
}

var (
	ErrNoIP    = errors.New("ip address is empty")
	ErrNoIPv4  = errors.New("ip address is not v4")
	ErrNoIPv6  = errors.New("ip address is not v6")
	ErrInvalid = errors.New("ip address is invalid")
	ErrRequest = errors.New("myip.com error")
	ErrStatus  = errors.New("unusual myip.com server response")
)

const (
	domain = "api.myip.com"
	Link   = "https://api.myip.com"
)

// IPv4 returns the clients online IP address.
func IPv4(ctx context.Context, cancel context.CancelFunc) (string, error) {
	s, err := Request(ctx, cancel, Link)
	if err != nil {
		return s, err
	}

	if s == "" {
		return "", nil
	}

	if err := Valid(false, s); err != nil {
		return s, err
	}

	return s, nil
}

// IPv6 returns the clients online IP address. Using this on a network
// that does not support IPv6 will result in an error.
func IPv6(ctx context.Context, cancel context.CancelFunc) (string, error) {
	s, err := Request(ctx, cancel, Link)
	if err != nil {
		return s, err
	}

	if s == "" {
		return "", nil
	}

	if err := Valid(true, s); err != nil {
		return s, err
	}

	return s, nil
}

// Request the myipcom API and return a valid IPv4 or IPv6 address.
func Request(ctx context.Context, cancel context.CancelFunc, url string) (string, error) {
	s, err := RequestS(ctx, cancel, url)
	if s == "" && err == nil && errors.Is(ctx.Err(), context.Canceled) {
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

	return s, err
}

// RequestS requests the myipcom API and return the parsed response body.
func RequestS(ctx context.Context, cancel context.CancelFunc, url string) (string, error) {
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s, %w", strings.ToLower(resp.Status), ErrStatus)
	}

	ip, err := parse(resp.Body)
	if err != nil {
		return "", err
	}

	return ip, nil
}

func parse(r io.Reader) (string, error) {
	var result Result
	jsonParser := json.NewDecoder(r)
	if err := jsonParser.Decode(&result); err != nil {
		return "", err
	}

	return result.IP, nil
}

// Valid returns nil if s is a valid textual representation of an IP address.
func Valid(ipv6 bool, s string) error {
	if s == "" {
		return ErrNoIP
	}

	ip := net.ParseIP(s)
	if ip == nil {
		return ErrInvalid
	}
	if ipv6 && ip.To16() == nil {
		return ErrNoIPv6
	}
	if ipv6 && ip.To16().Equal(ip.To4()) {
		return ErrNoIPv6
	}
	if !ip.To16().Equal(ip.To4()) {
		return nil
	}
	if ip.To4() == nil {
		return ErrNoIPv4
	}
	return nil
}

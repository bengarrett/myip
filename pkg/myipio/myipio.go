// Package myipio returns your Internet-facing IPv4 or IPv6
// address, sourced from the Workshell MyIP API.
// https://www.my-ip.io
// © Ben Garrett https://github.com/bengarrett/myip
package myipio

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

// https://api.my-ip.io/ip.json
//
// Output:
// {
// 	"success": true,
// 	"ip": "100.100.0.0",
// 	"type": "IPv4"
// }

// Result of query.
type Result struct {
	Success bool   `json:"success"`
	IP      string `json:"ip"`
	Type    string `json:"type"`
}

var (
	ErrNoIP      = errors.New("ip address is empty")
	ErrNoSuccess = errors.New("ip address is unsuccessful")
	ErrNoIPv4    = errors.New("ip address is not ipv4")
	ErrNoIPv6    = errors.New("ip address is not ipv6")
	ErrInvalid   = errors.New("ip address is invalid")
	ErrRequest   = errors.New("myip.com error")
	ErrStatus    = errors.New("unusual my-ip.io server response")
)

const (
	domain = "my-ip.io"
	Linkv4 = "https://api4.my-ip.io/ip.json"
	Linkv6 = "https://api6.my-ip.io/ip.json"
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

// Request the myipcom API and return a valid IPv4 or IPv6 address.
func Request(ctx context.Context, cancel context.CancelFunc, url string) (string, error) {
	r, err := RequestR(ctx, cancel, url)
	if r.IP == "" && err == nil && errors.Is(ctx.Err(), context.Canceled) {
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

	if err := r.Valid(url); err != nil {
		return r.IP, err
	}

	return r.IP, nil
}

// RequestR requests the myipcom API and return the parsed response body.
func RequestR(ctx context.Context, cancel context.CancelFunc, url string) (Result, error) {
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Result{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Result{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Result{}, fmt.Errorf("%s, %w", strings.ToLower(resp.Status), ErrStatus)
	}

	r, err := parse(resp.Body)
	if err != nil {
		return Result{}, err
	}

	return r, nil
}

func parse(r io.Reader) (Result, error) {
	var result Result
	jsonParser := json.NewDecoder(r)
	if err := jsonParser.Decode(&result); err != nil {
		return Result{}, err
	}

	return result, nil
}

// Valid returns nil if url is a valid textual representation of an IP address.
func (r Result) Valid(url string) error {
	if r.IP == "" {
		return ErrNoIP
	}

	if !r.Success {
		return ErrNoSuccess
	}

	if url == Linkv4 && !strings.EqualFold(r.Type, "ipv4") {
		return ErrNoIPv4
	}
	if url == Linkv6 && !strings.EqualFold(r.Type, "ipv6") {
		return ErrNoIPv6
	}

	if net.ParseIP(r.IP) == nil {
		return ErrInvalid
	}

	return nil
}

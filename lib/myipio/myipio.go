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

// Result of IPv4 query.
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
	ErrStatus    = errors.New("unusual my-ip.io server response")
)

const (
	domain = "my-ip.io"
	linkv4 = "https://api4.my-ip.io/ip.json"
	linkv6 = "https://api6.my-ip.io/ip.json"
)

// IPv4 returns the Internet facing IPv4 address using the free ipify.org service.
func IPv4(ctx context.Context, cancel context.CancelFunc) (string, error) {
	return Request(ctx, cancel, linkv4)
}

// IPv6 returns the Internet facing IPv6 address using the free ipify.org service.
func IPv6(ctx context.Context, cancel context.CancelFunc) (string, error) {
	return Request(ctx, cancel, linkv6)
}

// Request returns the Internet facing IP address of the free ipify.org service.
func Request(ctx context.Context, cancel context.CancelFunc, url string) (string, error) {
	r, err := request(ctx, cancel, url)
	if r.IP == "" && err == nil && ctx.Err() == context.Canceled {
		return "", nil
	}
	if err != nil {
		switch errors.Unwrap(err) {
		case context.DeadlineExceeded:
			fmt.Printf("\n%s: timeout", domain)
			return "", nil
		default:
			return "", fmt.Errorf("%s error: %s", domain, err)
		}
	}

	if ok, err := r.valid(url); !ok {
		return r.IP, err
	}

	return r.IP, nil
}

func request(ctx context.Context, cancel context.CancelFunc, url string) (Result, error) {
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

	//log.Printf("\nReceived %d from %s\n", resp.StatusCode, url)

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

func (r Result) valid(url string) (bool, error) {
	if r.IP == "" {
		return false, ErrNoIP
	}

	if !r.Success {
		return false, ErrNoSuccess
	}

	if url == linkv4 && !strings.EqualFold(r.Type, "ipv4") {
		return false, ErrNoIPv4
	}
	if url == linkv6 && !strings.EqualFold(r.Type, "ipv6") {
		return false, ErrNoIPv6
	}

	if net.ParseIP(r.IP) == nil {
		return false, ErrInvalid
	}

	return true, nil
}

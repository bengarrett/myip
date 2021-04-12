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

// Result of IPv4 query.
type Result struct {
	IP      string `json:"ip"`
	Country string `json:"country"`
	ISOCode string `json:"cc"`
}

var (
	ErrNoIP    = errors.New("ip address is empty")
	ErrNoIPv4  = errors.New("ip address is not v4")
	ErrInvalid = errors.New("ip address is invalid")
	ErrStatus  = errors.New("unusual myip.com server response")

	link = "https://api.myip.com"
)

const domain = "api.myip.com"

// IPv4 returns the Internet facing IP address using the free myip.com service.
func IPv4(ctx context.Context, cancel context.CancelFunc) (string, error) {
	s, err := request(ctx, cancel, link)
	if s == "" && ctx.Err() == context.Canceled {
		return "", nil
	}
	if err != nil {
		switch errors.Unwrap(err) {
		case context.DeadlineExceeded:
			fmt.Printf("\n%s: timeout\n", domain)
			return "", nil
		default:
			return "", fmt.Errorf("%s error: %s", domain, err)
		}
	}

	if ok, err := valid(s); !ok {
		return s, err
	}

	return s, nil
}

func request(ctx context.Context, cancel context.CancelFunc, url string) (string, error) {
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

	//log.Printf("\nReceived %d from %s\n", resp.StatusCode, url)

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

func valid(s string) (bool, error) {
	if s == "" {
		return false, ErrNoIP
	}
	ip := net.ParseIP(s)
	if ip == nil {
		return false, ErrInvalid
	}
	if ip.To4() == nil {
		return false, ErrNoIPv4
	}

	return true, nil
}

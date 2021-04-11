package myipio

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
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
	ErrInvalid   = errors.New("ip address is invalid")
	ErrStatus    = errors.New("unusual my-ip.io server response")
)

const (
	domain = "api4.my-ip.io"
	link   = "https://api4.my-ip.io/ip.json"
)

// IPv4 returns the Internet facing IP address of the free my-ip.io service.
func IPv4(ctx context.Context, cancel context.CancelFunc) (string, error) {
	s, err := request(ctx, cancel, link)
	if err != nil {
		if _, ok := err.(*url.Error); ok {
			if strings.Contains(err.Error(), "context deadline exceeded") {
				fmt.Printf("\n%s: timeout\n", domain)
				return "", nil
			}
			fmt.Printf("\n%s: %s\n", domain, err)
			return "", nil
		}
		return "", fmt.Errorf("%s error: %s", domain, err)
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

	r, err := parse(resp.Body)
	if err != nil {
		return "", err
	}

	if ok, err := r.valid(); !ok {
		return r.IP, err
	}

	return r.IP, nil
}

func parse(r io.Reader) (Result, error) {
	var result Result
	jsonParser := json.NewDecoder(r)
	if err := jsonParser.Decode(&result); err != nil {
		return Result{}, err
	}

	return result, nil
}

func (r Result) valid() (bool, error) {
	if r.IP == "" {
		return false, ErrNoIP
	}
	if !r.Success {
		return false, ErrNoSuccess
	}
	if !strings.EqualFold(r.Type, "ipv4") {
		return false, ErrNoIPv4
	}
	if net.ParseIP(r.IP) == nil {
		return false, ErrInvalid
	}

	return true, nil
}

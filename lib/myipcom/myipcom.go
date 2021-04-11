package myipcom

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
)

const (
	domain = "api.myip.com"
	link   = "https://api.myip.com"
)

// IPv4 returns the Internet facing IP address using the free myip.com service.
func IPv4(ctx context.Context, cancel context.CancelFunc) (string, error) {
	s, err := request(ctx, cancel, link)
	if errors.Is(err, ErrNoIPv4) {
		return "", nil
	}
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

// func get(d string) (string, error) {
// 	c := &http.Client{
// 		Timeout: timeout * time.Second,
// 	}
// 	ctx := context.Background()
// 	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://"+d, nil)
// 	if err != nil {
// 		return "", err
// 	}
// 	resp, err := c.Do(req)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer resp.Body.Close()
// 	if resp.StatusCode != http.StatusOK {
// 		return "", fmt.Errorf("%s, %w", strings.ToLower(resp.Status), ErrStatus)
// 	}
// 	r, err := parse(resp.Body)
// 	if err != nil {
// 		return "", err
// 	}
// 	if ok, err := r.valid(); !ok {
// 		return r.IP, err
// 	}
// 	return r.IP, nil
// }

// func parse(r io.Reader) (Result, error) {
// 	var result Result
// 	jsonParser := json.NewDecoder(r)
// 	if err := jsonParser.Decode(&result); err != nil {
// 		return Result{}, err
// 	}
// 	return result, nil
// }

func (r Result) valid() (bool, error) {
	if r.IP == "" {
		return false, ErrNoIP
	}
	ip := net.ParseIP(r.IP)
	if ip == nil {
		return false, ErrInvalid
	}
	if ip.To4() == nil {
		return false, ErrNoIPv4
	}
	return true, nil
}

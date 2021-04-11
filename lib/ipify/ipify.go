package ipify

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
)

// https://api.ipify.org
//
// Output:
// 1.1.1.1

var (
	ErrNoIP    = errors.New("ip address is empty")
	ErrInvalid = errors.New("ip address is invalid")
	ErrStatus  = errors.New("unusual ipify.org server response")
)

const (
	domain = "api.ipify.org"
	link   = "https://api.ipify.org"
)

// IPv4 returns the Internet facing IP address using the free ipify.org service.
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

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	ip := string(b)
	if ok, err := valid(ip); !ok {
		return ip, err
	}

	return ip, nil
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

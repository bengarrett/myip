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
	ErrStatus  = errors.New("unusual ipify.org server response")

	link = "https://api.ipify.org"
)

const domain = "api.ipify.org"

// IPv4 returns the Internet facing IP address using the free ipify.org service.
func IPv4(ctx context.Context, cancel context.CancelFunc) (string, error) {
	b, err := request(ctx, cancel, link)
	if err == nil && ctx.Err() == context.Canceled {
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

	//log.Printf("\nReceived %d from %s\n", resp.StatusCode, url)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s, %w", strings.ToLower(resp.Status), ErrStatus)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return b, err
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
